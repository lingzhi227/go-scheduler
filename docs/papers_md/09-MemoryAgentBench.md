::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_main
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_content
# Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions {#evaluating-memory-in-llm-agents-via-incremental-multi-turn-interactions .ltx_title .ltx_title_document}

::: ltx_authors
[ [Yuanzhe Hu^1^   , Yu Wang^[1∗]{.ltx_text .ltx_font_italic}^, Julian McAuley^1^\
^1^UC San Diego\
^1^[{yuh127,yuw164,jmcauley}@ucsd.edu\
[ ![\[Uncaptioned image\]](/html/2507.05257/assets/x1.png){#g1 .ltx_graphics .ltx_img_square width="27" height="27"}]{.ltx_text .ltx_font_serif style="position:relative; bottom:-2.5pt;"}]{.ltx_text .ltx_font_typewriter} [Datasets](https://huggingface.co/datasets/ai-hyz/MemoryAgentBench){.ltx_ref .ltx_href target="_blank"} [![\[Uncaptioned image\]](/html/2507.05257/assets/x2.png){#g2 .ltx_graphics .ltx_img_square width="20" height="20"}]{.ltx_text style="position:relative; bottom:-1.0pt;"} [Source Code](https://github.com/HUST-AI-HYZ/MemoryAgentBench){.ltx_ref .ltx_href target="_blank"} ]{.ltx_personname}[Y. Hu and Y. Wang contribute equally.]{.ltx_author_notes}]{.ltx_creator .ltx_role_author}
:::

::: ltx_abstract
###### Abstract {#abstract .ltx_title .ltx_title_abstract}

Recent benchmarks for Large Language Model (LLM) agents primarily focus on evaluating reasoning, planning, and execution capabilities, while another critical component---memory, encompassing how agents memorize, update, and retrieve long-term information---is under-evaluated due to the lack of benchmarks. We term agents with memory mechanisms as [memory agents]{.ltx_text .ltx_font_bold}. In this paper, we identify four core competencies essential for memory agents: accurate retrieval, test-time learning, long-range understanding, and conflict resolution. Existing datasets either rely on limited context lengths or are tailored for static, long-context settings like book-based QA, which do not reflect the interactive, multi-turn nature of memory agents that incrementally accumulate information. Furthermore, no existing benchmarks cover all four competencies. Therefore, we introduce [MemoryAgentBench]{.ltx_text .ltx_font_bold}, a new benchmark specifically designed for memory agents. Our benchmark combines reformulated existing datasets with newly constructed ones, covering the above four memory competencies, providing a systematic and challenging testbed for assessing memory quality. We evaluate a diverse set of memory agents, ranging from simple context-based and retrieval-augmented generation (RAG) systems to advanced agents with external memory modules and tool integration. Empirical results reveal that current methods fall short of mastering all four competencies, underscoring the need for further research into comprehensive memory mechanisms for LLM agents.
:::

::::::::: {#S1 .section .ltx_section}
## [1 ]{.ltx_tag .ltx_tag_section}Introduction {#introduction .ltx_title .ltx_title_section}

::: {#S1.p1 .ltx_para}
Large Language Model (LLM) agents have rapidly transitioned from proof-of-concept chatbots to end-to-end systems that can write software (Wang et al., [2024c](#bib.bib42){.ltx_ref}), control browsers (Müller and Žunič, [2024](#bib.bib34){.ltx_ref}), and reason over multi-modal inputs. Frameworks such as [Manus]{.ltx_text .ltx_font_smallcaps}, [OWL]{.ltx_text .ltx_font_smallcaps} (Hu et al., [2025](#bib.bib16){.ltx_ref}), [OpenHands]{.ltx_text .ltx_font_smallcaps} (Wang et al., [2024c](#bib.bib42){.ltx_ref}), and [Codex]{.ltx_text .ltx_font_smallcaps} routinely solve complex, tool-rich tasks and achieve state-of-the-art results on agentic benchmarks like GAIA (Mialon et al., [2023](#bib.bib31){.ltx_ref}) and SWE-Bench (Jimenez et al., [2023](#bib.bib18){.ltx_ref}). Yet these evaluations focus almost exclusively on *reasoning* (planning, tool using, code synthesis) and leave the equally important question of *memorization* (abstraction, storing, updating, retrieving) largely under-explored.
:::

::: {#S1.p2 .ltx_para}
Recent memory-centric architectures---ranging from parametric memory systems like MemoryLLM (Wang et al., [2024d](#bib.bib44){.ltx_ref}), SELF-PARAM ([Wang et al.,](#bib.bib43){.ltx_ref} ), and M+(Wang et al., [2025](#bib.bib46){.ltx_ref}) to commercial token-level memory solutions such as [MemGPT[(]{.ltx_text .ltx_font_upright}Packer et al.[, ]{.ltx_text .ltx_font_upright}[2023](#bib.bib36){.ltx_ref}; Lin et al.[, ]{.ltx_text .ltx_font_upright}[2025](#bib.bib27){.ltx_ref}[)]{.ltx_text .ltx_font_upright}]{.ltx_text .ltx_font_smallcaps}, [Mem0[(]{.ltx_text .ltx_font_upright}Chhikara et al.[, ]{.ltx_text .ltx_font_upright}[2025](#bib.bib7){.ltx_ref}[)]{.ltx_text .ltx_font_upright}]{.ltx_text .ltx_font_smallcaps}, [Cognee]{.ltx_text .ltx_font_smallcaps}, and [Zep[(]{.ltx_text .ltx_font_upright}Rasmussen et al.[, ]{.ltx_text .ltx_font_upright}[2025](#bib.bib37){.ltx_ref}[)]{.ltx_text .ltx_font_upright}]{.ltx_text .ltx_font_smallcaps}---employ diverse strategies for storing and retrieving past information. Despite growing interest, their real-world effectiveness remains largely anecdotal, and there is currently no unified benchmark for systematically evaluating the quality of memory in agents. In this paper, we refer to agents equipped with memory mechanisms as [Memory Agents]{.ltx_text .ltx_font_bold}, where memory can take various forms, including parameters, vectors, textual histories, or external databases. In this paper, we primarily focus on memory agents that utilize textual histories and external databases, as these approaches are most commonly deployed in real-world applications. In contrast, memory encoded in model parameters (Wang et al., [2024d](#bib.bib44){.ltx_ref}, [2025](#bib.bib46){.ltx_ref}; Yin et al., [2024](#bib.bib51){.ltx_ref}) remains largely within academic research and is typically less capable than proprietary memory systems equipped on closed-sourced API models.
:::

::: {#S1.p3 .ltx_para}
To evaluate memory agents, we identify four complementary competencies (Examples shown in Figure [[1]{.ltx_text .ltx_ref_tag}](#S1.F1 "Figure 1 ‣ 1 Introduction ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}): (1) [Accurate Retrieval]{.ltx_text .ltx_font_bold}: The ability to extract the correct snippet in response to a query. This can involve one-hop or multi-hop retrieval, as long as the relevant information can be accessed with a single query. (2) [Test-Time Learning]{.ltx_text .ltx_font_bold}: The capacity to incorporate new behaviors or acquire new skills during deployment, without additional training. (3) [Long-Range Understanding]{.ltx_text .ltx_font_bold}: The ability to integrate information distributed across extended contexts ($\geq$ 100k tokens) and answer questions requiring a global understanding of the entire sequence. (4) [Conflict Resolution]{.ltx_text .ltx_font_bold}: The skill to revise, overwrite, or remove previously stored information when faced with contradictory evidence, aligning with goals in model editing and knowledge unlearning tasks (Meng et al., [2023](#bib.bib30){.ltx_ref}; Wang et al., [2024e](#bib.bib45){.ltx_ref}).
:::

<figure id="S1.F1" class="ltx_figure">
<img src="/html/2507.05257/assets/figures/intro.png" id="S1.F1.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="586" height="319" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 1</span>: </span><span class="ltx_text" style="font-size:90%;">Four complementary competencies that memory agents should have.</span></figcaption>
</figure>

::: {#S1.p4 .ltx_para}
Previous datasets developed to evaluate memory in language models have notable limitations. Early benchmarks such as [LOCOMO]{.ltx_text .ltx_font_smallcaps} (Maharana et al., [2024](#bib.bib29){.ltx_ref}) ($\sim$ 9k tokens), LooGLE(Li et al., [2023](#bib.bib23){.ltx_ref}) ($\sim$ 24k tokens), and LongBench(Bai et al., [2023](#bib.bib3){.ltx_ref}) ($\sim$ 20k tokens) feature relatively short contexts that no longer challenge current models. More recent datasets like NovelQA(Wang et al., [2024a](#bib.bib40){.ltx_ref}) ($\sim$`<!-- -->`{=html}200k tokens), NOCHA(Karpinska et al., [2024](#bib.bib19){.ltx_ref}) ($\sim$`<!-- -->`{=html}127k tokens), Loong(Wang et al., [2024b](#bib.bib41){.ltx_ref}) ($\sim$`<!-- -->`{=html}100k tokens), and $\infty$-Bench(Zhang et al., [2024](#bib.bib53){.ltx_ref}) ($\sim$`<!-- -->`{=html}150k tokens) extend the context length to evaluate global reasoning and retrieval capabilities. However, these datasets were primarily designed for evaluating long-context language models rather than memory agents. The reason that long-context benchmarks cannot be directly used to evaluate memory agents is as follows. There is a fundamental distinction between memory and long context: memory serves as a compressed and distilled representation of past information. Rather than storing all historical content verbatim, memory selectively extracts salient details, removes irrelevant information, and often incorporates new inferences derived from prior experiences. Consequently, [memory agents are designed to process context incrementally]{.ltx_text .ltx_font_bold}---absorbing input piece by piece, abstracting and consolidating information over time, generating new inferences, and learning novel rules from accumulated history. For this reason, datasets that provide the entire context in a single block are not directly applicable to evaluating memory agents. A more recent effort, [LongMemEval]{.ltx_text .ltx_font_smallcaps} (Wu et al., [2024](#bib.bib48){.ltx_ref}), seeks to address this limitation by using synthetic long-form conversations, which can be injected into memory gradually, session by session. Nonetheless, its evaluation framework remains constrained by limited topical diversity and less realistic interaction patterns, reducing its applicability to real-world memory agent scenarios.
:::

::: {#S1.p5 .ltx_para}
To address these limitations, we introduce a unified benchmark framework, [MemoryAgentBench]{.ltx_text .ltx_font_bold}, specifically designed to evaluate a broad spectrum of memory mechanisms in agent systems. We also provide a framework for memory agent evaluation. In this framework, agents are presented with sequences of textual inputs that simulate multi-turn interactions with users. We repurpose existing datasets originally developed for long-context LLM evaluation by segmenting their inputs into multiple chunks and feeding them incrementally to the agent. However, since these datasets do not fully capture all four targeted memory competencies, we also introduce two new datasets: [EventQA]{.ltx_text .ltx_font_bold} and [FactConsolidation]{.ltx_text .ltx_font_bold}, designed to evaluate accurate retrieval and conflict resolution, respectively. Our benchmark includes evaluations of state-of-the-art commercial memory agents (such as Mem0 and MemGPT), long-context agents that treat the full input as memory, and RAG agents that extend their memory through retrieval methods. We examine how techniques developed for long-context models and RAG transfer to the memory agent setting, and how commercial memory agents perform under more challenging, competency-specific tests. By providing a consistent evaluation protocol across diverse agent architectures and datasets, [MemoryAgentBench]{.ltx_text .ltx_font_bold} delivers comprehensive insights into agent performance across the four core memory competencies.
:::

::: {#S1.p6 .ltx_para}
Our contributions are summarized as follows:

- [[•]{.ltx_tag .ltx_tag_item}]{#S1.I1.i1}

  ::: {#S1.I1.i1.p1 .ltx_para}
  [Datasets:]{.ltx_text .ltx_font_bold} We re-structure existing datasets and create two new datasets to construct a comprehensive benchmark, covering four distinct memory competencies.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S1.I1.i2}

  ::: {#S1.I1.i2.p1 .ltx_para}
  [Framework:]{.ltx_text .ltx_font_bold} We provide a unified evaluation framework, and open-source the codebase and datasets to encourage reproducibility and further research.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S1.I1.i3}

  ::: {#S1.I1.i3.p1 .ltx_para}
  [Empirical Study:]{.ltx_text .ltx_font_bold} We implement various simple agents with diverse memory mechanisms, adopt commercial agents, and evaluate these agents on our proposed benchmark. With our results, we show that existing memory agents, while effective in some tasks, still face significant challenges on some aspects.
  :::
:::
:::::::::

:::::::: {#S2 .section .ltx_section}
## [2 ]{.ltx_tag .ltx_tag_section}Related Work {#related-work .ltx_title .ltx_title_section}

:::: {#S2.SS1 .section .ltx_subsection}
### [2.1 ]{.ltx_tag .ltx_tag_subsection}Benchmarks with Long Input {#benchmarks-with-long-input .ltx_title .ltx_title_subsection}

::: {#S2.SS1.p1 .ltx_para}
In this section, we review prior work on long-context benchmarks. Early benchmarks designed for long-context evaluation include LongBench(Bai et al., [2023](#bib.bib3){.ltx_ref}) and LooGLE(Li et al., [2023](#bib.bib23){.ltx_ref}), with average input lengths of approximately 20k and 24k tokens, respectively. More recent benchmarks---such as $\infty$-Bench (Zhang et al., [2024](#bib.bib53){.ltx_ref}), HELMET(Yen et al., [2024](#bib.bib50){.ltx_ref}), RULER(Hsieh et al., [2024](#bib.bib15){.ltx_ref}), NOCHA(Karpinska et al., [2024](#bib.bib19){.ltx_ref}), NoLiMa (Modarressi et al., [2025](#bib.bib33){.ltx_ref}) and LongBench V2(Bai et al., [2024](#bib.bib4){.ltx_ref})---extend context lengths to over 100k tokens and are primarily intended to evaluate the capabilities of long-context models. However, despite their scale, these benchmarks are not designed to assess memory agents, and no prior work has repurposed them for that goal. More recently, LOCOMO (Maharana et al., [2024](#bib.bib29){.ltx_ref}) and LongMemEval (Wu et al., [2024](#bib.bib48){.ltx_ref}) have been proposed specifically for evaluating memory agents. While promising, LOCOMO still features relatively short conversations ($\sim$`<!-- -->`{=html}9k), and LongMemEval uses synthetic conversations with limited topical diversity, making the dialogues less realistic and potentially less representative of real-world memory use cases.
:::
::::

::::: {#S2.SS2 .section .ltx_subsection}
### [2.2 ]{.ltx_tag .ltx_tag_subsection}Agents with Memory Mechanisms {#agents-with-memory-mechanisms .ltx_title .ltx_title_subsection}

::: {#S2.SS2.p1 .ltx_para}
Memory mechanisms are attracting more and more attention lately (Wang et al., [2025/02](#bib.bib47){.ltx_ref}). Recent advancements in LLMs have demonstrated the capability to process extended context lengths, ranging from 100K to over 1 million tokens. For instance, models such as GPT-4o (OpenAI, [2025](#bib.bib35){.ltx_ref}) and Claude 3.7 (Anthropic, [2025](#bib.bib1){.ltx_ref}) can handle inputs of approximately 100K to 200K tokens, while models like Gemini 2.0 Pro (DeepMind, [2025](#bib.bib8){.ltx_ref}) and the GPT-4.1 series extend this capacity beyond 1 million tokens. These strong long-context capabilities enable a simple yet effective form of memory: storing information directly within the context window. However, this approach is inherently constrained by a hard limit---once the context window is exceeded, earlier information must be discarded.
:::

::: {#S2.SS2.p2 .ltx_para}
In parallel, RAG continues to serve as a dominant paradigm for managing excessive context. By retrieving relevant information from earlier context and feeding it to the LLM, RAG allows systems to overcome context length limitations. For example, OpenAI's recent memory functionality[^1^[[^1^[1]{.ltx_tag .ltx_tag_note}[https://openai.com/index/memory-and-new-controls-for-chatgpt/](https://openai.com/index/memory-and-new-controls-for-chatgpt/){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}]{.ltx_note_content}]{.ltx_note_outer}]{#footnote1 .ltx_note .ltx_role_footnote} combines explicit user preference tracking with retrieval-based methods that reference prior interactions. RAG methods can be broadly classified into three categories: [1. Simple RAG]{.ltx_text .ltx_font_bold}: These methods rely on string-matching techniques such as TF-IDF, BM25 (Robertson and Walker, [1994](#bib.bib38){.ltx_ref}), and BMX (Li et al., [2024](#bib.bib25){.ltx_ref}), which are entirely non-neural and operate on string-level similarity. [2. Embedding-based RAG]{.ltx_text .ltx_font_bold}: This class leverages neural encoders, primarily transformers, to map text into dense vector representations (Wu et al., [2022](#bib.bib49){.ltx_ref}). Early methods like DPR (Karpukhin et al., [2020](#bib.bib20){.ltx_ref}) and Contriever (Izacard et al., [2021](#bib.bib17){.ltx_ref}) are based on BERT (Devlin et al., [2019](#bib.bib9){.ltx_ref}), while more recent models such as NV-Embed-v2 (Lee et al., [2024](#bib.bib22){.ltx_ref}) utilize decoder-only backbones and achieve significantly improved retrieval performance. [3. Structure-Augmented RAG]{.ltx_text .ltx_font_bold}: These approaches enhance retrieval with structural representations such as graphs or trees. Representative systems include GraphRAG (Edge et al., [2024](#bib.bib10){.ltx_ref}), RAPTOR (Sarthi et al., [2024](#bib.bib39){.ltx_ref}), HippoRAG-V2 (Gutiérrez et al., [2025](#bib.bib12){.ltx_ref}), Cognee, Zep (Rasmussen et al., [2025](#bib.bib37){.ltx_ref}), and Mem0 (Chhikara et al., [2025](#bib.bib7){.ltx_ref}), where Mem0 also offers a graph-augmented variant, Mem0^[g]{.ltx_text .ltx_font_italic}^, built on structured factual knowledge. Despite their effectiveness, RAG-based methods face challenges with ambiguous queries, multi-hop reasoning, and long-range comprehension. When questions require integrating knowledge across an entire session or learning from long, skill-encoding inputs, the retrieval mechanism---limited to the top-k most relevant passages---may fail to surface the necessary information. To address these limitations, Agentic Memory Agents introduce an iterative, decision-driven framework. Rather than relying on a single-pass retrieval, these agents dynamically process the query, retrieve evidence, reflect, and iterate through multiple retrieval and reasoning cycles. Examples include MemGPT (Packer et al., [2023](#bib.bib36){.ltx_ref}), Self-RAG (Asai et al., [2023](#bib.bib2){.ltx_ref}), and Auto-RAG (Yu et al., [2024](#bib.bib52){.ltx_ref}). This agentic design is particularly effective for resolving ambiguous or multi-step queries. Nonetheless, these methods remain fundamentally constrained by the limitations of RAG---namely, the inability to fully understand or learn from long-range context that is inaccessible via retrieval alone.
:::
:::::
::::::::

::::::::::::::::::::::::::::::::::::: {#S3 .section .ltx_section}
## [3 ]{.ltx_tag .ltx_tag_section}MemoryAgentBench {#memoryagentbench .ltx_title .ltx_title_section}

:::::::::::: {#S3.SS1 .section .ltx_subsection}
### [3.1 ]{.ltx_tag .ltx_tag_subsection}Aspects of the Evaluation {#aspects-of-the-evaluation .ltx_title .ltx_title_subsection}

::: {#S3.SS1.p1 .ltx_para}
The evaluation of memory agents encompasses the following key dimensions:
:::

:::: {#S3.SS1.SSS0.Px1 .section .ltx_paragraph}
##### Accurate Retrieval (AR) {#accurate-retrieval-ar .ltx_title .ltx_title_paragraph}

::: {#S3.SS1.SSS0.Px1.p1 .ltx_para}
The task of accurately retrieving information has been extensively explored in prior work. In the domain of long-context modeling, the Needle-in-a-Haystack (NIAH) task is widely used to evaluate a model's ability to locate a specific value based on a given key within a lengthy input. Extensions such as multi-value NIAH further test the model's capacity to retrieve multiple values scattered across the input context. In the RAG setting, this corresponds to document-based QA, where the model must identify and extract relevant snippets from one or more documents to answer a query. These snippets might reside in a single location or be distributed across multiple documents. In this paper, we focus on agentic settings, where the "long context" or "multiple documents" become long-form conversations. We define Accurate Retrieval (AR) as the ability of an agent to identify and retrieve important information that may be dispersed throughout a long dialogue history.
:::
::::

:::: {#S3.SS1.SSS0.Px2 .section .ltx_paragraph}
##### Test-Time Learning (TTL) {#test-time-learning-ttl .ltx_title .ltx_title_paragraph}

::: {#S3.SS1.SSS0.Px2.p1 .ltx_para}
An essential capability for real-world agents is the ability to acquire new skills dynamically through interaction with users. This mirrors the concept of In-Context Learning (ICL) in LLMs, where the model learns from a prompt containing a small number of examples, often framed as few-shot classification tasks. Ideally, performance improves with additional examples in the prompt. In the conversational agent setting, prompts are replaced by dialogue histories. We define Test-Time Learning (TTL) as the agent's ability to learn to perform new tasks directly from the conversation. This property is crucial for enabling self-evolving agents that can continuously adapt and improve in real-world deployments.
:::
::::

:::: {#S3.SS1.SSS0.Px3 .section .ltx_paragraph}
##### Long-Range Understanding (LRU) {#long-range-understanding-lru .ltx_title .ltx_title_paragraph}

::: {#S3.SS1.SSS0.Px3.p1 .ltx_para}
Long-range understanding refers to the agent's ability to form abstract, high-level comprehension over extended conversations. For example, when a user narrates a long story, the agent should retain the content and derive a holistic understanding rather than just recall isolated facts. We define Long-Range Understanding (LRU) as the ability to reason about long-form inputs and answer high-level questions that require an understanding of the overall content, rather than detailed recall. An example question might be: "Summarize the main experiences of Harry Potter."
:::
::::

:::: {#S3.SS1.SSS0.Px4 .section .ltx_paragraph}
##### Conflict Resolution (CR) {#conflict-resolution-cr .ltx_title .ltx_title_paragraph}

::: {#S3.SS1.SSS0.Px4.p1 .ltx_para}
In long-term interactions, agents often face evolving or conflicting information---whether about the external world (e.g., changes in political leadership) or user-specific facts (e.g., a new occupation). This challenge is closely related to model editing (Meng et al., [2023](#bib.bib30){.ltx_ref}; Fang et al., [2024](#bib.bib11){.ltx_ref}) and knowledge unlearning (Wang et al., [2024e](#bib.bib45){.ltx_ref}), which focus on modifying or removing factual knowledge from language models. We define Conflict Resolution (CR) as the agent's ability to detect and resolve contradictions between existing knowledge and newly acquired information, ensuring the agent remains aligned with current realities and user states. CR is distinct from Abstractive Retrieval (AR) in two key ways. (1) Certain questions requiring CR cannot be answered solely through AR. As illustrated in Figure [[1]{.ltx_text .ltx_ref_tag}](#S1.F1 "Figure 1 ‣ 1 Introduction ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}, an agent that retrieves all facts related to [pears]{.ltx_text .ltx_font_typewriter} may fail to identify the updated information in the second message. (2) In AR, earlier messages remain relevant and should be retained, even when multiple pieces of evidence are required. In contrast, CR involves identifying outdated or incorrect information and discarding it. That is, AR requires preservation of all related content, whereas CR requires overwriting prior facts to reflect the most up-to-date truth.
:::

<figure id="S3.T1" class="ltx_table">
<table class="ltx_tabular ltx_centering ltx_align_middle">
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_tt"><span class="ltx_text ltx_font_bold" style="font-size:90%;">Capability</span></td>
<td class="ltx_td ltx_align_justify ltx_border_tt"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="ltx_text ltx_font_bold" style="font-size:90%;">Benchmarks / Tasks</span></span> </span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold" style="font-size:90%;"># of Sequences</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold" style="font-size:90%;"># of QAs</span></td>
<td class="ltx_td ltx_nopad_r ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold" style="font-size:90%;">Avg Len</span></td>
</tr>
<tr class="ltx_tr">
<td rowspan="5" class="ltx_td ltx_align_left ltx_border_t"><span class="ltx_text" style="font-size:90%;"> <span class="ltx_inline-block ltx_parbox ltx_align_middle" style="width:113.8pt;"> <span class="ltx_p">Accurate</span> <span class="ltx_p">Retrieval</span> </span></span></td>
<td class="ltx_td ltx_align_justify ltx_border_t"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="ltx_text" style="font-size:90%;">RULER-QA </span><span class="ltx_text" style="font-size:90%;">(</span>Hsieh et al.<span class="ltx_text" style="font-size:90%;">, </span><a href="#bib.bib15" class="ltx_ref">2024</a><span class="ltx_text" style="font-size:90%;">)</span></span> </span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">2</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">200</span></td>
<td class="ltx_td ltx_nopad_r ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">309K</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_justify"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="ltx_text" style="font-size:90%;">RULER-NIAH-MQ </span><span class="ltx_text" style="font-size:90%;">(</span>Hsieh et al.<span class="ltx_text" style="font-size:90%;">, </span><a href="#bib.bib15" class="ltx_ref">2024</a><span class="ltx_text" style="font-size:90%;">)</span></span> </span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">1</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">100</span></td>
<td class="ltx_td ltx_nopad_r ltx_align_center"><span class="ltx_text" style="font-size:90%;">448K</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_justify"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="math inline">∞</span><span class="ltx_text" style="font-size:90%;">Bench-QA </span><span class="ltx_text" style="font-size:90%;">(</span>Zhang et al.<span class="ltx_text" style="font-size:90%;">, </span><a href="#bib.bib53" class="ltx_ref">2024</a><span class="ltx_text" style="font-size:90%;">)</span></span> </span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">100</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">100</span></td>
<td class="ltx_td ltx_nopad_r ltx_align_center"><span class="ltx_text" style="font-size:90%;">183K</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_justify"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="ltx_text" style="font-size:90%;">LongMemEval (S*) </span><span class="ltx_text" style="font-size:90%;">(</span>Wu et al.<span class="ltx_text" style="font-size:90%;">, </span><a href="#bib.bib48" class="ltx_ref">2024</a><span class="ltx_text" style="font-size:90%;">)</span></span> </span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">5</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">300</span></td>
<td class="ltx_td ltx_nopad_r ltx_align_center"><span class="ltx_text" style="font-size:90%;">355K</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_justify" style="padding-bottom: 2.0pt"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="ltx_text" style="font-size:90%;">EventQA (</span><em>ours</em><span class="ltx_text" style="font-size:90%;">)</span></span> </span></td>
<td class="ltx_td ltx_align_center" style="padding-bottom: 2.0pt"><span class="ltx_text" style="font-size:90%;">5</span></td>
<td class="ltx_td ltx_align_center" style="padding-bottom: 2.0pt"><span class="ltx_text" style="font-size:90%;">500</span></td>
<td class="ltx_td ltx_nopad_r ltx_align_center" style="padding-bottom: 2.0pt"><span class="ltx_text" style="font-size:90%;">534K</span></td>
</tr>
<tr class="ltx_tr">
<td rowspan="5" class="ltx_td ltx_align_left ltx_border_t"><span class="ltx_text" style="font-size:90%;"> <span class="ltx_inline-block ltx_parbox ltx_align_middle" style="width:113.8pt;"> <span class="ltx_p">Test-Time</span> <span class="ltx_p">Learning</span> </span></span></td>
<td class="ltx_td ltx_align_justify ltx_border_t"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="ltx_text" style="font-size:90%;">BANKING-77</span></span> </span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">1</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">100</span></td>
<td class="ltx_td ltx_nopad_r ltx_border_t"></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_justify"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="ltx_text" style="font-size:90%;">CLINC-150</span></span> </span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">1</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">100</span></td>
<td class="ltx_td ltx_nopad_r"></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_justify"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="ltx_text" style="font-size:90%;">NLU</span></span> </span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">1</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">100</span></td>
<td class="ltx_td ltx_nopad_r ltx_align_center"><span class="ltx_text" style="font-size:90%;">103K</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_justify"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="ltx_text" style="font-size:90%;">TREC (Coarse)</span></span> </span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">1</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">100</span></td>
<td class="ltx_td ltx_nopad_r"></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_justify"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="ltx_text" style="font-size:90%;">TREC (Fine)</span></span> </span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">1</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">100</span></td>
<td class="ltx_td ltx_nopad_r"></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td" style="padding-bottom: 2.0pt"></td>
<td class="ltx_td ltx_align_justify" style="padding-bottom: 2.0pt"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="ltx_text" style="font-size:90%;">Movie-Rec Redial </span><span class="ltx_text" style="font-size:90%;">(</span>He et al.<span class="ltx_text" style="font-size:90%;">, </span><a href="#bib.bib14" class="ltx_ref">2023b</a><span class="ltx_text" style="font-size:90%;">)</span></span> </span></td>
<td class="ltx_td ltx_align_center" style="padding-bottom: 2.0pt"><span class="ltx_text" style="font-size:90%;">1</span></td>
<td class="ltx_td ltx_align_center" style="padding-bottom: 2.0pt"><span class="ltx_text" style="font-size:90%;">200</span></td>
<td class="ltx_td ltx_nopad_r ltx_align_center" style="padding-bottom: 2.0pt"><span class="ltx_text" style="font-size:90%;">1.44M</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_t" style="padding-bottom: 2.0pt"><span class="ltx_text" style="font-size:90%;"> <span class="ltx_inline-block ltx_parbox ltx_align_middle" style="width:113.8pt;"> <span class="ltx_p">Long-Range Understanding</span> </span></span></td>
<td class="ltx_td ltx_align_justify ltx_border_t" style="padding-bottom: 2.0pt"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="math inline">∞</span><span class="ltx_text" style="font-size:90%;">Bench-Sum </span><span class="ltx_text" style="font-size:90%;">(</span>Zhang et al.<span class="ltx_text" style="font-size:90%;">, </span><a href="#bib.bib53" class="ltx_ref">2024</a><span class="ltx_text" style="font-size:90%;">)</span></span> </span></td>
<td class="ltx_td ltx_align_center ltx_border_t" style="padding-bottom: 2.0pt"><span class="ltx_text" style="font-size:90%;">100</span></td>
<td class="ltx_td ltx_align_center ltx_border_t" style="padding-bottom: 2.0pt"><span class="ltx_text" style="font-size:90%;">100</span></td>
<td class="ltx_td ltx_nopad_r ltx_align_center ltx_border_t" style="padding-bottom: 2.0pt"><span class="ltx_text" style="font-size:90%;">172K</span></td>
</tr>
<tr class="ltx_tr">
<td rowspan="2" class="ltx_td ltx_align_left ltx_border_bb ltx_border_t"><span class="ltx_text" style="font-size:90%;"> <span class="ltx_inline-block ltx_parbox ltx_align_middle" style="width:113.8pt;"> <span class="ltx_p">Conflict Resolution</span> </span></span></td>
<td class="ltx_td ltx_align_justify ltx_border_t"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="ltx_text" style="font-size:90%;">FactConsolidation-SH (</span><em>ours</em><span class="ltx_text" style="font-size:90%;">)</span></span> </span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">1</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">100</span></td>
<td rowspan="2" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_bb ltx_border_t"><span class="ltx_text" style="font-size:90%;">262K</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_justify ltx_border_bb"><span class="ltx_inline-block ltx_align_top"> <span class="ltx_p"><span class="ltx_text" style="font-size:90%;">FactConsolidation-MH (</span><em>ours</em><span class="ltx_text" style="font-size:90%;">)</span></span> </span></td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text" style="font-size:90%;">1</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text" style="font-size:90%;">100</span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 1: </span>Datasets categorized by the specific aspects of evaluation.</figcaption>
</figure>
::::
::::::::::::

:::::::::::: {#S3.SS2 .section .ltx_subsection}
### [3.2 ]{.ltx_tag .ltx_tag_subsection}Dataset Preperation {#dataset-preperation .ltx_title .ltx_title_subsection}

::: {#S3.SS2.p1 .ltx_para}
In this section, we describe how we adopt existing datasets and construct new ones for evaluating each aspect introduced in Section [[3.1]{.ltx_text .ltx_ref_tag}](#S3.SS1 "3.1 Aspects of the Evaluation ‣ 3 MemoryAgentBench ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}. All datasets with their categories are shown in Table [[1]{.ltx_text .ltx_ref_tag}](#S3.T1 "Table 1 ‣ Conflict Resolution (CR) ‣ 3.1 Aspects of the Evaluation ‣ 3 MemoryAgentBench ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}. We introduce the details in datasets curation in Appendix [[A]{.ltx_text .ltx_ref_tag}](#A1 "Appendix A Details of Dataset ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}.
:::

:::: {#S3.SS2.SSS0.Px1 .section .ltx_paragraph}
##### Datasets for Accurate Retrieval (AR) {#datasets-for-accurate-retrieval-ar .ltx_title .ltx_title_paragraph}

::: {#S3.SS2.SSS0.Px1.p1 .ltx_para}
We adopt five datasets to evaluate the accurate retrieval capability of memory agents. Four are adapted from existing benchmarks, and one is newly constructed: (1) [RULER-QA]{.ltx_text .ltx_font_bold}: This is a NIAH-style QA task where a long passage contains single (QA-1) or multiple (QA-2) snippets answering the input question. The agent must identify and extract relevant snippets from the extended context. (2) [NIAH-MQ]{.ltx_text .ltx_font_bold}: We use the multiple-query (MQ) version of the NIAH dataset from RULER (Hsieh et al., [2024](#bib.bib15){.ltx_ref}), where each query seeks a different numeric value embedded in a long passage. The agent must accurately retrieve multiple distinct answers. (3) $\infty$[Bench-En.QA]{.ltx_text .ltx_font_bold}: This task from $\infty$Bench presents free-form QA questions based on entire books, with all entities replaced by fictitious names to avoid contamination from model pretraining. Compared to synthetic datasets like RULER-QA, this benchmark is more realistic and challenging due to the natural narrative structure of books. (4) [LongMemEval]{.ltx_text .ltx_font_bold}: This benchmark evaluates memory agents on long dialogue histories. Although task types like information extraction (IE) or multi-session reasoning are included, most tasks can be reformulated as single-retrieval problems requiring agents to retrieve the correct segments spanning a long multi-turn conversation. Among these, LongMemEval is already formatted for agent-based evaluation with session separation. We use the original LongMemEval(S) dataset ($\sim$`<!-- -->`{=html}110K tokens) and reformulated chat history into five long dialogues ($\sim$`<!-- -->`{=html}355K tokens) with 300 questions (LongMemEval (S\*) in Table  [[1]{.ltx_text .ltx_ref_tag}](#S3.T1 "Table 1 ‣ Conflict Resolution (CR) ‣ 3.1 Aspects of the Evaluation ‣ 3 MemoryAgentBench ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}). We create LongMemEval (S\*) specifically for increasing the number of questions per context, mitigating the exhaustive needs of reconstructing the memory for each quesiton. (5) [EventQA (ours)]{.ltx_text .ltx_font_bold}: We introduce EventQA this reasoning style NIAH task to evaluate agents' ability to recall and reason about temporal sequences in long-form narratives. In this dataset, the agent is required to read a novel and select the correct event from a series of candidates after receiving up-to five previous events. For these datasets, which are originally designed for long-context modeling, we split documents into chunks and sequentially inject them into the agent.
:::
::::

:::: {#S3.SS2.SSS0.Px2 .section .ltx_paragraph}
##### Datasets for Test-Time Learning (TTL) {#datasets-for-test-time-learning-ttl .ltx_title .ltx_title_paragraph}

::: {#S3.SS2.SSS0.Px2.p1 .ltx_para}
We evaluate TTL via two task categories: (1) [Multi-Class Classification (MCC)]{.ltx_text .ltx_font_bold}: We adopt five classification datasets used in prior TTL work (Bertsch et al., [2024](#bib.bib5){.ltx_ref}; Yen et al., [2024](#bib.bib50){.ltx_ref}): BANKING77 (Casanueva et al., [2020](#bib.bib6){.ltx_ref}), CLINC150 (Larson et al., [2019](#bib.bib21){.ltx_ref}), TREC-Coarse, TREC-Fine (Li and Roth, [2002](#bib.bib26){.ltx_ref}), and NLU (Liu et al., [2019](#bib.bib28){.ltx_ref}). Each task requires the agent to map sentences to class labels, leveraging previously seen labeled examples in context. (2) [Recommendation (Recom)]{.ltx_text .ltx_font_bold}: We use the Redial (Li et al., [2018](#bib.bib24){.ltx_ref}) dataset to evaluate movie recommendation via dialogue. Following the setup from He et al. ([2023a](#bib.bib13){.ltx_ref}), the agent is exposed to thousands of movie-related dialogue turns and is asked to recommend twenty relevant movies based on the long interaction history.
:::
::::

:::: {#S3.SS2.SSS0.Px3 .section .ltx_paragraph}
##### Datasets for Long Range Understanding (LRU) {#datasets-for-long-range-understanding-lru .ltx_title .ltx_title_paragraph}

::: {#S3.SS2.SSS0.Px3.p1 .ltx_para}
We adopt the Summarization task [En.Sum]{.ltx_text .ltx_font_typewriter} from $\infty$-Bench (Zhang et al., [2024](#bib.bib53){.ltx_ref}). The agent is required to analyze and organize the plot and characters of the novel, and then compose a summary of 1000 to 1200 words.
:::
::::

:::: {#S3.SS2.SSS0.Px4 .section .ltx_paragraph}
##### Datasets for Conflict Resolution (CR) {#datasets-for-conflict-resolution-cr .ltx_title .ltx_title_paragraph}

::: {#S3.SS2.SSS0.Px4.p1 .ltx_para}
To assess whether an agent can consolidate conflicting factual updates and reason over them, we construct a new dataset called FactConsolidation. Specifically, We build this benchmark using counterfactual edit pairs from [MQUAKE]{.ltx_text .ltx_font_smallcaps} (Zhong et al., [2023](#bib.bib54){.ltx_ref}). Each pair contains a true fact and a rewritten, contradictory version. These are ordered such that the rewritten (new) fact appears after the original, simulating a realistic update scenario. We concatenate multiple such edit pairs to create long contexts of length 32K, 64K, 262K. We then adpot MQUAKE's original questions and categorize them into: (1) [FactConsolidation-SH (Ours)]{.ltx_text .ltx_font_bold} (SH means Single-Hop), requiring direct factual recall (e.g., "Which country was tool $A$ created in?"), and (2) [FactConsolidation-MH (Ours)]{.ltx_text .ltx_font_bold} (MH refers to Multi-Hop), requiring inference over multiple facts (e.g., "What is the location of death of the spouse of person $B$?"). Agents are prompted to prioritize later information in case of conflict and reason based on the final memory state. This setup directly evaluates the strength and consistency of conflict resolution over long sequences.
:::
::::
::::::::::::

:::::::::: {#S3.SS3 .section .ltx_subsection}
### [3.3 ]{.ltx_tag .ltx_tag_subsection}Different Categories of Memory Agents {#different-categories-of-memory-agents .ltx_title .ltx_title_subsection}

::: {#S3.SS3.p1 .ltx_para}
We evaluate three major types of memory agents that reflect common strategies for handling long-term information: *Long-Context Agents*, *RAG Agents*, and *Agentic Memory Agents*. These approaches differ in how they store, retrieve, and reason over past inputs.
:::

:::: {#S3.SS3.SSS0.Px1 .section .ltx_paragraph}
##### *Long Context Agents* {#long-context-agents .ltx_title .ltx_title_paragraph}

::: {#S3.SS3.SSS0.Px1.p1 .ltx_para}
Modern language models often support extended context windows ranging from 128K to over 1M tokens. A straightforward strategy for memory is to maintain a context buffer of the most recent tokens. For example, in a model with a 128K-token limit, the agent concatenates all incoming chunks until the total exceeds the window size. Once the limit is reached, the earliest chunks are evicted in a FIFO (first-in, first-out) manner. This agent design relies solely on positional recency and assumes the model can attend effectively over the current context window.
:::
::::

:::: {#S3.SS3.SSS0.Px2 .section .ltx_paragraph}
##### *RAG Agents* {#rag-agents .ltx_title .ltx_title_paragraph}

::: {#S3.SS3.SSS0.Px2.p1 .ltx_para}
RAG-based agents address context limitations by storing past information in an external memory pool and retrieving relevant content as needed. We consider three RAG variants: (1) *Simple RAG Agents*: All input chunks are stored as raw text. During inference, a keyword or rule-based string matching mechanism retrieves relevant passages. (2) *Embedding-based RAG Agents*: Each input chunk is embedded and saved. At query time, the agent embeds the query and performs retrieval using cosine similarity between embeddings. (3) *Structure-Augmented RAG Agents*: After ingesting all input chunks, the agent constructs a structured representation (e.g., knowledge graph or event timeline). Subsequent queries are answered based on this structured memory.
:::
::::

:::: {#S3.SS3.SSS0.Px3 .section .ltx_paragraph}
##### *Agentic Memory Agents* {#agentic-memory-agents .ltx_title .ltx_title_paragraph}

::: {#S3.SS3.SSS0.Px3.p1 .ltx_para}
Agentic memory agents extend beyond static memory stores by employing agentic loops---iterative reasoning cycles in which the agent may reformulate questions, perform memory lookups, and update its working memory. These agents are designed to simulate a more human-like process of recalling, verifying, and integrating knowledge.
:::
::::
::::::::::

:::::::: {#S3.SS4 .section .ltx_subsection}
### [3.4 ]{.ltx_tag .ltx_tag_subsection}Datasets and Agents Formulation {#datasets-and-agents-formulation .ltx_title .ltx_title_subsection}

::::: {#S3.SS4.SSS0.Px1 .section .ltx_paragraph}
##### Datasets Formulation {#datasets-formulation .ltx_title .ltx_title_paragraph}

::: {#S3.SS4.SSS0.Px1.p1 .ltx_para}
We standardize all datasets into the format: ${c_{1},c_{2},\cdots,c_{n}}$ (chunks), ${q_{1},q_{2},\cdots,q_{m}}$ (questions), and ${a_{1},a_{2},\cdots,a_{m}}$ (answers), where $c_{i}$ denotes the $i$-th chunk wrapped to construct a user message with instructions of memorizing the content in a sequential input, and ${c_{1},c_{2},\cdots,c_{n}}$ represents a single conversation. Each chunk is accompanied by instructions prompting the agent to memorize its contents. Example prompts are provided in Appendix [[B.1]{.ltx_text .ltx_ref_tag}](#A2.SS1 "B.1 Instructions for Memory Construction ‣ Appendix B Prompts ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}.
:::

::: {#S3.SS4.SSS0.Px1.p2 .ltx_para}
When curating datasets like EventQA and FactConsolidation, we deliberately design scenarios where multiple questions follow a single context. This allows us to probe the model's memory multiple times with one sequential injection. For example, in LME (S\*), five contexts are paired with 300 questions (shown in Table [[1]{.ltx_text .ltx_ref_tag}](#S3.T1 "Table 1 ‣ Conflict Resolution (CR) ‣ 3.1 Aspects of the Evaluation ‣ 3 MemoryAgentBench ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}). This design choice reflects a key trend: as LLMs support increasingly long context windows and memory agents become more capable of handling extended inputs, evaluation datasets must also scale accordingly. Injecting 1M tokens for just one question is resource-inefficient, whereas associating the same input with many questions provides significantly higher utility.
:::
:::::

:::: {#S3.SS4.SSS0.Px2 .section .ltx_paragraph}
##### Agents Formulation {#agents-formulation .ltx_title .ltx_title_paragraph}

::: {#S3.SS4.SSS0.Px2.p1 .ltx_para}
In our framework, all agents are required to take the chunks one by one, absorb them into memory, and incrementally update the memory. After seeing all the chunks, we ask the agent to answer the related questions.
:::
::::
::::::::
:::::::::::::::::::::::::::::::::::::

:::::::::::::::::: {#S4 .section .ltx_section}
## [4 ]{.ltx_tag .ltx_tag_section}Experiments {#experiments .ltx_title .ltx_title_section}

::::: {#S4.SS1 .section .ltx_subsection}
### [4.1 ]{.ltx_tag .ltx_tag_subsection}Experimental Setup {#experimental-setup .ltx_title .ltx_title_subsection}

::: {#S4.SS1.p1 .ltx_para}
The datasets are split into four categories and the statistics of all datasets are also shown in Table [[1]{.ltx_text .ltx_ref_tag}](#S3.T1 "Table 1 ‣ Conflict Resolution (CR) ‣ 3.1 Aspects of the Evaluation ‣ 3 MemoryAgentBench ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}. The evaluation metrics for all datasets are shown in Table [[5]{.ltx_text .ltx_ref_tag}](#A1.T5 "Table 5 ‣ (5) EventQA ‣ A.1 Accurate Retrieval (AR) ‣ Appendix A Details of Dataset ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref} in Appendix [[A]{.ltx_text .ltx_ref_tag}](#A1 "Appendix A Details of Dataset ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}, along with more dataset details. Then for the agents, as described in Section [[3.3]{.ltx_text .ltx_ref_tag}](#S3.SS3 "3.3 Different Categories of Memory Agents ‣ 3 MemoryAgentBench ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}, we consider three categories of agents: *Long-Context Agents*, *RAG agents* and *Agentic Memory Agents*, where *RAG Agents* can be further split into *Simple RAG Agents*, *Embedding-based RAG Agents* and *Structure-Augmented RAG Agents*.
:::

::: {#S4.SS1.p2 .ltx_para}
For chunk size settings, we choose a chunk size of 512 for the RULER-QA, NIAH-MQ, and LME(S\*) tasks in AR, as well as for all tasks in CR. This is mainly because these tasks are composed of long texts synthesized from multiple short texts. For other tasks, we use a chunk size of 4096. Considering computational overhead, we uniformly use a chunk size of 4096 for the Mem0 and Cognee methods. We report the settings of the chunk size in Table  [[15]{.ltx_text .ltx_ref_tag}](#A4.T15 "Table 15 ‣ D.3 Settings of the Chunk Size ‣ Appendix D Experimental Settings ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref} in Appendix  [[C]{.ltx_text .ltx_ref_tag}](#A3 "Appendix C Detailed Experimental Results ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}.
:::

<figure id="S4.T2" class="ltx_table">
<div class="ltx_inline-block ltx_align_center ltx_transformed_outer" style="width:433.6pt;height:300.6pt;vertical-align:-300.6pt;">
<span class="ltx_transformed_inner" style="transform:translate(-125.7pt,0.0pt) scale(0.632943950254187,0.632943950254187) ;"> </span>
<table class="ltx_tabular ltx_align_middle">
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<td class="ltx_td ltx_border_r"></td>
<td colspan="5" class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold" style="background-color:#F0F0F0;">AR</span></td>
<td colspan="2" class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold" style="background-color:#F0F0F0;">TTL</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold" style="background-color:#F0F0F0;">LRU</span></td>
<td colspan="2" class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold" style="background-color:#F0F0F0;">CR</span></td>
</tr>
<tr class="ltx_tr" style="background-color:#F0F0F0;">
<td class="ltx_td ltx_align_left ltx_border_r"><span class="ltx_text ltx_font_bold" style="background-color:#F0F0F0;">Agent Type</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">RULER-QA</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">NIAH-MQ</span></td>
<td class="ltx_td ltx_align_center"><span class="math inline">∞</span><span class="ltx_text" style="background-color:#F0F0F0;">Bench-QA</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">LME(S*)</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F0F0F0;">EventQA</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">MCC</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F0F0F0;">Recom</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="math inline">∞</span><span class="ltx_text" style="background-color:#F0F0F0;">Bench-Sum</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">FactCon-SH</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">FactCon-MH</span></td>
</tr>
<tr class="ltx_tr" style="background-color:#DCEBFA;">
<td colspan="11" class="ltx_td ltx_align_center ltx_border_tt"><em>Long-Context Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">GPT-4o</td>
<td class="ltx_td ltx_align_center">61.5</td>
<td class="ltx_td ltx_align_center">25.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">55.4</span></td>
<td class="ltx_td ltx_align_center">32.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">77.2</td>
<td class="ltx_td ltx_align_center">87.6</td>
<td class="ltx_td ltx_align_center ltx_border_r">12.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">32.2</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">60.0</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">5.0</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">GPT-4o-mini</td>
<td class="ltx_td ltx_align_center">53.5</td>
<td class="ltx_td ltx_align_center">22.8</td>
<td class="ltx_td ltx_align_center">44.9</td>
<td class="ltx_td ltx_align_center">30.7</td>
<td class="ltx_td ltx_align_center ltx_border_r">59.0</td>
<td class="ltx_td ltx_align_center">82.4</td>
<td class="ltx_td ltx_align_center ltx_border_r">15.1</td>
<td class="ltx_td ltx_align_center ltx_border_r">28.9</td>
<td class="ltx_td ltx_align_center">45.0</td>
<td class="ltx_td ltx_align_center">5.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">GPT-4.1-mini</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">74.5</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">94.8</span></td>
<td class="ltx_td ltx_align_center">45.8</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">55.7</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text ltx_font_bold">82.6</span></td>
<td class="ltx_td ltx_align_center">75.6</td>
<td class="ltx_td ltx_align_center ltx_border_r">16.7</td>
<td class="ltx_td ltx_align_center ltx_border_r">41.9</td>
<td class="ltx_td ltx_align_center">36.0</td>
<td class="ltx_td ltx_align_center">5.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">Gemini-2.0-Flash</td>
<td class="ltx_td ltx_align_center">73.0</td>
<td class="ltx_td ltx_align_center">83.8</td>
<td class="ltx_td ltx_align_center">53.2</td>
<td class="ltx_td ltx_align_center">47.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">67.2</td>
<td class="ltx_td ltx_align_center">84.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">8.7</td>
<td class="ltx_td ltx_align_center ltx_border_r">23.9</td>
<td class="ltx_td ltx_align_center">30.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">Claude-3.7-Sonnet</td>
<td class="ltx_td ltx_align_center">65.0</td>
<td class="ltx_td ltx_align_center">38.0</td>
<td class="ltx_td ltx_align_center">50.6</td>
<td class="ltx_td ltx_align_center">34.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">74.6</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">89.4</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text ltx_font_bold">18.3</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text ltx_font_bold">52.5</span></td>
<td class="ltx_td ltx_align_center">43.0</td>
<td class="ltx_td ltx_align_center">2.0</td>
</tr>
<tr class="ltx_tr" style="background-color:#F0F0F0;">
<td class="ltx_td ltx_align_left ltx_border_r ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">GPT-4o-mini</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">53.5</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">22.8</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">44.9</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">30.7</span></td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">59.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">82.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">15.1</span></td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">28.9</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">45.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">5.0</span></td>
</tr>
<tr class="ltx_tr" style="background-color:#DCFADC;">
<td colspan="11" class="ltx_td ltx_align_center"><em>Simple RAG Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">BM25</td>
<td class="ltx_td ltx_align_center">61.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">100.0</span></td>
<td class="ltx_td ltx_align_center">45.6</td>
<td class="ltx_td ltx_align_center">45.3</td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text ltx_font_bold">74.6</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">75.4</span></td>
<td class="ltx_td ltx_align_center ltx_border_r">13.6</td>
<td class="ltx_td ltx_align_center ltx_border_r">20.9</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">56.0</span></td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr" style="background-color:#DCFADC;">
<td colspan="11" class="ltx_td ltx_align_center"><em>Embedding RAG Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">Contriever</td>
<td class="ltx_td ltx_align_center">26.5</td>
<td class="ltx_td ltx_align_center">2.5</td>
<td class="ltx_td ltx_align_center">38.1</td>
<td class="ltx_td ltx_align_center">15.7</td>
<td class="ltx_td ltx_align_center ltx_border_r">66.8</td>
<td class="ltx_td ltx_align_center">70.6</td>
<td class="ltx_td ltx_align_center ltx_border_r">15.2</td>
<td class="ltx_td ltx_align_center ltx_border_r">21.2</td>
<td class="ltx_td ltx_align_center">18.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">7.0</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">Text-Embed-3-Small</td>
<td class="ltx_td ltx_align_center">52.0</td>
<td class="ltx_td ltx_align_center">7.2</td>
<td class="ltx_td ltx_align_center">44.4</td>
<td class="ltx_td ltx_align_center">48.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">63.0</td>
<td class="ltx_td ltx_align_center">70.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">15.3</td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text ltx_font_bold">25.7</span></td>
<td class="ltx_td ltx_align_center">28.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">Text-Embed-3-Large</td>
<td class="ltx_td ltx_align_center">49.0</td>
<td class="ltx_td ltx_align_center">19.5</td>
<td class="ltx_td ltx_align_center">50.1</td>
<td class="ltx_td ltx_align_center">52.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">70.0</td>
<td class="ltx_td ltx_align_center">72.4</td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text ltx_font_bold">16.2</span></td>
<td class="ltx_td ltx_align_center ltx_border_r">21.6</td>
<td class="ltx_td ltx_align_center">28.0</td>
<td class="ltx_td ltx_align_center">4.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">NV-Embed-v2</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">83.0</span></td>
<td class="ltx_td ltx_align_center">73.5</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">51.4</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">55.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_r">72.8</td>
<td class="ltx_td ltx_align_center">69.4</td>
<td class="ltx_td ltx_align_center ltx_border_r">13.5</td>
<td class="ltx_td ltx_align_center ltx_border_r">20.7</td>
<td class="ltx_td ltx_align_center">55.0</td>
<td class="ltx_td ltx_align_center">6.0</td>
</tr>
<tr class="ltx_tr" style="background-color:#DCFADC;">
<td colspan="11" class="ltx_td ltx_align_center"><em>Structure-Augmented RAG Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">RAPTOR</td>
<td class="ltx_td ltx_align_center">33.5</td>
<td class="ltx_td ltx_align_center">15.8</td>
<td class="ltx_td ltx_align_center">31.3</td>
<td class="ltx_td ltx_align_center">34.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">45.8</td>
<td class="ltx_td ltx_align_center">59.4</td>
<td class="ltx_td ltx_align_center ltx_border_r">12.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">13.4</td>
<td class="ltx_td ltx_align_center">14.0</td>
<td class="ltx_td ltx_align_center">1.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">GraphRAG</td>
<td class="ltx_td ltx_align_center">47.0</td>
<td class="ltx_td ltx_align_center">38.3</td>
<td class="ltx_td ltx_align_center">35.8</td>
<td class="ltx_td ltx_align_center">35.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">34.4</td>
<td class="ltx_td ltx_align_center">39.8</td>
<td class="ltx_td ltx_align_center ltx_border_r">9.8</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4</td>
<td class="ltx_td ltx_align_center">14.0</td>
<td class="ltx_td ltx_align_center">2.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">HippoRAG-v2</td>
<td class="ltx_td ltx_align_center">71.0</td>
<td class="ltx_td ltx_align_center">67.5</td>
<td class="ltx_td ltx_align_center">45.7</td>
<td class="ltx_td ltx_align_center">50.7</td>
<td class="ltx_td ltx_align_center ltx_border_r">67.6</td>
<td class="ltx_td ltx_align_center">61.4</td>
<td class="ltx_td ltx_align_center ltx_border_r">10.2</td>
<td class="ltx_td ltx_align_center ltx_border_r">14.6</td>
<td class="ltx_td ltx_align_center">54.0</td>
<td class="ltx_td ltx_align_center">5.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">Mem0</td>
<td class="ltx_td ltx_align_center">28.0</td>
<td class="ltx_td ltx_align_center">4.8</td>
<td class="ltx_td ltx_align_center">22.4</td>
<td class="ltx_td ltx_align_center">36.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">37.5</td>
<td class="ltx_td ltx_align_center">3.4</td>
<td class="ltx_td ltx_align_center ltx_border_r">10.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.8</td>
<td class="ltx_td ltx_align_center">18.0</td>
<td class="ltx_td ltx_align_center">2.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">Cognee</td>
<td class="ltx_td ltx_align_center">33.5</td>
<td class="ltx_td ltx_align_center">4.0</td>
<td class="ltx_td ltx_align_center">19.7</td>
<td class="ltx_td ltx_align_center">29.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">26.8</td>
<td class="ltx_td ltx_align_center">35.4</td>
<td class="ltx_td ltx_align_center ltx_border_r">10.1</td>
<td class="ltx_td ltx_align_center ltx_border_r">2.3</td>
<td class="ltx_td ltx_align_center">28.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr" style="background-color:#FAE6F0;">
<td colspan="11" class="ltx_td ltx_align_center"><em>Agentic Memory Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_r">Self-RAG</td>
<td class="ltx_td ltx_align_center">38.5</td>
<td class="ltx_td ltx_align_center">8.0</td>
<td class="ltx_td ltx_align_center">28.5</td>
<td class="ltx_td ltx_align_center">25.7</td>
<td class="ltx_td ltx_align_center ltx_border_r">31.8</td>
<td class="ltx_td ltx_align_center">11.6</td>
<td class="ltx_td ltx_align_center ltx_border_r">12.8</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.9</td>
<td class="ltx_td ltx_align_center">19.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_bb ltx_border_r">MemGPT</td>
<td class="ltx_td ltx_align_center ltx_border_bb">39.5</td>
<td class="ltx_td ltx_align_center ltx_border_bb">8.8</td>
<td class="ltx_td ltx_align_center ltx_border_bb">20.8</td>
<td class="ltx_td ltx_align_center ltx_border_bb">32.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">26.2</td>
<td class="ltx_td ltx_align_center ltx_border_bb">67.6</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">14.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">2.5</td>
<td class="ltx_td ltx_align_center ltx_border_bb">28.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">3.0</td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 2</span>: </span><span class="ltx_text" style="font-size:90%;">Overall Performance Comparison. All RAG agents and commercial memory agents use GPT-4o-mini as the backbone. Thus we highlight the performance of GPT-4o-mini as the reference. FactCon-SH and FactCon-MH mean FactConsolidation Single Hop and FactConsolidation Multi Hop, respectively. We use the NV-Embed-v2 as dense retriever based on the open-source code of HippoRAG-v2. </span></figcaption>
</figure>
:::::

:::: {#S4.SS2 .section .ltx_subsection}
### [4.2 ]{.ltx_tag .ltx_tag_subsection}Overall Performance Comparison {#overall-performance-comparison .ltx_title .ltx_title_subsection}

::: {#S4.SS2.p1 .ltx_para}
Table [[2]{.ltx_text .ltx_ref_tag}](#S4.T2 "Table 2 ‣ 4.1 Experimental Setup ‣ 4 Experiments ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref} presents the overall performance across different benchmarks. We summarize the key findings as follows: (1) [Superiority of RAG methods in Accurate Retrieval Tasks.]{.ltx_text .ltx_font_bold} Most RAG Agents are better than the backbone model "GPT-4o-mini" in the tasks within the Accurate Retrieval Category. This matches our intuition where RAG agents typically excel at extracting a small snippet of text that is crucial for answering the question. (2) [Superiority of Long-Context Models in Test-Time Learning and Long-Range Understanding.]{.ltx_text .ltx_font_bold} Long-context models achieve the best performance on TTL and LRU. This highlights a fundamental limitation of RAG methods and commercial memory agents, which still follow an agentic RAG paradigm. These systems retrieve only partial information from the past context, lacking the ability to capture a holistic understanding of the input---let alone perform learning across it. (3) [Limitation of All Existing Methods on Conflict Resolution.]{.ltx_text .ltx_font_bold} Although being a well-discussed task in model-editing community (Mitchell et al., [2022](#bib.bib32){.ltx_ref}; Fang et al., [2024](#bib.bib11){.ltx_ref}), resolving conflict poses a significant challenge on memory agents. We observe that all methods fail on the multi-hop situation (with achieving at most 6% accuracy). Only long context agents can achieve fairly reasonable results on single-hop scenarios. In Section [[4.4.2]{.ltx_text .ltx_ref_tag}](#S4.SS4.SSS2 "4.4.2 Validation of Dataset FactConsolidation ‣ 4.4 Ablation Study on Input Chunk Size ‣ 4 Experiments ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}, we show that current reasoning models can have much better performance, while it does not change the conclusion that Conflict Resolution still poses a significant challenge on all memory mechanisms. (4) [Limited Performance of Commercial Memory Agents.]{.ltx_text .ltx_font_bold} Commercial memory agents such as MemGPT and Mem0 exhibit limited performance across a broad range of benchmarks. This shortfall can be attributed to three primary factors. First, these systems frequently fail to capture and preserve sufficient information when storing inputs into memory. For example, Mem0 depends on extracting factual knowledge from inputs, an approach that inherently discards a substantial portion of the original content. As a result, reconstructing inputs and supporting downstream tasks such as question answering becomes significantly more challenging. While Mem0 has demonstrated relatively strong performance on conversational tasks such as LOCOMO---where information density is comparatively low---it tends to perform poorly on benchmarks containing dense informational content, including RULER and $\infty$-Bench. For tasks emphasizing Time-to-Live (TTL) and Least Recently Used (LRU) retrieval, these limitations are often even more pronounced. Second, both MemGPT and Mem0 rely on retrieval mechanisms that only access a subset of stored information. In the case of Mem0, retrieval is typically performed a single time, similar to conventional RAG methods, constraining the breadth of information available for reasoning. MemGPT, although adopting a more agentic framework that permits multiple retrieval iterations, does not maintain temporal or structural metadata about the stored content. Consequently, the agent is unable to reconstruct longer documents in their original form, which adversely affects performance on LRU and other tasks requiring structured memory retrieval. Finally, methods such as MemGPT depend heavily on embedding-based retrieval mechanisms, which can be insufficient for fine-grained tasks like NIAH, where locating specific, precise information ("the needle in the haystack") is essential. These embedding-based approaches often struggle to distinguish subtle contextual nuances that are critical for accurate retrieval in such settings.
:::

<figure id="S4.F2" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S4.F2.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2507.05257/assets/x3.png" id="S4.F2.sf1.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="830" height="461" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(a)</span> </span><span class="ltx_text" style="font-size:90%;">RULER-QA performance</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S4.F2.sf2" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2507.05257/assets/x4.png" id="S4.F2.sf2.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="830" height="461" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">(b)</span> </span><span class="math inline">∞</span><span class="ltx_text" style="font-size:90%;">Bench-Sum performance</span></figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 2</span>: </span><span class="ltx_text" style="font-size:90%;">Performances on RULER-QA with different chunk sizes.</span></figcaption>
</figure>
::::

:::: {#S4.SS3 .section .ltx_subsection}
### [4.3 ]{.ltx_tag .ltx_tag_subsection}Ablation Study {#ablation-study .ltx_title .ltx_title_subsection}

::: {#S4.SS3.p1 .ltx_para}
In this Section, we conduct experiments and result analysis along four dimensions: Input Chunk Size, Retrieval TopK, Validation of Dataset, and Computational Latency. More detailed experimental results are provided in Appendix  [[C]{.ltx_text .ltx_ref_tag}](#A3 "Appendix C Detailed Experimental Results ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}.
:::
::::

:::::::::: {#S4.SS4 .section .ltx_subsection}
### [4.4 ]{.ltx_tag .ltx_tag_subsection}Ablation Study on Input Chunk Size {#ablation-study-on-input-chunk-size .ltx_title .ltx_title_subsection}

::: {#S4.SS4.p1 .ltx_para}
To understand how chunk size impacts performance, particularly for RAG methods and agentic memory agents, we conduct an additional analysis where we vary the chunk size while fixing the number of retrieved chunks to 10. The results are presented in Figure [[2]{.ltx_text .ltx_ref_tag}](#S4.F2 "Figure 2 ‣ 4.2 Overall Performance Comparison ‣ 4 Experiments ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}. From the figure, we observe the following: (1) In the RULER-QA task, reducing chunk size has little effect on BM25 performance. This is expected, as BM25 relies on term-frequency-based scoring and document-level ranking, and does not inherently benefit from finer-grained segmentation beyond the impact on term distributions. In contrast, embedding-based methods---including MemGPT, which uses [text-embedding-3-small]{.ltx_text .ltx_font_typewriter} as its retriever---consistently perform better with smaller chunks. This suggests that finer segmentation improves the granularity and relevance of retrieved results for models that rely on dense semantic representations. (2) In $\infty$Bench-Sum, however, smaller chunk sizes lead to worse performance. This task requires the agent to summarize an entire conversation, and smaller chunks correspond to fewer available tokens per retrieval. As a result, the agent has access to less context, which degrades summarization quality. The results suggest that, when resources permit, using smaller chunk sizes and increasing the number of retrieval calls during memory construction can improve performance on Accurate Retrieval (AR) tasks. Finer-grained segmentation enhances the relevance of retrieved information, particularly for embedding-based methods. However, for tasks requiring Long-Range Understanding (LRU), varying the chunk size hurts the performance. This is likely because RAG methods are inherently less suited for tasks that demand integration of information across a large, coherent context.
:::

<figure id="S4.F3" class="ltx_figure">
<img src="/html/2507.05257/assets/x5.png" id="S4.F3.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="830" height="230" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 3</span>: </span><span class="ltx_text" style="font-size:90%;">The accuracies on different benchmarks when varying the retrieval top-k to be 2, 5 and 10.</span></figcaption>
</figure>

:::: {#S4.SS4.SSS1 .section .ltx_subsubsection}
#### [4.4.1 ]{.ltx_tag .ltx_tag_subsubsection}Ablation Study on Retrieval TopK {#ablation-study-on-retrieval-topk .ltx_title .ltx_title_subsubsection}

::: {#S4.SS4.SSS1.p1 .ltx_para}
In our experiments, although we report results with the number of retrieved chunks set to 10 in Table [[2]{.ltx_text .ltx_ref_tag}](#S4.T2 "Table 2 ‣ 4.1 Experimental Setup ‣ 4 Experiments ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}, we also conducted ablation studies with varying retrieval sizes. A subset of these results is visualized in Figure [[3]{.ltx_text .ltx_ref_tag}](#S4.F3 "Figure 3 ‣ 4.4 Ablation Study on Input Chunk Size ‣ 4 Experiments ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}, with the full results provided in Table [[11]{.ltx_text .ltx_ref_tag}](#A3.T11 "Table 11 ‣ C.3 Detailed Results on Ablation Study ‣ Appendix C Detailed Experimental Results ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}. The results indicate that increasing the number of retrieved chunks generally improves performance across most tasks. It is worth noting that, with a chunk size of 4096 tokens, retrieving 10 chunks already yields an input of approximately 40k tokens. This places significant demands on model capacity. Due to this high token volume, we do not evaluate settings with 20 retrieved chunks.
:::

<figure id="S4.T3" class="ltx_table ltx_align_floatright">
<table class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<thead class="ltx_thead">
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_column ltx_th_row ltx_border_r ltx_border_tt"></th>
<th colspan="2" class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_r ltx_border_tt">FactCon-SH</th>
<th colspan="2" class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_tt">FactCon-MH</th>
</tr>
</thead>
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<td class="ltx_td ltx_th ltx_th_row ltx_border_r"></td>
<td class="ltx_td ltx_align_center">6K</td>
<td class="ltx_td ltx_align_center ltx_border_r">32K</td>
<td class="ltx_td ltx_align_center">6K</td>
<td class="ltx_td ltx_align_center">32K</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_th_row ltx_border_r ltx_border_t">GPT-4o</td>
<td class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_t">92.0</td>
<td class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_r ltx_border_t">88.0</td>
<td class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_t">28.0</td>
<td class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_t">10.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_bb ltx_border_r">O4-mini</td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text ltx_font_bold">100.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">61.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text ltx_font_bold">80.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb">14.0</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 3</span>: </span><span class="ltx_text" style="font-size:90%;">Performances of reasoning models on the dataset FactConsolidation.</span></figcaption>
</figure>
::::

:::: {#S4.SS4.SSS2 .section .ltx_subsubsection}
#### [4.4.2 ]{.ltx_tag .ltx_tag_subsubsection}Validation of Dataset FactConsolidation {#validation-of-dataset-factconsolidation .ltx_title .ltx_title_subsubsection}

::: {#S4.SS4.SSS2.p1 .ltx_para}
As the performance of different models on this dataset remains drastically low, we turn to the stronger reasoning model o4-mini and validate our dataset by checking the performance of o4-mini on a smaller version of this dataset. The results are shown in Table  [[3]{.ltx_text .ltx_ref_tag}](#S4.T3 "Table 3 ‣ 4.4.1 Ablation Study on Retrieval TopK ‣ 4.4 Ablation Study on Input Chunk Size ‣ 4 Experiments ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}.
:::

<figure id="S4.T4" class="ltx_table ltx_align_floatright">
<table class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_row ltx_border_r ltx_border_tt"></th>
<td colspan="2" class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="font-size:90%;">512</span></td>
<td colspan="2" class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="font-size:90%;">4096</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_row ltx_border_r"></th>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">M.C.</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">Q.E.</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">M.C.</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">Q.E.</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t"><span class="ltx_text" style="font-size:90%;">GPT-4o-mini</span></th>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">0.09</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">5.2</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">0.07</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">5.1</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t"><span class="ltx_text" style="font-size:90%;">BM25</span></th>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">0.11</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">0.79</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">0.10</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">1.8</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t"><span class="ltx_text" style="font-size:90%;">Contriever</span></th>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">7.2</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">0.76</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">1.7</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">2.0</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r"><span class="ltx_text" style="font-size:90%;">Text-Embed-3-Large</span></th>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">6.3</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">0.54</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">5.4</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">1.8</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r"><span class="ltx_text" style="font-size:90%;">NV-Embed-v2</span></th>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">93.4</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">0.83</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">42.9</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">1.8</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t"><span class="ltx_text" style="font-size:90%;">RAPTOR</span></th>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">151</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">0.51</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">133</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">0.60</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r"><span class="ltx_text" style="font-size:90%;">GraphRAG</span></th>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">123</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">9.9</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">90.4</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">9.4</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r"><span class="ltx_text" style="font-size:90%;">HippoRAG-v2</span></th>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">817</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">1.1</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">284</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">2.6</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r"><span class="ltx_text" style="font-size:90%;">Mem0</span></th>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">14644</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">1.2</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">2140</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">1.2</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r"><span class="ltx_text" style="font-size:90%;">Cognee</span></th>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">8309</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">33.2</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">962</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="font-size:90%;">4.5</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t"><span class="ltx_text" style="font-size:90%;">Self-RAG</span></th>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">8.4</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">2.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">6.7</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text" style="font-size:90%;">1.7</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_bb ltx_border_r"><span class="ltx_text" style="font-size:90%;">MemGPT</span></th>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text" style="font-size:90%;">413</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text" style="font-size:90%;">10.6</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text" style="font-size:90%;">93.3</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text" style="font-size:90%;">11.4</span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 4: </span>Computational Latency (in seconds). M.C.: Memory Construction. Q.E.: Query Execution.</figcaption>
</figure>
::::

:::: {#S4.SS4.SSS3 .section .ltx_subsubsection}
#### [4.4.3 ]{.ltx_tag .ltx_tag_subsubsection}Analysis of Computational Latency {#analysis-of-computational-latency .ltx_title .ltx_title_subsubsection}

::: {#S4.SS4.SSS3.p1 .ltx_para}
To illustrate the latency of various memory agents in terms of (1) Memory Construction (M.C.); (2) Query Execution (Q.E.), we randomly choose 5 examples from RULER-QA2 and LME (S\*), and report the latency of various memory agents. This part of experiments is done on a server with Four NVDIA L40 GPU and AMD EPYC 7713 64-Core CPU. We use the NV-Embed-v2 (7B) as the embedding model in HippoRAG-v2. We show the summarized results in Table [[4]{.ltx_text .ltx_ref_tag}](#S4.T4 "Table 4 ‣ 4.4.2 Validation of Dataset FactConsolidation ‣ 4.4 Ablation Study on Input Chunk Size ‣ 4 Experiments ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref} and put the full results in Table [[12]{.ltx_text .ltx_ref_tag}](#A3.T12 "Table 12 ‣ C.3 Detailed Results on Ablation Study ‣ Appendix C Detailed Experimental Results ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref} and [[14]{.ltx_text .ltx_ref_tag}](#A3.T14 "Table 14 ‣ C.3 Detailed Results on Ablation Study ‣ Appendix C Detailed Experimental Results ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}. From the table, we find that using a smaller chunk size requires significantly more time for memory construction, especially for methods such as HippoRAG-v2, Mem0, Cognee, and MemGPT. Meanwhile, methods such as Mem0, Cognee need extremely high resources when constructing the memory, which may pose challenges in real-world applications.
:::
::::
::::::::::
::::::::::::::::::

:::: {#S5 .section .ltx_section}
## [5 ]{.ltx_tag .ltx_tag_section}Conclusion and Future Work {#conclusion-and-future-work .ltx_title .ltx_title_section}

::: {#S5.p1 .ltx_para}
In this paper, we introduce [MemoryAgentBench]{.ltx_text .ltx_font_bold}, a unified benchmark designed to evaluate memory agents across four essential competencies: accurate retrieval, test-time learning, long-range understanding, and conflict resolution. While prior benchmarks focus largely on skill execution or long-context reasoning, MemoryAgentBench fills a critical gap by assessing how agents store, update, and utilize long-term information across multi-turn interactions. To build this benchmark, we restructure existing datasets and propose two new ones---[EventQA]{.ltx_text .ltx_font_bold} and [FactConsolidation]{.ltx_text .ltx_font_bold}---tailored to stress specific memory behaviors often overlooked in prior work. We evaluate a wide spectrum of agents, including long-context models, RAG-based systems, and commercial memory agents, under a consistent evaluation protocol. Our results reveal that, despite recent advances, current memory agents still exhibit substantial limitations when faced with tasks requiring dynamic memory updates and long-range consistency. One limitation of our work is that the datasets used in MemoryAgentBench are primarily synthetic, which may not fully reflect the characteristics of real-world user conversations. As future work, we aim to collect and curate more realistic, real-world datasets aligned with our four competencies to further enrich and diversify the benchmark and provide more comprehensive evaluations for memory agents.
:::
::::

::::: {#Sx1 .section .ltx_section}
## Acknowledgment {#acknowledgment .ltx_title .ltx_title_section}

::: {#Sx1.p1 .ltx_para}
We thank Kevin Lin for engaging in thoughtful discussions around the overall idea and latency evaluation. His input played a role in shaping the evaluation pipeline for the memory agents.
:::

::: {.ltx_pagination .ltx_role_newpage}
:::
:::::

::: {#bib .section .ltx_bibliography}
## References {#references .ltx_title .ltx_title_bibliography}

- [[Anthropic \[2025\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Anthropic. ]{.ltx_bibblock} [Claude 3.7 sonnet, 2025. ]{.ltx_bibblock} [URL [https://www.anthropic.com/news/claude-3-7-sonnet](https://www.anthropic.com/news/claude-3-7-sonnet){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock} [This announcement introduces Claude 3.7 Sonnet, described as Anthropic's most intelligent model to date and the first hybrid reasoning model generally available on the market. ]{.ltx_bibblock}]{#bib.bib1}
- [[Asai et al. \[2023\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Akari Asai, Zeqiu Wu, Yizhong Wang, Avirup Sil, and Hannaneh Hajishirzi. ]{.ltx_bibblock} [Self-rag: Learning to retrieve, generate, and critique through self-reflection. ]{.ltx_bibblock} [In *The Twelfth International Conference on Learning Representations*, 2023. ]{.ltx_bibblock}]{#bib.bib2}
- [[Bai et al. \[2023\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yushi Bai, Xin Lv, Jiajie Zhang, Hongchang Lyu, Jiankai Tang, Zhidian Huang, Zhengxiao Du, Xiao Liu, Aohan Zeng, Lei Hou, et al. ]{.ltx_bibblock} [Longbench: A bilingual, multitask benchmark for long context understanding. ]{.ltx_bibblock} [*arXiv preprint arXiv:2308.14508*, 2023. ]{.ltx_bibblock}]{#bib.bib3}
- [[Bai et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yushi Bai, Shangqing Tu, Jiajie Zhang, Hao Peng, Xiaozhi Wang, Xin Lv, Shulin Cao, Jiazheng Xu, Lei Hou, Yuxiao Dong, et al. ]{.ltx_bibblock} [Longbench v2: Towards deeper understanding and reasoning on realistic long-context multitasks. ]{.ltx_bibblock} [*arXiv preprint arXiv:2412.15204*, 2024. ]{.ltx_bibblock}]{#bib.bib4}
- [[Bertsch et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Amanda Bertsch, Maor Ivgi, Emily Xiao, Uri Alon, Jonathan Berant, Matthew R Gormley, and Graham Neubig. ]{.ltx_bibblock} [In-context learning with long-context models: An in-depth exploration. ]{.ltx_bibblock} [*arXiv preprint arXiv:2405.00200*, 2024. ]{.ltx_bibblock}]{#bib.bib5}
- [[Casanueva et al. \[2020\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Iñigo Casanueva, Tadas Temčinas, Daniela Gerz, Matthew Henderson, and Ivan Vulić. ]{.ltx_bibblock} [Efficient intent detection with dual sentence encoders. ]{.ltx_bibblock} [In Tsung-Hsien Wen, Asli Celikyilmaz, Zhou Yu, Alexandros Papangelis, Mihail Eric, Anuj Kumar, Iñigo Casanueva, and Rushin Shah, editors, *Proceedings of the 2nd Workshop on Natural Language Processing for Conversational AI*, pages 38--45, Online, July 2020. Association for Computational Linguistics. ]{.ltx_bibblock} [doi: [10.18653/v1/2020.nlp4convai-1.5]{.ltx_ref .ltx_nolink .ltx_Url .ltx_ref_self}. ]{.ltx_bibblock} [URL [https://aclanthology.org/2020.nlp4convai-1.5/](https://aclanthology.org/2020.nlp4convai-1.5/){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock}]{#bib.bib6}
- [[Chhikara et al. \[2025\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Prateek Chhikara, Dev Khant, Saket Aryan, Taranjeet Singh, and Deshraj Yadav. ]{.ltx_bibblock} [Mem0: Building production-ready ai agents with scalable long-term memory. ]{.ltx_bibblock} [*arXiv preprint arXiv:2504.19413*, 2025. ]{.ltx_bibblock}]{#bib.bib7}
- [[DeepMind \[2025\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ DeepMind. ]{.ltx_bibblock} [Gemini pro, 2025. ]{.ltx_bibblock} [URL [https://deepmind.google/technologies/gemini/pro/](https://deepmind.google/technologies/gemini/pro/){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock} [This page provides an overview of Gemini Pro, highlighting its advanced capabilities and applications in various fields. ]{.ltx_bibblock}]{#bib.bib8}
- [[Devlin et al. \[2019\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jacob Devlin, Ming-Wei Chang, Kenton Lee, and Kristina Toutanova. ]{.ltx_bibblock} [Bert: Pre-training of deep bidirectional transformers for language understanding. ]{.ltx_bibblock} [In *Proceedings of the 2019 conference of the North American chapter of the association for computational linguistics: human language technologies, volume 1 (long and short papers)*, pages 4171--4186, 2019. ]{.ltx_bibblock}]{#bib.bib9}
- [[Edge et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Darren Edge, Ha Trinh, Newman Cheng, Joshua Bradley, Alex Chao, Apurva Mody, Steven Truitt, Dasha Metropolitansky, Robert Osazuwa Ness, and Jonathan Larson. ]{.ltx_bibblock} [From local to global: A graph rag approach to query-focused summarization. ]{.ltx_bibblock} [*arXiv preprint arXiv:2404.16130*, 2024. ]{.ltx_bibblock}]{#bib.bib10}
- [[Fang et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Junfeng Fang, Houcheng Jiang, Kun Wang, Yunshan Ma, Shi Jie, Xiang Wang, Xiangnan He, and Tat-Seng Chua. ]{.ltx_bibblock} [Alphaedit: Null-space constrained knowledge editing for language models. ]{.ltx_bibblock} [*arXiv preprint arXiv:2410.02355*, 2024. ]{.ltx_bibblock}]{#bib.bib11}
- [[Gutiérrez et al. \[2025\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Bernal Jiménez Gutiérrez, Yiheng Shu, Weijian Qi, Sizhe Zhou, and Yu Su. ]{.ltx_bibblock} [From rag to memory: Non-parametric continual learning for large language models. ]{.ltx_bibblock} [*arXiv preprint arXiv:2502.14802*, 2025. ]{.ltx_bibblock}]{#bib.bib12}
- [[He et al. \[2023a\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zhankui He, Zhouhang Xie, Rahul Jha, Harald Steck, Dawen Liang, Yesu Feng, Bodhisattwa Prasad Majumder, Nathan Kallus, and Julian McAuley. ]{.ltx_bibblock} [Large language models as zero-shot conversational recommenders. ]{.ltx_bibblock} [In *Proceedings of the 32nd ACM international conference on information and knowledge management*, pages 720--730, 2023a. ]{.ltx_bibblock}]{#bib.bib13}
- [[He et al. \[2023b\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zhankui He, Zhouhang Xie, Rahul Jha, Harald Steck, Dawen Liang, Yesu Feng, Bodhisattwa Prasad Majumder, Nathan Kallus, and Julian McAuley. ]{.ltx_bibblock} [Large language models as zero-shot conversational recommenders. ]{.ltx_bibblock} [In *Proceedings of the 32nd ACM international conference on information and knowledge management*, pages 720--730, 2023b. ]{.ltx_bibblock}]{#bib.bib14}
- [[Hsieh et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Cheng-Ping Hsieh, Simeng Sun, Samuel Kriman, Shantanu Acharya, Dima Rekesh, Fei Jia, Yang Zhang, and Boris Ginsburg. ]{.ltx_bibblock} [RULER: What's the Real Context Size of Your Long-Context Language Models?, August 2024. ]{.ltx_bibblock} [URL [http://arxiv.org/abs/2404.06654](http://arxiv.org/abs/2404.06654){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock} [arXiv:2404.06654 \[cs\]. ]{.ltx_bibblock}]{#bib.bib15}
- [[Hu et al. \[2025\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Mengkang Hu, Yuhang Zhou, Wendong Fan, Yuzhou Nie, Bowei Xia, Tao Sun, Ziyu Ye, Zhaoxuan Jin, Yingru Li, Zeyu Zhang, Yifeng Wang, Qianshuo Ye, Ping Luo, and Guohao Li. ]{.ltx_bibblock} [Owl: Optimized workforce learning for general multi-agent assistance in real-world task automation, 2025. ]{.ltx_bibblock} [URL [https://github.com/camel-ai/owl](https://github.com/camel-ai/owl){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock}]{#bib.bib16}
- [[Izacard et al. \[2021\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Gautier Izacard, Mathilde Caron, Lucas Hosseini, Sebastian Riedel, Piotr Bojanowski, Armand Joulin, and Edouard Grave. ]{.ltx_bibblock} [Unsupervised dense information retrieval with contrastive learning. ]{.ltx_bibblock} [*arXiv preprint arXiv:2112.09118*, 2021. ]{.ltx_bibblock}]{#bib.bib17}
- [[Jimenez et al. \[2023\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Carlos E Jimenez, John Yang, Alexander Wettig, Shunyu Yao, Kexin Pei, Ofir Press, and Karthik Narasimhan. ]{.ltx_bibblock} [Swe-bench: Can language models resolve real-world github issues? ]{.ltx_bibblock} [*arXiv preprint arXiv:2310.06770*, 2023. ]{.ltx_bibblock}]{#bib.bib18}
- [[Karpinska et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Marzena Karpinska, Katherine Thai, Kyle Lo, Tanya Goyal, and Mohit Iyyer. ]{.ltx_bibblock} [One thousand and one pairs: A\" novel\" challenge for long-context language models. ]{.ltx_bibblock} [*arXiv preprint arXiv:2406.16264*, 2024. ]{.ltx_bibblock}]{#bib.bib19}
- [[Karpukhin et al. \[2020\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Vladimir Karpukhin, Barlas Oguz, Sewon Min, Patrick SH Lewis, Ledell Wu, Sergey Edunov, Danqi Chen, and Wen-tau Yih. ]{.ltx_bibblock} [Dense passage retrieval for open-domain question answering. ]{.ltx_bibblock} [In *EMNLP (1)*, pages 6769--6781, 2020. ]{.ltx_bibblock}]{#bib.bib20}
- [[Larson et al. \[2019\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Stefan Larson, Anish Mahendran, Joseph J. Peper, Christopher Clarke, Andrew Lee, Parker Hill, Jonathan K. Kummerfeld, Kevin Leach, Michael A. Laurenzano, Lingjia Tang, and Jason Mars. ]{.ltx_bibblock} [An evaluation dataset for intent classification and out-of-scope prediction. ]{.ltx_bibblock} [In Kentaro Inui, Jing Jiang, Vincent Ng, and Xiaojun Wan, editors, *Proceedings of the 2019 Conference on Empirical Methods in Natural Language Processing and the 9th International Joint Conference on Natural Language Processing (EMNLP-IJCNLP)*, pages 1311--1316, Hong Kong, China, November 2019. Association for Computational Linguistics. ]{.ltx_bibblock} [doi: [10.18653/v1/D19-1131]{.ltx_ref .ltx_nolink .ltx_Url .ltx_ref_self}. ]{.ltx_bibblock} [URL [https://aclanthology.org/D19-1131/](https://aclanthology.org/D19-1131/){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock}]{#bib.bib21}
- [[Lee et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Chankyu Lee, Rajarshi Roy, Mengyao Xu, Jonathan Raiman, Mohammad Shoeybi, Bryan Catanzaro, and Wei Ping. ]{.ltx_bibblock} [Nv-embed: Improved techniques for training llms as generalist embedding models. ]{.ltx_bibblock} [*arXiv preprint arXiv:2405.17428*, 2024. ]{.ltx_bibblock}]{#bib.bib22}
- [[Li et al. \[2023\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jiaqi Li, Mengmeng Wang, Zilong Zheng, and Muhan Zhang. ]{.ltx_bibblock} [Loogle: Can long-context language models understand long contexts? ]{.ltx_bibblock} [*arXiv preprint arXiv:2311.04939*, 2023. ]{.ltx_bibblock}]{#bib.bib23}
- [[Li et al. \[2018\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Raymond Li, Samira Ebrahimi Kahou, Hannes Schulz, Vincent Michalski, Laurent Charlin, and Chris Pal. ]{.ltx_bibblock} [Towards deep conversational recommendations. ]{.ltx_bibblock} [*Advances in neural information processing systems*, 31, 2018. ]{.ltx_bibblock}]{#bib.bib24}
- [[Li et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xianming Li, Julius Lipp, Aamir Shakir, Rui Huang, and Jing Li. ]{.ltx_bibblock} [Bmx: Entropy-weighted similarity and semantic-enhanced lexical search. ]{.ltx_bibblock} [*arXiv preprint arXiv:2408.06643*, 2024. ]{.ltx_bibblock}]{#bib.bib25}
- [[Li and Roth \[2002\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xin Li and Dan Roth. ]{.ltx_bibblock} [Learning question classifiers. ]{.ltx_bibblock} [In *COLING 2002: The 19th International Conference on Computational Linguistics*, 2002. ]{.ltx_bibblock} [URL [https://aclanthology.org/C02-1150/](https://aclanthology.org/C02-1150/){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock}]{#bib.bib26}
- [[Lin et al. \[2025\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Kevin Lin, Charlie Snell, Yu Wang, Charles Packer, Sarah Wooders, Ion Stoica, and Joseph E Gonzalez. ]{.ltx_bibblock} [Sleep-time compute: Beyond inference scaling at test-time. ]{.ltx_bibblock} [*arXiv preprint arXiv:2504.13171*, 2025. ]{.ltx_bibblock}]{#bib.bib27}
- [[Liu et al. \[2019\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xingkun Liu, Arash Eshghi, Pawel Swietojanski, and Verena Rieser. ]{.ltx_bibblock} [Benchmarking natural language understanding services for building conversational agents, 2019. ]{.ltx_bibblock} [URL [https://arxiv.org/abs/1903.05566](https://arxiv.org/abs/1903.05566){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock}]{#bib.bib28}
- [[Maharana et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Adyasha Maharana, Dong-Ho Lee, Sergey Tulyakov, Mohit Bansal, Francesco Barbieri, and Yuwei Fang. ]{.ltx_bibblock} [Evaluating very long-term conversational memory of llm agents. ]{.ltx_bibblock} [*arXiv preprint arXiv:2402.17753*, 2024. ]{.ltx_bibblock}]{#bib.bib29}
- [[Meng et al. \[2023\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Kevin Meng, Arnab Sen Sharma, Alex J. Andonian, Yonatan Belinkov, and David Bau. ]{.ltx_bibblock} [Mass-editing memory in a transformer. ]{.ltx_bibblock} [In *ICLR*. OpenReview.net, 2023. ]{.ltx_bibblock}]{#bib.bib30}
- [[Mialon et al. \[2023\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Grégoire Mialon, Clémentine Fourrier, Thomas Wolf, Yann LeCun, and Thomas Scialom. ]{.ltx_bibblock} [Gaia: a benchmark for general ai assistants. ]{.ltx_bibblock} [In *The Twelfth International Conference on Learning Representations*, 2023. ]{.ltx_bibblock}]{#bib.bib31}
- [[Mitchell et al. \[2022\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Eric Mitchell, Charles Lin, Antoine Bosselut, Christopher D. Manning, and Chelsea Finn. ]{.ltx_bibblock} [Memory-based model editing at scale. ]{.ltx_bibblock} [In *ICML*, volume 162 of *Proceedings of Machine Learning Research*, pages 15817--15831. PMLR, 2022. ]{.ltx_bibblock}]{#bib.bib32}
- [[Modarressi et al. \[2025\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ali Modarressi, Hanieh Deilamsalehy, Franck Dernoncourt, Trung Bui, Ryan A Rossi, Seunghyun Yoon, and Hinrich Schütze. ]{.ltx_bibblock} [Nolima: Long-context evaluation beyond literal matching. ]{.ltx_bibblock} [*arXiv preprint arXiv:2502.05167*, 2025. ]{.ltx_bibblock}]{#bib.bib33}
- [[Müller and Žunič \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Magnus Müller and Gregor Žunič. ]{.ltx_bibblock} [Browser use: Enable ai to control your browser, 2024. ]{.ltx_bibblock} [URL [https://github.com/browser-use/browser-use](https://github.com/browser-use/browser-use){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock}]{#bib.bib34}
- [[OpenAI \[2025\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ OpenAI. ]{.ltx_bibblock} [Gpt-4o system card, 2025. ]{.ltx_bibblock} [URL [https://openai.com/index/gpt-4o-system-card/](https://openai.com/index/gpt-4o-system-card/){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock} [This report outlines the safety work carried out prior to releasing GPT-4o including external red teaming, frontier risk evaluations according to our Preparedness Framework, and an overview of the mitigations we built in to address key risk areas. ]{.ltx_bibblock}]{#bib.bib35}
- [[Packer et al. \[2023\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Charles Packer, Vivian Fang, Shishir_G Patil, Kevin Lin, Sarah Wooders, and Joseph_E Gonzalez. ]{.ltx_bibblock} [Memgpt: Towards llms as operating systems. ]{.ltx_bibblock} [2023. ]{.ltx_bibblock}]{#bib.bib36}
- [[Rasmussen et al. \[2025\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Preston Rasmussen, Pavlo Paliychuk, Travis Beauvais, Jack Ryan, and Daniel Chalef. ]{.ltx_bibblock} [Zep: A temporal knowledge graph architecture for agent memory. ]{.ltx_bibblock} [*arXiv preprint arXiv:2501.13956*, 2025. ]{.ltx_bibblock}]{#bib.bib37}
- [[Robertson and Walker \[1994\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Stephen E Robertson and Steve Walker. ]{.ltx_bibblock} [Some simple effective approximations to the 2-poisson model for probabilistic weighted retrieval. ]{.ltx_bibblock} [In *SIGIR'94: Proceedings of the Seventeenth Annual International ACM-SIGIR Conference on Research and Development in Information Retrieval, organised by Dublin City University*, pages 232--241. Springer, 1994. ]{.ltx_bibblock}]{#bib.bib38}
- [[Sarthi et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Parth Sarthi, Salman Abdullah, Aditi Tuli, Shubh Khanna, Anna Goldie, and Christopher D Manning. ]{.ltx_bibblock} [Raptor: Recursive abstractive processing for tree-organized retrieval. ]{.ltx_bibblock} [In *The Twelfth International Conference on Learning Representations*, 2024. ]{.ltx_bibblock}]{#bib.bib39}
- [[Wang et al. \[2024a\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Cunxiang Wang, Ruoxi Ning, Boqi Pan, Tonghui Wu, Qipeng Guo, Cheng Deng, Guangsheng Bao, Qian Wang, and Yue Zhang. ]{.ltx_bibblock} [Novelqa: A benchmark for long-range novel question answering. ]{.ltx_bibblock} [*arXiv preprint arXiv:2403.12766*, 2024a. ]{.ltx_bibblock}]{#bib.bib40}
- [[Wang et al. \[2024b\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Minzheng Wang, Longze Chen, Fu Cheng, Shengyi Liao, Xinghua Zhang, Bingli Wu, Haiyang Yu, Nan Xu, Lei Zhang, Run Luo, et al. ]{.ltx_bibblock} [Leave no document behind: Benchmarking long-context llms with extended multi-doc qa. ]{.ltx_bibblock} [In *Proceedings of the 2024 Conference on Empirical Methods in Natural Language Processing*, pages 5627--5646, 2024b. ]{.ltx_bibblock}]{#bib.bib41}
- [[Wang et al. \[2024c\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xingyao Wang, Boxuan Li, Yufan Song, Frank F Xu, Xiangru Tang, Mingchen Zhuge, Jiayi Pan, Yueqi Song, Bowen Li, Jaskirat Singh, et al. ]{.ltx_bibblock} [Openhands: An open platform for ai software developers as generalist agents. ]{.ltx_bibblock} [In *The Thirteenth International Conference on Learning Representations*, 2024c. ]{.ltx_bibblock}]{#bib.bib42}
- [[\[43\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yu Wang, Xinshuang Liu, Xiusi Chen, Sean O'Brien, Junda Wu, and Julian McAuley. ]{.ltx_bibblock} [Self-updatable large language models by integrating context into model parameters. ]{.ltx_bibblock} [In *The Thirteenth International Conference on Learning Representations*. ]{.ltx_bibblock}]{#bib.bib43}
- [[Wang et al. \[2024d\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yu Wang, Yifan Gao, Xiusi Chen, Haoming Jiang, Shiyang Li, Jingfeng Yang, Qingyu Yin, Zheng Li, Xian Li, Bing Yin, et al. ]{.ltx_bibblock} [Memoryllm: Towards self-updatable large language models. ]{.ltx_bibblock} [*arXiv preprint arXiv:2402.04624*, 2024d. ]{.ltx_bibblock}]{#bib.bib44}
- [[Wang et al. \[2024e\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yu Wang, Ruihan Wu, Zexue He, Xiusi Chen, and Julian McAuley. ]{.ltx_bibblock} [Large scale knowledge washing. ]{.ltx_bibblock} [*arXiv preprint arXiv:2405.16720*, 2024e. ]{.ltx_bibblock}]{#bib.bib45}
- [[Wang et al. \[2025\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yu Wang, Dmitry Krotov, Yuanzhe Hu, Yifan Gao, Wangchunshu Zhou, Julian McAuley, Dan Gutfreund, Rogerio Feris, and Zexue He. ]{.ltx_bibblock} [M+: Extending memoryLLM with scalable long-term memory. ]{.ltx_bibblock} [In *Forty-second International Conference on Machine Learning*, 2025. ]{.ltx_bibblock} [URL [https://openreview.net/forum?id=OcqbkROe8J](https://openreview.net/forum?id=OcqbkROe8J){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}. ]{.ltx_bibblock}]{#bib.bib46}
- [[Wang et al. \[2025/02\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yu Wang, Chi Han, Tongtong Wu, Xiaoxin He, Wangchunshu Zhou, Nafis Sadeq, Xiusi Chen, Zexue He, Wei Wang, Gholamreza Haffari, Heng Ji, and Julian J. McAuley. ]{.ltx_bibblock} [Towards lifespan cognitive systems. ]{.ltx_bibblock} [*TMLR*, 2025/02. ]{.ltx_bibblock}]{#bib.bib47}
- [[Wu et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Di Wu, Hongwei Wang, Wenhao Yu, Yuwei Zhang, Kai-Wei Chang, and Dong Yu. ]{.ltx_bibblock} [Longmemeval: Benchmarking chat assistants on long-term interactive memory. ]{.ltx_bibblock} [*arXiv preprint arXiv:2410.10813*, 2024. ]{.ltx_bibblock}]{#bib.bib48}
- [[Wu et al. \[2022\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Qiyu Wu, Chongyang Tao, Tao Shen, Can Xu, Xiubo Geng, and Daxin Jiang. ]{.ltx_bibblock} [Pcl: Peer-contrastive learning with diverse augmentations for unsupervised sentence embeddings. ]{.ltx_bibblock} [*arXiv preprint arXiv:2201.12093*, 2022. ]{.ltx_bibblock}]{#bib.bib49}
- [[Yen et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Howard Yen, Tianyu Gao, Minmin Hou, Ke Ding, Daniel Fleischer, Peter Izsak, Moshe Wasserblat, and Danqi Chen. ]{.ltx_bibblock} [Helmet: How to evaluate long-context language models effectively and thoroughly. ]{.ltx_bibblock} [*arXiv preprint arXiv:2410.02694*, 2024. ]{.ltx_bibblock}]{#bib.bib50}
- [[Yin et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zhangyue Yin, Qiushi Sun, Qipeng Guo, Zhiyuan Zeng, Qinyuan Cheng, Xipeng Qiu, and Xuan-Jing Huang. ]{.ltx_bibblock} [Explicit memory learning with expectation maximization. ]{.ltx_bibblock} [In *Proceedings of the 2024 Conference on Empirical Methods in Natural Language Processing*, pages 16618--16635, 2024. ]{.ltx_bibblock}]{#bib.bib51}
- [[Yu et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tian Yu, Shaolei Zhang, and Yang Feng. ]{.ltx_bibblock} [Auto-rag: Autonomous retrieval-augmented generation for large language models. ]{.ltx_bibblock} [*arXiv preprint arXiv:2411.19443*, 2024. ]{.ltx_bibblock}]{#bib.bib52}
- [[Zhang et al. \[2024\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xinrong Zhang, Yingfa Chen, Shengding Hu, Zihang Xu, Junhao Chen, Moo Hao, Xu Han, Zhen Thai, Shuo Wang, Zhiyuan Liu, et al. ]{.ltx_bibblock} [$\infty$bench: Extending long context evaluation beyond 100k tokens. ]{.ltx_bibblock} [In *Proceedings of the 62nd Annual Meeting of the Association for Computational Linguistics (Volume 1: Long Papers)*, pages 15262--15277, 2024. ]{.ltx_bibblock}]{#bib.bib53}
- [[Zhong et al. \[2023\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zexuan Zhong, Zhengxuan Wu, Christopher D Manning, Christopher Potts, and Danqi Chen. ]{.ltx_bibblock} [Mquake: Assessing knowledge editing in language models via multi-hop questions. ]{.ltx_bibblock} [*arXiv preprint arXiv:2305.14795*, 2023. ]{.ltx_bibblock}]{#bib.bib54}
:::

::: {.ltx_pagination .ltx_role_newpage}
:::

:::::::::::::::::::::: {#A1 .section .ltx_appendix}
## [Appendix A ]{.ltx_tag .ltx_tag_appendix}Details of Dataset {#appendix-a-details-of-dataset .ltx_title .ltx_title_appendix}

::: {#A1.p1 .ltx_para}
Here we provide a detailed introduction to the datasets used for evaluating the four core competencies, including the dataset curation, corresponding metrics, average context length, and a brief description. Details are shown in Table [[5]{.ltx_text .ltx_ref_tag}](#A1.T5 "Table 5 ‣ (5) EventQA ‣ A.1 Accurate Retrieval (AR) ‣ Appendix A Details of Dataset ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}.
:::

:::::::::::::: {#A1.SS1 .section .ltx_subsection}
### [A.1 ]{.ltx_tag .ltx_tag_subsection}Accurate Retrieval (AR) {#a.1-accurate-retrieval-ar .ltx_title .ltx_title_subsection}

::: {#A1.SS1.p1 .ltx_para}
We use five datasets to evaluate the accurate retrieval capability of memory agents.
:::

:::: {#A1.SS1.SSS0.Px1 .section .ltx_paragraph}
##### (1) RULER-QA {#ruler-qa .ltx_title .ltx_title_paragraph}

::: {#A1.SS1.SSS0.Px1.p1 .ltx_para}
We adopt two QA datasets from \[Hsieh et al., [2024](#bib.bib15){.ltx_ref}\]. These datasets provide multiple synthetic contexts of varying lengths, ranging from 3K to over 200K tokens. We select 100 questions from the datasets with shorter contexts. For each of these 100 questions, we collect all the documents of them, removed duplicate content, and then shuffled and concatenated them to create new long contexts of 197K or 421K tokens, making sure the new contexts containing the gold passage. Since most answers are short informational entities, such as years, names, or yes/no responses, we use substring exact match (SubEM) as the evaluation metric. SubEM measures whether the predicted answer exactly matches the gold answer as a substring, which is a common standard in question answering systems.
:::
::::

:::: {#A1.SS1.SSS0.Px2 .section .ltx_paragraph}
##### (2) NIAH-MQ {#niah-mq .ltx_title .ltx_title_paragraph}

::: {#A1.SS1.SSS0.Px2.p1 .ltx_para}
We choose a context with a length of 448K tokens, containing 100 groups and a total of 400 queries. We first check whether these queries appeared in the context. Then, we randomly shuffl the queries and their corresponding numbers evenly to prevent them from clustering together in the long context. The primary evaluation criterion is whether the agent can successfully retrieve the correct numbers. Therefore, we use average recall as the evaluation metric.
:::
::::

:::: {#A1.SS1.SSS0.Px3 .section .ltx_paragraph}
##### (3) $\infty$Bench-En.QA {#inftybench-en.qa .ltx_title .ltx_title_paragraph}

::: {#A1.SS1.SSS0.Px3.p1 .ltx_para}
We borrow this dataset from Zhang et al. \[[2024](#bib.bib53){.ltx_ref}\]. It is a QA task using novels in which character names have been replaced. This makes the content more coherent and closer to real-world scenarios. We use ROUGE F1 for this dataset since answers are mostly entity names.
:::
::::

:::: {#A1.SS1.SSS0.Px4 .section .ltx_paragraph}
##### (4) LongMemEval {#longmemeval .ltx_title .ltx_title_paragraph}

::: {#A1.SS1.SSS0.Px4.p1 .ltx_para}
This is a dialogue-based QA dataset. For LME(S\*), we use multiple historical conversation data segments, arrange them in chronological order, and finally concatenate them to create five long conversation histories, each with a length of approximately 355K tokens. Since some of the questions have open-ended answers, we adopt the approach used in previous work and employ the GPT-4o model to assess whether the agent's responses meet the requirements. If a response is deemed satisfactory, it is marked as True. Finally, we calculate the proportion of satisfactory responses as the evaluation metric.
:::
::::

:::: {#A1.SS1.SSS0.Px5 .section .ltx_paragraph}
##### (5) EventQA {#eventqa .ltx_title .ltx_title_paragraph}

::: {#A1.SS1.SSS0.Px5.p1 .ltx_para}
Using five books from $\infty$Bench (each \>390K tokens, counted using the [gpt-4o-mini]{.ltx_text .ltx_font_typewriter} tokenizer), we identify the ten most frequently mentioned characters via [SpaCy]{.ltx_text .ltx_font_typewriter} NER. We extract 101 events experienced by key characters using [gpt-4o]{.ltx_text .ltx_font_typewriter}. For each event, we construct a 6-way multiple-choice question by pairing the true event with five distractors generated via [gpt-4o]{.ltx_text .ltx_font_typewriter}. The agent receives up-to five previous events and must identify the correct continuation. We report the mean accuracy over 100 such questions per book, and ultimately present the average accuracy across all five books.
:::

<figure id="A1.T5" class="ltx_table">
<div class="ltx_inline-block ltx_align_center ltx_transformed_outer" style="width:424.9pt;height:203.8pt;vertical-align:-203.8pt;">
<span class="ltx_transformed_inner" style="transform:translate(-151.4pt,0.0pt) scale(0.58397878545639,0.58397878545639) ;"> </span>
<table class="ltx_tabular ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_tt"><span class="ltx_text ltx_font_bold">Category</span></th>
<th class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_tt"><span class="ltx_text ltx_font_bold">Dataset</span></th>
<td class="ltx_td ltx_align_left ltx_border_tt"><span class="ltx_text ltx_font_bold">Metrics</span></td>
<td class="ltx_td ltx_align_left ltx_border_tt"><span class="ltx_text ltx_font_bold">Avg. Length</span></td>
<td class="ltx_td ltx_align_left ltx_border_tt"><span class="ltx_text ltx_font_bold">Description</span></td>
</tr>
<tr class="ltx_tr">
<th rowspan="7" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_t"><span class="ltx_text ltx_font_bold">Accurate Retrieval</span></th>
<th class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_t">RULER-QA1</th>
<td rowspan="2" class="ltx_td ltx_align_left ltx_border_t">SubEM</td>
<td class="ltx_td ltx_align_left ltx_border_t">197K</td>
<td rowspan="2" class="ltx_td ltx_align_left ltx_border_t">Gold passage retrieval QA.</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row">RULER-QA2</th>
<td class="ltx_td ltx_align_left">421K</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row">RULER-NIAH-MQ</th>
<td class="ltx_td ltx_align_left">Recall</td>
<td class="ltx_td ltx_align_left">448K</td>
<td class="ltx_td ltx_align_left">Retrieve multiple “needles” from the “haystack”.</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row"><span class="math inline">∞</span>Bench-QA</th>
<td class="ltx_td ltx_align_left">ROUGE F1</td>
<td class="ltx_td ltx_align_left">183K</td>
<td class="ltx_td ltx_align_left">Novel QA with entity replacement.</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row">LongMemEval (S)</th>
<td rowspan="2" class="ltx_td ltx_align_left">Model Based Acc.</td>
<td class="ltx_td ltx_align_left">110K</td>
<td rowspan="2" class="ltx_td ltx_align_left">Dialogues based QA.</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row">LongMemEval (S*)</th>
<td class="ltx_td ltx_align_left">355K</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row">EventQA (<em>ours</em>)</th>
<td class="ltx_td ltx_align_left">Accuracy</td>
<td class="ltx_td ltx_align_left">534K</td>
<td class="ltx_td ltx_align_left">Novel multiple-choice QA on characters events.</td>
</tr>
<tr class="ltx_tr">
<th rowspan="6" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_t"><span class="ltx_text ltx_font_bold">Test-time Learning</span></th>
<th class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_t">BANKING77</th>
<td rowspan="5" class="ltx_td ltx_align_left ltx_border_t">Accuracy</td>
<td rowspan="5" class="ltx_td ltx_align_left ltx_border_t">103K</td>
<td class="ltx_td ltx_align_left ltx_border_t">Banking intent classification, 77 labels</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row">CLINC150</th>
<td class="ltx_td ltx_align_left">Intent classification, 151 labels</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row">NLU</th>
<td class="ltx_td ltx_align_left">Task intent classification, 68 labels</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row">TREC Coarse</th>
<td class="ltx_td ltx_align_left">Question type classification, 6 labels</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row">TREC Fine</th>
<td class="ltx_td ltx_align_left">Question type classification, 50 labels</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row">Movie Recommendation</th>
<td class="ltx_td ltx_align_left">Recall@5</td>
<td class="ltx_td ltx_align_left">1.44M</td>
<td class="ltx_td ltx_align_left">Recommend movies based on provided dialogues examples.</td>
</tr>
<tr class="ltx_tr">
<th rowspan="2" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_t"><span class="ltx_text ltx_font_bold">Long Range Understanding</span></th>
<th class="ltx_td ltx_th ltx_th_row ltx_border_t"></th>
<td class="ltx_td ltx_border_t"></td>
<td class="ltx_td ltx_border_t"></td>
<td class="ltx_td ltx_border_t"></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row"><span class="math inline">∞</span>Bench-Sum</th>
<td class="ltx_td ltx_align_left">Model Based F1</td>
<td class="ltx_td ltx_align_left">172K</td>
<td class="ltx_td ltx_align_left">Novel summarization with entity replacement.</td>
</tr>
<tr class="ltx_tr">
<th rowspan="2" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_bb ltx_border_t"><span class="ltx_text ltx_font_bold">Conflict Resolving</span></th>
<th class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_t">FactConsolidation-SH (<em>ours</em>)</th>
<td rowspan="2" class="ltx_td ltx_align_left ltx_border_bb ltx_border_t">SubEM</td>
<td rowspan="2" class="ltx_td ltx_align_left ltx_border_bb ltx_border_t">262K</td>
<td class="ltx_td ltx_align_left ltx_border_t">Conflict solving in single hop reasoning.</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_bb">FactConsolidation-MH (<em>ours</em>)</th>
<td class="ltx_td ltx_align_left ltx_border_bb">Conflict solving in multiple hop reasoning.</td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 5</span>: </span><span class="ltx_text" style="font-size:90%;"> Overview of evaluation datasets. We select datasets that cover various important long-context capabilities. SubEM: substring exact match. In the table, we underline the datasets we constructed ourselves. Avg. Length: Average Context Length (measured using the GPT-4o-mini model’s tokenizer). </span></figcaption>
</figure>
::::
::::::::::::::

:::: {#A1.SS2 .section .ltx_subsection}
### [A.2 ]{.ltx_tag .ltx_tag_subsection}Test-time Learning (TTL) {#a.2-test-time-learning-ttl .ltx_title .ltx_title_subsection}

::: {#A1.SS2.p1 .ltx_para}
We evaluate TTL via two task categories: (1) [Multi-Class Classification (MCC)]{.ltx_text .ltx_font_bold}: We adopt five classification datasets used in prior TTL work. For dataset curation, we use thousands of sentence samples from different categories, with each type of sample assigned a number as its label. Following the format \"{sentence} \\n Label: {label} \\n\", we concatenate all the sentences into a long context and shuffle them to prevent samples of the same type from being too concentrated. In this task, the agent needs to refer to a long context and correctly classify the input content. Therefore, we use average accuracy as the evaluation metric. (2) [Recommendation (Recom)]{.ltx_text .ltx_font_bold}: We concatenate multiple short dialogues about movie recommendations from the original dataset, remove duplicate dialogues, and create a long context containing over a thousand recommendation instances. In this task, the agent is required to recommend 20 movies based on the content of the dialogue. We evaluate the recommendations by calculating Recall@5, which measures the overlap between the top 5 recommended movies and the ground truth.
:::
::::

:::: {#A1.SS3 .section .ltx_subsection}
### [A.3 ]{.ltx_tag .ltx_tag_subsection}Long-Range Understanding (LRU) {#a.3-long-range-understanding-lru .ltx_title .ltx_title_subsection}

::: {#A1.SS3.p1 .ltx_para}
We evaluate LRU via the Summarization task [En.Sum]{.ltx_text .ltx_font_typewriter} from $\infty$-Bench \[Zhang et al., [2024](#bib.bib53){.ltx_ref}\]. We follow the settings from  \[Yen et al., [2024](#bib.bib50){.ltx_ref}\] and use the GPT-4o model in evaluating the summarized text. In this process, we assess the fluency of the input text (scored as 0 or 1) and use the dot product of this score with the F1 score as the final evaluation metric.
:::
::::

:::: {#A1.SS4 .section .ltx_subsection}
### [A.4 ]{.ltx_tag .ltx_tag_subsection} Conflict Resolution (CR) {#a.4-conflict-resolution-cr .ltx_title .ltx_title_subsection}

::: {#A1.SS4.p1 .ltx_para}
We use counterfactual edit pairs from the MQUAKE \[Zhong et al., [2023](#bib.bib54){.ltx_ref}\] dataset. Each sentence containing information was assigned a number. For each edit pair, the sentence representing outdated information (the distractor) is given a smaller number, while the sentence representing more recent information (the one containing the answer) is given a larger number. We then concatenate these sentences into a long context in order according to their assigned numbers. We evaluate the CR via two datasets: [Single-Hop Editing]{.ltx_text .ltx_font_bold} and [Multi-Hop Editing]{.ltx_text .ltx_font_bold}. In these tasks, the agent's responses are mostly informational entities. Therefore, we also use SubEM (Substring Exact Match) as the evaluation metric.
:::
::::
::::::::::::::::::::::

:::::::::: {#A2 .section .ltx_appendix}
## [Appendix B ]{.ltx_tag .ltx_tag_appendix}Prompts {#appendix-b-prompts .ltx_title .ltx_title_appendix}

::: {#A2.p1 .ltx_para}
We introduce some example prompts used in this section.
:::

:::: {#A2.SS1 .section .ltx_subsection}
### [B.1 ]{.ltx_tag .ltx_tag_subsection}Instructions for Memory Construction {#b.1-instructions-for-memory-construction .ltx_title .ltx_title_subsection}

::: {#A2.SS1.p1 .ltx_para}
When processing long-context inputs, we split the content into chunks of a specified size and feed these chunks into the agent as memory. The agent can then extract relevant information from its memory based on the query to assist with query execution. This chunking approach helps organize and manage large amounts of contextual information, making retrieval and reasoning more efficient. In Figure  [[4]{.ltx_text .ltx_ref_tag}](#A2.F4 "Figure 4 ‣ B.1 Instructions for Memory Construction ‣ Appendix B Prompts ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}, we provide several example instructions that require the agent to memorize the corresponding context.
:::

<figure id="A2.F4" class="ltx_figure">
<div class="ltx_inline-block ltx_align_center ltx_transformed_outer" style="width:433.6pt;height:112.3pt;vertical-align:-112.3pt;">
<span class="ltx_transformed_inner" style="transform:translate(-43.4pt,0.0pt) scale(0.833335459359658,0.833335459359658) ;"><span class="ltx_inline-block"><img src="data:image/svg+xml;base64,PHN2ZyBpZD0iQTIuRjQucGljMSIgY2xhc3M9Imx0eF9waWN0dXJlIiBoZWlnaHQ9IjE4NS4wMiIgb3ZlcmZsb3c9InZpc2libGUiIHZlcnNpb249IjEuMSIgdmlld2JveD0iMCAwIDcyMCAxODUuMDIiIHdpZHRoPSI3MjAiPjxnIHRyYW5zZm9ybT0idHJhbnNsYXRlKDAsMTg1LjAyKSBtYXRyaXgoMSAwIDAgLTEgMCAwKSIgZmlsbD0iIzAwMDAwMCIgc3Ryb2tlPSIjMDAwMDAwIiBzdHJva2Utd2lkdGg9IjAuNHB0Ij48ZyBmaWxsPSIjMDAwMEJGIiBmaWxsLW9wYWNpdHk9IjEuMCI+PHBhdGggZD0iTSAwIDEyLjk5IEwgMCAxNzIuMDIgQyAwIDE3OS4yIDUuODIgMTg1LjAyIDEyLjk5IDE4NS4wMiBMIDcwNy4wMSAxODUuMDIgQyA3MTQuMTggMTg1LjAyIDcyMCAxNzkuMiA3MjAgMTcyLjAyIEwgNzIwIDEyLjk5IEMgNzIwIDUuODIgNzE0LjE4IDAgNzA3LjAxIDAgTCAxMi45OSAwIEMgNS44MiAwIDAgNS44MiAwIDEyLjk5IFoiIHN0eWxlPSJzdHJva2U6bm9uZSIgLz48L2c+PGcgZmlsbD0iI0Y5RjlGOSIgZmlsbC1vcGFjaXR5PSIxLjAiPjxwYXRoIGQ9Ik0gMS4xOCAxMi45OSBMIDEuMTggMTU1LjQ5IEwgNzE4LjgyIDE1NS40OSBMIDcxOC44MiAxMi45OSBDIDcxOC44MiA2LjQ3IDcxMy41MyAxLjE4IDcwNy4wMSAxLjE4IEwgMTIuOTkgMS4xOCBDIDYuNDcgMS4xOCAxLjE4IDYuNDcgMS4xOCAxMi45OSBaIiBzdHlsZT0ic3Ryb2tlOm5vbmUiIC8+PC9nPjxnIGZpbGwtb3BhY2l0eT0iMS4wIiB0cmFuc2Zvcm09Im1hdHJpeCgxLjAgMC4wIDAuMCAxLjAgMjAuODcgMTc5LjkpIj48Zm9yZWlnbm9iamVjdCB3aWR0aD0iNjc4LjI3IiBoZWlnaHQ9IjE5LjI5IiB0cmFuc2Zvcm09Im1hdHJpeCgxIDAgMCAtMSAwIDApIiBvdmVyZmxvdz0idmlzaWJsZSIgc3R5bGU9Ii0tZm9fd2lkdGggOjQ5LjAyZW07LS1mb19oZWlnaHQ6MGVtOy0tZm9fZGVwdGggOjEuMzllbTsiIGNvbG9yPSIjRkZGRkZGIj48c3BhbiBjbGFzcz0ibHR4X2ZvcmVpZ25vYmplY3RfY29udGFpbmVyIj48c3BhbiBjbGFzcz0ibHR4X2ZvcmVpZ25vYmplY3RfY29udGVudCI+CjxzcGFuIGNsYXNzPSJsdHhfaW5saW5lLWJsb2NrIGx0eF9taW5pcGFnZSBsdHhfYWxpZ25fYm90dG9tIiBzdHlsZT0id2lkdGg6NDkuMDJlbTsiPgo8c3BhbiBjbGFzcz0ibHR4X3AiPlByb21wdHMgVXNlZCBmb3IgTWVtb3J5IENvbnN0cnVjdGlvbiBvbiBWYXJpb3VzIFRhc2tzPC9zcGFuPgo8L3NwYW4+PC9zcGFuPjwvc3Bhbj48L2ZvcmVpZ25vYmplY3Q+PC9nPjxnIGZpbGwtb3BhY2l0eT0iMS4wIiB0cmFuc2Zvcm09Im1hdHJpeCgxLjAgMC4wIDAuMCAxLjAgMjAuODcgMTQzLjY3KSI+PGZvcmVpZ25vYmplY3Qgd2lkdGg9IjY3OC4yNyIgaGVpZ2h0PSIxMzAuNjgiIHRyYW5zZm9ybT0ibWF0cml4KDEgMCAwIC0xIDAgMCkiIG92ZXJmbG93PSJ2aXNpYmxlIiBzdHlsZT0iLS1mb193aWR0aCA6NDkuMDJlbTstLWZvX2hlaWdodDowZW07LS1mb19kZXB0aCA6OS40NGVtOyIgY29sb3I9IiMwMDAwMDAiPjxzcGFuIGNsYXNzPSJsdHhfZm9yZWlnbm9iamVjdF9jb250YWluZXIiPjxzcGFuIGNsYXNzPSJsdHhfZm9yZWlnbm9iamVjdF9jb250ZW50Ij4KPHNwYW4gY2xhc3M9Imx0eF9pbmxpbmUtYmxvY2sgbHR4X21pbmlwYWdlIGx0eF9hbGlnbl9ib3R0b20iIHN0eWxlPSJ3aWR0aDo0OS4wMmVtOyI+CjxzcGFuIGNsYXNzPSJsdHhfcCI+PHNwYW4gY2xhc3M9Imx0eF90ZXh0IGx0eF9mb250X2JvbGQiPklGIDxlbSBjbGFzcz0ibHR4X2VtcGggbHR4X2ZvbnRfaXRhbGljIj5Mb25nTWVtRXZhbDwvZW0+Ojwvc3Bhbj48L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+TWVtb3JpemUgdGhlIGZvbGxvd2luZyBjb252ZXJzYXRpb24gYmV0d2VlbiB0aGUgdXNlciBhbmQgdGhlIGFzc2lzdGFudDogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojMDAwMEZGOyI+Jmx0O2NodW5rJmd0Ozwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj48c3BhbiBjbGFzcz0ibHR4X3RleHQgbHR4X2ZvbnRfYm9sZCI+RUxJRiA8ZW0gY2xhc3M9Imx0eF9lbXBoIGx0eF9mb250X2l0YWxpYyI+TW92aWUgUmVjb21tZW5kYXRpb248L2VtPjo8L3NwYW4+PC9zcGFuPgo8c3BhbiBjbGFzcz0ibHR4X3AiPk1lbW9yaXplIHRoZSBmb2xsb3dpbmcgZGlhbG9ndWVzIGJldHdlZW4gYSB1c2VyIGFuZCByZWNvbW1lbmRlciBzeXN0ZW06IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6IzAwMDBGRjsiPiZsdDtjaHVuayZndDs8L3NwYW4+IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj48L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+PHNwYW4gY2xhc3M9Imx0eF90ZXh0IGx0eF9mb250X2JvbGQiPkVMSUYgPGVtIGNsYXNzPSJsdHhfZW1waCBsdHhfZm9udF9pdGFsaWMiPkZhY3QgQ29uc29saWRhdGlvbjwvZW0+Ojwvc3Bhbj48L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+TWVtb3JpemUgdGhlIHRoZXNlIGZvbGxvd2luZyBmYWN0czogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojMDAwMEZGOyI+Jmx0O2NodW5rJmd0Ozwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj48c3BhbiBjbGFzcz0ibHR4X3RleHQgbHR4X2ZvbnRfYm9sZCI+RUxTRTo8L3NwYW4+PC9zcGFuPgo8c3BhbiBjbGFzcz0ibHR4X3AiPk1lbW9yaXplIHRoZSBmb2xsb3dpbmcgY29udGVudDogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojMDAwMEZGOyI+Jmx0O2NodW5rJmd0Ozwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPjwvc3Bhbj4KPC9zcGFuPjwvc3Bhbj48L3NwYW4+PC9mb3JlaWdub2JqZWN0PjwvZz48L2c+PC9zdmc+" id="A2.F4.pic1" class="ltx_picture" /></span> </span>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 4</span>: </span><span class="ltx_text" style="font-size:90%;">The prompts we use for the agents to create the memory.</span></figcaption>
</figure>
::::

:::: {#A2.SS2 .section .ltx_subsection}
### [B.2 ]{.ltx_tag .ltx_tag_subsection}Instructions for Long-Context Agents {#b.2-instructions-for-long-context-agents .ltx_title .ltx_title_subsection}

::: {#A2.SS2.p1 .ltx_para}
In Figure  [[5]{.ltx_text .ltx_ref_tag}](#A2.F5 "Figure 5 ‣ B.2 Instructions for Long-Context Agents ‣ Appendix B Prompts ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}, we provide the examples of instructions used on different of datasets. For some existing datasets, we adjust the prompt settings from previous work such as  \[Hsieh et al., [2024](#bib.bib15){.ltx_ref}, Wu et al., [2024](#bib.bib48){.ltx_ref}\]. For example, for the dataset $\infty$[Bench-QA]{.ltx_text .ltx_font_bold} and $\infty$[Bench-Sum]{.ltx_text .ltx_font_bold}, we also insert two answer examples as [\<demo\>]{.ltx_text style="color:#004A54;"} in the prompt to help the agent better understand the questions and standardize its outputs.
:::

<figure id="A2.F5" class="ltx_figure">
<div class="ltx_inline-block ltx_align_center ltx_transformed_outer" style="width:433.6pt;height:522.2pt;vertical-align:-522.2pt;">
<span class="ltx_transformed_inner" style="transform:translate(-43.4pt,0.0pt) scale(0.833335459359658,0.833335459359658) ;"><span class="ltx_inline-block"><img src="data:image/svg+xml;base64,PHN2ZyBpZD0iQTIuRjUucGljMSIgY2xhc3M9Imx0eF9waWN0dXJlIiBoZWlnaHQ9Ijg2NS42NCIgb3ZlcmZsb3c9InZpc2libGUiIHZlcnNpb249IjEuMSIgdmlld2JveD0iMCAwIDcyMCA4NjUuNjQiIHdpZHRoPSI3MjAiPjxnIHRyYW5zZm9ybT0idHJhbnNsYXRlKDAsODY1LjY0KSBtYXRyaXgoMSAwIDAgLTEgMCAwKSIgZmlsbD0iIzAwMDAwMCIgc3Ryb2tlPSIjMDAwMDAwIiBzdHJva2Utd2lkdGg9IjAuNHB0Ij48ZyBmaWxsPSIjMDAwMEJGIiBmaWxsLW9wYWNpdHk9IjEuMCI+PHBhdGggZD0iTSAwIDEyLjk5IEwgMCA4NTIuNjUgQyAwIDg1OS44MyA1LjgyIDg2NS42NCAxMi45OSA4NjUuNjQgTCA3MDcuMDEgODY1LjY0IEMgNzE0LjE4IDg2NS42NCA3MjAgODU5LjgzIDcyMCA4NTIuNjUgTCA3MjAgMTIuOTkgQyA3MjAgNS44MiA3MTQuMTggMCA3MDcuMDEgMCBMIDEyLjk5IDAgQyA1LjgyIDAgMCA1LjgyIDAgMTIuOTkgWiIgc3R5bGU9InN0cm9rZTpub25lIiAvPjwvZz48ZyBmaWxsPSIjRjlGOUY5IiBmaWxsLW9wYWNpdHk9IjEuMCI+PHBhdGggZD0iTSAxLjE4IDEyLjk5IEwgMS4xOCA4MzYuMTEgTCA3MTguODIgODM2LjExIEwgNzE4LjgyIDEyLjk5IEMgNzE4LjgyIDYuNDcgNzEzLjUzIDEuMTggNzA3LjAxIDEuMTggTCAxMi45OSAxLjE4IEMgNi40NyAxLjE4IDEuMTggNi40NyAxLjE4IDEyLjk5IFoiIHN0eWxlPSJzdHJva2U6bm9uZSIgLz48L2c+PGcgZmlsbC1vcGFjaXR5PSIxLjAiIHRyYW5zZm9ybT0ibWF0cml4KDEuMCAwLjAgMC4wIDEuMCAyMC44NyA4NjAuNTMpIj48Zm9yZWlnbm9iamVjdCB3aWR0aD0iNjc4LjI3IiBoZWlnaHQ9IjE5LjI5IiB0cmFuc2Zvcm09Im1hdHJpeCgxIDAgMCAtMSAwIDApIiBvdmVyZmxvdz0idmlzaWJsZSIgc3R5bGU9Ii0tZm9fd2lkdGggOjQ5LjAyZW07LS1mb19oZWlnaHQ6MGVtOy0tZm9fZGVwdGggOjEuMzllbTsiIGNvbG9yPSIjRkZGRkZGIj48c3BhbiBjbGFzcz0ibHR4X2ZvcmVpZ25vYmplY3RfY29udGFpbmVyIj48c3BhbiBjbGFzcz0ibHR4X2ZvcmVpZ25vYmplY3RfY29udGVudCI+CjxzcGFuIGNsYXNzPSJsdHhfaW5saW5lLWJsb2NrIGx0eF9taW5pcGFnZSBsdHhfYWxpZ25fYm90dG9tIiBzdHlsZT0id2lkdGg6NDkuMDJlbTsiPgo8c3BhbiBjbGFzcz0ibHR4X3AiPlByb21wdHMgVXNlZCBmb3IgTG9uZy1Db250ZXh0IEFnZW50cyBvbiBWYXJpb3VzIFRhc2tzPC9zcGFuPgo8L3NwYW4+PC9zcGFuPjwvc3Bhbj48L2ZvcmVpZ25vYmplY3Q+PC9nPjxnIGZpbGwtb3BhY2l0eT0iMS4wIiB0cmFuc2Zvcm09Im1hdHJpeCgxLjAgMC4wIDAuMCAxLjAgMjAuODcgODI0LjMpIj48Zm9yZWlnbm9iamVjdCB3aWR0aD0iNjc4LjI3IiBoZWlnaHQ9IjgxMS4zMSIgdHJhbnNmb3JtPSJtYXRyaXgoMSAwIDAgLTEgMCAwKSIgb3ZlcmZsb3c9InZpc2libGUiIHN0eWxlPSItLWZvX3dpZHRoIDo0OS4wMmVtOy0tZm9faGVpZ2h0OjBlbTstLWZvX2RlcHRoIDo1OC42M2VtOyIgY29sb3I9IiMwMDAwMDAiPjxzcGFuIGNsYXNzPSJsdHhfZm9yZWlnbm9iamVjdF9jb250YWluZXIiPjxzcGFuIGNsYXNzPSJsdHhfZm9yZWlnbm9iamVjdF9jb250ZW50Ij4KPHNwYW4gY2xhc3M9Imx0eF9pbmxpbmUtYmxvY2sgbHR4X21pbmlwYWdlIGx0eF9hbGlnbl9ib3R0b20iIHN0eWxlPSJ3aWR0aDo0OS4wMmVtOyI+CjxzcGFuIGNsYXNzPSJsdHhfcCI+PGVtIGNsYXNzPSJsdHhfZW1waCBsdHhfZm9udF9ib2xkIGx0eF9mb250X2l0YWxpYyI+UlVMRVItUUE8L2VtPjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj5UaGUgY29udGV4dCBpcyBnaXZlbiBhcyBiZWxvdzogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6IzE3MUQ5MTsiPiZsdDttZW1vcnkmZ3Q7PC9zcGFuPi4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBwbGVhc2UgbWVtb3JpemUgaXQuPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBBbnN3ZXIgdGhlIHF1ZXN0aW9uIGJhc2VkIG9uIHRoZSBtZW1vcml6ZWQgZG9jdW1lbnRzLiBPbmx5IGdpdmUgbWUgdGhlIGFuc3dlciBhbmQgZG8gbm90IG91dHB1dCBhbnkgb3RoZXIgd29yZHMuIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gTm93IEFuc3dlciB0aGUgUXVlc3Rpb246IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjgwMDA7Ij4mbHQ7cXVlc3Rpb24mZ3Q7PC9zcGFuPiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IEFuc3dlcjo8L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+PGVtIGNsYXNzPSJsdHhfZW1waCBsdHhfZm9udF9ib2xkIGx0eF9mb250X2l0YWxpYyI+UlVMRVItTklBSC1NUTwvZW0+PC9zcGFuPgo8c3BhbiBjbGFzcz0ibHR4X3AiPlRoZSBjb250ZXh0IGlzIGdpdmVuIGFzIGJlbG93OiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojMTcxRDkxOyI+Jmx0O21lbW9yeSZndDs8L3NwYW4+LiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IFBsZWFzZSBtZW1vcml6ZSBpdC4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBTb21lIHNwZWNpYWwgbWFnaWMgbnVtYmVycyBhcmUgaGlkZGVuIHdpdGhpbiB0aGUgbWVtb3JpemVkIHRleHQuIE1ha2Ugc3VyZSB0byBtZW1vcml6ZSBpdC4gSSB3aWxsIHF1aXogeW91IGFib3V0IHRoZSBudW1iZXJzIGFmdGVyd2FyZHMuPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBOb3cgQW5zd2VyIHRoZSBRdWVzdGlvbjogV2hhdCBhcmUgYWxsIHRoZSBzcGVjaWFsIG1hZ2ljIG51bWJlcnMgZm9yIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjgwMDA7Ij4mbHQ7cXVlc3Rpb24mZ3Q7PC9zcGFuPiBtZW50aW9uZWQgaW4gdGhlIG1lbW9yaXplZCB0ZXh0PyA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IFRoZSBzcGVjaWFsIG1hZ2ljIG51bWJlcnMgZm9yIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjgwMDA7Ij4mbHQ7cXVlc3Rpb24mZ3Q7PC9zcGFuPiBtZW50aW9uZWQgaW4gdGhlIG1lbW9yaXplIHRleHQgYXJlOjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj48ZW0gY2xhc3M9Imx0eF9lbXBoIGx0eF9mb250X2l0YWxpYyI+PG1hdGggaWQ9IkEyLkY1LnBpYzEubTEiIGNsYXNzPSJsdHhfTWF0aCIgYWx0dGV4dD0iXGluZnR5IiBkaXNwbGF5PSJpbmxpbmUiPjxzZW1hbnRpY3M+PG1pIG1hdGh2YXJpYW50PSJub3JtYWwiPuKInjwvbWk+PGFubm90YXRpb24gZW5jb2Rpbmc9ImFwcGxpY2F0aW9uL3gtdGV4Ij5caW5mdHk8L2Fubm90YXRpb24+PC9zZW1hbnRpY3M+PC9tYXRoPjxzcGFuIGNsYXNzPSJsdHhfdGV4dCBsdHhfZm9udF9ib2xkIj5CZW5jaC1RQTwvc3Bhbj48L2VtPjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj5UaGUgY29udGV4dCBpcyBnaXZlbiBhcyBiZWxvdzogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6IzE3MUQ5MTsiPiZsdDttZW1vcnkmZ3Q7PC9zcGFuPi4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBQbGVhc2UgbWVtb3JpemUgaXQuIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gQmFzZWQgb24gdGhlIGNvbnRleHQgeW91IG1lbW9yaXplZCwgYW5zd2VyIHRoZSBxdWVzdGlvbiBhcyBjb25jaXNlbHkgYXMgeW91IGNhbiwgdXNpbmcgYSBzaW5nbGUgcGhyYXNlIGlmIHBvc3NpYmxlLiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiMwMDRBNTQ7Ij4mbHQ7ZGVtbyZndDs8L3NwYW4+LjxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gTm93IEFuc3dlciB0aGUgUXVlc3Rpb246IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjgwMDA7Ij4mbHQ7cXVlc3Rpb24mZ3Q7PC9zcGFuPi48c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IEFuc3dlcjo8L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+PGVtIGNsYXNzPSJsdHhfZW1waCBsdHhfZm9udF9ib2xkIGx0eF9mb250X2l0YWxpYyI+TG9uZ01lbUV2YWw8L2VtPjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj5IZXJlIGFyZSBzZXZlcmFsIGhpc3RvcnkgY2hhdHMgYmV0d2VlbiB5b3UgYW5kIGEgdXNlciA6IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiMxNzFEOTE7Ij4mbHQ7bWVtb3J5Jmd0Ozwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBQbGVhc2UgbWVtb3JpemUgdGhlbS4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBUaGUgaGlzdG9yeSBjaGF0cyBhcmUgYmV0d2VlbiB5b3UgYW5kIGEgdXNlci4gQmFzZWQgb24gdGhlIHJlbGV2YW50IGNoYXQgaGlzdG9yeSwgYW5zd2VyIHRoZSBxdWVzdGlvbiBhcyBjb25jaXNlbHkgYXMgeW91IGNhbiwgdXNpbmcgYSBzaW5nbGUgcGhyYXNlIGlmIHBvc3NpYmxlLjxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gQ3VycmVudCBEYXRlOiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojMEFBMzQ0OyI+Jmx0O3F1ZXN0aW9uX2RhdGUmZ3Q7PC9zcGFuPiwgPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBOb3cgQW5zd2VyIHRoZSBRdWVzdGlvbjogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGODAwMDsiPiZsdDtxdWVzdGlvbiZndDs8L3NwYW4+IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gQW5zd2VyOjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj48ZW0gY2xhc3M9Imx0eF9lbXBoIGx0eF9mb250X2JvbGQgbHR4X2ZvbnRfaXRhbGljIj5FdmVudFFBPC9lbT48L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+VGhlIGNvbnRleHQgaXMgZ2l2ZW4gYXMgYmVsb3c6IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiMxNzFEOTE7Ij4mbHQ7bWVtb3J5Jmd0Ozwvc3Bhbj4uIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gUGxlYXNlIG1lbW9yaXplIGl0LiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IEJhc2VkIG9uIHRoZSBjb250ZXh0IHlvdSBtZW1vcml6ZWQsIGNvbXBsZXRlIHRoZSB0YXNrIGJlbG93OiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IFRoZXNlIGFyZSB0aGUgZXZlbnRzIHRoYXQgaGF2ZSBhbHJlYWR5IG9jY3VycmVkOiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRkEzQTM7Ij4mbHQ7cHJldmlvdXNfZXZlbnRzJmd0Ozwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBCZWxvdyBpcyBhIGxpc3Qgb2YgcG9zc2libGUgc3Vic2VxdWVudCBldmVudHM6PHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkY4MDAwOyI+Jmx0O3F1ZXN0aW9uJmd0Ozwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBZb3VyIHRhc2sgaXMgdG8gY2hvb3NlIGZyb20gdGhlIGFib3ZlIGV2ZW50cyB3aGljaCBldmVudCBoYXBwZW5zIG5leHQgYmFzZWQgb24gdGhlIGJvb2sgZXhjZXJwdC4gSW4geW91ciByZXNwb25zZSB0byBtZSwgb25seSBpbmNsdWRlIHRoZSBhbnN3ZXIgd2l0aG91dCBhbnl0aGluZyBlbHNlLiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IFRoZSBldmVudCB0aGF0IGhhcHBlbnMgbmV4dCBpczo8L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+PGVtIGNsYXNzPSJsdHhfZW1waCBsdHhfZm9udF9ib2xkIGx0eF9mb250X2l0YWxpYyI+TGFiZWwgTWF0Y2hpbmcgKEJBTktJTkc3NywgZXRjLik8L2VtPjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj5UaGUgY29udGV4dCBpcyBnaXZlbiBhcyBiZWxvdzogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6IzE3MUQ5MTsiPiZsdDttZW1vcnkmZ3Q7PC9zcGFuPi4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBQbGVhc2UgbWVtb3JpemUgaXQuIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gVXNlIHRoZSBwcm92aWRlZCBtYXBwaW5nIGZyb20gdGhlIGNvbnRleHQgdG8gbnVtZXJpY2FsIGxhYmVsIHRvIGFzc2lnbiBhIG51bWVyaWNhbCBsYWJlbCB0byB0aGUgY29udGV4dC4gT25seSBvdXRwdXQgJnF1b3Q7bGFiZWw6IHt7bGFiZWx9fSZxdW90OyBhbmQgbm90aGluZyBlbHNlLiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IFF1ZXN0aW9uOiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkY4MDAwOyI+Jmx0O3F1ZXN0aW9uJmd0Ozwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBsYWJlbDo8L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+PGVtIGNsYXNzPSJsdHhfZW1waCBsdHhfZm9udF9ib2xkIGx0eF9mb250X2l0YWxpYyI+TW92aWUgUmVjb21tZW5kYXRpb248L2VtPjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj5IZXJlIGFyZSBkaWFsb2d1ZXMgYmV0d2VlbiBhIHVzZXIgYW5kIHJlY29tbWVuZGVyIHN5c3RlbTogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6IzE3MUQ5MTsiPiZsdDttZW1vcnkmZ3Q7PC9zcGFuPi4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBQbGVhc2UgbWVtb3JpemUgdGhlbS4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBQcmV0ZW5kIHlvdSBhcmUgYSBtb3ZpZSByZWNvbW1lbmRlciBzeXN0ZW0uIFlvdSBuZWVkIHRvIHJlY29tbWVuZCBtb3ZpZXMgYmFzZWQgb24gdGhlIGRpYWxvZ3VlcyB5b3UgaGF2ZSBtZW1vcml6ZWQuIE5vdyBJIHdpbGwgZ2l2ZSB5b3UgYSBuZXcgY29udmVyc2F0aW9uIGJldHdlZW4gYSB1c2VyIGFuZCB5b3UgKGEgcmVjb21tZW5kZXIgc3lzdGVtKS4gQmFzZWQgb24gdGhlIGNvbnZlcnNhdGlvbiwgeW91IHJlcGx5IG1lIHdpdGggMjAgcmVjb21tZW5kYXRpb25zIHdpdGhvdXQgZXh0cmEgc2VudGVuY2VzLiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IEZvciBFeGFtcGxlOjxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gW0NvbnZlcnNhdGlvbl0gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBUaGUgcmVjb21tZW5kYXRpb25zIGFyZTogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiAxLm1vdmllMSA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IDIubW92aWUyIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4g4oCmPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBIZXJlIGlzIHRoZSBjb252ZXJzYXRpb246IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjgwMDA7Ij4mbHQ7cXVlc3Rpb24mZ3Q7PC9zcGFuPiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IFRoZSByZWNvbW1lbmRhdGlvbnMgYXJlOjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj48ZW0gY2xhc3M9Imx0eF9lbXBoIGx0eF9mb250X2l0YWxpYyI+PG1hdGggaWQ9IkEyLkY1LnBpYzEubTIiIGNsYXNzPSJsdHhfTWF0aCIgYWx0dGV4dD0iXGluZnR5IiBkaXNwbGF5PSJpbmxpbmUiPjxzZW1hbnRpY3M+PG1pIG1hdGh2YXJpYW50PSJub3JtYWwiPuKInjwvbWk+PGFubm90YXRpb24gZW5jb2Rpbmc9ImFwcGxpY2F0aW9uL3gtdGV4Ij5caW5mdHk8L2Fubm90YXRpb24+PC9zZW1hbnRpY3M+PC9tYXRoPjxzcGFuIGNsYXNzPSJsdHhfdGV4dCBsdHhfZm9udF9ib2xkIj5CZW5jaC1TdW08L3NwYW4+PC9lbT48L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+VGhlIGJvb2sgaXMgZ2l2ZW4gYXMgYmVsb3c6IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiMxNzFEOTE7Ij4mbHQ7bWVtb3J5Jmd0Ozwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBQbGVhc2UgbWVtb3JpemUgaXQuIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gWW91IGFyZSBnaXZlbiBhIGJvb2sgYWJvdmUgYW5kIHlvdSBhcmUgdGFza2VkIHRvIHN1bW1hcml6ZSBpdC4gV3JpdGUgYSBzdW1tYXJ5IG9mIGFib3V0IDEwMDAgdG8gMTIwMCB3b3Jkcy4gT25seSB3cml0ZSBhYm91dCB0aGUgcGxvdCBhbmQgY2hhcmFjdGVycyBvZiB0aGUgc3RvcnkuIERvIG5vdCBkaXNjdXNzIHRoZSB0aGVtZXMgb3IgYmFja2dyb3VuZCBvZiB0aGUgYm9vay4gRG8gbm90IHByb3ZpZGUgYW55IGFuYWx5c2lzIG9yIGNvbW1lbnRhcnkuIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6IzAwNEE1NDsiPiZsdDtkZW1vJmd0Ozwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBOb3cgc3VtbWFyaXplIHRoZSBib29rLjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj48ZW0gY2xhc3M9Imx0eF9lbXBoIGx0eF9mb250X2JvbGQgbHR4X2ZvbnRfaXRhbGljIj5GYWN0IENvbnNvbGlkYXRpb248L2VtPjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj5IZXJlIGlzIGEga25vd2xlZGdlIHBvb2wgd2l0aCBsb3RzIG9mIG5ldyBmYWN0czogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6IzE3MUQ5MTsiPiZsdDttZW1vcnkmZ3Q7PC9zcGFuPi4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBQbGVhc2UgbWVtb3JpemUgaXQuIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gUHJldGVuZCB5b3UgYXJlIGEga25vd2xlZGdlIG1hbmFnZW1lbnQgc3lzdGVtLiBFYWNoIGZhY3QgaW4gdGhlIGtub3dsZWRnZSBwb29sIGlzIHByb3ZpZGVkIHdpdGggYSBzZXJpYWwgbnVtYmVyIGF0IHRoZSBiZWdpbm5pbmcsIGFuZCB0aGUgbmV3ZXIgZmFjdCBoYXMgbGFyZ2VyIHNlcmlhbCBudW1iZXIuIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gWW91IG5lZWQgdG8gc29sdmUgdGhlIGNvbmZsaWN0cyBvZiBmYWN0cyBpbiB0aGUga25vd2xlZGdlIHBvb2wgYnkgZmluZGluZyB0aGUgbmV3ZXN0IGZhY3QuIFlvdSBuZWVkIHRvIGFuc3dlciBhIHF1ZXN0aW9uIGJhc2VkIG9uIHRoaXMgcnVsZS4gWW91IHNob3VsZCBnaXZlIGEgdmVyeSBjb25jaXNlIGFuc3dlciB3aXRob3V0IHNheWluZyBvdGhlciB3b3JkcyBmb3IgdGhlIHF1ZXN0aW9uICoqb25seSoqIGZyb20gdGhlIGtub3dsZWRnZSBwb29sIHlvdSBoYXZlIG1lbW9yaXplZCByYXRoZXIgdGhhbiB0aGUgcmVhbCBmYWN0cyBpbiByZWFsIHdvcmxkLiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IEZvciBleGFtcGxlOiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IFtLbm93bGVkZ2UgUG9vbF0gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBRdWVzdGlvbjogQmFzZWQgb24gdGhlIHByb3ZpZGVkIEtub3dsZWRnZSBQb29sLCB3aGF0IGlzIHRoZSBuYW1lIG9mIHRoZSBjdXJyZW50IHByZXNpZGVudCBvZiBDb3VudHJ5IFI/IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gQW5zd2VyOiBQZXJzb24gRC4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBOb3cgQW5zd2VyIHRoZSBRdWVzdGlvbjogQmFzZWQgb24gdGhlIHByb3ZpZGVkIEtub3dsZWRnZSBQb29sLCA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkY4MDAwOyI+Jmx0O3F1ZXN0aW9uJmd0Ozwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBBbnN3ZXI6PC9zcGFuPgo8L3NwYW4+PC9zcGFuPjwvc3Bhbj48L2ZvcmVpZ25vYmplY3Q+PC9nPjwvZz48L3N2Zz4=" id="A2.F5.pic1" class="ltx_picture" /></span> </span>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 5</span>: </span><span class="ltx_text" style="font-size:90%;">The prompts we use for the <em>Long-Context Agents</em> in Table <a href="#S4.T2" class="ltx_ref" title="Table 2 ‣ 4.1 Experimental Setup ‣ 4 Experiments ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"><span class="ltx_text ltx_ref_tag">2</span></a>. Here &lt;memory&gt; refers to the accumulated text from the sequential inputs.</span></figcaption>
</figure>
::::

:::: {#A2.SS3 .section .ltx_subsection}
### [B.3 ]{.ltx_tag .ltx_tag_subsection}Instructions for RAG Agents {#b.3-instructions-for-rag-agents .ltx_title .ltx_title_subsection}

::: {#A2.SS3.p1 .ltx_para}
We provide examples of prompts used for the RAG based Agents in Figure  [[6]{.ltx_text .ltx_ref_tag}](#A2.F6 "Figure 6 ‣ B.3 Instructions for RAG Agents ‣ Appendix B Prompts ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}. For this type of agent, after storing the input long context in memory, we use [\<question\>]{.ltx_text style="color:#FF8000;"} as the memory retrieval query for most tasks. But for task [RULER-NIAH-MQ]{.ltx_text .ltx_font_bold}, we use the entire question \"What are all the special magic numbers for [\<question\>]{.ltx_text style="color:#FF8000;"} mentioned in the memorized text?\" as the query. And for $\infty$[Bench-Sum]{.ltx_text .ltx_font_bold}, we use the entire query without the [\<demo\>]{.ltx_text style="color:#004A54;"} for the memory retrieval.
:::

<figure id="A2.F6" class="ltx_figure">
<div class="ltx_inline-block ltx_align_center ltx_transformed_outer" style="width:433.6pt;height:522.2pt;vertical-align:-522.2pt;">
<span class="ltx_transformed_inner" style="transform:translate(-43.4pt,0.0pt) scale(0.833335459359658,0.833335459359658) ;"><span class="ltx_inline-block"><img src="data:image/svg+xml;base64,PHN2ZyBpZD0iQTIuRjYucGljMSIgY2xhc3M9Imx0eF9waWN0dXJlIiBoZWlnaHQ9Ijg2NS42NCIgb3ZlcmZsb3c9InZpc2libGUiIHZlcnNpb249IjEuMSIgdmlld2JveD0iMCAwIDcyMCA4NjUuNjQiIHdpZHRoPSI3MjAiPjxnIHRyYW5zZm9ybT0idHJhbnNsYXRlKDAsODY1LjY0KSBtYXRyaXgoMSAwIDAgLTEgMCAwKSIgZmlsbD0iIzAwMDAwMCIgc3Ryb2tlPSIjMDAwMDAwIiBzdHJva2Utd2lkdGg9IjAuNHB0Ij48ZyBmaWxsPSIjMDAwMEJGIiBmaWxsLW9wYWNpdHk9IjEuMCI+PHBhdGggZD0iTSAwIDEyLjk5IEwgMCA4NTIuNjUgQyAwIDg1OS44MyA1LjgyIDg2NS42NCAxMi45OSA4NjUuNjQgTCA3MDcuMDEgODY1LjY0IEMgNzE0LjE4IDg2NS42NCA3MjAgODU5LjgzIDcyMCA4NTIuNjUgTCA3MjAgMTIuOTkgQyA3MjAgNS44MiA3MTQuMTggMCA3MDcuMDEgMCBMIDEyLjk5IDAgQyA1LjgyIDAgMCA1LjgyIDAgMTIuOTkgWiIgc3R5bGU9InN0cm9rZTpub25lIiAvPjwvZz48ZyBmaWxsPSIjRjlGOUY5IiBmaWxsLW9wYWNpdHk9IjEuMCI+PHBhdGggZD0iTSAxLjE4IDEyLjk5IEwgMS4xOCA4MzYuMTEgTCA3MTguODIgODM2LjExIEwgNzE4LjgyIDEyLjk5IEMgNzE4LjgyIDYuNDcgNzEzLjUzIDEuMTggNzA3LjAxIDEuMTggTCAxMi45OSAxLjE4IEMgNi40NyAxLjE4IDEuMTggNi40NyAxLjE4IDEyLjk5IFoiIHN0eWxlPSJzdHJva2U6bm9uZSIgLz48L2c+PGcgZmlsbC1vcGFjaXR5PSIxLjAiIHRyYW5zZm9ybT0ibWF0cml4KDEuMCAwLjAgMC4wIDEuMCAyMC44NyA4NjAuNTMpIj48Zm9yZWlnbm9iamVjdCB3aWR0aD0iNjc4LjI3IiBoZWlnaHQ9IjE5LjI5IiB0cmFuc2Zvcm09Im1hdHJpeCgxIDAgMCAtMSAwIDApIiBvdmVyZmxvdz0idmlzaWJsZSIgc3R5bGU9Ii0tZm9fd2lkdGggOjQ5LjAyZW07LS1mb19oZWlnaHQ6MGVtOy0tZm9fZGVwdGggOjEuMzllbTsiIGNvbG9yPSIjRkZGRkZGIj48c3BhbiBjbGFzcz0ibHR4X2ZvcmVpZ25vYmplY3RfY29udGFpbmVyIj48c3BhbiBjbGFzcz0ibHR4X2ZvcmVpZ25vYmplY3RfY29udGVudCI+CjxzcGFuIGNsYXNzPSJsdHhfaW5saW5lLWJsb2NrIGx0eF9taW5pcGFnZSBsdHhfYWxpZ25fYm90dG9tIiBzdHlsZT0id2lkdGg6NDkuMDJlbTsiPgo8c3BhbiBjbGFzcz0ibHR4X3AiPlByb21wdHMgVXNlZCBmb3IgUkFHIEJhc2VkIEFnZW50cyBvbiBWYXJpb3VzIFRhc2tzPC9zcGFuPgo8L3NwYW4+PC9zcGFuPjwvc3Bhbj48L2ZvcmVpZ25vYmplY3Q+PC9nPjxnIGZpbGwtb3BhY2l0eT0iMS4wIiB0cmFuc2Zvcm09Im1hdHJpeCgxLjAgMC4wIDAuMCAxLjAgMjAuODcgODI0LjMpIj48Zm9yZWlnbm9iamVjdCB3aWR0aD0iNjc4LjI3IiBoZWlnaHQ9IjgxMS4zMSIgdHJhbnNmb3JtPSJtYXRyaXgoMSAwIDAgLTEgMCAwKSIgb3ZlcmZsb3c9InZpc2libGUiIHN0eWxlPSItLWZvX3dpZHRoIDo0OS4wMmVtOy0tZm9faGVpZ2h0OjBlbTstLWZvX2RlcHRoIDo1OC42M2VtOyIgY29sb3I9IiMwMDAwMDAiPjxzcGFuIGNsYXNzPSJsdHhfZm9yZWlnbm9iamVjdF9jb250YWluZXIiPjxzcGFuIGNsYXNzPSJsdHhfZm9yZWlnbm9iamVjdF9jb250ZW50Ij4KPHNwYW4gY2xhc3M9Imx0eF9pbmxpbmUtYmxvY2sgbHR4X21pbmlwYWdlIGx0eF9hbGlnbl9ib3R0b20iIHN0eWxlPSJ3aWR0aDo0OS4wMmVtOyI+CjxzcGFuIGNsYXNzPSJsdHhfcCI+PGVtIGNsYXNzPSJsdHhfZW1waCBsdHhfZm9udF9ib2xkIGx0eF9mb250X2l0YWxpYyI+UlVMRVItUUE8L2VtPjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj5IZXJlIGlzIHRoZSBjb250ZXh0IHJldHJpZXZlZCBmcm9tIG1lbW9yeTogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6IzE3MUQ5MTsiPiZsdDttZW1vcnkmZ3Q7PC9zcGFuPi48c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IEFuc3dlciB0aGUgcXVlc3Rpb24gYmFzZWQgb24gdGhlIHJldHJpZXZlZCBjb250ZXh0LiBPbmx5IGdpdmUgbWUgdGhlIGFuc3dlciBhbmQgZG8gbm90IG91dHB1dCBhbnkgb3RoZXIgd29yZHMuIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gTm93IEFuc3dlciB0aGUgUXVlc3Rpb246IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjgwMDA7Ij4mbHQ7cXVlc3Rpb24mZ3Q7PC9zcGFuPiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IEFuc3dlcjo8L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+PGVtIGNsYXNzPSJsdHhfZW1waCBsdHhfZm9udF9ib2xkIGx0eF9mb250X2l0YWxpYyI+UlVMRVItTklBSC1NUTwvZW0+PC9zcGFuPgo8c3BhbiBjbGFzcz0ibHR4X3AiPkhlcmUgaXMgdGhlIGNvbnRleHQgcmV0cmlldmVkIGZyb20gbWVtb3J5OiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojMTcxRDkxOyI+Jmx0O21lbW9yeSZndDs8L3NwYW4+LjxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gU29tZSBzcGVjaWFsIG1hZ2ljIG51bWJlcnMgYXJlIGhpZGRlbiB3aXRoaW4gdGhlIHJldHJpZXZlZCB0ZXh0LiBNYWtlIHN1cmUgdG8gbWVtb3JpemUgaXQuIEkgd2lsbCBxdWl6IHlvdSBhYm91dCB0aGUgbnVtYmVycyBhZnRlcndhcmRzLjxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gTm93IEFuc3dlciB0aGUgUXVlc3Rpb246IFdoYXQgYXJlIGFsbCB0aGUgc3BlY2lhbCBtYWdpYyBudW1iZXJzIGZvciA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkY4MDAwOyI+Jmx0O3F1ZXN0aW9uJmd0Ozwvc3Bhbj4gbWVudGlvbmVkIGluIHRoZSBtZW1vcml6ZWQgdGV4dD8gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBUaGUgc3BlY2lhbCBtYWdpYyBudW1iZXJzIGZvciA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkY4MDAwOyI+Jmx0O3F1ZXN0aW9uJmd0Ozwvc3Bhbj4gbWVudGlvbmVkIGluIHRoZSBtZW1vcml6ZSB0ZXh0IGFyZTo8L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+PGVtIGNsYXNzPSJsdHhfZW1waCBsdHhfZm9udF9pdGFsaWMiPjxtYXRoIGlkPSJBMi5GNi5waWMxLm0xIiBjbGFzcz0ibHR4X01hdGgiIGFsdHRleHQ9IlxpbmZ0eSIgZGlzcGxheT0iaW5saW5lIj48c2VtYW50aWNzPjxtaSBtYXRodmFyaWFudD0ibm9ybWFsIj7iiJ48L21pPjxhbm5vdGF0aW9uIGVuY29kaW5nPSJhcHBsaWNhdGlvbi94LXRleCI+XGluZnR5PC9hbm5vdGF0aW9uPjwvc2VtYW50aWNzPjwvbWF0aD48c3BhbiBjbGFzcz0ibHR4X3RleHQgbHR4X2ZvbnRfYm9sZCI+QmVuY2gtUUE8L3NwYW4+PC9lbT48L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+SGVyZSBpcyB0aGUgY29udGV4dCByZXRyaWV2ZWQgZnJvbSBtZW1vcnk6IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiMxNzFEOTE7Ij4mbHQ7bWVtb3J5Jmd0Ozwvc3Bhbj4uPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBCYXNlZCBvbiB0aGUgY29udGV4dCB5b3UgcmV0cmlldmVkLCBhbnN3ZXIgdGhlIHF1ZXN0aW9uIGFzIGNvbmNpc2VseSBhcyB5b3UgY2FuLCB1c2luZyBhIHNpbmdsZSBwaHJhc2UgaWYgcG9zc2libGUuIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6IzAwNEE1NDsiPiZsdDtkZW1vJmd0Ozwvc3Bhbj4uPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBOb3cgQW5zd2VyIHRoZSBRdWVzdGlvbjogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGODAwMDsiPiZsdDtxdWVzdGlvbiZndDs8L3NwYW4+LjxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gQW5zd2VyOjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj48ZW0gY2xhc3M9Imx0eF9lbXBoIGx0eF9mb250X2JvbGQgbHR4X2ZvbnRfaXRhbGljIj5Mb25nTWVtRXZhbDwvZW0+PC9zcGFuPgo8c3BhbiBjbGFzcz0ibHR4X3AiPkhlcmUgYXJlIHJldHJpZXZlZCBzZXZlcmFsIGhpc3RvcnkgY2hhdHMgYmV0d2VlbiB5b3UgYW5kIGEgdXNlciBmcm9tIG1lbW9yeTogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6IzE3MUQ5MTsiPiZsdDttZW1vcnkmZ3Q7PC9zcGFuPiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IFRoZSByZXRyaWV2ZWQgaGlzdG9yeSBjaGF0cyBhcmUgYmV0d2VlbiB5b3UgYW5kIGEgdXNlci4gQmFzZWQgb24gdGhlIHJlbGV2YW50IGNoYXQgaGlzdG9yeSwgYW5zd2VyIHRoZSBxdWVzdGlvbiBhcyBjb25jaXNlbHkgYXMgeW91IGNhbiwgdXNpbmcgYSBzaW5nbGUgcGhyYXNlIGlmIHBvc3NpYmxlLjxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gQ3VycmVudCBEYXRlOiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojMEFBMzQ0OyI+Jmx0O3F1ZXN0aW9uX2RhdGUmZ3Q7PC9zcGFuPiwgPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBOb3cgQW5zd2VyIHRoZSBRdWVzdGlvbjogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGODAwMDsiPiZsdDtxdWVzdGlvbiZndDs8L3NwYW4+IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gQW5zd2VyOjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj48ZW0gY2xhc3M9Imx0eF9lbXBoIGx0eF9mb250X2JvbGQgbHR4X2ZvbnRfaXRhbGljIj5FdmVudFFBPC9lbT48L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+SGVyZSBpcyB0aGUgY29udGV4dCByZXRyaWV2ZWQgZnJvbSBtZW1vcnk6IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiMxNzFEOTE7Ij4mbHQ7bWVtb3J5Jmd0Ozwvc3Bhbj4uPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBCYXNlZCBvbiB0aGUgY29udGV4dCB5b3UgcmV0cmlldmVkLCBjb21wbGV0ZSB0aGUgdGFzayBiZWxvdzogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBUaGVzZSBhcmUgdGhlIGV2ZW50cyB0aGF0IGhhdmUgYWxyZWFkeSBvY2N1cnJlZDogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkZBM0EzOyI+Jmx0O3ByZXZpb3VzX2V2ZW50cyZndDs8L3NwYW4+IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gQmVsb3cgaXMgYSBsaXN0IG9mIHBvc3NpYmxlIHN1YnNlcXVlbnQgZXZlbnRzOjxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGODAwMDsiPiZsdDtxdWVzdGlvbiZndDs8L3NwYW4+IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gWW91ciB0YXNrIGlzIHRvIGNob29zZSBmcm9tIHRoZSBhYm92ZSBldmVudHMgd2hpY2ggZXZlbnQgaGFwcGVucyBuZXh0IGJhc2VkIG9uIHRoZSBib29rIGV4Y2VycHQuIEluIHlvdXIgcmVzcG9uc2UgdG8gbWUsIG9ubHkgaW5jbHVkZSB0aGUgYW5zd2VyIHdpdGhvdXQgYW55dGhpbmcgZWxzZS4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBUaGUgZXZlbnQgdGhhdCBoYXBwZW5zIG5leHQgaXM6PC9zcGFuPgo8c3BhbiBjbGFzcz0ibHR4X3AiPjxlbSBjbGFzcz0ibHR4X2VtcGggbHR4X2ZvbnRfYm9sZCBsdHhfZm9udF9pdGFsaWMiPkxhYmVsIE1hdGNoaW5nIChCQU5LSU5HNzcsIGV0Yy4pPC9lbT48L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+SGVyZSBhcmUgdGhlIGV4YW1wbGVzIHJldHJpZXZlZCBmcm9tIG1lbW9yeTogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6IzE3MUQ5MTsiPiZsdDttZW1vcnkmZ3Q7PC9zcGFuPi4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBVc2UgdGhlIHJldHJpZXZlZCBtYXBwaW5nIGZyb20gdGhlIGNvbnRleHQgdG8gbnVtZXJpY2FsIGxhYmVsIHRvIGFzc2lnbiBhIG51bWVyaWNhbCBsYWJlbCB0byB0aGUgY29udGV4dC4gT25seSBvdXRwdXQgJnF1b3Q7bGFiZWw6IHt7bGFiZWx9fSZxdW90OyBhbmQgbm90aGluZyBlbHNlLiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IFF1ZXN0aW9uOiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkY4MDAwOyI+Jmx0O3F1ZXN0aW9uJmd0Ozwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBsYWJlbDo8L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+PGVtIGNsYXNzPSJsdHhfZW1waCBsdHhfZm9udF9ib2xkIGx0eF9mb250X2l0YWxpYyI+TW92aWUgUmVjb21tZW5kYXRpb248L2VtPjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj5IZXJlIGFyZSByZXRyaWV2ZWQgZGlhbG9ndWVzIGJldHdlZW4gYSB1c2VyIGFuZCByZWNvbW1lbmRlciBzeXN0ZW0gZnJvbSBtZW1vcnk6IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiMxNzFEOTE7Ij4mbHQ7bWVtb3J5Jmd0Ozwvc3Bhbj4uIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gUHJldGVuZCB5b3UgYXJlIGEgbW92aWUgcmVjb21tZW5kZXIgc3lzdGVtLiBZb3UgbmVlZCB0byByZWNvbW1lbmQgbW92aWVzIGJhc2VkIG9uIHRoZSBleGFtcGxlIGRpYWxvZ3VlcyB5b3UgaGF2ZSByZXRyaWV2ZWQuIE5vdyBJIHdpbGwgZ2l2ZSB5b3UgYSBuZXcgY29udmVyc2F0aW9uIGJldHdlZW4gYSB1c2VyIGFuZCB5b3UgKGEgcmVjb21tZW5kZXIgc3lzdGVtKS4gQmFzZWQgb24gdGhlIGNvbnZlcnNhdGlvbiwgeW91IHJlcGx5IG1lIHdpdGggMjAgcmVjb21tZW5kYXRpb25zIHdpdGhvdXQgZXh0cmEgc2VudGVuY2VzLiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IEZvciBFeGFtcGxlOjxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gW0NvbnZlcnNhdGlvbl0gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBUaGUgcmVjb21tZW5kYXRpb25zIGFyZTogPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiAxLm1vdmllMSA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IDIubW92aWUyIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4g4oCmPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBIZXJlIGlzIHRoZSBjb252ZXJzYXRpb246IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjgwMDA7Ij4mbHQ7cXVlc3Rpb24mZ3Q7PC9zcGFuPiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IFRoZSByZWNvbW1lbmRhdGlvbnMgYXJlOjwvc3Bhbj4KPHNwYW4gY2xhc3M9Imx0eF9wIj48ZW0gY2xhc3M9Imx0eF9lbXBoIGx0eF9mb250X2l0YWxpYyI+PG1hdGggaWQ9IkEyLkY2LnBpYzEubTIiIGNsYXNzPSJsdHhfTWF0aCIgYWx0dGV4dD0iXGluZnR5IiBkaXNwbGF5PSJpbmxpbmUiPjxzZW1hbnRpY3M+PG1pIG1hdGh2YXJpYW50PSJub3JtYWwiPuKInjwvbWk+PGFubm90YXRpb24gZW5jb2Rpbmc9ImFwcGxpY2F0aW9uL3gtdGV4Ij5caW5mdHk8L2Fubm90YXRpb24+PC9zZW1hbnRpY3M+PC9tYXRoPjxzcGFuIGNsYXNzPSJsdHhfdGV4dCBsdHhfZm9udF9ib2xkIj5CZW5jaC1TdW08L3NwYW4+PC9lbT48L3NwYW4+CjxzcGFuIGNsYXNzPSJsdHhfcCI+VGhlIGJvb2sgY29udGV4dCBpcyByZXRyaWV2ZWQgZnJvbSBtZW1vcnkgYW5kIGl0IGlzIGdpdmVuIGFzIGJlbG93OiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojMTcxRDkxOyI+Jmx0O21lbW9yeSZndDs8L3NwYW4+IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gWW91IGFyZSBnaXZlbiByZXRyaWV2ZWQgY29udGV4dCBhYm92ZSBhbmQgeW91IGFyZSB0YXNrZWQgdG8gc3VtbWFyaXplIGl0LiBXcml0ZSBhIHN1bW1hcnkgb2YgYWJvdXQgMTAwMCB0byAxMjAwIHdvcmRzLiBPbmx5IHdyaXRlIGFib3V0IHRoZSBwbG90IGFuZCBjaGFyYWN0ZXJzIG9mIHRoZSBzdG9yeS4gRG8gbm90IGRpc2N1c3MgdGhlIHRoZW1lcyBvciBiYWNrZ3JvdW5kIG9mIHRoZSBib29rLiBEbyBub3QgcHJvdmlkZSBhbnkgYW5hbHlzaXMgb3IgY29tbWVudGFyeS4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojMDA0QTU0OyI+Jmx0O2RlbW8mZ3Q7PC9zcGFuPiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IE5vdyBzdW1tYXJpemUgdGhlIGJvb2suPC9zcGFuPgo8c3BhbiBjbGFzcz0ibHR4X3AiPjxlbSBjbGFzcz0ibHR4X2VtcGggbHR4X2ZvbnRfYm9sZCBsdHhfZm9udF9pdGFsaWMiPkZhY3QgQ29uc29saWRhdGlvbjwvZW0+PC9zcGFuPgo8c3BhbiBjbGFzcz0ibHR4X3AiPkhlcmUgaXMgYSBsaXN0IG9mIGtub3dsZWRnZSByZXRyaWV2ZWQgZnJvbSBtZW1vcnk6IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiMxNzFEOTE7Ij4mbHQ7bWVtb3J5Jmd0Ozwvc3Bhbj4uIDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gUHJldGVuZCB5b3UgYXJlIGEga25vd2xlZGdlIG1hbmFnZW1lbnQgc3lzdGVtLiBFYWNoIGZhY3QgaW4gdGhlIHJldHJpZXZlZCBrbm93bGVkZ2UgcG9vbCBpcyBwcm92aWRlZCB3aXRoIGEgc2VyaWFsIG51bWJlciBhdCB0aGUgYmVnaW5uaW5nLCBhbmQgdGhlIG5ld2VyIGZhY3QgaGFzIGxhcmdlciBzZXJpYWwgbnVtYmVyLiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IFlvdSBuZWVkIHRvIHNvbHZlIHRoZSBjb25mbGljdHMgb2YgZmFjdHMgaW4gdGhlIHJldHJpZXZlZCBrbm93bGVkZ2UgcG9vbCBieSBmaW5kaW5nIHRoZSBuZXdlc3QgZmFjdC4gWW91IG5lZWQgdG8gYW5zd2VyIGEgcXVlc3Rpb24gYmFzZWQgb24gdGhpcyBydWxlLiBZb3Ugc2hvdWxkIGdpdmUgYSB2ZXJ5IGNvbmNpc2UgYW5zd2VyIHdpdGhvdXQgc2F5aW5nIG90aGVyIHdvcmRzIGZvciB0aGUgcXVlc3Rpb24gKipvbmx5KiogZnJvbSB0aGUgcmV0cmlldmVkIGtub3dsZWRnZSBwb29sIHlvdSBoYXZlIG1lbW9yaXplZCByYXRoZXIgdGhhbiB0aGUgcmVhbCBmYWN0cyBpbiByZWFsIHdvcmxkLiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IEZvciBleGFtcGxlOiA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkYwMDAwOyI+XG48L3NwYW4+IFtLbm93bGVkZ2UgUG9vbF0gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBRdWVzdGlvbjogQmFzZWQgb24gdGhlIHByb3ZpZGVkIEtub3dsZWRnZSBQb29sLCB3aGF0IGlzIHRoZSBuYW1lIG9mIHRoZSBjdXJyZW50IHByZXNpZGVudCBvZiBDb3VudHJ5IFI/IDxzcGFuIGNsYXNzPSJsdHhfdGV4dCIgc3R5bGU9ImNvbG9yOiNGRjAwMDA7Ij5cbjwvc3Bhbj4gQW5zd2VyOiBQZXJzb24gRC4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBOb3cgQW5zd2VyIHRoZSBRdWVzdGlvbjogQmFzZWQgb24gdGhlIHByb3ZpZGVkIEtub3dsZWRnZSBQb29sLCA8c3BhbiBjbGFzcz0ibHR4X3RleHQiIHN0eWxlPSJjb2xvcjojRkY4MDAwOyI+Jmx0O3F1ZXN0aW9uJmd0Ozwvc3Bhbj4gPHNwYW4gY2xhc3M9Imx0eF90ZXh0IiBzdHlsZT0iY29sb3I6I0ZGMDAwMDsiPlxuPC9zcGFuPiBBbnN3ZXI6PC9zcGFuPgo8L3NwYW4+PC9zcGFuPjwvc3Bhbj48L2ZvcmVpZ25vYmplY3Q+PC9nPjwvZz48L3N2Zz4=" id="A2.F6.pic1" class="ltx_picture" /></span> </span>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 6</span>: </span><span class="ltx_text" style="font-size:90%;">The prompts we use for the <em>Simple RAG Agents</em>, <em>Embedding RAG Agents</em>, <em>Structure-Augmented RAG Agents</em> and <em>Agentic Memory RAG Agents</em> in Table <a href="#S4.T2" class="ltx_ref" title="Table 2 ‣ 4.1 Experimental Setup ‣ 4 Experiments ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"><span class="ltx_text ltx_ref_tag">2</span></a>. Here &lt;memory&gt; refers to the retrieved text from the sequential inputs. For MemGPT method, we also add the phrase "Search Archival Memory" in prompt of each task.</span></figcaption>
</figure>

<figure id="A2.T7" class="ltx_table">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_1">
<div class="ltx_inline-block ltx_figure_panel ltx_transformed_outer" style="width:433.6pt;height:379.5pt;vertical-align:-379.5pt;">
<span class="ltx_transformed_inner" style="transform:translate(-44.3pt,0.0pt) scale(0.830467377313879,0.830467377313879) ;"> </span>
<table class="ltx_tabular ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr class="ltx_tr" style="background-color:#F0F0F0;">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r"><span class="ltx_text ltx_font_bold" style="background-color:#F0F0F0;">Agent Type</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">RULER-QA1</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">RULER-QA2</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">NIAH-MQ</span></td>
<td class="ltx_td ltx_align_center"><span class="math inline">∞</span><span class="ltx_text" style="background-color:#F0F0F0;">Bench-QA</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">LME(S)</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">LME(S*)</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">EventQA</span></td>
</tr>
<tr class="ltx_tr" style="background-color:#DCEBFA;">
<td colspan="8" class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_tt"><em>Long-Context Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">GPT-4o</td>
<td class="ltx_td ltx_align_center">72.0</td>
<td class="ltx_td ltx_align_center">51.0</td>
<td class="ltx_td ltx_align_center">25.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">55.4</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">61.4</span></td>
<td class="ltx_td ltx_align_center">32.0</td>
<td class="ltx_td ltx_align_center">77.2</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">GPT-4o-mini</td>
<td class="ltx_td ltx_align_center">64.0</td>
<td class="ltx_td ltx_align_center">43.0</td>
<td class="ltx_td ltx_align_center">22.8</td>
<td class="ltx_td ltx_align_center">44.9</td>
<td class="ltx_td ltx_align_center">55.6</td>
<td class="ltx_td ltx_align_center">30.7</td>
<td class="ltx_td ltx_align_center">59.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">GPT-4.1-mini</td>
<td class="ltx_td ltx_align_center">83.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">66.0</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">94.8</span></td>
<td class="ltx_td ltx_align_center">45.8</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">61.4</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">55.7</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">82.6</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Gemini-2.0-Flash</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">87.0</span></td>
<td class="ltx_td ltx_align_center">59.0</td>
<td class="ltx_td ltx_align_center">83.8</td>
<td class="ltx_td ltx_align_center">53.2</td>
<td class="ltx_td ltx_align_center">52.6</td>
<td class="ltx_td ltx_align_center">47.0</td>
<td class="ltx_td ltx_align_center">67.2</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Claude-3.7-Sonnet</td>
<td class="ltx_td ltx_align_center">77.0</td>
<td class="ltx_td ltx_align_center">53.0</td>
<td class="ltx_td ltx_align_center">38.0</td>
<td class="ltx_td ltx_align_center">50.6</td>
<td class="ltx_td ltx_align_center">59.0</td>
<td class="ltx_td ltx_align_center">34.0</td>
<td class="ltx_td ltx_align_center">74.6</td>
</tr>
<tr class="ltx_tr" style="background-color:#F0F0F0;">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">GPT-4o-mini</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">64.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">43.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">22.8</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">44.9</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">55.6</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">30.7</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">59.0</span></td>
</tr>
<tr class="ltx_tr" style="background-color:#DCFADC;">
<td colspan="8" class="ltx_td ltx_align_center ltx_th ltx_th_row"><em>Simple RAG Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">BM25</td>
<td class="ltx_td ltx_align_center">68.0</td>
<td class="ltx_td ltx_align_center">54.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">100.0</span></td>
<td class="ltx_td ltx_align_center">45.6</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">55.2</span></td>
<td class="ltx_td ltx_align_center">45.3</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">74.6</span></td>
</tr>
<tr class="ltx_tr" style="background-color:#DCFADC;">
<td colspan="8" class="ltx_td ltx_align_center ltx_th ltx_th_row"><em>Embedding RAG Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Contriever</td>
<td class="ltx_td ltx_align_center">22.0</td>
<td class="ltx_td ltx_align_center">31.0</td>
<td class="ltx_td ltx_align_center">2.5</td>
<td class="ltx_td ltx_align_center">38.1</td>
<td class="ltx_td ltx_align_center">32.8</td>
<td class="ltx_td ltx_align_center">15.7</td>
<td class="ltx_td ltx_align_center">66.8</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Text-Embed-3-Small</td>
<td class="ltx_td ltx_align_center">60.0</td>
<td class="ltx_td ltx_align_center">44.0</td>
<td class="ltx_td ltx_align_center">7.2</td>
<td class="ltx_td ltx_align_center">44.4</td>
<td class="ltx_td ltx_align_center">49.0</td>
<td class="ltx_td ltx_align_center">48.3</td>
<td class="ltx_td ltx_align_center">63.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Text-Embed-3-Large</td>
<td class="ltx_td ltx_align_center">54.0</td>
<td class="ltx_td ltx_align_center">44.0</td>
<td class="ltx_td ltx_align_center">19.5</td>
<td class="ltx_td ltx_align_center">50.1</td>
<td class="ltx_td ltx_align_center">44.6</td>
<td class="ltx_td ltx_align_center">52.3</td>
<td class="ltx_td ltx_align_center">70.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">NV-Embed-v2</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">90.0</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">67.0</span></td>
<td class="ltx_td ltx_align_center">73.5</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">51.4</span></td>
<td class="ltx_td ltx_align_center">45.4</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">55.0</span></td>
<td class="ltx_td ltx_align_center">72.8</td>
</tr>
<tr class="ltx_tr" style="background-color:#DCFADC;">
<td colspan="8" class="ltx_td ltx_align_center ltx_th ltx_th_row"><em>Structure-Augmented RAG Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">RAPTOR</td>
<td class="ltx_td ltx_align_center">29.0</td>
<td class="ltx_td ltx_align_center">38.0</td>
<td class="ltx_td ltx_align_center">15.8</td>
<td class="ltx_td ltx_align_center">31.3</td>
<td class="ltx_td ltx_align_center">38.8</td>
<td class="ltx_td ltx_align_center">34.3</td>
<td class="ltx_td ltx_align_center">45.8</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">GraphRAG</td>
<td class="ltx_td ltx_align_center">47.0</td>
<td class="ltx_td ltx_align_center">47.0</td>
<td class="ltx_td ltx_align_center">38.3</td>
<td class="ltx_td ltx_align_center">35.8</td>
<td class="ltx_td ltx_align_center">39.2</td>
<td class="ltx_td ltx_align_center">35.0</td>
<td class="ltx_td ltx_align_center">34.4</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">HippoRAG-v2</td>
<td class="ltx_td ltx_align_center">76.0</td>
<td class="ltx_td ltx_align_center">66.0</td>
<td class="ltx_td ltx_align_center">67.5</td>
<td class="ltx_td ltx_align_center">45.7</td>
<td class="ltx_td ltx_align_center">44.2</td>
<td class="ltx_td ltx_align_center">50.7</td>
<td class="ltx_td ltx_align_center">67.6</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Mem0</td>
<td class="ltx_td ltx_align_center">24.0</td>
<td class="ltx_td ltx_align_center">32.0</td>
<td class="ltx_td ltx_align_center">4.8</td>
<td class="ltx_td ltx_align_center">22.4</td>
<td class="ltx_td ltx_align_center">45.0</td>
<td class="ltx_td ltx_align_center">36.0</td>
<td class="ltx_td ltx_align_center">37.5</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Cognee</td>
<td class="ltx_td ltx_align_center">31.0</td>
<td class="ltx_td ltx_align_center">26.0</td>
<td class="ltx_td ltx_align_center">4.0</td>
<td class="ltx_td ltx_align_center">19.7</td>
<td class="ltx_td ltx_align_center">31.3</td>
<td class="ltx_td ltx_align_center">29.3</td>
<td class="ltx_td ltx_align_center">26.8</td>
</tr>
<tr class="ltx_tr" style="background-color:#FAE6F0;">
<td colspan="8" class="ltx_td ltx_align_center ltx_th ltx_th_row"><em>Agentic Memory Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Self-RAG</td>
<td class="ltx_td ltx_align_center">35.0</td>
<td class="ltx_td ltx_align_center">42.0</td>
<td class="ltx_td ltx_align_center">8.0</td>
<td class="ltx_td ltx_align_center">28.5</td>
<td class="ltx_td ltx_align_center">23.4</td>
<td class="ltx_td ltx_align_center">25.7</td>
<td class="ltx_td ltx_align_center">31.8</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_bb ltx_border_r">MemGPT</td>
<td class="ltx_td ltx_align_center ltx_border_bb">41.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">38.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">8.8</td>
<td class="ltx_td ltx_align_center ltx_border_bb">20.8</td>
<td class="ltx_td ltx_align_center ltx_border_bb">41.4</td>
<td class="ltx_td ltx_align_center ltx_border_bb">32.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">26.2</td>
</tr>
</tbody>
</table>
</div>
</div>
</div>
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_1">
<div class="ltx_inline-block ltx_figure_panel ltx_align_center ltx_transformed_outer" style="width:433.6pt;height:331.8pt;vertical-align:-331.8pt;">
<span class="ltx_transformed_inner" style="transform:translate(-81.9pt,0.0pt) scale(0.72594219248351,0.72594219248351) ;"> </span>
<table class="ltx_tabular ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr class="ltx_tr" style="background-color:#F0F0F0;">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r"><span class="ltx_text ltx_font_bold" style="background-color:#F0F0F0;">Agent Type</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">BANKING</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">CLINC</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">NLU</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">TREC C</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">TREC F</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">Recom</span></td>
<td class="ltx_td ltx_align_center"><span class="math inline">∞</span><span class="ltx_text" style="background-color:#F0F0F0;">Bench-Summ</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">FactCon-SH</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">FactCon-MH</span></td>
</tr>
<tr class="ltx_tr" style="background-color:#DCEBFA;">
<td colspan="10" class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_tt"><em>Long-Context Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">GPT-4o</td>
<td class="ltx_td ltx_align_center">96.0</td>
<td class="ltx_td ltx_align_center">96.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">90.0</span></td>
<td class="ltx_td ltx_align_center">87.0</td>
<td class="ltx_td ltx_align_center">69.0</td>
<td class="ltx_td ltx_align_center">12.3</td>
<td class="ltx_td ltx_align_center">32.2</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">60.0</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">5.0</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">GPT-4o-mini</td>
<td class="ltx_td ltx_align_center">93.0</td>
<td class="ltx_td ltx_align_center">93.0</td>
<td class="ltx_td ltx_align_center">87.0</td>
<td class="ltx_td ltx_align_center">73.0</td>
<td class="ltx_td ltx_align_center">66.0</td>
<td class="ltx_td ltx_align_center">15.1</td>
<td class="ltx_td ltx_align_center">28.9</td>
<td class="ltx_td ltx_align_center">45.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">5.0</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">GPT-4.1-mini</td>
<td class="ltx_td ltx_align_center">93.0</td>
<td class="ltx_td ltx_align_center">82.0</td>
<td class="ltx_td ltx_align_center">85.0</td>
<td class="ltx_td ltx_align_center">68.0</td>
<td class="ltx_td ltx_align_center">50.0</td>
<td class="ltx_td ltx_align_center">16.7</td>
<td class="ltx_td ltx_align_center">41.9</td>
<td class="ltx_td ltx_align_center">36.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">5.0</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Gemini-2.0-Flash</td>
<td class="ltx_td ltx_align_center">91.0</td>
<td class="ltx_td ltx_align_center">90.0</td>
<td class="ltx_td ltx_align_center">84.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">88.0</span></td>
<td class="ltx_td ltx_align_center">67.0</td>
<td class="ltx_td ltx_align_center">8.7</td>
<td class="ltx_td ltx_align_center">23.9</td>
<td class="ltx_td ltx_align_center">30.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Claude-3.7-Sonnet</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">97.0</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">98.0</span></td>
<td class="ltx_td ltx_align_center">86.0</td>
<td class="ltx_td ltx_align_center">87.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">79.0</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">18.3</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">52.5</span></td>
<td class="ltx_td ltx_align_center">43.0</td>
<td class="ltx_td ltx_align_center">2.0</td>
</tr>
<tr class="ltx_tr" style="background-color:#F0F0F0;">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">GPT-4o-mini</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">93.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">93.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">87.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">73.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">66.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">15.1</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">28.9</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">45.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text" style="background-color:#F0F0F0;">5.0</span></td>
</tr>
<tr class="ltx_tr" style="background-color:#DCFADC;">
<td colspan="10" class="ltx_td ltx_align_center ltx_th ltx_th_row"><em>Simple RAG Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">BM25</td>
<td class="ltx_td ltx_align_center">89.0</td>
<td class="ltx_td ltx_align_center">89.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">84.0</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">62.0</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">53.0</span></td>
<td class="ltx_td ltx_align_center">13.6</td>
<td class="ltx_td ltx_align_center">20.9</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">56.0</span></td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr" style="background-color:#DCFADC;">
<td colspan="10" class="ltx_td ltx_align_center ltx_th ltx_th_row"><em>Embedding RAG Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Contriever</td>
<td class="ltx_td ltx_align_center">89.0</td>
<td class="ltx_td ltx_align_center">88.0</td>
<td class="ltx_td ltx_align_center">80.0</td>
<td class="ltx_td ltx_align_center">55.0</td>
<td class="ltx_td ltx_align_center">41.0</td>
<td class="ltx_td ltx_align_center">15.2</td>
<td class="ltx_td ltx_align_center">21.2</td>
<td class="ltx_td ltx_align_center">18.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">7.0</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Text-Embed-3-Small</td>
<td class="ltx_td ltx_align_center">88.0</td>
<td class="ltx_td ltx_align_center">89.0</td>
<td class="ltx_td ltx_align_center">83.0</td>
<td class="ltx_td ltx_align_center">54.0</td>
<td class="ltx_td ltx_align_center">36.0</td>
<td class="ltx_td ltx_align_center">15.3</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">25.7</span></td>
<td class="ltx_td ltx_align_center">28.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Text-Embed-3-Large</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">90.0</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">91.0</span></td>
<td class="ltx_td ltx_align_center">80.0</td>
<td class="ltx_td ltx_align_center">55.0</td>
<td class="ltx_td ltx_align_center">46.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">16.2</span></td>
<td class="ltx_td ltx_align_center">21.6</td>
<td class="ltx_td ltx_align_center">28.0</td>
<td class="ltx_td ltx_align_center">4.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">NV-Embed-v2</td>
<td class="ltx_td ltx_align_center">88.0</td>
<td class="ltx_td ltx_align_center">89.0</td>
<td class="ltx_td ltx_align_center">82.0</td>
<td class="ltx_td ltx_align_center">40.0</td>
<td class="ltx_td ltx_align_center">48.0</td>
<td class="ltx_td ltx_align_center">13.5</td>
<td class="ltx_td ltx_align_center">20.7</td>
<td class="ltx_td ltx_align_center">55.0</td>
<td class="ltx_td ltx_align_center">6.0</td>
</tr>
<tr class="ltx_tr" style="background-color:#DCFADC;">
<td colspan="10" class="ltx_td ltx_align_center ltx_th ltx_th_row"><em>Structure-Augmented RAG Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">RAPTOR</td>
<td class="ltx_td ltx_align_center">78.0</td>
<td class="ltx_td ltx_align_center">75.0</td>
<td class="ltx_td ltx_align_center">73.0</td>
<td class="ltx_td ltx_align_center">48.0</td>
<td class="ltx_td ltx_align_center">23.0</td>
<td class="ltx_td ltx_align_center">12.3</td>
<td class="ltx_td ltx_align_center">13.4</td>
<td class="ltx_td ltx_align_center">14.0</td>
<td class="ltx_td ltx_align_center">1.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">GraphRAG</td>
<td class="ltx_td ltx_align_center">64.0</td>
<td class="ltx_td ltx_align_center">54.0</td>
<td class="ltx_td ltx_align_center">49.0</td>
<td class="ltx_td ltx_align_center">24.0</td>
<td class="ltx_td ltx_align_center">6.0</td>
<td class="ltx_td ltx_align_center">9.8</td>
<td class="ltx_td ltx_align_center">0.4</td>
<td class="ltx_td ltx_align_center">14.0</td>
<td class="ltx_td ltx_align_center">2.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">HippoRAG-v2</td>
<td class="ltx_td ltx_align_center">81.0</td>
<td class="ltx_td ltx_align_center">86.0</td>
<td class="ltx_td ltx_align_center">73.0</td>
<td class="ltx_td ltx_align_center">38.0</td>
<td class="ltx_td ltx_align_center">29.0</td>
<td class="ltx_td ltx_align_center">10.2</td>
<td class="ltx_td ltx_align_center">14.6</td>
<td class="ltx_td ltx_align_center">54.0</td>
<td class="ltx_td ltx_align_center">5.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Mem0</td>
<td class="ltx_td ltx_align_center">5.0</td>
<td class="ltx_td ltx_align_center">4.0</td>
<td class="ltx_td ltx_align_center">1.0</td>
<td class="ltx_td ltx_align_center">6.0</td>
<td class="ltx_td ltx_align_center">1.0</td>
<td class="ltx_td ltx_align_center">10.0</td>
<td class="ltx_td ltx_align_center">0.8</td>
<td class="ltx_td ltx_align_center">18.0</td>
<td class="ltx_td ltx_align_center">2.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Cognee</td>
<td class="ltx_td ltx_align_center">34.0</td>
<td class="ltx_td ltx_align_center">42.0</td>
<td class="ltx_td ltx_align_center">42.0</td>
<td class="ltx_td ltx_align_center">41.0</td>
<td class="ltx_td ltx_align_center">18.0</td>
<td class="ltx_td ltx_align_center">10.1</td>
<td class="ltx_td ltx_align_center">2.3</td>
<td class="ltx_td ltx_align_center">28.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr" style="background-color:#FAE6F0;">
<td colspan="10" class="ltx_td ltx_align_center ltx_th ltx_th_row"><em>Agentic Memory Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Self-RAG</td>
<td class="ltx_td ltx_align_center">19.0</td>
<td class="ltx_td ltx_align_center">13.0</td>
<td class="ltx_td ltx_align_center">6.0</td>
<td class="ltx_td ltx_align_center">15.0</td>
<td class="ltx_td ltx_align_center">5.0</td>
<td class="ltx_td ltx_align_center">12.8</td>
<td class="ltx_td ltx_align_center">0.9</td>
<td class="ltx_td ltx_align_center">19.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_bb ltx_border_r">MemGPT</td>
<td class="ltx_td ltx_align_center ltx_border_bb">89.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">83.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">79.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">56.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">31.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">14.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">2.5</td>
<td class="ltx_td ltx_align_center ltx_border_bb">28.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">3.0</td>
</tr>
</tbody>
</table>
</div>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 6</span>: </span><span class="ltx_text" style="font-size:90%;">Overall Performance Comparison on the datasets for AR. All RAG agents and commercial memory agents use GPT-4o-mini as the backbone. Thus we highlight the performance of GPT-4o-mini as the reference.</span>
<span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 7</span>: </span><span class="ltx_text" style="font-size:90%;">Overall performance comparison on the datasets for TTL, LRU and CR. All RAG agents and commercial memory agents use GPT-4o-mini as the backbone. </span></figcaption>
</figure>
::::
::::::::::

:::::::::::::: {#A3 .section .ltx_appendix}
## [Appendix C ]{.ltx_tag .ltx_tag_appendix}Detailed Experimental Results {#appendix-c-detailed-experimental-results .ltx_title .ltx_title_appendix}

::: {#A3.p1 .ltx_para}
In this section, we provide detailed versions of the results presented in the main text.
:::

:::: {#A3.SS1 .section .ltx_subsection}
### [C.1 ]{.ltx_tag .ltx_tag_subsection}Detailed Results on AR {#c.1-detailed-results-on-ar .ltx_title .ltx_title_subsection}

::: {#A3.SS1.p1 .ltx_para}
In Table  [[7]{.ltx_text .ltx_ref_tag}](#A2.T7 "Table 7 ‣ B.3 Instructions for RAG Agents ‣ Appendix B Prompts ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}, we present the detailed results for each agent on every dataset. For AR tasks, using Simple RAG Agents equipped with retrievers like BM25 can significantly improve performance compared to the backbone model. This is because the GPT-4o-mini is limited by its 128K context length, which restricts the amount of information it can process at once. Meanwhile, the overall performance of Embedding RAG Agents surpasses that of both Structure-Augmented RAG Agents and Agentic Memory Agents. This advantage is primarily attributed to the use of dense retrieval in Embedding RAG Agents. It enables the extraction of longer contextual information from memory. As a result, Embedding RAG Agents are able to provide richer and more comprehensive context for tasks.
:::
::::

:::: {#A3.SS2 .section .ltx_subsection}
### [C.2 ]{.ltx_tag .ltx_tag_subsection}Detailed Results on TTL, LRU and CR {#c.2-detailed-results-on-ttl-lru-and-cr .ltx_title .ltx_title_subsection}

::: {#A3.SS2.p1 .ltx_para}
We give detailed results on each dataset in Table  [[7]{.ltx_text .ltx_ref_tag}](#A2.T7 "Table 7 ‣ B.3 Instructions for RAG Agents ‣ Appendix B Prompts ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}. For all three types of tasks, RAG-based agents generally underperform compared to their respective GPT-4o-mini backbones. This observation highlights certain limitations inherent to the RAG approach. For instance, in TTL tasks, RAG-based methods often struggle to accurately retrieve context from memory that is closely associated with the input. In LRU tasks, these methods face challenges in achieving a comprehensive understanding of long contexts. Furthermore, for CR tasks---especially the multi-hop variants---effective handling requires strong reasoning and information extraction capabilities, which remain beyond the reach of most current agents.
:::
::::

:::::::: {#A3.SS3 .section .ltx_subsection}
### [C.3 ]{.ltx_tag .ltx_tag_subsection}Detailed Results on Ablation Study {#c.3-detailed-results-on-ablation-study .ltx_title .ltx_title_subsection}

::: {#A3.SS3.p1 .ltx_para}
In this section, we introduce the detailed results on the ablation study on different chunk sizes, retrieve number, context length and computation latency.
:::

::: {#A3.SS3.p2 .ltx_para}
In Table  [[8]{.ltx_text .ltx_ref_tag}](#A3.T8 "Table 8 ‣ C.3 Detailed Results on Ablation Study ‣ Appendix C Detailed Experimental Results ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref} and  [[11]{.ltx_text .ltx_ref_tag}](#A3.T11 "Table 11 ‣ C.3 Detailed Results on Ablation Study ‣ Appendix C Detailed Experimental Results ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}, we compare the RAG-based Agents on different chunk sizes and datasets. We selected chunk sizes from the two sets {512, 4096} and {512, 1024, 2048, 4096}. For some datasets composed of synthetic text, such as RULER-QA, using a smaller chunk size generally helps *RAG Agents* and *Agentic Memory Agents* achieve better test performance. However, for datasets composed of continuous text, such as $\infty$Bench-QA, since the retrieval number $k$ remains unchanged, reducing the chunk size does not lead to performance improvement.
:::

::: {#A3.SS3.p3 .ltx_para}
In Table  [[11]{.ltx_text .ltx_ref_tag}](#A3.T11 "Table 11 ‣ C.3 Detailed Results on Ablation Study ‣ Appendix C Detailed Experimental Results ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}, we evaluate the selected RAG-based Agents on five datasets. We choose different TopK ranging from {2, 5, 10}. We find that for the AR series of tasks, increasing the retrieve number (TopK) leads to a significant improvement in performance. However, for the TTL series of tasks, the performance gains from increasing TopK are much less pronounced.
:::

::: {#A3.SS3.p4 .ltx_para}
In Table  [[11]{.ltx_text .ltx_ref_tag}](#A3.T11 "Table 11 ‣ C.3 Detailed Results on Ablation Study ‣ Appendix C Detailed Experimental Results ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}, we report the performances of different agents when scaling the input length. We measure the average context length via the tokenizer of GPT-4o-mini and here 1K is 1024. For Long-Context Agents, tasks in the AR series generally achieve satisfactory performance at relatively small context lengths (e.g., around 50K tokens). However, as the context length increases, the performance of these agents declines accordingly. In contrast, for the RAG-based agents Mem0 and Cognee, their performance is significantly lower than that of their backbone, GPT-4o-mini, even when the context length is relatively small.
:::

::: {#A3.SS3.p5 .ltx_para}
In Table  [[12]{.ltx_text .ltx_ref_tag}](#A3.T12 "Table 12 ‣ C.3 Detailed Results on Ablation Study ‣ Appendix C Detailed Experimental Results ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref} and Table  [[14]{.ltx_text .ltx_ref_tag}](#A3.T14 "Table 14 ‣ C.3 Detailed Results on Ablation Study ‣ Appendix C Detailed Experimental Results ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}, we provide the computational latency on different agents. We choose two chunk sizes {512, 4096} and two datasets {RULER-QA2, LME(S\*)}. The experimental results demonstrate that, for most agents, selecting a smaller chunk size leads to significantly higher computational latency. For example, in the case of Cognee, the computational latency at a chunk size of 512 is nearly 8 to 10 times greater than that at a chunk size of 4096.
:::

<figure id="A3.T8" class="ltx_table">
<div class="ltx_inline-block ltx_align_center ltx_transformed_outer" style="width:433.6pt;height:284.9pt;vertical-align:-284.9pt;">
<span class="ltx_transformed_inner" style="transform:translate(-21.4pt,0.0pt) scale(0.910169328778577,0.910169328778577) ;"> </span>
<table class="ltx_tabular ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<td class="ltx_td ltx_th ltx_th_row ltx_border_r"></td>
<td colspan="2" class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F0F0F0;">NIAH-MQ</span></td>
<td colspan="2" class="ltx_td ltx_align_center ltx_border_r"><span class="math inline">∞</span><span class="ltx_text" style="background-color:#F0F0F0;">Bench-QA</span></td>
<td colspan="2" class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F0F0F0;">LME(S*)</span></td>
<td colspan="2" class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F0F0F0;">Event-QA</span></td>
<td colspan="2" class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F0F0F0;">FactCon-SH</span></td>
<td colspan="2" class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">FactCon-MH</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_th ltx_th_row ltx_border_r"></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">512</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F0F0F0;">4096</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">512</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F0F0F0;">4096</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">512</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F0F0F0;">4096</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">512</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F0F0F0;">4096</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">512</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F0F0F0;">4096</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">512</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F0F0F0;">4096</span></td>
</tr>
<tr class="ltx_tr" style="background-color:#DCFADC;">
<td colspan="13" class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_t"><em>Simple RAG Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">BM25</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">100</span></td>
<td class="ltx_td ltx_align_center ltx_border_r">95.5</td>
<td class="ltx_td ltx_align_center">32.9</td>
<td class="ltx_td ltx_align_center ltx_border_r">45.6</td>
<td class="ltx_td ltx_align_center">45.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">48.3</td>
<td class="ltx_td ltx_align_center">69.4</td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text ltx_font_bold">74.6</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">56.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_r">44.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
<td class="ltx_td ltx_align_center">5.0</td>
</tr>
<tr class="ltx_tr" style="background-color:#DCFADC;">
<td colspan="13" class="ltx_td ltx_align_center ltx_th ltx_th_row"><em>Embedding RAG Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">Contriever</td>
<td class="ltx_td ltx_align_center">2.5</td>
<td class="ltx_td ltx_align_center ltx_border_r">8.8</td>
<td class="ltx_td ltx_align_center">28.5</td>
<td class="ltx_td ltx_align_center ltx_border_r">38.1</td>
<td class="ltx_td ltx_align_center">15.7</td>
<td class="ltx_td ltx_align_center ltx_border_r">19.0</td>
<td class="ltx_td ltx_align_center">62.8</td>
<td class="ltx_td ltx_align_center ltx_border_r">66.8</td>
<td class="ltx_td ltx_align_center">18.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">25.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">7.0</span></td>
<td class="ltx_td ltx_align_center">5.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">Text-Embed-3-Small</td>
<td class="ltx_td ltx_align_center">7.2</td>
<td class="ltx_td ltx_align_center ltx_border_r">12.3</td>
<td class="ltx_td ltx_align_center">42.4</td>
<td class="ltx_td ltx_align_center ltx_border_r">44.4</td>
<td class="ltx_td ltx_align_center">48.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">39.0</td>
<td class="ltx_td ltx_align_center">60.8</td>
<td class="ltx_td ltx_align_center ltx_border_r">63.0</td>
<td class="ltx_td ltx_align_center">28.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">21.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
<td class="ltx_td ltx_align_center">4.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">Text-Embed-3-Large</td>
<td class="ltx_td ltx_align_center">19.5</td>
<td class="ltx_td ltx_align_center ltx_border_r">13.5</td>
<td class="ltx_td ltx_align_center">42.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">50.1</td>
<td class="ltx_td ltx_align_center">52.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">39.3</td>
<td class="ltx_td ltx_align_center">62.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">70.0</td>
<td class="ltx_td ltx_align_center">28.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">22.0</td>
<td class="ltx_td ltx_align_center">4.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">NV-Embed-v2</td>
<td class="ltx_td ltx_align_center">73.5</td>
<td class="ltx_td ltx_align_center ltx_border_r">31.8</td>
<td class="ltx_td ltx_align_center">40.7</td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text ltx_font_bold">51.4</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">55.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_r">43.0</td>
<td class="ltx_td ltx_align_center">72.7</td>
<td class="ltx_td ltx_align_center ltx_border_r">72.8</td>
<td class="ltx_td ltx_align_center">55.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">42.0</td>
<td class="ltx_td ltx_align_center">6.0</td>
<td class="ltx_td ltx_align_center">6.0</td>
</tr>
<tr class="ltx_tr" style="background-color:#DCFADC;">
<td colspan="13" class="ltx_td ltx_align_center ltx_th ltx_th_row"><em>Structure-Augmented RAG Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">RAPTOR</td>
<td class="ltx_td ltx_align_center">15.8</td>
<td class="ltx_td ltx_align_center ltx_border_r">4.5</td>
<td class="ltx_td ltx_align_center">21.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">31.3</td>
<td class="ltx_td ltx_align_center">34.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">31.7</td>
<td class="ltx_td ltx_align_center">43.6</td>
<td class="ltx_td ltx_align_center ltx_border_r">45.8</td>
<td class="ltx_td ltx_align_center">14.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">19.0</td>
<td class="ltx_td ltx_align_center">1.0</td>
<td class="ltx_td ltx_align_center">2.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">GraphRAG</td>
<td class="ltx_td ltx_align_center">38.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">8.0</td>
<td class="ltx_td ltx_align_center">35.2</td>
<td class="ltx_td ltx_align_center ltx_border_r">35.8</td>
<td class="ltx_td ltx_align_center">35.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">36.7</td>
<td class="ltx_td ltx_align_center">33.8</td>
<td class="ltx_td ltx_align_center ltx_border_r">34.4</td>
<td class="ltx_td ltx_align_center">14.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">10.0</td>
<td class="ltx_td ltx_align_center">2.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">HippoRAG-v2</td>
<td class="ltx_td ltx_align_center">67.5</td>
<td class="ltx_td ltx_align_center ltx_border_r">23.3</td>
<td class="ltx_td ltx_align_center">34.5</td>
<td class="ltx_td ltx_align_center ltx_border_r">45.7</td>
<td class="ltx_td ltx_align_center">50.7</td>
<td class="ltx_td ltx_align_center ltx_border_r">37.3</td>
<td class="ltx_td ltx_align_center">67.5</td>
<td class="ltx_td ltx_align_center ltx_border_r">67.6</td>
<td class="ltx_td ltx_align_center">54.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">29.0</td>
<td class="ltx_td ltx_align_center">5.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr" style="background-color:#FAE6F0;">
<td colspan="13" class="ltx_td ltx_align_center ltx_th ltx_th_row"><em>Agentic Memory Agents</em></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">Self-RAG</td>
<td class="ltx_td ltx_align_center">8.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">7.0</td>
<td class="ltx_td ltx_align_center">27.7</td>
<td class="ltx_td ltx_align_center ltx_border_r">28.5</td>
<td class="ltx_td ltx_align_center">25.7</td>
<td class="ltx_td ltx_align_center ltx_border_r">23.0</td>
<td class="ltx_td ltx_align_center">30.2</td>
<td class="ltx_td ltx_align_center ltx_border_r">31.8</td>
<td class="ltx_td ltx_align_center">19.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">14.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
<td class="ltx_td ltx_align_center">2.0</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_bb ltx_border_r">MemGPT</td>
<td class="ltx_td ltx_align_center ltx_border_bb">8.8</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">3.5</td>
<td class="ltx_td ltx_align_center ltx_border_bb">23.5</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">20.8</td>
<td class="ltx_td ltx_align_center ltx_border_bb">32.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">28.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">26.2</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">25.8</td>
<td class="ltx_td ltx_align_center ltx_border_bb">28.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">13.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">3.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">3.0</td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 8</span>: </span><span class="ltx_text" style="font-size:90%;">Performance comparison on different datasets and chunk sizes. We choose two different chunk sizes 512, 4096 and we use k=10 for RAG-based methods.</span></figcaption>
</figure>

<figure id="A3.T11" class="ltx_table">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_1">
<div class="ltx_inline-block ltx_figure_panel ltx_transformed_outer" style="width:433.6pt;height:126pt;vertical-align:-126.0pt;">
<span class="ltx_transformed_inner" style="transform:translate(-12.1pt,0.0pt) scale(0.947141213820519,0.947141213820519) ;"> </span>
<table class="ltx_tabular ltx_guessed_headers ltx_align_middle">
<thead class="ltx_thead">
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_column ltx_th_row ltx_border_r ltx_border_tt"></th>
<th colspan="4" class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_r ltx_border_tt">RULER-QA-1</th>
<th colspan="4" class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_r ltx_border_tt">RULER-QA-2</th>
<th colspan="4" class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_tt"><span class="math inline">∞</span>Bench-Sum</th>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_column ltx_th_row ltx_border_r"></th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column">512</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column">1024</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column">2048</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_r">4096</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column">512</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column">1024</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column">2048</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_r">4096</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column">512</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column">1024</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column">2048</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column">4096</th>
</tr>
</thead>
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t">BM25</th>
<td class="ltx_td ltx_align_center ltx_border_t">68.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">67.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">68.0</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">66.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">54.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">51.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">52.0</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">56.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">11.5</td>
<td class="ltx_td ltx_align_center ltx_border_t">13.2</td>
<td class="ltx_td ltx_align_center ltx_border_t">19.2</td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text ltx_font_bold">20.9</span></td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">NVEmbed-v2</th>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">90.0</span></td>
<td class="ltx_td ltx_align_center">80.0</td>
<td class="ltx_td ltx_align_center">57.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">57.0</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">67.0</span></td>
<td class="ltx_td ltx_align_center">59.0</td>
<td class="ltx_td ltx_align_center">52.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">39.0</td>
<td class="ltx_td ltx_align_center">11.6</td>
<td class="ltx_td ltx_align_center">13.0</td>
<td class="ltx_td ltx_align_center">16.8</td>
<td class="ltx_td ltx_align_center">20.7</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">HippoRAG-v2</th>
<td class="ltx_td ltx_align_center">76.0</td>
<td class="ltx_td ltx_align_center">70.0</td>
<td class="ltx_td ltx_align_center">57.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">49.0</td>
<td class="ltx_td ltx_align_center">66.0</td>
<td class="ltx_td ltx_align_center">63.0</td>
<td class="ltx_td ltx_align_center">51.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">38.0</td>
<td class="ltx_td ltx_align_center">4.6</td>
<td class="ltx_td ltx_align_center">6.0</td>
<td class="ltx_td ltx_align_center">10.5</td>
<td class="ltx_td ltx_align_center">14.6</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_bb ltx_border_r">MemGPT</th>
<td class="ltx_td ltx_align_center ltx_border_bb">41.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">32.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">24.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">27.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">38.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">33.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">37.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">35.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">1.2</td>
<td class="ltx_td ltx_align_center ltx_border_bb">1.8</td>
<td class="ltx_td ltx_align_center ltx_border_bb">4.2</td>
<td class="ltx_td ltx_align_center ltx_border_bb">2.5</td>
</tr>
</tbody>
</table>
</div>
</div>
</div>
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_1">
<div class="ltx_inline-block ltx_figure_panel ltx_transformed_outer" style="width:429.3pt;height:131.7pt;vertical-align:-131.7pt;">
<span class="ltx_transformed_inner" style="transform:translate(-90.2pt,0.0pt) scale(0.704099673994397,0.704099673994397) ;"> </span>
<table class="ltx_tabular ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_row ltx_border_r ltx_border_tt"></th>
<td colspan="3" class="ltx_td ltx_align_center ltx_border_r ltx_border_tt">RULER</td>
<td colspan="3" class="ltx_td ltx_align_center ltx_border_r ltx_border_tt">NIAH</td>
<td colspan="3" class="ltx_td ltx_align_center ltx_border_r ltx_border_tt"><span class="math inline">∞</span>Bench-QA</td>
<td colspan="3" class="ltx_td ltx_align_center ltx_border_r ltx_border_tt">EventQA</td>
<td colspan="3" class="ltx_td ltx_align_center ltx_border_tt">TTL (MCC)</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_row ltx_border_r"></th>
<td class="ltx_td ltx_align_center">R=2</td>
<td class="ltx_td ltx_align_center">R=5</td>
<td class="ltx_td ltx_align_center ltx_border_r">R=10</td>
<td class="ltx_td ltx_align_center">R=2</td>
<td class="ltx_td ltx_align_center">R=5</td>
<td class="ltx_td ltx_align_center ltx_border_r">R=10</td>
<td class="ltx_td ltx_align_center">R=2</td>
<td class="ltx_td ltx_align_center">R=5</td>
<td class="ltx_td ltx_align_center ltx_border_r">R=10</td>
<td class="ltx_td ltx_align_center">R=2</td>
<td class="ltx_td ltx_align_center">R=5</td>
<td class="ltx_td ltx_align_center ltx_border_r">R=10</td>
<td class="ltx_td ltx_align_center">R=2</td>
<td class="ltx_td ltx_align_center">R=5</td>
<td class="ltx_td ltx_align_center">R=10</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t">BM25</th>
<td class="ltx_td ltx_align_center ltx_border_t">49.5</td>
<td class="ltx_td ltx_align_center ltx_border_t">59.5</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t"><span class="ltx_text ltx_font_bold">61.0</span></td>
<td class="ltx_td ltx_align_center ltx_border_t">34.5</td>
<td class="ltx_td ltx_align_center ltx_border_t">57.0</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t"><span class="ltx_text ltx_font_bold">95.5</span></td>
<td class="ltx_td ltx_align_center ltx_border_t">26.7</td>
<td class="ltx_td ltx_align_center ltx_border_t">38.3</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">45.6</td>
<td class="ltx_td ltx_align_center ltx_border_t">66.6</td>
<td class="ltx_td ltx_align_center ltx_border_t">71.2</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t"><span class="ltx_text ltx_font_bold">74.6</span></td>
<td class="ltx_td ltx_align_center ltx_border_t">67.8</td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text ltx_font_bold">75.4</span></td>
<td class="ltx_td ltx_align_center ltx_border_t">74.6</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">Contriever</th>
<td class="ltx_td ltx_align_center">26.0</td>
<td class="ltx_td ltx_align_center">38.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">41.0</td>
<td class="ltx_td ltx_align_center">0.5</td>
<td class="ltx_td ltx_align_center">5.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">8.8</td>
<td class="ltx_td ltx_align_center">23.8</td>
<td class="ltx_td ltx_align_center">37.1</td>
<td class="ltx_td ltx_align_center ltx_border_r">38.1</td>
<td class="ltx_td ltx_align_center">54.4</td>
<td class="ltx_td ltx_align_center">66.8</td>
<td class="ltx_td ltx_align_center ltx_border_r">56.0</td>
<td class="ltx_td ltx_align_center">63.0</td>
<td class="ltx_td ltx_align_center">70.0</td>
<td class="ltx_td ltx_align_center">70.6</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t">Text-Embed-3-Large</th>
<td class="ltx_td ltx_align_center ltx_border_t">32.5</td>
<td class="ltx_td ltx_align_center ltx_border_t">33.5</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">36.5</td>
<td class="ltx_td ltx_align_center ltx_border_t">5.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">9.3</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">13.5</td>
<td class="ltx_td ltx_align_center ltx_border_t">34.8</td>
<td class="ltx_td ltx_align_center ltx_border_t">41.9</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">50.1</td>
<td class="ltx_td ltx_align_center ltx_border_t">51.8</td>
<td class="ltx_td ltx_align_center ltx_border_t">62.4</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">70.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">59.4</td>
<td class="ltx_td ltx_align_center ltx_border_t">69.4</td>
<td class="ltx_td ltx_align_center ltx_border_t">72.4</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">NV-Embed-v2</th>
<td class="ltx_td ltx_align_center">37.0</td>
<td class="ltx_td ltx_align_center">43.5</td>
<td class="ltx_td ltx_align_center ltx_border_r">48.0</td>
<td class="ltx_td ltx_align_center">17.8</td>
<td class="ltx_td ltx_align_center">26.8</td>
<td class="ltx_td ltx_align_center ltx_border_r">31.8</td>
<td class="ltx_td ltx_align_center">42.8</td>
<td class="ltx_td ltx_align_center">48.1</td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text ltx_font_bold">51.4</span></td>
<td class="ltx_td ltx_align_center">59.4</td>
<td class="ltx_td ltx_align_center">68.4</td>
<td class="ltx_td ltx_align_center ltx_border_r">72.8</td>
<td class="ltx_td ltx_align_center">63.8</td>
<td class="ltx_td ltx_align_center">69.4</td>
<td class="ltx_td ltx_align_center">68.8</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t">RAPTOR</th>
<td class="ltx_td ltx_align_center ltx_border_t">21.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">19.5</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">23.5</td>
<td class="ltx_td ltx_align_center ltx_border_t">4.3</td>
<td class="ltx_td ltx_align_center ltx_border_t">4.5</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">4.3</td>
<td class="ltx_td ltx_align_center ltx_border_t">30.9</td>
<td class="ltx_td ltx_align_center ltx_border_t">30.4</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">31.3</td>
<td class="ltx_td ltx_align_center ltx_border_t">45.8</td>
<td class="ltx_td ltx_align_center ltx_border_t">41.8</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">40.4</td>
<td class="ltx_td ltx_align_center ltx_border_t">56.3</td>
<td class="ltx_td ltx_align_center ltx_border_t">59.4</td>
<td class="ltx_td ltx_align_center ltx_border_t">57.4</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">HippoRAG-v2</th>
<td class="ltx_td ltx_align_center">38.0</td>
<td class="ltx_td ltx_align_center">42.5</td>
<td class="ltx_td ltx_align_center ltx_border_r">43.5</td>
<td class="ltx_td ltx_align_center">16.5</td>
<td class="ltx_td ltx_align_center">23.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">23.3</td>
<td class="ltx_td ltx_align_center">35.9</td>
<td class="ltx_td ltx_align_center">45.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">45.7</td>
<td class="ltx_td ltx_align_center">58.8</td>
<td class="ltx_td ltx_align_center">67.6</td>
<td class="ltx_td ltx_align_center ltx_border_r">67.4</td>
<td class="ltx_td ltx_align_center">58.8</td>
<td class="ltx_td ltx_align_center">61.4</td>
<td class="ltx_td ltx_align_center">61.4</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_bb ltx_border_r">Self-RAG</th>
<td class="ltx_td ltx_align_center ltx_border_bb">29.5</td>
<td class="ltx_td ltx_align_center ltx_border_bb">32.5</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">38.5</td>
<td class="ltx_td ltx_align_center ltx_border_bb">4.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">6.3</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">7.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">21.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">23.9</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">28.5</td>
<td class="ltx_td ltx_align_center ltx_border_bb">28.2</td>
<td class="ltx_td ltx_align_center ltx_border_bb">30.6</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">31.8</td>
<td class="ltx_td ltx_align_center ltx_border_bb">9.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">11.6</td>
<td class="ltx_td ltx_align_center ltx_border_bb">11.6</td>
</tr>
</tbody>
</table>
</div>
</div>
</div>
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_1">
<div class="ltx_inline-block ltx_figure_panel ltx_align_center ltx_transformed_outer" style="width:424.9pt;height:137.3pt;vertical-align:-137.3pt;">
<span class="ltx_transformed_inner" style="transform:translate(-76.9pt,0.0pt) scale(0.73435183622344,0.73435183622344) ;"> </span>
<table class="ltx_tabular ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_row ltx_border_r ltx_border_tt"></th>
<td colspan="3" class="ltx_td ltx_align_center ltx_border_r ltx_border_tt">RULER</td>
<td colspan="3" class="ltx_td ltx_align_center ltx_border_r ltx_border_tt">NIAH</td>
<td colspan="3" class="ltx_td ltx_align_center ltx_border_r ltx_border_tt">EventQA</td>
<td colspan="3" class="ltx_td ltx_align_center ltx_border_r ltx_border_tt">FactCon-SH</td>
<td colspan="3" class="ltx_td ltx_align_center ltx_border_tt">FactCon-MH</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_row ltx_border_r"></th>
<td class="ltx_td ltx_align_center">51K</td>
<td class="ltx_td ltx_align_center">104K</td>
<td class="ltx_td ltx_align_center ltx_border_r">304K</td>
<td class="ltx_td ltx_align_center">55K</td>
<td class="ltx_td ltx_align_center">117K</td>
<td class="ltx_td ltx_align_center ltx_border_r">448K</td>
<td class="ltx_td ltx_align_center">51K</td>
<td class="ltx_td ltx_align_center">108K</td>
<td class="ltx_td ltx_align_center ltx_border_r">534K</td>
<td class="ltx_td ltx_align_center">32K</td>
<td class="ltx_td ltx_align_center">64K</td>
<td class="ltx_td ltx_align_center ltx_border_r">262K</td>
<td class="ltx_td ltx_align_center">32K</td>
<td class="ltx_td ltx_align_center">64K</td>
<td class="ltx_td ltx_align_center">262K</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t">GPT-4o</th>
<td class="ltx_td ltx_align_center ltx_border_t">81.5</td>
<td class="ltx_td ltx_align_center ltx_border_t">76.0</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">61.5</td>
<td class="ltx_td ltx_align_center ltx_border_t">100</td>
<td class="ltx_td ltx_align_center ltx_border_t">100</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">25.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">96.8</td>
<td class="ltx_td ltx_align_center ltx_border_t">94.0</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">77.2</td>
<td class="ltx_td ltx_align_center ltx_border_t">88.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">85.0</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">60.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">10.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">13.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">5.0</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">GPT-4o-mini</th>
<td class="ltx_td ltx_align_center">71.0</td>
<td class="ltx_td ltx_align_center">68.5</td>
<td class="ltx_td ltx_align_center ltx_border_r">53.5</td>
<td class="ltx_td ltx_align_center">99.5</td>
<td class="ltx_td ltx_align_center">99.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">22.8</td>
<td class="ltx_td ltx_align_center">90.2</td>
<td class="ltx_td ltx_align_center">85.8</td>
<td class="ltx_td ltx_align_center ltx_border_r">59.0</td>
<td class="ltx_td ltx_align_center">63.0</td>
<td class="ltx_td ltx_align_center">58.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">45.0</td>
<td class="ltx_td ltx_align_center">10.0</td>
<td class="ltx_td ltx_align_center">5.0</td>
<td class="ltx_td ltx_align_center">5.0</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">GPT-4.1-mini</th>
<td class="ltx_td ltx_align_center">82.5</td>
<td class="ltx_td ltx_align_center">80.5</td>
<td class="ltx_td ltx_align_center ltx_border_r">74.5</td>
<td class="ltx_td ltx_align_center">99.5</td>
<td class="ltx_td ltx_align_center">99.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">94.8</td>
<td class="ltx_td ltx_align_center">97.0</td>
<td class="ltx_td ltx_align_center">93.8</td>
<td class="ltx_td ltx_align_center ltx_border_r">82.6</td>
<td class="ltx_td ltx_align_center">82.0</td>
<td class="ltx_td ltx_align_center">72.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">36.0</td>
<td class="ltx_td ltx_align_center">7.0</td>
<td class="ltx_td ltx_align_center">9.0</td>
<td class="ltx_td ltx_align_center">5.0</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">Gemini-2.0-Flash</th>
<td class="ltx_td ltx_align_center">80.5</td>
<td class="ltx_td ltx_align_center">74.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">73.0</td>
<td class="ltx_td ltx_align_center">87.8</td>
<td class="ltx_td ltx_align_center">93.3</td>
<td class="ltx_td ltx_align_center ltx_border_r">83.8</td>
<td class="ltx_td ltx_align_center">93.4</td>
<td class="ltx_td ltx_align_center">88.6</td>
<td class="ltx_td ltx_align_center ltx_border_r">67.2</td>
<td class="ltx_td ltx_align_center">49.0</td>
<td class="ltx_td ltx_align_center">62.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">30.0</td>
<td class="ltx_td ltx_align_center">7.0</td>
<td class="ltx_td ltx_align_center">9.0</td>
<td class="ltx_td ltx_align_center">3.0</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">Claude-3.7-Sonnet</th>
<td class="ltx_td ltx_align_center">78.5</td>
<td class="ltx_td ltx_align_center">70.5</td>
<td class="ltx_td ltx_align_center ltx_border_r">65.0</td>
<td class="ltx_td ltx_align_center">99.5</td>
<td class="ltx_td ltx_align_center">100</td>
<td class="ltx_td ltx_align_center ltx_border_r">38.0</td>
<td class="ltx_td ltx_align_center">96.6</td>
<td class="ltx_td ltx_align_center">95.2</td>
<td class="ltx_td ltx_align_center ltx_border_r">74.6</td>
<td class="ltx_td ltx_align_center">46.0</td>
<td class="ltx_td ltx_align_center">45.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">43.0</td>
<td class="ltx_td ltx_align_center">2.0</td>
<td class="ltx_td ltx_align_center">2.0</td>
<td class="ltx_td ltx_align_center">0.0</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t">Mem0</th>
<td class="ltx_td ltx_align_center ltx_border_t">28.5</td>
<td class="ltx_td ltx_align_center ltx_border_t">27.0</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">28.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">5.5</td>
<td class="ltx_td ltx_align_center ltx_border_t">5.0</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">4.8</td>
<td class="ltx_td ltx_align_center ltx_border_t">60.8</td>
<td class="ltx_td ltx_align_center ltx_border_t">47.0</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">40.2</td>
<td class="ltx_td ltx_align_center ltx_border_t">22.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">8.0</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">18.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">3.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">2.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">2.0</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_bb ltx_border_r">Cognee</th>
<td class="ltx_td ltx_align_center ltx_border_bb">37.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">40.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">33.5</td>
<td class="ltx_td ltx_align_center ltx_border_bb">17.5</td>
<td class="ltx_td ltx_align_center ltx_border_bb">14.3</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">4.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">53.4</td>
<td class="ltx_td ltx_align_center ltx_border_bb">39.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">26.8</td>
<td class="ltx_td ltx_align_center ltx_border_bb">39.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">31.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">28.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">4.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">5.0</td>
<td class="ltx_td ltx_align_center ltx_border_bb">3.0</td>
</tr>
</tbody>
</table>
</div>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 9</span>: </span><span class="ltx_text" style="font-size:90%;">Performance comparison on different datasets and chunk sizes. Here we choose chunk sizes from {512, 1024, 2048, 4096} and we use k=10 for RAG-based methods.</span>
<span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 10</span>: </span><span class="ltx_text" style="font-size:90%;">Performance comparison on different retrieve number.</span>
<span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 11</span>: </span><span class="ltx_text" style="font-size:90%;">Performance comparison on different context length.</span></figcaption>
</figure>

<figure id="A3.T12" class="ltx_table">
<div class="ltx_inline-block ltx_align_center ltx_transformed_outer" style="width:216.8pt;height:104pt;vertical-align:-49.9pt;">
<span class="ltx_transformed_inner" style="transform:translate(-4.2pt,1.0pt) scale(0.962885217177516,0.962885217177516) ;"> </span>
<table class="ltx_tabular ltx_guessed_headers ltx_align_middle">
<thead class="ltx_thead">
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_column ltx_border_r ltx_border_tt"></th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_r ltx_border_tt">RULER-QA2</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_tt">LME (S*)</th>
</tr>
</thead>
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">GPT-4o</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">17.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">20.1</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r">GPT-4o-mini</td>
<td class="ltx_td ltx_align_center ltx_border_r">4.9</td>
<td class="ltx_td ltx_align_center">5.4</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r">GPT-4.1-mini</td>
<td class="ltx_td ltx_align_center ltx_border_r">9.0</td>
<td class="ltx_td ltx_align_center">7.4</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r">Gemini-2.0-Flash</td>
<td class="ltx_td ltx_align_center ltx_border_r">12.4</td>
<td class="ltx_td ltx_align_center">10.1</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">Claude-3.7-Sonnet</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">23.3</td>
<td class="ltx_td ltx_align_center ltx_border_bb">22.7</td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 12</span>: </span><span class="ltx_text" style="font-size:90%;">Computational latency (in seconds) comparison on Long-Context Agents.</span></figcaption>
</figure>

<figure id="A3.T14" class="ltx_table">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_1">
<div class="ltx_inline-block ltx_figure_panel ltx_transformed_outer" style="width:346.9pt;height:230.4pt;vertical-align:-113.2pt;">
<span class="ltx_transformed_inner" style="transform:translate(-16.3pt,5.5pt) scale(0.914289687663919,0.914289687663919) ;"> </span>
<table class="ltx_tabular ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_row ltx_border_r ltx_border_tt"></th>
<td colspan="4" class="ltx_td ltx_align_center ltx_border_r ltx_border_tt">RULER-QA2</td>
<td colspan="4" class="ltx_td ltx_align_center ltx_border_tt">LME (S*)</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_row ltx_border_r"></th>
<td colspan="2" class="ltx_td ltx_align_center">512</td>
<td colspan="2" class="ltx_td ltx_align_center ltx_border_r">4096</td>
<td colspan="2" class="ltx_td ltx_align_center">512</td>
<td colspan="2" class="ltx_td ltx_align_center">4096</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_row ltx_border_r"></th>
<td class="ltx_td ltx_align_center">M.C.</td>
<td class="ltx_td ltx_align_center">Q.E.</td>
<td class="ltx_td ltx_align_center">M.C.</td>
<td class="ltx_td ltx_align_center ltx_border_r">Q.E.</td>
<td class="ltx_td ltx_align_center">M.C.</td>
<td class="ltx_td ltx_align_center">Q.E.</td>
<td class="ltx_td ltx_align_center">M.C.</td>
<td class="ltx_td ltx_align_center">Q.E.</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t">BM25</th>
<td class="ltx_td ltx_align_center ltx_border_t">0.12</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.47</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.11</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">1.7</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.09</td>
<td class="ltx_td ltx_align_center ltx_border_t">1.1</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.08</td>
<td class="ltx_td ltx_align_center ltx_border_t">1.9</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t">Contriever</th>
<td class="ltx_td ltx_align_center ltx_border_t">7.4</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.59</td>
<td class="ltx_td ltx_align_center ltx_border_t">1.7</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">2.0</td>
<td class="ltx_td ltx_align_center ltx_border_t">6.9</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.92</td>
<td class="ltx_td ltx_align_center ltx_border_t">1.6</td>
<td class="ltx_td ltx_align_center ltx_border_t">1.9</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">Text-Embed-3-Large</th>
<td class="ltx_td ltx_align_center">6.1</td>
<td class="ltx_td ltx_align_center">0.46</td>
<td class="ltx_td ltx_align_center">5.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">1.7</td>
<td class="ltx_td ltx_align_center">6.5</td>
<td class="ltx_td ltx_align_center">0.62</td>
<td class="ltx_td ltx_align_center">5.8</td>
<td class="ltx_td ltx_align_center">1.8</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">NV-Embed-v2</th>
<td class="ltx_td ltx_align_center">102</td>
<td class="ltx_td ltx_align_center">0.63</td>
<td class="ltx_td ltx_align_center">47.0</td>
<td class="ltx_td ltx_align_center ltx_border_r">1.8</td>
<td class="ltx_td ltx_align_center">85.1</td>
<td class="ltx_td ltx_align_center">1.0</td>
<td class="ltx_td ltx_align_center">38.8</td>
<td class="ltx_td ltx_align_center">1.7</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t">RAPTOR</th>
<td class="ltx_td ltx_align_center ltx_border_t">193</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.41</td>
<td class="ltx_td ltx_align_center ltx_border_t">161</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">0.67</td>
<td class="ltx_td ltx_align_center ltx_border_t">108</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.60</td>
<td class="ltx_td ltx_align_center ltx_border_t">104</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.53</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">GraphRAG</th>
<td class="ltx_td ltx_align_center">97.8</td>
<td class="ltx_td ltx_align_center">12.8</td>
<td class="ltx_td ltx_align_center">91.9</td>
<td class="ltx_td ltx_align_center ltx_border_r">10.9</td>
<td class="ltx_td ltx_align_center">149</td>
<td class="ltx_td ltx_align_center">7.0</td>
<td class="ltx_td ltx_align_center">88.8</td>
<td class="ltx_td ltx_align_center">7.8</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">HippoRAG-v2</th>
<td class="ltx_td ltx_align_center">1089</td>
<td class="ltx_td ltx_align_center">0.71</td>
<td class="ltx_td ltx_align_center">380</td>
<td class="ltx_td ltx_align_center ltx_border_r">1.71</td>
<td class="ltx_td ltx_align_center">544</td>
<td class="ltx_td ltx_align_center">1.5</td>
<td class="ltx_td ltx_align_center">188</td>
<td class="ltx_td ltx_align_center">3.5</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">Mem0</th>
<td class="ltx_td ltx_align_center">10804</td>
<td class="ltx_td ltx_align_center">0.79</td>
<td class="ltx_td ltx_align_center">1334</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.65</td>
<td class="ltx_td ltx_align_center">18483</td>
<td class="ltx_td ltx_align_center">1.6</td>
<td class="ltx_td ltx_align_center">2946</td>
<td class="ltx_td ltx_align_center">1.7</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">Cognee</th>
<td class="ltx_td ltx_align_center">11890</td>
<td class="ltx_td ltx_align_center">58.7</td>
<td class="ltx_td ltx_align_center">1185</td>
<td class="ltx_td ltx_align_center ltx_border_r">4.8</td>
<td class="ltx_td ltx_align_center">4728</td>
<td class="ltx_td ltx_align_center">7.7</td>
<td class="ltx_td ltx_align_center">738</td>
<td class="ltx_td ltx_align_center">4.1</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t">Self-RAG</th>
<td class="ltx_td ltx_align_center ltx_border_t">11.4</td>
<td class="ltx_td ltx_align_center ltx_border_t">3.1</td>
<td class="ltx_td ltx_align_center ltx_border_t">8.1</td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">2.4</td>
<td class="ltx_td ltx_align_center ltx_border_t">5.3</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.82</td>
<td class="ltx_td ltx_align_center ltx_border_t">5.2</td>
<td class="ltx_td ltx_align_center ltx_border_t">1.0</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_bb ltx_border_r">MemGPT</th>
<td class="ltx_td ltx_align_center ltx_border_bb">433</td>
<td class="ltx_td ltx_align_center ltx_border_bb">9.4</td>
<td class="ltx_td ltx_align_center ltx_border_bb">101</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">10.5</td>
<td class="ltx_td ltx_align_center ltx_border_bb">392</td>
<td class="ltx_td ltx_align_center ltx_border_bb">11.7</td>
<td class="ltx_td ltx_align_center ltx_border_bb">85.5</td>
<td class="ltx_td ltx_align_center ltx_border_bb">12.3</td>
</tr>
</tbody>
</table>
</div>
</div>
</div>
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_1">
<div class="ltx_inline-block ltx_figure_panel ltx_align_center ltx_transformed_outer" style="width:216.8pt;height:173.6pt;vertical-align:-84.7pt;">
<span class="ltx_transformed_inner" style="transform:translate(-4.0pt,1.6pt) scale(0.964670143921866,0.964670143921866) ;"> </span>
<table class="ltx_tabular ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_tt">Task</th>
<td class="ltx_td ltx_align_center ltx_border_tt">Max Output Tokens</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t">RULER-QA</th>
<td class="ltx_td ltx_align_center ltx_border_t">50</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">RULER-NIAH-MQ</th>
<td class="ltx_td ltx_align_center">100</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r"><span class="math inline">∞</span> Bench-QA</th>
<td class="ltx_td ltx_align_center">10</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">LongMemEval</th>
<td class="ltx_td ltx_align_center">100</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">EventQA</th>
<td class="ltx_td ltx_align_center">40</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t">ICL_Five</th>
<td class="ltx_td ltx_align_center ltx_border_t">20</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">Movie Recommendation</th>
<td class="ltx_td ltx_align_center">300</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r ltx_border_t"><span class="math inline">∞</span> Bench-Sum</th>
<td class="ltx_td ltx_align_center ltx_border_t">1,200</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_bb ltx_border_r ltx_border_t">FactConsolidation</th>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t">10</td>
</tr>
</tbody>
</table>
</div>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 13</span>: </span><span class="ltx_text" style="font-size:90%;">Computational latency (in seconds) comparison on RAG based agents. M.C. means Memory Construction and Q.E. means Query Execution.</span>
<span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 14</span>: </span><span class="ltx_text" style="font-size:90%;">Maximum output token limits for various tasks</span></figcaption>
</figure>
::::::::
::::::::::::::

:::::::::::: {#A4 .section .ltx_appendix}
## [Appendix D ]{.ltx_tag .ltx_tag_appendix}Experimental Settings {#appendix-d-experimental-settings .ltx_title .ltx_title_appendix}

::: {#A4.p1 .ltx_para}
In this section, we present the experimental settings.
:::

:::: {#A4.SS1 .section .ltx_subsection}
### [D.1 ]{.ltx_tag .ltx_tag_subsection}Max Output Tokens {#d.1-max-output-tokens .ltx_title .ltx_title_subsection}

::: {#A4.SS1.p1 .ltx_para}
We provide the token number limitation for each task in Table  [[14]{.ltx_text .ltx_ref_tag}](#A3.T14 "Table 14 ‣ C.3 Detailed Results on Ablation Study ‣ Appendix C Detailed Experimental Results ‣ Evaluating Memory in LLM Agents via Incremental Multi-Turn Interactions"){.ltx_ref}.
:::
::::

::::: {#A4.SS2 .section .ltx_subsection}
### [D.2 ]{.ltx_tag .ltx_tag_subsection}Settings of the RAG Agents {#d.2-settings-of-the-rag-agents .ltx_title .ltx_title_subsection}

::: {#A4.SS2.p1 .ltx_para}
For the embedding model selection in Structure-Augmented RAG Agents and Agentic Memory Agents, most approaches utilize OpenAI's embedding models, such as Text-Embed-3-Small. While for the HippoRAG-v2 method, we follow the same experimental setting as in Gutiérrez et al. \[[2025](#bib.bib12){.ltx_ref}\], employing the NV-Embed-v2 model.
:::

::: {#A4.SS2.p2 .ltx_para}
We implement three open-sourced memory agents in our main experiments. (1) For Mem0, we use [memory.add()]{.ltx_text .ltx_font_bold} function to add the message with the content from each context chunk into the agent's memory repository during memory consolidation. During query execution, the relevant memory elements are retrieved through [memory.search()]{.ltx_text .ltx_font_bold} function. The retrieved memories are then integrated into the query before being processed by the GPT-4o-mini backbone model to complete the requested tasks. (2) For MemGPT, we employ the [insert_passage()]{.ltx_text .ltx_font_bold} function during the memory consolidation phase to inject long context chunks into the Archival Memory structure. During query execution, this agent processes requests via the [send_message()]{.ltx_text .ltx_font_bold} function which generates appropriate responses based on the archived information. (3) For Cognee, we utilize the [cognee.add()]{.ltx_text .ltx_font_bold} and [cognee.cognify()]{.ltx_text .ltx_font_bold} functions to construct the memory graph from input chunks wherein the memory consolidation phase. During query execution, the [cognee.search()]{.ltx_text .ltx_font_bold} function is used to retrieve contextually relevant information from the memory graph based on the input query.
:::
:::::

::::: {#A4.SS3 .section .ltx_subsection}
### [D.3 ]{.ltx_tag .ltx_tag_subsection}Settings of the Chunk Size {#d.3-settings-of-the-chunk-size .ltx_title .ltx_title_subsection}

::: {#A4.SS3.p1 .ltx_para}
We use smaller chunk size (512) for synthetic context used in AR and CR. For some tasks based on continuous text, such as $\infty$Bench and EventQA, we used a larger chunk size (4096). For tasks such as MCC, Recom and LME(S), considering the characteristics of these tasks and the computational cost, we also chose a larger chunk size (4096). For the two memory construction methods that are more time-consuming, Mem0 and Cognee, we uniformly used a chunk size of 4096 across all datasets.
:::

<figure id="A4.T15" class="ltx_table">
<table class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<thead class="ltx_thead">
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_th_row ltx_border_r ltx_border_tt"><span class="ltx_text ltx_font_bold">Chunk Size</span></th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_tt">512</th>
<th class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_tt">4096</th>
</tr>
</thead>
<tbody class="ltx_tbody">
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_row ltx_border_r ltx_border_t"></th>
<td class="ltx_td ltx_align_center ltx_border_t">RULER-QA, NIAH-MQ</td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="math inline">∞</span>Bench-QA, <span class="math inline">∞</span>Bench-Sum</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_align_center ltx_th ltx_th_row ltx_border_r">Dataset</th>
<td class="ltx_td ltx_align_center">FactCon-SH, FactCon-MH</td>
<td class="ltx_td ltx_align_center">MCC, Recom</td>
</tr>
<tr class="ltx_tr">
<th class="ltx_td ltx_th ltx_th_row ltx_border_bb ltx_border_r"></th>
<td class="ltx_td ltx_align_center ltx_border_bb">LME(S*)</td>
<td class="ltx_td ltx_align_center ltx_border_bb">EventQA, LME(S)</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 15</span>: </span><span class="ltx_text" style="font-size:90%;">The choice of chunk size for different datasets.</span></figcaption>
</figure>

::: {.ltx_pagination .ltx_role_newpage}
:::
:::::
::::::::::::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

::: ar5iv-footer
[◄](/html/2507.05256){.ar5iv-nav-button .ar5iv-nav-button-prev} [![ar5iv homepage](/assets/ar5iv.png){height="40"}](/){.ar5iv-home-button} [Feeling\
lucky?](/feeling_lucky){.ar5iv-text-button} [](/land_of_honey_and_milk){rel="nofollow" aria-hidden="true" tabindex="-1"} [Conversion\
report](/log/2507.05257){.ar5iv-text-button .ar5iv-severity-warning} [Report\
an issue](https://github.com/dginev/ar5iv/issues/new?template=improve-article--arxiv-id-.md&title=Improve+article+2507.05257){.ar5iv-text-button target="_blank"} [View original\
on arXiv](https://arxiv.org/abs/2507.05257){.ar5iv-text-button .arxiv-ui-theme}[►](/html/2507.05258){.ar5iv-nav-button .ar5iv-nav-button-next}
:::

[[]{.color-scheme-icon}](javascript:toggleColorScheme() "Toggle ar5iv color scheme"){.ar5iv-toggle-color-scheme} [Copyright](https://arxiv.org/help/license){.ar5iv-footer-button target="_blank"} [Privacy Policy](https://arxiv.org/help/policies/privacy_policy){.ar5iv-footer-button target="_blank"}

::: ltx_page_logo
Generated on Tue Aug 5 15:54:09 2025 by [[L[a]{.ltx_font_smallcaps style="position:relative; bottom:2.2pt;"}T[e]{.ltx_font_smallcaps style="font-size:120%;position:relative; bottom:-0.2ex;"}]{style="letter-spacing:-0.2em; margin-right:0.1em;"}[XML]{style="font-size:90%; position:relative; bottom:-0.2ex;"}![Mascot Sammy](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAsAAAAOCAYAAAD5YeaVAAAAAXNSR0IArs4c6QAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAALEwAACxMBAJqcGAAAAAd0SU1FB9wKExQZLWTEaOUAAAAddEVYdENvbW1lbnQAQ3JlYXRlZCB3aXRoIFRoZSBHSU1Q72QlbgAAAdpJREFUKM9tkL+L2nAARz9fPZNCKFapUn8kyI0e4iRHSR1Kb8ng0lJw6FYHFwv2LwhOpcWxTjeUunYqOmqd6hEoRDhtDWdA8ApRYsSUCDHNt5ul13vz4w0vWCgUnnEc975arX6ORqN3VqtVZbfbTQC4uEHANM3jSqXymFI6yWazP2KxWAXAL9zCUa1Wy2tXVxheKA9YNoR8Pt+aTqe4FVVVvz05O6MBhqUIBGk8Hn8HAOVy+T+XLJfLS4ZhTiRJgqIoVBRFIoric47jPnmeB1mW/9rr9ZpSSn3Lsmir1fJZlqWlUonKsvwWwD8ymc/nXwVBeLjf7xEKhdBut9Hr9WgmkyGEkJwsy5eHG5vN5g0AKIoCAEgkEkin0wQAfN9/cXPdheu6P33fBwB4ngcAcByHJpPJl+fn54mD3Gg0NrquXxeLRQAAwzAYj8cwTZPwPH9/sVg8PXweDAauqqr2cDjEer1GJBLBZDJBs9mE4zjwfZ85lAGg2+06hmGgXq+j3+/DsixYlgVN03a9Xu8jgCNCyIegIAgx13Vfd7vdu+FweG8YRkjXdWy329+dTgeSJD3ieZ7RNO0VAXAPwDEAO5VKndi2fWrb9jWl9Esul6PZbDY9Go1OZ7PZ9z/lyuD3OozU2wAAAABJRU5ErkJggg==)](http://dlmf.nist.gov/LaTeXML/){.ltx_LaTeXML_logo target="_blank"}
:::
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
