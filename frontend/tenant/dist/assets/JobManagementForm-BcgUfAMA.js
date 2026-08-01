const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["assets/quill-BLmY9xB4.js","assets/rolldown-runtime-QTnfLwEv.js"])))=>i.map(i=>d[i]);
import{A as e,B as t,D as n,J as r,M as i,N as a,S as o,T as s,V as c,W as l,c as u,ct as d,dt as f,h as p,k as m,l as h,lt as g,m as _,o as v,r as y,s as b,u as x,ut as S,y as C}from"./runtime-core.esm-bundler-X_uJX_FV.js";import{a as w,d as T,t as E}from"./button-rY8vogVu.js";import{A as D,a as O}from"./ripple-BB-Blkgv.js";import{l as k,m as A,s as j,t as M,u as N}from"./index-BUZaJSAZ.js";import{t as P}from"./useI18n-BTKJxl68.js";import{n as F}from"./responseHandler-B5MnXl3B.js";import{t as I}from"./tag-BcP346rv.js";import{t as L}from"./FormRow-DknQZfFw.js";import{t as R}from"./baseeditableholder-COw1OOPE.js";import{t as z}from"./textarea-CM81eAS2.js";import{t as B}from"./TextInput-BSlvHRx7.js";import{i as V,n as H,o as U,t as W}from"./column-CkmiY7wH.js";import{t as G}from"./SelectLabel-BrW6kToF.js";import{t as K}from"./ConfirmDeleteDialog-BxMzmECz.js";import{t as ee}from"./SkeletonTable-5a0eXFby.js";import{t as te}from"./toggleswitch-rm0GD9m-.js";var ne={class:`space-y-4`},re={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},ie={class:`text-sm text-gray-500 dark:text-gray-400`},ae={class:`max-w-2xl space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},oe={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},se={class:`flex justify-end pt-2`},ce=`/api/v1/tenant/job-management/identifications`,le={__name:`JobIdentificationSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},orgGradingId:{type:String,default:``},orgJobFamilyId:{type:String,default:``},gradingOptions:{type:Array,default:()=>[]},jobFamilyOptions:{type:Array,default:()=>[]}},emits:[`saved`],setup(e,{emit:i}){let a=i,o=e,{t:c}=P(),u=A(),d=l(!1),m=l(``),g=l({}),_=l(``),y=l({grading_id:``}),S=v(()=>{let e=o.jobFamilyOptions.find(e=>e.value===o.orgJobFamilyId);return e?e.label:o.orgJobFamilyId||`-`});function C(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function T(){if(o.orgId)try{let e=(await w.get(ce,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];_.value=t.id,y.value.grading_id=t.grading_id||o.orgGradingId||``}else y.value.grading_id=o.orgGradingId||``}catch{y.value.grading_id=o.orgGradingId||``}}async function D(){if(m.value=``,g.value={},!y.value.grading_id){m.value=c(`job_management.grading_required`);return}d.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,grading_id:y.value.grading_id,organization_id:o.orgId};if(_.value)await w.put(`${ce}/${_.value}`,{grading_id:y.value.grading_id});else{let t=await w.post(ce,e);_.value=t.data?.data?.id||``}u.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=C(e);Object.keys(t).length>0?(g.value=t,m.value=Object.values(t).join(`, `)):m.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{d.value=!1}}return s(T),(i,a)=>(n(),x(`div`,ne,[b(`div`,null,[b(`h2`,re,f(r(c)(`job_management.identifications`)),1),b(`p`,ie,f(r(c)(`job_management.identification_description`)),1)]),b(`div`,ae,[p(L,{label:r(c)(`organization.nomenclature`)},{default:t(()=>[p(B,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:r(c)(`organization.full_code`)},{default:t(()=>[p(B,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:r(c)(`organization.job_family`)},{default:t(()=>[p(B,{"model-value":S.value,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:r(c)(`organization.grading`)},{default:t(()=>[p(r(U),{modelValue:y.value.grading_id,"onUpdate:modelValue":a[0]||=e=>y.value.grading_id=e,options:e.gradingOptions,"option-label":`label`,"option-value":`value`,placeholder:r(c)(`organization.select_grading`),class:`w-full`,size:`small`,invalid:!!g.value.grading_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`]),m.value?(n(),x(`div`,oe,f(m.value),1)):h(``,!0),b(`div`,se,[p(r(E),{label:r(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:d.value,disabled:!y.value.grading_id,onClick:D},null,8,[`label`,`loading`,`disabled`])])])]))}},ue={class:`space-y-4`},de={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},fe={class:`text-sm text-gray-500 dark:text-gray-400`},pe={class:`max-w-2xl space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},me={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},he={class:`flex justify-end gap-2 pt-2`},q=`/api/v1/tenant/job-management/objectives`,ge={__name:`JobObjectiveSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:i}){let a=i,o=e,{t:c}=P(),m=A(),g=l(!1),_=l(!1),v=l(``),y=l({}),S=l(``),C=l(!1),T=l(``),D=l({objective:``});function O(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function k(){if(o.orgId)try{let e=(await w.get(q,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];S.value=t.id,D.value.objective=t.objective||``}}catch{}}async function j(){v.value=``,y.value={},g.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,objective:D.value.objective||``,organization_id:o.orgId};if(S.value)await w.put(`${q}/${S.value}`,{objective:D.value.objective||``});else{let t=await w.post(q,e);S.value=t.data?.data?.id||``}m.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=O(e);Object.keys(t).length>0?(y.value=t,v.value=Object.values(t).join(`, `)):v.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{g.value=!1}}async function M(){if(S.value){_.value=!0,T.value=``;try{await w.delete(`${q}/${S.value}`),C.value=!1,S.value=``,D.value.objective=``,a(`saved`),m.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||c(`message.operation_failed`)}finally{_.value=!1}}}return s(k),(i,a)=>(n(),x(`div`,ue,[b(`div`,null,[b(`h2`,de,f(r(c)(`job_management.objectives`)),1),b(`p`,fe,f(r(c)(`job_management.objective_description`)),1)]),b(`div`,pe,[p(L,{label:r(c)(`organization.nomenclature`)},{default:t(()=>[p(B,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:r(c)(`organization.full_code`)},{default:t(()=>[p(B,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:r(c)(`job_management.objective`)},{default:t(()=>[p(r(z),{modelValue:D.value.objective,"onUpdate:modelValue":a[0]||=e=>D.value.objective=e,rows:`3`,class:d([`w-full`,{"p-invalid":y.value.objective}]),placeholder:r(c)(`job_management.objective`)+`...`},null,8,[`modelValue`,`class`,`placeholder`])]),_:1},8,[`label`]),v.value?(n(),x(`div`,me,f(v.value),1)):h(``,!0),b(`div`,he,[S.value?(n(),u(r(E),{key:0,label:r(c)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:a[1]||=e=>C.value=!0},null,8,[`label`])):h(``,!0),p(r(E),{label:S.value?r(c)(`common.update`):r(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:g.value,disabled:g.value,onClick:j},null,8,[`label`,`loading`,`disabled`])])]),p(K,{visible:C.value,"onUpdate:visible":a[2]||=e=>C.value=e,loading:_.value,"error-msg":T.value,onConfirm:M,onCancel:a[3]||=e=>C.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},_e={class:`space-y-4`},ve={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},ye={class:`text-sm text-gray-500 dark:text-gray-400`},be={class:`max-w-2xl space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},xe={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Se={class:`flex justify-end gap-2 pt-2`},J=`/api/v1/tenant/job-management/education-experiences`,Ce={__name:`JobEduExpSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:i}){let a=i,o=e,{t:c}=P(),d=A(),m=l(!1),g=l(!1),_=l(``),y=l({}),S=l(``),C=l(!1),T=l(``),D=l({job_management_value_education_id:``,job_management_value_experience_id:``}),O=v(()=>o.jobValueMap?.education||[]),k=v(()=>o.jobValueMap?.experience||[]);async function j(){if(o.orgId)try{let e=(await w.get(J,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];S.value=t.id,D.value.job_management_value_education_id=t.job_management_value_education_id||``,D.value.job_management_value_experience_id=t.job_management_value_experience_id||``}}catch{}}async function M(){_.value=``,y.value={},m.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,job_management_value_education_id:D.value.job_management_value_education_id||null,job_management_value_experience_id:D.value.job_management_value_experience_id||null,organization_id:o.orgId};if(S.value)await w.put(`${J}/${S.value}`,{job_management_value_education_id:D.value.job_management_value_education_id||null,job_management_value_experience_id:D.value.job_management_value_experience_id||null});else{let t=await w.post(J,e);S.value=t.data?.data?.id||``}d.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=F(e);Object.keys(t).length>0?(y.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{m.value=!1}}async function N(){if(S.value){g.value=!0,T.value=``;try{await w.delete(`${J}/${S.value}`),C.value=!1,S.value=``,D.value.job_management_value_education_id=``,D.value.job_management_value_experience_id=``,a(`saved`),d.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||c(`message.operation_failed`)}finally{g.value=!1}}}return s(j),(i,a)=>(n(),x(`div`,_e,[b(`div`,null,[b(`h2`,ve,f(r(c)(`job_management.education_experience`)),1),b(`p`,ye,f(r(c)(`job_management.education_experience_description`)),1)]),b(`div`,be,[p(L,{label:r(c)(`organization.nomenclature`)},{default:t(()=>[p(B,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:r(c)(`organization.full_code`)},{default:t(()=>[p(B,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),p(L,{label:r(c)(`job_management.education_level`),errors:y.value?.job_management_value_education_id},{default:t(()=>[p(r(U),{modelValue:D.value.job_management_value_education_id,"onUpdate:modelValue":a[0]||=e=>D.value.job_management_value_education_id=e,options:O.value,"option-label":`label`,"option-value":`value`,placeholder:r(c)(`common.select`),class:`w-full`,size:`small`,showClear:``,invalid:!!y.value.job_management_value_education_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(c)(`job_management.experience_level`),errors:y.value?.job_management_value_experience_id},{default:t(()=>[p(r(U),{modelValue:D.value.job_management_value_experience_id,"onUpdate:modelValue":a[1]||=e=>D.value.job_management_value_experience_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:r(c)(`common.select`),class:`w-full`,size:`small`,showClear:``,invalid:!!y.value.job_management_value_experience_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`]),_.value?(n(),x(`div`,xe,f(_.value),1)):h(``,!0),b(`div`,Se,[S.value?(n(),u(r(E),{key:0,label:r(c)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:a[2]||=e=>C.value=!0},null,8,[`label`])):h(``,!0),p(r(E),{label:S.value?r(c)(`common.update`):r(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:m.value,disabled:m.value,onClick:M},null,8,[`label`,`loading`,`disabled`])])]),p(K,{visible:C.value,"onUpdate:visible":a[3]||=e=>C.value=e,loading:g.value,"error-msg":T.value,onConfirm:N,onCancel:a[4]||=e=>C.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},we=O.extend({name:`editor`,style:`
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
`,classes:{root:function(e){return[`p-editor`,{"p-invalid":e.instance.$invalid}]},toolbar:`p-editor-toolbar`,content:`p-editor-content`}}),Te={name:`BaseEditor`,extends:R,props:{placeholder:String,readonly:Boolean,formats:Array,editorStyle:null,modules:null},style:we,provide:function(){return{$pcEditor:this,$parentInstance:this}}};function Y(e){"@babel/helpers - typeof";return Y=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},Y(e)}function Ee(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function De(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?Ee(Object(n),!0).forEach(function(t){Oe(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):Ee(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function Oe(e,t,n){return(t=ke(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function ke(e){var t=Ae(e,`string`);return Y(t)==`symbol`?t:t+``}function Ae(e,t){if(Y(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(Y(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var je=function(){try{return window.Quill}catch{return null}}(),X={name:`Editor`,extends:Te,inheritAttrs:!1,emits:[`text-change`,`selection-change`,`load`],quill:null,watch:{modelValue:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},d_value:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},readonly:function(){this.handleReadOnlyChange()}},mounted:function(){var e=this,t={modules:De({toolbar:this.$refs.toolbarElement},this.modules),readOnly:this.readonly,theme:`snow`,formats:this.formats,placeholder:this.placeholder};je?(this.quill=new je(this.$refs.editorElement,t),this.initQuill(),this.handleLoad()):j(()=>import(`./quill-BLmY9xB4.js`).then(function(n){n&&D(e.$refs.editorElement)&&(n.default?e.quill=new n.default(e.$refs.editorElement,t):e.quill=new n(e.$refs.editorElement,t),e.initQuill())}),__vite__mapDeps([0,1])).then(function(){e.handleLoad()})},beforeUnmount:function(){this.quill=null},methods:{renderValue:function(e){if(this.quill)if(e){var t=this.quill.clipboard.convert({html:e});this.quill.setContents(t)}else this.quill.setText(``)},initQuill:function(){var e=this;this.renderValue(this.d_value),this.quill.on(`text-change`,function(t,n,r){if(r===`user`){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();i===`<p><br></p>`&&(i=``),e.writeValue(i),e.$emit(`text-change`,{htmlValue:i,textValue:a,delta:t,source:r,instance:e.quill})}}),this.quill.on(`selection-change`,function(t,n,r){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();e.$emit(`selection-change`,{htmlValue:i,textValue:a,range:t,oldRange:n,source:r,instance:e.quill})})},handleLoad:function(){this.quill&&this.quill.getModule(`toolbar`)&&this.$emit(`load`,{instance:this.quill})},handleReadOnlyChange:function(){this.quill&&this.quill.enable(!this.readonly)}}};function Me(t,r,i,a,s,c){return n(),x(`div`,o({class:t.cx(`root`)},t.ptmi(`root`)),[b(`div`,o({ref:`toolbarElement`,class:t.cx(`toolbar`)},t.ptm(`toolbar`)),[e(t.$slots,`toolbar`,{},function(){return[b(`span`,o({class:`ql-formats`},t.ptm(`formats`)),[b(`select`,o({class:`ql-header`,defaultValue:`0`},t.ptm(`header`)),[b(`option`,o({value:`1`},t.ptm(`option`)),`Heading`,16),b(`option`,o({value:`2`},t.ptm(`option`)),`Subheading`,16),b(`option`,o({value:`0`},t.ptm(`option`)),`Normal`,16)],16),b(`select`,o({class:`ql-font`},t.ptm(`font`)),[b(`option`,g(C(t.ptm(`option`))),null,16),b(`option`,o({value:`serif`},t.ptm(`option`)),null,16),b(`option`,o({value:`monospace`},t.ptm(`option`)),null,16)],16)],16),b(`span`,o({class:`ql-formats`},t.ptm(`formats`)),[b(`button`,o({class:`ql-bold`,type:`button`},t.ptm(`bold`)),null,16),b(`button`,o({class:`ql-italic`,type:`button`},t.ptm(`italic`)),null,16),b(`button`,o({class:`ql-underline`,type:`button`},t.ptm(`underline`)),null,16)],16),b(`span`,o({class:`ql-formats`},t.ptm(`formats`)),[b(`select`,o({class:`ql-color`},t.ptm(`color`)),null,16),b(`select`,o({class:`ql-background`},t.ptm(`background`)),null,16)],16),b(`span`,o({class:`ql-formats`},t.ptm(`formats`)),[b(`button`,o({class:`ql-list`,value:`ordered`,type:`button`},t.ptm(`list`)),null,16),b(`button`,o({class:`ql-list`,value:`bullet`,type:`button`},t.ptm(`list`)),null,16),b(`select`,o({class:`ql-align`},t.ptm(`select`)),[b(`option`,o({defaultValue:``},t.ptm(`option`)),null,16),b(`option`,o({value:`center`},t.ptm(`option`)),null,16),b(`option`,o({value:`right`},t.ptm(`option`)),null,16),b(`option`,o({value:`justify`},t.ptm(`option`)),null,16)],16)],16),b(`span`,o({class:`ql-formats`},t.ptm(`formats`)),[b(`button`,o({class:`ql-link`,type:`button`},t.ptm(`link`)),null,16),b(`button`,o({class:`ql-image`,type:`button`},t.ptm(`image`)),null,16),b(`button`,o({class:`ql-code-block`,type:`button`},t.ptm(`codeBlock`)),null,16)],16),b(`span`,o({class:`ql-formats`},t.ptm(`formats`)),[b(`button`,o({class:`ql-clean`,type:`button`},t.ptm(`clean`)),null,16)],16)]})],16),b(`div`,o({ref:`editorElement`,class:t.cx(`content`),style:t.editorStyle},t.ptm(`content`)),null,16)],16)}X.render=Me;var Ne={key:0,class:`text-gray-500 dark:text-gray-400 text-xs`},Pe=[`innerHTML`],Fe={key:2,class:`text-gray-800 dark:text-gray-100`},Ie={class:`flex items-center gap-1`},Z={__name:`DataTableSection`,props:{items:Array,loading:Boolean,total:Number,columns:{type:Array,default:()=>[]},entity:String,orgId:String,onLoad:Function},emits:[`edit`,`delete`],setup(a){let o=a,{t:d}=P(),g=l(1),_=l(15),S=v(()=>(g.value-1)*_.value),C=v(()=>[...o.columns.map(e=>({type:`text`,width:`w-24`,headerWidth:`w-20`})),{type:`icons`,count:2,headerWidth:`w-16`}]);function w(e){g.value=e.page+1,_.value=e.rows,o.onLoad&&o.onLoad(g.value,_.value)}return s(()=>{o.onLoad&&o.onLoad(1,15)}),(o,s)=>{let l=i(`tooltip`);return a.loading?(n(),u(ee,{key:0,columns:C.value,rows:8},null,8,[`columns`])):(n(),u(r(H),{key:1,value:a.items,lazy:``,totalRecords:a.total,first:S.value,rows:_.value,onPage:w,paginator:``,paginatorTemplate:`FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown`,rowsPerPageOptions:[10,15,25,50],size:`small`,class:`!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden`},{empty:t(()=>[e(o.$slots,`empty`)]),default:t(()=>[(n(!0),x(y,null,m(a.columns,e=>(n(),u(r(W),{key:e.field,field:e.field,header:e.header,sortable:``},{body:t(({data:t})=>[e.field.startsWith(`_`)?(n(),x(`span`,Ne,f(t[e.field]||`-`),1)):h(``,!0),e.html?(n(),x(`div`,{key:1,class:`editor-content`,innerHTML:t[e.field]},null,8,Pe)):(n(),x(`span`,Fe,f(t[e.field]||`-`),1))]),_:2},1032,[`field`,`header`]))),128)),p(r(W),{header:r(d)(`common.actions`),style:{width:`90px`},frozen:``,alignFrozen:`right`},{body:t(({data:e})=>[b(`div`,Ie,[c(p(r(E),{icon:`pi pi-pencil`,size:`small`,text:``,severity:`secondary`,onClick:t=>o.$emit(`edit`,e)},null,8,[`onClick`]),[[l,r(d)(`common.edit`),void 0,{left:!0}]]),c(p(r(E),{icon:`pi pi-trash`,size:`small`,text:``,severity:`danger`,onClick:t=>o.$emit(`delete`,e)},null,8,[`onClick`]),[[l,r(d)(`common.delete`),void 0,{left:!0}]])])]),_:1},8,[`header`])]),_:3},8,[`value`,`totalRecords`,`first`,`rows`]))}}},Le={class:`space-y-4`},Re={key:0,class:`bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded p-3 text-xs text-red-700 dark:text-red-300`},Q={__name:`DialogForm`,props:{visible:Boolean,title:String,saving:Boolean,errors:{type:Object,default:()=>({})},width:{type:String,default:`480px`}},emits:[`save`,`cancel`],setup(i){let a=i,{t:o}=P(),s=v(()=>a.width===`maximize`?`90vw`:a.width);return(a,c)=>(n(),u(r(M),{visible:i.visible,"onUpdate:visible":c[2]||=e=>a.$emit(`update:visible`,e),header:i.title,modal:``,style:S({width:s.value}),class:`p-fluid`,closable:!i.saving},{footer:t(()=>[p(r(E),{label:r(o)(`common.cancel`),size:`small`,outlined:``,severity:`secondary`,disabled:i.saving,onClick:c[0]||=e=>a.$emit(`cancel`)},null,8,[`label`,`disabled`]),p(r(E),{label:r(o)(`common.save`),icon:`pi pi-check`,size:`small`,loading:i.saving,onClick:c[1]||=e=>a.$emit(`save`)},null,8,[`label`,`loading`])]),default:t(()=>[b(`div`,Le,[e(a.$slots,`default`),Object.keys(i.errors).length?(n(),x(`div`,Re,[(n(!0),x(y,null,m(i.errors,(e,t)=>(n(),x(`p`,{key:t,class:`mb-1`},[b(`strong`,null,f(t)+`:`,1),_(` `+f(Array.isArray(e)?e.join(`, `):e),1)]))),128))])):h(``,!0)])]),_:3},8,[`visible`,`header`,`style`,`closable`]))}},ze={class:`space-y-4`},Be={class:`flex items-center justify-between`},Ve={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},He={class:`text-sm text-gray-500 dark:text-gray-400`},Ue={class:`flex flex-col items-center justify-center py-10 text-gray-400`},We={class:`text-sm font-medium`},Ge=`/api/v1/tenant/job-management/responsibilities`,Ke={__name:`JobResponsibilitySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:i}){let a=e,o=i,{t:s}=P(),c=A(),u=l([]),m=l(!1),g=l(0),_=l(!1),S=l(!1),C=l(``),T=l(!1),D=l({}),O=l(!1),k=l(!1),j=l(``),M=l(null),N=l({main_task:``,activities:``,outputs:``,success_indicators:``}),I=v(()=>{let e=s(`job_management.responsibilities_title`);return S.value?`${e}`:`${s(`common.create`)} ${e}`}),R=v(()=>[{field:`main_task`,header:s(`job_management.main_task`),html:!0},{field:`activities`,header:s(`job_management.activities`),html:!0},{field:`outputs`,header:s(`job_management.outputs`),html:!0},{field:`success_indicators`,header:s(`job_management.success_indicators`),html:!0}]);async function z(e,t){m.value=!0;try{let n=await w.get(Ge,{params:{page:e,per_page:t,organization_id:a.orgId}}),r=n.data?.data||[];u.value=r.map(e=>({...e,main_task:e.main_task,activities:e.activities,outputs:e.outputs,success_indicators:e.success_indicators})),g.value=n.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{m.value=!1}}function B(){S.value=!1,C.value=``,N.value={main_task:``,activities:``,outputs:``,success_indicators:``},D.value={},_.value=!0}function V(e){S.value=!0,C.value=e.id,N.value={main_task:e.main_task||``,activities:e.activities||``,outputs:e.outputs||``,success_indicators:e.success_indicators||``},D.value={},_.value=!0}async function H(){T.value=!0,D.value={};try{let e={nomenclature:a.orgName||``,full_code:a.orgCode||``,...N.value,organization_id:a.orgId};S.value?await w.put(`${Ge}/${C.value}`,e):await w.post(Ge,e),_.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),z(1,15)}catch(e){let t=F(e);Object.keys(t).length?D.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{T.value=!1}}function U(e){M.value=e,j.value=``,O.value=!0}async function W(){if(M.value){k.value=!0,j.value=``;try{await w.delete(`${Ge}/${M.value.id}`),O.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),z(1,15)}catch(e){j.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{k.value=!1}}}return(i,a)=>(n(),x(`div`,ze,[b(`div`,Be,[b(`div`,null,[b(`h2`,Ve,f(r(s)(`job_management.responsibilities_title`)),1),b(`p`,He,f(r(s)(`job_management.responsibilities_description`)),1)]),p(r(E),{label:r(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:a[0]||=e=>B()},null,8,[`label`])]),p(Z,{items:u.value,loading:m.value,total:g.value,columns:R.value,entity:`responsibilities`,"org-id":e.orgId,"on-load":z,onEdit:V,onDelete:U},{empty:t(()=>[b(`div`,Ue,[a[9]||=b(`i`,{class:`pi pi-list-check text-3xl mb-2 opacity-50`},null,-1),b(`p`,We,f(r(s)(`job_management.empty_responsibilities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:_.value,"onUpdate:visible":a[5]||=e=>_.value=e,title:I.value,saving:T.value,errors:D.value,width:`maximize`,onSave:H,onCancel:a[6]||=e=>_.value=!1},{default:t(()=>[_.value?(n(),x(y,{key:0},[p(L,{label:r(s)(`job_management.main_task`),errors:D.value?.main_task},{default:t(()=>[p(r(X),{modelValue:N.value.main_task,"onUpdate:modelValue":a[1]||=e=>N.value.main_task=e,editorStyle:`height:120px`,class:d({"p-invalid":D.value?.main_task})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.activities`),errors:D.value?.activities},{default:t(()=>[p(r(X),{modelValue:N.value.activities,"onUpdate:modelValue":a[2]||=e=>N.value.activities=e,editorStyle:`height:120px`,class:d({"p-invalid":D.value?.activities})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.outputs`),errors:D.value?.outputs},{default:t(()=>[p(r(X),{modelValue:N.value.outputs,"onUpdate:modelValue":a[3]||=e=>N.value.outputs=e,editorStyle:`height:120px`,class:d({"p-invalid":D.value?.outputs})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.success_indicators`),errors:D.value?.success_indicators},{default:t(()=>[p(r(X),{modelValue:N.value.success_indicators,"onUpdate:modelValue":a[4]||=e=>N.value.success_indicators=e,editorStyle:`height:120px`,class:d({"p-invalid":D.value?.success_indicators})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])],64)):h(``,!0)]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:O.value,"onUpdate:visible":a[7]||=e=>O.value=e,loading:k.value,"error-msg":j.value,onConfirm:W,onCancel:a[8]||=e=>O.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},qe={class:`space-y-4`},Je={class:`flex items-center justify-between`},Ye={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Xe={class:`text-sm text-gray-500 dark:text-gray-400`},Ze={class:`flex flex-col items-center justify-center py-10 text-gray-400`},Qe={class:`text-sm font-medium`},$e=`/api/v1/tenant/job-management/hr-authorities`,et={__name:`JobHRAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:i}){let a=e,o=i,{t:s}=P(),c=A(),u=l([]),m=l(!1),h=l(0),g=l(!1),_=l(!1),y=l(``),S=l(!1),C=l({}),T=l(!1),D=l(!1),O=l(``),k=l(null),j=l({description:``}),M=v(()=>{let e=s(`job_management.hr_authorities`);return _.value?`${s(`common.edit`)} ${e}`:`${s(`common.create`)} ${e}`}),N=v(()=>[{field:`description`,header:s(`job_management.description`)}]);async function I(e,t){m.value=!0;try{let n=await w.get($e,{params:{page:e,per_page:t,organization_id:a.orgId}});u.value=n.data?.data||[],h.value=n.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{m.value=!1}}function R(){_.value=!1,y.value=``,j.value={nomenclature:``,full_code:``,description:``},C.value={},g.value=!0}function B(e){_.value=!0,y.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},C.value={},g.value=!0}async function V(){S.value=!0,C.value={};try{let e={...j.value,nomenclature:a.orgName||``,full_code:a.orgCode||``,organization_id:a.orgId};_.value?await w.put(`${$e}/${y.value}`,e):await w.post($e,e),g.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),I(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{S.value=!1}}function H(e){k.value=e,O.value=``,T.value=!0}async function U(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${$e}/${k.value.id}`),T.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),I(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(i,a)=>(n(),x(`div`,qe,[b(`div`,Je,[b(`div`,null,[b(`h2`,Ye,f(r(s)(`job_management.hr_authorities`)),1),b(`p`,Xe,f(r(s)(`job_management.authority_description`)),1)]),p(r(E),{label:r(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:a[0]||=e=>R()},null,8,[`label`])]),p(Z,{items:u.value,loading:m.value,total:h.value,columns:N.value,entity:`hr-authorities`,"org-id":e.orgId,"on-load":I,onEdit:B,onDelete:H},{empty:t(()=>[b(`div`,Ze,[a[6]||=b(`i`,{class:`pi pi-users text-3xl mb-2 opacity-50`},null,-1),b(`p`,Qe,f(r(s)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:g.value,"onUpdate:visible":a[2]||=e=>g.value=e,title:M.value,saving:S.value,errors:C.value,onSave:V,onCancel:a[3]||=e=>g.value=!1},{default:t(()=>[p(L,{label:r(s)(`job_management.description`),errors:C.value?.description},{default:t(()=>[p(r(z),{modelValue:j.value.description,"onUpdate:modelValue":a[1]||=e=>j.value.description=e,rows:`3`,class:d([`w-full`,{"p-invalid":C.value?.description}])},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":a[4]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:U,onCancel:a[5]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},tt={class:`space-y-4`},nt={class:`flex items-center justify-between`},rt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},it={class:`text-sm text-gray-500 dark:text-gray-400`},at={class:`flex flex-col items-center justify-center py-10 text-gray-400`},ot={class:`text-sm font-medium`},st=`/api/v1/tenant/job-management/operational-authorities`,ct={__name:`JobOpAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:i}){let a=e,o=i,{t:s}=P(),c=A(),u=l([]),m=l(!1),h=l(0),g=l(!1),_=l(!1),y=l(``),S=l(!1),C=l({}),T=l(!1),D=l(!1),O=l(``),k=l(null),j=l({description:``}),M=v(()=>{let e=s(`job_management.op_authorities`);return _.value?`${s(`common.edit`)} ${e}`:`${s(`common.create`)} ${e}`}),N=v(()=>[{field:`description`,header:s(`job_management.description`)}]);async function I(e,t){m.value=!0;try{let n=await w.get(st,{params:{page:e,per_page:t,organization_id:a.orgId}});u.value=n.data?.data||[],h.value=n.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{m.value=!1}}function R(){_.value=!1,y.value=``,j.value={nomenclature:``,full_code:``,description:``},C.value={},g.value=!0}function B(e){_.value=!0,y.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},C.value={},g.value=!0}async function V(){S.value=!0,C.value={};try{let e={...j.value,nomenclature:a.orgName||``,full_code:a.orgCode||``,organization_id:a.orgId};_.value?await w.put(`${st}/${y.value}`,e):await w.post(st,e),g.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),I(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{S.value=!1}}function H(e){k.value=e,O.value=``,T.value=!0}async function U(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${st}/${k.value.id}`),T.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),I(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(i,a)=>(n(),x(`div`,tt,[b(`div`,nt,[b(`div`,null,[b(`h2`,rt,f(r(s)(`job_management.op_authorities`)),1),b(`p`,it,f(r(s)(`job_management.authority_description`)),1)]),p(r(E),{label:r(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:a[0]||=e=>R()},null,8,[`label`])]),p(Z,{items:u.value,loading:m.value,total:h.value,columns:N.value,entity:`operational-authorities`,"org-id":e.orgId,"on-load":I,onEdit:B,onDelete:H},{empty:t(()=>[b(`div`,at,[a[6]||=b(`i`,{class:`pi pi-cog text-3xl mb-2 opacity-50`},null,-1),b(`p`,ot,f(r(s)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:g.value,"onUpdate:visible":a[2]||=e=>g.value=e,title:M.value,saving:S.value,errors:C.value,onSave:V,onCancel:a[3]||=e=>g.value=!1},{default:t(()=>[p(L,{label:r(s)(`job_management.description`),errors:C.value?.description},{default:t(()=>[p(r(z),{modelValue:j.value.description,"onUpdate:modelValue":a[1]||=e=>j.value.description=e,class:d([`w-full`,{"p-invalid":C.value?.description}]),rows:`3`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":a[4]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:U,onCancel:a[5]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},lt={class:`space-y-4`},ut={class:`flex items-center justify-between`},dt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},ft={class:`text-sm text-gray-500 dark:text-gray-400`},pt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},mt={class:`text-sm font-medium`},ht=`/api/v1/tenant/job-management/working-activities`,gt={__name:`JobActivitySection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:i}){let a=e,o=i,{t:s}=P(),c=A(),u=l([]),m=l(!1),h=l(0),g=l(!1),_=l(!1),y=l(``),S=l(!1),C=l({}),T=l(!1),D=l(!1),O=l(``),k=l(null),j=l({nomenclature:``,full_code:``,job_management_value_id:``});v(()=>Object.values(a.jobValueMap||{}).flat());let M=v(()=>[{field:`nomenclature`,header:s(`organization.nomenclature`)},{field:`full_code`,header:s(`organization.full_code`)}]);async function N(e,t){m.value=!0;try{let n=await w.get(ht,{params:{page:e,per_page:t,organization_id:a.orgId}});u.value=n.data?.data||[],h.value=n.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{m.value=!1}}function I(){_.value=!1,y.value=``,j.value={nomenclature:``,full_code:``,job_management_value_id:``},C.value={},g.value=!0}function R(e){_.value=!0,y.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_id:e.job_management_value_id||``},C.value={},g.value=!0}async function z(){S.value=!0,C.value={};try{let e={...j.value,organization_id:a.orgId};_.value?await w.put(`${ht}/${y.value}`,e):await w.post(ht,e),g.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),N(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{S.value=!1}}function V(e){k.value=e,O.value=``,T.value=!0}async function H(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${ht}/${k.value.id}`),T.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),N(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(i,a)=>(n(),x(`div`,lt,[b(`div`,ut,[b(`div`,null,[b(`h2`,dt,f(r(s)(`job_management.activities`)),1),b(`p`,ft,f(r(s)(`job_management.activity_description`)),1)]),p(r(E),{label:r(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:a[0]||=e=>I()},null,8,[`label`])]),p(Z,{items:u.value,loading:m.value,total:h.value,columns:M.value,entity:`working-activities`,"org-id":e.orgId,"on-load":N,onEdit:R,onDelete:V},{empty:t(()=>[b(`div`,pt,[a[8]||=b(`i`,{class:`pi pi-bolt text-3xl mb-2 opacity-50`},null,-1),b(`p`,mt,f(r(s)(`job_management.empty_activities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:g.value,"onUpdate:visible":a[4]||=e=>g.value=e,title:_.value?r(s)(`common.edit`):r(s)(`common.create`),saving:S.value,errors:C.value,onSave:z,onCancel:a[5]||=e=>g.value=!1},{default:t(()=>[p(L,{label:r(s)(`organization.nomenclature`),required:``,errors:C.value?.nomenclature},{default:t(()=>[p(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":a[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:d({"p-invalid":C.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`organization.full_code`),required:``,errors:C.value?.full_code},{default:t(()=>[p(B,{modelValue:j.value.full_code,"onUpdate:modelValue":a[2]||=e=>j.value.full_code=e,maxlength:`20`,class:d({"p-invalid":C.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.activity_type`),errors:C.value?.job_management_value_id},{default:t(()=>[p(G,{modelValue:j.value.job_management_value_id,"onUpdate:modelValue":a[3]||=e=>j.value.job_management_value_id=e,options:i.activityOptions,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":a[6]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:H,onCancel:a[7]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},_t={class:`space-y-4`},vt={class:`flex items-center justify-between`},yt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},bt={class:`text-sm text-gray-500 dark:text-gray-400`},xt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},St={class:`text-sm font-medium`},$=`/api/v1/tenant/job-management/working-risks`,Ct={__name:`JobRiskSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:i}){let a=e,o=i,{t:s}=P(),c=A(),u=l([]),m=l(!1),h=l(0),g=l(!1),_=l(!1),y=l(``),S=l(!1),C=l({}),T=l(!1),D=l(!1),O=l(``),k=l(null),j=l({nomenclature:``,full_code:``,job_management_value_environment_id:``,job_management_value_hazard_id:``}),M=v(()=>a.jobValueMap?.environment||[]),N=v(()=>a.jobValueMap?.hazard||[]),I=v(()=>[{field:`nomenclature`,header:s(`organization.nomenclature`)},{field:`full_code`,header:s(`organization.full_code`)}]);async function R(e,t){m.value=!0;try{let n=await w.get($,{params:{page:e,per_page:t,organization_id:a.orgId}});u.value=n.data?.data||[],h.value=n.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{m.value=!1}}function z(){_.value=!1,y.value=``,j.value={nomenclature:``,full_code:``,job_management_value_environment_id:``,job_management_value_hazard_id:``},C.value={},g.value=!0}function V(e){_.value=!0,y.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_environment_id:e.job_management_value_environment_id||``,job_management_value_hazard_id:e.job_management_value_hazard_id||``},C.value={},g.value=!0}async function H(){S.value=!0,C.value={};try{let e={...j.value,organization_id:a.orgId};_.value?await w.put(`${$}/${y.value}`,e):await w.post($,e),g.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),R(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{S.value=!1}}function U(e){k.value=e,O.value=``,T.value=!0}async function W(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${$}/${k.value.id}`),T.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),R(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(i,a)=>(n(),x(`div`,_t,[b(`div`,vt,[b(`div`,null,[b(`h2`,yt,f(r(s)(`job_management.risks`)),1),b(`p`,bt,f(r(s)(`job_management.risk_description`)),1)]),p(r(E),{label:r(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:a[0]||=e=>z()},null,8,[`label`])]),p(Z,{items:u.value,loading:m.value,total:h.value,columns:I.value,entity:`working-risks`,"org-id":e.orgId,"on-load":R,onEdit:V,onDelete:U},{empty:t(()=>[b(`div`,xt,[a[9]||=b(`i`,{class:`pi pi-exclamation-triangle text-3xl mb-2 opacity-50`},null,-1),b(`p`,St,f(r(s)(`job_management.empty_risks`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:g.value,"onUpdate:visible":a[5]||=e=>g.value=e,title:_.value?r(s)(`common.edit`):r(s)(`common.create`),saving:S.value,errors:C.value,onSave:H,onCancel:a[6]||=e=>g.value=!1},{default:t(()=>[p(L,{label:r(s)(`organization.nomenclature`),required:``,errors:C.value?.nomenclature},{default:t(()=>[p(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":a[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:d({"p-invalid":C.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`organization.full_code`),required:``,errors:C.value?.full_code},{default:t(()=>[p(B,{modelValue:j.value.full_code,"onUpdate:modelValue":a[2]||=e=>j.value.full_code=e,maxlength:`20`,class:d({"p-invalid":C.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.environment_risk`),errors:C.value?.job_management_value_environment_id},{default:t(()=>[p(G,{modelValue:j.value.job_management_value_environment_id,"onUpdate:modelValue":a[3]||=e=>j.value.job_management_value_environment_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.hazard_risk`),errors:C.value?.job_management_value_hazard_id},{default:t(()=>[p(G,{modelValue:j.value.job_management_value_hazard_id,"onUpdate:modelValue":a[4]||=e=>j.value.job_management_value_hazard_id=e,options:N.value,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":a[7]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:W,onCancel:a[8]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},wt={class:`space-y-4`},Tt={class:`flex items-center justify-between`},Et={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Dt={class:`text-sm text-gray-500 dark:text-gray-400`},Ot={class:`flex flex-col items-center justify-center py-10 text-gray-400`},kt={class:`text-sm font-medium`},At=`/api/v1/tenant/job-management/relationships`,jt={__name:`JobRelationshipSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:i}){let a=e,o=i,{t:s}=P(),c=A(),u=l([]),m=l(!1),h=l(0),g=l(!1),_=l(!1),y=l(``),S=l(!1),C=l({}),T=l(!1),D=l(!1),O=l(``),k=l(null),j=l({nomenclature:``,full_code:``,job_management_value_relationship_id:``,job_management_value_frequency_id:``}),M=v(()=>a.jobValueMap?.relationship||[]),N=v(()=>a.jobValueMap?.frequency||[]),I=v(()=>[{field:`nomenclature`,header:s(`organization.nomenclature`)},{field:`full_code`,header:s(`organization.full_code`)}]);async function R(e,t){m.value=!0;try{let n=await w.get(At,{params:{page:e,per_page:t,organization_id:a.orgId}});u.value=n.data?.data||[],h.value=n.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{m.value=!1}}function z(){_.value=!1,y.value=``,j.value={nomenclature:``,full_code:``,job_management_value_relationship_id:``,job_management_value_frequency_id:``},C.value={},g.value=!0}function V(e){_.value=!0,y.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_relationship_id:e.job_management_value_relationship_id||``,job_management_value_frequency_id:e.job_management_value_frequency_id||``},C.value={},g.value=!0}async function H(){S.value=!0,C.value={};try{let e={...j.value,organization_id:a.orgId};_.value?await w.put(`${At}/${y.value}`,e):await w.post(At,e),g.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),R(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{S.value=!1}}function U(e){k.value=e,O.value=``,T.value=!0}async function W(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${At}/${k.value.id}`),T.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),R(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(i,a)=>(n(),x(`div`,wt,[b(`div`,Tt,[b(`div`,null,[b(`h2`,Et,f(r(s)(`job_management.relationships`)),1),b(`p`,Dt,f(r(s)(`job_management.relationship_description`)),1)]),p(r(E),{label:r(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:a[0]||=e=>z()},null,8,[`label`])]),p(Z,{items:u.value,loading:m.value,total:h.value,columns:I.value,entity:`relationships`,"org-id":e.orgId,"on-load":R,onEdit:V,onDelete:U},{empty:t(()=>[b(`div`,Ot,[a[9]||=b(`i`,{class:`pi pi-share-alt text-3xl mb-2 opacity-50`},null,-1),b(`p`,kt,f(r(s)(`job_management.empty_relationships`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:g.value,"onUpdate:visible":a[5]||=e=>g.value=e,title:_.value?r(s)(`common.edit`):r(s)(`common.create`),saving:S.value,errors:C.value,onSave:H,onCancel:a[6]||=e=>g.value=!1},{default:t(()=>[p(L,{label:r(s)(`organization.nomenclature`),required:``,errors:C.value?.nomenclature},{default:t(()=>[p(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":a[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:d({"p-invalid":C.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`organization.full_code`),required:``,errors:C.value?.full_code},{default:t(()=>[p(B,{modelValue:j.value.full_code,"onUpdate:modelValue":a[2]||=e=>j.value.full_code=e,maxlength:`20`,class:d({"p-invalid":C.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.relationship_type`),errors:C.value?.job_management_value_relationship_id},{default:t(()=>[p(G,{modelValue:j.value.job_management_value_relationship_id,"onUpdate:modelValue":a[3]||=e=>j.value.job_management_value_relationship_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.frequency`),errors:C.value?.job_management_value_frequency_id},{default:t(()=>[p(G,{modelValue:j.value.job_management_value_frequency_id,"onUpdate:modelValue":a[4]||=e=>j.value.job_management_value_frequency_id=e,options:N.value,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":a[7]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:W,onCancel:a[8]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Mt={class:`space-y-4`},Nt={class:`flex items-center justify-between`},Pt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Ft={class:`text-sm text-gray-500 dark:text-gray-400`},It={class:`flex flex-col items-center justify-center py-10 text-gray-400`},Lt={class:`text-sm font-medium`},Rt=`/api/v1/tenant/job-management/subordinate-controls`,zt={__name:`JobSubordinateSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:i}){let a=e,o=i,{t:s}=P(),c=A(),u=l([]),m=l(!1),h=l(0),g=l(!1),_=l(!1),y=l(``),S=l(!1),C=l({}),T=l(!1),D=l(!1),O=l(``),k=l(null),j=l({nomenclature:``,full_code:``,job_management_value_id:``}),M=v(()=>Object.values(a.jobValueMap||{}).flat()),N=v(()=>[{field:`nomenclature`,header:s(`organization.nomenclature`)},{field:`full_code`,header:s(`organization.full_code`)}]);async function I(e,t){m.value=!0;try{let n=await w.get(Rt,{params:{page:e,per_page:t,organization_id:a.orgId}});u.value=n.data?.data||[],h.value=n.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{m.value=!1}}function R(){_.value=!1,y.value=``,j.value={nomenclature:``,full_code:``,job_management_value_id:``},C.value={},g.value=!0}function z(e){_.value=!0,y.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_id:e.job_management_value_id||``},C.value={},g.value=!0}async function V(){S.value=!0,C.value={};try{let e={...j.value,organization_id:a.orgId};_.value?await w.put(`${Rt}/${y.value}`,e):await w.post(Rt,e),g.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),I(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{S.value=!1}}function H(e){k.value=e,O.value=``,T.value=!0}async function U(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${Rt}/${k.value.id}`),T.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),I(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(i,a)=>(n(),x(`div`,Mt,[b(`div`,Nt,[b(`div`,null,[b(`h2`,Pt,f(r(s)(`job_management.subordinate_controls`)),1),b(`p`,Ft,f(r(s)(`job_management.subordinate_description`)),1)]),p(r(E),{label:r(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:a[0]||=e=>R()},null,8,[`label`])]),p(Z,{items:u.value,loading:m.value,total:h.value,columns:N.value,entity:`subordinate-controls`,"org-id":e.orgId,"on-load":I,onEdit:z,onDelete:H},{empty:t(()=>[b(`div`,It,[a[8]||=b(`i`,{class:`pi pi-sitemap text-3xl mb-2 opacity-50`},null,-1),b(`p`,Lt,f(r(s)(`job_management.empty_subordinates`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:g.value,"onUpdate:visible":a[4]||=e=>g.value=e,title:_.value?r(s)(`common.edit`):r(s)(`common.create`),saving:S.value,errors:C.value,onSave:V,onCancel:a[5]||=e=>g.value=!1},{default:t(()=>[p(L,{label:r(s)(`organization.nomenclature`),required:``,errors:C.value?.nomenclature},{default:t(()=>[p(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":a[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:d({"p-invalid":C.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`organization.full_code`),required:``,errors:C.value?.full_code},{default:t(()=>[p(B,{modelValue:j.value.full_code,"onUpdate:modelValue":a[2]||=e=>j.value.full_code=e,maxlength:`20`,class:d({"p-invalid":C.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.control_type`),errors:C.value?.job_management_value_id},{default:t(()=>[p(G,{modelValue:j.value.job_management_value_id,"onUpdate:modelValue":a[3]||=e=>j.value.job_management_value_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":a[6]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:U,onCancel:a[7]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Bt={class:`space-y-4`},Vt={class:`flex items-center justify-between`},Ht={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Ut={class:`text-sm text-gray-500 dark:text-gray-400`},Wt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},Gt={class:`text-sm font-medium`},Kt=`/api/v1/tenant/job-management/assets`,qt={__name:`JobAssetSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:i}){let a=e,o=i,{t:s}=P(),c=A(),u=l([]),m=l(!1),h=l(0),g=l(!1),_=l(!1),y=l(``),S=l(!1),C=l({}),T=l(!1),D=l(!1),O=l(``),k=l(null),j=l({nomenclature:``,full_code:``,job_management_value_asset_id:``,job_management_value_authority_id:``}),M=v(()=>a.jobValueMap?.asset||[]),N=v(()=>a.jobValueMap?.authority||[]),I=v(()=>[{field:`nomenclature`,header:s(`organization.nomenclature`)},{field:`full_code`,header:s(`organization.full_code`)}]);async function R(e,t){m.value=!0;try{let n=await w.get(Kt,{params:{page:e,per_page:t,organization_id:a.orgId}});u.value=n.data?.data||[],h.value=n.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{m.value=!1}}function z(){_.value=!1,y.value=``,j.value={nomenclature:``,full_code:``,job_management_value_asset_id:``,job_management_value_authority_id:``},C.value={},g.value=!0}function V(e){_.value=!0,y.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_asset_id:e.job_management_value_asset_id||``,job_management_value_authority_id:e.job_management_value_authority_id||``},C.value={},g.value=!0}async function H(){S.value=!0,C.value={};try{let e={...j.value,organization_id:a.orgId};_.value?await w.put(`${Kt}/${y.value}`,e):await w.post(Kt,e),g.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),R(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{S.value=!1}}function U(e){k.value=e,O.value=``,T.value=!0}async function W(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${Kt}/${k.value.id}`),T.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),R(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(i,a)=>(n(),x(`div`,Bt,[b(`div`,Vt,[b(`div`,null,[b(`h2`,Ht,f(r(s)(`job_management.assets`)),1),b(`p`,Ut,f(r(s)(`job_management.asset_description`)),1)]),p(r(E),{label:r(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:a[0]||=e=>z()},null,8,[`label`])]),p(Z,{items:u.value,loading:m.value,total:h.value,columns:I.value,entity:`assets`,"org-id":e.orgId,"on-load":R,onEdit:V,onDelete:U},{empty:t(()=>[b(`div`,Wt,[a[9]||=b(`i`,{class:`pi pi-box text-3xl mb-2 opacity-50`},null,-1),b(`p`,Gt,f(r(s)(`job_management.empty_assets`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:g.value,"onUpdate:visible":a[5]||=e=>g.value=e,title:_.value?r(s)(`common.edit`):r(s)(`common.create`),saving:S.value,errors:C.value,onSave:H,onCancel:a[6]||=e=>g.value=!1},{default:t(()=>[p(L,{label:r(s)(`organization.nomenclature`),required:``,errors:C.value?.nomenclature},{default:t(()=>[p(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":a[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:d({"p-invalid":C.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`organization.full_code`),required:``,errors:C.value?.full_code},{default:t(()=>[p(B,{modelValue:j.value.full_code,"onUpdate:modelValue":a[2]||=e=>j.value.full_code=e,maxlength:`20`,class:d({"p-invalid":C.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.asset_type`),errors:C.value?.job_management_value_asset_id},{default:t(()=>[p(G,{modelValue:j.value.job_management_value_asset_id,"onUpdate:modelValue":a[3]||=e=>j.value.job_management_value_asset_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.authority_level`),errors:C.value?.job_management_value_authority_id},{default:t(()=>[p(G,{modelValue:j.value.job_management_value_authority_id,"onUpdate:modelValue":a[4]||=e=>j.value.job_management_value_authority_id=e,options:N.value,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":a[7]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:W,onCancel:a[8]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Jt={class:`space-y-4`},Yt={class:`flex items-center justify-between`},Xt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Zt={class:`text-sm text-gray-500 dark:text-gray-400`},Qt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},$t={class:`text-sm font-medium`},en=`/api/v1/tenant/job-management/financials`,tn={__name:`JobFinancialSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:i}){let a=e,o=i,{t:s}=P(),c=A(),u=l([]),m=l(!1),h=l(0),g=l(!1),_=l(!1),y=l(``),S=l(!1),C=l({}),T=l(!1),D=l(!1),O=l(``),k=l(null),j=l({nomenclature:``,full_code:``,is_authorized:!1,job_management_value_cash_id:``,job_management_value_authority_id:``,job_management_value_impact_id:``}),M=v(()=>a.jobValueMap?.cash||[]),N=v(()=>a.jobValueMap?.authority||[]),I=v(()=>a.jobValueMap?.impact||[]),R=v(()=>[{field:`nomenclature`,header:s(`organization.nomenclature`)},{field:`full_code`,header:s(`organization.full_code`)},{field:`is_authorized`,header:s(`job_management.is_authorized`)}]);async function z(e,t){m.value=!0;try{let n=await w.get(en,{params:{page:e,per_page:t,organization_id:a.orgId}});u.value=n.data?.data||[],h.value=n.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{m.value=!1}}function V(){_.value=!1,y.value=``,j.value={nomenclature:``,full_code:``,is_authorized:!1,job_management_value_cash_id:``,job_management_value_authority_id:``,job_management_value_impact_id:``},C.value={},g.value=!0}function H(e){_.value=!0,y.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,is_authorized:!!e.is_authorized,job_management_value_cash_id:e.job_management_value_cash_id||``,job_management_value_authority_id:e.job_management_value_authority_id||``,job_management_value_impact_id:e.job_management_value_impact_id||``},C.value={},g.value=!0}async function U(){S.value=!0,C.value={};try{let e={...j.value,organization_id:a.orgId};_.value?await w.put(`${en}/${y.value}`,e):await w.post(en,e),g.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),z(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{S.value=!1}}function W(e){k.value=e,O.value=``,T.value=!0}async function ee(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${en}/${k.value.id}`),T.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),z(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(i,a)=>(n(),x(`div`,Jt,[b(`div`,Yt,[b(`div`,null,[b(`h2`,Xt,f(r(s)(`job_management.financials`)),1),b(`p`,Zt,f(r(s)(`job_management.financial_description`)),1)]),p(r(E),{label:r(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:a[0]||=e=>V()},null,8,[`label`])]),p(Z,{items:u.value,loading:m.value,total:h.value,columns:R.value,entity:`financials`,"org-id":e.orgId,"on-load":z,onEdit:H,onDelete:W},{empty:t(()=>[b(`div`,Qt,[a[11]||=b(`i`,{class:`pi pi-money-bill text-3xl mb-2 opacity-50`},null,-1),b(`p`,$t,f(r(s)(`job_management.empty_financials`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:g.value,"onUpdate:visible":a[7]||=e=>g.value=e,title:_.value?r(s)(`common.edit`):r(s)(`common.create`),saving:S.value,errors:C.value,onSave:U,onCancel:a[8]||=e=>g.value=!1},{default:t(()=>[p(L,{label:r(s)(`organization.nomenclature`),required:``,errors:C.value?.nomenclature},{default:t(()=>[p(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":a[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:d({"p-invalid":C.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`organization.full_code`),required:``,errors:C.value?.full_code},{default:t(()=>[p(B,{modelValue:j.value.full_code,"onUpdate:modelValue":a[2]||=e=>j.value.full_code=e,maxlength:`20`,class:d({"p-invalid":C.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.is_authorized`),class:`md:col-span-2`},{default:t(()=>[p(r(te),{modelValue:j.value.is_authorized,"onUpdate:modelValue":a[3]||=e=>j.value.is_authorized=e},null,8,[`modelValue`])]),_:1},8,[`label`]),p(L,{label:r(s)(`job_management.cash_level`),errors:C.value?.job_management_value_cash_id},{default:t(()=>[p(G,{modelValue:j.value.job_management_value_cash_id,"onUpdate:modelValue":a[4]||=e=>j.value.job_management_value_cash_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.authority_level`),errors:C.value?.job_management_value_authority_id},{default:t(()=>[p(G,{modelValue:j.value.job_management_value_authority_id,"onUpdate:modelValue":a[5]||=e=>j.value.job_management_value_authority_id=e,options:N.value,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.impact_level`),errors:C.value?.job_management_value_impact_id},{default:t(()=>[p(G,{modelValue:j.value.job_management_value_impact_id,"onUpdate:modelValue":a[6]||=e=>j.value.job_management_value_impact_id=e,options:I.value,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":a[9]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:ee,onCancel:a[10]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},nn={class:`space-y-4`},rn={class:`flex items-center justify-between`},an={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},on={class:`text-sm text-gray-500 dark:text-gray-400`},sn={class:`flex flex-col items-center justify-center py-10 text-gray-400`},cn={class:`text-sm font-medium`},ln=`/api/v1/tenant/job-management/potency-competencies`,un={__name:`JobPotencySection`,props:{orgId:String,jobValueMap:Object,competencyOptions:Array},emits:[`saved`],setup(e,{emit:i}){let a=e,o=i,{t:s}=P(),c=A(),u=l([]),m=l(!1),h=l(0),g=l(!1),_=l(!1),y=l(``),S=l(!1),C=l({}),T=l(!1),D=l(!1),O=l(``),k=l(null),j=l({competency_id:``,job_management_value_id:``,weight:null}),M=v(()=>Object.values(a.jobValueMap||{}).flat()),N=v(()=>[{field:`_competency`,header:s(`job_management.competency`)},{field:`weight`,header:s(`job_management.weight`)}]);async function I(e,t){m.value=!0;try{let n=await w.get(ln,{params:{page:e,per_page:t,organization_id:a.orgId}});u.value=(n.data?.data||[]).map(e=>{let t=a.competencyOptions?.find(t=>t.value===e.competency_id);return{...e,_competency:t?.label||e.competency_id}}),h.value=n.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{m.value=!1}}function R(){_.value=!1,y.value=``,j.value={competency_id:``,job_management_value_id:``,weight:null},C.value={},g.value=!0}function z(e){_.value=!0,y.value=e.id,j.value={competency_id:e.competency_id||``,job_management_value_id:e.job_management_value_id||``,weight:e.weight??null},C.value={},g.value=!0}async function B(){S.value=!0,C.value={};try{let e={...j.value,organization_id:a.orgId};_.value?await w.put(`${ln}/${y.value}`,e):await w.post(ln,e),g.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),I(1,15)}catch(e){let t=F(e);Object.keys(t).length?C.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{S.value=!1}}function H(e){k.value=e,O.value=``,T.value=!0}async function U(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${ln}/${k.value.id}`),T.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),I(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{D.value=!1}}}return(i,a)=>(n(),x(`div`,nn,[b(`div`,rn,[b(`div`,null,[b(`h2`,an,f(r(s)(`job_management.potency_competencies`)),1),b(`p`,on,f(r(s)(`job_management.potency_description`)),1)]),p(r(E),{label:r(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:a[0]||=e=>R()},null,8,[`label`])]),p(Z,{items:u.value,loading:m.value,total:h.value,columns:N.value,entity:`potency-competencies`,"org-id":e.orgId,"on-load":I,onEdit:z,onDelete:H},{empty:t(()=>[b(`div`,sn,[a[8]||=b(`i`,{class:`pi pi-star text-3xl mb-2 opacity-50`},null,-1),b(`p`,cn,f(r(s)(`job_management.empty_potency`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:g.value,"onUpdate:visible":a[4]||=e=>g.value=e,title:_.value?r(s)(`common.edit`):r(s)(`common.create`),saving:S.value,errors:C.value,onSave:B,onCancel:a[5]||=e=>g.value=!1},{default:t(()=>[p(L,{label:r(s)(`job_management.competency`),required:``,errors:C.value?.competency_id},{default:t(()=>[p(G,{modelValue:j.value.competency_id,"onUpdate:modelValue":a[1]||=e=>j.value.competency_id=e,options:e.competencyOptions,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.value_ref`),errors:C.value?.job_management_value_id},{default:t(()=>[p(G,{modelValue:j.value.job_management_value_id,"onUpdate:modelValue":a[2]||=e=>j.value.job_management_value_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.weight`),errors:C.value?.weight},{default:t(()=>[p(r(V),{modelValue:j.value.weight,"onUpdate:modelValue":a[3]||=e=>j.value.weight=e,min:0,max:100,class:d([{"p-invalid":C.value?.weight},`w-full`]),size:`small`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:T.value,"onUpdate:visible":a[6]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:U,onCancel:a[7]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},dn={class:`space-y-4`},fn={class:`flex items-center justify-between`},pn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},mn={class:`text-sm text-gray-500 dark:text-gray-400`},hn={class:`flex flex-col items-center justify-center py-10 text-gray-400`},gn={class:`text-sm font-medium`},_n=`/api/v1/tenant/job-management/competency-groups`,vn={__name:`JobCompetencyGroupSection`,props:{orgId:String},emits:[`saved`],setup(e,{emit:i}){let a=e,o=i,{t:s}=P(),c=A(),u=l([]),m=l(!1);l(0);let h=l(!1),g=l(!1),_=l(``),y=l(!1),S=l({}),C=l(!1),T=l(!1),D=l(``),O=l(null),k=l({category:``,weight:null}),j=v(()=>[{label:`${s(`job_management.technical`)} (${s(`job_management.category`)})`,value:`technical`},{label:`${s(`job_management.managerial`)} (${s(`job_management.category`)})`,value:`managerial`}]),M=v(()=>[{field:`category`,header:s(`job_management.category`)},{field:`weight`,header:s(`job_management.weight`)}]);async function N(){m.value=!0;try{let e=await w.get(_n,{params:{organization_id:a.orgId}});u.value=e.data?.data||(Array.isArray(e.data)?e.data:[])}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{m.value=!1}}function I(){g.value=!1,_.value=``,k.value={category:`technical`,weight:null},S.value={},h.value=!0}function R(e){g.value=!0,_.value=e.id,k.value={category:e.category||`technical`,weight:e.weight??null},S.value={},h.value=!0}async function z(){y.value=!0,S.value={};try{let e={...k.value,organization_id:a.orgId};g.value?await w.put(`${_n}/${_.value}`,e):await w.post(_n,e),h.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),N()}catch(e){let t=F(e);Object.keys(t).length?S.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{y.value=!1}}function B(e){O.value=e,D.value=``,C.value=!0}async function H(){if(O.value){T.value=!0,D.value=``;try{await w.delete(`${_n}/${O.value.id}`),C.value=!1,o(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),N()}catch(e){D.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{T.value=!1}}}return(i,a)=>(n(),x(`div`,dn,[b(`div`,fn,[b(`div`,null,[b(`h2`,pn,f(r(s)(`job_management.competency_groups`)),1),b(`p`,mn,f(r(s)(`job_management.competency_group_description`)),1)]),p(r(E),{label:r(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:a[0]||=e=>I()},null,8,[`label`])]),p(Z,{items:u.value,loading:m.value,total:u.value.length,columns:M.value,entity:`competency-groups`,"org-id":e.orgId,"on-load":N,onEdit:R,onDelete:B},{empty:t(()=>[b(`div`,hn,[a[7]||=b(`i`,{class:`pi pi-chart-pie text-3xl mb-2 opacity-50`},null,-1),b(`p`,gn,f(r(s)(`job_management.empty_competency_groups`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),p(Q,{visible:h.value,"onUpdate:visible":a[3]||=e=>h.value=e,title:g.value?r(s)(`common.edit`):r(s)(`common.create`),saving:y.value,errors:S.value,onSave:z,onCancel:a[4]||=e=>h.value=!1},{default:t(()=>[p(L,{label:r(s)(`job_management.category`),required:``,errors:S.value?.category},{default:t(()=>[p(G,{modelValue:k.value.category,"onUpdate:modelValue":a[1]||=e=>k.value.category=e,options:j.value,optionLabel:`label`,optionValue:`value`,placeholder:r(s)(`common.select`),class:d({"p-invalid":S.value?.category})},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),p(L,{label:r(s)(`job_management.weight`),required:``,errors:S.value?.weight},{default:t(()=>[p(r(V),{modelValue:k.value.weight,"onUpdate:modelValue":a[2]||=e=>k.value.weight=e,min:0,max:100,suffix:`%`,class:d([{"p-invalid":S.value?.weight},`w-full`]),size:`small`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),p(K,{visible:C.value,"onUpdate:visible":a[5]||=e=>C.value=e,loading:T.value,"error-msg":D.value,onConfirm:H,onCancel:a[6]||=e=>C.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},yn={class:`space-y-6`},bn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},xn={class:`text-sm text-gray-500 dark:text-gray-400`},Sn={key:0,class:`flex items-center justify-center py-12`},Cn={class:`grid grid-cols-1 md:grid-cols-3 gap-4`},wn={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Tn={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},En={class:`text-2xl font-bold text-emerald-600 dark:text-emerald-400`},Dn={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},On={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},kn={class:`text-2xl font-bold text-blue-600 dark:text-blue-400`},An={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},jn={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},Mn={key:0,class:`text-[10px] text-gray-400 mt-2`},Nn={key:0,class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden`},Pn={class:`px-5 py-3 border-b border-gray-200 dark:border-gray-700 font-semibold text-sm text-gray-700 dark:text-gray-300`},Fn={class:`p-5`},In={class:`text-sm text-gray-700 dark:text-gray-300 capitalize`},Ln={class:`text-sm font-semibold text-gray-900 dark:text-gray-100`},Rn={key:2},zn={class:`flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500`},Bn={class:`text-sm font-medium`},Vn={class:`text-xs mt-1`},Hn={class:`flex justify-end gap-3`},Un=`/api/v1/tenant/job-management/scores`,Wn={__name:`JobScoreSection`,props:{orgId:String},emits:[`saved`],setup(e,{emit:t}){let i=e,a=t,{t:o}=P(),c=A(),d=l(!1),g=l(!1),_=l(null),S=v(()=>{if(!_.value?.components)return null;try{return JSON.parse(_.value.components)}catch{return null}});function C(e){return e?.toLocaleString?.(`id-ID`)??`-`}async function T(){if(i.orgId){d.value=!0;try{let e=await w.get(`${Un}/${i.orgId}`);_.value=e.data?.data||null,a(`saved`)}catch{_.value=null}finally{d.value=!1}}}async function D(){if(i.orgId){g.value=!0;try{let e=await w.put(`${Un}/${i.orgId}`,{components:null});_.value=e.data?.data||null,c.add({severity:`success`,summary:o(`message.success`),detail:o(`job_management.score_calculated`),life:2e3})}catch(e){c.add({severity:`error`,summary:o(`message.error`),detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{g.value=!1}}}return s(T),(e,t)=>(n(),x(`div`,yn,[b(`div`,null,[b(`h2`,bn,f(r(o)(`job_management.scores`)),1),b(`p`,xn,f(r(o)(`job_management.score_description`)),1)]),d.value?(n(),x(`div`,Sn,[...t[0]||=[b(`i`,{class:`pi pi-spin pi-spinner text-emerald-500 text-2xl`},null,-1)]])):_.value?(n(),x(y,{key:1},[b(`div`,Cn,[b(`div`,wn,[b(`div`,Tn,f(r(o)(`job_management.value_with_financial`)),1),b(`div`,En,f(C(_.value.job_value_with_financial)),1)]),b(`div`,Dn,[b(`div`,On,f(r(o)(`job_management.value_without_financial`)),1),b(`div`,kn,f(C(_.value.job_value_without_financial)),1)]),b(`div`,An,[b(`div`,jn,f(r(o)(`job_management.has_financial_authority`)),1),p(r(I),{value:_.value.has_financial_authority?r(o)(`common.yes`):r(o)(`common.no`),severity:_.value.has_financial_authority?`success`:`danger`,class:`!text-xs`},null,8,[`value`,`severity`]),_.value.calculated_at?(n(),x(`div`,Mn,f(r(o)(`job_management.calculated_at`))+`: `+f(_.value.calculated_at),1)):h(``,!0)])]),S.value?(n(),x(`div`,Nn,[b(`div`,Pn,f(r(o)(`job_management.component_breakdown`)),1),b(`div`,Fn,[(n(!0),x(y,null,m(S.value,(e,t)=>(n(),x(`div`,{key:t,class:`flex items-center justify-between py-1.5 border-b border-gray-100 dark:border-gray-700 last:border-0`},[b(`span`,In,f(t.replace(/_/g,` `)),1),b(`span`,Ln,f(C(e)),1)]))),128))])])):h(``,!0)],64)):(n(),x(`div`,Rn,[b(`div`,zn,[t[1]||=b(`i`,{class:`pi pi-calculator text-4xl mb-3 opacity-50`},null,-1),b(`p`,Bn,f(r(o)(`job_management.no_score`)),1),b(`p`,Vn,f(r(o)(`job_management.score_hint`)),1)])])),b(`div`,Hn,[p(r(E),{label:r(o)(`common.refresh`),icon:`pi pi-refresh`,size:`small`,text:``,onClick:T},null,8,[`label`]),_.value?(n(),u(r(E),{key:0,label:r(o)(`job_management.recalculate`),icon:`pi pi-calculator`,size:`small`,severity:`info`,loading:g.value,onClick:D},null,8,[`label`,`loading`])):h(``,!0)])]))}},Gn={class:`max-w-full mx-auto`},Kn={key:0,class:`flex gap-6`},qn={class:`w-56 space-y-2`},Jn={class:`flex-1 space-y-3`},Yn={key:1,class:`flex gap-6`},Xn={class:`w-56 shrink-0 space-y-1`},Zn=[`onClick`,`onKeydown`],Qn={key:0,class:`pi pi-check text-xs`},$n={class:`flex-1 min-w-0`},er={key:0,class:`pi pi-check-circle text-emerald-400 text-xs shrink-0`},tr={class:`flex-1 min-w-0`},nr={__name:`JobManagementForm`,setup(e){let t=N(),i=k(),{t:o}=P(),c=A(),p=i.query.org_id||``,g=l(0),_=l(!0),S=l(Array(15).fill(!1)),C=l(``),E=l(``),D=l(``),O=l(``),j=l([]),M=l([]),F=l([]),I=l({}),L=l([]),R=[{labelKey:`job_management.identifications`,icon:`pi pi-id-card`,comp:le},{labelKey:`job_management.objectives`,icon:`pi pi-bullseye`,comp:ge},{labelKey:`job_management.responsibilities_title`,icon:`pi pi-list-check`,comp:Ke},{labelKey:`job_management.education_experience`,icon:`pi pi-graduation-cap`,comp:Ce},{labelKey:`job_management.hr_authorities`,icon:`pi pi-users`,comp:et},{labelKey:`job_management.op_authorities`,icon:`pi pi-cog`,comp:ct},{labelKey:`job_management.activities`,icon:`pi pi-bolt`,comp:gt},{labelKey:`job_management.risks`,icon:`pi pi-exclamation-triangle`,comp:Ct},{labelKey:`job_management.relationships`,icon:`pi pi-share-alt`,comp:jt},{labelKey:`job_management.subordinate_controls`,icon:`pi pi-sitemap`,comp:zt},{labelKey:`job_management.assets`,icon:`pi pi-box`,comp:qt},{labelKey:`job_management.financials`,icon:`pi pi-money-bill`,comp:tn},{labelKey:`job_management.potency_competencies`,icon:`pi pi-star`,comp:un},{labelKey:`job_management.competency_groups`,icon:`pi pi-chart-pie`,comp:vn},{labelKey:`job_management.scores`,icon:`pi pi-calculator`,comp:Wn}],z=v(()=>R[g.value]?.comp||null);function B(e){return g.value===e?`bg-emerald-50 dark:bg-emerald-900/20 ring-1 ring-emerald-300 dark:ring-emerald-700`:(S.value[e],`hover:bg-gray-50 dark:hover:bg-gray-800`)}function V(e){return g.value===e?`bg-emerald-600 text-white`:S.value[e]?`bg-emerald-100 dark:bg-emerald-800 text-emerald-600 dark:text-emerald-300`:`bg-gray-200 dark:bg-gray-600 text-gray-600 dark:text-gray-300`}function H(e){return g.value===e?`text-emerald-700 dark:text-emerald-300`:S.value[e]?`text-emerald-600 dark:text-emerald-400`:`text-gray-700 dark:text-gray-300`}function U(e){g.value=e,t.replace({query:{...i.query,section:String(e)}})}function W(e){typeof e==`number`&&(S.value[e]=!0)}async function G(){if(p)try{let e=(await w.get(`/api/v1/tenant/organizations/${p}`)).data?.data;e&&(C.value=e.nomenclature||``,E.value=e.full_code||e.code||``,D.value=e.grading_id||``,O.value=e.job_family_id||``)}catch{}}async function K(){try{let[e,t,n,r]=await Promise.all([w.get(`/api/v1/tenant/settings/gradings?per_page=100`),w.get(`/api/v1/tenant/job-management/values?per_page=200`),w.get(`/api/v1/tenant/competencies?per_page=200`).catch(()=>({data:{data:[]}})),w.get(`/api/v1/tenant/settings/job-families?per_page=100`)]);j.value=(e.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id})),M.value=(r.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id}));let i=t.data?.data||[];F.value=i.map(e=>({label:`${e.type}${e.level?` Lv.`+e.level:``}${e.descriptions?` — `+e.descriptions:``}`,value:e.id,type:e.type,level:e.level,descriptions:e.descriptions}));let a={};i.forEach(e=>{a[e.type]||(a[e.type]=[]),a[e.type].push({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level})}),I.value=a,L.value=(n.data?.data||[]).map(e=>({label:e.name||e.code,value:e.id}))}catch{}}return s(async()=>{try{await Promise.all([G(),K()]);let e=parseInt(i.query.section);!isNaN(e)&&e>=0&&e<R.length&&(g.value=e)}catch(e){c.add({severity:`error`,summary:o(`message.error`),detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{_.value=!1}}),(e,t)=>(n(),x(`div`,Gn,[_.value?(n(),x(`div`,Kn,[b(`div`,qn,[(n(),x(y,null,m(8,e=>b(`div`,{key:e,class:`h-12 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))]),b(`div`,Jn,[(n(),x(y,null,m(6,e=>b(`div`,{key:e,class:`h-9 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))])])):(n(),x(`div`,Yn,[b(`div`,Xn,[(n(),x(y,null,m(R,(e,t)=>b(`div`,{key:t,role:`button`,tabindex:0,class:d([`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150 cursor-pointer select-none`,B(t)]),onClick:e=>U(t),onKeydown:T(e=>U(t),[`enter`])},[b(`div`,{class:d([`w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0 transition-colors duration-150`,V(t)])},[S.value[t]?(n(),x(`i`,Qn)):(n(),x(`i`,{key:1,class:d(e.icon)},null,2))],2),b(`div`,$n,[b(`div`,{class:d([`text-sm font-medium truncate`,H(t)])},f(r(o)(e.labelKey)),3)]),S.value[t]?(n(),x(`i`,er)):h(``,!0)],42,Zn)),64))]),b(`div`,tr,[(n(),u(a(z.value),{key:g.value,"org-id":r(p),"org-name":C.value,"org-code":E.value,"org-grading-id":D.value,"org-job-family-id":O.value,"grading-options":j.value,"job-family-options":M.value,"job-value-options":F.value,"competency-options":L.value,"job-value-map":I.value,onSaved:W},null,40,[`org-id`,`org-name`,`org-code`,`org-grading-id`,`org-job-family-id`,`grading-options`,`job-family-options`,`job-value-options`,`competency-options`,`job-value-map`]))])]))]))}};export{nr as default};