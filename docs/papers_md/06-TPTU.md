:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_main
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_content
# TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage {#tptu-large-language-model-based-ai-agents-for-task-planning-and-tool-usage .ltx_title .ltx_title_document}

::: ltx_authors
[ [Jingqing Ruan^†‡^\
ruanjingqing@sensetime.com &Yihong Chen^†‡^\
chenyihong@sensetime.com &Bin Zhang^†‡^\
zhangbin11@sensetime.com &Zhiwei Xu^†‡^\
xuzhiwei@sensetime.com &Tianpeng Bao^†^\
baotianpeng@sensetime.com &Guoqing Du^†^\
duguoqing@sensetime.com &Shiwei Shi^†^\
shishiwei@sensetime.com &Hangyu Mao^†∗^\
maohangyu@sensetime.com &Ziyue Li ^+^\
zlibn@connect.ust.hk &Xingyu Zeng\
zengxingyu@sensetime.com &Rui Zhao\
zhaorui@sensetime.com\
\
SenseTime Research ]{.ltx_personname}]{.ltx_creator .ltx_role_author}
:::

::: ltx_abstract
###### Abstract {#abstract .ltx_title .ltx_title_abstract}

With recent advancements in natural language processing, Large Language Models (LLMs) have emerged as powerful tools for various real-world applications. Despite their powers, the intrinsic generative abilities of LLMs may prove insufficient for handling complex tasks, which necessitate a combination of task planning and the usage of external tools. In this paper, we first propose a structured framework tailored for LLM-based AI Agents and then discuss the crucial capabilities necessary for tackling intricate problems. Within this framework, we design two distinct types of agents (i.e., one-step agent and sequential agent) to execute the inference process. Subsequently, we instantiate the framework using various LLMs and evaluate their [T]{#id29.id1.1 .ltx_text .ltx_framed .ltx_framed_underline}ask [P]{#id29.id1.2 .ltx_text .ltx_framed .ltx_framed_underline}lanning and [T]{#id29.id1.3 .ltx_text .ltx_framed .ltx_framed_underline}ool [U]{#id29.id1.4 .ltx_text .ltx_framed .ltx_framed_underline}sage (TPTU) abilities on typical tasks. By highlighting key findings and challenges, our goal is to provide a helpful resource for researchers and practitioners to leverage the power of LLMs in their AI applications. Our study emphasizes the substantial potential of these models while also identifying areas that need more investigation and improvement. [The code and resources will be available on GitHub.]{#id29.id1.5 .ltx_text style="color:#000000;"}
:::

[^\${}\^{\\dagger}\$^[[^\${}\^{\\dagger}\$^[footnotetext: ]{.ltx_note_type}These authors contribute equally to this work.]{.ltx_note_content}]{.ltx_note_outer}]{#footnotex1 .ltx_note .ltx_role_footnotetext}[^\${}\^{+}\$^[[^\${}\^{+}\$^[footnotetext: ]{.ltx_note_type}External discussion and ideation.]{.ltx_note_content}]{.ltx_note_outer}]{#footnotex2 .ltx_note .ltx_role_footnotetext}[^\${}\^{\\ddagger}\$^[[^\${}\^{\\ddagger}\$^[footnotetext: ]{.ltx_note_type}These authors work as research interns at SenseTime Research.]{.ltx_note_content}]{.ltx_note_outer}]{#footnotex3 .ltx_note .ltx_role_footnotetext}[^\${}\^{\*}\$^[[^\${}\^{\*}\$^[footnotetext: ]{.ltx_note_type}The corresponding author.]{.ltx_note_content}]{.ltx_note_outer}]{#footnotex4 .ltx_note .ltx_role_footnotetext}

::::::: {#S1 .section .ltx_section}
## [1 ]{.ltx_tag .ltx_tag_section}Introduction {#introduction .ltx_title .ltx_title_section}

::: {#S1.p1 .ltx_para}
Large Language Model (LLM) \[[1](#bib.bib1){.ltx_ref}\] is a recent breakthrough in natural language processing (NLP) research. These models are trained on massive amounts of text data and can solve a wide range of tasks, even those that were not included in their training dataset, [known as "emerging" ability]{#S1.p1.1.1 .ltx_text style="color:#000000;"}. This ability is especially evident in [the tasks of]{#S1.p1.1.2 .ltx_text style="color:#000000;"} few-shot \[[2](#bib.bib2){.ltx_ref}\] and zero-shot \[[3](#bib.bib3){.ltx_ref}\] learning, where LLMs can perform well with minimal or even no [fine-tuning to adapt to a new task]{#S1.p1.1.3 .ltx_text style="color:#000000;"}.
:::

<figure id="S1.F1" class="ltx_figure">
<img src="/html/2308.03427/assets/x1.png" id="S1.F1.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="230" height="165" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 1: </span>Our LLM-based agents plan tasks and use tools.</figcaption>
</figure>

::: {#S1.p2 .ltx_para}
However, the application of LLMs in real-world settings presents unique challenges. On the one hand, LLMs have proved to be [incompetent in]{#S1.p2.1.1 .ltx_text style="color:#000000;"} solving logic problems such as mathematics, and their training data is also out of date [(e.g., the knowledge cutoff date for GPT-4 \[[4](#bib.bib4){.ltx_ref}\] is up to January 2022)]{#S1.p2.1.2 .ltx_text style="color:#000000;"}. Teaching LLMs to use tools such as calculators, [calendar]{#S1.p2.1.3 .ltx_text style="color:#000000;"}, or search engines can help prevent them from hallucinating \[[5](#bib.bib5){.ltx_ref}\]. On the other hand, despite their impressive problem-solving abilities, the successful integration of these models into complex systems often requires more than just task understanding - it requires the capacity to manipulate various tools and interact effectively with users. This is exemplified in systems like AutoGPT [^1^[[^1^[1]{.ltx_tag .ltx_tag_note}[https://github.com/Significant-Gravitas/Auto-GPT](https://github.com/Significant-Gravitas/Auto-GPT){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}]{.ltx_note_content}]{.ltx_note_outer}]{#footnote1 .ltx_note .ltx_role_footnote}, BabyAGI [^2^[[^2^[2]{.ltx_tag .ltx_tag_note}[https://github.com/yoheinakajima/babyagi](https://github.com/yoheinakajima/babyagi){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}]{.ltx_note_content}]{.ltx_note_outer}]{#footnote2 .ltx_note .ltx_role_footnote}, and ChatGPT-plugins [^3^[[^3^[3]{.ltx_tag .ltx_tag_note}[https://openai.com/blog/chatgpt-plugins](https://openai.com/blog/chatgpt-plugins){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}]{.ltx_note_content}]{.ltx_note_outer}]{#footnote3 .ltx_note .ltx_role_footnote}, which leverage LLMs' capabilities beyond merely generating well-written texts and programs. In these systems, LLMs operate as the central controller, manipulating different tools and interacting with humans, thus taking on the role of Artificial Intelligence Agents (AI Agents). In addition to being central planners, LLMs are often used as intermediaries between macro plans and low-level tool calls or as specific tools. As such, LLMs are seen as a crucial approximation of the linguistic world model in real-world systems.
:::

::: {#S1.p3 .ltx_para}
In this paper, we propose a structured framework [for LLM-based AI Agents to evaluate the existing LLMs' planning and tool-using ability]{#S1.p3.1.1 .ltx_text style="color:#000000;"} and discuss the necessary abilities of such LLM-based AI Agents. Furthermore, we instantiate the framework with different LLMs and evaluate their [T]{#S1.p3.1.2 .ltx_text .ltx_font_bold}ask [P]{#S1.p3.1.3 .ltx_text .ltx_font_bold}lanning and [T]{#S1.p3.1.4 .ltx_text .ltx_font_bold}ool [U]{#S1.p3.1.5 .ltx_text .ltx_font_bold}sage (TPTU) abilities on several tasks. [As shown in Figure [[1]{.ltx_text .ltx_ref_tag}](#S1.F1 "Figure 1 ‣ 1 Introduction ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, we use the Doraemon as an analogy of our LLM-based agents: Doraemon's magic 4D pocket consists of millions of gadgets (the Tool Set), and Doraemon needs to pick the right tools and solve tasks in a right order.]{#S1.p3.1.6 .ltx_text style="color:#000000;"} Our main contributions are summarized as follows:
:::

::: {#S1.p4 .ltx_para}
1.  [[1.]{.ltx_tag .ltx_tag_item}]{#S1.I1.i1}

    ::: {#S1.I1.i1.p1 .ltx_para}
    We propose a structured framework tailored for LLM-based AI Agents to evaluate the TPTU abilities of the existing open-source LLMs.
    :::
2.  [[2.]{.ltx_tag .ltx_tag_item}]{#S1.I1.i2}

    ::: {#S1.I1.i2.p1 .ltx_para}
    We design two distinct types of agents, namely, one-step agent and sequential agent, to execute the inference process of [conducting sub-tasks in a once-for-all or sequential manner, respectively]{#S1.I1.i2.p1.1.1 .ltx_text style="color:#000000;"}. We provide detailed empirical results and analysis.
    :::
3.  [[3.]{.ltx_tag .ltx_tag_item}]{#S1.I1.i3}

    ::: {#S1.I1.i3.p1 .ltx_para}
    Our study reveals significant potential in utilizing LLMs for complex tasks. [Furthermore, we conclude four following potential weaknesses of LLM-based agents: failing to output in a specific format, struggling to grasp task requirements, over-utilizing one tool, and lack of summary skills. These observations could spark some insights and shed light on the areas that deserve further investigation and improvement.]{#S1.I1.i3.p1.1.1 .ltx_text style="color:#000000;"}
    :::
:::
:::::::

::::::::::::::: {#S2 .section .ltx_section}
## [2 ]{.ltx_tag .ltx_tag_section}Method {#method .ltx_title .ltx_title_section}

::: {#S2.p1 .ltx_para}
To the best of our knowledge, the study of "Agent", "Autonomous Agent", "AI Agent\" and "Multi-Agent" has been a central part of AI research for decades \[[6](#bib.bib6){.ltx_ref}, [7](#bib.bib7){.ltx_ref}, [8](#bib.bib8){.ltx_ref}, [9](#bib.bib9){.ltx_ref}, [10](#bib.bib10){.ltx_ref}, [11](#bib.bib11){.ltx_ref}\], aimed at understanding and building intelligent and autonomous systems, but there is currently no standardized definition for AI Agents, particularly those that are based on LLMs.
:::

::: {#S2.p2 .ltx_para}
[In this paper, the Artificial Intelligence Agent (AI Agent) is defined as a program that employs AI techniques to perform tasks that typically require human-like intelligence]{#S2.p2.1.1 .ltx_text .ltx_font_bold}. AI Agents can take many forms, from simple chatbots to complex autonomous systems that interact with their environment and make decisions in real-time. They can be trained using a variety of machine learning techniques, including supervised, unsupervised, and reinforcement learning, and can be programmed to perform specific tasks or learn from their experiences in order to improve their performance over time.
:::

::::: {#S2.SS1 .section .ltx_subsection}
### [2.1 ]{.ltx_tag .ltx_tag_subsection}Agent Framework {#agent-framework .ltx_title .ltx_title_subsection}

<figure id="S2.F2" class="ltx_figure">
<img src="/html/2308.03427/assets/x2.png" id="S2.F2.g1" class="ltx_graphics ltx_centering ltx_img_portrait" width="368" height="458" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 2: </span>The proposed framework for LLM-based AI Agents.</figcaption>
</figure>

::: {#S2.SS1.p1 .ltx_para}
We are particularly interested in the AI Agent that employs the LLM techniques (i.e., LLM-based AI Agent), due to its high efficiency and flexibility in various tasks and domains. Specifically, we design our AI Agent framework with six components as shown in Figure [[2]{.ltx_text .ltx_ref_tag}](#S2.F2 "Figure 2 ‣ 2.1 Agent Framework ‣ 2 Method ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}:
:::

::: {#S2.SS1.p2 .ltx_para}
1.  [[1.]{.ltx_tag .ltx_tag_item}]{#S2.I1.i1}

    ::: {#S2.I1.i1.p1 .ltx_para}
    [Task Instruction]{#S2.I1.i1.p1.1.1 .ltx_text .ltx_font_bold}. This is the explicit input of the agent. In practical systems, the task instruction comes from human users of the systems. For example, in a human resources (HR) system, the user may give a task instruction: How much budget is required to provide a 100\$ incentive for each colleague who has worked for five years? In contrast, in a criminal investigation system, the user may give a task instruction: Deploy surveillance on a group of suspects.
    :::
2.  [[2.]{.ltx_tag .ltx_tag_item}]{#S2.I1.i2}

    ::: {#S2.I1.i2.p1 .ltx_para}
    [Designed Prompt]{#S2.I1.i2.p1.1.1 .ltx_text .ltx_font_bold}. This is an additional form of input for the agent, derived from tasks that the human users anticipate the AI Agent will complete. Humans can craft specific instructions or demonstrations to steer the LLM-based AI Agents toward generating suitable responses. These guiding inputs could encompass system instructions, tool descriptions, few-shot demonstrations, chat history, or even error output.
    :::
3.  [[3.]{.ltx_tag .ltx_tag_item}]{#S2.I1.i3}

    ::: {#S2.I1.i3.p1 .ltx_para}
    [Tool Set]{#S2.I1.i3.p1.1.1 .ltx_text .ltx_font_bold}. It is another input for the agent, which refers to the set of external resources, services, or subsystems that the AI Agent can utilize to aid in its tasks. This could include databases for information retrieval [\[]{#S2.I1.i3.p1.1.2.1 .ltx_text style="color:#000000;"}[12](#bib.bib12){.ltx_ref}[\]]{#S2.I1.i3.p1.1.3.2 .ltx_text style="color:#000000;"}, APIs for interacting with external systems [\[]{#S2.I1.i3.p1.1.4.1 .ltx_text style="color:#000000;"}[5](#bib.bib5){.ltx_ref}[\]]{#S2.I1.i3.p1.1.5.2 .ltx_text style="color:#000000;"}, other AI models specialized for tasks such as image recognition or sentiment analysis [\[]{#S2.I1.i3.p1.1.6.1 .ltx_text style="color:#000000;"}[13](#bib.bib13){.ltx_ref}[\]]{#S2.I1.i3.p1.1.7.2 .ltx_text style="color:#000000;"}, or even non-AI tools and resources such as web scraping tools or data visualization libraries [\[]{#S2.I1.i3.p1.1.8.1 .ltx_text style="color:#000000;"}[14](#bib.bib14){.ltx_ref}[\]]{#S2.I1.i3.p1.1.9.2 .ltx_text style="color:#000000;"}. The toolset expands the capabilities of the AI Agent, enabling it to access and process information beyond its internal knowledge, interact with other systems, or perform specialized tasks that it may not be capable of on its own. For example, an AI Agent might use a weather API to fetch current weather information, or a Python interpreter to solve the mathematical question.
    :::
4.  [[4.]{.ltx_tag .ltx_tag_item}]{#S2.I1.i4}

    ::: {#S2.I1.i4.p1 .ltx_para}
    [LLM]{#S2.I1.i4.p1.1.1 .ltx_text .ltx_font_bold}. This is the core component of the system that interprets the task instructions and prompts, interacts with the toolset, and generates the intermediate outputs and final answers. In this context, we utilize publicly available large language models such as ChatGPT, GPT-4 [\[]{#S2.I1.i4.p1.1.2.1 .ltx_text style="color:#000000;"}[4](#bib.bib4){.ltx_ref}[\]]{#S2.I1.i4.p1.1.3.2 .ltx_text style="color:#000000;"}, [InterLM \[[15](#bib.bib15){.ltx_ref}\]]{#S2.I1.i4.p1.1.4 .ltx_text style="color:#000000;"}, and others.
    :::
5.  [[5.]{.ltx_tag .ltx_tag_item}]{#S2.I1.i5}

    ::: {#S2.I1.i5.p1 .ltx_para}
    [Intermediate Output]{#S2.I1.i5.p1.1.1 .ltx_text .ltx_font_bold}. This represents the output generated by the LLM-based AI Agent after it processes the task instructions and prompts, and interacts with the toolset. There are three typical intermediate outputs: (1) the high-level plans to fulfill the original user instruction, (2) selected and created tools to fulfill each subtask in the plans, and (3) the results or errors produced after tool execution. The output can be reviewed and refined, either by the AI Agent itself or with human oversight, to ensure it is accurate and meets the requirements of the task instruction.
    :::
6.  [[6.]{.ltx_tag .ltx_tag_item}]{#S2.I1.i6}

    ::: {#S2.I1.i6.p1 .ltx_para}
    [Final Answer]{#S2.I1.i6.p1.1.1 .ltx_text .ltx_font_bold}. This is the output that the AI Agent summarizes and provides to the user after all processing (including task planning, tool usage, and possibly error feedback) has been completed.
    :::
:::
:::::

::::: {#S2.SS2 .section .ltx_subsection}
### [2.2 ]{.ltx_tag .ltx_tag_subsection}Agent Ability {#agent-ability .ltx_title .ltx_title_subsection}

::: {#S2.SS2.p1 .ltx_para}
To apply LLM-based AI Agents to augment or replace human decision-making in real-world applications, the agents typically require the following abilities:

1.  [[1.]{.ltx_tag .ltx_tag_item}]{#S2.I2.i1}

    ::: {#S2.I2.i1.p1 .ltx_para}
    [Perception Ability]{#S2.I2.i1.p1.1.1 .ltx_text .ltx_font_bold}: AI Agents must be able to perceive the task instruction from human and system specifications.
    :::
2.  [[2.]{.ltx_tag .ltx_tag_item}]{#S2.I2.i2}

    ::: {#S2.I2.i2.p1 .ltx_para}
    [Task Planing Ability]{#S2.I2.i2.p1.1.1 .ltx_text .ltx_font_bold}: AI Agents should have the capacity to create a step-by-step plan for complex task composition based on the perceived instruction and specifications. This usually involves the generation of critical subtask sequences, and the ability to adjust the plan dynamically in response to changes in the task or environment.
    :::
3.  [[3.]{.ltx_tag .ltx_tag_item}]{#S2.I2.i3}

    ::: {#S2.I2.i3.p1 .ltx_para}
    [Tool Usage Ability]{#S2.I2.i3.p1.1.1 .ltx_text .ltx_font_bold}: On the one hand, AI Agents should possess the capacity to [select]{#S2.I2.i3.p1.1.2 .ltx_text .ltx_font_bold} a variety of existing tools or resources to execute complex tasks. On the other hand, AI Agents should [create]{#S2.I2.i3.p1.1.3 .ltx_text .ltx_font_bold} new tools by converting the task requirements. This ability enables the AI Agent to extend its capabilities beyond LLM itself and the existing tools by leveraging the vast resources available in the digital world. Finally, AI Agents should be able to [execute]{#S2.I2.i3.p1.1.4 .ltx_text .ltx_font_bold} the selected or created tools for truly grounding the human request based on the resources and constraints of systems.
    :::
4.  [[4.]{.ltx_tag .ltx_tag_item}]{#S2.I2.i4}

    ::: {#S2.I2.i4.p1 .ltx_para}
    [Learning/Reflection/Memory (from Feedback)]{#S2.I2.i4.p1.1.1 .ltx_text .ltx_font_bold}: AI Agents should be capable of learning from feedback, including correct results and exception errors. They should incorporate memory, such as logging or chat history, and reflection to adapt their plans or decisions. This allows the agents to improve their performance and efficiency in task execution continuously.
    :::
5.  [[5.]{.ltx_tag .ltx_tag_item}]{#S2.I2.i5}

    ::: {#S2.I2.i5.p1 .ltx_para}
    [Summarization]{#S2.I2.i5.p1.1.1 .ltx_text .ltx_font_bold}: After several rounds of interaction with humans, tools, and systems, AI agents can ultimately complete the original task provided by the users. At this point, AI agents should be able to summarize the interaction history and provide a final answer that is concise and easy to understand for the users.
    :::
:::

::: {#S2.SS2.p2 .ltx_para}
To endow AI Agents with the aforementioned abilities, some techniques that can be used include chain-of-thought (CoT) and vector databases, as shown in Table [[1]{.ltx_text .ltx_ref_tag}](#S2.T1 "Table 1 ‣ 2.2 Agent Ability ‣ 2 Method ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::

<figure id="S2.T1" class="ltx_table">
<table id="S2.T1.3" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="S2.T1.3.1" class="ltx_tr">
<td id="S2.T1.3.1.1" class="ltx_td ltx_align_left ltx_border_tt" style="padding-top: 1.15pt; padding-bottom: 1.15pt"><span id="S2.T1.3.1.1.1" class="ltx_text ltx_font_bold" style="font-size:90%;">Ability</span></td>
<td id="S2.T1.3.1.2" class="ltx_td ltx_align_left ltx_border_tt" style="padding-top: 1.15pt; padding-bottom: 1.15pt"><span id="S2.T1.3.1.2.1" class="ltx_text ltx_font_bold" style="font-size:90%;">Possible Techniques</span></td>
</tr>
<tr id="S2.T1.3.2" class="ltx_tr">
<td id="S2.T1.3.2.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.15pt; padding-bottom: 1.15pt"><span id="S2.T1.3.2.1.1" class="ltx_text" style="font-size:90%;">Perception</span></td>
<td id="S2.T1.3.2.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.15pt; padding-bottom: 1.15pt"><span id="S2.T1.3.2.2.1" class="ltx_text" style="font-size:90%;">Multi-input Fusion</span></td>
</tr>
<tr id="S2.T1.3.3" class="ltx_tr">
<td id="S2.T1.3.3.1" class="ltx_td ltx_align_left" style="padding-top: 1.15pt; padding-bottom: 1.15pt"><span id="S2.T1.3.3.1.1" class="ltx_text" style="font-size:90%;">Task Planing</span></td>
<td id="S2.T1.3.3.2" class="ltx_td ltx_align_left" style="padding-top: 1.15pt; padding-bottom: 1.15pt"><span id="S2.T1.3.3.2.1" class="ltx_text" style="font-size:90%;">Zero-shot CoT and Few-shot CoT</span></td>
</tr>
<tr id="S2.T1.3.4" class="ltx_tr">
<td id="S2.T1.3.4.1" class="ltx_td ltx_align_left" style="padding-top: 1.15pt; padding-bottom: 1.15pt"><span id="S2.T1.3.4.1.1" class="ltx_text"></span><span id="S2.T1.3.4.1.2" class="ltx_text" style="font-size:90%;"> <span id="S2.T1.3.4.1.2.1" class="ltx_tabular ltx_align_middle"> <span id="S2.T1.3.4.1.2.1.1" class="ltx_tr"> <span id="S2.T1.3.4.1.2.1.1.1" class="ltx_td ltx_nopad_r ltx_align_left" style="padding-top:1.15pt;padding-bottom:1.15pt;">Tool Usage</span></span> <span id="S2.T1.3.4.1.2.1.2" class="ltx_tr"> <span id="S2.T1.3.4.1.2.1.2.1" class="ltx_td ltx_nopad_r ltx_align_left" style="padding-top:1.15pt;padding-bottom:1.15pt;">(Selection/Creation/Execution)</span></span> </span></span><span id="S2.T1.3.4.1.3" class="ltx_text"></span><span id="S2.T1.3.4.1.4" class="ltx_text" style="font-size:90%;"></span></td>
<td id="S2.T1.3.4.2" class="ltx_td ltx_align_left" style="padding-top: 1.15pt; padding-bottom: 1.15pt"><span id="S2.T1.3.4.2.1" class="ltx_text"></span><span id="S2.T1.3.4.2.2" class="ltx_text" style="font-size:90%;"> <span id="S2.T1.3.4.2.2.1" class="ltx_tabular ltx_align_middle"> <span id="S2.T1.3.4.2.2.1.1" class="ltx_tr"> <span id="S2.T1.3.4.2.2.1.1.1" class="ltx_td ltx_nopad_r ltx_align_left" style="padding-top:1.15pt;padding-bottom:1.15pt;">Text Matching/Code Generation/</span></span> <span id="S2.T1.3.4.2.2.1.2" class="ltx_tr"> <span id="S2.T1.3.4.2.2.1.2.1" class="ltx_td ltx_nopad_r ltx_align_left" style="padding-top:1.15pt;padding-bottom:1.15pt;">Action Grounding</span></span> </span></span><span id="S2.T1.3.4.2.3" class="ltx_text"></span><span id="S2.T1.3.4.2.4" class="ltx_text" style="font-size:90%;"></span></td>
</tr>
<tr id="S2.T1.3.5" class="ltx_tr">
<td id="S2.T1.3.5.1" class="ltx_td ltx_align_left" style="padding-top: 1.15pt; padding-bottom: 1.15pt"><span id="S2.T1.3.5.1.1" class="ltx_text" style="font-size:90%;">Learning/Reflection/Memory</span></td>
<td id="S2.T1.3.5.2" class="ltx_td ltx_align_left" style="padding-top: 1.15pt; padding-bottom: 1.15pt"><span id="S2.T1.3.5.2.1" class="ltx_text"></span><span id="S2.T1.3.5.2.2" class="ltx_text" style="font-size:90%;"> <span id="S2.T1.3.5.2.2.1" class="ltx_tabular ltx_align_middle"> <span id="S2.T1.3.5.2.2.1.1" class="ltx_tr"> <span id="S2.T1.3.5.2.2.1.1.1" class="ltx_td ltx_nopad_r ltx_align_left" style="padding-top:1.15pt;padding-bottom:1.15pt;">RLHF/Multi-agent Debate/</span></span> <span id="S2.T1.3.5.2.2.1.2" class="ltx_tr"> <span id="S2.T1.3.5.2.2.1.2.1" class="ltx_td ltx_nopad_r ltx_align_left" style="padding-top:1.15pt;padding-bottom:1.15pt;">Vector Database</span></span> </span></span><span id="S2.T1.3.5.2.3" class="ltx_text"></span><span id="S2.T1.3.5.2.4" class="ltx_text" style="font-size:90%;"></span></td>
</tr>
<tr id="S2.T1.3.6" class="ltx_tr">
<td id="S2.T1.3.6.1" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.15pt; padding-bottom: 1.15pt"><span id="S2.T1.3.6.1.1" class="ltx_text" style="font-size:90%;">Summarization</span></td>
<td id="S2.T1.3.6.2" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.15pt; padding-bottom: 1.15pt"><span id="S2.T1.3.6.2.1" class="ltx_text"></span><span id="S2.T1.3.6.2.2" class="ltx_text" style="font-size:90%;"> <span id="S2.T1.3.6.2.2.1" class="ltx_tabular ltx_align_middle"> <span id="S2.T1.3.6.2.2.1.1" class="ltx_tr"> <span id="S2.T1.3.6.2.2.1.1.1" class="ltx_td ltx_nopad_r ltx_align_left" style="padding-top:1.15pt;padding-bottom:1.15pt;">Attention Mechanism and</span></span> <span id="S2.T1.3.6.2.2.1.2" class="ltx_tr"> <span id="S2.T1.3.6.2.2.1.2.1" class="ltx_td ltx_nopad_r ltx_align_left" style="padding-top:1.15pt;padding-bottom:1.15pt;">Natural Language Generation</span></span> </span></span><span id="S2.T1.3.6.2.3" class="ltx_text"></span><span id="S2.T1.3.6.2.4" class="ltx_text" style="font-size:90%;"></span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 1: </span>A simple illustration of the techniques for endowing the key ability.</figcaption>
</figure>
:::::

:::::: {#S2.SS3 .section .ltx_subsection}
### [2.3 ]{.ltx_tag .ltx_tag_subsection}Agent Design {#agent-design .ltx_title .ltx_title_subsection}

::: {#S2.SS3.p1 .ltx_para}
Task planning and tool usage represent the cornerstone of LLM's abilities. Others like perception, learning/reflection/memory (from feedback), and summarization are indeed critical, but they primarily serve to enhance and support these two core competencies. Therefore, concentrating on these two key competencies - [T]{#S2.SS3.p1.1.1 .ltx_text .ltx_font_bold}ask [P]{#S2.SS3.p1.1.2 .ltx_text .ltx_font_bold}lanning and [T]{#S2.SS3.p1.1.3 .ltx_text .ltx_font_bold}ool [U]{#S2.SS3.p1.1.4 .ltx_text .ltx_font_bold}sage (TPTU for short) - we have devised two distinct types of AI agents, as depicted in Figure [[3]{.ltx_text .ltx_ref_tag}](#S2.F3 "Figure 3 ‣ 2.3 Agent Design ‣ 2 Method ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}:
:::

::: {#S2.SS3.p2 .ltx_para}
- [[•]{.ltx_tag .ltx_tag_item}]{#S2.I3.i1}

  ::: {#S2.I3.i1.p1 .ltx_para}
  The first one, named as the [O]{#S2.I3.i1.p1.1.1 .ltx_text .ltx_font_bold}ne-step [A]{#S2.I3.i1.p1.1.2 .ltx_text .ltx_font_bold}gent (TPTU-OA), adopts a global perspective to interpret the original problem, effectively breaking it down into a sequence of sub-tasks in a single instance. This strategy fully harnesses the model's comprehensive understanding capabilities to map out the problem-solving steps for all sub-tasks at once. This method underscores the significance of a holistic understanding and planning of the overall task, albeit it might lack flexibility when dealing with individual sub-tasks.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S2.I3.i2}

  ::: {#S2.I3.i2.p1 .ltx_para}
  The second type, referred to as the [S]{#S2.I3.i2.p1.1.1 .ltx_text .ltx_font_bold}equential [A]{#S2.I3.i2.p1.1.2 .ltx_text .ltx_font_bold}gent (TPTU-SA), emphasizes tackling the current sub-task at hand. Upon successfully resolving the ongoing sub-task, this agent requests the LLMs to provide the succeeding sub-task. This approach enables the model to maintain a clear and concentrated focus throughout the problem-solving journey, tackling issues incrementally. Such a methodology allows for continuous feedback and progress within the confines of addressing a broader problem.
  :::
:::

<figure id="S2.F3" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_1">
<figure id="S2.F3.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<div id="S2.F3.sf1.1" class="ltx_block ltx_minipage ltx_align_top" style="width:359.9pt;">
<img src="/html/2308.03427/assets/x3.png" id="S2.F3.sf1.1.g1" class="ltx_graphics ltx_img_landscape" width="461" height="121" alt="Refer to caption" />
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span id="S2.F3.sf1.3.1.1" class="ltx_text" style="font-size:80%;">(a)</span> </span><span id="S2.F3.sf1.4.2" class="ltx_text" style="font-size:80%;">One-step Agent (TPTU-OA)</span></figcaption>
</figure>
</div>
<div class="ltx_flex_break">

</div>
<div class="ltx_flex_cell ltx_flex_size_1">
<figure id="S2.F3.sf2" class="ltx_figure ltx_figure_panel ltx_align_center">
<div id="S2.F3.sf2.1" class="ltx_block ltx_minipage ltx_align_top" style="width:359.9pt;">
<img src="/html/2308.03427/assets/x4.png" id="S2.F3.sf2.1.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="461" height="215" alt="Refer to caption" />
</div>
<figcaption><span class="ltx_tag ltx_tag_figure"><span id="S2.F3.sf2.3.1.1" class="ltx_text" style="font-size:80%;">(b)</span> </span><span id="S2.F3.sf2.4.2" class="ltx_text" style="font-size:80%;">Sequential Agent (TPTU-SA)</span></figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 3: </span>The workflows of the One-step Agent and the Sequential Agent are specifically designed to assess the Task Planning and Tool Usage abilities of LLMs.</figcaption>
</figure>

::: {#S2.SS3.p3 .ltx_para}
These two distinct agent models represent two disparate problem-solving strategies - the [one-step and sequential]{#S2.SS3.p3.1.1 .ltx_text style="color:#000000;"} resolution [^4^[[^4^[4]{.ltx_tag .ltx_tag_note}One can also combine the two strategies to design a hierarchical agent, but this is beyond the scope of this paper.]{.ltx_note_content}]{.ltx_note_outer}]{#footnote4 .ltx_note .ltx_role_footnote}. In our subsequent experiments, we aim to understand their respective strengths and weaknesses and how they can be best utilized to leverage the capabilities of LLMs in real-world problem-solving scenarios.
:::
::::::
:::::::::::::::

::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: {#S3 .section .ltx_section}
## [3 ]{.ltx_tag .ltx_tag_section}Evaluation {#evaluation .ltx_title .ltx_title_section}

::: {#S3.p1 .ltx_para}
We instantiate the proposed LLM-based AI Agent framework (TPTU-OA and TPTU-SA) with different LLMs and evaluate their performance on typical tasks.
:::

:::::::::::::: {#S3.SS1 .section .ltx_subsection}
### [3.1 ]{.ltx_tag .ltx_tag_subsection}Preparations {#preparations .ltx_title .ltx_title_subsection}

::: {#S3.SS1.p1 .ltx_para}
Before beginning our evaluation, we first outline the preparations. We will give detailed descriptions of the datasets, available tools, and popular large language models.
:::

::::::: {#S3.SS1.SSS1 .section .ltx_subsubsection}
#### [3.1.1 ]{.ltx_tag .ltx_tag_subsubsection}Datasets {#datasets .ltx_title .ltx_title_subsubsection}

::: {#S3.SS1.SSS1.p1 .ltx_para}
We first clarify the motivations behind our choice of tools for evaluation. The selection was guided by two primary factors: [the number of tools]{#S3.SS1.SSS1.p1.1.1 .ltx_text .ltx_font_bold} to be evaluated and [the specific tools]{#S3.SS1.SSS1.p1.1.2 .ltx_text .ltx_font_bold} to be included.
:::

::: {#S3.SS1.SSS1.p2 .ltx_para}
Firstly, regarding the number of tools, it is important to state that our proposed evaluation framework is extensible. It can incorporate any number of tools as pluggable components to be managed by the LLM-based AI agents. Secondly, looking at the current work on tool-augmented LLMs, such as T-Bench \[[16](#bib.bib16){.ltx_ref}\] and ToolBench \[[17](#bib.bib17){.ltx_ref}\], we see that only a handful of tools are launched and executed in a single scenario. Meanwhile, API-Bank \[[18](#bib.bib18){.ltx_ref}\], in a single scenario, typically dispatches only one API tool and awaits its response. APIBench \[[19](#bib.bib19){.ltx_ref}\] and ToolAlpaca \[[20](#bib.bib20){.ltx_ref}\] do not even execute a tool response. Hence, for the sake of simplicity and focus in our evaluation, we have decided to primarily assess two tools (which can be called multiple times) within a single scenario.
:::

::: {#S3.SS1.SSS1.p3 .ltx_para}
Secondly, we also need to decide which specific tools should be used for evaluation. Consider a real-world scenario where we pose the question: "How much budget is required to offer a $\$$`<!-- -->`{=html}100 incentive to each employee who has been with the company for over five years?\". To answer this, we first need to retrieve the relevant data from a database, typically using SQL, to find the number of eligible employees. Then, we need to perform a mathematical calculation to estimate the total budget. Such scenarios are quite common in daily life where the formulation and resolution of a question often involve SQL and mathematical tools.
:::

::: {#S3.SS1.SSS1.p4 .ltx_para}
Recognizing the importance of these tools, we have chosen to focus our evaluation on SQL and Python generators, which represent the capabilities of database querying and mathematical computation, respectively. To this end, we have prepared 120 question-answer pairs that vary in complexity. These pairs provide a rigorous assessment of the LLM-based AI agents in understanding, generating, and utilizing these essential tools. For further information on these queries and their corresponding demonstrations, please refer to Appendix [[A]{.ltx_text .ltx_ref_tag}](#A1 "Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::
:::::::

::::: {#S3.SS1.SSS2 .section .ltx_subsubsection}
#### [3.1.2 ]{.ltx_tag .ltx_tag_subsubsection}Tools {#tools .ltx_title .ltx_title_subsubsection}

::: {#S3.SS1.SSS2.p1 .ltx_para}
We have defined a total of 12 available tools for the selection of the LLM-based AI agents for evaluation. They are defined as follows:
:::

::: {#S3.SS1.SSS2.p2 .ltx_para}
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I1.i1}

  ::: {#S3.I1.i1.p1 .ltx_para}
  SQL generator: Given an input question and a database, create a syntactically correct SQLite query statement.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I1.i2}

  ::: {#S3.I1.i2.p1 .ltx_para}
  Python generator: Given an input question and some information, generate a syntactically correct Python code.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I1.i3}

  ::: {#S3.I1.i3.p1 .ltx_para}
  Weather query tool: Given a location, output the current real-time weather at that location.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I1.i4}

  ::: {#S3.I1.i4.p1 .ltx_para}
  Image generator: Given a text description, generate a related image.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I1.i5}

  ::: {#S3.I1.i5.p1 .ltx_para}
  Text extractor: Given a link to an image, extract the corresponding text and its position coordinates.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I1.i6}

  ::: {#S3.I1.i6.p1 .ltx_para}
  Translator: Given a piece of text, translate it into other languages.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I1.i7}

  ::: {#S3.I1.i7.p1 .ltx_para}
  Bing Searcher: Given a piece of text, conduct a search on the Bing browser and return content.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I1.i8}

  ::: {#S3.I1.i8.p1 .ltx_para}
  Shell generator: Given an input question and some information, generate a syntactically correct Shell code.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I1.i9}

  ::: {#S3.I1.i9.p1 .ltx_para}
  Java generator: Given an input question and some information, generate a syntactically correct Java code.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I1.i10}

  ::: {#S3.I1.i10.p1 .ltx_para}
  Wikipedia searcher: Given a piece of text, conduct a search on Wikipedia and return content.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I1.i11}

  ::: {#S3.I1.i11.p1 .ltx_para}
  Office software: Given a text description, automatically generate corresponding long documents or spreadsheets or PPTs.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I1.i12}

  ::: {#S3.I1.i12.p1 .ltx_para}
  Movie player: Given a movie name, automatically play the corresponding movie resources.
  :::
:::
:::::

:::: {#S3.SS1.SSS3 .section .ltx_subsubsection}
#### [3.1.3 ]{.ltx_tag .ltx_tag_subsubsection}LLMs {#llms .ltx_title .ltx_title_subsubsection}

::: {#S3.SS1.SSS3.p1 .ltx_para}
The LLMs evaluated in this paper are listed in Table [[2]{.ltx_text .ltx_ref_tag}](#S3.T2 "Table 2 ‣ 3.1.3 LLMs ‣ 3.1 Preparations ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, elaborated as follows:

- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I2.i1}

  ::: {#S3.I2.i1.p1 .ltx_para}
  [GPT]{#S3.I2.i1.p1.1.1 .ltx_text .ltx_font_bold} series developed by OpenAI boasts a powerful language model with a vast number of parameters, enabling it to tackle intricate problems efficiently. This paper aims to evaluate the performance of ChatGPT, which balances the performance with costs (the number of OpenAI API calls).
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I2.i2}

  ::: {#S3.I2.i2.p1 .ltx_para}
  [Claude]{#S3.I2.i2.p1.1.1 .ltx_text .ltx_font_bold} is committed to maintaining honesty and ensuring user safety, which is developed by Anthropic. With its impressive size, Claude ranks among the largest language models globally and poses a formidable challenge to ChatGPT as a strong competitor.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I2.i3}

  ::: {#S3.I2.i3.p1 .ltx_para}
  [InternLM]{#S3.I2.i3.p1.1.1 .ltx_text .ltx_font_bold}, a sophisticated language model developed by Shanghai AI Lab, boasts a multi-round dialogue capability and an impressive ability to comprehend super-long text. This language model is meticulously designed to cater to the nuances of the Chinese language, enabling it to comprehensively understand and effectively process Chinese text. Here, we adopted the version with 120 billion parameters.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I2.i4}

  ::: {#S3.I2.i4.p1 .ltx_para}
  [Ziya]{#S3.I2.i4.p1.1.1 .ltx_text .ltx_font_bold} is an expansive and robust pre-training model developed by IDEA, derived from the LLaMa with 13 billion parameters. This comprehensive model exhibits a wide range of capabilities, including translation, programming, and mathematical calculations. Notably, it stands out as a bilingual LLM, highlighting its ability to effectively process and comprehend text in Chinese.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I2.i5}

  ::: {#S3.I2.i5.p1 .ltx_para}
  [ChatGLM]{#S3.I2.i5.p1.1.1 .ltx_text .ltx_font_bold}, developed by Tsinghua University, is an open-source dialogue language model that supports bilingual Q&A in Chinese and English, with a particular focus on Chinese optimization. Built on the General Language Model (GLM) architecture and utilizing model quantization technology, the ChatGLM can be easily deployed on consumer-grade graphics cards, enabling local implementation by users.
  :::
- [[•]{.ltx_tag .ltx_tag_item}]{#S3.I2.i6}

  ::: {#S3.I2.i6.p1 .ltx_para}
  [Chinese-Alpaca-Plus]{#S3.I2.i6.p1.1.1 .ltx_text .ltx_font_bold} is achieved by extending LLaMA's existing vocabulary with an additional 20,000 Chinese tokens from Meta AI (formerly known as Facebook AI Research Laboratory). In this version, we use a model with 33 billion parameters. The training text has been expanded to 120GB, and the fine-tuning instruction data has been increased to 4.3M.
  :::
:::

<figure id="S3.T2" class="ltx_table">
<table id="S3.T2.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="S3.T2.1.1" class="ltx_tr">
<td id="S3.T2.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T2.1.1.1.1" class="ltx_text ltx_font_bold">Organization</span></td>
<td id="S3.T2.1.1.2" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T2.1.1.2.1" class="ltx_text ltx_font_bold">Model Name</span></td>
<td id="S3.T2.1.1.3" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T2.1.1.3.1" class="ltx_text ltx_font_bold">Model Parameters</span></td>
</tr>
<tr id="S3.T2.1.2" class="ltx_tr">
<td id="S3.T2.1.2.1" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">OpenAI</td>
<td id="S3.T2.1.2.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><a href="https://chat.openai.com/" class="ltx_ref ltx_href" target="_blank">ChatGPT</a>[<a href="#bib.bib21" class="ltx_ref">21</a>]</td>
<td id="S3.T2.1.2.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">200B</td>
</tr>
<tr id="S3.T2.1.3" class="ltx_tr">
<td id="S3.T2.1.3.1" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Anthropic</td>
<td id="S3.T2.1.3.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><a href="https://claude.ai" class="ltx_ref ltx_href" target="_blank">Claude</a>[<a href="#bib.bib22" class="ltx_ref">22</a>]</td>
<td id="S3.T2.1.3.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">&gt;52B</td>
</tr>
<tr id="S3.T2.1.4" class="ltx_tr">
<td id="S3.T2.1.4.1" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Shanghai AI Lab</td>
<td id="S3.T2.1.4.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">InternLM</td>
<td id="S3.T2.1.4.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">120B</td>
</tr>
<tr id="S3.T2.1.5" class="ltx_tr">
<td id="S3.T2.1.5.1" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">IDEA</td>
<td id="S3.T2.1.5.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><a href="https://github.com/IDEA-CCNL/Fengshenbang-LM" class="ltx_ref ltx_href" target="_blank">Ziya-13B</a></td>
<td id="S3.T2.1.5.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">13B</td>
</tr>
<tr id="S3.T2.1.6" class="ltx_tr">
<td id="S3.T2.1.6.1" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Tsinghua University</td>
<td id="S3.T2.1.6.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><a href="https://github.com/THUDM/GLM-130B" class="ltx_ref ltx_href" target="_blank">ChatGLM-130B</a>[<a href="#bib.bib23" class="ltx_ref">23</a>]</td>
<td id="S3.T2.1.6.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">130B</td>
</tr>
<tr id="S3.T2.1.7" class="ltx_tr">
<td id="S3.T2.1.7.1" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">-</td>
<td id="S3.T2.1.7.2" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><a href="https://github.com/ymcui/Chinese-LLaMA-Alpaca/" class="ltx_ref ltx_href" target="_blank">Chinese-Alpaca-Plus-33B</a>[<a href="#bib.bib24" class="ltx_ref">24</a>, <a href="#bib.bib25" class="ltx_ref">25</a>]</td>
<td id="S3.T2.1.7.3" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">33B</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 2: </span>The LLMs evaluated in this paper.</figcaption>
</figure>
::::
::::::::::::::

:::::::::::::::::::::::::::::::: {#S3.SS2 .section .ltx_subsection}
### [3.2 ]{.ltx_tag .ltx_tag_subsection}Evaluation on Task Planning Ability {#evaluation-on-task-planning-ability .ltx_title .ltx_title_subsection}

::: {#S3.SS2.p1 .ltx_para}
In this section, to evaluate the planning capabilities of the LLM-based AI agents, we have structured the evaluations as follows.
:::

::: {#S3.SS2.p2 .ltx_para}
For TPTU-OA, we begin by examining the agents' ability to plan the order of tool use. This is followed by an evaluation of the agents' capacity to not only plan the sequence of tools but also the corresponding subtask descriptions. Subsequently, we conduct a specialized planning evaluation where the agents must generate multiple sequences of key-value pairs of the form {tool: subtask description} in complex problem teardowns. Moreover, we expand the toolset with additional, unrelated tools to further challenge and reassess the planning ability of the LLM-based AI agents.
:::

::: {#S3.SS2.p3 .ltx_para}
For TPTU-SA, we follow the regime that the agent should generate multiple sequences of key-value pairs of the form {tool: subtask description} for evaluation.
:::

:::::::: {#S3.SS2.SSS1 .section .ltx_subsubsection}
#### [3.2.1 ]{.ltx_tag .ltx_tag_subsubsection}TPTU-OA: Tool Order Planning {#tptu-oa-tool-order-planning .ltx_title .ltx_title_subsubsection}

::: {#S3.SS2.SSS1.p1 .ltx_para}
Here, we utilize two kinds of tools for problem-solving: the SQL generator, which retrieves data from databases, and the Python generator, adept at addressing mathematical questions.
:::

::: {#S3.SS2.SSS1.p2 .ltx_para}
To validate the capacity of the LLM-based AI agents to strategically plan for the tool order, we designed the prompt as shown in Figure [[8]{.ltx_text .ltx_ref_tag}](#A2.F8 "Figure 8 ‣ Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} of Appendix [[B]{.ltx_text .ltx_ref_tag}](#A2 "Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}. This design is motivated by the goal to assess the ability of LLM-based AI agents to understand complex problems, subsequently decomposing them into a sequence of simpler tasks executed by appropriately selected tools. Specifically, we require the LLM-based AI agent to follow our instructions, select tools from our pre-defined tool set with detailed function descriptions, conform to the given format strictly, and understand the demonstrations to learn from them.
:::

::: {#S3.SS2.SSS1.p3 .ltx_para}
Upon feeding these prompts into the LLM-based AI agents under evaluation, we obtained the following accuracy rates for the tool planning, as shown in Table [[3]{.ltx_text .ltx_ref_tag}](#S3.T3 "Table 3 ‣ 3.2.1 TPTU-OA: Tool Order Planning ‣ 3.2 Evaluation on Task Planning Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::

<figure id="S3.T3" class="ltx_table">
<table id="S3.T3.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="S3.T3.1.1" class="ltx_tr">
<td id="S3.T3.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T3.1.1.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T3.1.1.2" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGPT</td>
<td id="S3.T3.1.1.3" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Claude</td>
<td id="S3.T3.1.1.4" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Ziya</td>
</tr>
<tr id="S3.T3.1.2" class="ltx_tr">
<td id="S3.T3.1.2.1" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T3.1.2.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T3.1.2.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">100%</td>
<td id="S3.T3.1.2.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">100%</td>
<td id="S3.T3.1.2.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">45%</td>
</tr>
<tr id="S3.T3.1.3" class="ltx_tr">
<td id="S3.T3.1.3.1" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T3.1.3.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T3.1.3.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGLM</td>
<td id="S3.T3.1.3.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Chinese-Alpaca-Plus</td>
<td id="S3.T3.1.3.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">InternLM</td>
</tr>
<tr id="S3.T3.1.4" class="ltx_tr">
<td id="S3.T3.1.4.1" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T3.1.4.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T3.1.4.2" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">45%</td>
<td id="S3.T3.1.4.3" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">20%</td>
<td id="S3.T3.1.4.4" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">80%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 3: </span>The evaluation results for the planning of tool order generation.</figcaption>
</figure>

::: {#S3.SS2.SSS1.p4 .ltx_para}
The results of our experiments indicate that models, notably Ziya and ChatGLM, frequently grapple with the generation of lists in the correct format. For other models, the predominant challenges lie in generating tools in the correct sequence or in the occasional omission of necessary tools. Nonetheless, the issue of parsing list formats is generally negligible.
:::

::: {#S3.SS2.SSS1.p5 .ltx_para}
These findings suggest that the majority of LLM-based AI agents possess a fundamental capability to analyze the tool needs of a given problem and understand its task requirements. To further explore whether these LLM-based AI agents can effectively break down the original problem into sub-tasks, we proceed to the following section.
:::
::::::::

:::::: {#S3.SS2.SSS2 .section .ltx_subsubsection}
#### [3.2.2 ]{.ltx_tag .ltx_tag_subsubsection}TPTU-OA: Tool Order Planning and Subtask Description Generation {#tptu-oa-tool-order-planning-and-subtask-description-generation .ltx_title .ltx_title_subsubsection}

::: {#S3.SS2.SSS2.p1 .ltx_para}
Simply planning the order of tool usage is not sufficient to fully address a problem. To truly solve it, we need to provide a guide or instructions for the usage of each tool, that is, a decomposed subtask description. Therefore, we can decompose the original complex problem into two separate sequences. One sequence represents the order in which the tools are utilized, while the other sequence corresponds to the subtask descriptions that each tool in the tool sequence aims to resolve. A problem is only truly solved when both the tool and subtask description sequences have been successfully planned. In order to verify whether LLM-based AI agents truly have the ability to solve complex problems, we designed a new prompt as shown in Figure [[9]{.ltx_text .ltx_ref_tag}](#A2.F9 "Figure 9 ‣ Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} of Appendix [[B]{.ltx_text .ltx_ref_tag}](#A2 "Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}. The main improvement is to plan the corresponding subtask description for each tool after the tool planning is completed.
:::

<figure id="S3.T4" class="ltx_table">
<table id="S3.T4.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="S3.T4.1.1" class="ltx_tr">
<td id="S3.T4.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T4.1.1.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T4.1.1.2" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGPT</td>
<td id="S3.T4.1.1.3" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Claude</td>
<td id="S3.T4.1.1.4" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Ziya</td>
</tr>
<tr id="S3.T4.1.2" class="ltx_tr">
<td id="S3.T4.1.2.1" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T4.1.2.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T4.1.2.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">55%</td>
<td id="S3.T4.1.2.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">15%</td>
<td id="S3.T4.1.2.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">10%</td>
</tr>
<tr id="S3.T4.1.3" class="ltx_tr">
<td id="S3.T4.1.3.1" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T4.1.3.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T4.1.3.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGLM</td>
<td id="S3.T4.1.3.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Chinese-Alpaca-Plus</td>
<td id="S3.T4.1.3.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">InternLM</td>
</tr>
<tr id="S3.T4.1.4" class="ltx_tr">
<td id="S3.T4.1.4.1" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T4.1.4.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T4.1.4.2" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">10%</td>
<td id="S3.T4.1.4.3" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">0%</td>
<td id="S3.T4.1.4.4" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">45%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 4: </span>The evaluation results for the planning of tool order and subtask description generation.</figcaption>
</figure>

::: {#S3.SS2.SSS2.p2 .ltx_para}
After feeding the prompt to these LLM-based AI agents, we get results shown in Table [[4]{.ltx_text .ltx_ref_tag}](#S3.T4 "Table 4 ‣ 3.2.2 TPTU-OA: Tool Order Planning and Subtask Description Generation ‣ 3.2 Evaluation on Task Planning Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::

::: {#S3.SS2.SSS2.p3 .ltx_para}
Although the generation of tool sequences and their corresponding subtask descriptions might be an effective way to problem-solving, there is a significant decrease in accuracy for all LLMs as can be seen. We hypothesize that there are a few potential drawbacks to this method:

1.  [[1.]{.ltx_tag .ltx_tag_item}]{#S3.I3.i1}

    ::: {#S3.I3.i1.p1 .ltx_para}
    [Difficulty in Error Tracking and Debugging]{#S3.I3.i1.p1.1.1 .ltx_text .ltx_font_bold}. Generating the complete tool and subtask sequences may make it more challenging to track and debug errors. If an error arises within the sequence, it might require a total regeneration instead of a simple modification or repair to the erroneous part.
    :::
2.  [[2.]{.ltx_tag .ltx_tag_item}]{#S3.I3.i2}

    ::: {#S3.I3.i2.p1 .ltx_para}
    [Tool-Subtask Pairing Issue]{#S3.I3.i2.p1.1.1 .ltx_text .ltx_font_bold}. If all tool sequences and subtask descriptions are generated independently, there's an inherent risk of misalignment between the tools and their corresponding subtasks. This could potentially lead to an improper pairing, which, in turn, could result in a flawed or ineffective solution that fails to appropriately resolve the given problem.
    :::
3.  [[3.]{.ltx_tag .ltx_tag_item}]{#S3.I3.i3}

    ::: {#S3.I3.i3.p1 .ltx_para}
    [Lack of Flexibility]{#S3.I3.i3.p1.1.1 .ltx_text .ltx_font_bold}. The approach may lack this flexibility when facing complex problems requiring adjustments to the tool or subtask sequence.
    :::
4.  [[4.]{.ltx_tag .ltx_tag_item}]{#S3.I3.i4}

    ::: {#S3.I3.i4.p1 .ltx_para}
    [Dependency on Global Information]{#S3.I3.i4.p1.1.1 .ltx_text .ltx_font_bold}. Generating the entire tool and subtask sequences requires a global understanding and planning of the entire problem. However, in some instances, certain parts of the problem might not be clear at the early stages of problem-solving, which could pose challenges within this framework.
    :::
:::
::::::

:::::::: {#S3.SS2.SSS3 .section .ltx_subsubsection}
#### [3.2.3 ]{.ltx_tag .ltx_tag_subsubsection}TPTU-OA: The Planning of Tool-Subtask Pair {#tptu-oa-the-planning-of-tool-subtask-pair .ltx_title .ltx_title_subsubsection}

::: {#S3.SS2.SSS3.p1 .ltx_para}
To mitigate the aforementioned issue, we propose a novel approach to foster flexible problem-solving with the LLM-based AI agent. We prompt the agent to generate multiple sequences, each consisting of a key-value pair in the format of {tool: subtask description} that associates a tool with its respective subtask description. This allows us to simultaneously plan the tool choice and subtask without the risk of improper matching. Moreover, it offers the flexibility to update the planned sequences in real-time based on evolving problem feedback, enhancing adaptability and efficiency when addressing complex tasks.
:::

::: {#S3.SS2.SSS3.p2 .ltx_para}
With this consideration, we have designed a unique prompt that encourages this advanced problem-solving strategy. In the following section, we delve into the specifics of this prompt design in Figure [[10]{.ltx_text .ltx_ref_tag}](#A2.F10 "Figure 10 ‣ Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} of Appendix [[B]{.ltx_text .ltx_ref_tag}](#A2 "Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}. The key improvement in this prompt is its directive for the LLM-based AI agents to stringently adhere to the predefined dictionary format. To facilitate this, we offer several demonstrations in our desired format, serving as references for the language model to follow.
:::

<figure id="S3.T5" class="ltx_table">
<table id="S3.T5.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="S3.T5.1.1" class="ltx_tr">
<td id="S3.T5.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T5.1.1.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T5.1.1.2" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGPT</td>
<td id="S3.T5.1.1.3" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Claude</td>
<td id="S3.T5.1.1.4" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Ziya</td>
</tr>
<tr id="S3.T5.1.2" class="ltx_tr">
<td id="S3.T5.1.2.1" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T5.1.2.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T5.1.2.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">75%</td>
<td id="S3.T5.1.2.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">90%</td>
<td id="S3.T5.1.2.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">20%</td>
</tr>
<tr id="S3.T5.1.3" class="ltx_tr">
<td id="S3.T5.1.3.1" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T5.1.3.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T5.1.3.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGLM</td>
<td id="S3.T5.1.3.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Chinese-Alpaca-Plus</td>
<td id="S3.T5.1.3.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">InternLM</td>
</tr>
<tr id="S3.T5.1.4" class="ltx_tr">
<td id="S3.T5.1.4.1" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T5.1.4.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T5.1.4.2" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">0%</td>
<td id="S3.T5.1.4.3" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">5%</td>
<td id="S3.T5.1.4.4" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">55%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 5: </span>The evaluation results for the planning of Tool-Subtask pair.</figcaption>
</figure>

::: {#S3.SS2.SSS3.p3 .ltx_para}
After feeding the prompt to these LLM-based AI agents, we get results shown in Table [[5]{.ltx_text .ltx_ref_tag}](#S3.T5 "Table 5 ‣ 3.2.3 TPTU-OA: The Planning of Tool-Subtask Pair ‣ 3.2 Evaluation on Task Planning Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::

::: {#S3.SS2.SSS3.p4 .ltx_para}
Analyzing the results from Tables [[4]{.ltx_text .ltx_ref_tag}](#S3.T4 "Table 4 ‣ 3.2.2 TPTU-OA: Tool Order Planning and Subtask Description Generation ‣ 3.2 Evaluation on Task Planning Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} and [[5]{.ltx_text .ltx_ref_tag}](#S3.T5 "Table 5 ‣ 3.2.3 TPTU-OA: The Planning of Tool-Subtask Pair ‣ 3.2 Evaluation on Task Planning Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, we observe a marked improvement of 52.9% when the tool-subtask pairs are generated in a unified format compared to separate generation of tools and subtasks.
:::

::: {#S3.SS2.SSS3.p5 .ltx_para}
This significant performance enhancement can likely be attributed to the close coupling between tools and their associated subtasks in our unified generation strategy. When tools and subtasks are generated separately, there is a potential disconnect or lack of coherence between the two, which could lead to less accurate or efficient solutions. In contrast, by generating tool-subtask pairs together, we ensure that each tool is directly tied to its relevant subtask, leading to a more coordinated and effective problem-solving approach. This might explain the observed increase in overall performance.
:::
::::::::

:::::: {#S3.SS2.SSS4 .section .ltx_subsubsection}
#### [3.2.4 ]{.ltx_tag .ltx_tag_subsubsection}TPTU-OA: The Planning of Tool-Subtask Pair with Unrelated Tools {#tptu-oa-the-planning-of-tool-subtask-pair-with-unrelated-tools .ltx_title .ltx_title_subsubsection}

::: {#S3.SS2.SSS4.p1 .ltx_para}
So far, our analysis and evaluation have been primarily focused on the LLM-based AI agents' proficiency in planning with specific tools. However, we are also interested in how it would perform when faced with many irrelevant or similar tools. Therefore, for a more comprehensive assessment, we expanded the prompt in Table [[10]{.ltx_text .ltx_ref_tag}](#A2.F10 "Figure 10 ‣ Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} to include an additional ten unrelated tools, as illustrated in Figure [[11]{.ltx_text .ltx_ref_tag}](#A2.F11 "Figure 11 ‣ Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} of Appendix [[B]{.ltx_text .ltx_ref_tag}](#A2 "Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::

<figure id="S3.T6" class="ltx_table">
<table id="S3.T6.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="S3.T6.1.1" class="ltx_tr">
<td id="S3.T6.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T6.1.1.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T6.1.1.2" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGPT</td>
<td id="S3.T6.1.1.3" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Claude</td>
<td id="S3.T6.1.1.4" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Ziya</td>
</tr>
<tr id="S3.T6.1.2" class="ltx_tr">
<td id="S3.T6.1.2.1" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T6.1.2.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T6.1.2.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">70%</td>
<td id="S3.T6.1.2.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">90%</td>
<td id="S3.T6.1.2.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">10%</td>
</tr>
<tr id="S3.T6.1.3" class="ltx_tr">
<td id="S3.T6.1.3.1" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T6.1.3.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T6.1.3.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGLM</td>
<td id="S3.T6.1.3.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Chinese-Alpaca-Plus</td>
<td id="S3.T6.1.3.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">InternLM</td>
</tr>
<tr id="S3.T6.1.4" class="ltx_tr">
<td id="S3.T6.1.4.1" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T6.1.4.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T6.1.4.2" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">0%</td>
<td id="S3.T6.1.4.3" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">5%</td>
<td id="S3.T6.1.4.4" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">50%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 6: </span>The evaluation results for the planning of Tool-Subtask pair with unrelated tools.</figcaption>
</figure>

::: {#S3.SS2.SSS4.p2 .ltx_para}
After feeding the prompt to these LLM-based AI agents, we get results shown in Table [[6]{.ltx_text .ltx_ref_tag}](#S3.T6 "Table 6 ‣ 3.2.4 TPTU-OA: The Planning of Tool-Subtask Pair with Unrelated Tools ‣ 3.2 Evaluation on Task Planning Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}. The results from our expanded evaluation demonstrate that even when presented with irrelevant or similar tools and descriptions, LLM-based AI agents consistently avoid selecting these unrelated tools (i.e., the accuracy has remained unchanged or exhibited only a marginal decrease compared with Table [[5]{.ltx_text .ltx_ref_tag}](#S3.T5 "Table 5 ‣ 3.2.3 TPTU-OA: The Planning of Tool-Subtask Pair ‣ 3.2 Evaluation on Task Planning Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}). This outcome indicates the effectiveness of our designed prompt, which successfully guides the LLM-based agents to understand the appropriate tool sequence for complex problem decomposition.
:::

::: {#S3.SS2.SSS4.p3 .ltx_para}
This observation reinforces the notion that a well-structured and informative prompt can efficiently guide AI agents to understand the core essence of the problem, thereby enabling them to sift through irrelevant information and focus on key tasks. This successful discrimination against unrelated tools also points towards the models' ability to understand the specific context of a problem and select the appropriate tools, thereby enhancing the overall problem-solving process.
:::
::::::

:::::::: {#S3.SS2.SSS5 .section .ltx_subsubsection}
#### [3.2.5 ]{.ltx_tag .ltx_tag_subsubsection}TPTU-SA: The Planning of Tool-Subtask Pair Generation {#tptu-sa-the-planning-of-tool-subtask-pair-generation .ltx_title .ltx_title_subsubsection}

::: {#S3.SS2.SSS5.p1 .ltx_para}
Upon identifying the drawbacks of first generating a list of tools and then generating corresponding subtask descriptions, we decided to focus subsequent tests on the generation of tool-subtask pairs. Consequently, in this section, we evaluate the capability of TPTU-SA to generate these tool-subtask pairs.
:::

::: {#S3.SS2.SSS5.p2 .ltx_para}
To achieve the goal of recursively generating tool-subtask pairs, we have designed prompts as illustrated in Figure [[12]{.ltx_text .ltx_ref_tag}](#A2.F12 "Figure 12 ‣ Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} of Appendix [[B]{.ltx_text .ltx_ref_tag}](#A2 "Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::

<figure id="S3.T7" class="ltx_table">
<table id="S3.T7.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="S3.T7.1.1" class="ltx_tr">
<td id="S3.T7.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T7.1.1.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T7.1.1.2" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGPT</td>
<td id="S3.T7.1.1.3" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Claude</td>
<td id="S3.T7.1.1.4" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Ziya</td>
</tr>
<tr id="S3.T7.1.2" class="ltx_tr">
<td id="S3.T7.1.2.1" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T7.1.2.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T7.1.2.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">80%</td>
<td id="S3.T7.1.2.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">100%</td>
<td id="S3.T7.1.2.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">10%</td>
</tr>
<tr id="S3.T7.1.3" class="ltx_tr">
<td id="S3.T7.1.3.1" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T7.1.3.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T7.1.3.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGLM</td>
<td id="S3.T7.1.3.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Chinese-Alpaca-Plus</td>
<td id="S3.T7.1.3.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">InternLM</td>
</tr>
<tr id="S3.T7.1.4" class="ltx_tr">
<td id="S3.T7.1.4.1" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T7.1.4.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T7.1.4.2" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">0%</td>
<td id="S3.T7.1.4.3" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">0%</td>
<td id="S3.T7.1.4.4" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">65%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 7: </span>The evaluation results for the planning of Tool-Subtask with the sequential agent.</figcaption>
</figure>

::: {#S3.SS2.SSS5.p3 .ltx_para}
The evaluation results are shown in Table [[7]{.ltx_text .ltx_ref_tag}](#S3.T7 "Table 7 ‣ 3.2.5 TPTU-SA: The Planning of Tool-Subtask Pair Generation ‣ 3.2 Evaluation on Task Planning Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}. Compared with results shown in Table [[5]{.ltx_text .ltx_ref_tag}](#S3.T5 "Table 5 ‣ 3.2.3 TPTU-OA: The Planning of Tool-Subtask Pair ‣ 3.2 Evaluation on Task Planning Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, TPTU-SA generally performs better than TPTU-OA especially for high--performing LLMs (e.g., ChatGPT, Claude and InternLM). We propose the following potential reasons for this observation:
:::

::: {#S3.SS2.SSS5.p4 .ltx_para}
1.  [[1.]{.ltx_tag .ltx_tag_item}]{#S3.I4.i1}

    ::: {#S3.I4.i1.p1 .ltx_para}
    [Sequentiality Mimics Human Problem-Solving]{#S3.I4.i1.p1.1.1 .ltx_text .ltx_font_bold}: In real-world scenarios, humans tend to solve complex problems by breaking them down into smaller, manageable subtasks which are often handled sequentially. Sequential agents are designed to mimic this step-by-step approach, which might inherently suit complex problem-solving better.
    :::
2.  [[2.]{.ltx_tag .ltx_tag_item}]{#S3.I4.i2}

    ::: {#S3.I4.i2.p1 .ltx_para}
    [Richer Contextual Understanding]{#S3.I4.i2.p1.1.1 .ltx_text .ltx_font_bold}: Sequential agents are exposed to the outcome of each previous subtask before moving on to the next one. This iterative process could facilitate a richer understanding of the problem context, enabling more accurate task planning and tool usage.
    :::
3.  [[3.]{.ltx_tag .ltx_tag_item}]{#S3.I4.i3}

    ::: {#S3.I4.i3.p1 .ltx_para}
    [Flexibility in Task Management]{#S3.I4.i3.p1.1.1 .ltx_text .ltx_font_bold}: In comparison to one-step agents, sequential agents might have more flexibility in managing tasks. They have the opportunity to correct errors or adjust their strategy after each step, which can lead to improved overall performance.
    :::
4.  [[4.]{.ltx_tag .ltx_tag_item}]{#S3.I4.i4}

    ::: {#S3.I4.i4.p1 .ltx_para}
    [Improved Learning From History]{#S3.I4.i4.p1.1.1 .ltx_text .ltx_font_bold}: The sequential process provides a history of actions and results which can be beneficial in learning. The agent can use this history to make better predictions about what tool to use next or what subtask to tackle, leading to more accurate and efficient problem-solving.
    :::
:::

::: {#S3.SS2.SSS5.p5 .ltx_para}
These points of analysis suggest that the structure and operation of sequential agents inherently confer certain advantages in complex problem-solving scenarios, leading to their superior performance.
:::
::::::::
::::::::::::::::::::::::::::::::

::::::::::::::::::::::::::::: {#S3.SS3 .section .ltx_subsection}
### [3.3 ]{.ltx_tag .ltx_tag_subsection}Evaluation on Tool Usage Ability {#evaluation-on-tool-usage-ability .ltx_title .ltx_title_subsection}

::: {#S3.SS3.p1 .ltx_para}
Before evaluating the end-to-end multi-tool usage ability of LLM-based AI agents, we first evaluate the effectiveness of single-tool usage for SQL generation and mathematical code generation.
:::

::: {#S3.SS3.p2 .ltx_para}
Subsequently, to assess the end-to-end performance of LLMs across various tools, two types of agents (TPTU-OA and TPTU-SA) were developed and several LLMs were subjected to testing under these agents. The role of the agents is to break down complex questions into simpler sub-questions and plan corresponding tools to solve them, based on the available toolset and corresponding tool descriptions.
:::

:::::::::::::::::::: {#S3.SS3.SSS1 .section .ltx_subsubsection}
#### [3.3.1 ]{.ltx_tag .ltx_tag_subsubsection}The effectiveness of Single Tool Usage {#the-effectiveness-of-single-tool-usage .ltx_title .ltx_title_subsubsection}

::: {#S3.SS3.SSS1.p1 .ltx_para}
Our aim is to systematically assess how effectively these models can use various tools, focusing on their proficiency with SQL and other coding languages.
:::

:::::: {#S3.SS3.SSS1.Px1 .section .ltx_paragraph}
##### The Effectiveness of simple SQL Creation {#the-effectiveness-of-simple-sql-creation .ltx_title .ltx_font_bold .ltx_title_paragraph}

::: {#S3.SS3.SSS1.Px1.p1 .ltx_para}
Using the schemas provided in Table [[13]{.ltx_text .ltx_ref_tag}](#A1.T13 "Table 13 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} and Table [[13]{.ltx_text .ltx_ref_tag}](#A1.T13 "Table 13 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, we construct questions similar to those in Table [[14]{.ltx_text .ltx_ref_tag}](#A1.T14 "Table 14 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, and refer readers to Appendix [[A]{.ltx_text .ltx_ref_tag}](#A1 "Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}. These questions are posed to various LLMs using our specifically designed prompts in Appendix [[B]{.ltx_text .ltx_ref_tag}](#A2 "Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::

::: {#S3.SS3.SSS1.Px1.p2 .ltx_para}
Following the tailored prompts, the LLMs are evaluated based on their responses to the presented queries. The results of this comprehensive assessment are compiled and exhibited in Figure [[8]{.ltx_text .ltx_ref_tag}](#S3.T8 "Table 8 ‣ The Effectiveness of simple SQL Creation ‣ 3.3.1 The effectiveness of Single Tool Usage ‣ 3.3 Evaluation on Tool Usage Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::

<figure id="S3.T8" class="ltx_table">
<table id="S3.T8.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="S3.T8.1.1" class="ltx_tr">
<td id="S3.T8.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T8.1.1.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T8.1.1.2" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGPT</td>
<td id="S3.T8.1.1.3" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Claude</td>
<td id="S3.T8.1.1.4" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Ziya</td>
</tr>
<tr id="S3.T8.1.2" class="ltx_tr">
<td id="S3.T8.1.2.1" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T8.1.2.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T8.1.2.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">90%</td>
<td id="S3.T8.1.2.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">100%</td>
<td id="S3.T8.1.2.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">50%</td>
</tr>
<tr id="S3.T8.1.3" class="ltx_tr">
<td id="S3.T8.1.3.1" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T8.1.3.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T8.1.3.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGLM</td>
<td id="S3.T8.1.3.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Chinese-Alpaca-Plus</td>
<td id="S3.T8.1.3.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">InternLM</td>
</tr>
<tr id="S3.T8.1.4" class="ltx_tr">
<td id="S3.T8.1.4.1" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T8.1.4.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T8.1.4.2" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">30%</td>
<td id="S3.T8.1.4.3" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">20%</td>
<td id="S3.T8.1.4.4" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">90%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 8: </span>The evaluation results for simple SQL questions.</figcaption>
</figure>

::: {#S3.SS3.SSS1.Px1.p3 .ltx_para}
This verifies the capabilities of each LLM in handling varying simple single-table SQL queries, thus providing a basis for comparison and analysis.
:::
::::::

::::::::: {#S3.SS3.SSS1.Px2 .section .ltx_paragraph}
##### The Effectiveness of Complex Nested SQL Creation {#the-effectiveness-of-complex-nested-sql-creation .ltx_title .ltx_font_bold .ltx_title_paragraph}

::: {#S3.SS3.SSS1.Px2.p1 .ltx_para}
Using the schemas provided in Table [[16]{.ltx_text .ltx_ref_tag}](#A1.T16 "Table 16 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, [[16]{.ltx_text .ltx_ref_tag}](#A1.T16 "Table 16 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, [[18]{.ltx_text .ltx_ref_tag}](#A1.T18 "Table 18 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, and [[18]{.ltx_text .ltx_ref_tag}](#A1.T18 "Table 18 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, we construct questions similar to those in Table [[19]{.ltx_text .ltx_ref_tag}](#A1.T19 "Table 19 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, and refer readers to Appendix [[A]{.ltx_text .ltx_ref_tag}](#A1 "Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}. For complex nested SQL questions, to further verify the SQL tool creation capability of LLMs, we have designed two types of prompts. One is the direct-guidance type, which explicitly informs the model that it needs to generate nested SQL query statements, as shown in Figure [[14]{.ltx_text .ltx_ref_tag}](#A2.F14 "Figure 14 ‣ Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} in Appendix [[B]{.ltx_text .ltx_ref_tag}](#A2 "Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::

::: {#S3.SS3.SSS1.Px2.p2 .ltx_para}
The other is based on the Chain-of-Thought (CoT) \[[26](#bib.bib26){.ltx_ref}\] approach, which leverages the model's ability to reason step by step to comprehend and craft SQL tools, and the prompt is shown in Figure [[15]{.ltx_text .ltx_ref_tag}](#A2.F15 "Figure 15 ‣ Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} in Appendix [[B]{.ltx_text .ltx_ref_tag}](#A2 "Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}. This method guides the model to sequentially generate SQL query clauses based on the problem context, thus breaking down the complex query generation task into smaller and manageable subtasks. This approach provides the model with a structured way to handle complex SQL tasks and showcases its capacity to engage in incremental reasoning and problem-solving.
:::

::: {#S3.SS3.SSS1.Px2.p3 .ltx_para}
The design of these two types of prompts serves as the backbone of our evaluation for complex nested SQL questions. While the direct-guidance approach focuses on testing the model's raw ability to generate SQL queries when explicitly instructed, the CoT-based approach evaluates a more nuanced capability: the model's reasoning and problem-solving skills in a step-by-step manner. Both these methods present unique challenges and offer valuable insights into the strengths and potential areas of improvement for the large language model's SQL tool generation ability. Subsequently, we will explore these two dimensions based on our experimental evaluations shown in Table [[9]{.ltx_text .ltx_ref_tag}](#S3.T9 "Table 9 ‣ The Effectiveness of Complex Nested SQL Creation ‣ 3.3.1 The effectiveness of Single Tool Usage ‣ 3.3 Evaluation on Tool Usage Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::

<figure id="S3.T9" class="ltx_table">
<table id="S3.T9.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="S3.T9.1.1" class="ltx_tr">
<td id="S3.T9.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T9.1.1.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T9.1.1.2" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGPT</td>
<td id="S3.T9.1.1.3" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Claude</td>
<td id="S3.T9.1.1.4" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Ziya</td>
</tr>
<tr id="S3.T9.1.2" class="ltx_tr">
<td id="S3.T9.1.2.1" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T9.1.2.1.1" class="ltx_text ltx_font_bold">Direct-based</span></td>
<td id="S3.T9.1.2.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">80%</td>
<td id="S3.T9.1.2.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">100%</td>
<td id="S3.T9.1.2.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">50%</td>
</tr>
<tr id="S3.T9.1.3" class="ltx_tr">
<td id="S3.T9.1.3.1" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T9.1.3.1.1" class="ltx_text ltx_font_bold">CoT-based</span></td>
<td id="S3.T9.1.3.2" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt">80%</td>
<td id="S3.T9.1.3.3" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt">100%</td>
<td id="S3.T9.1.3.4" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt">40%</td>
</tr>
<tr id="S3.T9.1.4" class="ltx_tr">
<td id="S3.T9.1.4.1" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T9.1.4.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T9.1.4.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGLM</td>
<td id="S3.T9.1.4.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Chinese-Alpaca-Plus</td>
<td id="S3.T9.1.4.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">InternLM</td>
</tr>
<tr id="S3.T9.1.5" class="ltx_tr">
<td id="S3.T9.1.5.1" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T9.1.5.1.1" class="ltx_text ltx_font_bold">Direct-based</span></td>
<td id="S3.T9.1.5.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">60%</td>
<td id="S3.T9.1.5.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">0%</td>
<td id="S3.T9.1.5.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">60%</td>
</tr>
<tr id="S3.T9.1.6" class="ltx_tr">
<td id="S3.T9.1.6.1" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T9.1.6.1.1" class="ltx_text ltx_font_bold">CoT-based</span></td>
<td id="S3.T9.1.6.2" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">70%</td>
<td id="S3.T9.1.6.3" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">0%</td>
<td id="S3.T9.1.6.4" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">50%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 9: </span>The evaluation results for complex nested SQL questions.</figcaption>
</figure>

::: {#S3.SS3.SSS1.Px2.p4 .ltx_para}
From the above results in Table [[9]{.ltx_text .ltx_ref_tag}](#S3.T9 "Table 9 ‣ The Effectiveness of Complex Nested SQL Creation ‣ 3.3.1 The effectiveness of Single Tool Usage ‣ 3.3 Evaluation on Tool Usage Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, it is clear that different models possess varying levels of proficiency in handling complex nested SQL tasks. Some models, like Claude, exhibit a robust capability in SQL generation, no matter whether the approach is direct or CoT-based. Most of these models demonstrate the SQL tool usage capability.
:::

::: {#S3.SS3.SSS1.Px2.p5 .ltx_para}
Specifically, some models such as ChatGLM show a distinct preference for the CoT-based approach, their performance improves when problems are broken down into smaller, manageable sub-tasks. This suggests that these models may have a stronger ability in sequential problem-solving and benefit more from step-by-step guidance. Conversely, models like Ziya and InternLM show a drop in performance when tasks are guided in the CoT-based format. This might indicate challenges in managing dependencies between sub-tasks or handling the continuity in sequential problem-solving. Lastly, Chinese-Alpaca-Plus shows significant room for improvement in complex SQL generation tasks. This shows that not all models are equally suited to handle advanced problem-solving involving nested SQL queries.
:::

::: {#S3.SS3.SSS1.Px2.p6 .ltx_para}
Overall, these findings underscore the importance of tailoring evaluation and training methodologies to the individual strengths and weaknesses of each model. By adopting this approach, we can better understand the performance variations across different models and provide targeted improvements to enhance their problem-solving abilities. Furthermore, this analysis highlights the potential of LLM-based agents in real-world applications, and the need to push their boundaries through continued research and development.
:::
:::::::::

::::::: {#S3.SS3.SSS1.Px3 .section .ltx_paragraph}
##### The Effectiveness of Mathematical Code Creation {#the-effectiveness-of-mathematical-code-creation .ltx_title .ltx_font_bold .ltx_title_paragraph}

::: {#S3.SS3.SSS1.Px3.p1 .ltx_para}
Following our evaluation of the LLM's proficiency in creating complex SQL queries, we now shift our focus to another tool creation: the creation of mathematical code. To the best of our knowledge, while large language models possess significant capabilities, they often fall short of providing highly accurate solutions to mathematical problems. Guiding these LLMs to generate mathematical code, and subsequently leveraging external tools to execute and derive the solutions, could significantly enhance their ability to tackle mathematical challenges.
:::

::: {#S3.SS3.SSS1.Px3.p2 .ltx_para}
In the upcoming section, we will conduct a detailed evaluation of guiding these LLMs to generate mathematical code. We aim to shed light on the true capability of these models in generating mathematical code and to elucidate the extent to which they can be utilized to aid in mathematical problem-solving. The prompt about how to guide LLMs is shown in Figure [[16]{.ltx_text .ltx_ref_tag}](#A2.F16 "Figure 16 ‣ Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} in Appendix [[B]{.ltx_text .ltx_ref_tag}](#A2 "Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::

<figure id="S3.T10" class="ltx_table">
<table id="S3.T10.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="S3.T10.1.1" class="ltx_tr">
<td id="S3.T10.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T10.1.1.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T10.1.1.2" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGPT</td>
<td id="S3.T10.1.1.3" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Claude</td>
<td id="S3.T10.1.1.4" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Ziya</td>
</tr>
<tr id="S3.T10.1.2" class="ltx_tr">
<td id="S3.T10.1.2.1" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T10.1.2.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T10.1.2.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">90%</td>
<td id="S3.T10.1.2.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">85%</td>
<td id="S3.T10.1.2.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">50%</td>
</tr>
<tr id="S3.T10.1.3" class="ltx_tr">
<td id="S3.T10.1.3.1" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T10.1.3.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T10.1.3.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGLM</td>
<td id="S3.T10.1.3.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Chinese-Alpaca-Plus</td>
<td id="S3.T10.1.3.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">InternLM</td>
</tr>
<tr id="S3.T10.1.4" class="ltx_tr">
<td id="S3.T10.1.4.1" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T10.1.4.1.1" class="ltx_text ltx_font_bold">Accuracy</span></td>
<td id="S3.T10.1.4.2" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">0%</td>
<td id="S3.T10.1.4.3" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">55%</td>
<td id="S3.T10.1.4.4" class="ltx_td ltx_align_center ltx_border_bb ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">95%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 10: </span>The evaluation results for mathematical questions.</figcaption>
</figure>

::: {#S3.SS3.SSS1.Px3.p3 .ltx_para}
The results shown in Table [[10]{.ltx_text .ltx_ref_tag}](#S3.T10 "Table 10 ‣ The Effectiveness of Mathematical Code Creation ‣ 3.3.1 The effectiveness of Single Tool Usage ‣ 3.3 Evaluation on Tool Usage Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} indicate that the capabilities of LLM-based agents to generate mathematical code vary considerably. High-performing models like ChatGPT, Claude, and InternLM display excellent proficiency, suggesting their potent ability to solve complex mathematical tasks. Middle-tier models, such as Ziya, show moderate success, indicating the potential for improvement and adaptability with the right training and optimization. Surprisingly, Alpaca demonstrated a notable proficiency in mathematical tasks, despite its poor performance in SQL generation, suggesting a possible inclination towards mathematical problems. In contrast, ChatGLM struggles significantly with mathematical code generation, underlining a potential weak spot in its capabilities and the need for focused improvement in this area.
:::

::: {#S3.SS3.SSS1.Px3.p4 .ltx_para}
Overall, these results underscore the task-dependent nature of LLMs' capabilities and highlight the importance of recognizing their individual strengths and weaknesses for optimal model guidance and enhanced problem-solving.
:::
:::::::
::::::::::::::::::::

:::::::: {#S3.SS3.SSS2 .section .ltx_subsubsection}
#### [3.3.2 ]{.ltx_tag .ltx_tag_subsubsection}TPTU-OA  and TPTU-SA: Tool Usage for Multiple Tools {#tptu-oa-and-tptu-sa-tool-usage-for-multiple-tools .ltx_title .ltx_title_subsubsection}

::: {#S3.SS3.SSS2.p1 .ltx_para}
We now aim to utilize the one-step agent and sequential agent, which we designed, to conduct an evaluation involving multiple tools. Corresponding prompts for each agent type have been crafted and are presented in Figure [[17]{.ltx_text .ltx_ref_tag}](#A2.F17 "Figure 17 ‣ Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} and Figure [[18]{.ltx_text .ltx_ref_tag}](#A2.F18 "Figure 18 ‣ Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} of Appendix [[B]{.ltx_text .ltx_ref_tag}](#A2 "Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, respectively.
:::

::: {#S3.SS3.SSS2.p2 .ltx_para}
In this phase of the evaluation, we need to automatically invoke the respective tools through code and produce the results. Given that user interface-based LLMs lack the capability to call external tools, we will only utilize the following four API-based LLMs (ChatGPT, Ziya, Chinese-Alpaca, and InternLM) for this comprehensive evaluation of external tool usage ability.
:::

<figure id="S3.T11" class="ltx_table">
<table id="S3.T11.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="S3.T11.1.1" class="ltx_tr">
<td id="S3.T11.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T11.1.1.1.1" class="ltx_text ltx_font_bold">Model</span></td>
<td id="S3.T11.1.1.2" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ChatGPT</td>
<td id="S3.T11.1.1.3" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Ziya</td>
<td id="S3.T11.1.1.4" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Chinese-Alpaca-Plus</td>
<td id="S3.T11.1.1.5" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt">InternLM</td>
</tr>
<tr id="S3.T11.1.2" class="ltx_tr">
<td id="S3.T11.1.2.1" class="ltx_td ltx_align_center" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T11.1.2.1.1" class="ltx_text ltx_font_bold">TPTU-OA</span></td>
<td id="S3.T11.1.2.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">50%</td>
<td id="S3.T11.1.2.3" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">0%</td>
<td id="S3.T11.1.2.4" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">0%</td>
<td id="S3.T11.1.2.5" class="ltx_td ltx_align_center ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">15%</td>
</tr>
<tr id="S3.T11.1.3" class="ltx_tr">
<td id="S3.T11.1.3.1" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="S3.T11.1.3.1.1" class="ltx_text ltx_font_bold">TPTU-SA</span></td>
<td id="S3.T11.1.3.2" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">55%</td>
<td id="S3.T11.1.3.3" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">0%</td>
<td id="S3.T11.1.3.4" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">0%</td>
<td id="S3.T11.1.3.5" class="ltx_td ltx_align_center ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">20%</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 11: </span>The evaluation results for end-to-end ability of multiple tools.</figcaption>
</figure>

::: {#S3.SS3.SSS2.p3 .ltx_para}
With agents mentioned above, the final results are presented in Table [[11]{.ltx_text .ltx_ref_tag}](#S3.T11 "Table 11 ‣ 3.3.2 TPTU-OA and TPTU-SA: Tool Usage for Multiple Tools ‣ 3.3 Evaluation on Tool Usage Ability ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}. The evaluation results demonstrate varying levels of task planning and tool usage capabilities among the four API-based LLMs. In the TPTU-OA evaluation, ChatGPT achieved a performance rate of 50%, significantly outperforming the other models, with InternLM at 15%, while both Ziya and Chinese-Alpaca did not manage to complete any tasks successfully, resulting in a score of 0%. In the TPTU-SA evaluation, an overall slight improvement was observed. ChatGPT maintained its leading position, with a slightly improved performance rate of 55%. InternLM also exhibited better performance, achieving a score of 20%, whereas Ziya and Chinese-Alpaca-Plus again failed to register any successful task completion.
:::

::: {#S3.SS3.SSS2.p4 .ltx_para}
These results reflect a notable discrepancy in the performance of LLMs when it comes to using external tools. ChatGPT and InternLM have demonstrated some ability to navigate these tasks, but their performance rates suggest there is significant room for improvement. Ziya and Chinese-Alpaca-Plus' performance indicates a struggle to effectively utilize external tools in their current state.
:::

::: {#S3.SS3.SSS2.p5 .ltx_para}
The differential performance between the TPTU-OA and TPTU-SA evaluation hints at the possible impact of the agent design on the LLMs' task execution ability. In particular, the performance increase under the sequential agent framework suggests that breaking down tasks into sequential steps might help LLM-based AI agents better utilize external tools. This insight could prove valuable in future improvements and developments of LLM-based AI agents. However, even with this approach, it is clear that LLM-based AI agents are far from perfect when it comes to effectively using external tools for complex tasks. This finding underlines the importance of further investigation and improvement in this domain.
:::
::::::::
:::::::::::::::::::::::::::::

:::::: {#S3.SS4 .section .ltx_subsection}
### [3.4 ]{.ltx_tag .ltx_tag_subsection}Insightful Observations {#insightful-observations .ltx_title .ltx_title_subsection}

::: {#S3.SS4.p1 .ltx_para}
Upon closer observation of our experimental results, we have identified several phenomena that deserved further exploration. These findings serve to broaden our understanding of LLM-based agents' behavior and capabilities and provide essential insights that could shape future research in this field. In the following, we will dissect these phenomena as shown in Figure [[4]{.ltx_text .ltx_ref_tag}](#S3.F4 "Figure 4 ‣ item 1 ‣ 3.4 Insightful Observations ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} - [[7]{.ltx_text .ltx_ref_tag}](#S3.F7 "Figure 7 ‣ item 4 ‣ 3.4 Insightful Observations ‣ 3 Evaluation ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, casting light on the weaknesses of LLM-based agents in the context of task planning and tool usage.
:::

::: {#S3.SS4.p2 .ltx_para}
1.  [[1.]{.ltx_tag .ltx_tag_item}]{#S3.I5.i1}

    ::: {#S3.I5.i1.p1 .ltx_para}
    [Misunderstanding Output Formats]{#S3.I5.i1.p1.1.1 .ltx_text .ltx_font_bold}: LLMs frequently encounter difficulty when output is required in specific formats such as lists or dictionaries. One such example includes inconsistencies between the number of tools and corresponding subtasks, leading to formatting issues that hinder the correct execution of tasks.
    :::

    <figure id="S3.F4" class="ltx_figure">
    <img src="/html/2308.03427/assets/x5.png" id="S3.F4.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="415" height="143" alt="Refer to caption" />
    <figcaption><span class="ltx_tag ltx_tag_figure">Figure 4: </span><span id="S3.F4.2.1" class="ltx_text ltx_font_bold">Issue-1</span>: Inconsistencies between the number of tools and corresponding subtasks.</figcaption>
    </figure>
2.  [[2.]{.ltx_tag .ltx_tag_item}]{#S3.I5.i2}

    ::: {#S3.I5.i2.p1 .ltx_para}
    [Struggling to Grasp Task Requirements]{#S3.I5.i2.p1.1.1 .ltx_text .ltx_font_bold}: LLMs might incorrectly disintegrate subproblems or apply unsuitable tools to carry out the subproblem. For example, an LLM might attempt to solve a purely mathematical problem by employing an SQL tool or could misunderstand similar terms like cube extraction and cube roots.
    :::

    <figure id="S3.F5" class="ltx_figure">
    <img src="/html/2308.03427/assets/x6.png" id="S3.F5.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="415" height="141" alt="Refer to caption" />
    <figcaption><span class="ltx_tag ltx_tag_figure">Figure 5: </span><span id="S3.F5.2.1" class="ltx_text ltx_font_bold">Issue-2</span>:Solve a purely mathematical problem by employing a SQL generator.</figcaption>
    </figure>
3.  [[3.]{.ltx_tag .ltx_tag_item}]{#S3.I5.i3}

    ::: {#S3.I5.i3.p1 .ltx_para}
    [Endless Extensions]{#S3.I5.i3.p1.1.1 .ltx_text .ltx_font_bold}: LLMs tend to overutilize a particular tool, even in instances where a single use would suffice for the correct result. This issue can lead to extended and nonsensical planning, where the same subtask is repeatedly solved.
    :::

    <figure id="S3.F6" class="ltx_figure">
    <img src="/html/2308.03427/assets/x7.png" id="S3.F6.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="415" height="187" alt="Refer to caption" />
    <figcaption><span class="ltx_tag ltx_tag_figure">Figure 6: </span><span id="S3.F6.2.1" class="ltx_text ltx_font_bold">Issue-3</span>: Unnecessary repetition of subtasks.</figcaption>
    </figure>
4.  [[4.]{.ltx_tag .ltx_tag_item}]{#S3.I5.i4}

    ::: {#S3.I5.i4.p1 .ltx_para}
    [Lack of Summary Skills]{#S3.I5.i4.p1.1.1 .ltx_text .ltx_font_bold}: LLMs do not take into account the responses to subproblems, relying instead on their internalized knowledge to generate the final answer. This may lead to a scenario where the final response only addresses a portion of the original query.
    :::

    <figure id="S3.F7" class="ltx_figure">
    <img src="/html/2308.03427/assets/x8.png" id="S3.F7.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="415" height="110" alt="Refer to caption" />
    <figcaption><span class="ltx_tag ltx_tag_figure">Figure 7: </span><span id="S3.F7.2.1" class="ltx_text ltx_font_bold">Issue-4</span>: Answering questions using common sense instead of generating code.</figcaption>
    </figure>
:::

::: {#S3.SS4.p3 .ltx_para}
By identifying and addressing these common issues, we stand a better chance at improving and refining LLMs, thereby unlocking their full potential.
:::
::::::
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

:::::::::: {#S4 .section .ltx_section}
## [4 ]{.ltx_tag .ltx_tag_section}Related Work {#related-work .ltx_title .ltx_title_section}

::: {#S4.p1 .ltx_para}
The remarkable capacity for usage and creation of tools have facilitated the transcendence of our innate physical and cognitive constraints, thereby profoundly advancing the progress and prosperity of human civilization and society. The swift advancement of LLM has rendered it feasible to use and create tools like humans. The integration of specialized tools with LLM has unlocked substantial potential in addressing intricate tasks. In this section, we offer a concise synopsis of the relevant research pertaining to tool learning based on LLMs.
:::

:::::: {#S4.SS1 .section .ltx_subsection}
### [4.1 ]{.ltx_tag .ltx_tag_subsection}Tool Usage {#tool-usage .ltx_title .ltx_title_subsection}

::: {#S4.SS1.p1 .ltx_para}
The initial advancements in tool learning have been constrained by the capabilities of artificial intelligence (AI) models. \[[27](#bib.bib27){.ltx_ref}\] Traditional deep learning approaches exhibit limitations in terms of comprehension of tool functionality and user intentions, and common sense reasoning abilities. Consequently, these limitations directly result in a notable decline in the stability and precision of tool learning methodologies. Recently, the advent of LLM has marked a pivotal juncture in the realm of tool learning. LLMs encompass a broad spectrum of common sense cognitive capabilities and exhibit remarkable proficiencies in natural language processing, reasoning, and interactive decision-making \[[28](#bib.bib28){.ltx_ref}, [29](#bib.bib29){.ltx_ref}, [30](#bib.bib30){.ltx_ref}, [31](#bib.bib31){.ltx_ref}, [32](#bib.bib32){.ltx_ref}\]. These attributes furnish indispensable prerequisites for LLMs to comprehend user intentions and effectively employ tools in tackling intricate tasks \[[33](#bib.bib33){.ltx_ref}\]. Simultaneously, the advancement of fine-tuning \[[34](#bib.bib34){.ltx_ref}, [35](#bib.bib35){.ltx_ref}, [36](#bib.bib36){.ltx_ref}, [37](#bib.bib37){.ltx_ref}, [38](#bib.bib38){.ltx_ref}\] and in-context learning \[[39](#bib.bib39){.ltx_ref}, [40](#bib.bib40){.ltx_ref}\] technology has offered robust support to LLM in addressing increasingly intricate challenges. In addition, tool usage can mitigate the inherent limitations of LLMs, encompassing the acquisition of up-to-date information from real-world events, refined mathematical computational abilities, and the mitigation of potential hallucinatory phenomena. \[[41](#bib.bib41){.ltx_ref}\]
:::

::: {#S4.SS1.p2 .ltx_para}
Within the realm of embodied intelligence \[[42](#bib.bib42){.ltx_ref}, [43](#bib.bib43){.ltx_ref}, [44](#bib.bib44){.ltx_ref}\], LLM engages in direct interactions with tangible tools like robots in order to enhance their cognitive abilities, optimize work productivity, and expand functional capacities. LLM possesses the capability to automatically devise action steps based on user intentions, enabling the guidance of robots in the completion of tasks \[[45](#bib.bib45){.ltx_ref}, [46](#bib.bib46){.ltx_ref}, [47](#bib.bib47){.ltx_ref}, [48](#bib.bib48){.ltx_ref}, [49](#bib.bib49){.ltx_ref}, [50](#bib.bib50){.ltx_ref}, [51](#bib.bib51){.ltx_ref}, [52](#bib.bib52){.ltx_ref}, [53](#bib.bib53){.ltx_ref}\], or alternatively, to directly generate underlying code that can be executed by robots \[[54](#bib.bib54){.ltx_ref}, [55](#bib.bib55){.ltx_ref}, [56](#bib.bib56){.ltx_ref}, [57](#bib.bib57){.ltx_ref}, [58](#bib.bib58){.ltx_ref}\]. Palm-E \[[50](#bib.bib50){.ltx_ref}\] introduced a multimodal language model which seamlessly integrates sensor data into its framework, enabling efficient planning of robot actions and task completion. Code as Policies (CaP) \[[58](#bib.bib58){.ltx_ref}\] facilitates the transformation of natural language instructions into code fragments that can be directly compiled and executed on robots. As for Inner Monologue \[[48](#bib.bib48){.ltx_ref}\], LLM incorporates diverse environmental feedback to construct inner monologues, thereby formulating effective robot control strategies. Furthermore, LP-SLAM \[[45](#bib.bib45){.ltx_ref}\] proposes a simultaneous localization and mapping (SLAM) system empowered with language perception capabilities, exploiting the potential of ChatGPT. PromptCraft \[[57](#bib.bib57){.ltx_ref}\], on the other hand, devises a function library tailored to ChatGPT on the robot platform, streamlining the conversion of user intentions into executable tasks via the underlying backend API.
:::

::: {#S4.SS1.p3 .ltx_para}
In addition to directly changing the real environment through interaction with tools in the physical world, LLM can also utilize software tools such as search engines \[[59](#bib.bib59){.ltx_ref}, [60](#bib.bib60){.ltx_ref}, [61](#bib.bib61){.ltx_ref}, [62](#bib.bib62){.ltx_ref}, [63](#bib.bib63){.ltx_ref}, [64](#bib.bib64){.ltx_ref}, [65](#bib.bib65){.ltx_ref}, [66](#bib.bib66){.ltx_ref}, [67](#bib.bib67){.ltx_ref}\], mobile \[[68](#bib.bib68){.ltx_ref}, [69](#bib.bib69){.ltx_ref}\], Microsoft Office \[[70](#bib.bib70){.ltx_ref}, [71](#bib.bib71){.ltx_ref}\], calculators \[[72](#bib.bib72){.ltx_ref}, [73](#bib.bib73){.ltx_ref}, [74](#bib.bib74){.ltx_ref}\], deep models \[[19](#bib.bib19){.ltx_ref}, [75](#bib.bib75){.ltx_ref}, [76](#bib.bib76){.ltx_ref}, [77](#bib.bib77){.ltx_ref}, [78](#bib.bib78){.ltx_ref}, [79](#bib.bib79){.ltx_ref}, [13](#bib.bib13){.ltx_ref}, [80](#bib.bib80){.ltx_ref}, [81](#bib.bib81){.ltx_ref}\] and other versatile APIs \[[82](#bib.bib82){.ltx_ref}, [5](#bib.bib5){.ltx_ref}, [83](#bib.bib83){.ltx_ref}, [84](#bib.bib84){.ltx_ref}, [20](#bib.bib20){.ltx_ref}, [85](#bib.bib85){.ltx_ref}\] to enhance model performance or complete complex workflows through flexible control of the software. Toolformer \[[5](#bib.bib5){.ltx_ref}\] employs a self-supervised methodology to fine-tune the language model, enabling it to acquire the ability to automatically invoke APIs. ART \[[86](#bib.bib86){.ltx_ref}\] leverages CoT \[[26](#bib.bib26){.ltx_ref}\] and In-context Learning \[[81](#bib.bib81){.ltx_ref}, [41](#bib.bib41){.ltx_ref}\] techniques to automatically generate multi-step reasoning processes for new tasks, while also selecting and utilizing the most appropriate available tool at each step. ASH \[[62](#bib.bib62){.ltx_ref}\] utilizes LLM for sequence hierarchical decision-making to achieve web navigation tasks. WebGPT \[[66](#bib.bib66){.ltx_ref}\] and WebCPM \[[64](#bib.bib64){.ltx_ref}\] use network search to assist in implementing Question Answering tasks. In addition, RCI \[[87](#bib.bib87){.ltx_ref}\] recursively criticizes and improves itself to execute computer tasks guided by natural language according to the prompting scheme. To achieve the analysis and processing of tables, TableGPT \[[71](#bib.bib71){.ltx_ref}\] employs a table encoder to transform tabular data into vector representations, which are then fed into an LLM for inference in combination with user queries.
:::
::::::

:::: {#S4.SS2 .section .ltx_subsection}
### [4.2 ]{.ltx_tag .ltx_tag_subsection}Tool Creation {#tool-creation .ltx_title .ltx_title_subsection}

::: {#S4.SS2.p1 .ltx_para}
The usage of tools is contingent upon the accessibility of external tools. Recently, efforts have been made to employ LLM as a tool creator in order to generate tools that can be utilized for diverse requests \[[88](#bib.bib88){.ltx_ref}, [89](#bib.bib89){.ltx_ref}, [90](#bib.bib90){.ltx_ref}, [91](#bib.bib91){.ltx_ref}, [92](#bib.bib92){.ltx_ref}, [93](#bib.bib93){.ltx_ref}, [94](#bib.bib94){.ltx_ref}, [95](#bib.bib95){.ltx_ref}\]. This development has consequently raised the demands placed on LLM. And these created tools are typically implemented as Python or SQL functions. LATM \[[88](#bib.bib88){.ltx_ref}\], for example, leverages the prowess of GPT-4 to create tools, and the usage of more cost-effective models has shown potential in exhibiting performance on par with larger models for these tool applications. EVAPORATE \[[94](#bib.bib94){.ltx_ref}\] involves the synthesis of multiple functions, which are subsequently utilized at a large scale to efficiently process documents and generate structured views.
:::
::::
::::::::::

:::: {#S5 .section .ltx_section}
## [5 ]{.ltx_tag .ltx_tag_section}Conclusion {#conclusion .ltx_title .ltx_title_section}

::: {#S5.p1 .ltx_para}
In this paper, we have introduced a structured framework specially designed for LLM-based AI Agents, with an emphasis on their abilities in task planning and tool usage. This framework, coupled with our design of two distinct types of agents assigned for the inference process, allows for a comprehensive evaluation of the capabilities of current open-source LLMs, thereby yielding critical insights into their effectiveness. Furthermore, our research highlights the significant potential of LLMs in managing complex tasks, revealing the exciting prospects they hold for future research and development. As we continue to explore and improve upon these models, we move closer to unlocking their full potential in a wide range of real-world applications.
:::
::::

::::::::: {#Sx1 .section .ltx_section}
## Acknowledgements {#acknowledgements .ltx_title .ltx_title_section}

::: {#Sx1.p1 .ltx_para}
This work was conducted collaboratively among the authors.
:::

::: {#Sx1.p2 .ltx_para}
Hangyu Mao and Rui Zhao led the project, formulating the central idea and laying out the framework for the primary literature review.
:::

::: {#Sx1.p3 .ltx_para}
Regarding the literature review phase, the surveys were conducted by various team members. Guoqing Du and Jingqing Ruan explored DNN-based Tool Scheduling by LLMs; Tianpeng Bao and Yihong Chen investigated Physical/Robot Tool Scheduling by LLMs; and Shiwei Shi and Zhiwei Xu handled the survey of API or GUI-based Tool Scheduling by LLMs. Bin Zhang summarized these papers and synthesized an overarching summary.
:::

::: {#Sx1.p4 .ltx_para}
As for the evaluation phase, Yihong Chen, Tianpeng Bao, Jingqing Ruan, Guoqing Du, Zhiwei Xu, Shiwei Shi, and Bin Zhang performed the experiments and analyzed the data. Hangyu Mao assisted in the analysis of the experimental phenomena and offered constructive suggestions for improvements. Xingyu Zeng and Rui Zhao provided invaluable feedback, contributed to the direction of the research. All authors participated in the discussion.
:::

::: {#Sx1.p5 .ltx_para}
Regarding the manuscript phase, Hangyu Mao organized the overall chapters of the manuscript and mainly wrote the methodology part, and provided assistance in other parts. Jingqing Ruan and Yihong Chen wrote the evaluation section. Bin Zhang wrote the summary of the literature review. Each author read and approved the final manuscript.
:::

::: {#Sx1.p6 .ltx_para}
The authors would like to thank Feng Zhu, Kun Wang, Yuhang Ran, Mengying Xu, Pengfei Jia, and Shaobo Lin for their valuable feedback, discussion, and participation in this project.
:::
:::::::::

::: {#bib .section .ltx_bibliography}
## References {#references .ltx_title .ltx_title_bibliography}

- [[\[1\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ W. X. Zhao, K. Zhou, J. Li, T. Tang, X. Wang, Y. Hou, Y. Min, B. Zhang, J. Zhang, Z. Dong *et al.*, "A survey of large language models," *arXiv preprint arXiv:2303.18223*, 2023. ]{.ltx_bibblock}]{#bib.bib1}
- [[\[2\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ T. Brown, B. Mann, N. Ryder, M. Subbiah, J. D. Kaplan, P. Dhariwal, A. Neelakantan, P. Shyam, G. Sastry, A. Askell *et al.*, "Language models are few-shot learners," *Advances in neural information processing systems*, vol. 33, pp. 1877--1901, 2020. ]{.ltx_bibblock}]{#bib.bib2}
- [[\[3\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ J. Wei, M. Bosma, V. Y. Zhao, K. Guu, A. W. Yu, B. Lester, N. Du, A. M. Dai, and Q. V. Le, "Finetuned language models are zero-shot learners," *arXiv preprint arXiv:2109.01652*, 2021. ]{.ltx_bibblock}]{#bib.bib3}
- [[\[4\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ OpenAI, "Gpt-4 technical report," 2023. ]{.ltx_bibblock}]{#bib.bib4}
- [[\[5\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ T. Schick, J. Dwivedi-Yu, R. Dessì, R. Raileanu, M. Lomeli, L. Zettlemoyer, N. Cancedda, and T. Scialom, "Toolformer: Language models can teach themselves to use tools," *arXiv preprint arXiv:2302.04761*, 2023. ]{.ltx_bibblock}]{#bib.bib5}
- [[\[6\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ N. R. Jennings, K. Sycara, and M. Wooldridge, "A roadmap of agent research and development," *Autonomous agents and multi-agent systems*, vol. 1, pp. 7--38, 1998. ]{.ltx_bibblock}]{#bib.bib6}
- [[\[7\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ N. R. Jennings and M. Wooldridge, "Applying agent technology," *Applied Artificial Intelligence an International Journal*, vol. 9, no. 4, pp. 357--369, 1995. ]{.ltx_bibblock}]{#bib.bib7}
- [[\[8\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ S. Franklin and A. Graesser, "Is it an agent, or just a program?: A taxonomy for autonomous agents," in *International workshop on agent theories, architectures, and languages*.   Springer, 1996, pp. 21--35. ]{.ltx_bibblock}]{#bib.bib8}
- [[\[9\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ C. Castelfranchi, "Modelling social action for ai agents," *Artificial intelligence*, vol. 103, no. 1-2, pp. 157--182, 1998. ]{.ltx_bibblock}]{#bib.bib9}
- [[\[10\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ J. Ferber and G. Weiss, *Multi-agent systems: an introduction to distributed artificial intelligence*.   Addison-wesley Reading, 1999, vol. 1. ]{.ltx_bibblock}]{#bib.bib10}
- [[\[11\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ L. Panait and S. Luke, "Cooperative multi-agent learning: The state of the art," *Autonomous agents and multi-agent systems*, vol. 11, pp. 387--434, 2005. ]{.ltx_bibblock}]{#bib.bib11}
- [[\[12\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ M. Pourreza and D. Rafiei, "Din-sql: Decomposed in-context learning of text-to-sql with self-correction," *arXiv preprint arXiv:2304.11015*, 2023. ]{.ltx_bibblock}]{#bib.bib12}
- [[\[13\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ C. Wu, S. Yin, W. Qi, X. Wang, Z. Tang, and N. Duan, "Visual chatgpt: Talking, drawing and editing with visual foundation models," *arXiv preprint arXiv:2303.04671*, 2023. ]{.ltx_bibblock}]{#bib.bib13}
- [[\[14\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ J. Gorniak, Y. Kim, S. Gwon, D. Wei, and N. W. Kim, "Vizability: Multimodal accessible data visualization with keyboard navigation and conversational interaction," *arXiv preprint arXiv:2310.09611*, 2023. ]{.ltx_bibblock}]{#bib.bib14}
- [[\[15\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ I. Team, "Internlm: A multilingual language model with progressively enhanced capabilities," [https://github.com/InternLM/InternLM](https://github.com/InternLM/InternLM){.ltx_ref .ltx_url .ltx_font_typewriter target="_blank"}, 2023. ]{.ltx_bibblock}]{#bib.bib15}
- [[\[16\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Q. Xu, F. Hong, B. Li, C. Hu, Z. Chen, and J. Zhang, "On the tool manipulation capability of open-source large language models," *arXiv preprint arXiv:2305.16504*, 2023. ]{.ltx_bibblock}]{#bib.bib16}
- [[\[17\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Y. Qin, S. Liang, Y. Ye, K. Zhu, L. Yan, Y. Lu, Y. Lin, X. Cong, X. Tang, B. Qian *et al.*, "Toolllm: Facilitating large language models to master 16000+ real-world apis," *arXiv preprint arXiv:2307.16789*, 2023. ]{.ltx_bibblock}]{#bib.bib17}
- [[\[18\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ M. Li, F. Song, B. Yu, H. Yu, Z. Li, F. Huang, and Y. Li, "Api-bank: A benchmark for tool-augmented llms," *arXiv preprint arXiv:2304.08244*, 2023. ]{.ltx_bibblock}]{#bib.bib18}
- [[\[19\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ S. G. Patil, T. Zhang, X. Wang, and J. E. Gonzalez, "Gorilla: Large language model connected with massive apis," *arXiv preprint arXiv:2305.15334*, 2023. ]{.ltx_bibblock}]{#bib.bib19}
- [[\[20\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Q. Tang, Z. Deng, H. Lin, X. Han, Q. Liang, and L. Sun, "Toolalpaca: Generalized tool learning for language models with 3000 simulated cases," *arXiv preprint arXiv:2306.05301*, 2023. ]{.ltx_bibblock}]{#bib.bib20}
- [[\[21\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ L. Ouyang, J. Wu, X. Jiang, D. Almeida, C. Wainwright, P. Mishkin, C. Zhang, S. Agarwal, K. Slama, A. Ray *et al.*, "Training language models to follow instructions with human feedback," *Advances in Neural Information Processing Systems*, vol. 35, pp. 27 730--27 744, 2022. ]{.ltx_bibblock}]{#bib.bib21}
- [[\[22\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Y. Bai, S. Kadavath, S. Kundu, A. Askell, J. Kernion, A. Jones, A. Chen, A. Goldie, A. Mirhoseini, C. McKinnon *et al.*, "Constitutional ai: Harmlessness from ai feedback," *arXiv preprint arXiv:2212.08073*, 2022. ]{.ltx_bibblock}]{#bib.bib22}
- [[\[23\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ A. Zeng, X. Liu, Z. Du, Z. Wang, H. Lai, M. Ding, Z. Yang, Y. Xu, W. Zheng, X. Xia *et al.*, "Glm-130b: An open bilingual pre-trained model," *arXiv preprint arXiv:2210.02414*, 2022. ]{.ltx_bibblock}]{#bib.bib23}
- [[\[24\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ H. Touvron, T. Lavril, G. Izacard, X. Martinet, M.-A. Lachaux, T. Lacroix, B. Rozière, N. Goyal, E. Hambro, F. Azhar *et al.*, "Llama: Open and efficient foundation language models," *arXiv preprint arXiv:2302.13971*, 2023. ]{.ltx_bibblock}]{#bib.bib24}
- [[\[25\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Y. Cui, Z. Yang, and X. Yao, "Efficient and effective text encoding for chinese llama and alpaca," *arXiv preprint arXiv:2304.08177*, 2023. ]{.ltx_bibblock}]{#bib.bib25}
- [[\[26\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ J. Wei, X. Wang, D. Schuurmans, M. Bosma, B. Ichter, F. Xia, E. H. Chi, Q. V. Le, and D. Zhou, "Chain-of-thought prompting elicits reasoning in large language models," *Neural Information Processing Systems*, 2022. ]{.ltx_bibblock}]{#bib.bib26}
- [[\[27\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ R. Bommasani, D. A. Hudson, E. Adeli, R. Altman, S. Arora, S. von Arx, M. S. Bernstein, J. Bohg, A. Bosselut, E. Brunskill *et al.*, "On the opportunities and risks of foundation models," *arXiv preprint arXiv:2108.07258*, 2021. ]{.ltx_bibblock}]{#bib.bib27}
- [[\[28\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ M. Mosbach, T. Pimentel, S. Ravfogel, D. Klakow, and Y. Elazar, "Few-shot fine-tuning vs. in-context learning: A fair comparison and evaluation," *arXiv preprint arXiv:2305.16938*, 2023. ]{.ltx_bibblock}]{#bib.bib28}
- [[\[29\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ J. Yang, H. Jin, R. Tang, X. Han, Q. Feng, H. Jiang, B. Yin, and X. Hu, "Harnessing the power of llms in practice: A survey on chatgpt and beyond," *arXiv preprint arXiv:2304.13712*, 2023. ]{.ltx_bibblock}]{#bib.bib29}
- [[\[30\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ C. Zhang, C. Zhang, C. Li, Y. Qiao, S. Zheng, S. K. Dam, M. Zhang, J. U. Kim, S. T. Kim, J. Choi *et al.*, "One small step for generative ai, one giant leap for agi: A complete survey on chatgpt in aigc era," *arXiv preprint arXiv:2304.06488*, 2023. ]{.ltx_bibblock}]{#bib.bib30}
- [[\[31\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ F. Yu, H. Zhang, and B. Wang, "Nature language reasoning, a survey," *arXiv preprint arXiv:2303.14725*, 2023. ]{.ltx_bibblock}]{#bib.bib31}
- [[\[32\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Z. Wang, G. Zhang, K. Yang, N. Shi, W. Zhou, S. Hao, G. Xiong, Y. Li, M. Y. Sim, X. Chen *et al.*, "Interactive natural language processing," *arXiv preprint arXiv:2305.13246*, 2023. ]{.ltx_bibblock}]{#bib.bib32}
- [[\[33\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Y. Qin, S. Hu, Y. Lin, W. Chen, N. Ding, G. Cui, Z. Zeng, Y. Huang, C. Xiao, C. Han *et al.*, "Tool learning with foundation models," *arXiv preprint arXiv:2304.08354*, 2023. ]{.ltx_bibblock}]{#bib.bib33}
- [[\[34\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ W. Yu, C. Zhu, Z. Li, Z. Hu, Q. Wang, H. Ji, and M. Jiang, "A survey of knowledge-enhanced text generation," *ACM Computing Surveys*, vol. 54, no. 11s, pp. 1--38, 2022. ]{.ltx_bibblock}]{#bib.bib34}
- [[\[35\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ E. J. Hu, Y. Shen, P. Wallis, Z. Allen-Zhu, Y. Li, S. Wang, L. Wang, and W. Chen, "Lora: Low-rank adaptation of large language models," *arXiv preprint arXiv:2106.09685*, 2021. ]{.ltx_bibblock}]{#bib.bib35}
- [[\[36\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ N. Houlsby, A. Giurgiu, S. Jastrzebski, B. Morrone, Q. De Laroussilhe, A. Gesmundo, M. Attariyan, and S. Gelly, "Parameter-efficient transfer learning for nlp," in *International Conference on Machine Learning*.   PMLR, 2019, pp. 2790--2799. ]{.ltx_bibblock}]{#bib.bib36}
- [[\[37\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ X. L. Li and P. Liang, "Prefix-tuning: Optimizing continuous prompts for generation," *arXiv preprint arXiv:2101.00190*, 2021. ]{.ltx_bibblock}]{#bib.bib37}
- [[\[38\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ X. Liu, Y. Zheng, Z. Du, M. Ding, Y. Qian, Z. Yang, and J. Tang, "Gpt understands, too," *arXiv preprint arXiv:2103.10385*, 2021. ]{.ltx_bibblock}]{#bib.bib38}
- [[\[39\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ S. Yao, J. Zhao, D. Yu, N. Du, I. Shafran, K. Narasimhan, and Y. Cao, "React: Synergizing reasoning and acting in language models," *arXiv preprint arXiv:2210.03629*, 2022. ]{.ltx_bibblock}]{#bib.bib39}
- [[\[40\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ T. Khot, H. Trivedi, M. Finlayson, Y. Fu, K. Richardson, P. Clark, and A. Sabharwal, "Decomposed prompting: A modular approach for solving complex tasks," *arXiv preprint arXiv:2210.02406*, 2022. ]{.ltx_bibblock}]{#bib.bib40}
- [[\[41\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ G. Mialon, R. Dessì, M. Lomeli, C. Nalmpantis, R. Pasunuru, R. Raileanu, B. Rozière, T. Schick, J. Dwivedi-Yu, A. Celikyilmaz *et al.*, "Augmented language models: a survey," *arXiv preprint arXiv:2302.07842*, 2023. ]{.ltx_bibblock}]{#bib.bib41}
- [[\[42\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ J. Duan, S. Yu, H. L. Tan, H. Zhu, and C. Tan, "A survey of embodied ai: From simulators to research tasks," *IEEE Transactions on Emerging Topics in Computational Intelligence*, vol. 6, no. 2, pp. 230--244, 2022. ]{.ltx_bibblock}]{#bib.bib42}
- [[\[43\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ M. Savva, A. Kadian, O. Maksymets, Y. Zhao, E. Wijmans, B. Jain, J. Straub, J. Liu, V. Koltun, J. Malik *et al.*, "Habitat: A platform for embodied ai research," in *Proceedings of the IEEE/CVF international conference on computer vision*, 2019, pp. 9339--9347. ]{.ltx_bibblock}]{#bib.bib43}
- [[\[44\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ S. Franklin, "Autonomous agents as embodied ai," *Cybernetics & Systems*, vol. 28, no. 6, pp. 499--520, 1997. ]{.ltx_bibblock}]{#bib.bib44}
- [[\[45\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ W. Zhang, Y. Guo, L. Niu, P. Li, C. Zhang, Z. Wan, J. Yan, F. U. D. Farrukh, and D. Zhang, "Lp-slam: Language-perceptive rgb-d slam system based on large language model," *arXiv preprint arXiv:2303.10089*, 2023. ]{.ltx_bibblock}]{#bib.bib45}
- [[\[46\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ D. Shah, B. Osiński, S. Levine *et al.*, "Lm-nav: Robotic navigation with large pre-trained models of language, vision, and action," in *Conference on Robot Learning*.   PMLR, 2023, pp. 492--504. ]{.ltx_bibblock}]{#bib.bib46}
- [[\[47\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ A. Brohan, Y. Chebotar, C. Finn, K. Hausman, A. Herzog, D. Ho, J. Ibarz, A. Irpan, E. Jang, R. Julian *et al.*, "Do as i can, not as i say: Grounding language in robotic affordances," in *Conference on Robot Learning*.   PMLR, 2023, pp. 287--318. ]{.ltx_bibblock}]{#bib.bib47}
- [[\[48\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ W. Huang, F. Xia, T. Xiao, H. Chan, J. Liang, P. Florence, A. Zeng, J. Tompson, I. Mordatch, Y. Chebotar *et al.*, "Inner monologue: Embodied reasoning through planning with language models," *arXiv preprint arXiv:2207.05608*, 2022. ]{.ltx_bibblock}]{#bib.bib48}
- [[\[49\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ B. Chen, F. Xia, B. Ichter, K. Rao, K. Gopalakrishnan, M. S. Ryoo, A. Stone, and D. Kappler, "Open-vocabulary queryable scene representations for real world planning," in *2023 IEEE International Conference on Robotics and Automation (ICRA)*.   IEEE, 2023, pp. 11 509--11 522. ]{.ltx_bibblock}]{#bib.bib49}
- [[\[50\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ D. Driess, F. Xia, M. S. Sajjadi, C. Lynch, A. Chowdhery, B. Ichter, A. Wahid, J. Tompson, Q. Vuong, T. Yu *et al.*, "Palm-e: An embodied multimodal language model," *arXiv preprint arXiv:2303.03378*, 2023. ]{.ltx_bibblock}]{#bib.bib50}
- [[\[51\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ N. Wake, A. Kanehira, K. Sasabuchi, J. Takamatsu, and K. Ikeuchi, "Chatgpt empowered long-step robot control in various environments: A case application," *arXiv preprint arXiv:2304.03893*, 2023. ]{.ltx_bibblock}]{#bib.bib51}
- [[\[52\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ K. Rana, J. Haviland, S. Garg, J. Abou-Chakra, I. Reid, and N. Suenderhauf, "Sayplan: Grounding large language models using 3d scene graphs for scalable task planning," *arXiv preprint arXiv:2307.06135*, 2023. ]{.ltx_bibblock}]{#bib.bib52}
- [[\[53\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ C. H. Song, J. Wu, C. Washington, B. M. Sadler, W.-L. Chao, and Y. Su, "Llm-planner: Few-shot grounded planning for embodied agents with large language models," *arXiv preprint arXiv:2212.04088*, 2022. ]{.ltx_bibblock}]{#bib.bib53}
- [[\[54\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ A. Brohan, N. Brown, J. Carbajal, Y. Chebotar, J. Dabis, C. Finn, K. Gopalakrishnan, K. Hausman, A. Herzog, J. Hsu *et al.*, "Rt-1: Robotics transformer for real-world control at scale," *arXiv preprint arXiv:2212.06817*, 2022. ]{.ltx_bibblock}]{#bib.bib54}
- [[\[55\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ A. Stone, T. Xiao, Y. Lu, K. Gopalakrishnan, K.-H. Lee, Q. Vuong, P. Wohlhart, B. Zitkovich, F. Xia, C. Finn *et al.*, "Open-world object manipulation using pre-trained vision-language models," *arXiv preprint arXiv:2303.00905*, 2023. ]{.ltx_bibblock}]{#bib.bib55}
- [[\[56\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ S. Reed, K. Zolna, E. Parisotto, S. G. Colmenarejo, A. Novikov, G. Barth-Maron, M. Gimenez, Y. Sulsky, J. Kay, J. T. Springenberg *et al.*, "A generalist agent," *arXiv preprint arXiv:2205.06175*, 2022. ]{.ltx_bibblock}]{#bib.bib56}
- [[\[57\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ S. Vemprala, R. Bonatti, A. Bucker, and A. Kapoor, "Chatgpt for robotics: Design principles and model abilities," *Microsoft Auton. Syst. Robot. Res*, vol. 2, p. 20, 2023. ]{.ltx_bibblock}]{#bib.bib57}
- [[\[58\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ J. Liang, W. Huang, F. Xia, P. Xu, K. Hausman, B. Ichter, P. Florence, and A. Zeng, "Code as policies: Language model programs for embodied control," in *2023 IEEE International Conference on Robotics and Automation (ICRA)*.   IEEE, 2023, pp. 9493--9500. ]{.ltx_bibblock}]{#bib.bib58}
- [[\[59\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ K. Guu, K. Lee, Z. Tung, P. Pasupat, and M. Chang, "Retrieval augmented language model pre-training," in *International conference on machine learning*.   PMLR, 2020, pp. 3929--3938. ]{.ltx_bibblock}]{#bib.bib59}
- [[\[60\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ P. Lewis, E. Perez, A. Piktus, F. Petroni, V. Karpukhin, N. Goyal, H. Küttler, M. Lewis, W.-t. Yih, T. Rocktäschel *et al.*, "Retrieval-augmented generation for knowledge-intensive nlp tasks," *Advances in Neural Information Processing Systems*, vol. 33, pp. 9459--9474, 2020. ]{.ltx_bibblock}]{#bib.bib60}
- [[\[61\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ S. Borgeaud, A. Mensch, J. Hoffmann, T. Cai, E. Rutherford, K. Millican, G. B. Van Den Driessche, J.-B. Lespiau, B. Damoc, A. Clark *et al.*, "Improving language models by retrieving from trillions of tokens," in *International conference on machine learning*.   PMLR, 2022, pp. 2206--2240. ]{.ltx_bibblock}]{#bib.bib61}
- [[\[62\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ A. Sridhar, R. Lo, F. F. Xu, H. Zhu, and S. Zhou, "Hierarchical prompting assists large language model on web navigation," *arXiv preprint arXiv:2305.14257*, 2023. ]{.ltx_bibblock}]{#bib.bib62}
- [[\[63\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ H. Furuta, O. Nachum, K.-H. Lee, Y. Matsuo, S. S. Gu, and I. Gur, "Multimodal web navigation with instruction-finetuned foundation models," *arXiv preprint arXiv:2305.11854*, 2023. ]{.ltx_bibblock}]{#bib.bib63}
- [[\[64\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Y. Qin, Z. Cai, D. Jin, L. Yan, S. Liang, K. Zhu, Y. Lin, X. Han, N. Ding, H. Wang *et al.*, "Webcpm: Interactive web search for chinese long-form question answering," *arXiv preprint arXiv:2305.06849*, 2023. ]{.ltx_bibblock}]{#bib.bib64}
- [[\[65\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ S. Yao, H. Chen, J. Yang, and K. Narasimhan, "Webshop: Towards scalable real-world web interaction with grounded language agents," *Advances in Neural Information Processing Systems*, vol. 35, pp. 20 744--20 757, 2022. ]{.ltx_bibblock}]{#bib.bib65}
- [[\[66\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ R. Nakano, J. Hilton, S. Balaji, J. Wu, L. Ouyang, C. Kim, C. Hesse, S. Jain, V. Kosaraju, W. Saunders *et al.*, "Webgpt: Browser-assisted question-answering with human feedback," *arXiv preprint arXiv:2112.09332*, 2021. ]{.ltx_bibblock}]{#bib.bib66}
- [[\[67\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Z. Yang, P. Qi, S. Zhang, Y. Bengio, W. W. Cohen, R. Salakhutdinov, and C. D. Manning, "Hotpotqa: A dataset for diverse, explainable multi-hop question answering," *arXiv preprint arXiv:1809.09600*, 2018. ]{.ltx_bibblock}]{#bib.bib67}
- [[\[68\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ B. Wang, G. Li, and Y. Li, "Enabling conversational interaction with mobile ui using large language models," in *Proceedings of the 2023 CHI Conference on Human Factors in Computing Systems*, 2023, pp. 1--17. ]{.ltx_bibblock}]{#bib.bib68}
- [[\[69\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ D. Zhang, L. Chen, and K. Yu, "Mobile-env: A universal platform for training and evaluation of mobile interaction," *arXiv preprint arXiv:2305.08144*, 2023. ]{.ltx_bibblock}]{#bib.bib69}
- [[\[70\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ H. Li, J. Su, Y. Chen, Q. Li, and Z. Zhang, "Sheetcopilot: Bringing software productivity to the next level through large language models," *arXiv preprint arXiv:2305.19308*, 2023. ]{.ltx_bibblock}]{#bib.bib70}
- [[\[71\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ L. Zha, J. Zhou, L. Li, R. Wang, Q. Huang, S. Yang, J. Yuan, C. Su, X. Li, A. Su *et al.*, "Tablegpt: Towards unifying tables, nature language and commands into one gpt," *arXiv preprint arXiv:2307.08674*, 2023. ]{.ltx_bibblock}]{#bib.bib71}
- [[\[72\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Z. Chen, K. Zhou, B. Zhang, Z. Gong, W. X. Zhao, and J.-R. Wen, "Chatcot: Tool-augmented chain-of-thought reasoning on$\backslash$$\backslash$chat-based large language models," *arXiv preprint arXiv:2305.14323*, 2023. ]{.ltx_bibblock}]{#bib.bib72}
- [[\[73\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ A. Parisi, Y. Zhao, and N. Fiedel, "Talm: Tool augmented language models," *arXiv preprint arXiv:2205.12255*, 2022. ]{.ltx_bibblock}]{#bib.bib73}
- [[\[74\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ K. Cobbe, V. Kosaraju, M. Bavarian, M. Chen, H. Jun, L. Kaiser, M. Plappert, J. Tworek, J. Hilton, R. Nakano *et al.*, "Training verifiers to solve math word problems," *arXiv preprint arXiv:2110.14168*, 2021. ]{.ltx_bibblock}]{#bib.bib74}
- [[\[75\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Z. Yang, L. Li, J. Wang, K. Lin, E. Azarnasab, F. Ahmed, Z. Liu, C. Liu, M. Zeng, and L. Wang, "Mm-react: Prompting chatgpt for multimodal reasoning and action," *arXiv preprint arXiv:2303.11381*, 2023. ]{.ltx_bibblock}]{#bib.bib75}
- [[\[76\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Z. Liu, Y. He, W. Wang, W. Wang, Y. Wang, S. Chen, Q. Zhang, Y. Yang, Q. Li, J. Yu *et al.*, "Internchat: Solving vision-centric tasks by interacting with chatbots beyond language," *arXiv preprint arXiv:2305.05662*, 2023. ]{.ltx_bibblock}]{#bib.bib76}
- [[\[77\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Y. Ge, W. Hua, J. Ji, J. Tan, S. Xu, and Y. Zhang, "Openagi: When llm meets domain experts," *arXiv preprint arXiv:2304.04370*, 2023. ]{.ltx_bibblock}]{#bib.bib77}
- [[\[78\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Y. Shen, K. Song, X. Tan, D. Li, W. Lu, and Y. Zhuang, "Hugginggpt: Solving ai tasks with chatgpt and its friends in huggingface," *arXiv preprint arXiv:2303.17580*, 2023. ]{.ltx_bibblock}]{#bib.bib78}
- [[\[79\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ D. Surís, S. Menon, and C. Vondrick, "Vipergpt: Visual inference via python execution for reasoning," *arXiv preprint arXiv:2303.08128*, 2023. ]{.ltx_bibblock}]{#bib.bib79}
- [[\[80\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ T. Gupta and A. Kembhavi, "Visual programming: Compositional visual reasoning without training," in *Proceedings of the IEEE/CVF Conference on Computer Vision and Pattern Recognition*, 2023, pp. 14 953--14 962. ]{.ltx_bibblock}]{#bib.bib80}
- [[\[81\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ L. Chen, B. Li, S. Shen, J. Yang, C. Li, K. Keutzer, T. Darrell, and Z. Liu, "Language models are visual reasoning coordinators," in *ICLR 2023 Workshop on Mathematical and Empirical Understanding of Foundation Models*, 2023. ]{.ltx_bibblock}]{#bib.bib81}
- [[\[82\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ P. Lu, B. Peng, H. Cheng, M. Galley, K.-W. Chang, Y. N. Wu, S.-C. Zhu, and J. Gao, "Chameleon: Plug-and-play compositional reasoning with large language models," *arXiv preprint arXiv:2304.09842*, 2023. ]{.ltx_bibblock}]{#bib.bib82}
- [[\[83\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Z. Gou, Z. Shao, Y. Gong, Y. Shen, Y. Yang, N. Duan, and W. Chen, "Critic: Large language models can self-correct with tool-interactive critiquing," *arXiv preprint arXiv:2305.11738*, 2023. ]{.ltx_bibblock}]{#bib.bib83}
- [[\[84\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Y. Liang, C. Wu, T. Song, W. Wu, Y. Xia, Y. Liu, Y. Ou, S. Lu, L. Ji, S. Mao *et al.*, "Taskmatrix. ai: Completing tasks by connecting foundation models with millions of apis," *arXiv preprint arXiv:2303.16434*, 2023. ]{.ltx_bibblock}]{#bib.bib84}
- [[\[85\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ S. Hao, T. Liu, Z. Wang, and Z. Hu, "Toolkengpt: Augmenting frozen language models with massive tools via tool embeddings," *arXiv preprint arXiv:2305.11554*, 2023. ]{.ltx_bibblock}]{#bib.bib85}
- [[\[86\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ B. Paranjape, S. Lundberg, S. Singh, H. Hajishirzi, L. Zettlemoyer, and M. T. Ribeiro, "Art: Automatic multi-step reasoning and tool-use for large language models," *arXiv preprint arXiv:2303.09014*, 2023. ]{.ltx_bibblock}]{#bib.bib86}
- [[\[87\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ G. Kim, P. Baldi, and S. McAleer, "Language models can solve computer tasks," *arXiv preprint arXiv:2303.17491*, 2023. ]{.ltx_bibblock}]{#bib.bib87}
- [[\[88\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ T. Cai, X. Wang, T. Ma, X. Chen, and D. Zhou, "Large language models as tool makers," *arXiv preprint arXiv:2305.17126*, 2023. ]{.ltx_bibblock}]{#bib.bib88}
- [[\[89\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ R. H. Lewis and J. Jiao, "Computegpt: A computational chat model for numerical problems," *arXiv preprint arXiv:2305.06223*, 2023. ]{.ltx_bibblock}]{#bib.bib89}
- [[\[90\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ L. Gao, A. Madaan, S. Zhou, U. Alon, P. Liu, Y. Yang, J. Callan, and G. Neubig, "Pal: Program-aided language models," in *International Conference on Machine Learning*.   PMLR, 2023, pp. 10 764--10 799. ]{.ltx_bibblock}]{#bib.bib90}
- [[\[91\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ G. Wang, Y. Xie, Y. Jiang, A. Mandlekar, C. Xiao, Y. Zhu, L. Fan, and A. Anandkumar, "Voyager: An open-ended embodied agent with large language models," *arXiv preprint arXiv:2305.16291*, 2023. ]{.ltx_bibblock}]{#bib.bib91}
- [[\[92\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ C. Qian, C. Han, Y. R. Fung, Y. Qin, Z. Liu, and H. Ji, "Creator: Disentangling abstract and concrete reasonings of large language models through tool creation," *arXiv preprint arXiv:2305.14318*, 2023. ]{.ltx_bibblock}]{#bib.bib92}
- [[\[93\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Y. Cai, S. Mao, W. Wu, Z. Wang, Y. Liang, T. Ge, C. Wu, W. You, T. Song, Y. Xia *et al.*, "Low-code llm: Visual programming over llms," *arXiv preprint arXiv:2304.08103*, 2023. ]{.ltx_bibblock}]{#bib.bib93}
- [[\[94\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ S. Arora, B. Yang, S. Eyuboglu, A. Narayan, A. Hojel, I. Trummer, and C. Ré, "Language models enable simple systems for generating structured views of heterogeneous data lakes," *arXiv preprint arXiv:2304.09433*, 2023. ]{.ltx_bibblock}]{#bib.bib94}
- [[\[95\]]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ W. Zhang, Y. Shen, W. Lu, and Y. Zhuang, "Data-copilot: Bridging billions of data and humans with autonomous workflow," *arXiv preprint arXiv:2306.07209*, 2023. ]{.ltx_bibblock}]{#bib.bib95}
:::

::: {.ltx_pagination .ltx_role_newpage}
:::

::::::::: {#A1 .section .ltx_appendix}
## [Appendix A ]{.ltx_tag .ltx_tag_appendix}Detailed Dataset Description {#appendix-a-detailed-dataset-description .ltx_title .ltx_title_appendix}

::: {#A1.p1 .ltx_para}
[Simple SQL queries]{#A1.p1.1.1 .ltx_text .ltx_font_bold}: These queries typically involve basic operations such as SELECT, FROM, WHERE, GROUP BY, etc. They are used to retrieve, filter, group, and sort data from a single table. We give the Schema of two tables in the SQL database in Table [[13]{.ltx_text .ltx_ref_tag}](#A1.T13 "Table 13 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} and [[13]{.ltx_text .ltx_ref_tag}](#A1.T13 "Table 13 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} and list several examples in Table [[14]{.ltx_text .ltx_ref_tag}](#A1.T14 "Table 14 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::

<figure id="A1.T13" class="ltx_table">
<div class="ltx_flex_figure ltx_flex_table">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="A1.T13.fig1" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_middle" style="width:216.8pt;">
<table id="A1.T13.fig1.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A1.T13.fig1.1.1" class="ltx_tr">
<td colspan="2" id="A1.T13.fig1.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T13.fig1.1.1.1.1" class="ltx_text ltx_font_bold">Person</span></td>
</tr>
<tr id="A1.T13.fig1.1.2" class="ltx_tr">
<td id="A1.T13.fig1.1.2.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T13.fig1.1.2.1.1" class="ltx_text ltx_font_bold">Column Name</span></td>
<td id="A1.T13.fig1.1.2.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T13.fig1.1.2.2.1" class="ltx_text ltx_font_bold">Type</span></td>
</tr>
<tr id="A1.T13.fig1.1.3" class="ltx_tr">
<td id="A1.T13.fig1.1.3.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">id</td>
<td id="A1.T13.fig1.1.3.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T13.fig1.1.4" class="ltx_tr">
<td id="A1.T13.fig1.1.4.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">name</td>
<td id="A1.T13.fig1.1.4.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T13.fig1.1.5" class="ltx_tr">
<td id="A1.T13.fig1.1.5.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">age</td>
<td id="A1.T13.fig1.1.5.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
<tr id="A1.T13.fig1.1.6" class="ltx_tr">
<td id="A1.T13.fig1.1.6.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">sex</td>
<td id="A1.T13.fig1.1.6.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T13.fig1.1.7" class="ltx_tr">
<td id="A1.T13.fig1.1.7.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">school</td>
<td id="A1.T13.fig1.1.7.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T13.fig1.1.8" class="ltx_tr">
<td id="A1.T13.fig1.1.8.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">phone</td>
<td id="A1.T13.fig1.1.8.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T13.fig1.1.9" class="ltx_tr">
<td id="A1.T13.fig1.1.9.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">qualifications</td>
<td id="A1.T13.fig1.1.9.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T13.fig1.1.10" class="ltx_tr">
<td id="A1.T13.fig1.1.10.1" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">ability</td>
<td id="A1.T13.fig1.1.10.2" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_figure">Table 12: </span>Schema of the Person table</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="A1.T13.fig2" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_middle" style="width:216.8pt;">
<table id="A1.T13.fig2.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A1.T13.fig2.1.1" class="ltx_tr">
<td colspan="2" id="A1.T13.fig2.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T13.fig2.1.1.1.1" class="ltx_text ltx_font_bold">School</span></td>
</tr>
<tr id="A1.T13.fig2.1.2" class="ltx_tr">
<td id="A1.T13.fig2.1.2.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T13.fig2.1.2.1.1" class="ltx_text ltx_font_bold">Column Name</span></td>
<td id="A1.T13.fig2.1.2.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T13.fig2.1.2.2.1" class="ltx_text ltx_font_bold">Type</span></td>
</tr>
<tr id="A1.T13.fig2.1.3" class="ltx_tr">
<td id="A1.T13.fig2.1.3.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">id</td>
<td id="A1.T13.fig2.1.3.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T13.fig2.1.4" class="ltx_tr">
<td id="A1.T13.fig2.1.4.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">name</td>
<td id="A1.T13.fig2.1.4.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T13.fig2.1.5" class="ltx_tr">
<td id="A1.T13.fig2.1.5.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">info_985</td>
<td id="A1.T13.fig2.1.5.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T13.fig2.1.6" class="ltx_tr">
<td id="A1.T13.fig2.1.6.1" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">info_211</td>
<td id="A1.T13.fig2.1.6.2" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_figure">Table 13: </span>Schema of the School table</figcaption>
</figure>
</div>
</div>
</figure>

<figure id="A1.T14" class="ltx_table">
<table id="A1.T14.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A1.T14.1.1" class="ltx_tr">
<td id="A1.T14.1.1.1" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.1.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.1.1.1.1" class="ltx_p" style="width:71.1pt;"><span id="A1.T14.1.1.1.1.1.1" class="ltx_text ltx_font_bold">Table ID</span></span> </span></td>
<td id="A1.T14.1.1.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.1.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.1.2.1.1" class="ltx_p" style="width:71.1pt;"><span id="A1.T14.1.1.2.1.1.1" class="ltx_text ltx_font_bold">Question</span></span> </span></td>
<td id="A1.T14.1.1.3" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.1.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.1.3.1.1" class="ltx_p" style="width:42.7pt;"><span id="A1.T14.1.1.3.1.1.1" class="ltx_text ltx_font_bold">Answer</span></span> </span></td>
<td id="A1.T14.1.1.4" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.1.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.1.4.1.1" class="ltx_p" style="width:142.3pt;"><span id="A1.T14.1.1.4.1.1.1" class="ltx_text ltx_font_bold">SQL reference</span></span> </span></td>
</tr>
<tr id="A1.T14.1.2" class="ltx_tr">
<td id="A1.T14.1.2.1" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.2.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.2.1.1.1" class="ltx_p" style="width:71.1pt;">Person</span> </span></td>
<td id="A1.T14.1.2.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.2.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.2.2.1.1" class="ltx_p" style="width:71.1pt;">Average ages</span> </span></td>
<td id="A1.T14.1.2.3" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.2.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.2.3.1.1" class="ltx_p" style="width:42.7pt;">35.16</span> </span></td>
<td id="A1.T14.1.2.4" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.2.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.2.4.1.1" class="ltx_p" style="width:142.3pt;">select avg(age) from Person</span> </span></td>
</tr>
<tr id="A1.T14.1.3" class="ltx_tr">
<td id="A1.T14.1.3.1" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.3.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.3.1.1.1" class="ltx_p" style="width:71.1pt;">Person</span> </span></td>
<td id="A1.T14.1.3.2" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.3.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.3.2.1.1" class="ltx_p" style="width:71.1pt;">How many men</span> </span></td>
<td id="A1.T14.1.3.3" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.3.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.3.3.1.1" class="ltx_p" style="width:42.7pt;">12</span> </span></td>
<td id="A1.T14.1.3.4" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.3.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.3.4.1.1" class="ltx_p" style="width:142.3pt;">select count(*) from Person where sex = ’male’</span> </span></td>
</tr>
<tr id="A1.T14.1.4" class="ltx_tr">
<td id="A1.T14.1.4.1" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.4.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.4.1.1.1" class="ltx_p" style="width:71.1pt;">School</span> </span></td>
<td id="A1.T14.1.4.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.4.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.4.2.1.1" class="ltx_p" style="width:71.1pt;">How many schools are both ’985’ and ’211’ institutions?</span> </span></td>
<td id="A1.T14.1.4.3" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.4.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.4.3.1.1" class="ltx_p" style="width:42.7pt;">11</span> </span></td>
<td id="A1.T14.1.4.4" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T14.1.4.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T14.1.4.4.1.1" class="ltx_p" style="width:142.3pt;">select count(*) from School where info_985 = ’yes’ and info_211 = ’yes’;</span> </span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 14: </span>Demostrations of simple SQL queries.</figcaption>
</figure>

::: {#A1.p2 .ltx_para}
[Complex nested SQL queries]{#A1.p2.1.1 .ltx_text .ltx_font_bold}: These queries contain subqueries, which are SQL queries nested inside a larger query. Nested queries can be used in various clauses such as SELECT, FROM, WHERE, and HAVING. They provide a way to perform multiple operations or calculations across multiple tables. We give the Schema of two tables in the SQL database in Table [[16]{.ltx_text .ltx_ref_tag}](#A1.T16 "Table 16 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, [[16]{.ltx_text .ltx_ref_tag}](#A1.T16 "Table 16 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, [[18]{.ltx_text .ltx_ref_tag}](#A1.T18 "Table 18 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, and [[18]{.ltx_text .ltx_ref_tag}](#A1.T18 "Table 18 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} and list several examples in Table [[19]{.ltx_text .ltx_ref_tag}](#A1.T19 "Table 19 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}.
:::

<figure id="A1.T16" class="ltx_table">
<div class="ltx_flex_figure ltx_flex_table">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="A1.T16.fig1" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_middle" style="width:216.8pt;">
<table id="A1.T16.fig1.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A1.T16.fig1.1.1" class="ltx_tr">
<td colspan="2" id="A1.T16.fig1.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T16.fig1.1.1.1.1" class="ltx_text ltx_font_bold">GoldenMelodyAwards</span></td>
</tr>
<tr id="A1.T16.fig1.1.2" class="ltx_tr">
<td id="A1.T16.fig1.1.2.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T16.fig1.1.2.1.1" class="ltx_text ltx_font_bold">Column Name</span></td>
<td id="A1.T16.fig1.1.2.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T16.fig1.1.2.2.1" class="ltx_text ltx_font_bold">Type</span></td>
</tr>
<tr id="A1.T16.fig1.1.3" class="ltx_tr">
<td id="A1.T16.fig1.1.3.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Nominated_Count</td>
<td id="A1.T16.fig1.1.3.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
<tr id="A1.T16.fig1.1.4" class="ltx_tr">
<td id="A1.T16.fig1.1.4.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Competing_Count</td>
<td id="A1.T16.fig1.1.4.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
<tr id="A1.T16.fig1.1.5" class="ltx_tr">
<td id="A1.T16.fig1.1.5.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Awards_Count</td>
<td id="A1.T16.fig1.1.5.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
<tr id="A1.T16.fig1.1.6" class="ltx_tr">
<td id="A1.T16.fig1.1.6.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Award_Name</td>
<td id="A1.T16.fig1.1.6.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T16.fig1.1.7" class="ltx_tr">
<td id="A1.T16.fig1.1.7.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Host</td>
<td id="A1.T16.fig1.1.7.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T16.fig1.1.8" class="ltx_tr">
<td id="A1.T16.fig1.1.8.1" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Year</td>
<td id="A1.T16.fig1.1.8.2" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TIME</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_figure">Table 15: </span>Schema of GoldenMelodyAwards</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="A1.T16.fig2" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_middle" style="width:216.8pt;">
<table id="A1.T16.fig2.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A1.T16.fig2.1.1" class="ltx_tr">
<td colspan="2" id="A1.T16.fig2.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T16.fig2.1.1.1.1" class="ltx_text ltx_font_bold">AwardNominees</span></td>
</tr>
<tr id="A1.T16.fig2.1.2" class="ltx_tr">
<td id="A1.T16.fig2.1.2.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T16.fig2.1.2.1.1" class="ltx_text ltx_font_bold">Column Name</span></td>
<td id="A1.T16.fig2.1.2.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T16.fig2.1.2.2.1" class="ltx_text ltx_font_bold">Type</span></td>
</tr>
<tr id="A1.T16.fig2.1.3" class="ltx_tr">
<td id="A1.T16.fig2.1.3.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Singer_ID</td>
<td id="A1.T16.fig2.1.3.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
<tr id="A1.T16.fig2.1.4" class="ltx_tr">
<td id="A1.T16.fig2.1.4.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Nominated_Work</td>
<td id="A1.T16.fig2.1.4.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T16.fig2.1.5" class="ltx_tr">
<td id="A1.T16.fig2.1.5.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Award_Name</td>
<td id="A1.T16.fig2.1.5.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T16.fig2.1.6" class="ltx_tr">
<td id="A1.T16.fig2.1.6.1" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Award_Edition_ID</td>
<td id="A1.T16.fig2.1.6.2" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_figure">Table 16: </span>Schema of the AwardNominees table</figcaption>
</figure>
</div>
</div>
</figure>

<figure id="A1.T18" class="ltx_table">
<div class="ltx_flex_figure ltx_flex_table">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="A1.T18.fig1" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_middle" style="width:216.8pt;">
<table id="A1.T18.fig1.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A1.T18.fig1.1.1" class="ltx_tr">
<td colspan="2" id="A1.T18.fig1.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T18.fig1.1.1.1.1" class="ltx_text ltx_font_bold">Singers</span></td>
</tr>
<tr id="A1.T18.fig1.1.2" class="ltx_tr">
<td id="A1.T18.fig1.1.2.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T18.fig1.1.2.1.1" class="ltx_text ltx_font_bold">Column Name</span></td>
<td id="A1.T18.fig1.1.2.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T18.fig1.1.2.2.1" class="ltx_text ltx_font_bold">Type</span></td>
</tr>
<tr id="A1.T18.fig1.1.3" class="ltx_tr">
<td id="A1.T18.fig1.1.3.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Name</td>
<td id="A1.T18.fig1.1.3.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T18.fig1.1.4" class="ltx_tr">
<td id="A1.T18.fig1.1.4.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Song_Count</td>
<td id="A1.T18.fig1.1.4.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
<tr id="A1.T18.fig1.1.5" class="ltx_tr">
<td id="A1.T18.fig1.1.5.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Album_Count</td>
<td id="A1.T18.fig1.1.5.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
<tr id="A1.T18.fig1.1.6" class="ltx_tr">
<td id="A1.T18.fig1.1.6.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Fan_Count</td>
<td id="A1.T18.fig1.1.6.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
<tr id="A1.T18.fig1.1.7" class="ltx_tr">
<td id="A1.T18.fig1.1.7.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Gender</td>
<td id="A1.T18.fig1.1.7.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T18.fig1.1.8" class="ltx_tr">
<td id="A1.T18.fig1.1.8.1" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Singer_ID</td>
<td id="A1.T18.fig1.1.8.2" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_figure">Table 17: </span>Schema of the Singers table</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="A1.T18.fig2" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_middle" style="width:216.8pt;">
<table id="A1.T18.fig2.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A1.T18.fig2.1.1" class="ltx_tr">
<td colspan="2" id="A1.T18.fig2.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T18.fig2.1.1.1.1" class="ltx_text ltx_font_bold">RecordCompanies</span></td>
</tr>
<tr id="A1.T18.fig2.1.2" class="ltx_tr">
<td id="A1.T18.fig2.1.2.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T18.fig2.1.2.1.1" class="ltx_text ltx_font_bold">Column Name</span></td>
<td id="A1.T18.fig2.1.2.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T18.fig2.1.2.2.1" class="ltx_text ltx_font_bold">Type</span></td>
</tr>
<tr id="A1.T18.fig2.1.3" class="ltx_tr">
<td id="A1.T18.fig2.1.3.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Record_Company</td>
<td id="A1.T18.fig2.1.3.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T18.fig2.1.4" class="ltx_tr">
<td id="A1.T18.fig2.1.4.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Signing_Date</td>
<td id="A1.T18.fig2.1.4.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TIME</td>
</tr>
<tr id="A1.T18.fig2.1.5" class="ltx_tr">
<td id="A1.T18.fig2.1.5.1" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Singer_ID</td>
<td id="A1.T18.fig2.1.5.2" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_figure">Table 18: </span>Schema of the RecordCompanies table</figcaption>
</figure>
</div>
</div>
</figure>

<figure id="A1.T19" class="ltx_table">
<table id="A1.T19.1" class="ltx_tabular ltx_centering ltx_align_middle">
<colgroup>
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
<col style="width: 25%" />
</colgroup>
<tbody>
<tr id="A1.T19.1.1" class="ltx_tr">
<td id="A1.T19.1.1.1" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.1.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.1.1.1.1" class="ltx_p" style="width:71.1pt;"><span id="A1.T19.1.1.1.1.1.1" class="ltx_text ltx_font_bold">Table ID</span></span> </span></td>
<td id="A1.T19.1.1.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.1.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.1.2.1.1" class="ltx_p" style="width:79.7pt;"><span id="A1.T19.1.1.2.1.1.1" class="ltx_text ltx_font_bold">Question</span></span> </span></td>
<td id="A1.T19.1.1.3" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.1.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.1.3.1.1" class="ltx_p" style="width:56.9pt;"><span id="A1.T19.1.1.3.1.1.1" class="ltx_text ltx_font_bold">Answer</span></span> </span></td>
<td id="A1.T19.1.1.4" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.1.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.1.4.1.1" class="ltx_p" style="width:176.4pt;"><span id="A1.T19.1.1.4.1.1.1" class="ltx_text ltx_font_bold">SQL reference</span></span> </span></td>
</tr>
<tr id="A1.T19.1.2" class="ltx_tr">
<td id="A1.T19.1.2.1" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.2.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.2.1.1.1" class="ltx_p" style="width:71.1pt;">GoldenMelody-<br />
Awards</span> </span></td>
<td id="A1.T19.1.2.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.2.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.2.2.1.1" class="ltx_p" style="width:79.7pt;">Golden Melody hosts, excluding the two with the least awards.</span> </span></td>
<td id="A1.T19.1.2.3" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.2.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.2.3.1.1" class="ltx_p" style="width:56.9pt;">"26th Golden Melody", "27th Golden Melody"</span> </span></td>
<td id="A1.T19.1.2.4" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.2.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.2.4.1.1" class="ltx_p" style="width:176.4pt;">select Award_Name from GoldenMelodyAwards where Host not in ( select Host from GoldenMelodyAwards group by Host order by avg ( Awards_Count ) asc limit 2 )</span> </span></td>
</tr>
<tr id="A1.T19.1.3" class="ltx_tr">
<td id="A1.T19.1.3.1" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.3.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.3.1.1.1" class="ltx_p" style="width:71.1pt;">AwardNominees &amp; Singers</span> </span></td>
<td id="A1.T19.1.3.2" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.3.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.3.2.1.1" class="ltx_p" style="width:79.7pt;">Names of singers never nominated for Golden Melody Awards.</span> </span></td>
<td id="A1.T19.1.3.3" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.3.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.3.3.1.1" class="ltx_p" style="width:56.9pt;">"Jay Chou", "Jian Cui"</span> </span></td>
<td id="A1.T19.1.3.4" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.3.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.3.4.1.1" class="ltx_p" style="width:176.4pt;">select Name from Singers where Singer_ID not in ( select Singer_ID from AwardNominees )</span> </span></td>
</tr>
<tr id="A1.T19.1.4" class="ltx_tr">
<td id="A1.T19.1.4.1" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.4.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.4.1.1.1" class="ltx_p" style="width:71.1pt;">RecordCompanies &amp; Singers</span> </span></td>
<td id="A1.T19.1.4.2" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.4.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.4.2.1.1" class="ltx_p" style="width:79.7pt;">Name and gender of singers without a record company.</span> </span></td>
<td id="A1.T19.1.4.3" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.4.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.4.3.1.1" class="ltx_p" style="width:56.9pt;">"Penny Tai:Femal"</span> </span></td>
<td id="A1.T19.1.4.4" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.4.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.4.4.1.1" class="ltx_p" style="width:176.4pt;">select Name, Gender from Singers where Singer_ID not in ( select Singer_ID from RecordCompanies );</span> </span></td>
</tr>
<tr id="A1.T19.1.5" class="ltx_tr">
<td id="A1.T19.1.5.1" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.5.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.5.1.1.1" class="ltx_p" style="width:71.1pt;">GoldenMelody-<br />
Awards</span> </span></td>
<td id="A1.T19.1.5.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.5.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.5.2.1.1" class="ltx_p" style="width:79.7pt;">How many times is the 27th Golden Melody count of the 28th’s?</span> </span></td>
<td id="A1.T19.1.5.3" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.5.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.5.3.1.1" class="ltx_p" style="width:56.9pt;">1</span> </span></td>
<td id="A1.T19.1.5.4" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T19.1.5.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T19.1.5.4.1.1" class="ltx_p" style="width:176.4pt;">select a.Awards_Count / b.Awards_Count from ( select Awards_Count from GoldenMelodyAwards where Award_Name == ’27th Golden Melody’ ) a , ( select Awards_Count from GoldenMelodyAwards where Award_Name == ’28th Golden Melody’ ) b</span> </span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 19: </span>Demostrations of complex nested SQL queries.</figcaption>
</figure>

::: {#A1.p3 .ltx_para}
[Complex nested queries utilizing multiple tools]{#A1.p3.1.1 .ltx_text .ltx_font_bold}: These are advanced queries that involve multiple tools, such as SQL queries, python code generation, user-defined functions, etc. We give the Schema of two tables in the SQL database in Table [[21]{.ltx_text .ltx_ref_tag}](#A1.T21 "Table 21 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}, and [[21]{.ltx_text .ltx_ref_tag}](#A1.T21 "Table 21 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref} and list several examples in Table [[22]{.ltx_text .ltx_ref_tag}](#A1.T22 "Table 22 ‣ Appendix A Detailed Dataset Description ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"){.ltx_ref}. For verifying the planning ability of the LLM-based AI agents, we select this type of query.
:::

<figure id="A1.T21" class="ltx_table">
<div class="ltx_flex_figure ltx_flex_table">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="A1.T21.fig1" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_middle" style="width:216.8pt;">
<table id="A1.T21.fig1.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A1.T21.fig1.1.1" class="ltx_tr">
<td colspan="2" id="A1.T21.fig1.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T21.fig1.1.1.1.1" class="ltx_text ltx_font_bold">Journal</span></td>
</tr>
<tr id="A1.T21.fig1.1.2" class="ltx_tr">
<td id="A1.T21.fig1.1.2.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T21.fig1.1.2.1.1" class="ltx_text ltx_font_bold">Column Name</span></td>
<td id="A1.T21.fig1.1.2.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T21.fig1.1.2.2.1" class="ltx_text ltx_font_bold">Type</span></td>
</tr>
<tr id="A1.T21.fig1.1.3" class="ltx_tr">
<td id="A1.T21.fig1.1.3.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Name</td>
<td id="A1.T21.fig1.1.3.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T21.fig1.1.4" class="ltx_tr">
<td id="A1.T21.fig1.1.4.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">First_Issue_Date</td>
<td id="A1.T21.fig1.1.4.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TIME</td>
</tr>
<tr id="A1.T21.fig1.1.5" class="ltx_tr">
<td id="A1.T21.fig1.1.5.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Journal_ID</td>
<td id="A1.T21.fig1.1.5.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
<tr id="A1.T21.fig1.1.6" class="ltx_tr">
<td id="A1.T21.fig1.1.6.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Category</td>
<td id="A1.T21.fig1.1.6.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T21.fig1.1.7" class="ltx_tr">
<td id="A1.T21.fig1.1.7.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Sponsor_Organization</td>
<td id="A1.T21.fig1.1.7.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T21.fig1.1.8" class="ltx_tr">
<td id="A1.T21.fig1.1.8.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Country</td>
<td id="A1.T21.fig1.1.8.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T21.fig1.1.9" class="ltx_tr">
<td id="A1.T21.fig1.1.9.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">s Language</td>
<td id="A1.T21.fig1.1.9.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">TEXT</td>
</tr>
<tr id="A1.T21.fig1.1.10" class="ltx_tr">
<td id="A1.T21.fig1.1.10.1" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Publication_Count</td>
<td id="A1.T21.fig1.1.10.2" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_figure">Table 20: </span>Schema of the Journal table</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="A1.T21.fig2" class="ltx_figure ltx_figure_panel ltx_minipage ltx_align_center ltx_align_middle" style="width:216.8pt;">
<table id="A1.T21.fig2.1" class="ltx_tabular ltx_centering ltx_align_middle">
<tbody>
<tr id="A1.T21.fig2.1.1" class="ltx_tr">
<td colspan="2" id="A1.T21.fig2.1.1.1" class="ltx_td ltx_align_center ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T21.fig2.1.1.1.1" class="ltx_text ltx_font_bold">CoverPersonality</span></td>
</tr>
<tr id="A1.T21.fig2.1.2" class="ltx_tr">
<td id="A1.T21.fig2.1.2.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T21.fig2.1.2.1.1" class="ltx_text ltx_font_bold">Column Name</span></td>
<td id="A1.T21.fig2.1.2.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T21.fig2.1.2.2.1" class="ltx_text ltx_font_bold">Type</span></td>
</tr>
<tr id="A1.T21.fig2.1.3" class="ltx_tr">
<td id="A1.T21.fig2.1.3.1" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Person_ID</td>
<td id="A1.T21.fig2.1.3.2" class="ltx_td ltx_align_left ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
<tr id="A1.T21.fig2.1.4" class="ltx_tr">
<td id="A1.T21.fig2.1.4.1" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Journal_ID</td>
<td id="A1.T21.fig2.1.4.2" class="ltx_td ltx_align_left" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
<tr id="A1.T21.fig2.1.5" class="ltx_tr">
<td id="A1.T21.fig2.1.5.1" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">Count</td>
<td id="A1.T21.fig2.1.5.2" class="ltx_td ltx_align_left ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt">INTEGER</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_figure">Table 21: </span>Schema of the CoverPersonality table</figcaption>
</figure>
</div>
</div>
</figure>

:::: {#A1.T22 .ltx_table .ltx_transformed_outer style="width:335.9pt;height:689.9pt;vertical-align:-0.0pt;"}
::: {.ltx_transformed_inner style="width:689.9pt;transform:translate(-177.01pt,-176.51pt) rotate(-90deg) ;"}
<figure>
<table id="A1.T22.1" class="ltx_tabular ltx_centering ltx_align_middle" style="width:100%;">
<colgroup>
<col style="width: 16%" />
<col style="width: 16%" />
<col style="width: 16%" />
<col style="width: 16%" />
<col style="width: 16%" />
<col style="width: 16%" />
</colgroup>
<tbody>
<tr id="A1.T22.1.1" class="ltx_tr">
<td id="A1.T22.1.1.1" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.1.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.1.1.1.1" class="ltx_p" style="width:74.0pt;">Table ID</span> </span></td>
<td id="A1.T22.1.1.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.1.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.1.2.1.1" class="ltx_p" style="width:85.4pt;">Question</span> </span></td>
<td id="A1.T22.1.1.3" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.1.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.1.3.1.1" class="ltx_p" style="width:113.8pt;">Answer</span> </span></td>
<td id="A1.T22.1.1.4" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.1.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.1.4.1.1" class="ltx_p" style="width:85.4pt;">Planning Tools</span> </span></td>
<td id="A1.T22.1.1.5" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.1.5.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.1.5.1.1" class="ltx_p" style="width:142.3pt;">SQL reference</span> </span></td>
<td id="A1.T22.1.1.6" class="ltx_td ltx_align_justify ltx_align_top ltx_border_tt" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.1.6.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.1.6.1.1" class="ltx_p" style="width:113.8pt;">Code reference</span> </span></td>
</tr>
<tr id="A1.T22.1.2" class="ltx_tr">
<td id="A1.T22.1.2.1" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.2.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.2.1.1.1" class="ltx_p" style="width:74.0pt;">Journal &amp; CoverPersonality</span> </span></td>
<td id="A1.T22.1.2.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.2.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.2.2.1.1" class="ltx_p" style="width:85.4pt;">Calculate the exponential of 3 and list the names and languages of journals with no cover personality.</span> </span></td>
<td id="A1.T22.1.2.3" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.2.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.2.3.1.1" class="ltx_p" style="width:113.8pt;">[20.08,<br />
"The Economist: Chinese, Reader’s Digest: English."]</span> </span></td>
<td id="A1.T22.1.2.4" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.2.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.2.4.1.1" class="ltx_p" style="width:85.4pt;">["PythonREPL",<br />
"SQL Generator"]</span> </span></td>
<td id="A1.T22.1.2.5" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.2.5.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.2.5.1.1" class="ltx_p" style="width:142.3pt;">select Name, Language from Journal where Journal_ID not in ( select Journal_ID from CoverPersonality )</span> </span></td>
<td id="A1.T22.1.2.6" class="ltx_td ltx_align_justify ltx_align_top ltx_border_t" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.2.6.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.2.6.1.1" class="ltx_p" style="width:113.8pt;">import math;<br />
return math.exp(3)</span> </span></td>
</tr>
<tr id="A1.T22.1.3" class="ltx_tr">
<td id="A1.T22.1.3.1" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.3.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.3.1.1.1" class="ltx_p" style="width:74.0pt;">CoverPersonality &amp; Journal</span> </span></td>
<td id="A1.T22.1.3.2" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.3.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.3.2.1.1" class="ltx_p" style="width:85.4pt;">Compute 4’s factorial, compare with GCD of 212 and list the names and languages of journals with no cover personality.</span> </span></td>
<td id="A1.T22.1.3.3" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.3.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.3.3.1.1" class="ltx_p" style="width:113.8pt;">[4, "The Economist: Chinese, Reader’s Digest: English."]</span> </span></td>
<td id="A1.T22.1.3.4" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.3.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.3.4.1.1" class="ltx_p" style="width:85.4pt;">["PythonREPL",<br />
"SQL Generator"]</span> </span></td>
<td id="A1.T22.1.3.5" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.3.5.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.3.5.1.1" class="ltx_p" style="width:142.3pt;">select Name, Language from Journal where Journal_ID not in ( select Journal_ID from CoverPersonality )</span> </span></td>
<td id="A1.T22.1.3.6" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.3.6.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.3.6.1.1" class="ltx_p" style="width:113.8pt;">import math;<br />
math.gcd(math.factorial(4), 212)</span> </span></td>
</tr>
<tr id="A1.T22.1.4" class="ltx_tr">
<td id="A1.T22.1.4.1" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.4.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.4.1.1.1" class="ltx_p" style="width:74.0pt;">Journal</span> </span></td>
<td id="A1.T22.1.4.2" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.4.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.4.2.1.1" class="ltx_p" style="width:85.4pt;">Calculate the square root of 24, and query for the language whose average number of published issues exceeds the overall average</span> </span></td>
<td id="A1.T22.1.4.3" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.4.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.4.3.1.1" class="ltx_p" style="width:113.8pt;">[4.8989795, "English"]</span> </span></td>
<td id="A1.T22.1.4.4" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.4.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.4.4.1.1" class="ltx_p" style="width:85.4pt;">["PythonREPL",<br />
"SQL Generator"]</span> </span></td>
<td id="A1.T22.1.4.5" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.4.5.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.4.5.1.1" class="ltx_p" style="width:142.3pt;">select Language from Journal group by Language having avg ( Publication_Count ) &gt; ( select avg ( Publication_Count ) from Journal )</span> </span></td>
<td id="A1.T22.1.4.6" class="ltx_td ltx_align_justify ltx_align_top" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.4.6.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.4.6.1.1" class="ltx_p" style="width:113.8pt;">import math;<br />
math.sqrt(24)</span> </span></td>
</tr>
<tr id="A1.T22.1.5" class="ltx_tr">
<td id="A1.T22.1.5.1" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.5.1.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.5.1.1.1" class="ltx_p" style="width:74.0pt;">CoverPersonality</span> </span></td>
<td id="A1.T22.1.5.2" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.5.2.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.5.2.1.1" class="ltx_p" style="width:85.4pt;">Compute the log base 10 of 5, then identify cover figures appearing less than the overall max frequency across journals.</span> </span></td>
<td id="A1.T22.1.5.3" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.5.3.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.5.3.1.1" class="ltx_p" style="width:113.8pt;">[0.69897, "Qing Hai, Xiaoming Huang, Cristiano Ronaldo, Kobe Bryant"]</span> </span></td>
<td id="A1.T22.1.5.4" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.5.4.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.5.4.1.1" class="ltx_p" style="width:85.4pt;">["PythonREPL",<br />
"SQL Generator"]</span> </span></td>
<td id="A1.T22.1.5.5" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.5.5.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.5.5.1.1" class="ltx_p" style="width:142.3pt;">select Person_ID from CoverPersonality where Count &lt; ( select max ( Count ) from CoverPersonality )</span> </span></td>
<td id="A1.T22.1.5.6" class="ltx_td ltx_align_justify ltx_align_top ltx_border_bb" style="padding-top: 1.25pt; padding-bottom: 1.25pt"><span id="A1.T22.1.5.6.1" class="ltx_inline-block ltx_align_top"> <span id="A1.T22.1.5.6.1.1" class="ltx_p" style="width:113.8pt;">import math;<br />
math.log10(5)</span> </span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 22: </span>Demostrations of complex nested queries utilizing multiple tools.</figcaption>
</figure>
:::
::::

::: {.ltx_pagination .ltx_role_newpage}
:::
:::::::::

::: {#A2 .section .ltx_appendix}
## [Appendix B ]{.ltx_tag .ltx_tag_appendix}Prompts Design {#appendix-b-prompts-design .ltx_title .ltx_title_appendix}

<figure id="A2.F8" class="ltx_figure">
<img src="/html/2308.03427/assets/prompt/4_plan_sql_math_tools.png" id="A2.F8.g1" class="ltx_graphics ltx_centering ltx_img_square" width="598" height="598" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 8: </span>The evaluation prompt for tool order planning.</figcaption>
</figure>

<figure id="A2.F9" class="ltx_figure">
<img src="/html/2308.03427/assets/prompt/4_plan_sql_math_tools_query.png" id="A2.F9.g1" class="ltx_graphics ltx_centering ltx_img_square" width="598" height="717" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 9: </span>The evaluation prompt for tool order and subtask description planning.</figcaption>
</figure>

<figure id="A2.F10" class="ltx_figure">
<img src="/html/2308.03427/assets/prompt/4_plan_sql_math_oneshot.png" id="A2.F10.g1" class="ltx_graphics ltx_centering ltx_img_portrait" width="598" height="769" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 10: </span>The evaluation prompt for one-step tool-subtask pair planning.</figcaption>
</figure>

<figure id="A2.F11" class="ltx_figure">
<img src="/html/2308.03427/assets/prompt/5_plan_sql_math_others.png" id="A2.F11.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="598" height="420" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 11: </span>The prompt added to Figure <a href="#A2.F10" class="ltx_ref" title="Figure 10 ‣ Appendix B Prompts Design ‣ TPTU: Large Language Model-based AI Agents for Task Planning and Tool Usage"><span class="ltx_text ltx_ref_tag">10</span></a> for tool-subtask pair planning with other unrelated tools.</figcaption>
</figure>

<figure id="A2.F12" class="ltx_figure">
<img src="/html/2308.03427/assets/prompt/5_plan_seq_planning.png" id="A2.F12.g1" class="ltx_graphics ltx_centering ltx_img_portrait" width="598" height="953" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 12: </span>The prompt for the tool-subtask pair generation with TPTU-SA.</figcaption>
</figure>

<figure id="A2.F13" class="ltx_figure">
<img src="/html/2308.03427/assets/prompt/1_simple_SQL_prompt.png" id="A2.F13.g1" class="ltx_graphics ltx_centering ltx_img_square" width="598" height="660" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 13: </span>The evaluation prompt for simple SQL questions.</figcaption>
</figure>

<figure id="A2.F14" class="ltx_figure">
<img src="/html/2308.03427/assets/prompt/2_nested_SQL_prompt.png" id="A2.F14.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="598" height="469" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 14: </span>The evaluation prompt for complex nested SQL questions.</figcaption>
</figure>

<figure id="A2.F15" class="ltx_figure">
<img src="/html/2308.03427/assets/prompt/2_nested_SQL_prompt_CoT.png" id="A2.F15.g1" class="ltx_graphics ltx_centering ltx_img_square" width="598" height="724" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 15: </span>The evaluation CoT-based prompt for complex nested SQL questions.</figcaption>
</figure>

<figure id="A2.F16" class="ltx_figure">
<img src="/html/2308.03427/assets/prompt/3_math_prompt.png" id="A2.F16.g1" class="ltx_graphics ltx_centering ltx_img_portrait" width="598" height="877" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 16: </span>The evaluation prompt for mathematical questions.</figcaption>
</figure>

<figure id="A2.F17" class="ltx_figure">
<img src="/html/2308.03427/assets/prompt/6_onestep_agent_prompt_en.png" id="A2.F17.g1" class="ltx_graphics ltx_centering ltx_img_portrait" width="598" height="812" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 17: </span>The system prompt for one-step agent.</figcaption>
</figure>

<figure id="A2.F18" class="ltx_figure">
<img src="/html/2308.03427/assets/prompt/6_sequential_agent_prompt.png" id="A2.F18.g1" class="ltx_graphics ltx_centering ltx_img_square" width="598" height="673" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 18: </span>The system prompt for the sequential agent.</figcaption>
</figure>
:::
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

::: ar5iv-footer
[◄](/html/2308.03426){.ar5iv-nav-button .ar5iv-nav-button-prev} [![ar5iv homepage](/assets/ar5iv.png){height="40"}](/){.ar5iv-home-button} [Feeling\
lucky?](/feeling_lucky){.ar5iv-text-button} [](/land_of_honey_and_milk){rel="nofollow" aria-hidden="true" tabindex="-1"} [Conversion\
report](/log/2308.03427){.ar5iv-text-button .ar5iv-severity-error} [Report\
an issue](https://github.com/dginev/ar5iv/issues/new?template=improve-article--arxiv-id-.md&title=Improve+article+2308.03427){.ar5iv-text-button target="_blank"} [View original\
on arXiv](https://arxiv.org/abs/2308.03427){.ar5iv-text-button .arxiv-ui-theme}[►](/html/2308.03428){.ar5iv-nav-button .ar5iv-nav-button-next}
:::

[[]{.color-scheme-icon}](javascript:toggleColorScheme() "Toggle ar5iv color scheme"){.ar5iv-toggle-color-scheme} [Copyright](https://arxiv.org/help/license){.ar5iv-footer-button target="_blank"} [Privacy Policy](https://arxiv.org/help/policies/privacy_policy){.ar5iv-footer-button target="_blank"}

::: ltx_page_logo
Generated on Wed Feb 28 13:54:28 2024 by [[L[a]{.ltx_font_smallcaps style="position:relative; bottom:2.2pt;"}T[e]{.ltx_font_smallcaps style="font-size:120%;position:relative; bottom:-0.2ex;"}]{style="letter-spacing:-0.2em; margin-right:0.1em;"}[XML]{style="font-size:90%; position:relative; bottom:-0.2ex;"}![Mascot Sammy](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAsAAAAOCAYAAAD5YeaVAAAAAXNSR0IArs4c6QAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAALEwAACxMBAJqcGAAAAAd0SU1FB9wKExQZLWTEaOUAAAAddEVYdENvbW1lbnQAQ3JlYXRlZCB3aXRoIFRoZSBHSU1Q72QlbgAAAdpJREFUKM9tkL+L2nAARz9fPZNCKFapUn8kyI0e4iRHSR1Kb8ng0lJw6FYHFwv2LwhOpcWxTjeUunYqOmqd6hEoRDhtDWdA8ApRYsSUCDHNt5ul13vz4w0vWCgUnnEc975arX6ORqN3VqtVZbfbTQC4uEHANM3jSqXymFI6yWazP2KxWAXAL9zCUa1Wy2tXVxheKA9YNoR8Pt+aTqe4FVVVvz05O6MBhqUIBGk8Hn8HAOVy+T+XLJfLS4ZhTiRJgqIoVBRFIoric47jPnmeB1mW/9rr9ZpSSn3Lsmir1fJZlqWlUonKsvwWwD8ymc/nXwVBeLjf7xEKhdBut9Hr9WgmkyGEkJwsy5eHG5vN5g0AKIoCAEgkEkin0wQAfN9/cXPdheu6P33fBwB4ngcAcByHJpPJl+fn54mD3Gg0NrquXxeLRQAAwzAYj8cwTZPwPH9/sVg8PXweDAauqqr2cDjEer1GJBLBZDJBs9mE4zjwfZ85lAGg2+06hmGgXq+j3+/DsixYlgVN03a9Xu8jgCNCyIegIAgx13Vfd7vdu+FweG8YRkjXdWy329+dTgeSJD3ieZ7RNO0VAXAPwDEAO5VKndi2fWrb9jWl9Esul6PZbDY9Go1OZ7PZ9z/lyuD3OozU2wAAAABJRU5ErkJggg==)](http://dlmf.nist.gov/LaTeXML/){.ltx_LaTeXML_logo target="_blank"}
:::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
