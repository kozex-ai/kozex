# Long-running Python worker process managed by the coze-sandbox pool.
# Communicates via stdin/stdout using newline-delimited JSON:
#   stdin:  {"code": "...", "params": {...}}
#   stdout: {"result": {...}}  or  {"error": "ExcType: message"}
#
# Each execution gets a fresh namespace so top-level names from one
# call cannot bleed into the next. sys.modules is shared across
# executions within the same worker process; if user code imports a
# module that patches builtins or global state, that patch persists
# until this worker is recycled.
import asyncio
import json
import sys


class Args:
    def __init__(self, params):
        self.params = params


class Output(dict):
    pass


def run_code(code, params):
    ns = {
        '__builtins__': __builtins__,
        'Args': Args,
        'Output': Output,
        'asyncio': asyncio,
        'json': json,
        'sys': sys,
    }
    exec(compile(code, '<user_code>', 'exec'), ns)  # noqa: S102
    main_fn = ns.get('main')
    if main_fn is None:
        raise RuntimeError('main() function not defined')
    result = asyncio.run(main_fn(Args(params)))
    if result is None:
        return {}
    return dict(result)


if __name__ == '__main__':
    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue
        try:
            req = json.loads(line)
            result = run_code(req.get('code', ''), req.get('params', {}))
            sys.stdout.write(json.dumps({'result': result}) + '\n')
        except Exception as e:
            sys.stdout.write(json.dumps({'error': f'{type(e).__name__}: {e}'}) + '\n')
        sys.stdout.flush()
