const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["assets/quill-BLmY9xB4.js","assets/rolldown-runtime-QTnfLwEv.js"])))=>i.map(i=>d[i]);
import{A as e,C as t,E as n,H as r,N as i,O as a,P as o,U as s,V as c,W as l,Z as u,c as d,dt as f,ft as p,h as m,j as h,l as g,m as _,mt as v,o as y,pt as b,q as x,r as S,s as C,u as w,y as T}from"./runtime-core.esm-bundler-BmFWOnhO.js";import{l as E,t as D}from"./button-BeZOQkK0.js";import{A as O,a as k}from"./ripple-CVbnBGPL.js";import{_ as A,f as j,l as M,p as N,s as P,t as F}from"./index-DiEK1uSa.js";import{t as I}from"./useI18n-BwB6MR-x.js";import{r as L}from"./responseHandler-BJxA-JZj.js";import{t as R}from"./tag-DfpGb59G.js";import{t as z}from"./FormRow-Dj7iXuyh.js";import{t as B}from"./baseeditableholder-CNNL9sCP.js";import{t as V}from"./textarea-BWEuClvt.js";import{t as H}from"./TextInput-CA3K1nYt.js";import{n as U,t as W}from"./column-D6jaMNT_.js";import{t as G}from"./select-KmrqMMzj.js";import{t as K}from"./inputnumber-ETSzxLCc.js";import{t as q}from"./multiselect-R_DmlfFc.js";import{t as ee}from"./toggleswitch-Cakebgdp.js";import{t as te}from"./SkeletonTable-CikfhXMe.js";import{t as J}from"./ConfirmDeleteDialog-BXT5N2w0.js";import{t as Y}from"./SelectLabel-CY3qiDu1.js";import{t as X}from"./SkeletonCard-CGT909cs.js";var ne={class:`space-y-4`},Z={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},re={class:`text-sm text-gray-500 dark:text-gray-400`},ie={class:`max-w-2xl`},ae={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},oe={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},se={class:`flex justify-end pt-2`},Q=`/api/v1/tenant/job-management/identifications`,ce={__name:`JobIdentificationSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},orgGradingId:{type:String,default:``},orgJobFamilyId:{type:String,default:``},gradingOptions:{type:Array,default:()=>[]},jobFamilyOptions:{type:Array,default:()=>[]}},emits:[`saved`],setup(e,{emit:t}){let i=t,o=e,{t:s}=I(),c=A(),l=x(!1),f=x(!0),p=x(``),h=x({}),_=x(``),b=x({grading_id:``}),S=y(()=>{let e=o.jobFamilyOptions.find(e=>e.value===o.orgJobFamilyId);return e?e.label:o.orgJobFamilyId||`-`});function T(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function E(){if(!o.orgId){f.value=!1;return}try{let e=(await M.get(Q,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];_.value=t.id,b.value.grading_id=t.grading_id||o.orgGradingId||``}else b.value.grading_id=o.orgGradingId||``}catch{b.value.grading_id=o.orgGradingId||``}finally{f.value=!1}}async function O(){if(p.value=``,h.value={},!b.value.grading_id){p.value=s(`job_management.grading_required`);return}l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,grading_id:b.value.grading_id,organization_id:o.orgId};if(_.value)await M.put(`${Q}/${_.value}`,{grading_id:b.value.grading_id});else{let t=await M.post(Q,e);_.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),i(`saved`)}catch(e){let t=T(e);Object.keys(t).length>0?(h.value=t,p.value=Object.values(t).join(`, `)):p.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}return n(E),(t,n)=>(a(),w(`div`,ne,[C(`div`,null,[C(`h2`,Z,v(u(s)(`job_management.identifications`)),1),C(`p`,re,v(u(s)(`job_management.identification_description`)),1)]),C(`div`,ie,[f.value?(a(),d(X,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(a(),w(`div`,ae,[m(z,{label:u(s)(`organization.job_family`)},{default:r(()=>[m(H,{"model-value":S.value,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(z,{label:u(s)(`organization.grading`)},{default:r(()=>[m(u(G),{modelValue:b.value.grading_id,"onUpdate:modelValue":n[0]||=e=>b.value.grading_id=e,options:e.gradingOptions,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`organization.select_grading`),class:`w-full`,size:`small`,invalid:!!h.value.grading_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`]),p.value?(a(),w(`div`,oe,v(p.value),1)):g(``,!0),C(`div`,se,[m(u(D),{label:u(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:!b.value.grading_id,onClick:O},null,8,[`label`,`loading`,`disabled`])])]))])]))}},le={class:`space-y-4`},ue={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},de={class:`text-sm text-gray-500 dark:text-gray-400`},fe={class:`max-w-2xl`},pe={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},me={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},he={class:`flex justify-end gap-2 pt-2`},ge=`/api/v1/tenant/job-management/objectives`,_e={__name:`JobObjectiveSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let i=t,o=e,{t:s}=I(),c=A(),l=x(!1),p=x(!0),h=x(!1),_=x(``),y=x({}),b=x(``),S=x(!1),T=x(``),E=x({objective:``});function O(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function k(){if(!o.orgId){p.value=!1;return}try{let e=(await M.get(ge,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];b.value=t.id,E.value.objective=t.objective||``}}catch{}finally{p.value=!1}}async function j(){_.value=``,y.value={},l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,objective:E.value.objective||``,organization_id:o.orgId};if(b.value)await M.put(`${ge}/${b.value}`,{objective:E.value.objective||``});else{let t=await M.post(ge,e);b.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),i(`saved`)}catch(e){let t=O(e);Object.keys(t).length>0?(y.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}async function N(){if(b.value){h.value=!0,T.value=``;try{await M.delete(`${ge}/${b.value}`),S.value=!1,b.value=``,E.value.objective=``,i(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{h.value=!1}}}return n(k),(e,t)=>(a(),w(`div`,le,[C(`div`,null,[C(`h2`,ue,v(u(s)(`job_management.objectives`)),1),C(`p`,de,v(u(s)(`job_management.objective_description`)),1)]),C(`div`,fe,[p.value?(a(),d(X,{key:0,type:`detail`,count:1,rows:3,cols:`grid-cols-1`,padding:`p-5`})):(a(),w(`div`,pe,[m(z,{label:u(s)(`job_management.objective`)},{default:r(()=>[m(u(V),{modelValue:E.value.objective,"onUpdate:modelValue":t[0]||=e=>E.value.objective=e,rows:`3`,class:f([`w-full`,{"p-invalid":y.value.objective}]),placeholder:u(s)(`job_management.objective`)+`...`},null,8,[`modelValue`,`class`,`placeholder`])]),_:1},8,[`label`]),_.value?(a(),w(`div`,me,v(_.value),1)):g(``,!0),C(`div`,he,[b.value?(a(),d(u(D),{key:0,label:u(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[1]||=e=>S.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:b.value?u(s)(`common.update`):u(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:l.value,onClick:j},null,8,[`label`,`loading`,`disabled`])])]))]),m(J,{visible:S.value,"onUpdate:visible":t[2]||=e=>S.value=e,loading:h.value,"error-msg":T.value,onConfirm:N,onCancel:t[3]||=e=>S.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},ve={class:`space-y-4`},ye={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},be={class:`text-sm text-gray-500 dark:text-gray-400`},xe={class:`max-w-2xl`},Se={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Ce={class:`pt-1`},we={class:`flex items-center gap-2 mb-3`},Te={class:`text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400`},Ee={class:`space-y-4`},De={class:`pt-4 border-t border-gray-200 dark:border-gray-700`},Oe={class:`flex items-center gap-2 mb-3`},ke={class:`text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400`},Ae={class:`space-y-4`},je={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Me={class:`flex justify-end gap-2 pt-2`},Ne=`/api/v1/tenant/job-management/education-experiences`,Pe={__name:`JobEduExpSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:t}){let i=t,o=e,{t:s}=I(),c=A(),l=x(!1),p=x(!0),h=x(!1),_=x(``),y=x({}),b=x(``),S=x(!1),T=x(``),E=x({education_id:``,education_major_id:[],job_family_id:[],experience_id:``}),O=x([]),k=x([]),j=x([]),N=x([]);async function P(){try{let[e,t,n,r]=await Promise.all([M.get(`/api/v1/tenant/job-management/values`,{params:{type:`education`,per_page:100}}),M.get(`/api/v1/tenant/job-management/values`,{params:{type:`experience`,per_page:100}}),M.get(`/api/v1/tenant/settings/education-majors?per_page=200`),M.get(`/api/v1/tenant/settings/job-families?per_page=100`)]);k.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),O.value=(t.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),j.value=(n.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id})),N.value=(r.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id}))}catch{}}async function F(){if(o.orgId)try{let e=(await M.get(Ne,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];b.value=t.id,E.value.education_id=t.education_id||``,E.value.education_major_id=Array.isArray(t.education_major_id)?t.education_major_id:[],E.value.job_family_id=Array.isArray(t.job_family_id)?t.job_family_id:[],E.value.experience_id=t.experience_id||``}}catch{}}async function R(){_.value=``,y.value={},l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,education_id:E.value.education_id||null,education_major_id:E.value.education_major_id||[],job_family_id:E.value.job_family_id||[],experience_id:E.value.experience_id||null,organization_id:o.orgId};if(b.value)await M.put(`${Ne}/${b.value}`,{education_id:E.value.education_id||``,education_major_id:E.value.education_major_id||[],job_family_id:E.value.job_family_id||[],experience_id:E.value.experience_id||``});else{let t=await M.post(Ne,e);b.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),i(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(y.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}async function B(){if(b.value){h.value=!0,T.value=``;try{await M.delete(`${Ne}/${b.value}`),S.value=!1,b.value=``,E.value.education_id=``,E.value.education_major_id=[],E.value.job_family_id=[],E.value.experience_id=``,i(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{h.value=!1}}}return n(async()=>{try{await Promise.all([P(),F()])}finally{p.value=!1}}),(e,t)=>(a(),w(`div`,ve,[C(`div`,null,[C(`h2`,ye,v(u(s)(`job_management.education_experience`)),1),C(`p`,be,v(u(s)(`job_management.education_experience_description`)),1)]),C(`div`,xe,[p.value?(a(),d(X,{key:0,type:`detail`,count:1,rows:6,cols:`grid-cols-1`,padding:`p-5`})):(a(),w(`div`,Se,[C(`div`,Ce,[C(`div`,we,[t[7]||=C(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400`},[C(`i`,{class:`pi pi-graduation-cap text-sm`})],-1),C(`h3`,Te,v(u(s)(`job_management.group_education`)),1),t[8]||=C(`div`,{class:`flex-1 border-t border-gray-200 dark:border-gray-700`},null,-1)]),C(`div`,Ee,[m(z,{label:u(s)(`job_management.education_level`),errors:y.value?.education_id},{default:r(()=>[m(Y,{modelValue:E.value.education_id,"onUpdate:modelValue":t[0]||=e=>E.value.education_id=e,options:k.value,placeholder:u(s)(`job_values.select_education`),class:f({"p-invalid":y.value?.education_id})},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),m(z,{label:u(s)(`job_management.education_major`),errors:y.value?.education_major_id},{default:r(()=>[m(u(q),{modelValue:E.value.education_major_id,"onUpdate:modelValue":t[1]||=e=>E.value.education_major_id=e,options:j.value,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),class:`w-full`,size:`small`,filter:``,showClear:``,display:`chip`,maxSelectedLabels:2,invalid:!!y.value.education_major_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`])])]),C(`div`,De,[C(`div`,Oe,[t[9]||=C(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400`},[C(`i`,{class:`pi pi-briefcase text-sm`})],-1),C(`h3`,ke,v(u(s)(`job_management.group_experience`)),1),t[10]||=C(`div`,{class:`flex-1 border-t border-gray-200 dark:border-gray-700`},null,-1)]),C(`div`,Ae,[m(z,{label:u(s)(`job_management.experience_range`),errors:y.value?.experience_id},{default:r(()=>[m(Y,{modelValue:E.value.experience_id,"onUpdate:modelValue":t[2]||=e=>E.value.experience_id=e,options:O.value,placeholder:u(s)(`common.select`),class:f({"p-invalid":y.value?.experience_id})},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),m(z,{label:u(s)(`job_management.job_family`),errors:y.value?.job_family_id},{default:r(()=>[m(u(q),{modelValue:E.value.job_family_id,"onUpdate:modelValue":t[3]||=e=>E.value.job_family_id=e,options:N.value,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),class:`w-full`,size:`small`,filter:``,showClear:``,display:`chip`,maxSelectedLabels:2,invalid:!!y.value.job_family_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`])])]),_.value?(a(),w(`div`,je,v(_.value),1)):g(``,!0),C(`div`,Me,[b.value?(a(),d(u(D),{key:0,label:u(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[4]||=e=>S.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:b.value?u(s)(`common.update`):u(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:l.value,onClick:R},null,8,[`label`,`loading`,`disabled`])])]))]),m(J,{visible:S.value,"onUpdate:visible":t[5]||=e=>S.value=e,loading:h.value,"error-msg":T.value,onConfirm:B,onCancel:t[6]||=e=>S.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Fe=k.extend({name:`editor`,style:`
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
`,classes:{root:function(e){return[`p-editor`,{"p-invalid":e.instance.$invalid}]},toolbar:`p-editor-toolbar`,content:`p-editor-content`}}),Ie={name:`BaseEditor`,extends:B,props:{placeholder:String,readonly:Boolean,formats:Array,editorStyle:null,modules:null},style:Fe,provide:function(){return{$pcEditor:this,$parentInstance:this}}};function Le(e){"@babel/helpers - typeof";return Le=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},Le(e)}function Re(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function ze(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?Re(Object(n),!0).forEach(function(t){Be(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):Re(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function Be(e,t,n){return(t=Ve(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function Ve(e){var t=He(e,`string`);return Le(t)==`symbol`?t:t+``}function He(e,t){if(Le(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(Le(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var Ue=function(){try{return window.Quill}catch{return null}}(),We={name:`Editor`,extends:Ie,inheritAttrs:!1,emits:[`text-change`,`selection-change`,`load`],quill:null,watch:{modelValue:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},d_value:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},readonly:function(){this.handleReadOnlyChange()}},mounted:function(){var e=this,t={modules:ze({toolbar:this.$refs.toolbarElement},this.modules),readOnly:this.readonly,theme:`snow`,formats:this.formats,placeholder:this.placeholder};Ue?(this.quill=new Ue(this.$refs.editorElement,t),this.initQuill(),this.handleLoad()):P(()=>import(`./quill-BLmY9xB4.js`).then(function(n){n&&O(e.$refs.editorElement)&&(n.default?e.quill=new n.default(e.$refs.editorElement,t):e.quill=new n(e.$refs.editorElement,t),e.initQuill())}),__vite__mapDeps([0,1])).then(function(){e.handleLoad()})},beforeUnmount:function(){this.quill=null},methods:{renderValue:function(e){if(this.quill)if(e){var t=this.quill.clipboard.convert({html:e});this.quill.setContents(t)}else this.quill.setText(``)},initQuill:function(){var e=this;this.renderValue(this.d_value),this.quill.on(`text-change`,function(t,n,r){if(r===`user`){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();i===`<p><br></p>`&&(i=``),e.writeValue(i),e.$emit(`text-change`,{htmlValue:i,textValue:a,delta:t,source:r,instance:e.quill})}}),this.quill.on(`selection-change`,function(t,n,r){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();e.$emit(`selection-change`,{htmlValue:i,textValue:a,range:t,oldRange:n,source:r,instance:e.quill})})},handleLoad:function(){this.quill&&this.quill.getModule(`toolbar`)&&this.$emit(`load`,{instance:this.quill})},handleReadOnlyChange:function(){this.quill&&this.quill.enable(!this.readonly)}}};function Ge(e,n,r,i,o,s){return a(),w(`div`,t({class:e.cx(`root`)},e.ptmi(`root`)),[C(`div`,t({ref:`toolbarElement`,class:e.cx(`toolbar`)},e.ptm(`toolbar`)),[h(e.$slots,`toolbar`,{},function(){return[C(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[C(`select`,t({class:`ql-header`,defaultValue:`0`},e.ptm(`header`)),[C(`option`,t({value:`1`},e.ptm(`option`)),`Heading`,16),C(`option`,t({value:`2`},e.ptm(`option`)),`Subheading`,16),C(`option`,t({value:`0`},e.ptm(`option`)),`Normal`,16)],16),C(`select`,t({class:`ql-font`},e.ptm(`font`)),[C(`option`,p(T(e.ptm(`option`))),null,16),C(`option`,t({value:`serif`},e.ptm(`option`)),null,16),C(`option`,t({value:`monospace`},e.ptm(`option`)),null,16)],16)],16),C(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[C(`button`,t({class:`ql-bold`,type:`button`},e.ptm(`bold`)),null,16),C(`button`,t({class:`ql-italic`,type:`button`},e.ptm(`italic`)),null,16),C(`button`,t({class:`ql-underline`,type:`button`},e.ptm(`underline`)),null,16)],16),C(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[C(`select`,t({class:`ql-color`},e.ptm(`color`)),null,16),C(`select`,t({class:`ql-background`},e.ptm(`background`)),null,16)],16),C(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[C(`button`,t({class:`ql-list`,value:`ordered`,type:`button`},e.ptm(`list`)),null,16),C(`button`,t({class:`ql-list`,value:`bullet`,type:`button`},e.ptm(`list`)),null,16),C(`select`,t({class:`ql-align`},e.ptm(`select`)),[C(`option`,t({defaultValue:``},e.ptm(`option`)),null,16),C(`option`,t({value:`center`},e.ptm(`option`)),null,16),C(`option`,t({value:`right`},e.ptm(`option`)),null,16),C(`option`,t({value:`justify`},e.ptm(`option`)),null,16)],16)],16),C(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[C(`button`,t({class:`ql-link`,type:`button`},e.ptm(`link`)),null,16),C(`button`,t({class:`ql-image`,type:`button`},e.ptm(`image`)),null,16),C(`button`,t({class:`ql-code-block`,type:`button`},e.ptm(`codeBlock`)),null,16)],16),C(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[C(`button`,t({class:`ql-clean`,type:`button`},e.ptm(`clean`)),null,16)],16)]})],16),C(`div`,t({ref:`editorElement`,class:e.cx(`content`),style:e.editorStyle},e.ptm(`content`)),null,16)],16)}We.render=Ge;var Ke={key:0,class:`text-gray-500 dark:text-gray-400 text-xs`},qe=[`innerHTML`],Je={key:2,class:`text-gray-800 dark:text-gray-100`},Ye={class:`flex items-center gap-1`},Xe={__name:`DataTableSection`,props:{items:Array,loading:Boolean,total:Number,columns:{type:Array,default:()=>[]},entity:String,orgId:String,onLoad:Function},emits:[`edit`,`delete`],setup(t){let o=t,{t:c}=I(),l=x(1),f=x(15),p=y(()=>(l.value-1)*f.value),_=y(()=>[...o.columns.map(e=>({type:`text`,width:`w-24`,headerWidth:`w-20`})),{type:`icons`,count:2,headerWidth:`w-16`}]);function b(e){l.value=e.page+1,f.value=e.rows,o.onLoad&&o.onLoad(l.value,f.value)}return n(()=>{o.onLoad&&o.onLoad(1,15)}),(n,o)=>{let l=i(`tooltip`);return t.loading?(a(),d(te,{key:0,columns:_.value,rows:8},null,8,[`columns`])):(a(),d(u(U),{key:1,value:t.items,lazy:``,totalRecords:t.total,first:p.value,rows:f.value,onPage:b,paginator:``,paginatorTemplate:`FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown`,rowsPerPageOptions:[10,15,25,50],size:`small`,class:`!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden`},{empty:r(()=>[h(n.$slots,`empty`)]),default:r(()=>[(a(!0),w(S,null,e(t.columns,e=>(a(),d(u(W),{key:e.field,field:e.field,header:e.header,sortable:``},{body:r(({data:t})=>[e.field.startsWith(`_`)?(a(),w(`span`,Ke,v(t[e.field]||`-`),1)):g(``,!0),e.html?(a(),w(`div`,{key:1,class:`editor-content`,innerHTML:t[e.field]},null,8,qe)):(a(),w(`span`,Je,v(t[e.field]||`-`),1))]),_:2},1032,[`field`,`header`]))),128)),m(u(W),{header:u(c)(`common.actions`),style:{width:`90px`},frozen:``,alignFrozen:`right`},{body:r(({data:e})=>[C(`div`,Ye,[s(m(u(D),{icon:`pi pi-pencil`,size:`small`,text:``,severity:`secondary`,onClick:t=>n.$emit(`edit`,e)},null,8,[`onClick`]),[[l,u(c)(`common.edit`),void 0,{left:!0}]]),s(m(u(D),{icon:`pi pi-trash`,size:`small`,text:``,severity:`danger`,onClick:t=>n.$emit(`delete`,e)},null,8,[`onClick`]),[[l,u(c)(`common.delete`),void 0,{left:!0}]])])]),_:1},8,[`header`])]),_:3},8,[`value`,`totalRecords`,`first`,`rows`]))}}},Ze={class:`space-y-4`},Qe={key:0,class:`bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded p-3 text-xs text-red-700 dark:text-red-300`},$e={__name:`DialogForm`,props:{visible:Boolean,title:String,saving:Boolean,errors:{type:Object,default:()=>({})},width:{type:String,default:`480px`}},emits:[`save`,`cancel`],setup(t){let n=t,{t:i}=I(),o=y(()=>n.width===`maximize`?`90vw`:n.width);return(n,s)=>(a(),d(u(F),{visible:t.visible,"onUpdate:visible":s[2]||=e=>n.$emit(`update:visible`,e),header:t.title,modal:``,style:b({width:o.value}),class:`p-fluid`,closable:!t.saving},{footer:r(()=>[m(u(D),{label:u(i)(`common.cancel`),size:`small`,outlined:``,severity:`secondary`,disabled:t.saving,onClick:s[0]||=e=>n.$emit(`cancel`)},null,8,[`label`,`disabled`]),m(u(D),{label:u(i)(`common.save`),icon:`pi pi-check`,size:`small`,loading:t.saving,onClick:s[1]||=e=>n.$emit(`save`)},null,8,[`label`,`loading`])]),default:r(()=>[C(`div`,Ze,[h(n.$slots,`default`),Object.keys(t.errors).length?(a(),w(`div`,Qe,[(a(!0),w(S,null,e(t.errors,(e,t)=>(a(),w(`p`,{key:t,class:`mb-1`},[C(`strong`,null,v(t)+`:`,1),_(` `+v(Array.isArray(e)?e.join(`, `):e),1)]))),128))])):g(``,!0)])]),_:3},8,[`visible`,`header`,`style`,`closable`]))}},et={class:`space-y-4`},tt={class:`flex items-center justify-between`},nt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},rt={class:`text-sm text-gray-500 dark:text-gray-400`},it={class:`flex flex-col items-center justify-center py-10 text-gray-400`},at={class:`text-sm font-medium`},ot=`/api/v1/tenant/job-management/responsibilities`,st={__name:`JobResponsibilitySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let n=e,i=t,{t:o}=I(),s=A(),c=x([]),l=x(!1),d=x(0),p=x(!1),h=x(!1),_=x(``),b=x(!1),T=x({}),E=x(!1),O=x(!1),k=x(``),j=x(null),N=x({main_task:``,activities:``,outputs:``,success_indicators:``}),P=y(()=>{let e=o(`job_management.responsibilities_title`);return h.value?`${e}`:`${o(`common.create`)} ${e}`}),F=y(()=>[{field:`main_task`,header:o(`job_management.main_task`),html:!0},{field:`activities`,header:o(`job_management.activities`),html:!0},{field:`outputs`,header:o(`job_management.outputs`),html:!0},{field:`success_indicators`,header:o(`job_management.success_indicators`),html:!0}]);async function R(e,t){l.value=!0;try{let r=await M.get(ot,{params:{page:e,per_page:t,organization_id:n.orgId}}),i=r.data?.data||[];c.value=i.map(e=>({...e,main_task:e.main_task,activities:e.activities,outputs:e.outputs,success_indicators:e.success_indicators})),d.value=r.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function B(){h.value=!1,_.value=``,N.value={main_task:``,activities:``,outputs:``,success_indicators:``},T.value={},p.value=!0}function V(e){h.value=!0,_.value=e.id,N.value={main_task:e.main_task||``,activities:e.activities||``,outputs:e.outputs||``,success_indicators:e.success_indicators||``},T.value={},p.value=!0}async function H(){b.value=!0,T.value={};try{let e={nomenclature:n.orgName||``,full_code:n.orgCode||``,...N.value,organization_id:n.orgId};h.value?await M.put(`${ot}/${_.value}`,e):await M.post(ot,e),p.value=!1,i(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),R(1,15)}catch(e){let t=L(e);Object.keys(t).length?T.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{b.value=!1}}function U(e){j.value=e,k.value=``,E.value=!0}async function W(){if(j.value){O.value=!0,k.value=``;try{await M.delete(`${ot}/${j.value.id}`),E.value=!1,i(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),R(1,15)}catch(e){k.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{O.value=!1}}}return(t,n)=>(a(),w(`div`,et,[C(`div`,tt,[C(`div`,null,[C(`h2`,nt,v(u(o)(`job_management.responsibilities_title`)),1),C(`p`,rt,v(u(o)(`job_management.responsibilities_description`)),1)]),m(u(D),{label:u(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>B()},null,8,[`label`])]),m(Xe,{items:c.value,loading:l.value,total:d.value,columns:F.value,entity:`responsibilities`,"org-id":e.orgId,"on-load":R,onEdit:V,onDelete:U},{empty:r(()=>[C(`div`,it,[n[9]||=C(`i`,{class:`pi pi-list-check text-3xl mb-2 opacity-50`},null,-1),C(`p`,at,v(u(o)(`job_management.empty_responsibilities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),m($e,{visible:p.value,"onUpdate:visible":n[5]||=e=>p.value=e,title:P.value,saving:b.value,errors:T.value,width:`maximize`,onSave:H,onCancel:n[6]||=e=>p.value=!1},{default:r(()=>[p.value?(a(),w(S,{key:0},[m(z,{label:u(o)(`job_management.main_task`),errors:T.value?.main_task},{default:r(()=>[m(u(We),{modelValue:N.value.main_task,"onUpdate:modelValue":n[1]||=e=>N.value.main_task=e,editorStyle:`height:120px`,class:f({"p-invalid":T.value?.main_task})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),m(z,{label:u(o)(`job_management.activities`),errors:T.value?.activities},{default:r(()=>[m(u(We),{modelValue:N.value.activities,"onUpdate:modelValue":n[2]||=e=>N.value.activities=e,editorStyle:`height:120px`,class:f({"p-invalid":T.value?.activities})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),m(z,{label:u(o)(`job_management.outputs`),errors:T.value?.outputs},{default:r(()=>[m(u(We),{modelValue:N.value.outputs,"onUpdate:modelValue":n[3]||=e=>N.value.outputs=e,editorStyle:`height:120px`,class:f({"p-invalid":T.value?.outputs})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),m(z,{label:u(o)(`job_management.success_indicators`),errors:T.value?.success_indicators},{default:r(()=>[m(u(We),{modelValue:N.value.success_indicators,"onUpdate:modelValue":n[4]||=e=>N.value.success_indicators=e,editorStyle:`height:120px`,class:f({"p-invalid":T.value?.success_indicators})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])],64)):g(``,!0)]),_:1},8,[`visible`,`title`,`saving`,`errors`]),m(J,{visible:E.value,"onUpdate:visible":n[7]||=e=>E.value=e,loading:O.value,"error-msg":k.value,onConfirm:W,onCancel:n[8]||=e=>E.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},ct={class:`space-y-4`},lt={class:`flex items-center justify-between`},ut={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},dt={class:`text-sm text-gray-500 dark:text-gray-400`},ft={class:`flex flex-col items-center justify-center py-10 text-gray-400`},pt={class:`text-sm font-medium`},mt=`/api/v1/tenant/job-management/hr-authorities`,ht={__name:`JobHRAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let n=e,i=t,{t:o}=I(),s=A(),c=x([]),l=x(!1),d=x(0),p=x(!1),h=x(!1),g=x(``),_=x(!1),b=x({}),S=x(!1),T=x(!1),E=x(``),O=x(null),k=x({description:``}),j=y(()=>{let e=o(`job_management.hr_authorities`);return h.value?`${o(`common.edit`)} ${e}`:`${o(`common.create`)} ${e}`}),N=y(()=>[{field:`description`,header:o(`job_management.description`)}]);async function P(e,t){l.value=!0;try{let r=await M.get(mt,{params:{page:e,per_page:t,organization_id:n.orgId}});c.value=r.data?.data||[],d.value=r.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function F(){h.value=!1,g.value=``,k.value={nomenclature:``,full_code:``,description:``},b.value={},p.value=!0}function R(e){h.value=!0,g.value=e.id,k.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},b.value={},p.value=!0}async function B(){_.value=!0,b.value={};try{let e={...k.value,nomenclature:n.orgName||``,full_code:n.orgCode||``,organization_id:n.orgId};h.value?await M.put(`${mt}/${g.value}`,e):await M.post(mt,e),p.value=!1,i(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),P(1,15)}catch(e){let t=L(e);Object.keys(t).length?b.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{_.value=!1}}function H(e){O.value=e,E.value=``,S.value=!0}async function U(){if(O.value){T.value=!0,E.value=``;try{await M.delete(`${mt}/${O.value.id}`),S.value=!1,i(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),P(1,15)}catch(e){E.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{T.value=!1}}}return(t,n)=>(a(),w(`div`,ct,[C(`div`,lt,[C(`div`,null,[C(`h2`,ut,v(u(o)(`job_management.hr_authorities`)),1),C(`p`,dt,v(u(o)(`job_management.authority_description`)),1)]),m(u(D),{label:u(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>F()},null,8,[`label`])]),m(Xe,{items:c.value,loading:l.value,total:d.value,columns:N.value,entity:`hr-authorities`,"org-id":e.orgId,"on-load":P,onEdit:R,onDelete:H},{empty:r(()=>[C(`div`,ft,[n[6]||=C(`i`,{class:`pi pi-users text-3xl mb-2 opacity-50`},null,-1),C(`p`,pt,v(u(o)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),m($e,{visible:p.value,"onUpdate:visible":n[2]||=e=>p.value=e,title:j.value,saving:_.value,errors:b.value,onSave:B,onCancel:n[3]||=e=>p.value=!1},{default:r(()=>[m(z,{label:u(o)(`job_management.description`),errors:b.value?.description},{default:r(()=>[m(u(V),{modelValue:k.value.description,"onUpdate:modelValue":n[1]||=e=>k.value.description=e,rows:`3`,class:f([`w-full`,{"p-invalid":b.value?.description}])},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),m(J,{visible:S.value,"onUpdate:visible":n[4]||=e=>S.value=e,loading:T.value,"error-msg":E.value,onConfirm:U,onCancel:n[5]||=e=>S.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},gt={class:`space-y-4`},_t={class:`flex items-center justify-between`},vt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},yt={class:`text-sm text-gray-500 dark:text-gray-400`},bt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},xt={class:`text-sm font-medium`},St=`/api/v1/tenant/job-management/operational-authorities`,Ct={__name:`JobOpAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let n=e,i=t,{t:o}=I(),s=A(),c=x([]),l=x(!1),d=x(0),p=x(!1),h=x(!1),g=x(``),_=x(!1),b=x({}),S=x(!1),T=x(!1),E=x(``),O=x(null),k=x({description:``}),j=y(()=>{let e=o(`job_management.op_authorities`);return h.value?`${o(`common.edit`)} ${e}`:`${o(`common.create`)} ${e}`}),N=y(()=>[{field:`description`,header:o(`job_management.description`)}]);async function P(e,t){l.value=!0;try{let r=await M.get(St,{params:{page:e,per_page:t,organization_id:n.orgId}});c.value=r.data?.data||[],d.value=r.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function F(){h.value=!1,g.value=``,k.value={nomenclature:``,full_code:``,description:``},b.value={},p.value=!0}function R(e){h.value=!0,g.value=e.id,k.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},b.value={},p.value=!0}async function B(){_.value=!0,b.value={};try{let e={...k.value,nomenclature:n.orgName||``,full_code:n.orgCode||``,organization_id:n.orgId};h.value?await M.put(`${St}/${g.value}`,e):await M.post(St,e),p.value=!1,i(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),P(1,15)}catch(e){let t=L(e);Object.keys(t).length?b.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{_.value=!1}}function H(e){O.value=e,E.value=``,S.value=!0}async function U(){if(O.value){T.value=!0,E.value=``;try{await M.delete(`${St}/${O.value.id}`),S.value=!1,i(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),P(1,15)}catch(e){E.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{T.value=!1}}}return(t,n)=>(a(),w(`div`,gt,[C(`div`,_t,[C(`div`,null,[C(`h2`,vt,v(u(o)(`job_management.op_authorities`)),1),C(`p`,yt,v(u(o)(`job_management.authority_description`)),1)]),m(u(D),{label:u(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>F()},null,8,[`label`])]),m(Xe,{items:c.value,loading:l.value,total:d.value,columns:N.value,entity:`operational-authorities`,"org-id":e.orgId,"on-load":P,onEdit:R,onDelete:H},{empty:r(()=>[C(`div`,bt,[n[6]||=C(`i`,{class:`pi pi-cog text-3xl mb-2 opacity-50`},null,-1),C(`p`,xt,v(u(o)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),m($e,{visible:p.value,"onUpdate:visible":n[2]||=e=>p.value=e,title:j.value,saving:_.value,errors:b.value,onSave:B,onCancel:n[3]||=e=>p.value=!1},{default:r(()=>[m(z,{label:u(o)(`job_management.description`),errors:b.value?.description},{default:r(()=>[m(u(V),{modelValue:k.value.description,"onUpdate:modelValue":n[1]||=e=>k.value.description=e,class:f([`w-full`,{"p-invalid":b.value?.description}]),rows:`3`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),m(J,{visible:S.value,"onUpdate:visible":n[4]||=e=>S.value=e,loading:T.value,"error-msg":E.value,onConfirm:U,onCancel:n[5]||=e=>S.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},wt={class:`space-y-4`},Tt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Et={class:`text-sm text-gray-500 dark:text-gray-400`},Dt={class:`max-w-2xl`},Ot={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},kt={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},At={class:`flex justify-end gap-2 pt-2`},jt=`/api/v1/tenant/job-management/working-activities`,Mt={__name:`JobActivitySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:t}){let i=t,o=e,{t:s}=I(),c=A(),l=x(!1),p=x(!0),h=x(``),_=x({}),y=x(``),b=x(!1),S=x(!1),T=x(``),E=x({job_management_value_id:``}),O=x([]);async function k(){try{let e=await M.get(`/api/v1/tenant/job-management/values`,{params:{type:`activity`,per_page:100}});O.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id}))}catch{}}async function j(){if(!o.orgId){p.value=!1;return}try{let e=(await M.get(jt,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];y.value=t.id,E.value.job_management_value_id=t.job_management_value_id||``}}catch{}}async function N(){h.value=``,_.value={},l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,job_management_value_id:E.value.job_management_value_id||null,organization_id:o.orgId};if(y.value)await M.put(`${jt}/${y.value}`,{job_management_value_id:E.value.job_management_value_id||``});else{let t=await M.post(jt,e);y.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),i(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(_.value=t,h.value=Object.values(t).join(`, `)):h.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}async function P(){if(y.value){S.value=!0,T.value=``;try{await M.delete(`${jt}/${y.value}`),b.value=!1,y.value=``,E.value.job_management_value_id=``,i(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{S.value=!1}}}return n(async()=>{try{await Promise.all([k(),j()])}finally{p.value=!1}}),(e,t)=>(a(),w(`div`,wt,[C(`div`,null,[C(`h2`,Tt,v(u(s)(`job_management.activities`)),1),C(`p`,Et,v(u(s)(`job_management.activity_description`)),1)]),C(`div`,Dt,[p.value?(a(),d(X,{key:0,type:`detail`,count:1,rows:3,cols:`grid-cols-1`,padding:`p-5`})):(a(),w(`div`,Ot,[m(z,{label:u(s)(`job_values.types.activity`),errors:_.value?.job_management_value_id},{default:r(()=>[m(Y,{modelValue:E.value.job_management_value_id,"onUpdate:modelValue":t[0]||=e=>E.value.job_management_value_id=e,options:O.value,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),class:f({"p-invalid":_.value?.job_management_value_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h.value?(a(),w(`div`,kt,v(h.value),1)):g(``,!0),C(`div`,At,[y.value?(a(),d(u(D),{key:0,label:u(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[1]||=e=>b.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:y.value?u(s)(`common.update`):u(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:l.value,onClick:N},null,8,[`label`,`loading`,`disabled`])])]))]),m(J,{visible:b.value,"onUpdate:visible":t[2]||=e=>b.value=e,loading:S.value,"error-msg":T.value,onConfirm:P,onCancel:t[3]||=e=>b.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Nt={class:`space-y-4`},Pt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Ft={class:`text-sm text-gray-500 dark:text-gray-400`},It={class:`max-w-2xl`},Lt={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Rt={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},zt={class:`flex justify-end gap-2 pt-2`},Bt=`/api/v1/tenant/job-management/working-risks`,Vt={__name:`JobRiskSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:t}){let i=t,o=e,{t:s}=I(),c=A(),l=x(!1),p=x(!0),h=x(``),_=x({}),y=x(``),b=x(!1),S=x(!1),T=x(``),E=x({job_management_value_environment_id:``,job_management_value_hazard_id:``}),O=x([]),k=x([]);async function j(){try{let[e,t]=await Promise.all([M.get(`/api/v1/tenant/job-management/values`,{params:{type:`environment`,per_page:100}}),M.get(`/api/v1/tenant/job-management/values`,{params:{type:`risk`,per_page:100}})]);O.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),k.value=(t.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id}))}catch{}}async function N(){if(!o.orgId){p.value=!1;return}try{let e=(await M.get(Bt,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];y.value=t.id,E.value.job_management_value_environment_id=t.job_management_value_environment_id||``,E.value.job_management_value_hazard_id=t.job_management_value_hazard_id||``}}catch{}}async function P(){h.value=``,_.value={},l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,job_management_value_environment_id:E.value.job_management_value_environment_id||null,job_management_value_hazard_id:E.value.job_management_value_hazard_id||null,organization_id:o.orgId};if(y.value)await M.put(`${Bt}/${y.value}`,{job_management_value_environment_id:E.value.job_management_value_environment_id||``,job_management_value_hazard_id:E.value.job_management_value_hazard_id||``});else{let t=await M.post(Bt,e);y.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),i(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(_.value=t,h.value=Object.values(t).join(`, `)):h.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}async function F(){if(y.value){S.value=!0,T.value=``;try{await M.delete(`${Bt}/${y.value}`),b.value=!1,y.value=``,E.value.job_management_value_environment_id=``,E.value.job_management_value_hazard_id=``,i(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{S.value=!1}}}return n(async()=>{try{await Promise.all([j(),N()])}finally{p.value=!1}}),(e,t)=>(a(),w(`div`,Nt,[C(`div`,null,[C(`h2`,Pt,v(u(s)(`job_management.risks`)),1),C(`p`,Ft,v(u(s)(`job_management.risk_description`)),1)]),C(`div`,It,[p.value?(a(),d(X,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(a(),w(`div`,Lt,[m(z,{label:u(s)(`job_management.work_environment`),errors:_.value?.job_management_value_environment_id},{default:r(()=>[m(Y,{modelValue:E.value.job_management_value_environment_id,"onUpdate:modelValue":t[0]||=e=>E.value.job_management_value_environment_id=e,options:O.value,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),class:f({"p-invalid":_.value?.job_management_value_environment_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),m(z,{label:u(s)(`job_management.risk`),errors:_.value?.job_management_value_hazard_id},{default:r(()=>[m(Y,{modelValue:E.value.job_management_value_hazard_id,"onUpdate:modelValue":t[1]||=e=>E.value.job_management_value_hazard_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),class:f({"p-invalid":_.value?.job_management_value_hazard_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h.value?(a(),w(`div`,Rt,v(h.value),1)):g(``,!0),C(`div`,zt,[y.value?(a(),d(u(D),{key:0,label:u(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[2]||=e=>b.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:y.value?u(s)(`common.update`):u(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:l.value,onClick:P},null,8,[`label`,`loading`,`disabled`])])]))]),m(J,{visible:b.value,"onUpdate:visible":t[3]||=e=>b.value=e,loading:S.value,"error-msg":T.value,onConfirm:F,onCancel:t[4]||=e=>b.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Ht={class:`space-y-4`},Ut={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Wt={class:`text-sm text-gray-500 dark:text-gray-400`},Gt={class:`max-w-2xl`},Kt={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},qt={class:`pt-1`},Jt={class:`flex items-center gap-2 mb-3`},Yt={class:`text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400`},Xt={class:`space-y-4`},Zt={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Qt={class:`flex justify-end gap-2 pt-2`},$t={class:`max-w-3xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5 space-y-4`},en={class:`flex items-center justify-between gap-2 flex-wrap`},tn={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},nn={class:`text-sm text-gray-500 dark:text-gray-400`},rn={key:0,class:`text-sm text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg px-3 py-2`},an={key:1,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},on={key:2,class:`overflow-x-auto border border-gray-200 dark:border-gray-700 rounded-lg`},sn={class:`w-full text-sm`},cn={class:`bg-gray-50 dark:bg-gray-700/40 text-left`},ln={class:`px-3 py-2 font-semibold text-gray-600 dark:text-gray-300`},un={class:`px-3 py-2 font-semibold text-gray-600 dark:text-gray-300`},dn={class:`px-3 py-2 align-top text-gray-500 dark:text-gray-400`},fn={class:`px-3 py-2`},pn={class:`px-3 py-2`},mn={class:`px-3 py-2 align-top`},hn={key:3,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},gn={key:4,class:`flex justify-end gap-2 pt-2`},$=`/api/v1/tenant/job-management/relationships`,_n={__name:`JobRelationshipSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},orgSummaryId:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(t,{emit:i}){let o=i,s=t,{t:c}=I(),l=A(),p=x(!1),h=x(!0),_=x(``),y=x({}),b=x(``),T=x(!1),E=x(!1),O=x(``),k=x({job_management_value_relationship_id:``,job_management_value_frequency_id:``}),j=x([]),N=x([]),P=x([]),F=x([]),R=x(!1),B=x(``);async function V(){try{let[e,t,n]=await Promise.all([M.get(`/api/v1/tenant/job-management/values`,{params:{type:`relationship`,per_page:100}}),M.get(`/api/v1/tenant/job-management/values`,{params:{type:`frequency`,per_page:100}}),s.orgSummaryId?M.get(`/api/v1/tenant/organizations`,{params:{summary_id:s.orgSummaryId,per_page:100}}):Promise.resolve({data:{data:[]}})]);j.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),N.value=(t.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),P.value=(n.data?.data||[]).filter(e=>e.id!==s.orgId).map(e=>({label:e.full_code?`${e.full_code} - ${e.nomenclature}`:e.nomenclature,value:e.id}))}catch{}}async function U(){if(!s.orgId){h.value=!1;return}try{let e=(await M.get($,{params:{organization_id:s.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];b.value=t.id,k.value.job_management_value_relationship_id=t.job_management_value_relationship_id||``,k.value.job_management_value_frequency_id=t.job_management_value_frequency_id||``,await ne()}}catch{}}async function W(){_.value=``,y.value={},p.value=!0;try{let e={nomenclature:s.orgName||``,full_code:s.orgCode||``,job_management_value_relationship_id:k.value.job_management_value_relationship_id||null,job_management_value_frequency_id:k.value.job_management_value_frequency_id||null,organization_id:s.orgId};if(b.value)await M.put(`${$}/${b.value}`,{job_management_value_relationship_id:k.value.job_management_value_relationship_id||``,job_management_value_frequency_id:k.value.job_management_value_frequency_id||``});else{let t=await M.post($,e);b.value=t.data?.data?.id||``}l.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),o(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(y.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{p.value=!1}}async function G(){if(b.value){E.value=!0,O.value=``;try{await M.delete(`${$}/${b.value}`),T.value=!1,b.value=``,k.value.job_management_value_relationship_id=``,k.value.job_management_value_frequency_id=``,F.value=[],o(`saved`),l.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){O.value=e?.response?.data?.error?.message||c(`message.operation_failed`)}finally{E.value=!1}}}let K=0;function q(){b.value&&F.value.push({_key:`new-${++K}`,id:``,organization_id:``,activity:``})}function ee(e){let t=F.value[e];t&&(t.id?te(t.id,e):F.value.splice(e,1))}async function te(e,t){try{await M.delete(`${$}/${b.value}/details/${e}`),F.value.splice(t,1),l.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){l.add({severity:`error`,summary:c(`message.error`),detail:e?.response?.data?.error?.message||c(`message.operation_failed`),life:4e3})}}async function ne(){if(b.value)try{let e=await M.get(`${$}/${b.value}/details`);F.value=(e.data?.data||[]).map(e=>({_key:`db-${++K}`,id:e.id,organization_id:e.organization_id||``,activity:e.activity||``}))}catch{}}async function Z(){if(!(!b.value||R.value)){B.value=``,R.value=!0;try{for(let e of F.value){let t={organization_id:e.organization_id||``,activity:e.activity||``};e.id?await M.put(`${$}/${b.value}/details/${e.id}`,t):e.id=(await M.post(`${$}/${b.value}/details`,t)).data?.data?.id||``}await ne(),l.add({severity:`success`,summary:c(`message.success`),detail:c(`job_management.relationship_details_saved`),life:2e3})}catch(e){let t=L(e);Object.keys(t).length>0?B.value=Object.values(t).join(`, `):B.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{R.value=!1}}}return n(async()=>{try{await Promise.all([V(),U()])}finally{h.value=!1}}),(t,n)=>(a(),w(`div`,Ht,[C(`div`,null,[C(`h2`,Ut,v(u(c)(`job_management.relationships`)),1),C(`p`,Wt,v(u(c)(`job_management.relationship_description`)),1)]),C(`div`,Gt,[h.value?(a(),d(X,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(a(),w(`div`,Kt,[C(`div`,qt,[C(`div`,Jt,[n[5]||=C(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400`},[C(`i`,{class:`pi pi-compass text-sm`})],-1),C(`h3`,Yt,v(u(c)(`job_management.relationship_group_scope`)),1),n[6]||=C(`div`,{class:`flex-1 border-t border-gray-200 dark:border-gray-700`},null,-1)]),C(`div`,Xt,[m(z,{label:u(c)(`job_management.relationship_type`),errors:y.value?.job_management_value_relationship_id},{default:r(()=>[m(Y,{modelValue:k.value.job_management_value_relationship_id,"onUpdate:modelValue":n[0]||=e=>k.value.job_management_value_relationship_id=e,options:j.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),class:f({"p-invalid":y.value?.job_management_value_relationship_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),m(z,{label:u(c)(`job_management.frequency`),errors:y.value?.job_management_value_frequency_id},{default:r(()=>[m(Y,{modelValue:k.value.job_management_value_frequency_id,"onUpdate:modelValue":n[1]||=e=>k.value.job_management_value_frequency_id=e,options:N.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),class:f({"p-invalid":y.value?.job_management_value_frequency_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`])])]),_.value?(a(),w(`div`,Zt,v(_.value),1)):g(``,!0),C(`div`,Qt,[b.value?(a(),d(u(D),{key:0,label:u(c)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:n[2]||=e=>T.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:b.value?u(c)(`common.update`):u(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:p.value,disabled:p.value,onClick:W},null,8,[`label`,`loading`,`disabled`])])]))]),C(`div`,$t,[C(`div`,en,[C(`div`,null,[C(`h3`,tn,v(u(c)(`job_management.relationship_details`)),1),C(`p`,nn,v(u(c)(`job_management.relationship_details_description`)),1)]),m(u(D),{label:u(c)(`job_management.add_relationship_detail`),icon:`pi pi-plus`,size:`small`,outlined:``,disabled:!b.value||R.value,onClick:q},null,8,[`label`,`disabled`])]),b.value?F.value.length===0?(a(),w(`div`,an,v(u(c)(`job_management.no_relationship_details`)),1)):g(``,!0):(a(),w(`div`,rn,v(u(c)(`job_management.save_relationship_first`)),1)),F.value.length>0?(a(),w(`div`,on,[C(`table`,sn,[C(`thead`,null,[C(`tr`,cn,[n[7]||=C(`th`,{class:`px-3 py-2 w-10 font-semibold text-gray-600 dark:text-gray-300`},`#`,-1),C(`th`,ln,v(u(c)(`job_management.relationship_organization`)),1),C(`th`,un,v(u(c)(`job_management.relationship_activity`)),1),n[8]||=C(`th`,{class:`px-3 py-2 w-12`},null,-1)])]),C(`tbody`,null,[(a(!0),w(S,null,e(F.value,(e,t)=>(a(),w(`tr`,{key:e._key,class:`border-t border-gray-200 dark:border-gray-700`},[C(`td`,dn,v(t+1),1),C(`td`,fn,[m(Y,{modelValue:e.organization_id,"onUpdate:modelValue":t=>e.organization_id=t,options:P.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),C(`td`,pn,[m(H,{modelValue:e.activity,"onUpdate:modelValue":t=>e.activity=t,placeholder:u(c)(`job_management.relationship_activity`)},null,8,[`modelValue`,`onUpdate:modelValue`,`placeholder`])]),C(`td`,mn,[m(u(D),{icon:`pi pi-trash`,severity:`danger`,size:`small`,text:``,rounded:``,"aria-label":`Remove`,onClick:e=>ee(t)},null,8,[`onClick`])])]))),128))])])])):g(``,!0),B.value?(a(),w(`div`,hn,v(B.value),1)):g(``,!0),F.value.length>0?(a(),w(`div`,gn,[m(u(D),{label:u(c)(`job_management.save_relationship_details`),icon:`pi pi-save`,size:`small`,loading:R.value,disabled:R.value||!b.value,onClick:Z},null,8,[`label`,`loading`,`disabled`])])):g(``,!0)]),m(J,{visible:T.value,"onUpdate:visible":n[3]||=e=>T.value=e,loading:E.value,"error-msg":O.value,onConfirm:G,onCancel:n[4]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},vn={class:`space-y-4`},yn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},bn={class:`text-sm text-gray-500 dark:text-gray-400`},xn={class:`max-w-2xl`},Sn={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Cn={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},wn={class:`flex justify-end gap-2 pt-2`},Tn=`/api/v1/tenant/job-management/subordinate-controls`,En={__name:`JobSubordinateSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:t}){let i=t,o=e,{t:s}=I(),c=A(),l=x(!1),p=x(!0),h=x(``),_=x({}),y=x(``),b=x(!1),S=x(!1),T=x(``),E=x({job_management_value_id:``}),O=x([]);async function k(){try{let e=await M.get(`/api/v1/tenant/job-management/values`,{params:{type:`subordinate`,per_page:100}});O.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id}))}catch{}}async function j(){if(!o.orgId){p.value=!1;return}try{let e=(await M.get(Tn,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];y.value=t.id,E.value.job_management_value_id=t.job_management_value_id||``}}catch{}}async function N(){h.value=``,_.value={},l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,job_management_value_id:E.value.job_management_value_id||null,organization_id:o.orgId};if(y.value)await M.put(`${Tn}/${y.value}`,{job_management_value_id:E.value.job_management_value_id||``});else{let t=await M.post(Tn,e);y.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),i(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(_.value=t,h.value=Object.values(t).join(`, `)):h.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}async function P(){if(y.value){S.value=!0,T.value=``;try{await M.delete(`${Tn}/${y.value}`),b.value=!1,y.value=``,E.value.job_management_value_id=``,i(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{S.value=!1}}}return n(async()=>{try{await Promise.all([k(),j()])}finally{p.value=!1}}),(e,t)=>(a(),w(`div`,vn,[C(`div`,null,[C(`h2`,yn,v(u(s)(`job_management.subordinate_controls`)),1),C(`p`,bn,v(u(s)(`job_management.subordinate_description`)),1)]),C(`div`,xn,[p.value?(a(),d(X,{key:0,type:`detail`,count:1,rows:3,cols:`grid-cols-1`,padding:`p-5`})):(a(),w(`div`,Sn,[m(z,{label:u(s)(`job_management.control_type`),errors:_.value?.job_management_value_id},{default:r(()=>[m(Y,{modelValue:E.value.job_management_value_id,"onUpdate:modelValue":t[0]||=e=>E.value.job_management_value_id=e,options:O.value,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),class:f({"p-invalid":_.value?.job_management_value_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h.value?(a(),w(`div`,Cn,v(h.value),1)):g(``,!0),C(`div`,wn,[y.value?(a(),d(u(D),{key:0,label:u(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[1]||=e=>b.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:y.value?u(s)(`common.update`):u(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:l.value,onClick:N},null,8,[`label`,`loading`,`disabled`])])]))]),m(J,{visible:b.value,"onUpdate:visible":t[2]||=e=>b.value=e,loading:S.value,"error-msg":T.value,onConfirm:P,onCancel:t[3]||=e=>b.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Dn={class:`space-y-4`},On={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},kn={class:`text-sm text-gray-500 dark:text-gray-400`},An={class:`max-w-2xl`},jn={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Mn={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Nn={class:`flex justify-end gap-2 pt-2`},Pn=`/api/v1/tenant/job-management/assets`,Fn={__name:`JobAssetSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:t}){let i=t,o=e,{t:s}=I(),c=A(),l=x(!1),p=x(!0),h=x(``),_=x({}),y=x(``),b=x(!1),S=x(!1),T=x(``),E=x({job_management_value_asset_id:``,job_management_value_authority_id:``}),O=x([]),k=x([]);async function j(){try{let[e,t]=await Promise.all([M.get(`/api/v1/tenant/job-management/values`,{params:{type:`asset`,per_page:100}}),M.get(`/api/v1/tenant/job-management/values`,{params:{type:`asset_authority`,per_page:100}})]);O.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),k.value=(t.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id}))}catch{}}async function N(){if(!o.orgId){p.value=!1;return}try{let e=(await M.get(Pn,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];y.value=t.id,E.value.job_management_value_asset_id=t.job_management_value_asset_id||``,E.value.job_management_value_authority_id=t.job_management_value_authority_id||``}}catch{}}async function P(){h.value=``,_.value={},l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,job_management_value_asset_id:E.value.job_management_value_asset_id||null,job_management_value_authority_id:E.value.job_management_value_authority_id||null,organization_id:o.orgId};if(y.value)await M.put(`${Pn}/${y.value}`,{job_management_value_asset_id:E.value.job_management_value_asset_id||``,job_management_value_authority_id:E.value.job_management_value_authority_id||``});else{let t=await M.post(Pn,e);y.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),i(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(_.value=t,h.value=Object.values(t).join(`, `)):h.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}async function F(){if(y.value){S.value=!0,T.value=``;try{await M.delete(`${Pn}/${y.value}`),b.value=!1,y.value=``,E.value.job_management_value_asset_id=``,E.value.job_management_value_authority_id=``,i(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{S.value=!1}}}return n(async()=>{try{await Promise.all([j(),N()])}finally{p.value=!1}}),(e,t)=>(a(),w(`div`,Dn,[C(`div`,null,[C(`h2`,On,v(u(s)(`job_management.assets`)),1),C(`p`,kn,v(u(s)(`job_management.asset_description`)),1)]),C(`div`,An,[p.value?(a(),d(X,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(a(),w(`div`,jn,[m(z,{label:u(s)(`job_management.asset_type`),errors:_.value?.job_management_value_asset_id},{default:r(()=>[m(Y,{modelValue:E.value.job_management_value_asset_id,"onUpdate:modelValue":t[0]||=e=>E.value.job_management_value_asset_id=e,options:O.value,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),class:f({"p-invalid":_.value?.job_management_value_asset_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),m(z,{label:u(s)(`job_management.authority_level`),errors:_.value?.job_management_value_authority_id},{default:r(()=>[m(Y,{modelValue:E.value.job_management_value_authority_id,"onUpdate:modelValue":t[1]||=e=>E.value.job_management_value_authority_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),class:f({"p-invalid":_.value?.job_management_value_authority_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h.value?(a(),w(`div`,Mn,v(h.value),1)):g(``,!0),C(`div`,Nn,[y.value?(a(),d(u(D),{key:0,label:u(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[2]||=e=>b.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:y.value?u(s)(`common.update`):u(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:l.value,onClick:P},null,8,[`label`,`loading`,`disabled`])])]))]),m(J,{visible:b.value,"onUpdate:visible":t[3]||=e=>b.value=e,loading:S.value,"error-msg":T.value,onConfirm:F,onCancel:t[4]||=e=>b.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},In={class:`space-y-4`},Ln={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Rn={class:`text-sm text-gray-500 dark:text-gray-400`},zn={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Bn={class:`flex items-center justify-between gap-4 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3`},Vn={class:`min-w-0`},Hn={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},Un={class:`text-xs text-gray-500 dark:text-gray-400 mt-0.5`},Wn={class:`space-y-4 pt-4 border-t border-gray-200 dark:border-gray-700`},Gn={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Kn={class:`flex justify-end gap-2 pt-2`},qn=`/api/v1/tenant/job-management/financials`,Jn={__name:`JobFinancialSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:t}){let i=t,o=e,{t:s}=I(),l=A(),p=x(!1),h=x(!0),_=x(``),b=x({}),S=x(``),T=x(!1),E=x(!1),O=x(``),k=x({is_authorized:!1,job_management_value_cash_id:``,job_management_value_authority_id:``,job_management_value_impact_id:``}),j=x([]),N=x([]),P=x([]),F=x([]),R=x([]),B=y(()=>k.value.is_authorized?N.value:P.value),V=y(()=>k.value.is_authorized?F.value:R.value);async function H(){try{let[e,t,n,r,i]=await Promise.all([M.get(`/api/v1/tenant/job-management/values`,{params:{type:`cash`,per_page:100}}),M.get(`/api/v1/tenant/job-management/values`,{params:{type:`authority`,per_page:100}}),M.get(`/api/v1/tenant/job-management/values`,{params:{type:`authority_unauthorized`,per_page:100}}),M.get(`/api/v1/tenant/job-management/values`,{params:{type:`impact`,per_page:100}}),M.get(`/api/v1/tenant/job-management/values`,{params:{type:`impact_unauthorized`,per_page:100}})]);j.value=(e.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),N.value=(t.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),P.value=(n.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),F.value=(r.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),R.value=(i.data?.data||[]).map(e=>({label:e.descriptions,value:e.id}))}catch{}}let U=!1;c(()=>k.value.is_authorized,(e,t)=>{U||e===t||(k.value.job_management_value_cash_id=``,k.value.job_management_value_authority_id=``,k.value.job_management_value_impact_id=``)},{flush:`sync`});async function W(){if(!o.orgId){h.value=!1;return}try{let e=(await M.get(qn,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];U=!0,S.value=t.id,k.value.is_authorized=!!t.is_authorized,k.value.job_management_value_cash_id=t.job_management_value_cash_id||``,k.value.job_management_value_authority_id=t.job_management_value_authority_id||``,k.value.job_management_value_impact_id=t.job_management_value_impact_id||``,U=!1}}catch{}}async function G(){_.value=``,b.value={},p.value=!0;try{let e=!!k.value.is_authorized,t={nomenclature:o.orgName||``,full_code:o.orgCode||``,is_authorized:e,job_management_value_cash_id:e&&k.value.job_management_value_cash_id||null,job_management_value_authority_id:k.value.job_management_value_authority_id||null,job_management_value_impact_id:k.value.job_management_value_impact_id||null,organization_id:o.orgId};if(S.value)await M.put(`${qn}/${S.value}`,{is_authorized:e,job_management_value_cash_id:e&&k.value.job_management_value_cash_id||``,job_management_value_authority_id:k.value.job_management_value_authority_id||``,job_management_value_impact_id:k.value.job_management_value_impact_id||``});else{let e=await M.post(qn,t);S.value=e.data?.data?.id||``}l.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),i(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(b.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{p.value=!1}}async function K(){if(S.value){E.value=!0,O.value=``;try{await M.delete(`${qn}/${S.value}`),T.value=!1,S.value=``,k.value.is_authorized=!1,k.value.job_management_value_cash_id=``,k.value.job_management_value_authority_id=``,k.value.job_management_value_impact_id=``,i(`saved`),l.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){O.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{E.value=!1}}}return n(async()=>{try{await Promise.all([H(),W()])}finally{h.value=!1}}),(e,t)=>(a(),w(`div`,In,[C(`div`,null,[C(`h2`,Ln,v(u(s)(`job_management.financials`)),1),C(`p`,Rn,v(u(s)(`job_management.financial_description`)),1)]),C(`div`,null,[h.value?(a(),d(X,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(a(),w(`div`,zn,[C(`div`,Bn,[C(`div`,Vn,[C(`p`,Hn,v(u(s)(`job_management.is_authorized`)),1),C(`p`,Un,v(u(s)(`job_management.is_authorized_description`)),1)]),m(u(ee),{modelValue:k.value.is_authorized,"onUpdate:modelValue":t[0]||=e=>k.value.is_authorized=e},null,8,[`modelValue`])]),C(`div`,Wn,[k.value.is_authorized?(a(),d(z,{key:0,label:u(s)(`job_management.cash_level`),errors:b.value?.job_management_value_cash_id},{default:r(()=>[m(Y,{modelValue:k.value.job_management_value_cash_id,"onUpdate:modelValue":t[1]||=e=>k.value.job_management_value_cash_id=e,options:j.value,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),class:f({"p-invalid":b.value?.job_management_value_cash_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`])):g(``,!0),m(z,{label:u(s)(`job_management.authority_level`),errors:b.value?.job_management_value_authority_id},{default:r(()=>[m(Y,{modelValue:k.value.job_management_value_authority_id,"onUpdate:modelValue":t[2]||=e=>k.value.job_management_value_authority_id=e,options:B.value,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),class:f({"p-invalid":b.value?.job_management_value_authority_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),m(z,{label:u(s)(`job_management.impact_level`),errors:b.value?.job_management_value_impact_id},{default:r(()=>[m(Y,{modelValue:k.value.job_management_value_impact_id,"onUpdate:modelValue":t[3]||=e=>k.value.job_management_value_impact_id=e,options:V.value,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),class:f({"p-invalid":b.value?.job_management_value_impact_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`])]),_.value?(a(),w(`div`,Gn,v(_.value),1)):g(``,!0),C(`div`,Kn,[S.value?(a(),d(u(D),{key:0,label:u(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[4]||=e=>T.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:S.value?u(s)(`common.update`):u(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:p.value,disabled:p.value,onClick:G},null,8,[`label`,`loading`,`disabled`])])]))]),m(J,{visible:T.value,"onUpdate:visible":t[5]||=e=>T.value=e,loading:E.value,"error-msg":O.value,onConfirm:K,onCancel:t[6]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Yn=`/api/v1/tenant/job-management/potency-competencies`;function Xn({orgId:e,rows:t,afterDelete:n,onSaved:r,matchBy:i=`value`,descriptionField:a=`descriptions`}){let{t:o}=I(),s=A(),c=x(!1),l=x(``),u=x(!1),d=x(!1),f=x(``),p=x(null),m=x([]);function h(e){let t=(e.levelOptions||[]).find(t=>t.value===e.job_management_value_id);return t&&t[a]||``}function g(e){if(i===`competency`)return e.competency_id&&m.value.find(t=>t.competency_id&&t.competency_id===e.competency_id)||null;let t=new Set((e.levelOptions||[]).map(e=>e.value));return m.value.find(e=>e.job_management_value_id&&t.has(e.job_management_value_id))||null}function _(){t.value.forEach(e=>{let t=g(e);e.recordId=t?t.id:``,e.job_management_value_id=t&&t.job_management_value_id||``,e.weight!==void 0&&(e.weight=t?t.weight??e.weight:e.weight)})}async function v(){if(!e.value){m.value=[];return}try{let t=await M.get(Yn,{params:{organization_id:e.value,per_page:100}});m.value=t.data?.data||[]}catch{m.value=[]}}function y(e){p.value=e,f.value=``,u.value=!0}async function b(){let e=p.value;if(e){d.value=!0,f.value=``;try{e.recordId&&await M.delete(`${Yn}/${e.recordId}`),n&&n(e),u.value=!1,await v(),_(),s.add({severity:`success`,summary:o(`message.success`),detail:o(`message.deleted`),life:2e3}),r&&r()}catch(e){f.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{d.value=!1,p.value=null}}}async function S(){l.value=``,c.value=!0;try{for(let n of t.value)if(n.job_management_value_id){let t=n.competency_id?{competency_id:n.competency_id,job_management_value_id:n.job_management_value_id}:{job_management_value_id:n.job_management_value_id};n.weight!==void 0&&n.weight!==null&&n.weight!==``&&(t.weight=n.weight),n.recordId?await M.put(`${Yn}/${n.recordId}`,t):n.recordId=(await M.post(Yn,{organization_id:e.value,...t})).data?.data?.id||``}else n.recordId&&=(await M.delete(`${Yn}/${n.recordId}`),``);await v(),_(),s.add({severity:`success`,summary:o(`message.success`),detail:o(`common.saved`),life:2e3}),r&&r()}catch(e){let t=L(e);Object.keys(t).length>0?l.value=Object.values(t).join(`, `):l.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{c.value=!1}}return{savingCard:c,errorMsg:l,deleteVisible:u,deleting:d,deleteError:f,deleteTarget:p,records:m,levelDescription:h,hydrateRows:_,loadData:v,askDeleteRow:y,handleDelete:b,handleSave:S}}var Zn={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Qn={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},$n={class:`text-sm text-gray-500 dark:text-gray-400`},er={key:0,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},tr={key:1,class:`overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700`},nr={class:`w-full text-sm`},rr={class:`bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400`},ir={class:`px-4 py-3 font-semibold min-w-[220px]`},ar={class:`px-4 py-3 font-semibold min-w-[260px]`},or={class:`px-4 py-3 font-semibold min-w-[260px]`},sr={class:`px-4 py-3 font-semibold w-16 text-right`},cr={class:`px-4 py-3`},lr={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},ur={class:`px-4 py-3`},dr={class:`px-4 py-3 text-sm text-gray-600 dark:text-gray-300`},fr={class:`px-4 py-3 text-right`},pr={key:2,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},mr={key:3,class:`flex justify-end gap-2 pt-1`},hr={__name:`SelectablePotencyCard`,props:{orgId:String,typeGroup:{type:String,required:!0},skeletonRows:{type:Number,default:5},titleKey:{type:String,required:!0},descriptionKey:{type:String,required:!0},emptyKey:{type:String,required:!0},saveLabelKey:{type:String,required:!0},deleteTitleKey:{type:String,required:!0},deleteMessageKey:{type:String,required:!0}},emits:[`saved`],setup(t,{emit:r}){let i=r,o=t,{t:s}=I(),f=x(!0),p=x([]),h=x([]),_=x([]),{savingCard:b,errorMsg:T,deleteVisible:E,deleting:O,deleteError:k,deleteTarget:A,records:j,levelDescription:N,hydrateRows:P,loadData:F,askDeleteRow:L,handleDelete:R,handleSave:z}=Xn({orgId:y(()=>o.orgId),rows:h,afterDelete:e=>{let t=Array.isArray(_.value)?_.value:[];_.value=t.filter(t=>t!==e.type)},onSaved:()=>i(`saved`)}),B=y(()=>(p.value||[]).find(e=>e.type_group===o.typeGroup));function V(e){let t=`job_values.types.${e.type}`,n=s(t);return n===t?e.description_group||e.type:n}let H=y(()=>(B.value?.types||[]).map(e=>({label:V(e),value:e.type})));function U(){let e={};(B.value?.types||[]).forEach(t=>{e[t.type]=t});let t=Array.isArray(_.value)?_.value:_.value?[_.value]:[];h.value=t.filter(t=>e[t]).map(t=>{let n=e[t];return{competency_id:``,competency_name:V(n),competency_definition:``,type:n.type,levelOptions:(n.options||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``})),recordId:``,job_management_value_id:``}})}async function W(){try{let e=await M.get(`/api/v1/tenant/job-management/values/tree`);p.value=e.data?.data||[],U()}catch{p.value=[],h.value=[]}}function G(){let e={};(B.value?.types||[]).forEach(t=>{(t.options||[]).forEach(n=>{e[n.id]=t.type})});let t=[];j.value.forEach(n=>{let r=n.job_management_value_id&&e[n.job_management_value_id];r&&!t.includes(r)&&t.push(r)}),_.value=t,U(),P()}return c(_,()=>{U(),P()}),n(async()=>{try{await Promise.all([W(),F()])}finally{G(),f.value=!1}}),(n,r)=>(a(),w(`div`,Zn,[C(`div`,null,[C(`h3`,Qn,v(u(s)(t.titleKey)),1),C(`p`,$n,v(u(s)(t.descriptionKey)),1)]),f.value?(a(),d(X,{key:0,type:`detail`,count:1,rows:t.skeletonRows,cols:`grid-cols-1`,padding:`p-5`},null,8,[`rows`])):(a(),w(S,{key:1},[m(Y,{modelValue:_.value,"onUpdate:modelValue":r[0]||=e=>_.value=e,options:H.value,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),showClear:``,multiple:``},null,8,[`modelValue`,`options`,`placeholder`]),h.value.length===0?(a(),w(`div`,er,v(u(s)(t.emptyKey)),1)):(a(),w(`div`,tr,[C(`table`,nr,[C(`thead`,null,[C(`tr`,rr,[C(`th`,ir,v(u(s)(`job_management.potency_table_name`)),1),C(`th`,ar,v(u(s)(`job_management.potency_table_level`)),1),C(`th`,or,v(u(s)(`job_management.potency_table_description`)),1),C(`th`,sr,v(u(s)(`common.actions`)),1)])]),C(`tbody`,null,[(a(!0),w(S,null,e(h.value,e=>(a(),w(`tr`,{key:e.type,class:`border-t border-gray-100 dark:border-gray-700 align-top`},[C(`td`,cr,[C(`div`,lr,v(e.competency_name),1)]),C(`td`,ur,[m(Y,{modelValue:e.job_management_value_id,"onUpdate:modelValue":t=>e.job_management_value_id=t,options:e.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),C(`td`,dr,v(u(N)(e)),1),C(`td`,fr,[m(u(D),{icon:`pi pi-trash`,severity:`danger`,text:``,rounded:``,size:`small`,disabled:u(b),"aria-label":u(s)(`common.delete`),onClick:t=>u(L)(e)},null,8,[`disabled`,`aria-label`,`onClick`])])]))),128))])])])),u(T)?(a(),w(`div`,pr,v(u(T)),1)):g(``,!0),h.value.length>0?(a(),w(`div`,mr,[m(u(D),{label:u(s)(t.saveLabelKey),icon:`pi pi-check`,size:`small`,loading:u(b),disabled:u(b)||!t.orgId,onClick:u(z)},null,8,[`label`,`loading`,`disabled`,`onClick`])])):g(``,!0)],64)),m(J,{visible:u(E),"onUpdate:visible":r[1]||=e=>l(E)?E.value=e:null,title:u(s)(t.deleteTitleKey),message:u(s)(t.deleteMessageKey,{name:u(A)?.competency_name||``}),loading:u(O),"error-msg":u(k),onConfirm:u(R),onCancel:r[2]||=e=>E.value=!1},null,8,[`visible`,`title`,`message`,`loading`,`error-msg`,`onConfirm`])]))}},gr={__name:`PsychologicalPotencyCard`,props:{orgId:String},emits:[`saved`],setup(e,{emit:t}){let n=t;return(t,r)=>(a(),d(hr,{"org-id":e.orgId,"type-group":`psychological`,"skeleton-rows":5,"title-key":`job_management.potency_required_title`,"description-key":`job_management.potency_required_description`,"empty-key":`job_management.potency_required_empty`,"save-label-key":`job_management.save_potency_levels`,"delete-title-key":`job_management.potency_confirm_delete_title`,"delete-message-key":`job_management.potency_confirm_delete`,onSaved:r[0]||=e=>n(`saved`)},null,8,[`org-id`]))}},_r={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},vr={class:`flex items-start justify-between gap-4`},yr={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},br={class:`text-sm text-gray-500 dark:text-gray-400`},xr={class:`flex flex-col items-end gap-1 shrink-0`},Sr={class:`flex items-center gap-2`},Cr={class:`text-xs font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap`},wr={class:`w-24 shrink-0`},Tr={key:0,class:`text-xs text-red-500 dark:text-red-400`},Er={key:0,class:`text-sm text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg px-3 py-2`},Dr={key:1,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},Or={key:2,class:`overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700`},kr={class:`w-full text-sm`},Ar={class:`bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400`},jr={class:`px-4 py-3 font-semibold min-w-[220px]`},Mr={class:`px-4 py-3 font-semibold min-w-[260px]`},Nr={class:`px-4 py-3 font-semibold min-w-[130px]`},Pr={class:`px-4 py-3 font-semibold min-w-[260px]`},Fr={class:`px-4 py-3 font-semibold w-16 text-right`},Ir={class:`px-4 py-3`},Lr={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},Rr={key:0,class:`text-xs text-gray-400 dark:text-gray-500`},zr={class:`px-4 py-3`},Br={class:`px-4 py-3`},Vr={class:`px-4 py-3 text-sm text-gray-600 dark:text-gray-300`},Hr={class:`px-4 py-3 text-right`},Ur={key:3,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Wr={key:4,class:`flex justify-end gap-2 pt-1`},Gr={__name:`TechnicalPotencyCard`,props:{orgId:String},emits:[`saved`,`weight-saved`],setup(t,{emit:r}){let i=r,o=t,{t:s}=I(),f=A(),p=x(!0),h=x([]),_=x([]),b=x([]),T=x([]),E=x([]),O=y(()=>E.value.length>0),k=x(``),j=x(``),N=x(``),P=x(!1),F=x(``),{savingCard:L,errorMsg:R,deleteVisible:z,deleting:B,deleteError:V,deleteTarget:H,records:U,levelDescription:W,hydrateRows:G,loadData:q,askDeleteRow:ee,handleDelete:te,handleSave:ne}=Xn({orgId:y(()=>o.orgId),rows:b,matchBy:`competency`,descriptionField:`note`,afterDelete:e=>{let t=Array.isArray(T.value)?T.value:[];T.value=t.filter(t=>t!==e.competency_id)},onSaved:()=>i(`saved`)}),Z=y(()=>(h.value||[]).map(e=>({id:e.id,name:e.name,cluster:e.cluster||``}))),re=y(()=>{let e={};return(Z.value||[]).forEach(t=>{(e[t.cluster]=e[t.cluster]||[]).push(t)}),Object.keys(e).sort().map(t=>({label:t,items:e[t].sort((e,t)=>e.name.localeCompare(t.name))}))}),ie=y(()=>(_.value||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``,note:e.note||``})));function ae(){let e={};(Z.value||[]).forEach(t=>{e[t.id]=t});let t=(Array.isArray(T.value)?T.value:T.value?[T.value]:[]).filter(t=>e[t]),n=t.length>0?Math.round(100/t.length*100)/100:0;b.value=t.map(t=>{let r=e[t];return{competency_id:t,competency_name:r.name,cluster:r.cluster,levelOptions:ie.value,recordId:``,job_management_value_id:``,weight:n}})}async function oe(){try{let[e,t]=await Promise.all([M.get(`/api/v1/tenant/settings/competencies`,{params:{per_page:500}}),M.get(`/api/v1/tenant/job-management/values/clusters/technical`)]);E.value=t.data?.data?.clusters||[];let n=new Set(E.value),r=e.data?.data||[];h.value=r.filter(e=>e.cluster&&n.has(e.cluster))}catch{h.value=[]}}async function se(){try{let e=await M.get(`/api/v1/tenant/job-management/values`,{params:{type:`technical`,per_page:100}});_.value=e.data?.data||[]}catch{_.value=[]}}async function Q(){if(!o.orgId){k.value=``,j.value=``;return}try{let e=((await M.get(`/api/v1/tenant/job-management/competency-groups`,{params:{organization_id:o.orgId}})).data?.data||[]).find(e=>e.category===`technical`);j.value=e?e.id:``,k.value=e?e.weight:``,N.value=k.value}catch{k.value=``,j.value=``,N.value=``}}async function ce(){if(k.value===``||k.value===null||k.value===void 0){F.value=s(`job_management.potency_technical_weight_required`);return}if(!(N.value!==``&&k.value===N.value)){P.value=!0,F.value=``;try{let e=j.value;if(!e)try{let t=((await M.get(`/api/v1/tenant/job-management/competency-groups`,{params:{organization_id:o.orgId}})).data?.data||[]).find(e=>e.category===`technical`);t&&(e=t.id)}catch{}let t={weight:k.value};e?await M.put(`/api/v1/tenant/job-management/competency-groups/${e}`,t):await M.post(`/api/v1/tenant/job-management/competency-groups`,{organization_id:o.orgId,category:`technical`,weight:k.value}),N.value=k.value,f.add({severity:`success`,summary:s(`message.success`),detail:s(`job_management.potency_technical_weight_saved`),life:2e3}),i(`saved`),i(`weight-saved`,k.value),await Q()}catch(e){F.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{P.value=!1}}}function le(){let e={};(Z.value||[]).forEach(t=>{e[t.id]=t});let t=[];U.value.forEach(n=>{n.competency_id&&e[n.competency_id]&&!t.includes(n.competency_id)&&t.push(n.competency_id)}),T.value=t,ae(),G()}return c(T,()=>{ae(),G()}),n(async()=>{try{await Promise.all([oe(),se(),q(),Q()])}finally{le(),p.value=!1}}),(n,r)=>(a(),w(`div`,_r,[C(`div`,vr,[C(`div`,null,[C(`h3`,yr,v(u(s)(`job_management.potency_technical_title`)),1),C(`p`,br,v(u(s)(`job_management.potency_technical_description`)),1)]),C(`div`,xr,[C(`div`,Sr,[C(`label`,Cr,v(u(s)(`job_management.potency_technical_weight_label`)),1),C(`div`,wr,[m(u(K),{modelValue:k.value,"onUpdate:modelValue":r[0]||=e=>k.value=e,fluid:``,min:0,max:100,suffix:`%`,size:`small`,disabled:P.value||!t.orgId,onBlur:ce},null,8,[`modelValue`,`disabled`])])]),F.value?(a(),w(`div`,Tr,v(F.value),1)):g(``,!0)])]),p.value?(a(),d(X,{key:0,type:`detail`,count:1,rows:8,cols:`grid-cols-1`,padding:`p-5`})):(a(),w(S,{key:1},[m(Y,{modelValue:T.value,"onUpdate:modelValue":r[1]||=e=>T.value=e,options:re.value,"option-label":`name`,"option-value":`id`,"option-group-label":`label`,"option-group-children":`items`,placeholder:u(s)(`job_management.potency_technical_placeholder`),showClear:``,multiple:``},null,8,[`modelValue`,`options`,`placeholder`]),O.value?b.value.length===0?(a(),w(`div`,Dr,v(u(s)(`job_management.potency_technical_empty`)),1)):(a(),w(`div`,Or,[C(`table`,kr,[C(`thead`,null,[C(`tr`,Ar,[C(`th`,jr,v(u(s)(`job_management.potency_table_name`)),1),C(`th`,Mr,v(u(s)(`job_management.potency_table_level`)),1),C(`th`,Nr,v(u(s)(`job_management.potency_table_weight`)),1),C(`th`,Pr,v(u(s)(`job_management.potency_table_description`)),1),C(`th`,Fr,v(u(s)(`common.actions`)),1)])]),C(`tbody`,null,[(a(!0),w(S,null,e(b.value,e=>(a(),w(`tr`,{key:e.competency_id,class:`border-t border-gray-100 dark:border-gray-700 align-top`},[C(`td`,Ir,[C(`div`,Lr,v(e.competency_name),1),e.cluster?(a(),w(`div`,Rr,v(e.cluster),1)):g(``,!0)]),C(`td`,zr,[m(Y,{modelValue:e.job_management_value_id,"onUpdate:modelValue":t=>e.job_management_value_id=t,options:e.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),C(`td`,Br,[m(u(K),{modelValue:e.weight,"onUpdate:modelValue":t=>e.weight=t,class:`!w-full`,min:0,max:100,suffix:`%`,size:`small`},null,8,[`modelValue`,`onUpdate:modelValue`])]),C(`td`,Vr,v(u(W)(e)),1),C(`td`,Hr,[m(u(D),{icon:`pi pi-trash`,severity:`danger`,text:``,rounded:``,size:`small`,disabled:u(L),"aria-label":u(s)(`common.delete`),onClick:t=>u(ee)(e)},null,8,[`disabled`,`aria-label`,`onClick`])])]))),128))])])])):(a(),w(`div`,Er,v(u(s)(`job_management.potency_technical_no_mapping`)),1)),u(R)?(a(),w(`div`,Ur,v(u(R)),1)):g(``,!0),b.value.length>0?(a(),w(`div`,Wr,[m(u(D),{label:u(s)(`job_management.save_technical`),icon:`pi pi-check`,size:`small`,loading:u(L),disabled:u(L)||!t.orgId,onClick:u(ne)},null,8,[`label`,`loading`,`disabled`,`onClick`])])):g(``,!0)],64)),m(J,{visible:u(z),"onUpdate:visible":r[2]||=e=>l(z)?z.value=e:null,title:u(s)(`job_management.potency_confirm_delete_title`),message:u(s)(`job_management.potency_confirm_delete`,{name:u(H)?.competency_name||``}),loading:u(B),"error-msg":u(V),onConfirm:u(te),onCancel:r[3]||=e=>z.value=!1},null,8,[`visible`,`title`,`message`,`loading`,`error-msg`,`onConfirm`])]))}},Kr={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},qr={class:`flex items-start justify-between gap-4`},Jr={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},Yr={class:`text-sm text-gray-500 dark:text-gray-400`},Xr={class:`flex flex-col items-end gap-1 shrink-0`},Zr={class:`flex items-center gap-2`},Qr={class:`text-xs font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap`},$r={class:`w-24 shrink-0 text-right`},ei={key:0,class:`text-sm font-semibold text-gray-700 dark:text-gray-200`},ti={key:1,class:`text-sm text-gray-400 dark:text-gray-500`},ni={key:0,class:`pi pi-spin pi-spinner text-sm text-gray-400`},ri={key:0,class:`text-xs text-red-500 dark:text-red-400`},ii={key:0,class:`text-sm text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg px-3 py-2`},ai={key:1,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},oi={key:2,class:`overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700`},si={class:`w-full text-sm`},ci={class:`bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400`},li={class:`px-4 py-3 font-semibold min-w-[220px]`},ui={class:`px-4 py-3 font-semibold min-w-[260px]`},di={class:`px-4 py-3 font-semibold min-w-[130px]`},fi={class:`px-4 py-3 font-semibold min-w-[260px]`},pi={class:`px-4 py-3 font-semibold w-16 text-right`},mi={class:`px-4 py-3`},hi={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},gi={key:0,class:`text-xs text-gray-400 dark:text-gray-500`},_i={class:`px-4 py-3`},vi={class:`px-4 py-3`},yi={class:`px-4 py-3 text-sm text-gray-600 dark:text-gray-300`},bi={class:`px-4 py-3 text-right`},xi={key:3,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Si={key:4,class:`flex justify-end gap-2 pt-1`},Ci={__name:`ManagerialPotencyCard`,props:{orgId:String,technicalWeight:{type:Number,default:null}},emits:[`saved`],setup(t,{emit:r}){let i=r,o=t,{t:s}=I(),f=A(),p=x(!0),h=x([]),_=x([]),b=x([]),T=x([]),E=x([]),O=y(()=>E.value.length>0),k=x(``),j=x(``),N=x(``),P=x(!1),F=x(``),L=y(()=>{let e=k.value;return e===``||e==null?null:Math.round((100-e)*100)/100}),{savingCard:R,errorMsg:z,deleteVisible:B,deleting:V,deleteError:H,deleteTarget:U,records:W,levelDescription:G,hydrateRows:q,loadData:ee,askDeleteRow:te,handleDelete:ne,handleSave:Z}=Xn({orgId:y(()=>o.orgId),rows:b,matchBy:`competency`,descriptionField:`note`,afterDelete:e=>{let t=Array.isArray(T.value)?T.value:[];T.value=t.filter(t=>t!==e.competency_id)},onSaved:()=>i(`saved`)}),re=y(()=>(h.value||[]).map(e=>({id:e.id,name:e.name,cluster:e.cluster||``}))),ie=y(()=>{let e={};return(re.value||[]).forEach(t=>{(e[t.cluster]=e[t.cluster]||[]).push(t)}),Object.keys(e).sort().map(t=>({label:t,items:e[t].sort((e,t)=>e.name.localeCompare(t.name))}))}),ae=y(()=>(_.value||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``,note:e.note||``})));function oe(){let e={};(re.value||[]).forEach(t=>{e[t.id]=t});let t=(Array.isArray(T.value)?T.value:T.value?[T.value]:[]).filter(t=>e[t]),n=t.length>0?Math.round(100/t.length*100)/100:0;b.value=t.map(t=>{let r=e[t];return{competency_id:t,competency_name:r.name,cluster:r.cluster,levelOptions:ae.value,recordId:``,job_management_value_id:``,weight:n}})}async function se(){try{let[e,t]=await Promise.all([M.get(`/api/v1/tenant/settings/competencies`,{params:{per_page:500}}),M.get(`/api/v1/tenant/job-management/values/clusters/managerial`)]);E.value=t.data?.data?.clusters||[];let n=new Set(E.value),r=e.data?.data||[];h.value=r.filter(e=>e.cluster&&n.has(e.cluster))}catch{h.value=[]}}async function Q(){try{let e=await M.get(`/api/v1/tenant/job-management/values`,{params:{type:`managerial`,per_page:100}});_.value=e.data?.data||[]}catch{_.value=[]}}async function ce(){if(!o.orgId){k.value=``,j.value=``;return}try{let e=(await M.get(`/api/v1/tenant/job-management/competency-groups`,{params:{organization_id:o.orgId}})).data?.data||[],t=e.find(e=>e.category===`technical`),n=e.find(e=>e.category===`managerial`);k.value=t?t.weight:``,j.value=n?n.id:``,N.value=n?n.weight:``}catch{k.value=``,j.value=``,N.value=``}}async function le({silent:e=!1}={}){let t=L.value;if(!(t===null||!o.orgId)){P.value=!0,F.value=``;try{let n=j.value;if(!n)try{let e=((await M.get(`/api/v1/tenant/job-management/competency-groups`,{params:{organization_id:o.orgId}})).data?.data||[]).find(e=>e.category===`managerial`);e&&(n=e.id)}catch{}let r={weight:t};n?await M.put(`/api/v1/tenant/job-management/competency-groups/${n}`,r):await M.post(`/api/v1/tenant/job-management/competency-groups`,{organization_id:o.orgId,category:`managerial`,weight:t}),j.value=n||``,N.value=t,e||f.add({severity:`success`,summary:s(`message.success`),detail:s(`job_management.potency_managerial_weight_saved`),life:2e3}),i(`saved`)}catch(e){F.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{P.value=!1}}}function ue(){let e={};(re.value||[]).forEach(t=>{e[t.id]=t});let t=[];W.value.forEach(n=>{n.competency_id&&e[n.competency_id]&&!t.includes(n.competency_id)&&t.push(n.competency_id)}),T.value=t,oe(),q()}return c(T,()=>{oe(),q()}),c(()=>o.technicalWeight,e=>{e!=null&&e!==``&&(k.value=e,le())}),n(async()=>{try{await Promise.all([se(),Q(),ee(),ce()])}finally{ue(),p.value=!1;let e=L.value;if(e!==null&&o.orgId){let t=N.value;(t===``||Math.abs(t-e)>.005)&&le({silent:!0})}}}),(n,r)=>(a(),w(`div`,Kr,[C(`div`,qr,[C(`div`,null,[C(`h3`,Jr,v(u(s)(`job_management.potency_managerial_title`)),1),C(`p`,Yr,v(u(s)(`job_management.potency_managerial_description`)),1)]),C(`div`,Xr,[C(`div`,Zr,[C(`label`,Qr,v(u(s)(`job_management.potency_managerial_weight_label`)),1),C(`div`,$r,[L.value===null?(a(),w(`span`,ti,`—`)):(a(),w(`span`,ei,v(L.value)+`%`,1))]),P.value?(a(),w(`i`,ni)):g(``,!0)]),F.value?(a(),w(`div`,ri,v(F.value),1)):g(``,!0)])]),p.value?(a(),d(X,{key:0,type:`detail`,count:1,rows:5,cols:`grid-cols-1`,padding:`p-5`})):(a(),w(S,{key:1},[m(Y,{modelValue:T.value,"onUpdate:modelValue":r[0]||=e=>T.value=e,options:ie.value,"option-label":`name`,"option-value":`id`,"option-group-label":`label`,"option-group-children":`items`,placeholder:u(s)(`job_management.potency_managerial_placeholder`),showClear:``,multiple:``},null,8,[`modelValue`,`options`,`placeholder`]),O.value?b.value.length===0?(a(),w(`div`,ai,v(u(s)(`job_management.potency_managerial_empty`)),1)):(a(),w(`div`,oi,[C(`table`,si,[C(`thead`,null,[C(`tr`,ci,[C(`th`,li,v(u(s)(`job_management.potency_table_name`)),1),C(`th`,ui,v(u(s)(`job_management.potency_table_level`)),1),C(`th`,di,v(u(s)(`job_management.potency_table_weight`)),1),C(`th`,fi,v(u(s)(`job_management.potency_table_description`)),1),C(`th`,pi,v(u(s)(`common.actions`)),1)])]),C(`tbody`,null,[(a(!0),w(S,null,e(b.value,e=>(a(),w(`tr`,{key:e.competency_id,class:`border-t border-gray-100 dark:border-gray-700 align-top`},[C(`td`,mi,[C(`div`,hi,v(e.competency_name),1),e.cluster?(a(),w(`div`,gi,v(e.cluster),1)):g(``,!0)]),C(`td`,_i,[m(Y,{modelValue:e.job_management_value_id,"onUpdate:modelValue":t=>e.job_management_value_id=t,options:e.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:u(s)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),C(`td`,vi,[m(u(K),{modelValue:e.weight,"onUpdate:modelValue":t=>e.weight=t,class:`!w-full`,min:0,max:100,suffix:`%`,size:`small`},null,8,[`modelValue`,`onUpdate:modelValue`])]),C(`td`,yi,v(u(G)(e)),1),C(`td`,bi,[m(u(D),{icon:`pi pi-trash`,severity:`danger`,text:``,rounded:``,size:`small`,disabled:u(R),"aria-label":u(s)(`common.delete`),onClick:t=>u(te)(e)},null,8,[`disabled`,`aria-label`,`onClick`])])]))),128))])])])):(a(),w(`div`,ii,v(u(s)(`job_management.potency_managerial_no_mapping`)),1)),u(z)?(a(),w(`div`,xi,v(u(z)),1)):g(``,!0),b.value.length>0?(a(),w(`div`,Si,[m(u(D),{label:u(s)(`job_management.save_managerial`),icon:`pi pi-check`,size:`small`,loading:u(R),disabled:u(R)||!t.orgId,onClick:u(Z)},null,8,[`label`,`loading`,`disabled`,`onClick`])])):g(``,!0)],64)),m(J,{visible:u(B),"onUpdate:visible":r[1]||=e=>l(B)?B.value=e:null,title:u(s)(`job_management.potency_confirm_delete_title`),message:u(s)(`job_management.potency_confirm_delete`,{name:u(U)?.competency_name||``}),loading:u(V),"error-msg":u(H),onConfirm:u(ne),onCancel:r[2]||=e=>B.value=!1},null,8,[`visible`,`title`,`message`,`loading`,`error-msg`,`onConfirm`])]))}},wi={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Ti={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},Ei={class:`text-sm text-gray-500 dark:text-gray-400`},Di={key:0,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},Oi={key:1,class:`overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700`},ki={class:`w-full text-sm`},Ai={class:`bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400`},ji={class:`px-4 py-3 font-semibold min-w-[220px]`},Mi={class:`px-4 py-3 font-semibold min-w-[260px]`},Ni={class:`px-4 py-3 font-semibold min-w-[260px]`},Pi={class:`px-4 py-3 font-semibold w-16 text-right`},Fi={class:`px-4 py-3`},Ii={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},Li={key:0,class:`mt-0.5 text-xs text-gray-500 dark:text-gray-400 leading-relaxed`},Ri={class:`px-4 py-3`},zi={class:`px-4 py-3 text-sm text-gray-600 dark:text-gray-300`},Bi={class:`px-4 py-3 text-right`},Vi={key:2,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Hi={key:3,class:`flex justify-end gap-2 pt-1`},Ui={__name:`PotencyLevelsCard`,props:{orgId:String,rows:{type:Array,default:()=>[]},optionsReady:{type:Boolean,default:!1},skeletonRows:{type:Number,default:2},titleKey:{type:String,required:!0},descriptionKey:{type:String,required:!0},emptyKey:{type:String,required:!0},saveLabelKey:{type:String,required:!0},deleteTitleKey:{type:String,required:!0},deleteMessageKey:{type:String,required:!0}},emits:[`saved`],setup(t,{emit:n}){let r=n,i=t,{t:o}=I(),s=x(!0),f=y(()=>i.rows),{savingCard:p,errorMsg:h,deleteVisible:_,deleting:b,deleteError:T,deleteTarget:E,levelDescription:O,hydrateRows:k,loadData:A,askDeleteRow:j,handleDelete:M,handleSave:N}=Xn({orgId:y(()=>i.orgId),rows:f,onSaved:()=>r(`saved`)}),P=!1;return c(()=>i.optionsReady,async e=>{if(!(!e||P)){P=!0;try{await A()}finally{k(),s.value=!1}}},{immediate:!0}),(n,r)=>(a(),w(`div`,wi,[C(`div`,null,[C(`h3`,Ti,v(u(o)(t.titleKey)),1),C(`p`,Ei,v(u(o)(t.descriptionKey)),1)]),s.value?(a(),d(X,{key:0,type:`detail`,count:1,rows:t.skeletonRows,cols:`grid-cols-1`,padding:`p-5`},null,8,[`rows`])):(a(),w(S,{key:1},[t.rows.length===0?(a(),w(`div`,Di,v(u(o)(t.emptyKey)),1)):(a(),w(`div`,Oi,[C(`table`,ki,[C(`thead`,null,[C(`tr`,Ai,[C(`th`,ji,v(u(o)(`job_management.potency_table_name`)),1),C(`th`,Mi,v(u(o)(`job_management.potency_table_level`)),1),C(`th`,Ni,v(u(o)(`job_management.potency_table_description`)),1),C(`th`,Pi,v(u(o)(`common.actions`)),1)])]),C(`tbody`,null,[(a(!0),w(S,null,e(t.rows,e=>(a(),w(`tr`,{key:e.type,class:`border-t border-gray-100 dark:border-gray-700 align-top`},[C(`td`,Fi,[C(`div`,Ii,v(e.competency_name),1),e.competency_definition?(a(),w(`div`,Li,v(e.competency_definition),1)):g(``,!0)]),C(`td`,Ri,[m(Y,{modelValue:e.job_management_value_id,"onUpdate:modelValue":t=>e.job_management_value_id=t,options:e.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:u(o)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),C(`td`,zi,v(u(O)(e)),1),C(`td`,Bi,[m(u(D),{icon:`pi pi-trash`,severity:`danger`,text:``,rounded:``,size:`small`,disabled:u(p),"aria-label":u(o)(`common.delete`),onClick:t=>u(j)(e)},null,8,[`disabled`,`aria-label`,`onClick`])])]))),128))])])])),u(h)?(a(),w(`div`,Vi,v(u(h)),1)):g(``,!0),t.rows.length>0?(a(),w(`div`,Hi,[m(u(D),{label:u(o)(t.saveLabelKey),icon:`pi pi-check`,size:`small`,loading:u(p),disabled:u(p)||!t.orgId,onClick:u(N)},null,8,[`label`,`loading`,`disabled`,`onClick`])])):g(``,!0)],64)),m(J,{visible:u(_),"onUpdate:visible":r[0]||=e=>l(_)?_.value=e:null,title:u(o)(t.deleteTitleKey),message:u(o)(t.deleteMessageKey,{name:u(E)?.competency_name||``}),loading:u(b),"error-msg":u(T),onConfirm:u(M),onCancel:r[1]||=e=>_.value=!1},null,8,[`visible`,`title`,`message`,`loading`,`error-msg`,`onConfirm`])]))}},Wi={__name:`TypesPotencyCard`,props:{orgId:String,types:{type:Array,required:!0},skeletonRows:{type:Number,default:2},titleKey:{type:String,required:!0},descriptionKey:{type:String,required:!0},emptyKey:{type:String,required:!0},saveLabelKey:{type:String,required:!0},deleteTitleKey:{type:String,required:!0},deleteMessageKey:{type:String,required:!0}},emits:[`saved`],setup(e,{emit:t}){let r=t,i=e,{t:o}=I(),s=x([]),c=x(!1);function l(e){s.value=i.types.filter(t=>(e[t.type]||[]).length>0).map(t=>({competency_id:``,competency_name:o(t.nameKey),competency_definition:``,type:t.type,levelOptions:e[t.type]||[],recordId:``,job_management_value_id:``}))}async function u(){try{let e=await M.get(`/api/v1/tenant/job-management/values/tree`),t={};(e.data?.data||[]).forEach(e=>{(e.types||[]).forEach(e=>{i.types.some(t=>t.type===e.type)&&(t[e.type]=(e.options||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``})))})}),l(t)}catch{s.value=[]}}return n(async()=>{await u(),c.value=!0}),(t,n)=>(a(),d(Ui,{"org-id":e.orgId,rows:s.value,"options-ready":c.value,"skeleton-rows":e.skeletonRows,"title-key":e.titleKey,"description-key":e.descriptionKey,"empty-key":e.emptyKey,"save-label-key":e.saveLabelKey,"delete-title-key":e.deleteTitleKey,"delete-message-key":e.deleteMessageKey,onSaved:n[0]||=e=>r(`saved`)},null,8,[`org-id`,`rows`,`options-ready`,`skeleton-rows`,`title-key`,`description-key`,`empty-key`,`save-label-key`,`delete-title-key`,`delete-message-key`]))}},Gi={__name:`ProblemSolvingPotencyCard`,props:{orgId:String},emits:[`saved`],setup(e,{emit:t}){let n=t,r=[{type:`thinking_environment`,nameKey:`job_management.problem_solving_environment`},{type:`thinking_chalenge`,nameKey:`job_management.problem_solving_challenge`}];return(t,i)=>(a(),d(Wi,{"org-id":e.orgId,types:r,"skeleton-rows":2,"title-key":`job_management.problem_solving_title`,"description-key":`job_management.problem_solving_description`,"empty-key":`job_management.problem_solving_empty`,"save-label-key":`job_management.save_problem_solving`,"delete-title-key":`job_management.problem_solving_confirm_delete_title`,"delete-message-key":`job_management.problem_solving_confirm_delete`,onSaved:i[0]||=e=>n(`saved`)},null,8,[`org-id`]))}},Ki={__name:`SkillPotencyCard`,props:{orgId:String},emits:[`saved`],setup(e,{emit:t}){let n=t,r=[{type:`communicating_influencing_skill`,nameKey:`job_management.skill_communicating_influencing`}];return(t,i)=>(a(),d(Wi,{"org-id":e.orgId,types:r,"skeleton-rows":2,"title-key":`job_management.skill_communicating_influencing_title`,"description-key":`job_management.skill_communicating_influencing_description`,"empty-key":`job_management.skill_communicating_influencing_empty`,"save-label-key":`job_management.save_skill`,"delete-title-key":`job_management.skill_confirm_delete_title`,"delete-message-key":`job_management.skill_confirm_delete`,onSaved:i[0]||=e=>n(`saved`)},null,8,[`org-id`]))}},qi={class:`space-y-4`},Ji={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Yi={class:`text-sm text-gray-500 dark:text-gray-400`},Xi={__name:`JobPotencySection`,props:{orgId:String,jobValueMap:{type:Object,default:()=>({})},competencyOptions:{type:Array,default:()=>[]}},emits:[`saved`],setup(e,{emit:t}){let n=t,{t:r}=I(),i=x(null);return(t,o)=>(a(),w(`div`,qi,[C(`div`,null,[C(`h2`,Ji,v(u(r)(`job_management.potency_competencies`)),1),C(`p`,Yi,v(u(r)(`job_management.potency_description`)),1)]),m(gr,{"org-id":e.orgId,onSaved:o[0]||=e=>n(`saved`)},null,8,[`org-id`]),m(Gr,{"org-id":e.orgId,onSaved:o[1]||=e=>n(`saved`),onWeightSaved:o[2]||=e=>i.value=e},null,8,[`org-id`]),m(Ci,{"org-id":e.orgId,"technical-weight":i.value,onSaved:o[3]||=e=>n(`saved`)},null,8,[`org-id`,`technical-weight`]),m(Gi,{"org-id":e.orgId,onSaved:o[4]||=e=>n(`saved`)},null,8,[`org-id`]),m(Ki,{"org-id":e.orgId,onSaved:o[5]||=e=>n(`saved`)},null,8,[`org-id`])]))}},Zi={class:`space-y-6`},Qi={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},$i={class:`text-sm text-gray-500 dark:text-gray-400`},ea={key:0,class:`flex items-center justify-center py-12`},ta={key:0,class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden`},na={class:`px-5 py-3 border-b border-gray-200 dark:border-gray-700 font-semibold text-sm text-gray-700 dark:text-gray-300`},ra={class:`divide-y divide-gray-100 dark:divide-gray-700`},ia={class:`hidden md:grid grid-cols-[minmax(0,2fr)_minmax(0,3fr)_auto] gap-4 px-5 py-2.5 bg-gray-50 dark:bg-gray-900/40 text-[11px] uppercase tracking-wider text-gray-400 dark:text-gray-500 font-medium`},aa={class:`text-right`},oa={class:`grid grid-cols-1 md:grid-cols-[minmax(0,2fr)_minmax(0,3fr)_auto] md:items-center gap-2`},sa={class:`text-sm font-semibold text-gray-700 dark:text-gray-200`},ca={class:`flex flex-wrap gap-1.5`},la={class:`font-medium`},ua={key:0,class:`font-mono`},da={class:`font-bold text-emerald-600 dark:text-emerald-400`},fa={class:`text-right`},pa={class:`text-sm font-bold text-gray-900 dark:text-gray-100`},ma={class:`px-5 py-4 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/40`},ha={class:`grid grid-cols-1 sm:grid-cols-2 gap-4`},ga={class:`flex items-center justify-between`},_a={class:`text-xs text-gray-500 dark:text-gray-400`},va={class:`text-sm font-bold text-emerald-600 dark:text-emerald-400`},ya={class:`flex items-center justify-between`},ba={class:`text-xs text-gray-500 dark:text-gray-400`},xa={class:`text-sm font-bold text-blue-600 dark:text-blue-400`},Sa={key:2},Ca={class:`flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500`},wa={class:`text-sm font-medium`},Ta={class:`text-xs mt-1`},Ea={class:`flex justify-end gap-3`},Da=`/api/v1/tenant/job-management/scores/org`,Oa={__name:`JobScoreSection`,props:{orgId:String},emits:[`saved`],setup(t,{emit:r}){let i=t,o=r,{t:s}=I(),c=A(),l=x(!1),p=x(!1),h=x(null),b=[{key:`education_experience`,labelKey:`job_management.education_experience`,points:[{labelKey:`job_management.group_education`,level:`education_level`,pts:`education_points`},{labelKey:`job_management.group_experience`,level:`experience_level`,pts:`experience_points`}]},{key:`potentials`,labelKey:`job_management.score_potentials`,points:[{labelKey:`job_management.average_level`,level:`average_level`,pts:null}]},{key:`competencies`,labelKey:`job_management.potency_competencies`,score:`base_score`,points:[{labelKey:`job_management.potency_technical_title`,level:`technical_average_level`,pts:`technical_points`},{labelKey:`job_management.potency_managerial_title`,level:`managerial_average_level`,pts:`managerial_points`},{labelKey:`job_management.skill_communicating_influencing`,level:`communication_level`,pts:`communication_points`}]},{key:`problem_solving`,labelKey:`job_management.problem_solving_title`,points:[{labelKey:`job_management.problem_solving_environment`,level:`environment_level`,pts:`environment_points`},{labelKey:`job_management.problem_solving_challenge`,level:`challenge_level`,pts:`challenge_points`}]},{key:`financial_authority`,labelKey:`job_management.financials`,points:[{labelKey:`job_management.cash_level`,level:`money_level`,pts:`money_points`},{labelKey:`job_management.authority_level`,level:`authority_level`,pts:`authority_points`},{labelKey:`job_management.impact_level`,level:`impact_level`,pts:`impact_points`}]},{key:`asset_authority`,labelKey:`job_management.assets`,points:[{labelKey:`job_management.asset_type`,level:`asset_value_level`,pts:`asset_value_points`},{labelKey:`job_management.authority_level`,level:`asset_authority_level`,pts:`asset_authority_points`}]},{key:`subordinate_control`,labelKey:`job_management.subordinate_controls`,points:[{labelKey:`job_management.score_level`,level:`level`,pts:`points`}]},{key:`work_scope`,labelKey:`job_management.relationships`,points:[{labelKey:`job_management.relationship_group_scope`,level:`scope_level`,pts:`scope_points`},{labelKey:`job_management.frequency`,level:`frequency_level`,pts:`frequency_points`}]},{key:`work_activity`,labelKey:`job_management.activities`,points:[{labelKey:`job_management.score_level`,level:`level`,pts:`points`}]},{key:`work_risk`,labelKey:`job_management.risks`,points:[{labelKey:`job_management.environment_risk`,level:`environment_level`,pts:`environment_points`},{labelKey:`job_management.hazard`,level:`hazard_level`,pts:`hazard_points`}]}],T=y(()=>{if(!h.value?.components)return null;try{return JSON.parse(h.value.components)}catch{return null}}),E=y(()=>T.value?b.map(e=>{let t=T.value[e.key]||{};return{key:e.key,labelKey:e.labelKey,score:t[e.score||`score`]??0,points:e.points.map(e=>({labelKey:e.labelKey,level:t[e.level]??null,points:e.pts==null?null:t[e.pts]??0}))}}):[]);function O(e){return e?.toLocaleString?.(`id-ID`)??`-`}function k(e){return e==null?`-`:String(e)}async function j(){if(i.orgId){l.value=!0;try{let e=await M.get(`${Da}/${i.orgId}`);h.value=e.data?.data||null,o(`saved`)}catch{h.value=null}finally{l.value=!1}}}async function N(){if(i.orgId){p.value=!0;try{let e=await M.put(`${Da}/${i.orgId}`,{components:null});h.value=e.data?.data||null,c.add({severity:`success`,summary:s(`message.success`),detail:s(`job_management.score_calculated`),life:2e3})}catch(e){c.add({severity:`error`,summary:s(`message.error`),detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{p.value=!1}}}return n(j),(t,n)=>(a(),w(`div`,Zi,[C(`div`,null,[C(`h2`,Qi,v(u(s)(`job_management.scores`)),1),C(`p`,$i,v(u(s)(`job_management.score_description`)),1)]),l.value?(a(),w(`div`,ea,[...n[0]||=[C(`i`,{class:`pi pi-spin pi-spinner text-emerald-500 text-2xl`},null,-1)]])):h.value?(a(),w(S,{key:1},[E.value.length?(a(),w(`div`,ta,[C(`div`,na,v(u(s)(`job_management.component_breakdown`)),1),C(`div`,ra,[C(`div`,ia,[C(`span`,null,v(u(s)(`job_management.score_component`)),1),C(`span`,null,v(u(s)(`job_management.score_points`)),1),C(`span`,aa,v(u(s)(`job_management.score_score`)),1)]),(a(!0),w(S,null,e(E.value,t=>(a(),w(`div`,{key:t.key,class:`px-5 py-1`},[C(`div`,oa,[C(`div`,sa,v(u(s)(t.labelKey)),1),C(`div`,ca,[(a(!0),w(S,null,e(t.points,e=>(a(),w(`span`,{key:e.labelKey,class:f([`inline-flex items-center gap-1 text-[11px] px-2 py-0.5 rounded-md border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-900/40 text-gray-600 dark:text-gray-300`,{"opacity-50":e.level==null}])},[C(`span`,la,v(u(s)(e.labelKey)),1),e.level==null?g(``,!0):(a(),w(`span`,ua,`Lv.`+v(k(e.level)),1)),e.points==null?e.level==null?(a(),w(S,{key:2},[_(`—`)],64)):g(``,!0):(a(),w(S,{key:1},[n[1]||=C(`i`,{class:`pi pi-arrow-right text-[8px] opacity-60`},null,-1),C(`span`,da,v(O(e.points)),1)],64))],2))),128))]),C(`div`,fa,[C(`span`,pa,v(O(t.score)),1)])])]))),128))]),C(`div`,ma,[C(`div`,ha,[C(`div`,ga,[C(`span`,_a,v(u(s)(`job_management.value_with_financial`)),1),C(`span`,va,v(O(h.value.job_value_with_financial)),1)]),C(`div`,ya,[C(`span`,ba,v(u(s)(`job_management.value_without_financial`)),1),C(`span`,xa,v(O(h.value.job_value_without_financial)),1)])])])])):g(``,!0)],64)):(a(),w(`div`,Sa,[C(`div`,Ca,[n[2]||=C(`i`,{class:`pi pi-calculator text-4xl mb-3 opacity-50`},null,-1),C(`p`,wa,v(u(s)(`job_management.no_score`)),1),C(`p`,Ta,v(u(s)(`job_management.score_hint`)),1)])])),C(`div`,Ea,[m(u(D),{label:u(s)(`common.refresh`),icon:`pi pi-refresh`,size:`small`,text:``,onClick:j},null,8,[`label`]),h.value?(a(),d(u(D),{key:0,label:u(s)(`job_management.recalculate`),icon:`pi pi-calculator`,size:`small`,severity:`info`,loading:p.value,onClick:N},null,8,[`label`,`loading`])):g(``,!0)])]))}},ka={class:`grid grid-cols-1 md:grid-cols-3 gap-4`},Aa={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4`},ja={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},Ma={class:`text-2xl font-bold text-emerald-600 dark:text-emerald-400`},Na={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4`},Pa={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},Fa={class:`text-2xl font-bold text-blue-600 dark:text-blue-400`},Ia={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4`},La={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},Ra={class:`mt-3 flex flex-wrap items-center gap-2`},za={key:1,class:`text-[10px] text-gray-400`},Ba={key:0,class:`text-[10px] text-gray-400 mt-2`},Va=`/api/v1/tenant/job-management/scores/org`,Ha={__name:`JobScoreSummary`,props:{orgId:String},setup(t,{expose:r}){let i=t,{t:o}=I(),s=x(!0),l=x(null);function f(e){return e?.toLocaleString?.(`id-ID`)??`-`}async function p(){if(i.orgId){s.value=!0;try{let e=await M.get(`${Va}/${i.orgId}`);l.value=e.data?.data||null}catch{l.value=null}finally{s.value=!1}}}return r({refresh:p}),c(()=>i.orgId,p),n(p),(t,n)=>(a(),w(`div`,ka,[s.value&&!l.value?(a(),w(S,{key:0},e(3,e=>C(`div`,{key:e,class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4 animate-pulse`},[...n[0]||=[C(`div`,{class:`h-3 w-24 bg-gray-200 dark:bg-gray-700 rounded mb-2`},null,-1),C(`div`,{class:`h-7 w-16 bg-gray-200 dark:bg-gray-700 rounded`},null,-1)]])),64)):(a(),w(S,{key:1},[C(`div`,Aa,[C(`div`,ja,v(u(o)(`job_management.value_with_financial`)),1),C(`div`,Ma,v(f(l.value?.job_value_with_financial)),1)]),C(`div`,Na,[C(`div`,Pa,v(u(o)(`job_management.value_without_financial`)),1),C(`div`,Fa,v(f(l.value?.job_value_without_financial)),1)]),C(`div`,Ia,[C(`div`,La,v(u(o)(`job_management.has_financial_authority`)),1),m(u(R),{value:l.value?l.value.has_financial_authority?u(o)(`common.yes`):u(o)(`common.no`):`-`,severity:l.value?.has_financial_authority?`success`:`danger`,class:`!text-xs`},null,8,[`value`,`severity`]),C(`div`,Ra,[l.value?(a(),d(u(R),{key:0,value:l.value.is_complete?u(o)(`job_management.score_complete`):u(o)(`job_management.score_incomplete`),severity:l.value.is_complete?`success`:`warning`,icon:l.value.is_complete?`pi pi-check-circle`:`pi pi-exclamation-triangle`,class:`!text-xs`},null,8,[`value`,`severity`,`icon`])):g(``,!0),l.value?.is_complete&&l.value.completed_at?(a(),w(`span`,za,v(u(o)(`job_management.completed_at`))+`: `+v(l.value.completed_at),1)):g(``,!0)]),l.value?.calculated_at?(a(),w(`div`,Ba,v(u(o)(`job_management.calculated_at`))+`: `+v(l.value.calculated_at),1)):g(``,!0)])],64))]))}},Ua={class:`max-w-full mx-auto`},Wa={key:0,class:`flex gap-6`},Ga={class:`w-56 space-y-2`},Ka={class:`flex-1 space-y-3`},qa={key:1,class:`flex gap-6`},Ja={class:`w-56 shrink-0 space-y-1`},Ya=[`onClick`,`onKeydown`],Xa={key:0,class:`pi pi-check text-xs`},Za={class:`flex-1 min-w-0`},Qa={key:0,class:`pi pi-check-circle text-emerald-400 text-xs shrink-0`},$a={class:`flex-1 min-w-0 space-y-4`},eo={class:`flex flex-col md:flex-row gap-4`},to={class:`md:w-72 shrink-0 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4`},no={class:`flex items-center gap-2 mb-3`},ro={class:`text-sm font-semibold text-gray-800 dark:text-gray-100 truncate`},io={class:`flex items-center justify-between gap-2`},ao={class:`text-[10px] uppercase tracking-wider text-gray-400 dark:text-gray-500`},oo={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 font-mono truncate`},so={class:`flex-1 min-w-0`},co={__name:`JobManagementForm`,setup(t){let r=N(),i=j(),{t:s}=I(),l=A(),p=y(()=>i.query.org_id||``),h=x(0),_=x(!0),b=x(Array(15).fill(!1)),T=x(``),D=x(``),O=x(``),k=x(``),P=x(``),F=x([]),L=x([]),R=x([]),z=x({}),B=x([]),V=x(null),H=[{labelKey:`job_management.identifications`,icon:`pi pi-id-card`,comp:ce},{labelKey:`job_management.objectives`,icon:`pi pi-bullseye`,comp:_e},{labelKey:`job_management.responsibilities_title`,icon:`pi pi-list-check`,comp:st},{labelKey:`job_management.education_experience`,icon:`pi pi-graduation-cap`,comp:Pe},{labelKey:`job_management.potency_competencies`,icon:`pi pi-star`,comp:Xi},{labelKey:`job_management.financials`,icon:`pi pi-money-bill`,comp:Jn},{labelKey:`job_management.assets`,icon:`pi pi-box`,comp:Fn},{labelKey:`job_management.subordinate_controls`,icon:`pi pi-sitemap`,comp:En},{labelKey:`job_management.relationships`,icon:`pi pi-share-alt`,comp:_n},{labelKey:`job_management.activities`,icon:`pi pi-bolt`,comp:Mt},{labelKey:`job_management.risks`,icon:`pi pi-exclamation-triangle`,comp:Vt},{labelKey:`job_management.hr_authorities`,icon:`pi pi-users`,comp:ht},{labelKey:`job_management.op_authorities`,icon:`pi pi-cog`,comp:Ct},{labelKey:`job_management.scores`,icon:`pi pi-calculator`,comp:Oa}],U=y(()=>H[h.value]?.comp||null);function W(e){return h.value===e?`bg-emerald-50 dark:bg-emerald-900/20 ring-1 ring-emerald-300 dark:ring-emerald-700`:(b.value[e],`hover:bg-gray-50 dark:hover:bg-gray-800`)}function G(e){return h.value===e?`bg-emerald-600 text-white`:b.value[e]?`bg-emerald-100 dark:bg-emerald-800 text-emerald-600 dark:text-emerald-300`:`bg-gray-200 dark:bg-gray-600 text-gray-600 dark:text-gray-300`}function K(e){return h.value===e?`text-emerald-700 dark:text-emerald-300`:b.value[e]?`text-emerald-600 dark:text-emerald-400`:`text-gray-700 dark:text-gray-300`}function q(e){h.value=e,r.replace({query:{...i.query,section:String(e)}})}function ee(e){typeof e==`number`&&(b.value[e]=!0),V.value?.refresh()}async function te(){if(p.value)try{let e=(await M.get(`/api/v1/tenant/organizations/${p.value}`)).data?.data;e&&(T.value=e.nomenclature||``,D.value=e.full_code||e.code||``,O.value=e.organization_summary_id||``,k.value=e.grading_id||``,P.value=e.job_family_id||``)}catch{}}async function J(){try{let[e,t,n,r]=await Promise.all([M.get(`/api/v1/tenant/settings/gradings?per_page=100`),M.get(`/api/v1/tenant/job-management/values?per_page=500`),M.get(`/api/v1/tenant/competency/competencies?per_page=200`).catch(()=>({data:{data:[]}})),M.get(`/api/v1/tenant/settings/job-families?per_page=100`)]);F.value=(e.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id})),L.value=(r.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id}));let i=t.data?.data||[];R.value=i.map(e=>({label:`${e.type}${e.level?` Lv.`+e.level:``}${e.descriptions?` — `+e.descriptions:``}`,value:e.id,type:e.type,level:e.level,descriptions:e.descriptions}));let a={};i.forEach(e=>{a[e.type]||(a[e.type]=[]),a[e.type].push({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``,type_group:e.type_group||``,description_group:e.description_group||``})}),z.value=a,B.value=(n.data?.data||[]).map(e=>({label:e.name||e.code,value:e.id,field:e.field||``,definition:e.definition||``}))}catch{}}return c(p,(e,t)=>{e!==t&&(b.value=Array(H.length).fill(!1),T.value=``,D.value=``,O.value=``,k.value=``,P.value=``,te())}),n(async()=>{try{await Promise.all([te(),J()]);let e=parseInt(i.query.section);!isNaN(e)&&e>=0&&e<H.length&&(h.value=e)}catch(e){l.add({severity:`error`,summary:s(`message.error`),detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{_.value=!1}}),(t,n)=>(a(),w(`div`,Ua,[_.value?(a(),w(`div`,Wa,[C(`div`,Ga,[(a(),w(S,null,e(8,e=>C(`div`,{key:e,class:`h-12 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))]),C(`div`,Ka,[(a(),w(S,null,e(6,e=>C(`div`,{key:e,class:`h-9 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))])])):(a(),w(`div`,qa,[C(`div`,Ja,[(a(),w(S,null,e(H,(e,t)=>C(`div`,{key:t,role:`button`,tabindex:0,class:f([`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150 cursor-pointer select-none`,W(t)]),onClick:e=>q(t),onKeydown:E(e=>q(t),[`enter`])},[C(`div`,{class:f([`w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0 transition-colors duration-150`,G(t)])},[b.value[t]?(a(),w(`i`,Xa)):(a(),w(`i`,{key:1,class:f(e.icon)},null,2))],2),C(`div`,Za,[C(`div`,{class:f([`text-sm font-medium truncate`,K(t)])},v(u(s)(e.labelKey)),3)]),b.value[t]?(a(),w(`i`,Qa)):g(``,!0)],42,Ya)),64))]),C(`div`,$a,[(a(),w(`div`,{key:`summary-${p.value}`,class:`sticky top-0 z-10 bg-white dark:bg-gray-900 pt-1 pb-3`},[C(`div`,eo,[C(`div`,to,[C(`div`,no,[n[0]||=C(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400`},[C(`i`,{class:`pi pi-briefcase text-sm`})],-1),C(`h3`,ro,v(T.value||u(s)(`job_management.job_info_untitled`)),1)]),C(`div`,io,[C(`span`,ao,v(u(s)(`organization.full_code`)),1),C(`span`,oo,v(D.value||`-`),1)])]),C(`div`,so,[m(Ha,{ref_key:`scoreSummaryRef`,ref:V,"org-id":p.value},null,8,[`org-id`])])])])),(a(),d(o(U.value),{key:`${h.value}-${p.value}`,"org-id":p.value,"org-name":T.value,"org-code":D.value,"org-summary-id":O.value,"org-grading-id":k.value,"org-job-family-id":P.value,"grading-options":F.value,"job-family-options":L.value,"job-value-options":R.value,"competency-options":B.value,"job-value-map":z.value,onSaved:ee},null,40,[`org-id`,`org-name`,`org-code`,`org-summary-id`,`org-grading-id`,`org-job-family-id`,`grading-options`,`job-family-options`,`job-value-options`,`competency-options`,`job-value-map`]))])]))]))}};export{co as default};