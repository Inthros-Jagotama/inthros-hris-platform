const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["assets/quill-BLmY9xB4.js","assets/rolldown-runtime-QTnfLwEv.js"])))=>i.map(i=>d[i]);
import{A as e,C as t,E as n,H as r,K as i,N as a,O as o,P as s,U as c,X as l,c as u,dt as d,ft as f,h as p,j as m,l as h,m as g,o as _,pt as v,r as y,s as b,u as x,ut as S,y as C}from"./runtime-core.esm-bundler-CFm0BMYx.js";import{a as w,d as T,t as E}from"./button-DVe0lZgG.js";import{A as D,a as O}from"./ripple-FLjJmYYY.js";import{l as k,m as A,s as j,t as M,u as N}from"./index-DJqKIIqm.js";import{t as P}from"./useI18n-NFTbW2yU.js";import{n as F}from"./responseHandler-B5MnXl3B.js";import{t as I}from"./tag--2hRZhYy.js";import{t as L}from"./FormRow-D-Jmd1__.js";import{t as R}from"./baseeditableholder-CnB2pcI8.js";import{t as z}from"./textarea-CZKEJ-vM.js";import{t as B}from"./TextInput-DT_yhCva.js";import{a as V,n as H,s as U,t as W}from"./column-CRD4EM5m.js";import{t as G}from"./SelectLabel-BUp5q6Ob.js";import{t as K}from"./ConfirmDeleteDialog-3BZ6Ov3U.js";import{t as ee}from"./SkeletonTable-BS8Q7JD0.js";import{t as te}from"./toggleswitch-BSOZ6l02.js";var ne={class:`space-y-4`},re={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},ie={class:`text-sm text-gray-500 dark:text-gray-400`},ae={class:`max-w-2xl space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},oe={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},se={class:`flex justify-end pt-2`},ce=`/api/v1/tenant/job-management/identifications`,le={__name:`JobIdentificationSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},orgGradingId:{type:String,default:``},orgJobFamilyId:{type:String,default:``},gradingOptions:{type:Array,default:()=>[]},jobFamilyOptions:{type:Array,default:()=>[]}},emits:[`saved`],setup(e,{emit:t}){let a=t,s=e,{t:c}=P(),u=A(),d=i(!1),f=i(``),m=i({}),g=i(``),y=i({grading_id:``}),S=_(()=>{let e=s.jobFamilyOptions.find(e=>e.value===s.orgJobFamilyId);return e?e.label:s.orgJobFamilyId||`-`});function C(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function T(){if(s.orgId)try{let e=(await w.get(ce,{params:{organization_id:s.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];g.value=t.id,y.value.grading_id=t.grading_id||s.orgGradingId||``}else y.value.grading_id=s.orgGradingId||``}catch{y.value.grading_id=s.orgGradingId||``}}async function D(){if(f.value=``,m.value={},!y.value.grading_id){f.value=c(`job_management.grading_required`);return}d.value=!0;try{let e={nomenclature:s.orgName||``,full_code:s.orgCode||``,grading_id:y.value.grading_id,organization_id:s.orgId};if(g.value)await w.put(`${ce}/${g.value}`,{grading_id:y.value.grading_id});else{let t=await w.post(ce,e);g.value=t.data?.data?.id||``}u.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=C(e);Object.keys(t).length>0?(m.value=t,f.value=Object.values(t).join(`, `)):f.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{d.value=!1}}return n(T),(t,n)=>(o(),x(`div`,ne,[b(`div`,null,[b(`h2`,re,v(l(c)(`job_management.identifications`)),1),b(`p`,ie,v(l(c)(`job_management.identification_description`)),1)]),b(`div`,ae,[p(L,{label:l(c)(`organization.nomenclature`)},{default:r(()=>[p(B,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:l(c)(`organization.full_code`)},{default:r(()=>[p(B,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:l(c)(`organization.job_family`)},{default:r(()=>[p(B,{"model-value":S.value,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:l(c)(`organization.grading`)},{default:r(()=>[p(l(U),{modelValue:y.value.grading_id,"onUpdate:modelValue":n[0]||=e=>y.value.grading_id=e,options:e.gradingOptions,"option-label":`label`,"option-value":`value`,placeholder:l(c)(`organization.select_grading`),class:`w-full`,size:`small`,invalid:!!m.value.grading_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`]),f.value?(o(),x(`div`,oe,v(f.value),1)):h(``,!0),b(`div`,se,[p(l(E),{label:l(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:d.value,disabled:!y.value.grading_id,onClick:D},null,8,[`label`,`loading`,`disabled`])])])]))}},ue={class:`space-y-4`},de={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},fe={class:`text-sm text-gray-500 dark:text-gray-400`},pe={class:`max-w-2xl space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},me={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},he={class:`flex justify-end gap-2 pt-2`},q=`/api/v1/tenant/job-management/objectives`,ge={__name:`JobObjectiveSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let a=t,s=e,{t:c}=P(),d=A(),f=i(!1),m=i(!1),g=i(``),_=i({}),y=i(``),C=i(!1),T=i(``),D=i({objective:``});function O(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function k(){if(s.orgId)try{let e=(await w.get(q,{params:{organization_id:s.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];y.value=t.id,D.value.objective=t.objective||``}}catch{}}async function j(){g.value=``,_.value={},f.value=!0;try{let e={nomenclature:s.orgName||``,full_code:s.orgCode||``,objective:D.value.objective||``,organization_id:s.orgId};if(y.value)await w.put(`${q}/${y.value}`,{objective:D.value.objective||``});else{let t=await w.post(q,e);y.value=t.data?.data?.id||``}d.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=O(e);Object.keys(t).length>0?(_.value=t,g.value=Object.values(t).join(`, `)):g.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{f.value=!1}}async function M(){if(y.value){m.value=!0,T.value=``;try{await w.delete(`${q}/${y.value}`),C.value=!1,y.value=``,D.value.objective=``,a(`saved`),d.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||c(`message.operation_failed`)}finally{m.value=!1}}}return n(k),(t,n)=>(o(),x(`div`,ue,[b(`div`,null,[b(`h2`,de,v(l(c)(`job_management.objectives`)),1),b(`p`,fe,v(l(c)(`job_management.objective_description`)),1)]),b(`div`,pe,[p(L,{label:l(c)(`organization.nomenclature`)},{default:r(()=>[p(B,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:l(c)(`organization.full_code`)},{default:r(()=>[p(B,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:l(c)(`job_management.objective`)},{default:r(()=>[p(l(z),{modelValue:D.value.objective,"onUpdate:modelValue":n[0]||=e=>D.value.objective=e,rows:`3`,class:S([`w-full`,{"p-invalid":_.value.objective}]),placeholder:l(c)(`job_management.objective`)+`...`},null,8,[`modelValue`,`class`,`placeholder`])]),_:1},8,[`label`]),g.value?(o(),x(`div`,me,v(g.value),1)):h(``,!0),b(`div`,he,[y.value?(o(),u(l(E),{key:0,label:l(c)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:n[1]||=e=>C.value=!0},null,8,[`label`])):h(``,!0),p(l(E),{label:y.value?l(c)(`common.update`):l(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:f.value,disabled:f.value,onClick:j},null,8,[`label`,`loading`,`disabled`])])]),p(K,{visible:C.value,"onUpdate:visible":n[2]||=e=>C.value=e,loading:m.value,"error-msg":T.value,onConfirm:M,onCancel:n[3]||=e=>C.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},_e={class:`space-y-4`},ve={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},ye={class:`text-sm text-gray-500 dark:text-gray-400`},be={class:`max-w-2xl space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},xe={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Se={class:`flex justify-end gap-2 pt-2`},J=`/api/v1/tenant/job-management/education-experiences`,Ce={__name:`JobEduExpSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:t}){let a=t,s=e,{t:c}=P(),d=A(),f=i(!1),m=i(!1),g=i(``),_=i({}),y=i(``),S=i(!1),C=i(``),T=i({education_id:``,education_major_id:``,job_family_id:``,experience_range:``}),D=[{label:`0-2 Tahun`,value:`0-2 Tahun`},{label:`3-5 Tahun`,value:`3-5 Tahun`},{label:`6-8 Tahun`,value:`6-8 Tahun`},{label:`9-11 Tahun`,value:`9-11 Tahun`},{label:`> 12 Tahun`,value:`> 12 Tahun`}],O=i([]),k=i([]),j=i([]);async function M(){try{let[e,t,n]=await Promise.all([w.get(`/api/v1/tenant/settings/educations?per_page=100`),w.get(`/api/v1/tenant/settings/education-majors?per_page=200`),w.get(`/api/v1/tenant/settings/job-families?per_page=100`)]);O.value=(e.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id})),k.value=(t.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id})),j.value=(n.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id}))}catch{}}async function N(){if(s.orgId)try{let e=(await w.get(J,{params:{organization_id:s.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];y.value=t.id,T.value.education_id=t.education_id||``,T.value.education_major_id=t.education_major_id||``,T.value.job_family_id=t.job_family_id||``,T.value.experience_range=t.experience_range||``}}catch{}}async function I(){g.value=``,_.value={},f.value=!0;try{let e={nomenclature:s.orgName||``,full_code:s.orgCode||``,education_id:T.value.education_id||null,education_major_id:T.value.education_major_id||null,job_family_id:T.value.job_family_id||null,experience_range:T.value.experience_range||null,organization_id:s.orgId};if(y.value)await w.put(`${J}/${y.value}`,{education_id:T.value.education_id||``,education_major_id:T.value.education_major_id||``,job_family_id:T.value.job_family_id||``,experience_range:T.value.experience_range||``});else{let t=await w.post(J,e);y.value=t.data?.data?.id||``}d.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=F(e);Object.keys(t).length>0?(_.value=t,g.value=Object.values(t).join(`, `)):g.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{f.value=!1}}async function R(){if(y.value){m.value=!0,C.value=``;try{await w.delete(`${J}/${y.value}`),S.value=!1,y.value=``,T.value.education_id=``,T.value.education_major_id=``,T.value.job_family_id=``,T.value.experience_range=``,a(`saved`),d.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){C.value=e?.response?.data?.error?.message||c(`message.operation_failed`)}finally{m.value=!1}}}return n(async()=>{await Promise.all([M(),N()])}),(t,n)=>(o(),x(`div`,_e,[b(`div`,null,[b(`h2`,ve,v(l(c)(`job_management.education_experience`)),1),b(`p`,ye,v(l(c)(`job_management.education_experience_description`)),1)]),b(`div`,be,[p(L,{label:l(c)(`organization.nomenclature`)},{default:r(()=>[p(B,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:l(c)(`organization.full_code`)},{default:r(()=>[p(B,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:l(c)(`job_management.education_level`),errors:_.value?.education_id},{default:r(()=>[p(l(U),{modelValue:T.value.education_id,"onUpdate:modelValue":n[0]||=e=>T.value.education_id=e,options:O.value,"option-label":`label`,"option-value":`value`,placeholder:l(c)(`common.select`),class:`w-full`,size:`small`,showClear:``,invalid:!!_.value.education_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(c)(`job_management.education_major`),errors:_.value?.education_major_id},{default:r(()=>[p(l(U),{modelValue:T.value.education_major_id,"onUpdate:modelValue":n[1]||=e=>T.value.education_major_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:l(c)(`common.select`),class:`w-full`,size:`small`,showClear:``,invalid:!!_.value.education_major_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(c)(`job_management.job_family`),errors:_.value?.job_family_id},{default:r(()=>[p(l(U),{modelValue:T.value.job_family_id,"onUpdate:modelValue":n[2]||=e=>T.value.job_family_id=e,options:j.value,"option-label":`label`,"option-value":`value`,placeholder:l(c)(`common.select`),class:`w-full`,size:`small`,showClear:``,invalid:!!_.value.job_family_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(c)(`job_management.experience_range`),errors:_.value?.experience_range},{default:r(()=>[p(l(U),{modelValue:T.value.experience_range,"onUpdate:modelValue":n[3]||=e=>T.value.experience_range=e,options:D,"option-label":`label`,"option-value":`value`,placeholder:l(c)(`common.select`),class:`w-full`,size:`small`,showClear:``,invalid:!!_.value.experience_range},null,8,[`modelValue`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`]),g.value?(o(),x(`div`,xe,v(g.value),1)):h(``,!0),b(`div`,Se,[y.value?(o(),u(l(E),{key:0,label:l(c)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:n[4]||=e=>S.value=!0},null,8,[`label`])):h(``,!0),p(l(E),{label:y.value?l(c)(`common.update`):l(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:f.value,disabled:f.value,onClick:I},null,8,[`label`,`loading`,`disabled`])])]),p(K,{visible:S.value,"onUpdate:visible":n[5]||=e=>S.value=e,loading:m.value,"error-msg":C.value,onConfirm:R,onCancel:n[6]||=e=>S.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},we=O.extend({name:`editor`,style:`
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
`,classes:{root:function(e){return[`p-editor`,{"p-invalid":e.instance.$invalid}]},toolbar:`p-editor-toolbar`,content:`p-editor-content`}}),Te={name:`BaseEditor`,extends:R,props:{placeholder:String,readonly:Boolean,formats:Array,editorStyle:null,modules:null},style:we,provide:function(){return{$pcEditor:this,$parentInstance:this}}};function Y(e){"@babel/helpers - typeof";return Y=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},Y(e)}function Ee(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function De(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?Ee(Object(n),!0).forEach(function(t){Oe(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):Ee(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function Oe(e,t,n){return(t=ke(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function ke(e){var t=Ae(e,`string`);return Y(t)==`symbol`?t:t+``}function Ae(e,t){if(Y(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(Y(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var je=function(){try{return window.Quill}catch{return null}}(),X={name:`Editor`,extends:Te,inheritAttrs:!1,emits:[`text-change`,`selection-change`,`load`],quill:null,watch:{modelValue:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},d_value:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},readonly:function(){this.handleReadOnlyChange()}},mounted:function(){var e=this,t={modules:De({toolbar:this.$refs.toolbarElement},this.modules),readOnly:this.readonly,theme:`snow`,formats:this.formats,placeholder:this.placeholder};je?(this.quill=new je(this.$refs.editorElement,t),this.initQuill(),this.handleLoad()):j(()=>import(`./quill-BLmY9xB4.js`).then(function(n){n&&D(e.$refs.editorElement)&&(n.default?e.quill=new n.default(e.$refs.editorElement,t):e.quill=new n(e.$refs.editorElement,t),e.initQuill())}),__vite__mapDeps([0,1])).then(function(){e.handleLoad()})},beforeUnmount:function(){this.quill=null},methods:{renderValue:function(e){if(this.quill)if(e){var t=this.quill.clipboard.convert({html:e});this.quill.setContents(t)}else this.quill.setText(``)},initQuill:function(){var e=this;this.renderValue(this.d_value),this.quill.on(`text-change`,function(t,n,r){if(r===`user`){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();i===`<p><br></p>`&&(i=``),e.writeValue(i),e.$emit(`text-change`,{htmlValue:i,textValue:a,delta:t,source:r,instance:e.quill})}}),this.quill.on(`selection-change`,function(t,n,r){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();e.$emit(`selection-change`,{htmlValue:i,textValue:a,range:t,oldRange:n,source:r,instance:e.quill})})},handleLoad:function(){this.quill&&this.quill.getModule(`toolbar`)&&this.$emit(`load`,{instance:this.quill})},handleReadOnlyChange:function(){this.quill&&this.quill.enable(!this.readonly)}}};function Me(e,n,r,i,a,s){return o(),x(`div`,t({class:e.cx(`root`)},e.ptmi(`root`)),[b(`div`,t({ref:`toolbarElement`,class:e.cx(`toolbar`)},e.ptm(`toolbar`)),[m(e.$slots,`toolbar`,{},function(){return[b(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[b(`select`,t({class:`ql-header`,defaultValue:`0`},e.ptm(`header`)),[b(`option`,t({value:`1`},e.ptm(`option`)),`Heading`,16),b(`option`,t({value:`2`},e.ptm(`option`)),`Subheading`,16),b(`option`,t({value:`0`},e.ptm(`option`)),`Normal`,16)],16),b(`select`,t({class:`ql-font`},e.ptm(`font`)),[b(`option`,d(C(e.ptm(`option`))),null,16),b(`option`,t({value:`serif`},e.ptm(`option`)),null,16),b(`option`,t({value:`monospace`},e.ptm(`option`)),null,16)],16)],16),b(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[b(`button`,t({class:`ql-bold`,type:`button`},e.ptm(`bold`)),null,16),b(`button`,t({class:`ql-italic`,type:`button`},e.ptm(`italic`)),null,16),b(`button`,t({class:`ql-underline`,type:`button`},e.ptm(`underline`)),null,16)],16),b(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[b(`select`,t({class:`ql-color`},e.ptm(`color`)),null,16),b(`select`,t({class:`ql-background`},e.ptm(`background`)),null,16)],16),b(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[b(`button`,t({class:`ql-list`,value:`ordered`,type:`button`},e.ptm(`list`)),null,16),b(`button`,t({class:`ql-list`,value:`bullet`,type:`button`},e.ptm(`list`)),null,16),b(`select`,t({class:`ql-align`},e.ptm(`select`)),[b(`option`,t({defaultValue:``},e.ptm(`option`)),null,16),b(`option`,t({value:`center`},e.ptm(`option`)),null,16),b(`option`,t({value:`right`},e.ptm(`option`)),null,16),b(`option`,t({value:`justify`},e.ptm(`option`)),null,16)],16)],16),b(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[b(`button`,t({class:`ql-link`,type:`button`},e.ptm(`link`)),null,16),b(`button`,t({class:`ql-image`,type:`button`},e.ptm(`image`)),null,16),b(`button`,t({class:`ql-code-block`,type:`button`},e.ptm(`codeBlock`)),null,16)],16),b(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[b(`button`,t({class:`ql-clean`,type:`button`},e.ptm(`clean`)),null,16)],16)]})],16),b(`div`,t({ref:`editorElement`,class:e.cx(`content`),style:e.editorStyle},e.ptm(`content`)),null,16)],16)}X.render=Me;var Ne={key:0,class:`text-gray-500 dark:text-gray-400 text-xs`},Pe=[`innerHTML`],Fe={key:2,class:`text-gray-800 dark:text-gray-100`},Ie={class:`flex items-center gap-1`},Z={__name:`DataTableSection`,props:{items:Array,loading:Boolean,total:Number,columns:{type:Array,default:()=>[]},entity:String,orgId:String,onLoad:Function},emits:[`edit`,`delete`],setup(t){let s=t,{t:d}=P(),f=i(1),g=i(15),S=_(()=>(f.value-1)*g.value),C=_(()=>[...s.columns.map(e=>({type:`text`,width:`w-24`,headerWidth:`w-20`})),{type:`icons`,count:2,headerWidth:`w-16`}]);function w(e){f.value=e.page+1,g.value=e.rows,s.onLoad&&s.onLoad(f.value,g.value)}return n(()=>{s.onLoad&&s.onLoad(1,15)}),(n,i)=>{let s=a(`tooltip`);return t.loading?(o(),u(ee,{key:0,columns:C.value,rows:8},null,8,[`columns`])):(o(),u(l(H),{key:1,value:t.items,lazy:``,totalRecords:t.total,first:S.value,rows:g.value,onPage:w,paginator:``,paginatorTemplate:`FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown`,rowsPerPageOptions:[10,15,25,50],size:`small`,class:`!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden`},{empty:r(()=>[m(n.$slots,`empty`)]),default:r(()=>[(o(!0),x(y,null,e(t.columns,e=>(o(),u(l(W),{key:e.field,field:e.field,header:e.header,sortable:``},{body:r(({data:t})=>[e.field.startsWith(`_`)?(o(),x(`span`,Ne,v(t[e.field]||`-`),1)):h(``,!0),e.html?(o(),x(`div`,{key:1,class:`editor-content`,innerHTML:t[e.field]},null,8,Pe)):(o(),x(`span`,Fe,v(t[e.field]||`-`),1))]),_:2},1032,[`field`,`header`]))),128)),p(l(W),{header:l(d)(`common.actions`),style:{width:`90px`},frozen:``,alignFrozen:`right`},{body:r(({data:e})=>[b(`div`,Ie,[c(p(l(E),{icon:`pi pi-pencil`,size:`small`,text:``,severity:`secondary`,onClick:t=>n.$emit(`edit`,e)},null,8,[`onClick`]),[[s,l(d)(`common.edit`),void 0,{left:!0}]]),c(p(l(E),{icon:`pi pi-trash`,size:`small`,text:``,severity:`danger`,onClick:t=>n.$emit(`delete`,e)},null,8,[`onClick`]),[[s,l(d)(`common.delete`),void 0,{left:!0}]])])]),_:1},8,[`header`])]),_:3},8,[`value`,`totalRecords`,`first`,`rows`]))}}},Le={class:`space-y-4`},Re={key:0,class:`bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded p-3 text-xs text-red-700 dark:text-red-300`},Q={__name:`DialogForm`,props:{visible:Boolean,title:String,saving:Boolean,errors:{type:Object,default:()=>({})},width:{type:String,default:`480px`}},emits:[`save`,`cancel`],setup(t){let n=t,{t:i}=P(),a=_(()=>n.width===`maximize`?`90vw`:n.width);return(n,s)=>(o(),u(l(M),{visible:t.visible,"onUpdate:visible":s[2]||=e=>n.$emit(`update:visible`,e),header:t.title,modal:``,style:f({width:a.value}),class:`p-fluid`,closable:!t.saving},{footer:r(()=>[p(l(E),{label:l(i)(`common.cancel`),size:`small`,outlined:``,severity:`secondary`,disabled:t.saving,onClick:s[0]||=e=>n.$emit(`cancel`)},null,8,[`label`,`disabled`]),p(l(E),{label:l(i)(`common.save`),icon:`pi pi-check`,size:`small`,loading:t.saving,onClick:s[1]||=e=>n.$emit(`save`)},null,8,[`label`,`loading`])]),default:r(()=>[b(`div`,Le,[m(n.$slots,`default`),Object.keys(t.errors).length?(o(),x(`div`,Re,[(o(!0),x(y,null,e(t.errors,(e,t)=>(o(),x(`p`,{key:t,class:`mb-1`},[b(`strong`,null,v(t)+`:`,1),g(` `+v(Array.isArray(e)?e.join(`, `):e),1)]))),128))])):h(``,!0)])]),_:3},8,[`visible`,`header`,`style`,`closable`]))}},ze={class:`space-y-4`},Be={class:`flex items-center justify-between`},Ve={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},He={class:`text-sm text-gray-500 dark:text-gray-400`},Ue={class:`flex flex-col items-center justify-center py-10 text-gray-400`},We={class:`text-sm font-medium`},Ge=`/api/v1/tenant/job-management/responsibilities`,Ke={__name:`JobResponsibilitySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=P(),c=A(),u=i([]),d=i(!1),f=i(0),m=i(!1),g=i(!1),C=i(``),T=i(!1),D=i({}),O=i(!1),k=i(!1),j=i(``),M=i(null),N=i({main_task:``,activities:``,outputs:``,success_indicators:``}),I=_(()=>{let e=s(`job_management.responsibilities_title`);return g.value?`${e}`:`${s(`common.create`)} ${e}`}),R=_(()=>[{field:`main_task`,header:s(`job_management.main_task`),html:!0},{field:`activities`,header:s(`job_management.activities`),html:!0},{field:`outputs`,header:s(`job_management.outputs`),html:!0},{field:`success_indicators`,header:s(`job_management.success_indicators`),html:!0}]);async function z(e,t){d.value=!0;try{let r=await w.get(Ge,{params:{page:e,per_page:t,organization_id:n.orgId}}),i=r.data?.data||[];u.value=i.map(e=>({...e,main_task:e.main_task,activities:e.activities,outputs:e.outputs,success_indicators:e.success_indicators})),f.value=r.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function B(){g.value=!1,C.value=``,N.value={main_task:``,activities:``,outputs:``,success_indicators:``},D.value={},m.value=!0}function V(e){g.value=!0,C.value=e.id,N.value={main_task:e.main_task||``,activities:e.activities||``,outputs:e.outputs||``,success_indicators:e.success_indicators||``},D.value={},m.value=!0}async function H(){T.value=!0,D.value={};try{let e={nomenclature:n.orgName||``,full_code:n.orgCode||``,...N.value,organization_id:n.orgId};g.value?await w.put(`${Ge}/${C.value}`,e):await w.post(Ge,e),m.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),z(1,15)}catch(e){let t=F(e);Object.keys(t).length?D.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{T.value=!1}}function U(e){M.value=e,j.value=``,O.value=!0}async function W(){if(M.value){k.value=!0,j.value=``;try{await w.delete(`${Ge}/${M.value.id}`),O.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),z(1,15)}catch(e){j.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{k.value=!1}}}return(t,n)=>(o(),x(`div`,ze,[b(`div`,Be,[b(`div`,null,[b(`h2`,Ve,v(l(s)(`job_management.responsibilities_title`)),1),b(`p`,He,v(l(s)(`job_management.responsibilities_description`)),1)]),p(l(E),{label:l(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>B()},null,8,[`label`])]),p(Z,{items:u.value,loading:d.value,total:f.value,columns:R.value,entity:`responsibilities`,"org-id":e.orgId,"on-load":z,onEdit:V,onDelete:U},{empty:r(()=>[b(`div`,Ue,[n[9]||=b(`i`,{class:`pi pi-list-check text-3xl mb-2 opacity-50`},null,-1),b(`p`,We,v(l(s)(`job_management.empty_responsibilities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:m.value,"onUpdate:visible":n[5]||=e=>m.value=e,title:I.value,saving:T.value,errors:D.value,width:`maximize`,onSave:H,onCancel:n[6]||=e=>m.value=!1},{default:r(()=>[m.value?(o(),x(y,{key:0},[p(L,{label:l(s)(`job_management.main_task`),errors:D.value?.main_task},{default:r(()=>[p(l(X),{modelValue:N.value.main_task,"onUpdate:modelValue":n[1]||=e=>N.value.main_task=e,editorStyle:`height:120px`,class:S({"p-invalid":D.value?.main_task})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.activities`),errors:D.value?.activities},{default:r(()=>[p(l(X),{modelValue:N.value.activities,"onUpdate:modelValue":n[2]||=e=>N.value.activities=e,editorStyle:`height:120px`,class:S({"p-invalid":D.value?.activities})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.outputs`),errors:D.value?.outputs},{default:r(()=>[p(l(X),{modelValue:N.value.outputs,"onUpdate:modelValue":n[3]||=e=>N.value.outputs=e,editorStyle:`height:120px`,class:S({"p-invalid":D.value?.outputs})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.success_indicators`),errors:D.value?.success_indicators},{default:r(()=>[p(l(X),{modelValue:N.value.success_indicators,"onUpdate:modelValue":n[4]||=e=>N.value.success_indicators=e,editorStyle:`height:120px`,class:S({"p-invalid":D.value?.success_indicators})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])],64)):h(``,!0)]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:O.value,"onUpdate:visible":n[7]||=e=>O.value=e,loading:k.value,"error-msg":j.value,onConfirm:W,onCancel:n[8]||=e=>O.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},qe={class:`space-y-4`},Je={class:`flex items-center justify-between`},Ye={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Xe={class:`text-sm text-gray-500 dark:text-gray-400`},Ze={class:`flex flex-col items-center justify-center py-10 text-gray-400`},Qe={class:`text-sm font-medium`},$e=`/api/v1/tenant/job-management/hr-authorities`,et={__name:`JobHRAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=P(),c=A(),u=i([]),d=i(!1),f=i(0),m=i(!1),h=i(!1),g=i(``),y=i(!1),C=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({description:``}),M=_(()=>{let e=s(`job_management.hr_authorities`);return h.value?`${s(`common.edit`)} ${e}`:`${s(`common.create`)} ${e}`}),N=_(()=>[{field:`description`,header:s(`job_management.description`)}]);async function I(e,t){d.value=!0;try{let r=await w.get($e,{params:{page:e,per_page:t,organization_id:n.orgId}});u.value=r.data?.data||[],f.value=r.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function R(){h.value=!1,g.value=``,j.value={nomenclature:``,full_code:``,description:``},C.value={},m.value=!0}function B(e){h.value=!0,g.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},C.value={},m.value=!0}async function V(){y.value=!0,C.value={};try{let e={...j.value,nomenclature:n.orgName||``,full_code:n.orgCode||``,organization_id:n.orgId};h.value?await w.put(`${$e}/${g.value}`,e):await w.post($e,e),m.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),I(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{y.value=!1}}function H(e){k.value=e,O.value=``,T.value=!0}async function U(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${$e}/${k.value.id}`),T.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),I(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(t,n)=>(o(),x(`div`,qe,[b(`div`,Je,[b(`div`,null,[b(`h2`,Ye,v(l(s)(`job_management.hr_authorities`)),1),b(`p`,Xe,v(l(s)(`job_management.authority_description`)),1)]),p(l(E),{label:l(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>R()},null,8,[`label`])]),p(Z,{items:u.value,loading:d.value,total:f.value,columns:N.value,entity:`hr-authorities`,"org-id":e.orgId,"on-load":I,onEdit:B,onDelete:H},{empty:r(()=>[b(`div`,Ze,[n[6]||=b(`i`,{class:`pi pi-users text-3xl mb-2 opacity-50`},null,-1),b(`p`,Qe,v(l(s)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:m.value,"onUpdate:visible":n[2]||=e=>m.value=e,title:M.value,saving:y.value,errors:C.value,onSave:V,onCancel:n[3]||=e=>m.value=!1},{default:r(()=>[p(L,{label:l(s)(`job_management.description`),errors:C.value?.description},{default:r(()=>[p(l(z),{modelValue:j.value.description,"onUpdate:modelValue":n[1]||=e=>j.value.description=e,rows:`3`,class:S([`w-full`,{"p-invalid":C.value?.description}])},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":n[4]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:U,onCancel:n[5]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},tt={class:`space-y-4`},nt={class:`flex items-center justify-between`},rt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},it={class:`text-sm text-gray-500 dark:text-gray-400`},at={class:`flex flex-col items-center justify-center py-10 text-gray-400`},ot={class:`text-sm font-medium`},st=`/api/v1/tenant/job-management/operational-authorities`,ct={__name:`JobOpAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=P(),c=A(),u=i([]),d=i(!1),f=i(0),m=i(!1),h=i(!1),g=i(``),y=i(!1),C=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({description:``}),M=_(()=>{let e=s(`job_management.op_authorities`);return h.value?`${s(`common.edit`)} ${e}`:`${s(`common.create`)} ${e}`}),N=_(()=>[{field:`description`,header:s(`job_management.description`)}]);async function I(e,t){d.value=!0;try{let r=await w.get(st,{params:{page:e,per_page:t,organization_id:n.orgId}});u.value=r.data?.data||[],f.value=r.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function R(){h.value=!1,g.value=``,j.value={nomenclature:``,full_code:``,description:``},C.value={},m.value=!0}function B(e){h.value=!0,g.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},C.value={},m.value=!0}async function V(){y.value=!0,C.value={};try{let e={...j.value,nomenclature:n.orgName||``,full_code:n.orgCode||``,organization_id:n.orgId};h.value?await w.put(`${st}/${g.value}`,e):await w.post(st,e),m.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),I(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{y.value=!1}}function H(e){k.value=e,O.value=``,T.value=!0}async function U(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${st}/${k.value.id}`),T.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),I(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(t,n)=>(o(),x(`div`,tt,[b(`div`,nt,[b(`div`,null,[b(`h2`,rt,v(l(s)(`job_management.op_authorities`)),1),b(`p`,it,v(l(s)(`job_management.authority_description`)),1)]),p(l(E),{label:l(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>R()},null,8,[`label`])]),p(Z,{items:u.value,loading:d.value,total:f.value,columns:N.value,entity:`operational-authorities`,"org-id":e.orgId,"on-load":I,onEdit:B,onDelete:H},{empty:r(()=>[b(`div`,at,[n[6]||=b(`i`,{class:`pi pi-cog text-3xl mb-2 opacity-50`},null,-1),b(`p`,ot,v(l(s)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:m.value,"onUpdate:visible":n[2]||=e=>m.value=e,title:M.value,saving:y.value,errors:C.value,onSave:V,onCancel:n[3]||=e=>m.value=!1},{default:r(()=>[p(L,{label:l(s)(`job_management.description`),errors:C.value?.description},{default:r(()=>[p(l(z),{modelValue:j.value.description,"onUpdate:modelValue":n[1]||=e=>j.value.description=e,class:S([`w-full`,{"p-invalid":C.value?.description}]),rows:`3`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":n[4]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:U,onCancel:n[5]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},lt={class:`space-y-4`},ut={class:`flex items-center justify-between`},dt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},ft={class:`text-sm text-gray-500 dark:text-gray-400`},pt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},mt={class:`text-sm font-medium`},ht=`/api/v1/tenant/job-management/working-activities`,gt={__name:`JobActivitySection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=P(),c=A(),u=i([]),d=i(!1),f=i(0),m=i(!1),h=i(!1),g=i(``),y=i(!1),C=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({nomenclature:``,full_code:``,job_management_value_id:``});_(()=>Object.values(n.jobValueMap||{}).flat());let M=_(()=>[{field:`nomenclature`,header:s(`organization.nomenclature`)},{field:`full_code`,header:s(`organization.full_code`)}]);async function N(e,t){d.value=!0;try{let r=await w.get(ht,{params:{page:e,per_page:t,organization_id:n.orgId}});u.value=r.data?.data||[],f.value=r.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function I(){h.value=!1,g.value=``,j.value={nomenclature:``,full_code:``,job_management_value_id:``},C.value={},m.value=!0}function R(e){h.value=!0,g.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_id:e.job_management_value_id||``},C.value={},m.value=!0}async function z(){y.value=!0,C.value={};try{let e={...j.value,organization_id:n.orgId};h.value?await w.put(`${ht}/${g.value}`,e):await w.post(ht,e),m.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),N(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{y.value=!1}}function V(e){k.value=e,O.value=``,T.value=!0}async function H(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${ht}/${k.value.id}`),T.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),N(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(t,n)=>(o(),x(`div`,lt,[b(`div`,ut,[b(`div`,null,[b(`h2`,dt,v(l(s)(`job_management.activities`)),1),b(`p`,ft,v(l(s)(`job_management.activity_description`)),1)]),p(l(E),{label:l(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>I()},null,8,[`label`])]),p(Z,{items:u.value,loading:d.value,total:f.value,columns:M.value,entity:`working-activities`,"org-id":e.orgId,"on-load":N,onEdit:R,onDelete:V},{empty:r(()=>[b(`div`,pt,[n[8]||=b(`i`,{class:`pi pi-bolt text-3xl mb-2 opacity-50`},null,-1),b(`p`,mt,v(l(s)(`job_management.empty_activities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:m.value,"onUpdate:visible":n[4]||=e=>m.value=e,title:h.value?l(s)(`common.edit`):l(s)(`common.create`),saving:y.value,errors:C.value,onSave:z,onCancel:n[5]||=e=>m.value=!1},{default:r(()=>[p(L,{label:l(s)(`organization.nomenclature`),required:``,errors:C.value?.nomenclature},{default:r(()=>[p(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":n[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:S({"p-invalid":C.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`organization.full_code`),required:``,errors:C.value?.full_code},{default:r(()=>[p(B,{modelValue:j.value.full_code,"onUpdate:modelValue":n[2]||=e=>j.value.full_code=e,maxlength:`20`,class:S({"p-invalid":C.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.activity_type`),errors:C.value?.job_management_value_id},{default:r(()=>[p(G,{modelValue:j.value.job_management_value_id,"onUpdate:modelValue":n[3]||=e=>j.value.job_management_value_id=e,options:t.activityOptions,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":n[6]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:H,onCancel:n[7]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},_t={class:`space-y-4`},vt={class:`flex items-center justify-between`},yt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},bt={class:`text-sm text-gray-500 dark:text-gray-400`},xt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},St={class:`text-sm font-medium`},$=`/api/v1/tenant/job-management/working-risks`,Ct={__name:`JobRiskSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=P(),c=A(),u=i([]),d=i(!1),f=i(0),m=i(!1),h=i(!1),g=i(``),y=i(!1),C=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({nomenclature:``,full_code:``,job_management_value_environment_id:``,job_management_value_hazard_id:``}),M=_(()=>n.jobValueMap?.environment||[]),N=_(()=>n.jobValueMap?.risk||[]),I=_(()=>[{field:`nomenclature`,header:s(`organization.nomenclature`)},{field:`full_code`,header:s(`organization.full_code`)}]);async function R(e,t){d.value=!0;try{let r=await w.get($,{params:{page:e,per_page:t,organization_id:n.orgId}});u.value=r.data?.data||[],f.value=r.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function z(){h.value=!1,g.value=``,j.value={nomenclature:``,full_code:``,job_management_value_environment_id:``,job_management_value_hazard_id:``},C.value={},m.value=!0}function V(e){h.value=!0,g.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_environment_id:e.job_management_value_environment_id||``,job_management_value_hazard_id:e.job_management_value_hazard_id||``},C.value={},m.value=!0}async function H(){y.value=!0,C.value={};try{let e={...j.value,organization_id:n.orgId};h.value?await w.put(`${$}/${g.value}`,e):await w.post($,e),m.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),R(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{y.value=!1}}function U(e){k.value=e,O.value=``,T.value=!0}async function W(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${$}/${k.value.id}`),T.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),R(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(t,n)=>(o(),x(`div`,_t,[b(`div`,vt,[b(`div`,null,[b(`h2`,yt,v(l(s)(`job_management.risks`)),1),b(`p`,bt,v(l(s)(`job_management.risk_description`)),1)]),p(l(E),{label:l(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>z()},null,8,[`label`])]),p(Z,{items:u.value,loading:d.value,total:f.value,columns:I.value,entity:`working-risks`,"org-id":e.orgId,"on-load":R,onEdit:V,onDelete:U},{empty:r(()=>[b(`div`,xt,[n[9]||=b(`i`,{class:`pi pi-exclamation-triangle text-3xl mb-2 opacity-50`},null,-1),b(`p`,St,v(l(s)(`job_management.empty_risks`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:m.value,"onUpdate:visible":n[5]||=e=>m.value=e,title:h.value?l(s)(`common.edit`):l(s)(`common.create`),saving:y.value,errors:C.value,onSave:H,onCancel:n[6]||=e=>m.value=!1},{default:r(()=>[p(L,{label:l(s)(`organization.nomenclature`),required:``,errors:C.value?.nomenclature},{default:r(()=>[p(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":n[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:S({"p-invalid":C.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`organization.full_code`),required:``,errors:C.value?.full_code},{default:r(()=>[p(B,{modelValue:j.value.full_code,"onUpdate:modelValue":n[2]||=e=>j.value.full_code=e,maxlength:`20`,class:S({"p-invalid":C.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.environment_risk`),errors:C.value?.job_management_value_environment_id},{default:r(()=>[p(G,{modelValue:j.value.job_management_value_environment_id,"onUpdate:modelValue":n[3]||=e=>j.value.job_management_value_environment_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.risk`),errors:C.value?.job_management_value_hazard_id},{default:r(()=>[p(G,{modelValue:j.value.job_management_value_hazard_id,"onUpdate:modelValue":n[4]||=e=>j.value.job_management_value_hazard_id=e,options:N.value,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":n[7]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:W,onCancel:n[8]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},wt={class:`space-y-4`},Tt={class:`flex items-center justify-between`},Et={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Dt={class:`text-sm text-gray-500 dark:text-gray-400`},Ot={class:`flex flex-col items-center justify-center py-10 text-gray-400`},kt={class:`text-sm font-medium`},At=`/api/v1/tenant/job-management/relationships`,jt={__name:`JobRelationshipSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=P(),c=A(),u=i([]),d=i(!1),f=i(0),m=i(!1),h=i(!1),g=i(``),y=i(!1),C=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({nomenclature:``,full_code:``,job_management_value_relationship_id:``,job_management_value_frequency_id:``}),M=_(()=>n.jobValueMap?.relationship||[]),N=_(()=>n.jobValueMap?.frequency||[]),I=_(()=>[{field:`nomenclature`,header:s(`organization.nomenclature`)},{field:`full_code`,header:s(`organization.full_code`)}]);async function R(e,t){d.value=!0;try{let r=await w.get(At,{params:{page:e,per_page:t,organization_id:n.orgId}});u.value=r.data?.data||[],f.value=r.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function z(){h.value=!1,g.value=``,j.value={nomenclature:``,full_code:``,job_management_value_relationship_id:``,job_management_value_frequency_id:``},C.value={},m.value=!0}function V(e){h.value=!0,g.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_relationship_id:e.job_management_value_relationship_id||``,job_management_value_frequency_id:e.job_management_value_frequency_id||``},C.value={},m.value=!0}async function H(){y.value=!0,C.value={};try{let e={...j.value,organization_id:n.orgId};h.value?await w.put(`${At}/${g.value}`,e):await w.post(At,e),m.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),R(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{y.value=!1}}function U(e){k.value=e,O.value=``,T.value=!0}async function W(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${At}/${k.value.id}`),T.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),R(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(t,n)=>(o(),x(`div`,wt,[b(`div`,Tt,[b(`div`,null,[b(`h2`,Et,v(l(s)(`job_management.relationships`)),1),b(`p`,Dt,v(l(s)(`job_management.relationship_description`)),1)]),p(l(E),{label:l(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>z()},null,8,[`label`])]),p(Z,{items:u.value,loading:d.value,total:f.value,columns:I.value,entity:`relationships`,"org-id":e.orgId,"on-load":R,onEdit:V,onDelete:U},{empty:r(()=>[b(`div`,Ot,[n[9]||=b(`i`,{class:`pi pi-share-alt text-3xl mb-2 opacity-50`},null,-1),b(`p`,kt,v(l(s)(`job_management.empty_relationships`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:m.value,"onUpdate:visible":n[5]||=e=>m.value=e,title:h.value?l(s)(`common.edit`):l(s)(`common.create`),saving:y.value,errors:C.value,onSave:H,onCancel:n[6]||=e=>m.value=!1},{default:r(()=>[p(L,{label:l(s)(`organization.nomenclature`),required:``,errors:C.value?.nomenclature},{default:r(()=>[p(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":n[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:S({"p-invalid":C.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`organization.full_code`),required:``,errors:C.value?.full_code},{default:r(()=>[p(B,{modelValue:j.value.full_code,"onUpdate:modelValue":n[2]||=e=>j.value.full_code=e,maxlength:`20`,class:S({"p-invalid":C.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.relationship_type`),errors:C.value?.job_management_value_relationship_id},{default:r(()=>[p(G,{modelValue:j.value.job_management_value_relationship_id,"onUpdate:modelValue":n[3]||=e=>j.value.job_management_value_relationship_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.frequency`),errors:C.value?.job_management_value_frequency_id},{default:r(()=>[p(G,{modelValue:j.value.job_management_value_frequency_id,"onUpdate:modelValue":n[4]||=e=>j.value.job_management_value_frequency_id=e,options:N.value,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":n[7]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:W,onCancel:n[8]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Mt={class:`space-y-4`},Nt={class:`flex items-center justify-between`},Pt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Ft={class:`text-sm text-gray-500 dark:text-gray-400`},It={class:`flex flex-col items-center justify-center py-10 text-gray-400`},Lt={class:`text-sm font-medium`},Rt=`/api/v1/tenant/job-management/subordinate-controls`,zt={__name:`JobSubordinateSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=P(),c=A(),u=i([]),d=i(!1),f=i(0),m=i(!1),h=i(!1),g=i(``),y=i(!1),C=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({nomenclature:``,full_code:``,job_management_value_id:``}),M=_(()=>Object.values(n.jobValueMap||{}).flat()),N=_(()=>[{field:`nomenclature`,header:s(`organization.nomenclature`)},{field:`full_code`,header:s(`organization.full_code`)}]);async function I(e,t){d.value=!0;try{let r=await w.get(Rt,{params:{page:e,per_page:t,organization_id:n.orgId}});u.value=r.data?.data||[],f.value=r.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function R(){h.value=!1,g.value=``,j.value={nomenclature:``,full_code:``,job_management_value_id:``},C.value={},m.value=!0}function z(e){h.value=!0,g.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_id:e.job_management_value_id||``},C.value={},m.value=!0}async function V(){y.value=!0,C.value={};try{let e={...j.value,organization_id:n.orgId};h.value?await w.put(`${Rt}/${g.value}`,e):await w.post(Rt,e),m.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),I(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{y.value=!1}}function H(e){k.value=e,O.value=``,T.value=!0}async function U(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${Rt}/${k.value.id}`),T.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),I(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(t,n)=>(o(),x(`div`,Mt,[b(`div`,Nt,[b(`div`,null,[b(`h2`,Pt,v(l(s)(`job_management.subordinate_controls`)),1),b(`p`,Ft,v(l(s)(`job_management.subordinate_description`)),1)]),p(l(E),{label:l(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>R()},null,8,[`label`])]),p(Z,{items:u.value,loading:d.value,total:f.value,columns:N.value,entity:`subordinate-controls`,"org-id":e.orgId,"on-load":I,onEdit:z,onDelete:H},{empty:r(()=>[b(`div`,It,[n[8]||=b(`i`,{class:`pi pi-sitemap text-3xl mb-2 opacity-50`},null,-1),b(`p`,Lt,v(l(s)(`job_management.empty_subordinates`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:m.value,"onUpdate:visible":n[4]||=e=>m.value=e,title:h.value?l(s)(`common.edit`):l(s)(`common.create`),saving:y.value,errors:C.value,onSave:V,onCancel:n[5]||=e=>m.value=!1},{default:r(()=>[p(L,{label:l(s)(`organization.nomenclature`),required:``,errors:C.value?.nomenclature},{default:r(()=>[p(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":n[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:S({"p-invalid":C.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`organization.full_code`),required:``,errors:C.value?.full_code},{default:r(()=>[p(B,{modelValue:j.value.full_code,"onUpdate:modelValue":n[2]||=e=>j.value.full_code=e,maxlength:`20`,class:S({"p-invalid":C.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.control_type`),errors:C.value?.job_management_value_id},{default:r(()=>[p(G,{modelValue:j.value.job_management_value_id,"onUpdate:modelValue":n[3]||=e=>j.value.job_management_value_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":n[6]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:U,onCancel:n[7]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Bt={class:`space-y-4`},Vt={class:`flex items-center justify-between`},Ht={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Ut={class:`text-sm text-gray-500 dark:text-gray-400`},Wt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},Gt={class:`text-sm font-medium`},Kt=`/api/v1/tenant/job-management/assets`,qt={__name:`JobAssetSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=P(),c=A(),u=i([]),d=i(!1),f=i(0),m=i(!1),h=i(!1),g=i(``),y=i(!1),C=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({nomenclature:``,full_code:``,job_management_value_asset_id:``,job_management_value_authority_id:``}),M=_(()=>n.jobValueMap?.asset||[]),N=_(()=>n.jobValueMap?.authority||[]),I=_(()=>[{field:`nomenclature`,header:s(`organization.nomenclature`)},{field:`full_code`,header:s(`organization.full_code`)}]);async function R(e,t){d.value=!0;try{let r=await w.get(Kt,{params:{page:e,per_page:t,organization_id:n.orgId}});u.value=r.data?.data||[],f.value=r.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function z(){h.value=!1,g.value=``,j.value={nomenclature:``,full_code:``,job_management_value_asset_id:``,job_management_value_authority_id:``},C.value={},m.value=!0}function V(e){h.value=!0,g.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_asset_id:e.job_management_value_asset_id||``,job_management_value_authority_id:e.job_management_value_authority_id||``},C.value={},m.value=!0}async function H(){y.value=!0,C.value={};try{let e={...j.value,organization_id:n.orgId};h.value?await w.put(`${Kt}/${g.value}`,e):await w.post(Kt,e),m.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),R(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{y.value=!1}}function U(e){k.value=e,O.value=``,T.value=!0}async function W(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${Kt}/${k.value.id}`),T.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),R(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(t,n)=>(o(),x(`div`,Bt,[b(`div`,Vt,[b(`div`,null,[b(`h2`,Ht,v(l(s)(`job_management.assets`)),1),b(`p`,Ut,v(l(s)(`job_management.asset_description`)),1)]),p(l(E),{label:l(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>z()},null,8,[`label`])]),p(Z,{items:u.value,loading:d.value,total:f.value,columns:I.value,entity:`assets`,"org-id":e.orgId,"on-load":R,onEdit:V,onDelete:U},{empty:r(()=>[b(`div`,Wt,[n[9]||=b(`i`,{class:`pi pi-box text-3xl mb-2 opacity-50`},null,-1),b(`p`,Gt,v(l(s)(`job_management.empty_assets`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:m.value,"onUpdate:visible":n[5]||=e=>m.value=e,title:h.value?l(s)(`common.edit`):l(s)(`common.create`),saving:y.value,errors:C.value,onSave:H,onCancel:n[6]||=e=>m.value=!1},{default:r(()=>[p(L,{label:l(s)(`organization.nomenclature`),required:``,errors:C.value?.nomenclature},{default:r(()=>[p(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":n[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:S({"p-invalid":C.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`organization.full_code`),required:``,errors:C.value?.full_code},{default:r(()=>[p(B,{modelValue:j.value.full_code,"onUpdate:modelValue":n[2]||=e=>j.value.full_code=e,maxlength:`20`,class:S({"p-invalid":C.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.asset_type`),errors:C.value?.job_management_value_asset_id},{default:r(()=>[p(G,{modelValue:j.value.job_management_value_asset_id,"onUpdate:modelValue":n[3]||=e=>j.value.job_management_value_asset_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.authority_level`),errors:C.value?.job_management_value_authority_id},{default:r(()=>[p(G,{modelValue:j.value.job_management_value_authority_id,"onUpdate:modelValue":n[4]||=e=>j.value.job_management_value_authority_id=e,options:N.value,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":n[7]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:W,onCancel:n[8]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Jt={class:`space-y-4`},Yt={class:`flex items-center justify-between`},Xt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Zt={class:`text-sm text-gray-500 dark:text-gray-400`},Qt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},$t={class:`text-sm font-medium`},en=`/api/v1/tenant/job-management/financials`,tn={__name:`JobFinancialSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=P(),c=A(),u=i([]),d=i(!1),f=i(0),m=i(!1),h=i(!1),g=i(``),y=i(!1),C=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({nomenclature:``,full_code:``,is_authorized:!1,job_management_value_cash_id:``,job_management_value_authority_id:``,job_management_value_impact_id:``}),M=_(()=>n.jobValueMap?.cash||[]),N=_(()=>n.jobValueMap?.authority||[]),I=_(()=>n.jobValueMap?.impact||[]),R=_(()=>[{field:`nomenclature`,header:s(`organization.nomenclature`)},{field:`full_code`,header:s(`organization.full_code`)},{field:`is_authorized`,header:s(`job_management.is_authorized`)}]);async function z(e,t){d.value=!0;try{let r=await w.get(en,{params:{page:e,per_page:t,organization_id:n.orgId}});u.value=r.data?.data||[],f.value=r.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function V(){h.value=!1,g.value=``,j.value={nomenclature:``,full_code:``,is_authorized:!1,job_management_value_cash_id:``,job_management_value_authority_id:``,job_management_value_impact_id:``},C.value={},m.value=!0}function H(e){h.value=!0,g.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,is_authorized:!!e.is_authorized,job_management_value_cash_id:e.job_management_value_cash_id||``,job_management_value_authority_id:e.job_management_value_authority_id||``,job_management_value_impact_id:e.job_management_value_impact_id||``},C.value={},m.value=!0}async function U(){y.value=!0,C.value={};try{let e={...j.value,organization_id:n.orgId};h.value?await w.put(`${en}/${g.value}`,e):await w.post(en,e),m.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),z(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{y.value=!1}}function W(e){k.value=e,O.value=``,T.value=!0}async function ee(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${en}/${k.value.id}`),T.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),z(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(t,n)=>(o(),x(`div`,Jt,[b(`div`,Yt,[b(`div`,null,[b(`h2`,Xt,v(l(s)(`job_management.financials`)),1),b(`p`,Zt,v(l(s)(`job_management.financial_description`)),1)]),p(l(E),{label:l(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>V()},null,8,[`label`])]),p(Z,{items:u.value,loading:d.value,total:f.value,columns:R.value,entity:`financials`,"org-id":e.orgId,"on-load":z,onEdit:H,onDelete:W},{empty:r(()=>[b(`div`,Qt,[n[11]||=b(`i`,{class:`pi pi-money-bill text-3xl mb-2 opacity-50`},null,-1),b(`p`,$t,v(l(s)(`job_management.empty_financials`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:m.value,"onUpdate:visible":n[7]||=e=>m.value=e,title:h.value?l(s)(`common.edit`):l(s)(`common.create`),saving:y.value,errors:C.value,onSave:U,onCancel:n[8]||=e=>m.value=!1},{default:r(()=>[p(L,{label:l(s)(`organization.nomenclature`),required:``,errors:C.value?.nomenclature},{default:r(()=>[p(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":n[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:S({"p-invalid":C.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`organization.full_code`),required:``,errors:C.value?.full_code},{default:r(()=>[p(B,{modelValue:j.value.full_code,"onUpdate:modelValue":n[2]||=e=>j.value.full_code=e,maxlength:`20`,class:S({"p-invalid":C.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.is_authorized`),class:`md:col-span-2`},{default:r(()=>[p(l(te),{modelValue:j.value.is_authorized,"onUpdate:modelValue":n[3]||=e=>j.value.is_authorized=e},null,8,[`modelValue`])]),_:1},8,[`label`]),p(L,{label:l(s)(`job_management.cash_level`),errors:C.value?.job_management_value_cash_id},{default:r(()=>[p(G,{modelValue:j.value.job_management_value_cash_id,"onUpdate:modelValue":n[4]||=e=>j.value.job_management_value_cash_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.authority_level`),errors:C.value?.job_management_value_authority_id},{default:r(()=>[p(G,{modelValue:j.value.job_management_value_authority_id,"onUpdate:modelValue":n[5]||=e=>j.value.job_management_value_authority_id=e,options:N.value,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.impact_level`),errors:C.value?.job_management_value_impact_id},{default:r(()=>[p(G,{modelValue:j.value.job_management_value_impact_id,"onUpdate:modelValue":n[6]||=e=>j.value.job_management_value_impact_id=e,options:I.value,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":n[9]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:ee,onCancel:n[10]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},nn={class:`space-y-4`},rn={class:`flex items-center justify-between`},an={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},on={class:`text-sm text-gray-500 dark:text-gray-400`},sn={class:`flex flex-col items-center justify-center py-10 text-gray-400`},cn={class:`text-sm font-medium`},ln=`/api/v1/tenant/job-management/potency-competencies`,un={__name:`JobPotencySection`,props:{orgId:String,jobValueMap:Object,competencyOptions:Array},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=P(),c=A(),u=i([]),d=i(!1),f=i(0),m=i(!1),h=i(!1),g=i(``),y=i(!1),C=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({competency_id:``,job_management_value_id:``,weight:null}),M=_(()=>Object.values(n.jobValueMap||{}).flat()),N=_(()=>[{field:`_competency`,header:s(`job_management.competency`)},{field:`weight`,header:s(`job_management.weight`)}]);async function I(e,t){d.value=!0;try{let r=await w.get(ln,{params:{page:e,per_page:t,organization_id:n.orgId}});u.value=(r.data?.data||[]).map(e=>{let t=n.competencyOptions?.find(t=>t.value===e.competency_id);return{...e,_competency:t?.label||e.competency_id}}),f.value=r.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function R(){h.value=!1,g.value=``,j.value={competency_id:``,job_management_value_id:``,weight:null},C.value={},m.value=!0}function z(e){h.value=!0,g.value=e.id,j.value={competency_id:e.competency_id||``,job_management_value_id:e.job_management_value_id||``,weight:e.weight??null},C.value={},m.value=!0}async function B(){y.value=!0,C.value={};try{let e={...j.value,organization_id:n.orgId};h.value?await w.put(`${ln}/${g.value}`,e):await w.post(ln,e),m.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),I(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{y.value=!1}}function H(e){k.value=e,O.value=``,T.value=!0}async function U(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${ln}/${k.value.id}`),T.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),I(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(t,n)=>(o(),x(`div`,nn,[b(`div`,rn,[b(`div`,null,[b(`h2`,an,v(l(s)(`job_management.potency_competencies`)),1),b(`p`,on,v(l(s)(`job_management.potency_description`)),1)]),p(l(E),{label:l(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>R()},null,8,[`label`])]),p(Z,{items:u.value,loading:d.value,total:f.value,columns:N.value,entity:`potency-competencies`,"org-id":e.orgId,"on-load":I,onEdit:z,onDelete:H},{empty:r(()=>[b(`div`,sn,[n[8]||=b(`i`,{class:`pi pi-star text-3xl mb-2 opacity-50`},null,-1),b(`p`,cn,v(l(s)(`job_management.empty_potency`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:m.value,"onUpdate:visible":n[4]||=e=>m.value=e,title:h.value?l(s)(`common.edit`):l(s)(`common.create`),saving:y.value,errors:C.value,onSave:B,onCancel:n[5]||=e=>m.value=!1},{default:r(()=>[p(L,{label:l(s)(`job_management.competency`),required:``,errors:C.value?.competency_id},{default:r(()=>[p(G,{modelValue:j.value.competency_id,"onUpdate:modelValue":n[1]||=e=>j.value.competency_id=e,options:e.competencyOptions,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.value_ref`),errors:C.value?.job_management_value_id},{default:r(()=>[p(G,{modelValue:j.value.job_management_value_id,"onUpdate:modelValue":n[2]||=e=>j.value.job_management_value_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.weight`),errors:C.value?.weight},{default:r(()=>[p(l(V),{modelValue:j.value.weight,"onUpdate:modelValue":n[3]||=e=>j.value.weight=e,min:0,max:100,class:S([{"p-invalid":C.value?.weight},`w-full`]),size:`small`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":n[6]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:U,onCancel:n[7]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},dn={class:`space-y-4`},fn={class:`flex items-center justify-between`},pn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},mn={class:`text-sm text-gray-500 dark:text-gray-400`},hn={class:`flex flex-col items-center justify-center py-10 text-gray-400`},gn={class:`text-sm font-medium`},_n=`/api/v1/tenant/job-management/competency-groups`,vn={__name:`JobCompetencyGroupSection`,props:{orgId:String},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=P(),c=A(),u=i([]),d=i(!1);i(0);let f=i(!1),m=i(!1),h=i(``),g=i(!1),y=i({}),C=i(!1),T=i(!1),D=i(``),O=i(null),k=i({category:``,weight:null}),j=_(()=>[{label:`${s(`job_management.technical`)} (${s(`job_management.category`)})`,value:`technical`},{label:`${s(`job_management.managerial`)} (${s(`job_management.category`)})`,value:`managerial`}]),M=_(()=>[{field:`category`,header:s(`job_management.category`)},{field:`weight`,header:s(`job_management.weight`)}]);async function N(){d.value=!0;try{let e=await w.get(_n,{params:{organization_id:n.orgId}});u.value=e.data?.data||(Array.isArray(e.data)?e.data:[])}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function I(){m.value=!1,h.value=``,k.value={category:`technical`,weight:null},y.value={},f.value=!0}function R(e){m.value=!0,h.value=e.id,k.value={category:e.category||`technical`,weight:e.weight??null},y.value={},f.value=!0}async function z(){g.value=!0,y.value={};try{let e={...k.value,organization_id:n.orgId};m.value?await w.put(`${_n}/${h.value}`,e):await w.post(_n,e),f.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),N()}catch(e){let t=F(e);Object.keys(t).length?y.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{g.value=!1}}function B(e){O.value=e,D.value=``,C.value=!0}async function H(){if(O.value){T.value=!0,D.value=``;try{await w.delete(`${_n}/${O.value.id}`),C.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),N()}catch(e){D.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{T.value=!1}}}return(t,n)=>(o(),x(`div`,dn,[b(`div`,fn,[b(`div`,null,[b(`h2`,pn,v(l(s)(`job_management.competency_groups`)),1),b(`p`,mn,v(l(s)(`job_management.competency_group_description`)),1)]),p(l(E),{label:l(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>I()},null,8,[`label`])]),p(Z,{items:u.value,loading:d.value,total:u.value.length,columns:M.value,entity:`competency-groups`,"org-id":e.orgId,"on-load":N,onEdit:R,onDelete:B},{empty:r(()=>[b(`div`,hn,[n[7]||=b(`i`,{class:`pi pi-chart-pie text-3xl mb-2 opacity-50`},null,-1),b(`p`,gn,v(l(s)(`job_management.empty_competency_groups`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:f.value,"onUpdate:visible":n[3]||=e=>f.value=e,title:m.value?l(s)(`common.edit`):l(s)(`common.create`),saving:g.value,errors:y.value,onSave:z,onCancel:n[4]||=e=>f.value=!1},{default:r(()=>[p(L,{label:l(s)(`job_management.category`),required:``,errors:y.value?.category},{default:r(()=>[p(G,{modelValue:k.value.category,"onUpdate:modelValue":n[1]||=e=>k.value.category=e,options:j.value,optionLabel:`label`,optionValue:`value`,placeholder:l(s)(`common.select`),class:S({"p-invalid":y.value?.category})},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:l(s)(`job_management.weight`),required:``,errors:y.value?.weight},{default:r(()=>[p(l(V),{modelValue:k.value.weight,"onUpdate:modelValue":n[2]||=e=>k.value.weight=e,min:0,max:100,suffix:`%`,class:S([{"p-invalid":y.value?.weight},`w-full`]),size:`small`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:C.value,"onUpdate:visible":n[5]||=e=>C.value=e,loading:T.value,"error-msg":D.value,onConfirm:H,onCancel:n[6]||=e=>C.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},yn={class:`space-y-6`},bn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},xn={class:`text-sm text-gray-500 dark:text-gray-400`},Sn={key:0,class:`flex items-center justify-center py-12`},Cn={class:`grid grid-cols-1 md:grid-cols-3 gap-4`},wn={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Tn={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},En={class:`text-2xl font-bold text-emerald-600 dark:text-emerald-400`},Dn={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},On={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},kn={class:`text-2xl font-bold text-blue-600 dark:text-blue-400`},An={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},jn={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},Mn={key:0,class:`text-[10px] text-gray-400 mt-2`},Nn={key:0,class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden`},Pn={class:`px-5 py-3 border-b border-gray-200 dark:border-gray-700 font-semibold text-sm text-gray-700 dark:text-gray-300`},Fn={class:`p-5`},In={class:`text-sm text-gray-700 dark:text-gray-300 capitalize`},Ln={class:`text-sm font-semibold text-gray-900 dark:text-gray-100`},Rn={key:2},zn={class:`flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500`},Bn={class:`text-sm font-medium`},Vn={class:`text-xs mt-1`},Hn={class:`flex justify-end gap-3`},Un=`/api/v1/tenant/job-management/scores`,Wn={__name:`JobScoreSection`,props:{orgId:String},emits:[`saved`],setup(t,{emit:r}){let a=t,s=r,{t:c}=P(),d=A(),f=i(!1),m=i(!1),g=i(null),S=_(()=>{if(!g.value?.components)return null;try{return JSON.parse(g.value.components)}catch{return null}});function C(e){return e?.toLocaleString?.(`id-ID`)??`-`}async function T(){if(a.orgId){f.value=!0;try{let e=await w.get(`${Un}/${a.orgId}`);g.value=e.data?.data||null,s(`saved`)}catch{g.value=null}finally{f.value=!1}}}async function D(){if(a.orgId){m.value=!0;try{let e=await w.put(`${Un}/${a.orgId}`,{components:null});g.value=e.data?.data||null,d.add({severity:`success`,summary:c(`message.success`),detail:c(`job_management.score_calculated`),life:2e3})}catch(e){d.add({severity:`error`,summary:c(`message.error`),detail:e.response?.data?.error?.message||c(`message.operation_failed`),life:4e3})}finally{m.value=!1}}}return n(T),(t,n)=>(o(),x(`div`,yn,[b(`div`,null,[b(`h2`,bn,v(l(c)(`job_management.scores`)),1),b(`p`,xn,v(l(c)(`job_management.score_description`)),1)]),f.value?(o(),x(`div`,Sn,[...n[0]||=[b(`i`,{class:`pi pi-spin pi-spinner text-emerald-500 text-2xl`},null,-1)]])):g.value?(o(),x(y,{key:1},[b(`div`,Cn,[b(`div`,wn,[b(`div`,Tn,v(l(c)(`job_management.value_with_financial`)),1),b(`div`,En,v(C(g.value.job_value_with_financial)),1)]),b(`div`,Dn,[b(`div`,On,v(l(c)(`job_management.value_without_financial`)),1),b(`div`,kn,v(C(g.value.job_value_without_financial)),1)]),b(`div`,An,[b(`div`,jn,v(l(c)(`job_management.has_financial_authority`)),1),p(l(I),{value:g.value.has_financial_authority?l(c)(`common.yes`):l(c)(`common.no`),severity:g.value.has_financial_authority?`success`:`danger`,class:`!text-xs`},null,8,[`value`,`severity`]),g.value.calculated_at?(o(),x(`div`,Mn,v(l(c)(`job_management.calculated_at`))+`: `+v(g.value.calculated_at),1)):h(``,!0)])]),S.value?(o(),x(`div`,Nn,[b(`div`,Pn,v(l(c)(`job_management.component_breakdown`)),1),b(`div`,Fn,[(o(!0),x(y,null,e(S.value,(e,t)=>(o(),x(`div`,{key:t,class:`flex items-center justify-between py-1.5 border-b border-gray-100 dark:border-gray-700 last:border-0`},[b(`span`,In,v(t.replace(/_/g,` `)),1),b(`span`,Ln,v(C(e)),1)]))),128))])])):h(``,!0)],64)):(o(),x(`div`,Rn,[b(`div`,zn,[n[1]||=b(`i`,{class:`pi pi-calculator text-4xl mb-3 opacity-50`},null,-1),b(`p`,Bn,v(l(c)(`job_management.no_score`)),1),b(`p`,Vn,v(l(c)(`job_management.score_hint`)),1)])])),b(`div`,Hn,[p(l(E),{label:l(c)(`common.refresh`),icon:`pi pi-refresh`,size:`small`,text:``,onClick:T},null,8,[`label`]),g.value?(o(),u(l(E),{key:0,label:l(c)(`job_management.recalculate`),icon:`pi pi-calculator`,size:`small`,severity:`info`,loading:m.value,onClick:D},null,8,[`label`,`loading`])):h(``,!0)])]))}},Gn={class:`max-w-full mx-auto`},Kn={key:0,class:`flex gap-6`},qn={class:`w-56 space-y-2`},Jn={class:`flex-1 space-y-3`},Yn={key:1,class:`flex gap-6`},Xn={class:`w-56 shrink-0 space-y-1`},Zn=[`onClick`,`onKeydown`],Qn={key:0,class:`pi pi-check text-xs`},$n={class:`flex-1 min-w-0`},er={key:0,class:`pi pi-check-circle text-emerald-400 text-xs shrink-0`},tr={class:`flex-1 min-w-0`},nr={__name:`JobManagementForm`,setup(t){let r=N(),a=k(),{t:c}=P(),d=A(),f=a.query.org_id||``,p=i(0),m=i(!0),g=i(Array(15).fill(!1)),C=i(``),E=i(``),D=i(``),O=i(``),j=i([]),M=i([]),F=i([]),I=i({}),L=i([]),R=[{labelKey:`job_management.identifications`,icon:`pi pi-id-card`,comp:le},{labelKey:`job_management.objectives`,icon:`pi pi-bullseye`,comp:ge},{labelKey:`job_management.responsibilities_title`,icon:`pi pi-list-check`,comp:Ke},{labelKey:`job_management.education_experience`,icon:`pi pi-graduation-cap`,comp:Ce},{labelKey:`job_management.hr_authorities`,icon:`pi pi-users`,comp:et},{labelKey:`job_management.op_authorities`,icon:`pi pi-cog`,comp:ct},{labelKey:`job_management.activities`,icon:`pi pi-bolt`,comp:gt},{labelKey:`job_management.risks`,icon:`pi pi-exclamation-triangle`,comp:Ct},{labelKey:`job_management.relationships`,icon:`pi pi-share-alt`,comp:jt},{labelKey:`job_management.subordinate_controls`,icon:`pi pi-sitemap`,comp:zt},{labelKey:`job_management.assets`,icon:`pi pi-box`,comp:qt},{labelKey:`job_management.financials`,icon:`pi pi-money-bill`,comp:tn},{labelKey:`job_management.potency_competencies`,icon:`pi pi-star`,comp:un},{labelKey:`job_management.competency_groups`,icon:`pi pi-chart-pie`,comp:vn},{labelKey:`job_management.scores`,icon:`pi pi-calculator`,comp:Wn}],z=_(()=>R[p.value]?.comp||null);function B(e){return p.value===e?`bg-emerald-50 dark:bg-emerald-900/20 ring-1 ring-emerald-300 dark:ring-emerald-700`:(g.value[e],`hover:bg-gray-50 dark:hover:bg-gray-800`)}function V(e){return p.value===e?`bg-emerald-600 text-white`:g.value[e]?`bg-emerald-100 dark:bg-emerald-800 text-emerald-600 dark:text-emerald-300`:`bg-gray-200 dark:bg-gray-600 text-gray-600 dark:text-gray-300`}function H(e){return p.value===e?`text-emerald-700 dark:text-emerald-300`:g.value[e]?`text-emerald-600 dark:text-emerald-400`:`text-gray-700 dark:text-gray-300`}function U(e){p.value=e,r.replace({query:{...a.query,section:String(e)}})}function W(e){typeof e==`number`&&(g.value[e]=!0)}async function G(){if(f)try{let e=(await w.get(`/api/v1/tenant/organizations/${f}`)).data?.data;e&&(C.value=e.nomenclature||``,E.value=e.full_code||e.code||``,D.value=e.grading_id||``,O.value=e.job_family_id||``)}catch{}}async function K(){try{let[e,t,n,r]=await Promise.all([w.get(`/api/v1/tenant/settings/gradings?per_page=100`),w.get(`/api/v1/tenant/job-management/values?per_page=200`),w.get(`/api/v1/tenant/competencies?per_page=200`).catch(()=>({data:{data:[]}})),w.get(`/api/v1/tenant/settings/job-families?per_page=100`)]);j.value=(e.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id})),M.value=(r.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id}));let i=t.data?.data||[];F.value=i.map(e=>({label:`${e.type}${e.level?` Lv.`+e.level:``}${e.descriptions?` — `+e.descriptions:``}`,value:e.id,type:e.type,level:e.level,descriptions:e.descriptions}));let a={};i.forEach(e=>{a[e.type]||(a[e.type]=[]),a[e.type].push({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level})}),I.value=a,L.value=(n.data?.data||[]).map(e=>({label:e.name||e.code,value:e.id}))}catch{}}return n(async()=>{try{await Promise.all([G(),K()]);let e=parseInt(a.query.section);!isNaN(e)&&e>=0&&e<R.length&&(p.value=e)}catch(e){d.add({severity:`error`,summary:c(`message.error`),detail:e.response?.data?.error?.message||c(`message.failed_to_load`),life:4e3})}finally{m.value=!1}}),(t,n)=>(o(),x(`div`,Gn,[m.value?(o(),x(`div`,Kn,[b(`div`,qn,[(o(),x(y,null,e(8,e=>b(`div`,{key:e,class:`h-12 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))]),b(`div`,Jn,[(o(),x(y,null,e(6,e=>b(`div`,{key:e,class:`h-9 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))])])):(o(),x(`div`,Yn,[b(`div`,Xn,[(o(),x(y,null,e(R,(e,t)=>b(`div`,{key:t,role:`button`,tabindex:0,class:S([`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150 cursor-pointer select-none`,B(t)]),onClick:e=>U(t),onKeydown:T(e=>U(t),[`enter`])},[b(`div`,{class:S([`w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0 transition-colors duration-150`,V(t)])},[g.value[t]?(o(),x(`i`,Qn)):(o(),x(`i`,{key:1,class:S(e.icon)},null,2))],2),b(`div`,$n,[b(`div`,{class:S([`text-sm font-medium truncate`,H(t)])},v(l(c)(e.labelKey)),3)]),g.value[t]?(o(),x(`i`,er)):h(``,!0)],42,Zn)),64))]),b(`div`,tr,[(o(),u(s(z.value),{key:p.value,"org-id":l(f),"org-name":C.value,"org-code":E.value,"org-grading-id":D.value,"org-job-family-id":O.value,"grading-options":j.value,"job-family-options":M.value,"job-value-options":F.value,"competency-options":L.value,"job-value-map":I.value,onSaved:W},null,40,[`org-id`,`org-name`,`org-code`,`org-grading-id`,`org-job-family-id`,`grading-options`,`job-family-options`,`job-value-options`,`competency-options`,`job-value-map`]))])]))]))}};export{nr as default};