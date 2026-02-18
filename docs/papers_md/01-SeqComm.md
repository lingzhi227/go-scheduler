::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_main
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::: ltx_page_content
# Multi-Agent Sequential Decision-Making via Communication {#multi-agent-sequential-decision-making-via-communication .ltx_title .ltx_title_document}

::: ltx_authors
[ [Ziluo Ding     Kefan Su     Weixin Hong     Liwen Zhu     Tiejun Huang     Zongqing Lu\
Peking University ]{.ltx_personname}]{.ltx_creator .ltx_role_author}
:::

::: ltx_abstract
###### Abstract {#abstract .ltx_title .ltx_title_abstract}

Communication helps agents to obtain information about others so that better coordinated behavior can be learned. Some existing work communicates predicted future trajectory with others, hoping to get clues about what others would do for better coordination. However, circular dependencies sometimes can occur when agents are treated synchronously so it is hard to coordinate decision-making. In this paper, we propose a novel communication scheme, [Sequential Communication]{#id1.id1.1 .ltx_text .ltx_font_italic} (SeqComm). SeqComm treats agents asynchronously (the upper-level agents make decisions before the lower-level ones) and has two communication phases. In negotiation phase, agents determine the priority of decision-making by communicating hidden states of observations and comparing the value of intention, which is obtained by modeling the environment dynamics. In launching phase, the upper-level agents take the lead in making decisions and communicate their actions with the lower-level agents. Theoretically, we prove the policies learned by SeqComm are guaranteed to improve monotonically and converge. Empirically, we show that SeqComm outperforms existing methods in various multi-agent cooperative tasks.
:::

:::::::: {#S1 .section .ltx_section}
## [1 ]{.ltx_tag .ltx_tag_section}Introduction {#introduction .ltx_title .ltx_title_section}

::: {#S1.p1 .ltx_para .ltx_noindent}
The partial observability and stochasticity inherent to the nature of multi-agent systems can easily impede the cooperation among agents and lead to catastrophic miscoordination (Ding et al., [2020](#bib.bib6){.ltx_ref}). Communication has been exploited to help agents obtain extra information during both training and execution to mitigate such problems (Foerster et al., [2016](#bib.bib10){.ltx_ref}; Sukhbaatar et al., [2016](#bib.bib45){.ltx_ref}; Peng et al., [2017](#bib.bib31){.ltx_ref}). Specifically, agents can share their information with others via a trainable communication channel.
:::

::: {#S1.p2 .ltx_para .ltx_noindent}
Centralized training with decentralized execution (CTDE) is a popular learning paradigm in cooperative multi-agent reinforcement learning (MARL). Although the centralized value function can be learned to evaluate the joint policy of agents, the decentralized policies of agents are essentially independent. Therefore, a coordination problem arises. That is, agents may make sub-optimal actions by mistakenly assuming others' actions when there exist multiple optimal joint actions (Busoniu et al., [2008](#bib.bib4){.ltx_ref}). Communication allows agents to obtain information about others to avoid miscoordination. However, most existing work only focuses on communicating messages, *e.g.,* the information of agents' current observation or historical trajectory (Jiang & Lu, [2018](#bib.bib18){.ltx_ref}; Singh et al., [2019](#bib.bib42){.ltx_ref}; Das et al., [2019](#bib.bib5){.ltx_ref}; Ding et al., [2020](#bib.bib6){.ltx_ref}). It is impossible for an agent to acquire other's actions before making decisions since the game model is usually synchronous, [i.e.]{#S1.p2.1.2 .ltx_text .ltx_font_italic}, agents make decisions and execute actions simultaneously. Recently, intention or imagination, depicted by a combination of predicted actions and observations of many future steps, has been proposed as part of messages (Kim et al., [2021](#bib.bib21){.ltx_ref}; Pretorius et al., [2021](#bib.bib33){.ltx_ref}). However, circular dependencies can still occur, so it may be hard to coordinate decision-making under synchronous settings.
:::

::: {#S1.p3 .ltx_para .ltx_noindent}
A general approach to solving the coordination problem is to make sure that ties between equally good actions are broken by all agents. One simple mechanism for doing so is to know exactly what others will do and adjust the behavior accordingly under a unique ordering of agents and actions (Busoniu et al., [2008](#bib.bib4){.ltx_ref}). Inspired by this, we reconsider the cooperative game from an asynchronous perspective. In other words, each agent is assigned a priority (*i.e.,* order) of decision-making each step in both training and execution, thus the Stackelberg equilibrium (SE) (Von Stackelberg, [2010](#bib.bib48){.ltx_ref}) is naturally set up as the learning objective. Specifically, the upper-level agents make decisions before the lower-level agents. Therefore, the lower-level agents can acquire the actual actions of the upper-level agents by communication and make their decisions conditioned on what the upper-level agents would do. Under this setting, the SE is likely to be Pareto superior to the average Nash equilibrium (NE) in games that require a high cooperation level (Zhang et al., [2020](#bib.bib55){.ltx_ref}). However, [is it necessary to decide a specific priority of decision-making for each agent?]{#S1.p3.1.2 .ltx_text .ltx_font_italic} Ideally, the optimal joint policy can be decomposed by any orders (Wen et al., [2019](#bib.bib53){.ltx_ref}), *e.g.,* ${\pi^{\ast}\hspace{0pt}{(a_{1},\left. a_{2} \middle| s \right.)}} = {\pi^{\ast}\hspace{0pt}{(\left. a_{1} \middle| s \right.)}\hspace{0pt}\pi^{\ast}\hspace{0pt}{(\left. a_{2} \middle| {s,a_{1}} \right.)}} = {\pi^{\ast}\hspace{0pt}{(\left. a_{2} \middle| s \right.)}\hspace{0pt}\pi^{\ast}\hspace{0pt}{(\left. a_{1} \middle| {s,a_{2}} \right.)}}$. But during the learning process, it is unlikely for agents to use the optimal actions of other agents for gradient calculation, making it still vulnerable to the relative overgeneralization problem (Wei et al., [2018](#bib.bib52){.ltx_ref}). Overall, there is no guarantee that the above equation will hold in the learning process, thus ordering should be carefully concerned.
:::

::: {#S1.p4 .ltx_para .ltx_noindent}
In this paper, we propose a novel model-based multi-round communication scheme for cooperative MARL, [Sequential Communication]{#S1.p4.1.1 .ltx_text .ltx_font_italic} (SeqComm), to enable agents to explicitly coordinate with each other. Specifically, SeqComm has two-phase communication, negotiation phase and launching phase. In the negotiation phase, agents communicate their hidden states of observations with others simultaneously. Then they are able to generate multiple predicted trajectories, called [intention]{#S1.p4.1.2 .ltx_text .ltx_font_italic}, by modeling the environmental dynamics and other agents' actions. In addition, the priority of decision-making is determined by communicating and comparing the corresponding values of agents' intentions. The value of each intention represents the rewards obtained by letting that agent take the upper-level position of the order sequence. The sequence of others follows the same procedure as aforementioned with the upper-level agents fixed. In the launching phase, the upper-level agents take the lead in decision-making and communicate their actual actions with the lower-level agents. Note that the actual actions will be executed simultaneously in the environment without any changes.
:::

::: {#S1.p5 .ltx_para .ltx_noindent}
SeqComm adopts the CTDE paradigm (Lowe et al., [2017](#bib.bib25){.ltx_ref}) and is currently built on MAPPO (Yu et al., [2021](#bib.bib54){.ltx_ref}). Theoretically, we prove the policies learned by SeqComm are guaranteed to improve monotonically and converge. Empirically, we evaluate SeqComm on a set of tasks in multi-agent particle environment (MPE) (Lowe et al., [2017](#bib.bib25){.ltx_ref}) and StarCraft multi-agent challenge (SMAC) (Samvelyan et al., [2019](#bib.bib40){.ltx_ref}). In all these tasks, we demonstrate that SeqComm outperforms prior communication-free and communication-based methods. By ablation studies, we confirm that treating agents asynchronously is a more effective way to promote coordination and SeqComm can provide the proper priority of decision-making for agents to develop better coordination.
:::
::::::::

:::::: {#S2 .section .ltx_section}
## [2 ]{.ltx_tag .ltx_tag_section}Related Work {#related-work .ltx_title .ltx_title_section}

::: {#S2.p1 .ltx_para .ltx_noindent}
[Communication.]{#S2.p1.1.1 .ltx_text .ltx_font_bold} Existing studies (Jiang & Lu, [2018](#bib.bib18){.ltx_ref}; Kim et al., [2019](#bib.bib20){.ltx_ref}; Singh et al., [2019](#bib.bib42){.ltx_ref}; Das et al., [2019](#bib.bib5){.ltx_ref}; Zhang et al., [2019](#bib.bib56){.ltx_ref}; Jiang et al., [2020](#bib.bib19){.ltx_ref}; Ding et al., [2020](#bib.bib6){.ltx_ref}; Konan et al., [2022](#bib.bib22){.ltx_ref}) in this realm mainly focus on how to extract valuable messages. ATOC (Jiang & Lu, [2018](#bib.bib18){.ltx_ref}) and IC3Net (Singh et al., [2019](#bib.bib42){.ltx_ref}) utilize gate mechanisms to decide when to communicate with other agents. Many works (Das et al., [2019](#bib.bib5){.ltx_ref}; Konan et al., [2022](#bib.bib22){.ltx_ref}) employ multi-round communication to fully reason the intentions of others and establish complex collaboration strategies. Social influence (Jaques et al., [2019](#bib.bib17){.ltx_ref}) uses communication to influence the behaviors of others. I2C (Ding et al., [2020](#bib.bib6){.ltx_ref}) only communicates with agents that are relevant and influential which are determined by causal inference. However, all these methods focus on how to exploit valuable information from current or past partial observations effectively and properly. More recently, some studies (Kim et al., [2021](#bib.bib21){.ltx_ref}; Du et al., [2021](#bib.bib7){.ltx_ref}; Pretorius et al., [2021](#bib.bib33){.ltx_ref}) begin to answer the question: can we favor cooperation beyond sharing partial observation? They allow agents to imagine their future states with a world model and communicate those with others. IS (Pretorius et al., [2021](#bib.bib33){.ltx_ref}), as the representation of this line of research, enables each agent to share its intention with other agents in the form of the encoded imagined trajectory and use the attention module to figure out the importance of the received intention. However, two concerns arise. On one hand, circular dependencies can lead to inaccurate predicted future trajectories as long as the multi-agent system treats agents synchronously. On the other hand, MARL struggles in extracting useful information from numerous messages, not to mention more complex and dubious messages, *i.e.,* predicted future trajectories.
:::

::: {#S2.p2 .ltx_para .ltx_noindent}
Unlike these works, we treat the agents from an asynchronously perspective therefore circular dependencies can be naturally resolved. Furthermore, agents only send actions to lower-level agents besides partial observations to make sure the messages are compact as well as informative.
:::

::: {#S2.p3 .ltx_para .ltx_noindent}
[Coordination.]{#S2.p3.1.1 .ltx_text .ltx_font_bold} The agents are essentially independent decision makers in execution and may break ties between equally good actions randomly. Thus, in the absence of additional mechanisms, different agents may break ties in different ways, and the resulting joint actions may be suboptimal. Coordination graphs (Guestrin et al., [2002](#bib.bib14){.ltx_ref}; Böhmer et al., [2020](#bib.bib2){.ltx_ref}; Wang et al., [2021b](#bib.bib51){.ltx_ref}) simplify the coordination when the global Q-function can be additively decomposed into local Q-functions that only depend on the actions of a subset of agents. Typically, a coordination graph expresses a higher-order value decomposition among agents. This improves the representational capacity to distinguish other agents' effects on local utility functions, which addresses the miscoordination problems caused by partial observability. Another general approach to solving the coordination problem is to make sure that ties are broken by all agents in the same way, requiring that random action choices are somehow coordinated or negotiated. Social conventions (Boutilier, [1996](#bib.bib3){.ltx_ref}) or role assignments (Prasad et al., [1998](#bib.bib32){.ltx_ref}) encode prior preferences towards certain joint actions and help break ties during action selection. Communication (Fischer et al., [2004](#bib.bib9){.ltx_ref}; Vlassis, [2007](#bib.bib47){.ltx_ref}) can be used to negotiate action choices, either alone or in combination with the aforementioned techniques. Our method follows this line of research by utilizing the ordering of agents and actions to break the ties, other than the enhanced representational capacity of the local value function.
:::
::::::

:::::::::::::::::: {#S3 .section .ltx_section}
## [3 ]{.ltx_tag .ltx_tag_section}Problem Formulation {#problem-formulation .ltx_title .ltx_title_section}

::: {#S3.p1 .ltx_para .ltx_noindent}
[Cost-Free Communication.]{#S3.p1.1.1 .ltx_text .ltx_font_bold} The decentralized partially observable Markov decision process (Dec-POMDP) can be extended to explicitly incorporate communicated messages and observations. The resulting model is called Dec-POMDP-Com (Goldman & Zilberstein, [2003](#bib.bib11){.ltx_ref}; [2004](#bib.bib12){.ltx_ref})..
:::

::: {#S3.p2 .ltx_para .ltx_noindent}
Pynadath & Tambe ([2002](#bib.bib34){.ltx_ref}) showed that under cost-free communication, a joint communication policy that shares local observations at each stage is optimal. Many studies have also investigated sharing local observations in models that are similar to Dec-POMDP-Com (Pynadath & Tambe, [2002](#bib.bib34){.ltx_ref}; Ooi & Wornell, [1996](#bib.bib29){.ltx_ref}; Nair et al., [2004](#bib.bib27){.ltx_ref}; Roth et al., [2005a](#bib.bib38){.ltx_ref}; [b](#bib.bib39){.ltx_ref}; Spaan et al., [2006](#bib.bib44){.ltx_ref}; Oliehoek et al., [2007](#bib.bib28){.ltx_ref}; Becker et al., [2004](#bib.bib1){.ltx_ref}). These works focus on issues other than communication cost. Although real world rarely exhibits such ideal conditions. They intend to model some domains as having approximately free communication to a sufficient degree and believe analyzing such cases gives us some insight to the benefit of communication.
:::

::: {#S3.p3 .ltx_para .ltx_noindent}
[Multi-Agent Sequential Decision-Making.]{#S3.p3.1.1 .ltx_text .ltx_font_bold} We consider fully cooperative multi-agent tasks that are modeled as Dec-POMDP-Com, where $n$ agents interact with the environment according to the following procedure, which we refer to as [multi-agent sequential decision-making]{#S3.p3.1.2 .ltx_text .ltx_font_italic}.
:::

::: {#S3.p4 .ltx_para .ltx_noindent}
At each timestep $t$, assume the priority (*i.e.,* order) of decision-making for all agents is given and each priority level has only one agent ([i.e.]{#S3.p4.27.2 .ltx_text .ltx_font_italic}, agents make decisions one by one). Note that the smaller the level index, the higher priority of decision-making is. The agent at each level $k$ gets its own observation $o_{t}^{k}$ drawn from the state $s_{t}$, and receives messages ${\mathbf{m}}_{t}^{- k}$ from all other agents, where ${\mathbf{m}}_{t}^{- k} \triangleq {\{{\{ o_{t}^{1},a_{t}^{1}\}},\ldots,{\{ o_{t}^{k - 1},a_{t}^{k - 1}\}},o_{t}^{k + 1},\ldots,o_{t}^{n}\}}$. Equivalently, ${\mathbf{m}}_{t}^{- k}$ can be written as $\{{\mathbf{o}}_{\mathbf{t}}^{- k},{\mathbf{a}}_{t}^{1:{k - 1}}\}$, where ${}{\mathbf{o}}_{\mathbf{t}}^{- k}$ denotes the joint observations of all agents except $k$, and ${\mathbf{a}}_{t}^{1:{k - 1}}$ denotes the joint actions of agents $1$ to $k - 1$. For the agent at the first level ([i.e.]{#S3.p4.27.3 .ltx_text .ltx_font_italic}, $k = 1$), ${\mathbf{a}}_{t}^{1:{k - 1}} = \varnothing$. Then, the agent determines its action $a_{t}^{k}$ sampled from its policy $\pi_{k}{( \cdot |o_{t}^{k},{\mathbf{m}}_{t}^{- k})}$ or equivalently $\pi_{k}{( \cdot |{\mathbf{o}}_{t},{\mathbf{a}}_{t}^{1:{k - 1}})}$ and sends it to the lower-level agents. After all agents have determined their actions, they perform the joint actions ${\mathbf{a}}_{t}$, which can be seen as sampled from the joint policy ${\mathbf{π}}{( \cdot |s_{t})}$ [factorized]{#S3.p4.27.4 .ltx_text .ltx_font_italic} as $\prod_{k = 1}^{n}\pi_{k}{( \cdot |{\mathbf{o}}_{t},{\mathbf{a}}_{t}^{1:{k - 1}})}$, in the environment and get a shared reward $r\hspace{0pt}{(s_{t},{\mathbf{a}}_{t})}$ and the state transitions to next state $s'$ according to the transition probability $p\hspace{0pt}{(\left. s' \middle| {s_{t},{\mathbf{a}}_{t}} \right.)}$. All agents aim to maximize the expected return $\sum_{t = 0}^{\infty}{\gamma^{t}\hspace{0pt}r_{t}}$, where $\gamma$ is the discount factor. The state-value function and action-value function of the level-$k$ agent are defined as follows:

  -- --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $${V_{\pi_{k}}\hspace{0pt}{(s,{\mathbf{a}}^{1:{k - 1}})}} \triangleq {\underset{\begin{matrix}                                                                                                            
     s_{1:\infty} \\                                                                                                                                                                                           
     {{\mathbf{a}}_{0}^{k:n} \sim {\mathbf{π}}_{k:n}} \\                                                                                                                                                       
     {{\mathbf{a}}_{1:\infty} \sim {\mathbf{π}}}                                                                                                                                                               
     \end{matrix}}{\mathbb{E}}{\lbrack{{\left. {\sum\limits_{t = 0}^{\infty}{\gamma^{t}\hspace{0pt}r_{t}}} \middle| s_{0} \right. = s},{{\mathbf{a}}_{0}^{1:{k - 1}} = {\mathbf{a}}^{1:{k - 1}}}}\rbrack}}$$   
     $${{Q_{\pi_{k}}\hspace{0pt}{(s,{\mathbf{a}}^{1:k})}} \triangleq {\underset{\begin{matrix}                                                                                                                 
     s_{1:\infty} \\                                                                                                                                                                                           
     {{\mathbf{a}}_{0}^{{k + 1}:n} \sim {\mathbf{π}}_{{k + 1}:n}} \\                                                                                                                                           
     {{\mathbf{a}}_{1:\infty} \sim {\mathbf{π}}}                                                                                                                                                               
     \end{matrix}}{\mathbb{E}}{\lbrack{{\left. {\sum\limits_{t = 0}^{\infty}{\gamma^{t}\hspace{0pt}r_{t}}} \middle| s_{0} \right. = s},{{\mathbf{a}}_{0}^{1:k} = {\mathbf{a}}^{1:k}}}\rbrack}}}.$$             
  -- --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
:::

::: {#S3.p5 .ltx_para .ltx_noindent}
For the setting of multi-agent sequential decision-making discussed above, we have the following proposition.
:::

:::::: {#Thmproposition1 .ltx_theorem .ltx_theorem_proposition}
###### [[Proposition 1]{#Thmproposition1.2.1.1 .ltx_text .ltx_font_bold}]{.ltx_tag .ltx_tag_theorem}[.]{#Thmproposition1.3.2 .ltx_text .ltx_font_bold} {#proposition-1. .ltx_title .ltx_runin .ltx_title_theorem}

::: {#Thmproposition1.p1 .ltx_para}
[If all the agents update its policy with individual TRPO (Schulman et al., [2015](#bib.bib41){.ltx_ref}) sequentially in multi-agent sequential decision-making, then the joint policy of all agents is guaranteed to improve monotonically and converge.]{#Thmproposition1.p1.1.1 .ltx_text .ltx_font_italic}
:::

:::: {#Thmproposition1.1 .ltx_proof}
###### Proof. {#proof. .ltx_title .ltx_runin .ltx_font_italic .ltx_title_proof}

::: {#Thmproposition1.1.p1 .ltx_para}
[The proof is given in Appendix [[A]{.ltx_text .ltx_ref_tag}](#A1 "Appendix A Proofs of Proposition 1 and Proposition 2 ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}. ∎]{#Thmproposition1.1.p1.1.1 .ltx_text}
:::
::::
::::::

::: {#S3.p6 .ltx_para .ltx_noindent}
Proposition [[1]{.ltx_text .ltx_ref_tag}](#Thmproposition1 "Proposition 1. ‣ 3 Problem Formulation ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} indicates that SeqComm has the performance guarantee regardless of the priority of decision-making in multi-agent sequential decision-making. However, the priority of decision-making indeed affects the optimality of the converged joint policy, and we have the following claim.
:::

:::: {#Thmclaim1 .ltx_theorem .ltx_theorem_claim}
###### [[Claim 1]{#Thmclaim1.1.1.1 .ltx_text .ltx_font_bold}]{.ltx_tag .ltx_tag_theorem}[.]{#Thmclaim1.2.2 .ltx_text .ltx_font_bold} {#claim-1. .ltx_title .ltx_runin .ltx_title_theorem}

::: {#Thmclaim1.p1 .ltx_para}
[The different priorities of decision-making affect the optimality of the convergence of the learning algorithm due to the relative overgeneralization problem.]{#Thmclaim1.p1.1.1 .ltx_text .ltx_font_italic}
:::
::::

<figure id="S3.F1" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S3.T1.st1" class="ltx_table ltx_figure_panel ltx_align_center">
<table id="S3.T1.st1.17" class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<thead class="ltx_thead">
<tr id="S3.T1.st1.1.1" class="ltx_tr">
<th id="S3.T1.st1.1.1.2" class="ltx_td ltx_th ltx_th_column ltx_th_row"></th>
<th id="S3.T1.st1.1.1.3" class="ltx_td ltx_th ltx_th_column ltx_th_row"></th>
<th colspan="3" id="S3.T1.st1.1.1.1" class="ltx_td ltx_align_center ltx_th ltx_th_column"><span id="S3.T1.st1.1.1.1.1" class="ltx_text" style="color:#00009B;">Agent <span class="math inline"><em>B</em></span></span></th>
</tr>
<tr id="S3.T1.st1.4.4" class="ltx_tr">
<th id="S3.T1.st1.4.4.4" class="ltx_td ltx_th ltx_th_column ltx_th_row"></th>
<th id="S3.T1.st1.4.4.5" class="ltx_td ltx_th ltx_th_column ltx_th_row"></th>
<th id="S3.T1.st1.2.2.1" class="ltx_td ltx_align_center ltx_th ltx_th_column"><span class="math inline"><em>b</em><sub>1</sub></span></th>
<th id="S3.T1.st1.3.3.2" class="ltx_td ltx_align_center ltx_th ltx_th_column"><span id="S3.T1.st1.3.3.2.1" class="ltx_text" style="color:#00009B;"><span class="math inline"><em>b</em><sub>2</sub></span></span></th>
<th id="S3.T1.st1.4.4.3" class="ltx_td ltx_align_center ltx_th ltx_th_column"><span class="math inline"><em>b</em><sub>3</sub></span></th>
</tr>
</thead>
<tbody class="ltx_tbody">
<tr id="S3.T1.st1.8.8" class="ltx_tr">
<th id="S3.T1.st1.8.8.5" class="ltx_td ltx_th ltx_th_row"></th>
<th id="S3.T1.st1.5.5.1" class="ltx_td ltx_align_left ltx_th ltx_th_row"><span class="math inline"><em>a</em><sub>1</sub></span></th>
<td id="S3.T1.st1.6.6.2" class="ltx_td ltx_align_center ltx_border_l ltx_border_r ltx_border_t"><span class="math inline">12</span></td>
<td id="S3.T1.st1.7.7.3" class="ltx_td ltx_align_center ltx_border_r ltx_border_t"><span class="math inline">6</span></td>
<td id="S3.T1.st1.8.8.4" class="ltx_td ltx_align_center ltx_border_r ltx_border_t"><span class="math inline">6</span></td>
</tr>
<tr id="S3.T1.st1.12.12" class="ltx_tr">
<th id="S3.T1.st1.12.12.5" class="ltx_td ltx_th ltx_th_row"></th>
<th id="S3.T1.st1.9.9.1" class="ltx_td ltx_align_left ltx_th ltx_th_row"><span class="math inline"><em>a</em><sub>2</sub></span></th>
<td id="S3.T1.st1.10.10.2" class="ltx_td ltx_align_center ltx_border_l ltx_border_r ltx_border_t"><span class="math inline">−6</span></td>
<td id="S3.T1.st1.11.11.3" class="ltx_td ltx_align_center ltx_border_r ltx_border_t"><span class="math inline">8</span></td>
<td id="S3.T1.st1.12.12.4" class="ltx_td ltx_align_center ltx_border_r ltx_border_t"><span class="math inline">0</span></td>
</tr>
<tr id="S3.T1.st1.17.17" class="ltx_tr">
<th id="S3.T1.st1.13.13.1" class="ltx_td ltx_align_center ltx_th ltx_th_row"><span id="S3.T1.st1.13.13.1.1" class="ltx_text"> <span id="S3.T1.st1.13.13.1.1.1" class="ltx_inline-block ltx_transformed_outer" style="width:8.8pt;height:37.2pt;vertical-align:-0.0pt;"><span class="ltx_transformed_inner" style="width:37.2pt;transform:translate(-14.22pt,-13.25pt) rotate(-90deg) ;"> <span id="S3.T1.st1.13.13.1.1.1.1" class="ltx_p"><span id="S3.T1.st1.13.13.1.1.1.1.1" class="ltx_text" style="color:#9A0000;">Agent <span class="math inline"><em>A</em></span></span></span> </span></span></span></th>
<th id="S3.T1.st1.14.14.2" class="ltx_td ltx_align_left ltx_th ltx_th_row"><span class="math inline"><em>a</em><sub>3</sub></span></th>
<td id="S3.T1.st1.15.15.3" class="ltx_td ltx_align_center ltx_border_b ltx_border_l ltx_border_r ltx_border_t"><span class="math inline">−6</span></td>
<td id="S3.T1.st1.16.16.4" class="ltx_td ltx_align_center ltx_border_b ltx_border_r ltx_border_t"><span class="math inline">0</span></td>
<td id="S3.T1.st1.17.17.5" class="ltx_td ltx_align_center ltx_border_b ltx_border_r ltx_border_t"><span class="math inline">8</span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">(a) </span>payoff matrix of the game</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S3.F1.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x1.png" id="S3.F1.sf1.g1" class="ltx_graphics ltx_centering ltx_img_square" width="424" height="354" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(a) </span>evaluations of different methods</figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 1: </span><a href="#S3.T1.st1" class="ltx_ref" title="In Figure 1 ‣ 3 Problem Formulation ‣ Multi-Agent Sequential Decision-Making via Communication"><span class="ltx_text ltx_ref_tag">1(a)</span></a> Payoff matrix for a one-step game. There are multiple local optima. <a href="#S3.F1.sf1" class="ltx_ref" title="In Figure 1 ‣ 3 Problem Formulation ‣ Multi-Agent Sequential Decision-Making via Communication"><span class="ltx_text ltx_ref_tag">1(a)</span></a> Evaluations of different methods for the game in terms of the mean reward and standard deviation of ten runs. <span class="math inline"><em>A</em> → <em>B</em></span>, <span class="math inline"><em>B</em> → <em>A</em></span>, <em>Simultaneous</em>, and <em>Learned</em> represent that agent <span class="math inline"><em>A</em></span> makes decisions first, agent <span class="math inline"><em>B</em></span> makes decisions first, two agents make decisions simultaneously, and there is another learned policy determining the priority of decision making, respectively. MAPPO (Yu et al., <a href="#bib.bib54" class="ltx_ref">2021</a>) is used as the backbone.</figcaption>
</figure>

::: {#S3.p7 .ltx_para .ltx_noindent}
We use a one-step matrix game as an example, as illustrated in Figure [[1(a)]{.ltx_text .ltx_ref_tag}](#S3.T1.st1 "In Figure 1 ‣ 3 Problem Formulation ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}, to demonstrate the influence of the priority of decision-making on the learning process. Due to relative overgeneralization, agent $B$ tends to choose $b_{2}$ or $b_{3}$ even under the CTDE paradigm (Wei et al., [2018](#bib.bib52){.ltx_ref}). Specifically, $b_{2}$ or $b_{3}$ in the suboptimal equilibrium is a better choice than $b_{1}$ in the optimal equilibrium when matched with arbitrary actions from agent $A$. Therefore, as shown in Figure [[1(a)]{.ltx_text .ltx_ref_tag}](#S3.F1.sf1 "In Figure 1 ‣ 3 Problem Formulation ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}, $B\rightarrow A$ ([i.e.]{#S3.p7.18.1 .ltx_text .ltx_font_italic}, agent $B$ makes decisions before $A$, and $A$'s policy conditions on the action of $B$) and [Simultaneous]{#S3.p7.18.2 .ltx_text .ltx_font_italic} ([i.e.]{#S3.p7.18.3 .ltx_text .ltx_font_italic}, two agents make decisions simultaneously and independently) are easily trapped into local optima. However, things can be different if agent $A$ goes first, as $A\rightarrow B$ achieves the optimum. As long as agent $A$ does not suffer from relative overgeneralization, it can help agent $B$ get rid of local optima by narrowing down the search space of $B$. Besides, a policy that determines the priority of decision-making can be learned under the guidance of the state-value function, denoted as [Learned]{#S3.p7.18.4 .ltx_text .ltx_font_italic}. It obtains better performance than $B\rightarrow A$ and [Simultaneous]{#S3.p7.18.5 .ltx_text .ltx_font_italic}, which indicates that dynamically determining the order during policy learning can be beneficial as we do not know the optimal priority in advance.
:::

:::: {#Thmremark1 .ltx_theorem .ltx_theorem_remark}
###### [[Remark 1]{#Thmremark1.1.1.1 .ltx_text .ltx_font_bold}]{.ltx_tag .ltx_tag_theorem}[.]{#Thmremark1.2.2 .ltx_text .ltx_font_bold} {#remark-1. .ltx_title .ltx_runin .ltx_title_theorem}

::: {#Thmremark1.p1 .ltx_para}
The priority (*i.e.,* order) of decision-making affects the optimality of the converged joint policy in multi-agent sequential decision-making, thus it is critical to determine the order. However, learning the order directly requires an additional centralized policy in execution, which does not fit the CTDE paradigm and is not generalizable in the scenario where the number of agents varies. Moreover, its learning complexity exponentially increases with the number of agents, making it infeasible in many cases.
:::
::::
::::::::::::::::::

:::::::::::::::::::::::::::: {#S4 .section .ltx_section}
## [4 ]{.ltx_tag .ltx_tag_section}Sequential Communication {#sequential-communication .ltx_title .ltx_title_section}

::: {#S4.p1 .ltx_para .ltx_noindent}
In this paper, we cast our eyes in another direction and resort to the world model. Ideally, we can randomly sample candidate order sequences, evaluate them under the world model (see Section [[4.1]{.ltx_text .ltx_ref_tag}](#S4.SS1 "4.1 Negotiation Phase ‣ 4 Sequential Communication ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}), and choose the order sequence that is deemed the most promising under the true dynamic. SeqComm is designed based on this principle to determine the priority of decision-making via communication.
:::

::: {#S4.p2 .ltx_para .ltx_noindent}
SeqComm adopts a multi-round communication mechanism, [i.e.]{#S4.p2.1.1 .ltx_text .ltx_font_italic}, agents are allowed to communicate with others in multiple rounds. Importantly, communication is separated into phases serving different purposes. One is the [negotiation]{#S4.p2.1.2 .ltx_text .ltx_font_italic} phase for agents to determine the priority of decision-making. Another is the [launching]{#S4.p2.1.3 .ltx_text .ltx_font_italic} phase for agents to act conditioning on actual actions upper-level agents will take to implement [explicit coordination via communication]{#S4.p2.1.4 .ltx_text .ltx_font_italic}. The overview of SeqComm is illustrated in Figure [[2]{.ltx_text .ltx_ref_tag}](#S4.F2 "Figure 2 ‣ 4.1 Negotiation Phase ‣ 4 Sequential Communication ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}. Each SeqComm agent consists of a policy, a critic, and a world model, as illustrated in Figure [[3]{.ltx_text .ltx_ref_tag}](#S4.F3 "Figure 3 ‣ 4.1 Negotiation Phase ‣ 4 Sequential Communication ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}, and the parameters of all networks are shared across agents (Gupta et al., [2017](#bib.bib15){.ltx_ref}).
:::

::::::::: {#S4.SS1 .section .ltx_subsection}
### [4.1 ]{.ltx_tag .ltx_tag_subsection}Negotiation Phase {#negotiation-phase .ltx_title .ltx_title_subsection}

::: {#S4.SS1.p1 .ltx_para .ltx_noindent}
In the negotiation phase, the observation encoder first takes $o_{t}$ as input and outputs a hidden state $h_{t}$, which is used to communicate with others. Agents then determine the priority of decision-making by [intention]{#S4.SS1.p1.2.1 .ltx_text .ltx_font_italic} which is established and evaluated based on the world model.
:::

::: {#S4.SS1.p2 .ltx_para .ltx_noindent}
[World Model.]{#S4.SS1.p2.3.1 .ltx_text .ltx_font_bold} The world model is needed to predict and evaluate future trajectories. SeqComm, unlike previous works (Kim et al., [2021](#bib.bib21){.ltx_ref}; Du et al., [2021](#bib.bib7){.ltx_ref}; Pretorius et al., [2021](#bib.bib33){.ltx_ref}), can utilize received hidden states of other agents in the first round of communication to model more precise environment dynamics for the explicit coordination in the next round of communication. Once an agent can access other agents' hidden states, it shall have adequate information to estimate their actions since all agents are homogeneous and parameter-sharing. Therefore, the world model $\mathcal{M}\hspace{0pt}{( \cdot )}$ takes as input the joint hidden states ${\mathbf{h}}_{t} = {\{ h_{t}^{1},\ldots,h_{t}^{n}\}}$ and actions ${\mathbf{a}}_{t}$, and predicts the next joint observations and reward,

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------ --
     $${{{\hat{\mathbf{o}}}_{t + 1},{\hat{r}}_{t + 1}} = {\mathcal{M}_{i}\hspace{0pt}{({{AM}_{w}\hspace{0pt}{({\mathbf{h}}_{t},{\mathbf{a}}_{t})}})}}},$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------ --

where ${AM}_{w}$ is the attention module. The reason that we adopt the attention module is to entitle the world model to be generalizable in the scenarios where additional agents are introduced or existing agents are removed.
:::

<figure id="S4.F2" class="ltx_figure">
<img src="/html/2209.12713/assets/x2.png" id="S4.F2.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="422" height="128" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 2: </span>Overview of SeqComm. SeqComm has two communication phases, the negotiation phase (<span id="S4.F2.3.1" class="ltx_text ltx_font_italic">left</span>) and the launching phase (<span id="S4.F2.4.2" class="ltx_text ltx_font_italic">right</span>). In the negotiation phase, agents communicate hidden states of observations with others and obtain their own intention. The priority of decision-making is determined by sharing and comparing the value of all the intentions. In the launching phase, the agents who hold the upper-level positions will make decisions prior to the lower-level agents. Besides, their actions will be shared with anyone that has not yet made decisions.</figcaption>
</figure>

<figure id="S4.F3" class="ltx_figure">
<img src="/html/2209.12713/assets/x3.png" id="S4.F3.g1" class="ltx_graphics ltx_centering ltx_img_square" width="194" height="213" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 3: </span>Architecture of SeqComm. The critic and policy of each agent take input as its own observation and received messages. The world model takes as input the joint hidden states and predicted joint actions.</figcaption>
</figure>

::: {#S4.SS1.p3 .ltx_para .ltx_noindent}
[Priority of Decision-Making.]{#S4.SS1.p3.1.1 .ltx_text .ltx_font_bold} The intention is the key element to determine the priority of decision-making. The notion of intention is described as an agent's future behavior in previous works (Rabinowitz et al., [2018](#bib.bib35){.ltx_ref}; Raileanu et al., [2018](#bib.bib36){.ltx_ref}; Kim et al., [2021](#bib.bib21){.ltx_ref}). However, we define the [intention]{#S4.SS1.p3.1.2 .ltx_text .ltx_font_italic} as an agent's future behavior [without considering others]{#S4.SS1.p3.1.3 .ltx_text .ltx_font_italic}.
:::

::: {#S4.SS1.p4 .ltx_para .ltx_noindent}
As mentioned before, an agent's intention considering others can lead to circular dependencies and cause miscoordination. By our definition, the intention of an agent should be depicted as all future trajectories considering that agent as the first-mover and ignoring the others. However, there are many possible future trajectories as the priority of the rest agents is *unfixed*. In practice, we use the Monte Carlo method to evaluate intention.
:::

::: {#S4.SS1.p5 .ltx_para .ltx_noindent}
Taking agent $i$ at timestep $t$ to illustrate, it firstly considers itself as the first-mover and produces its action only based on the joint hidden states, ${\hat{a}}_{t}^{i} \sim \pi_{i}{( \cdot |{AM}_{a}{({\mathbf{h}}_{t})})}$, where we again use an attention module ${AM}_{a}$ to handle the input. For the order sequence of lower-level agents, we randomly sample a set of order sequences from unfixed agents. Assume agent $j$ is the second-mover, agent $i$ models $j$'s action by considering the upper-level action following its own policy ${\hat{a}}_{t}^{j} \sim \pi_{i}{( \cdot |{AM}_{a}{({\mathbf{h}}_{t},{\hat{a}}_{t}^{i})})}$. The same procedure is applied to predict the actions of all other agents following the sampled order sequence. Based on the joint hidden states and predicted actions, the next joint observations ${\hat{\mathbf{o}}}_{t + 1}$ and corresponding reward ${\hat{r}}_{t + 1}$ can be predicted by the world model. The length of the predicted future trajectory is $H$ and it can then be written as $\tau^{t} = {\{{\hat{\mathbf{o}}}_{t + 1},{\hat{\mathbf{a}}}_{t + 1},\ldots,{\hat{\mathbf{o}}}_{t + H},{\hat{\mathbf{a}}}_{t + H}\}}$ by repeating the procedure aforementioned and the value of one trajectory is defined as the return of that trajectory $v_{\tau^{t}} = {\sum_{t' = {t + 1}}^{t + H}{{\gamma^{t' - t - 1}\hspace{0pt}{\hat{r}}_{t'}}/H}}$. In addition, the intention value is defined as the average value of $F$ future trajectories with different sampled order sequences. The choice of $F$ is a tradeoff between the computation overhead and the accuracy of the estimation.
:::

::: {#S4.SS1.p6 .ltx_para .ltx_noindent}
After all the agents have computed their own intention and the corresponding value, they again communicate their intention values to others. Then agents would compare and choose the agent with the highest intention value to be the first-mover. The priority of lower-level decision-making follows the same procedure with the upper-level agents fixed. Note that some agents are required to communicate intention values with others multiple times until the priority of decision-making is finally determined.
:::
:::::::::

::::: {#S4.SS2 .section .ltx_subsection}
### [4.2 ]{.ltx_tag .ltx_tag_subsection}Launching Phase {#launching-phase .ltx_title .ltx_title_subsection}

::: {#S4.SS2.p1 .ltx_para .ltx_noindent}
As for the launching phase, agents communicate for obtaining additional information to make decisions. Apart from the received hidden states from the last phase, we allow agents to get what [actual]{#S4.SS2.p1.3.1 .ltx_text .ltx_font_italic} actions the upper-level agents will take in execution, while other studies can only infer others' actions by opponent modeling (Rabinowitz et al., [2018](#bib.bib35){.ltx_ref}; Raileanu et al., [2018](#bib.bib36){.ltx_ref}) or communicating intentions (Kim et al., [2021](#bib.bib21){.ltx_ref}). Therefore, miscoordination can be naturally avoided and a better cooperation strategy is possible since lower-level agents can adjust their behaviors accordingly. A lower-level agent $i$ make a decision following the policy $\pi_{i}{( \cdot |{AM}_{a}{({\mathbf{h}}_{t},{\mathbf{a}}_{t}^{u\hspace{0pt}p\hspace{0pt}p\hspace{0pt}e\hspace{0pt}r})})}$, where ${\mathbf{a}}_{t}^{u\hspace{0pt}p\hspace{0pt}p\hspace{0pt}e\hspace{0pt}r}$ means received actual actions from all upper-level agents. As long as the agent has decided its action, it will send its action to all other lower-level agents by the communication channel. [Note that the actions are executed simultaneously and distributedly in execution, though agents make decisions sequentially. ]{#S4.SS2.p1.3.2 .ltx_text .ltx_font_italic}
:::

::: {#S4.SS2.p2 .ltx_para .ltx_noindent}
[Communication Overhead.]{#S4.SS2.p2.5.1 .ltx_text .ltx_font_bold} Two communication phases alternate until all agents determine their levels and get upper-level actions. Note that many previous works also adopt the multi-round communication scheme (Das et al., [2019](#bib.bib5){.ltx_ref}; Singh et al., [2019](#bib.bib42){.ltx_ref}). As for implementation in practice, compared with communicating high-dimensional hidden states/observations by multiple rounds (Das et al., [2019](#bib.bib5){.ltx_ref}; Singh et al., [2019](#bib.bib42){.ltx_ref}), or transferring multi-step trajectory (Kim et al., [2021](#bib.bib21){.ltx_ref}), SeqComm needs more rounds, but it only transmits hidden states for one time. For the rest $n - 1$ round communication with total ${({n - 1})}/2$ broadcasts per agent, only a single intention value and an action will be exchanged. Considering there are $n!$ permutations of different order choices for $n$ agents, our method has greatly reduced computation overhead since each agent needs to calculate up to $n$ times to search for a satisfying order. Although SeqComm is more suitable for latency-tolerate MARL tasks, *e.g.,* power dispatch (minutes) (Wang et al., [2021a](#bib.bib49){.ltx_ref}), inventory management (hours) (Feng et al., [2021](#bib.bib8){.ltx_ref}), maritime transportation (days) (Li et al., [2019](#bib.bib24){.ltx_ref}), it is possible for SeqComm to have a wider range of applications given the rapid development of the communication technology, *e.g.,* 5G.
:::
:::::

::::::::::::::: {#S4.SS3 .section .ltx_subsection}
### [4.3 ]{.ltx_tag .ltx_tag_subsection}Theoretical Analysis {#theoretical-analysis .ltx_title .ltx_title_subsection}

::: {#S4.SS3.p1 .ltx_para .ltx_noindent}
As the priority of decision-making is determined by intention values, SeqComm is likely to choose different orders [at different timesteps]{#S4.SS3.p1.1.1 .ltx_text .ltx_font_italic} during training. However, we have the following proposition that theoretically guarantees the performance of the learned joint policy under SeqComm.
:::

:::::: {#Thmproposition2 .ltx_theorem .ltx_theorem_proposition}
###### [[Proposition 2]{#Thmproposition2.2.1.1 .ltx_text .ltx_font_bold}]{.ltx_tag .ltx_tag_theorem}[.]{#Thmproposition2.3.2 .ltx_text .ltx_font_bold} {#proposition-2. .ltx_title .ltx_runin .ltx_title_theorem}

::: {#Thmproposition2.p1 .ltx_para}
[The monotonic improvement and convergence of the joint policy in SeqComm are independent of the priority of decision-making of agents at each timestep.]{#Thmproposition2.p1.1.1 .ltx_text .ltx_font_italic}
:::

:::: {#Thmproposition2.1 .ltx_proof}
###### Proof. {#proof.-1 .ltx_title .ltx_runin .ltx_font_italic .ltx_title_proof}

::: {#Thmproposition2.1.p1 .ltx_para}
[The proof is given in Appendix [[A]{.ltx_text .ltx_ref_tag}](#A1 "Appendix A Proofs of Proposition 1 and Proposition 2 ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}. ∎]{#Thmproposition2.1.p1.1.1 .ltx_text}
:::
::::
::::::

::: {#S4.SS3.p2 .ltx_para .ltx_noindent}
The priority of decision-making is chosen under the world model, thus the compounding errors in the world model can result in discrepancies between the predicted returns of the same order under the world model and the true dynamics. We then analyze the monotonic improvement for the joint policy under the world model based on Janner et al. ([2019](#bib.bib16){.ltx_ref}).
:::

:::::: {#Thmtheorem1 .ltx_theorem .ltx_theorem_theorem}
###### [[Theorem 1]{#Thmtheorem1.2.1.1 .ltx_text .ltx_font_bold}]{.ltx_tag .ltx_tag_theorem}[.]{#Thmtheorem1.3.2 .ltx_text .ltx_font_bold} {#theorem-1. .ltx_title .ltx_runin .ltx_title_theorem}

::: {#Thmtheorem1.p1 .ltx_para}
[Let the expected total variation between two transition distributions be bounded at each timestep as $\max_{t}{\mathbb{E}}_{s \sim {\mathbf{π}}_{\beta,t}}{\lbrack D_{T\hspace{0pt}V}{(p{(s'|s,\mathbf{a})}||\hat{p}{(s'|s,\mathbf{a})})}\rbrack} \leq \epsilon_{m}$, and the policy divergences at level $k$ be bounded as $\max_{s,\mathbf{a}^{1:{k - 1}}}D_{T\hspace{0pt}V}{(\pi_{\beta,k}{(a^{k}|s,\mathbf{a}^{1:{k - 1}})}||\pi_{k}{(a^{k}|s,\mathbf{a}^{1:{k - 1}})})} \leq \epsilon_{\pi_{k}}$, where ${\mathbf{π}}_{\beta}$ is the data collecting policy for the model and $\hat{p}\hspace{0pt}{(\left. s' \middle| {s,\mathbf{a}} \right.)}$ is the transition distribution under the model. Then the model return $\hat{\eta}$ and true return $\eta$ of the policy $\mathbf{π}$ are bounded as:]{#Thmtheorem1.p1.8.8 .ltx_text .ltx_font_italic}

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $${\hat{\eta}\hspace{0pt}{\lbrack{\mathbf{π}}\rbrack}} \geq {{\eta\hspace{0pt}{\lbrack{\mathbf{π}}\rbrack}} - \underset{C\hspace{0pt}{(\epsilon_{m},\mathbf{\epsilon}_{\pi_{1:n}})}}{\underbrace{\lbrack{\frac{2\hspace{0pt}\gamma\hspace{0pt}r_{\max}\hspace{0pt}{({\epsilon_{m}+{2\hspace{0pt}{\sum_{k=1}^{n}\epsilon_{\pi_{k}}}}})}}{{({1-\gamma})}^{2}}+\frac{4\hspace{0pt}r_{\max}\hspace{0pt}{\sum_{k=1}^{n}\epsilon_{\pi_{k}}}}{({1-\gamma})}}\rbrack}}}$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
:::

:::: {#Thmtheorem1.1 .ltx_proof}
###### Proof. {#proof.-2 .ltx_title .ltx_runin .ltx_font_italic .ltx_title_proof}

::: {#Thmtheorem1.1.p1 .ltx_para}
[The proof is given in Appendix [[B]{.ltx_text .ltx_ref_tag}](#A2 "Appendix B Proofs of Theorem 1 ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}. ∎]{#Thmtheorem1.1.p1.1.1 .ltx_text}
:::
::::
::::::

:::: {#Thmremark2 .ltx_theorem .ltx_theorem_remark}
###### [[Remark 2]{#Thmremark2.1.1.1 .ltx_text .ltx_font_bold}]{.ltx_tag .ltx_tag_theorem}[.]{#Thmremark2.2.2 .ltx_text .ltx_font_bold} {#remark-2. .ltx_title .ltx_runin .ltx_title_theorem}

::: {#Thmremark2.p1 .ltx_para}
Theorem [[1]{.ltx_text .ltx_ref_tag}](#Thmtheorem1 "Theorem 1. ‣ 4.3 Theoretical Analysis ‣ 4 Sequential Communication ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} provides a useful relationship between the compounding errors and the policy update. As long as we improve the return under the true dynamic by more than the gap, $C\hspace{0pt}{(\epsilon_{m},\mathbf{\epsilon}_{\pi_{1:n}})}$, we can guarantee the policy improvement under the world model. If no such policy exists to overcome the gap, it implies the model error is too high, that is, there is a large discrepancy between the world model and true dynamics. Thus the order sequence obtained under the world model is not reliable. Such an order sequence is almost the same as a random one. Though a random order sequence also has the theoretical guarantee of Proposition [[2]{.ltx_text .ltx_ref_tag}](#Thmproposition2 "Proposition 2. ‣ 4.3 Theoretical Analysis ‣ 4 Sequential Communication ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}, we will show in Section [[5.2]{.ltx_text .ltx_ref_tag}](#S5.SS2 "5.2 Ablation Studies ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} that a random order sequence leads to a poor local optimum empirically.
:::
::::
:::::::::::::::
::::::::::::::::::::::::::::

:::::::::::::::: {#S5 .section .ltx_section}
## [5 ]{.ltx_tag .ltx_tag_section}Experiments {#experiments .ltx_title .ltx_title_section}

::: {#S5.p1 .ltx_para .ltx_noindent}
Sequential communication (SeqComm) is currently instantiated based on MAPPO (Yu et al., [2021](#bib.bib54){.ltx_ref}). We evaluate SeqComm on three tasks in multi-agent particle environment (MPE) (Lowe et al., [2017](#bib.bib25){.ltx_ref}) and four maps in StarCraft multi-agent challenge (SMAC) (Samvelyan et al., [2019](#bib.bib40){.ltx_ref}).
:::

::: {#S5.p2 .ltx_para .ltx_noindent}
For these experiments, we compare SeqComm against the following communication-free and communication-based CTDE baselines: MAPPO (Yu et al., [2021](#bib.bib54){.ltx_ref}), QMIX (Rashid et al., [2018](#bib.bib37){.ltx_ref}), IS (Kim et al., [2021](#bib.bib21){.ltx_ref}), TarMAC (Das et al., [2019](#bib.bib5){.ltx_ref}), and I2C (Ding et al., [2020](#bib.bib6){.ltx_ref}). In more detail, IS communicates predicted future trajectories (observations and actions), and predictions are made by the environment model. TarMAC uses the attention model to focus more on important incoming messages (the hidden states of observations). TarMAC is reproduced based on MAPPO instead of A2C in the original paper for better performance. I2C infers one-to-one communication to reduce the redundancy of messages (also conditioning on observations).
:::

::: {#S5.p3 .ltx_para .ltx_noindent}
In the experiments, all the methods are parameter-sharing for fast convergence. We have fine-tuned the baselines for a fair comparison. Please refer to Appendix [[E]{.ltx_text .ltx_ref_tag}](#A5 "Appendix E Experimental Settings ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} for experimental settings and Appendix [[F]{.ltx_text .ltx_ref_tag}](#A6 "Appendix F Implementation Details ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} for implementation details. All results are presented in terms of the mean and standard deviation of five runs with different random seeds.
:::

::::::: {#S5.SS1 .section .ltx_subsection}
### [5.1 ]{.ltx_tag .ltx_tag_subsection}Results {#results .ltx_title .ltx_title_subsection}

<figure id="S5.F4" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S5.F4.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x4.png" id="S5.F4.sf1.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="338" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(a) </span>PP</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S5.F4.sf2" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x5.png" id="S5.F4.sf2.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="338" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(b) </span>CN</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S5.F4.sf3" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x6.png" id="S5.F4.sf3.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="338" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(c) </span>KA</figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 4: </span>Learning curves in terms of the mean reward averaged over timesteps of SeqComm and baselines on three MPE tasks: <a href="#S5.F4.sf1" class="ltx_ref" title="In Figure 4 ‣ 5.1 Results ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"><span class="ltx_text ltx_ref_tag">4(a)</span></a> predator-prey, <a href="#S5.F4.sf2" class="ltx_ref" title="In Figure 4 ‣ 5.1 Results ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"><span class="ltx_text ltx_ref_tag">4(b)</span></a> cooperative navigation, and <a href="#S5.F4.sf3" class="ltx_ref" title="In Figure 4 ‣ 5.1 Results ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"><span class="ltx_text ltx_ref_tag">4(c)</span></a> keep-away.</figcaption>
</figure>

<figure id="S5.F5" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_4">
<figure id="S5.F5.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x7.png" id="S5.F5.sf1.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="372" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(a) </span>6h<span class="math inline">_</span>vs<span class="math inline">_</span>8z</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_4">
<figure id="S5.F5.sf2" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x8.png" id="S5.F5.sf2.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="372" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(b) </span>MMM2</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_4">
<figure id="S5.F5.sf3" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x9.png" id="S5.F5.sf3.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="372" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(c) </span>10m<span class="math inline">_</span>vs<span class="math inline">_</span>11m</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_4">
<figure id="S5.F5.sf4" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x10.png" id="S5.F5.sf4.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="372" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(d) </span>8m<span class="math inline">_</span>vs<span class="math inline">_</span>9m</figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 5: </span>Learning curves in terms of the win rate of SeqComm and baselines on four customized SMAC maps: <a href="#S5.F5.sf1" class="ltx_ref" title="In Figure 5 ‣ 5.1 Results ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"><span class="ltx_text ltx_ref_tag">5(a)</span></a> 6h<span class="math inline">_</span>vs<span class="math inline">_</span>8z, <a href="#S5.F5.sf2" class="ltx_ref" title="In Figure 5 ‣ 5.1 Results ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"><span class="ltx_text ltx_ref_tag">5(b)</span></a> MMM2, <a href="#S5.F5.sf3" class="ltx_ref" title="In Figure 5 ‣ 5.1 Results ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"><span class="ltx_text ltx_ref_tag">5(c)</span></a> 10m<span class="math inline">_</span>vs<span class="math inline">_</span>11m, and <a href="#S5.F5.sf4" class="ltx_ref" title="In Figure 5 ‣ 5.1 Results ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"><span class="ltx_text ltx_ref_tag">5(d)</span></a> 8m<span class="math inline">_</span>vs<span class="math inline">_</span>9m.</figcaption>
</figure>

::: {#S5.SS1.p1 .ltx_para .ltx_noindent}
[MPE.]{#S5.SS1.p1.1.1 .ltx_text .ltx_font_bold} We experiment on predator-prey (PP), cooperative navigation (CN), and keep-away (KA) in MPE. In PP, five predators (agents) try to capture three prey. In CN, five agents try to occupy five landmarks. In KA, three attackers (agents) try to occupy three landmarks, however, there are three defenders to push them away. In all three tasks, the size of agents is set to be larger than the original settings so that collisions occur more easily, following the settings in (Kim et al., [2021](#bib.bib21){.ltx_ref}). In addition, agents cannot observe any other agents, and this makes the task more difficult and communication more important. We can observe similar modifications in previous works (Foerster et al., [2016](#bib.bib10){.ltx_ref}; Ding et al., [2020](#bib.bib6){.ltx_ref}). After all, we want to demonstrate the superior over communication-based baselines, and communication-based methods are more suitable for scenarios with limited vision. More details about experimental settings are available in Appendix [[E]{.ltx_text .ltx_ref_tag}](#A5 "Appendix E Experimental Settings ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}.
:::

::: {#S5.SS1.p2 .ltx_para .ltx_noindent}
Figure [[4]{.ltx_text .ltx_ref_tag}](#S5.F4 "Figure 4 ‣ 5.1 Results ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} shows the learning curves of all the methods in terms of the mean reward averaged over timesteps in PP, CN, and KA. We can see that SeqComm converges to the highest mean reward compared with all the baselines. The results demonstrate the superiority of SeqComm. In more detail, all communication-based methods outperform MAPPO, indicating the necessity of communication in these difficult tasks. Apart from MAPPO, IS performs the worst since it may access inaccurate predicted information due to the circular dependencies. The substantial improvement SeqComm over I2C and TarMAC is attributed to that SeqComm allows agents to get more valuable action information for explicit coordination. The agents learned by SeqComm show sophisticated coordination strategies induced by the priority of decision-making, which can be witnessed by the visualization of agent behaviors. More details are given in Appendix [[C]{.ltx_text .ltx_ref_tag}](#A3 "Appendix C Additional Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}. Note that QMIX is omitted in the comparison for clear presentation since Yu et al. ([2021](#bib.bib54){.ltx_ref}) have shown QMIX and MAPPO exhibit similar performance in various MPE tasks.
:::

<figure id="S5.F6" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S5.F6.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x11.png" id="S5.F6.sf1.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="338" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(a) </span>PP</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S5.F6.sf2" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x12.png" id="S5.F6.sf2.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="338" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(b) </span>CN</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_3">
<figure id="S5.F6.sf3" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x13.png" id="S5.F6.sf3.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="338" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(c) </span>KA</figcaption>
</figure>
</div>
<div class="ltx_flex_break">

</div>
<div class="ltx_flex_cell ltx_flex_size_4">
<figure id="S5.F6.sf4" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x14.png" id="S5.F6.sf4.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="372" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(d) </span>6h<span class="math inline">_</span>vs<span class="math inline">_</span>8z</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_4">
<figure id="S5.F6.sf5" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x15.png" id="S5.F6.sf5.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="372" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(e) </span>MMM2</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_4">
<figure id="S5.F6.sf6" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x16.png" id="S5.F6.sf6.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="372" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(f) </span>10m<span class="math inline">_</span>vs<span class="math inline">_</span>11m</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_4">
<figure id="S5.F6.sf7" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x17.png" id="S5.F6.sf7.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="372" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(g) </span>8m<span class="math inline">_</span>vs<span class="math inline">_</span>9m</figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 6: </span>Ablation studies on the priority of decision-making in all the tasks. Fix-C: the priority of decision-making is fixed throughout one episode. Random-C: the priority of decision-making is determined randomly. TarMAC is also compared as a reference without explicit action coordination.</figcaption>
</figure>

::: {#S5.SS1.p3 .ltx_para .ltx_noindent}
[SMAC.]{#S5.SS1.p3.8.1 .ltx_text .ltx_font_bold} We also evaluate SeqComm against the baselines on four customized maps in SMAC: 6h$\_$vs$\_$`<!-- -->`{=html}8z, MMM2, 10m$\_$vs$\_$`<!-- -->`{=html}11m, and 8m$\_$vs$\_$`<!-- -->`{=html}9m, where we have made some minor changes to the observation part of agents to make it more difficult. Specifically, the sight range of agents is reduced from $9$ to $2$, and agents cannot perceive any information about their allies even if they are within the sight range. NDQ (Wang et al., [2020](#bib.bib50){.ltx_ref}) adopts a similar change to increase the difficulty of action coordination and demonstrates that the miscoordination problem is widespread in multi-agent learning. The rest settings remain the same as the default.
:::

::: {#S5.SS1.p4 .ltx_para .ltx_noindent}
The learning curves of SeqComm and the baselines in terms of the win rate are illustrated in Figure [[5]{.ltx_text .ltx_ref_tag}](#S5.F5 "Figure 5 ‣ 5.1 Results ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}. IS and I2C fail in this task and get a zero win rate because these two methods are built on MADDPG. However, MADDPG cannot work well in SMAC, especially when we reduce the sight range of agents, which is also supported by other studies (Papoudakis et al., [2021](#bib.bib30){.ltx_ref}). SeqComm and TarMAC converge to better performances than MAPPO and QMIX, which demonstrates the benefit of communication. Moreover, SeqComm outperforms TarMAC, which again verifies the gain of explicit action coordination.
:::
:::::::

::::::: {#S5.SS2 .section .ltx_subsection}
### [5.2 ]{.ltx_tag .ltx_tag_subsection}Ablation Studies {#ablation-studies .ltx_title .ltx_title_subsection}

::: {#S5.SS2.p1 .ltx_para .ltx_noindent}
[Priority of Decision-Making.]{#S5.SS2.p1.1.1 .ltx_text .ltx_font_bold} We compare SeqComm with two ablation baselines with only a difference in the priority of decision-making: the priority of decision-making is fixed throughout one episode, denoted as Fix-C, and the priority of decision-making is determined randomly at each timestep, denoted as Random-C. TarMAC is also compared as a reference without explicit action coordination.
:::

::: {#S5.SS2.p2 .ltx_para .ltx_noindent}
As depicted in Figure [[6]{.ltx_text .ltx_ref_tag}](#S5.F6 "Figure 6 ‣ 5.1 Results ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}, SeqComm achieves a higher mean reward or win rate than Fix-C, Random-C, and TarMAC in all the tasks. These results verify the importance of the priority of decision-making and the necessity to continuously adjust it during one episode. It is also demonstrated that SeqComm can provide a proper priority of decision-making. As discussed in Section [[4.3]{.ltx_text .ltx_ref_tag}](#S4.SS3 "4.3 Theoretical Analysis ‣ 4 Sequential Communication ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}, although Fix-C and Random-C also have the theoretical guarantee, they converge to poor local optima in practice. Moreover, Fix-C and Random-C show better performance than TarMAC in most tasks. This result accords with the hypothesis that the SE is likely to be Pareto superior to the average NE in games with a high cooperation level. Additionally, the learned policy of SeqComm can generalize well to the same task with a different number of agents in MPE, which is detailed in Appendix [[C]{.ltx_text .ltx_ref_tag}](#A3 "Appendix C Additional Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}.
:::

::: {#S5.SS2.p3 .ltx_para .ltx_noindent}
[Communication Range.]{#S5.SS2.p3.4.1 .ltx_text .ltx_font_bold} We also carry out ablation studies on communication range in MPE tasks. Note that communication range means how many nearest neighbors each agent is allowed to communicate with, following the setting in Ding et al. ([2020](#bib.bib6){.ltx_ref}). We reduce the communication range of SeqComm from $4$ to $2$ and $0$. As there are only three agents in KA, it is omitted in this study. The results are shown in Figure [[7]{.ltx_text .ltx_ref_tag}](#S5.F7 "Figure 7 ‣ 5.2 Ablation Studies ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}. Communication-based agents perform better than communication-free agents, which accords with the results of many previous studies. More importantly, the superiority of SeqComm with communication range $2$ over the corresponding TarMAC again demonstrates the effectiveness of sequential communication even in reduced communication ranges.
:::

<figure id="S5.F7" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S5.F7.sf1" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x18.png" id="S5.F7.sf1.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="338" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(a) </span>PP</figcaption>
</figure>
</div>
<div class="ltx_flex_cell ltx_flex_size_2">
<figure id="S5.F7.sf2" class="ltx_figure ltx_figure_panel ltx_align_center">
<img src="/html/2209.12713/assets/x19.png" id="S5.F7.sf2.g1" class="ltx_graphics ltx_centering ltx_img_landscape" width="465" height="338" alt="Refer to caption" />
<figcaption><span class="ltx_tag ltx_tag_figure">(b) </span>CN</figcaption>
</figure>
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 7: </span> Ablation studies on reduced communication range in <a href="#S5.F7.sf1" class="ltx_ref" title="In Figure 7 ‣ 5.2 Ablation Studies ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"><span class="ltx_text ltx_ref_tag">7(a)</span></a> predator-prey and <a href="#S5.F7.sf2" class="ltx_ref" title="In Figure 7 ‣ 5.2 Ablation Studies ‣ 5 Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"><span class="ltx_text ltx_ref_tag">7(b)</span></a> cooperative navigation.</figcaption>
</figure>

::: {#S5.SS2.p4 .ltx_para .ltx_noindent}
However, as the communication range decreases from $4$ to $2$, there is no performance reduction in these two MPE tasks. On the contrary, the agents with communication range $2$ perform the best. It accords with the results in I2C (Ding et al., [2020](#bib.bib6){.ltx_ref}) and ATOC (Jiang & Lu, [2018](#bib.bib18){.ltx_ref}) that redundant information can impair the learning process sometimes. In other settings, this conclusion might not be true. Moreover, since under our communication scheme agents can obtain more information, *i.e.,* the actual actions of others, it is more reasonable that SeqComm can still outperform other methods in reduced communication ranges.
:::
:::::::
::::::::::::::::

:::: {#S6 .section .ltx_section}
## [6 ]{.ltx_tag .ltx_tag_section}Conclusions {#conclusions .ltx_title .ltx_title_section}

::: {#S6.p1 .ltx_para .ltx_noindent}
We have proposed SeqComm, which enables agents explicitly coordinate with each other. SeqComm from an asynchronous perspective allows agents to make decisions sequentially. A two-phase communication scheme has been adopted for determining the priority of decision-making and communicating messages accordingly. Theoretically, we prove the policies learned by SeqComm are guaranteed to improve monotonically and converge. Empirically, it is demonstrated that SeqComm outperforms baselines in a variety of cooperative multi-agent tasks and SeqComm can provide a proper priority of decision-making.
:::
::::

::: {#bib .section .ltx_bibliography}
## References {#references .ltx_title .ltx_title_bibliography}

- [[Becker et al. (2004)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Raphen Becker, Shlomo Zilberstein, Victor R. Lesser, and Claudia V. Goldman. ]{.ltx_bibblock} [Solving transition independent decentralized markov decision processes. ]{.ltx_bibblock} [*J. Artif. Intell. Res.*, 22:423--455, 2004. ]{.ltx_bibblock}]{#bib.bib1}
- [[Böhmer et al. (2020)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Wendelin Böhmer, Vitaly Kurin, and Shimon Whiteson. ]{.ltx_bibblock} [Deep coordination graphs. ]{.ltx_bibblock} [In *International Conference on Machine Learning (ICML)*, 2020. ]{.ltx_bibblock}]{#bib.bib2}
- [[Boutilier (1996)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Craig Boutilier. ]{.ltx_bibblock} [Planning, learning and coordination in multiagent decision processes. ]{.ltx_bibblock} [In *Conference on Theoretical Aspects of Rationality and Knowledge*, 1996. ]{.ltx_bibblock}]{#bib.bib3}
- [[Busoniu et al. (2008)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Lucian Busoniu, Robert Babuska, and Bart De Schutter. ]{.ltx_bibblock} [A comprehensive survey of multiagent reinforcement learning. ]{.ltx_bibblock} [*IEEE Transactions on Systems, Man, and Cybernetics, Part C (Applications and Reviews)*, 38(2):156--172, 2008. ]{.ltx_bibblock}]{#bib.bib4}
- [[Das et al. (2019)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Abhishek Das, Théophile Gervet, Joshua Romoff, Dhruv Batra, Devi Parikh, Mike Rabbat, and Joelle Pineau. ]{.ltx_bibblock} [Tarmac: Targeted multi-agent communication. ]{.ltx_bibblock} [In *International Conference on Machine Learning (ICML)*, 2019. ]{.ltx_bibblock}]{#bib.bib5}
- [[Ding et al. (2020)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ziluo Ding, Tiejun Huang, and Zongqing Lu. ]{.ltx_bibblock} [Learning individually inferred communication for multi-agent cooperation. ]{.ltx_bibblock} [In *Advances in Neural Information Processing Systems (NeurIPS)*, 2020. ]{.ltx_bibblock}]{#bib.bib6}
- [[Du et al. (2021)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Yali Du, Yifan Zhao, Meng Fang, Jun Wang, Gangyan Xu, and Haifeng Zhang. ]{.ltx_bibblock} [Learning predictive communication by imagination in networked system control, 2021. ]{.ltx_bibblock}]{#bib.bib7}
- [[Feng et al. (2021)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Mingxiao Feng, Guozi Liu, Li Zhao, Lei Song, Jiang Bian, Tao Qin, Wengang Zhou, Houqiang Li, and Tie-Yan Liu. ]{.ltx_bibblock} [Multi-agent reinforcement learning with shared resource in inventory management. ]{.ltx_bibblock} [2021. ]{.ltx_bibblock}]{#bib.bib8}
- [[Fischer et al. (2004)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Felix Fischer, Michael Rovatsos, and Gerhard Weiss. ]{.ltx_bibblock} [Hierarchical reinforcement learning in communication-mediated multiagent coordination. ]{.ltx_bibblock} [In *International Joint Conference on Autonomous Agents and Multiagent Systems (AAMAS)*, 2004. ]{.ltx_bibblock}]{#bib.bib9}
- [[Foerster et al. (2016)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jakob Foerster, Ioannis Alexandros Assael, Nando de Freitas, and Shimon Whiteson. ]{.ltx_bibblock} [Learning to communicate with deep multi-agent reinforcement learning. ]{.ltx_bibblock} [In *Advances in Neural Information Processing Systems (NeurIPS)*, 2016. ]{.ltx_bibblock}]{#bib.bib10}
- [[Goldman & Zilberstein (2003)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Claudia V. Goldman and Shlomo Zilberstein. ]{.ltx_bibblock} [Optimizing information exchange in cooperative multi-agent systems. ]{.ltx_bibblock} [In *International Joint Conference on Autonomous Agents and Multiagent Systems (AAMAS)*, 2003. ]{.ltx_bibblock}]{#bib.bib11}
- [[Goldman & Zilberstein (2004)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Claudia V. Goldman and Shlomo Zilberstein. ]{.ltx_bibblock} [Decentralized control of cooperative systems: Categorization and complexity analysis. ]{.ltx_bibblock} [*J. Artif. Intell. Res.*, 22:143--174, 2004. ]{.ltx_bibblock}]{#bib.bib12}
- [[Greenwald et al. (2003)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Amy Greenwald, Keith Hall, and Roberto Serrano. ]{.ltx_bibblock} [Correlated q-learning. ]{.ltx_bibblock} [In *A comprehensive survey of multiagent reinforcement learning*, 2003. ]{.ltx_bibblock}]{#bib.bib13}
- [[Guestrin et al. (2002)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Carlos Guestrin, Michail Lagoudakis, and Ronald Parr. ]{.ltx_bibblock} [Coordinated reinforcement learning. ]{.ltx_bibblock} [In *International Conference on Machine Learning (ICML)*, 2002. ]{.ltx_bibblock}]{#bib.bib14}
- [[Gupta et al. (2017)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jayesh K Gupta, Maxim Egorov, and Mykel Kochenderfer. ]{.ltx_bibblock} [Cooperative multi-agent control using deep reinforcement learning. ]{.ltx_bibblock} [In *International Conference on Autonomous Agents and Multiagent Systems (AAMAS)*, 2017. ]{.ltx_bibblock}]{#bib.bib15}
- [[Janner et al. (2019)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Michael Janner, Justin Fu, Marvin Zhang, and Sergey Levine. ]{.ltx_bibblock} [When to trust your model: Model-based policy optimization. ]{.ltx_bibblock} [In *Advances in Neural Information Processing Systems (NeurIPS)*, 2019. ]{.ltx_bibblock}]{#bib.bib16}
- [[Jaques et al. (2019)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Natasha Jaques, Angeliki Lazaridou, Edward Hughes, Caglar Gulcehre, Pedro Ortega, Dj Strouse, Joel Z Leibo, and Nando De Freitas. ]{.ltx_bibblock} [Social influence as intrinsic motivation for multi-agent deep reinforcement learning. ]{.ltx_bibblock} [In *International Conference on Machine Learning (ICML)*, 2019. ]{.ltx_bibblock}]{#bib.bib17}
- [[Jiang & Lu (2018)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jiechuan Jiang and Zongqing Lu. ]{.ltx_bibblock} [Learning attentional communication for multi-agent cooperation. ]{.ltx_bibblock} [*Advances in Neural Information Processing Systems (NeurIPS)*, 2018. ]{.ltx_bibblock}]{#bib.bib18}
- [[Jiang et al. (2020)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jiechuan Jiang, Chen Dun, Tiejun Huang, and Zongqing Lu. ]{.ltx_bibblock} [Graph convolutional reinforcement learning. ]{.ltx_bibblock} [In *International Conference on Learning Representation (ICLR)*, 2020. ]{.ltx_bibblock}]{#bib.bib19}
- [[Kim et al. (2019)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Daewoo Kim, Sangwoo Moon, David Hostallero, Wan Ju Kang, Taeyoung Lee, Kyunghwan Son, and Yung Yi. ]{.ltx_bibblock} [Learning to schedule communication in multi-agent reinforcement learning. ]{.ltx_bibblock} [In *International Conference on Learning Representations (ICLR)*, 2019. ]{.ltx_bibblock}]{#bib.bib20}
- [[Kim et al. (2021)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Woojun Kim, Jongeui Park, and Youngchul Sung. ]{.ltx_bibblock} [Communication in multi-agent reinforcement learning: Intention sharing. ]{.ltx_bibblock} [In *International Conference on Learning Representations (ICLR)*, 2021. ]{.ltx_bibblock}]{#bib.bib21}
- [[Konan et al. (2022)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Sachin Konan, Esmaeil Seraj, and Matthew Gombolay. ]{.ltx_bibblock} [Iterated reasoning with mutual information in cooperative and byzantine decentralized teaming. ]{.ltx_bibblock} [In *International Conference on Learning Representations (ICLR)*, 2022. ]{.ltx_bibblock}]{#bib.bib22}
- [[Könönen (2004)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ville Könönen. ]{.ltx_bibblock} [Asymmetric multiagent reinforcement learning. ]{.ltx_bibblock} [*Web Intelligence and Agent Systems: An international journal*, 2(2):105--121, 2004. ]{.ltx_bibblock}]{#bib.bib23}
- [[Li et al. (2019)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Xihan Li, Jia Zhang, Jiang Bian, Yunhai Tong, and Tie-Yan Liu. ]{.ltx_bibblock} [A cooperative multi-agent reinforcement learning framework for resource balancing in complex logistics network. ]{.ltx_bibblock} [*arXiv preprint arXiv:1903.00714*, 2019. ]{.ltx_bibblock}]{#bib.bib24}
- [[Lowe et al. (2017)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ryan Lowe, Yi Wu, Aviv Tamar, Jean Harb, OpenAI Pieter Abbeel, and Igor Mordatch. ]{.ltx_bibblock} [Multi-agent actor-critic for mixed cooperative-competitive environments. ]{.ltx_bibblock} [In *Advances in Neural Information Processing Systems (NeurIPS)*, 2017. ]{.ltx_bibblock}]{#bib.bib25}
- [[Ma et al. (2019)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Hang Ma, Daniel Harabor, Peter J Stuckey, Jiaoyang Li, and Sven Koenig. ]{.ltx_bibblock} [Searching with consistent prioritization for multi-agent path finding. ]{.ltx_bibblock} [In *AAAI Conference on Artificial Intelligence (AAAI)*, 2019. ]{.ltx_bibblock}]{#bib.bib26}
- [[Nair et al. (2004)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ranjit Nair, Milind Tambe, Maayan Roth, and Makoto Yokoo. ]{.ltx_bibblock} [Communication for improving policy computation in distributed pomdps. ]{.ltx_bibblock} [In *International Joint Conference on Autonomous Agents and Multiagent Systems (AAMAS)*, 2004. ]{.ltx_bibblock}]{#bib.bib27}
- [[Oliehoek et al. (2007)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Frans A. Oliehoek, Matthijs T. J. Spaan, and Nikos A. Vlassis. ]{.ltx_bibblock} [Dec-pomdps with delayed communication. ]{.ltx_bibblock} [In *AAMAS Workshop on Multi-agent Sequential Decision Making in Uncertain Domains*, 2007. ]{.ltx_bibblock}]{#bib.bib28}
- [[Ooi & Wornell (1996)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ J.M. Ooi and G.W. Wornell. ]{.ltx_bibblock} [Decentralized control of a multiple access broadcast channel: performance bounds. ]{.ltx_bibblock} [In *IEEE Conference on Decision and Control*, 1996. ]{.ltx_bibblock}]{#bib.bib29}
- [[Papoudakis et al. (2021)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Georgios Papoudakis, Filippos Christianos, Lukas Schäfer, and Stefano V. Albrecht. ]{.ltx_bibblock} [Benchmarking multi-agent deep reinforcement learning algorithms in cooperative tasks, 2021. ]{.ltx_bibblock}]{#bib.bib30}
- [[Peng et al. (2017)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Peng Peng, Ying Wen, Yaodong Yang, Quan Yuan, Zhenkun Tang, Haitao Long, and Jun Wang. ]{.ltx_bibblock} [Multiagent bidirectionally-coordinated nets: Emergence of human-level coordination in learning to play starcraft combat games. ]{.ltx_bibblock} [*arXiv preprint arXiv:1703.10069*, 2017. ]{.ltx_bibblock}]{#bib.bib31}
- [[Prasad et al. (1998)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ MV Nagendra Prasad, Victor R Lesser, and Susan E Lander. ]{.ltx_bibblock} [Learning organizational roles for negotiated search in a multiagent system. ]{.ltx_bibblock} [*International Journal of Human-Computer Studies*, 48(1):51--67, 1998. ]{.ltx_bibblock}]{#bib.bib32}
- [[Pretorius et al. (2021)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Arnu Pretorius, Scott Cameron, Andries Petrus Smit, Elan van Biljon, Lawrence Francis, Femi Azeez, Alexandre Laterre, and Karim Beguir. ]{.ltx_bibblock} [Learning to communicate through imagination with model-based deep multi-agent reinforcement learning, 2021. ]{.ltx_bibblock}]{#bib.bib33}
- [[Pynadath & Tambe (2002)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ David V. Pynadath and Milind Tambe. ]{.ltx_bibblock} [The communicative multiagent team decision problem: Analyzing teamwork theories and models. ]{.ltx_bibblock} [*J. Artif. Intell. Res.*, 16:389--423, 2002. ]{.ltx_bibblock}]{#bib.bib34}
- [[Rabinowitz et al. (2018)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Neil Rabinowitz, Frank Perbet, Francis Song, Chiyuan Zhang, SM Ali Eslami, and Matthew Botvinick. ]{.ltx_bibblock} [Machine theory of mind. ]{.ltx_bibblock} [In *International Conference on Machine Learning (ICML)*, 2018. ]{.ltx_bibblock}]{#bib.bib35}
- [[Raileanu et al. (2018)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Roberta Raileanu, Emily Denton, Arthur Szlam, and Rob Fergus. ]{.ltx_bibblock} [Modeling others using oneself in multi-agent reinforcement learning. ]{.ltx_bibblock} [In *International Conference on Machine Learning (ICML)*, 2018. ]{.ltx_bibblock}]{#bib.bib36}
- [[Rashid et al. (2018)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tabish Rashid, Mikayel Samvelyan, Christian Schroeder de Witt, Gregory Farquhar, Jakob Foerster, and Shimon Whiteson. ]{.ltx_bibblock} [Qmix: Monotonic value function factorisation for deep multi-agent reinforcement learning. ]{.ltx_bibblock} [In *International Conference on Machine Learning (ICML)*, 2018. ]{.ltx_bibblock}]{#bib.bib37}
- [[Roth et al. (2005a)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Maayan Roth, Reid Simmons, and Manuela Veloso. ]{.ltx_bibblock} [Decentralized communication strategies for coordinated multi-agent policies. ]{.ltx_bibblock} [In Lynne E. Parker, Frank E. Schneider, and Alan C. Schultz (eds.), *Multi-Robot Systems. From Swarms to Intelligent Automata Volume III*. Springer, 2005a. ]{.ltx_bibblock}]{#bib.bib38}
- [[Roth et al. (2005b)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Maayan Roth, Reid G. Simmons, and Manuela M. Veloso. ]{.ltx_bibblock} [Reasoning about joint beliefs for execution-time communication decisions. ]{.ltx_bibblock} [In *International Joint Conference on Autonomous Agents and Multiagent Systems (AAMAS)*, 2005b. ]{.ltx_bibblock}]{#bib.bib39}
- [[Samvelyan et al. (2019)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Mikayel Samvelyan, Tabish Rashid, Christian Schroeder de Witt, Gregory Farquhar, Nantas Nardelli, Tim G. J. Rudner, Chia-Man Hung, Philiph H. S. Torr, Jakob Foerster, and Shimon Whiteson. ]{.ltx_bibblock} [The StarCraft Multi-Agent Challenge. ]{.ltx_bibblock} [*arXiv preprint arXiv:1902.04043*, 2019. ]{.ltx_bibblock}]{#bib.bib40}
- [[Schulman et al. (2015)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ John Schulman, Sergey Levine, Pieter Abbeel, Michael Jordan, and Philipp Moritz. ]{.ltx_bibblock} [Trust region policy optimization. ]{.ltx_bibblock} [In *International Conference on Machine Learning (ICML)*, 2015. ]{.ltx_bibblock}]{#bib.bib41}
- [[Singh et al. (2019)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Amanpreet Singh, Tushar Jain, and Sainbayar Sukhbaatar. ]{.ltx_bibblock} [Individualized controlled continuous communication model for multiagent cooperative and competitive tasks. ]{.ltx_bibblock} [In *International Conference on Learning Representations (ICLR)*, 2019. ]{.ltx_bibblock}]{#bib.bib42}
- [[Sodomka et al. (2013)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Eric Sodomka, Elizabeth Hilliard, Michael Littman, and Amy Greenwald. ]{.ltx_bibblock} [Coco-q: Learning in stochastic games with side payments. ]{.ltx_bibblock} [In *International Conference on Machine Learning (ICML)*, 2013. ]{.ltx_bibblock}]{#bib.bib43}
- [[Spaan et al. (2006)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Matthijs T. J. Spaan, Geoffrey J. Gordon, and Nikos Vlassis. ]{.ltx_bibblock} [Decentralized planning under uncertainty for teams of communicating agents. ]{.ltx_bibblock} [In *International Joint Conference on Autonomous Agents and Multiagent Systems (AAMAS)*, 2006. ]{.ltx_bibblock}]{#bib.bib44}
- [[Sukhbaatar et al. (2016)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Sainbayar Sukhbaatar, Rob Fergus, et al. ]{.ltx_bibblock} [Learning multiagent communication with backpropagation. ]{.ltx_bibblock} [In *Advances in Neural Information Processing Systems (NeurIPS)*, 2016. ]{.ltx_bibblock}]{#bib.bib45}
- [[Van Den Berg & Overmars (2005)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jur P Van Den Berg and Mark H Overmars. ]{.ltx_bibblock} [Prioritized motion planning for multiple robots. ]{.ltx_bibblock} [In *IEEE/RSJ International Conference on Intelligent Robots and Systems (IROS)*, 2005. ]{.ltx_bibblock}]{#bib.bib46}
- [[Vlassis (2007)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Nikos Vlassis. ]{.ltx_bibblock} [A concise introduction to multiagent systems and distributed artificial intelligence. ]{.ltx_bibblock} [*Synthesis Lectures on Artificial Intelligence and Machine Learning*, 1(1):1--71, 2007. ]{.ltx_bibblock}]{#bib.bib47}
- [[Von Stackelberg (2010)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Heinrich Von Stackelberg. ]{.ltx_bibblock} [*Market structure and equilibrium*. ]{.ltx_bibblock} [Springer Science & Business Media, 2010. ]{.ltx_bibblock}]{#bib.bib48}
- [[Wang et al. (2021a)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Jianhong Wang, Wangkun Xu, Yunjie Gu, Wenbin Song, and Tim C Green. ]{.ltx_bibblock} [Multi-agent reinforcement learning for active voltage control on power distribution networks. ]{.ltx_bibblock} [*Advances in Neural Information Processing Systems (NeurIPS)*, 2021a. ]{.ltx_bibblock}]{#bib.bib49}
- [[Wang et al. (2020)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tonghan Wang, Jianhao Wang, Chongyi Zheng, and Chongjie Zhang. ]{.ltx_bibblock} [Learning nearly decomposable value functions via communication minimization. ]{.ltx_bibblock} [In *International Conference on Learning Representation (ICLR)*, 2020. ]{.ltx_bibblock}]{#bib.bib50}
- [[Wang et al. (2021b)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Tonghan Wang, Liang Zeng, Weijun Dong, Qianlan Yang, Yang Yu, and Chongjie Zhang. ]{.ltx_bibblock} [Context-aware sparse deep coordination graphs. ]{.ltx_bibblock} [*arXiv preprint arXiv:2106.02886*, 2021b. ]{.ltx_bibblock}]{#bib.bib51}
- [[Wei et al. (2018)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ermo Wei, Drew Wicke, David Freelan, and Sean Luke. ]{.ltx_bibblock} [Multiagent soft q-learning. ]{.ltx_bibblock} [In *AAAI Spring Symposium Series*, 2018. ]{.ltx_bibblock}]{#bib.bib52}
- [[Wen et al. (2019)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Ying Wen, Yaodong Yang, Rui Luo, Jun Wang, and W Pan. ]{.ltx_bibblock} [Probabilistic recursive reasoning for multi-agent reinforcement learning. ]{.ltx_bibblock} [In *International Conference on Learning Representations (ICLR)*, 2019. ]{.ltx_bibblock}]{#bib.bib53}
- [[Yu et al. (2021)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Chao Yu, Akash Velu, Eugene Vinitsky, Yu Wang, Alexandre Bayen, and Yi Wu. ]{.ltx_bibblock} [The surprising effectiveness of mappo in cooperative, multi-agent games. ]{.ltx_bibblock} [*arXiv preprint arXiv:2103.01955*, 2021. ]{.ltx_bibblock}]{#bib.bib54}
- [[Zhang et al. (2020)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Haifeng Zhang, Weizhe Chen, Zeren Huang, Minne Li, Yaodong Yang, Weinan Zhang, and Jun Wang. ]{.ltx_bibblock} [Bi-level actor-critic for multi-agent coordination. ]{.ltx_bibblock} [In *AAAI Conference on Artificial Intelligence (AAAI)*, 2020. ]{.ltx_bibblock}]{#bib.bib55}
- [[Zhang et al. (2019)]{.ltx_tag .ltx_role_refnum .ltx_tag_bibitem} [ Sai Qian Zhang, Qi Zhang, and Jieyu Lin. ]{.ltx_bibblock} [Efficient communication in multi-agent reinforcement learning via variance based control. ]{.ltx_bibblock} [In *Advances in Neural Information Processing Systems (NeurIPS)*, 2019. ]{.ltx_bibblock}]{#bib.bib56}
:::

::::::::::::::::: {#A1 .section .ltx_appendix}
## [Appendix A ]{.ltx_tag .ltx_tag_appendix}Proofs of Proposition [[1]{.ltx_text .ltx_ref_tag}](#Thmproposition1 "Proposition 1. ‣ 3 Problem Formulation ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} and Proposition [[2]{.ltx_text .ltx_ref_tag}](#Thmproposition2 "Proposition 2. ‣ 4.3 Theoretical Analysis ‣ 4 Sequential Communication ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} {#appendix-a-proofs-of-proposition-1-and-proposition-2 .ltx_title .ltx_title_appendix}

:::: {#Thmlemma1 .ltx_theorem .ltx_theorem_lemma}
###### [[Lemma 1]{#Thmlemma1.1.1.1 .ltx_text .ltx_font_bold}]{.ltx_tag .ltx_tag_theorem}[ ]{#Thmlemma1.2.2 .ltx_text .ltx_font_bold}(Agent-by-Agent PPO)[.]{#Thmlemma1.3.3 .ltx_text .ltx_font_bold} {#lemma-1-agent-by-agent-ppo. .ltx_title .ltx_runin .ltx_title_theorem}

::: {#Thmlemma1.p1 .ltx_para}
[If we update the policy of each agent $i$ with TRPO Schulman et al. ([2015](#bib.bib41){.ltx_ref}) (or approximately PPO) when fixing all the other agent's policies, then the joint policy will improve monotonically.]{#Thmlemma1.p1.1.1 .ltx_text .ltx_font_italic}
:::
::::

:::::: {#A1.3 .ltx_proof}
###### Proof. {#proof.-3 .ltx_title .ltx_runin .ltx_font_italic .ltx_title_proof}

::: {#A1.1.p1 .ltx_para .ltx_noindent}
We consider the joint surrogate objective in TRPO $L_{{\mathbf{π}}_{old}}\hspace{0pt}{({\mathbf{π}}_{new})}$ where ${\mathbf{π}}_{old}$ is the joint policy before updating and ${\mathbf{π}}_{new}$ is the joint policy after updating.
:::

::: {#A1.2.p2 .ltx_para .ltx_noindent}
Given that $\pi_{new}^{- i} = \pi_{old}^{- i}$, we have:

  -- ------------------------------------------------------------ ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $L_{{\mathbf{π}}_{old}}\hspace{0pt}{({\mathbf{π}}_{new})}$   $= {{\mathbb{E}}_{a \sim {\mathbf{π}}_{new}}\hspace{0pt}{\lbrack{A_{{\mathbf{π}}_{old}}\hspace{0pt}{(s,{\mathbf{a}})}}\rbrack}}$                                                                                                                                                                                                                                                            
                                                                  $= {{\mathbb{E}}_{a \sim {\mathbf{π}}_{old}}\hspace{0pt}{\lbrack{\frac{{\mathbf{π}}_{new}\hspace{0pt}{(\left. {\mathbf{a}} \middle| s \right.)}}{{\mathbf{π}}_{old}\hspace{0pt}{(\left. {\mathbf{a}} \middle| s \right.)}}\hspace{0pt}A_{{\mathbf{π}}_{old}}\hspace{0pt}{(s,{\mathbf{a}})}}\rbrack}}$                                                                                       
                                                                  $= {{\mathbb{E}}_{a \sim {\mathbf{π}}_{old}}\hspace{0pt}{\lbrack{\frac{\pi_{new}^{i}\hspace{0pt}{(\left. a^{i} \middle| s \right.)}}{\pi_{old}^{i}\hspace{0pt}{(\left. a^{i} \middle| s \right.)}}\hspace{0pt}A_{{\mathbf{π}}_{old}}\hspace{0pt}{(s,{\mathbf{a}})}}\rbrack}}$                                                                                                               
                                                                  $= {{\mathbb{E}}_{a^{i} \sim \pi_{old}^{i}}\hspace{0pt}\left\lbrack {\frac{\pi_{new}^{i}\hspace{0pt}{(\left. a^{i} \middle| s \right.)}}{\pi_{old}^{i}\hspace{0pt}{(\left. a^{i} \middle| s \right.)}}\hspace{0pt}{\mathbb{E}}_{a^{- i} \sim \pi_{o\hspace{0pt}l\hspace{0pt}d}^{- i}}\hspace{0pt}{\lbrack{A_{{\mathbf{π}}_{old}}\hspace{0pt}{(s,a^{i},a^{- i})}}\rbrack}} \right\rbrack}$   
                                                                  $= {{\mathbb{E}}_{a^{i} \sim \pi_{old}^{i}}\hspace{0pt}\left\lbrack {\frac{\pi_{new}^{i}\hspace{0pt}{(\left. a^{i} \middle| s \right.)}}{\pi_{old}^{i}\hspace{0pt}{(\left. a^{i} \middle| s \right.)}}\hspace{0pt}A_{{\mathbf{π}}_{old}}^{i}\hspace{0pt}{(s,a^{i})}} \right\rbrack}$                                                                                                        
                                                                  ${= {L_{\pi_{old}^{i}}\hspace{0pt}{(\pi_{new}^{i})}}},$                                                                                                                                                                                                                                                                                                                                     
  -- ------------------------------------------------------------ ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --

where ${A_{{\mathbf{π}}_{old}}^{i}\hspace{0pt}{(s,a^{i})}} = {{\mathbb{E}}_{a^{- i} \sim \pi_{old}^{- i}}\hspace{0pt}{\lbrack{A_{{\mathbf{π}}_{old}}\hspace{0pt}{(s,a^{i},a^{- i})}}\rbrack}}$ is the individual advantage of agent $i$, and the third equation is from the condition $\pi_{new}^{- i} = \pi_{old}^{- i}$.
:::

::: {#A1.3.p3 .ltx_para}
With the result of TRPO, we have the following conclusion:

  -- --------------------------------------------------------------- ------------------------------------------------------------------------------------------------------------------------------------------------- --
     ${J\hspace{0pt}{(\pi_{new})}} - {J\hspace{0pt}{(\pi_{old})}}$   $\geq L_{{\mathbf{π}}_{old}}{({\mathbf{π}}_{new})} - {CD}_{KL}^{\max}{({\mathbf{π}}_{new}||{\mathbf{π}}_{old})}$                                  
                                                                     $= L_{\pi_{old}^{i}}{(\pi_{new}^{i})} - {CD}_{KL}^{\max}{(\pi_{new}^{i}||\pi_{old}^{i})}\quad{(\text{from~}\pi_{new}^{- i} = \pi_{old}^{- i})}$   
  -- --------------------------------------------------------------- ------------------------------------------------------------------------------------------------------------------------------------------------- --

This means the individual objective is the same as the joint objective so the monotonic improvement is guaranteed. ∎
:::
::::::

::: {#A1.p1 .ltx_para}
Then we can show the proof of Proposition [[1]{.ltx_text .ltx_ref_tag}](#Thmproposition1 "Proposition 1. ‣ 3 Problem Formulation ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}.
:::

:::::: {#A1.6 .ltx_proof}
###### Proof. {#proof.-4 .ltx_title .ltx_runin .ltx_font_italic .ltx_title_proof}

::: {#A1.4.p1 .ltx_para .ltx_noindent}
We will build a new MDP $\overset{\sim}{M}$ based on the original MDP. We keep the action space $\overset{\sim}{A} = A = \times_{i = 1}^{n}A^{i}$, where $A^{i}$ is the original action space of agent $i$. The new state space contains multiple layers. We define ${\overset{\sim}{S}}^{k} = S \times {( \times_{i = 1}^{k}A^{i})}$ for $k = {1,2,\cdots,{n - 1}}$ and ${\overset{\sim}{S}}^{0} = S$, where $S$ is the original state space. Then a new state ${\overset{\sim}{s}}^{k} \in {\overset{\sim}{S}}^{k}$ means that ${\overset{\sim}{s}}^{k} = {(s,a^{1},a^{2},\cdots,a^{k})}$. The total new state space is defined as $\overset{\sim}{S} = {\cup_{i = 0}^{n - 1}{\overset{\sim}{S}}^{i}}$. Next we define the transition probability $\overset{\sim}{P}$ as following:

  -- --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     ${{\overset{\sim}{P}\hspace{0pt}{(\left. {\overset{\sim}{s}}' \middle| {{\overset{\sim}{s}}^{k},a^{k + 1},a^{- {({k + 1})}}} \right.)}} = {\mathbb{1}\hspace{0pt}\left( {{\overset{\sim}{s}}' = {({\overset{\sim}{s}}^{k},a^{k + 1})}} \right)}},{k < {n - 1}}$                                                                                                       
     ${{{\overset{\sim}{P}\hspace{0pt}{(\left. {\overset{\sim}{s}}' \middle| {{\overset{\sim}{s}}^{k},a^{k + 1},a^{- {({k + 1})}}} \right.)}} = {\mathbb{1}\hspace{0pt}\left( {{\overset{\sim}{s}}' \in {\overset{\sim}{S}}^{0}} \right)\hspace{0pt}P\hspace{0pt}{(\left. {\overset{\sim}{s}}' \middle| {{\overset{\sim}{s}}^{k},a^{k + 1}} \right.)}}},{k = {n - 1}}}.$   
  -- --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --

This means that the state in the layer $k$ can only transition to the state in the layer $k + 1$ with the corresponding action, and the state in the layer $n - 1$ will transition to the layer 0 with the probability $P$ in the original MDP. The reward function $\overset{\sim}{r}$ is defined as following:

  -- ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     ${{\overset{\sim}{r}\hspace{0pt}{(\overset{\sim}{s},{\mathbf{a}})}} = {\mathbb{1}\hspace{0pt}\left( {\overset{\sim}{s} \in {\overset{\sim}{S}}_{0}} \right)\hspace{0pt}r\hspace{0pt}{(\overset{\sim}{s},{\mathbf{a}})}}}.$   
  -- ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --

This means the reward is only obtained when the state in layer 0 and the value is the same as the original reward function. Now we obtain the total definition of the new MDP $\overset{\sim}{M} = {\{\overset{\sim}{S},\overset{\sim}{A},\overset{\sim}{P},\overset{\sim}{r},\gamma\}}$.
:::

::: {#A1.5.p2 .ltx_para .ltx_noindent}
Then we claim that if all agents learn in multi-agent sequential decision-making by PPO, they are actually taking agent-by-agent PPO in the new MDP $\overset{\sim}{M}$. To be precise, one update of multi-agent sequential decision-making in the original MDP $M$ equals to a round of update from agent 1 to agent $n$ by agent-by-agent PPO in the new MDP $\overset{\sim}{M}$. Moreover, the total reward of a round in the new MDP $\overset{\sim}{M}$ is the same as the reward in one timestep in the original MDP $M$. With this conclusion and Lemma [[1]{.ltx_text .ltx_ref_tag}](#Thmlemma1 "Lemma 1 (Agent-by-Agent PPO). ‣ Appendix A Proofs of Proposition 1 and Proposition 2 ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}, we complete the proof.
:::

::: {#A1.6.p3 .ltx_para}
∎
:::
::::::

::: {#A1.p2 .ltx_para .ltx_noindent}
The proof of Proposition [[2]{.ltx_text .ltx_ref_tag}](#Thmproposition2 "Proposition 2. ‣ 4.3 Theoretical Analysis ‣ 4 Sequential Communication ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} can be seen as a corollary of the proof of Proposition [[1]{.ltx_text .ltx_ref_tag}](#Thmproposition1 "Proposition 1. ‣ 3 Problem Formulation ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}.
:::

:::: {#A1.7 .ltx_proof}
###### Proof. {#proof.-5 .ltx_title .ltx_runin .ltx_font_italic .ltx_title_proof}

::: {#A1.7.p1 .ltx_para}
From Lemma [[1]{.ltx_text .ltx_ref_tag}](#Thmlemma1 "Lemma 1 (Agent-by-Agent PPO). ‣ Appendix A Proofs of Proposition 1 and Proposition 2 ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} we know that the monotonic improvement of the joint policy in the new MDP $\overset{\sim}{M}$ is guaranteed for each update of one single agent's policy. So even if the different round of updates in the new MDP $\overset{\sim}{M}$ is with different order of the decision-making, the monotonic improvement of the joint policy is still guaranteed. Finally, from the proof of Proposition [[1]{.ltx_text .ltx_ref_tag}](#Thmproposition1 "Proposition 1. ‣ 3 Problem Formulation ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}, we know that the monotonic improvement in the new MDP $\overset{\sim}{M}$ equals to the monotonic improvement in the original MDP $M$. These complete the proof. ∎
:::
::::
:::::::::::::::::

:::::::::::::::::::::::: {#A2 .section .ltx_appendix}
## [Appendix B ]{.ltx_tag .ltx_tag_appendix}Proofs of Theorem [[1]{.ltx_text .ltx_ref_tag}](#Thmtheorem1 "Theorem 1. ‣ 4.3 Theoretical Analysis ‣ 4 Sequential Communication ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} {#appendix-b-proofs-of-theorem-1 .ltx_title .ltx_title_appendix}

:::::: {#Thmlemma2 .ltx_theorem .ltx_theorem_lemma}
###### [[Lemma 2]{#Thmlemma2.2.1.1 .ltx_text .ltx_font_bold}]{.ltx_tag .ltx_tag_theorem}[ ]{#Thmlemma2.3.2 .ltx_text .ltx_font_bold}(TVD of the joint distributions)[.]{#Thmlemma2.4.3 .ltx_text .ltx_font_bold} {#lemma-2-tvd-of-the-joint-distributions. .ltx_title .ltx_runin .ltx_title_theorem}

::: {#Thmlemma2.p1 .ltx_para}
[Suppose we have two distribution ${p_{1}\hspace{0pt}{(x,y)}} = {p_{1}\hspace{0pt}{(x)}\hspace{0pt}p_{1}\hspace{0pt}{(\left. x \middle| y \right.)}}$ and ${p_{2}\hspace{0pt}{(x,y)}} = {p_{2}\hspace{0pt}{(x)}\hspace{0pt}p_{2}\hspace{0pt}{(\left. x \middle| y \right.)}}$. We can bound the total variation distance of the joint as:]{#Thmlemma2.p1.2.2 .ltx_text .ltx_font_italic}

  -- -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $$D_{T\hspace{0pt}V}{(p_{1}{(x,y)}||p_{2}{(x,y)})} \leq D_{T\hspace{0pt}V}{(p_{1}{(x)}||p_{2}{(x)})} + \max\limits_{x}D_{T\hspace{0pt}V}{(p_{1}{(y|x)}||p_{2}{(y|x)})}$$   
  -- -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
:::

:::: {#Thmlemma2.1 .ltx_proof}
###### Proof. {#proof.-6 .ltx_title .ltx_runin .ltx_font_italic .ltx_title_proof}

::: {#Thmlemma2.1.p1 .ltx_para}
[See (Janner et al., [2019](#bib.bib16){.ltx_ref}) (Lemma B.1). ∎]{#Thmlemma2.1.p1.1.1 .ltx_text}
:::
::::
::::::

::::::: {#Thmlemma3 .ltx_theorem .ltx_theorem_lemma}
###### [[Lemma 3]{#Thmlemma3.2.1.1 .ltx_text .ltx_font_bold}]{.ltx_tag .ltx_tag_theorem}[ ]{#Thmlemma3.3.2 .ltx_text .ltx_font_bold}(Markov chain TVD bound, time-varing)[.]{#Thmlemma3.4.3 .ltx_text .ltx_font_bold} {#lemma-3-markov-chain-tvd-bound-time-varing. .ltx_title .ltx_runin .ltx_title_theorem}

::: {#Thmlemma3.p1 .ltx_para .ltx_noindent}
[Suppose the expected KL-divergence between two transition is bounded as $\max_{t}{\mathbb{E}}_{s \sim {p_{1,t}\hspace{0pt}{(s)}}}D_{K\hspace{0pt}L}{(p_{1}{(s'|s)}||p_{2}{(s'|s)})} \leq \delta$, and the initial state distributions are the same ${p_{{1,t} = 0}\hspace{0pt}{(s)}} = {p_{{2,t} = 0}\hspace{0pt}{(s)}}$. Then the distance in the state marginal is bounded as:]{#Thmlemma3.p1.2.2 .ltx_text .ltx_font_italic}
:::

::: {#Thmlemma3.p2 .ltx_para .ltx_noindent}
  -- ------------------------------------------------------------------- --
     $$D_{T\hspace{0pt}V}{(p_{1,t}{(s)}||p_{2,t}{(s)})} \leq t\delta$$   
  -- ------------------------------------------------------------------- --
:::

:::: {#Thmlemma3.1 .ltx_proof}
###### Proof. {#proof.-7 .ltx_title .ltx_runin .ltx_font_italic .ltx_title_proof}

::: {#Thmlemma3.1.p1 .ltx_para}
[See (Janner et al., [2019](#bib.bib16){.ltx_ref}) (Lemma B.2). ∎]{#Thmlemma3.1.p1.1.1 .ltx_text}
:::
::::
:::::::

:::::::: {#Thmlemma4 .ltx_theorem .ltx_theorem_lemma}
###### [[Lemma 4]{#Thmlemma4.4.1.1 .ltx_text .ltx_font_bold}]{.ltx_tag .ltx_tag_theorem}[ ]{#Thmlemma4.5.2 .ltx_text .ltx_font_bold}(Branched Returns Bound)[.]{#Thmlemma4.6.3 .ltx_text .ltx_font_bold} {#lemma-4-branched-returns-bound. .ltx_title .ltx_runin .ltx_title_theorem}

::: {#Thmlemma4.p1 .ltx_para .ltx_noindent}
[Suppose the expected KL-divergence between two dynamics distributions is bounded as $\max_{t}{\mathbb{E}}_{s \sim {p_{1,t}\hspace{0pt}{(s)}}}{\lbrack D_{T\hspace{0pt}V}{(p_{1}{(s'|s,\mathbf{a})}||p_{2}{(s'|s,\mathbf{a})})}\rbrack}$, and the policy divergences at level $k$ are bounded as $\max_{s,\mathbf{a}^{1:{k - 1}}}D_{T\hspace{0pt}V}{(\pi_{1}{(a^{k}|s,\mathbf{a}^{1:{k - 1}})}||\pi_{2}{(a^{k}|s,\mathbf{a}^{1:{k - 1}})})} \leq \epsilon_{\pi_{k}}$. Then the returns are bounded as:]{#Thmlemma4.p1.3.3 .ltx_text .ltx_font_italic}

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --
     $${{|{\eta_{1} - \eta_{2}}|} \leq {\frac{2\hspace{0pt}r_{\max}\hspace{0pt}\gamma\hspace{0pt}{({\epsilon_{m} + {\sum_{k = 1}^{n}\epsilon_{\pi_{k}}}})}}{{({1 - \gamma})}^{2}} + \frac{2\hspace{0pt}r_{\max}\hspace{0pt}{\sum_{k = 1}^{n}\epsilon_{\pi_{k}}}}{1 - \gamma}}},$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --

[where $r_{\max}$ is the upper bound of the reward function.]{#Thmlemma4.p1.4.1 .ltx_text .ltx_font_italic}
:::

:::::: {#Thmlemma4.3 .ltx_proof}
###### Proof. {#proof.-8 .ltx_title .ltx_runin .ltx_font_italic .ltx_title_proof}

::: {#Thmlemma4.1.p1 .ltx_para .ltx_noindent}
[Here, $\eta_{1}$ denotes the returns of ${\mathbf{π}}_{1}$ under dynamics $p_{1}\hspace{0pt}{(\left. s' \middle| {s,{\mathbf{a}}} \right.)}$, and $\eta_{2}$ denotes the returns of ${\mathbf{π}}_{2}$ under dynamics $p_{2}\hspace{0pt}{(\left. s' \middle| {s,{\mathbf{a}}} \right.)}$. Then we have]{#Thmlemma4.1.p1.6.6 .ltx_text}

  -- --------------------------- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $|{\eta_{1} - \eta_{2}}|$   $= {|{\sum\limits_{s,{\mathbf{a}}}{{({{p_{1}\hspace{0pt}{(s,{\mathbf{a}})}} - {p_{2}\hspace{0pt}{(s,{\mathbf{a}})}}})}\hspace{0pt}r\hspace{0pt}{(s,{\mathbf{a}})}}}|}$                                              
                                 $= {|{\sum\limits_{t}{\sum\limits_{s,{\mathbf{a}}}{\gamma^{t}\hspace{0pt}{({{p_{1,t}\hspace{0pt}{(s,{\mathbf{a}})}} - {p_{2,t}\hspace{0pt}{(s,{\mathbf{a}})}}})}\hspace{0pt}r\hspace{0pt}{(s,{\mathbf{a}})}}}}|}$   
                                 $\leq {\sum\limits_{t}{\sum\limits_{s,{\mathbf{a}}}{\gamma^{t}\hspace{0pt}{|{{p_{1,t}\hspace{0pt}{(s,{\mathbf{a}})}} - {p_{2,t}\hspace{0pt}{(s,{\mathbf{a}})}}}|}\hspace{0pt}r\hspace{0pt}{(s,{\mathbf{a}})}}}}$    
                                 ${\leq {r_{\max}\hspace{0pt}{\sum\limits_{t}{\sum\limits_{s,{\mathbf{a}}}{\gamma^{t}\hspace{0pt}{|{{p_{1,t}\hspace{0pt}{(s,{\mathbf{a}})}} - {p_{2,t}\hspace{0pt}{(s,{\mathbf{a}})}}}|}}}}}}.$                      
  -- --------------------------- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --

[By Lemma [[2]{.ltx_text .ltx_ref_tag}](#Thmlemma2 "Lemma 2 (TVD of the joint distributions). ‣ Appendix B Proofs of Theorem 1 ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}, we get]{#Thmlemma4.1.p1.7.1 .ltx_text}

  -- --------------------------------------------------------------------------------------------- -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $\max\limits_{s}D_{T\hspace{0pt}V}{(\pi_{1}{({\mathbf{a}}|s)}||\pi_{2}{({\mathbf{a}}|s)})}$   $\leq \max\limits_{s,a_{1}}D_{T\hspace{0pt}V}{(\pi_{1}{({\mathbf{a}}^{- 1}|s,a^{1})}||\pi_{2}{({\mathbf{a}}^{- 1}|s,a^{1})})}$                                                         
                                                                                                   $+ \max\limits_{s}D_{T\hspace{0pt}V}{(\pi_{1}{(a^{1}|s)}||\pi_{2}{(a^{1}|s)})}$                                                                                                        
                                                                                                   $\leq \cdots$                                                                                                                                                                          
                                                                                                   $\leq \sum\limits_{k = 1}^{n}\max\limits_{s,{\mathbf{a}}^{1:{k - 1}}}D_{T\hspace{0pt}V}{(\pi_{1}{(a^{k}|s,{\mathbf{a}}^{1:{k - 1}})}||\pi_{2}{(a^{k}|s,{\mathbf{a}}^{1:{k - 1}})})}$   
                                                                                                   ${\leq {\sum\limits_{k = 1}^{n}\epsilon_{\pi_{k}}}}.$                                                                                                                                  
  -- --------------------------------------------------------------------------------------------- -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
:::

::: {#Thmlemma4.2.p2 .ltx_para .ltx_noindent}
[We then apply Lemma [[3]{.ltx_text .ltx_ref_tag}](#Thmlemma3 "Lemma 3 (Markov chain TVD bound, time-varing). ‣ Appendix B Proofs of Theorem 1 ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}, using $\delta = {\epsilon_{m} + {\sum_{k = 1}^{n}\epsilon_{\pi_{k}}}}$ (via Lemma [[3]{.ltx_text .ltx_ref_tag}](#Thmlemma3 "Lemma 3 (Markov chain TVD bound, time-varing). ‣ Appendix B Proofs of Theorem 1 ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} and [[2]{.ltx_text .ltx_ref_tag}](#Thmlemma2 "Lemma 2 (TVD of the joint distributions). ‣ Appendix B Proofs of Theorem 1 ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}) to get]{#Thmlemma4.2.p2.1.1 .ltx_text}

  -- ---------------------------------------------------- -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
     $D_{T\hspace{0pt}V}{(p_{1,t}{(s)}||p_{2,t}{(s)})}$   $\leq t\max\limits_{t}E_{s \sim {p_{1,t}\hspace{0pt}{(s)}}}D_{T\hspace{0pt}V}{(p_{1,t}{(s'|s)}||p_{2,t}{(s'|s)})}$                                                         
                                                          $\leq t\max\limits_{t}E_{s \sim {p_{1,t}\hspace{0pt}{(s)}}}D_{T\hspace{0pt}V}{(p_{1,t}{(s',{\mathbf{a}}|s)}||p_{2,t}{(s',{\mathbf{a}}|s)})}$                               
                                                          $\leq t{(\max\limits_{t}E_{s \sim {p_{1,t}\hspace{0pt}{(s)}}}D_{T\hspace{0pt}V}{(p_{1,t}{(s'|s,{\mathbf{a}})}||p_{2,t}{(s'|s,{\mathbf{a}})})}}$                            
                                                          $+ \max\limits_{t}E_{s \sim {p_{1,t}\hspace{0pt}{(s)}}}\max\limits_{s}D_{T\hspace{0pt}V}{({\mathbf{π}}_{1,t}{({\mathbf{a}}|s)}||{\mathbf{π}}_{2,t}{({\mathbf{a}}|s)})})$   
                                                          $\leq {t\hspace{0pt}{({\epsilon_{m} + {\sum\limits_{k = 1}^{n}\epsilon_{\pi_{k}}}})}}$                                                                                     
  -- ---------------------------------------------------- -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- --
:::

::: {#Thmlemma4.3.p3 .ltx_para}
[And we also get $D_{T\hspace{0pt}V}{(p_{1,t}{(s,{\mathbf{a}})}||p_{2,t}{(s,{\mathbf{a}})})} \leq t{(\epsilon_{m} + \sum_{k = 1}^{n}\epsilon_{\pi_{k}})} + \sum_{k = 1}^{n}\epsilon_{\pi_{k}}$ by Lemma [[2]{.ltx_text .ltx_ref_tag}](#Thmlemma2 "Lemma 2 (TVD of the joint distributions). ‣ Appendix B Proofs of Theorem 1 ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}. Thus, by plugging this back, we get:]{#Thmlemma4.3.p3.1.1 .ltx_text}

  -- --------------------------- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --
     $|{\eta_{1} - \eta_{2}}|$   $\leq {r_{\max}\hspace{0pt}{\sum\limits_{t}{\sum\limits_{s,{\mathbf{a}}}{\gamma^{t}\hspace{0pt}{|{{p_{1,t}\hspace{0pt}{(s,{\mathbf{a}})}} - {p_{2,t}\hspace{0pt}{(s,{\mathbf{a}})}}}|}}}}}$                              
                                 $\leq {2\hspace{0pt}r_{\max}\hspace{0pt}{\sum\limits_{t}{\gamma^{t}\hspace{0pt}{({{t\hspace{0pt}{({\epsilon_{m} + {\sum\limits_{k = 1}^{n}\epsilon_{\pi_{k}}}})}} + {\sum\limits_{k = 1}^{n}\epsilon_{\pi_{k}}}})}}}}$   
                                 $\leq {2\hspace{0pt}r_{\max}\hspace{0pt}{({\frac{\gamma{(\epsilon_{m} + \sum_{k = 1}^{n}\epsilon_{\pi_{k}})})}{{({1 - \gamma})}^{2}} + \frac{\sum_{k = 1}^{n}\epsilon_{\pi_{k}}}{1 - \gamma}})}}$                        
  -- --------------------------- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --

[∎]{#Thmlemma4.3.p3.2.1 .ltx_text}
:::
::::::
::::::::

::: {#A2.p1 .ltx_para}
Then we can show the proof of Theorem [[1]{.ltx_text .ltx_ref_tag}](#Thmtheorem1 "Theorem 1. ‣ 4.3 Theoretical Analysis ‣ 4 Sequential Communication ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}.
:::

::::::: {#A2.4 .ltx_proof}
###### Proof. {#proof.-9 .ltx_title .ltx_runin .ltx_font_italic .ltx_title_proof}

::: {#A2.1.p1 .ltx_para .ltx_noindent}
Let ${\mathbf{π}}_{\beta}$ denote the data collecting policy. We use Lemma [[4]{.ltx_text .ltx_ref_tag}](#Thmlemma4 "Lemma 4 (Branched Returns Bound). ‣ Appendix B Proofs of Theorem 1 ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} to bound the returns, but it will require bounded model error under the new policy $\mathbf{π}$. Thus, we need to introduce ${\mathbf{π}}_{\beta}$ by adding and subtracting $\eta\hspace{0pt}{\lbrack{\mathbf{π}}_{\beta}\rbrack}$, to get:

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --
     $${{{\hat{\eta}\hspace{0pt}{\lbrack{\mathbf{π}}\rbrack}} - {\eta\hspace{0pt}{\lbrack{\mathbf{π}}\rbrack}}} = {{{{\hat{\eta}\hspace{0pt}{\lbrack{\mathbf{π}}\rbrack}} - {\eta\hspace{0pt}{\lbrack{\mathbf{π}}_{\beta}\rbrack}}} + {\eta\hspace{0pt}{\lbrack{\mathbf{π}}_{\beta}\rbrack}}} - {\eta\hspace{0pt}{\lbrack{\mathbf{π}}\rbrack}}}}.$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --
:::

::: {#A2.2.p2 .ltx_para .ltx_noindent}
we can bound $L_{1}$ and $L_{2}$ both using Lemma [[4]{.ltx_text .ltx_ref_tag}](#Thmlemma4 "Lemma 4 (Branched Returns Bound). ‣ Appendix B Proofs of Theorem 1 ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} by using $\delta = {\sum_{k = 1}^{n}\epsilon_{\pi_{k}}}$ and $\delta = {\epsilon_{m} + {\sum_{k = 1}^{n}\epsilon_{\pi_{k}}}}$ respectively, and obtain:

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --
     $$L_{1} \geq {{- \frac{2\hspace{0pt}\gamma\hspace{0pt}r_{\max}\hspace{0pt}{\sum_{k = 1}^{n}\epsilon_{\pi_{k}}}}{{({1 - \gamma})}^{2}}} - \frac{2\hspace{0pt}r_{\max}\hspace{0pt}{\sum_{k = 1}^{n}\epsilon_{\pi_{k}}}}{({1 - \gamma})}}$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --
:::

::: {#A2.3.p3 .ltx_para .ltx_noindent}
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --
     $${L_{2} \geq {{- \frac{2\hspace{0pt}\gamma\hspace{0pt}r_{\max}\hspace{0pt}{({\epsilon_{\pi_{m}} + {\sum_{k = 1}^{n}\epsilon_{\pi_{k}}}})}}{{({1 - \gamma})}^{2}}} - \frac{2\hspace{0pt}r_{\max}\hspace{0pt}{\sum_{k = 1}^{n}\epsilon_{\pi_{k}}}}{({1 - \gamma})}}}.$$   
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ --
:::

::: {#A2.4.p4 .ltx_para}
Adding these two bounds together yields the conclusion. ∎
:::
:::::::
::::::::::::::::::::::::

::::::::: {#A3 .section .ltx_appendix}
## [Appendix C ]{.ltx_tag .ltx_tag_appendix}Additional Experiments {#appendix-c-additional-experiments .ltx_title .ltx_title_appendix}

::::: {#A3.SS1 .section .ltx_subsection}
### [C.1 ]{.ltx_tag .ltx_tag_subsection}Illustration of Learned Priority of Decision-Making {#c.1-illustration-of-learned-priority-of-decision-making .ltx_title .ltx_title_subsection}

::: {#A3.SS1.p1 .ltx_para .ltx_noindent}
Figure [[8]{.ltx_text .ltx_ref_tag}](#A3.F8 "Figure 8 ‣ C.1 Illustration of Learned Priority of Decision-Making ‣ Appendix C Additional Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} (upper panel from $a$ to $e$) shows the priority order of decision-making determined by SeqComm in PP. Agent $2$ that is far away from other preys and predators is chosen to be the first-mover. If agents want to encircle and capture the preys, the agents ([e.g.]{#A3.SS1.p1.7.1 .ltx_text .ltx_font_italic}, agent 2 and 5) that are on the periphery of the encircling circle should hold upper-level positions since they are able to decide how to narrow the encirclement. In addition, agent $3$ makes decisions prior to agent $5$ so that collision can be avoided after agent $5$ obtains the intention of agent $3$.
:::

::: {#A3.SS1.p2 .ltx_para .ltx_noindent}
For CN, as illustrated in Figure [[8]{.ltx_text .ltx_ref_tag}](#A3.F8 "Figure 8 ‣ C.1 Illustration of Learned Priority of Decision-Making ‣ Appendix C Additional Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} (lower panel from $a$ to $e$), agent $2$ is far away from all the landmarks and all other agents are in a better position to occupy landmarks. Therefore, agents $2$ is chosen to be the first-mover, which is similar to the phenomenon observed in PP. Once it has determined the target to occupy, other agents (agent $5$ and $3$) can adjust their actions accordingly and avoid conflict of goals. Otherwise, if agent $5$ makes a decision first and chooses to occupy the closest landmark, then agent $2$ has to approach to a further landmark which would take more steps.
:::

<figure id="A3.F8" class="ltx_figure">
<div class="ltx_flex_figure">
<div class="ltx_flex_cell ltx_flex_size_1">
<img src="/html/2209.12713/assets/x20.png" id="A3.F8.g1" class="ltx_graphics ltx_centering ltx_figure_panel ltx_img_landscape" width="422" height="86" alt="Refer to caption" />
</div>
<div class="ltx_flex_break">

</div>
<div class="ltx_flex_cell ltx_flex_size_1">
<img src="/html/2209.12713/assets/x21.png" id="A3.F8.g2" class="ltx_graphics ltx_centering ltx_figure_panel ltx_img_landscape" width="422" height="86" alt="Refer to caption" />
</div>
</div>
<figcaption><span class="ltx_tag ltx_tag_figure">Figure 8: </span>Illustration of learned priority of decision making in PP (<span id="A3.F8.5.1" class="ltx_text ltx_font_italic">upper panel</span>) and CN (<span id="A3.F8.6.2" class="ltx_text ltx_font_italic">lower panel</span>). Preys (landmarks) are viewed in black and predators (agents) are viewed in grey in PP (CN). From <span id="A3.F8.7.3" class="ltx_text ltx_font_italic">a</span> to <span id="A3.F8.8.4" class="ltx_text ltx_font_italic">e</span>, shown is the priority order. The smaller the level index, the higher priority of decision-making is.</figcaption>
</figure>
:::::

::::: {#A3.SS2 .section .ltx_subsection}
### [C.2 ]{.ltx_tag .ltx_tag_subsection}Generalization {#c.2-generalization .ltx_title .ltx_title_subsection}

::: {#A3.SS2.p1 .ltx_para .ltx_noindent}
Generalization to different numbers of agents has always been a key problem in MARL. For most algorithms in communication, once the model is trained in one scenario, it is unlikely for agents to maintain relatively competitive performance in other scenarios with different numbers of agents. However, as we employ attention modules to process communicated messages so that agents can handle messages of different lengths. In addition, the module used to determine the priority of decision-making is also not restricted by the number of agents. Thus, we investigate whether SeqComm generalizes well to different numbers of agents in CN and PP.
:::

<figure id="S3.T1" class="ltx_table">
<table id="S3.T1.12" class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<thead class="ltx_thead">
<tr id="S3.T1.12.13.1" class="ltx_tr">
<th id="S3.T1.12.13.1.1" class="ltx_td ltx_th ltx_th_row ltx_border_tt" style="padding-left: 3.0pt; padding-right: 3.0pt"></th>
<th id="S3.T1.12.13.1.2" class="ltx_td ltx_align_center ltx_th ltx_th_column ltx_border_tt" style="padding-left: 3.0pt; padding-right: 3.0pt">Fix-C</th>
<th id="S3.T1.12.13.1.3" class="ltx_td ltx_nopad_r ltx_align_center ltx_th ltx_th_column ltx_border_tt" style="padding-left: 3.0pt; padding-right: 3.0pt">SeqComm</th>
</tr>
</thead>
<tbody class="ltx_tbody">
<tr id="S3.T1.4.4" class="ltx_tr">
<th id="S3.T1.4.4.5" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_t" style="padding-left: 3.0pt; padding-right: 3.0pt">3-agent in CN</th>
<td id="S3.T1.2.2.2" class="ltx_td ltx_align_center ltx_border_t" style="padding-left: 3.0pt; padding-right: 3.0pt"><span class="math inline">−0.83</span> <span class="math inline">±0.17</span></td>
<td id="S3.T1.4.4.4" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_t" style="padding-left: 3.0pt; padding-right: 3.0pt"><span class="math inline">−0.76</span> <span class="math inline">±0.08</span></td>
</tr>
<tr id="S3.T1.8.8" class="ltx_tr">
<th id="S3.T1.8.8.5" class="ltx_td ltx_align_left ltx_th ltx_th_row" style="padding-left: 3.0pt; padding-right: 3.0pt">7-agent in CN</th>
<td id="S3.T1.6.6.2" class="ltx_td ltx_align_center" style="padding-left: 3.0pt; padding-right: 3.0pt"><span class="math inline">−1.79</span> <span class="math inline">±0.15</span></td>
<td id="S3.T1.8.8.4" class="ltx_td ltx_nopad_r ltx_align_center" style="padding-left: 3.0pt; padding-right: 3.0pt"><span class="math inline">−1.57</span> <span class="math inline">±0.10</span></td>
</tr>
<tr id="S3.T1.12.12" class="ltx_tr">
<th id="S3.T1.12.12.5" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_bb" style="padding-left: 3.0pt; padding-right: 3.0pt">7-agent in PP</th>
<td id="S3.T1.10.10.2" class="ltx_td ltx_align_center ltx_border_bb" style="padding-left: 3.0pt; padding-right: 3.0pt"><span class="math inline">−1.89</span> <span class="math inline">±0.45</span></td>
<td id="S3.T1.12.12.4" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_bb" style="padding-left: 3.0pt; padding-right: 3.0pt"><span class="math inline">−1.31</span> <span class="math inline">±0.60</span></td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 1: </span>Mean reward in different tasks, averaged over timesteps, with 200 test trials.</figcaption>
</figure>

::: {#A3.SS2.p2 .ltx_para .ltx_noindent}
For both tasks, SeqComm is trained on 5-agent settings. Then, we test SeqComm in 3-agent and 7-agent settings of CN and 7-agent setting of PP. We use Fix-C trained [directly]{#A3.SS2.p2.1.1 .ltx_text .ltx_font_italic} on these test tasks to illustrate the performance of SeqComm. Note that the quantity of both landmarks and preys is adjusted according to the number of agents in CN and PP. The test results are shown in Table [[1]{.ltx_text .ltx_ref_tag}](#S3.T1 "Table 1 ‣ C.2 Generalization ‣ Appendix C Additional Experiments ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref}. SeqComm exhibits the superiority in CN and PP, demonstrating that SeqComm may have a good generalization to the number of agents. A thorough study of the generalization of SeqComm is left to future work.
:::
:::::
:::::::::

::::: {#A4 .section .ltx_appendix}
## [Appendix D ]{.ltx_tag .ltx_tag_appendix}Additional Related Work {#appendix-d-additional-related-work .ltx_title .ltx_title_appendix}

::: {#A4.p1 .ltx_para .ltx_noindent}
[Multi-Agent Path Finding (MAPF).]{#A4.p1.1.1 .ltx_text .ltx_font_bold} MAPF aims to plan collision-free paths for multiple agents on a given graph from their given start vertices to target vertices. In MAPF, prioritized planning is deeply coupled with collision avoidance (Van Den Berg & Overmars, [2005](#bib.bib46){.ltx_ref}; Ma et al., [2019](#bib.bib26){.ltx_ref}), where collision is used to design constraints or heuristics for planning. Unlike MAPF, our method couples the priority of decision-making with the learning objective and thus is more general. In addition, the different motivations and problem settings may lead to the incompatibility of the methods in the two fields.
:::

::: {#A4.p2 .ltx_para .ltx_noindent}
[Reinforcement Learning in Stackelberg Game.]{#A4.p2.1.1 .ltx_text .ltx_font_bold} Many studies (Könönen, [2004](#bib.bib23){.ltx_ref}; Sodomka et al., [2013](#bib.bib43){.ltx_ref}; Greenwald et al., [2003](#bib.bib13){.ltx_ref}; Zhang et al., [2020](#bib.bib55){.ltx_ref}) have investigated reinforcement learning in finding the Stackelberg equilibrium. Bi-AC (Zhang et al., [2020](#bib.bib55){.ltx_ref}) is a bi-level actor-critic method that allows agents to have different knowledge bases so that the Stackelberg equilibrium (SE) is possible to find. The actions still can be executed simultaneously and distributedly. It empirically studies the relationship between the cooperation level and the superiority of the SE over the Nash equilibrium. AQL (Könönen, [2004](#bib.bib23){.ltx_ref}) updates the Q-value by solving the SE in each iteration and can be regarded as the value-based version of Bi-AC. Existing work mainly focuses on two-agent settings and their order is fixed in advance. However, the fixed order can hardly be an optimal solution as we will show in the next section. To address this issue, we exploit agents' intentions to dynamically determine the priority of decision-making along the way of interacting with each other.
:::
:::::

:::: {#A5 .section .ltx_appendix}
## [Appendix E ]{.ltx_tag .ltx_tag_appendix}Experimental Settings {#appendix-e-experimental-settings .ltx_title .ltx_title_appendix}

::: {#A5.p1 .ltx_para .ltx_noindent}
In cooperative navigation, there are 5 agents and the size of each is 0.15. They need to occupy 5 landmarks with the size of 0.05. The acceleration of agents is 7. In predator-prey, the number of predators (agents) and prey is set to 5 and 3, respectively, and their sizes are 0.15 and 0.05. The acceleration is 5 for predators and 7 for prey. In keep away, the number of attackers (agents) and defenders is set to 3, and their sizes are respectively 0.15 and 0.05. Besides, the acceleration is 6 for attackers and 4 for defenders. The three landmarks are located at $(0.00,0.30)$, $(0.25,{- 0.15})$, and $({- 0.25},{- 0.15})$. Note that each agent is allowed to communicate with all other agents in all three tasks. The team reward is similar across tasks. At a timestep $t$, it can be written as $r_{\text{team}}^{t} = {{- {\sum_{i = 1}^{n}d_{i}^{t}}} + {C^{t}\hspace{0pt}r_{\text{collision}}}}$, where $d_{i}^{t}$ is the distance of landmark/prey $i$ to its nearest agent/predator, $C^{t}$ is the number of collisions (when the distance between two agents is less than the sum of their sizes) occurred at timestep $t$, and $r_{\text{collision}} = {- 1}$. In addition, agents act discretely and have 5 actions (stay and move up, down, left, right). The length of each episode is 20, 30, and 20 in cooperative navigation, predator-prey, and keep-away, respectively.
:::
::::

::::::::::::::::: {#A6 .section .ltx_appendix}
## [Appendix F ]{.ltx_tag .ltx_tag_appendix}Implementation Details {#appendix-f-implementation-details .ltx_title .ltx_title_appendix}

:::::: {#A6.SS1 .section .ltx_subsection}
### [F.1 ]{.ltx_tag .ltx_tag_subsection}Architecture and Hyperparameters {#f.1-architecture-and-hyperparameters .ltx_title .ltx_title_subsection}

::: {#A6.SS1.p1 .ltx_para .ltx_noindent}
Our models, including SeqComm, Fix-C, and Random-C are trained based on MAPPO. The critic and policy network are realized by two fully connected layers. As for the attention module, key, query, and value have one fully connected layer each. The size of hidden layers is 100. Tanh functions are used as nonlinearity. For I2C, we use their official code with default settings of basic hyperparameters and networks. As there is no released code of IS and TarMAC, we implement IS and TarMAC by ourselves, following the instructions mentioned in the original papers (Kim et al., [2021](#bib.bib21){.ltx_ref}; Das et al., [2019](#bib.bib5){.ltx_ref}).
:::

::: {#A6.SS1.p2 .ltx_para .ltx_noindent}
For the world model, observations and actions are firstly encoded by a fully connected layer. The output size for the observation encoder is 48, and the output size for the action encoder is 16. Then the outputs of the encoder will be passed into the attention module with the same structure aforementioned. Finally, we use a fully connected layer to decode. In these layers, Tanh is used as the nonlinearity.
:::

::: {#A6.SS1.p3 .ltx_para .ltx_noindent}
Table [[2]{.ltx_text .ltx_ref_tag}](#A6.T2 "Table 2 ‣ F.1 Architecture and Hyperparameters ‣ Appendix F Implementation Details ‣ Multi-Agent Sequential Decision-Making via Communication"){.ltx_ref} summarize the hyperparameters used by SeqComm and the baselines in the experiments.
:::

<figure id="A6.T2" class="ltx_table">
<table id="A6.T2.18" class="ltx_tabular ltx_centering ltx_guessed_headers ltx_align_middle">
<tbody class="ltx_tbody">
<tr id="A6.T2.18.19.1" class="ltx_tr">
<th id="A6.T2.18.19.1.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_tt" style="padding: 1pt 2.0pt">Hyperparameter</th>
<td id="A6.T2.18.19.1.2" class="ltx_td ltx_align_center ltx_border_tt" style="padding: 1pt 2.0pt">SeqComm</td>
<td id="A6.T2.18.19.1.3" class="ltx_td ltx_align_center ltx_border_tt" style="padding: 1pt 2.0pt">Random-C</td>
<td id="A6.T2.18.19.1.4" class="ltx_td ltx_align_center ltx_border_tt" style="padding: 1pt 2.0pt">Fix-C</td>
<td id="A6.T2.18.19.1.5" class="ltx_td ltx_align_center ltx_border_tt" style="padding: 1pt 2.0pt">TarMAC</td>
<td id="A6.T2.18.19.1.6" class="ltx_td ltx_align_center ltx_border_tt" style="padding: 1pt 2.0pt">MAPPO</td>
<td id="A6.T2.18.19.1.7" class="ltx_td ltx_align_center ltx_border_tt" style="padding: 1pt 2.0pt">I2C</td>
<td id="A6.T2.18.19.1.8" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_tt" style="padding: 1pt 2.0pt">IS</td>
</tr>
<tr id="A6.T2.1.1" class="ltx_tr">
<th id="A6.T2.1.1.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_t" style="padding: 1pt 2.0pt">discount (<span class="math inline"><em>γ</em></span>)</th>
<td colspan="7" id="A6.T2.1.1.2" class="ltx_td ltx_align_center ltx_border_t" style="padding: 1pt 2.0pt">0.95,0.95,0.95,0.99</td>
</tr>
<tr id="A6.T2.18.20.2" class="ltx_tr">
<th id="A6.T2.18.20.2.1" class="ltx_td ltx_align_left ltx_th ltx_th_row" style="padding: 1pt 2.0pt">batch size</th>
<td id="A6.T2.18.20.2.2" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.18.20.2.3" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.18.20.2.4" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.18.20.2.5" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.18.20.2.6" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.18.20.2.7" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">800</td>
<td id="A6.T2.18.20.2.8" class="ltx_td ltx_nopad_r ltx_align_center" style="padding: 1pt 2.0pt">1024</td>
</tr>
<tr id="A6.T2.2.2" class="ltx_tr">
<th id="A6.T2.2.2.2" class="ltx_td ltx_align_left ltx_th ltx_th_row" style="padding: 1pt 2.0pt">buffer capacity</th>
<td id="A6.T2.2.2.3" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.2.2.4" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.2.2.5" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.2.2.6" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.2.2.7" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td colspan="2" id="A6.T2.2.2.1" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt"><span class="math inline">1 <em>e</em><sup>6</sup></span></td>
</tr>
<tr id="A6.T2.18.21.3" class="ltx_tr">
<th id="A6.T2.18.21.3.1" class="ltx_td ltx_align_left ltx_th ltx_th_row" style="padding: 1pt 2.0pt">number of processes</th>
<td colspan="5" id="A6.T2.18.21.3.2" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">16,16,16,8</td>
<td id="A6.T2.18.21.3.3" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.18.21.3.4" class="ltx_td ltx_nopad_r ltx_align_center" style="padding: 1pt 2.0pt">–</td>
</tr>
<tr id="A6.T2.11.11" class="ltx_tr">
<th id="A6.T2.11.11.10" class="ltx_td ltx_align_left ltx_th ltx_th_row" style="padding: 1pt 2.0pt">learning rate</th>
<td colspan="5" id="A6.T2.6.6.4" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt"><span class="math inline">1.5 <em>e</em><sup>−5</sup></span>, <span class="math inline">1 <em>e</em><sup>−5</sup></span>, <span class="math inline">4 <em>e</em><sup>−5</sup></span>,<span class="math inline">5 <em>e</em><sup>−5</sup></span></td>
<td id="A6.T2.10.10.8" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt"><span class="math inline">1 <em>e</em><sup>−2</sup></span>, <span class="math inline">1 <em>e</em><sup>−3</sup></span>, <span class="math inline">1 <em>e</em><sup>−3</sup></span>,<span class="math inline">5 <em>e</em><sup>−5</sup></span></td>
<td id="A6.T2.11.11.9" class="ltx_td ltx_nopad_r ltx_align_center" style="padding: 1pt 2.0pt"><span class="math inline">1 <em>e</em><sup>−2</sup></span></td>
</tr>
<tr id="A6.T2.16.16" class="ltx_tr">
<th id="A6.T2.12.12.1" class="ltx_td ltx_align_left ltx_th ltx_th_row" style="padding: 1pt 2.0pt"><span class="math inline"><em>H</em></span></th>
<td id="A6.T2.16.16.5" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt"><span class="math inline">10</span>,<span class="math inline">10</span>,<span class="math inline">20</span>,<span class="math inline">5</span></td>
<td id="A6.T2.16.16.6" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.16.16.7" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.16.16.8" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.16.16.9" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.16.16.10" class="ltx_td ltx_align_center" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.16.16.11" class="ltx_td ltx_nopad_r ltx_align_center" style="padding: 1pt 2.0pt">–</td>
</tr>
<tr id="A6.T2.18.18" class="ltx_tr">
<th id="A6.T2.17.17.1" class="ltx_td ltx_align_left ltx_th ltx_th_row ltx_border_bb" style="padding: 1pt 2.0pt"><span class="math inline"><em>F</em></span></th>
<td colspan="2" id="A6.T2.18.18.2" class="ltx_td ltx_align_center ltx_border_bb" style="padding: 1pt 2.0pt"><span class="math inline">2, 2, 1, 1</span></td>
<td id="A6.T2.18.18.3" class="ltx_td ltx_align_center ltx_border_bb" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.18.18.4" class="ltx_td ltx_align_center ltx_border_bb" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.18.18.5" class="ltx_td ltx_align_center ltx_border_bb" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.18.18.6" class="ltx_td ltx_align_center ltx_border_bb" style="padding: 1pt 2.0pt">–</td>
<td id="A6.T2.18.18.7" class="ltx_td ltx_nopad_r ltx_align_center ltx_border_bb" style="padding: 1pt 2.0pt">–</td>
</tr>
</tbody>
</table>
<figcaption><span class="ltx_tag ltx_tag_table">Table 2: </span>Hyperparameters for predator-prey, cooperative navigation, keep-away, and SMAC</figcaption>
</figure>
::::::

:::::: {#A6.SS2 .section .ltx_subsection}
### [F.2 ]{.ltx_tag .ltx_tag_subsection}Attention Module {#f.2-attention-module .ltx_title .ltx_title_subsection}

::: {#A6.SS2.p1 .ltx_para .ltx_noindent}
Attention module (AM) is applied to process messages in the world model, critic network, and policy network. AM consists of three components: query, key, and values. The output of AM is the weighted sum of values, where the weight of value is determined by the dot product of the query and the corresponding key.
:::

::: {#A6.SS2.p2 .ltx_para .ltx_noindent}
For AM in the world model denoted as ${AM}_{w}$, agent $i$ gets messages ${\mathbf{m}}_{t}^{- i} = {\mathbf{h}}_{t}^{- i}$ from all other agents at timestep $t$ in negotiation phase, and predicts a query vector $q_{t}^{i}$ following ${AM}_{w,q}^{i}\hspace{0pt}{(h_{t}^{i})}$. The query is used to compute a dot product with keys ${\mathbf{k}}_{t} = {\lbrack k_{t}^{1},\cdots,k_{t}^{n}\rbrack}$. Note that $k_{t}^{j}$ is obtained by the message from agent $j$ following ${AM}_{a,k}^{i}\hspace{0pt}{(h_{t}^{j})}$ for $j \neq i$, and $k_{t}^{i}$ is from ${AM}_{{neg},k}^{i}\hspace{0pt}{(h_{t}^{i})}$. Besides, it is scaled by $1/\sqrt{d_{k}}$ followed by a softmax to obtain attention weights $\alpha$ for each value vector:

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ -- ----------------------------------------------------
     $$\alpha_{i} = {{softmax}\left\lbrack {\frac{q_{t}^{iT}\hspace{0pt}k_{t}^{1}}{\sqrt{d_{k}}}\hspace{0pt}\cdots\hspace{0pt}\underset{\alpha_{i\hspace{0pt}j}}{\underbrace{\frac{q_{t}^{iT}\hspace{0pt}k_{t}^{j}}{\sqrt{d_{k}}}}}\hspace{0pt}\cdots\hspace{0pt}\frac{q_{t}^{iT}\hspace{0pt}k_{t}^{n}}{\sqrt{d_{k}}}} \right\rbrack}$$      [(1)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ -- ----------------------------------------------------

The output of attention module is defined as: $c_{t}^{i} = {\sum_{j = 1}^{n}{\alpha_{i\hspace{0pt}j}\hspace{0pt}v_{t}^{j}}}$, where $v_{t}^{j}$ is obtained from messages or its own hidden state of observation following ${AM}_{w,v}^{i}\hspace{0pt}{( \cdot )}$.
:::

::: {#A6.SS2.p3 .ltx_para .ltx_noindent}
As for AM in the policy and critic network denoted as ${AM}_{a}$ , agent $i$ gets additional messages from upper-level agent in the launching phase. The message from upper-level and lower-level agent can be expanded as ${\mathbf{m}}_{t}^{u\hspace{0pt}p\hspace{0pt}p\hspace{0pt}e\hspace{0pt}r} = {\lbrack{\mathbf{h}}_{t}^{u\hspace{0pt}p\hspace{0pt}p\hspace{0pt}e\hspace{0pt}r},{\mathbf{a}}_{t}^{u\hspace{0pt}p\hspace{0pt}p\hspace{0pt}e\hspace{0pt}r}\rbrack}$ and ${\mathbf{m}}_{t}^{l\hspace{0pt}o\hspace{0pt}w\hspace{0pt}e\hspace{0pt}r} = {\lbrack{\mathbf{h}}_{t}^{l\hspace{0pt}o\hspace{0pt}w\hspace{0pt}e\hspace{0pt}r},0\rbrack}$, respectively. In addition, the query depends on agent's own hidden state of observation $h_{t}^{i}$, but keys and values are only from messages of other agents.
:::
::::::

:::::::: {#A6.SS3 .section .ltx_subsection}
### [F.3 ]{.ltx_tag .ltx_tag_subsection}Training {#f.3-training .ltx_title .ltx_title_subsection}

::: {#A6.SS3.p1 .ltx_para .ltx_noindent}
The training of SeqComm is an extension of MAPPO. The observation encoder $e$, the critic $V$, and the policy $\pi$ are respectively parameterized by $\theta_{e}$, $\theta_{v}$, $\theta_{\pi}$. Besides, the attention module ${AM}_{a}$ is parameterized by $\theta_{a}$ and takes as input the agent's hidden state, the messages (hidden states of other agents) in the negotiation phase, and the messages (the actions of upper-level agents) in launching phase. Let $\mathcal{D} = {\{\tau_{k}\}}_{k = 1}^{K}$ be a set of trajectories by running policy in the environment. Note that we drop time $t$ in the following notations for simplicity.
:::

::: {#A6.SS3.p2 .ltx_para .ltx_noindent}
The value function is fitted by regression on mean-squared error:

  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
     $${\mathcal{L}\hspace{0pt}{(\theta_{v},\theta_{a},\theta_{e})}} = {\frac{1}{K\hspace{0pt}T}\hspace{0pt}{\sum\limits_{\tau \in \mathcal{D}}{\sum\limits_{t = 0}^{T - 1}\left\| {{V\hspace{0pt}{({{AM}_{a}\hspace{0pt}{({e\hspace{0pt}{({\mathbf{o}})}},{\mathbf{a}}^{u\hspace{0pt}p\hspace{0pt}p\hspace{0pt}e\hspace{0pt}r})}})}} - \hat{R}} \right\|_{2}^{2}}}}$$      [(2)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
  -- ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------

where $\hat{R}$ is the discount rewards-to-go.
:::

::: {#A6.SS3.p3 .ltx_para .ltx_noindent}
We update the policy by maximizing the PPO-Clip objective:

  -- ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
     $$\begin{aligned}                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              [(3)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
     {\mathcal{L}\hspace{0pt}{(\theta_{\pi},\theta_{a},\theta_{e})}} & {= {\frac{1}{K\hspace{0pt}T}\hspace{0pt}{\sum\limits_{\tau \in \mathcal{D}}{\sum\limits_{t = 0}^{T - 1}{\min{({\frac{\pi\hspace{0pt}{(\left. a \middle| {{AM}_{a}\hspace{0pt}{({e\hspace{0pt}{({\mathbf{o}})}},{\mathbf{a}}^{u\hspace{0pt}p\hspace{0pt}p\hspace{0pt}e\hspace{0pt}r})}} \right.)}}{\pi_{o\hspace{0pt}l\hspace{0pt}d}\hspace{0pt}{(\left. a \middle| {{AM}_{a}\hspace{0pt}{({e\hspace{0pt}{({\mathbf{o}})}},{\mathbf{a}}^{u\hspace{0pt}p\hspace{0pt}p\hspace{0pt}e\hspace{0pt}r})}} \right.)}}\hspace{0pt}A_{\pi_{o\hspace{0pt}l\hspace{0pt}d}}},{g\hspace{0pt}{(\epsilon,A_{\pi_{o\hspace{0pt}l\hspace{0pt}d}})}})}}}}}}      
     \end{aligned}$$                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                
  -- ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------

where $g{(\epsilon,A)} = \left\{ \begin{array}{rcr}
{{({1 + \epsilon})}\hspace{0pt}A} & & {A \geq 0} \\
{{({1 - \epsilon})}\hspace{0pt}A} & & {A \leq 0}
\end{array} \right.$, and $A_{\pi_{o\hspace{0pt}l\hspace{0pt}d}}\hspace{0pt}{({\mathbf{o}},{\mathbf{a}}^{u\hspace{0pt}p\hspace{0pt}p\hspace{0pt}e\hspace{0pt}r},a)}$ is computed using the GAE method.
:::

::: {#A6.SS3.p4 .ltx_para .ltx_noindent}
The world model $\mathcal{M}$ is parameterized by $\theta_{w}$ is trained as a regression model using the training data set $\mathcal{S}$. It is updated with the loss:

  -- ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
     $${{\mathcal{L}\hspace{0pt}{(\theta_{w})}} = {\frac{1}{|\mathcal{S}|}\hspace{0pt}{\sum\limits_{{{\mathbf{o}},{\mathbf{a}},{\mathbf{o}}',r} \in \mathcal{S}}\left\| {{({\mathbf{o}}',r)} - {\mathcal{M}\hspace{0pt}{({{AM}_{w}\hspace{0pt}{({e\hspace{0pt}{({\mathbf{o}})}},{\mathbf{a}})}})}}} \right\|_{2}^{2}}}}.$$      [(4)]{.ltx_tag .ltx_tag_equation .ltx_align_right}
  -- ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- -- ----------------------------------------------------
:::

::: {#A6.SS3.p5 .ltx_para .ltx_noindent}
We trained our model on one GeForce GTX 1050 Ti and Intel(R) Core(TM) i9-9900K CPU @ 3.60GHz.
:::
::::::::
:::::::::::::::::
::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

::: ar5iv-footer
[◄](/html/2209.12712){.ar5iv-nav-button .ar5iv-nav-button-prev} [![ar5iv homepage](/assets/ar5iv.png){height="40"}](/){.ar5iv-home-button} [Feeling\
lucky?](/feeling_lucky){.ar5iv-text-button} [](/land_of_honey_and_milk){rel="nofollow" aria-hidden="true" tabindex="-1"} [Conversion\
report](/log/2209.12713){.ar5iv-text-button .ar5iv-severity-warning} [Report\
an issue](https://github.com/dginev/ar5iv/issues/new?template=improve-article--arxiv-id-.md&title=Improve+article+2209.12713){.ar5iv-text-button target="_blank"} [View original\
on arXiv](https://arxiv.org/abs/2209.12713){.ar5iv-text-button .arxiv-ui-theme}[►](/html/2209.12714){.ar5iv-nav-button .ar5iv-nav-button-next}
:::

[[]{.color-scheme-icon}](javascript:toggleColorScheme() "Toggle ar5iv color scheme"){.ar5iv-toggle-color-scheme} [Copyright](https://arxiv.org/help/license){.ar5iv-footer-button target="_blank"} [Privacy Policy](https://arxiv.org/help/policies/privacy_policy){.ar5iv-footer-button target="_blank"}

::: ltx_page_logo
Generated on Wed Mar 13 23:43:35 2024 by [[L[a]{.ltx_font_smallcaps style="position:relative; bottom:2.2pt;"}T[e]{.ltx_font_smallcaps style="font-size:120%;position:relative; bottom:-0.2ex;"}]{style="letter-spacing:-0.2em; margin-right:0.1em;"}[XML]{style="font-size:90%; position:relative; bottom:-0.2ex;"}![Mascot Sammy](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAsAAAAOCAYAAAD5YeaVAAAAAXNSR0IArs4c6QAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAALEwAACxMBAJqcGAAAAAd0SU1FB9wKExQZLWTEaOUAAAAddEVYdENvbW1lbnQAQ3JlYXRlZCB3aXRoIFRoZSBHSU1Q72QlbgAAAdpJREFUKM9tkL+L2nAARz9fPZNCKFapUn8kyI0e4iRHSR1Kb8ng0lJw6FYHFwv2LwhOpcWxTjeUunYqOmqd6hEoRDhtDWdA8ApRYsSUCDHNt5ul13vz4w0vWCgUnnEc975arX6ORqN3VqtVZbfbTQC4uEHANM3jSqXymFI6yWazP2KxWAXAL9zCUa1Wy2tXVxheKA9YNoR8Pt+aTqe4FVVVvz05O6MBhqUIBGk8Hn8HAOVy+T+XLJfLS4ZhTiRJgqIoVBRFIoric47jPnmeB1mW/9rr9ZpSSn3Lsmir1fJZlqWlUonKsvwWwD8ymc/nXwVBeLjf7xEKhdBut9Hr9WgmkyGEkJwsy5eHG5vN5g0AKIoCAEgkEkin0wQAfN9/cXPdheu6P33fBwB4ngcAcByHJpPJl+fn54mD3Gg0NrquXxeLRQAAwzAYj8cwTZPwPH9/sVg8PXweDAauqqr2cDjEer1GJBLBZDJBs9mE4zjwfZ85lAGg2+06hmGgXq+j3+/DsixYlgVN03a9Xu8jgCNCyIegIAgx13Vfd7vdu+FweG8YRkjXdWy329+dTgeSJD3ieZ7RNO0VAXAPwDEAO5VKndi2fWrb9jWl9Esul6PZbDY9Go1OZ7PZ9z/lyuD3OozU2wAAAABJRU5ErkJggg==)](http://dlmf.nist.gov/LaTeXML/){.ltx_LaTeXML_logo target="_blank"}
:::
:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
