const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["assets/quill-BLmY9xB4.js","assets/rolldown-runtime-QTnfLwEv.js"])))=>i.map(i=>d[i]);
import{C as e,E as t,F as n,G as r,H as i,J as a,M as o,P as s,Q as c,U as l,W as u,c as d,ft as f,h as p,ht as m,j as h,k as g,l as _,m as v,mt as y,o as b,pt as x,r as S,s as C,u as w,y as T}from"./runtime-core.esm-bundler-huI9Rd5Y.js";import{D as E,r as D}from"./basecomponent-DqsHbXkj.js";import{t as O}from"./button-D_kLtROP.js";import{E as k,h as A,m as j,s as M,t as N,u as P,y as F}from"./index-CZhlZQzp.js";import{t as I}from"./useI18n-C-t5nx3y.js";import{r as L}from"./responseHandler-BJxA-JZj.js";import{t as R}from"./tag-DiTIWH6q.js";import{t as z}from"./FormRow-7wyT--0X.js";import{t as B}from"./baseeditableholder-9FLZGK7-.js";import{t as V}from"./textarea-DGRt0261.js";import{t as H}from"./TextInput-CpjWOgI4.js";import{n as U,t as W}from"./column-CS1IQFbn.js";import{t as G}from"./select-C_ZXC6tO.js";import{t as K}from"./inputnumber-CXGvz2dB.js";import{t as q}from"./multiselect-BpzM-UXG.js";import{t as ee}from"./toggleswitch-DM2U_kot.js";import{t as te}from"./SkeletonTable-BXI63Axg.js";import{t as J}from"./ConfirmDeleteDialog-DJ131IHz.js";import{t as Y}from"./SelectLabel-_CYmrEvY.js";import{t as X}from"./SkeletonCard-DbG--FTJ.js";var ne={class:`space-y-4`},Z={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},re={class:`text-sm text-gray-500 dark:text-gray-400`},ie={class:`max-w-2xl`},ae={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},oe={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},se={class:`flex justify-end pt-2`},Q=`/api/v1/tenant/job-management/identifications`,ce={__name:`JobIdentificationSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},orgGradingId:{type:String,default:``},orgJobFamilyId:{type:String,default:``},gradingOptions:{type:Array,default:()=>[]},jobFamilyOptions:{type:Array,default:()=>[]}},emits:[`saved`],setup(e,{emit:n}){let r=n,i=e,{t:o}=I(),s=F(),u=a(!1),f=a(!0),h=a(``),v=a({}),y=a(``),x=a({grading_id:``}),S=b(()=>{let e=i.jobFamilyOptions.find(e=>e.value===i.orgJobFamilyId);return e?e.label:i.orgJobFamilyId||`-`});function T(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function E(){if(!i.orgId){f.value=!1;return}try{let e=(await P.get(Q,{params:{organization_id:i.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];y.value=t.id,x.value.grading_id=t.grading_id||i.orgGradingId||``}else x.value.grading_id=i.orgGradingId||``}catch{x.value.grading_id=i.orgGradingId||``}finally{f.value=!1}}async function D(){if(h.value=``,v.value={},!x.value.grading_id){h.value=o(`job_management.grading_required`);return}u.value=!0;try{let e={nomenclature:i.orgName||``,full_code:i.orgCode||``,grading_id:x.value.grading_id,organization_id:i.orgId};if(y.value)await P.put(`${Q}/${y.value}`,{grading_id:x.value.grading_id});else{let t=await P.post(Q,e);y.value=t.data?.data?.id||``}s.add({severity:`success`,summary:o(`message.success`),detail:o(`common.saved`),life:2e3}),r(`saved`)}catch(e){let t=T(e);Object.keys(t).length>0?(v.value=t,h.value=Object.values(t).join(`, `)):h.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{u.value=!1}}return t(E),(t,n)=>(g(),w(`div`,ne,[C(`div`,null,[C(`h2`,Z,m(c(o)(`job_management.identifications`)),1),C(`p`,re,m(c(o)(`job_management.identification_description`)),1)]),C(`div`,ie,[f.value?(g(),d(X,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(g(),w(`div`,ae,[p(z,{label:c(o)(`organization.job_family`)},{default:l(()=>[p(H,{"model-value":S.value,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(z,{label:c(o)(`organization.grading`)},{default:l(()=>[p(c(G),{modelValue:x.value.grading_id,"onUpdate:modelValue":n[0]||=e=>x.value.grading_id=e,options:e.gradingOptions,"option-label":`label`,"option-value":`value`,placeholder:c(o)(`organization.select_grading`),class:`w-full`,size:`small`,invalid:!!v.value.grading_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`]),h.value?(g(),w(`div`,oe,m(h.value),1)):_(``,!0),C(`div`,se,[p(c(O),{label:c(o)(`common.save`),icon:`pi pi-check`,size:`small`,loading:u.value,disabled:!x.value.grading_id,onClick:D},null,8,[`label`,`loading`,`disabled`])])]))])]))}},le={class:`space-y-4`},ue={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},de={class:`text-sm text-gray-500 dark:text-gray-400`},fe={class:`max-w-2xl`},pe={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},me={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},he={class:`flex justify-end gap-2 pt-2`},ge=`/api/v1/tenant/job-management/objectives`,_e={__name:`JobObjectiveSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:n}){let r=n,i=e,{t:o}=I(),s=F(),u=a(!1),h=a(!0),v=a(!1),y=a(``),b=a({}),x=a(``),S=a(!1),T=a(``),E=a({objective:``});function D(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function k(){if(!i.orgId){h.value=!1;return}try{let e=(await P.get(ge,{params:{organization_id:i.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];x.value=t.id,E.value.objective=t.objective||``}}catch{}finally{h.value=!1}}async function A(){y.value=``,b.value={},u.value=!0;try{let e={nomenclature:i.orgName||``,full_code:i.orgCode||``,objective:E.value.objective||``,organization_id:i.orgId};if(x.value)await P.put(`${ge}/${x.value}`,{objective:E.value.objective||``});else{let t=await P.post(ge,e);x.value=t.data?.data?.id||``}s.add({severity:`success`,summary:o(`message.success`),detail:o(`common.saved`),life:2e3}),r(`saved`)}catch(e){let t=D(e);Object.keys(t).length>0?(b.value=t,y.value=Object.values(t).join(`, `)):y.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{u.value=!1}}async function j(){if(x.value){v.value=!0,T.value=``;try{await P.delete(`${ge}/${x.value}`),S.value=!1,x.value=``,E.value.objective=``,r(`saved`),s.add({severity:`success`,summary:o(`message.success`),detail:o(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||o(`message.operation_failed`)}finally{v.value=!1}}}return t(k),(e,t)=>(g(),w(`div`,le,[C(`div`,null,[C(`h2`,ue,m(c(o)(`job_management.objectives`)),1),C(`p`,de,m(c(o)(`job_management.objective_description`)),1)]),C(`div`,fe,[h.value?(g(),d(X,{key:0,type:`detail`,count:1,rows:3,cols:`grid-cols-1`,padding:`p-5`})):(g(),w(`div`,pe,[p(z,{label:c(o)(`job_management.objective`)},{default:l(()=>[p(c(V),{modelValue:E.value.objective,"onUpdate:modelValue":t[0]||=e=>E.value.objective=e,rows:`3`,class:f([`w-full`,{"p-invalid":b.value.objective}]),placeholder:c(o)(`job_management.objective`)+`...`},null,8,[`modelValue`,`class`,`placeholder`])]),_:1},8,[`label`]),y.value?(g(),w(`div`,me,m(y.value),1)):_(``,!0),C(`div`,he,[x.value?(g(),d(c(O),{key:0,label:c(o)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[1]||=e=>S.value=!0},null,8,[`label`])):_(``,!0),p(c(O),{label:x.value?c(o)(`common.update`):c(o)(`common.save`),icon:`pi pi-check`,size:`small`,loading:u.value,disabled:u.value,onClick:A},null,8,[`label`,`loading`,`disabled`])])]))]),p(J,{visible:S.value,"onUpdate:visible":t[2]||=e=>S.value=e,loading:v.value,"error-msg":T.value,onConfirm:j,onCancel:t[3]||=e=>S.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},ve={class:`space-y-4`},ye={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},be={class:`text-sm text-gray-500 dark:text-gray-400`},xe={class:`max-w-2xl`},Se={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Ce={class:`pt-1`},we={class:`flex items-center gap-2 mb-3`},Te={class:`text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400`},Ee={class:`space-y-4`},De={class:`pt-4 border-t border-gray-200 dark:border-gray-700`},Oe={class:`flex items-center gap-2 mb-3`},ke={class:`text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400`},Ae={class:`space-y-4`},je={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Me={class:`flex justify-end gap-2 pt-2`},Ne=`/api/v1/tenant/job-management/education-experiences`,Pe={__name:`JobEduExpSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:n}){let r=n,i=e,{t:o}=I(),s=F(),u=a(!1),h=a(!0),v=a(!1),y=a(``),b=a({}),x=a(``),S=a(!1),T=a(``),E=a({education_id:``,education_major_id:[],job_family_id:[],experience_id:``}),D=a([]),k=a([]),A=a([]),j=a([]);async function M(){try{let[e,t,n,r]=await Promise.all([P.get(`/api/v1/tenant/job-management/values`,{params:{type:`education`,per_page:100}}),P.get(`/api/v1/tenant/job-management/values`,{params:{type:`experience`,per_page:100}}),P.get(`/api/v1/tenant/settings/education-majors?per_page=200`),P.get(`/api/v1/tenant/settings/job-families?per_page=100`)]);k.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),D.value=(t.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),A.value=(n.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id})),j.value=(r.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id}))}catch{}}async function N(){if(i.orgId)try{let e=(await P.get(Ne,{params:{organization_id:i.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];x.value=t.id,E.value.education_id=t.education_id||``,E.value.education_major_id=Array.isArray(t.education_major_id)?t.education_major_id:[],E.value.job_family_id=Array.isArray(t.job_family_id)?t.job_family_id:[],E.value.experience_id=t.experience_id||``}}catch{}}async function R(){y.value=``,b.value={},u.value=!0;try{let e={nomenclature:i.orgName||``,full_code:i.orgCode||``,education_id:E.value.education_id||null,education_major_id:E.value.education_major_id||[],job_family_id:E.value.job_family_id||[],experience_id:E.value.experience_id||null,organization_id:i.orgId};if(x.value)await P.put(`${Ne}/${x.value}`,{education_id:E.value.education_id||``,education_major_id:E.value.education_major_id||[],job_family_id:E.value.job_family_id||[],experience_id:E.value.experience_id||``});else{let t=await P.post(Ne,e);x.value=t.data?.data?.id||``}s.add({severity:`success`,summary:o(`message.success`),detail:o(`common.saved`),life:2e3}),r(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(b.value=t,y.value=Object.values(t).join(`, `)):y.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{u.value=!1}}async function B(){if(x.value){v.value=!0,T.value=``;try{await P.delete(`${Ne}/${x.value}`),S.value=!1,x.value=``,E.value.education_id=``,E.value.education_major_id=[],E.value.job_family_id=[],E.value.experience_id=``,r(`saved`),s.add({severity:`success`,summary:o(`message.success`),detail:o(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||o(`message.operation_failed`)}finally{v.value=!1}}}return t(async()=>{try{await Promise.all([M(),N()])}finally{h.value=!1}}),(e,t)=>(g(),w(`div`,ve,[C(`div`,null,[C(`h2`,ye,m(c(o)(`job_management.education_experience`)),1),C(`p`,be,m(c(o)(`job_management.education_experience_description`)),1)]),C(`div`,xe,[h.value?(g(),d(X,{key:0,type:`detail`,count:1,rows:6,cols:`grid-cols-1`,padding:`p-5`})):(g(),w(`div`,Se,[C(`div`,Ce,[C(`div`,we,[t[7]||=C(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400`},[C(`i`,{class:`pi pi-graduation-cap text-sm`})],-1),C(`h3`,Te,m(c(o)(`job_management.group_education`)),1),t[8]||=C(`div`,{class:`flex-1 border-t border-gray-200 dark:border-gray-700`},null,-1)]),C(`div`,Ee,[p(z,{label:c(o)(`job_management.education_level`),errors:b.value?.education_id},{default:l(()=>[p(Y,{modelValue:E.value.education_id,"onUpdate:modelValue":t[0]||=e=>E.value.education_id=e,options:k.value,placeholder:c(o)(`job_values.select_education`),class:f({"p-invalid":b.value?.education_id})},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),p(z,{label:c(o)(`job_management.education_major`),errors:b.value?.education_major_id},{default:l(()=>[p(c(q),{modelValue:E.value.education_major_id,"onUpdate:modelValue":t[1]||=e=>E.value.education_major_id=e,options:A.value,"option-label":`label`,"option-value":`value`,placeholder:c(o)(`common.select`),class:`w-full`,size:`small`,filter:``,showClear:``,display:`chip`,maxSelectedLabels:2,invalid:!!b.value.education_major_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`])])]),C(`div`,De,[C(`div`,Oe,[t[9]||=C(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400`},[C(`i`,{class:`pi pi-briefcase text-sm`})],-1),C(`h3`,ke,m(c(o)(`job_management.group_experience`)),1),t[10]||=C(`div`,{class:`flex-1 border-t border-gray-200 dark:border-gray-700`},null,-1)]),C(`div`,Ae,[p(z,{label:c(o)(`job_management.experience_range`),errors:b.value?.experience_id},{default:l(()=>[p(Y,{modelValue:E.value.experience_id,"onUpdate:modelValue":t[2]||=e=>E.value.experience_id=e,options:D.value,placeholder:c(o)(`common.select`),class:f({"p-invalid":b.value?.experience_id})},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),p(z,{label:c(o)(`job_management.job_family`),errors:b.value?.job_family_id},{default:l(()=>[p(c(q),{modelValue:E.value.job_family_id,"onUpdate:modelValue":t[3]||=e=>E.value.job_family_id=e,options:j.value,"option-label":`label`,"option-value":`value`,placeholder:c(o)(`common.select`),class:`w-full`,size:`small`,filter:``,showClear:``,display:`chip`,maxSelectedLabels:2,invalid:!!b.value.job_family_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`])])]),y.value?(g(),w(`div`,je,m(y.value),1)):_(``,!0),C(`div`,Me,[x.value?(g(),d(c(O),{key:0,label:c(o)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[4]||=e=>S.value=!0},null,8,[`label`])):_(``,!0),p(c(O),{label:x.value?c(o)(`common.update`):c(o)(`common.save`),icon:`pi pi-check`,size:`small`,loading:u.value,disabled:u.value,onClick:R},null,8,[`label`,`loading`,`disabled`])])]))]),p(J,{visible:S.value,"onUpdate:visible":t[5]||=e=>S.value=e,loading:v.value,"error-msg":T.value,onConfirm:B,onCancel:t[6]||=e=>S.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Fe=D.extend({name:`editor`,style:`
    /*!
* Quill Editor v1.3.3
* https://quilljs.com/
* Copyright (c) 2014, Jason Chen
* Copyright (c) 2013, salesforce.com
*/
    .ql-container {
        box-sizing: border-box;
        font-family: Helvetica, Arial, sans-serif;
        font-size: 13px;
        height: 100%;
        margin: 0;
        position: relative;
    }
    .ql-container.ql-disabled .ql-tooltip {
        visibility: hidden;
    }
    .ql-container.ql-disabled .ql-editor ul[data-checked] > li::before {
        pointer-events: none;
    }
    .ql-clipboard {
        inset-inline-start: -100000px;
        height: 1px;
        overflow-y: hidden;
        position: absolute;
        top: 50%;
    }
    .ql-clipboard p {
        margin: 0;
        padding: 0;
    }
    .ql-editor {
        box-sizing: border-box;
        line-height: 1.42;
        height: 100%;
        outline: none;
        overflow-y: auto;
        padding: 12px 15px;
        tab-size: 4;
        -moz-tab-size: 4;
        text-align: left;
        white-space: pre-wrap;
        word-wrap: break-word;
    }
    .ql-editor > * {
        cursor: text;
    }
    .ql-editor p,
    .ql-editor ol,
    .ql-editor ul,
    .ql-editor pre,
    .ql-editor blockquote,
    .ql-editor h1,
    .ql-editor h2,
    .ql-editor h3,
    .ql-editor h4,
    .ql-editor h5,
    .ql-editor h6 {
        margin: 0;
        padding: 0;
        counter-reset: list-1 list-2 list-3 list-4 list-5 list-6 list-7 list-8 list-9;
    }
    .ql-editor ol,
    .ql-editor ul {
        padding-inline-start: 1.5rem;
    }
    .ql-editor ol > li,
    .ql-editor ul > li {
        list-style-type: none;
    }
    .ql-editor ul > li::before {
        content: '\\2022';
    }
    .ql-editor ul[data-checked='true'],
    .ql-editor ul[data-checked='false'] {
        pointer-events: none;
    }
    .ql-editor ul[data-checked='true'] > li *,
    .ql-editor ul[data-checked='false'] > li * {
        pointer-events: all;
    }
    .ql-editor ul[data-checked='true'] > li::before,
    .ql-editor ul[data-checked='false'] > li::before {
        color: #777;
        cursor: pointer;
        pointer-events: all;
    }
    .ql-editor ul[data-checked='true'] > li::before {
        content: '\\2611';
    }
    .ql-editor ul[data-checked='false'] > li::before {
        content: '\\2610';
    }
    .ql-editor li::before {
        display: inline-block;
        white-space: nowrap;
        width: 1.2rem;
    }
    .ql-editor li:not(.ql-direction-rtl)::before {
        margin-inline-start: -1.5rem;
        margin-inline-end: 0.3rem;
        text-align: right;
    }
    .ql-editor li.ql-direction-rtl::before {
        margin-inline-start: 0.3rem;
        margin-inline-end: -1.5rem;
    }
    .ql-editor ol li:not(.ql-direction-rtl),
    .ql-editor ul li:not(.ql-direction-rtl) {
        padding-inline-start: 1.5rem;
    }
    .ql-editor ol li.ql-direction-rtl,
    .ql-editor ul li.ql-direction-rtl {
        padding-inline-end: 1.5rem;
    }
    .ql-editor ol li {
        counter-reset: list-1 list-2 list-3 list-4 list-5 list-6 list-7 list-8 list-9;
        counter-increment: list-0;
    }
    .ql-editor ol li:before {
        content: counter(list-0, decimal) '. ';
    }
    .ql-editor ol li.ql-indent-1 {
        counter-increment: list-1;
    }
    .ql-editor ol li.ql-indent-1:before {
        content: counter(list-1, lower-alpha) '. ';
    }
    .ql-editor ol li.ql-indent-1 {
        counter-reset: list-2 list-3 list-4 list-5 list-6 list-7 list-8 list-9;
    }
    .ql-editor ol li.ql-indent-2 {
        counter-increment: list-2;
    }
    .ql-editor ol li.ql-indent-2:before {
        content: counter(list-2, lower-roman) '. ';
    }
    .ql-editor ol li.ql-indent-2 {
        counter-reset: list-3 list-4 list-5 list-6 list-7 list-8 list-9;
    }
    .ql-editor ol li.ql-indent-3 {
        counter-increment: list-3;
    }
    .ql-editor ol li.ql-indent-3:before {
        content: counter(list-3, decimal) '. ';
    }
    .ql-editor ol li.ql-indent-3 {
        counter-reset: list-4 list-5 list-6 list-7 list-8 list-9;
    }
    .ql-editor ol li.ql-indent-4 {
        counter-increment: list-4;
    }
    .ql-editor ol li.ql-indent-4:before {
        content: counter(list-4, lower-alpha) '. ';
    }
    .ql-editor ol li.ql-indent-4 {
        counter-reset: list-5 list-6 list-7 list-8 list-9;
    }
    .ql-editor ol li.ql-indent-5 {
        counter-increment: list-5;
    }
    .ql-editor ol li.ql-indent-5:before {
        content: counter(list-5, lower-roman) '. ';
    }
    .ql-editor ol li.ql-indent-5 {
        counter-reset: list-6 list-7 list-8 list-9;
    }
    .ql-editor ol li.ql-indent-6 {
        counter-increment: list-6;
    }
    .ql-editor ol li.ql-indent-6:before {
        content: counter(list-6, decimal) '. ';
    }
    .ql-editor ol li.ql-indent-6 {
        counter-reset: list-7 list-8 list-9;
    }
    .ql-editor ol li.ql-indent-7 {
        counter-increment: list-7;
    }
    .ql-editor ol li.ql-indent-7:before {
        content: counter(list-7, lower-alpha) '. ';
    }
    .ql-editor ol li.ql-indent-7 {
        counter-reset: list-8 list-9;
    }
    .ql-editor ol li.ql-indent-8 {
        counter-increment: list-8;
    }
    .ql-editor ol li.ql-indent-8:before {
        content: counter(list-8, lower-roman) '. ';
    }
    .ql-editor ol li.ql-indent-8 {
        counter-reset: list-9;
    }
    .ql-editor ol li.ql-indent-9 {
        counter-increment: list-9;
    }
    .ql-editor ol li.ql-indent-9:before {
        content: counter(list-9, decimal) '. ';
    }
    .ql-editor .ql-video {
        display: block;
        max-width: 100%;
    }
    .ql-editor .ql-video.ql-align-center {
        margin: 0 auto;
    }
    .ql-editor .ql-video.ql-align-right {
        margin: 0 0 0 auto;
    }
    .ql-editor .ql-bg-black {
        background: #000;
    }
    .ql-editor .ql-bg-red {
        background: #e60000;
    }
    .ql-editor .ql-bg-orange {
        background: #f90;
    }
    .ql-editor .ql-bg-yellow {
        background: #ff0;
    }
    .ql-editor .ql-bg-green {
        background: #008a00;
    }
    .ql-editor .ql-bg-blue {
        background: #06c;
    }
    .ql-editor .ql-bg-purple {
        background: #93f;
    }
    .ql-editor .ql-color-white {
        color: #fff;
    }
    .ql-editor .ql-color-red {
        color: #e60000;
    }
    .ql-editor .ql-color-orange {
        color: #f90;
    }
    .ql-editor .ql-color-yellow {
        color: #ff0;
    }
    .ql-editor .ql-color-green {
        color: #008a00;
    }
    .ql-editor .ql-color-blue {
        color: #06c;
    }
    .ql-editor .ql-color-purple {
        color: #93f;
    }
    .ql-editor .ql-font-serif {
        font-family:
            Georgia,
            Times New Roman,
            serif;
    }
    .ql-editor .ql-font-monospace {
        font-family:
            Monaco,
            Courier New,
            monospace;
    }
    .ql-editor .ql-size-small {
        font-size: 0.75rem;
    }
    .ql-editor .ql-size-large {
        font-size: 1.5rem;
    }
    .ql-editor .ql-size-huge {
        font-size: 2.5rem;
    }
    .ql-editor .ql-direction-rtl {
        direction: rtl;
        text-align: inherit;
    }
    .ql-editor .ql-align-center {
        text-align: center;
    }
    .ql-editor .ql-align-justify {
        text-align: justify;
    }
    .ql-editor .ql-align-right {
        text-align: right;
    }
    .ql-editor.ql-blank::before {
        color: dt('form.field.placeholder.color');
        content: attr(data-placeholder);
        font-style: italic;
        inset-inline-start: 15px;
        pointer-events: none;
        position: absolute;
        inset-inline-end: 15px;
    }
    .ql-snow.ql-toolbar:after,
    .ql-snow .ql-toolbar:after {
        clear: both;
        content: '';
        display: table;
    }
    .ql-snow.ql-toolbar button,
    .ql-snow .ql-toolbar button {
        background: none;
        border: none;
        cursor: pointer;
        display: inline-block;
        float: left;
        height: 24px;
        padding-block: 3px;
        padding-inline: 5px;
        width: 28px;
    }
    .ql-snow.ql-toolbar button svg,
    .ql-snow .ql-toolbar button svg {
        float: left;
        height: 100%;
    }
    .ql-snow.ql-toolbar button:active:hover,
    .ql-snow .ql-toolbar button:active:hover {
        outline: none;
    }
    .ql-snow.ql-toolbar input.ql-image[type='file'],
    .ql-snow .ql-toolbar input.ql-image[type='file'] {
        display: none;
    }
    .ql-snow.ql-toolbar button:hover,
    .ql-snow .ql-toolbar button:hover,
    .ql-snow.ql-toolbar button:focus,
    .ql-snow .ql-toolbar button:focus,
    .ql-snow.ql-toolbar button.ql-active,
    .ql-snow .ql-toolbar button.ql-active,
    .ql-snow.ql-toolbar .ql-picker-label:hover,
    .ql-snow .ql-toolbar .ql-picker-label:hover,
    .ql-snow.ql-toolbar .ql-picker-label.ql-active,
    .ql-snow .ql-toolbar .ql-picker-label.ql-active,
    .ql-snow.ql-toolbar .ql-picker-item:hover,
    .ql-snow .ql-toolbar .ql-picker-item:hover,
    .ql-snow.ql-toolbar .ql-picker-item.ql-selected,
    .ql-snow .ql-toolbar .ql-picker-item.ql-selected {
        color: #06c;
    }
    .ql-snow.ql-toolbar button:hover .ql-fill,
    .ql-snow .ql-toolbar button:hover .ql-fill,
    .ql-snow.ql-toolbar button:focus .ql-fill,
    .ql-snow .ql-toolbar button:focus .ql-fill,
    .ql-snow.ql-toolbar button.ql-active .ql-fill,
    .ql-snow .ql-toolbar button.ql-active .ql-fill,
    .ql-snow.ql-toolbar .ql-picker-label:hover .ql-fill,
    .ql-snow .ql-toolbar .ql-picker-label:hover .ql-fill,
    .ql-snow.ql-toolbar .ql-picker-label.ql-active .ql-fill,
    .ql-snow .ql-toolbar .ql-picker-label.ql-active .ql-fill,
    .ql-snow.ql-toolbar .ql-picker-item:hover .ql-fill,
    .ql-snow .ql-toolbar .ql-picker-item:hover .ql-fill,
    .ql-snow.ql-toolbar .ql-picker-item.ql-selected .ql-fill,
    .ql-snow .ql-toolbar .ql-picker-item.ql-selected .ql-fill,
    .ql-snow.ql-toolbar button:hover .ql-stroke.ql-fill,
    .ql-snow .ql-toolbar button:hover .ql-stroke.ql-fill,
    .ql-snow.ql-toolbar button:focus .ql-stroke.ql-fill,
    .ql-snow .ql-toolbar button:focus .ql-stroke.ql-fill,
    .ql-snow.ql-toolbar button.ql-active .ql-stroke.ql-fill,
    .ql-snow .ql-toolbar button.ql-active .ql-stroke.ql-fill,
    .ql-snow.ql-toolbar .ql-picker-label:hover .ql-stroke.ql-fill,
    .ql-snow .ql-toolbar .ql-picker-label:hover .ql-stroke.ql-fill,
    .ql-snow.ql-toolbar .ql-picker-label.ql-active .ql-stroke.ql-fill,
    .ql-snow .ql-toolbar .ql-picker-label.ql-active .ql-stroke.ql-fill,
    .ql-snow.ql-toolbar .ql-picker-item:hover .ql-stroke.ql-fill,
    .ql-snow .ql-toolbar .ql-picker-item:hover .ql-stroke.ql-fill,
    .ql-snow.ql-toolbar .ql-picker-item.ql-selected .ql-stroke.ql-fill,
    .ql-snow .ql-toolbar .ql-picker-item.ql-selected .ql-stroke.ql-fill {
        fill: #06c;
    }
    .ql-snow.ql-toolbar button:hover .ql-stroke,
    .ql-snow .ql-toolbar button:hover .ql-stroke,
    .ql-snow.ql-toolbar button:focus .ql-stroke,
    .ql-snow .ql-toolbar button:focus .ql-stroke,
    .ql-snow.ql-toolbar button.ql-active .ql-stroke,
    .ql-snow .ql-toolbar button.ql-active .ql-stroke,
    .ql-snow.ql-toolbar .ql-picker-label:hover .ql-stroke,
    .ql-snow .ql-toolbar .ql-picker-label:hover .ql-stroke,
    .ql-snow.ql-toolbar .ql-picker-label.ql-active .ql-stroke,
    .ql-snow .ql-toolbar .ql-picker-label.ql-active .ql-stroke,
    .ql-snow.ql-toolbar .ql-picker-item:hover .ql-stroke,
    .ql-snow .ql-toolbar .ql-picker-item:hover .ql-stroke,
    .ql-snow.ql-toolbar .ql-picker-item.ql-selected .ql-stroke,
    .ql-snow .ql-toolbar .ql-picker-item.ql-selected .ql-stroke,
    .ql-snow.ql-toolbar button:hover .ql-stroke-miter,
    .ql-snow .ql-toolbar button:hover .ql-stroke-miter,
    .ql-snow.ql-toolbar button:focus .ql-stroke-miter,
    .ql-snow .ql-toolbar button:focus .ql-stroke-miter,
    .ql-snow.ql-toolbar button.ql-active .ql-stroke-miter,
    .ql-snow.ql-toolbar button.ql-active .ql-stroke-miter,
    .ql-snow.ql-toolbar .ql-picker-label:hover .ql-stroke-miter,
    .ql-snow .ql-toolbar .ql-picker-label:hover .ql-stroke-miter,
    .ql-snow.ql-toolbar .ql-picker-label.ql-active .ql-stroke-miter,
    .ql-snow .ql-toolbar .ql-picker-label.ql-active .ql-stroke-miter,
    .ql-snow.ql-toolbar .ql-picker-item:hover .ql-stroke-miter,
    .ql-snow .ql-toolbar .ql-picker-item:hover .ql-stroke-miter,
    .ql-snow.ql-toolbar .ql-picker-item.ql-selected .ql-stroke-miter,
    .ql-snow .ql-toolbar .ql-picker-item.ql-selected .ql-stroke-miter {
        stroke: #06c;
    }
    @media (pointer: coarse) {
        .ql-snow.ql-toolbar button:hover:not(.ql-active),
        .ql-snow .ql-toolbar button:hover:not(.ql-active) {
            color: #444;
        }
        .ql-snow.ql-toolbar button:hover:not(.ql-active) .ql-fill,
        .ql-snow .ql-toolbar button:hover:not(.ql-active) .ql-fill,
        .ql-snow.ql-toolbar button:hover:not(.ql-active) .ql-stroke.ql-fill,
        .ql-snow .ql-toolbar button:hover:not(.ql-active) .ql-stroke.ql-fill {
            fill: #444;
        }
        .ql-snow.ql-toolbar button:hover:not(.ql-active) .ql-stroke,
        .ql-snow .ql-toolbar button:hover:not(.ql-active) .ql-stroke,
        .ql-snow.ql-toolbar button:hover:not(.ql-active) .ql-stroke-miter,
        .ql-snow .ql-toolbar button:hover:not(.ql-active) .ql-stroke-miter {
            stroke: #444;
        }
    }
    .ql-snow {
        box-sizing: border-box;
    }
    .ql-snow * {
        box-sizing: border-box;
    }
    .ql-snow .ql-hidden {
        display: none;
    }
    .ql-snow .ql-out-bottom,
    .ql-snow .ql-out-top {
        visibility: hidden;
    }
    .ql-snow .ql-tooltip {
        position: absolute;
        transform: translateY(10px);
    }
    .ql-snow .ql-tooltip a {
        cursor: pointer;
        text-decoration: none;
    }
    .ql-snow .ql-tooltip.ql-flip {
        transform: translateY(-10px);
    }
    .ql-snow .ql-formats {
        display: inline-block;
        vertical-align: middle;
    }
    .ql-snow .ql-formats:after {
        clear: both;
        content: '';
        display: table;
    }
    .ql-snow .ql-stroke {
        fill: none;
        stroke: #444;
        stroke-linecap: round;
        stroke-linejoin: round;
        stroke-width: 2;
    }
    .ql-snow .ql-stroke-miter {
        fill: none;
        stroke: #444;
        stroke-miterlimit: 10;
        stroke-width: 2;
    }
    .ql-snow .ql-fill,
    .ql-snow .ql-stroke.ql-fill {
        fill: #444;
    }
    .ql-snow .ql-empty {
        fill: none;
    }
    .ql-snow .ql-even {
        fill-rule: evenodd;
    }
    .ql-snow .ql-thin,
    .ql-snow .ql-stroke.ql-thin {
        stroke-width: 1;
    }
    .ql-snow .ql-transparent {
        opacity: 0.4;
    }
    .ql-snow .ql-direction svg:last-child {
        display: none;
    }
    .ql-snow .ql-direction.ql-active svg:last-child {
        display: inline;
    }
    .ql-snow .ql-direction.ql-active svg:first-child {
        display: none;
    }
    .ql-snow .ql-editor h1 {
        font-size: 2rem;
    }
    .ql-snow .ql-editor h2 {
        font-size: 1.5rem;
    }
    .ql-snow .ql-editor h3 {
        font-size: 1.17rem;
    }
    .ql-snow .ql-editor h4 {
        font-size: 1rem;
    }
    .ql-snow .ql-editor h5 {
        font-size: 0.83rem;
    }
    .ql-snow .ql-editor h6 {
        font-size: 0.67rem;
    }
    .ql-snow .ql-editor a {
        text-decoration: underline;
    }
    .ql-snow .ql-editor blockquote {
        border-inline-start: 4px solid #ccc;
        margin-block-end: 5px;
        margin-block-start: 5px;
        padding-inline-start: 16px;
    }
    .ql-snow .ql-editor code,
    .ql-snow .ql-editor pre {
        background: #f0f0f0;
        border-radius: 3px;
    }
    .ql-snow .ql-editor pre {
        white-space: pre-wrap;
        margin-block-end: 5px;
        margin-block-start: 5px;
        padding: 5px 10px;
    }
    .ql-snow .ql-editor code {
        font-size: 85%;
        padding: 2px 4px;
    }
    .ql-snow .ql-editor pre.ql-syntax {
        background: #23241f;
        color: #f8f8f2;
        overflow: visible;
    }
    .ql-snow .ql-editor img {
        max-width: 100%;
    }
    .ql-snow .ql-picker {
        color: #444;
        display: inline-block;
        float: left;
        inset-inline-start: 0;
        font-size: 14px;
        font-weight: 500;
        height: 24px;
        position: relative;
        vertical-align: middle;
    }
    .ql-snow .ql-picker-label {
        cursor: pointer;
        display: inline-block;
        height: 100%;
        padding-inline-start: 8px;
        padding-inline-end: 2px;
        position: relative;
        width: 100%;
    }
    .ql-snow .ql-picker-label::before {
        display: inline-block;
        line-height: 22px;
    }
    .ql-snow .ql-picker-options {
        background: #fff;
        display: none;
        min-width: 100%;
        padding: 4px 8px;
        position: absolute;
        white-space: nowrap;
    }
    .ql-snow .ql-picker-options .ql-picker-item {
        cursor: pointer;
        display: block;
        padding-block-end: 5px;
        padding-block-start: 5px;
    }
    .ql-snow .ql-picker.ql-expanded .ql-picker-label {
        color: #ccc;
        z-index: 2;
    }
    .ql-snow .ql-picker.ql-expanded .ql-picker-label .ql-fill {
        fill: #ccc;
    }
    .ql-snow .ql-picker.ql-expanded .ql-picker-label .ql-stroke {
        stroke: #ccc;
    }
    .ql-snow .ql-picker.ql-expanded .ql-picker-options {
        display: block;
        margin-block-start: -1px;
        top: 100%;
        z-index: 1;
    }
    .ql-snow .ql-color-picker,
    .ql-snow .ql-icon-picker {
        width: 28px;
    }
    .ql-snow .ql-color-picker .ql-picker-label,
    .ql-snow .ql-icon-picker .ql-picker-label {
        padding: 2px 4px;
    }
    .ql-snow .ql-color-picker .ql-picker-label svg,
    .ql-snow .ql-icon-picker .ql-picker-label svg {
        inset-inline-end: 4px;
    }
    .ql-snow .ql-icon-picker .ql-picker-options {
        padding: 4px 0;
    }
    .ql-snow .ql-icon-picker .ql-picker-item {
        height: 24px;
        width: 24px;
        padding: 2px 4px;
    }
    .ql-snow .ql-color-picker .ql-picker-options {
        padding: 3px 5px;
        width: 152px;
    }
    .ql-snow .ql-color-picker .ql-picker-item {
        border: 1px solid transparent;
        float: left;
        height: 16px;
        margin: 2px;
        padding: 0;
        width: 16px;
    }
    .ql-snow .ql-picker:not(.ql-color-picker):not(.ql-icon-picker) svg {
        position: absolute;
        margin-block-start: -9px;
        inset-inline-end: 0;
        top: 50%;
        width: 18px;
    }
    .ql-snow .ql-picker.ql-header .ql-picker-label[data-label]:not([data-label=''])::before,
    .ql-snow .ql-picker.ql-font .ql-picker-label[data-label]:not([data-label=''])::before,
    .ql-snow .ql-picker.ql-size .ql-picker-label[data-label]:not([data-label=''])::before,
    .ql-snow .ql-picker.ql-header .ql-picker-item[data-label]:not([data-label=''])::before,
    .ql-snow .ql-picker.ql-font .ql-picker-item[data-label]:not([data-label=''])::before,
    .ql-snow .ql-picker.ql-size .ql-picker-item[data-label]:not([data-label=''])::before {
        content: attr(data-label);
    }
    .ql-snow .ql-picker.ql-header {
        width: 98px;
    }
    .ql-snow .ql-picker.ql-header .ql-picker-label::before,
    .ql-snow .ql-picker.ql-header .ql-picker-item::before {
        content: 'Normal';
    }
    .ql-snow .ql-picker.ql-header .ql-picker-label[data-value='1']::before,
    .ql-snow .ql-picker.ql-header .ql-picker-item[data-value='1']::before {
        content: 'Heading 1';
    }
    .ql-snow .ql-picker.ql-header .ql-picker-label[data-value='2']::before,
    .ql-snow .ql-picker.ql-header .ql-picker-item[data-value='2']::before {
        content: 'Heading 2';
    }
    .ql-snow .ql-picker.ql-header .ql-picker-label[data-value='3']::before,
    .ql-snow .ql-picker.ql-header .ql-picker-item[data-value='3']::before {
        content: 'Heading 3';
    }
    .ql-snow .ql-picker.ql-header .ql-picker-label[data-value='4']::before,
    .ql-snow .ql-picker.ql-header .ql-picker-item[data-value='4']::before {
        content: 'Heading 4';
    }
    .ql-snow .ql-picker.ql-header .ql-picker-label[data-value='5']::before,
    .ql-snow .ql-picker.ql-header .ql-picker-item[data-value='5']::before {
        content: 'Heading 5';
    }
    .ql-snow .ql-picker.ql-header .ql-picker-label[data-value='6']::before,
    .ql-snow .ql-picker.ql-header .ql-picker-item[data-value='6']::before {
        content: 'Heading 6';
    }
    .ql-snow .ql-picker.ql-header .ql-picker-item[data-value='1']::before {
        font-size: 2rem;
    }
    .ql-snow .ql-picker.ql-header .ql-picker-item[data-value='2']::before {
        font-size: 1.5rem;
    }
    .ql-snow .ql-picker.ql-header .ql-picker-item[data-value='3']::before {
        font-size: 1.17rem;
    }
    .ql-snow .ql-picker.ql-header .ql-picker-item[data-value='4']::before {
        font-size: 1rem;
    }
    .ql-snow .ql-picker.ql-header .ql-picker-item[data-value='5']::before {
        font-size: 0.83rem;
    }
    .ql-snow .ql-picker.ql-header .ql-picker-item[data-value='6']::before {
        font-size: 0.67rem;
    }
    .ql-snow .ql-picker.ql-font {
        width: 108px;
    }
    .ql-snow .ql-picker.ql-font .ql-picker-label::before,
    .ql-snow .ql-picker.ql-font .ql-picker-item::before {
        content: 'Sans Serif';
    }
    .ql-snow .ql-picker.ql-font .ql-picker-label[data-value='serif']::before,
    .ql-snow .ql-picker.ql-font .ql-picker-item[data-value='serif']::before {
        content: 'Serif';
    }
    .ql-snow .ql-picker.ql-font .ql-picker-label[data-value='monospace']::before,
    .ql-snow .ql-picker.ql-font .ql-picker-item[data-value='monospace']::before {
        content: 'Monospace';
    }
    .ql-snow .ql-picker.ql-font .ql-picker-item[data-value='serif']::before {
        font-family:
            Georgia,
            Times New Roman,
            serif;
    }
    .ql-snow .ql-picker.ql-font .ql-picker-item[data-value='monospace']::before {
        font-family:
            Monaco,
            Courier New,
            monospace;
    }
    .ql-snow .ql-picker.ql-size {
        width: 98px;
    }
    .ql-snow .ql-picker.ql-size .ql-picker-label::before,
    .ql-snow .ql-picker.ql-size .ql-picker-item::before {
        content: 'Normal';
    }
    .ql-snow .ql-picker.ql-size .ql-picker-label[data-value='small']::before,
    .ql-snow .ql-picker.ql-size .ql-picker-item[data-value='small']::before {
        content: 'Small';
    }
    .ql-snow .ql-picker.ql-size .ql-picker-label[data-value='large']::before,
    .ql-snow .ql-picker.ql-size .ql-picker-item[data-value='large']::before {
        content: 'Large';
    }
    .ql-snow .ql-picker.ql-size .ql-picker-label[data-value='huge']::before,
    .ql-snow .ql-picker.ql-size .ql-picker-item[data-value='huge']::before {
        content: 'Huge';
    }
    .ql-snow .ql-picker.ql-size .ql-picker-item[data-value='small']::before {
        font-size: 10px;
    }
    .ql-snow .ql-picker.ql-size .ql-picker-item[data-value='large']::before {
        font-size: 18px;
    }
    .ql-snow .ql-picker.ql-size .ql-picker-item[data-value='huge']::before {
        font-size: 32px;
    }
    .ql-snow .ql-color-picker.ql-background .ql-picker-item {
        background: #fff;
    }
    .ql-snow .ql-color-picker.ql-color .ql-picker-item {
        background: #000;
    }
    .ql-toolbar.ql-snow {
        border: 1px solid #ccc;
        box-sizing: border-box;
        font-family: 'Helvetica Neue', 'Helvetica', 'Arial', sans-serif;
        padding: 8px;
    }
    .ql-toolbar.ql-snow .ql-formats {
        margin-inline-end: 15px;
    }
    .ql-toolbar.ql-snow .ql-picker-label {
        border: 1px solid transparent;
    }
    .ql-toolbar.ql-snow .ql-picker-options {
        border: 1px solid transparent;
        box-shadow: rgba(0, 0, 0, 0.2) 0 2px 8px;
    }
    .ql-toolbar.ql-snow .ql-picker.ql-expanded .ql-picker-label {
        border-color: #ccc;
    }
    .ql-toolbar.ql-snow .ql-picker.ql-expanded .ql-picker-options {
        border-color: #ccc;
    }
    .ql-toolbar.ql-snow .ql-color-picker .ql-picker-item.ql-selected,
    .ql-toolbar.ql-snow .ql-color-picker .ql-picker-item:hover {
        border-color: #000;
    }
    .ql-toolbar.ql-snow + .ql-container.ql-snow {
        border-block-start: 0;
    }
    .ql-snow .ql-tooltip {
        background: #fff;
        border: 1px solid #ccc;
        box-shadow: 0 0 5px #ddd;
        color: #444;
        padding: 5px 12px;
        white-space: nowrap;
    }
    .ql-snow .ql-tooltip::before {
        content: 'Visit URL:';
        line-height: 26px;
        margin-inline-end: 8px;
    }
    .ql-snow .ql-tooltip input[type='text'] {
        display: none;
        border: 1px solid #ccc;
        font-size: 13px;
        height: 26px;
        margin: 0;
        padding: 3px 5px;
        width: 170px;
    }
    .ql-snow .ql-tooltip a.ql-preview {
        display: inline-block;
        max-width: 200px;
        overflow-x: hidden;
        text-overflow: ellipsis;
        vertical-align: top;
    }
    .ql-snow .ql-tooltip a.ql-action::after {
        border-inline-end: 1px solid #ccc;
        content: 'Edit';
        margin-inline-start: 16px;
        padding-inline-end: 8px;
    }
    .ql-snow .ql-tooltip a.ql-remove::before {
        content: 'Remove';
        margin-inline-start: 8px;
    }
    .ql-snow .ql-tooltip a {
        line-height: 26px;
    }
    .ql-snow .ql-tooltip.ql-editing a.ql-preview,
    .ql-snow .ql-tooltip.ql-editing a.ql-remove {
        display: none;
    }
    .ql-snow .ql-tooltip.ql-editing input[type='text'] {
        display: inline-block;
    }
    .ql-snow .ql-tooltip.ql-editing a.ql-action::after {
        border-inline-end: 0;
        content: 'Save';
        padding-inline-end: 0;
    }
    .ql-snow .ql-tooltip[data-mode='link']::before {
        content: 'Enter link:';
    }
    .ql-snow .ql-tooltip[data-mode='formula']::before {
        content: 'Enter formula:';
    }
    .ql-snow .ql-tooltip[data-mode='video']::before {
        content: 'Enter video:';
    }
    .ql-snow a {
        color: #06c;
    }
    .ql-container.ql-snow {
        border: 1px solid #ccc;
    }

    .p-editor {
        display: block;
    }

    .p-editor .p-editor-toolbar {
        background: dt('editor.toolbar.background');
        border-start-end-radius: dt('editor.toolbar.border.radius');
        border-start-start-radius: dt('editor.toolbar.border.radius');
    }

    .p-editor .p-editor-toolbar.ql-snow {
        border: 1px solid dt('editor.toolbar.border.color');
    }

    .p-editor .p-editor-toolbar.ql-snow .ql-stroke {
        stroke: dt('editor.toolbar.item.color');
    }

    .p-editor .p-editor-toolbar.ql-snow .ql-fill {
        fill: dt('editor.toolbar.item.color');
    }

    .p-editor .p-editor-toolbar.ql-snow .ql-picker .ql-picker-label {
        border: 0 none;
        color: dt('editor.toolbar.item.color');
    }

    .p-editor .p-editor-toolbar.ql-snow .ql-picker .ql-picker-label:hover {
        color: dt('editor.toolbar.item.hover.color');
    }

    .p-editor .p-editor-toolbar.ql-snow .ql-picker .ql-picker-label:hover .ql-stroke {
        stroke: dt('editor.toolbar.item.hover.color');
    }

    .p-editor .p-editor-toolbar.ql-snow .ql-picker .ql-picker-label:hover .ql-fill {
        fill: dt('editor.toolbar.item.hover.color');
    }

    .p-editor .p-editor-toolbar.ql-snow .ql-picker.ql-expanded .ql-picker-label {
        color: dt('editor.toolbar.item.active.color');
    }

    .p-editor .p-editor-toolbar.ql-snow .ql-picker.ql-expanded .ql-picker-label .ql-stroke {
        stroke: dt('editor.toolbar.item.active.color');
    }

    .p-editor .p-editor-toolbar.ql-snow .ql-picker.ql-expanded .ql-picker-label .ql-fill {
        fill: dt('editor.toolbar.item.active.color');
    }

    .p-editor .p-editor-toolbar.ql-snow .ql-picker.ql-expanded .ql-picker-options {
        background: dt('editor.overlay.background');
        border: 1px solid dt('editor.overlay.border.color');
        box-shadow: dt('editor.overlay.shadow');
        border-radius: dt('editor.overlay.border.radius');
        padding: dt('editor.overlay.padding');
    }

    .p-editor .p-editor-toolbar.ql-snow .ql-picker.ql-expanded .ql-picker-options .ql-picker-item {
        color: dt('editor.overlay.option.color');
        border-radius: dt('editor.overlay.option.border.radius');
    }

    .p-editor .p-editor-toolbar.ql-snow .ql-picker.ql-expanded .ql-picker-options .ql-picker-item:hover {
        background: dt('editor.overlay.option.focus.background');
        color: dt('editor.overlay.option.focus.color');
    }

    .p-editor .p-editor-toolbar.ql-snow .ql-picker.ql-expanded:not(.ql-color-picker, .ql-icon-picker) .ql-picker-item {
        padding: dt('editor.overlay.option.padding');
    }

    .p-editor .p-editor-content {
        border-end-end-radius: dt('editor.content.border.radius');
        border-end-start-radius: dt('editor.content.border.radius');
    }

    .p-editor .p-editor-content.ql-snow {
        border: 1px solid dt('editor.content.border.color');
    }

    .p-editor .p-editor-content .ql-editor {
        background: dt('editor.content.background');
        color: dt('editor.content.color');
        border-end-end-radius: dt('editor.content.border.radius');
        border-end-start-radius: dt('editor.content.border.radius');
    }

    .p-editor .ql-snow.ql-toolbar button:hover,
    .p-editor .ql-snow.ql-toolbar button:focus {
        color: dt('editor.toolbar.item.hover.color');
    }

    .p-editor .ql-snow.ql-toolbar button:hover .ql-stroke,
    .p-editor .ql-snow.ql-toolbar button:focus .ql-stroke {
        stroke: dt('editor.toolbar.item.hover.color');
    }

    .p-editor .ql-snow.ql-toolbar button:hover .ql-fill,
    .p-editor .ql-snow.ql-toolbar button:focus .ql-fill {
        fill: dt('editor.toolbar.item.hover.color');
    }

    .p-editor .ql-snow.ql-toolbar button.ql-active,
    .p-editor .ql-snow.ql-toolbar .ql-picker-label.ql-active,
    .p-editor .ql-snow.ql-toolbar .ql-picker-item.ql-selected {
        color: dt('editor.toolbar.item.active.color');
    }

    .p-editor .ql-snow.ql-toolbar button.ql-active .ql-stroke,
    .p-editor .ql-snow.ql-toolbar .ql-picker-label.ql-active .ql-stroke,
    .p-editor .ql-snow.ql-toolbar .ql-picker-item.ql-selected .ql-stroke {
        stroke: dt('editor.toolbar.item.active.color');
    }

    .p-editor .ql-snow.ql-toolbar button.ql-active .ql-fill,
    .p-editor .ql-snow.ql-toolbar .ql-picker-label.ql-active .ql-fill,
    .p-editor .ql-snow.ql-toolbar .ql-picker-item.ql-selected .ql-fill {
        fill: dt('editor.toolbar.item.active.color');
    }

    .p-editor .ql-snow.ql-toolbar button.ql-active .ql-picker-label,
    .p-editor .ql-snow.ql-toolbar .ql-picker-label.ql-active .ql-picker-label,
    .p-editor .ql-snow.ql-toolbar .ql-picker-item.ql-selected .ql-picker-label {
        color: dt('editor.toolbar.item.active.color');
    }
`,classes:{root:function(e){return[`p-editor`,{"p-invalid":e.instance.$invalid}]},toolbar:`p-editor-toolbar`,content:`p-editor-content`}}),Ie={name:`BaseEditor`,extends:B,props:{placeholder:String,readonly:Boolean,formats:Array,editorStyle:null,modules:null},style:Fe,provide:function(){return{$pcEditor:this,$parentInstance:this}}};function Le(e){"@babel/helpers - typeof";return Le=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},Le(e)}function Re(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function ze(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?Re(Object(n),!0).forEach(function(t){Be(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):Re(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function Be(e,t,n){return(t=Ve(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function Ve(e){var t=He(e,`string`);return Le(t)==`symbol`?t:t+``}function He(e,t){if(Le(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(Le(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var Ue=function(){try{return window.Quill}catch{return null}}(),We={name:`Editor`,extends:Ie,inheritAttrs:!1,emits:[`text-change`,`selection-change`,`load`],quill:null,watch:{modelValue:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},d_value:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},readonly:function(){this.handleReadOnlyChange()}},mounted:function(){var e=this,t={modules:ze({toolbar:this.$refs.toolbarElement},this.modules),readOnly:this.readonly,theme:`snow`,formats:this.formats,placeholder:this.placeholder};Ue?(this.quill=new Ue(this.$refs.editorElement,t),this.initQuill(),this.handleLoad()):M(()=>import(`./quill-BLmY9xB4.js`).then(function(n){n&&E(e.$refs.editorElement)&&(n.default?e.quill=new n.default(e.$refs.editorElement,t):e.quill=new n(e.$refs.editorElement,t),e.initQuill())}),__vite__mapDeps([0,1])).then(function(){e.handleLoad()})},beforeUnmount:function(){this.quill=null},methods:{renderValue:function(e){if(this.quill)if(e){var t=this.quill.clipboard.convert({html:e});this.quill.setContents(t)}else this.quill.setText(``)},initQuill:function(){var e=this;this.renderValue(this.d_value),this.quill.on(`text-change`,function(t,n,r){if(r===`user`){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();i===`<p><br></p>`&&(i=``),e.writeValue(i),e.$emit(`text-change`,{htmlValue:i,textValue:a,delta:t,source:r,instance:e.quill})}}),this.quill.on(`selection-change`,function(t,n,r){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();e.$emit(`selection-change`,{htmlValue:i,textValue:a,range:t,oldRange:n,source:r,instance:e.quill})})},handleLoad:function(){this.quill&&this.quill.getModule(`toolbar`)&&this.$emit(`load`,{instance:this.quill})},handleReadOnlyChange:function(){this.quill&&this.quill.enable(!this.readonly)}}};function Ge(t,n,r,i,a,s){return g(),w(`div`,e({class:t.cx(`root`)},t.ptmi(`root`)),[C(`div`,e({ref:`toolbarElement`,class:t.cx(`toolbar`)},t.ptm(`toolbar`)),[o(t.$slots,`toolbar`,{},function(){return[C(`span`,e({class:`ql-formats`},t.ptm(`formats`)),[C(`select`,e({class:`ql-header`,defaultValue:`0`},t.ptm(`header`)),[C(`option`,e({value:`1`},t.ptm(`option`)),`Heading`,16),C(`option`,e({value:`2`},t.ptm(`option`)),`Subheading`,16),C(`option`,e({value:`0`},t.ptm(`option`)),`Normal`,16)],16),C(`select`,e({class:`ql-font`},t.ptm(`font`)),[C(`option`,x(T(t.ptm(`option`))),null,16),C(`option`,e({value:`serif`},t.ptm(`option`)),null,16),C(`option`,e({value:`monospace`},t.ptm(`option`)),null,16)],16)],16),C(`span`,e({class:`ql-formats`},t.ptm(`formats`)),[C(`button`,e({class:`ql-bold`,type:`button`},t.ptm(`bold`)),null,16),C(`button`,e({class:`ql-italic`,type:`button`},t.ptm(`italic`)),null,16),C(`button`,e({class:`ql-underline`,type:`button`},t.ptm(`underline`)),null,16)],16),C(`span`,e({class:`ql-formats`},t.ptm(`formats`)),[C(`select`,e({class:`ql-color`},t.ptm(`color`)),null,16),C(`select`,e({class:`ql-background`},t.ptm(`background`)),null,16)],16),C(`span`,e({class:`ql-formats`},t.ptm(`formats`)),[C(`button`,e({class:`ql-list`,value:`ordered`,type:`button`},t.ptm(`list`)),null,16),C(`button`,e({class:`ql-list`,value:`bullet`,type:`button`},t.ptm(`list`)),null,16),C(`select`,e({class:`ql-align`},t.ptm(`select`)),[C(`option`,e({defaultValue:``},t.ptm(`option`)),null,16),C(`option`,e({value:`center`},t.ptm(`option`)),null,16),C(`option`,e({value:`right`},t.ptm(`option`)),null,16),C(`option`,e({value:`justify`},t.ptm(`option`)),null,16)],16)],16),C(`span`,e({class:`ql-formats`},t.ptm(`formats`)),[C(`button`,e({class:`ql-link`,type:`button`},t.ptm(`link`)),null,16),C(`button`,e({class:`ql-image`,type:`button`},t.ptm(`image`)),null,16),C(`button`,e({class:`ql-code-block`,type:`button`},t.ptm(`codeBlock`)),null,16)],16),C(`span`,e({class:`ql-formats`},t.ptm(`formats`)),[C(`button`,e({class:`ql-clean`,type:`button`},t.ptm(`clean`)),null,16)],16)]})],16),C(`div`,e({ref:`editorElement`,class:t.cx(`content`),style:t.editorStyle},t.ptm(`content`)),null,16)],16)}We.render=Ge;var Ke={key:0,class:`text-gray-500 dark:text-gray-400 text-xs`},qe=[`innerHTML`],Je={key:2,class:`text-gray-800 dark:text-gray-100`},Ye={class:`flex items-center gap-1`},Xe={__name:`DataTableSection`,props:{items:Array,loading:Boolean,total:Number,columns:{type:Array,default:()=>[]},entity:String,orgId:String,onLoad:Function},emits:[`edit`,`delete`],setup(e){let n=e,{t:r}=I(),i=a(1),f=a(15),v=b(()=>(i.value-1)*f.value),y=b(()=>[...n.columns.map(e=>({type:`text`,width:`w-24`,headerWidth:`w-20`})),{type:`icons`,count:2,headerWidth:`w-16`}]);function x(e){i.value=e.page+1,f.value=e.rows,n.onLoad&&n.onLoad(i.value,f.value)}return t(()=>{n.onLoad&&n.onLoad(1,15)}),(t,n)=>{let i=s(`tooltip`);return e.loading?(g(),d(te,{key:0,columns:y.value,rows:8},null,8,[`columns`])):(g(),d(c(U),{key:1,value:e.items,lazy:``,totalRecords:e.total,first:v.value,rows:f.value,onPage:x,paginator:``,paginatorTemplate:`FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown`,rowsPerPageOptions:[10,15,25,50],size:`small`,class:`!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden`},{empty:l(()=>[o(t.$slots,`empty`)]),default:l(()=>[(g(!0),w(S,null,h(e.columns,e=>(g(),d(c(W),{key:e.field,field:e.field,header:e.header,sortable:``},{body:l(({data:t})=>[e.field.startsWith(`_`)?(g(),w(`span`,Ke,m(t[e.field]||`-`),1)):_(``,!0),e.html?(g(),w(`div`,{key:1,class:`editor-content`,innerHTML:t[e.field]},null,8,qe)):(g(),w(`span`,Je,m(t[e.field]||`-`),1))]),_:2},1032,[`field`,`header`]))),128)),p(c(W),{header:c(r)(`common.actions`),style:{width:`90px`},frozen:``,alignFrozen:`right`},{body:l(({data:e})=>[C(`div`,Ye,[u(p(c(O),{icon:`pi pi-pencil`,size:`small`,text:``,severity:`secondary`,onClick:n=>t.$emit(`edit`,e)},null,8,[`onClick`]),[[i,c(r)(`common.edit`),void 0,{left:!0}]]),u(p(c(O),{icon:`pi pi-trash`,size:`small`,text:``,severity:`danger`,onClick:n=>t.$emit(`delete`,e)},null,8,[`onClick`]),[[i,c(r)(`common.delete`),void 0,{left:!0}]])])]),_:1},8,[`header`])]),_:3},8,[`value`,`totalRecords`,`first`,`rows`]))}}},Ze={class:`space-y-4`},Qe={key:0,class:`bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded p-3 text-xs text-red-700 dark:text-red-300`},$e={__name:`DialogForm`,props:{visible:Boolean,title:String,saving:Boolean,errors:{type:Object,default:()=>({})},width:{type:String,default:`480px`}},emits:[`save`,`cancel`],setup(e){let t=e,{t:n}=I(),r=b(()=>t.width===`maximize`?`90vw`:t.width);return(t,i)=>(g(),d(c(N),{visible:e.visible,"onUpdate:visible":i[2]||=e=>t.$emit(`update:visible`,e),header:e.title,modal:``,style:y({width:r.value}),class:`p-fluid`,closable:!e.saving},{footer:l(()=>[p(c(O),{label:c(n)(`common.cancel`),size:`small`,outlined:``,severity:`secondary`,disabled:e.saving,onClick:i[0]||=e=>t.$emit(`cancel`)},null,8,[`label`,`disabled`]),p(c(O),{label:c(n)(`common.save`),icon:`pi pi-check`,size:`small`,loading:e.saving,onClick:i[1]||=e=>t.$emit(`save`)},null,8,[`label`,`loading`])]),default:l(()=>[C(`div`,Ze,[o(t.$slots,`default`),Object.keys(e.errors).length?(g(),w(`div`,Qe,[(g(!0),w(S,null,h(e.errors,(e,t)=>(g(),w(`p`,{key:t,class:`mb-1`},[C(`strong`,null,m(t)+`:`,1),v(` `+m(Array.isArray(e)?e.join(`, `):e),1)]))),128))])):_(``,!0)])]),_:3},8,[`visible`,`header`,`style`,`closable`]))}},et={class:`space-y-4`},tt={class:`flex items-center justify-between`},nt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},rt={class:`text-sm text-gray-500 dark:text-gray-400`},it={class:`flex flex-col items-center justify-center py-10 text-gray-400`},at={class:`text-sm font-medium`},ot=`/api/v1/tenant/job-management/responsibilities`,st={__name:`JobResponsibilitySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let n=e,r=t,{t:i}=I(),o=F(),s=a([]),u=a(!1),d=a(0),h=a(!1),v=a(!1),y=a(``),x=a(!1),T=a({}),E=a(!1),D=a(!1),k=a(``),A=a(null),j=a({main_task:``,activities:``,outputs:``,success_indicators:``}),M=b(()=>{let e=i(`job_management.responsibilities_title`);return v.value?`${e}`:`${i(`common.create`)} ${e}`}),N=b(()=>[{field:`main_task`,header:i(`job_management.main_task`),html:!0},{field:`activities`,header:i(`job_management.activities`),html:!0},{field:`outputs`,header:i(`job_management.outputs`),html:!0},{field:`success_indicators`,header:i(`job_management.success_indicators`),html:!0}]);async function R(e,t){u.value=!0;try{let r=await P.get(ot,{params:{page:e,per_page:t,organization_id:n.orgId}}),i=r.data?.data||[];s.value=i.map(e=>({...e,main_task:e.main_task,activities:e.activities,outputs:e.outputs,success_indicators:e.success_indicators})),d.value=r.data?.total||0}catch(e){o.add({severity:`error`,detail:e.response?.data?.error?.message||i(`message.failed_to_load`),life:4e3})}finally{u.value=!1}}function B(){v.value=!1,y.value=``,j.value={main_task:``,activities:``,outputs:``,success_indicators:``},T.value={},h.value=!0}function V(e){v.value=!0,y.value=e.id,j.value={main_task:e.main_task||``,activities:e.activities||``,outputs:e.outputs||``,success_indicators:e.success_indicators||``},T.value={},h.value=!0}async function H(){x.value=!0,T.value={};try{let e={nomenclature:n.orgName||``,full_code:n.orgCode||``,...j.value,organization_id:n.orgId};v.value?await P.put(`${ot}/${y.value}`,e):await P.post(ot,e),h.value=!1,r(`saved`),o.add({severity:`success`,detail:i(`message.saved`),life:2e3}),R(1,15)}catch(e){let t=L(e);Object.keys(t).length?T.value=t:o.add({severity:`error`,detail:e.response?.data?.error?.message||i(`message.operation_failed`),life:4e3})}finally{x.value=!1}}function U(e){A.value=e,k.value=``,E.value=!0}async function W(){if(A.value){D.value=!0,k.value=``;try{await P.delete(`${ot}/${A.value.id}`),E.value=!1,r(`saved`),o.add({severity:`success`,detail:i(`message.deleted`),life:2e3}),R(1,15)}catch(e){k.value=e.response?.data?.error?.message||i(`message.operation_failed`)}finally{D.value=!1}}}return(t,n)=>(g(),w(`div`,et,[C(`div`,tt,[C(`div`,null,[C(`h2`,nt,m(c(i)(`job_management.responsibilities_title`)),1),C(`p`,rt,m(c(i)(`job_management.responsibilities_description`)),1)]),p(c(O),{label:c(i)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>B()},null,8,[`label`])]),p(Xe,{items:s.value,loading:u.value,total:d.value,columns:N.value,entity:`responsibilities`,"org-id":e.orgId,"on-load":R,onEdit:V,onDelete:U},{empty:l(()=>[C(`div`,it,[n[9]||=C(`i`,{class:`pi pi-list-check text-3xl mb-2 opacity-50`},null,-1),C(`p`,at,m(c(i)(`job_management.empty_responsibilities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p($e,{visible:h.value,"onUpdate:visible":n[5]||=e=>h.value=e,title:M.value,saving:x.value,errors:T.value,width:`maximize`,onSave:H,onCancel:n[6]||=e=>h.value=!1},{default:l(()=>[h.value?(g(),w(S,{key:0},[p(z,{label:c(i)(`job_management.main_task`),errors:T.value?.main_task},{default:l(()=>[p(c(We),{modelValue:j.value.main_task,"onUpdate:modelValue":n[1]||=e=>j.value.main_task=e,editorStyle:`height:120px`,class:f({"p-invalid":T.value?.main_task})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(z,{label:c(i)(`job_management.activities`),errors:T.value?.activities},{default:l(()=>[p(c(We),{modelValue:j.value.activities,"onUpdate:modelValue":n[2]||=e=>j.value.activities=e,editorStyle:`height:120px`,class:f({"p-invalid":T.value?.activities})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(z,{label:c(i)(`job_management.outputs`),errors:T.value?.outputs},{default:l(()=>[p(c(We),{modelValue:j.value.outputs,"onUpdate:modelValue":n[3]||=e=>j.value.outputs=e,editorStyle:`height:120px`,class:f({"p-invalid":T.value?.outputs})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(z,{label:c(i)(`job_management.success_indicators`),errors:T.value?.success_indicators},{default:l(()=>[p(c(We),{modelValue:j.value.success_indicators,"onUpdate:modelValue":n[4]||=e=>j.value.success_indicators=e,editorStyle:`height:120px`,class:f({"p-invalid":T.value?.success_indicators})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])],64)):_(``,!0)]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(J,{visible:E.value,"onUpdate:visible":n[7]||=e=>E.value=e,loading:D.value,"error-msg":k.value,onConfirm:W,onCancel:n[8]||=e=>E.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},ct={class:`space-y-4`},lt={class:`flex items-center justify-between`},ut={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},dt={class:`text-sm text-gray-500 dark:text-gray-400`},ft={class:`flex flex-col items-center justify-center py-10 text-gray-400`},pt={class:`text-sm font-medium`},mt=`/api/v1/tenant/job-management/hr-authorities`,ht={__name:`JobHRAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let n=e,r=t,{t:i}=I(),o=F(),s=a([]),u=a(!1),d=a(0),h=a(!1),_=a(!1),v=a(``),y=a(!1),x=a({}),S=a(!1),T=a(!1),E=a(``),D=a(null),k=a({description:``}),A=b(()=>{let e=i(`job_management.hr_authorities`);return _.value?`${i(`common.edit`)} ${e}`:`${i(`common.create`)} ${e}`}),j=b(()=>[{field:`description`,header:i(`job_management.description`)}]);async function M(e,t){u.value=!0;try{let r=await P.get(mt,{params:{page:e,per_page:t,organization_id:n.orgId}});s.value=r.data?.data||[],d.value=r.data?.total||0}catch(e){o.add({severity:`error`,detail:e.response?.data?.error?.message||i(`message.failed_to_load`),life:4e3})}finally{u.value=!1}}function N(){_.value=!1,v.value=``,k.value={nomenclature:``,full_code:``,description:``},x.value={},h.value=!0}function R(e){_.value=!0,v.value=e.id,k.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},x.value={},h.value=!0}async function B(){y.value=!0,x.value={};try{let e={...k.value,nomenclature:n.orgName||``,full_code:n.orgCode||``,organization_id:n.orgId};_.value?await P.put(`${mt}/${v.value}`,e):await P.post(mt,e),h.value=!1,r(`saved`),o.add({severity:`success`,detail:i(`message.saved`),life:2e3}),M(1,15)}catch(e){let t=L(e);Object.keys(t).length?x.value=t:o.add({severity:`error`,detail:e.response?.data?.error?.message||i(`message.operation_failed`),life:4e3})}finally{y.value=!1}}function H(e){D.value=e,E.value=``,S.value=!0}async function U(){if(D.value){T.value=!0,E.value=``;try{await P.delete(`${mt}/${D.value.id}`),S.value=!1,r(`saved`),o.add({severity:`success`,detail:i(`message.deleted`),life:2e3}),M(1,15)}catch(e){E.value=e.response?.data?.error?.message||i(`message.operation_failed`)}finally{T.value=!1}}}return(t,n)=>(g(),w(`div`,ct,[C(`div`,lt,[C(`div`,null,[C(`h2`,ut,m(c(i)(`job_management.hr_authorities`)),1),C(`p`,dt,m(c(i)(`job_management.authority_description`)),1)]),p(c(O),{label:c(i)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>N()},null,8,[`label`])]),p(Xe,{items:s.value,loading:u.value,total:d.value,columns:j.value,entity:`hr-authorities`,"org-id":e.orgId,"on-load":M,onEdit:R,onDelete:H},{empty:l(()=>[C(`div`,ft,[n[6]||=C(`i`,{class:`pi pi-users text-3xl mb-2 opacity-50`},null,-1),C(`p`,pt,m(c(i)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p($e,{visible:h.value,"onUpdate:visible":n[2]||=e=>h.value=e,title:A.value,saving:y.value,errors:x.value,onSave:B,onCancel:n[3]||=e=>h.value=!1},{default:l(()=>[p(z,{label:c(i)(`job_management.description`),errors:x.value?.description},{default:l(()=>[p(c(V),{modelValue:k.value.description,"onUpdate:modelValue":n[1]||=e=>k.value.description=e,rows:`3`,class:f([`w-full`,{"p-invalid":x.value?.description}])},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(J,{visible:S.value,"onUpdate:visible":n[4]||=e=>S.value=e,loading:T.value,"error-msg":E.value,onConfirm:U,onCancel:n[5]||=e=>S.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},gt={class:`space-y-4`},_t={class:`flex items-center justify-between`},vt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},yt={class:`text-sm text-gray-500 dark:text-gray-400`},bt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},xt={class:`text-sm font-medium`},St=`/api/v1/tenant/job-management/operational-authorities`,Ct={__name:`JobOpAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let n=e,r=t,{t:i}=I(),o=F(),s=a([]),u=a(!1),d=a(0),h=a(!1),_=a(!1),v=a(``),y=a(!1),x=a({}),S=a(!1),T=a(!1),E=a(``),D=a(null),k=a({description:``}),A=b(()=>{let e=i(`job_management.op_authorities`);return _.value?`${i(`common.edit`)} ${e}`:`${i(`common.create`)} ${e}`}),j=b(()=>[{field:`description`,header:i(`job_management.description`)}]);async function M(e,t){u.value=!0;try{let r=await P.get(St,{params:{page:e,per_page:t,organization_id:n.orgId}});s.value=r.data?.data||[],d.value=r.data?.total||0}catch(e){o.add({severity:`error`,detail:e.response?.data?.error?.message||i(`message.failed_to_load`),life:4e3})}finally{u.value=!1}}function N(){_.value=!1,v.value=``,k.value={nomenclature:``,full_code:``,description:``},x.value={},h.value=!0}function R(e){_.value=!0,v.value=e.id,k.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},x.value={},h.value=!0}async function B(){y.value=!0,x.value={};try{let e={...k.value,nomenclature:n.orgName||``,full_code:n.orgCode||``,organization_id:n.orgId};_.value?await P.put(`${St}/${v.value}`,e):await P.post(St,e),h.value=!1,r(`saved`),o.add({severity:`success`,detail:i(`message.saved`),life:2e3}),M(1,15)}catch(e){let t=L(e);Object.keys(t).length?x.value=t:o.add({severity:`error`,detail:e.response?.data?.error?.message||i(`message.operation_failed`),life:4e3})}finally{y.value=!1}}function H(e){D.value=e,E.value=``,S.value=!0}async function U(){if(D.value){T.value=!0,E.value=``;try{await P.delete(`${St}/${D.value.id}`),S.value=!1,r(`saved`),o.add({severity:`success`,detail:i(`message.deleted`),life:2e3}),M(1,15)}catch(e){E.value=e.response?.data?.error?.message||i(`message.operation_failed`)}finally{T.value=!1}}}return(t,n)=>(g(),w(`div`,gt,[C(`div`,_t,[C(`div`,null,[C(`h2`,vt,m(c(i)(`job_management.op_authorities`)),1),C(`p`,yt,m(c(i)(`job_management.authority_description`)),1)]),p(c(O),{label:c(i)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>N()},null,8,[`label`])]),p(Xe,{items:s.value,loading:u.value,total:d.value,columns:j.value,entity:`operational-authorities`,"org-id":e.orgId,"on-load":M,onEdit:R,onDelete:H},{empty:l(()=>[C(`div`,bt,[n[6]||=C(`i`,{class:`pi pi-cog text-3xl mb-2 opacity-50`},null,-1),C(`p`,xt,m(c(i)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p($e,{visible:h.value,"onUpdate:visible":n[2]||=e=>h.value=e,title:A.value,saving:y.value,errors:x.value,onSave:B,onCancel:n[3]||=e=>h.value=!1},{default:l(()=>[p(z,{label:c(i)(`job_management.description`),errors:x.value?.description},{default:l(()=>[p(c(V),{modelValue:k.value.description,"onUpdate:modelValue":n[1]||=e=>k.value.description=e,class:f([`w-full`,{"p-invalid":x.value?.description}]),rows:`3`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(J,{visible:S.value,"onUpdate:visible":n[4]||=e=>S.value=e,loading:T.value,"error-msg":E.value,onConfirm:U,onCancel:n[5]||=e=>S.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},wt={class:`space-y-4`},Tt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Et={class:`text-sm text-gray-500 dark:text-gray-400`},Dt={class:`max-w-2xl`},Ot={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},kt={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},At={class:`flex justify-end gap-2 pt-2`},jt=`/api/v1/tenant/job-management/working-activities`,Mt={__name:`JobActivitySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:n}){let r=n,i=e,{t:o}=I(),s=F(),u=a(!1),h=a(!0),v=a(``),y=a({}),b=a(``),x=a(!1),S=a(!1),T=a(``),E=a({job_management_value_id:``}),D=a([]);async function k(){try{let e=await P.get(`/api/v1/tenant/job-management/values`,{params:{type:`activity`,per_page:100}});D.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id}))}catch{}}async function A(){if(!i.orgId){h.value=!1;return}try{let e=(await P.get(jt,{params:{organization_id:i.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];b.value=t.id,E.value.job_management_value_id=t.job_management_value_id||``}}catch{}}async function j(){v.value=``,y.value={},u.value=!0;try{let e={nomenclature:i.orgName||``,full_code:i.orgCode||``,job_management_value_id:E.value.job_management_value_id||null,organization_id:i.orgId};if(b.value)await P.put(`${jt}/${b.value}`,{job_management_value_id:E.value.job_management_value_id||``});else{let t=await P.post(jt,e);b.value=t.data?.data?.id||``}s.add({severity:`success`,summary:o(`message.success`),detail:o(`common.saved`),life:2e3}),r(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(y.value=t,v.value=Object.values(t).join(`, `)):v.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{u.value=!1}}async function M(){if(b.value){S.value=!0,T.value=``;try{await P.delete(`${jt}/${b.value}`),x.value=!1,b.value=``,E.value.job_management_value_id=``,r(`saved`),s.add({severity:`success`,summary:o(`message.success`),detail:o(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||o(`message.operation_failed`)}finally{S.value=!1}}}return t(async()=>{try{await Promise.all([k(),A()])}finally{h.value=!1}}),(e,t)=>(g(),w(`div`,wt,[C(`div`,null,[C(`h2`,Tt,m(c(o)(`job_management.activities`)),1),C(`p`,Et,m(c(o)(`job_management.activity_description`)),1)]),C(`div`,Dt,[h.value?(g(),d(X,{key:0,type:`detail`,count:1,rows:3,cols:`grid-cols-1`,padding:`p-5`})):(g(),w(`div`,Ot,[p(z,{label:c(o)(`job_values.types.activity`),errors:y.value?.job_management_value_id},{default:l(()=>[p(Y,{modelValue:E.value.job_management_value_id,"onUpdate:modelValue":t[0]||=e=>E.value.job_management_value_id=e,options:D.value,"option-label":`label`,"option-value":`value`,placeholder:c(o)(`common.select`),class:f({"p-invalid":y.value?.job_management_value_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),v.value?(g(),w(`div`,kt,m(v.value),1)):_(``,!0),C(`div`,At,[b.value?(g(),d(c(O),{key:0,label:c(o)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[1]||=e=>x.value=!0},null,8,[`label`])):_(``,!0),p(c(O),{label:b.value?c(o)(`common.update`):c(o)(`common.save`),icon:`pi pi-check`,size:`small`,loading:u.value,disabled:u.value,onClick:j},null,8,[`label`,`loading`,`disabled`])])]))]),p(J,{visible:x.value,"onUpdate:visible":t[2]||=e=>x.value=e,loading:S.value,"error-msg":T.value,onConfirm:M,onCancel:t[3]||=e=>x.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Nt={class:`space-y-4`},Pt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Ft={class:`text-sm text-gray-500 dark:text-gray-400`},It={class:`max-w-2xl`},Lt={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Rt={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},zt={class:`flex justify-end gap-2 pt-2`},Bt=`/api/v1/tenant/job-management/working-risks`,Vt={__name:`JobRiskSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:n}){let r=n,i=e,{t:o}=I(),s=F(),u=a(!1),h=a(!0),v=a(``),y=a({}),b=a(``),x=a(!1),S=a(!1),T=a(``),E=a({job_management_value_environment_id:``,job_management_value_hazard_id:``}),D=a([]),k=a([]);async function A(){try{let[e,t]=await Promise.all([P.get(`/api/v1/tenant/job-management/values`,{params:{type:`environment`,per_page:100}}),P.get(`/api/v1/tenant/job-management/values`,{params:{type:`risk`,per_page:100}})]);D.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),k.value=(t.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id}))}catch{}}async function j(){if(!i.orgId){h.value=!1;return}try{let e=(await P.get(Bt,{params:{organization_id:i.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];b.value=t.id,E.value.job_management_value_environment_id=t.job_management_value_environment_id||``,E.value.job_management_value_hazard_id=t.job_management_value_hazard_id||``}}catch{}}async function M(){v.value=``,y.value={},u.value=!0;try{let e={nomenclature:i.orgName||``,full_code:i.orgCode||``,job_management_value_environment_id:E.value.job_management_value_environment_id||null,job_management_value_hazard_id:E.value.job_management_value_hazard_id||null,organization_id:i.orgId};if(b.value)await P.put(`${Bt}/${b.value}`,{job_management_value_environment_id:E.value.job_management_value_environment_id||``,job_management_value_hazard_id:E.value.job_management_value_hazard_id||``});else{let t=await P.post(Bt,e);b.value=t.data?.data?.id||``}s.add({severity:`success`,summary:o(`message.success`),detail:o(`common.saved`),life:2e3}),r(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(y.value=t,v.value=Object.values(t).join(`, `)):v.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{u.value=!1}}async function N(){if(b.value){S.value=!0,T.value=``;try{await P.delete(`${Bt}/${b.value}`),x.value=!1,b.value=``,E.value.job_management_value_environment_id=``,E.value.job_management_value_hazard_id=``,r(`saved`),s.add({severity:`success`,summary:o(`message.success`),detail:o(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||o(`message.operation_failed`)}finally{S.value=!1}}}return t(async()=>{try{await Promise.all([A(),j()])}finally{h.value=!1}}),(e,t)=>(g(),w(`div`,Nt,[C(`div`,null,[C(`h2`,Pt,m(c(o)(`job_management.risks`)),1),C(`p`,Ft,m(c(o)(`job_management.risk_description`)),1)]),C(`div`,It,[h.value?(g(),d(X,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(g(),w(`div`,Lt,[p(z,{label:c(o)(`job_management.work_environment`),errors:y.value?.job_management_value_environment_id},{default:l(()=>[p(Y,{modelValue:E.value.job_management_value_environment_id,"onUpdate:modelValue":t[0]||=e=>E.value.job_management_value_environment_id=e,options:D.value,"option-label":`label`,"option-value":`value`,placeholder:c(o)(`common.select`),class:f({"p-invalid":y.value?.job_management_value_environment_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),p(z,{label:c(o)(`job_management.risk`),errors:y.value?.job_management_value_hazard_id},{default:l(()=>[p(Y,{modelValue:E.value.job_management_value_hazard_id,"onUpdate:modelValue":t[1]||=e=>E.value.job_management_value_hazard_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:c(o)(`common.select`),class:f({"p-invalid":y.value?.job_management_value_hazard_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),v.value?(g(),w(`div`,Rt,m(v.value),1)):_(``,!0),C(`div`,zt,[b.value?(g(),d(c(O),{key:0,label:c(o)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[2]||=e=>x.value=!0},null,8,[`label`])):_(``,!0),p(c(O),{label:b.value?c(o)(`common.update`):c(o)(`common.save`),icon:`pi pi-check`,size:`small`,loading:u.value,disabled:u.value,onClick:M},null,8,[`label`,`loading`,`disabled`])])]))]),p(J,{visible:x.value,"onUpdate:visible":t[3]||=e=>x.value=e,loading:S.value,"error-msg":T.value,onConfirm:N,onCancel:t[4]||=e=>x.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Ht={class:`space-y-4`},Ut={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Wt={class:`text-sm text-gray-500 dark:text-gray-400`},Gt={class:`max-w-2xl`},Kt={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},qt={class:`pt-1`},Jt={class:`flex items-center gap-2 mb-3`},Yt={class:`text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400`},Xt={class:`space-y-4`},Zt={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Qt={class:`flex justify-end gap-2 pt-2`},$t={class:`max-w-3xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5 space-y-4`},en={class:`flex items-center justify-between gap-2 flex-wrap`},tn={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},nn={class:`text-sm text-gray-500 dark:text-gray-400`},rn={key:0,class:`text-sm text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg px-3 py-2`},an={key:1,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},on={key:2,class:`overflow-x-auto border border-gray-200 dark:border-gray-700 rounded-lg`},sn={class:`w-full text-sm`},cn={class:`bg-gray-50 dark:bg-gray-700/40 text-left`},ln={class:`px-3 py-2 font-semibold text-gray-600 dark:text-gray-300`},un={class:`px-3 py-2 font-semibold text-gray-600 dark:text-gray-300`},dn={class:`px-3 py-2 align-top text-gray-500 dark:text-gray-400`},fn={class:`px-3 py-2`},pn={class:`px-3 py-2`},mn={class:`px-3 py-2 align-top`},hn={key:3,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},gn={key:4,class:`flex justify-end gap-2 pt-2`},$=`/api/v1/tenant/job-management/relationships`,_n={__name:`JobRelationshipSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},orgSummaryId:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:n}){let r=n,i=e,{t:o}=I(),s=F(),u=a(!1),v=a(!0),y=a(``),b=a({}),x=a(``),T=a(!1),E=a(!1),D=a(``),k=a({job_management_value_relationship_id:``,job_management_value_frequency_id:``}),A=a([]),j=a([]),M=a([]),N=a([]),R=a(!1),B=a(``);async function V(){try{let[e,t,n]=await Promise.all([P.get(`/api/v1/tenant/job-management/values`,{params:{type:`relationship`,per_page:100}}),P.get(`/api/v1/tenant/job-management/values`,{params:{type:`frequency`,per_page:100}}),i.orgSummaryId?P.get(`/api/v1/tenant/organizations`,{params:{summary_id:i.orgSummaryId,per_page:100}}):Promise.resolve({data:{data:[]}})]);A.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),j.value=(t.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),M.value=(n.data?.data||[]).filter(e=>e.id!==i.orgId).map(e=>({label:e.full_code?`${e.full_code} - ${e.nomenclature}`:e.nomenclature,value:e.id}))}catch{}}async function U(){if(!i.orgId){v.value=!1;return}try{let e=(await P.get($,{params:{organization_id:i.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];x.value=t.id,k.value.job_management_value_relationship_id=t.job_management_value_relationship_id||``,k.value.job_management_value_frequency_id=t.job_management_value_frequency_id||``,await ne()}}catch{}}async function W(){y.value=``,b.value={},u.value=!0;try{let e={nomenclature:i.orgName||``,full_code:i.orgCode||``,job_management_value_relationship_id:k.value.job_management_value_relationship_id||null,job_management_value_frequency_id:k.value.job_management_value_frequency_id||null,organization_id:i.orgId};if(x.value)await P.put(`${$}/${x.value}`,{job_management_value_relationship_id:k.value.job_management_value_relationship_id||``,job_management_value_frequency_id:k.value.job_management_value_frequency_id||``});else{let t=await P.post($,e);x.value=t.data?.data?.id||``}s.add({severity:`success`,summary:o(`message.success`),detail:o(`common.saved`),life:2e3}),r(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(b.value=t,y.value=Object.values(t).join(`, `)):y.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{u.value=!1}}async function G(){if(x.value){E.value=!0,D.value=``;try{await P.delete(`${$}/${x.value}`),T.value=!1,x.value=``,k.value.job_management_value_relationship_id=``,k.value.job_management_value_frequency_id=``,N.value=[],r(`saved`),s.add({severity:`success`,summary:o(`message.success`),detail:o(`message.deleted`),life:2e3})}catch(e){D.value=e?.response?.data?.error?.message||o(`message.operation_failed`)}finally{E.value=!1}}}let K=0;function q(){x.value&&N.value.push({_key:`new-${++K}`,id:``,organization_id:``,activity:``})}function ee(e){let t=N.value[e];t&&(t.id?te(t.id,e):N.value.splice(e,1))}async function te(e,t){try{await P.delete(`${$}/${x.value}/details/${e}`),N.value.splice(t,1),s.add({severity:`success`,summary:o(`message.success`),detail:o(`message.deleted`),life:2e3})}catch(e){s.add({severity:`error`,summary:o(`message.error`),detail:e?.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}}async function ne(){if(x.value)try{let e=await P.get(`${$}/${x.value}/details`);N.value=(e.data?.data||[]).map(e=>({_key:`db-${++K}`,id:e.id,organization_id:e.organization_id||``,activity:e.activity||``}))}catch{}}async function Z(){if(!(!x.value||R.value)){B.value=``,R.value=!0;try{for(let e of N.value){let t={organization_id:e.organization_id||``,activity:e.activity||``};e.id?await P.put(`${$}/${x.value}/details/${e.id}`,t):e.id=(await P.post(`${$}/${x.value}/details`,t)).data?.data?.id||``}await ne(),s.add({severity:`success`,summary:o(`message.success`),detail:o(`job_management.relationship_details_saved`),life:2e3})}catch(e){let t=L(e);Object.keys(t).length>0?B.value=Object.values(t).join(`, `):B.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{R.value=!1}}}return t(async()=>{try{await Promise.all([V(),U()])}finally{v.value=!1}}),(e,t)=>(g(),w(`div`,Ht,[C(`div`,null,[C(`h2`,Ut,m(c(o)(`job_management.relationships`)),1),C(`p`,Wt,m(c(o)(`job_management.relationship_description`)),1)]),C(`div`,Gt,[v.value?(g(),d(X,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(g(),w(`div`,Kt,[C(`div`,qt,[C(`div`,Jt,[t[5]||=C(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400`},[C(`i`,{class:`pi pi-compass text-sm`})],-1),C(`h3`,Yt,m(c(o)(`job_management.relationship_group_scope`)),1),t[6]||=C(`div`,{class:`flex-1 border-t border-gray-200 dark:border-gray-700`},null,-1)]),C(`div`,Xt,[p(z,{label:c(o)(`job_management.relationship_type`),errors:b.value?.job_management_value_relationship_id},{default:l(()=>[p(Y,{modelValue:k.value.job_management_value_relationship_id,"onUpdate:modelValue":t[0]||=e=>k.value.job_management_value_relationship_id=e,options:A.value,"option-label":`label`,"option-value":`value`,placeholder:c(o)(`common.select`),class:f({"p-invalid":b.value?.job_management_value_relationship_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),p(z,{label:c(o)(`job_management.frequency`),errors:b.value?.job_management_value_frequency_id},{default:l(()=>[p(Y,{modelValue:k.value.job_management_value_frequency_id,"onUpdate:modelValue":t[1]||=e=>k.value.job_management_value_frequency_id=e,options:j.value,"option-label":`label`,"option-value":`value`,placeholder:c(o)(`common.select`),class:f({"p-invalid":b.value?.job_management_value_frequency_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`])])]),y.value?(g(),w(`div`,Zt,m(y.value),1)):_(``,!0),C(`div`,Qt,[x.value?(g(),d(c(O),{key:0,label:c(o)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[2]||=e=>T.value=!0},null,8,[`label`])):_(``,!0),p(c(O),{label:x.value?c(o)(`common.update`):c(o)(`common.save`),icon:`pi pi-check`,size:`small`,loading:u.value,disabled:u.value,onClick:W},null,8,[`label`,`loading`,`disabled`])])]))]),C(`div`,$t,[C(`div`,en,[C(`div`,null,[C(`h3`,tn,m(c(o)(`job_management.relationship_details`)),1),C(`p`,nn,m(c(o)(`job_management.relationship_details_description`)),1)]),p(c(O),{label:c(o)(`job_management.add_relationship_detail`),icon:`pi pi-plus`,size:`small`,outlined:``,disabled:!x.value||R.value,onClick:q},null,8,[`label`,`disabled`])]),x.value?N.value.length===0?(g(),w(`div`,an,m(c(o)(`job_management.no_relationship_details`)),1)):_(``,!0):(g(),w(`div`,rn,m(c(o)(`job_management.save_relationship_first`)),1)),N.value.length>0?(g(),w(`div`,on,[C(`table`,sn,[C(`thead`,null,[C(`tr`,cn,[t[7]||=C(`th`,{class:`px-3 py-2 w-10 font-semibold text-gray-600 dark:text-gray-300`},`#`,-1),C(`th`,ln,m(c(o)(`job_management.relationship_organization`)),1),C(`th`,un,m(c(o)(`job_management.relationship_activity`)),1),t[8]||=C(`th`,{class:`px-3 py-2 w-12`},null,-1)])]),C(`tbody`,null,[(g(!0),w(S,null,h(N.value,(e,t)=>(g(),w(`tr`,{key:e._key,class:`border-t border-gray-200 dark:border-gray-700`},[C(`td`,dn,m(t+1),1),C(`td`,fn,[p(Y,{modelValue:e.organization_id,"onUpdate:modelValue":t=>e.organization_id=t,options:M.value,"option-label":`label`,"option-value":`value`,placeholder:c(o)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),C(`td`,pn,[p(H,{modelValue:e.activity,"onUpdate:modelValue":t=>e.activity=t,placeholder:c(o)(`job_management.relationship_activity`)},null,8,[`modelValue`,`onUpdate:modelValue`,`placeholder`])]),C(`td`,mn,[p(c(O),{icon:`pi pi-trash`,severity:`danger`,size:`small`,text:``,rounded:``,"aria-label":`Remove`,onClick:e=>ee(t)},null,8,[`onClick`])])]))),128))])])])):_(``,!0),B.value?(g(),w(`div`,hn,m(B.value),1)):_(``,!0),N.value.length>0?(g(),w(`div`,gn,[p(c(O),{label:c(o)(`job_management.save_relationship_details`),icon:`pi pi-save`,size:`small`,loading:R.value,disabled:R.value||!x.value,onClick:Z},null,8,[`label`,`loading`,`disabled`])])):_(``,!0)]),p(J,{visible:T.value,"onUpdate:visible":t[3]||=e=>T.value=e,loading:E.value,"error-msg":D.value,onConfirm:G,onCancel:t[4]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},vn={class:`space-y-4`},yn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},bn={class:`text-sm text-gray-500 dark:text-gray-400`},xn={class:`max-w-2xl`},Sn={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Cn={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},wn={class:`flex justify-end gap-2 pt-2`},Tn=`/api/v1/tenant/job-management/subordinate-controls`,En={__name:`JobSubordinateSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:n}){let r=n,i=e,{t:o}=I(),s=F(),u=a(!1),h=a(!0),v=a(``),y=a({}),b=a(``),x=a(!1),S=a(!1),T=a(``),E=a({job_management_value_id:``}),D=a([]);async function k(){try{let e=await P.get(`/api/v1/tenant/job-management/values`,{params:{type:`subordinate`,per_page:100}});D.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id}))}catch{}}async function A(){if(!i.orgId){h.value=!1;return}try{let e=(await P.get(Tn,{params:{organization_id:i.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];b.value=t.id,E.value.job_management_value_id=t.job_management_value_id||``}}catch{}}async function j(){v.value=``,y.value={},u.value=!0;try{let e={nomenclature:i.orgName||``,full_code:i.orgCode||``,job_management_value_id:E.value.job_management_value_id||null,organization_id:i.orgId};if(b.value)await P.put(`${Tn}/${b.value}`,{job_management_value_id:E.value.job_management_value_id||``});else{let t=await P.post(Tn,e);b.value=t.data?.data?.id||``}s.add({severity:`success`,summary:o(`message.success`),detail:o(`common.saved`),life:2e3}),r(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(y.value=t,v.value=Object.values(t).join(`, `)):v.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{u.value=!1}}async function M(){if(b.value){S.value=!0,T.value=``;try{await P.delete(`${Tn}/${b.value}`),x.value=!1,b.value=``,E.value.job_management_value_id=``,r(`saved`),s.add({severity:`success`,summary:o(`message.success`),detail:o(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||o(`message.operation_failed`)}finally{S.value=!1}}}return t(async()=>{try{await Promise.all([k(),A()])}finally{h.value=!1}}),(e,t)=>(g(),w(`div`,vn,[C(`div`,null,[C(`h2`,yn,m(c(o)(`job_management.subordinate_controls`)),1),C(`p`,bn,m(c(o)(`job_management.subordinate_description`)),1)]),C(`div`,xn,[h.value?(g(),d(X,{key:0,type:`detail`,count:1,rows:3,cols:`grid-cols-1`,padding:`p-5`})):(g(),w(`div`,Sn,[p(z,{label:c(o)(`job_management.control_type`),errors:y.value?.job_management_value_id},{default:l(()=>[p(Y,{modelValue:E.value.job_management_value_id,"onUpdate:modelValue":t[0]||=e=>E.value.job_management_value_id=e,options:D.value,"option-label":`label`,"option-value":`value`,placeholder:c(o)(`common.select`),class:f({"p-invalid":y.value?.job_management_value_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),v.value?(g(),w(`div`,Cn,m(v.value),1)):_(``,!0),C(`div`,wn,[b.value?(g(),d(c(O),{key:0,label:c(o)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[1]||=e=>x.value=!0},null,8,[`label`])):_(``,!0),p(c(O),{label:b.value?c(o)(`common.update`):c(o)(`common.save`),icon:`pi pi-check`,size:`small`,loading:u.value,disabled:u.value,onClick:j},null,8,[`label`,`loading`,`disabled`])])]))]),p(J,{visible:x.value,"onUpdate:visible":t[2]||=e=>x.value=e,loading:S.value,"error-msg":T.value,onConfirm:M,onCancel:t[3]||=e=>x.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Dn={class:`space-y-4`},On={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},kn={class:`text-sm text-gray-500 dark:text-gray-400`},An={class:`max-w-2xl`},jn={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Mn={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Nn={class:`flex justify-end gap-2 pt-2`},Pn=`/api/v1/tenant/job-management/assets`,Fn={__name:`JobAssetSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:n}){let r=n,i=e,{t:o}=I(),s=F(),u=a(!1),h=a(!0),v=a(``),y=a({}),b=a(``),x=a(!1),S=a(!1),T=a(``),E=a({job_management_value_asset_id:``,job_management_value_authority_id:``}),D=a([]),k=a([]);async function A(){try{let[e,t]=await Promise.all([P.get(`/api/v1/tenant/job-management/values`,{params:{type:`asset`,per_page:100}}),P.get(`/api/v1/tenant/job-management/values`,{params:{type:`asset_authority`,per_page:100}})]);D.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),k.value=(t.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id}))}catch{}}async function j(){if(!i.orgId){h.value=!1;return}try{let e=(await P.get(Pn,{params:{organization_id:i.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];b.value=t.id,E.value.job_management_value_asset_id=t.job_management_value_asset_id||``,E.value.job_management_value_authority_id=t.job_management_value_authority_id||``}}catch{}}async function M(){v.value=``,y.value={},u.value=!0;try{let e={nomenclature:i.orgName||``,full_code:i.orgCode||``,job_management_value_asset_id:E.value.job_management_value_asset_id||null,job_management_value_authority_id:E.value.job_management_value_authority_id||null,organization_id:i.orgId};if(b.value)await P.put(`${Pn}/${b.value}`,{job_management_value_asset_id:E.value.job_management_value_asset_id||``,job_management_value_authority_id:E.value.job_management_value_authority_id||``});else{let t=await P.post(Pn,e);b.value=t.data?.data?.id||``}s.add({severity:`success`,summary:o(`message.success`),detail:o(`common.saved`),life:2e3}),r(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(y.value=t,v.value=Object.values(t).join(`, `)):v.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{u.value=!1}}async function N(){if(b.value){S.value=!0,T.value=``;try{await P.delete(`${Pn}/${b.value}`),x.value=!1,b.value=``,E.value.job_management_value_asset_id=``,E.value.job_management_value_authority_id=``,r(`saved`),s.add({severity:`success`,summary:o(`message.success`),detail:o(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||o(`message.operation_failed`)}finally{S.value=!1}}}return t(async()=>{try{await Promise.all([A(),j()])}finally{h.value=!1}}),(e,t)=>(g(),w(`div`,Dn,[C(`div`,null,[C(`h2`,On,m(c(o)(`job_management.assets`)),1),C(`p`,kn,m(c(o)(`job_management.asset_description`)),1)]),C(`div`,An,[h.value?(g(),d(X,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(g(),w(`div`,jn,[p(z,{label:c(o)(`job_management.asset_type`),errors:y.value?.job_management_value_asset_id},{default:l(()=>[p(Y,{modelValue:E.value.job_management_value_asset_id,"onUpdate:modelValue":t[0]||=e=>E.value.job_management_value_asset_id=e,options:D.value,"option-label":`label`,"option-value":`value`,placeholder:c(o)(`common.select`),class:f({"p-invalid":y.value?.job_management_value_asset_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),p(z,{label:c(o)(`job_management.authority_level`),errors:y.value?.job_management_value_authority_id},{default:l(()=>[p(Y,{modelValue:E.value.job_management_value_authority_id,"onUpdate:modelValue":t[1]||=e=>E.value.job_management_value_authority_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:c(o)(`common.select`),class:f({"p-invalid":y.value?.job_management_value_authority_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),v.value?(g(),w(`div`,Mn,m(v.value),1)):_(``,!0),C(`div`,Nn,[b.value?(g(),d(c(O),{key:0,label:c(o)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[2]||=e=>x.value=!0},null,8,[`label`])):_(``,!0),p(c(O),{label:b.value?c(o)(`common.update`):c(o)(`common.save`),icon:`pi pi-check`,size:`small`,loading:u.value,disabled:u.value,onClick:M},null,8,[`label`,`loading`,`disabled`])])]))]),p(J,{visible:x.value,"onUpdate:visible":t[3]||=e=>x.value=e,loading:S.value,"error-msg":T.value,onConfirm:N,onCancel:t[4]||=e=>x.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},In={class:`space-y-4`},Ln={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Rn={class:`text-sm text-gray-500 dark:text-gray-400`},zn={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Bn={class:`flex items-center justify-between gap-4 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3`},Vn={class:`min-w-0`},Hn={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},Un={class:`text-xs text-gray-500 dark:text-gray-400 mt-0.5`},Wn={class:`space-y-4 pt-4 border-t border-gray-200 dark:border-gray-700`},Gn={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Kn={class:`flex justify-end gap-2 pt-2`},qn=`/api/v1/tenant/job-management/financials`,Jn={__name:`JobFinancialSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:n}){let r=n,o=e,{t:s}=I(),u=F(),h=a(!1),v=a(!0),y=a(``),x=a({}),S=a(``),T=a(!1),E=a(!1),D=a(``),k=a({is_authorized:!1,job_management_value_cash_id:``,job_management_value_authority_id:``,job_management_value_impact_id:``}),A=a([]),j=a([]),M=a([]),N=a([]),R=a([]),B=b(()=>k.value.is_authorized?j.value:M.value),V=b(()=>k.value.is_authorized?N.value:R.value);async function H(){try{let[e,t,n,r,i]=await Promise.all([P.get(`/api/v1/tenant/job-management/values`,{params:{type:`cash`,per_page:100}}),P.get(`/api/v1/tenant/job-management/values`,{params:{type:`authority`,per_page:100}}),P.get(`/api/v1/tenant/job-management/values`,{params:{type:`authority_unauthorized`,per_page:100}}),P.get(`/api/v1/tenant/job-management/values`,{params:{type:`impact`,per_page:100}}),P.get(`/api/v1/tenant/job-management/values`,{params:{type:`impact_unauthorized`,per_page:100}})]);A.value=(e.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),j.value=(t.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),M.value=(n.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),N.value=(r.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),R.value=(i.data?.data||[]).map(e=>({label:e.descriptions,value:e.id}))}catch{}}let U=!1;i(()=>k.value.is_authorized,(e,t)=>{U||e===t||(k.value.job_management_value_cash_id=``,k.value.job_management_value_authority_id=``,k.value.job_management_value_impact_id=``)},{flush:`sync`});async function W(){if(!o.orgId){v.value=!1;return}try{let e=(await P.get(qn,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];U=!0,S.value=t.id,k.value.is_authorized=!!t.is_authorized,k.value.job_management_value_cash_id=t.job_management_value_cash_id||``,k.value.job_management_value_authority_id=t.job_management_value_authority_id||``,k.value.job_management_value_impact_id=t.job_management_value_impact_id||``,U=!1}}catch{}}async function G(){y.value=``,x.value={},h.value=!0;try{let e=!!k.value.is_authorized,t={nomenclature:o.orgName||``,full_code:o.orgCode||``,is_authorized:e,job_management_value_cash_id:e&&k.value.job_management_value_cash_id||null,job_management_value_authority_id:k.value.job_management_value_authority_id||null,job_management_value_impact_id:k.value.job_management_value_impact_id||null,organization_id:o.orgId};if(S.value)await P.put(`${qn}/${S.value}`,{is_authorized:e,job_management_value_cash_id:e&&k.value.job_management_value_cash_id||``,job_management_value_authority_id:k.value.job_management_value_authority_id||``,job_management_value_impact_id:k.value.job_management_value_impact_id||``});else{let e=await P.post(qn,t);S.value=e.data?.data?.id||``}u.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),r(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(x.value=t,y.value=Object.values(t).join(`, `)):y.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{h.value=!1}}async function K(){if(S.value){E.value=!0,D.value=``;try{await P.delete(`${qn}/${S.value}`),T.value=!1,S.value=``,k.value.is_authorized=!1,k.value.job_management_value_cash_id=``,k.value.job_management_value_authority_id=``,k.value.job_management_value_impact_id=``,r(`saved`),u.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){D.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{E.value=!1}}}return t(async()=>{try{await Promise.all([H(),W()])}finally{v.value=!1}}),(e,t)=>(g(),w(`div`,In,[C(`div`,null,[C(`h2`,Ln,m(c(s)(`job_management.financials`)),1),C(`p`,Rn,m(c(s)(`job_management.financial_description`)),1)]),C(`div`,null,[v.value?(g(),d(X,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(g(),w(`div`,zn,[C(`div`,Bn,[C(`div`,Vn,[C(`p`,Hn,m(c(s)(`job_management.is_authorized`)),1),C(`p`,Un,m(c(s)(`job_management.is_authorized_description`)),1)]),p(c(ee),{modelValue:k.value.is_authorized,"onUpdate:modelValue":t[0]||=e=>k.value.is_authorized=e},null,8,[`modelValue`])]),C(`div`,Wn,[k.value.is_authorized?(g(),d(z,{key:0,label:c(s)(`job_management.cash_level`),errors:x.value?.job_management_value_cash_id},{default:l(()=>[p(Y,{modelValue:k.value.job_management_value_cash_id,"onUpdate:modelValue":t[1]||=e=>k.value.job_management_value_cash_id=e,options:A.value,"option-label":`label`,"option-value":`value`,placeholder:c(s)(`common.select`),class:f({"p-invalid":x.value?.job_management_value_cash_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`])):_(``,!0),p(z,{label:c(s)(`job_management.authority_level`),errors:x.value?.job_management_value_authority_id},{default:l(()=>[p(Y,{modelValue:k.value.job_management_value_authority_id,"onUpdate:modelValue":t[2]||=e=>k.value.job_management_value_authority_id=e,options:B.value,"option-label":`label`,"option-value":`value`,placeholder:c(s)(`common.select`),class:f({"p-invalid":x.value?.job_management_value_authority_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),p(z,{label:c(s)(`job_management.impact_level`),errors:x.value?.job_management_value_impact_id},{default:l(()=>[p(Y,{modelValue:k.value.job_management_value_impact_id,"onUpdate:modelValue":t[3]||=e=>k.value.job_management_value_impact_id=e,options:V.value,"option-label":`label`,"option-value":`value`,placeholder:c(s)(`common.select`),class:f({"p-invalid":x.value?.job_management_value_impact_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`])]),y.value?(g(),w(`div`,Gn,m(y.value),1)):_(``,!0),C(`div`,Kn,[S.value?(g(),d(c(O),{key:0,label:c(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[4]||=e=>T.value=!0},null,8,[`label`])):_(``,!0),p(c(O),{label:S.value?c(s)(`common.update`):c(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:h.value,disabled:h.value,onClick:G},null,8,[`label`,`loading`,`disabled`])])]))]),p(J,{visible:T.value,"onUpdate:visible":t[5]||=e=>T.value=e,loading:E.value,"error-msg":D.value,onConfirm:K,onCancel:t[6]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Yn=`/api/v1/tenant/job-management/potency-competencies`;function Xn({orgId:e,rows:t,afterDelete:n,onSaved:r,matchBy:i=`value`,descriptionField:o=`descriptions`}){let{t:s}=I(),c=F(),l=a(!1),u=a(``),d=a(!1),f=a(!1),p=a(``),m=a(null),h=a([]);function g(e){let t=(e.levelOptions||[]).find(t=>t.value===e.job_management_value_id);return t&&t[o]||``}function _(e){if(i===`competency`)return e.competency_id&&h.value.find(t=>t.competency_id&&t.competency_id===e.competency_id)||null;let t=new Set((e.levelOptions||[]).map(e=>e.value));return h.value.find(e=>e.job_management_value_id&&t.has(e.job_management_value_id))||null}function v(){t.value.forEach(e=>{let t=_(e);e.recordId=t?t.id:``,e.job_management_value_id=t&&t.job_management_value_id||``,e.weight!==void 0&&(e.weight=t?t.weight??e.weight:e.weight)})}async function y(){if(!e.value){h.value=[];return}try{let t=await P.get(Yn,{params:{organization_id:e.value,per_page:100}});h.value=t.data?.data||[]}catch{h.value=[]}}function b(e){m.value=e,p.value=``,d.value=!0}async function x(){let e=m.value;if(e){f.value=!0,p.value=``;try{e.recordId&&await P.delete(`${Yn}/${e.recordId}`),n&&n(e),d.value=!1,await y(),v(),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3}),r&&r()}catch(e){p.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{f.value=!1,m.value=null}}}async function S(){u.value=``,l.value=!0;try{for(let n of t.value)if(n.job_management_value_id){let t=n.competency_id?{competency_id:n.competency_id,job_management_value_id:n.job_management_value_id}:{job_management_value_id:n.job_management_value_id};n.weight!==void 0&&n.weight!==null&&n.weight!==``&&(t.weight=n.weight),n.recordId?await P.put(`${Yn}/${n.recordId}`,t):n.recordId=(await P.post(Yn,{organization_id:e.value,...t})).data?.data?.id||``}else n.recordId&&=(await P.delete(`${Yn}/${n.recordId}`),``);await y(),v(),c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),r&&r()}catch(e){let t=L(e);Object.keys(t).length>0?u.value=Object.values(t).join(`, `):u.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}return{savingCard:l,errorMsg:u,deleteVisible:d,deleting:f,deleteError:p,deleteTarget:m,records:h,levelDescription:g,hydrateRows:v,loadData:y,askDeleteRow:b,handleDelete:x,handleSave:S}}var Zn={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Qn={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},$n={class:`text-sm text-gray-500 dark:text-gray-400`},er={key:0,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},tr={key:1,class:`overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700`},nr={class:`w-full text-sm`},rr={class:`bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400`},ir={class:`px-4 py-3 font-semibold min-w-[220px]`},ar={class:`px-4 py-3 font-semibold min-w-[260px]`},or={class:`px-4 py-3 font-semibold min-w-[260px]`},sr={class:`px-4 py-3 font-semibold w-16 text-right`},cr={class:`px-4 py-3`},lr={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},ur={class:`px-4 py-3`},dr={class:`px-4 py-3 text-sm text-gray-600 dark:text-gray-300`},fr={class:`px-4 py-3 text-right`},pr={key:2,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},mr={key:3,class:`flex justify-end gap-2 pt-1`},hr={__name:`SelectablePotencyCard`,props:{orgId:String,typeGroup:{type:String,required:!0},skeletonRows:{type:Number,default:5},titleKey:{type:String,required:!0},descriptionKey:{type:String,required:!0},emptyKey:{type:String,required:!0},saveLabelKey:{type:String,required:!0},deleteTitleKey:{type:String,required:!0},deleteMessageKey:{type:String,required:!0}},emits:[`saved`],setup(e,{emit:n}){let o=n,s=e,{t:l}=I(),u=a(!0),f=a([]),v=a([]),y=a([]),{savingCard:x,errorMsg:T,deleteVisible:E,deleting:D,deleteError:k,deleteTarget:A,records:j,levelDescription:M,hydrateRows:N,loadData:F,askDeleteRow:L,handleDelete:R,handleSave:z}=Xn({orgId:b(()=>s.orgId),rows:v,afterDelete:e=>{let t=Array.isArray(y.value)?y.value:[];y.value=t.filter(t=>t!==e.type)},onSaved:()=>o(`saved`)}),B=b(()=>(f.value||[]).find(e=>e.type_group===s.typeGroup));function V(e){let t=`job_values.types.${e.type}`,n=l(t);return n===t?e.description_group||e.type:n}let H=b(()=>(B.value?.types||[]).map(e=>({label:V(e),value:e.type})));function U(){let e={};(B.value?.types||[]).forEach(t=>{e[t.type]=t});let t=Array.isArray(y.value)?y.value:y.value?[y.value]:[];v.value=t.filter(t=>e[t]).map(t=>{let n=e[t];return{competency_id:``,competency_name:V(n),competency_definition:``,type:n.type,levelOptions:(n.options||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``})),recordId:``,job_management_value_id:``}})}async function W(){try{let e=await P.get(`/api/v1/tenant/job-management/values/tree`);f.value=e.data?.data||[],U()}catch{f.value=[],v.value=[]}}function G(){let e={};(B.value?.types||[]).forEach(t=>{(t.options||[]).forEach(n=>{e[n.id]=t.type})});let t=[];j.value.forEach(n=>{let r=n.job_management_value_id&&e[n.job_management_value_id];r&&!t.includes(r)&&t.push(r)}),y.value=t,U(),N()}return i(y,()=>{U(),N()}),t(async()=>{try{await Promise.all([W(),F()])}finally{G(),u.value=!1}}),(t,n)=>(g(),w(`div`,Zn,[C(`div`,null,[C(`h3`,Qn,m(c(l)(e.titleKey)),1),C(`p`,$n,m(c(l)(e.descriptionKey)),1)]),u.value?(g(),d(X,{key:0,type:`detail`,count:1,rows:e.skeletonRows,cols:`grid-cols-1`,padding:`p-5`},null,8,[`rows`])):(g(),w(S,{key:1},[p(Y,{modelValue:y.value,"onUpdate:modelValue":n[0]||=e=>y.value=e,options:H.value,"option-label":`label`,"option-value":`value`,placeholder:c(l)(`common.select`),showClear:``,multiple:``},null,8,[`modelValue`,`options`,`placeholder`]),v.value.length===0?(g(),w(`div`,er,m(c(l)(e.emptyKey)),1)):(g(),w(`div`,tr,[C(`table`,nr,[C(`thead`,null,[C(`tr`,rr,[C(`th`,ir,m(c(l)(`job_management.potency_table_name`)),1),C(`th`,ar,m(c(l)(`job_management.potency_table_level`)),1),C(`th`,or,m(c(l)(`job_management.potency_table_description`)),1),C(`th`,sr,m(c(l)(`common.actions`)),1)])]),C(`tbody`,null,[(g(!0),w(S,null,h(v.value,e=>(g(),w(`tr`,{key:e.type,class:`border-t border-gray-100 dark:border-gray-700 align-top`},[C(`td`,cr,[C(`div`,lr,m(e.competency_name),1)]),C(`td`,ur,[p(Y,{modelValue:e.job_management_value_id,"onUpdate:modelValue":t=>e.job_management_value_id=t,options:e.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:c(l)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),C(`td`,dr,m(c(M)(e)),1),C(`td`,fr,[p(c(O),{icon:`pi pi-trash`,severity:`danger`,text:``,rounded:``,size:`small`,disabled:c(x),"aria-label":c(l)(`common.delete`),onClick:t=>c(L)(e)},null,8,[`disabled`,`aria-label`,`onClick`])])]))),128))])])])),c(T)?(g(),w(`div`,pr,m(c(T)),1)):_(``,!0),v.value.length>0?(g(),w(`div`,mr,[p(c(O),{label:c(l)(e.saveLabelKey),icon:`pi pi-check`,size:`small`,loading:c(x),disabled:c(x)||!e.orgId,onClick:c(z)},null,8,[`label`,`loading`,`disabled`,`onClick`])])):_(``,!0)],64)),p(J,{visible:c(E),"onUpdate:visible":n[1]||=e=>r(E)?E.value=e:null,title:c(l)(e.deleteTitleKey),message:c(l)(e.deleteMessageKey,{name:c(A)?.competency_name||``}),loading:c(D),"error-msg":c(k),onConfirm:c(R),onCancel:n[2]||=e=>E.value=!1},null,8,[`visible`,`title`,`message`,`loading`,`error-msg`,`onConfirm`])]))}},gr={__name:`PsychologicalPotencyCard`,props:{orgId:String},emits:[`saved`],setup(e,{emit:t}){let n=t;return(t,r)=>(g(),d(hr,{"org-id":e.orgId,"type-group":`psychological`,"skeleton-rows":5,"title-key":`job_management.potency_required_title`,"description-key":`job_management.potency_required_description`,"empty-key":`job_management.potency_required_empty`,"save-label-key":`job_management.save_potency_levels`,"delete-title-key":`job_management.potency_confirm_delete_title`,"delete-message-key":`job_management.potency_confirm_delete`,onSaved:r[0]||=e=>n(`saved`)},null,8,[`org-id`]))}},_r={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},vr={class:`flex items-start justify-between gap-4`},yr={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},br={class:`text-sm text-gray-500 dark:text-gray-400`},xr={class:`flex flex-col items-end gap-1 shrink-0`},Sr={class:`flex items-center gap-2`},Cr={class:`text-xs font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap`},wr={class:`w-24 shrink-0`},Tr={key:0,class:`text-xs text-red-500 dark:text-red-400`},Er={key:0,class:`text-sm text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg px-3 py-2`},Dr={key:1,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},Or={key:2,class:`overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700`},kr={class:`w-full text-sm`},Ar={class:`bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400`},jr={class:`px-4 py-3 font-semibold min-w-[220px]`},Mr={class:`px-4 py-3 font-semibold min-w-[260px]`},Nr={class:`px-4 py-3 font-semibold min-w-[130px]`},Pr={class:`px-4 py-3 font-semibold min-w-[260px]`},Fr={class:`px-4 py-3 font-semibold w-16 text-right`},Ir={class:`px-4 py-3`},Lr={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},Rr={key:0,class:`text-xs text-gray-400 dark:text-gray-500`},zr={class:`px-4 py-3`},Br={class:`px-4 py-3`},Vr={class:`px-4 py-3 text-sm text-gray-600 dark:text-gray-300`},Hr={class:`px-4 py-3 text-right`},Ur={key:3,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Wr={key:4,class:`flex justify-end gap-2 pt-1`},Gr={__name:`TechnicalPotencyCard`,props:{orgId:String},emits:[`saved`,`weight-saved`],setup(e,{emit:n}){let o=n,s=e,{t:l}=I(),u=F(),f=a(!0),v=a([]),y=a([]),x=a([]),T=a([]),E=a([]),D=b(()=>E.value.length>0),k=a(``),A=a(``),j=a(``),M=a(!1),N=a(``),{savingCard:L,errorMsg:R,deleteVisible:z,deleting:B,deleteError:V,deleteTarget:H,records:U,levelDescription:W,hydrateRows:G,loadData:q,askDeleteRow:ee,handleDelete:te,handleSave:ne}=Xn({orgId:b(()=>s.orgId),rows:x,matchBy:`competency`,descriptionField:`note`,afterDelete:e=>{let t=Array.isArray(T.value)?T.value:[];T.value=t.filter(t=>t!==e.competency_id)},onSaved:()=>o(`saved`)}),Z=b(()=>(v.value||[]).map(e=>({id:e.id,name:e.name,cluster:e.cluster||``}))),re=b(()=>{let e={};return(Z.value||[]).forEach(t=>{(e[t.cluster]=e[t.cluster]||[]).push(t)}),Object.keys(e).sort().map(t=>({label:t,items:e[t].sort((e,t)=>e.name.localeCompare(t.name))}))}),ie=b(()=>(y.value||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``,note:e.note||``})));function ae(){let e={};(Z.value||[]).forEach(t=>{e[t.id]=t});let t=(Array.isArray(T.value)?T.value:T.value?[T.value]:[]).filter(t=>e[t]),n=t.length>0?Math.round(100/t.length*100)/100:0;x.value=t.map(t=>{let r=e[t];return{competency_id:t,competency_name:r.name,cluster:r.cluster,levelOptions:ie.value,recordId:``,job_management_value_id:``,weight:n}})}async function oe(){try{let[e,t]=await Promise.all([P.get(`/api/v1/tenant/settings/competencies`,{params:{per_page:500}}),P.get(`/api/v1/tenant/job-management/values/clusters/technical`)]);E.value=t.data?.data?.clusters||[];let n=new Set(E.value),r=e.data?.data||[];v.value=r.filter(e=>e.cluster&&n.has(e.cluster))}catch{v.value=[]}}async function se(){try{let e=await P.get(`/api/v1/tenant/job-management/values`,{params:{type:`technical`,per_page:100}});y.value=e.data?.data||[]}catch{y.value=[]}}async function Q(){if(!s.orgId){k.value=``,A.value=``;return}try{let e=((await P.get(`/api/v1/tenant/job-management/competency-groups`,{params:{organization_id:s.orgId}})).data?.data||[]).find(e=>e.category===`technical`);A.value=e?e.id:``,k.value=e?e.weight:``,j.value=k.value}catch{k.value=``,A.value=``,j.value=``}}async function ce(){if(k.value===``||k.value===null||k.value===void 0){N.value=l(`job_management.potency_technical_weight_required`);return}if(!(j.value!==``&&k.value===j.value)){M.value=!0,N.value=``;try{let e=A.value;if(!e)try{let t=((await P.get(`/api/v1/tenant/job-management/competency-groups`,{params:{organization_id:s.orgId}})).data?.data||[]).find(e=>e.category===`technical`);t&&(e=t.id)}catch{}let t={weight:k.value};e?await P.put(`/api/v1/tenant/job-management/competency-groups/${e}`,t):await P.post(`/api/v1/tenant/job-management/competency-groups`,{organization_id:s.orgId,category:`technical`,weight:k.value}),j.value=k.value,u.add({severity:`success`,summary:l(`message.success`),detail:l(`job_management.potency_technical_weight_saved`),life:2e3}),o(`saved`),o(`weight-saved`,k.value),await Q()}catch(e){N.value=e?.response?.data?.error?.message||e.message||l(`message.operation_failed`)}finally{M.value=!1}}}function le(){let e={};(Z.value||[]).forEach(t=>{e[t.id]=t});let t=[];U.value.forEach(n=>{n.competency_id&&e[n.competency_id]&&!t.includes(n.competency_id)&&t.push(n.competency_id)}),T.value=t,ae(),G()}return i(T,()=>{ae(),G()}),t(async()=>{try{await Promise.all([oe(),se(),q(),Q()])}finally{le(),f.value=!1}}),(t,n)=>(g(),w(`div`,_r,[C(`div`,vr,[C(`div`,null,[C(`h3`,yr,m(c(l)(`job_management.potency_technical_title`)),1),C(`p`,br,m(c(l)(`job_management.potency_technical_description`)),1)]),C(`div`,xr,[C(`div`,Sr,[C(`label`,Cr,m(c(l)(`job_management.potency_technical_weight_label`)),1),C(`div`,wr,[p(c(K),{modelValue:k.value,"onUpdate:modelValue":n[0]||=e=>k.value=e,fluid:``,min:0,max:100,suffix:`%`,size:`small`,disabled:M.value||!e.orgId,onBlur:ce},null,8,[`modelValue`,`disabled`])])]),N.value?(g(),w(`div`,Tr,m(N.value),1)):_(``,!0)])]),f.value?(g(),d(X,{key:0,type:`detail`,count:1,rows:8,cols:`grid-cols-1`,padding:`p-5`})):(g(),w(S,{key:1},[p(Y,{modelValue:T.value,"onUpdate:modelValue":n[1]||=e=>T.value=e,options:re.value,"option-label":`name`,"option-value":`id`,"option-group-label":`label`,"option-group-children":`items`,placeholder:c(l)(`job_management.potency_technical_placeholder`),showClear:``,multiple:``},null,8,[`modelValue`,`options`,`placeholder`]),D.value?x.value.length===0?(g(),w(`div`,Dr,m(c(l)(`job_management.potency_technical_empty`)),1)):(g(),w(`div`,Or,[C(`table`,kr,[C(`thead`,null,[C(`tr`,Ar,[C(`th`,jr,m(c(l)(`job_management.potency_table_name`)),1),C(`th`,Mr,m(c(l)(`job_management.potency_table_level`)),1),C(`th`,Nr,m(c(l)(`job_management.potency_table_weight`)),1),C(`th`,Pr,m(c(l)(`job_management.potency_table_description`)),1),C(`th`,Fr,m(c(l)(`common.actions`)),1)])]),C(`tbody`,null,[(g(!0),w(S,null,h(x.value,e=>(g(),w(`tr`,{key:e.competency_id,class:`border-t border-gray-100 dark:border-gray-700 align-top`},[C(`td`,Ir,[C(`div`,Lr,m(e.competency_name),1),e.cluster?(g(),w(`div`,Rr,m(e.cluster),1)):_(``,!0)]),C(`td`,zr,[p(Y,{modelValue:e.job_management_value_id,"onUpdate:modelValue":t=>e.job_management_value_id=t,options:e.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:c(l)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),C(`td`,Br,[p(c(K),{modelValue:e.weight,"onUpdate:modelValue":t=>e.weight=t,class:`!w-full`,min:0,max:100,suffix:`%`,size:`small`},null,8,[`modelValue`,`onUpdate:modelValue`])]),C(`td`,Vr,m(c(W)(e)),1),C(`td`,Hr,[p(c(O),{icon:`pi pi-trash`,severity:`danger`,text:``,rounded:``,size:`small`,disabled:c(L),"aria-label":c(l)(`common.delete`),onClick:t=>c(ee)(e)},null,8,[`disabled`,`aria-label`,`onClick`])])]))),128))])])])):(g(),w(`div`,Er,m(c(l)(`job_management.potency_technical_no_mapping`)),1)),c(R)?(g(),w(`div`,Ur,m(c(R)),1)):_(``,!0),x.value.length>0?(g(),w(`div`,Wr,[p(c(O),{label:c(l)(`job_management.save_technical`),icon:`pi pi-check`,size:`small`,loading:c(L),disabled:c(L)||!e.orgId,onClick:c(ne)},null,8,[`label`,`loading`,`disabled`,`onClick`])])):_(``,!0)],64)),p(J,{visible:c(z),"onUpdate:visible":n[2]||=e=>r(z)?z.value=e:null,title:c(l)(`job_management.potency_confirm_delete_title`),message:c(l)(`job_management.potency_confirm_delete`,{name:c(H)?.competency_name||``}),loading:c(B),"error-msg":c(V),onConfirm:c(te),onCancel:n[3]||=e=>z.value=!1},null,8,[`visible`,`title`,`message`,`loading`,`error-msg`,`onConfirm`])]))}},Kr={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},qr={class:`flex items-start justify-between gap-4`},Jr={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},Yr={class:`text-sm text-gray-500 dark:text-gray-400`},Xr={class:`flex flex-col items-end gap-1 shrink-0`},Zr={class:`flex items-center gap-2`},Qr={class:`text-xs font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap`},$r={class:`w-24 shrink-0 text-right`},ei={key:0,class:`text-sm font-semibold text-gray-700 dark:text-gray-200`},ti={key:1,class:`text-sm text-gray-400 dark:text-gray-500`},ni={key:0,class:`pi pi-spin pi-spinner text-sm text-gray-400`},ri={key:0,class:`text-xs text-red-500 dark:text-red-400`},ii={key:0,class:`text-sm text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg px-3 py-2`},ai={key:1,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},oi={key:2,class:`overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700`},si={class:`w-full text-sm`},ci={class:`bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400`},li={class:`px-4 py-3 font-semibold min-w-[220px]`},ui={class:`px-4 py-3 font-semibold min-w-[260px]`},di={class:`px-4 py-3 font-semibold min-w-[130px]`},fi={class:`px-4 py-3 font-semibold min-w-[260px]`},pi={class:`px-4 py-3 font-semibold w-16 text-right`},mi={class:`px-4 py-3`},hi={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},gi={key:0,class:`text-xs text-gray-400 dark:text-gray-500`},_i={class:`px-4 py-3`},vi={class:`px-4 py-3`},yi={class:`px-4 py-3 text-sm text-gray-600 dark:text-gray-300`},bi={class:`px-4 py-3 text-right`},xi={key:3,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Si={key:4,class:`flex justify-end gap-2 pt-1`},Ci={__name:`ManagerialPotencyCard`,props:{orgId:String,technicalWeight:{type:Number,default:null}},emits:[`saved`],setup(e,{emit:n}){let o=n,s=e,{t:l}=I(),u=F(),f=a(!0),v=a([]),y=a([]),x=a([]),T=a([]),E=a([]),D=b(()=>E.value.length>0),k=a(``),A=a(``),j=a(``),M=a(!1),N=a(``),L=b(()=>{let e=k.value;return e===``||e==null?null:Math.round((100-e)*100)/100}),{savingCard:R,errorMsg:z,deleteVisible:B,deleting:V,deleteError:H,deleteTarget:U,records:W,levelDescription:G,hydrateRows:q,loadData:ee,askDeleteRow:te,handleDelete:ne,handleSave:Z}=Xn({orgId:b(()=>s.orgId),rows:x,matchBy:`competency`,descriptionField:`note`,afterDelete:e=>{let t=Array.isArray(T.value)?T.value:[];T.value=t.filter(t=>t!==e.competency_id)},onSaved:()=>o(`saved`)}),re=b(()=>(v.value||[]).map(e=>({id:e.id,name:e.name,cluster:e.cluster||``}))),ie=b(()=>{let e={};return(re.value||[]).forEach(t=>{(e[t.cluster]=e[t.cluster]||[]).push(t)}),Object.keys(e).sort().map(t=>({label:t,items:e[t].sort((e,t)=>e.name.localeCompare(t.name))}))}),ae=b(()=>(y.value||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``,note:e.note||``})));function oe(){let e={};(re.value||[]).forEach(t=>{e[t.id]=t});let t=(Array.isArray(T.value)?T.value:T.value?[T.value]:[]).filter(t=>e[t]),n=t.length>0?Math.round(100/t.length*100)/100:0;x.value=t.map(t=>{let r=e[t];return{competency_id:t,competency_name:r.name,cluster:r.cluster,levelOptions:ae.value,recordId:``,job_management_value_id:``,weight:n}})}async function se(){try{let[e,t]=await Promise.all([P.get(`/api/v1/tenant/settings/competencies`,{params:{per_page:500}}),P.get(`/api/v1/tenant/job-management/values/clusters/managerial`)]);E.value=t.data?.data?.clusters||[];let n=new Set(E.value),r=e.data?.data||[];v.value=r.filter(e=>e.cluster&&n.has(e.cluster))}catch{v.value=[]}}async function Q(){try{let e=await P.get(`/api/v1/tenant/job-management/values`,{params:{type:`managerial`,per_page:100}});y.value=e.data?.data||[]}catch{y.value=[]}}async function ce(){if(!s.orgId){k.value=``,A.value=``;return}try{let e=(await P.get(`/api/v1/tenant/job-management/competency-groups`,{params:{organization_id:s.orgId}})).data?.data||[],t=e.find(e=>e.category===`technical`),n=e.find(e=>e.category===`managerial`);k.value=t?t.weight:``,A.value=n?n.id:``,j.value=n?n.weight:``}catch{k.value=``,A.value=``,j.value=``}}async function le({silent:e=!1}={}){let t=L.value;if(!(t===null||!s.orgId)){M.value=!0,N.value=``;try{let n=A.value;if(!n)try{let e=((await P.get(`/api/v1/tenant/job-management/competency-groups`,{params:{organization_id:s.orgId}})).data?.data||[]).find(e=>e.category===`managerial`);e&&(n=e.id)}catch{}let r={weight:t};n?await P.put(`/api/v1/tenant/job-management/competency-groups/${n}`,r):await P.post(`/api/v1/tenant/job-management/competency-groups`,{organization_id:s.orgId,category:`managerial`,weight:t}),A.value=n||``,j.value=t,e||u.add({severity:`success`,summary:l(`message.success`),detail:l(`job_management.potency_managerial_weight_saved`),life:2e3}),o(`saved`)}catch(e){N.value=e?.response?.data?.error?.message||e.message||l(`message.operation_failed`)}finally{M.value=!1}}}function ue(){let e={};(re.value||[]).forEach(t=>{e[t.id]=t});let t=[];W.value.forEach(n=>{n.competency_id&&e[n.competency_id]&&!t.includes(n.competency_id)&&t.push(n.competency_id)}),T.value=t,oe(),q()}return i(T,()=>{oe(),q()}),i(()=>s.technicalWeight,e=>{e!=null&&e!==``&&(k.value=e,le())}),t(async()=>{try{await Promise.all([se(),Q(),ee(),ce()])}finally{ue(),f.value=!1;let e=L.value;if(e!==null&&s.orgId){let t=j.value;(t===``||Math.abs(t-e)>.005)&&le({silent:!0})}}}),(t,n)=>(g(),w(`div`,Kr,[C(`div`,qr,[C(`div`,null,[C(`h3`,Jr,m(c(l)(`job_management.potency_managerial_title`)),1),C(`p`,Yr,m(c(l)(`job_management.potency_managerial_description`)),1)]),C(`div`,Xr,[C(`div`,Zr,[C(`label`,Qr,m(c(l)(`job_management.potency_managerial_weight_label`)),1),C(`div`,$r,[L.value===null?(g(),w(`span`,ti,`—`)):(g(),w(`span`,ei,m(L.value)+`%`,1))]),M.value?(g(),w(`i`,ni)):_(``,!0)]),N.value?(g(),w(`div`,ri,m(N.value),1)):_(``,!0)])]),f.value?(g(),d(X,{key:0,type:`detail`,count:1,rows:5,cols:`grid-cols-1`,padding:`p-5`})):(g(),w(S,{key:1},[p(Y,{modelValue:T.value,"onUpdate:modelValue":n[0]||=e=>T.value=e,options:ie.value,"option-label":`name`,"option-value":`id`,"option-group-label":`label`,"option-group-children":`items`,placeholder:c(l)(`job_management.potency_managerial_placeholder`),showClear:``,multiple:``},null,8,[`modelValue`,`options`,`placeholder`]),D.value?x.value.length===0?(g(),w(`div`,ai,m(c(l)(`job_management.potency_managerial_empty`)),1)):(g(),w(`div`,oi,[C(`table`,si,[C(`thead`,null,[C(`tr`,ci,[C(`th`,li,m(c(l)(`job_management.potency_table_name`)),1),C(`th`,ui,m(c(l)(`job_management.potency_table_level`)),1),C(`th`,di,m(c(l)(`job_management.potency_table_weight`)),1),C(`th`,fi,m(c(l)(`job_management.potency_table_description`)),1),C(`th`,pi,m(c(l)(`common.actions`)),1)])]),C(`tbody`,null,[(g(!0),w(S,null,h(x.value,e=>(g(),w(`tr`,{key:e.competency_id,class:`border-t border-gray-100 dark:border-gray-700 align-top`},[C(`td`,mi,[C(`div`,hi,m(e.competency_name),1),e.cluster?(g(),w(`div`,gi,m(e.cluster),1)):_(``,!0)]),C(`td`,_i,[p(Y,{modelValue:e.job_management_value_id,"onUpdate:modelValue":t=>e.job_management_value_id=t,options:e.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:c(l)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),C(`td`,vi,[p(c(K),{modelValue:e.weight,"onUpdate:modelValue":t=>e.weight=t,class:`!w-full`,min:0,max:100,suffix:`%`,size:`small`},null,8,[`modelValue`,`onUpdate:modelValue`])]),C(`td`,yi,m(c(G)(e)),1),C(`td`,bi,[p(c(O),{icon:`pi pi-trash`,severity:`danger`,text:``,rounded:``,size:`small`,disabled:c(R),"aria-label":c(l)(`common.delete`),onClick:t=>c(te)(e)},null,8,[`disabled`,`aria-label`,`onClick`])])]))),128))])])])):(g(),w(`div`,ii,m(c(l)(`job_management.potency_managerial_no_mapping`)),1)),c(z)?(g(),w(`div`,xi,m(c(z)),1)):_(``,!0),x.value.length>0?(g(),w(`div`,Si,[p(c(O),{label:c(l)(`job_management.save_managerial`),icon:`pi pi-check`,size:`small`,loading:c(R),disabled:c(R)||!e.orgId,onClick:c(Z)},null,8,[`label`,`loading`,`disabled`,`onClick`])])):_(``,!0)],64)),p(J,{visible:c(B),"onUpdate:visible":n[1]||=e=>r(B)?B.value=e:null,title:c(l)(`job_management.potency_confirm_delete_title`),message:c(l)(`job_management.potency_confirm_delete`,{name:c(U)?.competency_name||``}),loading:c(V),"error-msg":c(H),onConfirm:c(ne),onCancel:n[2]||=e=>B.value=!1},null,8,[`visible`,`title`,`message`,`loading`,`error-msg`,`onConfirm`])]))}},wi={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Ti={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},Ei={class:`text-sm text-gray-500 dark:text-gray-400`},Di={key:0,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},Oi={key:1,class:`overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700`},ki={class:`w-full text-sm`},Ai={class:`bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400`},ji={class:`px-4 py-3 font-semibold min-w-[220px]`},Mi={class:`px-4 py-3 font-semibold min-w-[260px]`},Ni={class:`px-4 py-3 font-semibold min-w-[260px]`},Pi={class:`px-4 py-3 font-semibold w-16 text-right`},Fi={class:`px-4 py-3`},Ii={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},Li={key:0,class:`mt-0.5 text-xs text-gray-500 dark:text-gray-400 leading-relaxed`},Ri={class:`px-4 py-3`},zi={class:`px-4 py-3 text-sm text-gray-600 dark:text-gray-300`},Bi={class:`px-4 py-3 text-right`},Vi={key:2,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Hi={key:3,class:`flex justify-end gap-2 pt-1`},Ui={__name:`PotencyLevelsCard`,props:{orgId:String,rows:{type:Array,default:()=>[]},optionsReady:{type:Boolean,default:!1},skeletonRows:{type:Number,default:2},titleKey:{type:String,required:!0},descriptionKey:{type:String,required:!0},emptyKey:{type:String,required:!0},saveLabelKey:{type:String,required:!0},deleteTitleKey:{type:String,required:!0},deleteMessageKey:{type:String,required:!0}},emits:[`saved`],setup(e,{emit:t}){let n=t,o=e,{t:s}=I(),l=a(!0),u=b(()=>o.rows),{savingCard:f,errorMsg:v,deleteVisible:y,deleting:x,deleteError:T,deleteTarget:E,levelDescription:D,hydrateRows:k,loadData:A,askDeleteRow:j,handleDelete:M,handleSave:N}=Xn({orgId:b(()=>o.orgId),rows:u,onSaved:()=>n(`saved`)}),P=!1;return i(()=>o.optionsReady,async e=>{if(!(!e||P)){P=!0;try{await A()}finally{k(),l.value=!1}}},{immediate:!0}),(t,n)=>(g(),w(`div`,wi,[C(`div`,null,[C(`h3`,Ti,m(c(s)(e.titleKey)),1),C(`p`,Ei,m(c(s)(e.descriptionKey)),1)]),l.value?(g(),d(X,{key:0,type:`detail`,count:1,rows:e.skeletonRows,cols:`grid-cols-1`,padding:`p-5`},null,8,[`rows`])):(g(),w(S,{key:1},[e.rows.length===0?(g(),w(`div`,Di,m(c(s)(e.emptyKey)),1)):(g(),w(`div`,Oi,[C(`table`,ki,[C(`thead`,null,[C(`tr`,Ai,[C(`th`,ji,m(c(s)(`job_management.potency_table_name`)),1),C(`th`,Mi,m(c(s)(`job_management.potency_table_level`)),1),C(`th`,Ni,m(c(s)(`job_management.potency_table_description`)),1),C(`th`,Pi,m(c(s)(`common.actions`)),1)])]),C(`tbody`,null,[(g(!0),w(S,null,h(e.rows,e=>(g(),w(`tr`,{key:e.type,class:`border-t border-gray-100 dark:border-gray-700 align-top`},[C(`td`,Fi,[C(`div`,Ii,m(e.competency_name),1),e.competency_definition?(g(),w(`div`,Li,m(e.competency_definition),1)):_(``,!0)]),C(`td`,Ri,[p(Y,{modelValue:e.job_management_value_id,"onUpdate:modelValue":t=>e.job_management_value_id=t,options:e.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:c(s)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),C(`td`,zi,m(c(D)(e)),1),C(`td`,Bi,[p(c(O),{icon:`pi pi-trash`,severity:`danger`,text:``,rounded:``,size:`small`,disabled:c(f),"aria-label":c(s)(`common.delete`),onClick:t=>c(j)(e)},null,8,[`disabled`,`aria-label`,`onClick`])])]))),128))])])])),c(v)?(g(),w(`div`,Vi,m(c(v)),1)):_(``,!0),e.rows.length>0?(g(),w(`div`,Hi,[p(c(O),{label:c(s)(e.saveLabelKey),icon:`pi pi-check`,size:`small`,loading:c(f),disabled:c(f)||!e.orgId,onClick:c(N)},null,8,[`label`,`loading`,`disabled`,`onClick`])])):_(``,!0)],64)),p(J,{visible:c(y),"onUpdate:visible":n[0]||=e=>r(y)?y.value=e:null,title:c(s)(e.deleteTitleKey),message:c(s)(e.deleteMessageKey,{name:c(E)?.competency_name||``}),loading:c(x),"error-msg":c(T),onConfirm:c(M),onCancel:n[1]||=e=>y.value=!1},null,8,[`visible`,`title`,`message`,`loading`,`error-msg`,`onConfirm`])]))}},Wi={__name:`TypesPotencyCard`,props:{orgId:String,types:{type:Array,required:!0},skeletonRows:{type:Number,default:2},titleKey:{type:String,required:!0},descriptionKey:{type:String,required:!0},emptyKey:{type:String,required:!0},saveLabelKey:{type:String,required:!0},deleteTitleKey:{type:String,required:!0},deleteMessageKey:{type:String,required:!0}},emits:[`saved`],setup(e,{emit:n}){let r=n,i=e,{t:o}=I(),s=a([]),c=a(!1);function l(e){s.value=i.types.filter(t=>(e[t.type]||[]).length>0).map(t=>({competency_id:``,competency_name:o(t.nameKey),competency_definition:``,type:t.type,levelOptions:e[t.type]||[],recordId:``,job_management_value_id:``}))}async function u(){try{let e=await P.get(`/api/v1/tenant/job-management/values/tree`),t={};(e.data?.data||[]).forEach(e=>{(e.types||[]).forEach(e=>{i.types.some(t=>t.type===e.type)&&(t[e.type]=(e.options||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``})))})}),l(t)}catch{s.value=[]}}return t(async()=>{await u(),c.value=!0}),(t,n)=>(g(),d(Ui,{"org-id":e.orgId,rows:s.value,"options-ready":c.value,"skeleton-rows":e.skeletonRows,"title-key":e.titleKey,"description-key":e.descriptionKey,"empty-key":e.emptyKey,"save-label-key":e.saveLabelKey,"delete-title-key":e.deleteTitleKey,"delete-message-key":e.deleteMessageKey,onSaved:n[0]||=e=>r(`saved`)},null,8,[`org-id`,`rows`,`options-ready`,`skeleton-rows`,`title-key`,`description-key`,`empty-key`,`save-label-key`,`delete-title-key`,`delete-message-key`]))}},Gi={__name:`ProblemSolvingPotencyCard`,props:{orgId:String},emits:[`saved`],setup(e,{emit:t}){let n=t,r=[{type:`thinking_environment`,nameKey:`job_management.problem_solving_environment`},{type:`thinking_chalenge`,nameKey:`job_management.problem_solving_challenge`}];return(t,i)=>(g(),d(Wi,{"org-id":e.orgId,types:r,"skeleton-rows":2,"title-key":`job_management.problem_solving_title`,"description-key":`job_management.problem_solving_description`,"empty-key":`job_management.problem_solving_empty`,"save-label-key":`job_management.save_problem_solving`,"delete-title-key":`job_management.problem_solving_confirm_delete_title`,"delete-message-key":`job_management.problem_solving_confirm_delete`,onSaved:i[0]||=e=>n(`saved`)},null,8,[`org-id`]))}},Ki={__name:`SkillPotencyCard`,props:{orgId:String},emits:[`saved`],setup(e,{emit:t}){let n=t,r=[{type:`communicating_influencing_skill`,nameKey:`job_management.skill_communicating_influencing`}];return(t,i)=>(g(),d(Wi,{"org-id":e.orgId,types:r,"skeleton-rows":2,"title-key":`job_management.skill_communicating_influencing_title`,"description-key":`job_management.skill_communicating_influencing_description`,"empty-key":`job_management.skill_communicating_influencing_empty`,"save-label-key":`job_management.save_skill`,"delete-title-key":`job_management.skill_confirm_delete_title`,"delete-message-key":`job_management.skill_confirm_delete`,onSaved:i[0]||=e=>n(`saved`)},null,8,[`org-id`]))}},qi={class:`space-y-4`},Ji={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Yi={class:`text-sm text-gray-500 dark:text-gray-400`},Xi={__name:`JobPotencySection`,props:{orgId:String,jobValueMap:{type:Object,default:()=>({})},competencyOptions:{type:Array,default:()=>[]}},emits:[`saved`],setup(e,{emit:t}){let n=t,{t:r}=I(),i=a(null);return(t,a)=>(g(),w(`div`,qi,[C(`div`,null,[C(`h2`,Ji,m(c(r)(`job_management.potency_competencies`)),1),C(`p`,Yi,m(c(r)(`job_management.potency_description`)),1)]),p(gr,{"org-id":e.orgId,onSaved:a[0]||=e=>n(`saved`)},null,8,[`org-id`]),p(Gr,{"org-id":e.orgId,onSaved:a[1]||=e=>n(`saved`),onWeightSaved:a[2]||=e=>i.value=e},null,8,[`org-id`]),p(Ci,{"org-id":e.orgId,"technical-weight":i.value,onSaved:a[3]||=e=>n(`saved`)},null,8,[`org-id`,`technical-weight`]),p(Gi,{"org-id":e.orgId,onSaved:a[4]||=e=>n(`saved`)},null,8,[`org-id`]),p(Ki,{"org-id":e.orgId,onSaved:a[5]||=e=>n(`saved`)},null,8,[`org-id`])]))}},Zi={class:`space-y-6`},Qi={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},$i={class:`text-sm text-gray-500 dark:text-gray-400`},ea={key:0,class:`flex items-center justify-center py-12`},ta={key:0,class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden`},na={class:`px-5 py-3 border-b border-gray-200 dark:border-gray-700 font-semibold text-sm text-gray-700 dark:text-gray-300`},ra={class:`divide-y divide-gray-100 dark:divide-gray-700`},ia={class:`hidden md:grid grid-cols-[minmax(0,2fr)_minmax(0,3fr)_auto] gap-4 px-5 py-2.5 bg-gray-50 dark:bg-gray-900/40 text-[11px] uppercase tracking-wider text-gray-400 dark:text-gray-500 font-medium`},aa={class:`text-right`},oa={class:`grid grid-cols-1 md:grid-cols-[minmax(0,2fr)_minmax(0,3fr)_auto] md:items-center gap-2`},sa={class:`text-sm font-semibold text-gray-700 dark:text-gray-200`},ca={class:`flex flex-wrap gap-1.5`},la={class:`font-medium`},ua={key:0,class:`font-mono`},da={class:`font-bold text-emerald-600 dark:text-emerald-400`},fa={class:`text-right`},pa={class:`text-sm font-bold text-gray-900 dark:text-gray-100`},ma={class:`px-5 py-4 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/40`},ha={class:`grid grid-cols-1 sm:grid-cols-2 gap-4`},ga={class:`flex items-center justify-between`},_a={class:`text-xs text-gray-500 dark:text-gray-400`},va={class:`text-sm font-bold text-emerald-600 dark:text-emerald-400`},ya={class:`flex items-center justify-between`},ba={class:`text-xs text-gray-500 dark:text-gray-400`},xa={class:`text-sm font-bold text-blue-600 dark:text-blue-400`},Sa={key:2},Ca={class:`flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500`},wa={class:`text-sm font-medium`},Ta={class:`text-xs mt-1`},Ea={class:`flex justify-end gap-3`},Da=`/api/v1/tenant/job-management/scores/org`,Oa={__name:`JobScoreSection`,props:{orgId:String},emits:[`saved`],setup(e,{emit:n}){let r=e,i=n,{t:o}=I(),s=F(),l=a(!1),u=a(!1),y=a(null),x=[{key:`education_experience`,labelKey:`job_management.education_experience`,points:[{labelKey:`job_management.group_education`,level:`education_level`,pts:`education_points`},{labelKey:`job_management.group_experience`,level:`experience_level`,pts:`experience_points`}]},{key:`potentials`,labelKey:`job_management.score_potentials`,points:[{labelKey:`job_management.average_level`,level:`average_level`,pts:null}]},{key:`competencies`,labelKey:`job_management.potency_competencies`,score:`base_score`,points:[{labelKey:`job_management.potency_technical_title`,level:`technical_average_level`,pts:`technical_points`},{labelKey:`job_management.potency_managerial_title`,level:`managerial_average_level`,pts:`managerial_points`},{labelKey:`job_management.skill_communicating_influencing`,level:`communication_level`,pts:`communication_points`}]},{key:`problem_solving`,labelKey:`job_management.problem_solving_title`,points:[{labelKey:`job_management.problem_solving_environment`,level:`environment_level`,pts:`environment_points`},{labelKey:`job_management.problem_solving_challenge`,level:`challenge_level`,pts:`challenge_points`}]},{key:`financial_authority`,labelKey:`job_management.financials`,points:[{labelKey:`job_management.cash_level`,level:`money_level`,pts:`money_points`},{labelKey:`job_management.authority_level`,level:`authority_level`,pts:`authority_points`},{labelKey:`job_management.impact_level`,level:`impact_level`,pts:`impact_points`}]},{key:`asset_authority`,labelKey:`job_management.assets`,points:[{labelKey:`job_management.asset_type`,level:`asset_value_level`,pts:`asset_value_points`},{labelKey:`job_management.authority_level`,level:`asset_authority_level`,pts:`asset_authority_points`}]},{key:`subordinate_control`,labelKey:`job_management.subordinate_controls`,points:[{labelKey:`job_management.score_level`,level:`level`,pts:`points`}]},{key:`work_scope`,labelKey:`job_management.relationships`,points:[{labelKey:`job_management.relationship_group_scope`,level:`scope_level`,pts:`scope_points`},{labelKey:`job_management.frequency`,level:`frequency_level`,pts:`frequency_points`}]},{key:`work_activity`,labelKey:`job_management.activities`,points:[{labelKey:`job_management.score_level`,level:`level`,pts:`points`}]},{key:`work_risk`,labelKey:`job_management.risks`,points:[{labelKey:`job_management.environment_risk`,level:`environment_level`,pts:`environment_points`},{labelKey:`job_management.hazard`,level:`hazard_level`,pts:`hazard_points`}]}],T=b(()=>{if(!y.value?.components)return null;try{return JSON.parse(y.value.components)}catch{return null}}),E=b(()=>T.value?x.map(e=>{let t=T.value[e.key]||{};return{key:e.key,labelKey:e.labelKey,score:t[e.score||`score`]??0,points:e.points.map(e=>({labelKey:e.labelKey,level:t[e.level]??null,points:e.pts==null?null:t[e.pts]??0}))}}):[]);function D(e){return e?.toLocaleString?.(`id-ID`)??`-`}function k(e){return e==null?`-`:String(e)}async function A(){if(r.orgId){l.value=!0;try{let e=await P.get(`${Da}/${r.orgId}`);y.value=e.data?.data||null,i(`saved`)}catch{y.value=null}finally{l.value=!1}}}async function j(){if(r.orgId){u.value=!0;try{let e=await P.put(`${Da}/${r.orgId}`,{components:null});y.value=e.data?.data||null,s.add({severity:`success`,summary:o(`message.success`),detail:o(`job_management.score_calculated`),life:2e3})}catch(e){s.add({severity:`error`,summary:o(`message.error`),detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{u.value=!1}}}return t(A),(e,t)=>(g(),w(`div`,Zi,[C(`div`,null,[C(`h2`,Qi,m(c(o)(`job_management.scores`)),1),C(`p`,$i,m(c(o)(`job_management.score_description`)),1)]),l.value?(g(),w(`div`,ea,[...t[0]||=[C(`i`,{class:`pi pi-spin pi-spinner text-emerald-500 text-2xl`},null,-1)]])):y.value?(g(),w(S,{key:1},[E.value.length?(g(),w(`div`,ta,[C(`div`,na,m(c(o)(`job_management.component_breakdown`)),1),C(`div`,ra,[C(`div`,ia,[C(`span`,null,m(c(o)(`job_management.score_component`)),1),C(`span`,null,m(c(o)(`job_management.score_points`)),1),C(`span`,aa,m(c(o)(`job_management.score_score`)),1)]),(g(!0),w(S,null,h(E.value,e=>(g(),w(`div`,{key:e.key,class:`px-5 py-1`},[C(`div`,oa,[C(`div`,sa,m(c(o)(e.labelKey)),1),C(`div`,ca,[(g(!0),w(S,null,h(e.points,e=>(g(),w(`span`,{key:e.labelKey,class:f([`inline-flex items-center gap-1 text-[11px] px-2 py-0.5 rounded-md border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-900/40 text-gray-600 dark:text-gray-300`,{"opacity-50":e.level==null}])},[C(`span`,la,m(c(o)(e.labelKey)),1),e.level==null?_(``,!0):(g(),w(`span`,ua,`Lv.`+m(k(e.level)),1)),e.points==null?e.level==null?(g(),w(S,{key:2},[v(`—`)],64)):_(``,!0):(g(),w(S,{key:1},[t[1]||=C(`i`,{class:`pi pi-arrow-right text-[8px] opacity-60`},null,-1),C(`span`,da,m(D(e.points)),1)],64))],2))),128))]),C(`div`,fa,[C(`span`,pa,m(D(e.score)),1)])])]))),128))]),C(`div`,ma,[C(`div`,ha,[C(`div`,ga,[C(`span`,_a,m(c(o)(`job_management.value_with_financial`)),1),C(`span`,va,m(D(y.value.job_value_with_financial)),1)]),C(`div`,ya,[C(`span`,ba,m(c(o)(`job_management.value_without_financial`)),1),C(`span`,xa,m(D(y.value.job_value_without_financial)),1)])])])])):_(``,!0)],64)):(g(),w(`div`,Sa,[C(`div`,Ca,[t[2]||=C(`i`,{class:`pi pi-calculator text-4xl mb-3 opacity-50`},null,-1),C(`p`,wa,m(c(o)(`job_management.no_score`)),1),C(`p`,Ta,m(c(o)(`job_management.score_hint`)),1)])])),C(`div`,Ea,[p(c(O),{label:c(o)(`common.refresh`),icon:`pi pi-refresh`,size:`small`,text:``,onClick:A},null,8,[`label`]),y.value?(g(),d(c(O),{key:0,label:c(o)(`job_management.recalculate`),icon:`pi pi-calculator`,size:`small`,severity:`info`,loading:u.value,onClick:j},null,8,[`label`,`loading`])):_(``,!0)])]))}},ka={class:`grid grid-cols-1 md:grid-cols-3 gap-4`},Aa={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4`},ja={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},Ma={class:`text-2xl font-bold text-emerald-600 dark:text-emerald-400`},Na={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4`},Pa={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},Fa={class:`text-2xl font-bold text-blue-600 dark:text-blue-400`},Ia={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4`},La={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},Ra={class:`mt-3 flex flex-wrap items-center gap-2`},za={key:1,class:`text-[10px] text-gray-400`},Ba={key:0,class:`text-[10px] text-gray-400 mt-2`},Va=`/api/v1/tenant/job-management/scores/org`,Ha={__name:`JobScoreSummary`,props:{orgId:String},setup(e,{expose:n}){let r=e,{t:o}=I(),s=a(!0),l=a(null);function u(e){return e?.toLocaleString?.(`id-ID`)??`-`}async function f(){if(r.orgId){s.value=!0;try{let e=await P.get(`${Va}/${r.orgId}`);l.value=e.data?.data||null}catch{l.value=null}finally{s.value=!1}}}return n({refresh:f}),i(()=>r.orgId,f),t(f),(e,t)=>(g(),w(`div`,ka,[s.value&&!l.value?(g(),w(S,{key:0},h(3,e=>C(`div`,{key:e,class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4 animate-pulse`},[...t[0]||=[C(`div`,{class:`h-3 w-24 bg-gray-200 dark:bg-gray-700 rounded mb-2`},null,-1),C(`div`,{class:`h-7 w-16 bg-gray-200 dark:bg-gray-700 rounded`},null,-1)]])),64)):(g(),w(S,{key:1},[C(`div`,Aa,[C(`div`,ja,m(c(o)(`job_management.value_with_financial`)),1),C(`div`,Ma,m(u(l.value?.job_value_with_financial)),1)]),C(`div`,Na,[C(`div`,Pa,m(c(o)(`job_management.value_without_financial`)),1),C(`div`,Fa,m(u(l.value?.job_value_without_financial)),1)]),C(`div`,Ia,[C(`div`,La,m(c(o)(`job_management.has_financial_authority`)),1),p(c(R),{value:l.value?l.value.has_financial_authority?c(o)(`common.yes`):c(o)(`common.no`):`-`,severity:l.value?.has_financial_authority?`success`:`danger`,class:`!text-xs`},null,8,[`value`,`severity`]),C(`div`,Ra,[l.value?(g(),d(c(R),{key:0,value:l.value.is_complete?c(o)(`job_management.score_complete`):c(o)(`job_management.score_incomplete`),severity:l.value.is_complete?`success`:`warning`,icon:l.value.is_complete?`pi pi-check-circle`:`pi pi-exclamation-triangle`,class:`!text-xs`},null,8,[`value`,`severity`,`icon`])):_(``,!0),l.value?.is_complete&&l.value.completed_at?(g(),w(`span`,za,m(c(o)(`job_management.completed_at`))+`: `+m(l.value.completed_at),1)):_(``,!0)]),l.value?.calculated_at?(g(),w(`div`,Ba,m(c(o)(`job_management.calculated_at`))+`: `+m(l.value.calculated_at),1)):_(``,!0)])],64))]))}},Ua={class:`max-w-full mx-auto`},Wa={key:0,class:`flex gap-6`},Ga={class:`w-56 space-y-2`},Ka={class:`flex-1 space-y-3`},qa={key:1,class:`flex gap-6`},Ja={class:`w-56 shrink-0 space-y-1`},Ya=[`onClick`,`onKeydown`],Xa={key:0,class:`pi pi-check text-xs`},Za={class:`flex-1 min-w-0`},Qa={key:0,class:`pi pi-check-circle text-emerald-400 text-xs shrink-0`},$a={class:`flex-1 min-w-0 space-y-4`},eo={class:`flex flex-col md:flex-row gap-4`},to={class:`md:w-72 shrink-0 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4`},no={class:`flex items-center gap-2 mb-3`},ro={class:`text-sm font-semibold text-gray-800 dark:text-gray-100 truncate`},io={class:`flex items-center justify-between gap-2`},ao={class:`text-[10px] uppercase tracking-wider text-gray-400 dark:text-gray-500`},oo={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 font-mono truncate`},so={class:`flex-1 min-w-0`},co={__name:`JobManagementForm`,setup(e){let r=A(),o=j(),{t:s}=I(),l=F(),u=b(()=>o.query.org_id||``),v=a(0),y=a(!0),x=a(Array(15).fill(!1)),T=a(``),E=a(``),D=a(``),O=a(``),M=a(``),N=a([]),L=a([]),R=a([]),z=a({}),B=a([]),V=a(null),H=[{labelKey:`job_management.identifications`,icon:`pi pi-id-card`,comp:ce},{labelKey:`job_management.objectives`,icon:`pi pi-bullseye`,comp:_e},{labelKey:`job_management.responsibilities_title`,icon:`pi pi-list-check`,comp:st},{labelKey:`job_management.education_experience`,icon:`pi pi-graduation-cap`,comp:Pe},{labelKey:`job_management.potency_competencies`,icon:`pi pi-star`,comp:Xi},{labelKey:`job_management.financials`,icon:`pi pi-money-bill`,comp:Jn},{labelKey:`job_management.assets`,icon:`pi pi-box`,comp:Fn},{labelKey:`job_management.subordinate_controls`,icon:`pi pi-sitemap`,comp:En},{labelKey:`job_management.relationships`,icon:`pi pi-share-alt`,comp:_n},{labelKey:`job_management.activities`,icon:`pi pi-bolt`,comp:Mt},{labelKey:`job_management.risks`,icon:`pi pi-exclamation-triangle`,comp:Vt},{labelKey:`job_management.hr_authorities`,icon:`pi pi-users`,comp:ht},{labelKey:`job_management.op_authorities`,icon:`pi pi-cog`,comp:Ct},{labelKey:`job_management.scores`,icon:`pi pi-calculator`,comp:Oa}],U=b(()=>H[v.value]?.comp||null);function W(e){return v.value===e?`bg-emerald-50 dark:bg-emerald-900/20 ring-1 ring-emerald-300 dark:ring-emerald-700`:(x.value[e],`hover:bg-gray-50 dark:hover:bg-gray-800`)}function G(e){return v.value===e?`bg-emerald-600 text-white`:x.value[e]?`bg-emerald-100 dark:bg-emerald-800 text-emerald-600 dark:text-emerald-300`:`bg-gray-200 dark:bg-gray-600 text-gray-600 dark:text-gray-300`}function K(e){return v.value===e?`text-emerald-700 dark:text-emerald-300`:x.value[e]?`text-emerald-600 dark:text-emerald-400`:`text-gray-700 dark:text-gray-300`}function q(e){v.value=e,r.replace({query:{...o.query,section:String(e)}})}function ee(e){typeof e==`number`&&(x.value[e]=!0),V.value?.refresh()}async function te(){if(u.value)try{let e=(await P.get(`/api/v1/tenant/organizations/${u.value}`)).data?.data;e&&(T.value=e.nomenclature||``,E.value=e.full_code||e.code||``,D.value=e.organization_summary_id||``,O.value=e.grading_id||``,M.value=e.job_family_id||``)}catch{}}async function J(){try{let[e,t,n,r]=await Promise.all([P.get(`/api/v1/tenant/settings/gradings?per_page=100`),P.get(`/api/v1/tenant/job-management/values?per_page=500`),P.get(`/api/v1/tenant/competency/competencies?per_page=200`).catch(()=>({data:{data:[]}})),P.get(`/api/v1/tenant/settings/job-families?per_page=100`)]);N.value=(e.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id})),L.value=(r.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id}));let i=t.data?.data||[];R.value=i.map(e=>({label:`${e.type}${e.level?` Lv.`+e.level:``}${e.descriptions?` — `+e.descriptions:``}`,value:e.id,type:e.type,level:e.level,descriptions:e.descriptions}));let a={};i.forEach(e=>{a[e.type]||(a[e.type]=[]),a[e.type].push({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``,type_group:e.type_group||``,description_group:e.description_group||``})}),z.value=a,B.value=(n.data?.data||[]).map(e=>({label:e.name||e.code,value:e.id,field:e.field||``,definition:e.definition||``}))}catch{}}return i(u,(e,t)=>{e!==t&&(x.value=Array(H.length).fill(!1),T.value=``,E.value=``,D.value=``,O.value=``,M.value=``,te())}),t(async()=>{try{await Promise.all([te(),J()]);let e=parseInt(o.query.section);!isNaN(e)&&e>=0&&e<H.length&&(v.value=e)}catch(e){l.add({severity:`error`,summary:s(`message.error`),detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{y.value=!1}}),(e,t)=>(g(),w(`div`,Ua,[y.value?(g(),w(`div`,Wa,[C(`div`,Ga,[(g(),w(S,null,h(8,e=>C(`div`,{key:e,class:`h-12 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))]),C(`div`,Ka,[(g(),w(S,null,h(6,e=>C(`div`,{key:e,class:`h-9 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))])])):(g(),w(`div`,qa,[C(`div`,Ja,[(g(),w(S,null,h(H,(e,t)=>C(`div`,{key:t,role:`button`,tabindex:0,class:f([`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150 cursor-pointer select-none`,W(t)]),onClick:e=>q(t),onKeydown:k(e=>q(t),[`enter`])},[C(`div`,{class:f([`w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0 transition-colors duration-150`,G(t)])},[x.value[t]?(g(),w(`i`,Xa)):(g(),w(`i`,{key:1,class:f(e.icon)},null,2))],2),C(`div`,Za,[C(`div`,{class:f([`text-sm font-medium truncate`,K(t)])},m(c(s)(e.labelKey)),3)]),x.value[t]?(g(),w(`i`,Qa)):_(``,!0)],42,Ya)),64))]),C(`div`,$a,[(g(),w(`div`,{key:`summary-${u.value}`,class:`sticky top-0 z-10 bg-white dark:bg-gray-900 pt-1 pb-3`},[C(`div`,eo,[C(`div`,to,[C(`div`,no,[t[0]||=C(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400`},[C(`i`,{class:`pi pi-briefcase text-sm`})],-1),C(`h3`,ro,m(T.value||c(s)(`job_management.job_info_untitled`)),1)]),C(`div`,io,[C(`span`,ao,m(c(s)(`organization.full_code`)),1),C(`span`,oo,m(E.value||`-`),1)])]),C(`div`,so,[p(Ha,{ref_key:`scoreSummaryRef`,ref:V,"org-id":u.value},null,8,[`org-id`])])])])),(g(),d(n(U.value),{key:`${v.value}-${u.value}`,"org-id":u.value,"org-name":T.value,"org-code":E.value,"org-summary-id":D.value,"org-grading-id":O.value,"org-job-family-id":M.value,"grading-options":N.value,"job-family-options":L.value,"job-value-options":R.value,"competency-options":B.value,"job-value-map":z.value,onSaved:ee},null,40,[`org-id`,`org-name`,`org-code`,`org-summary-id`,`org-grading-id`,`org-job-family-id`,`grading-options`,`job-family-options`,`job-value-options`,`competency-options`,`job-value-map`]))])]))]))}};export{co as default};