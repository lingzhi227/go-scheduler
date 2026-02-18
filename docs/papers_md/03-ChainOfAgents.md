::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_main
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_content
# Chain of Agents: Large Language Models Collaborating on Long-Context Tasks {#chain-of-agents-large-language-models-collaborating-on-long-context-tasks .ltx_title .ltx_title_document}

::: ltx_authors
[ [Yusen Zhang^[♣]{#id9.9.id1.1 .ltx_text .ltx_font_italic}^ , Ruoxi Sun^♢^, Yanfei Chen^♢^, Tomas Pfister^♢^, Rui Zhang^[♣†]{#id13.13.id5.1 .ltx_text .ltx_font_italic}^, Sercan Ö. Arik^♢^\
^♣^ Penn State University, ^♢^ Google Cloud AI Research\
[{yfz5488, rmz5227}@psu.edu, {ruoxis, yanfeichen, tpfister, soarik}@google.com]{#id17.17.id9 .ltx_text .ltx_font_typewriter}\
]{.ltx_personname}[Work done while the author was a student researcher at Google Cloud AI Research. Last authors]{.ltx_author_notes}]{.ltx_creator .ltx_role_author}
:::

::: ltx_abstract
###### Abstract {#abstract .ltx_title .ltx_title_abstract}

Addressing the challenge of effectively processing long contexts has become a critical issue for Large Language Models (LLMs). Two common strategies have emerged: 1) reducing the input length, such as retrieving relevant chunks by Retrieval-Augmented Generation (RAG), and 2) expanding the context window limit of LLMs. However, both strategies have drawbacks: input reduction has no guarantee of covering the part with needed information, while window extension struggles with focusing on the pertinent information for solving the task. To mitigate these limitations, we propose [Chain-of-Agents (CoA)]{#id18.id1.1 .ltx_text .ltx_font_italic}, a novel framework that harnesses multi-agent collaboration through natural language to enable information aggregation and context reasoning across various LLMs over long-context tasks. CoA consists of multiple worker agents who sequentially communicate to handle different segmented portions of the text, followed by a manager agent who synthesizes these contributions into a coherent final output. CoA processes the entire input by interleaving reading and reasoning, and it mitigates long context focus issues by assigning each agent a short context. We perform comprehensive evaluation of CoA on a wide range of long-context tasks in question answering, summarization, and code completion, demonstrating significant improvements by up to 10% over strong baselines of RAG, Full-Context, and multi-agent LLMs.
:::

:::::::: {#S1 .section .ltx_section}
## [1 ]{.ltx_tag .ltx_tag_section}Introduction {#introduction .ltx_title .ltx_title_section}

::: {#S1.p1 .ltx_para}
Despite their impressive performance across a wide range of scenarios, LLMs struggle with tasks that involve long contexts \[[8](#bib.bib8){.ltx_ref}, [57](#bib.bib57){.ltx_ref}, [52](#bib.bib52){.ltx_ref}\]. Numerous application scenarios demand extremely long contexts, such as question answering \[[78](#bib.bib78){.ltx_ref}, [20](#bib.bib20){.ltx_ref}, [63](#bib.bib63){.ltx_ref}\], document and dialogue summarization \[[23](#bib.bib23){.ltx_ref}, [84](#bib.bib84){.ltx_ref}, [83](#bib.bib83){.ltx_ref}, [82](#bib.bib82){.ltx_ref}, [12](#bib.bib12){.ltx_ref}\], and code completion \[[18](#bib.bib18){.ltx_ref}, [39](#bib.bib39){.ltx_ref}\], where the inputs contain entire books \[[29](#bib.bib29){.ltx_ref}, [30](#bib.bib30){.ltx_ref}\] and long articles \[[14](#bib.bib14){.ltx_ref}\].
:::

::: {#S1.p2 .ltx_para}
To tackle the challenge with long context tasks, two major directions have been explored as shown in Table [[1]{.ltx_text .ltx_ref_tag}](#S1.T1 "Table 1 ‣ 1 Introduction ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}: [input reduction]{#S1.p2.1.1 .ltx_text .ltx_font_italic} and [window extension]{#S1.p2.1.2 .ltx_text .ltx_font_italic}. [Input reduction]{#S1.p2.1.3 .ltx_text .ltx_font_italic} reduces the length of the input context before feeding to downstream LLMs. Truncation approaches \[[1](#bib.bib1){.ltx_ref}, [61](#bib.bib61){.ltx_ref}\] directly truncate the input. Retrieval Augmented Generation (RAG) \[[74](#bib.bib74){.ltx_ref}\] extends this direction by retrieving the most relevant chunks through embedding similarity. However, because of low retrieval accuracy, LLMs could receive an incomplete context for solving the task, hurting performance. [Window extension]{#S1.p2.1.4 .ltx_text .ltx_font_italic} extends the context window of LLMs via finetuning to consume the whole input \[[13](#bib.bib13){.ltx_ref}, [40](#bib.bib40){.ltx_ref}, [43](#bib.bib43){.ltx_ref}\]. For example, Claude-3 \[[5](#bib.bib5){.ltx_ref}\] directly allows reading 200k tokens for each input. However, when the window becomes longer, LLMs struggle to focus on the needed information to solve the task, suffering from ineffective context utilization such as the ["lost in the middle"]{#S1.p2.1.5 .ltx_text .ltx_font_italic} issue \[[33](#bib.bib33){.ltx_ref}, [3](#bib.bib3){.ltx_ref}, [38](#bib.bib38){.ltx_ref}\].
:::

<figure id="S1.T1" class="ltx_table">
<div id="S1.T1.9" class="ltx_inline-block ltx_transformed_outer" style="width:433.6pt;height:70.2pt;vertical-align:-0.0pt;">
<span class="ltx_transformed_inner" style="transform:translate(-116.7pt,18.9pt) scale(0.650165833043877,0.650165833043877) ;"> </span>
<table id="S1.T1.9.5" class="ltx_tabular ltx_align_middle">
<tbody>
<tr id="S1.T1.9.5.6" class="ltx_tr">
<td id="S1.T1.9.5.6.1" class="ltx_td ltx_align_left ltx_border_tt">Category</td>
<td id="S1.T1.9.5.6.2" class="ltx_td ltx_align_left ltx_border_tt">Example Work</td>
<td id="S1.T1.9.5.6.3" class="ltx_td ltx_align_center ltx_border_tt">Rec.</td>
<td id="S1.T1.9.5.6.4" class="ltx_td ltx_align_center ltx_border_tt">Foc.</td>
<td id="S1.T1.9.5.6.5" class="ltx_td ltx_align_center ltx_border_tt">No Train</td>
<td id="S1.T1.9.5.6.6" class="ltx_td ltx_align_center ltx_border_tt">Read</td>
<td id="S1.T1.9.5.6.7" class="ltx_td ltx_align_center ltx_border_tt">Agent</td>
<td id="S1.T1.9.5.6.8" class="ltx_td ltx_align_center ltx_border_tt">Applicability</td>
<td id="S1.T1.9.5.6.9" class="ltx_td ltx_align_center ltx_border_tt">Inter.</td>
</tr>
<tr id="S1.T1.5.1.1" class="ltx_tr">
<td rowspan="2" id="S1.T1.5.1.1.2" class="ltx_td ltx_align_left ltx_border_t"><span id="S1.T1.5.1.1.2.1" class="ltx_text">Input Reduction</span></td>
<td id="S1.T1.5.1.1.3" class="ltx_td ltx_align_left ltx_border_t">Truncation [<a href="#bib.bib45" class="ltx_ref">45</a>]</td>
<td id="S1.T1.5.1.1.4" class="ltx_td ltx_align_center ltx_border_t">✗</td>
<td id="S1.T1.5.1.1.5" class="ltx_td ltx_align_center ltx_border_t">✓</td>
<td id="S1.T1.5.1.1.6" class="ltx_td ltx_align_center ltx_border_t">✓</td>
<td id="S1.T1.5.1.1.1" class="ltx_td ltx_align_center ltx_border_t"><span class="math inline"><em>k</em></span></td>
<td id="S1.T1.5.1.1.7" class="ltx_td ltx_align_center ltx_border_t">Single</td>
<td id="S1.T1.5.1.1.8" class="ltx_td ltx_align_center ltx_border_t">Generic</td>
<td id="S1.T1.5.1.1.9" class="ltx_td ltx_align_center ltx_border_t">Low</td>
</tr>
<tr id="S1.T1.6.2.2" class="ltx_tr">
<td id="S1.T1.6.2.2.2" class="ltx_td ltx_align_left">RAG [<a href="#bib.bib74" class="ltx_ref">74</a>]</td>
<td id="S1.T1.6.2.2.3" class="ltx_td ltx_align_center">✗</td>
<td id="S1.T1.6.2.2.4" class="ltx_td ltx_align_center">✓</td>
<td id="S1.T1.6.2.2.5" class="ltx_td ltx_align_center">✗</td>
<td id="S1.T1.6.2.2.1" class="ltx_td ltx_align_center"><span class="math inline"><em>n</em> + <em>k</em></span></td>
<td id="S1.T1.6.2.2.6" class="ltx_td ltx_align_center">Single</td>
<td id="S1.T1.6.2.2.7" class="ltx_td ltx_align_center">Query-based</td>
<td id="S1.T1.6.2.2.8" class="ltx_td ltx_align_center">Medium</td>
</tr>
<tr id="S1.T1.7.3.3" class="ltx_tr">
<td rowspan="2" id="S1.T1.7.3.3.2" class="ltx_td ltx_align_left ltx_border_t"><span id="S1.T1.7.3.3.2.1" class="ltx_text">Window Extension</span></td>
<td id="S1.T1.7.3.3.3" class="ltx_td ltx_align_left ltx_border_t">Position Interpolation [<a href="#bib.bib13" class="ltx_ref">13</a>]</td>
<td id="S1.T1.7.3.3.4" class="ltx_td ltx_align_center ltx_border_t">✓</td>
<td id="S1.T1.7.3.3.5" class="ltx_td ltx_align_center ltx_border_t">✗</td>
<td id="S1.T1.7.3.3.6" class="ltx_td ltx_align_center ltx_border_t">✗</td>
<td id="S1.T1.7.3.3.1" class="ltx_td ltx_align_center ltx_border_t"><span class="math inline"><em>n</em></span></td>
<td id="S1.T1.7.3.3.7" class="ltx_td ltx_align_center ltx_border_t">Single</td>
<td id="S1.T1.7.3.3.8" class="ltx_td ltx_align_center ltx_border_t">Generic</td>
<td id="S1.T1.7.3.3.9" class="ltx_td ltx_align_center ltx_border_t">Low</td>
</tr>
<tr id="S1.T1.8.4.4" class="ltx_tr">
<td id="S1.T1.8.4.4.2" class="ltx_td ltx_align_left">Long Context [<a href="#bib.bib5" class="ltx_ref">5</a>]</td>
<td id="S1.T1.8.4.4.3" class="ltx_td ltx_align_center">✓</td>
<td id="S1.T1.8.4.4.4" class="ltx_td ltx_align_center">✗</td>
<td id="S1.T1.8.4.4.5" class="ltx_td ltx_align_center">✗</td>
<td id="S1.T1.8.4.4.1" class="ltx_td ltx_align_center"><span class="math inline"><em>n</em></span></td>
<td id="S1.T1.8.4.4.6" class="ltx_td ltx_align_center">Single</td>
<td id="S1.T1.8.4.4.7" class="ltx_td ltx_align_center">Generic</td>
<td id="S1.T1.8.4.4.8" class="ltx_td ltx_align_center">Low</td>
</tr>
<tr id="S1.T1.9.5.5" class="ltx_tr">
<td id="S1.T1.9.5.5.2" class="ltx_td ltx_align_left ltx_border_bb ltx_border_t">Multi-agent LLMs</td>
<td id="S1.T1.9.5.5.3" class="ltx_td ltx_align_left ltx_border_bb ltx_border_t">Chain-of-Agents (Ours)</td>
<td id="S1.T1.9.5.5.4" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t">✓</td>
<td id="S1.T1.9.5.5.5" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t">✓</td>
<td id="S1.T1.9.5.5.6" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t">✓</td>
<td id="S1.T1.9.5.5.1" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t"><span class="math inline"><em>n</em></span></td>
<td id="S1.T1.9.5.5.7" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t">Multiple</td>
<td id="S1.T1.9.5.5.8" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t">Generic</td>
<td id="S1.T1.9.5.5.9" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t">High</td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span id="S1.T1.11.3.1" class="ltx_text" style="font-size:90%;">Table 1</span>: </span><span id="S1.T1.4.2" class="ltx_text" style="font-size:90%;">Comparison between Chain-of-Agents and prior methods for long-context tasks. Rec./Foc.: being able to mitigate inaccurate receptive field/long context focusing issues. Read: the number of tokens as model input, where <span class="math inline"><em>n</em></span> is the total input length, <span class="math inline"><em>k</em></span> is the context window limit of LLMs. Inter.: the interpretability of the approach. Note that RAG is ‘medium interpretable’ because of the re-ranked chunks.</span></figcaption>
</figure>

<figure id="S1.F1" class="ltx_figure">
<img src="/html/2406.02818/assets/figures/CoA.png" id="S1.F1.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="598" height="274" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span id="S1.F1.2.1.1" class="ltx_text" style="font-size:90%;">Figure 1</span>: </span><span id="S1.F1.3.2" class="ltx_text" style="font-size:90%;">Overview of Chain-of-Agents, a training free, task agnostic, and highly-interpretable framework that harnesses multi-agent collaboration for long-context tasks. It consists of multiple worker agents who sequentially communicate to handle different segmented portions of the text, followed by a manager agent who synthesizes these contributions into a coherent final output.</span></figcaption>
</figure>

::: {#S1.p3 .ltx_para}
Motivated by the aforementioned challenges, we propose a novel framework, [Chain-of-Agents (CoA)]{#S1.p3.1.1 .ltx_text .ltx_font_italic}, inspired by human-like processing of long-context tasks. The key idea of CoA is to harness multi-agent communication to enable information aggregation and context reasoning capabilities across different LLMs. As shown in Figure [[1]{.ltx_text .ltx_ref_tag}](#S1.F1 "Figure 1 ‣ 1 Introduction ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}, CoA contains two stages. In stage 1, a series of worker agents in charge of different chunks of long context collaborate and aggregate evidence for answering the given query. To this end, the workers read and process sequentially, each receiving the message from previous worker and transferring the useful updated information to the next. In stage 2, the manager agent receives the complete evidence from last worker agent and generates the final response.
:::

::: {#S1.p4 .ltx_para}
As shown in Table [[1]{.ltx_text .ltx_ref_tag}](#S1.T1 "Table 1 ‣ 1 Introduction ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}, CoA is a training free, task agnostic, and highly interpretable framework processing entire "receptive field" by interleaved reading-processing and mitigating the long context focusing issue by assigning each agent a short context. Different from [input reduction]{#S1.p4.4.1 .ltx_text .ltx_font_italic} where LLMs need to start processing with low receptive field over reduced inputs ("read-then-process"), workers in CoA start to process each chunk before reading all input ("interleaved read-process"), tackling the problems that input reduction struggles with, such as, generic summarization or counting of passages \[[6](#bib.bib6){.ltx_ref}\]. Different from [context extension]{#S1.p4.4.2 .ltx_text .ltx_font_italic}, CoA leverages the capability of communication rather than trying to feed many tokens into an LLM. This is a more natural solution for complex tasks because we assume that each LLM has its limit and there are always complex context tasks surpassing its limit. Compared with Full-Context, CoA is also cost effective by reducing time complexity from $n^{2}$ to $n\hspace{0pt}k$, where $n$ is input tokens and $k$ is the context limit of LLMs.
:::

::: {#S1.p5 .ltx_para}
We conduct intensive experiments on [nine datasets]{#S1.p5.1.1 .ltx_text .ltx_font_italic}, including question answering, summarization, and code completion tasks with [six LLMs]{#S1.p5.1.2 .ltx_text .ltx_font_italic}, with PaLM 2 \[[4](#bib.bib4){.ltx_ref}\], Gemini \[[61](#bib.bib61){.ltx_ref}\], and Claude 3 \[[5](#bib.bib5){.ltx_ref}\] models. We compare CoA with two strong baselines chosen from [input reduction]{#S1.p5.1.3 .ltx_text .ltx_font_italic} and [window extension]{#S1.p5.1.4 .ltx_text .ltx_font_italic} approaches, respectively: (i) RAG, which uses a state-of-the-art retriever to obtain the most relevant information to feed into the LLM and (ii) Full-Context (Vanilla), which feeds all input into the LLM until reaching the window limit. Our results show that on all nine datasets, CoA obtains significant improvement over all baselines by up to 10%. Noting that there is not enough research on multi-agent for long context tasks, we carefully create two multi-agent baselines, including a hierarchical structure and result merging approach to further demonstrate that CoA is superior among other possible multi-agent frameworks.
:::
::::::::

::::::::: {#S2 .section .ltx_section}
## [2 ]{.ltx_tag .ltx_tag_section}Related work {#related-work .ltx_title .ltx_title_section}

:::: {#S2.SS0.SSS0.Px1 .section .ltx_paragraph}
#### Multi-agent LLMs. {#multi-agent-llms. .ltx_title .ltx_title_paragraph}

::: {#S2.SS0.SSS0.Px1.p1 .ltx_para}
Multi-agent LLMs has become a popular topic \[[19](#bib.bib19){.ltx_ref}\]. A large proportion of works focus on social simulation. "Generative agents" is a sandbox environment allowing 25 agents to communicate with each other \[[47](#bib.bib47){.ltx_ref}\], while "Social simulacra" extends to 1000 agents \[[48](#bib.bib48){.ltx_ref}\]. Based on the success of them, some works explore the game settings \[[35](#bib.bib35){.ltx_ref}, [68](#bib.bib68){.ltx_ref}, [75](#bib.bib75){.ltx_ref}, [76](#bib.bib76){.ltx_ref}, [42](#bib.bib42){.ltx_ref}\], world wars \[[21](#bib.bib21){.ltx_ref}\], economy markets \[[32](#bib.bib32){.ltx_ref}, [70](#bib.bib70){.ltx_ref}\], recommendation systems \[[81](#bib.bib81){.ltx_ref}\], and pandemics \[[17](#bib.bib17){.ltx_ref}\]. Others advance problem solving, focusing on reasoning of short text via multi-agent debating \[[16](#bib.bib16){.ltx_ref}, [73](#bib.bib73){.ltx_ref}, [10](#bib.bib10){.ltx_ref}, [60](#bib.bib60){.ltx_ref}\] and discussing \[[11](#bib.bib11){.ltx_ref}, [53](#bib.bib53){.ltx_ref}\] for different tasks in reasoning \[[16](#bib.bib16){.ltx_ref}, [60](#bib.bib60){.ltx_ref}\], mechanics problems \[[44](#bib.bib44){.ltx_ref}\], paper review \[[77](#bib.bib77){.ltx_ref}\], knowledge graph construction \[[79](#bib.bib79){.ltx_ref}\], and code intelligence \[[64](#bib.bib64){.ltx_ref}, [22](#bib.bib22){.ltx_ref}\]. Different from the above works, we improve problem solving on long context tasks. To the best of our knowledge, the closest work utilizes a tree structure to do single-hop QA over long context \[[9](#bib.bib9){.ltx_ref}\]. However, it is not designed for multi-hop reasoning or other tasks without communication between sibling agents.
:::
::::

:::: {#S2.SS0.SSS0.Px2 .section .ltx_paragraph}
#### Long Context Modeling for LLMs. {#long-context-modeling-for-llms. .ltx_title .ltx_title_paragraph}

::: {#S2.SS0.SSS0.Px2.p1 .ltx_para}
[Input Reduction]{#S2.SS0.SSS0.Px2.p1.1.1 .ltx_text .ltx_font_italic}: RAG is broadly leveraged to solve long context query-based tasks \[[74](#bib.bib74){.ltx_ref}, [2](#bib.bib2){.ltx_ref}\]. Combined with a strong retriever \[[72](#bib.bib72){.ltx_ref}, [37](#bib.bib37){.ltx_ref}, [67](#bib.bib67){.ltx_ref}\], LLMs are expected to handle long context questions in open domains. Previous studies have augmented LLMs during pretraining \[[24](#bib.bib24){.ltx_ref}, [66](#bib.bib66){.ltx_ref}\], finetuning \[[31](#bib.bib31){.ltx_ref}\], inference \[[80](#bib.bib80){.ltx_ref}\], or directly integrating \[[26](#bib.bib26){.ltx_ref}, [55](#bib.bib55){.ltx_ref}\]. Moreover, some token-level retrieval approaches are proposed \[[34](#bib.bib34){.ltx_ref}\]. Longllmlingua \[[25](#bib.bib25){.ltx_ref}\] removes tokens from long prompt to compress long context prompt to a desired budget. [Window Extension]{#S2.SS0.SSS0.Px2.p1.1.2 .ltx_text .ltx_font_italic}: The context windows of LLMs are getting longer and longer thanks to the development of GPUs. For instance, the context window increases from 1024 (GPT-2 \[[51](#bib.bib51){.ltx_ref}\]), 2048 (GPT-3 \[[7](#bib.bib7){.ltx_ref}\]), to 128k (GPT-4 \[[45](#bib.bib45){.ltx_ref}\]). Moreover, the newest version of Claude-3 \[[5](#bib.bib5){.ltx_ref}\] supports 200k context windows. To save the cost of LLM training, some continue learning or finetuning approaches are proposed to extend the context window of pretrained LLMs \[[41](#bib.bib41){.ltx_ref}, [50](#bib.bib50){.ltx_ref}, [40](#bib.bib40){.ltx_ref}, [43](#bib.bib43){.ltx_ref}\]. For instance, position interpolation \[[13](#bib.bib13){.ltx_ref}\] modifies rotary position encoding \[[58](#bib.bib58){.ltx_ref}\] and extends the context length of LLaMA \[[62](#bib.bib62){.ltx_ref}\] to 32k. Different from the above works, CoA does not reduce the input length or extend the window length of LLMs, but rather leverages multi-agent collaboration and communication to obtain the full receptive field.
:::
::::

:::: {#S2.SS0.SSS0.Px3 .section .ltx_paragraph}
#### Complex Task Reasoning. {#complex-task-reasoning. .ltx_title .ltx_title_paragraph}

::: {#S2.SS0.SSS0.Px3.p1 .ltx_para}
Previous works on complex reasoning have focused on decomposing the complex question into sub-questions to solve them step-by-step. \[[49](#bib.bib49){.ltx_ref}\] decompose the questions with an unsupervised model and answer them separately with another model. Decomposed Prompting \[[27](#bib.bib27){.ltx_ref}\] leverages some predefined modules to classify each decomposed sub-question, then further decompose if needed. Additionally, decomposing is used for human-computer interaction \[[71](#bib.bib71){.ltx_ref}\], and prompter training \[[65](#bib.bib65){.ltx_ref}\]. Recently, many work has been proposed for LLMs, such as Chain-of-thought \[[69](#bib.bib69){.ltx_ref}\] Least-to-most prompting \[[85](#bib.bib85){.ltx_ref}\] and Pearl \[[59](#bib.bib59){.ltx_ref}\]. However, the length of the prompt does not exceed the context limit of a single agent. By contrast, our Chain of Agents framework is proposed to effectively reason across multiple agents to support the unlimited length of source text.
:::

<figure id="S2.SS0.SSS0.Px3.fig1" class="ltx_figure">
<figure id="alg1" class="ltx_float ltx_minipage ltx_align_middle ltx_framed ltx_framed_top" style="width:216.8pt;">
<div id="alg1.5" class="ltx_listing">
<div id="alg1.l0" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span id="alg1.l0.1.1.1" class="ltx_text" style="font-size:80%;">0:</span></span><span id="alg1.l0.2" class="ltx_text" style="font-size:90%;">  Source input </span><span class="math inline"><em>x</em></span><span id="alg1.l0.3" class="ltx_text" style="font-size:90%;">, query </span><span class="math inline"><em>q</em></span><span id="alg1.l0.4" class="ltx_text" style="font-size:90%;">, agent window size </span><span class="math inline"><em>k</em></span><span id="alg1.l0.5" class="ltx_text" style="font-size:90%;">,     large language model </span><span class="math inline">LLM (*)</span><span id="alg1.l0.6" class="ltx_text" style="font-size:90%;">. </span>
</div>
<div id="alg1.l0a" class="ltx_listingline">
<span class="ltx_tag ltx_tag_listingline"><span id="alg1.l0a.1.1.1" class="ltx_text" style="font-size:80%;">0:</span></span><span id="alg1.l0a.2" class="ltx_text" style="font-size:90%;">  Answer to the query. </span>
</div>
<div id="alg1.l1" class="ltx_listingline">
<span id="alg1.l1.1" class="ltx_text" style="font-size:90%;">  Split </span><span class="math inline"><em>x</em></span><span id="alg1.l1.2" class="ltx_text" style="font-size:90%;"> into </span><span class="math inline"><em>l</em></span><span id="alg1.l1.3" class="ltx_text" style="font-size:90%;"> chunks </span><span class="math inline">{<em>c</em><sub>1</sub>, <em>c</em><sub>2</sub>, ⋯, <em>c</em><sub><em>l</em></sub>}</span><span id="alg1.l1.4" class="ltx_text" style="font-size:90%;"> </span>
</div>
<div id="alg1.l2" class="ltx_listingline">
<span id="alg1.l2.1" class="ltx_text" style="font-size:90%;">  where </span><span class="math inline"><em>c</em><sub><em>i</em></sub></span><span id="alg1.l2.2" class="ltx_text" style="font-size:90%;"> is shorter than </span><span class="math inline"><em>k</em></span><span id="alg1.l2.3" class="ltx_text" style="font-size:90%;"> </span>
</div>
<div id="alg1.l3" class="ltx_listingline">
<span id="alg1.l3.1" class="ltx_text" style="font-size:90%;">  Initialize </span><span class="math inline"><em>C</em> <em>U</em><sub>0</sub></span><span id="alg1.l3.2" class="ltx_text" style="font-size:90%;"> </span><span class="math inline"> ← empty string</span><span id="alg1.l3.3" class="ltx_text" style="font-size:90%;">. </span>
</div>
<div id="alg1.l4" class="ltx_listingline">
<span id="alg1.l4.1" class="ltx_text" style="font-size:90%;">  </span><span id="alg1.l4.2" class="ltx_text ltx_font_bold" style="font-size:90%;">for</span><span id="alg1.l4.3" class="ltx_text" style="font-size:90%;"> </span><span class="math inline"><em>i</em></span><span id="alg1.l4.4" class="ltx_text" style="font-size:90%;"> in </span><span class="math inline">1, 2, ⋯, <em>l</em></span><span id="alg1.l4.5" class="ltx_text" style="font-size:90%;"> </span><span id="alg1.l4.6" class="ltx_text ltx_font_bold" style="font-size:90%;">do</span><span id="alg1.l4.7" class="ltx_text" style="font-size:90%;"> </span>
</div>
<div id="alg1.l5" class="ltx_listingline">
<span id="alg1.l5.1" class="ltx_text" style="font-size:90%;">     </span><span class="math inline"><em>C</em> <em>U</em><sub><em>i</em></sub> ← LLM<sub><em>W</em><sub><em>i</em></sub></sub> (<em>I</em><sub><em>W</em></sub>, <em>C</em> <em>U</em><sub><em>i</em> − 1</sub>, <em>c</em><sub><em>i</em></sub>, <em>q</em>)</span><span id="alg1.l5.2" class="ltx_text" style="font-size:90%;"> </span>
</div>
<div id="alg1.l6" class="ltx_listingline">
<span id="alg1.l6.1" class="ltx_text" style="font-size:90%;">  </span><span id="alg1.l6.2" class="ltx_text ltx_font_bold" style="font-size:90%;">end</span><span id="alg1.l6.3" class="ltx_text" style="font-size:90%;"> </span><span id="alg1.l6.4" class="ltx_text ltx_font_bold" style="font-size:90%;">for</span>
</div>
<div id="alg1.l7" class="ltx_listingline">
<span id="alg1.l7.1" class="ltx_text" style="font-size:90%;">  </span><span id="alg1.l7.2" class="ltx_text ltx_font_bold" style="font-size:90%;">return</span><span id="alg1.l7.3" class="ltx_text" style="font-size:90%;">  </span><span class="math inline">LLM<sub><em>M</em></sub> (<em>I</em><sub><em>M</em></sub>, <em>C</em> <em>U</em><sub><em>l</em></sub>, <em>q</em>)</span><span id="alg1.l7.4" class="ltx_text" style="font-size:90%;"> </span>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_float"><span id="alg1.4.1.1" class="ltx_text ltx_font_bold">Algorithm 1</span> </span> Chain of Agents (CoA).</figcaption>
</figure>
</figure>
::::
:::::::::

:::::::::::::: {#S3 .section .ltx_section}
## [3 ]{.ltx_tag .ltx_tag_section}Method {#method .ltx_title .ltx_title_section}

::: {#S3.p1 .ltx_para}
Figure [[1]{.ltx_text .ltx_ref_tag}](#S1.F1 "Figure 1 ‣ 1 Introduction ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref} shows the overview of our Chain-of-Agents (CoA) framework, containing two stages. In stage 1, long context is split into chunks where each chunk can be processed by a worker agent. Then, the worker agents communicate sequentially to produce evidence over the entire context. In stage 2, a manager agent consumes the knowledge from the chain of workers to generate the final answer.
:::

::: {#S3.p2 .ltx_para}
To formulate the task, we denote a long-context sample as $(x,y,q)$, where $x$ is the input of $n$ tokens, $y$ is the output of $m$ tokens, $q$ is an optional query. Given a LLM with $k$ tokens (usually $k \ll n$) as the context window limit, the target is to generate $y$ with the limited input context window. Therefore, we divide each source text $x$ into chunks $x = {\{ c_{1},{c_{2}\hspace{0pt}\ldots\hspace{0pt}c_{l}}\}}$, so that each chunk can be completely fed into the LLM agent backbone model.
:::

::::::: {#S3.SS1 .section .ltx_subsection}
### [3.1 ]{.ltx_tag .ltx_tag_subsection}Stage 1: Worker Agent: Segment Comprehension and Chain-Communication {#stage-1-worker-agent-segment-comprehension-and-chain-communication .ltx_title .ltx_title_subsection}

::: {#S3.SS1.p1 .ltx_para}
In Stage 1, CoA contains a sequence of $l$ number of worker agents. Each worker $W_{i}$ inputs the concatenation of a chunk $c_{i}$ from source text $x$, a query $q$, instruction for a specific task for worker agent $I_{W}$, and the message passed from the previous agent, denoted as "communication unit" $C\hspace{0pt}U_{i - 1}$. The worker agents process them and output the message $C\hspace{0pt}U_{i}$ for next worker, expressed as:

  -- ---------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
     $${{C\hspace{0pt}U_{i}} = {\text{LLM}_{W_{i}}\hspace{0pt}{(I_{W},{C\hspace{0pt}U_{i - 1}},c_{i},q)}}},$$      [(1)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
  -- ---------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
:::

::: {#S3.SS1.p2 .ltx_para}
CUs produced by worker agents vary across different tasks. For question answering, CU contains the evidence for the manager to answer the question. For summarization, CU contains the summary of the previous texts. For code completion, CU contains the code summary with function/class names and explanation. Effectiveness on diverse tasks demonstrates the flexibility of CoA (Appendix [[C]{.ltx_text .ltx_ref_tag}](#A3 "Appendix C Example Outputs ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}).
:::

::: {#S3.SS1.p3 .ltx_para}
The multi-step worker communication in CoA expands the model context to the full receptive field, meaning that the last worker can read the full input no matter how long the input is. Therefore, CoA is extensible to inputs with different lengths by adjusting the number of worker agents.
:::

::: {#S3.SS1.p4 .ltx_para}
The left side of Figure [[1]{.ltx_text .ltx_ref_tag}](#S1.F1 "Figure 1 ‣ 1 Introduction ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref} underscores the necessity of collaborative communication among workers to effectively address complex, long-context reasoning tasks. We observe that 1) Although the question is unanswerable given $c_{1}$, $W_{1}$ generates related evidence that is useful for answering the question; 2) with the partial answer from the previous worker, $W_{2}$ further reasons with the current source to complete the full reasoning chain across agents and generate the interpretative reasoning chain; 3) $W_{3}$ finds no related information in the chunk 3, it directly rewrites $C\hspace{0pt}U_{2}$ by putting the correct answer as the first token of $C\hspace{0pt}U_{3}$ without adding any unrelated information. This shows that if workers are independent (such as tree structure communication), it is impossible to answer hop two while the answer of hop one is held by another worker.
:::
:::::::

:::: {#S3.SS2 .section .ltx_subsection}
### [3.2 ]{.ltx_tag .ltx_tag_subsection}Stage 2: Manager Agent: Information Integration and Response Generation {#stage-2-manager-agent-information-integration-and-response-generation .ltx_title .ltx_title_subsection}

::: {#S3.SS2.p1 .ltx_para}
In Stage 2, after multiple steps of information extraction and comprehension by worker agents, the manager agent produces the final solution. While worker agents extract relevant information in a long-context source, the manager agent synthesizes relevant information accumulated by the end of ["worker-agent-chain"]{#S3.SS2.p1.4.1 .ltx_text .ltx_font_italic} to generate the final answer. Specifically, given the instruction for manager $I_{M}$ and query $q$, the manager agent consumes accumulated knowledge from last worker $C\hspace{0pt}U_{l}$ and generates the final answer $R\hspace{0pt}e\hspace{0pt}s\hspace{0pt}p\hspace{0pt}o\hspace{0pt}n\hspace{0pt}s\hspace{0pt}e$:

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
     $${R\hspace{0pt}e\hspace{0pt}s\hspace{0pt}p\hspace{0pt}o\hspace{0pt}n\hspace{0pt}s\hspace{0pt}e} = {\text{LLM}_{M}\hspace{0pt}{(I_{M},{C\hspace{0pt}U_{l}},q)}}$$      [(2)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------

The benefit of using a separate LLM as the manager agent is to decompose the duty of analyzing chunks in the long-context source (["worker agents"]{#S3.SS2.p1.5.1 .ltx_text .ltx_font_italic}) and producing the final answer (["manager agent"]{#S3.SS2.p1.5.2 .ltx_text .ltx_font_italic}), so that every agent can fulfill its duty to the most[^1^[[^1^[1]{.ltx_tag .ltx_tag_note}Other design choices: Our experiments show that using the last worker $W_{l}$ to directly generate the final result leads to a performance drop. Besides, feeding the manager with all $C\hspace{0pt}U_{i}$ or some $C\hspace{0pt}U$ that is related to the answer (decided by $W_{i}$) also hurts the performance because of confusion led by conflicting $C\hspace{0pt}U_{i}$.]{.ltx_note_content}]{.ltx_note_outer}]{#footnote1 .ltx_note .ltx_role_footnote}.
:::
::::

:::: {#S3.SS3 .section .ltx_subsection}
### [3.3 ]{.ltx_tag .ltx_tag_subsection}Time Complexity Analysis {#time-complexity-analysis .ltx_title .ltx_title_subsection}

<figure id="S3.T2" class="ltx_table ltx_align_floatright">
<table id="S3.T2.4" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="S3.T2.4.5" class="ltx_tr">
<td id="S3.T2.4.5.1" class="ltx_td ltx_border_tt"></td>
<td id="S3.T2.4.5.2" class="ltx_td ltx_align_left ltx_border_tt">Encode</td>
<td id="S3.T2.4.5.3" class="ltx_td ltx_nopad_r ltx_align_left ltx_border_tt">Decode</td>
</tr>
<tr id="S3.T2.2.2" class="ltx_tr">
<td id="S3.T2.2.2.3" class="ltx_td ltx_align_left ltx_border_t">Full-Context</td>
<td id="S3.T2.1.1.1" class="ltx_td ltx_align_left ltx_border_t"><span class="math inline">𝒪 (<em>n</em><sup>2</sup>)</span></td>
<td id="S3.T2.2.2.2" class="ltx_td ltx_nopad_r ltx_align_left ltx_border_t"><span class="math inline">𝒪 (<em>n</em> <em>r</em>)</span></td>
</tr>
<tr id="S3.T2.4.4" class="ltx_tr">
<td id="S3.T2.4.4.3" class="ltx_td ltx_align_left ltx_border_bb">CoA</td>
<td id="S3.T2.3.3.1" class="ltx_td ltx_align_left ltx_border_bb"><span class="math inline">𝒪 (<em>n</em> <em>k</em>)</span></td>
<td id="S3.T2.4.4.2" class="ltx_td ltx_nopad_r ltx_align_left ltx_border_bb"><span class="math inline">𝒪 (<em>n</em> <em>r</em>)</span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table"><span id="S3.T2.6.1.1" class="ltx_text" style="font-size:90%;">Table 2</span>: </span><span id="S3.T2.7.2" class="ltx_text" style="font-size:90%;">Time complexity.</span></figcaption>
</figure>

::: {#S3.SS3.p1 .ltx_para}
We compare the time cost of full-context input and Chain-of-Agents theoretically in a decoder-only setting. We assume the response generated by LLMs contains $r$ tokens on average, the input has $n$ tokens, and the context limit of LLM is $k$. The time complexity is shown in Table [[2]{.ltx_text .ltx_ref_tag}](#S3.T2 "Table 2 ‣ 3.3 Time Complexity Analysis ‣ 3 Method ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref} (Appendix [[A]{.ltx_text .ltx_ref_tag}](#A1 "Appendix A Proof of Time Complexity ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}). As can be seen, the encoding time of CoA is less than Full-Context because $k \ll n$ in long context tasks, while they have the same decoding time. This demonstrates the efficiency of CoA compared with the Full-Context baseline.
:::
::::
::::::::::::::

:::::::::::::::::::::: {#S4 .section .ltx_section}
## [4 ]{.ltx_tag .ltx_tag_section}Experiment {#experiment .ltx_title .ltx_title_section}

:::::::::::: {#S4.SS1 .section .ltx_subsection}
### [4.1 ]{.ltx_tag .ltx_tag_subsection}Experiment Setup {#experiment-setup .ltx_title .ltx_title_subsection}

<figure id="S4.T3" class="ltx_table">
<div id="S4.T3.5" class="ltx_inline-block ltx_transformed_outer" style="width:433.6pt;height:66pt;vertical-align:-0.0pt;">
<span class="ltx_transformed_inner" style="transform:translate(-78.9pt,12.0pt) scale(0.733170118944337,0.733170118944337) ;"> </span>
<table id="S4.T3.5.1" class="ltx_tabular ltx_align_middle">
<tbody>
<tr id="S4.T3.5.1.1" class="ltx_tr">
<td id="S4.T3.5.1.1.1" class="ltx_td ltx_border_tt"></td>
<td colspan="5" id="S4.T3.5.1.1.2" class="ltx_td ltx_align_center ltx_border_tt">Question Answering</td>
<td colspan="3" id="S4.T3.5.1.1.3" class="ltx_td ltx_align_center ltx_border_tt">Summarization</td>
<td id="S4.T3.5.1.1.4" class="ltx_td ltx_align_center ltx_border_tt">Code</td>
</tr>
<tr id="S4.T3.5.1.2" class="ltx_tr">
<td id="S4.T3.5.1.2.1" class="ltx_td"></td>
<td id="S4.T3.5.1.2.2" class="ltx_td ltx_align_center ltx_border_t">HotpotQA</td>
<td id="S4.T3.5.1.2.3" class="ltx_td ltx_align_center ltx_border_t">MuSiQue</td>
<td id="S4.T3.5.1.2.4" class="ltx_td ltx_align_center ltx_border_t">NarrativeQA</td>
<td id="S4.T3.5.1.2.5" class="ltx_td ltx_align_center ltx_border_t">Qasper</td>
<td id="S4.T3.5.1.2.6" class="ltx_td ltx_align_center ltx_border_t">QuALITY</td>
<td id="S4.T3.5.1.2.7" class="ltx_td ltx_align_center ltx_border_t">QMSum</td>
<td id="S4.T3.5.1.2.8" class="ltx_td ltx_align_center ltx_border_t">GovReport</td>
<td id="S4.T3.5.1.2.9" class="ltx_td ltx_align_center ltx_border_t">BookSum</td>
<td id="S4.T3.5.1.2.10" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_t">RepoBench-P</td>
</tr>
<tr id="S4.T3.5.1.3" class="ltx_tr">
<td id="S4.T3.5.1.3.1" class="ltx_td ltx_align_left ltx_border_t">Avg. Input</td>
<td id="S4.T3.5.1.3.2" class="ltx_td ltx_align_center ltx_border_t">10603</td>
<td id="S4.T3.5.1.3.3" class="ltx_td ltx_align_center ltx_border_t">12975</td>
<td id="S4.T3.5.1.3.4" class="ltx_td ltx_align_center ltx_border_t">71787</td>
<td id="S4.T3.5.1.3.5" class="ltx_td ltx_align_center ltx_border_t">4236</td>
<td id="S4.T3.5.1.3.6" class="ltx_td ltx_align_center ltx_border_t">4936</td>
<td id="S4.T3.5.1.3.7" class="ltx_td ltx_align_center ltx_border_t">12524</td>
<td id="S4.T3.5.1.3.8" class="ltx_td ltx_align_center ltx_border_t">9239</td>
<td id="S4.T3.5.1.3.9" class="ltx_td ltx_align_center ltx_border_t">108478</td>
<td id="S4.T3.5.1.3.10" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_t">7105</td>
</tr>
<tr id="S4.T3.5.1.4" class="ltx_tr">
<td id="S4.T3.5.1.4.1" class="ltx_td ltx_align_left">Avg. Agents</td>
<td id="S4.T3.5.1.4.2" class="ltx_td ltx_align_center">2.35</td>
<td id="S4.T3.5.1.4.3" class="ltx_td ltx_align_center">2.88</td>
<td id="S4.T3.5.1.4.4" class="ltx_td ltx_align_center">12.45</td>
<td id="S4.T3.5.1.4.5" class="ltx_td ltx_align_center">1.12</td>
<td id="S4.T3.5.1.4.6" class="ltx_td ltx_align_center">1.31</td>
<td id="S4.T3.5.1.4.7" class="ltx_td ltx_align_center">2.57</td>
<td id="S4.T3.5.1.4.8" class="ltx_td ltx_align_center">2.03</td>
<td id="S4.T3.5.1.4.9" class="ltx_td ltx_align_center">18.63</td>
<td id="S4.T3.5.1.4.10" class="ltx_td ltx_nopad_r ltx_align_center">1.69</td>
</tr>
<tr id="S4.T3.5.1.5" class="ltx_tr">
<td id="S4.T3.5.1.5.1" class="ltx_td ltx_align_left ltx_border_bb">Query-based</td>
<td id="S4.T3.5.1.5.2" class="ltx_td ltx_align_center ltx_border_bb">✓</td>
<td id="S4.T3.5.1.5.3" class="ltx_td ltx_align_center ltx_border_bb">✓</td>
<td id="S4.T3.5.1.5.4" class="ltx_td ltx_align_center ltx_border_bb">✓</td>
<td id="S4.T3.5.1.5.5" class="ltx_td ltx_align_center ltx_border_bb">✓</td>
<td id="S4.T3.5.1.5.6" class="ltx_td ltx_align_center ltx_border_bb">✓</td>
<td id="S4.T3.5.1.5.7" class="ltx_td ltx_align_center ltx_border_bb">✓</td>
<td id="S4.T3.5.1.5.8" class="ltx_td ltx_align_center ltx_border_bb">✗</td>
<td id="S4.T3.5.1.5.9" class="ltx_td ltx_align_center ltx_border_bb">✗</td>
<td id="S4.T3.5.1.5.10" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_bb">✓</td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span id="S4.T3.3.1.1" class="ltx_text" style="font-size:90%;">Table 3</span>: </span><span id="S4.T3.4.2" class="ltx_text ltx_font_bold" style="font-size:90%;">Dataset Statistics<span id="S4.T3.4.2.1" class="ltx_text ltx_font_medium">. Avg. Input/Agents is the average words/agents (8k) for source input.</span></span></figcaption>
</figure>

:::: {#S4.SS1.SSS0.Px1 .section .ltx_paragraph}
#### Datasets. {#datasets. .ltx_title .ltx_title_paragraph}

::: {#S4.SS1.SSS0.Px1.p1 .ltx_para}
We conduct experiments on nine long context datasets across three task types (Table [[3]{.ltx_text .ltx_ref_tag}](#S4.T3 "Table 3 ‣ 4.1 Experiment Setup ‣ 4 Experiment ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}):

- [[•]{.ltx_tag .ltx_tag_item}]{#S4.I1.i1}

  ::: {#S4.I1.i1.p1 .ltx_para}
  [Question Answering.]{#S4.I1.i1.p1.1.1 .ltx_text .ltx_font_bold} We consider five QA datasets from the LongBench \[[6](#bib.bib6){.ltx_ref}\] and SCROLL \[[54](#bib.bib54){.ltx_ref}\]. [HotpotQA]{#S4.I1.i1.p1.1.2 .ltx_text .ltx_font_bold} \[[78](#bib.bib78){.ltx_ref}\] is a Wikipedia-based multi-hop QA dataset. It requires reasoning across multiple passages to find the answer. [MuSiQue]{#S4.I1.i1.p1.1.3 .ltx_text .ltx_font_bold} \[[63](#bib.bib63){.ltx_ref}\] is a multi-hop QA dataset. It is much more difficult than HotpotQA as it contains more hops in one sample, unanswerable questions, and harder distracting content. [NarrativeQA]{#S4.I1.i1.p1.1.4 .ltx_text .ltx_font_bold} \[[28](#bib.bib28){.ltx_ref}\] is a QA dataset over entire books or movie transcripts. The answers can be abstract or extractive, yes/no, and unanswerable. [Qasper]{#S4.I1.i1.p1.1.5 .ltx_text .ltx_font_bold} \[[15](#bib.bib15){.ltx_ref}\] is a question answering dataset over NLP papers. It also contains extractive, abstractive, yes/no, and unanswerable questions. [QuALITY]{#S4.I1.i1.p1.1.6 .ltx_text .ltx_font_bold} \[[46](#bib.bib46){.ltx_ref}\] is a dataset based on stories and articles with multiple-choice questions for each sample. The model needs to select the correct answer among choices.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S4.I1.i2}

  ::: {#S4.I1.i2.p1 .ltx_para}
  [Summarization.]{#S4.I1.i2.p1.1.1 .ltx_text .ltx_font_bold} We pick two summarization datasets from SCROLLS. [QMSum]{#S4.I1.i2.p1.1.2 .ltx_text .ltx_font_bold} \[[84](#bib.bib84){.ltx_ref}\] is a query-based summarization dataset, formed by meeting transcripts from multiple domains such as academic and industrial products. [GovReport]{#S4.I1.i2.p1.1.3 .ltx_text .ltx_font_bold} \[[23](#bib.bib23){.ltx_ref}\] is a generic summarization dataset containing long reports published by the U.S. Government Accountability Office. We also use one dataset for long context memorization tasks. [BookSum]{#S4.I1.i2.p1.1.4 .ltx_text .ltx_font_bold} \[[30](#bib.bib30){.ltx_ref}\] is a collection of datasets for long-form narrative summarization, including novels, plays, and stories. We use the book-level partition of the BookSum dataset for experiments.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S4.I1.i3}

  ::: {#S4.I1.i3.p1 .ltx_para}
  [Code Completion.]{#S4.I1.i3.p1.1.1 .ltx_text .ltx_font_bold} We pick one code completion dataset from LongBench. [RepoBench-P]{#S4.I1.i3.p1.1.2 .ltx_text .ltx_font_bold} \[[39](#bib.bib39){.ltx_ref}\] is collected from GitHub repositories, and the model needs to generate the next line of code given the long code base.
  :::
:::
::::

:::: {#S4.SS1.SSS0.Px2 .section .ltx_paragraph}
#### Metrics. {#metrics. .ltx_title .ltx_title_paragraph}

::: {#S4.SS1.SSS0.Px2.p1 .ltx_para}
We report the geometric mean of ROUGE \[[36](#bib.bib36){.ltx_ref}\] for Summarization tasks, code similarity score \[[6](#bib.bib6){.ltx_ref}\] for Code Completion task, exact match for QuALITY dataset \[[54](#bib.bib54){.ltx_ref}\], and F1 score for the rest of the Question Answering datasets \[[6](#bib.bib6){.ltx_ref}\].
:::
::::

:::: {#S4.SS1.SSS0.Px3 .section .ltx_paragraph}
#### LLMs. {#llms. .ltx_title .ltx_title_paragraph}

::: {#S4.SS1.SSS0.Px3.p1 .ltx_para}
We use six LLMs in total as the backbone of CoA across all experiments. [PaLM 2]{#S4.SS1.SSS0.Px3.p1.2.1 .ltx_text .ltx_font_bold} \[[4](#bib.bib4){.ltx_ref}\] is a series of models with a dense left-to-right, decoder-only language model pretrained on a high-quality corpus of 780 billion tokens. We use [text-bison@001]{#S4.SS1.SSS0.Px3.p1.2.2 .ltx_text .ltx_font_bold} and [text-unicorn@001]{#S4.SS1.SSS0.Px3.p1.2.3 .ltx_text .ltx_font_bold} for the experiments with an 8k maximum context window. [Gemini 1.0]{#S4.SS1.SSS0.Px3.p1.2.4 .ltx_text .ltx_font_bold} \[[61](#bib.bib61){.ltx_ref}\] is a family of LLMs proposed by Google. We use [gemini-ultra]{#S4.SS1.SSS0.Px3.p1.2.5 .ltx_text .ltx_font_bold} for experiments. The input limit is 32k tokens for Gemini. [Claude 3]{#S4.SS1.SSS0.Px3.p1.2.6 .ltx_text .ltx_font_bold} \[[5](#bib.bib5){.ltx_ref}\] is a family of large language models developed by Anthropic. The family includes three state-of-the-art models in ascending order of capability: [claude-3-haiku]{#S4.SS1.SSS0.Px3.p1.2.7 .ltx_text .ltx_font_bold}, [claude-3-sonnet]{#S4.SS1.SSS0.Px3.p1.2.8 .ltx_text .ltx_font_bold}, and [claude-3-opus]{#S4.SS1.SSS0.Px3.p1.2.9 .ltx_text .ltx_font_bold}. These models are capable of consuming 200k tokens in the context window, providing a strong baseline for long context tasks. Although our framework is flexible to use diverse types of LLMs as workers and manager, we use the same model for each $W_{i}$ and $M$ if not specified.
:::
::::

::::: {#S4.SS1.SSS0.Px4 .section .ltx_paragraph}
#### Baselines. {#baselines. .ltx_title .ltx_title_paragraph}

::: {#S4.SS1.SSS0.Px4.p1 .ltx_para}
Our principle of choosing baselines is to find the strongest and most typical approaches from [input reduction]{#S4.SS1.SSS0.Px4.p1.1.1 .ltx_text .ltx_font_italic} and [window extension]{#S4.SS1.SSS0.Px4.p1.1.2 .ltx_text .ltx_font_italic}. The first baseline is [Vanilla]{#S4.SS1.SSS0.Px4.p1.1.3 .ltx_text .ltx_font_bold}. It directly consumes tokens until the context window of LLM is fully utilized, implying a 200k window LLM if using Claude 3. The other one is Retrieval-Augmented Generation ([RAG]{#S4.SS1.SSS0.Px4.p1.1.4 .ltx_text .ltx_font_bold}). We use the state-of-the-art retriever \[[72](#bib.bib72){.ltx_ref}\]. Following \[[74](#bib.bib74){.ltx_ref}\], we first segment the source text into chunks of 300 words and re-rank them using a retriever. Top-n chunks are then fed into the downstream LLM until the context window is fully utilized. GovReport dataset does not contain a query initially, we create a pseudo query "What is the summary of the whole government report?" as the query to rerank.
:::

::: {#S4.SS1.SSS0.Px4.p2 .ltx_para}
To evaluate the performance of CoA compared with possible multi-agent approaches, we carefully construct two multi-agent approaches. For these two approaches, similar to CoA, we also assign each chunk $c_{i}$ to $W_{i}$ using similar instructions to generate $C\hspace{0pt}U_{i}$. In these approaches, worker agents are parallel and independent while CoA is sequential. [Multi-Agent Voting]{#S4.SS1.SSS0.Px4.p2.13.1 .ltx_text .ltx_font_bold} ([Merge]{#S4.SS1.SSS0.Px4.p2.13.2 .ltx_text .ltx_font_bold}) Each agent directly generate an answer $a_{i}$ according to $c_{i}$. A majority voting is applied to all answers $a_{i}$ to decide the final answer. [Multi-Agent Hierarchical Structure]{#S4.SS1.SSS0.Px4.p2.13.3 .ltx_text .ltx_font_bold} ([Hierarchical]{#S4.SS1.SSS0.Px4.p2.13.4 .ltx_text .ltx_font_bold}). Inspired by \[[9](#bib.bib9){.ltx_ref}\], we propose a hierarchical framework, where the communication forms a tree structure between workers $W_{i}$ and manager $M$. For each worker, it first judges whether $c_{i}$ contains useful information. If true, it generates a communication unit $C\hspace{0pt}U_{i}$. Then, all $C\hspace{0pt}U_{i}$ are sent to the manager $M$ to come up with a final answer. Besides, we append an integer number $L$ at the end of every approach to clearly remind the window size limit of LLM. For instance, "CoA (8K)" refers to the base LLM used in CoA with window size 8K.
:::
:::::
::::::::::::

::::::::::: {#S4.SS2 .section .ltx_subsection}
### [4.2 ]{.ltx_tag .ltx_tag_subsection}Overall Results of CoA {#overall-results-of-coa .ltx_title .ltx_title_subsection}

<figure id="S4.T4" class="ltx_table">
<div id="S4.T4.5" class="ltx_inline-block ltx_transformed_outer" style="width:433.6pt;height:148.2pt;vertical-align:-0.0pt;">
<span class="ltx_transformed_inner" style="transform:translate(-99.1pt,33.9pt) scale(0.686332698525283,0.686332698525283) ;"> </span>
<table id="S4.T4.5.1" class="ltx_tabular ltx_align_middle">
<tbody>
<tr id="S4.T4.5.1.1" class="ltx_tr">
<td id="S4.T4.5.1.1.1" class="ltx_td ltx_border_tt"></td>
<td id="S4.T4.5.1.1.2" class="ltx_td ltx_border_tt"></td>
<td colspan="5" id="S4.T4.5.1.1.3" class="ltx_td ltx_align_center ltx_border_tt">Question Answering</td>
<td colspan="2" id="S4.T4.5.1.1.4" class="ltx_td ltx_align_center ltx_border_tt">Summarization</td>
<td id="S4.T4.5.1.1.5" class="ltx_td ltx_align_center ltx_border_tt">Code</td>
</tr>
<tr id="S4.T4.5.1.2" class="ltx_tr">
<td id="S4.T4.5.1.2.1" class="ltx_td ltx_align_left">LLMs</td>
<td id="S4.T4.5.1.2.2" class="ltx_td ltx_align_left">Baselines</td>
<td id="S4.T4.5.1.2.3" class="ltx_td ltx_align_left ltx_border_t">HotpotQA</td>
<td id="S4.T4.5.1.2.4" class="ltx_td ltx_align_left ltx_border_t">MuSiQue</td>
<td id="S4.T4.5.1.2.5" class="ltx_td ltx_align_left ltx_border_t">NarrativeQA</td>
<td id="S4.T4.5.1.2.6" class="ltx_td ltx_align_left ltx_border_t">Qasper</td>
<td id="S4.T4.5.1.2.7" class="ltx_td ltx_align_left ltx_border_t">QuALITY</td>
<td id="S4.T4.5.1.2.8" class="ltx_td ltx_align_left ltx_border_t">QMSum</td>
<td id="S4.T4.5.1.2.9" class="ltx_td ltx_align_left ltx_border_t">GovReport</td>
<td id="S4.T4.5.1.2.10" class="ltx_td ltx_align_left ltx_border_t">RepoBench-P</td>
</tr>
<tr id="S4.T4.5.1.3" class="ltx_tr">
<td rowspan="3" id="S4.T4.5.1.3.1" class="ltx_td ltx_align_center ltx_border_t"><span id="S4.T4.5.1.3.1.1" class="ltx_text">text-bison</span></td>
<td id="S4.T4.5.1.3.2" class="ltx_td ltx_align_left ltx_border_t">Vanilla (8k)</td>
<td id="S4.T4.5.1.3.3" class="ltx_td ltx_align_center ltx_border_t">45.57</td>
<td id="S4.T4.5.1.3.4" class="ltx_td ltx_align_center ltx_border_t">26.87</td>
<td id="S4.T4.5.1.3.5" class="ltx_td ltx_align_center ltx_border_t">11.96</td>
<td id="S4.T4.5.1.3.6" class="ltx_td ltx_align_center ltx_border_t">26.56</td>
<td id="S4.T4.5.1.3.7" class="ltx_td ltx_align_center ltx_border_t">61.86</td>
<td id="S4.T4.5.1.3.8" class="ltx_td ltx_align_center ltx_border_t">15.45</td>
<td id="S4.T4.5.1.3.9" class="ltx_td ltx_align_center ltx_border_t">20.60</td>
<td id="S4.T4.5.1.3.10" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_t">56.30</td>
</tr>
<tr id="S4.T4.5.1.4" class="ltx_tr">
<td id="S4.T4.5.1.4.1" class="ltx_td ltx_align_left">RAG (8k)</td>
<td id="S4.T4.5.1.4.2" class="ltx_td ltx_align_center">51.91</td>
<td id="S4.T4.5.1.4.3" class="ltx_td ltx_align_center">33.83</td>
<td id="S4.T4.5.1.4.4" class="ltx_td ltx_align_center">14.20</td>
<td id="S4.T4.5.1.4.5" class="ltx_td ltx_align_center">27.20</td>
<td id="S4.T4.5.1.4.6" class="ltx_td ltx_align_center">55.28</td>
<td id="S4.T4.5.1.4.7" class="ltx_td ltx_align_center">15.59</td>
<td id="S4.T4.5.1.4.8" class="ltx_td ltx_align_center">20.83</td>
<td id="S4.T4.5.1.4.9" class="ltx_td ltx_nopad_r ltx_align_center">55.63</td>
</tr>
<tr id="S4.T4.5.1.5" class="ltx_tr">
<td id="S4.T4.5.1.5.1" class="ltx_td ltx_align_left">CoA (8k)</td>
<td id="S4.T4.5.1.5.2" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.5.2.1" class="ltx_text ltx_font_bold">53.62</span></td>
<td id="S4.T4.5.1.5.3" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.5.3.1" class="ltx_text ltx_font_bold">37.09</span></td>
<td id="S4.T4.5.1.5.4" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.5.4.1" class="ltx_text ltx_font_bold">25.26</span></td>
<td id="S4.T4.5.1.5.5" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.5.5.1" class="ltx_text ltx_font_bold">37.17</span></td>
<td id="S4.T4.5.1.5.6" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.5.6.1" class="ltx_text ltx_font_bold">65.42</span></td>
<td id="S4.T4.5.1.5.7" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.5.7.1" class="ltx_text ltx_font_bold">16.77</span></td>
<td id="S4.T4.5.1.5.8" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.5.8.1" class="ltx_text ltx_font_bold">26.11</span></td>
<td id="S4.T4.5.1.5.9" class="ltx_td ltx_nopad_r ltx_align_center"><span id="S4.T4.5.1.5.9.1" class="ltx_text ltx_font_bold">58.25</span></td>
</tr>
<tr id="S4.T4.5.1.6" class="ltx_tr">
<td rowspan="3" id="S4.T4.5.1.6.1" class="ltx_td ltx_align_center ltx_border_t"><span id="S4.T4.5.1.6.1.1" class="ltx_text">text-unicorn</span></td>
<td id="S4.T4.5.1.6.2" class="ltx_td ltx_align_left ltx_border_t">Vanilla (8k)</td>
<td id="S4.T4.5.1.6.3" class="ltx_td ltx_align_center ltx_border_t">51.09</td>
<td id="S4.T4.5.1.6.4" class="ltx_td ltx_align_center ltx_border_t">29.67</td>
<td id="S4.T4.5.1.6.5" class="ltx_td ltx_align_center ltx_border_t">14.45</td>
<td id="S4.T4.5.1.6.6" class="ltx_td ltx_align_center ltx_border_t">27.81</td>
<td id="S4.T4.5.1.6.7" class="ltx_td ltx_align_center ltx_border_t">83.40</td>
<td id="S4.T4.5.1.6.8" class="ltx_td ltx_align_center ltx_border_t">16.61</td>
<td id="S4.T4.5.1.6.9" class="ltx_td ltx_align_center ltx_border_t">23.50</td>
<td id="S4.T4.5.1.6.10" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_t">53.87</td>
</tr>
<tr id="S4.T4.5.1.7" class="ltx_tr">
<td id="S4.T4.5.1.7.1" class="ltx_td ltx_align_left">RAG (8k)</td>
<td id="S4.T4.5.1.7.2" class="ltx_td ltx_align_center">58.01</td>
<td id="S4.T4.5.1.7.3" class="ltx_td ltx_align_center">40.38</td>
<td id="S4.T4.5.1.7.4" class="ltx_td ltx_align_center">19.12</td>
<td id="S4.T4.5.1.7.5" class="ltx_td ltx_align_center">24.44</td>
<td id="S4.T4.5.1.7.6" class="ltx_td ltx_align_center">83.00</td>
<td id="S4.T4.5.1.7.7" class="ltx_td ltx_align_center">16.83</td>
<td id="S4.T4.5.1.7.8" class="ltx_td ltx_align_center">21.43</td>
<td id="S4.T4.5.1.7.9" class="ltx_td ltx_nopad_r ltx_align_center">50.49</td>
</tr>
<tr id="S4.T4.5.1.8" class="ltx_tr">
<td id="S4.T4.5.1.8.1" class="ltx_td ltx_align_left">CoA (8k)</td>
<td id="S4.T4.5.1.8.2" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.8.2.1" class="ltx_text ltx_font_bold">62.04</span></td>
<td id="S4.T4.5.1.8.3" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.8.3.1" class="ltx_text ltx_font_bold">42.49</span></td>
<td id="S4.T4.5.1.8.4" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.8.4.1" class="ltx_text ltx_font_bold">20.37</span></td>
<td id="S4.T4.5.1.8.5" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.8.5.1" class="ltx_text ltx_font_bold">38.01</span></td>
<td id="S4.T4.5.1.8.6" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.8.6.1" class="ltx_text ltx_font_bold">83.80</span></td>
<td id="S4.T4.5.1.8.7" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.8.7.1" class="ltx_text ltx_font_bold">17.67</span></td>
<td id="S4.T4.5.1.8.8" class="ltx_td ltx_align_center"><span id="S4.T4.5.1.8.8.1" class="ltx_text ltx_font_bold">26.48</span></td>
<td id="S4.T4.5.1.8.9" class="ltx_td ltx_nopad_r ltx_align_center"><span id="S4.T4.5.1.8.9.1" class="ltx_text ltx_font_bold">60.39</span></td>
</tr>
<tr id="S4.T4.5.1.9" class="ltx_tr">
<td rowspan="4" id="S4.T4.5.1.9.1" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t"><span id="S4.T4.5.1.9.1.1" class="ltx_text">gemini-ultra</span></td>
<td id="S4.T4.5.1.9.2" class="ltx_td ltx_align_left ltx_border_t">Vanilla (8k)</td>
<td id="S4.T4.5.1.9.3" class="ltx_td ltx_align_center ltx_border_t">40.62</td>
<td id="S4.T4.5.1.9.4" class="ltx_td ltx_align_center ltx_border_t">23.61</td>
<td id="S4.T4.5.1.9.5" class="ltx_td ltx_align_center ltx_border_t">7.71</td>
<td id="S4.T4.5.1.9.6" class="ltx_td ltx_align_center ltx_border_t">20.59</td>
<td id="S4.T4.5.1.9.7" class="ltx_td ltx_align_center ltx_border_t">57.40</td>
<td id="S4.T4.5.1.9.8" class="ltx_td ltx_align_center ltx_border_t">12.10</td>
<td id="S4.T4.5.1.9.9" class="ltx_td ltx_align_center ltx_border_t">26.18</td>
<td id="S4.T4.5.1.9.10" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_t">49.09</td>
</tr>
<tr id="S4.T4.5.1.10" class="ltx_tr">
<td id="S4.T4.5.1.10.1" class="ltx_td ltx_align_left">Vanilla (32k)</td>
<td id="S4.T4.5.1.10.2" class="ltx_td ltx_align_center">45.09</td>
<td id="S4.T4.5.1.10.3" class="ltx_td ltx_align_center">27.93</td>
<td id="S4.T4.5.1.10.4" class="ltx_td ltx_align_center">7.21</td>
<td id="S4.T4.5.1.10.5" class="ltx_td ltx_align_center">21.71</td>
<td id="S4.T4.5.1.10.6" class="ltx_td ltx_align_center">58.60</td>
<td id="S4.T4.5.1.10.7" class="ltx_td ltx_align_center">10.24</td>
<td id="S4.T4.5.1.10.8" class="ltx_td ltx_align_center">26.96</td>
<td id="S4.T4.5.1.10.9" class="ltx_td ltx_nopad_r ltx_align_center">73.04</td>
</tr>
<tr id="S4.T4.5.1.11" class="ltx_tr">
<td id="S4.T4.5.1.11.1" class="ltx_td ltx_align_left">RAG (8k)</td>
<td id="S4.T4.5.1.11.2" class="ltx_td ltx_align_center">51.13</td>
<td id="S4.T4.5.1.11.3" class="ltx_td ltx_align_center">31.56</td>
<td id="S4.T4.5.1.11.4" class="ltx_td ltx_align_center">14.51</td>
<td id="S4.T4.5.1.11.5" class="ltx_td ltx_align_center">18.70</td>
<td id="S4.T4.5.1.11.6" class="ltx_td ltx_align_center">62.40</td>
<td id="S4.T4.5.1.11.7" class="ltx_td ltx_align_center">12.70</td>
<td id="S4.T4.5.1.11.8" class="ltx_td ltx_align_center">25.87</td>
<td id="S4.T4.5.1.11.9" class="ltx_td ltx_nopad_r ltx_align_center">72.94</td>
</tr>
<tr id="S4.T4.5.1.12" class="ltx_tr">
<td id="S4.T4.5.1.12.1" class="ltx_td ltx_align_left ltx_border_bb">CoA (8k)</td>
<td id="S4.T4.5.1.12.2" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T4.5.1.12.2.1" class="ltx_text ltx_font_bold">54.26</span></td>
<td id="S4.T4.5.1.12.3" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T4.5.1.12.3.1" class="ltx_text ltx_font_bold">35.09</span></td>
<td id="S4.T4.5.1.12.4" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T4.5.1.12.4.1" class="ltx_text ltx_font_bold">25.26</span></td>
<td id="S4.T4.5.1.12.5" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T4.5.1.12.5.1" class="ltx_text ltx_font_bold">35.10</span></td>
<td id="S4.T4.5.1.12.6" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T4.5.1.12.6.1" class="ltx_text ltx_font_bold">80.60</span></td>
<td id="S4.T4.5.1.12.7" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T4.5.1.12.7.1" class="ltx_text ltx_font_bold">12.84</span></td>
<td id="S4.T4.5.1.12.8" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T4.5.1.12.8.1" class="ltx_text ltx_font_bold">26.98</span></td>
<td id="S4.T4.5.1.12.9" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_bb"><span id="S4.T4.5.1.12.9.1" class="ltx_text ltx_font_bold">73.05</span></td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span id="S4.T4.3.1.1" class="ltx_text" style="font-size:90%;">Table 4</span>: </span><span id="S4.T4.4.2" class="ltx_text ltx_font_bold" style="font-size:90%;">Overall results of CoA<span id="S4.T4.4.2.1" class="ltx_text ltx_font_medium">. CoA significantly outperforms Vanilla and RAG using various backbone LLMs on all datasets.</span></span></figcaption>
</figure>

:::: {#S4.SS2.SSS0.Px1 .section .ltx_paragraph}
#### Question Answering. {#question-answering. .ltx_title .ltx_title_paragraph}

::: {#S4.SS2.SSS0.Px1.p1 .ltx_para}
Table [[4]{.ltx_text .ltx_ref_tag}](#S4.T4 "Table 4 ‣ 4.2 Overall Results of CoA ‣ 4 Experiment ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref} shows the results of Question Answering tasks on all three models. CoA (8k) outperforms Vanilla (8k) by a large margin on [all 8 datasets]{#S4.SS2.SSS0.Px1.p1.1.1 .ltx_text .ltx_font_italic}, including 13.30% on NarrativeQA, 12.82% on MuSiQue, and 22.00% on Quality, for text-bison, text-unicorn, and gemini-ultra, respectively. Also, CoA (8k) outperforms RAG (8k) model for all 8 datasets using all three LLMs, demonstrating that CoA achieves higher performance than RAG. In other words, [using multi-agent LLMs outperforms RAG models]{#S4.SS2.SSS0.Px1.p1.1.2 .ltx_text .ltx_font_bold}. It is also worth noting that for gemini-ultra, Vanilla (32k) improves the Vanilla (8k) baseline, yet it is still lower than CoA (8k).
:::
::::

:::: {#S4.SS2.SSS0.Px2 .section .ltx_paragraph}
#### Summarization and Code Completion. {#summarization-and-code-completion. .ltx_title .ltx_title_paragraph}

::: {#S4.SS2.SSS0.Px2.p1 .ltx_para}
Table [[4]{.ltx_text .ltx_ref_tag}](#S4.T4 "Table 4 ‣ 4.2 Overall Results of CoA ‣ 4 Experiment ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref} shows the results of Summarization and Code Completion tasks. Similarly, CoA (8k) also outperforms all Vanilla (8k) and (32k) baselines on all three datasets, demonstrating the strong capability of CoA on various tasks. It is worth noting that for GovReport, RAG fails to improve the baseline with pseudo query. By contrast, CoA improves the performance significantly, showing that [CoA can also be applied in non-query tasks]{#S4.SS2.SSS0.Px2.p1.1.1 .ltx_text .ltx_font_bold}.
:::
::::

:::: {#S4.SS2.SSS0.Px3 .section .ltx_paragraph}
#### Long Context LLMs. {#long-context-llms. .ltx_title .ltx_title_paragraph}

::: {#S4.SS2.SSS0.Px3.p1 .ltx_para}
As Claude 3 models support 200k of tokens, we call these models long context models (LCM). Table [[5]{.ltx_text .ltx_ref_tag}](#S4.T5 "Table 5 ‣ Long Context LLMs. ‣ 4.2 Overall Results of CoA ‣ 4 Experiment ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref} shows the performance of the LCM on two datasets. As can be seen, CoA (8k) outperforms Vanilla (200k) significantly, showing that with only an 8k context window, [CoA can achieve a much higher performance than LCM with a 200k context window]{#S4.SS2.SSS0.Px3.p1.1.1 .ltx_text .ltx_font_bold}. Also, CoA improves the performance with the samples that can be fed into a 200k context window (no truncation). Moreover, the improvements over the Vanilla (200k) and RAG (8k) become higher when the model size increases from Haiku to Opus (11.63/11.36/17.4 for NarrativeQA, 1.66/2.86/3.47 for BookSum). This demonstrates that [CoA benefits from stronger models to achieve higher improvements]{#S4.SS2.SSS0.Px3.p1.1.2 .ltx_text .ltx_font_bold}.
:::

<figure id="S4.T5" class="ltx_table">
<div id="S4.T5.4" class="ltx_inline-block ltx_transformed_outer" style="width:433.6pt;height:136.7pt;vertical-align:-0.0pt;">
<span class="ltx_transformed_inner" style="transform:translate(-40.1pt,12.6pt) scale(0.843845626931395,0.843845626931395) ;"> </span>
<table id="S4.T5.4.1" class="ltx_tabular ltx_align_middle">
<tbody>
<tr id="S4.T5.4.1.1" class="ltx_tr">
<td id="S4.T5.4.1.1.1" class="ltx_td ltx_border_tt"></td>
<td colspan="3" id="S4.T5.4.1.1.2" class="ltx_td ltx_align_center ltx_border_tt">claude-3-haiku</td>
<td colspan="3" id="S4.T5.4.1.1.3" class="ltx_td ltx_align_center ltx_border_tt">claude-3-sonnet</td>
<td colspan="3" id="S4.T5.4.1.1.4" class="ltx_td ltx_align_center ltx_border_tt">claude-3-opus</td>
</tr>
<tr id="S4.T5.4.1.2" class="ltx_tr">
<td id="S4.T5.4.1.2.1" class="ltx_td"></td>
<td id="S4.T5.4.1.2.2" class="ltx_td ltx_align_left ltx_border_t">No Trun.</td>
<td id="S4.T5.4.1.2.3" class="ltx_td ltx_align_left ltx_border_t">Trun.</td>
<td id="S4.T5.4.1.2.4" class="ltx_td ltx_align_left ltx_border_t">Average</td>
<td id="S4.T5.4.1.2.5" class="ltx_td ltx_align_left ltx_border_t">No Trun.</td>
<td id="S4.T5.4.1.2.6" class="ltx_td ltx_align_left ltx_border_t">Trun.</td>
<td id="S4.T5.4.1.2.7" class="ltx_td ltx_align_left ltx_border_t">Average</td>
<td id="S4.T5.4.1.2.8" class="ltx_td ltx_align_left ltx_border_t">No Trun.</td>
<td id="S4.T5.4.1.2.9" class="ltx_td ltx_align_left ltx_border_t">Trun.</td>
<td id="S4.T5.4.1.2.10" class="ltx_td ltx_align_left ltx_border_t">Average</td>
</tr>
<tr id="S4.T5.4.1.3" class="ltx_tr">
<td id="S4.T5.4.1.3.1" class="ltx_td ltx_align_left ltx_border_t"><span id="S4.T5.4.1.3.1.1" class="ltx_text ltx_font_bold">NarrativeQA</span></td>
<td id="S4.T5.4.1.3.2" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.3.3" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.3.4" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.3.5" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.3.6" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.3.7" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.3.8" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.3.9" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.3.10" class="ltx_td ltx_border_t"></td>
</tr>
<tr id="S4.T5.4.1.4" class="ltx_tr">
<td id="S4.T5.4.1.4.1" class="ltx_td ltx_align_left">Vanilla (200k)</td>
<td id="S4.T5.4.1.4.2" class="ltx_td ltx_align_center">8.00</td>
<td id="S4.T5.4.1.4.3" class="ltx_td ltx_align_center">2.50</td>
<td id="S4.T5.4.1.4.4" class="ltx_td ltx_align_center">7.17</td>
<td id="S4.T5.4.1.4.5" class="ltx_td ltx_align_center">5.58</td>
<td id="S4.T5.4.1.4.6" class="ltx_td ltx_align_center">2.44</td>
<td id="S4.T5.4.1.4.7" class="ltx_td ltx_align_center">5.15</td>
<td id="S4.T5.4.1.4.8" class="ltx_td ltx_align_center">7.23</td>
<td id="S4.T5.4.1.4.9" class="ltx_td ltx_align_center">2.35</td>
<td id="S4.T5.4.1.4.10" class="ltx_td ltx_nopad_r ltx_align_center">6.56</td>
</tr>
<tr id="S4.T5.4.1.5" class="ltx_tr">
<td id="S4.T5.4.1.5.1" class="ltx_td ltx_align_left">RAG (8k)</td>
<td id="S4.T5.4.1.5.2" class="ltx_td ltx_align_center">5.94</td>
<td id="S4.T5.4.1.5.3" class="ltx_td ltx_align_center">4.22</td>
<td id="S4.T5.4.1.5.4" class="ltx_td ltx_align_center">5.71</td>
<td id="S4.T5.4.1.5.5" class="ltx_td ltx_align_center">9.09</td>
<td id="S4.T5.4.1.5.6" class="ltx_td ltx_align_center">5.17</td>
<td id="S4.T5.4.1.5.7" class="ltx_td ltx_align_center">8.50</td>
<td id="S4.T5.4.1.5.8" class="ltx_td ltx_align_center">6.13</td>
<td id="S4.T5.4.1.5.9" class="ltx_td ltx_align_center">4.29</td>
<td id="S4.T5.4.1.5.10" class="ltx_td ltx_nopad_r ltx_align_center">5.86</td>
</tr>
<tr id="S4.T5.4.1.6" class="ltx_tr">
<td id="S4.T5.4.1.6.1" class="ltx_td ltx_align_left">CoA (8k)</td>
<td id="S4.T5.4.1.6.2" class="ltx_td ltx_align_center"><span id="S4.T5.4.1.6.2.1" class="ltx_text ltx_font_bold">18.31</span></td>
<td id="S4.T5.4.1.6.3" class="ltx_td ltx_align_center"><span id="S4.T5.4.1.6.3.1" class="ltx_text ltx_font_bold">21.34</span></td>
<td id="S4.T5.4.1.6.4" class="ltx_td ltx_align_center"><span id="S4.T5.4.1.6.4.1" class="ltx_text ltx_font_bold">18.80</span></td>
<td id="S4.T5.4.1.6.5" class="ltx_td ltx_align_center"><span id="S4.T5.4.1.6.5.1" class="ltx_text ltx_font_bold">16.63</span></td>
<td id="S4.T5.4.1.6.6" class="ltx_td ltx_align_center"><span id="S4.T5.4.1.6.6.1" class="ltx_text ltx_font_bold">16.47</span></td>
<td id="S4.T5.4.1.6.7" class="ltx_td ltx_align_center"><span id="S4.T5.4.1.6.7.1" class="ltx_text ltx_font_bold">16.51</span></td>
<td id="S4.T5.4.1.6.8" class="ltx_td ltx_align_center"><span id="S4.T5.4.1.6.8.1" class="ltx_text ltx_font_bold">24.38</span></td>
<td id="S4.T5.4.1.6.9" class="ltx_td ltx_align_center"><span id="S4.T5.4.1.6.9.1" class="ltx_text ltx_font_bold">21.26</span></td>
<td id="S4.T5.4.1.6.10" class="ltx_td ltx_nopad_r ltx_align_center"><span id="S4.T5.4.1.6.10.1" class="ltx_text ltx_font_bold">23.96</span></td>
</tr>
<tr id="S4.T5.4.1.7" class="ltx_tr">
<td id="S4.T5.4.1.7.1" class="ltx_td ltx_align_left ltx_border_t"><span id="S4.T5.4.1.7.1.1" class="ltx_text ltx_font_bold">BookSum</span></td>
<td id="S4.T5.4.1.7.2" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.7.3" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.7.4" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.7.5" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.7.6" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.7.7" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.7.8" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.7.9" class="ltx_td ltx_border_t"></td>
<td id="S4.T5.4.1.7.10" class="ltx_td ltx_border_t"></td>
</tr>
<tr id="S4.T5.4.1.8" class="ltx_tr">
<td id="S4.T5.4.1.8.1" class="ltx_td ltx_align_left">Vanilla (200k)</td>
<td id="S4.T5.4.1.8.2" class="ltx_td ltx_align_center">11.98</td>
<td id="S4.T5.4.1.8.3" class="ltx_td ltx_align_center">11.70</td>
<td id="S4.T5.4.1.8.4" class="ltx_td ltx_align_center">12.04</td>
<td id="S4.T5.4.1.8.5" class="ltx_td ltx_align_center">12.17</td>
<td id="S4.T5.4.1.8.6" class="ltx_td ltx_align_center">11.90</td>
<td id="S4.T5.4.1.8.7" class="ltx_td ltx_align_center">12.10</td>
<td id="S4.T5.4.1.8.8" class="ltx_td ltx_align_center">14.11</td>
<td id="S4.T5.4.1.8.9" class="ltx_td ltx_align_center">13.67</td>
<td id="S4.T5.4.1.8.10" class="ltx_td ltx_nopad_r ltx_align_center">14.00</td>
</tr>
<tr id="S4.T5.4.1.9" class="ltx_tr">
<td id="S4.T5.4.1.9.1" class="ltx_td ltx_align_left ltx_border_bb">CoA (8k)</td>
<td id="S4.T5.4.1.9.2" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T5.4.1.9.2.1" class="ltx_text ltx_font_bold">13.28</span></td>
<td id="S4.T5.4.1.9.3" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T5.4.1.9.3.1" class="ltx_text ltx_font_bold">13.73</span></td>
<td id="S4.T5.4.1.9.4" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T5.4.1.9.4.1" class="ltx_text ltx_font_bold">13.70</span></td>
<td id="S4.T5.4.1.9.5" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T5.4.1.9.5.1" class="ltx_text ltx_font_bold">14.92</span></td>
<td id="S4.T5.4.1.9.6" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T5.4.1.9.6.1" class="ltx_text ltx_font_bold">15.05</span></td>
<td id="S4.T5.4.1.9.7" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T5.4.1.9.7.1" class="ltx_text ltx_font_bold">14.96</span></td>
<td id="S4.T5.4.1.9.8" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T5.4.1.9.8.1" class="ltx_text ltx_font_bold">17.74</span></td>
<td id="S4.T5.4.1.9.9" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T5.4.1.9.9.1" class="ltx_text ltx_font_bold">16.68</span></td>
<td id="S4.T5.4.1.9.10" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_bb"><span id="S4.T5.4.1.9.10.1" class="ltx_text ltx_font_bold">17.47</span></td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span id="S4.T5.2.1.1" class="ltx_text" style="font-size:90%;">Table 5</span>: </span><span id="S4.T5.3.2" class="ltx_text" style="font-size:90%;">Comparison with long context LLMs on NarrativeQA and BookSum. CoA significantly outperforms Claude 3 with 200k context limits. No Trun./Trun. indicates the source text in the sample is less/more than 200k tokens which does not need/needs truncation for vanilla (200k) baseline. Average is the mean value across all samples.</span></figcaption>
</figure>
::::

:::: {#S4.SS2.SSS0.Px4 .section .ltx_paragraph}
#### Other Multi-Agent Frameworks. {#other-multi-agent-frameworks. .ltx_title .ltx_title_paragraph}

::: {#S4.SS2.SSS0.Px4.p1 .ltx_para}
As shown in Table [[6]{.ltx_text .ltx_ref_tag}](#S4.T6 "Table 6 ‣ Other Multi-Agent Frameworks. ‣ 4.2 Overall Results of CoA ‣ 4 Experiment ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}, Hierarchical (8k) outperforms Vanilla (8k) on five out of eight datasets, demonstrating the hierarchical approach can also improve the vanilla baselines significantly. Merge (8k) is lower than Vanilla (8k) except for GovReport, showing that merging is especially effective for long summarization tasks such as GovReport. As can be seen, CoA outperforms Hierarchical and Merge on all eight datasets. The reason behind the results is because Hierarchical and Merge do not allow workers to communicate with each other due to their parallel designs. Thus, each worker can only maintain the information in its own chunk $c_{i}$ which blocks the understanding of the whole text, hurting the performance greatly.
:::

<figure id="S4.T6" class="ltx_table">
<div id="S4.T6.4" class="ltx_inline-block ltx_transformed_outer" style="width:433.6pt;height:67.7pt;vertical-align:-0.0pt;">
<span class="ltx_transformed_inner" style="transform:translate(-71.3pt,11.1pt) scale(0.752539096377486,0.752539096377486) ;"> </span>
<table id="S4.T6.4.1" class="ltx_tabular ltx_align_middle">
<tbody>
<tr id="S4.T6.4.1.1" class="ltx_tr">
<td id="S4.T6.4.1.1.1" class="ltx_td ltx_border_tt"></td>
<td id="S4.T6.4.1.1.2" class="ltx_td ltx_align_left ltx_border_tt">HotpotQA</td>
<td id="S4.T6.4.1.1.3" class="ltx_td ltx_align_left ltx_border_tt">MuSiQue</td>
<td id="S4.T6.4.1.1.4" class="ltx_td ltx_align_left ltx_border_tt">NarrativeQA</td>
<td id="S4.T6.4.1.1.5" class="ltx_td ltx_align_left ltx_border_tt">Qasper</td>
<td id="S4.T6.4.1.1.6" class="ltx_td ltx_align_left ltx_border_tt">QuALITY</td>
<td id="S4.T6.4.1.1.7" class="ltx_td ltx_align_left ltx_border_tt">QMSum</td>
<td id="S4.T6.4.1.1.8" class="ltx_td ltx_align_left ltx_border_tt">GovReport</td>
<td id="S4.T6.4.1.1.9" class="ltx_td ltx_align_left ltx_border_tt">RepoBench-P</td>
</tr>
<tr id="S4.T6.4.1.2" class="ltx_tr">
<td id="S4.T6.4.1.2.1" class="ltx_td ltx_align_left ltx_border_t">Vanilla (8k)</td>
<td id="S4.T6.4.1.2.2" class="ltx_td ltx_align_center ltx_border_t">45.57</td>
<td id="S4.T6.4.1.2.3" class="ltx_td ltx_align_center ltx_border_t">26.87</td>
<td id="S4.T6.4.1.2.4" class="ltx_td ltx_align_center ltx_border_t">11.96</td>
<td id="S4.T6.4.1.2.5" class="ltx_td ltx_align_center ltx_border_t">26.56</td>
<td id="S4.T6.4.1.2.6" class="ltx_td ltx_align_center ltx_border_t">61.86</td>
<td id="S4.T6.4.1.2.7" class="ltx_td ltx_align_center ltx_border_t">15.45</td>
<td id="S4.T6.4.1.2.8" class="ltx_td ltx_align_center ltx_border_t">20.60</td>
<td id="S4.T6.4.1.2.9" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_t">56.30</td>
</tr>
<tr id="S4.T6.4.1.3" class="ltx_tr">
<td id="S4.T6.4.1.3.1" class="ltx_td ltx_align_left">Merge (8k)</td>
<td id="S4.T6.4.1.3.2" class="ltx_td ltx_align_center">42.96</td>
<td id="S4.T6.4.1.3.3" class="ltx_td ltx_align_center">26.66</td>
<td id="S4.T6.4.1.3.4" class="ltx_td ltx_align_center">11.27</td>
<td id="S4.T6.4.1.3.5" class="ltx_td ltx_align_center">26.78</td>
<td id="S4.T6.4.1.3.6" class="ltx_td ltx_align_center">59.30</td>
<td id="S4.T6.4.1.3.7" class="ltx_td ltx_align_center">9.42</td>
<td id="S4.T6.4.1.3.8" class="ltx_td ltx_align_center">25.38</td>
<td id="S4.T6.4.1.3.9" class="ltx_td ltx_nopad_r ltx_align_center">33.66</td>
</tr>
<tr id="S4.T6.4.1.4" class="ltx_tr">
<td id="S4.T6.4.1.4.1" class="ltx_td ltx_align_left">Hierarchical (8k)</td>
<td id="S4.T6.4.1.4.2" class="ltx_td ltx_align_center">50.62</td>
<td id="S4.T6.4.1.4.3" class="ltx_td ltx_align_center">29.40</td>
<td id="S4.T6.4.1.4.4" class="ltx_td ltx_align_center">17.04</td>
<td id="S4.T6.4.1.4.5" class="ltx_td ltx_align_center">31.39</td>
<td id="S4.T6.4.1.4.6" class="ltx_td ltx_align_center">64.20</td>
<td id="S4.T6.4.1.4.7" class="ltx_td ltx_align_center">15.19</td>
<td id="S4.T6.4.1.4.8" class="ltx_td ltx_align_center">16.54</td>
<td id="S4.T6.4.1.4.9" class="ltx_td ltx_nopad_r ltx_align_center">27.96</td>
</tr>
<tr id="S4.T6.4.1.5" class="ltx_tr">
<td id="S4.T6.4.1.5.1" class="ltx_td ltx_align_left ltx_border_bb">CoA (8k)</td>
<td id="S4.T6.4.1.5.2" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T6.4.1.5.2.1" class="ltx_text ltx_font_bold">53.62</span></td>
<td id="S4.T6.4.1.5.3" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T6.4.1.5.3.1" class="ltx_text ltx_font_bold">37.09</span></td>
<td id="S4.T6.4.1.5.4" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T6.4.1.5.4.1" class="ltx_text ltx_font_bold">25.26</span></td>
<td id="S4.T6.4.1.5.5" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T6.4.1.5.5.1" class="ltx_text ltx_font_bold">37.17</span></td>
<td id="S4.T6.4.1.5.6" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T6.4.1.5.6.1" class="ltx_text ltx_font_bold">65.42</span></td>
<td id="S4.T6.4.1.5.7" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T6.4.1.5.7.1" class="ltx_text ltx_font_bold">16.77</span></td>
<td id="S4.T6.4.1.5.8" class="ltx_td ltx_align_center ltx_border_bb"><span id="S4.T6.4.1.5.8.1" class="ltx_text ltx_font_bold">26.11</span></td>
<td id="S4.T6.4.1.5.9" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_bb"><span id="S4.T6.4.1.5.9.1" class="ltx_text ltx_font_bold">58.25</span></td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span id="S4.T6.2.1.1" class="ltx_text" style="font-size:90%;">Table 6</span>: </span><span id="S4.T6.3.2" class="ltx_text" style="font-size:90%;">Comparison between CoA and other multi-agent frameworks. CoA with sequential agents outperforms other designs with multiple parallel agents including Merge and Hierarchical.</span></figcaption>
</figure>
::::
:::::::::::
::::::::::::::::::::::

::::::::::::::: {#S5 .section .ltx_section}
## [5 ]{.ltx_tag .ltx_tag_section}Analyses {#analyses .ltx_title .ltx_title_section}

:::: {#S5.SS1 .section .ltx_subsection}
### [5.1 ]{.ltx_tag .ltx_tag_subsection}CoA Improvement is More Obvious When RAG Fails to Retrieve Gold Answer {#coa-improvement-is-more-obvious-when-rag-fails-to-retrieve-gold-answer .ltx_title .ltx_title_subsection}

::: {#S5.SS1.p1 .ltx_para}
To demonstrate this, we first classify the samples in NarrativeQA dataset into different bins according to the position (index) of the chunk in RAG processed input that contains the gold answer. Then, we compute the average score of the CoA and RAG results of different bins. Figure [[4]{.ltx_text .ltx_ref_tag}](#S5.F4 "Figure 4 ‣ 5.2 CoA Improvement is More Obvious When Long Context Models Meet Longer Inputs ‣ 5 Analyses ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref} shows the results. As shown in the figure, RAG performs better when the index is smaller (the gold answer is nearer to the top), showing that downstream LLMs rely significantly on the quality of RAG re-ranking. Besides, the performance of RAG is positively correlated to CoA's when it successfully retrieves the gold answer. However, when RAG fails, CoA can greatly improve the performance (much higher than the tendency line).
:::
::::

:::: {#S5.SS2 .section .ltx_subsection}
### [5.2 ]{.ltx_tag .ltx_tag_subsection}CoA Improvement is More Obvious When Long Context Models Meet Longer Inputs {#coa-improvement-is-more-obvious-when-long-context-models-meet-longer-inputs .ltx_title .ltx_title_subsection}

<figure id="S5.F2" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S5.F2.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2406.02818/assets/figures/haiku.png" id="S5.F2.sf1.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="598" height="306" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span id="S5.F2.sf1.2.1.1" class="ltx_text" style="font-size:90%;">(a)</span> </span><span id="S5.F2.sf1.3.2" class="ltx_text" style="font-size:90%;">Claude 3 Haiku</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S5.F2.sf2" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2406.02818/assets/figures/sonnet.png" id="S5.F2.sf2.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="598" height="306" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span id="S5.F2.sf2.2.1.1" class="ltx_text" style="font-size:90%;">(b)</span> </span><span id="S5.F2.sf2.3.2" class="ltx_text" style="font-size:90%;">Claude 3 Sonnet</span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S5.F2.sf3" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2406.02818/assets/figures/opus.png" id="S5.F2.sf3.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="598" height="306" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span id="S5.F2.sf3.2.1.1" class="ltx_text" style="font-size:90%;">(c)</span> </span><span id="S5.F2.sf3.3.2" class="ltx_text" style="font-size:90%;">Claude 3 Opus</span></figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span id="S5.F2.2.1.1" class="ltx_text" style="font-size:90%;">Figure 2</span>: </span><span id="S5.F2.3.2" class="ltx_text" style="font-size:90%;">Performance of Claude 3 on BookSum. Improvement is more obvious for longer inputs.</span></figcaption>
</figure>

::: {#S5.SS2.p1 .ltx_para}
We compare the performance of CoA and Vanilla with Claude 3 on BookSum. As shown in Figure [[2]{.ltx_text .ltx_ref_tag}](#S5.F2 "Figure 2 ‣ 5.2 CoA Improvement is More Obvious When Long Context Models Meet Longer Inputs ‣ 5 Analyses ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}, CoA can outperform the vanilla baseline by a large margin on various source lengths. It is worth noting that, when the length of the sample increases, the performance even increases for CoA, and the improvement over Vanilla (200k) baseline becomes more significant. The improvement of CoA reaches around 100% when the length is larger than 400k. Thus, we can conclude that 1) CoA can still enhance the LLM performance even though the model has a very long context window limit; and 2) CoA delivers more performance gains when the input is longer.
:::

<figure id="S5.F4" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S5.F4.1" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_top" style="width:212.5pt;">
<img src="/html/2406.02818/assets/x1.png" id="S5.F4.1.g1" class="ltx_graphics ltx_img_landscape" width="461" height="267" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span id="S5.F4.1.1.1.1" class="ltx_text" style="font-size:90%;">Figure 3</span>: </span><span id="S5.F4.1.2.2" class="ltx_text" style="font-size:90%;">Comparison on NarrativeQA. X-axis/Y-axis indicate RAG/CoA performance while each point represents a bin. The number indicates the chunk index of gold answer (ratio of number of samples in bracket), and the size of the point indicates the improvement of CoA over RAG. </span></figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S5.F4.2" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_top" style="width:212.5pt;">
<img src="/html/2406.02818/assets/figures/lostinthemiddle.png" id="S5.F4.2.g1" class="ltx_graphics ltx_img_landscape" width="598" height="338" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span id="S5.F4.2.1.1.1" class="ltx_text" style="font-size:90%;">Figure 4</span>: </span><span id="S5.F4.2.2.2" class="ltx_text" style="font-size:90%;">Performance of CoA and Full on Natural Questions. CoA mitigates the lost-in-the-middle issue. X-axis is the index of document with gold answer where small number indicates gold answer is closer to start.</span></figcaption>
</figure>
</div>
</div>
</figure>
::::

:::: {#S5.SS3 .section .ltx_subsection}
### [5.3 ]{.ltx_tag .ltx_tag_subsection}CoA Mitigates "Lost-in-the-Middle" Phenomenon {#coa-mitigates-lost-in-the-middle-phenomenon .ltx_title .ltx_title_subsection}

::: {#S5.SS3.p1 .ltx_para}
To assess the "lost-in-the-middle" \[[38](#bib.bib38){.ltx_ref}\] effect on Vanilla and CoA models, we replicated the original study by randomly selecting 500 samples from their dataset to create a QA dataset. The results are displayed in Figure [[4]{.ltx_text .ltx_ref_tag}](#S5.F4 "Figure 4 ‣ 5.2 CoA Improvement is More Obvious When Long Context Models Meet Longer Inputs ‣ 5 Analyses ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}. The Vanilla model exhibits a significant \"lost-in-the-middle\" issue, with a performance range of 6.13 ($\pm 2.17$). In contrast, CoA shows resilience against this issue, with a narrower performance gap of 4.89 ($\pm 1.91$), demonstrating that CoA effectively mitigates this problem by providing each agent a shorter context to focus on.
:::
::::

:::: {#S5.SS4 .section .ltx_subsection}
### [5.4 ]{.ltx_tag .ltx_tag_subsection}Multi-agent Collaboration in CoA Enables Complex Reasoning over Long Context {#multi-agent-collaboration-in-coa-enables-complex-reasoning-over-long-context .ltx_title .ltx_title_subsection}

<figure id="S5.F5" class="ltx_figure">
<img src="/html/2406.02818/assets/figures/case_study.png" id="S5.F5.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="598" height="223" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span id="S5.F5.2.1.1" class="ltx_text" style="font-size:90%;">Figure 5</span>: </span><span id="S5.F5.3.2" class="ltx_text" style="font-size:90%;">A case study of RAG (left) and CoA (right) on HotpotQA. The sequential agent communication enables CoA to perform complex multi-hop reasoning over long contexts.</span></figcaption>
</figure>

::: {#S5.SS4.p1 .ltx_para}
Figure [[5]{.ltx_text .ltx_ref_tag}](#S5.F5 "Figure 5 ‣ 5.4 Multi-agent Collaboration in CoA Enables Complex Reasoning over Long Context ‣ 5 Analyses ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref} displays a sample prediction from HotpotQA. To find the correct answer, RAG retrieves text chunks with high semantic similarity with the [query]{#S5.SS4.p1.1.1 .ltx_text .ltx_font_italic}. However, conducting multi-hop reasoning is challenging as the critical first-hop [answer]{#S5.SS4.p1.1.2 .ltx_text .ltx_font_italic} often lacks semantic relevance to the [query]{#S5.SS4.p1.1.3 .ltx_text .ltx_font_italic}. In contrast, CoA operates differently: the first agent explores related topics without knowing the query's answer, aiding subsequent inference. The second agent, also unaware of the answer, broadens the topic scope by incorporating new information. The third agent finally discovers the answer, synthesizing information from earlier agents and new data to complete the reasoning chain. This collaborative approach highlights CoA's ability to facilitate complex reasoning across long context tasks.
:::
::::

:::: {#S5.SS5 .section .ltx_subsection}
### [5.5 ]{.ltx_tag .ltx_tag_subsection}Ablation Study: Effectiveness of Manager and Alternative Design Choices {#ablation-study-effectiveness-of-manager-and-alternative-design-choices .ltx_title .ltx_title_subsection}

<figure id="S5.T7" class="ltx_table">
<div id="S5.T7.4" class="ltx_inline-block ltx_transformed_outer" style="width:433.6pt;height:83.3pt;vertical-align:-0.0pt;">
<span class="ltx_transformed_inner" style="transform:translate(-17.3pt,3.3pt) scale(0.926097034795006,0.926097034795006) ;"> </span>
<table id="S5.T7.4.1" class="ltx_tabular ltx_align_middle">
<tbody>
<tr id="S5.T7.4.1.1" class="ltx_tr">
<td id="S5.T7.4.1.1.1" class="ltx_td ltx_border_tt"></td>
<td id="S5.T7.4.1.1.2" class="ltx_td ltx_align_center ltx_border_tt">HotpotQA</td>
<td id="S5.T7.4.1.1.3" class="ltx_td ltx_align_center ltx_border_tt">MuSiQue</td>
<td id="S5.T7.4.1.1.4" class="ltx_td ltx_align_center ltx_border_tt">NarrativeQA</td>
<td id="S5.T7.4.1.1.5" class="ltx_td ltx_align_center ltx_border_tt">Qasper</td>
<td id="S5.T7.4.1.1.6" class="ltx_td ltx_align_center ltx_border_tt">QuALITY</td>
<td id="S5.T7.4.1.1.7" class="ltx_td ltx_align_center ltx_border_tt">QMSum</td>
<td id="S5.T7.4.1.1.8" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_tt">RepoBench-P</td>
</tr>
<tr id="S5.T7.4.1.2" class="ltx_tr">
<td id="S5.T7.4.1.2.1" class="ltx_td ltx_align_left ltx_border_t">CoA</td>
<td id="S5.T7.4.1.2.2" class="ltx_td ltx_align_center ltx_border_t">53.62</td>
<td id="S5.T7.4.1.2.3" class="ltx_td ltx_align_center ltx_border_t"><span id="S5.T7.4.1.2.3.1" class="ltx_text ltx_font_bold">37.09</span></td>
<td id="S5.T7.4.1.2.4" class="ltx_td ltx_align_center ltx_border_t"><span id="S5.T7.4.1.2.4.1" class="ltx_text ltx_font_bold">25.26</span></td>
<td id="S5.T7.4.1.2.5" class="ltx_td ltx_align_center ltx_border_t"><span id="S5.T7.4.1.2.5.1" class="ltx_text ltx_font_bold">37.17</span></td>
<td id="S5.T7.4.1.2.6" class="ltx_td ltx_align_center ltx_border_t"><span id="S5.T7.4.1.2.6.1" class="ltx_text ltx_font_bold">65.42</span></td>
<td id="S5.T7.4.1.2.7" class="ltx_td ltx_align_center ltx_border_t"><span id="S5.T7.4.1.2.7.1" class="ltx_text ltx_font_bold">16.77</span></td>
<td id="S5.T7.4.1.2.8" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_t">58.25</td>
</tr>
<tr id="S5.T7.4.1.3" class="ltx_tr">
<td id="S5.T7.4.1.3.1" class="ltx_td ltx_align_left">w/o Manager</td>
<td id="S5.T7.4.1.3.2" class="ltx_td ltx_align_center">48.58</td>
<td id="S5.T7.4.1.3.3" class="ltx_td ltx_align_center">26.79</td>
<td id="S5.T7.4.1.3.4" class="ltx_td ltx_align_center">20.80</td>
<td id="S5.T7.4.1.3.5" class="ltx_td ltx_align_center">29.66</td>
<td id="S5.T7.4.1.3.6" class="ltx_td ltx_align_center">58.80</td>
<td id="S5.T7.4.1.3.7" class="ltx_td ltx_align_center">16.50</td>
<td id="S5.T7.4.1.3.8" class="ltx_td ltx_nopad_r ltx_align_center">56.16</td>
</tr>
<tr id="S5.T7.4.1.4" class="ltx_tr">
<td id="S5.T7.4.1.4.1" class="ltx_td ltx_align_left">Right-to-Left</td>
<td id="S5.T7.4.1.4.2" class="ltx_td ltx_align_center">51.83</td>
<td id="S5.T7.4.1.4.3" class="ltx_td ltx_align_center">29.77</td>
<td id="S5.T7.4.1.4.4" class="ltx_td ltx_align_center">21.57</td>
<td id="S5.T7.4.1.4.5" class="ltx_td ltx_align_center">36.60</td>
<td id="S5.T7.4.1.4.6" class="ltx_td ltx_align_center">62.80</td>
<td id="S5.T7.4.1.4.7" class="ltx_td ltx_align_center">15.91</td>
<td id="S5.T7.4.1.4.8" class="ltx_td ltx_nopad_r ltx_align_center">55.10</td>
</tr>
<tr id="S5.T7.4.1.5" class="ltx_tr">
<td id="S5.T7.4.1.5.1" class="ltx_td ltx_align_left ltx_border_bb">Permutation</td>
<td id="S5.T7.4.1.5.2" class="ltx_td ltx_align_center ltx_border_bb"><span id="S5.T7.4.1.5.2.1" class="ltx_text ltx_font_bold">56.05</span></td>
<td id="S5.T7.4.1.5.3" class="ltx_td ltx_align_center ltx_border_bb">34.55</td>
<td id="S5.T7.4.1.5.4" class="ltx_td ltx_align_center ltx_border_bb">23.60</td>
<td id="S5.T7.4.1.5.5" class="ltx_td ltx_align_center ltx_border_bb">37.42</td>
<td id="S5.T7.4.1.5.6" class="ltx_td ltx_align_center ltx_border_bb">64.60</td>
<td id="S5.T7.4.1.5.7" class="ltx_td ltx_align_center ltx_border_bb">16.50</td>
<td id="S5.T7.4.1.5.8" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_bb"><span id="S5.T7.4.1.5.8.1" class="ltx_text ltx_font_bold">58.43</span></td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span id="S5.T7.2.1.1" class="ltx_text" style="font-size:90%;">Table 7</span>: </span><span id="S5.T7.3.2" class="ltx_text" style="font-size:90%;">Ablation on CoA. Manager plays an important role in CoA, and left-to-right yields the best performance among other reading orders including Right-to-Left and Permutation.</span></figcaption>
</figure>

<figure id="S5.T8" class="ltx_table">
<div id="S5.T8.4" class="ltx_inline-block ltx_transformed_outer" style="width:433.6pt;height:212.2pt;vertical-align:-0.0pt;">
<span class="ltx_transformed_inner" style="transform:translate(-22.2pt,10.9pt) scale(0.906995211242792,0.906995211242792) ;"> </span>
<table id="S5.T8.4.1" class="ltx_tabular ltx_align_middle">
<tbody>
<tr id="S5.T8.4.1.1" class="ltx_tr">
<td id="S5.T8.4.1.1.1" class="ltx_td ltx_border_tt"></td>
<td id="S5.T8.4.1.1.2" class="ltx_td ltx_align_left ltx_border_tt">HotpotQA</td>
<td id="S5.T8.4.1.1.3" class="ltx_td ltx_align_left ltx_border_tt">MuSiQue</td>
<td id="S5.T8.4.1.1.4" class="ltx_td ltx_align_left ltx_border_tt">NarrativeQA</td>
<td id="S5.T8.4.1.1.5" class="ltx_td ltx_align_left ltx_border_tt">Qasper</td>
<td id="S5.T8.4.1.1.6" class="ltx_td ltx_align_left ltx_border_tt">QuALITY</td>
<td id="S5.T8.4.1.1.7" class="ltx_td ltx_align_left ltx_border_tt">QMSum</td>
<td id="S5.T8.4.1.1.8" class="ltx_td ltx_align_left ltx_border_tt">RepoBench-P</td>
</tr>
<tr id="S5.T8.4.1.2" class="ltx_tr">
<td colspan="8" id="S5.T8.4.1.2.1" class="ltx_td ltx_align_left ltx_border_t"><span id="S5.T8.4.1.2.1.1" class="ltx_text ltx_font_bold">Bi-direction: left-to-right and right-to-left paths (2-way)</span></td>
</tr>
<tr id="S5.T8.4.1.3" class="ltx_tr">
<td id="S5.T8.4.1.3.1" class="ltx_td ltx_align_left">w/ judge</td>
<td id="S5.T8.4.1.3.2" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.3.2.1" class="ltx_text" style="background-color:#D9EAD3;">54.11</span></td>
<td id="S5.T8.4.1.3.3" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.3.3.1" class="ltx_text" style="background-color:#D9EAD3;">36.97</span></td>
<td id="S5.T8.4.1.3.4" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.3.4.1" class="ltx_text" style="background-color:#D9EAD3;">24.47</span></td>
<td id="S5.T8.4.1.3.5" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.3.5.1" class="ltx_text" style="background-color:#D9EAD3;">39.23</span></td>
<td id="S5.T8.4.1.3.6" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.3.6.1" class="ltx_text" style="background-color:#D9EAD3;">65.00</span></td>
<td id="S5.T8.4.1.3.7" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.3.7.1" class="ltx_text" style="background-color:#D9EAD3;">16.41</span></td>
<td id="S5.T8.4.1.3.8" class="ltx_td ltx_nopad_r ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.3.8.1" class="ltx_text" style="background-color:#D9EAD3;">49.69</span></td>
</tr>
<tr id="S5.T8.4.1.4" class="ltx_tr">
<td id="S5.T8.4.1.4.1" class="ltx_td ltx_align_left">w/ vote</td>
<td id="S5.T8.4.1.4.2" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.4.2.1" class="ltx_text" style="background-color:#F4CCCC;">51.46</span></td>
<td id="S5.T8.4.1.4.3" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.4.3.1" class="ltx_text" style="background-color:#F4CCCC;">34.16</span></td>
<td id="S5.T8.4.1.4.4" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.4.4.1" class="ltx_text" style="background-color:#F4CCCC;">21.59</span></td>
<td id="S5.T8.4.1.4.5" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.4.5.1" class="ltx_text" style="background-color:#F4CCCC;">37.42</span></td>
<td id="S5.T8.4.1.4.6" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.4.6.1" class="ltx_text" style="background-color:#F4CCCC;">64.60</span></td>
<td id="S5.T8.4.1.4.7" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.4.7.1" class="ltx_text" style="background-color:#F4CCCC;">9.51</span></td>
<td id="S5.T8.4.1.4.8" class="ltx_td ltx_nopad_r ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.4.8.1" class="ltx_text" style="background-color:#F4CCCC;">38.70</span></td>
</tr>
<tr id="S5.T8.4.1.5" class="ltx_tr">
<td id="S5.T8.4.1.5.1" class="ltx_td ltx_align_left">oracle</td>
<td id="S5.T8.4.1.5.2" class="ltx_td ltx_align_center">62.12</td>
<td id="S5.T8.4.1.5.3" class="ltx_td ltx_align_center">48.02</td>
<td id="S5.T8.4.1.5.4" class="ltx_td ltx_align_center">32.84</td>
<td id="S5.T8.4.1.5.5" class="ltx_td ltx_align_center">46.37</td>
<td id="S5.T8.4.1.5.6" class="ltx_td ltx_align_center">71.80</td>
<td id="S5.T8.4.1.5.7" class="ltx_td ltx_align_center">18.83</td>
<td id="S5.T8.4.1.5.8" class="ltx_td ltx_nopad_r ltx_align_center">59.67</td>
</tr>
<tr id="S5.T8.4.1.6" class="ltx_tr">
<td colspan="8" id="S5.T8.4.1.6.1" class="ltx_td ltx_align_left ltx_border_t"><span id="S5.T8.4.1.6.1.1" class="ltx_text ltx_font_bold">Self-Consistency: five left-to-right reasoning paths (5-way)</span></td>
</tr>
<tr id="S5.T8.4.1.7" class="ltx_tr">
<td id="S5.T8.4.1.7.1" class="ltx_td ltx_align_left">w/ judge</td>
<td id="S5.T8.4.1.7.2" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.7.2.1" class="ltx_text" style="background-color:#F4CCCC;">57.17</span></td>
<td id="S5.T8.4.1.7.3" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.7.3.1" class="ltx_text" style="background-color:#F4CCCC;">38.82</span></td>
<td id="S5.T8.4.1.7.4" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.7.4.1" class="ltx_text" style="background-color:#F4CCCC;">21.58</span></td>
<td id="S5.T8.4.1.7.5" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.7.5.1" class="ltx_text" style="background-color:#F4CCCC;">36.24</span></td>
<td id="S5.T8.4.1.7.6" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.7.6.1" class="ltx_text" style="background-color:#F4CCCC;">62.80</span></td>
<td id="S5.T8.4.1.7.7" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.7.7.1" class="ltx_text" style="background-color:#D9EAD3;">17.06</span></td>
<td id="S5.T8.4.1.7.8" class="ltx_td ltx_nopad_r ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.7.8.1" class="ltx_text" style="background-color:#D9EAD3;">46.97</span></td>
</tr>
<tr id="S5.T8.4.1.8" class="ltx_tr">
<td id="S5.T8.4.1.8.1" class="ltx_td ltx_align_left">w/ vote</td>
<td id="S5.T8.4.1.8.2" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.8.2.1" class="ltx_text" style="background-color:#D9EAD3;">57.49</span></td>
<td id="S5.T8.4.1.8.3" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.8.3.1" class="ltx_text" style="background-color:#D9EAD3;">40.78</span></td>
<td id="S5.T8.4.1.8.4" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.8.4.1" class="ltx_text" style="background-color:#D9EAD3;">25.56</span></td>
<td id="S5.T8.4.1.8.5" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.8.5.1" class="ltx_text" style="background-color:#D9EAD3;">39.15</span></td>
<td id="S5.T8.4.1.8.6" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.8.6.1" class="ltx_text" style="background-color:#D9EAD3;">68.60</span></td>
<td id="S5.T8.4.1.8.7" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.8.7.1" class="ltx_text" style="background-color:#F4CCCC;">8.35</span></td>
<td id="S5.T8.4.1.8.8" class="ltx_td ltx_nopad_r ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.8.8.1" class="ltx_text" style="background-color:#F4CCCC;">35.56</span></td>
</tr>
<tr id="S5.T8.4.1.9" class="ltx_tr">
<td id="S5.T8.4.1.9.1" class="ltx_td ltx_align_left">oracle</td>
<td id="S5.T8.4.1.9.2" class="ltx_td ltx_align_center">67.07</td>
<td id="S5.T8.4.1.9.3" class="ltx_td ltx_align_center">55.74</td>
<td id="S5.T8.4.1.9.4" class="ltx_td ltx_align_center">39.89</td>
<td id="S5.T8.4.1.9.5" class="ltx_td ltx_align_center">52.74</td>
<td id="S5.T8.4.1.9.6" class="ltx_td ltx_align_center">80.40</td>
<td id="S5.T8.4.1.9.7" class="ltx_td ltx_align_center">20.81</td>
<td id="S5.T8.4.1.9.8" class="ltx_td ltx_nopad_r ltx_align_center">63.52</td>
</tr>
<tr id="S5.T8.4.1.10" class="ltx_tr">
<td colspan="8" id="S5.T8.4.1.10.1" class="ltx_td ltx_align_left ltx_border_t"><span id="S5.T8.4.1.10.1.1" class="ltx_text ltx_font_bold">Permutation: five random order reasoning paths (5-way)</span></td>
</tr>
<tr id="S5.T8.4.1.11" class="ltx_tr">
<td id="S5.T8.4.1.11.1" class="ltx_td ltx_align_left">w/ judge</td>
<td id="S5.T8.4.1.11.2" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.11.2.1" class="ltx_text" style="background-color:#D9EAD3;">59.17</span></td>
<td id="S5.T8.4.1.11.3" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.11.3.1" class="ltx_text" style="background-color:#D9EAD3;">42.37</span></td>
<td id="S5.T8.4.1.11.4" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.11.4.1" class="ltx_text" style="background-color:#F4CCCC;">25.47</span></td>
<td id="S5.T8.4.1.11.5" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.11.5.1" class="ltx_text" style="background-color:#F4CCCC;">37.65</span></td>
<td id="S5.T8.4.1.11.6" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.11.6.1" class="ltx_text" style="background-color:#F4CCCC;">63.40</span></td>
<td id="S5.T8.4.1.11.7" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.11.7.1" class="ltx_text" style="background-color:#D9EAD3;">17.81</span></td>
<td id="S5.T8.4.1.11.8" class="ltx_td ltx_nopad_r ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.11.8.1" class="ltx_text" style="background-color:#D9EAD3;">52.45</span></td>
</tr>
<tr id="S5.T8.4.1.12" class="ltx_tr">
<td id="S5.T8.4.1.12.1" class="ltx_td ltx_align_left">w/ vote</td>
<td id="S5.T8.4.1.12.2" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.12.2.1" class="ltx_text" style="background-color:#F4CCCC;">58.29</span></td>
<td id="S5.T8.4.1.12.3" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.12.3.1" class="ltx_text" style="background-color:#F4CCCC;">39.17</span></td>
<td id="S5.T8.4.1.12.4" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.12.4.1" class="ltx_text" style="background-color:#D9EAD3;">26.58</span></td>
<td id="S5.T8.4.1.12.5" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.12.5.1" class="ltx_text" style="background-color:#D9EAD3;">38.09</span></td>
<td id="S5.T8.4.1.12.6" class="ltx_td ltx_align_center" style="background-color: #D9EAD3"><span id="S5.T8.4.1.12.6.1" class="ltx_text" style="background-color:#D9EAD3;">67.60</span></td>
<td id="S5.T8.4.1.12.7" class="ltx_td ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.12.7.1" class="ltx_text" style="background-color:#F4CCCC;">8.31</span></td>
<td id="S5.T8.4.1.12.8" class="ltx_td ltx_nopad_r ltx_align_center" style="background-color: #F4CCCC"><span id="S5.T8.4.1.12.8.1" class="ltx_text" style="background-color:#F4CCCC;">35.44</span></td>
</tr>
<tr id="S5.T8.4.1.13" class="ltx_tr">
<td id="S5.T8.4.1.13.1" class="ltx_td ltx_align_left ltx_border_bb">oracle</td>
<td id="S5.T8.4.1.13.2" class="ltx_td ltx_align_center ltx_border_bb">75.73</td>
<td id="S5.T8.4.1.13.3" class="ltx_td ltx_align_center ltx_border_bb">60.16</td>
<td id="S5.T8.4.1.13.4" class="ltx_td ltx_align_center ltx_border_bb">39.58</td>
<td id="S5.T8.4.1.13.5" class="ltx_td ltx_align_center ltx_border_bb">52.22</td>
<td id="S5.T8.4.1.13.6" class="ltx_td ltx_align_center ltx_border_bb">79.80</td>
<td id="S5.T8.4.1.13.7" class="ltx_td ltx_align_center ltx_border_bb">20.88</td>
<td id="S5.T8.4.1.13.8" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_bb">67.80</td>
</tr>
</tbody>
</table>
</div>
<figcaption><span class="ltx_tag ltx_tag_table"><span id="S5.T8.2.1.1" class="ltx_text" style="font-size:90%;">Table 8</span>: </span><span id="S5.T8.3.2" class="ltx_text" style="font-size:90%;">Comparison of three multi-path augmentation through judge or voting. Multi-path CoA furthers enhance the performance. </span></figcaption>
</figure>

::: {#S5.SS5.p1 .ltx_para}
To demonstrate the effect of the manager, we conduct an ablation study that uses the last worker to generate results directly. As shown in Table [[7]{.ltx_text .ltx_ref_tag}](#S5.T7 "Table 7 ‣ 5.5 Ablation Study: Effectiveness of Manager and Alternative Design Choices ‣ 5 Analyses ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}, "w/o Manager" hurts the performance significantly, dropping more than 10% on MuSiQue. This demonstrates the important role of the manager. Next, to empirically verify that left-to-right yields the best performance, we evaluate other orders of reading, including Right-to-Left by reading from the last chunk to the first one and Permutation which reads in random order. As shown in Table [[7]{.ltx_text .ltx_ref_tag}](#S5.T7 "Table 7 ‣ 5.5 Ablation Study: Effectiveness of Manager and Alternative Design Choices ‣ 5 Analyses ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}, on most of the datasets, left-to-right yields the highest score, demonstrating the advantages of natural reading order.
:::
::::

:::: {#S5.SS6 .section .ltx_subsection}
### [5.6 ]{.ltx_tag .ltx_tag_subsection}Multi-path Chain-of-Agents Further Enhances Performance {#multi-path-chain-of-agents-further-enhances-performance .ltx_title .ltx_title_subsection}

::: {#S5.SS6.p1 .ltx_para}
We manually investigated the results over these three orders (left-to-right, right-to-left, permutation), and we found that other orders sometimes can produce better answers than left-to-right. Inspired by this observation, we explore two approaches to select the best result among multiple paths. w/ vote applies majority voting over the final results while w/ judge uses an LLM to judge the most reliable $C\hspace{0pt}U_{l}$ of diverse paths and generate the final answer. Oracle picks the best path by evaluating score of each path, yielding the upper bound performance. Table [[8]{.ltx_text .ltx_ref_tag}](#S5.T8 "Table 8 ‣ 5.5 Ablation Study: Effectiveness of Manager and Alternative Design Choices ‣ 5 Analyses ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref} compares three multi-path augmentation approaches. Surprisingly, results show that 1) all ensemble approaches (Bi-direction, Self-consistency, and Permutation) can further enhance the performance of CoA and 5-way Permutation yields the best improvement, 2) majority voting (w/ vote) of final answer is better than using an LLM as judge (w/ judge) in Self-consistency, but worse in Bi-direction, 3) using LLM judge (w/ judge) works well on long result generation tasks (QMSum, RepoBench-P), and 4) there is large space to improve because oracle (choose as answer the one with highest performance) is much higher than either w/ judge or w/ vote. We leave the direction of multi-path reasoning to future study.
:::
::::
:::::::::::::::

:::::: {#S6 .section .ltx_section}
## [6 ]{.ltx_tag .ltx_tag_section}Conclusion {#conclusion .ltx_title .ltx_title_section}

::: {#S6.p1 .ltx_para}
In this paper, we propose Chain-of-Agents, a multi-agent LLM collaboration framework for solving long context tasks. It is a training free, task/length agnostic, interpretable, and cost-effective framework. Experiments show that Chain-of-Agents outperforms RAG and Long Context LLMs by a large margin despite of its simple design. Analysis shows that by integrating information aggregation and context reasoning, CoA mitigates lost-in-the-middle and performs better on longer samples.
:::

:::: {#S6.SS0.SSS0.Px1 .section .ltx_paragraph}
#### Limitations. {#limitations. .ltx_title .ltx_title_paragraph}

::: {#S6.SS0.SSS0.Px1.p1 .ltx_para}
While CoA features with a simple and effective design, future directions can address the following limitations to further improve its prowess and efficiency. First, communication effectiveness can be further improved via finetuning or in-context learning because current LLMs are aligned with human norms which is not optimal for communication between LLMs. Second, CoA does not explore other forms of communication approaches, such as debating or complex discussions. Third, the cost and latency of running CoA can be further reduced, such as replacing some LLMs with more effective models via model routing \[[56](#bib.bib56){.ltx_ref}\].
:::
::::
::::::

::: {#bib .section .ltx_bibliography}
## References {#references .ltx_title .ltx_title_bibliography}

- [[\[1\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Josh Achiam, Steven Adler, Sandhini Agarwal, Lama Ahmad, Ilge Akkaya, Florencia Leoni Aleman, Diogo Almeida, Janko Altenschmidt, Sam Altman, Shyamal Anadkat, et al. ]{.ltx_bibblock} [Gpt-4 technical report. ]{.ltx_bibblock} [[arXiv preprint arXiv:2303.08774]{#bib.bib1.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib1}
- [[\[2\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Qingyao Ai, Ting Bai, Zhao Cao, Yi Chang, Jiawei Chen, Zhumin Chen, Zhiyong Cheng, Shoubin Dong, Zhicheng Dou, Fuli Feng, et al. ]{.ltx_bibblock} [Information retrieval meets large language models: a strategic report from chinese ir community. ]{.ltx_bibblock} [[AI Open]{#bib.bib2.1.1 .ltx_text .ltx_font_italic}, 4:80--90, 2023. ]{.ltx_bibblock}]{#bib.bib2}
- [[\[3\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Shengnan An, Zexiong Ma, Zeqi Lin, Nanning Zheng, and Jian-Guang Lou. ]{.ltx_bibblock} [Make your llm fully utilize the context. ]{.ltx_bibblock} [[arXiv preprint arXiv:2404.16811]{#bib.bib3.1.1 .ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib3}
- [[\[4\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Rohan Anil, Andrew M Dai, Orhan Firat, Melvin Johnson, Dmitry Lepikhin, Alexandre Passos, Siamak Shakeri, Emanuel Taropa, Paige Bailey, Zhifeng Chen, et al. ]{.ltx_bibblock} [Palm 2 technical report. ]{.ltx_bibblock} [[arXiv preprint arXiv:2305.10403]{#bib.bib4.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib4}
- [[\[5\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ AI Anthropic. ]{.ltx_bibblock} [The claude 3 model family: Opus, sonnet, haiku. ]{.ltx_bibblock} [[Claude-3 Model Card]{#bib.bib5.1.1 .ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib5}
- [[\[6\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yushi Bai, Xin Lv, Jiajie Zhang, Hongchang Lyu, Jiankai Tang, Zhidian Huang, Zhengxiao Du, Xiao Liu, Aohan Zeng, Lei Hou, Yuxiao Dong, Jie Tang, and Juanzi Li. ]{.ltx_bibblock} [Longbench: A bilingual, multitask benchmark for long context understanding. ]{.ltx_bibblock} [[arXiv preprint arXiv:2308.14508]{#bib.bib6.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib6}
- [[\[7\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tom Brown, Benjamin Mann, Nick Ryder, Melanie Subbiah, Jared D Kaplan, Prafulla Dhariwal, Arvind Neelakantan, Pranav Shyam, Girish Sastry, Amanda Askell, Sandhini Agarwal, Ariel Herbert-Voss, Gretchen Krueger, Tom Henighan, Rewon Child, Aditya Ramesh, Daniel Ziegler, Jeffrey Wu, Clemens Winter, Chris Hesse, Mark Chen, Eric Sigler, Mateusz Litwin, Scott Gray, Benjamin Chess, Jack Clark, Christopher Berner, Sam McCandlish, Alec Radford, Ilya Sutskever, and Dario Amodei. ]{.ltx_bibblock} [Language models are few-shot learners. ]{.ltx_bibblock} [In H. Larochelle, M. Ranzato, R. Hadsell, M.F. Balcan, and H. Lin, editors, [Advances in Neural Information Processing Systems]{#bib.bib7.1.1 .ltx_text .ltx_font_italic}, volume 33, pages 1877--1901. Curran Associates, Inc., 2020. ]{.ltx_bibblock}]{#bib.bib7}
- [[\[8\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tom Brown, Benjamin Mann, Nick Ryder, Melanie Subbiah, Jared D Kaplan, Prafulla Dhariwal, Arvind Neelakantan, Pranav Shyam, Girish Sastry, Amanda Askell, et al. ]{.ltx_bibblock} [Language models are few-shot learners. ]{.ltx_bibblock} [[Advances in neural information processing systems]{#bib.bib8.1.1 .ltx_text .ltx_font_italic}, 33:1877--1901, 2020. ]{.ltx_bibblock}]{#bib.bib8}
- [[\[9\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Howard Chen, Ramakanth Pasunuru, Jason Weston, and Asli Celikyilmaz. ]{.ltx_bibblock} [Walking down the memory maze: Beyond context limit through interactive reading. ]{.ltx_bibblock} [[arXiv preprint arXiv:2310.05029]{#bib.bib9.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib9}
- [[\[10\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Huaben Chen, Wenkang Ji, Lufeng Xu, and Shiyu Zhao. ]{.ltx_bibblock} [Multi-agent consensus seeking via large language models. ]{.ltx_bibblock} [[arXiv preprint arXiv:2310.20151]{#bib.bib10.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib10}
- [[\[11\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Justin Chih-Yao Chen, Swarnadeep Saha, and Mohit Bansal. ]{.ltx_bibblock} [Reconcile: Round-table conference improves reasoning via consensus among diverse llms. ]{.ltx_bibblock} [[arXiv preprint arXiv:2309.13007]{#bib.bib11.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib11}
- [[\[12\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Mingda Chen, Zewei Chu, Sam Wiseman, and Kevin Gimpel. ]{.ltx_bibblock} [SummScreen: A dataset for abstractive screenplay summarization. ]{.ltx_bibblock} [In Smaranda Muresan, Preslav Nakov, and Aline Villavicencio, editors, [Proceedings of the 60th Annual Meeting of the Association for Computational Linguistics (Volume 1: Long Papers)]{#bib.bib12.1.1 .ltx_text .ltx_font_italic}, pages 8602--8615, Dublin, Ireland, May 2022. Association for Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib12}
- [[\[13\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Shouyuan Chen, Sherman Wong, Liangjian Chen, and Yuandong Tian. ]{.ltx_bibblock} [Extending context window of large language models via positional interpolation. ]{.ltx_bibblock} [[arXiv preprint arXiv:2306.15595]{#bib.bib13.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib13}
- [[\[14\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Pradeep Dasigi, Kyle Lo, Iz Beltagy, Arman Cohan, Noah A. Smith, and Matt Gardner. ]{.ltx_bibblock} [A dataset of information-seeking questions and answers anchored in research papers. ]{.ltx_bibblock} [In Kristina Toutanova, Anna Rumshisky, Luke Zettlemoyer, Dilek Hakkani-Tur, Iz Beltagy, Steven Bethard, Ryan Cotterell, Tanmoy Chakraborty, and Yichao Zhou, editors, [Proceedings of the 2021 Conference of the North American Chapter of the Association for Computational Linguistics: Human Language Technologies]{#bib.bib14.1.1 .ltx_text .ltx_font_italic}, pages 4599--4610, Online, June 2021. Association for Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib14}
- [[\[15\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Pradeep Dasigi, Kyle Lo, Iz Beltagy, Arman Cohan, Noah A Smith, and Matt Gardner. ]{.ltx_bibblock} [A dataset of information-seeking questions and answers anchored in research papers. ]{.ltx_bibblock} [[arXiv preprint arXiv:2105.03011]{#bib.bib15.1.1 .ltx_text .ltx_font_italic}, 2021. ]{.ltx_bibblock}]{#bib.bib15}
- [[\[16\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yilun Du, Shuang Li, Antonio Torralba, Joshua B Tenenbaum, and Igor Mordatch. ]{.ltx_bibblock} [Improving factuality and reasoning in language models through multiagent debate. ]{.ltx_bibblock} [[arXiv preprint arXiv:2305.14325]{#bib.bib16.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib16}
- [[\[17\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Navid Ghaffarzadegan, Aritra Majumdar, Ross Williams, and Niyousha Hosseinichimeh. ]{.ltx_bibblock} [Generative agent-based modeling: Unveiling social system dynamics through coupling mechanistic models with generative artificial intelligence. ]{.ltx_bibblock} [[arXiv preprint arXiv:2309.11456]{#bib.bib17.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib17}
- [[\[18\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Daya Guo, Canwen Xu, Nan Duan, Jian Yin, and Julian McAuley. ]{.ltx_bibblock} [Longcoder: A long-range pre-trained language model for code completion. ]{.ltx_bibblock} [In [International Conference on Machine Learning]{#bib.bib18.1.1 .ltx_text .ltx_font_italic}, pages 12098--12107. PMLR, 2023. ]{.ltx_bibblock}]{#bib.bib18}
- [[\[19\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Taicheng Guo, Xiuying Chen, Yaqi Wang, Ruidi Chang, Shichao Pei, Nitesh V Chawla, Olaf Wiest, and Xiangliang Zhang. ]{.ltx_bibblock} [Large language model based multi-agents: A survey of progress and challenges. ]{.ltx_bibblock} [[arXiv preprint arXiv:2402.01680]{#bib.bib19.1.1 .ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib19}
- [[\[20\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xanh Ho, Anh-Khoa Duong Nguyen, Saku Sugawara, and Akiko Aizawa. ]{.ltx_bibblock} [Constructing a multi-hop QA dataset for comprehensive evaluation of reasoning steps. ]{.ltx_bibblock} [In Donia Scott, Nuria Bel, and Chengqing Zong, editors, [Proceedings of the 28th International Conference on Computational Linguistics]{#bib.bib20.1.1 .ltx_text .ltx_font_italic}, pages 6609--6625, Barcelona, Spain (Online), December 2020. International Committee on Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib20}
- [[\[21\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Wenyue Hua, Lizhou Fan, Lingyao Li, Kai Mei, Jianchao Ji, Yingqiang Ge, Libby Hemphill, and Yongfeng Zhang. ]{.ltx_bibblock} [War and peace (waragent): Large language model-based multi-agent simulation of world wars. ]{.ltx_bibblock} [[arXiv preprint arXiv:2311.17227]{#bib.bib21.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib21}
- [[\[22\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Dong Huang, Qingwen Bu, Jie M Zhang, Michael Luck, and Heming Cui. ]{.ltx_bibblock} [Agentcoder: Multi-agent-based code generation with iterative testing and optimisation. ]{.ltx_bibblock} [[arXiv preprint arXiv:2312.13010]{#bib.bib22.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib22}
- [[\[23\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Luyang Huang, Shuyang Cao, Nikolaus Parulian, Heng Ji, and Lu Wang. ]{.ltx_bibblock} [Efficient attentions for long document summarization. ]{.ltx_bibblock} [In Kristina Toutanova, Anna Rumshisky, Luke Zettlemoyer, Dilek Hakkani-Tur, Iz Beltagy, Steven Bethard, Ryan Cotterell, Tanmoy Chakraborty, and Yichao Zhou, editors, [Proceedings of the 2021 Conference of the North American Chapter of the Association for Computational Linguistics: Human Language Technologies]{#bib.bib23.1.1 .ltx_text .ltx_font_italic}, pages 1419--1436, Online, June 2021. Association for Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib23}
- [[\[24\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Gautier Izacard, Patrick Lewis, Maria Lomeli, Lucas Hosseini, Fabio Petroni, Timo Schick, Jane Dwivedi-Yu, Armand Joulin, Sebastian Riedel, and Edouard Grave. ]{.ltx_bibblock} [Atlas: Few-shot learning with retrieval augmented language models. ]{.ltx_bibblock} [[Journal of Machine Learning Research]{#bib.bib24.1.1 .ltx_text .ltx_font_italic}, 24(251):1--43, 2023. ]{.ltx_bibblock}]{#bib.bib24}
- [[\[25\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Huiqiang Jiang, Qianhui Wu, Xufang Luo, Dongsheng Li, Chin-Yew Lin, Yuqing Yang, and Lili Qiu. ]{.ltx_bibblock} [Longllmlingua: Accelerating and enhancing llms in long context scenarios via prompt compression. ]{.ltx_bibblock} [[arXiv preprint arXiv:2310.06839]{#bib.bib25.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib25}
- [[\[26\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zhengbao Jiang, Luyu Gao, Zhiruo Wang, Jun Araki, Haibo Ding, Jamie Callan, and Graham Neubig. ]{.ltx_bibblock} [Retrieval as attention: End-to-end learning of retrieval and reading within a single transformer. ]{.ltx_bibblock} [In Yoav Goldberg, Zornitsa Kozareva, and Yue Zhang, editors, [Proceedings of the 2022 Conference on Empirical Methods in Natural Language Processing]{#bib.bib26.1.1 .ltx_text .ltx_font_italic}, pages 2336--2349, Abu Dhabi, United Arab Emirates, December 2022. Association for Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib26}
- [[\[27\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tushar Khot, H. Trivedi, Matthew Finlayson, Yao Fu, Kyle Richardson, Peter Clark, and Ashish Sabharwal. ]{.ltx_bibblock} [Decomposed prompting: A modular approach for solving complex tasks. ]{.ltx_bibblock} [[ArXiv]{#bib.bib27.1.1 .ltx_text .ltx_font_italic}, abs/2210.02406, 2022. ]{.ltx_bibblock}]{#bib.bib27}
- [[\[28\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tomáš Kočiskỳ, Jonathan Schwarz, Phil Blunsom, Chris Dyer, Karl Moritz Hermann, Gábor Melis, and Edward Grefenstette. ]{.ltx_bibblock} [The narrativeqa reading comprehension challenge. ]{.ltx_bibblock} [[Transactions of the Association for Computational Linguistics]{#bib.bib28.1.1 .ltx_text .ltx_font_italic}, 6:317--328, 2018. ]{.ltx_bibblock}]{#bib.bib28}
- [[\[29\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tomáš Kočiský, Jonathan Schwarz, Phil Blunsom, Chris Dyer, Karl Moritz Hermann, Gábor Melis, and Edward Grefenstette. ]{.ltx_bibblock} [The NarrativeQA Reading Comprehension Challenge. ]{.ltx_bibblock} [[Transactions of the Association for Computational Linguistics]{#bib.bib29.1.1 .ltx_text .ltx_font_italic}, 6:317--328, 05 2018. ]{.ltx_bibblock}]{#bib.bib29}
- [[\[30\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Wojciech Kryściński, Nazneen Rajani, Divyansh Agarwal, Caiming Xiong, and Dragomir Radev. ]{.ltx_bibblock} [Booksum: A collection of datasets for long-form narrative summarization. ]{.ltx_bibblock} [[arXiv preprint arXiv:2105.08209]{#bib.bib30.1.1 .ltx_text .ltx_font_italic}, 2021. ]{.ltx_bibblock}]{#bib.bib30}
- [[\[31\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Patrick Lewis, Ethan Perez, Aleksandra Piktus, Fabio Petroni, Vladimir Karpukhin, Naman Goyal, Heinrich Küttler, Mike Lewis, Wen-tau Yih, Tim Rocktäschel, et al. ]{.ltx_bibblock} [Retrieval-augmented generation for knowledge-intensive nlp tasks. ]{.ltx_bibblock} [[Advances in Neural Information Processing Systems]{#bib.bib31.1.1 .ltx_text .ltx_font_italic}, 33:9459--9474, 2020. ]{.ltx_bibblock}]{#bib.bib31}
- [[\[32\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Nian Li, Chen Gao, Yong Li, and Qingmin Liao. ]{.ltx_bibblock} [Large language model-empowered agents for simulating macroeconomic activities. ]{.ltx_bibblock} [[arXiv preprint arXiv:2310.10436]{#bib.bib32.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib32}
- [[\[33\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tianle Li, Ge Zhang, Quy Duc Do, Xiang Yue, and Wenhu Chen. ]{.ltx_bibblock} [Long-context llms struggle with long in-context learning. ]{.ltx_bibblock} [[arXiv preprint arXiv:2404.02060]{#bib.bib33.1.1 .ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib33}
- [[\[34\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yucheng Li. ]{.ltx_bibblock} [Unlocking context constraints of llms: Enhancing context efficiency of llms with self-information-based content filtering. ]{.ltx_bibblock} [[arXiv preprint arXiv:2304.12102]{#bib.bib34.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib34}
- [[\[35\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jonathan Light, Min Cai, Sheng Shen, and Ziniu Hu. ]{.ltx_bibblock} [Avalonbench: Evaluating llms playing the game of avalon. ]{.ltx_bibblock} [In [NeurIPS 2023 Foundation Models for Decision Making Workshop]{#bib.bib35.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib35}
- [[\[36\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Chin-Yew Lin. ]{.ltx_bibblock} [ROUGE: A package for automatic evaluation of summaries. ]{.ltx_bibblock} [In [Text Summarization Branches Out]{#bib.bib36.1.1 .ltx_text .ltx_font_italic}, pages 74--81, Barcelona, Spain, 2004. Association for Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib36}
- [[\[37\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Sheng-Chieh Lin, Akari Asai, Minghan Li, Barlas Oguz, Jimmy Lin, Yashar Mehdad, Wen-tau Yih, and Xilun Chen. ]{.ltx_bibblock} [How to train your dragon: Diverse augmentation towards generalizable dense retrieval. ]{.ltx_bibblock} [In Houda Bouamor, Juan Pino, and Kalika Bali, editors, [Findings of the Association for Computational Linguistics: EMNLP 2023]{#bib.bib37.1.1 .ltx_text .ltx_font_italic}, pages 6385--6400, Singapore, December 2023. Association for Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib37}
- [[\[38\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Nelson F Liu, Kevin Lin, John Hewitt, Ashwin Paranjape, Michele Bevilacqua, Fabio Petroni, and Percy Liang. ]{.ltx_bibblock} [Lost in the middle: How language models use long contexts. ]{.ltx_bibblock} [[Transactions of the Association for Computational Linguistics]{#bib.bib38.1.1 .ltx_text .ltx_font_italic}, 12:157--173, 2024. ]{.ltx_bibblock}]{#bib.bib38}
- [[\[39\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tianyang Liu, Canwen Xu, and Julian McAuley. ]{.ltx_bibblock} [Repobench: Benchmarking repository-level code auto-completion systems. ]{.ltx_bibblock} [[arXiv preprint arXiv:2306.03091]{#bib.bib39.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib39}
- [[\[40\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xuezhe Ma, Xiaomeng Yang, Wenhan Xiong, Beidi Chen, Lili Yu, Hao Zhang, Jonathan May, Luke Zettlemoyer, Omer Levy, and Chunting Zhou. ]{.ltx_bibblock} [Megalodon: Efficient llm pretraining and inference with unlimited context length. ]{.ltx_bibblock} [[arXiv preprint arXiv:2404.08801]{#bib.bib40.1.1 .ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib40}
- [[\[41\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Amirkeivan Mohtashami and Martin Jaggi. ]{.ltx_bibblock} [Landmark attention: Random-access infinite context length for transformers. ]{.ltx_bibblock} [[arXiv preprint arXiv:2305.16300]{#bib.bib41.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib41}
- [[\[42\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Gabriel Mukobi, Hannah Erlebach, Niklas Lauffer, Lewis Hammond, Alan Chan, and Jesse Clifton. ]{.ltx_bibblock} [Welfare diplomacy: Benchmarking language model cooperation. ]{.ltx_bibblock} [[arXiv preprint arXiv:2310.08901]{#bib.bib42.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib42}
- [[\[43\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tsendsuren Munkhdalai, Manaal Faruqui, and Siddharth Gopal. ]{.ltx_bibblock} [Leave no context behind: Efficient infinite context transformers with infini-attention. ]{.ltx_bibblock} [[arXiv preprint arXiv:2404.07143]{#bib.bib43.1.1 .ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib43}
- [[\[44\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Bo Ni and Markus J Buehler. ]{.ltx_bibblock} [Mechagents: Large language model multi-agent collaborations can solve mechanics problems, generate new data, and integrate knowledge. ]{.ltx_bibblock} [[Extreme Mechanics Letters]{#bib.bib44.1.1 .ltx_text .ltx_font_italic}, page 102131, 2024. ]{.ltx_bibblock}]{#bib.bib44}
- [[\[45\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ OpenAI. ]{.ltx_bibblock} [Gpt-4 technical report. ]{.ltx_bibblock} [[ArXiv]{#bib.bib45.1.1 .ltx_text .ltx_font_italic}, abs/2303.08774, 2023. ]{.ltx_bibblock}]{#bib.bib45}
- [[\[46\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Richard Yuanzhe Pang, Alicia Parrish, Nitish Joshi, Nikita Nangia, Jason Phang, Angelica Chen, Vishakh Padmakumar, Johnny Ma, Jana Thompson, He He, and Samuel Bowman. ]{.ltx_bibblock} [QuALITY: Question answering with long input texts, yes! ]{.ltx_bibblock} [In [Proceedings of the 2022 Conference of the North American Chapter of the Association for Computational Linguistics: Human Language Technologies]{#bib.bib46.1.1 .ltx_text .ltx_font_italic}, pages 5336--5358, Seattle, United States, July 2022. Association for Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib46}
- [[\[47\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Joon Sung Park, Joseph O'Brien, Carrie Jun Cai, Meredith Ringel Morris, Percy Liang, and Michael S Bernstein. ]{.ltx_bibblock} [Generative agents: Interactive simulacra of human behavior. ]{.ltx_bibblock} [In [Proceedings of the 36th Annual ACM Symposium on User Interface Software and Technology]{#bib.bib47.1.1 .ltx_text .ltx_font_italic}, pages 1--22, 2023. ]{.ltx_bibblock}]{#bib.bib47}
- [[\[48\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Joon Sung Park, Lindsay Popowski, Carrie Cai, Meredith Ringel Morris, Percy Liang, and Michael S Bernstein. ]{.ltx_bibblock} [Social simulacra: Creating populated prototypes for social computing systems. ]{.ltx_bibblock} [In [Proceedings of the 35th Annual ACM Symposium on User Interface Software and Technology]{#bib.bib48.1.1 .ltx_text .ltx_font_italic}, pages 1--18, 2022. ]{.ltx_bibblock}]{#bib.bib48}
- [[\[49\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ethan Perez, Patrick Lewis, Wen-tau Yih, Kyunghyun Cho, and Douwe Kiela. ]{.ltx_bibblock} [Unsupervised question decomposition for question answering. ]{.ltx_bibblock} [In Bonnie Webber, Trevor Cohn, Yulan He, and Yang Liu, editors, [Proceedings of the 2020 Conference on Empirical Methods in Natural Language Processing (EMNLP)]{#bib.bib49.1.1 .ltx_text .ltx_font_italic}, pages 8864--8880, Online, November 2020. Association for Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib49}
- [[\[50\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ofir Press, Noah A Smith, and Mike Lewis. ]{.ltx_bibblock} [Train short, test long: Attention with linear biases enables input length extrapolation. ]{.ltx_bibblock} [[arXiv preprint arXiv:2108.12409]{#bib.bib50.1.1 .ltx_text .ltx_font_italic}, 2021. ]{.ltx_bibblock}]{#bib.bib50}
- [[\[51\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Alec Radford, Jeff Wu, Rewon Child, David Luan, Dario Amodei, and Ilya Sutskever. ]{.ltx_bibblock} [Language models are unsupervised multitask learners, 2019. ]{.ltx_bibblock}]{#bib.bib51}
- [[\[52\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Machel Reid, Nikolay Savinov, Denis Teplyashin, Dmitry Lepikhin, Timothy Lillicrap, Jean-baptiste Alayrac, Radu Soricut, Angeliki Lazaridou, Orhan Firat, Julian Schrittwieser, et al. ]{.ltx_bibblock} [Gemini 1.5: Unlocking multimodal understanding across millions of tokens of context. ]{.ltx_bibblock} [[arXiv preprint arXiv:2403.05530]{#bib.bib52.1.1 .ltx_text .ltx_font_italic}, 2024. ]{.ltx_bibblock}]{#bib.bib52}
- [[\[53\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Swarnadeep Saha, Omer Levy, Asli Celikyilmaz, Mohit Bansal, Jason Weston, and Xian Li. ]{.ltx_bibblock} [Branch-solve-merge improves large language model evaluation and generation. ]{.ltx_bibblock} [[arXiv preprint arXiv:2310.15123]{#bib.bib53.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib53}
- [[\[54\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Uri Shaham, Elad Segal, Maor Ivgi, Avia Efrat, Ori Yoran, Adi Haviv, Ankit Gupta, Wenhan Xiong, Mor Geva, Jonathan Berant, and Omer Levy. ]{.ltx_bibblock} [SCROLLS: Standardized CompaRison over long language sequences. ]{.ltx_bibblock} [In Yoav Goldberg, Zornitsa Kozareva, and Yue Zhang, editors, [Proceedings of the 2022 Conference on Empirical Methods in Natural Language Processing]{#bib.bib54.1.1 .ltx_text .ltx_font_italic}, pages 12007--12021, Abu Dhabi, United Arab Emirates, December 2022. Association for Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib54}
- [[\[55\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Weijia Shi, Sewon Min, Michihiro Yasunaga, Minjoon Seo, Rich James, Mike Lewis, Luke Zettlemoyer, and Wen-tau Yih. ]{.ltx_bibblock} [Replug: Retrieval-augmented black-box language models. ]{.ltx_bibblock} [[arXiv preprint arXiv:2301.12652]{#bib.bib55.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib55}
- [[\[56\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tal Shnitzer, Anthony Ou, Mírian Silva, Kate Soule, Yuekai Sun, Justin Solomon, Neil Thompson, and Mikhail Yurochkin. ]{.ltx_bibblock} [Large language model routing with benchmark datasets. ]{.ltx_bibblock} [[arXiv preprint arXiv:2309.15789]{#bib.bib56.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib56}
- [[\[57\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Aarohi Srivastava, Abhinav Rastogi, Abhishek Rao, Abu Awal Md Shoeb, Abubakar Abid, Adam Fisch, Adam R Brown, Adam Santoro, Aditya Gupta, Adrià Garriga-Alonso, et al. ]{.ltx_bibblock} [Beyond the imitation game: Quantifying and extrapolating the capabilities of language models. ]{.ltx_bibblock} [[arXiv preprint arXiv:2206.04615]{#bib.bib57.1.1 .ltx_text .ltx_font_italic}, 2022. ]{.ltx_bibblock}]{#bib.bib57}
- [[\[58\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jianlin Su, Murtadha Ahmed, Yu Lu, Shengfeng Pan, Wen Bo, and Yunfeng Liu. ]{.ltx_bibblock} [Roformer: Enhanced transformer with rotary position embedding. ]{.ltx_bibblock} [[Neurocomputing]{#bib.bib58.1.1 .ltx_text .ltx_font_italic}, 568:127063, 2024. ]{.ltx_bibblock}]{#bib.bib58}
- [[\[59\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Simeng Sun, Y. Liu, Shuo Wang, Chenguang Zhu, and Mohit Iyyer. ]{.ltx_bibblock} [Pearl: Prompting large language models to plan and execute actions over long documents. ]{.ltx_bibblock} [[ArXiv]{#bib.bib59.1.1 .ltx_text .ltx_font_italic}, abs/2305.14564, 2023. ]{.ltx_bibblock}]{#bib.bib59}
- [[\[60\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xiangru Tang, Anni Zou, Zhuosheng Zhang, Yilun Zhao, Xingyao Zhang, Arman Cohan, and Mark Gerstein. ]{.ltx_bibblock} [Medagents: Large language models as collaborators for zero-shot medical reasoning. ]{.ltx_bibblock} [[arXiv preprint arXiv:2311.10537]{#bib.bib60.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib60}
- [[\[61\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Gemini Team, Rohan Anil, Sebastian Borgeaud, Yonghui Wu, Jean-Baptiste Alayrac, Jiahui Yu, Radu Soricut, Johan Schalkwyk, Andrew M Dai, Anja Hauth, et al. ]{.ltx_bibblock} [Gemini: a family of highly capable multimodal models. ]{.ltx_bibblock} [[arXiv preprint arXiv:2312.11805]{#bib.bib61.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib61}
- [[\[62\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Hugo Touvron, Thibaut Lavril, Gautier Izacard, Xavier Martinet, Marie-Anne Lachaux, Timothée Lacroix, Baptiste Rozière, Naman Goyal, Eric Hambro, Faisal Azhar, et al. ]{.ltx_bibblock} [Llama: Open and efficient foundation language models. ]{.ltx_bibblock} [[arXiv preprint arXiv:2302.13971]{#bib.bib62.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib62}
- [[\[63\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Harsh Trivedi, Niranjan Balasubramanian, Tushar Khot, and Ashish Sabharwal. ]{.ltx_bibblock} [MuSiQue: Multihop questions via single-hop question composition. ]{.ltx_bibblock} [[Transactions of the Association for Computational Linguistics]{#bib.bib63.1.1 .ltx_text .ltx_font_italic}, 10:539--554, 2022. ]{.ltx_bibblock}]{#bib.bib63}
- [[\[64\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Bing Wang, Changyu Ren, Jian Yang, Xinnian Liang, Jiaqi Bai, Qian-Wen Zhang, Zhao Yan, and Zhoujun Li. ]{.ltx_bibblock} [Mac-sql: Multi-agent collaboration for text-to-sql. ]{.ltx_bibblock} [[arXiv preprint arXiv:2312.11242]{#bib.bib64.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib64}
- [[\[65\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Boshi Wang, Xiang Deng, and Huan Sun. ]{.ltx_bibblock} [Iteratively prompt pre-trained language models for chain of thought. ]{.ltx_bibblock} [In Yoav Goldberg, Zornitsa Kozareva, and Yue Zhang, editors, [Proceedings of the 2022 Conference on Empirical Methods in Natural Language Processing]{#bib.bib65.1.1 .ltx_text .ltx_font_italic}, pages 2714--2730, Abu Dhabi, United Arab Emirates, December 2022. Association for Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib65}
- [[\[66\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Boxin Wang, Wei Ping, Peng Xu, Lawrence McAfee, Zihan Liu, Mohammad Shoeybi, Yi Dong, Oleksii Kuchaiev, Bo Li, Chaowei Xiao, Anima Anandkumar, and Bryan Catanzaro. ]{.ltx_bibblock} [Shall we pretrain autoregressive language models with retrieval? a comprehensive study. ]{.ltx_bibblock} [In Houda Bouamor, Juan Pino, and Kalika Bali, editors, [Proceedings of the 2023 Conference on Empirical Methods in Natural Language Processing]{#bib.bib66.1.1 .ltx_text .ltx_font_italic}, pages 7763--7786, Singapore, December 2023. Association for Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib66}
- [[\[67\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Liang Wang, Nan Yang, Xiaolong Huang, Binxing Jiao, Linjun Yang, Daxin Jiang, Rangan Majumder, and Furu Wei. ]{.ltx_bibblock} [Text embeddings by weakly-supervised contrastive pre-training. ]{.ltx_bibblock} [[arXiv preprint arXiv:2212.03533]{#bib.bib67.1.1 .ltx_text .ltx_font_italic}, 2022. ]{.ltx_bibblock}]{#bib.bib67}
- [[\[68\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Shenzhi Wang, Chang Liu, Zilong Zheng, Siyuan Qi, Shuo Chen, Qisen Yang, Andrew Zhao, Chaofei Wang, Shiji Song, and Gao Huang. ]{.ltx_bibblock} [Avalon's game of thoughts: Battle against deception through recursive contemplation. ]{.ltx_bibblock} [[arXiv preprint arXiv:2310.01320]{#bib.bib68.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib68}
- [[\[69\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jason Wei, Xuezhi Wang, Dale Schuurmans, Maarten Bosma, Ed Huai hsin Chi, F. Xia, Quoc Le, and Denny Zhou. ]{.ltx_bibblock} [Chain of thought prompting elicits reasoning in large language models. ]{.ltx_bibblock} [[ArXiv]{#bib.bib69.1.1 .ltx_text .ltx_font_italic}, abs/2201.11903, 2022. ]{.ltx_bibblock}]{#bib.bib69}
- [[\[70\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Martin Weiss, Nasim Rahaman, Manuel Wuthrich, Yoshua Bengio, Li Erran Li, Bernhard Schölkopf, and Christopher Pal. ]{.ltx_bibblock} [Rethinking the buyer's inspection paradox in information markets with language agents, 2023. ]{.ltx_bibblock}]{#bib.bib70}
- [[\[71\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tongshuang Wu, Michael Terry, and Carrie Jun Cai. ]{.ltx_bibblock} [Ai chains: Transparent and controllable human-ai interaction by chaining large language model prompts. ]{.ltx_bibblock} [In [Proceedings of the 2022 CHI Conference on Human Factors in Computing Systems]{#bib.bib71.1.1 .ltx_text .ltx_font_italic}, CHI '22, New York, NY, USA, 2022. Association for Computing Machinery. ]{.ltx_bibblock}]{#bib.bib71}
- [[\[72\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Shitao Xiao, Zheng Liu, Peitian Zhang, and Niklas Muennighoff. ]{.ltx_bibblock} [C-pack: Packaged resources to advance general chinese embedding, 2023. ]{.ltx_bibblock}]{#bib.bib72}
- [[\[73\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Kai Xiong, Xiao Ding, Yixin Cao, Ting Liu, and Bing Qin. ]{.ltx_bibblock} [Examining inter-consistency of large language models collaboration: An in-depth analysis via debate. ]{.ltx_bibblock} [In [The 2023 Conference on Empirical Methods in Natural Language Processing]{#bib.bib73.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib73}
- [[\[74\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Peng Xu, Wei Ping, Xianchao Wu, Lawrence McAfee, Chen Zhu, Zihan Liu, Sandeep Subramanian, Evelina Bakhturina, Mohammad Shoeybi, and Bryan Catanzaro. ]{.ltx_bibblock} [Retrieval meets long context large language models. ]{.ltx_bibblock} [[arXiv preprint arXiv:2310.03025]{#bib.bib74.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib74}
- [[\[75\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yuzhuang Xu, Shuo Wang, Peng Li, Fuwen Luo, Xiaolong Wang, Weidong Liu, and Yang Liu. ]{.ltx_bibblock} [Exploring large language models for communication games: An empirical study on werewolf. ]{.ltx_bibblock} [[arXiv preprint arXiv:2309.04658]{#bib.bib75.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib75}
- [[\[76\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zelai Xu, Chao Yu, Fei Fang, Yu Wang, and Yi Wu. ]{.ltx_bibblock} [Language agents with reinforcement learning for strategic play in the werewolf game. ]{.ltx_bibblock} [[arXiv preprint arXiv:2310.18940]{#bib.bib76.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib76}
- [[\[77\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zhenran Xu, Senbao Shi, Baotian Hu, Jindi Yu, Dongfang Li, Min Zhang, and Yuxiang Wu. ]{.ltx_bibblock} [Towards reasoning in large language models via multi-agent peer review collaboration. ]{.ltx_bibblock} [[arXiv preprint arXiv:2311.08152]{#bib.bib77.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib77}
- [[\[78\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Zhilin Yang, Peng Qi, Saizheng Zhang, Yoshua Bengio, William Cohen, Ruslan Salakhutdinov, and Christopher D. Manning. ]{.ltx_bibblock} [HotpotQA: A dataset for diverse, explainable multi-hop question answering. ]{.ltx_bibblock} [In Ellen Riloff, David Chiang, Julia Hockenmaier, and Jun'ichi Tsujii, editors, [Proceedings of the 2018 Conference on Empirical Methods in Natural Language Processing]{#bib.bib78.1.1 .ltx_text .ltx_font_italic}, pages 2369--2380, Brussels, Belgium, October-November 2018. Association for Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib78}
- [[\[79\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Hongbin Ye, Honghao Gui, Aijia Zhang, Tong Liu, Wei Hua, and Weiqiang Jia. ]{.ltx_bibblock} [Beyond isolation: Multi-agent synergy for improving knowledge graph construction. ]{.ltx_bibblock} [[arXiv preprint arXiv:2312.03022]{#bib.bib79.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib79}
- [[\[80\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Dani Yogatama, Cyprien de Masson d'Autume, and Lingpeng Kong. ]{.ltx_bibblock} [Adaptive semiparametric language models. ]{.ltx_bibblock} [[Transactions of the Association for Computational Linguistics]{#bib.bib80.1.1 .ltx_text .ltx_font_italic}, 9:362--373, 2021. ]{.ltx_bibblock}]{#bib.bib80}
- [[\[81\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ An Zhang, Leheng Sheng, Yuxin Chen, Hao Li, Yang Deng, Xiang Wang, and Tat-Seng Chua. ]{.ltx_bibblock} [On generative agents in recommendation. ]{.ltx_bibblock} [[arXiv preprint arXiv:2310.10108]{#bib.bib81.1.1 .ltx_text .ltx_font_italic}, 2023. ]{.ltx_bibblock}]{#bib.bib81}
- [[\[82\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yusen Zhang, Ansong Ni, Ziming Mao, Chen Henry Wu, Chenguang Zhu, Budhaditya Deb, Ahmed H Awadallah, Dragomir Radev, and Rui Zhang. ]{.ltx_bibblock} [Summ^[n]{#bib.bib82.2.1.1 .ltx_text .ltx_font_italic}^: A multi-stage summarization framework for long input dialogues and documents. ]{.ltx_bibblock} [[arXiv preprint arXiv:2110.10150]{#bib.bib82.3.1 .ltx_text .ltx_font_italic}, 2021. ]{.ltx_bibblock}]{#bib.bib82}
- [[\[83\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yusen Zhang, Ansong Ni, Tao Yu, Rui Zhang, Chenguang Zhu, Budhaditya Deb, Asli Celikyilmaz, Ahmed Hassan Awadallah, and Dragomir Radev. ]{.ltx_bibblock} [An exploratory study on long dialogue summarization: What works and what's next. ]{.ltx_bibblock} [[arXiv preprint arXiv:2109.04609]{#bib.bib83.1.1 .ltx_text .ltx_font_italic}, 2021. ]{.ltx_bibblock}]{#bib.bib83}
- [[\[84\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ming Zhong, Da Yin, Tao Yu, Ahmad Zaidi, Mutethia Mutuma, Rahul Jha, Ahmed Hassan Awadallah, Asli Celikyilmaz, Yang Liu, Xipeng Qiu, and Dragomir Radev. ]{.ltx_bibblock} [QMSum: A new benchmark for query-based multi-domain meeting summarization. ]{.ltx_bibblock} [In Kristina Toutanova, Anna Rumshisky, Luke Zettlemoyer, Dilek Hakkani-Tur, Iz Beltagy, Steven Bethard, Ryan Cotterell, Tanmoy Chakraborty, and Yichao Zhou, editors, [Proceedings of the 2021 Conference of the North American Chapter of the Association for Computational Linguistics: Human Language Technologies]{#bib.bib84.1.1 .ltx_text .ltx_font_italic}, pages 5905--5921, Online, June 2021. Association for Computational Linguistics. ]{.ltx_bibblock}]{#bib.bib84}
- [[\[85\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Denny Zhou, Nathanael Scharli, Le Hou, Jason Wei, Nathan Scales, Xuezhi Wang, Dale Schuurmans, Olivier Bousquet, Quoc Le, and Ed Huai hsin Chi. ]{.ltx_bibblock} [Least-to-most prompting enables complex reasoning in large language models. ]{.ltx_bibblock} [[ArXiv]{#bib.bib85.1.1 .ltx_text .ltx_font_italic}, abs/2205.10625, 2022. ]{.ltx_bibblock}]{#bib.bib85}
:::

::: {.ltx_pagination .ltx_role_newpage}
:::

::::: {#A1 .section .ltx_appendix}
## [Appendix A ]{.ltx_tag .ltx_tag_appendix}Proof of Time Complexity {#appendix-a-proof-of-time-complexity .ltx_title .ltx_title_appendix}

::: {#A1.p1 .ltx_para}
Assuming the source text containing $n$ tokens, window limit of LLM is $k$ tokens, and the responses contain $r$ tokens in average. For decoder-only LLM, we grasp the operations for attention calculation as the time cost unit. Then, for Full-Context LLM, total operation for encoding input source text $T_{\text{Full}}$ is:

  -- ------------------------------------------------------------------------------------------------------------------------------------------ -- ----------------------------------------------------
     $$T_{E\hspace{0pt}n\hspace{0pt}c} = {({1 + 2 + \cdots + n})} = \frac{{({n + 1})}\hspace{0pt}n}{2} = {\mathcal{O}\hspace{0pt}{(n^{2})}}$$      [(3)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
  -- ------------------------------------------------------------------------------------------------------------------------------------------ -- ----------------------------------------------------

Similarly, decoding starts when the model already generate all input. Thus, the first decoded token attends to $n$ positions. Total operation for decoding response is ($r \ll n$):

  -- ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
     $$T_{D\hspace{0pt}e\hspace{0pt}c} = {({n + 1 + n + 2 + \cdots + n + r})} = \frac{{({n + 1 + n + r})}\hspace{0pt}r}{2} = {\mathcal{O}\hspace{0pt}{({{n\hspace{0pt}r} + r^{2}})}} = {\mathcal{O}\hspace{0pt}{({n\hspace{0pt}r})}}$$      [(4)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
  -- ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
:::

::: {#A1.p2 .ltx_para}
For Chain-of-Agents, we first split the source into $\lceil{n/k}\rceil$ chunks. Thus, total encoding time for all input is:

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
     $$T_{E\hspace{0pt}n\hspace{0pt}c} = {{({1 + 2 + \cdots + k})} \times {\lceil{n/k}\rceil}} = \frac{{{({k + 1})}\hspace{0pt}k} \times {\lceil{n/k}\rceil}}{2} = {\mathcal{O}\hspace{0pt}{({{k^{2} \times n}/k})}} = {\mathcal{O}\hspace{0pt}{({n\hspace{0pt}k})}}$$      [(5)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------

Decoding starts when the model already generate $k$ tokens. Thus, the first decoded token attends to $k$ positions. Total operation for decoding response is ($r \ll k$):

  -- ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
     $$T_{D\hspace{0pt}e\hspace{0pt}c} = {{({k + 1 + k + 2 + \cdots + k + r})} \times {\lceil{n/k}\rceil}} = \frac{{{({k + 1 + k + r})}\hspace{0pt}r} \times {\lceil{n/k}\rceil}}{2} = {\mathcal{O}\hspace{0pt}{({{n\hspace{0pt}r} + {{n\hspace{0pt}r^{2}}/k}})}} = {\mathcal{O}\hspace{0pt}{({n\hspace{0pt}r})}}$$      [(6)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
  -- ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
:::
:::::

:::: {#A2 .section .ltx_appendix}
## [Appendix B ]{.ltx_tag .ltx_tag_appendix}Implementation Details {#appendix-b-implementation-details .ltx_title .ltx_title_appendix}

::: {#A2.p1 .ltx_para}
For all experiments, we use Vertex model garden [^2^[[^2^[2]{.ltx_tag .ltx_tag_note}[https://cloud.google.com/model-garden](https://cloud.google.com/model-garden){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}]{.ltx_note_content}]{.ltx_note_outer}]{#footnote2 .ltx_note .ltx_role_footnote} API to use all six models. Maximum generation token is set to 2048 for gemini-ultra and set to 1024 for the rest of the models. We set temperature to 0 for all experiments except for Self-consistency setting. Table [[9]{.ltx_text .ltx_ref_tag}](#A2.T9 "Table 9 ‣ Appendix B Implementation Details ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref} shows the prompt for all models. for task specific requirement of 9 datasets, we follow the original LongBench \[[6](#bib.bib6){.ltx_ref}\] and SCROLLS \[[54](#bib.bib54){.ltx_ref}\]. For RAG model, we use the model provided by Huggingface[^3^[[^3^[3]{.ltx_tag .ltx_tag_note}[https://huggingface.co/](https://huggingface.co/){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}]{.ltx_note_content}]{.ltx_note_outer}]{#footnote3 .ltx_note .ltx_role_footnote} and run on A100 GPUs to rerank the chunks.
:::

<figure id="A2.T9" class="ltx_table">
<table id="A2.T9.10" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A2.T9.2.2" class="ltx_tr">
<td id="A2.T9.2.2.3" class="ltx_td ltx_align_left ltx_border_tt"><span id="A2.T9.2.2.3.1" class="ltx_text" style="font-size:90%;">Vanilla</span></td>
<td id="A2.T9.2.2.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt"><span id="A2.T9.2.2.2.2" class="ltx_inline-block ltx_align_top"> <span id="A2.T9.2.2.2.2.2" class="ltx_p" style="width:369.9pt;"><span id="A2.T9.2.2.2.2.2.3" class="ltx_text"></span><span id="A2.T9.2.2.2.2.2.2" class="ltx_text" style="font-size:90%;"> <span id="A2.T9.2.2.2.2.2.2.2" class="ltx_tabular ltx_align_middle"> <span id="A2.T9.2.2.2.2.2.2.2.3" class="ltx_tr"> <span id="A2.T9.2.2.2.2.2.2.2.3.1" class="ltx_td ltx_nopad_r ltx_align_left">{Task specific requirement}</span></span> <span id="A2.T9.1.1.1.1.1.1.1.1" class="ltx_tr"> <span id="A2.T9.1.1.1.1.1.1.1.1.1" class="ltx_td ltx_nopad_r ltx_align_left">{Source Input <span class="math inline"><em>x</em></span> with truncation if needed}</span></span> <span id="A2.T9.2.2.2.2.2.2.2.2" class="ltx_tr"> <span id="A2.T9.2.2.2.2.2.2.2.2.1" class="ltx_td ltx_nopad_r ltx_align_left">Question: {Question <span class="math inline"><em>q</em></span>}</span></span> <span id="A2.T9.2.2.2.2.2.2.2.4" class="ltx_tr"> <span id="A2.T9.2.2.2.2.2.2.2.4.1" class="ltx_td ltx_nopad_r ltx_align_left">Answer:</span></span> </span></span><span id="A2.T9.2.2.2.2.2.4" class="ltx_text"></span><span id="A2.T9.2.2.2.2.2.5" class="ltx_text" style="font-size:90%;"></span></span> </span></td>
</tr>
<tr id="A2.T9.4.4" class="ltx_tr">
<td id="A2.T9.4.4.3" class="ltx_td ltx_align_left ltx_border_t"><span id="A2.T9.4.4.3.1" class="ltx_text" style="font-size:90%;">RAG</span></td>
<td id="A2.T9.4.4.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A2.T9.4.4.2.2" class="ltx_inline-block ltx_align_top"> <span id="A2.T9.4.4.2.2.2" class="ltx_p" style="width:369.9pt;"><span id="A2.T9.4.4.2.2.2.3" class="ltx_text"></span><span id="A2.T9.4.4.2.2.2.2" class="ltx_text" style="font-size:90%;"> <span id="A2.T9.4.4.2.2.2.2.2" class="ltx_tabular ltx_align_middle"> <span id="A2.T9.4.4.2.2.2.2.2.3" class="ltx_tr"> <span id="A2.T9.4.4.2.2.2.2.2.3.1" class="ltx_td ltx_nopad_r ltx_align_left">{Task specific requirement}</span></span> <span id="A2.T9.3.3.1.1.1.1.1.1" class="ltx_tr"> <span id="A2.T9.3.3.1.1.1.1.1.1.1" class="ltx_td ltx_nopad_r ltx_align_left">{Retrieved Chunks of Source Input <span class="math inline"><em>x</em></span>}</span></span> <span id="A2.T9.4.4.2.2.2.2.2.2" class="ltx_tr"> <span id="A2.T9.4.4.2.2.2.2.2.2.1" class="ltx_td ltx_nopad_r ltx_align_left">Question: {Question <span class="math inline"><em>q</em></span>}</span></span> <span id="A2.T9.4.4.2.2.2.2.2.4" class="ltx_tr"> <span id="A2.T9.4.4.2.2.2.2.2.4.1" class="ltx_td ltx_nopad_r ltx_align_left">Answer:</span></span> </span></span><span id="A2.T9.4.4.2.2.2.4" class="ltx_text"></span><span id="A2.T9.4.4.2.2.2.5" class="ltx_text" style="font-size:90%;"></span></span> </span></td>
</tr>
<tr id="A2.T9.10.10" class="ltx_tr">
<td id="A2.T9.10.10.7" class="ltx_td ltx_align_left ltx_border_bb ltx_border_t"><span id="A2.T9.10.10.7.1" class="ltx_text" style="font-size:90%;">CoA</span></td>
<td id="A2.T9.10.10.6" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb ltx_border_t"><span id="A2.T9.10.10.6.6" class="ltx_inline-block ltx_align_top"> <span id="A2.T9.10.10.6.6.6" class="ltx_p" style="width:369.9pt;"><span id="A2.T9.10.10.6.6.6.7" class="ltx_text"></span><span id="A2.T9.10.10.6.6.6.6" class="ltx_text" style="font-size:90%;"> <span id="A2.T9.10.10.6.6.6.6.6" class="ltx_tabular ltx_align_middle"> <span id="A2.T9.5.5.1.1.1.1.1.1" class="ltx_tr"> <span id="A2.T9.5.5.1.1.1.1.1.1.1" class="ltx_td ltx_nopad_r ltx_align_left"><span id="A2.T9.5.5.1.1.1.1.1.1.1.1" class="ltx_text ltx_font_bold">Worker <span class="math inline"><em>W</em><sub><em>i</em></sub></span>:</span></span></span> <span id="A2.T9.6.6.2.2.2.2.2.2" class="ltx_tr"> <span id="A2.T9.6.6.2.2.2.2.2.2.1" class="ltx_td ltx_nopad_r ltx_align_left">{Input Chunk <span class="math inline"><em>c</em><sub><em>i</em></sub></span>}</span></span> <span id="A2.T9.7.7.3.3.3.3.3.3" class="ltx_tr"> <span id="A2.T9.7.7.3.3.3.3.3.3.1" class="ltx_td ltx_nopad_r ltx_align_left">Here is the summary of the previous source text: {Previous Communication Unit (<span class="math inline"><em>C</em> <em>U</em><sub><em>i</em> − 1</sub></span>)}</span></span> <span id="A2.T9.8.8.4.4.4.4.4.4" class="ltx_tr"> <span id="A2.T9.8.8.4.4.4.4.4.4.1" class="ltx_td ltx_nopad_r ltx_align_left">Question: {Query <span class="math inline"><em>q</em></span>}</span></span> <span id="A2.T9.10.10.6.6.6.6.6.7" class="ltx_tr"> <span id="A2.T9.10.10.6.6.6.6.6.7.1" class="ltx_td ltx_nopad_r ltx_align_left">You need to read current source text and summary of previous source text (if any) and generate a</span></span> <span id="A2.T9.10.10.6.6.6.6.6.8" class="ltx_tr"> <span id="A2.T9.10.10.6.6.6.6.6.8.1" class="ltx_td ltx_nopad_r ltx_align_left">summary to include them both. Later, this summary will be used for other agents to answer the</span></span> <span id="A2.T9.10.10.6.6.6.6.6.9" class="ltx_tr"> <span id="A2.T9.10.10.6.6.6.6.6.9.1" class="ltx_td ltx_nopad_r ltx_align_left">Query, if any. So please write the summary that can include the evidence for answering the Query:</span></span> <span id="A2.T9.9.9.5.5.5.5.5.5" class="ltx_tr"> <span id="A2.T9.9.9.5.5.5.5.5.5.1" class="ltx_td ltx_nopad_r ltx_align_left"><span id="A2.T9.9.9.5.5.5.5.5.5.1.1" class="ltx_text ltx_font_bold">Manager <span class="math inline"><em>M</em></span>:</span></span></span> <span id="A2.T9.10.10.6.6.6.6.6.10" class="ltx_tr"> <span id="A2.T9.10.10.6.6.6.6.6.10.1" class="ltx_td ltx_nopad_r ltx_align_left">{Task specific requirement}</span></span> <span id="A2.T9.10.10.6.6.6.6.6.11" class="ltx_tr"> <span id="A2.T9.10.10.6.6.6.6.6.11.1" class="ltx_td ltx_nopad_r ltx_align_left">The following are given passages. However, the source text is too long and has been summarized. You</span></span> <span id="A2.T9.10.10.6.6.6.6.6.12" class="ltx_tr"> <span id="A2.T9.10.10.6.6.6.6.6.12.1" class="ltx_td ltx_nopad_r ltx_align_left">need to answer based on the summary:</span></span> <span id="A2.T9.10.10.6.6.6.6.6.6" class="ltx_tr"> <span id="A2.T9.10.10.6.6.6.6.6.6.1" class="ltx_td ltx_nopad_r ltx_align_left">{Previous Communication Unit <span class="math inline"><em>C</em> <em>U</em><sub><em>l</em></sub></span>}</span></span> <span id="A2.T9.10.10.6.6.6.6.6.13" class="ltx_tr"> <span id="A2.T9.10.10.6.6.6.6.6.13.1" class="ltx_td ltx_nopad_r ltx_align_left">Question: {question}</span></span> <span id="A2.T9.10.10.6.6.6.6.6.14" class="ltx_tr"> <span id="A2.T9.10.10.6.6.6.6.6.14.1" class="ltx_td ltx_nopad_r ltx_align_left">Answer:</span></span> </span></span><span id="A2.T9.10.10.6.6.6.8" class="ltx_text"></span><span id="A2.T9.10.10.6.6.6.9" class="ltx_text" style="font-size:90%;"></span></span> </span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 9: </span>Prompt of all models for query-based tasks.</figcaption>
</figure>

<figure id="A2.T10" class="ltx_table">
<table id="A2.T10.7" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A2.T10.1.1" class="ltx_tr">
<td id="A2.T10.1.1.2" class="ltx_td ltx_align_left ltx_border_tt"><span id="A2.T10.1.1.2.1" class="ltx_text" style="font-size:90%;">Vanilla</span></td>
<td id="A2.T10.1.1.1" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt"><span id="A2.T10.1.1.1.1" class="ltx_inline-block ltx_align_top"> <span id="A2.T10.1.1.1.1.1" class="ltx_p" style="width:369.9pt;"><span id="A2.T10.1.1.1.1.1.2" class="ltx_text"></span><span id="A2.T10.1.1.1.1.1.1" class="ltx_text" style="font-size:90%;"> <span id="A2.T10.1.1.1.1.1.1.1" class="ltx_tabular ltx_align_middle"> <span id="A2.T10.1.1.1.1.1.1.1.2" class="ltx_tr"> <span id="A2.T10.1.1.1.1.1.1.1.2.1" class="ltx_td ltx_nopad_r ltx_align_left">{Task specific requirement}</span></span> <span id="A2.T10.1.1.1.1.1.1.1.1" class="ltx_tr"> <span id="A2.T10.1.1.1.1.1.1.1.1.1" class="ltx_td ltx_nopad_r ltx_align_left">{Source Input <span class="math inline"><em>x</em></span> with truncation if needed}</span></span> <span id="A2.T10.1.1.1.1.1.1.1.3" class="ltx_tr"> <span id="A2.T10.1.1.1.1.1.1.1.3.1" class="ltx_td ltx_nopad_r ltx_align_left">Answer:</span></span> </span></span><span id="A2.T10.1.1.1.1.1.3" class="ltx_text"></span><span id="A2.T10.1.1.1.1.1.4" class="ltx_text" style="font-size:90%;"></span></span> </span></td>
</tr>
<tr id="A2.T10.2.2" class="ltx_tr">
<td id="A2.T10.2.2.2" class="ltx_td ltx_align_left ltx_border_t"><span id="A2.T10.2.2.2.1" class="ltx_text" style="font-size:90%;">RAG</span></td>
<td id="A2.T10.2.2.1" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A2.T10.2.2.1.1" class="ltx_inline-block ltx_align_top"> <span id="A2.T10.2.2.1.1.1" class="ltx_p" style="width:369.9pt;"><span id="A2.T10.2.2.1.1.1.2" class="ltx_text"></span><span id="A2.T10.2.2.1.1.1.1" class="ltx_text" style="font-size:90%;"> <span id="A2.T10.2.2.1.1.1.1.1" class="ltx_tabular ltx_align_middle"> <span id="A2.T10.2.2.1.1.1.1.1.2" class="ltx_tr"> <span id="A2.T10.2.2.1.1.1.1.1.2.1" class="ltx_td ltx_nopad_r ltx_align_left">{Task specific requirement}</span></span> <span id="A2.T10.2.2.1.1.1.1.1.1" class="ltx_tr"> <span id="A2.T10.2.2.1.1.1.1.1.1.1" class="ltx_td ltx_nopad_r ltx_align_left">{Retrieved Chunks of Source Input <span class="math inline"><em>x</em></span>}</span></span> <span id="A2.T10.2.2.1.1.1.1.1.3" class="ltx_tr"> <span id="A2.T10.2.2.1.1.1.1.1.3.1" class="ltx_td ltx_nopad_r ltx_align_left">Answer:</span></span> </span></span><span id="A2.T10.2.2.1.1.1.3" class="ltx_text"></span><span id="A2.T10.2.2.1.1.1.4" class="ltx_text" style="font-size:90%;"></span></span> </span></td>
</tr>
<tr id="A2.T10.7.7" class="ltx_tr">
<td id="A2.T10.7.7.6" class="ltx_td ltx_align_left ltx_border_bb ltx_border_t"><span id="A2.T10.7.7.6.1" class="ltx_text" style="font-size:90%;">CoA</span></td>
<td id="A2.T10.7.7.5" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb ltx_border_t"><span id="A2.T10.7.7.5.5" class="ltx_inline-block ltx_align_top"> <span id="A2.T10.7.7.5.5.5" class="ltx_p" style="width:369.9pt;"><span id="A2.T10.7.7.5.5.5.6" class="ltx_text"></span><span id="A2.T10.7.7.5.5.5.5" class="ltx_text" style="font-size:90%;"> <span id="A2.T10.7.7.5.5.5.5.5" class="ltx_tabular ltx_align_middle"> <span id="A2.T10.3.3.1.1.1.1.1.1" class="ltx_tr"> <span id="A2.T10.3.3.1.1.1.1.1.1.1" class="ltx_td ltx_nopad_r ltx_align_left"><span id="A2.T10.3.3.1.1.1.1.1.1.1.1" class="ltx_text ltx_font_bold">Worker <span class="math inline"><em>W</em><sub><em>i</em></sub></span>:</span></span></span> <span id="A2.T10.4.4.2.2.2.2.2.2" class="ltx_tr"> <span id="A2.T10.4.4.2.2.2.2.2.2.1" class="ltx_td ltx_nopad_r ltx_align_left">{Input Chunk <span class="math inline"><em>c</em><sub><em>i</em></sub></span>}</span></span> <span id="A2.T10.5.5.3.3.3.3.3.3" class="ltx_tr"> <span id="A2.T10.5.5.3.3.3.3.3.3.1" class="ltx_td ltx_nopad_r ltx_align_left">Here is the summary of the previous source text: {Previous Communication Unit (<span class="math inline"><em>C</em> <em>U</em><sub><em>i</em> − 1</sub></span>)}</span></span> <span id="A2.T10.7.7.5.5.5.5.5.6" class="ltx_tr"> <span id="A2.T10.7.7.5.5.5.5.5.6.1" class="ltx_td ltx_nopad_r ltx_align_left">You need to read the current source text and summary of previous source text (if any) and generate a</span></span> <span id="A2.T10.7.7.5.5.5.5.5.7" class="ltx_tr"> <span id="A2.T10.7.7.5.5.5.5.5.7.1" class="ltx_td ltx_nopad_r ltx_align_left">summary to include them both. Later, this summary will be used for other agents to generate a summary</span></span> <span id="A2.T10.7.7.5.5.5.5.5.8" class="ltx_tr"> <span id="A2.T10.7.7.5.5.5.5.5.8.1" class="ltx_td ltx_nopad_r ltx_align_left">for the whole text. Thus, your generated summary should be relatively long.</span></span> <span id="A2.T10.6.6.4.4.4.4.4.4" class="ltx_tr"> <span id="A2.T10.6.6.4.4.4.4.4.4.1" class="ltx_td ltx_nopad_r ltx_align_left"><span id="A2.T10.6.6.4.4.4.4.4.4.1.1" class="ltx_text ltx_font_bold">Manager <span class="math inline"><em>M</em></span>:</span></span></span> <span id="A2.T10.7.7.5.5.5.5.5.9" class="ltx_tr"> <span id="A2.T10.7.7.5.5.5.5.5.9.1" class="ltx_td ltx_nopad_r ltx_align_left">{Task specific requirement}</span></span> <span id="A2.T10.7.7.5.5.5.5.5.10" class="ltx_tr"> <span id="A2.T10.7.7.5.5.5.5.5.10.1" class="ltx_td ltx_nopad_r ltx_align_left">The following are given passages. However, the source text is too long and has been summarized. You</span></span> <span id="A2.T10.7.7.5.5.5.5.5.11" class="ltx_tr"> <span id="A2.T10.7.7.5.5.5.5.5.11.1" class="ltx_td ltx_nopad_r ltx_align_left">need to answer based on the summary:</span></span> <span id="A2.T10.7.7.5.5.5.5.5.5" class="ltx_tr"> <span id="A2.T10.7.7.5.5.5.5.5.5.1" class="ltx_td ltx_nopad_r ltx_align_left">{Previous Communication Unit <span class="math inline"><em>C</em> <em>U</em><sub><em>l</em></sub></span>}</span></span> <span id="A2.T10.7.7.5.5.5.5.5.12" class="ltx_tr"> <span id="A2.T10.7.7.5.5.5.5.5.12.1" class="ltx_td ltx_nopad_r ltx_align_left">Answer:</span></span> </span></span><span id="A2.T10.7.7.5.5.5.7" class="ltx_text"></span><span id="A2.T10.7.7.5.5.5.8" class="ltx_text" style="font-size:90%;"></span></span> </span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 10: </span>Prompt of all models for non-query-based tasks.</figcaption>
</figure>
::::

:::: {#A3 .section .ltx_appendix}
## [Appendix C ]{.ltx_tag .ltx_tag_appendix}Example Outputs {#appendix-c-example-outputs .ltx_title .ltx_title_appendix}

::: {#A3.p1 .ltx_para}
We pick one sample from HotpotQA in Table [[11]{.ltx_text .ltx_ref_tag}](#A3.T11 "Table 11 ‣ Appendix C Example Outputs ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}, QMSum in Table [[12]{.ltx_text .ltx_ref_tag}](#A3.T12 "Table 12 ‣ Appendix C Example Outputs ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}, and RepoBench-P in Table [[13]{.ltx_text .ltx_ref_tag}](#A3.T13 "Table 13 ‣ Appendix C Example Outputs ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref} to show the prediction of Chain-of-Agents, in the order of question $q$, Communication Unit $C\hspace{0pt}U$ of agents, and final prediction of manager.
:::

<figure id="A3.T11" class="ltx_table">
<table id="A3.T11.4" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A3.T11.4.1" class="ltx_tr">
<td id="A3.T11.4.1.1" class="ltx_td ltx_align_left ltx_border_tt"><span id="A3.T11.4.1.1.1" class="ltx_text" style="font-size:90%;">Question</span></td>
<td id="A3.T11.4.1.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt"><span id="A3.T11.4.1.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T11.4.1.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T11.4.1.2.1.1.1" class="ltx_text" style="font-size:90%;">Gary L. Bennett was a part of the space missions that have a primary destination of what celestial body?</span></span> </span></td>
</tr>
<tr id="A3.T11.4.2" class="ltx_tr">
<td id="A3.T11.4.2.1" class="ltx_td ltx_align_left ltx_border_t"><span id="A3.T11.4.2.1.1" class="ltx_text" style="font-size:90%;">Worker 1</span></td>
<td id="A3.T11.4.2.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A3.T11.4.2.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T11.4.2.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T11.4.2.2.1.1.1" class="ltx_text" style="font-size:90%;">Gary L. Bennett, a scientist and engineer, has contributed to various space missions, including Voyager, Galileo, and Ulysses. He has worked on advanced space power and propulsion systems and has been involved in planetary protection measures. Bennett’s expertise has been instrumental in ensuring the scientific integrity of celestial bodies and preventing harmful contamination. He has received numerous awards and accolades for his contributions to space exploration and is recognized as a leading expert in the field of planetary protection.</span></span> </span></td>
</tr>
<tr id="A3.T11.4.3" class="ltx_tr">
<td id="A3.T11.4.3.1" class="ltx_td ltx_align_left ltx_border_t"><span id="A3.T11.4.3.1.1" class="ltx_text" style="font-size:90%;">Worker 2</span></td>
<td id="A3.T11.4.3.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A3.T11.4.3.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T11.4.3.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T11.4.3.2.1.1.1" class="ltx_text" style="font-size:90%;">Gary L. Bennett, a renowned scientist and engineer, has made significant contributions to space missions, including Voyager, Galileo, and Ulysses. His expertise in advanced space power and propulsion systems, as well as planetary protection measures, has been crucial in ensuring the scientific integrity of celestial bodies and preventing harmful contamination. Bennett has received numerous accolades for his work, including the NASA Exceptional Service Medal and the COSPAR Distinguished Service Award.</span></span> </span></td>
</tr>
<tr id="A3.T11.4.4" class="ltx_tr">
<td id="A3.T11.4.4.1" class="ltx_td ltx_align_left ltx_border_t"><span id="A3.T11.4.4.1.1" class="ltx_text" style="font-size:90%;">Worker 3</span></td>
<td id="A3.T11.4.4.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A3.T11.4.4.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T11.4.4.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T11.4.4.2.1.1.1" class="ltx_text" style="font-size:90%;">Gary L. Bennett, a distinguished scientist and engineer, played a pivotal role in various space missions, particularly Voyager, Galileo, and Ulysses. His expertise in advanced space power and propulsion systems, coupled with his focus on planetary protection measures, has been instrumental in safeguarding the scientific integrity of celestial bodies. Ulysses, launched in 1990, embarked on a unique trajectory to explore both the southern and northern polar regions of the Sun. During its extended mission, Ulysses provided invaluable data on the Sun’s magnetic field, solar wind, and the presence of dust in the Solar System.</span></span> </span></td>
</tr>
<tr id="A3.T11.4.5" class="ltx_tr">
<td id="A3.T11.4.5.1" class="ltx_td ltx_align_left ltx_border_t"><span id="A3.T11.4.5.1.1" class="ltx_text" style="font-size:90%;">Manager</span></td>
<td id="A3.T11.4.5.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A3.T11.4.5.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T11.4.5.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T11.4.5.2.1.1.1" class="ltx_text" style="font-size:90%;">Sun</span></span> </span></td>
</tr>
<tr id="A3.T11.4.6" class="ltx_tr">
<td id="A3.T11.4.6.1" class="ltx_td ltx_align_left ltx_border_t"><span id="A3.T11.4.6.1.1" class="ltx_text" style="font-size:90%;">RAG</span></td>
<td id="A3.T11.4.6.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A3.T11.4.6.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T11.4.6.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T11.4.6.2.1.1.1" class="ltx_text" style="font-size:90%;">Mars</span></span> </span></td>
</tr>
<tr id="A3.T11.4.7" class="ltx_tr">
<td id="A3.T11.4.7.1" class="ltx_td ltx_align_left ltx_border_bb ltx_border_t"><span id="A3.T11.4.7.1.1" class="ltx_text" style="font-size:90%;">Gold</span></td>
<td id="A3.T11.4.7.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb ltx_border_t"><span id="A3.T11.4.7.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T11.4.7.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T11.4.7.2.1.1.1" class="ltx_text" style="font-size:90%;">Sun</span></span> </span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 11: </span>Example output on HotpotQA dataset.</figcaption>
</figure>

<figure id="A3.T12" class="ltx_table">
<table id="A3.T12.4" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A3.T12.4.1" class="ltx_tr">
<td id="A3.T12.4.1.1" class="ltx_td ltx_align_left ltx_border_tt"><span id="A3.T12.4.1.1.1" class="ltx_text" style="font-size:90%;">Question</span></td>
<td id="A3.T12.4.1.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt"><span id="A3.T12.4.1.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T12.4.1.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T12.4.1.2.1.1.1" class="ltx_text" style="font-size:90%;">Summarize the discussion about industrial components.</span></span> </span></td>
</tr>
<tr id="A3.T12.4.2" class="ltx_tr">
<td id="A3.T12.4.2.1" class="ltx_td ltx_align_left ltx_border_t"><span id="A3.T12.4.2.1.1" class="ltx_text" style="font-size:90%;">Worker 1</span></td>
<td id="A3.T12.4.2.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A3.T12.4.2.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T12.4.2.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T12.4.2.2.1.1.1" class="ltx_text" style="font-size:90%;">The industrial designer proposed that the remote control should be made of titanium or rubber, and the buttons should be simple. The marketing expert said that the remote control should be fancy, innovative, and easy to use. They also discussed the trend of making products look like fruit and vegetables.</span></span> </span></td>
</tr>
<tr id="A3.T12.4.3" class="ltx_tr">
<td id="A3.T12.4.3.1" class="ltx_td ltx_align_left ltx_border_t"><span id="A3.T12.4.3.1.1" class="ltx_text" style="font-size:90%;">Worker 2</span></td>
<td id="A3.T12.4.3.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A3.T12.4.3.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T12.4.3.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T12.4.3.2.1.1.1" class="ltx_text" style="font-size:90%;">The industrial designer proposed that the remote control should be made of titanium or rubber, and the buttons should be simple. The marketing expert said that the remote control should be fancy, innovative, and easy to use. They also discussed the trend of making products look like fruit and vegetables. The industrial designer said that they would explore the two options of titanium and rubber, and see if rubber is expensive. The user interface designer said that they would pretend that the modelling clay is titanium and paint it afterwards. The project manager said that they would stick with two batteries and not reinvent the wheel.</span></span> </span></td>
</tr>
<tr id="A3.T12.4.4" class="ltx_tr">
<td id="A3.T12.4.4.1" class="ltx_td ltx_align_left ltx_border_t"><span id="A3.T12.4.4.1.1" class="ltx_text" style="font-size:90%;">Manager</span></td>
<td id="A3.T12.4.4.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A3.T12.4.4.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T12.4.4.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T12.4.4.2.1.1.1" class="ltx_text" style="font-size:90%;">The industrial designer proposed that the remote control should be made of titanium or rubber, and the buttons should be simple. The marketing expert said that the remote control should be fancy, innovative, and easy to use. They also discussed the trend of making products look like fruit and vegetables. The industrial designer said that they would explore the two options of titanium and rubber, and see if rubber is expensive. The user interface designer said that they would pretend that the modelling clay is titanium and paint it afterwards. The project manager said that they would stick with two batteries and not reinvent the wheel. (ROUGE Score: 21.38)</span></span> </span></td>
</tr>
<tr id="A3.T12.4.5" class="ltx_tr">
<td id="A3.T12.4.5.1" class="ltx_td ltx_align_left ltx_border_t"><span id="A3.T12.4.5.1.1" class="ltx_text" style="font-size:90%;">RAG</span></td>
<td id="A3.T12.4.5.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A3.T12.4.5.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T12.4.5.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T12.4.5.2.1.1.1" class="ltx_text" style="font-size:90%;">The industrial designer and the marketing manager discussed the industrial components of the remote control. They agreed that the remote control should have a simple design, with a few buttons and a speech recognition technology. They also discussed the possibility of adding a L_C_D_ display, but decided that it would be too expensive. The designer suggested that the remote control could be made of titanium or rubber, and the marketing manager agreed that a spongy design would be more appealing to consumers. (ROUGE Score: 16.54)</span></span> </span></td>
</tr>
<tr id="A3.T12.4.6" class="ltx_tr">
<td id="A3.T12.4.6.1" class="ltx_td ltx_align_left ltx_border_bb ltx_border_t"><span id="A3.T12.4.6.1.1" class="ltx_text" style="font-size:90%;">Gold</span></td>
<td id="A3.T12.4.6.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb ltx_border_t"><span id="A3.T12.4.6.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T12.4.6.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T12.4.6.2.1.1.1" class="ltx_text" style="font-size:90%;">The industrial designer provided several options respectively for energy, material and interface. Among these options, the industrial designer preferred traditional battery, titanium and simple push buttons so that they would have enough money for speech recognition. After that, the user interface designer proposed an LCD display only for output, which might be taken into account. The group also talked about how to enable the users to find their controllers.</span></span> </span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 12: </span>Example output on QMSum dataset.</figcaption>
</figure>

<figure id="A3.T13" class="ltx_table">
<table id="A3.T13.4" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A3.T13.4.1" class="ltx_tr">
<td id="A3.T13.4.1.1" class="ltx_td ltx_align_left ltx_border_tt"><span id="A3.T13.4.1.1.1" class="ltx_text" style="font-size:90%;">Question</span></td>
<td id="A3.T13.4.1.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt"><span id="A3.T13.4.1.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T13.4.1.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T13.4.1.2.1.1.1" class="ltx_text" style="font-size:90%;">What is the next line of code: … if(type == 5 || type == 0) {</span></span> </span></td>
</tr>
<tr id="A3.T13.4.2" class="ltx_tr">
<td id="A3.T13.4.2.1" class="ltx_td ltx_align_left ltx_border_t"><span id="A3.T13.4.2.1.1" class="ltx_text" style="font-size:90%;">Worker 1</span></td>
<td id="A3.T13.4.2.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A3.T13.4.2.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T13.4.2.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T13.4.2.2.1.1.1" class="ltx_text" style="font-size:90%;">The method SensorDataCollectorManager.flushSensorDataCache(int type, String deviceID) flushes the database cache for the given sensor type. If the type is 0, all sensor types are flushed. The method first checks if the type is valid, and if it is, it calls the flushDBCache method for the appropriate sensor collector. SensorDataCollectorManager.flushSensorDataCache(int type, String deviceID) flushes the database cache for the given sensor type. If the type is 0, all sensor types are flushed. The method first checks if the type is valid, and if it is, it calls the flushDBCache method for the appropriate sensor collector.</span></span> </span></td>
</tr>
<tr id="A3.T13.4.3" class="ltx_tr">
<td id="A3.T13.4.3.1" class="ltx_td ltx_align_left ltx_border_t"><span id="A3.T13.4.3.1.1" class="ltx_text" style="font-size:90%;">Worker 2</span></td>
<td id="A3.T13.4.3.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A3.T13.4.3.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T13.4.3.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T13.4.3.2.1.1.1" class="ltx_text" style="font-size:90%;">The method SensorDataCollectorManager.flushSensorDataCache(int type, String deviceID) flushes the database cache for the given sensor type. If the type is 0, all sensor types are flushed. The method first checks if the type is valid, and if it is, it calls the flushDBCache method for the appropriate sensor collector.</span></span> </span></td>
</tr>
<tr id="A3.T13.4.4" class="ltx_tr">
<td id="A3.T13.4.4.1" class="ltx_td ltx_align_left ltx_border_t"><span id="A3.T13.4.4.1.1" class="ltx_text" style="font-size:90%;">Manager</span></td>
<td id="A3.T13.4.4.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A3.T13.4.4.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T13.4.4.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T13.4.4.2.1.1.1" class="ltx_text" style="font-size:90%;">LightSensorCollector.flushDBCache(deviceID);</span></span> </span></td>
</tr>
<tr id="A3.T13.4.5" class="ltx_tr">
<td id="A3.T13.4.5.1" class="ltx_td ltx_align_left ltx_border_t"><span id="A3.T13.4.5.1.1" class="ltx_text" style="font-size:90%;">RAG</span></td>
<td id="A3.T13.4.5.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t"><span id="A3.T13.4.5.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T13.4.5.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T13.4.5.2.1.1.1" class="ltx_text" style="font-size:90%;">AccelerometerSensorCollector.flushDBCache(deviceID);</span></span> </span></td>
</tr>
<tr id="A3.T13.4.6" class="ltx_tr">
<td id="A3.T13.4.6.1" class="ltx_td ltx_align_left ltx_border_bb ltx_border_t"><span id="A3.T13.4.6.1.1" class="ltx_text" style="font-size:90%;">Gold</span></td>
<td id="A3.T13.4.6.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb ltx_border_t"><span id="A3.T13.4.6.2.1" class="ltx_inline-block ltx_align_top"> <span id="A3.T13.4.6.2.1.1" class="ltx_p" style="width:369.9pt;"><span id="A3.T13.4.6.2.1.1.1" class="ltx_text" style="font-size:90%;">LightSensorCollector.flushDBCache(deviceID);</span></span> </span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 13: </span>Example output on RepoBench-P dataset.</figcaption>
</figure>
::::

:::: {#A4 .section .ltx_appendix}
## [Appendix D ]{.ltx_tag .ltx_tag_appendix}Robustness against Context Window Size {#appendix-d-robustness-against-context-window-size .ltx_title .ltx_title_appendix}

<figure id="A4.F6" class="ltx_figure">
<img src="/html/2406.02818/assets/figures/contextwindow.png" id="A4.F6.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="299" height="169" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure"><span id="A4.F6.2.1.1" class="ltx_text" style="font-size:90%;">Figure 6</span>: </span><span id="A4.F6.3.2" class="ltx_text" style="font-size:90%;">Performance of CoA on Claude 3 Haiku on the NarrativeQA dataset with various context window sizes of an agent. Results show the robustness of CoA towards different choices of context lengths.</span></figcaption>
</figure>

::: {#A4.p1 .ltx_para}
We set the default context window of CoA to 8k due to the limitation of text-bison and unicorn models. To test the influence of CoA against context window change, we set window size to 4k, 8k, 16k, 32k, 64k, and 128k of Claude 3 Haiku model and evaluate on NarrativeQA dataset and see the performance change. As shown in Figure [[6]{.ltx_text .ltx_ref_tag}](#A4.F6 "Figure 6 ‣ Appendix D Robustness against Context Window Size ‣ Chain of Agents: Large Language Models Collaborating on Long-Context Tasks"){.ltx_ref}, the performance of the model increases from 4k to 16k and stabilize to around 20 with context window further increases. This result shows that CoA will benefit from increasing length and keep stable when the length touches a bound.
:::
::::

:::: {#A5 .section .ltx_appendix}
## [Appendix E ]{.ltx_tag .ltx_tag_appendix}Broader Impacts {#appendix-e-broader-impacts .ltx_title .ltx_title_appendix}

::: {#A5.p1 .ltx_para}
Chain-of-Agents is a generic framework for long context tasks. users can apply this to diverse tasks not restricting to the mentioned ones. It will greatly increase the efficiency of individuals or companies to solve complex long context tasks. Besides, the interpretablity of such approach can reduce the misuse of the LLMs because users can check the correctness of results and decrease the possibility of making faults. However, similar to all prompt based approaches, this framework requires careful prompt design for unseen large language models, users may not get optimal solution on certain newly proposed LLMs. Besides, it may increase the number of the calls for API, causing higher network traffic and higher latency for user pools.
:::
::::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

::: ar5iv-footer
[◄](/html/2406.02817){.ar5iv-nav-button .ar5iv-nav-button-prev} [![ar5iv homepage](/assets/ar5iv.png){height="40"}](/){.ar5iv-home-button} [Feeling\
lucky?](/feeling_lucky){.ar5iv-text-button} [](/land_of_honey_and_milk){rel="nofollow" aria-hidden="true" tabindex="-1"} [Conversion\
report](/log/2406.02818){.ar5iv-text-button .ar5iv-severity-ok} [Report\
an issue](https://github.com/dginev/ar5iv/issues/new?template=improve-article--arxiv-id-.md&title=Improve+article+2406.02818){.ar5iv-text-button target="_blank"} [View original\
on arXiv](https://arxiv.org/abs/2406.02818){.ar5iv-text-button .arxiv-ui-theme}[►](/html/2406.02819){.ar5iv-nav-button .ar5iv-nav-button-next}
:::

[[]{.color-scheme-icon}](javascript:toggleColorScheme() "Toggle ar5iv color scheme"){.ar5iv-toggle-color-scheme} [Copyright](https://arxiv.org/help/license){.ar5iv-footer-button target="_blank"} [Privacy Policy](https://arxiv.org/help/policies/privacy_policy){.ar5iv-footer-button target="_blank"}

::: ltx_page_logo
Generated on Fri Jul 5 18:19:43 2024 by [[L[a]{.ltx_font_smallcaps style="position:relative; bottom:2.2pt;"}T[e]{.ltx_font_smallcaps style="font-size:120%;position:relative; bottom:-0.2ex;"}]{style="letter-spacing:-0.2em; margin-right:0.1em;"}[XML]{style="font-size:90%; position:relative; bottom:-0.2ex;"}![Mascot Sammy](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAsAAAAOCAYAAAD5YeaVAAAAAXNSR0IArs4c6QAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAALEwAACxMBAJqcGAAAAAd0SU1FB9wKExQZLWTEaOUAAAAddEVYdENvbW1lbnQAQ3JlYXRlZCB3aXRoIFRoZSBHSU1Q72QlbgAAAdpJREFUKM9tkL+L2nAARz9fPZNCKFapUn8kyI0e4iRHSR1Kb8ng0lJw6FYHFwv2LwhOpcWxTjeUunYqOmqd6hEoRDhtDWdA8ApRYsSUCDHNt5ul13vz4w0vWCgUnnEc975arX6ORqN3VqtVZbfbTQC4uEHANM3jSqXymFI6yWazP2KxWAXAL9zCUa1Wy2tXVxheKA9YNoR8Pt+aTqe4FVVVvz05O6MBhqUIBGk8Hn8HAOVy+T+XLJfLS4ZhTiRJgqIoVBRFIoric47jPnmeB1mW/9rr9ZpSSn3Lsmir1fJZlqWlUonKsvwWwD8ymc/nXwVBeLjf7xEKhdBut9Hr9WgmkyGEkJwsy5eHG5vN5g0AKIoCAEgkEkin0wQAfN9/cXPdheu6P33fBwB4ngcAcByHJpPJl+fn54mD3Gg0NrquXxeLRQAAwzAYj8cwTZPwPH9/sVg8PXweDAauqqr2cDjEer1GJBLBZDJBs9mE4zjwfZ85lAGg2+06hmGgXq+j3+/DsixYlgVN03a9Xu8jgCNCyIegIAgx13Vfd7vdu+FweG8YRkjXdWy329+dTgeSJD3ieZ7RNO0VAXAPwDEAO5VKndi2fWrb9jWl9Esul6PZbDY9Go1OZ7PZ9z/lyuD3OozU2wAAAABJRU5ErkJggg==)](http://dlmf.nist.gov/LaTeXML/){.ltx_LaTeXML_logo target="_blank"}
:::
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
