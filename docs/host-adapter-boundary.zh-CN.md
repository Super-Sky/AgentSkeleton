# 宿主适配层边界

## 目的

本文档用于定义 `AgentSkeleton` 与宿主侧接入层之间的产品边界，例如 `Codex` skill/plugin 与 `Claude Code` skill/plugin。

目标是在实现零命令感知、完全会话驱动体验的同时，保持产品核心仍然聚焦于文档编排，而不是演变成一个通用对话 agent。

## 产品方向

目标中的最终形态应当是：

- 用户只通过对话工作
- Codex 与 Claude Code 等宿主决定何时调用 AgentSkeleton
- AgentSkeleton 保持为结构化引擎，而不是扩展成通用自然语言 agent

换句话说：

- 用户体验应该是 agent-native 的
- 产品内核应该是 protocol-oriented 的

## 核心决策

`AgentSkeleton` 不应扩展成一个完整的自然语言 agent runtime。

相反，应采用如下分工：

- 宿主 skill/plugin 负责理解用户意图
- 宿主模型在必要时负责追问和澄清
- AgentSkeleton 只接收和返回结构化工作流状态

这个边界能够避免项目变成第二套职责重叠的 agent 系统。

## 为什么需要这个边界

### 1. 保持产品聚焦

如果 AgentSkeleton 自己承担自然语言理解，它就会从“文档工作流引擎”滑向“agent 平台”。

这会显著增加以下复杂度：

- prompting
- 对话状态管理
- 歧义消解
- 长会话记忆

### 2. 复用宿主优势

Codex 与 Claude Code 已经天然擅长：

- 理解自然语言
- 结合仓库上下文进行推理
- 在对话中主动澄清问题
- 维持会话记忆

AgentSkeleton 应复用这些能力，而不是在核心里重复建设。

### 3. 保持共享协议稳定

结构化核心更容易做到：

- 测试
- 版本化
- 校验
- 跨宿主保持一致

相比之下，把宿主专属的交互逻辑塞进核心，只会让系统更脆弱。

## 职责分工

### 宿主 Skill / Plugin 的职责

宿主层应负责：

- 从对话中理解用户请求
- 判断当前仓库是否应启用 AgentSkeleton
- 判断何时调用 `init-docs` 或 `reshape-docs`
- 判断何时调用 `plan`、`focus-doc` 和 `response --apply`
- 向用户追问缺失的业务上下文
- 判断是推进当前 priority，还是回看 review targets
- 以自然对话的方式向用户呈现进度

### AgentSkeleton 的职责

AgentSkeleton 应负责：

- 仓库文档状态管理
- 结构化工作流推进
- 文档优先级计算
- 起草上下文包生成
- 回看目标计算
- 响应校验与 apply 行为
- apply 之后的下一步状态输出

### 共享边界

宿主可以决定“用户想做什么”。

AgentSkeleton 应决定“当前文档工作流状态意味着什么”。

## AgentSkeleton 应避免承担的内容

为保持边界清晰，AgentSkeleton 应避免自己承担：

- 开放式聊天交互
- 通用自然语言命令解析
- 面向任意任务的自由意图分类
- 宿主专属 UX 策略
- 脱离仓库状态的会话记忆

如果这些职责进入核心，产品边界会迅速模糊，也会让宿主接入更难维护。

## 期望的宿主接入模型

推荐的接入链路应是：

1. 用户通过对话表达目标
2. 宿主 skill/plugin 将目标映射为 AgentSkeleton 动作
3. AgentSkeleton 返回结构化工作流引导
4. 宿主模型负责起草或更新仓库文档
5. 宿主通过 `response --apply` 将结果写回
6. 宿主继续从返回的 `post_apply_plan` 推进

用户不应需要知道这条链路背后的命令名。

## 最小协议要求

为了让这个边界长期成立，AgentSkeleton 需要继续把输出设计成“机器可消费”的形式。

以下信息应保持显式、稳定：

- 当前工作流状态
- 当前 priority
- review candidates
- drafting package 上下文
- missing context
- change batch identity
- post-apply continuation state

只要可以用结构化字段表达，就不要把关键流程语义藏进长段自由文本里。

## Plugin 与 Skill 的定位

Codex 与 Claude Code 的接入形式可以是：

- 原生插件
- 可复用 skill
- 宿主侧 wrapper
- 集成式命令策略

具体封装方式可以因宿主而异，但边界应保持不变：

- 宿主接入层负责意图处理
- AgentSkeleton 负责工作流编排

## 战略含义

这意味着产品更适合被定义为：

- 文档编排引擎
- CLI / protocol core
- 可被宿主集成的工作流系统

而不应被定义为：

- 一个面向任意自然语言任务的独立对话 agent

## 当前阶段建议

现阶段项目应优先投入在：

- 更清晰的 host adapter 规则
- 更强的结构化输出契约
- 更稳定的工作流推进语义
- 更薄的 Codex / Claude Code 接入层

自然语言体验的提升，优先应来自宿主 skill/plugin 的改进，而不是把核心改造成独立 agent runtime。

## 结论

零命令感知体验应该发生在宿主层。

长期稳定的产品逻辑应该保留在 AgentSkeleton 内核中。

这种分层最有利于让项目同时具备：

- 更好的易用性
- 更好的演进性
- 更好的可测试性
- 更好的 Codex / Claude Code 对齐能力
