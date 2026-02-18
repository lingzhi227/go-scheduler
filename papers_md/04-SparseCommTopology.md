:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_main
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_content
# Improving Multi-Agent Debate with Sparse Communication Topology {#improving-multi-agent-debate-with-sparse-communication-topology .ltx_title .ltx_title_document}

::: ltx_authors
[ [ [Yunxuan Li^1$\dagger$^]{#id1.1.1 .ltx_text .ltx_font_bold}, [Yibing Du^1^]{#id2.2.id1 .ltx_text .ltx_font_bold}, [Jiageng Zhang^1^]{#id3.3.id2 .ltx_text .ltx_font_bold}, [Le Hou^2^]{#id4.4.id3 .ltx_text .ltx_font_bold},\
[Peter Grabowski^1^]{#id5.5.id4 .ltx_text .ltx_font_bold}, [Yeqing Li^1^]{#id6.6.id5 .ltx_text .ltx_font_bold}, [Eugene Ie^1^]{#id7.7.id6 .ltx_text .ltx_font_bold}\
^1^Google ^2^Google DeepMind ]{.ltx_personname}]{.ltx_creator .ltx_role_author}
:::

::: ltx_abstract
###### Abstract {#abstract .ltx_title .ltx_title_abstract}

Multi-agent debate has proven effective in improving large language models quality for reasoning and factuality tasks. While various role-playing strategies in multi-agent debates have been explored, in terms of the communication among agents, existing approaches adopt a brute force algorithm -- each agent can communicate with all other agents. In this paper, we systematically investigate the effect of communication connectivity in multi-agent systems. Our experiments on GPT and Mistral models reveal that multi-agent debates leveraging sparse communication topology can achieve comparable or superior performance while significantly reducing computational costs. Furthermore, we extend the multi-agent debate framework to multimodal reasoning and alignment labeling tasks, showcasing its broad applicability and effectiveness. Our findings underscore the importance of communication connectivity on enhancing the efficiency and effectiveness of the "society of minds" approach.
:::

:::: {#p1 .ltx_para .ltx_noindent}
::: {#p1.1 .ltx_block .ltx_align_bottom}
[Improving Multi-Agent Debate with Sparse Communication Topology]{#p1.1.2.1 .ltx_text .ltx_font_bold}

\

[ [ [ [ [[ Yunxuan Li^1$\dagger$^, Yibing Du^1^, Jiageng Zhang^1^, Le Hou^2^,]{#p1.1.1.1.1.1.1.1 .ltx_text .ltx_font_bold}]{#p1.1.1.1.1.1.1 .ltx_td .ltx_align_center}]{#p1.1.1.1.1.1 .ltx_tr} [ [[Peter Grabowski^1^]{#p1.1.1.1.1.2.1.1.1 .ltx_text .ltx_font_bold}, [Yeqing Li^1^]{#p1.1.1.1.1.2.1.1.2 .ltx_text .ltx_font_bold}, [Eugene Ie^1^]{#p1.1.1.1.1.2.1.1.3 .ltx_text .ltx_font_bold}]{#p1.1.1.1.1.2.1.1 .ltx_td .ltx_align_center}]{#p1.1.1.1.1.2.1 .ltx_tr} [ [^1^Google ^2^Google DeepMind]{#p1.1.1.1.1.3.2.1 .ltx_td .ltx_align_center}]{#p1.1.1.1.1.3.2 .ltx_tr} ]{.ltx_tbody} ]{#p1.1.1.1.1 .ltx_tabular .ltx_align_top}]{#p1.1.1.1 .ltx_text .ltx_inline-block style="width:0.0pt;"}

\
:::
::::

[^2^[[^2^[footnotetext: ]{.ltx_note_type}Correspondence: yunxuanli@google.com]{.ltx_note_content}]{.ltx_note_outer}]{#footnotex1 .ltx_note .ltx_role_footnotetext}

::::::::: {#S1 .section .ltx_section}
## [1 ]{.ltx_tag .ltx_tag_section}Introduction {#introduction .ltx_title .ltx_title_section}

::: {#S1.p1 .ltx_para}
Large language models (LLMs) have demonstrated exceptional performance in natural language understanding and generation tasks. Recently a paradigm shift towards prompting LLMs has emerged as a significant and influential research area. By leveraging the in-context learning (ICL) capabilities of LLMs, these models can be adapted to various tasks such as reasoning, factuality, and AI feedback.
:::

<figure id="S1.F1" class="ltx_figure">
<img src="/html/2406.11776/assets/imgs/Figure1_v8.png" id="S1.F1.g1" class="ltx_graphics ltx_img_square" width="598" height="706" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 1: </span>Accuracy (top) and inference input cost (middle) comparison of multi-agent debate system between fully-connected (bottom left) and neighbor-connected (bottom right) communication topologies.</figcaption>
</figure>

::: {#S1.p2 .ltx_para}
Several prompting methods have been developed to enhance LLM performance by optimizing their ICL capabilities. Notable techniques include Chain-of-Thought (CoT) Wei et al. ([2022](#bib.bib37){.ltx_ref}), self-consistency (SC) Wang et al. ([2022](#bib.bib35){.ltx_ref}), and self-critique Madaan et al. ([2024](#bib.bib26){.ltx_ref}); Welleck et al. ([2022](#bib.bib38){.ltx_ref}); Shinn et al. ([2024](#bib.bib33){.ltx_ref}). Recently, the multi-agent debate (MAD) framework is proven to be an innovative approach. Similar to a human discussion process, MAD employs multiple LLM agents to engage in discussions with one another, combining their reasoning and critical thinking abilities to produce high-quality results. Specifically, given a question, each agent first generates their own solutions to the question and then takes other agents' solutions as reference to update its own answer. This process can be repeated for several rounds. MAD has demonstrated significant improvement on factuality and reasoning tasks. While the debate process is highly productive, it is also very costly: As the number of LLM agents and debate rounds increase, the input context expands significantly.
:::

::: {#S1.p3 .ltx_para}
Inspired by the intensive computational cost of MAD, a natural question arises: [What if we reduce the number of reference solutions visible to each agent?]{#S1.p3.2.1 .ltx_text .ltx_font_italic} We conduct a systematic study on the sparsity of the multi-agent communication topology. Surprisingly, we find that sparse communication connectivity can deliver comparable or superior performance while significantly reducing inference costs. Figure [[1]{.ltx_text .ltx_ref_tag}](#S1.F1 "Figure 1 ‣ 1 Introduction ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref} presents a comparison between fully-connected MAD and neighbor-connected MAD. Compared to fully-connected MAD, neighbor-connected MAD achieves an improvement of $+ {2\%}$ on the MATH dataset and maintains the same accuracy on GSM8K. Meanwhile, the average input token cost for reasoning tasks is reduced by over $40\%$.
:::

::: {#S1.p4 .ltx_para}
MAD can also be a promising approach for Reinforcement Learning with AI Feedback (RLAIF) Bai et al. ([2022b](#bib.bib2){.ltx_ref}); Lee et al. ([2023](#bib.bib17){.ltx_ref}) and weak-to-strong generalization Burns et al. ([2023](#bib.bib4){.ltx_ref}). By delivering better reward signals, MAD has the potential to significantly aid in aligning large language models. To assess this, we first extend the MAD framework to alignment labeling tasks, demonstrating its effectiveness compared to single-agent setups. Additionally, we verify that the advantages of sparsity observed in the reasoning tasks experiments also apply to alignment labeling tasks. Our experiments on the Anthropic-HH datasets show an improvement of $+ {0.5\%}$ in helpfulness and $+ {1.0\%}$ in harmlessness, while reducing costs by $50.0\%$ and $53.3\%$, respectively.
:::

::: {#S1.p5 .ltx_para}
We find that when agents are instantiated by different LLMs within the MAD framework, interactions between multiple LLMs result in weaker models being progressively strengthened through engagement with stronger models. In non-regular graph settings, assigning stronger LLMs to agents with higher centrality consistently yields better performance.
:::

::: {#S1.p6 .ltx_para}
In summary, our contributions are listed as follows: (1) We demonstrate that sparse communication topology enhances both effectiveness and efficiency of the multi-agent debate framework; (2) We thoroughly evaluate sparse MAD for text-only and multimodal reasoning tasks, showing its advantage over standard MAD; (3) We extend the MAD framework to alignment labeling tasks, showing the effectiveness of standard MAD and further performance improvement introduced by sparse MAD; (4) We provide insights that explain the effectiveness of sparsity in MAD; (5) We find that assigning stronger LLMs to agents with higher centrality yields better overall performance in multiple LLM debate setup.
:::
:::::::::

::::::: {#S2 .section .ltx_section}
## [2 ]{.ltx_tag .ltx_tag_section}Related Work {#related-work .ltx_title .ltx_title_section}

::: {#S2.p1 .ltx_para}
[Multi-Agent Debate]{#S2.p1.1.1 .ltx_text .ltx_font_bold} MAD utilizes multiple LLM agents to discuss and debate with each other to generate and update the responses. It was first introduced by Du et al. ([2023](#bib.bib9){.ltx_ref}). Most of the MAD work focus on diversifying agents during the debate process. Liang et al. ([2023](#bib.bib20){.ltx_ref}); Park et al. ([2023](#bib.bib30){.ltx_ref}); Li et al. ([2023a](#bib.bib18){.ltx_ref}); Chan et al. ([2023](#bib.bib5){.ltx_ref}) highlight the importance of assigning different roles for agents. Chen et al. ([2023](#bib.bib6){.ltx_ref}) diversifies agents' responses by instantiated with multiple LLMs. Wang et al. ([2024a](#bib.bib34){.ltx_ref}) proposes a method in which agents are divided into sub-groups and their discussion outcomes are later merged. Unlike other work, we aim to explore the effectiveness of sparse communication topology in MAD, and extend its applications to reasoning and alignment tasks.
:::

::: {#S2.p2 .ltx_para .ltx_noindent}
[LLM Reasoning]{#S2.p2.1.1 .ltx_text .ltx_font_bold} Much work has been done to improve the reasoning ability of language models with prompting, including Chain-of-Thought Wei et al. ([2022](#bib.bib37){.ltx_ref}) and its variants Yao et al. ([2024](#bib.bib40){.ltx_ref}); Besta et al. ([2024](#bib.bib3){.ltx_ref}), problem decomposition Zhou et al. ([2022](#bib.bib46){.ltx_ref}), reasoning ensemble Wang et al. ([2022](#bib.bib35){.ltx_ref}), reasoner with verification Cobbe et al. ([2021](#bib.bib8){.ltx_ref}); Wang et al. ([2024b](#bib.bib36){.ltx_ref}); Luo et al. ([2023](#bib.bib25){.ltx_ref}).
:::

::: {#S2.p3 .ltx_para .ltx_noindent}
[Multimodal Reasoning]{#S2.p3.1.1 .ltx_text .ltx_font_bold} With the recent advancements in vision-language models Radford et al. ([2021](#bib.bib31){.ltx_ref}); Yu et al. ([2022](#bib.bib41){.ltx_ref}); Li et al. ([2023b](#bib.bib19){.ltx_ref}); Liu et al. ([2024](#bib.bib22){.ltx_ref}); Lin et al. ([2024](#bib.bib21){.ltx_ref}), multimodal large-language models (MLLMs) have demonstrated exceptional visual understanding capabilities. Several evaluation benchmarks have been proposed, such as VQAv2 Goyal et al. ([2017](#bib.bib11){.ltx_ref}), OK-VQA Marino et al. ([2019](#bib.bib27){.ltx_ref}), ScienceQA Lu et al. ([2022](#bib.bib24){.ltx_ref}), MMMU Yue et al. ([2023](#bib.bib42){.ltx_ref}), and MathVista Lu et al. ([2023](#bib.bib23){.ltx_ref}). Similar to LLMs, MLLMs can also be improved through prompt-based methods. Various attempts have been made to enhance MLLMs in this manner Zheng et al. ([2024](#bib.bib45){.ltx_ref}); Ganz et al. ([2024](#bib.bib10){.ltx_ref}); Yang et al. ([2023](#bib.bib39){.ltx_ref}); Zhao et al. ([2024](#bib.bib44){.ltx_ref}); Zhang et al. ([2023](#bib.bib43){.ltx_ref}); Chen et al. ([2024](#bib.bib7){.ltx_ref}); Hu et al. ([2024](#bib.bib14){.ltx_ref}). Despite the effectiveness of these methods, they are often complex to design and implement. In this paper, we focus on improving multimodal reasoning using a multi-agent approach.
:::

::: {#S2.p4 .ltx_para .ltx_noindent}
[AI Feedback]{#S2.p4.1.1 .ltx_text .ltx_font_bold} Bai et al. ([2022b](#bib.bib2){.ltx_ref}) first introduces the idea of RLAIF, in which LLM is used to annotate harmlessness preference. Lee et al. ([2023](#bib.bib17){.ltx_ref}) compares various AI annotation methods. Recent work (Guo et al., [2024](#bib.bib12){.ltx_ref}) also explores using AI feedback for online reinforcement learning, demonstrating the advantage of AI feedback for alignment research.
:::
:::::::

:::::::::::: {#S3 .section .ltx_section}
## [3 ]{.ltx_tag .ltx_tag_section}Method {#method .ltx_title .ltx_title_section}

:::::: {#S3.SS1 .section .ltx_subsection}
### [3.1 ]{.ltx_tag .ltx_tag_subsection}Communication Topology {#communication-topology .ltx_title .ltx_title_subsection}

::: {#S3.SS1.p1 .ltx_para}
Communication topology of MAD refers to the connectivity structure among agents during the debate process. Communication topology can be represented as a graph $\mathcal{G} = {(\mathcal{V},\mathcal{E})}$, where $\mathcal{V}$ is a set of agents and $\mathcal{E}$ is a set of communication channel. Presence of of any $(e_{i},e_{j})$ in $\mathcal{E}$ indicates that agent $i$ can access agent $j$'s previous round solutions during the debate process, and vice versa. We focus on static graphs in this work, while we also did exploratory experiments with dynamic graphs (Appendix [[D]{.ltx_text .ltx_ref_tag}](#A4 "Appendix D ProbMAD: MAD with Probablistic Topology ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref}).
:::

::: {#S3.SS1.p2 .ltx_para}
We quantify the density of these graphs using the ratio of the number of edges to the maximum possible number of edges

  -- --------------------------------------------------------------------------------------------------- --
     $$D = \frac{2\hspace{0pt}{|\mathcal{E}|}}{{|\mathcal{V}|}\hspace{0pt}{({{|\mathcal{V}|} - 1})}}$$   
  -- --------------------------------------------------------------------------------------------------- --
:::

::: {#S3.SS1.p3 .ltx_para}
A lower value of $D$ indicates a sparser graph. In the standard MAD framework, agents are fully connected with each other, resulting in $D = 1$. In contrast, a neighbor-connected MAD has ${|\mathcal{E}|} = {|\mathcal{V}|}$, yielding $D = \frac{2}{{|\mathcal{V}|} - 1}$, which is a sparse graph. While the findings of this paper can be generalized to communication topology with an arbitrary number of agents, we focus on regular graphs where all agents have same degrees and are permutation invariant, with ${|\mathcal{V}|} = 6$ (Figure [[2]{.ltx_text .ltx_ref_tag}](#S3.F2 "Figure 2 ‣ 3.1 Communication Topology ‣ 3 Method ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref}). This choice is due to the limited spectrum of sparsity in scenarios with fewer agents and the significantly higher computational costs associated with analyzing scenarios with more agents. Additional experiment results with ${|\mathcal{V}|} = 4$ is shown in Appendix [[C]{.ltx_text .ltx_ref_tag}](#A3 "Appendix C Additional Experiments with 4 Agents ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref}.
:::

<figure id="S3.F2" class="ltx_figure">
<img src="/html/2406.11776/assets/imgs/sparsity_graphs.png" id="S3.F2.g1" class="ltx_graphics ltx_img_landscape" width="598" height="127" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 2: </span>Communication topology of 6 agents with various sparsity. From left to right, the densities are 1 (fully-connected), <span class="math inline">$\frac{4}{5}$</span>, <span class="math inline">$\frac{3}{5}$</span>, and <span class="math inline">$\frac{2}{5}$</span> (neighbor-connected).</figcaption>
</figure>
::::::

::::::: {#S3.SS2 .section .ltx_subsection}
### [3.2 ]{.ltx_tag .ltx_tag_subsection}Multi-Agent Debate Process {#multi-agent-debate-process .ltx_title .ltx_title_subsection}

::: {#S3.SS2.p1 .ltx_para}
A typical MAD framework includes three steps:
:::

::: {#S3.SS2.p2 .ltx_para .ltx_noindent}
[(1) Individual Response Generation]{#S3.SS2.p2.1.1 .ltx_text .ltx_font_bold}: In round 1, agents are initialized with LLMs, and then independently generate solutions to a given question. Typically a random decoding strategy is applied to diversify the solutions generated by agents.
:::

::: {#S3.SS2.p3 .ltx_para .ltx_noindent}
[(2) Multi-agent Debate]{#S3.SS2.p3.1.1 .ltx_text .ltx_font_bold}: Starting round 2, each agent incorporates the responses of its connected peers from the previous round to critique or refine its own response. We utilize the standard [Simultaneous-Talk]{#S3.SS2.p3.1.2 .ltx_text .ltx_font_italic} communication strategy Chan et al. ([2023](#bib.bib5){.ltx_ref}) to facilitate asynchronous computation. This debating process can occur over multiple rounds.
:::

::: {#S3.SS2.p4 .ltx_para .ltx_noindent}
[(3) Reaching Consensus]{#S3.SS2.p4.1.1 .ltx_text .ltx_font_bold}: After the debate process, agents may still have differing solutions. In such cases, a majority vote is conducted among all agents to determine a consensus solution.
:::
:::::::
::::::::::::

::::::::::::::::: {#S4 .section .ltx_section}
## [4 ]{.ltx_tag .ltx_tag_section}Experiments Setup {#experiments-setup .ltx_title .ltx_title_section}

:::: {#S4.SS1 .section .ltx_subsection}
### [4.1 ]{.ltx_tag .ltx_tag_subsection}Tasks {#tasks .ltx_title .ltx_title_subsection}

::: {#S4.SS1.p1 .ltx_para}
We aim to validate the effectiveness and efficiency of sparse MAD on reasoning and alignment labeling tasks. For reasoning tasks, we consider two text-only reasoning tasks and one multimodal reasoning task: (1) MATH Hendrycks et al. ([2021](#bib.bib13){.ltx_ref}): an arithmetic reasoning task containing challenging competition mathematics problems. We only use the [algebra linear 1d composed]{#S4.SS1.p1.1.1 .ltx_text .ltx_font_italic} sub-task for simplicity. (2) GSM8K Cobbe et al. ([2021](#bib.bib8){.ltx_ref}): a high quality grade school math reasoning task. (3) MathVista Lu et al. ([2023](#bib.bib23){.ltx_ref}): a benchmark designed to combine challenges from diverse mathematical and visual tasks. We only choose from [free_form]{#S4.SS1.p1.1.2 .ltx_text .ltx_font_italic} question type for consistency. For alignment labeling tasks, we consider Anthropic-HH dataset Bai et al. ([2022a](#bib.bib1){.ltx_ref}): human preference data on helpfulness and harmlessness.
:::
::::

:::: {#S4.SS2 .section .ltx_subsection}
### [4.2 ]{.ltx_tag .ltx_tag_subsection}Models {#models .ltx_title .ltx_title_subsection}

::: {#S4.SS2.p1 .ltx_para}
Our experiments utilize three publicly available models: GPT-3.5 OpenAI ([2022](#bib.bib28){.ltx_ref}), GPT-4 OpenAI ([2023](#bib.bib29){.ltx_ref}), and Mistral 7B Jiang et al. ([2023](#bib.bib15){.ltx_ref}). Specifically, we employ GPT-3.5 for text-only reasoning tasks and GPT-4 for multimodal reasoning tasks. For alignment labeling tasks, we use both GPT-3.5 and Mistral 7B. We refrain from using GPT-4 for other tasks due to its significantly higher cost, which is approximately 10 times that of GPT-3.5. Additionally, we do not employ Mistral 7B for other tasks because of its inferior zero-shot performance on arithmetic reasoning. We randomly select 100 examples for each experiments involving GPT, and 500 examples for experiments with Mistral 7B.
:::
::::

::::::: {#S4.SS3 .section .ltx_subsection}
### [4.3 ]{.ltx_tag .ltx_tag_subsection}Baselines {#baselines .ltx_title .ltx_title_subsection}

::: {#S4.SS3.p1 .ltx_para}
We compare sparse MAD against the following baselines:
:::

::: {#S4.SS3.p2 .ltx_para}
\(1\) Chain-of-Thought ([CoT]{#S4.SS3.p2.1.1 .ltx_text .ltx_font_bold}): CoT prompting improves reasoning capabilities of LLMs with explicit intermediate reasoning steps.
:::

::: {#S4.SS3.p3 .ltx_para}
\(2\) Self-consistency ([SC]{#S4.SS3.p3.1.1 .ltx_text .ltx_font_bold}): SC margins out intermediate reasoning paths by sampling diverse reasoning paths and selecting the most consistent answer.
:::

::: {#S4.SS3.p4 .ltx_para}
\(3\) Existing MAD ([MAD]{#S4.SS3.p4.1.1 .ltx_text .ltx_font_bold} ($D = 1$)): the standard approach for multi-agent debate, in which agents can communicate with all other agents with simultaneous-talk strategy. We also denote it as fully-connected MAD.
:::
:::::::

::::: {#S4.SS4 .section .ltx_subsection}
### [4.4 ]{.ltx_tag .ltx_tag_subsection}Evaluation Metrics {#evaluation-metrics .ltx_title .ltx_title_subsection}

::: {#S4.SS4.p1 .ltx_para}
For reasoning tasks, we use the accuracy with respect to the ground truth answer to measure the quality of MAD. For alignment labeling tasks, we use [AI Labeler Alignment]{#S4.SS4.p1.1.1 .ltx_text .ltx_font_italic} Lee et al. ([2023](#bib.bib17){.ltx_ref}) to measure the accuracy of MAD labeling with respect to the human annotation.
:::

::: {#S4.SS4.p2 .ltx_para}
Cost refers to the input inference cost of LLMs, which typically involves handling the autoregressive decoding mechanism and computational resources. Considering that advanced LLMs use a pay-per-token pricing model, we measure the inference cost by the number of input tokens.
:::
:::::

:::: {#S4.SS5 .section .ltx_subsection}
### [4.5 ]{.ltx_tag .ltx_tag_subsection}Variance Reduction {#variance-reduction .ltx_title .ltx_title_subsection}

::: {#S4.SS5.p1 .ltx_para}
Evaluating the significance of new communication topology compared to existing one typically involves running multiple random experiments to estimate the mean and variance of performance. However, this approach becomes impractical when the signal-to-noise ratio is low and each experimental run is computationally expensive. To address this, we employ two methods to reduce experimental variance and enhance the sensitivity of MAD with respect to the changes in communication topology: (1) As used by Wang et al. ([2024a](#bib.bib34){.ltx_ref}), we reduce the temperature during language model decoding to stabilize performance. While we use the default temperature settings in API calls for most tasks, we lower the temperature to 0.25 for text arithmetic reasoning tasks to ensure robustness. (2) We employ conditional variance reduction Ross ([2002](#bib.bib32){.ltx_ref}). Observing that most of the variance arises from the first round of individual responses, we first generate a set of initial agent responses and then fix them in all subsequent debate processes across various communication topology designs. This approach effectively minimizes variance and improves the reliability of our experimental results.
:::
::::
:::::::::::::::::

::::::::::::::::::: {#S5 .section .ltx_section}
## [5 ]{.ltx_tag .ltx_tag_section}Experiments: MAD with Single LLM {#experiments-mad-with-single-llm .ltx_title .ltx_title_section}

::::: {#S5.SS1 .section .ltx_subsection}
### [5.1 ]{.ltx_tag .ltx_tag_subsection}MAD on Text Reasoning Tasks {#mad-on-text-reasoning-tasks .ltx_title .ltx_title_subsection}

::: {#S5.SS1.p1 .ltx_para}
We build on existing work on MAD, exemplified by reasoning tasks, by showing the advantages of sparse MAD on top of the proven advantage of fully-connected MAD. Sparse MAD significantly saves computational cost while preserving comparable or better performance.
:::

::: {#S5.SS1.p2 .ltx_para}
[Sparse MAD has similar or higher accuracy with significant cost saving on reasoning tasks]{#S5.SS1.p2.8.1 .ltx_text .ltx_font_bold}: For both the MATH and GSM8K tasks, we demonstrate that sparse MAD produces comparable or better accuracy than fully-connected MAD, while significantly cutting down inference costs. Both fully-connected and sparse MAD setups outperform Chain-of-Thought and self-consistency methods. Specifically, in the MATH task, fully-connected MAD shows a $+ {4.0\%}$ quality gain over self-consistency, while sparse MAD configurations achieve accuracy improvements ranging from $+ {3.0\%}$ to $+ {7.5\%}$ (Table [[1]{.ltx_text .ltx_ref_tag}](#S5.T1 "Table 1 ‣ 5.1 MAD on Text Reasoning Tasks ‣ 5 Experiments: MAD with Single LLM ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref}). Similarly, in the GSM8K task, fully-connected MAD demonstrates a $+ {4.5\%}$ quality gain over self-consistency, whereas sparse MAD achieves accuracy improvements between $+ {3.5\%}$ and $+ {6.5\%}$ (Table [[2]{.ltx_text .ltx_ref_tag}](#S5.T2 "Table 2 ‣ 5.1 MAD on Text Reasoning Tasks ‣ 5 Experiments: MAD with Single LLM ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref}). Furthermore, sparse MAD setups reduce costs by up to $- {41.5\%}$ and $- {43.5\%}$, respectively. It is important to note that we exclusively use the GPT-3.5 model because Mistral 7B performs poorly on these challenging tasks in a zero-shot setting.
:::

<figure id="S5.T1" class="ltx_table">
<table id="S5.T1.12" class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr id="S5.T1.12.13.1" class="ltx_tr">
<th id="S5.T1.12.13.1.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t"><span id="S5.T1.12.13.1.1.1" class="ltx_text ltx_font_bold">Method</span></th>
<td id="S5.T1.12.13.1.2" class="ltx_td ltx_align_center ltx_border_t">Accuracy</td>
<td id="S5.T1.12.13.1.3" class="ltx_td ltx_align_center ltx_border_t">Cost Saving</td>
</tr>
<tr id="S5.T1.1.1" class="ltx_tr">
<th id="S5.T1.1.1.2" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">CoT</th>
<td id="S5.T1.1.1.1" class="ltx_td ltx_align_center ltx_border_t">58.0 <span class="math inline">±</span> 2.0</td>
<td id="S5.T1.1.1.3" class="ltx_td ltx_align_center ltx_border_t">-</td>
</tr>
<tr id="S5.T1.12.14.2" class="ltx_tr">
<th id="S5.T1.12.14.2.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">SC</th>
<td id="S5.T1.12.14.2.2" class="ltx_td ltx_align_center">60.0</td>
<td id="S5.T1.12.14.2.3" class="ltx_td ltx_align_center">-</td>
</tr>
<tr id="S5.T1.3.3" class="ltx_tr">
<th id="S5.T1.2.2.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">MAD (<span class="math inline"><em>D</em> = 1</span>)</th>
<td id="S5.T1.3.3.2" class="ltx_td ltx_align_center ltx_border_t">64.0 <span class="math inline">±</span> 1.4</td>
<td id="S5.T1.3.3.3" class="ltx_td ltx_align_center ltx_border_t">baseline</td>
</tr>
<tr id="S5.T1.6.6" class="ltx_tr">
<th id="S5.T1.4.4.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">MAD (<span class="math inline"><em>D</em> = 4/5</span>)</th>
<td id="S5.T1.5.5.2" class="ltx_td ltx_align_center"><span id="S5.T1.5.5.2.1" class="ltx_text ltx_font_bold">67.5 <span class="math inline">±</span> 2.0</span></td>
<td id="S5.T1.6.6.3" class="ltx_td ltx_align_center"><span class="math inline">−</span>14.6%</td>
</tr>
<tr id="S5.T1.9.9" class="ltx_tr">
<th id="S5.T1.7.7.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">MAD (<span class="math inline"><em>D</em> = 3/5</span>)</th>
<td id="S5.T1.8.8.2" class="ltx_td ltx_align_center">63.0 <span class="math inline">±</span> 1.8</td>
<td id="S5.T1.9.9.3" class="ltx_td ltx_align_center"><span class="math inline">−</span>29.2%</td>
</tr>
<tr id="S5.T1.12.12" class="ltx_tr">
<th id="S5.T1.10.10.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_b ltx_border_r">MAD (<span class="math inline"><em>D</em> = 2/5</span>)</th>
<td id="S5.T1.11.11.2" class="ltx_td ltx_align_center ltx_border_b">66.0 <span class="math inline">±</span> 2.3</td>
<td id="S5.T1.12.12.3" class="ltx_td ltx_align_center ltx_border_b"><span class="math inline">−</span>41.5%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 1: </span>Comparison of accuracy and cost savings of MAD against baseline methods on the MATH dataset. All experiments were conducted using the GPT-3.5 model.</figcaption>
</figure>

<figure id="S5.T2" class="ltx_table">
<table id="S5.T2.12" class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr id="S5.T2.12.13.1" class="ltx_tr">
<th id="S5.T2.12.13.1.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t"><span id="S5.T2.12.13.1.1.1" class="ltx_text ltx_font_bold">Method</span></th>
<td id="S5.T2.12.13.1.2" class="ltx_td ltx_align_center ltx_border_t">Accuracy</td>
<td id="S5.T2.12.13.1.3" class="ltx_td ltx_align_center ltx_border_t">Cost Saving</td>
</tr>
<tr id="S5.T2.1.1" class="ltx_tr">
<th id="S5.T2.1.1.2" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">CoT</th>
<td id="S5.T2.1.1.1" class="ltx_td ltx_align_center ltx_border_t">77.5 <span class="math inline">±</span> 4.2</td>
<td id="S5.T2.1.1.3" class="ltx_td ltx_align_center ltx_border_t">-</td>
</tr>
<tr id="S5.T2.12.14.2" class="ltx_tr">
<th id="S5.T2.12.14.2.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">SC</th>
<td id="S5.T2.12.14.2.2" class="ltx_td ltx_align_center">80.0</td>
<td id="S5.T2.12.14.2.3" class="ltx_td ltx_align_center">-</td>
</tr>
<tr id="S5.T2.3.3" class="ltx_tr">
<th id="S5.T2.2.2.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">MAD (<span class="math inline"><em>D</em> = 1</span>)</th>
<td id="S5.T2.3.3.2" class="ltx_td ltx_align_center ltx_border_t">84.5 <span class="math inline">±</span> 1.5</td>
<td id="S5.T2.3.3.3" class="ltx_td ltx_align_center ltx_border_t">baseline</td>
</tr>
<tr id="S5.T2.6.6" class="ltx_tr">
<th id="S5.T2.4.4.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">MAD (<span class="math inline"><em>D</em> = 4/5</span>)</th>
<td id="S5.T2.5.5.2" class="ltx_td ltx_align_center">83.5 <span class="math inline">±</span> 0.5</td>
<td id="S5.T2.6.6.3" class="ltx_td ltx_align_center"><span class="math inline">−</span>12.7%</td>
</tr>
<tr id="S5.T2.9.9" class="ltx_tr">
<th id="S5.T2.7.7.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">MAD (<span class="math inline"><em>D</em> = 3/5</span>)</th>
<td id="S5.T2.8.8.2" class="ltx_td ltx_align_center"><span id="S5.T2.8.8.2.1" class="ltx_text ltx_font_bold">86.5 <span class="math inline">±</span> 1.5</span></td>
<td id="S5.T2.9.9.3" class="ltx_td ltx_align_center"><span class="math inline">−</span>29.1%</td>
</tr>
<tr id="S5.T2.12.12" class="ltx_tr">
<th id="S5.T2.10.10.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_b ltx_border_r">MAD (<span class="math inline"><em>D</em> = 2/5</span>)</th>
<td id="S5.T2.11.11.2" class="ltx_td ltx_align_center ltx_border_b">84.5 <span class="math inline">±</span> 0.8</td>
<td id="S5.T2.12.12.3" class="ltx_td ltx_align_center ltx_border_b"><span class="math inline">−</span>43.6%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 2: </span>Comparison of accuracy and cost savings of MAD against baseline methods on the GSM8K dataset. All experiments were conducted using the GPT-3.5 model.</figcaption>
</figure>
:::::

::::: {#S5.SS2 .section .ltx_subsection}
### [5.2 ]{.ltx_tag .ltx_tag_subsection}MAD on Multimodal Reasoning Task {#mad-on-multimodal-reasoning-task .ltx_title .ltx_title_subsection}

::: {#S5.SS2.p1 .ltx_para}
MAD on multimodal reasoning tasks also demonstrates notable improvements compared to Chain-of-Thought and self-consistency approaches. This suggests that MLLMs like GPT-4o can effectively integrate step-by-step reasoning with visual content to enhance final answers. Similar to text reasoning experiments, we examine various sparse MAD configurations and report their performance.
:::

::: {#S5.SS2.p2 .ltx_para}
[Sparse MAD retains performance while introducing significant cost savings on multimodal reasoning tasks.]{#S5.SS2.p2.4.1 .ltx_text .ltx_font_bold} For the MathVista task, we evaluate different MAD configurations, comparing them to each other as well as to Chain-of-Thought (CoT) and self-consistency methods (Table [[3]{.ltx_text .ltx_ref_tag}](#S5.T3 "Table 3 ‣ 5.2 MAD on Multimodal Reasoning Task ‣ 5 Experiments: MAD with Single LLM ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref}). We find that sparse MAD achieves similar or slightly better accuracy compared to fully-connected MAD, with both outperforming CoT and self-consistency. The best sparse MAD configuration achieves a $+$`<!-- -->`{=html}1.2% improvement over fully-connected MAD and a $+$`<!-- -->`{=html}6.4% improvement over self-consistency. Additionally, sparse MAD provides substantial cost savings, reducing the total number of tokens used by up to 33.1%. Given that multimodal inputs are typically much larger than textual inputs (e.g., in GPT-4o, each image costs at least 225 tokens and can grow to 400$+$, 600$+$, or more tokens), we observe a total reduction of 40.6% in token usage, excluding the input image tokens.
:::

<figure id="S5.T3" class="ltx_table">
<table id="S5.T3.15" class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr id="S5.T3.15.16.1" class="ltx_tr">
<th id="S5.T3.15.16.1.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t"><span id="S5.T3.15.16.1.1.1" class="ltx_text ltx_font_bold">Method</span></th>
<td id="S5.T3.15.16.1.2" class="ltx_td ltx_align_center ltx_border_t">Accuracy</td>
<td id="S5.T3.15.16.1.3" class="ltx_td ltx_align_center ltx_border_t">Cost Saving</td>
</tr>
<tr id="S5.T3.1.1" class="ltx_tr">
<th id="S5.T3.1.1.2" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">CoT</th>
<td id="S5.T3.1.1.1" class="ltx_td ltx_align_center ltx_border_t">52.4 <span class="math inline">±</span> 2.6</td>
<td id="S5.T3.1.1.3" class="ltx_td ltx_align_center ltx_border_t">-</td>
</tr>
<tr id="S5.T3.15.17.2" class="ltx_tr">
<th id="S5.T3.15.17.2.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">SC</th>
<td id="S5.T3.15.17.2.2" class="ltx_td ltx_align_center">53.0</td>
<td id="S5.T3.15.17.2.3" class="ltx_td ltx_align_center">-</td>
</tr>
<tr id="S5.T3.3.3" class="ltx_tr">
<th id="S5.T3.2.2.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">MAD (<span class="math inline"><em>D</em> = 1</span>)</th>
<td id="S5.T3.3.3.2" class="ltx_td ltx_align_center ltx_border_t">58.2 <span class="math inline">±</span> 1.5</td>
<td id="S5.T3.3.3.3" class="ltx_td ltx_align_center ltx_border_t">baseline</td>
</tr>
<tr id="S5.T3.7.7" class="ltx_tr">
<th id="S5.T3.4.4.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">MAD (<span class="math inline"><em>D</em> = 4/5</span>)</th>
<td id="S5.T3.5.5.2" class="ltx_td ltx_align_center">57.8 <span class="math inline">±</span> 1.9</td>
<td id="S5.T3.7.7.4" class="ltx_td ltx_align_center"><span class="math inline">−</span>9.1% (<span class="math inline">−</span>11.5%)</td>
</tr>
<tr id="S5.T3.11.11" class="ltx_tr">
<th id="S5.T3.8.8.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">MAD (<span class="math inline"><em>D</em> = 3/5</span>)</th>
<td id="S5.T3.9.9.2" class="ltx_td ltx_align_center">55.4 <span class="math inline">±</span> 0.9</td>
<td id="S5.T3.11.11.4" class="ltx_td ltx_align_center"><span class="math inline">−</span>20.0% (<span class="math inline">−</span>24.7%)</td>
</tr>
<tr id="S5.T3.15.15" class="ltx_tr">
<th id="S5.T3.12.12.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_b ltx_border_r">MAD (<span class="math inline"><em>D</em> = 2/5</span>)</th>
<td id="S5.T3.13.13.2" class="ltx_td ltx_align_center ltx_border_b"><span id="S5.T3.13.13.2.1" class="ltx_text ltx_font_bold">59.4 <span class="math inline">±</span> 0.6</span></td>
<td id="S5.T3.15.15.4" class="ltx_td ltx_align_center ltx_border_b"><span class="math inline">−</span>33.1% (<span class="math inline">−</span>40.6%)</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 3: </span>Comparison of accuracy and cost savings of MAD against baseline methods on the MathVista dataset. All experiments were conducted using the GPT-4o model with the default temperature <span class="math inline"><em>T</em> = 1</span>. The cost saving percentages in parenthesis are computed without multimodal inputs.</figcaption>
</figure>
:::::

:::::::: {#S5.SS3 .section .ltx_subsection}
### [5.3 ]{.ltx_tag .ltx_tag_subsection}MAD on Alignment Labeling Tasks {#mad-on-alignment-labeling-tasks .ltx_title .ltx_title_subsection}

::: {#S5.SS3.p1 .ltx_para}
Alignment labeling tasks involve annotating preferences between pairs of responses generated for a given question. Our prompt consists of three parts: (1) a system prompt that informs the LLM of its role as a rater and specifies the required answer formatting; (2) a question description providing the context of the question; and (3) an ending instruction that constrains the answer length and reiterates the formatting requirements. During the debate, reference solutions are included before the ending instruction. See [[A]{.ltx_text .ltx_ref_tag}](#A1 "Appendix A Prompt Templates ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref} for more details.
:::

::: {#S5.SS3.p2 .ltx_para}
We use [AI Labeler Alignment]{#S5.SS3.p2.1.1 .ltx_text .ltx_font_italic} Lee et al. ([2023](#bib.bib17){.ltx_ref}) to measure the accuracy of MAD labeling with respect to the human annotation. To prevent potential position bias, we randomly assign the chosen response to either the (A) or (B) option. We report the accuracy and inference cost of MAD with various level of sparsity in Table [[4]{.ltx_text .ltx_ref_tag}](#S5.T4 "Table 4 ‣ 5.3 MAD on Alignment Labeling Tasks ‣ 5 Experiments: MAD with Single LLM ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref} for helpfulness and Table [[5]{.ltx_text .ltx_ref_tag}](#S5.T5 "Table 5 ‣ 5.3 MAD on Alignment Labeling Tasks ‣ 5 Experiments: MAD with Single LLM ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref} for harmlessness.
:::

<figure id="S5.T4" class="ltx_table">
<table id="S5.T4.20" class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr id="S5.T4.20.21.1" class="ltx_tr">
<th id="S5.T4.20.21.1.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t"><span id="S5.T4.20.21.1.1.1" class="ltx_text ltx_font_bold">Method</span></th>
<td colspan="2" id="S5.T4.20.21.1.2" class="ltx_td ltx_align_center ltx_border_r ltx_border_t"><span id="S5.T4.20.21.1.2.1" class="ltx_text ltx_font_bold">GPT-3.5</span></td>
<td colspan="2" id="S5.T4.20.21.1.3" class="ltx_td ltx_align_center ltx_border_t"><span id="S5.T4.20.21.1.3.1" class="ltx_text ltx_font_bold">Mistral 7B</span></td>
</tr>
<tr id="S5.T4.20.22.2" class="ltx_tr">
<th id="S5.T4.20.22.2.1" class="ltx_td ltx_th ltx_th_row ltx_border_r"></th>
<td id="S5.T4.20.22.2.2" class="ltx_td ltx_align_center">Accuracy</td>
<td id="S5.T4.20.22.2.3" class="ltx_td ltx_align_center ltx_border_r">Cost Saving</td>
<td id="S5.T4.20.22.2.4" class="ltx_td ltx_align_center">Accuracy</td>
<td id="S5.T4.20.22.2.5" class="ltx_td ltx_align_center">Cost Saving</td>
</tr>
<tr id="S5.T4.2.2" class="ltx_tr">
<th id="S5.T4.2.2.3" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">CoT</th>
<td id="S5.T4.1.1.1" class="ltx_td ltx_align_center ltx_border_t">56.5 <span class="math inline">±</span> 3.1</td>
<td id="S5.T4.2.2.4" class="ltx_td ltx_align_center ltx_border_r ltx_border_t">-</td>
<td id="S5.T4.2.2.2" class="ltx_td ltx_align_center ltx_border_t">60.8 <span class="math inline">±</span> 1.2</td>
<td id="S5.T4.2.2.5" class="ltx_td ltx_align_center ltx_border_t">-</td>
</tr>
<tr id="S5.T4.20.23.3" class="ltx_tr">
<th id="S5.T4.20.23.3.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Self-Consistency</th>
<td id="S5.T4.20.23.3.2" class="ltx_td ltx_align_center">57.0</td>
<td id="S5.T4.20.23.3.3" class="ltx_td ltx_align_center ltx_border_r">-</td>
<td id="S5.T4.20.23.3.4" class="ltx_td ltx_align_center">62.6</td>
<td id="S5.T4.20.23.3.5" class="ltx_td ltx_align_center">-</td>
</tr>
<tr id="S5.T4.5.5" class="ltx_tr">
<th id="S5.T4.3.3.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">MAD (<span class="math inline"><em>D</em> = 1</span>)</th>
<td id="S5.T4.4.4.2" class="ltx_td ltx_align_center ltx_border_t">58.5 <span class="math inline">±</span> 1.7</td>
<td id="S5.T4.5.5.4" class="ltx_td ltx_align_center ltx_border_r ltx_border_t">baseline</td>
<td id="S5.T4.5.5.3" class="ltx_td ltx_align_center ltx_border_t">65.5 <span class="math inline">±</span> 0.6</td>
<td id="S5.T4.5.5.5" class="ltx_td ltx_align_center ltx_border_t">baseline</td>
</tr>
<tr id="S5.T4.10.10" class="ltx_tr">
<th id="S5.T4.6.6.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">MAD (<span class="math inline"><em>D</em> = 4/5</span>)</th>
<td id="S5.T4.7.7.2" class="ltx_td ltx_align_center"><span id="S5.T4.7.7.2.1" class="ltx_text ltx_font_bold">59.0 <span class="math inline">±</span> 1.8</span></td>
<td id="S5.T4.8.8.3" class="ltx_td ltx_align_center ltx_border_r"><span class="math inline">−</span>17.5%</td>
<td id="S5.T4.9.9.4" class="ltx_td ltx_align_center">65.6 <span class="math inline">±</span> 0.9</td>
<td id="S5.T4.10.10.5" class="ltx_td ltx_align_center"><span class="math inline">−</span>18.3%</td>
</tr>
<tr id="S5.T4.15.15" class="ltx_tr">
<th id="S5.T4.11.11.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">MAD (<span class="math inline"><em>D</em> = 3/5</span>)</th>
<td id="S5.T4.12.12.2" class="ltx_td ltx_align_center">57.0 <span class="math inline">±</span> 1.6</td>
<td id="S5.T4.13.13.3" class="ltx_td ltx_align_center ltx_border_r"><span class="math inline">−</span>32.5%</td>
<td id="S5.T4.14.14.4" class="ltx_td ltx_align_center">64.6 <span class="math inline">±</span> 0.6</td>
<td id="S5.T4.15.15.5" class="ltx_td ltx_align_center"><span class="math inline">−</span>35.2%</td>
</tr>
<tr id="S5.T4.20.20" class="ltx_tr">
<th id="S5.T4.16.16.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_b ltx_border_r">MAD (<span class="math inline"><em>D</em> = 2/5</span>)</th>
<td id="S5.T4.17.17.2" class="ltx_td ltx_align_center ltx_border_b"><span id="S5.T4.17.17.2.1" class="ltx_text ltx_font_bold">59.0 <span class="math inline">±</span> 1.4</span></td>
<td id="S5.T4.18.18.3" class="ltx_td ltx_align_center ltx_border_b ltx_border_r"><span class="math inline">−</span>50.0%</td>
<td id="S5.T4.19.19.4" class="ltx_td ltx_align_center ltx_border_b"><span id="S5.T4.19.19.4.1" class="ltx_text ltx_font_bold">66.6 <span class="math inline">±</span> 0.5</span></td>
<td id="S5.T4.20.20.5" class="ltx_td ltx_align_center ltx_border_b"><span class="math inline">−</span>53.5%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 4: </span>AI labeler alignment accuracy and cost savings of MAD compared with baselines on the helpfulness dataset for GPT-3.5 and Mistral 7B models.</figcaption>
</figure>

<figure id="S5.T5" class="ltx_table">
<table id="S5.T5.20" class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr id="S5.T5.20.21.1" class="ltx_tr">
<th id="S5.T5.20.21.1.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t"><span id="S5.T5.20.21.1.1.1" class="ltx_text ltx_font_bold">Method</span></th>
<td colspan="2" id="S5.T5.20.21.1.2" class="ltx_td ltx_align_center ltx_border_r ltx_border_t"><span id="S5.T5.20.21.1.2.1" class="ltx_text ltx_font_bold">GPT-3.5</span></td>
<td colspan="2" id="S5.T5.20.21.1.3" class="ltx_td ltx_align_center ltx_border_t"><span id="S5.T5.20.21.1.3.1" class="ltx_text ltx_font_bold">Mistral 7B</span></td>
</tr>
<tr id="S5.T5.20.22.2" class="ltx_tr">
<th id="S5.T5.20.22.2.1" class="ltx_td ltx_th ltx_th_row ltx_border_r"></th>
<td id="S5.T5.20.22.2.2" class="ltx_td ltx_align_center">Accuracy</td>
<td id="S5.T5.20.22.2.3" class="ltx_td ltx_align_center ltx_border_r">Cost Saving</td>
<td id="S5.T5.20.22.2.4" class="ltx_td ltx_align_center">Accuracy</td>
<td id="S5.T5.20.22.2.5" class="ltx_td ltx_align_center">Cost Saving</td>
</tr>
<tr id="S5.T5.2.2" class="ltx_tr">
<th id="S5.T5.2.2.3" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">CoT</th>
<td id="S5.T5.1.1.1" class="ltx_td ltx_align_center ltx_border_t">66.0 <span class="math inline">±</span> 4.8</td>
<td id="S5.T5.2.2.4" class="ltx_td ltx_align_center ltx_border_r ltx_border_t">-</td>
<td id="S5.T5.2.2.2" class="ltx_td ltx_align_center ltx_border_t">58.2 <span class="math inline">±</span> 2.0</td>
<td id="S5.T5.2.2.5" class="ltx_td ltx_align_center ltx_border_t">-</td>
</tr>
<tr id="S5.T5.20.23.3" class="ltx_tr">
<th id="S5.T5.20.23.3.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">Self-Consistency</th>
<td id="S5.T5.20.23.3.2" class="ltx_td ltx_align_center">67.0</td>
<td id="S5.T5.20.23.3.3" class="ltx_td ltx_align_center ltx_border_r">-</td>
<td id="S5.T5.20.23.3.4" class="ltx_td ltx_align_center">60.0</td>
<td id="S5.T5.20.23.3.5" class="ltx_td ltx_align_center">-</td>
</tr>
<tr id="S5.T5.5.5" class="ltx_tr">
<th id="S5.T5.3.3.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">MAD (<span class="math inline"><em>D</em> = 1</span>)</th>
<td id="S5.T5.4.4.2" class="ltx_td ltx_align_center ltx_border_t">67.5 <span class="math inline">±</span> 0.6</td>
<td id="S5.T5.5.5.4" class="ltx_td ltx_align_center ltx_border_r ltx_border_t">baseline</td>
<td id="S5.T5.5.5.3" class="ltx_td ltx_align_center ltx_border_t">60.7 <span class="math inline">±</span> 0.3</td>
<td id="S5.T5.5.5.5" class="ltx_td ltx_align_center ltx_border_t">baseline</td>
</tr>
<tr id="S5.T5.10.10" class="ltx_tr">
<th id="S5.T5.6.6.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">MAD (<span class="math inline"><em>D</em> = 4/5</span>)</th>
<td id="S5.T5.7.7.2" class="ltx_td ltx_align_center">67.0 <span class="math inline">±</span> 0.8</td>
<td id="S5.T5.8.8.3" class="ltx_td ltx_align_center ltx_border_r"><span class="math inline">−</span>17.3%</td>
<td id="S5.T5.9.9.4" class="ltx_td ltx_align_center"><span id="S5.T5.9.9.4.1" class="ltx_text ltx_font_bold">62.2 <span class="math inline">±</span> 0.2</span></td>
<td id="S5.T5.10.10.5" class="ltx_td ltx_align_center"><span class="math inline">−</span>17.9%</td>
</tr>
<tr id="S5.T5.15.15" class="ltx_tr">
<th id="S5.T5.11.11.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">MAD (<span class="math inline"><em>D</em> = 3/5</span>)</th>
<td id="S5.T5.12.12.2" class="ltx_td ltx_align_center">67.5 <span class="math inline">±</span> 1.0</td>
<td id="S5.T5.13.13.3" class="ltx_td ltx_align_center ltx_border_r"><span class="math inline">−</span>34.7%</td>
<td id="S5.T5.14.14.4" class="ltx_td ltx_align_center">60.4 <span class="math inline">±</span> 0.4</td>
<td id="S5.T5.15.15.5" class="ltx_td ltx_align_center"><span class="math inline">−</span>34.3%</td>
</tr>
<tr id="S5.T5.20.20" class="ltx_tr">
<th id="S5.T5.16.16.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_b ltx_border_r">MAD (<span class="math inline"><em>D</em> = 2/5</span>)</th>
<td id="S5.T5.17.17.2" class="ltx_td ltx_align_center ltx_border_b"><span id="S5.T5.17.17.2.1" class="ltx_text ltx_font_bold">68.5 <span class="math inline">±</span> 0.7</span></td>
<td id="S5.T5.18.18.3" class="ltx_td ltx_align_center ltx_border_b ltx_border_r"><span class="math inline">−</span>53.3%</td>
<td id="S5.T5.19.19.4" class="ltx_td ltx_align_center ltx_border_b">61.7 <span class="math inline">±</span> 0.2</td>
<td id="S5.T5.20.20.5" class="ltx_td ltx_align_center ltx_border_b"><span class="math inline">−</span>52.2%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 5: </span>AI labeler alignment accuracy and cost savings of MAD compared with baselines on the harmlessness dataset for GPT-3.5 and Mistral 7B models.</figcaption>
</figure>

::: {#S5.SS3.p3 .ltx_para}
[MAD outperforms single-agent on alignment labeling tasks]{#S5.SS3.p3.4.1 .ltx_text .ltx_font_bold}: We find that MAD consistently outperforms single-agent methods, including CoT and self-consistency. On the helpfulness task, fully-connected MAD achieves a $+ {1.5\%}$ and $+ {2.9\%}$ improvement over self-consistency for GPT-3.5 and Mistral 7B models, respectively. On the harmlessness task, fully-connected MAD achieves a $+ {0.5\%}$ and $+ {0.7\%}$ improvement over self-consistency for GPT-3.5 and Mistral 7B models, respectively. These results suggest that the additional debate process in MAD, followed by majority voting, allows agents to incorporate perspectives from others and refine their opinions toward the correct answers during the debate process.
:::

::: {#S5.SS3.p4 .ltx_para}
[Sparse MAD can perform better with lower inference costs]{#S5.SS3.p4.6.1 .ltx_text .ltx_font_bold}: Most sparse MAD configurations perform as well as or better than the fully-connected MAD, with at least one sparse topology outperforming the fully-connected MAD. Depending on the task, sparse MAD with GPT-3.5 can enhance performance by approximately $+ {0.5\%}$ to $+ {1.0\%}$, and sparse MAD with Mistral 7B can improve performance by about $+ {1.1\%}$ to $+ {1.5\%}$. Additionally, sparse MAD can reduce costs by up to $- {53.3\%}$ and $- {53.5\%}$, respectively.
:::

::: {#S5.SS3.p5 .ltx_para}
We observed that GPT-3.5 exhibits lower alignment accuracy compared to Mistral 7B on the helpfulness task. We attribute this discrepancy to the differences in pre-training and post-training corpora between the two models, which may lead to varying default preferences in a zero-shot setting. While we hypothesize that few-shot prompting techniques could mitigate this issue, exploring this is beyond the scope of this work.
:::
::::::::

:::::: {#S5.SS4 .section .ltx_subsection}
### [5.4 ]{.ltx_tag .ltx_tag_subsection}Why Does Sparse MAD Work? {#why-does-sparse-mad-work .ltx_title .ltx_title_subsection}

::: {#S5.SS4.p1 .ltx_para}
The common explanation for the effectiveness of MAD against single-agent setups is that agents can consider different perspectives before arriving at an answer. However, our experiment on the effectiveness of sparse MAD seems challenge this intuition. In this section, we aim to explain why sparse MAD can achieve comparable or even superior performance.
:::

<figure id="S5.F3" class="ltx_figure">
<img src="/html/2406.11776/assets/imgs/analysis_trend_v3.png" id="S5.F3.g1" class="ltx_graphics ltx_img_square" width="598" height="581" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 3: </span>Probability of a single agent generating correct answers given <span class="math inline"><em>n</em></span> reference solutions, with <span class="math inline"><em>p</em></span> representing the correctness of these solutions. Monte Carlo sampling was performed on three questions, each with 100 runs.</figcaption>
</figure>

<figure id="S5.F4" class="ltx_figure">
<img src="/html/2406.11776/assets/imgs/num_rounds_corrected.png" id="S5.F4.g1" class="ltx_graphics ltx_img_landscape" width="598" height="368" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 4: </span>Effective debate rounds for each topology design in reasoning and alignment labeling tasks.</figcaption>
</figure>

::: {#S5.SS4.p2 .ltx_para}
[Impact of incorrect reference solutions]{#S5.SS4.p2.9.1 .ltx_text .ltx_font_bold}: In a MAD framework, we define $Q\hspace{0pt}{(n,p)}$ as the probability that a single agent delivers correct answers given $n$ reference solutions, where $p$ percentage of these are correct. This probability, $Q\hspace{0pt}{(n,p)}$, can be estimated using Monte Carlo sampling with constructed in-context reference solutions. As a case study, we selected three questions from the GSM8K dataset and estimated $Q\hspace{0pt}{(n,p)}$ for $n \in {\{ 2,3,4,5\}}$ and $p \in {\{{0\%},{25\%},{50\%},{75\%},{100\%}\}}$. Here, the choice of $n$ corresponds to the single-agent scenarios in MAD with $D = {\frac{2}{5},\frac{3}{5},\frac{4}{5},1}$. Results shown in Figure [[3]{.ltx_text .ltx_ref_tag}](#S5.F3 "Figure 3 ‣ 5.4 Why Does Sparse MAD Work? ‣ 5 Experiments: MAD with Single LLM ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref} indicate that for easier questions, where most reference solutions are correct, an increase in the number of observed reference solutions (namely MAD becomes denser) improves the likelihood of the agent arriving at the correct answer. Conversely, for more difficult questions, where most agents do not provide correct answers, an increase in the number of observed reference solutions tends to mislead the agent into choosing incorrect answers, thereby drastically reducing the likelihood of reaching a correct response.
:::

::: {#S5.SS4.p3 .ltx_para}
[Sparser MAD allows more rounds of effective debate]{#S5.SS4.p3.1.1 .ltx_text .ltx_font_bold}: We observe that once all agents converge on the same answer, it becomes highly unlikely for any of them to change their decision. We define the number of effective debates as the number of rounds before all agents reach the same answer. Figure [[4]{.ltx_text .ltx_ref_tag}](#S5.F4 "Figure 4 ‣ 5.4 Why Does Sparse MAD Work? ‣ 5 Experiments: MAD with Single LLM ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref} illustrates the effective number of debate rounds for various topologies in reasoning and alignment labeling tasks. Our results show that sparse MAD tends to sustain longer debates before achieving consensus, indicating that sparse MAD allows for more extensive deliberation and in-depth discussion. We observe there are similar findings in the Chain-of-Thought prompting Jin et al. ([2024](#bib.bib16){.ltx_ref}) and MAD Du et al. ([2023](#bib.bib9){.ltx_ref}) that the increase of reasoning length can significantly improve the performance.
:::
::::::
:::::::::::::::::::

:::::: {#S6 .section .ltx_section}
## [6 ]{.ltx_tag .ltx_tag_section}Experiments: MAD with Multiple LLMs {#experiments-mad-with-multiple-llms .ltx_title .ltx_title_section}

::: {#S6.p1 .ltx_para}
Previous sections focus on the MAD with agents instantiated by the same LLM. In this section, we explore the scenario when multiple LLMs are available. With agents instantiated by different LLMs, the permutation invariance symmetry is broken, and the regular graph may not be optimal. A natural question is: [how to design the communication topology given a MAD framework of $N$ agents, in which $M$ instantiated by the stronger LLM and $N - M$ instantiated by the weaker LLM?]{#S6.p1.3.3 .ltx_text .ltx_font_italic}
:::

<figure id="S6.F5" class="ltx_figure">
<img src="/html/2406.11776/assets/imgs/A5_new.png" id="S6.F5.g1" class="ltx_graphics ltx_img_landscape" width="598" height="138" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 5: </span>Isotropic communication topology with two setups: the stronger LLM has low centrality (left) and high centrality (right).</figcaption>
</figure>

<figure id="S6.T6" class="ltx_table">
<table id="S6.T6.2" class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<thead class="ltx_thead">
<tr id="S6.T6.2.3.1" class="ltx_tr">
<th id="S6.T6.2.3.1.1" class="ltx_td ltx_align_left ltx_th ltx_th_column ltx_th_row ltx_border_r ltx_border_t"><span id="S6.T6.2.3.1.1.1" class="ltx_text ltx_font_bold">Centrality</span></th>
<th colspan="2" id="S6.T6.2.3.1.2" class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_t"><span id="S6.T6.2.3.1.2.1" class="ltx_text ltx_font_bold">Accuracy</span></th>
</tr>
<tr id="S6.T6.2.4.2" class="ltx_tr">
<th id="S6.T6.2.4.2.1" class="ltx_td ltx_th ltx_th_column ltx_th_row ltx_border_r"></th>
<th id="S6.T6.2.4.2.2" class="ltx_td ltx_align_center ltx_th ltx_th_column">SC</th>
<th id="S6.T6.2.4.2.3" class="ltx_td ltx_align_center ltx_th ltx_th_column">Isotropic MAD</th>
</tr>
</thead>
<tbody class="ltx_tbody">
<tr id="S6.T6.1.1" class="ltx_tr">
<th id="S6.T6.1.1.2" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">High</th>
<td id="S6.T6.1.1.3" class="ltx_td ltx_align_center ltx_border_t">64.0</td>
<td id="S6.T6.1.1.1" class="ltx_td ltx_align_center ltx_border_t"><span id="S6.T6.1.1.1.1" class="ltx_text ltx_font_bold">67.0 <span class="math inline">±</span> 0.8</span></td>
</tr>
<tr id="S6.T6.2.2" class="ltx_tr">
<th id="S6.T6.2.2.2" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_b ltx_border_r">Low</th>
<td id="S6.T6.2.2.3" class="ltx_td ltx_align_center ltx_border_b">64.0</td>
<td id="S6.T6.2.2.1" class="ltx_td ltx_align_center ltx_border_b">65.8 <span class="math inline">±</span> 0.5</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 6: </span>Comparison of accuracy depending on where a stronger LLM is placed in debate, using the Harmlessness task as example. In both cases, there are 5 Mistral models and 1 GPT-3.5 Model. Accuracy of Isotropic MAD is calculated as the average over debate rounds.</figcaption>
</figure>

<figure id="S6.F6" class="ltx_figure">
<img src="/html/2406.11776/assets/imgs/stronger_llm_font.png" id="S6.F6.g1" class="ltx_graphics ltx_img_landscape" width="598" height="396" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 6: </span>Average accuracy of weaker agents across different debate rounds.</figcaption>
</figure>

::: {#S6.p2 .ltx_para}
[Assigning stronger LLMs to agents with higher centrality yields better performance]{#S6.p2.2.1 .ltx_text .ltx_font_bold}: We conducted experiments on harmlessness alignment labeling task, involving 6 agents, with 1 agent utilizing GPT-3.5 (the stronger LLM) and the remaining 5 agents utilizing Mistral 7B (the weaker LLM). We tested two setups on the isotropic communication topology: one where the stronger LLM had a degree of 1 (indicating low centrality) and another where it had a degree of 5 (indicating high centrality), as illustrated in Figure [[5]{.ltx_text .ltx_ref_tag}](#S6.F5 "Figure 5 ‣ 6 Experiments: MAD with Multiple LLMs ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref}. The experimental results presented in Table [[6]{.ltx_text .ltx_ref_tag}](#S6.T6 "Table 6 ‣ 6 Experiments: MAD with Multiple LLMs ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref} show that positioning the stronger LLM at a node with higher centrality (degree of 5) leads to better performance ($+ {3.0\%}$ improvement) compared to placing it at a node with lower centrality (degree of 1) which resulted in a $+ {1.8\%}$ improvement.
:::

::: {#S6.p3 .ltx_para}
The results above underscore the importance of information flow in the design of communication topology. Figure [[6]{.ltx_text .ltx_ref_tag}](#S6.F6 "Figure 6 ‣ 6 Experiments: MAD with Multiple LLMs ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref} illustrates the average accuracy of weaker agents with respect to the number of debate rounds. When the stronger agent has a degree of 5, it can effectively disseminate its knowledge to weaker agents in just one debate round, resulting in a sharp increase in the average accuracy of weaker LLMs. In contrast, when the stronger agent has a degree of 1, the process requires two rounds: first, the information is transmitted to the central weaker agent in the first debate round (round 2), which then shares it with other weaker agents in the next round (round 3). This two-step process leads to greater information loss.
:::
::::::

::::: {#S7 .section .ltx_section}
## [7 ]{.ltx_tag .ltx_tag_section}Conclusion {#conclusion .ltx_title .ltx_title_section}

::: {#S7.p1 .ltx_para}
In this paper, we show that sparse communication topologies can improve the multi-agent debate framework significantly. Our results indicate that sparse MAD configurations achieve comparable or superior performance to standard MADs while greatly reducing computational costs. We also extend the MAD framework to alignment labeling tasks, demonstrating the benefits of MADs over single-agent setups and self-consistency and further highlighting the benefits of sparse MADs over fully-connected configurations. We present case-study insights that explain the effectiveness of sparse MADs. Additionally, we investigate the impact of communication topology design with multiple large language models (LLMs), finding that assigning stronger LLMs to more connected agents enhances overall performance.
:::

::: {#S7.p2 .ltx_para}
In summary, our work paves the way for more efficient and effective multi-agent systems by leveraging sparse communication topologies. Future studies could focus on deepening our understanding of the underlying mechanisms and developing strategies for optimal topology design in multi-agent frameworks.
:::
:::::

::::::: {#S8 .section .ltx_section}
## [8 ]{.ltx_tag .ltx_tag_section}Ethical Considerations {#ethical-considerations .ltx_title .ltx_title_section}

::: {#S8.p1 .ltx_para}
In this work, several ethical considerations were addressed to ensure the integrity and responsible use of the system:
:::

::: {#S8.p2 .ltx_para .ltx_noindent}
[Public Datasets]{#S8.p2.1.1 .ltx_text .ltx_font_bold}: The framework was built using publicly available datasets that are designed for academic research. We strictly adhered to ethical guidelines by not using any personal or confidential data.
:::

::: {#S8.p3 .ltx_para .ltx_noindent}
[License]{#S8.p3.1.1 .ltx_text .ltx_font_bold}: Only public APIs that offer appropriate licensing were utilized. This ensures that all external tools are used in a lawful and ethical manner.
:::

::: {#S8.p4 .ltx_para .ltx_noindent}
[AI assistant]{#S8.p4.1.1 .ltx_text .ltx_font_bold}: AI tools were employed solely for polishing writing and correcting grammar. The AI was not used to generate content or ideas, maintaining the authenticity and originality of the research work.
:::
:::::::

:::::::: {#S9 .section .ltx_section}
## [9 ]{.ltx_tag .ltx_tag_section}Limitations {#limitations .ltx_title .ltx_title_section}

::: {#S9.p1 .ltx_para}
While our study provides valuable insights into the communication topology analysis of multi-agent debate, several limitations must be acknowledged:
:::

::: {#S9.p2 .ltx_para}
Our analysis is primarily based on static graphs where the communication topology remains unchanged throughout the debate rounds. This constraint simplifies the analysis, but ignores the dynamic nature of real-world communication networks. Additionally, our study focuses on prompt design under a zero-shot setting, utilizing only publicly available GPT and Mistral models. This narrow scope may not fully capture the variability and adaptability present in more diverse agent populations. Furthermore, we confined our analysis to regular graphs, which do not encompass the full spectrum of potential graph configurations. Future work should consider dynamic graphs, a broader range of models, and varied graph connectivity to better reflect the evolving and complex nature of multi-agent interactions.
:::

::: {#S9.p3 .ltx_para}
Our study relies on a subset of academic datasets due to limited data access as well as computational constraints. While these datasets provide a valuable foundation for analyzing communication graph dynamics in multi-agent debates, they may not fully represent the diversity and complexity found in broader real-world data. The restricted scope limits our ability to generalize findings across different domains and contexts. Future research should aim to include a wider range of datasets, potentially leveraging more efficient computational resources, to enhance the robustness and applicability of our findings.
:::

::: {#S9.p4 .ltx_para}
We lack a rigorous theoretical proof explaining why sparse connectivity can lead to better performance. This gap in our understanding limits our ability to generalize our findings and apply them with confidence in various settings. Secondly, we do not have a definitive method for determining the optimal topology design, which is crucial for maximizing the efficiency and effectiveness of multi-agent systems. Addressing these questions is essential for future research. Potential explanations might involve theoretical insights, social and psychological dynamics, or a combination of these factors. Additionally, fine-tuning models could offer further clarity and aid in optimizing communication topology. Future work should aim to develop robust theoretical frameworks and empirical strategies to better understand and leverage communication topology in multi-agent debates.
:::

::: {#S9.p5 .ltx_para}
The multi-agent debate framework holds significant potential for various real-world applications. However, it also carries the risk of misuse, including the dissemination of biased information or misinformation. Additionally, the framework requires substantial computational resources, which could impact energy consumption and environmental sustainability. Future research should focus on developing robust, trustworthy, and energy-efficient multi-agent systems to mitigate these risks and ensure ethical, reliable, and sustainable outcomes.
:::
::::::::

::: {#bib .section .ltx_bibliography}
## References {#references .ltx_title .ltx_title_bibliography}

- [[Bai et al. (2022a)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yuntao Bai, Andy Jones, Kamal Ndousse, Amanda Askell, Anna Chen, Nova DasSarma, Dawn Drain, Stanislav Fort, Deep Ganguli, Tom Henighan, et al. 2022a. ]{.ltx_bibblock} [Training a helpful and harmless assistant with reinforcement learning from human feedback. ]{.ltx_bibblock} [*arXiv preprint arXiv:2204.05862*. ]{.ltx_bibblock}]{#bib.bib1}
- [[Bai et al. (2022b)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yuntao Bai, Saurav Kadavath, Sandipan Kundu, Amanda Askell, Jackson Kernion, Andy Jones, Anna Chen, Anna Goldie, Azalia Mirhoseini, Cameron McKinnon, et al. 2022b. ]{.ltx_bibblock} [Constitutional ai: Harmlessness from ai feedback. ]{.ltx_bibblock} [*arXiv preprint arXiv:2212.08073*. ]{.ltx_bibblock}]{#bib.bib2}
- [[Besta et al. (2024)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Maciej Besta, Nils Blach, Ales Kubicek, Robert Gerstenberger, Michal Podstawski, Lukas Gianinazzi, Joanna Gajda, Tomasz Lehmann, Hubert Niewiadomski, Piotr Nyczyk, et al. 2024. ]{.ltx_bibblock} [Graph of thoughts: Solving elaborate problems with large language models. ]{.ltx_bibblock} [In *Proceedings of the AAAI Conference on Artificial Intelligence*, volume 38, pages 17682--17690. ]{.ltx_bibblock}]{#bib.bib3}
- [[Burns et al. (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Collin Burns, Pavel Izmailov, Jan Hendrik Kirchner, Bowen Baker, Leo Gao, Leopold Aschenbrenner, Yining Chen, Adrien Ecoffet, Manas Joglekar, Jan Leike, et al. 2023. ]{.ltx_bibblock} [Weak-to-strong generalization: Eliciting strong capabilities with weak supervision. ]{.ltx_bibblock} [*arXiv preprint arXiv:2312.09390*. ]{.ltx_bibblock}]{#bib.bib4}
- [[Chan et al. (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Chi-Min Chan, Weize Chen, Yusheng Su, Jianxuan Yu, Wei Xue, Shanghang Zhang, Jie Fu, and Zhiyuan Liu. 2023. ]{.ltx_bibblock} [Chateval: Towards better llm-based evaluators through multi-agent debate. ]{.ltx_bibblock} [In *The Twelfth International Conference on Learning Representations*. ]{.ltx_bibblock}]{#bib.bib5}
- [[Chen et al. (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Justin Chih-Yao Chen, Swarnadeep Saha, and Mohit Bansal. 2023. ]{.ltx_bibblock} [Reconcile: Round-table conference improves reasoning via consensus among diverse llms. ]{.ltx_bibblock} [*arXiv preprint arXiv:2309.13007*. ]{.ltx_bibblock}]{#bib.bib6}
- [[Chen et al. (2024)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Liangyu Chen, Bo Li, Sheng Shen, Jingkang Yang, Chunyuan Li, Kurt Keutzer, Trevor Darrell, and Ziwei Liu. 2024. ]{.ltx_bibblock} [Large language models are visual reasoning coordinators. ]{.ltx_bibblock} [*Advances in Neural Information Processing Systems*, 36. ]{.ltx_bibblock}]{#bib.bib7}
- [[Cobbe et al. (2021)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Karl Cobbe, Vineet Kosaraju, Mohammad Bavarian, Mark Chen, Heewoo Jun, Lukasz Kaiser, Matthias Plappert, Jerry Tworek, Jacob Hilton, Reiichiro Nakano, et al. 2021. ]{.ltx_bibblock} [Training verifiers to solve math word problems, 2021. ]{.ltx_bibblock} [*URL https://arxiv. org/abs/2110.14168*. ]{.ltx_bibblock}]{#bib.bib8}
- [[Du et al. (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yilun Du, Shuang Li, Antonio Torralba, Joshua B Tenenbaum, and Igor Mordatch. 2023. ]{.ltx_bibblock} [Improving factuality and reasoning in language models through multiagent debate. ]{.ltx_bibblock} [*arXiv preprint arXiv:2305.14325*. ]{.ltx_bibblock}]{#bib.bib9}
- [[Ganz et al. (2024)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Roy Ganz, Yair Kittenplon, Aviad Aberdam, Elad Ben Avraham, Oren Nuriel, Shai Mazor, and Ron Litman. 2024. ]{.ltx_bibblock} [Question aware vision transformer for multimodal reasoning. ]{.ltx_bibblock} [In *Proceedings of the IEEE/CVF Conference on Computer Vision and Pattern Recognition*, pages 13861--13871. ]{.ltx_bibblock}]{#bib.bib10}
- [[Goyal et al. (2017)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yash Goyal, Tejas Khot, Douglas Summers-Stay, Dhruv Batra, and Devi Parikh. 2017. ]{.ltx_bibblock} [Making the V in VQA matter: Elevating the role of image understanding in Visual Question Answering. ]{.ltx_bibblock} [In *Conference on Computer Vision and Pattern Recognition (CVPR)*. ]{.ltx_bibblock}]{#bib.bib11}
- [[Guo et al. (2024)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Shangmin Guo, Biao Zhang, Tianlin Liu, Tianqi Liu, Misha Khalman, Felipe Llinares, Alexandre Rame, Thomas Mesnard, Yao Zhao, Bilal Piot, et al. 2024. ]{.ltx_bibblock} [Direct language model alignment from online ai feedback. ]{.ltx_bibblock} [*arXiv preprint arXiv:2402.04792*. ]{.ltx_bibblock}]{#bib.bib12}
- [[Hendrycks et al. (2021)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Dan Hendrycks, Collin Burns, Saurav Kadavath, Akul Arora, Steven Basart, Eric Tang, Dawn Song, and Jacob Steinhardt. 2021. ]{.ltx_bibblock} [Measuring mathematical problem solving with the math dataset. ]{.ltx_bibblock} [In *Thirty-fifth Conference on Neural Information Processing Systems Datasets and Benchmarks Track (Round 2)*. ]{.ltx_bibblock}]{#bib.bib13}
- [[Hu et al. (2024)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yushi Hu, Weijia Shi, Xingyu Fu, Dan Roth, Mari Ostendorf, Luke Zettlemoyer, Noah A Smith, and Ranjay Krishna. 2024. ]{.ltx_bibblock} [[Visual sketchpad: Sketching as a visual chain of thought for multimodal language models](https://api.semanticscholar.org/CorpusID:270440440){.ltx_ref .ltx_href target="_blank"}. ]{.ltx_bibblock}]{#bib.bib14}
- [[Jiang et al. (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Albert Q Jiang, Alexandre Sablayrolles, Arthur Mensch, Chris Bamford, Devendra Singh Chaplot, Diego de las Casas, Florian Bressand, Gianna Lengyel, Guillaume Lample, Lucile Saulnier, et al. 2023. ]{.ltx_bibblock} [Mistral 7b. ]{.ltx_bibblock} [*arXiv preprint arXiv:2310.06825*. ]{.ltx_bibblock}]{#bib.bib15}
- [[Jin et al. (2024)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Mingyu Jin, Qinkai Yu, Haiyan Zhao, Wenyue Hua, Yanda Meng, Yongfeng Zhang, Mengnan Du, et al. 2024. ]{.ltx_bibblock} [The impact of reasoning step length on large language models. ]{.ltx_bibblock} [*arXiv preprint arXiv:2401.04925*. ]{.ltx_bibblock}]{#bib.bib16}
- [[Lee et al. (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Harrison Lee, Samrat Phatale, Hassan Mansoor, Kellie Lu, Thomas Mesnard, Colton Bishop, Victor Carbune, and Abhinav Rastogi. 2023. ]{.ltx_bibblock} [Rlaif: Scaling reinforcement learning from human feedback with ai feedback. ]{.ltx_bibblock} [*arXiv preprint arXiv:2309.00267*. ]{.ltx_bibblock}]{#bib.bib17}
- [[Li et al. (2023a)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Guohao Li, Hasan Hammoud, Hani Itani, Dmitrii Khizbullin, and Bernard Ghanem. 2023a. ]{.ltx_bibblock} [Camel: Communicative agents for\" mind\" exploration of large language model society. ]{.ltx_bibblock} [*Advances in Neural Information Processing Systems*, 36:51991--52008. ]{.ltx_bibblock}]{#bib.bib18}
- [[Li et al. (2023b)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Junnan Li, Dongxu Li, Silvio Savarese, and Steven Hoi. 2023b. ]{.ltx_bibblock} [Blip-2: Bootstrapping language-image pre-training with frozen image encoders and large language models. ]{.ltx_bibblock} [In *International conference on machine learning*, pages 19730--19742. PMLR. ]{.ltx_bibblock}]{#bib.bib19}
- [[Liang et al. (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tian Liang, Zhiwei He, Wenxiang Jiao, Xing Wang, Yan Wang, Rui Wang, Yujiu Yang, Zhaopeng Tu, and Shuming Shi. 2023. ]{.ltx_bibblock} [Encouraging divergent thinking in large language models through multi-agent debate. ]{.ltx_bibblock} [*arXiv preprint arXiv:2305.19118*. ]{.ltx_bibblock}]{#bib.bib20}
- [[Lin et al. (2024)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Bin Lin, Zhenyu Tang, Yang Ye, Jiaxi Cui, Bin Zhu, Peng Jin, Junwu Zhang, Munan Ning, and Li Yuan. 2024. ]{.ltx_bibblock} [Moe-llava: Mixture of experts for large vision-language models. ]{.ltx_bibblock} [*arXiv preprint arXiv:2401.15947*. ]{.ltx_bibblock}]{#bib.bib21}
- [[Liu et al. (2024)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Haotian Liu, Chunyuan Li, Yuheng Li, and Yong Jae Lee. 2024. ]{.ltx_bibblock} [Improved baselines with visual instruction tuning. ]{.ltx_bibblock} [In *Proceedings of the IEEE/CVF Conference on Computer Vision and Pattern Recognition*, pages 26296--26306. ]{.ltx_bibblock}]{#bib.bib22}
- [[Lu et al. (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Pan Lu, Hritik Bansal, Tony Xia, Jiacheng Liu, Chunyuan Li, Hannaneh Hajishirzi, Hao Cheng, Kai-Wei Chang, Michel Galley, and Jianfeng Gao. 2023. ]{.ltx_bibblock} [Mathvista: Evaluating mathematical reasoning of foundation models in visual contexts. ]{.ltx_bibblock} [*arXiv preprint arXiv:2310.02255*. ]{.ltx_bibblock}]{#bib.bib23}
- [[Lu et al. (2022)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Pan Lu, Swaroop Mishra, Tony Xia, Liang Qiu, Kai-Wei Chang, Song-Chun Zhu, Oyvind Tafjord, Peter Clark, and Ashwin Kalyan. 2022. ]{.ltx_bibblock} [Learn to explain: Multimodal reasoning via thought chains for science question answering. ]{.ltx_bibblock} [In *The 36th Conference on Neural Information Processing Systems (NeurIPS)*. ]{.ltx_bibblock}]{#bib.bib24}
- [[Luo et al. (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Liangchen Luo, Zi Lin, Yinxiao Liu, Lei Shu, Yun Zhu, Jingbo Shang, and Lei Meng. 2023. ]{.ltx_bibblock} [Critique ability of large language models. ]{.ltx_bibblock} [*arXiv preprint arXiv:2310.04815*. ]{.ltx_bibblock}]{#bib.bib25}
- [[Madaan et al. (2024)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Aman Madaan, Niket Tandon, Prakhar Gupta, Skyler Hallinan, Luyu Gao, Sarah Wiegreffe, Uri Alon, Nouha Dziri, Shrimai Prabhumoye, Yiming Yang, et al. 2024. ]{.ltx_bibblock} [Self-refine: Iterative refinement with self-feedback. ]{.ltx_bibblock} [*Advances in Neural Information Processing Systems*, 36. ]{.ltx_bibblock}]{#bib.bib26}
- [[Marino et al. (2019)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Kenneth Marino, Mohammad Rastegari, Ali Farhadi, and Roozbeh Mottaghi. 2019. ]{.ltx_bibblock} [Ok-vqa: A visual question answering benchmark requiring external knowledge. ]{.ltx_bibblock} [In *Conference on Computer Vision and Pattern Recognition (CVPR)*. ]{.ltx_bibblock}]{#bib.bib27}
- [[OpenAI (2022)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ OpenAI. 2022. ]{.ltx_bibblock} [[Chatgpt](https://openai.com/blog/chatgpt){.ltx_ref .ltx_href target="_blank"}. ]{.ltx_bibblock}]{#bib.bib28}
- [[OpenAI (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ OpenAI. 2023. ]{.ltx_bibblock} [Gpt 4 technical report. ]{.ltx_bibblock} [*arXiv preprint arXiv:2303.08774*. ]{.ltx_bibblock}]{#bib.bib29}
- [[Park et al. (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Joon Sung Park, Joseph O'Brien, Carrie Jun Cai, Meredith Ringel Morris, Percy Liang, and Michael S Bernstein. 2023. ]{.ltx_bibblock} [Generative agents: Interactive simulacra of human behavior. ]{.ltx_bibblock} [In *Proceedings of the 36th Annual ACM Symposium on User Interface Software and Technology*, pages 1--22. ]{.ltx_bibblock}]{#bib.bib30}
- [[Radford et al. (2021)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Alec Radford, Jong Wook Kim, Chris Hallacy, Aditya Ramesh, Gabriel Goh, Sandhini Agarwal, Girish Sastry, Amanda Askell, Pamela Mishkin, Jack Clark, et al. 2021. ]{.ltx_bibblock} [Learning transferable visual models from natural language supervision. ]{.ltx_bibblock} [In *International conference on machine learning*, pages 8748--8763. PMLR. ]{.ltx_bibblock}]{#bib.bib31}
- [[Ross (2002)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ S.M. Ross. 2002. ]{.ltx_bibblock} [[*Simulation*](https://books.google.com/books?id=DApvQgAACAAJ){.ltx_ref .ltx_href target="_blank"}. ]{.ltx_bibblock} [Academic Press. ]{.ltx_bibblock}]{#bib.bib32}
- [[Shinn et al. (2024)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Noah Shinn, Federico Cassano, Ashwin Gopinath, Karthik Narasimhan, and Shunyu Yao. 2024. ]{.ltx_bibblock} [Reflexion: Language agents with verbal reinforcement learning. ]{.ltx_bibblock} [*Advances in Neural Information Processing Systems*, 36. ]{.ltx_bibblock}]{#bib.bib33}
- [[Wang et al. (2024a)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Qineng Wang, Zihao Wang, Ying Su, Hanghang Tong, and Yangqiu Song. 2024a. ]{.ltx_bibblock} [Rethinking the bounds of llm reasoning: Are multi-agent discussions the key? ]{.ltx_bibblock} [*arXiv preprint arXiv:2402.18272*. ]{.ltx_bibblock}]{#bib.bib34}
- [[Wang et al. (2022)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xuezhi Wang, Jason Wei, Dale Schuurmans, Quoc V Le, Ed H Chi, Sharan Narang, Aakanksha Chowdhery, and Denny Zhou. 2022. ]{.ltx_bibblock} [Self-consistency improves chain of thought reasoning in language models. ]{.ltx_bibblock} [In *The Eleventh International Conference on Learning Representations*. ]{.ltx_bibblock}]{#bib.bib35}
- [[Wang et al. (2024b)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zihan Wang, Yunxuan Li, Yuexin Wu, Liangchen Luo, Le Hou, Hongkun Yu, and Jingbo Shang. 2024b. ]{.ltx_bibblock} [Multi-step problem solving through a verifier: An empirical analysis on model-induced process supervision. ]{.ltx_bibblock} [*arXiv preprint arXiv:2402.02658*. ]{.ltx_bibblock}]{#bib.bib36}
- [[Wei et al. (2022)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jason Wei, Xuezhi Wang, Dale Schuurmans, Maarten Bosma, Fei Xia, Ed Chi, Quoc V Le, Denny Zhou, et al. 2022. ]{.ltx_bibblock} [Chain-of-thought prompting elicits reasoning in large language models. ]{.ltx_bibblock} [*Advances in neural information processing systems*, 35:24824--24837. ]{.ltx_bibblock}]{#bib.bib37}
- [[Welleck et al. (2022)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Sean Welleck, Ximing Lu, Peter West, Faeze Brahman, Tianxiao Shen, Daniel Khashabi, and Yejin Choi. 2022. ]{.ltx_bibblock} [Generating sequences by learning to self-correct. ]{.ltx_bibblock} [*arXiv preprint arXiv:2211.00053*. ]{.ltx_bibblock}]{#bib.bib38}
- [[Yang et al. (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zhengyuan Yang, Linjie Li, Jianfeng Wang, Kevin Lin, Ehsan Azarnasab, Faisal Ahmed, Zicheng Liu, Ce Liu, Michael Zeng, and Lijuan Wang. 2023. ]{.ltx_bibblock} [Mm-react: Prompting chatgpt for multimodal reasoning and action. ]{.ltx_bibblock} [*arXiv preprint arXiv:2303.11381*. ]{.ltx_bibblock}]{#bib.bib39}
- [[Yao et al. (2024)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Shunyu Yao, Dian Yu, Jeffrey Zhao, Izhak Shafran, Tom Griffiths, Yuan Cao, and Karthik Narasimhan. 2024. ]{.ltx_bibblock} [Tree of thoughts: Deliberate problem solving with large language models. ]{.ltx_bibblock} [*Advances in Neural Information Processing Systems*, 36. ]{.ltx_bibblock}]{#bib.bib40}
- [[Yu et al. (2022)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jiahui Yu, Zirui Wang, Vijay Vasudevan, Legg Yeung, Mojtaba Seyedhosseini, and Yonghui Wu. 2022. ]{.ltx_bibblock} [Coca: Contrastive captioners are image-text foundation models. ]{.ltx_bibblock} [*arXiv preprint arXiv:2205.01917*. ]{.ltx_bibblock}]{#bib.bib41}
- [[Yue et al. (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xiang Yue, Yuansheng Ni, Kai Zhang, Tianyu Zheng, Ruoqi Liu, Ge Zhang, Samuel Stevens, Dongfu Jiang, Weiming Ren, Yuxuan Sun, et al. 2023. ]{.ltx_bibblock} [Mmmu: A massive multi-discipline multimodal understanding and reasoning benchmark for expert agi. ]{.ltx_bibblock} [*arXiv preprint arXiv:2311.16502*. ]{.ltx_bibblock}]{#bib.bib42}
- [[Zhang et al. (2023)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zhuosheng Zhang, Aston Zhang, Mu Li, Hai Zhao, George Karypis, and Alex Smola. 2023. ]{.ltx_bibblock} [Multimodal chain-of-thought reasoning in language models. ]{.ltx_bibblock} [*arXiv preprint arXiv:2302.00923*. ]{.ltx_bibblock}]{#bib.bib43}
- [[Zhao et al. (2024)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xueliang Zhao, Xinting Huang, Tingchen Fu, Qintong Li, Shansan Gong, Lemao Liu, Wei Bi, and Lingpeng Kong. 2024. ]{.ltx_bibblock} [Bba: Bi-modal behavioral alignment for reasoning with large vision-language models. ]{.ltx_bibblock} [*arXiv preprint arXiv:2402.13577*. ]{.ltx_bibblock}]{#bib.bib44}
- [[Zheng et al. (2024)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Changmeng Zheng, Dayong Liang, Wengyu Zhang, Xiao-Yong Wei, Tat-Seng Chua, and Qing Li. 2024. ]{.ltx_bibblock} [A picture is worth a graph: Blueprint debate on graph for multimodal reasoning. ]{.ltx_bibblock} [*arXiv preprint arXiv:2403.14972*. ]{.ltx_bibblock}]{#bib.bib45}
- [[Zhou et al. (2022)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Denny Zhou, Nathanael Schärli, Le Hou, Jason Wei, Nathan Scales, Xuezhi Wang, Dale Schuurmans, Claire Cui, Olivier Bousquet, Quoc V Le, et al. 2022. ]{.ltx_bibblock} [Least-to-most prompting enables complex reasoning in large language models. ]{.ltx_bibblock} [In *The Eleventh International Conference on Learning Representations*. ]{.ltx_bibblock}]{#bib.bib46}
:::

::::::::::::::::::::::::::::::::::::::::::::::::::: {#A1 .section .ltx_appendix}
## [Appendix A ]{.ltx_tag .ltx_tag_appendix}Prompt Templates {#appendix-a-prompt-templates .ltx_title .ltx_title_appendix}

:::::::::::::: {#A1.SS1 .section .ltx_subsection}
### [A.1 ]{.ltx_tag .ltx_tag_subsection}Text Reasoning Tasks {#a.1-text-reasoning-tasks .ltx_title .ltx_title_subsection}

::: {#A1.SS1.p1 .ltx_para .ltx_noindent}
[System Prompt]{#A1.SS1.p1.1.1 .ltx_text .ltx_font_bold}:
:::

::: {#A1.SS1.p2 .ltx_para .ltx_noindent}
You are a helpful assistant with expertise in mathematics and reasoning. Your task is to assist in solving a math reasoning problem by providing a clear and detailed solution. Limit your output within 100 words, and your final answer should be a single numerical number, in the form of {{[answer]{#A1.SS1.p2.1.1 .ltx_text .ltx_font_italic}}}, at the end of your response.
:::

::: {#A1.SS1.p3 .ltx_para .ltx_noindent}
[Starting Prompt]{#A1.SS1.p3.1.1 .ltx_text .ltx_font_bold}:
:::

::: {#A1.SS1.p4 .ltx_para .ltx_noindent}
Can you solve the following math problem? {[question]{#A1.SS1.p4.1.1 .ltx_text .ltx_font_italic}} Explain your reasoning. Your final answer should be a single numerical number, in the form of {{[answer]{#A1.SS1.p4.1.2 .ltx_text .ltx_font_italic}}}, at the end of your response.
:::

::: {#A1.SS1.p5 .ltx_para .ltx_noindent}
[Debate Prompt]{#A1.SS1.p5.1.1 .ltx_text .ltx_font_bold}:
:::

::: {#A1.SS1.p6 .ltx_para .ltx_noindent}
These are the solutions to the problem from other agents:
:::

::: {#A1.SS1.p7 .ltx_para .ltx_noindent}
One agent solution: {[reference solution]{#A1.SS1.p7.1.1 .ltx_text .ltx_font_italic}}
:::

::: {#A1.SS1.p8 .ltx_para .ltx_noindent}
One agent solution: {[reference solution]{#A1.SS1.p8.1.1 .ltx_text .ltx_font_italic}}
:::

::: {#A1.SS1.p9 .ltx_para .ltx_noindent}
One agent solution: {[reference solution]{#A1.SS1.p9.1.1 .ltx_text .ltx_font_italic}}
:::

::: {#A1.SS1.p10 .ltx_para .ltx_noindent}
...
:::

::: {#A1.SS1.p11 .ltx_para .ltx_noindent}
Using the solutions from other agents as additional information, can you provide your answer to the math problem? The original math problem is {[question]{#A1.SS1.p11.1.1 .ltx_text .ltx_font_italic}}. Your final answer should be a single numerical number, in the form of {{[answer]{#A1.SS1.p11.1.2 .ltx_text .ltx_font_italic}}}, at the end of your response.
:::
::::::::::::::

::::::::::::::::::: {#A1.SS2 .section .ltx_subsection}
### [A.2 ]{.ltx_tag .ltx_tag_subsection}Multimodal Reasoning Tasks {#a.2-multimodal-reasoning-tasks .ltx_title .ltx_title_subsection}

::: {#A1.SS2.p1 .ltx_para .ltx_noindent}
[System Prompt]{#A1.SS2.p1.1.1 .ltx_text .ltx_font_bold}: Please answer the question requiring an integer answer or a floating-point number with one decimal place and provide the final value, e.g., 1, 2, 3, 1.1, 1.2, 1.3, at the end.
:::

::: {#A1.SS2.p2 .ltx_para .ltx_noindent}
When providing an answer,
:::

::: {#A1.SS2.p3 .ltx_para}
1\. Reason through the question step by step.
:::

::: {#A1.SS2.p4 .ltx_para}
2\. Keep each step concise, ensuring the total reasoning is under 400 words. Conclude with the final answer in the format \"\*\*FINAL ANSWER:\*\* (X)\" where X should be the numerical answer. Note that the answer has to be surrounded by the parenthesis, e.g., \*\*FINAL ANSWER:\*\* (1).
:::

::: {#A1.SS2.p5 .ltx_para .ltx_noindent}
[Starting Prompt]{#A1.SS2.p5.1.1 .ltx_text .ltx_font_bold}:
:::

::: {#A1.SS2.p6 .ltx_para .ltx_noindent}
{[question]{#A1.SS2.p6.1.1 .ltx_text .ltx_font_italic}}
:::

::: {#A1.SS2.p7 .ltx_para .ltx_noindent}
[Debate Prompt]{#A1.SS2.p7.1.1 .ltx_text .ltx_font_bold}:
:::

::: {#A1.SS2.p8 .ltx_para .ltx_noindent}
Below are responses from {[number of visible agents]{#A1.SS2.p8.1.1 .ltx_text .ltx_font_italic}} other agents:
:::

::: {#A1.SS2.p9 .ltx_para .ltx_noindent}
Response {[agent index]{#A1.SS2.p9.1.1 .ltx_text .ltx_font_italic}}: {[reference solution]{#A1.SS2.p9.1.2 .ltx_text .ltx_font_italic}}
:::

::: {#A1.SS2.p10 .ltx_para .ltx_noindent}
Response {[agent index]{#A1.SS2.p10.1.1 .ltx_text .ltx_font_italic}}: {[reference solution]{#A1.SS2.p10.1.2 .ltx_text .ltx_font_italic}}
:::

::: {#A1.SS2.p11 .ltx_para .ltx_noindent}
Response {[agent index]{#A1.SS2.p11.1.1 .ltx_text .ltx_font_italic}}: {[reference solution]{#A1.SS2.p11.1.2 .ltx_text .ltx_font_italic}}
:::

::: {#A1.SS2.p12 .ltx_para .ltx_noindent}
...
:::

::: {#A1.SS2.p13 .ltx_para .ltx_noindent}
These responses contain step-by-step reasoning. Use the insights from all the responses to refine and update your answer. When providing an answer:
:::

::: {#A1.SS2.p14 .ltx_para}
1\. Take a close look at the question again.
:::

::: {#A1.SS2.p15 .ltx_para}
2\. Reason through the question, your previous response and responses from other agents step by step.
:::

::: {#A1.SS2.p16 .ltx_para}
3\. Keep each step concise, ensuring the total reasoning is under 400 words. Conclude with the final answer in the format \"\*\*FINAL ANSWER:\*\* (X)\" where X should be the numerical answer. Note that the answer has to be surrounded by the parenthesis, e.g., \*\*FINAL ANSWER:\*\* (1).
:::
:::::::::::::::::::

::::::::::::::::::::: {#A1.SS3 .section .ltx_subsection}
### [A.3 ]{.ltx_tag .ltx_tag_subsection}Alignment labeling Tasks {#a.3-alignment-labeling-tasks .ltx_title .ltx_title_subsection}

::: {#A1.SS3.p1 .ltx_para .ltx_noindent}
[System Prompt]{#A1.SS3.p1.1.1 .ltx_text .ltx_font_bold}:
:::

::: {#A1.SS3.p2 .ltx_para .ltx_noindent}
You are an expert rater of helpful and honest assistant responses. The user will give you two dialogues, please help user choose the one with the most helpful and honest response. Limit your answer within 100 words, and put your final answer (either (A) or (B)) at the end of your response.
:::

::: {#A1.SS3.p3 .ltx_para .ltx_noindent}
[Starting Prompt]{#A1.SS3.p3.1.1 .ltx_text .ltx_font_bold}:
:::

::: {#A1.SS3.p4 .ltx_para .ltx_noindent}
The following are two dialogues. Which one is more helpful and honest?  \
:::

::: {#A1.SS3.p5 .ltx_para .ltx_noindent}
\[Start of Dialogue A\]
:::

::: {#A1.SS3.p6 .ltx_para .ltx_noindent}
{[Dialogue A]{#A1.SS3.p6.1.1 .ltx_text .ltx_font_italic}}
:::

::: {#A1.SS3.p7 .ltx_para .ltx_noindent}
\[End of Dialogue A\]  \
:::

::: {#A1.SS3.p8 .ltx_para .ltx_noindent}
\[Start of Dialogue B\]
:::

::: {#A1.SS3.p9 .ltx_para .ltx_noindent}
{[Dialogue B]{#A1.SS3.p9.1.1 .ltx_text .ltx_font_italic}}
:::

::: {#A1.SS3.p10 .ltx_para .ltx_noindent}
\[End of Dialogue B\]
:::

::: {#A1.SS3.p11 .ltx_para .ltx_noindent}
Limit your answer within 100 words, and put your final answer (either (A) or (B)) at the end of your response.
:::

::: {#A1.SS3.p12 .ltx_para .ltx_noindent}
[Debate Prompt]{#A1.SS3.p12.1.1 .ltx_text .ltx_font_bold}:
:::

::: {#A1.SS3.p13 .ltx_para .ltx_noindent}
These are the solutions to the problem from other agents:
:::

::: {#A1.SS3.p14 .ltx_para .ltx_noindent}
One agent solution: {[reference solution]{#A1.SS3.p14.1.1 .ltx_text .ltx_font_italic}}
:::

::: {#A1.SS3.p15 .ltx_para .ltx_noindent}
One agent solution: {[reference solution]{#A1.SS3.p15.1.1 .ltx_text .ltx_font_italic}}
:::

::: {#A1.SS3.p16 .ltx_para .ltx_noindent}
One agent solution: {[reference solution]{#A1.SS3.p16.1.1 .ltx_text .ltx_font_italic}}
:::

::: {#A1.SS3.p17 .ltx_para .ltx_noindent}
...
:::

::: {#A1.SS3.p18 .ltx_para .ltx_noindent}
Using the reasoning from other agents as additional advice, can you provide an updated answer? Examine your solution and those of other agents step by step. Limit your answer within 100 words, and put your final answer (either (A) or (B)) at the end of your response.
:::
:::::::::::::::::::::
:::::::::::::::::::::::::::::::::::::::::::::::::::

:::: {#A2 .section .ltx_appendix}
## [Appendix B ]{.ltx_tag .ltx_tag_appendix}Additional Experiments with Different Temperature {#appendix-b-additional-experiments-with-different-temperature .ltx_title .ltx_title_appendix}

::: {#A2.p1 .ltx_para}
For multimodal experiments, we also examined how different temperatures affect the performance of MAD. We compared the accuracy and cost savings between the default temperature $T = 1$ for GPT-4o and a more conservative temperature $T = 0.25$, aiming to generate more consistent answers. While Table [[3]{.ltx_text .ltx_ref_tag}](#S5.T3 "Table 3 ‣ 5.2 MAD on Multimodal Reasoning Task ‣ 5 Experiments: MAD with Single LLM ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref} reports performance at $T = 1$, we observed almost no difference in accuracy with $T = 0.25$. However, $T = 0.25$ resulted in slightly greater cost savings, as shown in Table [[7]{.ltx_text .ltx_ref_tag}](#A2.T7 "Table 7 ‣ Appendix B Additional Experiments with Different Temperature ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref}.
:::

<figure id="A2.T7" class="ltx_table">
<table id="A2.T7.14" class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<thead class="ltx_thead">
<tr id="A2.T7.14.15.1" class="ltx_tr">
<th id="A2.T7.14.15.1.1" class="ltx_td ltx_align_left ltx_th ltx_th_column ltx_th_row ltx_border_r ltx_border_t"><span id="A2.T7.14.15.1.1.1" class="ltx_text ltx_font_bold">Method</span></th>
<th id="A2.T7.14.15.1.2" class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_t">Accuracy</th>
<th id="A2.T7.14.15.1.3" class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_t">Cost Saving</th>
</tr>
</thead>
<tbody class="ltx_tbody">
<tr id="A2.T7.2.2" class="ltx_tr">
<th id="A2.T7.1.1.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">MAD (<span class="math inline"><em>D</em> = 1</span>)</th>
<td id="A2.T7.2.2.2" class="ltx_td ltx_align_center ltx_border_t">57.8 <span class="math inline">±</span> 1.0</td>
<td id="A2.T7.2.2.3" class="ltx_td ltx_align_center ltx_border_t">baseline</td>
</tr>
<tr id="A2.T7.6.6" class="ltx_tr">
<th id="A2.T7.3.3.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">MAD (<span class="math inline"><em>D</em> = 4/5</span>)</th>
<td id="A2.T7.4.4.2" class="ltx_td ltx_align_center">57.4 <span class="math inline">±</span> 0.6</td>
<td id="A2.T7.6.6.4" class="ltx_td ltx_align_center"><span class="math inline">−</span>11.8% (<span class="math inline">−</span>14.3%)</td>
</tr>
<tr id="A2.T7.10.10" class="ltx_tr">
<th id="A2.T7.7.7.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">MAD (<span class="math inline"><em>D</em> = 3/5</span>)</th>
<td id="A2.T7.8.8.2" class="ltx_td ltx_align_center">57.4 <span class="math inline">±</span> 3.5</td>
<td id="A2.T7.10.10.4" class="ltx_td ltx_align_center"><span class="math inline">−</span>21.1% (<span class="math inline">−</span>26.0%)</td>
</tr>
<tr id="A2.T7.14.14" class="ltx_tr">
<th id="A2.T7.11.11.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_b ltx_border_r">MAD (<span class="math inline"><em>D</em> = 2/5</span>)</th>
<td id="A2.T7.12.12.2" class="ltx_td ltx_align_center ltx_border_b"><span id="A2.T7.12.12.2.1" class="ltx_text ltx_font_bold">59.0 <span class="math inline">±</span> 1.0</span></td>
<td id="A2.T7.14.14.4" class="ltx_td ltx_align_center ltx_border_b"><span class="math inline">−</span>37.6% (<span class="math inline">−</span>46.5%)</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 7: </span>Comparison of accuracy and cost savings of different MADs on the MathVista dataset. All experiments were conducted using the GPT-4o model with temperature set to <span class="math inline">0.25</span>. The cost saving percentages in parenthesis are computed without multimodal inputs.</figcaption>
</figure>
::::

:::: {#A3 .section .ltx_appendix}
## [Appendix C ]{.ltx_tag .ltx_tag_appendix}Additional Experiments with 4 Agents {#appendix-c-additional-experiments-with-4-agents .ltx_title .ltx_title_appendix}

::: {#A3.p1 .ltx_para}
Regular graph with 4 agents only have two configurations (as shown in Figure [[7]{.ltx_text .ltx_ref_tag}](#A3.F7 "Figure 7 ‣ Appendix C Additional Experiments with 4 Agents ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref}). Our experiments on GSM8K shows similar pattern in accuracy between these two setup, shown in Table [[8]{.ltx_text .ltx_ref_tag}](#A3.T8 "Table 8 ‣ Appendix C Additional Experiments with 4 Agents ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref}.
:::

<figure id="A3.F7" class="ltx_figure">
<img src="/html/2406.11776/assets/imgs/agent_4_graph.png" id="A3.F7.g1" class="ltx_graphics ltx_img_landscape" width="598" height="273" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 7: </span>Regular graph with 4 agents.</figcaption>
</figure>

<figure id="A3.T8" class="ltx_table">
<table id="A3.T8.5" class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<thead class="ltx_thead">
<tr id="A3.T8.5.6.1" class="ltx_tr">
<th id="A3.T8.5.6.1.1" class="ltx_td ltx_align_left ltx_th ltx_th_column ltx_th_row ltx_border_r ltx_border_t"><span id="A3.T8.5.6.1.1.1" class="ltx_text ltx_font_bold">Method</span></th>
<th id="A3.T8.5.6.1.2" class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_t">Accuracy</th>
<th id="A3.T8.5.6.1.3" class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_t">Cost</th>
</tr>
<tr id="A3.T8.5.7.2" class="ltx_tr">
<th id="A3.T8.5.7.2.1" class="ltx_td ltx_align_left ltx_th ltx_th_column ltx_th_row ltx_border_r ltx_border_t">SC</th>
<th id="A3.T8.5.7.2.2" class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_t">81.0</th>
<th id="A3.T8.5.7.2.3" class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_t">-</th>
</tr>
</thead>
<tbody class="ltx_tbody">
<tr id="A3.T8.2.2" class="ltx_tr">
<th id="A3.T8.1.1.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t"><span class="math inline"><em>D</em> = 1</span></th>
<td id="A3.T8.2.2.2" class="ltx_td ltx_align_center ltx_border_t">81.7 <span class="math inline">±</span> 0.9</td>
<td id="A3.T8.2.2.3" class="ltx_td ltx_align_center ltx_border_t">baseline</td>
</tr>
<tr id="A3.T8.5.5" class="ltx_tr">
<th id="A3.T8.3.3.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_b ltx_border_r"><span class="math inline"><em>D</em> = 2/3</span></th>
<td id="A3.T8.4.4.2" class="ltx_td ltx_align_center ltx_border_b"><span id="A3.T8.4.4.2.1" class="ltx_text ltx_font_bold">82.7</span> <span class="math inline">±</span> 1.2</td>
<td id="A3.T8.5.5.3" class="ltx_td ltx_align_center ltx_border_b"><span class="math inline">−</span>25.6%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 8: </span>Accuracy comparison of MAD against baseline methods on the GSM8K dataset. Experiments were conducted using the GPT-3.5model.</figcaption>
</figure>
::::

:::: {#A4 .section .ltx_appendix}
## [Appendix D ]{.ltx_tag .ltx_tag_appendix}ProbMAD: MAD with Probablistic Topology {#appendix-d-probmad-mad-with-probablistic-topology .ltx_title .ltx_title_appendix}

::: {#A4.p1 .ltx_para}
While we primarily focus on sparse MADs with fixed communication topology, we also investigate ProbMAD where communication is probablistic. For any MAD with a given $D$, the ProbMAD counterpart is a topology where the probability that a given agent sees any reference solution from previous round is $D$. In Table [[9]{.ltx_text .ltx_ref_tag}](#A4.T9 "Table 9 ‣ Appendix D ProbMAD: MAD with Probablistic Topology ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref}, we use GPT-3.5 on GSM8K to show that the performance of ProbMAD is comparable to fully-connected MAD and its cost-saving ability is similar to sparse MAD topologies we discuss earlier. More work is to be done to compare deterministic and probablistic sparsity and explain the mechanism. In the meantime, we show that the probablistic way of thinking about communication topology allows our approach to be even more generally applicable to any number of agents.
:::

<figure id="A4.T9" class="ltx_table">
<table id="A4.T9.12" class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr id="A4.T9.12.13.1" class="ltx_tr">
<th id="A4.T9.12.13.1.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t"><span id="A4.T9.12.13.1.1.1" class="ltx_text ltx_font_bold">Method</span></th>
<td id="A4.T9.12.13.1.2" class="ltx_td ltx_align_center ltx_border_t">Accuracy</td>
<td id="A4.T9.12.13.1.3" class="ltx_td ltx_align_center ltx_border_t">Cost Saving</td>
</tr>
<tr id="A4.T9.1.1" class="ltx_tr">
<th id="A4.T9.1.1.2" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">CoT</th>
<td id="A4.T9.1.1.1" class="ltx_td ltx_align_center ltx_border_t">77.5 <span class="math inline">±</span> 4.2</td>
<td id="A4.T9.1.1.3" class="ltx_td ltx_align_center ltx_border_t">-</td>
</tr>
<tr id="A4.T9.12.14.2" class="ltx_tr">
<th id="A4.T9.12.14.2.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">SC</th>
<td id="A4.T9.12.14.2.2" class="ltx_td ltx_align_center">80.0</td>
<td id="A4.T9.12.14.2.3" class="ltx_td ltx_align_center">-</td>
</tr>
<tr id="A4.T9.3.3" class="ltx_tr">
<th id="A4.T9.2.2.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r ltx_border_t">MAD (<span class="math inline"><em>D</em> = 1</span>)</th>
<td id="A4.T9.3.3.2" class="ltx_td ltx_align_center ltx_border_t"><span id="A4.T9.3.3.2.1" class="ltx_text ltx_font_bold">84.5 <span class="math inline">±</span> 1.5</span></td>
<td id="A4.T9.3.3.3" class="ltx_td ltx_align_center ltx_border_t">baseline</td>
</tr>
<tr id="A4.T9.6.6" class="ltx_tr">
<th id="A4.T9.4.4.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">ProbMAD (<span class="math inline"><em>D</em> = 4/5</span>)</th>
<td id="A4.T9.5.5.2" class="ltx_td ltx_align_center"><span id="A4.T9.5.5.2.1" class="ltx_text ltx_font_bold">84.5 <span class="math inline">±</span> 0.7</span></td>
<td id="A4.T9.6.6.3" class="ltx_td ltx_align_center"><span class="math inline">−</span>14.3%</td>
</tr>
<tr id="A4.T9.9.9" class="ltx_tr">
<th id="A4.T9.7.7.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_r">ProbMAD (<span class="math inline"><em>D</em> = 3/5</span>)</th>
<td id="A4.T9.8.8.2" class="ltx_td ltx_align_center">83.5 <span class="math inline">±</span> 0.7</td>
<td id="A4.T9.9.9.3" class="ltx_td ltx_align_center"><span class="math inline">−</span>29.6%</td>
</tr>
<tr id="A4.T9.12.12" class="ltx_tr">
<th id="A4.T9.10.10.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_b ltx_border_r">ProbMAD (<span class="math inline"><em>D</em> = 2/5</span>)</th>
<td id="A4.T9.11.11.2" class="ltx_td ltx_align_center ltx_border_b">84.0 <span class="math inline">±</span> 1.7</td>
<td id="A4.T9.12.12.3" class="ltx_td ltx_align_center ltx_border_b"><span class="math inline">−</span>47.1%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 9: </span>Comparison of accuracy and cost savings of probabilistic MAD against baseline methods on the GSM8K dataset. All experiments were conducted using the GPT-3.5 model.</figcaption>
</figure>
::::

:::: {#A5 .section .ltx_appendix}
## [Appendix E ]{.ltx_tag .ltx_tag_appendix}Rounds of Effective Debate for Mistral 7B {#appendix-e-rounds-of-effective-debate-for-mistral-7b .ltx_title .ltx_title_appendix}

::: {#A5.p1 .ltx_para}
Similar to what we observe on GPT-3.5, the rounds of effective debate using Mistral 7B model also increases on both preference tasks when MAD becomes sparse (Figure [[8]{.ltx_text .ltx_ref_tag}](#A5.F8 "Figure 8 ‣ Appendix E Rounds of Effective Debate for Mistral 7B ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref}.
:::

<figure id="A5.F8" class="ltx_figure">
<img src="/html/2406.11776/assets/imgs/num_rounds_mistral_corrected.png" id="A5.F8.g1" class="ltx_graphics ltx_img_landscape" width="598" height="368" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 8: </span>Effective debate rounds for each topology design in alignment labeling tasks using the Mistral 7B model.</figcaption>
</figure>
::::

:::: {#A6 .section .ltx_appendix}
## [Appendix F ]{.ltx_tag .ltx_tag_appendix}Types of Agent Behaviors {#appendix-f-types-of-agent-behaviors .ltx_title .ltx_title_appendix}

::: {#A6.p1 .ltx_para}
During the multi-agent debate process, we observe four common types of agent responses to reference solutions (Figure [[9]{.ltx_text .ltx_ref_tag}](#A6.F9 "Figure 9 ‣ Appendix F Types of Agent Behaviors ‣ Improving Multi-Agent Debate with Sparse Communication Topology"){.ltx_ref}). Agents may learn from other agents' reasoning, correct a mistake made by another agent, act as an arbitrator to evaluate others' solutions, or occasionally be misled by the input of their peers.
:::

<figure id="A6.F9" class="ltx_figure">
<img src="/html/2406.11776/assets/imgs/agent_example.png" id="A6.F9.g1" class="ltx_graphics ltx_img_landscape" width="598" height="448" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 9: </span>Common types (with nicknames) of agent behaviors when given reference solutions.</figcaption>
</figure>
::::
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

::: ar5iv-footer
[◄](/html/2406.11775){.ar5iv-nav-button .ar5iv-nav-button-prev} [![ar5iv homepage](/assets/ar5iv.png){height="40"}](/){.ar5iv-home-button} [Feeling\
lucky?](/feeling_lucky){.ar5iv-text-button} [](/land_of_honey_and_milk){rel="nofollow" aria-hidden="true" tabindex="-1"} [Conversion\
report](/log/2406.11776){.ar5iv-text-button .ar5iv-severity-warning} [Report\
an issue](https://github.com/dginev/ar5iv/issues/new?template=improve-article--arxiv-id-.md&title=Improve+article+2406.11776){.ar5iv-text-button target="_blank"} [View original\
on arXiv](https://arxiv.org/abs/2406.11776){.ar5iv-text-button .arxiv-ui-theme}[►](/html/2406.11777){.ar5iv-nav-button .ar5iv-nav-button-next}
:::

[[]{.color-scheme-icon}](javascript:toggleColorScheme() "Toggle ar5iv color scheme"){.ar5iv-toggle-color-scheme} [Copyright](https://arxiv.org/help/license){.ar5iv-footer-button target="_blank"} [Privacy Policy](https://arxiv.org/help/policies/privacy_policy){.ar5iv-footer-button target="_blank"}

::: ltx_page_logo
Generated on Sat Jul 6 00:52:19 2024 by [[L[a]{.ltx_font_smallcaps style="position:relative; bottom:2.2pt;"}T[e]{.ltx_font_smallcaps style="font-size:120%;position:relative; bottom:-0.2ex;"}]{style="letter-spacing:-0.2em; margin-right:0.1em;"}[XML]{style="font-size:90%; position:relative; bottom:-0.2ex;"}![Mascot Sammy](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAsAAAAOCAYAAAD5YeaVAAAAAXNSR0IArs4c6QAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAALEwAACxMBAJqcGAAAAAd0SU1FB9wKExQZLWTEaOUAAAAddEVYdENvbW1lbnQAQ3JlYXRlZCB3aXRoIFRoZSBHSU1Q72QlbgAAAdpJREFUKM9tkL+L2nAARz9fPZNCKFapUn8kyI0e4iRHSR1Kb8ng0lJw6FYHFwv2LwhOpcWxTjeUunYqOmqd6hEoRDhtDWdA8ApRYsSUCDHNt5ul13vz4w0vWCgUnnEc975arX6ORqN3VqtVZbfbTQC4uEHANM3jSqXymFI6yWazP2KxWAXAL9zCUa1Wy2tXVxheKA9YNoR8Pt+aTqe4FVVVvz05O6MBhqUIBGk8Hn8HAOVy+T+XLJfLS4ZhTiRJgqIoVBRFIoric47jPnmeB1mW/9rr9ZpSSn3Lsmir1fJZlqWlUonKsvwWwD8ymc/nXwVBeLjf7xEKhdBut9Hr9WgmkyGEkJwsy5eHG5vN5g0AKIoCAEgkEkin0wQAfN9/cXPdheu6P33fBwB4ngcAcByHJpPJl+fn54mD3Gg0NrquXxeLRQAAwzAYj8cwTZPwPH9/sVg8PXweDAauqqr2cDjEer1GJBLBZDJBs9mE4zjwfZ85lAGg2+06hmGgXq+j3+/DsixYlgVN03a9Xu8jgCNCyIegIAgx13Vfd7vdu+FweG8YRkjXdWy329+dTgeSJD3ieZ7RNO0VAXAPwDEAO5VKndi2fWrb9jWl9Esul6PZbDY9Go1OZ7PZ9z/lyuD3OozU2wAAAABJRU5ErkJggg==)](http://dlmf.nist.gov/LaTeXML/){.ltx_LaTeXML_logo target="_blank"}
:::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
