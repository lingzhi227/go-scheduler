A curated list of **Go-based open-source projects** (suitable as scheduling/orchestration building blocks) and **related top-venue papers**.
(Note: many of these projects are strictly *runtime/actor/workflow/k8s scheduler plugins*, but all can serve as the skeleton or key components of a MAS scheduler.)

---

## 1) Go Open-Source Projects (GitHub)

### A. Multi-Agent Orchestration / Scheduling

| Project | Description |
|---------|-------------|
| [clawe](https://github.com/getclawe/clawe) | Multi-agent coordination system where each agent has its own identity, workspace, and cron-based heartbeat, with kanban-style task management |
| [go-agent](https://github.com/Protocol-Lattice/go-agent) | "Lattice" -- production-grade AI agent library with pluggable LLMs, tool calling, RAG memory, and multi-agent coordination |
| [go-swarm](https://github.com/dipendra-sharma/go-swarm) | Multi-agent framework with built-in OpenAI LLM support, tool integration, and inter-agent message passing |
| [swarm-go](https://github.com/feiskyer/swarm-go) | Lightweight multi-agent orchestration framework (inspired by OpenAI Swarm) with event-driven workflows and composable agent architecture |
| [trpc-agent-go](https://github.com/trpc-group/trpc-agent-go) | From tRPC Group -- hierarchical planners, multi-agent orchestration, persistent memory, and rich tool ecosystem |
| [eino](https://github.com/cloudwego/eino) | Go LLM application framework from CloudWeGo (LangChain-like), providing reusable components: ChatModel, Tool, Retriever, etc. |
| [blades](https://github.com/go-kratos/blades) | Multimodal AI agent framework from go-kratos with pluggable models, tools, memory, and Kratos-style middleware |
| [agentfield](https://github.com/Agent-Field/agentfield) | Kubernetes-inspired infrastructure platform for deploying, scaling, observing, and governing AI agents in production |

### B. Actor / Distributed Actor (high-throughput message-driven scheduling substrate)

| Project | Description |
|---------|-------------|
| [protoactor-go](https://github.com/asynkron/protoactor-go) | Cross-platform actor framework (Go + C#) with Protobuf messaging and virtual actor model |
| [ergo](https://github.com/ergo-services/ergo) | Erlang/OTP patterns in Go (actors, supervision trees, network transparency) for distributed fault-tolerant systems |
| [goakt](https://github.com/Tochemey/goakt) | Distributed Go actor framework with Protobuf messages, designed for scalable high-availability cluster mode |
| [hollywood](https://github.com/anthdm/hollywood) | Ultra-fast, low-latency actor engine -- 10M msg/sec, targeting game servers, trading engines, and ad systems |
| [gam](https://github.com/utrack/gam) | Minimal Go Actor Model with gRPC remoting + Consul clustering, 800k+ msg/sec across nodes |

### C. Durable Workflow Engine (reliable orchestration: timeouts, retries, idempotency, state recovery)

| Project | Description |
|---------|-------------|
| [temporal](https://github.com/temporalio/temporal) | Durable execution platform (successor to Uber Cadence) for fault-tolerant, long-running workflow orchestration |
| [go-workflows](https://github.com/cschleiden/go-workflows) | Temporal-inspired durable workflow engine in pure Go -- write fault-tolerant workflows as ordinary Go code |

### D. Kubernetes Scheduler (placement / resource scheduling plugins)

| Project | Description |
|---------|-------------|
| [scheduler-plugins](https://github.com/kubernetes-sigs/scheduler-plugins) | Official K8s SIG out-of-tree scheduler plugins, production-tested, vendorable as Go SDK or deployable via Helm |
| [kube-scheduler-simulator](https://github.com/kubernetes-sigs/kube-scheduler-simulator) | K8s scheduler simulator -- test custom scheduling configs and plugins without a real cluster, shows per-node plugin decisions |
| [DRS](https://github.com/JolyonJian/DRS) | Deep Reinforcement Learning (DQN) enhanced K8s scheduler for optimizing microservice placement across nodes |

### E. Cluster Membership / Failure Detection (essential membership layer for cross-node MAS)

| Project | Description |
|---------|-------------|
| [memberlist](https://github.com/hashicorp/memberlist) | Hashicorp library for cluster membership management and node failure detection using SWIM gossip protocol |
| [go-grpc-config-gossiping-cluster](https://github.com/davidmontoyago/go-grpc-config-gossiping-cluster) | Demo: eventually consistent config replication across a gRPC service cluster via hashicorp/memberlist gossip |

---

## 2) Top-Venue Papers

### A. Multi-Agent Coordination / Communication Structure

* **NeurIPS 2024 -- SeqComm (Sequential Communication / multi-level communication)** ([NeurIPS][1])

```text
https://neurips.cc/virtual/2024/poster/96719
https://proceedings.neurips.cc/paper_files/paper/2024/file/d6be51e667e0b263e89a23294b57f8cf-Paper-Conference.pdf
https://arxiv.org/abs/2209.12713
```

* **NeurIPS 2024 -- COPPER (Reflective Multi-Agent Collaboration)** ([NeurIPS Proceedings][2])

```text
https://neurips.cc/virtual/2024/poster/93147
https://proceedings.neurips.cc/paper_files/paper/2024/file/fa54b0edce5eef0bb07654e8ee800cb4-Paper-Conference.pdf
```

* **NeurIPS 2024/2025 -- Chain-of-Agents (CoA)** ([OpenReview][3])

```text
https://openreview.net/forum?id=LuCLf4BJsr
https://research.google/blog/chain-of-agents-large-language-models-collaborating-on-long-context-tasks/
```

### B. Multi-Agent Debate and Sparse Communication

* **EMNLP 2024 Findings -- Improving Multi-Agent Debate with Sparse Communication Topology** ([ACL Anthology][4])

```text
https://aclanthology.org/2024.findings-emnlp.427/
https://aclanthology.org/2024.findings-emnlp.427.pdf
https://arxiv.org/abs/2406.11776
```

* **NeurIPS 2025 -- Debate or Vote** ([NeurIPS][5])

```text
https://neurips.cc/virtual/2025/poster/116557
```

### C. Tool-Use / Agent Execution Pipeline (planning -> tools -> feedback)

* **NeurIPS 2023 -- TPTU (Task Planning and Tool Usage)** ([arXiv][6])

```text
https://arxiv.org/abs/2308.03427
https://openreview.net/pdf?id=GrkgKtOjaH
```

* **NeurIPS 2024 -- AVATAR (Optimizing LLM Agents for Tool Usage)** ([NeurIPS][7])

```text
https://neurips.cc/virtual/2024/poster/95465
https://arxiv.org/abs/2406.11200
https://cs.stanford.edu/people/jure/pubs/avatar-neurips24.pdf
```

* **NeurIPS 2025 -- Self-Challenging Language Model Agents** ([NeurIPS][8])

```text
https://neurips.cc/virtual/2025/poster/119495
https://arxiv.org/abs/2506.01716
```

### D. Memory (hierarchical memory evaluation / mechanisms)

* **ICLR 2026 -- MemoryAgentBench (Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions)** ([ICLR][9])
  > Corresponding open-source project: benchmark suite for evaluating LLM agent memory (retrieval, test-time learning, long-range understanding, conflict resolution) via incremental multi-turn interactions

```text
https://iclr.cc/virtual/2026/poster/10010781
https://openreview.net/forum?id=DT7JyQC3MR
https://arxiv.org/abs/2507.05257
https://github.com/HUST-AI-HYZ/MemoryAgentBench
```

* **ICLR 2026 -- MemGen / Weaving Generative Latent Memory** ([OpenReview][10])

```text
https://openreview.net/pdf/aeb0aad4039600443f998b8b2bc72dc61d172605.pdf
https://arxiv.org/pdf/2509.24704
```

### E. Scheduling (system scheduling / topology alignment -- applicable as reference)

* **AAAI 2021 -- Scheduling of Time-Varying Workloads Using Reinforcement Learning** ([AAAI][11])

```text
https://aaai.org/papers/09000-scheduling-of-time-varying-workloads-using-reinforcement-learning/
```

* **NeurIPS 2025 -- Arnold (Topology-Aware Communication Alignment / scheduling on >9600 GPUs)** ([NeurIPS][12])

```text
https://neurips.cc/virtual/2025/poster/115232
https://arxiv.org/abs/2509.15940
```

---

For the next step toward the goal of a **high-speed MAS scheduler across vLLM servers**, the projects above can be evaluated along three axes: **(1) message bus / actor substrate, (2) durable orchestration, (3) K8s placement** -- e.g. `protoactor-go + memberlist + custom vLLM router`, or `Temporal + custom streaming router`, etc.

[1]: https://neurips.cc/virtual/2024/poster/96719 "Multi-Agent Coordination via Multi-Level Communication"
[2]: https://proceedings.neurips.cc/paper_files/paper/2024/hash/fa54b0edce5eef0bb07654e8ee800cb4-Abstract-Conference.html "Reflective Multi-Agent Collaboration based on Large ..."
[3]: https://openreview.net/forum?id=LuCLf4BJsr "Chain of Agents: Large Language Models Collaborating on..."
[4]: https://aclanthology.org/2024.findings-emnlp.427/ "Improving Multi-Agent Debate with Sparse Communication ..."
[5]: https://neurips.cc/virtual/2025/poster/116557 "Debate or Vote: Which Yields Better Decisions in Multi ..."
[6]: https://arxiv.org/abs/2308.03427 "TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"
[7]: https://neurips.cc/virtual/2024/poster/95465 "Optimizing LLM Agents for Tool Usage via Contrastive ..."
[8]: https://neurips.cc/virtual/2025/poster/119495 "Self-Challenging Language Model Agents"
[9]: https://iclr.cc/virtual/2026/poster/10010781 "Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"
[10]: https://openreview.net/pdf/aeb0aad4039600443f998b8b2bc72dc61d172605.pdf "MEMGEN: WEAVING GENERATIVE LATENT MEMORY"
[11]: https://aaai.org/papers/09000-scheduling-of-time-varying-workloads-using-reinforcement-learning/ "Scheduling of Time-Varying Workloads Using ..."
[12]: https://neurips.cc/virtual/2025/poster/115232 "Efficient Pre-Training of LLMs via Topology-Aware ..."
