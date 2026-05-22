import { defineConfig } from 'vitepress'

const zhSidebar = [
  {
    text: '快速开始',
    items: [
      { text: '什么是 kozex', link: '/zh/what-is-kozex' },
      { text: '快速开始', link: '/zh/quickstart' },
      { text: '常见问题', link: '/zh/faq' },
    ],
  },
  {
    text: '配置',
    items: [
      { text: '模型配置', link: '/zh/model-configuration' },
      { text: '插件配置', link: '/zh/plugin-configuration' },
      { text: '基础组件配置', link: '/zh/basic-component-configuration' },
    ],
  },
  {
    text: 'API 与 SDK',
    items: [
      { text: 'API 参考', link: '/zh/api-reference' },
      { text: 'Chat SDK', link: '/zh/chat-sdk' },
    ],
  },
  {
    text: '开发指南',
    items: [
      { text: '开发规范', link: '/zh/development-standards' },
      { text: '新增工作流节点（前端）', link: '/zh/add-workflow-nodes-frontend' },
      { text: '新增工作流节点（后端）', link: '/zh/add-workflow-nodes-backend' },
      { text: '新增 API 接口', link: '/zh/add-apis' },
    ],
  },
]

const enSidebar = [
  {
    text: 'Getting Started',
    items: [
      { text: 'What is kozex', link: '/en/what-is-kozex' },
      { text: 'Quickstart', link: '/en/quickstart' },
      { text: 'FAQ', link: '/en/faq' },
    ],
  },
  {
    text: 'Configuration',
    items: [
      { text: 'Model Configuration', link: '/en/model-configuration' },
      { text: 'Plugin Configuration', link: '/en/plugin-configuration' },
      { text: 'Basic Component Configuration', link: '/en/basic-component-configuration' },
    ],
  },
  {
    text: 'API & SDK',
    items: [
      { text: 'API Reference', link: '/en/api-reference' },
      { text: 'Chat SDK', link: '/en/chat-sdk' },
    ],
  },
  {
    text: 'Developer Guide',
    items: [
      { text: 'Development Standards', link: '/en/development-standards' },
      { text: 'Add Workflow Nodes (Frontend)', link: '/en/add-workflow-nodes-frontend' },
      { text: 'Add Workflow Nodes (Backend)', link: '/en/add-workflow-nodes-backend' },
      { text: 'Add APIs', link: '/en/add-apis' },
    ],
  },
]

export default defineConfig({
  base: '/kozex/',
  title: 'Kozex',
  description: 'Enterprise-grade AI Agent Platform',

  locales: {
    zh: {
      label: '简体中文',
      lang: 'zh-CN',
      link: '/zh/',
      description: '企业级 AI Agent 开发平台',
      themeConfig: {
        nav: [
          { text: '快速开始', link: '/zh/quickstart' },
          { text: 'API 参考', link: '/zh/api-reference' },
          { text: '开发指南', link: '/zh/development-standards' },
        ],
        sidebar: zhSidebar,
        outline: { label: '本页目录' },
        docFooter: { prev: '上一页', next: '下一页' },
        darkModeSwitchLabel: '主题',
        lightModeSwitchTitle: '切换到浅色模式',
        darkModeSwitchTitle: '切换到深色模式',
        sidebarMenuLabel: '菜单',
        returnToTopLabel: '回到顶部',
        langMenuLabel: '切换语言',
      },
    },
    en: {
      label: 'English',
      lang: 'en-US',
      link: '/en/',
      themeConfig: {
        nav: [
          { text: 'Quickstart', link: '/en/quickstart' },
          { text: 'API Reference', link: '/en/api-reference' },
          { text: 'Developer Guide', link: '/en/development-standards' },
        ],
        sidebar: enSidebar,
      },
    },
  },

  themeConfig: {
    logo: { light: '/logo.png', dark: '/logo.png', alt: 'Kozex' },
    logoLink: '/kozex/',
    socialLinks: [
      { icon: 'github', link: 'https://github.com/kozex-ai/kozex' },
    ],
    search: {
      provider: 'local',
    },
  },

  head: [
    ['link', { rel: 'icon', href: '/favicon.ico' }],
    ['script', { async: '', src: 'https://www.googletagmanager.com/gtag/js?id=G-CG0KC2WMWT' }],
    ['script', {}, `window.dataLayer = window.dataLayer || [];
function gtag(){dataLayer.push(arguments);}
gtag('js', new Date());
gtag('config', 'G-CG0KC2WMWT');`],
    ['script', {}, `var _hmt = _hmt || [];
(function() {
  var hm = document.createElement("script");
  hm.src = "https://hm.baidu.com/hm.js?fbc77fe3711a0548fb5c914e99daea6a";
  var s = document.getElementsByTagName("script")[0];
  s.parentNode.insertBefore(hm, s);
})();`],
  ],

  ignoreDeadLinks: [
    /localhost/,     // runtime app links in quickstart/config docs
    /\/guides\//,    // upstream coze.cn doc links not migrated
    /coze\.cn/,
  ],
})
