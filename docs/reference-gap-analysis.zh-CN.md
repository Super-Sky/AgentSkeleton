# 参考样本差距分析

## 目的

本文档用于把 `secure_digital_platform` 作为参考样本进行能力对照，明确：

- 它已经具备哪些成熟的 AI 友好资产
- AgentSkeleton 已经抽象了哪些部分
- AgentSkeleton 还缺哪些关键能力

这份文档只关注可复用的文档化和协作模式，不复制业务实现，也不引入私有情报。

## 参考样本观察

当前参考项目已经具备以下成熟资产：

- 仓库入口文档
  - `README.md`
  - `AGENTS.md`
  - `CLAUDE.md`
- 方法型文档资产
  - `docs/README.md`
  - `docs/CODE_STYLE_GUIDE.md`
  - `docs/DOC_TEMPLATES.md`
  - `docs/REPO_WORKFLOW_GUIDE.md`
  - `docs/TASK_DELIVERY_GUIDE.md`
  - `docs/TASK_EXECUTION_TEMPLATE.md`
- 版本化文档目录
  - `docs/vX.Y.Z/README.md`
  - `docs/vX.Y.Z/features/...`
- 双宿主 skill 落位
  - `.agents/skills/`
  - `.claude/skills/`
- README 中对项目现状、历史结构、技能位置、文档职责的显式说明

这说明参考项目已经从“代码仓库”演进为“可被 agent 理解和协作的仓库”。

## 已覆盖能力

AgentSkeleton 当前已经覆盖的部分：

- 中英文双语入口文档与产品文档
- `AGENTS.md` / `CLAUDE.md` 双宿主协作入口
- 文档蓝图目录与基础模板
- 版本化文档蓝图
- 老项目文档化重塑蓝图
- CLI 上下文模型
- `plan` / `next` / `prompt` / `response` / `workflow` 最小闭环
- 模型返回校验、重试与 repair prompt
- `<output-dir>/.agentskeleton` 过程产物模型
- workflow trace 持久化

这些能力已经足够支撑“对话式文档引导工具”的最小产品形态。

## 仍然存在的差距

### 1. 文档生成链路还不完整

当前虽已支持 `workflow --write-plan-files`，但仍偏基础：

- 只覆盖部分核心文档
- 尚未形成“按文档逐个推进”的生成状态机
- 对已生成文档、待生成文档、待补充上下文的跟踪仍不够细

这是最优先需要补强的能力。

### 2. 方法型文档资产还未完全产品化

虽然我们已经抽象了 `CODE_STYLE_GUIDE`、`DOC_TEMPLATES`、`REPO_WORKFLOW_GUIDE`、`TASK_DELIVERY_GUIDE`、`TASK_EXECUTION_TEMPLATE` 这类蓝图，但还缺：

- 哪类项目默认生成哪些方法文档
- 新项目与老项目在方法文档上的不同优先级
- CLI 如何根据上下文推荐这些方法文档

也就是说，文档模板已经存在，但调度逻辑还不够强。

### 3. skills 集成层还不够明确

参考项目已经形成：

- `.agents/skills/`
- `.claude/skills/`

这种双落位结构。

AgentSkeleton 当前已经在原则层说明 skill 是外部能力、应由标准工具创建，但还没有完整定义：

- 哪些项目需要 skills
- skills 何时应被纳入生成计划
- 文档蓝图与 skill 落位如何联动

### 4. 版本文档与任务交付的联动还不够强

参考项目的版本目录已经在真实迭代中被使用。AgentSkeleton 当前只有蓝图，缺少：

- 某次迭代如何进入 `docs/vX.Y.Z/`
- feature 文档何时创建
- review checklist 何时出现
- 它们如何与 `workflow` 状态绑定

### 5. 缺少“参考样本到通用模型”的映射文档

现在我们知道参考项目里有哪些资产，但还没有一份正式文档说明：

- 参考项目中的哪个资产被抽象成了 AgentSkeleton 的哪一类蓝图
- 哪些被刻意舍弃
- 舍弃原因是什么

这会让后续继续抽象时容易反复讨论。

## 不应照搬的部分

以下内容不应被 AgentSkeleton 直接模板化：

- 参考项目的具体业务实体
- 参考项目的接口路径和服务拆分
- 参考项目的版本兼容历史细节
- 参考项目的具体代码分层实现
- 任意真实数据、私有字段或业务情报

AgentSkeleton 应抽象的是“如何让仓库自解释”和“如何让文档协作稳定发生”，而不是复制某个业务项目。

## 建议的下一步

按优先级建议如下：

1. 让 `workflow` 输出“当前优先生成文档”及其所需上下文。
2. 把已生成文档状态回写到上下文，形成稳定的文档推进闭环。
3. 明确方法型文档在不同项目模式下的生成优先级。
4. 定义 skill 集成触发点，而不是只停留在原则说明。
5. 增加一份“参考样本映射表”，说明参考资产与通用蓝图的对应关系。

## 当前判断

结论是：

AgentSkeleton 的方向与 `secure_digital_platform` 的文档化成熟形态是一致的，但两者不是同一个产品。

- `secure_digital_platform` 是业务仓库中的成熟实践样本
- AgentSkeleton 是将这些实践抽象成通用 CLI 与文档蓝图的产品

因此，AgentSkeleton 应向它学习“文档和协作资产如何成熟”，而不是复制它的业务结构。
