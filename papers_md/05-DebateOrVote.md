:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_main
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_content
# Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models? {#debate-or-vote-which-yields-better-decisions-in-multi-agent-large-language-models .ltx_title .ltx_title_document}

::: ltx_authors
[ [Hyeong Kyu Choi    Xiaojin Zhu    Yixuan Li\
Department of Computer Sciences, University of Wisconsin-Madison\
[{froilanchoi, jerryzhu, sharonli}@cs.wisc.edu]{.ltx_text .ltx_font_typewriter}\
]{.ltx_personname}[Corresponding author]{.ltx_author_notes}]{.ltx_creator .ltx_role_author}
:::

::: ltx_abstract
###### Abstract {#abstract .ltx_title .ltx_title_abstract}

Multi-Agent Debate (MAD) has emerged as a promising paradigm for improving the performance of large language models through collaborative reasoning. Despite recent advances, the key factors driving MAD's effectiveness remain unclear. In this work, we disentangle MAD into two key components--Majority Voting and inter-agent Debate--and assess their respective contributions. Through extensive experiments across seven NLP benchmarks, we find that Majority Voting alone accounts for most of the performance gains typically attributed to MAD. To explain this, we propose a theoretical framework that models debate as a stochastic process. We prove that it induces a martingale over agents' belief trajectories, implying that debate alone does not improve expected correctness. Guided by these insights, we demonstrate that targeted interventions, by biasing the belief update toward correction, can meaningfully enhance debate effectiveness. Overall, our findings suggest that while MAD has potential, simple ensembling methods remain strong and more reliable alternatives in many practical settings. Code is released in [https://github.com/deeplearning-wisc/debate-or-vote](https://github.com/deeplearning-wisc/debate-or-vote){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}.
:::

:::::::::: {#S1 .section .ltx_section}
### [1 ]{.ltx_tag .ltx_tag_section}Introduction {#introduction .ltx_title .ltx_title_section}

::: {#S1.p1 .ltx_para}
> ["Out of intense complexities, intense simplicities emerge."]{.ltx_text .ltx_font_italic}
>
> --- [W. Churchill\
> ]{.ltx_text .ltx_font_smallcaps}
:::

::: {#S1.p2 .ltx_para}
Throughout history, humans have relied on deliberation to resolve ambiguity, challenge assumptions, and seek better answers. From courtrooms and panels to scientific collaborations, group reasoning plays a central role in decision-making. This process---where individuals reflect, revise, and converge through interaction---has long been seen as a hallmark of intelligent behavior. Inspired by this, recent work has explored whether large language models (LLMs) might similarly benefit from structured interaction. Multi-Agent Debate (MAD) has emerged as a popular framework: multiple LLM agents are prompted to discuss a shared question, each updating their answer based on the responses of their peers \[[1](#bib.bib1){.ltx_ref}, [2](#bib.bib2){.ltx_ref}, [3](#bib.bib3){.ltx_ref}, [4](#bib.bib4){.ltx_ref}, [5](#bib.bib5){.ltx_ref}, [6](#bib.bib6){.ltx_ref}\]. The hope is that, like human deliberation, such interaction will improve reasoning and lead to better outcomes.
:::

::: {#S1.p3 .ltx_para}
At its core, MAD integrates two key ingredients: the use of multiple agents ("Multi-Agent") and their interaction through iterative discussions ("Debate"). Recent work has introduced increasingly sophisticated variants---ranging from diverse communication protocols \[[3](#bib.bib3){.ltx_ref}, [7](#bib.bib7){.ltx_ref}, [8](#bib.bib8){.ltx_ref}\], designing efficient and effective system architectures \[[1](#bib.bib1){.ltx_ref}, [2](#bib.bib2){.ltx_ref}, [9](#bib.bib9){.ltx_ref}, [10](#bib.bib10){.ltx_ref}\], and assigning varied roles or personas to agents \[[11](#bib.bib11){.ltx_ref}, [12](#bib.bib12){.ltx_ref}, [13](#bib.bib13){.ltx_ref}\]. Despite these advances, the underlying mechanisms behind MAD's effectiveness remain unclear. A natural step toward understanding MAD's performance is to disentangle the contribution of each component---*are the gains primarily due to meaningful communication between agents, or simply the result of aggregating multiple outputs?* Answering this question is important because it informs whether the growing complexity of MAD design is justified by tangible benefits. If most of the performance gain stems from ensembling---[i.e.]{.ltx_text .ltx_font_italic}, aggregating diverse outputs from independent agents---then simpler methods like majority voting may suffice, avoiding additional computational and architectural overhead (see Figure [[2]{.ltx_text .ltx_ref_tag}](#S1.F2 "Figure 2 ‣ 1 Introduction ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref} for visual comparison).
:::

<figure id="S1.F2" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S1.F2.fig1" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_top" style="width:144.9pt;">
<img src="/html/2508.17536/assets/x1.png" id="S1.F2.g1" class="ltx_graphics ltx_centering ltx_img_square" width="830" height="696" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 1</span>: </span><span class="ltx_text" style="font-size:90%;">Majority Voting vs. MAD overview.</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S1.F2.fig2" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_top" style="width:189.8pt;">
<img src="/html/2508.17536/assets/x2.png" id="S1.F2.g2" class="ltx_graphics ltx_centering ltx_img_landscape" width="830" height="494" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 2</span>: </span><span class="ltx_text" style="font-size:90%;">Majority Voting is the main contributor to MAD.</span></figcaption>
</figure>
</div>
</div>
</figure>

::: {#S1.p4 .ltx_para}
To better understand the relative contributions of ensembling *vs.* interaction, we conduct an extensive empirical study quantifying each component's effect. Specifically, we measure the contribution of the "Multi-Agent" component using the performance achieved through Majority Voting, [i.e.]{.ltx_text .ltx_font_italic}, the aggregated output of agents before any debate rounds occur. We then compare this baseline to the final performance after multiple rounds of "Debate", allowing us to isolate the additional benefit introduced by inter-agent communications. Surprisingly, we find that Majority Voting accounts for most of the performance gains in MAD. In fact, in most cases, *majority voting [without any debate]{.ltx_text .ltx_font_bold} performs on par with MAD*, as seen in Figure [[2]{.ltx_text .ltx_ref_tag}](#S1.F2 "Figure 2 ‣ 1 Introduction ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}. To ensure the broad applicability of our findings, our evaluation spans seven diverse benchmarks across multiple tasks and models.
:::

::: {#S1.p5 .ltx_para}
Beyond empirical observations, we introduce a theoretical framework in Section [[4]{.ltx_text .ltx_ref_tag}](#S4 "4 Theoretical Analysis ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref} that rigorously explains how agents' uncertainty and belief updates shape collective decision-making in both voting and debate. At its core, the framework models each agent as a stochastic process governed by a Dirichlet-Compound-Multinomial (DCM) distribution, capturing internal uncertainty through a Dirichlet belief prior and output randomness via Multinomial sampling. This closely mirrors the behavior of real-world LLMs, which produce different outputs for the same question due to uncertainty and stochastic generation process (*e.g.*, via temperature or nucleus sampling). Within this framework, we characterize MAD as a Bayesian posterior belief update process and prove that it induces a *martingale* over agents' belief in the correct answer---meaning the expected belief remains unchanged over debate rounds. This implies that debate itself does not systematically improve or degrade beliefs on average; rather, belief evolution is driven by stochastic peer influence. In other words, we prove formally that *majority vote does essentially all the work*, which explains our empirical findings.
:::

::: {#S1.p6 .ltx_para}
Our theoretical framework further sheds light on the new design principles to improve MAD (Section [[5]{.ltx_text .ltx_ref_tag}](#S5 "5 How Does Theory Inform Improved Design of MAD? ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}). In particular, it highlights the importance of controlling the martingale by biasing belief updates toward correct signals during debate. We operationalize this insight through several interventions, where correct responses exert more influence than misleading ones, leading to improvements over standard MAD. We summarize our contributions and significance as follows:
:::

::: {#S1.p7 .ltx_para}
1.  [[1.]{.ltx_tag .ltx_tag_item}]{#S1.I1.i1}

    ::: {#S1.I1.i1.p1 .ltx_para}
    We comprehensively demonstrate that Majority Voting is as effective as Multi-Agent Debate, when evaluated across seven representative benchmark datasets. We further expand our investigation to more general MAD settings, including configurations with larger and more capable agents, heterogeneous agent populations, and open-ended natural language tasks.
    :::
2.  [[2.]{.ltx_tag .ltx_tag_item}]{#S1.I1.i2}

    ::: {#S1.I1.i2.p1 .ltx_para}
    We develop a [new theoretical framework]{.ltx_text .ltx_font_italic} that reveals majority voting's success probability, and rigorously characterizes multi-agent debate as a martingale process. This framework lays a principled foundation for future work to better understand MAD systems.
    :::
3.  [[3.]{.ltx_tag .ltx_tag_item}]{#S1.I1.i3}

    ::: {#S1.I1.i3.p1 .ltx_para}
    Our theoretical analysis informs that debate alone does not improve beyond majority voting. By designing strategies that help preserve correct responses across debate rounds, we achieve notable improvements in multi-agent debate performance. This sheds light on future research to effectively improve MAD systems.
    :::
:::
::::::::::

:::::::::: {#S2 .section .ltx_section}
### [2 ]{.ltx_tag .ltx_tag_section}Preliminaries {#preliminaries .ltx_title .ltx_title_section}

:::: {#S2.SS0.SSS0.Px1 .section .ltx_paragraph}
##### Multi-Agent Debate {#multi-agent-debate .ltx_title .ltx_title_paragraph}

::: {#S2.SS0.SSS0.Px1.p1 .ltx_para}
is a collaborative framework in which multiple language model agents engage in structured interaction---typically in the form of iterative exchanges or discussions---to solve a task such as question answering or text generation \[[1](#bib.bib1){.ltx_ref}, [2](#bib.bib2){.ltx_ref}, [3](#bib.bib3){.ltx_ref}, [4](#bib.bib4){.ltx_ref}, [5](#bib.bib5){.ltx_ref}, [6](#bib.bib6){.ltx_ref}\]. In a typical MAD protocol, each agent independently generates an initial response and then engages in a series of debate rounds. At round $t$, an agent receives the original question along with responses from its peers at round $t-1$, prompting the model to update its response accordingly. This iterative process is designed to leverage diverse reasoning paths and peer wisdom, potentially enhancing the overall decision quality. After all rounds of debate, the final answer is typically derived through an aggregation mechanism, such as majority voting. Specific prompts are provided in Appendix [[B.1]{.ltx_text .ltx_ref_tag}](#A2.SS1 "B.1 MAD Templates ‣ Appendix B Prompt Templates ‣ Appendix ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}.
:::
::::

::::::: {#S2.SS0.SSS0.Px2 .section .ltx_paragraph}
##### Debate *vs.* Voting: Formalization. {#debate-vs.-voting-formalization. .ltx_title .ltx_title_paragraph}

::: {#S2.SS0.SSS0.Px2.p1 .ltx_para}
Let $\mathcal{X}$ denote the input space (e.g., natural language questions), and $\mathcal{Y}$ the output space (e.g., free-form or multiple-choice answers). We consider a set of $N$ language model agents, denoted by $\{a_{1},\dots,a_{N}\}$, where each agent defines a stochastic function $f_{i}:\mathcal{X}\rightarrow\mathcal{Y}$ that produces an initial response ${y}_{i,0}\sim f_{i}(x)$ for input $x\in\mathcal{X}$.
:::

::: {#S2.SS0.SSS0.Px2.p2 .ltx_para}
In the *Majority Voting* setting, the initial responses $\{{y}_{i,0}\}_{i=1}^{N}$ are directly aggregated using a voting function $\mathcal{V}:\mathcal{Y}^{N}\rightarrow\mathcal{Y}$ to obtain the final prediction, typically returning the most frequent answer:

  -- ----------------------------------------------------- --
     $${y}_{0}=\mathcal{V}({y}_{1,0},\dots,{y}_{N,0}).$$   
  -- ----------------------------------------------------- --
:::

::: {#S2.SS0.SSS0.Px2.p3 .ltx_para}
In contrast, *Multi-Agent Debate* introduces $T$ rounds of iterative communication. We formalize the communication structure of debate as an undirected graph $\mathcal{G}$, where each node corresponds to an agent and edges indicate which agents observe one another. At round $t\geq 1$, each agent $a_{i}$ observes responses from a set of neighboring agents at the previous round and updates its answer accordingly. We define the response set of neighbors available to agent $a_{i}$ at round $t$ as:

  -- ------------------------------------------------------------------- --
     $$\mathcal{R}_{i}^{(t)}=\{{y}_{j,t-1}\mid j\in\mathcal{N}(i)\},$$   
  -- ------------------------------------------------------------------- --

where $\mathcal{N}(i)\subseteq\{1,\dots,N\}$ is the index set of neighbors, observable to agent $a_{i}$ ([e.g.]{.ltx_text .ltx_font_italic}, in the fully connected setting, $\mathcal{N}(i)=\{1,\dots,N\}$). The response update is given by:

  -- ---------------------------------------------------------------- --
     $${y}_{i,t}=\mathcal{D}\left(x;\mathcal{R}_{i}^{(t)}\right),$$   
  -- ---------------------------------------------------------------- --

where $\mathcal{D}$ denotes a single round of debate. The iterative debate process over $T$ rounds can be expressed as a function composition:

  -- --------------------------------------------------------------------------------------------------------------------------------- --
     $${y}_{i,T}=(\mathcal{D}\circ\mathcal{D}\circ\cdots\circ\mathcal{D})(x;\mathcal{R}_{i})=\mathcal{D}^{(T)}(x;\mathcal{R}_{i}).$$   
  -- --------------------------------------------------------------------------------------------------------------------------------- --

The final aggregated output after $T$ rounds is:

  -- ---------------------------------------------------------------- --
     $${y}_{T}=\mathcal{V}({y}_{1,T},{y}_{2,T},\cdots,{y}_{N,T}).$$   
  -- ---------------------------------------------------------------- --
:::

::: {#S2.SS0.SSS0.Px2.p4 .ltx_para}
We adopt the simultaneous-talk protocol \[[3](#bib.bib3){.ltx_ref}\], where all agents update in parallel based on the previous round's responses. Following the common setup in prior works, we focus on homogeneous agent settings, [i.e.]{.ltx_text .ltx_font_italic}, all agents share the same underlying model architecture or behavior. This allows us to isolate the effect of inter-agent communication and contrast MAD directly with simple majority voting. Our goal is to contrast the performance of MAD against simple majority voting and assess whether iterative inter-agent communication provides measurable improvements beyond ensembling alone. We will extend to the heterogeneous setting in Section [[6]{.ltx_text .ltx_ref_tag}](#S6 "6 Extended Experiments to General Settings ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}.
:::
:::::::
::::::::::

::::::::::::: {#S3 .section .ltx_section}
### [3 ]{.ltx_tag .ltx_tag_section}Is Debate Really Necessary? A Closer Look at Debate *vs.* Voting {#is-debate-really-necessary-a-closer-look-at-debate-vs.-voting .ltx_title .ltx_title_section}

<figure id="S3.T1" class="ltx_table">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_1">
<div class="ltx_inline-block ltx_figure_panel ltx_align_center ltx_transformed_outer" style="width:433.6pt;height:135.5pt;vertical-align:-66.0pt;">
<span class="ltx_transformed_inner" style="transform:translate(-93.8pt,29.3pt) scale(0.698045210970388,0.698045210970388) ;"> </span>
<table class="ltx_tabular ltx_align_middle">
<tbody>
<tr class="ltx_tr">
<td rowspan="2" class="ltx_td ltx_align_left ltx_border_tt" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">Methods</span></td>
<td colspan="7" class="ltx_td ltx_align_center ltx_border_tt" style="padding-left: 2.5pt; padding-right: 2.5pt">Qwen2.5-7B-Instruct</td>
<td class="ltx_td ltx_border_tt" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">Arithmetics</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">GSM8K</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">MMLU</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">MMLU</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">HellaSwag</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">CommonSense</span></td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">HH-RLHF</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">Average</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
<td class="ltx_td" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
<td class="ltx_td" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">(Pro.Med.)</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">(Form.Log.)</span></td>
<td class="ltx_td" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">QA</span></td>
<td class="ltx_td ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
<td class="ltx_td" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
</tr>
<tr class="ltx_tr">
<td colspan="9" class="ltx_td ltx_align_center ltx_border_t" style="background-color: #F2F2F2; padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold" style="background-color:#F2F2F2;">Single-Agent</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Single-agent baseline</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8140 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .04</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8713 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .00</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7868 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .01</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4905 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .03</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7880 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .01</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8153 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .01</span></td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4773 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .01</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7205</td>
</tr>
<tr class="ltx_tr">
<td colspan="9" class="ltx_td ltx_align_center ltx_border_t" style="background-color: #F2F2F2; padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold" style="background-color:#F2F2F2;">Multi-Agent</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Decentralized MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7600</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8867</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8051</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.5556</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8033</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.8567</span></td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4967</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7377</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Decentralized MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6700</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8533</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8051</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5000</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8000</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8500</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5000</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7112</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Decentralized MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6700</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8333</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8051</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4762</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8000</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8433</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.5067</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7050</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Sparse MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8400</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.9033</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8051</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4762</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7967</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8367</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4733</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7330</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Sparse MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8100</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8833</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.8162</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4365</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7967</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8367</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4733</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7218</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Sparse MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7900</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8700</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8088</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4365</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7900</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8333</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4833</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7160</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Centralized MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4300</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7300</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.8162</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4762</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8100</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.8567</span></td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4667</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6551</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Centralized MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4800</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7367</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.8162</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4603</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8100</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8500</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4733</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6609</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Centralized MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5500</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7200</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8125</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4444</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.8133</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8467</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4833</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6672</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">Majority Voting</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.9900</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.9400</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7941</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5397</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8033</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8300</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4867</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.7691</span></td>
</tr>
</tbody>
</table>
</div>
</div>
<div class="ltx_flex_break">

</div>
<div class="ltx_flex_cell ltx_flex_size_1">
<div class="ltx_inline-block ltx_figure_panel ltx_align_center ltx_transformed_outer" style="width:433.6pt;height:135.5pt;vertical-align:-66.0pt;">
<span class="ltx_transformed_inner" style="transform:translate(-93.8pt,29.3pt) scale(0.698045210970388,0.698045210970388) ;"> </span>
<table class="ltx_tabular ltx_align_middle">
<tbody>
<tr class="ltx_tr">
<td rowspan="2" class="ltx_td ltx_align_left ltx_border_tt" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">Methods</span></td>
<td colspan="7" class="ltx_td ltx_align_center ltx_border_tt" style="padding-left: 2.5pt; padding-right: 2.5pt">Llama3.1-8B-Instruct</td>
<td class="ltx_td ltx_border_tt" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">Arithmetics</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">GSM8K</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">MMLU</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">MMLU</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">HellaSwag</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">CommonSense</span></td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">HH-RLHF</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">Average</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
<td class="ltx_td" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
<td class="ltx_td" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">(Pro.Med.)</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">(Form.Log.)</span></td>
<td class="ltx_td" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">QA</span></td>
<td class="ltx_td ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
<td class="ltx_td" style="padding-left: 2.5pt; padding-right: 2.5pt"></td>
</tr>
<tr class="ltx_tr">
<td colspan="9" class="ltx_td ltx_align_center ltx_border_t" style="background-color: #F2F2F2; padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold" style="background-color:#F2F2F2;">Single-Agent</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Single-agent baseline</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7320 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .03</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7393 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .01</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7441 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .01</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.3794 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .02</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6267 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .03</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6767 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .01</span></td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4440 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .02</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6203</td>
</tr>
<tr class="ltx_tr">
<td colspan="9" class="ltx_td ltx_align_center ltx_border_t" style="background-color: #F2F2F2; padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold" style="background-color:#F2F2F2;">Multi-Agent</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Decentralized MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8200</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7933</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7868</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.5238</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6767</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7267</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5233</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6929</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Decentralized MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8300</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7933</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7684</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5000</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6300</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7033</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5267</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6788</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Decentralized MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8500</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7800</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7463</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5000</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6300</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7000</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5267</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6761</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Sparse MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8500</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8133</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8015</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4683</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6767</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7567</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5267</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6990</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Sparse MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8500</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7967</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7831</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4206</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6233</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7467</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5333</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6791</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Sparse MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.8500</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7667</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7868</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4365</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6233</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7233</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.5400</span></td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6752</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Centralized MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7300</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6400</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6949</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.3810</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6000</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7400</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4800</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6094</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Centralized MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7700</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6200</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6507</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.3730</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6133</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7200</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4900</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6053</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left" style="padding-left: 2.5pt; padding-right: 2.5pt">Centralized MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.7300</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5933</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5846</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.3413</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5800</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.6967</td>
<td class="ltx_td ltx_align_center ltx_border_r" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4433</td>
<td class="ltx_td ltx_align_center" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5670</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">Majority Voting</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.8900</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.8733</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.8199</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt">0.5159</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.6933</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.7800</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt">0.4967</td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-left: 2.5pt; padding-right: 2.5pt"><span class="ltx_text ltx_font_bold">0.7242</span></td>
</tr>
</tbody>
</table>
</div>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 1</span>: </span><span class="ltx_text ltx_font_bold" style="font-size:90%;">Majority Voting vs. Multi-Agent Debate.<span class="ltx_text ltx_font_medium"> Benchmark performances are measured in Accuracy.</span></span></figcaption>
</figure>

::: {#S3.p1 .ltx_para}
Multi-agent debate is often regarded as a promising mechanism for enhancing LLM performance via collaborative deliberation. *But how much of its effectiveness truly comes from the debate itself---and how much is simply due to aggregating multiple answers*? To address this question, we dissect MAD into two components---multi-agent ensembling and inter-agent communication---and present empirical evidence revealing that simple majority voting accounts for most of the observed gains. We begin by explaining the experimental setup in the following section.
:::

::::::: {#S3.SS1 .section .ltx_subsection}
#### [3.1 ]{.ltx_tag .ltx_tag_subsection}Experimental Setup {#experimental-setup .ltx_title .ltx_title_subsection}

:::: {#S3.SS1.SSS0.Px1 .section .ltx_paragraph}
##### Baselines. {#baselines. .ltx_title .ltx_title_paragraph}

::: {#S3.SS1.SSS0.Px1.p1 .ltx_para}
The key distinction among multi-agent debate methods typically lies in the design of the debate function $\mathcal{D}$, particularly in how agents communicate and the roles they assume. To comprehensively evaluate these variations, we consider the following representative approaches: (1) [Decentralized MAD]{.ltx_text .ltx_font_italic} \[[2](#bib.bib2){.ltx_ref}\], where each agent observes all other agents' responses from the previous round. (2) [Sparse MAD]{.ltx_text .ltx_font_italic} \[[10](#bib.bib10){.ltx_ref}\], a variant of Decentralized MAD with a sparse communication topology to enhance efficiency. (3) [Centralized MAD]{.ltx_text .ltx_font_italic} \[[14](#bib.bib14){.ltx_ref}\], where a central agent aggregates peer responses and generates the updated response at each round. (4) [Majority Voting]{.ltx_text .ltx_font_italic}, which selects the final answer by aggregating initial responses from multiple agents [without any debate]{.ltx_text .ltx_font_bold}. This can be viewed as a special case with $T=0$. For all the multi-agent approaches, we adopt $N=5$ in our main comparison and will ablate on the effect of $N$. For single-agent baselines, we average across 5 independent runs.
:::
::::

:::: {#S3.SS1.SSS0.Px2 .section .ltx_paragraph}
##### Benchmarks. {#benchmarks. .ltx_title .ltx_title_paragraph}

::: {#S3.SS1.SSS0.Px2.p1 .ltx_para}
Following previous MAD literature, we focus on solving six natural language question answering tasks by conducting extensive evaluations across benchmark datasets: (1) [Arithmetics]{.ltx_text .ltx_font_italic}, (2) [Mathematical Reasoning]{.ltx_text .ltx_font_italic} (Grade School Math 8k \[[15](#bib.bib15){.ltx_ref}\]), (3) [Factual Question Answering]{.ltx_text .ltx_font_italic} (MMLU Professional Medicine and Formal Logics \[[16](#bib.bib16){.ltx_ref}, [17](#bib.bib17){.ltx_ref}\]), (4) [Natural Langauge Inference]{.ltx_text .ltx_font_italic} (HellaSwag \[[18](#bib.bib18){.ltx_ref}\]), (5) [Commonsense Reasoning]{.ltx_text .ltx_font_italic} (CommonsenseQA \[[19](#bib.bib19){.ltx_ref}\]), and (6) [Alignment Labeling]{.ltx_text .ltx_font_italic} (HH-RLHF \[[20](#bib.bib20){.ltx_ref}\]), where we adopt the "AI labeler alignment" practice \[[21](#bib.bib21){.ltx_ref}\], similar to \[[10](#bib.bib10){.ltx_ref}\]. For fairness in comparison, all baselines are evaluated on the same data subsets. More details are provided in Appendix [[A.2]{.ltx_text .ltx_ref_tag}](#A1.SS2 "A.2 Dataset Details ‣ Appendix A Experimental Details ‣ Appendix ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}.
:::
::::
:::::::

:::::: {#S3.SS2 .section .ltx_subsection}
#### [3.2 ]{.ltx_tag .ltx_tag_subsection}Key Observations {#key-observations .ltx_title .ltx_title_subsection}

::::: {#S3.SS2.SSS0.Px1 .section .ltx_paragraph}
##### Majority voting is surprisingly strong. {#majority-voting-is-surprisingly-strong. .ltx_title .ltx_title_paragraph}

::: {#S3.SS2.SSS0.Px1.p1 .ltx_para}
In Table [[1]{.ltx_text .ltx_ref_tag}](#S3.T1 "Table 1 ‣ 3 Is Debate Really Necessary? A Closer Look at Debate vs. Voting ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}, we compare the performance of single-agent, MAD, and majority voting approaches across seven benchmark datasets using the Qwen2.5-7B-Instruct \[[22](#bib.bib22){.ltx_ref}\] and Llama3.1-8B-Instruct \[[23](#bib.bib23){.ltx_ref}\] models. Consistent with the typical choice of existing literature, we compare 2- and 3-round, along with a prolonged 5-round debate setting among five agents. Interestingly, while MAD consistently outperforms the single-agent baseline, it does not reliably surpass the much simpler majority voting strategy. *Notably, in most cases, majority voting performs on par with MAD*. To further assess the impact of model capacity, we additionally evaluate the more capable Qwen2.5-32B-Instruct model in Section [[6]{.ltx_text .ltx_ref_tag}](#S6 "6 Extended Experiments to General Settings ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}. Although overall performance improves in the MAD setting, the majority voting strategy continues to account for most of the performance gains. These findings suggest that the effectiveness of MAD is largely driven by model ensembling, rather than the iterative debate process itself.
:::

<figure id="S3.F3" class="ltx_figure ltx_align_floatright">
<img src="/html/2508.17536/assets/x3.png" id="S3.F3.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="224" height="156" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 3</span>: </span><span class="ltx_text" style="font-size:90%;">Accuracy improves with more agents.</span></figcaption>
</figure>

::: {#S3.SS2.SSS0.Px1.p2 .ltx_para}
To gain deeper insight into the effect of MAD components, we present an ablation study in Figure [[3]{.ltx_text .ltx_ref_tag}](#S3.F3 "Figure 3 ‣ Majority voting is surprisingly strong. ‣ 3.2 Key Observations ‣ 3 Is Debate Really Necessary? A Closer Look at Debate vs. Voting ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}. It illustrates the effect of varying the number of Qwen2.5 agents participating in each round of debate, from $N=1$ to $N=5$. Overall, increasing the number of agents generally leads to improved performance. The trend suggests that MAD's effectiveness may stem primarily from the ensemble effect of multiple agents. Our next section formalizes our observations.
:::
:::::
::::::
:::::::::::::

:::::::::::::::::::::::::: {#S4 .section .ltx_section}
### [4 ]{.ltx_tag .ltx_tag_section}Theoretical Analysis {#theoretical-analysis .ltx_title .ltx_title_section}

::: {#S4.p1 .ltx_para}
To better understand the dynamics underlying our empirical findings in Section [[3]{.ltx_text .ltx_ref_tag}](#S3 "3 Is Debate Really Necessary? A Closer Look at Debate vs. Voting ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}, we now turn to a formal analysis of multi-agent debate and majority voting. Our theoretical framework allows us to capture how agent uncertainty and belief updates shape the collective decision-making in both debate and voting, grounded in Bayesian principles. Specifically, for a given input question, we consider a population of $N$ agents, each generating a response from a finite set $\mathcal{A}$, which may represent multiple-choice question options or a set of plausible completions for open-ended tasks. We model each agent as an idealized generative model governed by a Dirichlet-Compound-Multinomial (DCM) distribution. *This closely mirrors how practical LLM systems produce different outputs for the same question due to uncertainty and sampling variability*. In particular, the Dirichlet prior captures the agent's internal belief over possible answers, while the Multinomial models the stochastic generation process ([e.g.]{.ltx_text .ltx_font_italic}, via temperature or nucleus sampling). This distribution is thus a natural choice because it encapsulates both internal uncertainty and output randomness, while also providing a principled Bayesian framework for belief updates across debate rounds---enabling analytical study of dynamics during the debate process. Below we provide the mathematical details of the DCM model.
:::

::::: {#S4.SS0.SSS0.Px1 .section .ltx_paragraph}
##### Definition 1. (Agent Response Generation via DCM) {#definition-1.-agent-response-generation-via-dcm .ltx_title .ltx_title_paragraph}

::: {#S4.SS0.SSS0.Px1.p1 .ltx_para}
[At round $t$, each agent $i$ is associated with a belief vector $\boldsymbol{\alpha}_{i,t}=(\alpha_{i,t}^{(1)},\ldots,\alpha_{i,t}^{(K)})\in\mathbb{R}^{K}_{+}$, where each entry $\alpha_{i,t}^{(k)}$ reflects the agent's belief in response option $k\in\mathcal{A}$. To generate a response $y_{i,t}$, the agent follows a two-step process:]{.ltx_text .ltx_font_italic}

  -- -------------------------------------------------------------------------------------- -------------------------------------------------------------------------------------------- --
     [(Belief sampling)]{.ltx_text .ltx_markedasmath .ltx_font_bold .ltx_font_italic}       $\displaystyle\boldsymbol{\theta}_{i,t}\sim\mathrm{Dirichlet}(\boldsymbol{\alpha}_{i,t}),$   
     [(Response generation)]{.ltx_text .ltx_markedasmath .ltx_font_bold .ltx_font_italic}   $\displaystyle y_{i,t}\sim\mathrm{Categorical}(\boldsymbol{\theta}_{i,t}).$                  
  -- -------------------------------------------------------------------------------------- -------------------------------------------------------------------------------------------- --

[The marginal probability of generating any particular response $y_{i,t}\in\mathcal{A}$---after integrating out the randomness in $\boldsymbol{\theta}_{i,t}$---is given by $P(y_{i,t}=k\mid\boldsymbol{\alpha}_{i,t})=\alpha_{i,t}^{(k)}/\sum_{j\in\mathcal{A}}\alpha_{i,t}^{(j)}.$]{.ltx_text .ltx_font_italic}
:::

::: {#S4.SS0.SSS0.Px1.p2 .ltx_para}
Before analyzing the dynamics of debate, we first consider the base case where agents respond independently without interaction. The following characterizes the success probability of majority voting under this condition, based solely on the agents' homogeneous initial beliefs, $\boldsymbol{\alpha}_{i,0}=\boldsymbol{\alpha}=(\alpha^{(1)},\ldots,\alpha^{(K)}$). Without loss of generality, let answer index 1 be the correct option. We assume that the correct answer has the largest belief $\alpha^{(1)}$, and all the other beliefs are ordered such that $\alpha^{(1)}>\alpha^{(2)}\geq\cdots\geq\alpha^{(K)}$.
:::
:::::

:::: {#S4.SS0.SSS0.Px2 .section .ltx_paragraph}
##### Theorem 1. (Majority Voting Success Probability) {#theorem-1.-majority-voting-success-probability .ltx_title .ltx_title_paragraph}

::: {#S4.SS0.SSS0.Px2.p1 .ltx_para}
[Let $\bar{\boldsymbol{\theta}}=(\bar{\theta}^{(1)},\ldots,\bar{\theta}^{(K)})=\boldsymbol{\alpha}/\sum_{j=1}^{K}\alpha_{j}$ denote the mean of the Dirichlet distribution, $\mathrm{Dirichlet}(\alpha)$, and define the margin $\Delta:=\bar{\theta_{1}}-\bar{\theta_{2}}$. If $N>K/\Delta^{2}$, then the probability that majority voting selects the answer 1 is lower bounded as:]{.ltx_text .ltx_font_italic}

  -- ------------------------------------------------------------------------------------------------------------------------- --
     $$\mathbb{P}(y_{\mathrm{mv}}=1)\geq 1-\exp\left(-N\left(\frac{\Delta}{\sqrt{K}}-\frac{1}{\sqrt{N}}\right)^{2}\right).$$   
  -- ------------------------------------------------------------------------------------------------------------------------- --
:::
::::

::::: {#S4.SS0.SSS0.Px3 .section .ltx_paragraph}
##### Remark 1. {#remark-1. .ltx_title .ltx_title_paragraph}

::: {#S4.SS0.SSS0.Px3.p1 .ltx_para}
This result highlights the *magnifying effect* of majority voting: even when the correct answer is only marginally more likely than the alternative answers, the lower bound on the success probability asymptotically approaches 1 as $N$ scales. Notably, this holds even when $\bar{\theta}_{1}\ll\frac{1}{2}$, as long as it remains the most probable choice option. In practice, we recognize that MAD systems often operate with a small number of agents due to computational constraints. We provide a sharper analysis specialized to this realistic regime that applies to arbitrary $N$ in Theorem 1.A (Appendix [[D]{.ltx_text .ltx_ref_tag}](#A4 "Appendix D Special Case of Theorem 1 ‣ Appendix ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}), without constraining its size. This complementary result provides insights into the reliability of majority voting in more practical, resource-constrained settings.
:::

::: {#S4.SS0.SSS0.Px3.p2 .ltx_para}
Next, we analyze the multi-agent debate performance by formalizing how each agent's belief $\boldsymbol{\alpha}_{i,t}$ evolves through debate. Specifically, each agent observes its neighbors' responses and performs a Bayesian posterior update of its belief accordingly.
:::
:::::

::::: {#S4.SS0.SSS0.Px4 .section .ltx_paragraph}
##### Definition 2. (Bayesian Belief Update from Neighbor Responses) {#definition-2.-bayesian-belief-update-from-neighbor-responses .ltx_title .ltx_title_paragraph}

::: {#S4.SS0.SSS0.Px4.p1 .ltx_para}
[Let $\{y_{j,t-1}\mid j\in\mathcal{N}(i)\}$ be the set of responses observed by agent $i$ from its neighbors $\mathcal{N}(i)$ at round $t$. These responses induce a count vector $\mathbf{c}_{i,t}=(c_{i,t}^{(1)},\ldots,c_{i,t}^{(K)})\in\mathbb{N}^{K}$, where $c_{i,t}^{(k)}$ denotes the number of neighbors who selected response $k$. Then, the agent updates its Dirichlet parameter as: $\boldsymbol{\alpha}_{i,t}=\boldsymbol{\alpha}_{i,t-1}+\mathbf{c}_{i,t}.$]{.ltx_text .ltx_font_italic}
:::

::: {#S4.SS0.SSS0.Px4.p2 .ltx_para}
Under this formulation, each round of multi-agent debate corresponds to a Bayesian update step under the conjugacy of the Dirichlet-Multinomial model.
:::
:::::

::::: {#S4.SS0.SSS0.Px5 .section .ltx_paragraph}
##### Lemma 1. (Bayesian Conjugacy in Multi-Agent Debate) {#lemma-1.-bayesian-conjugacy-in-multi-agent-debate .ltx_title .ltx_title_paragraph}

::: {#S4.SS0.SSS0.Px5.p1 .ltx_para}
[At round $t$, after observing responses from its neighbors $\mathcal{N}(i)$, the agent $i$ aggregates these into a count vector $\mathbf{c}_{i,t}$ as in Definition 2. Then, by Bayesian conjugacy, the posterior distribution over $\boldsymbol{\theta}_{i,t}$ remains Dirichlet:]{.ltx_text .ltx_font_italic}

  -- ------------------------------------------------------------------------------------------------------------------------------------------ --
     $$\boldsymbol{\theta}_{i,t}\mid\{y_{j,t-1}\}_{j\in\mathcal{N}(i)}\sim\mathrm{Dirichlet}(\boldsymbol{\alpha}_{i,t-1}+\mathbf{c}_{i,t}).$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------ --
:::

::: {#S4.SS0.SSS0.Px5.p2 .ltx_para}
As a result, we can prove that the evolution of agent beliefs forms a martingale process, with full proof in Appendix [[C]{.ltx_text .ltx_ref_tag}](#A3 "Appendix C Proofs ‣ Appendix ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}.
:::
:::::

:::: {#S4.SS0.SSS0.Px6 .section .ltx_paragraph}
##### Theorem 2. (Martingale Behavior of Multi-Agent Debate) {#theorem-2.-martingale-behavior-of-multi-agent-debate .ltx_title .ltx_title_paragraph}

::: {#S4.SS0.SSS0.Px6.p1 .ltx_para}
[For agent $i$, let $p_{t}:=\bar{{\theta}}_{i,t}^{(1)}$ denote its belief in the correct answer at debate round $t$. Under Bayesian conjugacy,]{.ltx_text .ltx_font_italic}

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --
     $$p_{t}:=\bar{\theta}_{i,t}^{(1)}=\frac{\alpha_{i,t}^{(1)}}{\sum_{j=1}^{K}\alpha_{i,t}^{(j)}}=\frac{\alpha_{i,t-1}^{(1)}+c_{i,t}^{(1)}}{\sum_{j=1}^{K}(\alpha_{i,t-1}^{(j)}+c_{i,t}^{(j)})}.$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --

[Then, sequence $\{p_{t}\}_{t\geq 0}$ forms a martingale. That is, the expected belief at the next round equals the current belief:]{.ltx_text .ltx_font_italic}

  -- ----------------------------------------------------------------------------- --
     $$\mathbb{E}[p_{t}\mid p_{t-1},\ldots,p_{0}]=p_{t-1},\;\forall_{t\geq 0}.$$   
  -- ----------------------------------------------------------------------------- --
:::

<figure id="S4.F4" class="ltx_figure">
<img src="/html/2508.17536/assets/x4.png" id="S4.F4.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="660" height="283" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Figure 4</span>: </span><span class="ltx_text" style="font-size:90%;">Martingale process of the mean agent accuracy across debate rounds.</span></figcaption>
</figure>
::::

:::: {#S4.SS0.SSS0.Px7 .section .ltx_paragraph}
##### Theoretical insights: Majority vote does essentially all the work. {#theoretical-insights-majority-vote-does-essentially-all-the-work. .ltx_title .ltx_title_paragraph}

::: {#S4.SS0.SSS0.Px7.p1 .ltx_para}
This theorem highlights a fundamental property of multi-agent debate under Bayesian update of the DCM model: the agent's belief in the correct answer behaves like a [martingale]{.ltx_text .ltx_font_italic}---that is, its expected value remains unchanged across rounds. This result is closely related to the classical Pólya Urn scheme \[[24](#bib.bib24){.ltx_ref}\]. Intuitively, this means that debate itself does not systematically improve or degrade an agent's belief on average; instead, belief updates are driven entirely by the stochastic influence of peer responses. While some debate trajectories may lead to stronger belief in the correct answer ([i.e.]{.ltx_text .ltx_font_italic}, [correction]{.ltx_text .ltx_font_italic}), others may lead to weaker belief ([i.e.]{.ltx_text .ltx_font_italic}, [subversion]{.ltx_text .ltx_font_italic}). While these local fluctuations affect the posterior counts, the expected belief in the correct answer remains equal to the initial $p_{0}=\bar{\theta}_{i,0}^{(1)}$, without any debate. This implies that, under our theoretical model, debate alone does not necessarily improve the initial accuracy--*majority voting accounts for the primary performance gains*. Our theory thus aligns well with our empirical findings in Section [[3]{.ltx_text .ltx_ref_tag}](#S3 "3 Is Debate Really Necessary? A Closer Look at Debate vs. Voting ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}.
:::
::::

:::: {#S4.SS0.SSS0.Px8 .section .ltx_paragraph}
##### Martingale behavior is supported empirically. {#martingale-behavior-is-supported-empirically. .ltx_title .ltx_title_paragraph}

::: {#S4.SS0.SSS0.Px8.p1 .ltx_para}
We empirically examine whether the sequence $\{p_{t}\}_{t\geq 0}$ exhibits martingale behavior. For each benchmark and debate round $t$, we estimate $p_{t}$ as the mean accuracy of the five agents[^1^[[^1^[1]{.ltx_tag .ltx_tag_note}Note that this quantity is different from the accuracy in Table [[1]{.ltx_text .ltx_ref_tag}](#S3.T1 "Table 1 ‣ 3 Is Debate Really Necessary? A Closer Look at Debate vs. Voting ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}, which is the ratio of correct responses [after]{.ltx_text .ltx_font_italic} voting.]{.ltx_note_content}]{.ltx_note_outer}]{#footnote1 .ltx_note .ltx_role_footnote}. As shown in Figure [[4]{.ltx_text .ltx_ref_tag}](#S4.F4 "Figure 4 ‣ Theorem 2. (Martingale Behavior of Multi-Agent Debate) ‣ 4 Theoretical Analysis ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}, *the resulting trajectories are essentially flat, which is consistent with the theoretical property that the expected value of a martingale remains unchanged over time*. Raw values for mean accuracy are provided in Table [[6]{.ltx_text .ltx_ref_tag}](#A5.T6 "Table 6 ‣ Appendix E Martingale Process Empirical Investigation ‣ Appendix ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref} (Appendix [[E]{.ltx_text .ltx_ref_tag}](#A5 "Appendix E Martingale Process Empirical Investigation ‣ Appendix ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}).
:::
::::

:::: {#S4.SS0.SSS0.Px9 .section .ltx_paragraph}
##### Generalized interpretation of the Bayesian update step. {#generalized-interpretation-of-the-bayesian-update-step. .ltx_title .ltx_title_paragraph}

::: {#S4.SS0.SSS0.Px9.p1 .ltx_para}
In Lemma 1 and Theorem 2, we define the update dynamics in terms of the count vector $\mathbf{c}_{i,t}$. Our framework can be generalized to open-ended tasks and capture the heterogeneous influence of each agent's response, with minor interpretational adjustments. For open-ended tasks, although responses are not strictly countable in categorical form, they can be represented in distributional or similarity-based spaces. Here, the count vector can be viewed more broadly, [e.g.]{.ltx_text .ltx_font_italic}, as a soft histogram over clustered response types, an embedding-based semantic agreement measure, or a weighted similarity score between textual outputs. In such settings, "consensus" is better defined semantically rather than symbolically: if agents independently produce explanations or rationales that are semantically aligned, this can be considered consensus even when surface forms differ. This notion can be operationalized using thresholds on embedding similarity, Levenshtein distance, or overlap in reasoning chains, among other metrics.
:::
::::
::::::::::::::::::::::::::

::::::::: {#S5 .section .ltx_section}
### [5 ]{.ltx_tag .ltx_tag_section}How Does Theory Inform Improved Design of MAD? {#how-does-theory-inform-improved-design-of-mad .ltx_title .ltx_title_section}

::: {#S5.p1 .ltx_para}
The martingale property in our theory reflects a neutral expectation over time, underscoring that without additional bias towards correct signals, debate alone does not guarantee convergence to the truth. Any additional benefit arises from local asymmetries in the observed stochastic process $\{p_{t}\}_{t\geq 0}$. Hence, to improve the effectiveness of MAD, we explore alternative designs that facilitate correction and/or suppress subversion in the debate process.
:::

::::: {#S5.SS1 .section .ltx_subsection}
#### [5.1 ]{.ltx_tag .ltx_tag_subsection}Belief Update by Biasing Towards Correct Signal {#belief-update-by-biasing-towards-correct-signal .ltx_title .ltx_title_subsection}

::: {#S5.SS1.p1 .ltx_para}
To investigate how targeted intervention in the belief update process can promote convergence to the correct answer, we first consider an oracle-style method that explicitly biases updates toward the correct signal. In this variant, once an agent produces the correct answer in any debate round, it becomes "locked" in that state; that is, its belief vector is no longer updated by subsequent peer responses. Formally, if agent $i$ outputs the correct answer at round $t$, we use this response in all subsequent rounds $t^{\prime}>t$. This update mechanism amplifies the asymmetry in favor of correction: correct signals persist and accumulate over rounds, while incorrect signals can still be revised. Consequently, the system's dynamics depart from the neutral martingale behavior discussed in Section [[4]{.ltx_text .ltx_ref_tag}](#S4 "4 Theoretical Analysis ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref} and instead exhibit a directional drift toward the correct answer.
:::

::: {#S5.SS1.p2 .ltx_para}
We refer to this method as [MAD-oracle]{.ltx_text .ltx_font_bold}, and report its performance in Table [[2]{.ltx_text .ltx_ref_tag}](#S5.T2 "Table 2 ‣ 5.1 Belief Update by Biasing Towards Correct Signal ‣ 5 How Does Theory Inform Improved Design of MAD? ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}. This variant yields substantial improvements over standard MAD, and always surpasses the Majority Voting baseline by a large margin. For instance, in Decentralized MAD with $T=5$ rounds, accuracy on MMLU (Form. Log.) increases from 0.5000 to 0.6825. Although this approach is not feasible in practice---since the true answer is not available---it reveals an upper bound on the benefit achievable by incorporating bias toward correct signals in the belief update process. In the next subsection, we investigate a more realistic alternative that aims to suppress subversion without direct access to ground truth.
:::

<figure id="S5.T2" class="ltx_table">
<div class="ltx_inline-block ltx_align_center ltx_transformed_outer" style="width:433.6pt;height:267.4pt;vertical-align:-132.0pt;">
<span class="ltx_transformed_inner" style="transform:translate(-97.2pt,59.9pt) scale(0.690507580070741,0.690507580070741) ;"> </span>
<table class="ltx_tabular ltx_align_middle">
<tbody>
<tr class="ltx_tr">
<td rowspan="2" class="ltx_td ltx_align_left ltx_border_tt"><span class="ltx_text ltx_font_bold">Methods</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">Arithmetics</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">GSM8K</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">MMLU</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">MMLU</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">HellaSwag</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">CommonSense</span></td>
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_tt"><span class="ltx_text ltx_font_bold">HH-RLHF</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">Average</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td"></td>
<td class="ltx_td"></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">(Pro.Med.)</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">(Form.Log.)</span></td>
<td class="ltx_td"></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">QA</span></td>
<td class="ltx_td ltx_border_r"></td>
<td class="ltx_td"></td>
</tr>
<tr class="ltx_tr">
<td colspan="9" class="ltx_td ltx_align_center ltx_border_t">Decentralized MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-vanilla</td>
<td class="ltx_td ltx_align_center">0.7600</td>
<td class="ltx_td ltx_align_center">0.8867</td>
<td class="ltx_td ltx_align_center">0.8051</td>
<td class="ltx_td ltx_align_center">0.5238</td>
<td class="ltx_td ltx_align_center">0.8033</td>
<td class="ltx_td ltx_align_center">0.8567</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4967</td>
<td class="ltx_td ltx_align_center">0.7332</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-Conformist</td>
<td class="ltx_td ltx_align_center">0.9200</td>
<td class="ltx_td ltx_align_center">0.9200</td>
<td class="ltx_td ltx_align_center">0.8015</td>
<td class="ltx_td ltx_align_center">0.5397</td>
<td class="ltx_td ltx_align_center">0.8000</td>
<td class="ltx_td ltx_align_center">0.8600</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4967</td>
<td class="ltx_td ltx_align_center">0.7625</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-Follower</td>
<td class="ltx_td ltx_align_center">0.9200</td>
<td class="ltx_td ltx_align_center">0.9300</td>
<td class="ltx_td ltx_align_center">0.8088</td>
<td class="ltx_td ltx_align_center">0.5317</td>
<td class="ltx_td ltx_align_center">0.8100</td>
<td class="ltx_td ltx_align_center">0.8500</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4900</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">0.7629</span></td>
</tr>
<tr class="ltx_tr" style="background-color:#F2F2F2;">
<td class="ltx_td ltx_align_left"><span class="ltx_text" style="background-color:#F2F2F2;">MAD-oracle</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9400</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9667</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8897</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.6587</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8333</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8933</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F2F2F2;">0.5833</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8236</span></td>
</tr>
<tr class="ltx_tr">
<td colspan="9" class="ltx_td ltx_align_center ltx_border_t">Decentralized MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-vanilla</td>
<td class="ltx_td ltx_align_center">0.6700</td>
<td class="ltx_td ltx_align_center">0.8533</td>
<td class="ltx_td ltx_align_center">0.8051</td>
<td class="ltx_td ltx_align_center">0.5000</td>
<td class="ltx_td ltx_align_center">0.8000</td>
<td class="ltx_td ltx_align_center">0.8500</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.5000</td>
<td class="ltx_td ltx_align_center">0.7112</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-Conformist</td>
<td class="ltx_td ltx_align_center">0.9000</td>
<td class="ltx_td ltx_align_center">0.9133</td>
<td class="ltx_td ltx_align_center">0.8015</td>
<td class="ltx_td ltx_align_center">0.5635</td>
<td class="ltx_td ltx_align_center">0.8033</td>
<td class="ltx_td ltx_align_center">0.8567</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4933</td>
<td class="ltx_td ltx_align_center">0.7617</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-Follower</td>
<td class="ltx_td ltx_align_center">0.9100</td>
<td class="ltx_td ltx_align_center">0.9200</td>
<td class="ltx_td ltx_align_center">0.8015</td>
<td class="ltx_td ltx_align_center">0.5476</td>
<td class="ltx_td ltx_align_center">0.8067</td>
<td class="ltx_td ltx_align_center">0.8567</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.5000</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">0.7632</span></td>
</tr>
<tr class="ltx_tr" style="background-color:#F2F2F2;">
<td class="ltx_td ltx_align_left"><span class="ltx_text" style="background-color:#F2F2F2;">MAD-oracle</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9400</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9667</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8897</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.6746</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8333</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8933</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F2F2F2;">0.5833</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8259</span></td>
</tr>
<tr class="ltx_tr">
<td colspan="9" class="ltx_td ltx_align_center ltx_border_t">Decentralized MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-vanilla</td>
<td class="ltx_td ltx_align_center">0.6700</td>
<td class="ltx_td ltx_align_center">0.8333</td>
<td class="ltx_td ltx_align_center">0.8051</td>
<td class="ltx_td ltx_align_center">0.5000</td>
<td class="ltx_td ltx_align_center">0.8000</td>
<td class="ltx_td ltx_align_center">0.8433</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.5067</td>
<td class="ltx_td ltx_align_center">0.7084</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-Conformist</td>
<td class="ltx_td ltx_align_center">0.8900</td>
<td class="ltx_td ltx_align_center">0.9133</td>
<td class="ltx_td ltx_align_center">0.8088</td>
<td class="ltx_td ltx_align_center">0.5079</td>
<td class="ltx_td ltx_align_center">0.8000</td>
<td class="ltx_td ltx_align_center">0.8500</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4967</td>
<td class="ltx_td ltx_align_center">0.7524</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-Follower</td>
<td class="ltx_td ltx_align_center">0.9000</td>
<td class="ltx_td ltx_align_center">0.9133</td>
<td class="ltx_td ltx_align_center">0.7978</td>
<td class="ltx_td ltx_align_center">0.5397</td>
<td class="ltx_td ltx_align_center">0.8033</td>
<td class="ltx_td ltx_align_center">0.8533</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4967</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">0.7577</span></td>
</tr>
<tr class="ltx_tr" style="background-color:#F2F2F2;">
<td class="ltx_td ltx_align_left"><span class="ltx_text" style="background-color:#F2F2F2;">MAD-oracle</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9400</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9667</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8897</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.6825</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8333</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8933</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F2F2F2;">0.5967</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8289</span></td>
</tr>
<tr class="ltx_tr">
<td colspan="9" class="ltx_td ltx_align_center ltx_border_t">Sparse MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-vanilla</td>
<td class="ltx_td ltx_align_center">0.8400</td>
<td class="ltx_td ltx_align_center">0.9033</td>
<td class="ltx_td ltx_align_center">0.8051</td>
<td class="ltx_td ltx_align_center">0.4683</td>
<td class="ltx_td ltx_align_center">0.7967</td>
<td class="ltx_td ltx_align_center">0.8367</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4733</td>
<td class="ltx_td ltx_align_center">0.7319</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-Conformist</td>
<td class="ltx_td ltx_align_center">0.9100</td>
<td class="ltx_td ltx_align_center">0.9233</td>
<td class="ltx_td ltx_align_center">0.8015</td>
<td class="ltx_td ltx_align_center">0.5238</td>
<td class="ltx_td ltx_align_center">0.8000</td>
<td class="ltx_td ltx_align_center">0.8300</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4833</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">0.7531</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-Follower</td>
<td class="ltx_td ltx_align_center">0.9200</td>
<td class="ltx_td ltx_align_center">0.9233</td>
<td class="ltx_td ltx_align_center">0.7941</td>
<td class="ltx_td ltx_align_center">0.5079</td>
<td class="ltx_td ltx_align_center">0.8000</td>
<td class="ltx_td ltx_align_center">0.8267</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4833</td>
<td class="ltx_td ltx_align_center">0.7508</td>
</tr>
<tr class="ltx_tr" style="background-color:#F2F2F2;">
<td class="ltx_td ltx_align_left"><span class="ltx_text" style="background-color:#F2F2F2;">MAD-oracle</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9200</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9733</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9007</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.6111</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8267</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9000</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F2F2F2;">0.6333</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8236</span></td>
</tr>
<tr class="ltx_tr">
<td colspan="9" class="ltx_td ltx_align_center ltx_border_t">Sparse MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-vanilla</td>
<td class="ltx_td ltx_align_center">0.8100</td>
<td class="ltx_td ltx_align_center">0.8833</td>
<td class="ltx_td ltx_align_center">0.8162</td>
<td class="ltx_td ltx_align_center">0.4206</td>
<td class="ltx_td ltx_align_center">0.7967</td>
<td class="ltx_td ltx_align_center">0.8367</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4733</td>
<td class="ltx_td ltx_align_center">0.7195</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-Conformist</td>
<td class="ltx_td ltx_align_center">0.9100</td>
<td class="ltx_td ltx_align_center">0.9233</td>
<td class="ltx_td ltx_align_center">0.8125</td>
<td class="ltx_td ltx_align_center">0.5000</td>
<td class="ltx_td ltx_align_center">0.8033</td>
<td class="ltx_td ltx_align_center">0.8367</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4633</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">0.7499</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-Follower</td>
<td class="ltx_td ltx_align_center">0.9100</td>
<td class="ltx_td ltx_align_center">0.9267</td>
<td class="ltx_td ltx_align_center">0.8015</td>
<td class="ltx_td ltx_align_center">0.5159</td>
<td class="ltx_td ltx_align_center">0.8000</td>
<td class="ltx_td ltx_align_center">0.8267</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4667</td>
<td class="ltx_td ltx_align_center">0.7496</td>
</tr>
<tr class="ltx_tr" style="background-color:#F2F2F2;">
<td class="ltx_td ltx_align_left"><span class="ltx_text" style="background-color:#F2F2F2;">MAD-oracle</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9200</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9733</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9007</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.6429</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8267</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.9000</span></td>
<td class="ltx_td ltx_align_center ltx_border_r"><span class="ltx_text" style="background-color:#F2F2F2;">0.6467</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text" style="background-color:#F2F2F2;">0.8300</span></td>
</tr>
<tr class="ltx_tr">
<td colspan="9" class="ltx_td ltx_align_center ltx_border_t">Sparse MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-vanilla</td>
<td class="ltx_td ltx_align_center">0.7900</td>
<td class="ltx_td ltx_align_center">0.8700</td>
<td class="ltx_td ltx_align_center">0.8088</td>
<td class="ltx_td ltx_align_center">0.4365</td>
<td class="ltx_td ltx_align_center">0.7900</td>
<td class="ltx_td ltx_align_center">0.8333</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4833</td>
<td class="ltx_td ltx_align_center">0.7160</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-Conformist</td>
<td class="ltx_td ltx_align_center">0.9200</td>
<td class="ltx_td ltx_align_center">0.9200</td>
<td class="ltx_td ltx_align_center">0.8162</td>
<td class="ltx_td ltx_align_center">0.4444</td>
<td class="ltx_td ltx_align_center">0.7967</td>
<td class="ltx_td ltx_align_center">0.8333</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4767</td>
<td class="ltx_td ltx_align_center">0.7439</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">MAD-Follower</td>
<td class="ltx_td ltx_align_center">0.9200</td>
<td class="ltx_td ltx_align_center">0.9233</td>
<td class="ltx_td ltx_align_center">0.8125</td>
<td class="ltx_td ltx_align_center">0.5000</td>
<td class="ltx_td ltx_align_center">0.8033</td>
<td class="ltx_td ltx_align_center">0.8333</td>
<td class="ltx_td ltx_align_center ltx_border_r">0.4767</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">0.7527</span></td>
</tr>
<tr class="ltx_tr" style="background-color:#F2F2F2;">
<td class="ltx_td ltx_align_left ltx_border_bb"><span class="ltx_text" style="background-color:#F2F2F2;">MAD-oracle</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text" style="background-color:#F2F2F2;">0.9400</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text" style="background-color:#F2F2F2;">0.9767</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text" style="background-color:#F2F2F2;">0.9007</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text" style="background-color:#F2F2F2;">0.6508</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text" style="background-color:#F2F2F2;">0.8267</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text" style="background-color:#F2F2F2;">0.9000</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r"><span class="ltx_text" style="background-color:#F2F2F2;">0.6600</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb"><span class="ltx_text" style="background-color:#F2F2F2;">0.8364</span></td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 2</span>: </span><span class="ltx_text ltx_font_bold" style="font-size:90%;">Improved design of MAD<span class="ltx_text ltx_font_medium">. Model is Qwen2.5-7B-Instruct.</span></span></figcaption>
</figure>
:::::

:::: {#S5.SS2 .section .ltx_subsection}
#### [5.2 ]{.ltx_tag .ltx_tag_subsection}Belief Update Guided by the Majority Vote {#belief-update-guided-by-the-majority-vote .ltx_title .ltx_title_subsection}

::: {#S5.SS2.p1 .ltx_para}
In practice, one would not have access to the oracle setting where correct answers are preserved by assumption. To approximate this behavior, we introduce simple modifications to the MAD update rule that leverage the positive signals from majority voting. Our design rationale is guided by our theoretical analysis---which shows that majority voting provides a more reliable estimate of the correct answer than any single agent, as it aggregates marginal advantages across the population. This suggests that using the majority response as a proxy for the ground truth can help steer belief updates in the right direction---effectively biasing the system toward correction without needing oracle access. To explore this idea, we propose two lightweight interventions that incorporate the majority vote into agents' belief dynamics. Specifically, we evaluate two strategies: (1) [MAD-Conformist]{.ltx_text .ltx_font_bold}: if an agent's response matches the majority vote in the previous round, it retains that response; (2) [MAD-Follower]{.ltx_text .ltx_font_bold}: with 30% probability, the agent adopts the majority response from the previous round, and otherwise samples a new one. As shown in Table [[2]{.ltx_text .ltx_ref_tag}](#S5.T2 "Table 2 ‣ 5.1 Belief Update by Biasing Towards Correct Signal ‣ 5 How Does Theory Inform Improved Design of MAD? ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}, these strategies consistently outperform the MAD-vanilla baseline. While they do not reach the oracle's upper bound performance, they demonstrate that simple, theory-informed modifications can yield meaningful improvements---pointing to a promising direction for future work to close the gap.
:::
::::
:::::::::

:::::::::: {#S6 .section .ltx_section}
### [6 ]{.ltx_tag .ltx_tag_section}Extended Experiments to General Settings {#extended-experiments-to-general-settings .ltx_title .ltx_title_section}

::: {#S6.p1 .ltx_para}
In this section, we broaden the scope of our investigation in Section [[3]{.ltx_text .ltx_ref_tag}](#S3 "3 Is Debate Really Necessary? A Closer Look at Debate vs. Voting ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref} to more general settings, evaluating whether the key observation (*i.e.*, majority vote is as effective as debate) holds on larger model size, heterogeneous agents, and open-ended question formats.
:::

:::: {#S6.SS0.SSS0.Px1 .section .ltx_paragraph}
##### Consistent observations in a larger and more capable model. {#consistent-observations-in-a-larger-and-more-capable-model. .ltx_title .ltx_title_paragraph}

::: {#S6.SS0.SSS0.Px1.p1 .ltx_para}
To assess the generality of our findings, we extend our evaluation to more capable language models. Specifically, we test our setup on the Qwen2.5-32B-Instruct model \[[22](#bib.bib22){.ltx_ref}\], using two representative tasks: GSM8K and HellaSwag. As shown in Table [[4]{.ltx_text .ltx_ref_tag}](#S6.T4 "Table 4 ‣ Heterogeneous Agents. ‣ 6 Extended Experiments to General Settings ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}, the results confirm our earlier observations that the performance of Majority Voting remains comparable to that of multi-agent methods. This suggests that our claim is not limited to smaller models, but also holds in high-capacity LLMs.
:::
::::

:::: {#S6.SS0.SSS0.Px2 .section .ltx_paragraph}
##### Heterogeneous Agents. {#heterogeneous-agents. .ltx_title .ltx_title_paragraph}

::: {#S6.SS0.SSS0.Px2.p1 .ltx_para}
While our primary focus has been on homogeneous agent settings, an important question remains: Do our findings also extend to heterogeneous agent configurations? To investigate this, we evaluate MAD systems composed of agents with distinct personas, as shown in Table [[4]{.ltx_text .ltx_ref_tag}](#S6.T4 "Table 4 ‣ Heterogeneous Agents. ‣ 6 Extended Experiments to General Settings ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}. Following the optimal persona sets identified for "college mathematics" and "clinical knowledge" via the "agent selection algorithm" introduced in \[[9](#bib.bib9){.ltx_ref}\], we construct diverse agent roles for each task. For GSM8K, the team includes a general-purpose "Assistant" alongside specialized roles: "Mathematician", "Lawyer", "Economist", and "Programmer". For the MMLU Professional Medicine subset, we include "Doctor", "Psychologist", "Mathematician", and "Programmer". In practice, we implement this by assigning each agent a system prompt that encodes a specific role or persona---[e.g.]{.ltx_text .ltx_font_italic}, Mathematician, Programmer, Lawyer---using the same prompt templates provided in \[[9](#bib.bib9){.ltx_ref}\] (see Appendix [[B.3]{.ltx_text .ltx_ref_tag}](#A2.SS3 "B.3 Persona Prompts ‣ Appendix B Prompt Templates ‣ Appendix ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref} for the prompts). Even in these heterogeneous settings, Majority Voting mostly paralleled MAD variants. However, several MAD results on Pro. Med. show larger gains, suggesting the potential benefit of assigning diverse personas in task-specific MAD systems.
:::

<figure id="S6.T4" class="ltx_table">
<div class="ltx_flex_figure ltx_flex_table">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S6.T4.fig1" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_top" style="width:158.7pt;">
<div class="ltx_inline-block ltx_align_center ltx_transformed_outer" style="width:433.6pt;height:294.1pt;vertical-align:-143.0pt;">
<span class="ltx_transformed_inner" style="transform:translate(82.7pt,-56.1pt) scale(1.61620585085022,1.61620585085022) ;"> </span>
<table class="ltx_tabular ltx_align_middle">
<tbody>
<tr class="ltx_tr">
<td rowspan="2" class="ltx_td ltx_align_left ltx_border_tt"><span class="ltx_text ltx_font_bold">Methods</span></td>
<td colspan="2" class="ltx_td ltx_align_center ltx_border_tt">Qwen2.5-32B-Instruct</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">GSM8K</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">HellaSwag</span></td>
</tr>
<tr class="ltx_tr">
<td colspan="3" class="ltx_td ltx_align_center ltx_border_t" style="background-color: #F2F2F2"><span class="ltx_text ltx_font_bold" style="background-color:#F2F2F2;">Single-Agent</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Single-agent baseline</td>
<td class="ltx_td ltx_align_center">0.7566 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .01</span></td>
<td class="ltx_td ltx_align_center">0.8620 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .01</span></td>
</tr>
<tr class="ltx_tr">
<td colspan="3" class="ltx_td ltx_align_center ltx_border_t" style="background-color: #F2F2F2"><span class="ltx_text ltx_font_bold" style="background-color:#F2F2F2;">Multi-Agent</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Decentralized MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
<td class="ltx_td ltx_align_center">0.9400</td>
<td class="ltx_td ltx_align_center">0.8633</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Decentralized MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
<td class="ltx_td ltx_align_center">0.9367</td>
<td class="ltx_td ltx_align_center">0.8600</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Decentralized MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
<td class="ltx_td ltx_align_center">0.9367</td>
<td class="ltx_td ltx_align_center">0.8600</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Sparse MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
<td class="ltx_td ltx_align_center">0.8433</td>
<td class="ltx_td ltx_align_center">0.8667</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Sparse MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
<td class="ltx_td ltx_align_center">0.9367</td>
<td class="ltx_td ltx_align_center">0.8633</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Sparse MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
<td class="ltx_td ltx_align_center">0.9333</td>
<td class="ltx_td ltx_align_center">0.8667</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Centralized MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
<td class="ltx_td ltx_align_center">0.8000</td>
<td class="ltx_td ltx_align_center">0.8667</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Centralized MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
<td class="ltx_td ltx_align_center">0.8333</td>
<td class="ltx_td ltx_align_center">0.8667</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Centralized MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
<td class="ltx_td ltx_align_center">0.8333</td>
<td class="ltx_td ltx_align_center">0.8667</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_bb ltx_border_t"><span class="ltx_text ltx_font_bold">Majority Voting</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t"><span class="ltx_text ltx_font_bold">0.9433</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t"><span class="ltx_text ltx_font_bold">0.8667</span></td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Table 3</span>: </span><span class="ltx_text ltx_font_bold" style="font-size:90%;">Results on a larger model.</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S6.T4.fig2" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_top" style="width:165.6pt;">
<div class="ltx_inline-block ltx_align_center ltx_transformed_outer" style="width:433.6pt;height:254pt;vertical-align:-123.5pt;">
<span class="ltx_transformed_inner" style="transform:translate(61.5pt,-36.0pt) scale(1.3956958531787,1.3956958531787) ;"> </span>
<table class="ltx_tabular ltx_align_middle">
<tbody>
<tr class="ltx_tr">
<td rowspan="2" class="ltx_td ltx_align_left ltx_border_tt"><span class="ltx_text ltx_font_bold">Methods</span></td>
<td colspan="2" class="ltx_td ltx_align_center ltx_border_tt">Qwen2.5-7B-Instruct</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">GSM8K</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">MMLU (Pro.Med.)</span></td>
</tr>
<tr class="ltx_tr">
<td colspan="3" class="ltx_td ltx_align_center ltx_border_t" style="background-color: #F2F2F2"><span class="ltx_text ltx_font_bold" style="background-color:#F2F2F2;">Single-Agent</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Single-agent baseline</td>
<td class="ltx_td ltx_align_center">0.8047<span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .05</span></td>
<td class="ltx_td ltx_align_center">0.7890 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .01</span></td>
</tr>
<tr class="ltx_tr">
<td colspan="3" class="ltx_td ltx_align_center ltx_border_t" style="background-color: #F2F2F2"><span class="ltx_text ltx_font_bold" style="background-color:#F2F2F2;">Multi-Agent</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Decentralized MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
<td class="ltx_td ltx_align_center">0.8033</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">0.8419</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Decentralized MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
<td class="ltx_td ltx_align_center">0.7733</td>
<td class="ltx_td ltx_align_center">0.8382</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Decentralized MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
<td class="ltx_td ltx_align_center">0.7433</td>
<td class="ltx_td ltx_align_center">0.8382</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Sparse MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
<td class="ltx_td ltx_align_center">0.8900</td>
<td class="ltx_td ltx_align_center">0.8346</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Sparse MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
<td class="ltx_td ltx_align_center">0.8667</td>
<td class="ltx_td ltx_align_center">0.8346</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Sparse MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
<td class="ltx_td ltx_align_center">0.8567</td>
<td class="ltx_td ltx_align_center">0.8382</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Centralized MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
<td class="ltx_td ltx_align_center">0.5933</td>
<td class="ltx_td ltx_align_center">0.8051</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Centralized MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
<td class="ltx_td ltx_align_center">0.5900</td>
<td class="ltx_td ltx_align_center">0.7978</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Centralized MAD (<span class="math inline"><em>T</em> = 5</span>)</td>
<td class="ltx_td ltx_align_center">0.6000</td>
<td class="ltx_td ltx_align_center">0.7978</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_bb ltx_border_t"><span class="ltx_text ltx_font_bold">Majority Voting</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t"><span class="ltx_text ltx_font_bold">0.9367</span></td>
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_t">0.8235</td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span class="ltx_text" style="font-size:90%;">Table 4</span>: </span><span class="ltx_text ltx_font_bold" style="font-size:90%;">Heterogeneous persona agents.<span class="ltx_text ltx_font_medium"> Single-agent is averaged over 5 personas (prompts in <a href="#A2.SS3" class="ltx_ref" title="B.3 Persona Prompts ‣ Appendix B Prompt Templates ‣ Appendix ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"><span class="ltx_text ltx_ref_tag">B.3</span></a>).</span></span></figcaption>
</figure>
</div>
</div>
</figure>

<figure id="S6.T5" class="ltx_table ltx_align_floatright">
<div class="ltx_inline-block ltx_align_center ltx_transformed_outer" style="width:256.9pt;height:61.2pt;vertical-align:-28.1pt;">
<span class="ltx_transformed_inner" style="transform:translate(0.0pt,0.0pt) scale(1,1) ;"> </span>
<p><span class="ltx_tabular ltx_align_middle"> <span class="ltx_tr"> <span class="ltx_td ltx_align_center ltx_border_tt" style="padding-left:5.0pt;padding-right:5.0pt;"><span class="ltx_text ltx_font_bold">Methods</span></span> <span class="ltx_td ltx_align_center ltx_border_tt" style="padding-left:5.0pt;padding-right:5.0pt;"><span class="ltx_text ltx_font_bold">Rouge-1</span></span> <span class="ltx_td ltx_align_center ltx_border_tt" style="padding-left:5.0pt;padding-right:5.0pt;"><span class="ltx_text ltx_font_bold">Rouge-L</span></span></span> <span class="ltx_tr"> <span class="ltx_td ltx_align_left ltx_border_t" style="padding-left:5.0pt;padding-right:5.0pt;">Best Single-agent</span> <span class="ltx_td ltx_align_center ltx_border_t" style="padding-left:5.0pt;padding-right:5.0pt;">0.2760</span> <span class="ltx_td ltx_align_center ltx_border_t" style="padding-left:5.0pt;padding-right:5.0pt;">0.1871</span></span> <span class="ltx_tr"> <span class="ltx_td ltx_align_left" style="padding-left:5.0pt;padding-right:5.0pt;">MAD (<span class="math inline"><em>T</em> = 1</span>)</span> <span class="ltx_td ltx_align_center" style="padding-left:5.0pt;padding-right:5.0pt;">0.2686</span> <span class="ltx_td ltx_align_center" style="padding-left:5.0pt;padding-right:5.0pt;">0.1814</span></span> <span class="ltx_tr"> <span class="ltx_td ltx_align_left" style="padding-left:5.0pt;padding-right:5.0pt;">MAD (<span class="math inline"><em>T</em> = 2</span>)</span> <span class="ltx_td ltx_align_center" style="padding-left:5.0pt;padding-right:5.0pt;">0.2773</span> <span class="ltx_td ltx_align_center" style="padding-left:5.0pt;padding-right:5.0pt;">0.1867</span></span> <span class="ltx_tr"> <span class="ltx_td ltx_align_left ltx_border_bb" style="padding-left:5.0pt;padding-right:5.0pt;">MAD (<span class="math inline"><em>T</em> = 3</span>)</span> <span class="ltx_td ltx_align_center ltx_border_bb" style="padding-left:5.0pt;padding-right:5.0pt;">0.2825</span> <span class="ltx_td ltx_align_center ltx_border_bb" style="padding-left:5.0pt;padding-right:5.0pt;">0.1852</span></span> </span></p>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 5</span>: </span><span class="ltx_text ltx_font_bold" style="font-size:90%;">Open-ended text generation.<span class="ltx_text ltx_font_medium"> Qwen2.5-7B-Instruct on Decentralized MAD.</span></span></figcaption>
</figure>
::::

:::: {#S6.SS0.SSS0.Px3 .section .ltx_paragraph}
##### Evaluation on open-ended text generation tasks. {#evaluation-on-open-ended-text-generation-tasks. .ltx_title .ltx_title_paragraph}

::: {#S6.SS0.SSS0.Px3.p1 .ltx_para}
In our previous experiments, we mainly focused on closed-ended question answering tasks, which are the primary focus of previous works on MAD. A natural follow-up question is whether our findings will hold in open-ended tasks, such as free-form text generation. To explore this, we evaluate MAD on a text summarization task using a subset of CNN/DailyMail dataset \[[25](#bib.bib25){.ltx_ref}\]. Unlike classification tasks, applying Majority Voting in summarization is not straightforward due to the lack of discrete answer choices. Instead, we report the best-performing agent at each debate round, as shown in Table [[5]{.ltx_text .ltx_ref_tag}](#S6.T5 "Table 5 ‣ Heterogeneous Agents. ‣ 6 Extended Experiments to General Settings ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}. Interestingly, we observe that ROUGE-1 and ROUGE-L scores remain relatively invariant across rounds, suggesting that the key observations from closed-ended tasks may also extend to open-ended tasks like summarization.
:::
::::
::::::::::

:::::::: {#S7 .section .ltx_section}
### [7 ]{.ltx_tag .ltx_tag_section}Related Works {#related-works .ltx_title .ltx_title_section}

::: {#S7.p1 .ltx_para}
Recently, there has been growing interest in multi-agent systems (MAS). Several survey papers have reviewed state-of-the-art LLM-based MAS approaches \[[14](#bib.bib14){.ltx_ref}, [26](#bib.bib26){.ltx_ref}, [27](#bib.bib27){.ltx_ref}, [28](#bib.bib28){.ltx_ref}\]. Within MAS, MAD has emerged as a particularly promising approach for enhancing the performance of single-agent benchmarks. In the following, we discuss the strengths and limitations of current MAD systems.
:::

:::: {#S7.SS0.SSS0.Px1 .section .ltx_paragraph}
##### Pros of Multi-Agent Debate. {#pros-of-multi-agent-debate. .ltx_title .ltx_title_paragraph}

::: {#S7.SS0.SSS0.Px1.p1 .ltx_para}
A key strength of MAD lies in its iterative discussion process, which has the potential to enhance both factual accuracy and reasoning quality. Building on this paradigm, several works have proposed MAD-based approaches for a variety of tasks \[[1](#bib.bib1){.ltx_ref}, [2](#bib.bib2){.ltx_ref}, [3](#bib.bib3){.ltx_ref}, [4](#bib.bib4){.ltx_ref}, [5](#bib.bib5){.ltx_ref}, [6](#bib.bib6){.ltx_ref}\]. To further advance MAD systems, \[[7](#bib.bib7){.ltx_ref}\] introduced enhancements grounded in debate theory, while \[[29](#bib.bib29){.ltx_ref}\] developed the Peer Rank and Peer Discussion mechanisms to select appropriate agent pairs for debate. Many studies have focused on designing effective communication architectures and protocols to improve efficiency and effectiveness \[[3](#bib.bib3){.ltx_ref}, [8](#bib.bib8){.ltx_ref}, [9](#bib.bib9){.ltx_ref}, [10](#bib.bib10){.ltx_ref}, [30](#bib.bib30){.ltx_ref}, [31](#bib.bib31){.ltx_ref}\]. Other works have emphasized the importance of diversity in MAD systems, leveraging heterogeneous LLM agents \[[11](#bib.bib11){.ltx_ref}\], injecting distinct personas into each agent \[[9](#bib.bib9){.ltx_ref}, [12](#bib.bib12){.ltx_ref}, [13](#bib.bib13){.ltx_ref}\], or enabling text generation with controlled diversity \[[32](#bib.bib32){.ltx_ref}, [33](#bib.bib33){.ltx_ref}\]. Additionally, learning-based methods have been explored to optimize MAD dynamics \[[9](#bib.bib9){.ltx_ref}, [34](#bib.bib34){.ltx_ref}, [35](#bib.bib35){.ltx_ref}\].
:::
::::

:::: {#S7.SS0.SSS0.Px2 .section .ltx_paragraph}
##### Cons of Multi-Agent Debate. {#cons-of-multi-agent-debate. .ltx_title .ltx_title_paragraph}

::: {#S7.SS0.SSS0.Px2.p1 .ltx_para}
While MAD systems are widely used as effective tools for solving various tasks, recent studies have raised concerns about their actual effectiveness. For instance, \[[36](#bib.bib36){.ltx_ref}\] conducted an in-depth analysis identifying 14 distinct failure modes in MAD systems, and \[[37](#bib.bib37){.ltx_ref}\] found that MAD does not consistently outperform single-agent approaches. Similarly, \[[38](#bib.bib38){.ltx_ref}\] showed that LLM agents are not self-corrective enough for MAD to be successful, and \[[39](#bib.bib39){.ltx_ref}\] reported that MAD performs no better than advanced single-agent reasoning methods and highlighted its sensitivity to hyperparameters. \[[40](#bib.bib40){.ltx_ref}\] echoed these concerns, showing that well-prompted single agents can sometimes outperform MAD. More specifically, \[[7](#bib.bib7){.ltx_ref}\] and \[[37](#bib.bib37){.ltx_ref}\] observed occurrences of subverted or incorrect answers in MAD debates, while \[[41](#bib.bib41){.ltx_ref}\] showed that MAD systems often converge to the majority opinion, even when that opinion reflects common misconceptions. On a slightly different gear, \[[42](#bib.bib42){.ltx_ref}\] compared various decision protocols and showed that multiple rounds of MAD actually decreases performance. In this work, we *perform a systematic comparison between MAD and simple majority vote, and provide theoretical foundation to understand how the success probability evolves throughout debate*, shedding light on future design of improved MAD.
:::
::::
::::::::

:::::: {#S8 .section .ltx_section}
### [8 ]{.ltx_tag .ltx_tag_section}Conclusion {#conclusion .ltx_title .ltx_title_section}

::: {#S8.p1 .ltx_para}
In this study, we provided comprehensive analysis of MAD and its core components. To investigate this, we conducted extensive experiments on seven benchmarks. Contrary to prevailing assumptions, we observe that most performance gains of MAD stem from majority voting rather than the debate process itself. To support this finding, we introduce a theoretical framework that characterizes debate dynamics as a martingale process, which preserves the expected success probability of each agent over time. These insights suggest that ensembling strategies like majority voting remain strong and often more reliable, highlighting the need to preserve correct answers during inter-agent debate. Overall, our work sheds light on the key mechanisms underlying MAD and offers concrete directions for improving its design.
:::

:::: {#S8.SS0.SSS0.Px1 .section .ltx_paragraph}
##### Broader Impact. {#broader-impact. .ltx_title .ltx_title_paragraph}

::: {#S8.SS0.SSS0.Px1.p1 .ltx_para}
Our findings highlight an important perspective on Multi-Agent Debate, showing that much of its effectiveness can be achieved through simpler, more accessible methods like Majority Voting. This opens the door to building more efficient and scalable collaborative AI systems without sacrificing performance. Moreover, by identifying the MAD as a martingale process, we offer actionable insights that can help make future MAD systems more robust and trustworthy. We believe our work contributes to a new perspective on debate-based AI frameworks that are both principled and practical. Ultimately, this supports the broader goal of making AI systems more reliable, collaborative, and aligned with human reasoning.
:::
::::
::::::

::::: {#Sx1 .section .ltx_section}
### Acknowledgement {#acknowledgement .ltx_title .ltx_title_section}

::: {#Sx1.p1 .ltx_para}
The authors would like to thank Leitian Tao and Xuanming Zhang for their valuable comments on the manuscript. Hyeong Kyu Choi and Yixuan Li are supported in part by the AFOSR Young Investigator Program under award number FA9550-23-1-0184, National Science Foundation under awards IIS-2237037 and IIS-2331669, Office of Naval Research under grant number N00014-23-1-2643, Schmidt Sciences Foundation, Open Philanthropy, Alfred P. Sloan Fellowship, and gifts from Google and Amazon. Xiaojin Zhu was supported in part by NSF grants 2202457, 2331669, 1836978, 2023239, ARO MURI W911NF2110317, and AF CoE FA9550-18-1-0166.
:::

::: {.ltx_pagination .ltx_role_newpage}
:::
:::::

::: {#bib .section .ltx_bibliography}
### References {#references .ltx_title .ltx_title_bibliography}

- [[\[1\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xiaohe Bo, Zeyu Zhang, Quanyu Dai, Xueyang Feng, Lei Wang, Rui Li, Xu Chen, and Ji-Rong Wen. ]{.ltx_bibblock} [Reflective multi-agent collaboration based on large language models. ]{.ltx_bibblock} [[Advances in Neural Information Processing Systems]{.ltx_text .ltx_font_italic}, 37:138595--138631, 2024. ]{.ltx_bibblock}]{#bib.bib1}
- [[\[2\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yilun Du, Shuang Li, Antonio Torralba, Joshua B Tenenbaum, and Igor Mordatch. ]{.ltx_bibblock} [Improving factuality and reasoning in language models through multiagent debate. ]{.ltx_bibblock} [In [International Conference on Machine Learning]{.ltx_text .ltx_font_italic}, pages 11733--11763. PMLR, 2024. ]{.ltx_bibblock}]{#bib.bib2}
- [[\[3\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Chi-Min Chan, Weize Chen, Yusheng Su, Jianxuan Yu, Wei Xue, Shanghang Zhang, Jie Fu, and Zhiyuan Liu. ]{.ltx_bibblock} [Chateval: Towards better llm-based evaluators through multi-agent debate. ]{.ltx_bibblock} [In [The Twelfth International Conference on Learning Representations]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib3}
- [[\[4\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xiangru Tang, Anni Zou, Zhuosheng Zhang, Ziming Li, Yilun Zhao, Xingyao Zhang, Arman Cohan, and Mark Gerstein. ]{.ltx_bibblock} [Medagents: Large language models as collaborators for zero-shot medical reasoning. ]{.ltx_bibblock} [In [Findings of the Association for Computational Linguistics ACL 2024]{.ltx_text .ltx_font_italic}, pages 599--621, 2024. ]{.ltx_bibblock}]{#bib.bib4}
- [[\[5\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Qingyun Wu, Gagan Bansal, Jieyu Zhang, Yiran Wu, Beibin Li, Erkang Zhu, Li Jiang, Xiaoyun Zhang, Shaokun Zhang, Jiale Liu, et al. ]{.ltx_bibblock} [Autogen: Enabling next-gen llm applications via multi-agent conversations. ]{.ltx_bibblock} [In [First Conference on Language Modeling]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib5}
- [[\[6\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Weize Chen, Yusheng Su, Jingwei Zuo, Cheng Yang, Chenfei Yuan, Chi-Min Chan, Heyang Yu, Yaxi Lu, Yi-Hsin Hung, Chen Qian, et al. ]{.ltx_bibblock} [Agentverse: Facilitating multi-agent collaboration and exploring emergent behaviors. ]{.ltx_bibblock} [In [The Twelfth International Conference on Learning Representations]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib6}
- [[\[7\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Kai Xiong, Xiao Ding, Yixin Cao, Ting Liu, and Bing Qin. ]{.ltx_bibblock} [Examining inter-consistency of large language models collaboration: An in-depth analysis via debate. ]{.ltx_bibblock} [In [Findings of the Association for Computational Linguistics: EMNLP 2023]{.ltx_text .ltx_font_italic}, pages 7572--7590, 2023. ]{.ltx_bibblock}]{#bib.bib7}
- [[\[8\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tongxuan Liu, Xingyu Wang, Weizhe Huang, Wenjiang Xu, Yuting Zeng, Lei Jiang, Hailong Yang, and Jing Li. ]{.ltx_bibblock} [Groupdebate: Enhancing the efficiency of multi-agent debate using group discussion. ]{.ltx_bibblock} [[arXiv preprint arXiv:2409.14051]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib8}
- [[\[9\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zijun Liu, Yanzhe Zhang, Peng Li, Yang Liu, and Diyi Yang. ]{.ltx_bibblock} [Dynamic llm-agent network: An llm-agent collaboration framework with agent team optimization. ]{.ltx_bibblock} [In [COLM]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib9}
- [[\[10\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yunxuan Li, Yibing Du, Jiageng Zhang, Le Hou, Peter Grabowski, Yeqing Li, and Eugene Ie. ]{.ltx_bibblock} [Improving multi-agent debate with sparse communication topology. ]{.ltx_bibblock} [In [Findings of the Association for Computational Linguistics: EMNLP 2024]{.ltx_text .ltx_font_italic}, pages 7281--7294, 2024. ]{.ltx_bibblock}]{#bib.bib10}
- [[\[11\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Justin Chen, Swarnadeep Saha, and Mohit Bansal. ]{.ltx_bibblock} [Reconcile: Round-table conference improves reasoning via consensus among diverse llms. ]{.ltx_bibblock} [In [Proceedings of the 62nd Annual Meeting of the Association for Computational Linguistics (Volume 1: Long Papers)]{.ltx_text .ltx_font_italic}, pages 7066--7085, 2024. ]{.ltx_bibblock}]{#bib.bib11}
- [[\[12\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tian Liang, Zhiwei He, Wenxiang Jiao, Xing Wang, Yan Wang, Rui Wang, Yujiu Yang, Shuming Shi, and Zhaopeng Tu. ]{.ltx_bibblock} [Encouraging divergent thinking in large language models through multi-agent debate. ]{.ltx_bibblock} [In [Proceedings of the 2024 Conference on Empirical Methods in Natural Language Processing]{.ltx_text .ltx_font_italic}, pages 17889--17904, 2024. ]{.ltx_bibblock}]{#bib.bib12}
- [[\[13\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zhenhailong Wang, Shaoguang Mao, Wenshan Wu, Tao Ge, Furu Wei, and Heng Ji. ]{.ltx_bibblock} [Unleashing the emergent cognitive synergy in large language models: A task-solving agent through multi-persona self-collaboration. ]{.ltx_bibblock} [In [Proceedings of the 2024 Conference of the North American Chapter of the Association for Computational Linguistics: Human Language Technologies (Volume 1: Long Papers)]{.ltx_text .ltx_font_italic}, pages 257--279, 2024. ]{.ltx_bibblock}]{#bib.bib13}
- [[\[14\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Taicheng Guo, Xiuying Chen, Yaqi Wang, Ruidi Chang, Shichao Pei, Nitesh V Chawla, Olaf Wiest, and Xiangliang Zhang. ]{.ltx_bibblock} [Large language model based multi-agents: a survey of progress and challenges. ]{.ltx_bibblock} [In [Proceedings of the Thirty-Third International Joint Conference on Artificial Intelligence]{.ltx_text .ltx_font_italic}, pages 8048--8057, 2024. ]{.ltx_bibblock}]{#bib.bib14}
- [[\[15\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Karl Cobbe, Vineet Kosaraju, Mohammad Bavarian, Mark Chen, Heewoo Jun, Lukasz Kaiser, Matthias Plappert, Jerry Tworek, Jacob Hilton, Reiichiro Nakano, Christopher Hesse, and John Schulman. ]{.ltx_bibblock} [Training verifiers to solve math word problems. ]{.ltx_bibblock} [[arXiv preprint arXiv:2110.14168]{.ltx_text .ltx_font_italic}, 2021. ]{.ltx_bibblock}]{#bib.bib15}
- [[\[16\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Dan Hendrycks, Collin Burns, Steven Basart, Andy Zou, Mantas Mazeika, Dawn Song, and Jacob Steinhardt. ]{.ltx_bibblock} [Measuring massive multitask language understanding. ]{.ltx_bibblock} [[Proceedings of the International Conference on Learning Representations (ICLR)]{.ltx_text .ltx_font_italic}, 2021. ]{.ltx_bibblock}]{#bib.bib16}
- [[\[17\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Dan Hendrycks, Collin Burns, Steven Basart, Andrew Critch, Jerry Li, Dawn Song, and Jacob Steinhardt. ]{.ltx_bibblock} [Aligning ai with shared human values. ]{.ltx_bibblock} [[Proceedings of the International Conference on Learning Representations (ICLR)]{.ltx_text .ltx_font_italic}, 2021. ]{.ltx_bibblock}]{#bib.bib17}
- [[\[18\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Rowan Zellers, Ari Holtzman, Yonatan Bisk, Ali Farhadi, and Yejin Choi. ]{.ltx_bibblock} [Hellaswag: Can a machine really finish your sentence? ]{.ltx_bibblock} [In [Proceedings of the 57th Annual Meeting of the Association for Computational Linguistics]{.ltx_text .ltx_font_italic}, 2019. ]{.ltx_bibblock}]{#bib.bib18}
- [[\[19\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Alon Talmor, Jonathan Herzig, Nicholas Lourie, and Jonathan Berant. ]{.ltx_bibblock} [Commonsenseqa: A question answering challenge targeting commonsense knowledge. ]{.ltx_bibblock} [In [Proceedings of the 2019 Conference of the North American Chapter of the Association for Computational Linguistics: Human Language Technologies, Volume 1 (Long and Short Papers)]{.ltx_text .ltx_font_italic}, pages 4149--4158, 2019. ]{.ltx_bibblock}]{#bib.bib19}
- [[\[20\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yuntao Bai, Andy Jones, Kamal Ndousse, Amanda Askell, Anna Chen, Nova DasSarma, Dawn Drain, Stanislav Fort, Deep Ganguli, Tom Henighan, et al. ]{.ltx_bibblock} [Training a helpful and harmless assistant with reinforcement learning from human feedback. ]{.ltx_bibblock} [[arXiv preprint arXiv:2204.05862]{.ltx_text .ltx_font_italic}, 2022. ]{.ltx_bibblock}]{#bib.bib20}
- [[\[21\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Harrison Lee, Samrat Phatale, Hassan Mansoor, Thomas Mesnard, Johan Ferret, Kellie Ren Lu, Colton Bishop, Ethan Hall, Victor Carbune, Abhinav Rastogi, et al. ]{.ltx_bibblock} [Rlaif vs. rlhf: Scaling reinforcement learning from human feedback with ai feedback. ]{.ltx_bibblock} [In [International Conference on Machine Learning]{.ltx_text .ltx_font_italic}, pages 26874--26901. PMLR, 2024. ]{.ltx_bibblock}]{#bib.bib21}
- [[\[22\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ An Yang, Baosong Yang, Beichen Zhang, Binyuan Hui, Bo Zheng, Bowen Yu, Chengyuan Li, Dayiheng Liu, Fei Huang, Haoran Wei, et al. ]{.ltx_bibblock} [Qwen2. 5 technical report. ]{.ltx_bibblock} [[arXiv preprint arXiv:2412.15115]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib22}
- [[\[23\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Aaron Grattafiori, Abhimanyu Dubey, Abhinav Jauhri, Abhinav Pandey, Abhishek Kadian, Ahmad Al-Dahle, Aiesha Letman, Akhil Mathur, Alan Schelten, Alex Vaughan, et al. ]{.ltx_bibblock} [The llama 3 herd of models. ]{.ltx_bibblock} [[arXiv preprint arXiv:2407.21783]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib23}
- [[\[24\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Robin Pemantle. ]{.ltx_bibblock} [A survey of random processes with reinforcement. ]{.ltx_bibblock} [[Probability Surveys]{.ltx_text .ltx_font_italic}, 4:1--79, 2007. ]{.ltx_bibblock}]{#bib.bib24}
- [[\[25\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Abigail See, Peter J. Liu, and Christopher D. Manning. ]{.ltx_bibblock} [Get to the point: Summarization with pointer-generator networks. ]{.ltx_bibblock} [In [Proceedings of the 55th Annual Meeting of the Association for Computational Linguistics (Volume 1: Long Papers)]{.ltx_text .ltx_font_italic}, pages 1073--1083. Association for Computational Linguistics, 2017. ]{.ltx_bibblock}]{#bib.bib25}
- [[\[26\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Khanh-Tung Tran, Dung Dao, Minh-Duong Nguyen, Quoc-Viet Pham, Barry O'Sullivan, and Hoang D Nguyen. ]{.ltx_bibblock} [Multi-agent collaboration mechanisms: A survey of llms. ]{.ltx_bibblock} [[arXiv preprint arXiv:2501.06322]{.ltx_text .ltx_font_italic}, 2025. ]{.ltx_bibblock}]{#bib.bib26}
- [[\[27\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Bingyu Yan, Xiaoming Zhang, Litian Zhang, Lian Zhang, Ziyi Zhou, Dezhuang Miao, and Chaozhuo Li. ]{.ltx_bibblock} [Beyond self-talk: A communication-centric survey of llm-based multi-agent systems. ]{.ltx_bibblock} [[arXiv preprint arXiv:2502.14321]{.ltx_text .ltx_font_italic}, 2025. ]{.ltx_bibblock}]{#bib.bib27}
- [[\[28\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xinyi Li, Sai Wang, Siqi Zeng, Yu Wu, and Yi Yang. ]{.ltx_bibblock} [A survey on llm-based multi-agent systems: workflow, infrastructure, and challenges. ]{.ltx_bibblock} [[Vicinagearth]{.ltx_text .ltx_font_italic}, 1(1):9, 2024. ]{.ltx_bibblock}]{#bib.bib28}
- [[\[29\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ruosen Li, Teerth Patel, and Xinya Du. ]{.ltx_bibblock} [Prd: Peer rank and discussion improve large language model based evaluations. ]{.ltx_bibblock} [[Transactions on Machine Learning Research]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib29}
- [[\[30\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Chau Pham, Boyi Liu, Yingxiang Yang, Zhengyu Chen, Tianyi Liu, Jianbo Yuan, Bryan A Plummer, Zhaoran Wang, and Hongxia Yang. ]{.ltx_bibblock} [Let models speak ciphers: Multiagent debate through embeddings. ]{.ltx_bibblock} [In [The Twelfth International Conference on Learning Representations]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib30}
- [[\[31\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Guibin Zhang, Yanwei Yue, Zhixun Li, Sukwon Yun, Guancheng Wan, Kun Wang, Dawei Cheng, Jeffrey Xu Yu, and Tianlong Chen. ]{.ltx_bibblock} [Cut the crap: An economical communication pipeline for llm-based multi-agent systems. ]{.ltx_bibblock} [[arXiv preprint arXiv:2410.02506]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib31}
- [[\[32\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yexiang Liu, Jie Cao, Zekun Li, Ran He, and Tieniu Tan. ]{.ltx_bibblock} [Breaking mental set to improve reasoning through diverse multi-agent debate. ]{.ltx_bibblock} [In [The Thirteenth International Conference on Learning Representations]{.ltx_text .ltx_font_italic}, 2025. ]{.ltx_bibblock}]{#bib.bib32}
- [[\[33\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ KuanChao Chu, Yi-Pei Chen, and Hideki Nakayama. ]{.ltx_bibblock} [Exploring and controlling diversity in llm-agent conversation. ]{.ltx_bibblock} [[arXiv preprint arXiv:2412.21102]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib33}
- [[\[34\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Andrew Estornell, Jean-Francois Ton, Yuanshun Yao, and Yang Liu. ]{.ltx_bibblock} [Acc-debate: An actor-critic approach to multi-agent debate. ]{.ltx_bibblock} [In [The Thirteenth International Conference on Learning Representations]{.ltx_text .ltx_font_italic}, 2025. ]{.ltx_bibblock}]{#bib.bib34}
- [[\[35\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Weize Chen, Jiarui Yuan, Chen Qian, Cheng Yang, Zhiyuan Liu, and Maosong Sun. ]{.ltx_bibblock} [Optima: Optimizing effectiveness and efficiency for llm-based multi-agent system. ]{.ltx_bibblock} [[arXiv preprint arXiv:2410.08115]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib35}
- [[\[36\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Mert Cemri, Melissa Z Pan, Shuyi Yang, Lakshya A Agrawal, Bhavya Chopra, Rishabh Tiwari, Kurt Keutzer, Aditya Parameswaran, Dan Klein, Kannan Ramchandran, et al. ]{.ltx_bibblock} [Why do multi-agent llm systems fail? ]{.ltx_bibblock} [[arXiv preprint arXiv:2503.13657]{.ltx_text .ltx_font_italic}, 2025. ]{.ltx_bibblock}]{#bib.bib36}
- [[\[37\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Hangfan Zhang, Zhiyao Cui, Xinrun Wang, Qiaosheng Zhang, Zhen Wang, Dinghao Wu, and Shuyue Hu. ]{.ltx_bibblock} [If multi-agent debate is the answer, what is the question? ]{.ltx_bibblock} [[arXiv preprint arXiv:2502.08788]{.ltx_text .ltx_font_italic}, 2025. ]{.ltx_bibblock}]{#bib.bib37}
- [[\[38\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jie Huang, Xinyun Chen, Swaroop Mishra, Huaixiu Steven Zheng, Adams Wei Yu, Xinying Song, and Denny Zhou. ]{.ltx_bibblock} [Large language models cannot self-correct reasoning yet. ]{.ltx_bibblock} [In [The Twelfth International Conference on Learning Representations]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib38}
- [[\[39\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Andries Petrus Smit, Nathan Grinsztajn, Paul Duckworth, Thomas D Barrett, and Arnu Pretorius. ]{.ltx_bibblock} [Should we be going mad? a look at multi-agent debate strategies for llms. ]{.ltx_bibblock} [In [International Conference on Machine Learning]{.ltx_text .ltx_font_italic}, pages 45883--45905. PMLR, 2024. ]{.ltx_bibblock}]{#bib.bib39}
- [[\[40\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Qineng Wang, Zihao Wang, Ying Su, Hanghang Tong, and Yangqiu Song. ]{.ltx_bibblock} [Rethinking the bounds of llm reasoning: Are multi-agent discussions the key? ]{.ltx_bibblock} [In [62nd Annual Meeting of the Association for Computational Linguistics, ACL 2024]{.ltx_text .ltx_font_italic}, pages 6106--6131. Association for Computational Linguistics (ACL), 2024. ]{.ltx_bibblock}]{#bib.bib40}
- [[\[41\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Andrew Estornell and Yang Liu. ]{.ltx_bibblock} [Multi-llm debate: Framework, principals, and interventions. ]{.ltx_bibblock} [[Advances in Neural Information Processing Systems]{.ltx_text .ltx_font_italic}, 37:28938--28964, 2024. ]{.ltx_bibblock}]{#bib.bib41}
- [[\[42\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Lars Benedikt Kaesberg, Jonas Becker, Jan Philip Wahle, Terry Ruas, and Bela Gipp. ]{.ltx_bibblock} [Voting or consensus? decision-making in multi-agent debate. ]{.ltx_bibblock} [[arXiv e-prints]{.ltx_text .ltx_font_italic}, pages arXiv--2502, 2025. ]{.ltx_bibblock}]{#bib.bib42}
- [[\[43\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Hyeong Kyu Choi, Weijie Xu, Chi Xue, Stephanie Eckman, and Chandan K Reddy. ]{.ltx_bibblock} [Mitigating selection bias with node pruning and auxiliary options. ]{.ltx_bibblock} [[arXiv preprint arXiv:2409.18857]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib43}
- [[\[44\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Sheng-Lun Wei, Cheng-Kuang Wu, Hen-Hsen Huang, and Hsin-Hsi Chen. ]{.ltx_bibblock} [Unveiling selection biases: Exploring order and token sensitivity in large language models. ]{.ltx_bibblock} [In [Findings of the Association for Computational Linguistics ACL 2024]{.ltx_text .ltx_font_italic}, pages 5598--5621, 2024. ]{.ltx_bibblock}]{#bib.bib44}
- [[\[45\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Chujie Zheng, Hao Zhou, Fandong Meng, Jie Zhou, and Minlie Huang. ]{.ltx_bibblock} [Large language models are not robust multiple choice selectors. ]{.ltx_bibblock} [In [The Twelfth International Conference on Learning Representations]{.ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib45}
- [[\[46\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Daniel Hsu, Sham M Kakade, and Tong Zhang. ]{.ltx_bibblock} [A spectral algorithm for learning hidden markov models. ]{.ltx_bibblock} [[Journal of Computer and System Sciences]{.ltx_text .ltx_font_italic}, 78(5):1460--1480, 2012. ]{.ltx_bibblock}]{#bib.bib46}
:::

::: {.ltx_pagination .ltx_role_newpage}
:::

::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: {#Pt1 .section .ltx_part}
## Appendix {#appendix .ltx_title .ltx_title_part}

:::::::::::::::: {#A1 .section .ltx_appendix}
### [Appendix A ]{.ltx_tag .ltx_tag_appendix}Experimental Details {#appendix-a-experimental-details .ltx_title .ltx_title_appendix}

::::: {#A1.SS1 .section .ltx_subsection}
#### [A.1 ]{.ltx_tag .ltx_tag_subsection}Hyperparameters and Resources {#a.1-hyperparameters-and-resources .ltx_title .ltx_title_subsection}

::: {#A1.SS1.p1 .ltx_para}
[Hyperparameters.]{.ltx_text .ltx_font_bold} To enable stochastic sampling from homogeneous agents, we set the sampling temperature to 1.0, and use nucleus sampling probability of 0.9, which means that sampling is done a dynamic set of likely tokens that together account for 90% of the total probability. Furthermore, we generate a maximum of 512 tokens for all models experimented in the paper.
:::

::: {#A1.SS1.p2 .ltx_para}
[Resources.]{.ltx_text .ltx_font_bold} All experiments were conducted using either RTX A6000 or RTX A100 GPUs.
:::
:::::

:::::::::::: {#A1.SS2 .section .ltx_subsection}
#### [A.2 ]{.ltx_tag .ltx_tag_subsection}Dataset Details {#a.2-dataset-details .ltx_title .ltx_title_subsection}

::: {#A1.SS2.p1 .ltx_para}
Here, we provide dataset details and the number of samples utilized in our experiments.
:::

::: {#A1.SS2.p2 .ltx_para}
[Arithmetics]{.ltx_text .ltx_font_bold} comprises 100 arithmetic questions, in the form of "[What is the result of $a+b*c+d-e/f$?]{.ltx_text .ltx_font_italic}\" The values $a$ to $f$ are randomly sampled from integers between 0 to 30.
:::

::: {#A1.SS2.p3 .ltx_para}
[GSM8K]{.ltx_text .ltx_font_bold} \[[15](#bib.bib15){.ltx_ref}\] comprises high-quality grade school math word problems intended to test mathematical multi-step reasoning. We randomly subsample 300 questions from the original test split.
:::

::: {#A1.SS2.p4 .ltx_para}
[MMLU (Professional Medicine)]{.ltx_text .ltx_font_bold} \[[16](#bib.bib16){.ltx_ref}, [17](#bib.bib17){.ltx_ref}\] is a benchmark specialized in evaluating reasoning abilities in medical domains at a professional level. Specifically, the dataset requires medical concepts, clinical reasoning, and biomedical knowledge to answer questions. We use the entire test split comprised of 272 questions.
:::

::: {#A1.SS2.p5 .ltx_para}
[MMLU (Formal Logic)]{.ltx_text .ltx_font_bold} \[[16](#bib.bib16){.ltx_ref}, [17](#bib.bib17){.ltx_ref}\] is designed to evaluate a model's proficiency in formal reasoning, symbolic manipulation, and logical analysis. We use the entire test split comprised of 126 question for evaluation.
:::

::: {#A1.SS2.p6 .ltx_para}
[HellaSwag]{.ltx_text .ltx_font_bold} \[[18](#bib.bib18){.ltx_ref}\] is a natural language inference (NLI) benchmark dataset, in the context of sentence completion tasks. That is, the benchmark tests whether a model can choose the most plausible continuation of a given context from multiple options, which is a task that demands not just linguistic proficiency but also real-world knowledge and reasoning. We randomly subsample 300 questions from the original test split.
:::

::: {#A1.SS2.p7 .ltx_para}
[CommonsenseQA]{.ltx_text .ltx_font_bold} \[[19](#bib.bib19){.ltx_ref}\] is a multiple-choice question answering dataset designed to evaluate a model's ability to apply commonsense knowledge in natural language understanding. We randomly subsample 300 questions from the original validation split.
:::

::: {#A1.SS2.p8 .ltx_para}
[HH-RLHF]{.ltx_text .ltx_font_bold} \[[20](#bib.bib20){.ltx_ref}\] is a collection of human-annotated data designed to train and evaluate language models for alignment with human preferences, focusing on helpfulness and harmlessness. The dataset is annotated with relative preferences, comprising 'chosen' and 'rejected' sample pairs. Similar to the "AI labeler alignment" practice \[[21](#bib.bib21){.ltx_ref}\], we ask the LLM agent to select the example that is more helpful and less harmful. To avoid selection bias \[[43](#bib.bib43){.ltx_ref}, [44](#bib.bib44){.ltx_ref}, [45](#bib.bib45){.ltx_ref}\], we randomly shuffle the order of "chosen\" and "rejected\" in the input prompt. We use a random subset of 300 pairs from the original test split.
:::

::: {#A1.SS2.p9 .ltx_para}
[CNN/DailyMails]{.ltx_text .ltx_font_bold} \[[25](#bib.bib25){.ltx_ref}\] is a dataset for abstractive text summarization. It was originally constructed from news articles published by CNN and the Daily Mail, aimed for evaluating models that generate concise summaries of long-form text. We use a random subset of 30 samples from the test split of dataset version 3.0.0.
:::
::::::::::::
::::::::::::::::

::::::::::::::::::: {#A2 .section .ltx_appendix}
### [Appendix B ]{.ltx_tag .ltx_tag_appendix}Prompt Templates {#appendix-b-prompt-templates .ltx_title .ltx_title_appendix}

::: {#A2.p1 .ltx_para}
Here, we provide all the prompt templates used for our experiments.
:::

:::::: {#A2.SS1 .section .ltx_subsection}
#### [B.1 ]{.ltx_tag .ltx_tag_subsection}MAD Templates {#b.1-mad-templates .ltx_title .ltx_title_subsection}

::: {#A2.SS1.p1 .ltx_para}
The following is the prompt template for multi-agent debate. For brevity, assume 3 agents are in the debate.
:::

::: {#A2.SS1.p2 .ltx_para}
[ [These are the recent opinions from other agents:]{.ltx_p} [One of the agents' response:]{.ltx_p} [[\<agent 2's response from the previous round\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [One of the agents' response:]{.ltx_p} [[\<agent 3's response from the previous round\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [This was your most recent opinion:]{.ltx_p} [[\<agent 1's response from the previous round\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [Use these opinions carefully as additional advice to revise your recent opinion to give your final answer to the question:]{.ltx_p} [[\<question\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [Make sure to state your final answer in curly brackets at the very end of your response, just like: \"{final answer: 12.34}\".]{.ltx_p} ]{.ltx_inline-block .ltx_framed .ltx_framed_rectangle style="border-color: #000000;"}
:::

::: {#A2.SS1.p3 .ltx_para}
For GSM8K, we give a slightly different demonstrative example: "{final answer: 123}\", and for the other MCQ datasets, we give: "{final answer: (A)}\". For the CNN/Daily Mail dataset, we instruct the model to 'Make sure to provide your summary after stating \"# Summary #\".'
:::
::::::

:::::::::: {#A2.SS2 .section .ltx_subsection}
#### [B.2 ]{.ltx_tag .ltx_tag_subsection}Task Templates {#b.2-task-templates .ltx_title .ltx_title_subsection}

::: {#A2.SS2.p1 .ltx_para}
For completeness, we provide the exact input format used for each dataset. These templates correspond to the [\<question\>]{.ltx_text .ltx_font_typewriter} field in the MAD prompt structure. Datasets not listed below follow the original question format as provided in their respective sources without modification.
:::

::: {#A2.SS2.p2 .ltx_para}
[ [[Arithmetics]{.ltx_text .ltx_framed .ltx_framed_underline}]{.ltx_p} [What is the result of $a$+$b$\*$c$+$d$-$e$÷$f$?]{.ltx_p} ]{.ltx_inline-block .ltx_framed .ltx_framed_rectangle style="border-color: #000000;"}
:::

::: {#A2.SS2.p3 .ltx_para}
[ [[MMLU Professional Medicine & Formal Logic]{.ltx_text .ltx_framed .ltx_framed_underline}]{.ltx_p} [[\<question\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [(A) [\<option 1\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [(B) [\<option 2\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [(C) [\<option 3\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [(D) [\<option 4\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} ]{.ltx_inline-block .ltx_framed .ltx_framed_rectangle style="border-color: #000000;"}
:::

::: {#A2.SS2.p4 .ltx_para}
[ [[HellaSwag]{.ltx_text .ltx_framed .ltx_framed_underline}]{.ltx_p} [Can you choose the option that best follows:]{.ltx_p} [\"[\<context\>]{.ltx_text .ltx_font_typewriter}\"?]{.ltx_p} [(A) [\<option 1\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [(B) [\<option 2\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [(C) [\<option 3\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [(D) [\<option 4\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} ]{.ltx_inline-block .ltx_framed .ltx_framed_rectangle style="border-color: #000000;"}
:::

::: {#A2.SS2.p5 .ltx_para}
[ [[CommonsenseQA]{.ltx_text .ltx_framed .ltx_framed_underline}]{.ltx_p} [[\<question\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [(A) [\<option 1\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [(B) [\<option 2\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [(C) [\<option 3\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [(D) [\<option 4\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [(E) [\<option 5\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} ]{.ltx_inline-block .ltx_framed .ltx_framed_rectangle style="border-color: #000000;"}
:::

::: {#A2.SS2.p6 .ltx_para}
[ [[HH-RLHF]{.ltx_text .ltx_framed .ltx_framed_underline}]{.ltx_p} [Which of the following is relatively more helpful and/or harmless data?]{.ltx_p} [(A) [\<chosen\>]{.ltx_text .ltx_font_typewriter} [or]{.ltx_text .ltx_font_italic} [\<rejected\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} [(B) [\<rejected\>]{.ltx_text .ltx_font_typewriter} [or]{.ltx_text .ltx_font_italic} [\<chosen\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} ]{.ltx_inline-block .ltx_framed .ltx_framed_rectangle style="border-color: #000000;"}
:::

::: {#A2.SS2.p7 .ltx_para}
[ [[CNN/DailyMails]{.ltx_text .ltx_framed .ltx_framed_underline}]{.ltx_p} [Summarize the following in three sentences:]{.ltx_p} [[\<article\>]{.ltx_text .ltx_font_typewriter}]{.ltx_p} ]{.ltx_inline-block .ltx_framed .ltx_framed_rectangle style="border-color: #000000;"}
:::
::::::::::

::::: {#A2.SS3 .section .ltx_subsection}
#### [B.3 ]{.ltx_tag .ltx_tag_subsection}Persona Prompts {#b.3-persona-prompts .ltx_title .ltx_title_subsection}

::: {#A2.SS3.p1 .ltx_para}
To assign heterogeneous personas to agents, we use a system prompt that specifies the role each agent should assume. We adopt the persona descriptions from \[[9](#bib.bib9){.ltx_ref}\], which are provided below:
:::

::: {#A2.SS3.p2 .ltx_para}
- [[•]{.ltx_tag .ltx_tag_item}]{#A2.I1.i1}

  ::: {#A2.I1.i1.p1 .ltx_para}
  [Assistant]{.ltx_text .ltx_framed .ltx_framed_underline}: You are a super-intelligent AI assistant capable of performing tasks more effectively than humans.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#A2.I1.i2}

  ::: {#A2.I1.i2.p1 .ltx_para}
  [Mathematician]{.ltx_text .ltx_framed .ltx_framed_underline}: You are a mathematician. You are good at math games, arithmetic calculation, and long-term planning.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#A2.I1.i3}

  ::: {#A2.I1.i3.p1 .ltx_para}
  [Economist]{.ltx_text .ltx_framed .ltx_framed_underline}: You are an economist. You are good at economics, finance, and business. You have experience on understanding charts while interpreting the macroeconomic environment prevailing across world economies.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#A2.I1.i4}

  ::: {#A2.I1.i4.p1 .ltx_para}
  [Programmer]{.ltx_text .ltx_framed .ltx_framed_underline}: You are a programmer. You are good at computer science, engineering, and physics. You have experience in designing and developing computer software and hardware.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#A2.I1.i5}

  ::: {#A2.I1.i5.p1 .ltx_para}
  [Lawyer]{.ltx_text .ltx_framed .ltx_framed_underline}: You are a lawyer. You are good at law, politics, and history.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#A2.I1.i6}

  ::: {#A2.I1.i6.p1 .ltx_para}
  [Psychologist]{.ltx_text .ltx_framed .ltx_framed_underline}: You are a psychologist. You are good at psychology, sociology, and philosophy. You give people scientific suggestions that will make them feel better.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#A2.I1.i7}

  ::: {#A2.I1.i7.p1 .ltx_para}
  [Doctor]{.ltx_text .ltx_framed .ltx_framed_underline}: You are a doctor and come up with creative treatments for illnesses or diseases. You are able to recommend conventional medicines, herbal remedies and other natural alternatives. You also consider the patient's age, lifestyle and medical history when providing your recommendations.
  :::
:::
:::::
:::::::::::::::::::

:::::::::::::::::::::::::::::::: {#A3 .section .ltx_appendix}
### [Appendix C ]{.ltx_tag .ltx_tag_appendix}Proofs {#appendix-c-proofs .ltx_title .ltx_title_appendix}

::::: {#A3.SS1 .section .ltx_subsection}
#### [C.1 ]{.ltx_tag .ltx_tag_subsection}Proof of Theorem 1 {#c.1-proof-of-theorem-1 .ltx_title .ltx_title_subsection}

::: {#A3.SS1.p1 .ltx_para}
[Theorem 1. (Majority Voting Success Probability)]{.ltx_text .ltx_font_bold} [ Let $\bar{\boldsymbol{\theta}}=(\bar{\theta}^{(1)},\ldots,\bar{\theta}^{(K)})=\boldsymbol{\alpha}/\sum_{j=1}^{K}\alpha_{j}$ denote the mean of the Dirichlet distribution, $\mathrm{Dirichlet}(\alpha)$, and define the margin $\Delta:=\bar{\theta_{1}}-\bar{\theta_{2}}$. If $N>k/\Delta^{2}$, then the probability that majority voting selects the answer 1 is lower bounded as:]{.ltx_text .ltx_font_italic}

  -- ------------------------------------------------------------------------------------------------------------------------- --
     $$\mathbb{P}(y_{\mathrm{mv}}=1)\geq 1-\exp\left(-N\left(\frac{\Delta}{\sqrt{k}}-\frac{1}{\sqrt{N}}\right)^{2}\right).$$   
  -- ------------------------------------------------------------------------------------------------------------------------- --
:::

::: {#A3.SS1.p2 .ltx_para}
[Proof of Theorem 1.]{.ltx_text .ltx_font_italic} We proceed in several steps to establish the result.
:::
:::::

:::: {#A3.SSx1 .section .ltx_subsection}
#### Step 1: Define the Empirical Distribution {#step-1-define-the-empirical-distribution .ltx_title .ltx_title_subsection}

::: {#A3.SSx1.p1 .ltx_para}
The empirical distribution of votes is given by $\hat{\mathbf{p}}=\mathbf{c}/N=(\hat{p}_{1},\hat{p}_{2},\ldots,\hat{p}_{K})$, where $\hat{p}_{j}=c_{j}/N$ is the fraction of agents selecting answer $j$. The true distribution is $\mathbf{p}=\bar{\boldsymbol{\theta}}$, so $\mathbb{E}[\hat{p}_{j}]=\bar{\theta}_{j}$. The majority-voted answer corresponds to the index with the largest $\hat{p}_{j}$:

  -- --------------------------------------------------------------------------- --
     $$y_{\mathrm{mv}}=\underset{j\in\{1,\ldots,K\}}{\arg\max}\;\hat{p}_{j}.$$   
  -- --------------------------------------------------------------------------- --

Our goal is to show that $\hat{p}_{1}>\hat{p}_{j}$ for all $j\neq 1$ with high probability under the given conditions.
:::
::::

::::::: {#A3.SSx2 .section .ltx_subsection}
#### Step 2: Establish a Key Lemma {#step-2-establish-a-key-lemma .ltx_title .ltx_title_subsection}

::: {#A3.SSx2.p1 .ltx_para}
To connect the empirical distribution to majority voting, we introduce the following lemma.
:::

::: {#A3.SSx2.p2 .ltx_para}
[Lemma 2.]{.ltx_text .ltx_font_bold} [If $\|\hat{\mathbf{p}}-\mathbf{p}\|_{1}<\Delta$, then $y_{\mathrm{mv}}=1$.]{.ltx_text .ltx_font_italic}
:::

::: {#A3.SSx2.p3 .ltx_para}
[Proof of Lemma 2.]{.ltx_text .ltx_font_italic} We prove this by contrapositive. Suppose $y_{\mathrm{mv}}\neq 1$. Then, there exists some $j\neq 1$ such that $\hat{p}_{j}\geq\hat{p}_{1}$. Compute the $L_{1}$ norm contribution involving classes 1 and $j$:

  -- ----------------------------------------------------------------------------------------------------------------------------------- --
     $$|p_{1}-\hat{p}_{1}|+|p_{j}-\hat{p}_{j}|\geq|(p_{1}-\hat{p}_{1})-(p_{j}-\hat{p}_{j})|=|p_{1}-p_{j}-(\hat{p}_{1}-\hat{p}_{j})|.$$   
  -- ----------------------------------------------------------------------------------------------------------------------------------- --

Since $p_{1}=\bar{\theta}_{1}$, $p_{j}=\bar{\theta}_{j}$, and $\Delta=\bar{\theta}_{1}-\bar{\theta}_{2}$, and given $\bar{\theta}_{1}\geq\bar{\theta}_{j}$ for all $j\geq 2$, we have $p_{1}-p_{j}\geq\bar{\theta}_{1}-\bar{\theta}_{2}=\Delta$ (since $\bar{\theta}_{2}\geq\bar{\theta}_{j}$ for $j\geq 2$). Given $\hat{p}_{j}\geq\hat{p}_{1}$, compute:

  -- ------------------------------------------------------------------------------ --
     $$\hat{p}_{1}-\hat{p}_{j}\leq 0\quad\text{and}\quad p_{1}-p_{j}\geq\Delta,$$   
  -- ------------------------------------------------------------------------------ --

so:

  -- --------------------------------------------------------------- --
     $$p_{1}-p_{j}-(\hat{p}_{1}-\hat{p}_{j})\geq\Delta-0=\Delta.$$   
  -- --------------------------------------------------------------- --

Thus:

  -- ---------------------------------------------------------- --
     $$|(p_{1}-\hat{p}_{1})-(p_{j}-\hat{p}_{j})|\geq\Delta,$$   
  -- ---------------------------------------------------------- --

implying:

  -- -------------------------------------------------------- --
     $$|p_{1}-\hat{p}_{1}|+|p_{j}-\hat{p}_{j}|\geq\Delta.$$   
  -- -------------------------------------------------------- --

Since $\|\hat{\mathbf{p}}-\mathbf{p}\|_{1}=\sum_{i=1}^{K}|p_{i}-\hat{p}_{i}|\geq|p_{1}-\hat{p}_{1}|+|p_{j}-\hat{p}_{j}|$, it follows that:

  -- ---------------------------------------------------- --
     $$\|\hat{\mathbf{p}}-\mathbf{p}\|_{1}\geq\Delta.$$   
  -- ---------------------------------------------------- --

Hence, if $y_{\mathrm{mv}}\neq 1$, then $\|\hat{\mathbf{p}}-\mathbf{p}\|_{1}\geq\Delta$. Conversely, if $\|\hat{\mathbf{p}}-\mathbf{p}\|_{1}<\Delta$, then $y_{\mathrm{mv}}=1$. Proof of Lemma 2 completed. $\hfill\square$
:::

::: {#A3.SSx2.p4 .ltx_para}
This lemma implies:

  -- ---------------------------------------------------------------------------------------------- --
     $$\mathbb{P}(y_{\mathrm{mv}}=1)\geq\mathbb{P}(\|\hat{\mathbf{p}}-\mathbf{p}\|_{1}<\Delta).$$   
  -- ---------------------------------------------------------------------------------------------- --
:::
:::::::

::::::: {#A3.SSx3 .section .ltx_subsection}
#### Step 3: Bound the $L_{1}$ Deviation {#step-3-bound-the-l_1-deviation .ltx_title .ltx_title_subsection}

::: {#A3.SSx3.p1 .ltx_para}
Since $\bar{\boldsymbol{\theta}}=\boldsymbol{\alpha}/\sum_{j=1}^{K}\alpha_{j}$ is the mean of the categorical distribution induced by the Dirichlet model, and agents draw answers independently from $\bar{\boldsymbol{\theta}}$, the counts $\mathbf{c}\sim\mathrm{Multinomial}(N,\bar{\boldsymbol{\theta}})$. Note that the Dirichlet-multinomial assumption in the original proof simplifies to a multinomial distribution here because $\bar{\boldsymbol{\theta}}$ is fixed as the mean.
:::

::: {#A3.SSx3.p2 .ltx_para}
We need to bound $\mathbb{P}(\|\hat{\mathbf{p}}-\mathbf{p}\|_{1}\geq\Delta)$. We use a concentration inequality for multinomial distributions. A known result (e.g., Proposition 19 in \[[46](#bib.bib46){.ltx_ref}\]) provides:

  -- --------------------------------------------------------------------------------------------------------------------------------------------- --
     $$\mathbb{P}\left(\|\hat{\mathbf{p}}-\mathbf{p}\|_{1}\geq\sqrt{K}\left(\frac{1}{\sqrt{N}}+\epsilon\right)\right)\leq\exp(-N\epsilon^{2}),$$   
  -- --------------------------------------------------------------------------------------------------------------------------------------------- --

for some $\epsilon>0$. This bound accounts for the $K$-dimensional nature of the distribution.
:::

::: {#A3.SSx3.p3 .ltx_para}
Set the threshold equal to $\Delta$:

  -- -------------------------------------------------------------- --
     $$\sqrt{K}\left(\frac{1}{\sqrt{N}}+\epsilon\right)=\Delta.$$   
  -- -------------------------------------------------------------- --

Solve for $\epsilon$:

  -- ---------------------------------------------------------- --
     $$\frac{1}{\sqrt{N}}+\epsilon=\frac{\Delta}{\sqrt{K}},$$   
  -- ---------------------------------------------------------- --

  -- ---------------------------------------------------------- --
     $$\epsilon=\frac{\Delta}{\sqrt{K}}-\frac{1}{\sqrt{N}}.$$   
  -- ---------------------------------------------------------- --

The condition $N>K/\Delta^{2}$ ensures $\epsilon>0$:

  -- --------------------------------------------------------------------------------------------------------------------------------- --
     $$\frac{\Delta}{\sqrt{K}}>\frac{1}{\sqrt{N}}\quad\text{implies}\quad\sqrt{N}\Delta>\sqrt{K}\quad\text{or}\quad N\Delta^{2}>K,$$   
  -- --------------------------------------------------------------------------------------------------------------------------------- --

which matches the theorem's condition. Since $\Delta>0$, this condition guarantees the bound is meaningful.
:::

::: {#A3.SSx3.p4 .ltx_para}
Substitute $\epsilon$ into the inequality:

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $$\mathbb{P}\left(\|\hat{\mathbf{p}}-\mathbf{p}\|_{1}\geq\Delta\right)\leq\exp\left(-N\left(\frac{\Delta}{\sqrt{K}}-\frac{1}{\sqrt{N}}\right)^{2}\right).$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------- --

Thus:

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $$\mathbb{P}\left(\|\hat{\mathbf{p}}-\mathbf{p}\|_{1}<\Delta\right)\geq 1-\exp\left(-N\left(\frac{\Delta}{\sqrt{K}}-\frac{1}{\sqrt{N}}\right)^{2}\right).$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------- --
:::
:::::::

:::: {#A3.SSx4 .section .ltx_subsection}
#### Step 4: Conclusion {#step-4-conclusion .ltx_title .ltx_title_subsection}

::: {#A3.SSx4.p1 .ltx_para}
From Lemma 2:

  -- ---------------------------------------------------------------------------------------------- --
     $$\mathbb{P}(y_{\mathrm{mv}}=1)\geq\mathbb{P}(\|\hat{\mathbf{p}}-\mathbf{p}\|_{1}<\Delta).$$   
  -- ---------------------------------------------------------------------------------------------- --

Combining with the probability bound, we derive:

  -- ------------------------------------------------------------------------------------------------------------------------- --
     $$\mathbb{P}(y_{\mathrm{mv}}=1)\geq 1-\exp\left(-N\left(\frac{\Delta}{\sqrt{K}}-\frac{1}{\sqrt{N}}\right)^{2}\right).$$   
  -- ------------------------------------------------------------------------------------------------------------------------- --

This matches the statement of Theorem 1. Proof completed. $\hfill\blacksquare$
:::
::::

::::::::: {#A3.SS2 .section .ltx_subsection}
#### [C.2 ]{.ltx_tag .ltx_tag_subsection}Proof of Lemma 1 {#c.2-proof-of-lemma-1 .ltx_title .ltx_title_subsection}

::: {#A3.SS2.p1 .ltx_para}
[Lemma 1. (Bayesian Conjugacy in Multi-Agent Debate)]{.ltx_text .ltx_font_bold} [At round $t$, after observing responses from its neighbors $\mathcal{N}(i)$, the agent $i$ aggregates these into a count vector $\mathbf{c}_{i,t}$ as in Definition 2. Then, by Bayesian conjugacy, the posterior distribution over $\boldsymbol{\theta}_{i,t}$ remains Dirichlet:]{.ltx_text .ltx_font_italic}

  -- ------------------------------------------------------------------------------------------------------------------------------------------ --
     $$\boldsymbol{\theta}_{i,t}\mid\{y_{j,t-1}\}_{j\in\mathcal{N}(i)}\sim\mathrm{Dirichlet}(\boldsymbol{\alpha}_{i,t-1}+\mathbf{c}_{i,t}).$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------ --
:::

::: {#A3.SS2.p2 .ltx_para}
[Proof of Lemma 1.]{.ltx_text .ltx_font_italic} We aim to show that the posterior distribution of $\boldsymbol{\theta}_{i,t}$, the agent $i$'s belief over the answer space at round $t$, given the observed responses $\{y_{j,t-1}\}_{j\in\mathcal{N}(i)}$ from its neighbors, is a Dirichlet distribution with updated parameters $\boldsymbol{\alpha}_{i,t-1}+\mathbf{c}_{i,t}$.
:::

::: {#A3.SS2.p3 .ltx_para}
First, consider the prior belief of agent $i$ at the start of round $t$, which is based on its belief from the previous round $t-1$. The prior distribution over the parameter $\boldsymbol{\theta}_{i,t}$ is:

  -- -------------------------------------------------------------------------------------------------------- --
     $$\boldsymbol{\theta}_{i,t}\sim\mathrm{Dirichlet}(\alpha_{i,t-1}^{(1)},\ldots,\alpha_{i,t-1}^{(K)}),$$   
  -- -------------------------------------------------------------------------------------------------------- --

where $\boldsymbol{\alpha}_{i,t-1}=(\alpha_{i,t-1}^{(1)},\ldots,\alpha_{i,t-1}^{(K)})$ are positive real numbers representing the prior parameters inherited from the previous round, and the density is proportional to:

  -- ---------------------------------------------------------------------------------------------------- --
     $$p(\boldsymbol{\theta}_{i,t})\propto\prod_{m=1}^{K}\theta_{i,t}^{(m)^{\alpha_{i,t-1}^{(m)}-1}},$$   
  -- ---------------------------------------------------------------------------------------------------- --

with $\sum_{m=1}^{K}\theta_{i,t}^{(m)}=1$.
:::

::: {#A3.SS2.p4 .ltx_para}
At round $t$, agent $i$ observes the responses $\{y_{j,t-1}\}_{j\in\mathcal{N}(i)}$ from its neighbors, where $y_{j,t-1}\in\{1,\ldots,K\}$ represents the answer chosen by neighbor $j$ at round $t-1$. According to Definition 2, these responses are aggregated into a count vector $\mathbf{c}_{i,t}=(c_{i,t}^{(1)},\ldots,c_{i,t}^{(k)})$, where:

  -- --------------------------------------------------------------------- --
     $$c_{i,t}^{(m)}=\sum_{j\in\mathcal{N}(i)}\mathbf{1}[y_{j,t-1}=m],$$   
  -- --------------------------------------------------------------------- --

and $\sum_{m=1}^{k}c_{i,t}^{(m)}=|\mathcal{N}(i)|$, the number of neighbors. Assuming each neighbor's response $y_{j,t-1}$ is drawn independently from a categorical distribution parameterized by $\boldsymbol{\theta}_{i,t}$, the likelihood of observing the count vector $\mathbf{c}_{i,t}$ given $\boldsymbol{\theta}_{i,t}$ follows a multinomial distribution:

  -- ---------------------------------------------------------------------------------------------------------------------- --
     $$\mathbf{c}_{i,t}\mid\boldsymbol{\theta}_{i,t}\sim\text{Multinomial}(|\mathcal{N}(i)|,\boldsymbol{\theta}_{i,t}),$$   
  -- ---------------------------------------------------------------------------------------------------------------------- --

with the likelihood proportional to:

  -- --------------------------------------------------------------------------------------------------------------- --
     $$p(\mathbf{c}_{i,t}\mid\boldsymbol{\theta}_{i,t})\propto\prod_{m=1}^{K}\theta_{i,t}^{(m)^{c_{i,t}^{(m)}}}.$$   
  -- --------------------------------------------------------------------------------------------------------------- --
:::

::: {#A3.SS2.p5 .ltx_para}
By Bayesian conjugacy, the posterior distribution of $\boldsymbol{\theta}_{i,t}$ given the observed counts is proportional to the product of the prior and the likelihood:

  -- ------------------------------------------------------------------------------------------------------------------------------------------------- --
     $$p(\boldsymbol{\theta}_{i,t}\mid\mathbf{c}_{i,t})\propto p(\boldsymbol{\theta}_{i,t})\cdot p(\mathbf{c}_{i,t}\mid\boldsymbol{\theta}_{i,t}).$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------- --

Substituting the expressions:

  -- --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $$p(\boldsymbol{\theta}_{i,t}\mid\mathbf{c}_{i,t})\propto\left(\prod_{m=1}^{K}\theta_{i,t}^{(m)^{\alpha_{i,t-1}^{(m)}-1}}\right)\cdot\left(\prod_{m=1}^{K}\theta_{i,t}^{(m)^{c_{i,t}^{(m)}}}\right)=\prod_{m=1}^{K}\theta_{i,t}^{(m)^{\alpha_{i,t-1}^{(m)}+c_{i,t}^{(m)}-1}}.$$   
  -- --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --

This form is characteristic of a Dirichlet distribution. Thus, the posterior is:

  -- -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $$\boldsymbol{\theta}_{i,t}\mid\{y_{j,t-1}\}_{j\in\mathcal{N}(i)}\sim\mathrm{Dirichlet}(\alpha_{i,t-1}^{(1)}+c_{i,t}^{(1)},\ldots,\alpha_{i,t-1}^{(K)}+c_{i,t}^{(K)}),$$   
  -- -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --

or equivalently:

  -- ------------------------------------------------------------------------------------------------------------------------------------------ --
     $$\boldsymbol{\theta}_{i,t}\mid\{y_{j,t-1}\}_{j\in\mathcal{N}(i)}\sim\mathrm{Dirichlet}(\boldsymbol{\alpha}_{i,t-1}+\mathbf{c}_{i,t}).$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------ --
:::

::: {#A3.SS2.p6 .ltx_para}
Since $\mathbf{c}_{i,t}$ is derived from the neighbor responses $\{y_{j,t-1}\}_{j\in\mathcal{N}(i)}$ as specified, and the update follows from the conjugacy property of the Dirichlet and multinomial distributions, the lemma holds. Proof completed. $\hfill\blacksquare$
:::
:::::::::

::::::: {#A3.SS3 .section .ltx_subsection}
#### [C.3 ]{.ltx_tag .ltx_tag_subsection}Proof of Theorem 2 {#c.3-proof-of-theorem-2 .ltx_title .ltx_title_subsection}

::: {#A3.SS3.p1 .ltx_para}
[Theorem 2. (Martingale Behavior of Multi-Agent Debate)]{.ltx_text .ltx_font_bold} [ For agent $i$, let $p_{t}:=\bar{{\theta}}_{i,t}^{(1)}$ denote its belief in the correct answer at debate round $t$. Under Bayesian conjugacy,]{.ltx_text .ltx_font_italic}

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --
     $$p_{t}:=\bar{\theta}_{i,t}^{(1)}=\frac{\alpha_{i,t}^{(1)}}{\sum_{j=1}^{K}\alpha_{i,t}^{(j)}}=\frac{\alpha_{i,t-1}^{(1)}+c_{i,t}^{(1)}}{\sum_{j=1}^{K}(\alpha_{i,t-1}^{(j)}+c_{i,t}^{(j)})}.$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --

[Then, sequence $\{p_{t}\}_{t\geq 0}$ forms a martingale. That is, the expected belief at the next round equals the current belief:]{.ltx_text .ltx_font_italic}

  -- ----------------------------------------------------------------------------- --
     $$\mathbb{E}[p_{t}\mid p_{t-1},\ldots,p_{0}]=p_{t-1},\;\forall_{t\geq 0}.$$   
  -- ----------------------------------------------------------------------------- --
:::

::: {#A3.SS3.p2 .ltx_para}
[Proof of Theorem 2.]{.ltx_text .ltx_font_italic} We prove that the sequence $\{p_{t}\}_{t\geq 0}$, where $p_{t}=\bar{\theta}_{i,t}^{(1)}$, is a martingale by showing that the conditional expectation of the belief in the correct answer at round $t+1$ given the filtration up to round $t$ equals the current belief. A sequence $\{X_{t}\}_{t\geq 0}$ is a martingale if $\mathbb{E}[X_{t+1}\mid\mathcal{F}_{t}]=X_{t}$, where $\mathcal{F}_{t}=\sigma(\bar{\theta}_{i,0}^{(1)},\ldots,\bar{\theta}_{i,t}^{(1)})$ represents the history of beliefs.
:::

::: {#A3.SS3.p3 .ltx_para}
By the problem statement, agent $i$ initializes its belief with $\boldsymbol{\theta}_{i,0}\sim\mathrm{Dirichlet}(\boldsymbol{\alpha}_{i,0})$, where $\boldsymbol{\alpha}_{i,0}=(\alpha_{1},\ldots,\alpha_{K})$ with $\alpha_{j}>0$. The initial expected belief is:

  -- ------------------------------------------------------------------------------------------------- --
     $$p_{0}=\bar{\theta}_{i,0}^{(1)}=\frac{\alpha_{i,0}^{(1)}}{\sum_{j=1}^{K}\alpha_{i,0}^{(j)}}.$$   
  -- ------------------------------------------------------------------------------------------------- --

At each round $t>0$, agent $i$ observes the responses $\{y_{j,t-1}\mid j\in\mathcal{N}(i)\}$ from its neighbors, aggregates them into a count vector $\mathbf{c}_{i,t}=(c_{i,t}^{(1)},\ldots,c_{i,t}^{(K)})$, where $c_{i,t}^{(m)}=\sum_{j\in\mathcal{N}(i)}\mathbf{1}[y_{j,t-1}=m]$, and updates its parameter vector as:

  -- ----------------------------------------------------------------------------- --
     $$\boldsymbol{\alpha}_{i,t}=\boldsymbol{\alpha}_{i,t-1}+\mathbf{c}_{i,t}.$$   
  -- ----------------------------------------------------------------------------- --

The belief in the correct answer (assumed to be class 1) at round $t$ is a random variable $\bar{\theta}_{i,t}^{(1)}$, and its expectation is:

  -- ------------------------------------------------------------------------------------------------- --
     $$p_{t}=\bar{\theta}_{i,t}^{(1)}=\frac{\alpha_{i,t}^{(1)}}{\sum_{j=1}^{K}\alpha_{i,t}^{(j)}}.$$   
  -- ------------------------------------------------------------------------------------------------- --
:::

::: {#A3.SS3.p4 .ltx_para}
To verify the martingale property, we compute the conditional expectation of $\bar{\theta}_{i,t+1}^{(1)}$ given $\mathcal{F}_{t}$. At round $t+1$, agent $i$ observes $\{y_{j,t}\mid j\in\mathcal{N}(i)\}$, forming a new count vector $\mathbf{c}_{i,t+1}$, where $c_{i,t+1}^{(1)}=\sum_{j\in\mathcal{N}(i)}\mathbf{1}[y_{j,t}=1]$. The updated parameter is:

  -- ------------------------------------------------------------------------------- --
     $$\boldsymbol{\alpha}_{i,t+1}=\boldsymbol{\alpha}_{i,t}+\mathbf{c}_{i,t+1},$$   
  -- ------------------------------------------------------------------------------- --

and the belief is:

  -- ------------------------------------------------------------------------------------------------------------------------------- --
     $$\bar{\theta}_{i,t+1}^{(1)}=\frac{\alpha_{i,t}^{(1)}+c_{i,t+1}^{(1)}}{\sum_{j=1}^{K}(\alpha_{i,t}^{(j)}+c_{i,t+1}^{(j)})}.$$   
  -- ------------------------------------------------------------------------------------------------------------------------------- --

The conditional expectation is:

  -- -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $$\mathbb{E}[\bar{\theta}_{i,t+1}^{(1)}\mid\mathcal{F}_{t}]=\mathbb{E}\left[\frac{\alpha_{i,t}^{(1)}+c_{i,t+1}^{(1)}}{\sum_{j=1}^{K}(\alpha_{i,t}^{(j)}+c_{i,t+1}^{(j)})}\mid\mathcal{F}_{t}\right].$$   
  -- -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --

Since $\boldsymbol{\alpha}_{i,t}$ is deterministic given $\mathcal{F}_{t}$, and $c_{i,t+1}^{(j)}$ depends on the random responses $y_{j,t}$, we need the expectation of the counts. Each $y_{j,t}\sim\text{Categorical}(\bar{\theta}_{j,t})$, and $\bar{\theta}_{j,t}^{(1)}$ is the belief of neighbor $j$ at round $t$, which is in $\mathcal{F}_{t}$ (assuming $\bar{\theta}_{j,t}$ is updated and observable up to $t$). Thus:

  -- ------------------------------------------------------------------------------------ --
     $$\mathbb{E}[\mathbf{1}[y_{j,t}=1]\mid\mathcal{F}_{t}]=\bar{\theta}_{j,t}^{(1)},$$   
  -- ------------------------------------------------------------------------------------ --

and:

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --
     $$\mathbb{E}[c_{i,t+1}^{(1)}\mid\mathcal{F}_{t}]=\mathbb{E}\left[\sum_{j\in\mathcal{N}(i)}\mathbf{1}[y_{j,t}=1]\mid\mathcal{F}_{t}\right]=\sum_{j\in\mathcal{N}(i)}\bar{\theta}_{j,t}^{(1)},$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --

  -- ------------------------------------------------------------------------------------------------------- --
     $$\mathbb{E}[c_{i,t+1}^{(j)}\mid\mathcal{F}_{t}]=\sum_{j\in\mathcal{N}(i)}\bar{\theta}_{j,t}^{(j)},$$   
  -- ------------------------------------------------------------------------------------------------------- --

and the denominator's expected value is:

  -- ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $$\mathbb{E}\left[\sum_{j=1}^{K}(\alpha_{i,t}^{(j)}+c_{i,t+1}^{(j)})\mid\mathcal{F}_{t}\right]=\sum_{j=1}^{K}\alpha_{i,t}^{(j)}+\sum_{j=1}^{K}\sum_{l\in\mathcal{N}(i)}\bar{\theta}_{l,t}^{(j)}.$$   
  -- ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --

However, since $\bar{\theta}_{j,t}$ is a probability vector summing to 1, and $c_{i,t+1}^{(j)}$ counts over $\mathcal{N}(i)$, the total expected count is $|\mathcal{N}(i)|$. The correct approach is to recognize that:

  -- -------------------------------------------------------------------------------------------------------------- --
     $$\bar{\theta}_{i,t+1}^{(1)}=\frac{\alpha_{i,t}^{(1)}+c_{i,t+1}^{(1)}}{\alpha_{i,t}^{0}+|\mathcal{N}(i)|},$$   
  -- -------------------------------------------------------------------------------------------------------------- --

where $\alpha_{i,t}^{0}=\sum_{j=1}^{K}\alpha_{i,t}^{(j)}$. The conditional expectation becomes:

  -- ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $$\mathbb{E}[\bar{\theta}_{i,t+1}^{(1)}\mid\mathcal{F}_{t}]=\frac{\alpha_{i,t}^{(1)}+\mathbb{E}[c_{i,t+1}^{(1)}\mid\mathcal{F}_{t}]}{\alpha_{i,t}^{0}+|\mathcal{N}(i)|}.$$   
  -- ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --

Substituting the expected count:

  -- ------------------------------------------------------------------------------------------------------- --
     $$\mathbb{E}[c_{i,t+1}^{(1)}\mid\mathcal{F}_{t}]=\sum_{j\in\mathcal{N}(i)}\bar{\theta}_{j,t}^{(1)},$$   
  -- ------------------------------------------------------------------------------------------------------- --

so:

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $$\mathbb{E}[\bar{\theta}_{i,t+1}^{(1)}\mid\mathcal{F}_{t}]=\frac{\alpha_{i,t}^{(1)}+\sum_{j\in\mathcal{N}(i)}\bar{\theta}_{j,t}^{(1)}}{\alpha_{i,t}^{0}+|\mathcal{N}(i)|}.$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --

For this to equal $\bar{\theta}_{i,t}^{(1)}=\frac{\alpha_{i,t}^{(1)}}{\alpha_{i,t}^{0}}$, we need:

  -- ----------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $$\frac{\alpha_{i,t}^{(1)}+\sum_{j\in\mathcal{N}(i)}\bar{\theta}_{j,t}^{(1)}}{\alpha_{i,t}^{0}+|\mathcal{N}(i)|}=\frac{\alpha_{i,t}^{(1)}}{\alpha_{i,t}^{0}},$$   
  -- ----------------------------------------------------------------------------------------------------------------------------------------------------------------- --

which implies:

  -- ---------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $$\alpha_{i,t}^{(1)}+\sum_{j\in\mathcal{N}(i)}\bar{\theta}_{j,t}^{(1)}=\frac{\alpha_{i,t}^{(1)}(\alpha_{i,t}^{0}+|\mathcal{N}(i)|)}{\alpha_{i,t}^{0}},$$   
  -- ---------------------------------------------------------------------------------------------------------------------------------------------------------- --

  -- ------------------------------------------------------------------------------------------------------------------------- --
     $$\sum_{j\in\mathcal{N}(i)}\bar{\theta}_{j,t}^{(1)}=\alpha_{i,t}^{(1)}\cdot\frac{|\mathcal{N}(i)|}{\alpha_{i,t}^{0}},$$   
  -- ------------------------------------------------------------------------------------------------------------------------- --

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $$\frac{1}{|\mathcal{N}(i)|}\sum_{j\in\mathcal{N}(i)}\bar{\theta}_{j,t}^{(1)}=\frac{\alpha_{i,t}^{(1)}}{\alpha_{i,t}^{0}}=\bar{\theta}_{i,t}^{(1)}.$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------- --

This condition holds if the average belief of the neighbors equals the agent's current belief, which is satisfied when agents start with homogeneous beliefs and the debate process preserves consistency. Thus, under this assumption, $\mathbb{E}[\bar{\theta}_{i,t+1}^{(1)}\mid\mathcal{F}_{t}]=\bar{\theta}_{i,t}^{(1)}$, and since $p_{t}=\bar{\theta}_{i,t}^{(1)}$, the sequence $\{p_{t}\}_{t\geq 0}$ is a martingale. Since $p_{t}$ is the expectation of a Dirichlet random variable and the update preserves the expected value under the specified condition, the theorem holds. Proof completed. $\hfill\blacksquare$
:::
:::::::
::::::::::::::::::::::::::::::::

::::::::::::::::::::::::: {#A4 .section .ltx_appendix}
### [Appendix D ]{.ltx_tag .ltx_tag_appendix}Special Case of Theorem 1 {#appendix-d-special-case-of-theorem-1 .ltx_title .ltx_title_appendix}

::: {#A4.p1 .ltx_para}
In this section, we analyze a special case of Theorem 1, where the probability of the correct answer, denoted by $\theta_{1}$, exceeds $\frac{1}{2}$. In other words, each agent independently makes a correct decision with probability at least 0.5. It naturally follows that $\frac{1}{2}\geq 1-\theta_{1}=\sum_{i=2}^{k}\theta_{i}\geq\max_{i\in\{2,\ldots,k\}}\theta_{i}$, indicating that the margin term $\Delta$ in Theorem 1 requires $\Delta\geq 0$, which does not rely on the number of agents $N$ or the set size of possible answers $k$.
:::

::: {#A4.p2 .ltx_para}
To formalize this, we define $p_{0}:=\theta_{1}$ to represent the initial probability that an individual agent selects the correct answer. Also, for simplicity, let $X_{1},X_{2},\ldots,X_{N}$ be independent Bernoulli random variables, where $X_{i}\sim\text{Bernoulli}(p_{0})$. The average correctness across agents is then given by the empirical mean $\bar{X}=\frac{1}{N}\sum_{i=1}^{N}X_{i}$, which corresponds to the fraction of agents who vote correctly. Under this setup, we analyze the lower bound of the Majority Voting success probability as follows:
:::

::: {#A4.p3 .ltx_para}
[Theorem 1.A (Majority Voting Success Probability)]{.ltx_text .ltx_font_bold}. [Let $X_{1},X_{2},\ldots,X_{n}$ be independent Bernoulli random variables with $X_{i}\sim\text{Bernoulli}(p_{0})$, where $p_{0}\in(0,1)$ is the probability that each agent is correct. Let $\bar{X}=\frac{1}{N}\sum_{i=1}^{N}X_{i}$ be the fraction of correct agents among $N$ independent agents. If $p_{0}>\frac{1}{2}$, the probability that a majority vote is successful is lower bounded as follows: ]{.ltx_text .ltx_font_italic}

  -- ----------------------------------------------------------------------------------------------------- --
     $$P\left(\bar{X}>\frac{1}{2}\right)\geq 1-\exp\left(-2N\left(p_{0}-\frac{1}{2}\right)^{2}\right).$$   
  -- ----------------------------------------------------------------------------------------------------- --
:::

::: {#A4.p4 .ltx_para}
[Proof of Theorem 1.A.]{.ltx_text .ltx_font_italic} Since $X_{i}\sim\text{Bernoulli}(p)$, we have $E[X_{i}]=p$, and the variables are independent and bounded: $0\leq X_{i}\leq 1$. Define $S_{n}=\sum_{i=1}^{n}X_{i}$, so $E[S_{n}]=np$, and $\bar{X}=\frac{S_{n}}{n}$, so $E[\bar{X}]=p$.
:::

::: {#A4.p5 .ltx_para}
We apply Hoeffding's Inequality to bound the deviation of $S_{n}$ from its mean. Hoeffding's Inequality states that for independent random variables $X_{i}$ with $a_{i}\leq X_{i}\leq b_{i}$, and $S_{n}=\sum_{i=1}^{n}X_{i}$:
:::

::: {#A4.p6 .ltx_para}
  -- ------------------------------------------------------------------------------------------------- --
     $$P(S_{n}-E[S_{n}]\geq t)\leq\exp\left(-\frac{2t^{2}}{\sum_{i=1}^{n}(b_{i}-a_{i})^{2}}\right)$$   
  -- ------------------------------------------------------------------------------------------------- --

  -- ------------------------------------------------------------------------------------------------- --
     $$P(S_{n}-E[S_{n}]\leq-t)\leq\exp\left(-\frac{2t^{2}}{\sum_{i=1}^{n}(b_{i}-a_{i})^{2}}\right)$$   
  -- ------------------------------------------------------------------------------------------------- --
:::

::: {#A4.p7 .ltx_para}
Here, $a_{i}=0$, $b_{i}=1$, so $(b_{i}-a_{i})^{2}=1$, and $\sum_{i=1}^{n}(b_{i}-a_{i})^{2}=n$. Let $t=\epsilon n$:
:::

::: {#A4.p8 .ltx_para}
  -- ------------------------------------------------------------------------------------------------------- --
     $$P(S_{n}-np\leq-\epsilon n)\leq\exp\left(-\frac{2(\epsilon n)^{2}}{n}\right)=\exp(-2n\epsilon^{2})$$   
  -- ------------------------------------------------------------------------------------------------------- --
:::

::: {#A4.p9 .ltx_para}
Rewrite in terms of $\bar{X}$:
:::

::: {#A4.p10 .ltx_para}
  -- ------------------------------------------------------------------------------------------------------- --
     $$P(S_{n}-np\leq-\epsilon n)=P\left(\frac{S_{n}}{n}-p\leq-\epsilon\right)=P(\bar{X}\leq p-\epsilon)$$   
  -- ------------------------------------------------------------------------------------------------------- --
:::

::: {#A4.p11 .ltx_para}
  -- -------------------------------------------------------- --
     $$P(\bar{X}\leq p-\epsilon)\leq\exp(-2n\epsilon^{2})$$   
  -- -------------------------------------------------------- --
:::

::: {#A4.p12 .ltx_para}
We want the majority vote to be successful, i.e., $\bar{X}>\frac{1}{2}$.
:::

::: {#A4.p13 .ltx_para}
First, consider the complementary event:
:::

::: {#A4.p14 .ltx_para}
  -- ------------------------------------------------------------------------------------------------- --
     $$P\left(\bar{X}\leq\frac{1}{2}\right)=P\left(\bar{X}\leq p-\left(p-\frac{1}{2}\right)\right)$$   
  -- ------------------------------------------------------------------------------------------------- --
:::

::: {#A4.p15 .ltx_para}
and set $\epsilon=p-\frac{1}{2}$:
:::

::: {#A4.p16 .ltx_para}
  -- ---------------------------------------------------------- --
     $$p-\epsilon=\frac{1}{2}\implies\epsilon=p-\frac{1}{2}$$   
  -- ---------------------------------------------------------- --
:::

::: {#A4.p17 .ltx_para}
Since $p>\frac{1}{2}$, $\epsilon>0$. Substitute $\epsilon$:
:::

::: {#A4.p18 .ltx_para}
  -- ------------------------------------------------------------------------------------------------ --
     $$P\left(\bar{X}\leq\frac{1}{2}\right)\leq\exp\left(-2n\left(p-\frac{1}{2}\right)^{2}\right)$$   
  -- ------------------------------------------------------------------------------------------------ --
:::

::: {#A4.p19 .ltx_para}
Thus, the probability of a successful majority vote is:
:::

::: {#A4.p20 .ltx_para}
  -- --------------------------------------------------------------------------------------------------------------------------------------- --
     $$P\left(\bar{X}>\frac{1}{2}\right)=1-P\left(\bar{X}\leq\frac{1}{2}\right)\geq 1-\exp\left(-2n\left(p-\frac{1}{2}\right)^{2}\right)$$   
  -- --------------------------------------------------------------------------------------------------------------------------------------- --
:::

::: {#A4.p21 .ltx_para}
Proof completed. $\hfill\blacksquare$
:::

::: {#A4.p22 .ltx_para}
Under this alternative formulation, the term $p-\frac{1}{2}$ serves as a conservative approximation of the margin $\Delta$, since all competing answer choices are necessarily assigned probabilities less than $\frac{1}{2}$. Then, note that the dominant term in the resulting lower bound remains consistent with that of the orignal Theorem 1, scaling as $\mathcal{O}(N\cdot\Delta^{2})$.
:::
:::::::::::::::::::::::::

:::: {#A5 .section .ltx_appendix}
### [Appendix E ]{.ltx_tag .ltx_tag_appendix}Martingale Process Empirical Investigation {#appendix-e-martingale-process-empirical-investigation .ltx_title .ltx_title_appendix}

::: {#A5.p1 .ltx_para}
Here, we provide the raw mean accuracy values used for Figure [[4]{.ltx_text .ltx_ref_tag}](#S4.F4 "Figure 4 ‣ Theorem 2. (Martingale Behavior of Multi-Agent Debate) ‣ 4 Theoretical Analysis ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref} is provided in Table [[6]{.ltx_text .ltx_ref_tag}](#A5.T6 "Table 6 ‣ Appendix E Martingale Process Empirical Investigation ‣ Appendix ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}.
:::

<figure id="A5.T6" class="ltx_table">
<div class="ltx_inline-block ltx_transformed_outer" style="width:411.9pt;height:109.8pt;vertical-align:-53.0pt;">
<span class="ltx_transformed_inner" style="transform:translate(-67.8pt,18.1pt) scale(0.752385300019829,0.752385300019829) ;"> </span>
<table class="ltx_tabular ltx_align_middle">
<tbody>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_tt">Rounds</td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">Arithmetics</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">GSM8K</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">Pro. Med.</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">Formal Logic</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">HellaSwag</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">CommonsenseQA</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">HH-RLHF</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">1</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.5400</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.7740</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.8022</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.4937</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.7900</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.8233</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.4613</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r">2</td>
<td class="ltx_td ltx_align_center">0.4880</td>
<td class="ltx_td ltx_align_center">0.7267</td>
<td class="ltx_td ltx_align_center">0.8015</td>
<td class="ltx_td ltx_align_center">0.4492</td>
<td class="ltx_td ltx_align_center">0.7880</td>
<td class="ltx_td ltx_align_center">0.8280</td>
<td class="ltx_td ltx_align_center">0.4600</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r">3</td>
<td class="ltx_td ltx_align_center">0.4880</td>
<td class="ltx_td ltx_align_center">0.7240</td>
<td class="ltx_td ltx_align_center">0.8029</td>
<td class="ltx_td ltx_align_center">0.4254</td>
<td class="ltx_td ltx_align_center">0.7880</td>
<td class="ltx_td ltx_align_center">0.8333</td>
<td class="ltx_td ltx_align_center">0.4547</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r">4</td>
<td class="ltx_td ltx_align_center">0.4860</td>
<td class="ltx_td ltx_align_center">0.7287</td>
<td class="ltx_td ltx_align_center">0.8015</td>
<td class="ltx_td ltx_align_center">0.4159</td>
<td class="ltx_td ltx_align_center">0.7880</td>
<td class="ltx_td ltx_align_center">0.8340</td>
<td class="ltx_td ltx_align_center">0.4560</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r">5</td>
<td class="ltx_td ltx_align_center">0.4940</td>
<td class="ltx_td ltx_align_center">0.7293</td>
<td class="ltx_td ltx_align_center">0.8015</td>
<td class="ltx_td ltx_align_center">0.4048</td>
<td class="ltx_td ltx_align_center">0.7873</td>
<td class="ltx_td ltx_align_center">0.8340</td>
<td class="ltx_td ltx_align_center">0.4560</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_tt">Rounds</td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">Arithmetics</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">GSM8K</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">Pro. Med.</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">Formal Logic</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">HellaSwag</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">CommonsenseQA</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">HH-RLHF</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r ltx_border_t">1</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.6340</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.6553</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.7199</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.3667</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.6487</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.7167</td>
<td class="ltx_td ltx_align_center ltx_border_t">0.5060</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r">2</td>
<td class="ltx_td ltx_align_center">0.6980</td>
<td class="ltx_td ltx_align_center">0.6533</td>
<td class="ltx_td ltx_align_center">0.6978</td>
<td class="ltx_td ltx_align_center">0.3492</td>
<td class="ltx_td ltx_align_center">0.6287</td>
<td class="ltx_td ltx_align_center">0.7307</td>
<td class="ltx_td ltx_align_center">0.5173</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r">3</td>
<td class="ltx_td ltx_align_center">0.6920</td>
<td class="ltx_td ltx_align_center">0.6520</td>
<td class="ltx_td ltx_align_center">0.6728</td>
<td class="ltx_td ltx_align_center">0.3397</td>
<td class="ltx_td ltx_align_center">0.5953</td>
<td class="ltx_td ltx_align_center">0.7220</td>
<td class="ltx_td ltx_align_center">0.5060</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_r">4</td>
<td class="ltx_td ltx_align_center">0.7300</td>
<td class="ltx_td ltx_align_center">0.6413</td>
<td class="ltx_td ltx_align_center">0.6662</td>
<td class="ltx_td ltx_align_center">0.3254</td>
<td class="ltx_td ltx_align_center">0.5820</td>
<td class="ltx_td ltx_align_center">0.7067</td>
<td class="ltx_td ltx_align_center">0.4987</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_center ltx_border_bb ltx_border_r">5</td>
<td class="ltx_td ltx_align_center ltx_border_bb">0.7280</td>
<td class="ltx_td ltx_align_center ltx_border_bb">0.6380</td>
<td class="ltx_td ltx_align_center ltx_border_bb">0.6581</td>
<td class="ltx_td ltx_align_center ltx_border_bb">0.3175</td>
<td class="ltx_td ltx_align_center ltx_border_bb">0.5647</td>
<td class="ltx_td ltx_align_center ltx_border_bb">0.6953</td>
<td class="ltx_td ltx_align_center ltx_border_bb">0.4973</td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 6</span>: </span><span class="ltx_text" style="font-size:90%;">Raw mean accuracy values for Figure <a href="#S4.F4" class="ltx_ref" title="Figure 4 ‣ Theorem 2. (Martingale Behavior of Multi-Agent Debate) ‣ 4 Theoretical Analysis ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"><span class="ltx_text ltx_ref_tag">4</span></a>. Top row is from Qwen2.5-7B-Instruct, and the bottom row is from Llama3.1-8B-Instruct.</span></figcaption>
</figure>
::::

:::::: {#A6 .section .ltx_appendix}
### [Appendix F ]{.ltx_tag .ltx_tag_appendix}Proper Evaluation Matters {#appendix-f-proper-evaluation-matters .ltx_title .ltx_title_appendix}

::: {#A6.p1 .ltx_para}
Another key takeaway from our study is that careful evaluation is critical for accurately assessing the utility of MAD. We find that the method used to extract final answers from free-form model responses can significantly affect measured performance---sometimes even reversing conclusions. While prior works have reported consistent gains from MAD over majority voting \[[2](#bib.bib2){.ltx_ref}, [3](#bib.bib3){.ltx_ref}, [10](#bib.bib10){.ltx_ref}\], we find these results may be partially driven by error-prone answer extraction, where rule-based parsing can fail even when the model's response is correct. For example, a model's output may be correct, but incorrectly marked as incorrect purely due to failures in parsing, rather than actual reasoning mistakes.
:::

<figure id="A6.T7" class="ltx_table ltx_align_floatright">
<div class="ltx_inline-block ltx_align_center ltx_transformed_outer" style="width:385.0pt;height:61.2pt;vertical-align:-28.1pt;">
<span class="ltx_transformed_inner" style="transform:translate(0.0pt,0.0pt) scale(1,1) ;"> </span>
<p><span class="ltx_tabular ltx_align_middle"> <span class="ltx_tr"> <span class="ltx_td ltx_border_tt" style="padding-left:5.0pt;padding-right:5.0pt;"></span> <span class="ltx_td ltx_align_center ltx_border_tt ltx_colspan ltx_colspan_2" style="padding-left:5.0pt;padding-right:5.0pt;"><span class="ltx_text ltx_font_bold">GSM8K</span></span></span> <span class="ltx_tr"> <span class="ltx_td ltx_align_center" style="padding-left:5.0pt;padding-right:5.0pt;"><span class="ltx_text ltx_font_bold">Methods</span></span> <span class="ltx_td ltx_align_center" style="padding-left:5.0pt;padding-right:5.0pt;"><span class="ltx_text ltx_font_bold">Ours</span></span> <span class="ltx_td ltx_align_center" style="padding-left:5.0pt;padding-right:5.0pt;"><span class="ltx_text ltx_font_bold">Prior</span> [<a href="#bib.bib2" class="ltx_ref">2</a>]</span></span> <span class="ltx_tr"> <span class="ltx_td ltx_align_left ltx_border_t" style="padding-left:5.0pt;padding-right:5.0pt;">Single-agent</span> <span class="ltx_td ltx_align_center ltx_border_t" style="padding-left:5.0pt;padding-right:5.0pt;">0.8713 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .00</span></span> <span class="ltx_td ltx_align_center ltx_border_t" style="padding-left:5.0pt;padding-right:5.0pt;">0.6620 <span class="math inline">±</span><span class="ltx_text" style="font-size:70%;"> .01</span></span></span> <span class="ltx_tr"> <span class="ltx_td ltx_align_left" style="padding-left:5.0pt;padding-right:5.0pt;">MAD (<span class="math inline"><em>T</em> = 2</span>)</span> <span class="ltx_td ltx_align_center" style="padding-left:5.0pt;padding-right:5.0pt;">0.8867</span> <span class="ltx_td ltx_align_center" style="padding-left:5.0pt;padding-right:5.0pt;"><span class="ltx_text ltx_font_bold">0.7533</span></span></span> <span class="ltx_tr"> <span class="ltx_td ltx_align_left ltx_border_bb" style="padding-left:5.0pt;padding-right:5.0pt;">Majority Voting</span> <span class="ltx_td ltx_align_center ltx_border_bb" style="padding-left:5.0pt;padding-right:5.0pt;"><span class="ltx_text ltx_font_bold">0.9400</span></span> <span class="ltx_td ltx_align_center ltx_border_bb" style="padding-left:5.0pt;padding-right:5.0pt;">0.6700</span></span> </span></p>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 7</span>: </span><span class="ltx_text ltx_font_bold" style="font-size:90%;">Effect of different answer extractors.<span class="ltx_text ltx_font_medium"> Qwen2.5-7B-Instruct on Decentralized MAD.</span></span></figcaption>
</figure>

::: {#A6.p2 .ltx_para}
To improve the reliability of answer extraction, we explicitly instructed each agent to append its final answer using a standardized format--for example, "{final answer: $\hat{y}$}\". This strategy substantially reduces parsing failures and yields more reliable evaluations (see Appendix [[B.1]{.ltx_text .ltx_ref_tag}](#A2.SS1 "B.1 MAD Templates ‣ Appendix B Prompt Templates ‣ Appendix ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref} for prompt details). Once final answers are extracted from each agent using the same protocol, we select the majority answer as the final response.
:::

::: {#A6.p3 .ltx_para}
Table [[7]{.ltx_text .ltx_ref_tag}](#A6.T7 "Table 7 ‣ Appendix F Proper Evaluation Matters ‣ Appendix ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref} shows that the extraction strategy significantly impacts performance. Our evaluation protocol improves single-agent accuracy and, on GSM8K, even outperforms MAD when the latter is evaluated using the prior strategy from \[[2](#bib.bib2){.ltx_ref}\]. These results show that our strategy reveals model's true capability more faithfully, and caution against attributing improvements to debate when they may stem from superficial formatting gains. Without rigorous and consistent evaluation, we may incorrectly estimate MAD's benefits and obscure whether inter-agent communication truly enhances decision quality.
:::
::::::

:::: {#A7 .section .ltx_appendix}
### [Appendix G ]{.ltx_tag .ltx_tag_appendix}Closed-source LLM Evaluation {#appendix-g-closed-source-llm-evaluation .ltx_title .ltx_title_appendix}

::: {#A7.p1 .ltx_para}
We extend our experiments to a closed-source LLM setting. Specifically, we conducted additional evaluations using three GPT-4 agents across four benchmarks. In Table [[8]{.ltx_text .ltx_ref_tag}](#A7.T8 "Table 8 ‣ Appendix G Closed-source LLM Evaluation ‣ Appendix ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}, the overall trends remain consistent with those observed in our open-source model experiments, supporting the generality of our findings.
:::

<figure id="A7.T8" class="ltx_table">
<div class="ltx_inline-block ltx_align_center ltx_transformed_outer" style="width:325.2pt;height:46.6pt;vertical-align:-21.4pt;">
<span class="ltx_transformed_inner" style="transform:translate(-51.0pt,7.3pt) scale(0.761100749613152,0.761100749613152) ;"> </span>
<table class="ltx_tabular ltx_align_middle">
<tbody>
<tr class="ltx_tr">
<td class="ltx_td ltx_border_tt"></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">Arithmetics</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">CSQA</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">HellaSwag</span></td>
<td class="ltx_td ltx_align_center ltx_border_tt"><span class="ltx_text ltx_font_bold">HH-RLHF</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_t">Majority Voting</td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text ltx_font_bold">0.9967</span></td>
<td class="ltx_td ltx_align_center ltx_border_t">0.8721</td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text ltx_font_bold">0.9078</span></td>
<td class="ltx_td ltx_align_center ltx_border_t"><span class="ltx_text ltx_font_bold">0.5612</span></td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Decentralized MAD (<span class="math inline"><em>T</em> = 1</span>)</td>
<td class="ltx_td ltx_align_center">0.9867</td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">0.8788</span></td>
<td class="ltx_td ltx_align_center"><span class="ltx_text ltx_font_bold">0.9078</span></td>
<td class="ltx_td ltx_align_center">0.5580</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left">Decentralized MAD (<span class="math inline"><em>T</em> = 2</span>)</td>
<td class="ltx_td ltx_align_center">0.9867</td>
<td class="ltx_td ltx_align_center">0.8784</td>
<td class="ltx_td ltx_align_center">0.9044</td>
<td class="ltx_td ltx_align_center">0.5577</td>
</tr>
<tr class="ltx_tr">
<td class="ltx_td ltx_align_left ltx_border_bb">Decentralized MAD (<span class="math inline"><em>T</em> = 3</span>)</td>
<td class="ltx_td ltx_align_center ltx_border_bb">0.9833</td>
<td class="ltx_td ltx_align_center ltx_border_bb">0.8780</td>
<td class="ltx_td ltx_align_center ltx_border_bb">0.9044</td>
<td class="ltx_td ltx_align_center ltx_border_bb">0.5459</td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span class="ltx_text" style="font-size:90%;">Table 8</span>: </span><span class="ltx_text ltx_font_bold" style="font-size:90%;">Further experiments on GPT-4</span></figcaption>
</figure>
::::

:::: {#A8 .section .ltx_appendix}
### [Appendix H ]{.ltx_tag .ltx_tag_appendix}Limitations and Future Works {#appendix-h-limitations-and-future-works .ltx_title .ltx_title_appendix}

::: {#A8.p1 .ltx_para}
As discussed in Section [[2]{.ltx_text .ltx_ref_tag}](#S2 "2 Preliminaries ‣ Debate or Vote: Which Yields Better Decisions in Multi-Agent Large Language Models?"){.ltx_ref}, our study primarily focuses on the Simultaneous Talk protocol \[[3](#bib.bib3){.ltx_ref}\], where all agents generate and share their responses concurrently in each debate round. While this setting is widely adopted in prior work, it does not capture the full spectrum of possible communication strategies within multi-agent systems. Alternative protocols, such as One-by-One, where agents respond sequentially, or Simultaneous-Talk-with-Summarizer, where a summarizer agent oversees and summarizes the state of the debate, introduce different dynamics that may influence the rates of subversion and correction. Investigating these alternative protocols in depth remains an important direction for future work. Furthermore, our theoretical framework relies on the assumption of agent homogeneity and may not directly generalize to heterogeneous settings. A promising direction for future work is to extend the martingale analysis to account for heterogeneous agents with differing prior beliefs or reasoning capabilities.
:::
::::
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

::: ar5iv-footer
[◄](/html/2508.17535){.ar5iv-nav-button .ar5iv-nav-button-prev} [![ar5iv homepage](/assets/ar5iv.png){height="40"}](/){.ar5iv-home-button} [Feeling\
lucky?](/feeling_lucky){.ar5iv-text-button} [](/land_of_honey_and_milk){rel="nofollow" aria-hidden="true" tabindex="-1"} [Conversion\
report](/log/2508.17536){.ar5iv-text-button .ar5iv-severity-warning} [Report\
an issue](https://github.com/dginev/ar5iv/issues/new?template=improve-article--arxiv-id-.md&title=Improve+article+2508.17536){.ar5iv-text-button target="_blank"} [View original\
on arXiv](https://arxiv.org/abs/2508.17536){.ar5iv-text-button .arxiv-ui-theme}[►](/html/2508.17537){.ar5iv-nav-button .ar5iv-nav-button-next}
:::

[[]{.color-scheme-icon}](javascript:toggleColorScheme() "Toggle ar5iv color scheme"){.ar5iv-toggle-color-scheme} [Copyright](https://arxiv.org/help/license){.ar5iv-footer-button target="_blank"} [Privacy Policy](https://arxiv.org/help/policies/privacy_policy){.ar5iv-footer-button target="_blank"}

::: ltx_page_logo
Generated on Fri Sep 5 13:28:26 2025 by [[L[a]{.ltx_font_smallcaps style="position:relative; bottom:2.2pt;"}T[e]{.ltx_font_smallcaps style="font-size:120%;position:relative; bottom:-0.2ex;"}]{style="letter-spacing:-0.2em; margin-right:0.1em;"}[XML]{style="font-size:90%; position:relative; bottom:-0.2ex;"}![Mascot Sammy](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAsAAAAOCAYAAAD5YeaVAAAAAXNSR0IArs4c6QAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAALEwAACxMBAJqcGAAAAAd0SU1FB9wKExQZLWTEaOUAAAAddEVYdENvbW1lbnQAQ3JlYXRlZCB3aXRoIFRoZSBHSU1Q72QlbgAAAdpJREFUKM9tkL+L2nAARz9fPZNCKFapUn8kyI0e4iRHSR1Kb8ng0lJw6FYHFwv2LwhOpcWxTjeUunYqOmqd6hEoRDhtDWdA8ApRYsSUCDHNt5ul13vz4w0vWCgUnnEc975arX6ORqN3VqtVZbfbTQC4uEHANM3jSqXymFI6yWazP2KxWAXAL9zCUa1Wy2tXVxheKA9YNoR8Pt+aTqe4FVVVvz05O6MBhqUIBGk8Hn8HAOVy+T+XLJfLS4ZhTiRJgqIoVBRFIoric47jPnmeB1mW/9rr9ZpSSn3Lsmir1fJZlqWlUonKsvwWwD8ymc/nXwVBeLjf7xEKhdBut9Hr9WgmkyGEkJwsy5eHG5vN5g0AKIoCAEgkEkin0wQAfN9/cXPdheu6P33fBwB4ngcAcByHJpPJl+fn54mD3Gg0NrquXxeLRQAAwzAYj8cwTZPwPH9/sVg8PXweDAauqqr2cDjEer1GJBLBZDJBs9mE4zjwfZ85lAGg2+06hmGgXq+j3+/DsixYlgVN03a9Xu8jgCNCyIegIAgx13Vfd7vdu+FweG8YRkjXdWy329+dTgeSJD3ieZ7RNO0VAXAPwDEAO5VKndi2fWrb9jWl9Esul6PZbDY9Go1OZ7PZ9z/lyuD3OozU2wAAAABJRU5ErkJggg==)](http://dlmf.nist.gov/LaTeXML/){.ltx_LaTeXML_logo target="_blank"}
:::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
