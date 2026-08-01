const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["assets/quill-BLmY9xB4.js","assets/rolldown-runtime-QTnfLwEv.js"])))=>i.map(i=>d[i]);
import{B as e,E as t,M as n,O as r,U as i,c as a,ct as o,j as s,k as c,l,lt as u,m as d,o as f,p,q as m,r as h,s as g,st as _,u as v,ut as y,v as b,w as x,x as S,z as C}from"./runtime-core.esm-bundler-CHUeMfpT.js";import{a as w,d as T,t as E}from"./button-CvK-Qhr7.js";import{A as D,a as O}from"./ripple-BFhyWgdB.js";import{l as k,m as A,s as j,t as M,u as N}from"./index-BBw7ZjWm.js";import{t as P}from"./useI18n-DxyS4lE5.js";import{n as F}from"./responseHandler-B5MnXl3B.js";import{t as I}from"./tag-DghSMZF_.js";import{t as L}from"./FormRow-CUHrbwQr.js";import{t as R}from"./baseeditableholder-BnNy24a5.js";import{t as z}from"./textarea-DFx4Mmyi.js";import{t as B}from"./TextInput-RytNn3EW.js";import{i as V,n as H,o as U,t as W}from"./column-B-oLVtB2.js";import{t as G}from"./SelectLabel-CH3pe4id.js";import{t as K}from"./ConfirmDeleteDialog-0ot3hrgz.js";import{t as ee}from"./SkeletonTable-DpqTHC0p.js";import{t as te}from"./toggleswitch-Cge-VV6w.js";var ne={class:`space-y-4`},re={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},ie={class:`text-sm text-gray-500 dark:text-gray-400`},ae={class:`max-w-2xl space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},oe={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},se={class:`flex justify-end pt-2`},ce=`/api/v1/tenant/job-management/identifications`,le={__name:`JobIdentificationSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},orgGradingId:{type:String,default:``},orgJobFamilyId:{type:String,default:``},gradingOptions:{type:Array,default:()=>[]},jobFamilyOptions:{type:Array,default:()=>[]}},emits:[`saved`],setup(e,{emit:n}){let r=n,a=e,{t:o}=P(),s=A(),c=i(!1),u=i(``),p=i({}),h=i(``),_=i({grading_id:``}),b=f(()=>{let e=a.jobFamilyOptions.find(e=>e.value===a.orgJobFamilyId);return e?e.label:a.orgJobFamilyId||`-`});function S(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function T(){if(a.orgId)try{let e=(await w.get(ce,{params:{organization_id:a.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];h.value=t.id,_.value.grading_id=t.grading_id||a.orgGradingId||``}else _.value.grading_id=a.orgGradingId||``}catch{_.value.grading_id=a.orgGradingId||``}}async function D(){if(u.value=``,p.value={},!_.value.grading_id){u.value=o(`job_management.grading_required`);return}c.value=!0;try{let e={nomenclature:a.orgName||``,full_code:a.orgCode||``,grading_id:_.value.grading_id,organization_id:a.orgId};if(h.value)await w.put(`${ce}/${h.value}`,{grading_id:_.value.grading_id});else{let t=await w.post(ce,e);h.value=t.data?.data?.id||``}s.add({severity:`success`,summary:o(`message.success`),detail:o(`common.saved`),life:2e3}),r(`saved`)}catch(e){let t=S(e);Object.keys(t).length>0?(p.value=t,u.value=Object.values(t).join(`, `)):u.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{c.value=!1}}return x(T),(n,r)=>(t(),v(`div`,ne,[g(`div`,null,[g(`h2`,re,y(m(o)(`job_management.identifications`)),1),g(`p`,ie,y(m(o)(`job_management.identification_description`)),1)]),g(`div`,ae,[d(L,{label:m(o)(`organization.nomenclature`)},{default:C(()=>[d(B,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),d(L,{label:m(o)(`organization.full_code`)},{default:C(()=>[d(B,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),d(L,{label:m(o)(`organization.job_family`)},{default:C(()=>[d(B,{"model-value":b.value,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),d(L,{label:m(o)(`organization.grading`)},{default:C(()=>[d(m(U),{modelValue:_.value.grading_id,"onUpdate:modelValue":r[0]||=e=>_.value.grading_id=e,options:e.gradingOptions,"option-label":`label`,"option-value":`value`,placeholder:m(o)(`organization.select_grading`),class:`w-full`,size:`small`,invalid:!!p.value.grading_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`]),u.value?(t(),v(`div`,oe,y(u.value),1)):l(``,!0),g(`div`,se,[d(m(E),{label:m(o)(`common.save`),icon:`pi pi-check`,size:`small`,loading:c.value,disabled:!_.value.grading_id,onClick:D},null,8,[`label`,`loading`,`disabled`])])])]))}},ue={class:`space-y-4`},de={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},fe={class:`text-sm text-gray-500 dark:text-gray-400`},pe={class:`max-w-2xl space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},me={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},he={class:`flex justify-end gap-2 pt-2`},q=`/api/v1/tenant/job-management/objectives`,ge={__name:`JobObjectiveSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:n}){let r=n,o=e,{t:s}=P(),c=A(),u=i(!1),f=i(!1),p=i(``),h=i({}),b=i(``),S=i(!1),T=i(``),D=i({objective:``});function O(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function k(){if(o.orgId)try{let e=(await w.get(q,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];b.value=t.id,D.value.objective=t.objective||``}}catch{}}async function j(){p.value=``,h.value={},u.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,objective:D.value.objective||``,organization_id:o.orgId};if(b.value)await w.put(`${q}/${b.value}`,{objective:D.value.objective||``});else{let t=await w.post(q,e);b.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),r(`saved`)}catch(e){let t=O(e);Object.keys(t).length>0?(h.value=t,p.value=Object.values(t).join(`, `)):p.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{u.value=!1}}async function M(){if(b.value){f.value=!0,T.value=``;try{await w.delete(`${q}/${b.value}`),S.value=!1,b.value=``,D.value.objective=``,r(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{f.value=!1}}}return x(k),(n,r)=>(t(),v(`div`,ue,[g(`div`,null,[g(`h2`,de,y(m(s)(`job_management.objectives`)),1),g(`p`,fe,y(m(s)(`job_management.objective_description`)),1)]),g(`div`,pe,[d(L,{label:m(s)(`organization.nomenclature`)},{default:C(()=>[d(B,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),d(L,{label:m(s)(`organization.full_code`)},{default:C(()=>[d(B,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),d(L,{label:m(s)(`job_management.objective`)},{default:C(()=>[d(m(z),{modelValue:D.value.objective,"onUpdate:modelValue":r[0]||=e=>D.value.objective=e,rows:`3`,class:_([`w-full`,{"p-invalid":h.value.objective}]),placeholder:m(s)(`job_management.objective`)+`...`},null,8,[`modelValue`,`class`,`placeholder`])]),_:1},8,[`label`]),p.value?(t(),v(`div`,me,y(p.value),1)):l(``,!0),g(`div`,he,[b.value?(t(),a(m(E),{key:0,label:m(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:r[1]||=e=>S.value=!0},null,8,[`label`])):l(``,!0),d(m(E),{label:b.value?m(s)(`common.update`):m(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:u.value,disabled:u.value,onClick:j},null,8,[`label`,`loading`,`disabled`])])]),d(K,{visible:S.value,"onUpdate:visible":r[2]||=e=>S.value=e,loading:f.value,"error-msg":T.value,onConfirm:M,onCancel:r[3]||=e=>S.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},_e={class:`space-y-4`},ve={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},ye={class:`text-sm text-gray-500 dark:text-gray-400`},be={class:`max-w-2xl space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},xe={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Se={class:`flex justify-end gap-2 pt-2`},J=`/api/v1/tenant/job-management/education-experiences`,Ce={__name:`JobEduExpSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:n}){let r=n,o=e,{t:s}=P(),c=A(),u=i(!1),p=i(!1),h=i(``),_=i({}),b=i(``),S=i(!1),T=i(``),D=i({job_management_value_education_id:``,job_management_value_experience_id:``}),O=f(()=>o.jobValueMap?.education||[]),k=f(()=>o.jobValueMap?.experience||[]);async function j(){if(o.orgId)try{let e=(await w.get(J,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];b.value=t.id,D.value.job_management_value_education_id=t.job_management_value_education_id||``,D.value.job_management_value_experience_id=t.job_management_value_experience_id||``}}catch{}}async function M(){h.value=``,_.value={},u.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,job_management_value_education_id:D.value.job_management_value_education_id||null,job_management_value_experience_id:D.value.job_management_value_experience_id||null,organization_id:o.orgId};if(b.value)await w.put(`${J}/${b.value}`,{job_management_value_education_id:D.value.job_management_value_education_id||null,job_management_value_experience_id:D.value.job_management_value_experience_id||null});else{let t=await w.post(J,e);b.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),r(`saved`)}catch(e){let t=F(e);Object.keys(t).length>0?(_.value=t,h.value=Object.values(t).join(`, `)):h.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{u.value=!1}}async function N(){if(b.value){p.value=!0,T.value=``;try{await w.delete(`${J}/${b.value}`),S.value=!1,b.value=``,D.value.job_management_value_education_id=``,D.value.job_management_value_experience_id=``,r(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{p.value=!1}}}return x(j),(n,r)=>(t(),v(`div`,_e,[g(`div`,null,[g(`h2`,ve,y(m(s)(`job_management.education_experience`)),1),g(`p`,ye,y(m(s)(`job_management.education_experience_description`)),1)]),g(`div`,be,[d(L,{label:m(s)(`organization.nomenclature`)},{default:C(()=>[d(B,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),d(L,{label:m(s)(`organization.full_code`)},{default:C(()=>[d(B,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),d(L,{label:m(s)(`job_management.education_level`),errors:_.value?.job_management_value_education_id},{default:C(()=>[d(m(U),{modelValue:D.value.job_management_value_education_id,"onUpdate:modelValue":r[0]||=e=>D.value.job_management_value_education_id=e,options:O.value,"option-label":`label`,"option-value":`value`,placeholder:m(s)(`common.select`),class:`w-full`,size:`small`,showClear:``,invalid:!!_.value.job_management_value_education_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(s)(`job_management.experience_level`),errors:_.value?.job_management_value_experience_id},{default:C(()=>[d(m(U),{modelValue:D.value.job_management_value_experience_id,"onUpdate:modelValue":r[1]||=e=>D.value.job_management_value_experience_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:m(s)(`common.select`),class:`w-full`,size:`small`,showClear:``,invalid:!!_.value.job_management_value_experience_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`]),h.value?(t(),v(`div`,xe,y(h.value),1)):l(``,!0),g(`div`,Se,[b.value?(t(),a(m(E),{key:0,label:m(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:r[2]||=e=>S.value=!0},null,8,[`label`])):l(``,!0),d(m(E),{label:b.value?m(s)(`common.update`):m(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:u.value,disabled:u.value,onClick:M},null,8,[`label`,`loading`,`disabled`])])]),d(K,{visible:S.value,"onUpdate:visible":r[3]||=e=>S.value=e,loading:p.value,"error-msg":T.value,onConfirm:N,onCancel:r[4]||=e=>S.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},we=O.extend({name:`editor`,style:`
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
`,classes:{root:function(e){return[`p-editor`,{"p-invalid":e.instance.$invalid}]},toolbar:`p-editor-toolbar`,content:`p-editor-content`}}),Te={name:`BaseEditor`,extends:R,props:{placeholder:String,readonly:Boolean,formats:Array,editorStyle:null,modules:null},style:we,provide:function(){return{$pcEditor:this,$parentInstance:this}}};function Y(e){"@babel/helpers - typeof";return Y=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},Y(e)}function Ee(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function De(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?Ee(Object(n),!0).forEach(function(t){Oe(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):Ee(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function Oe(e,t,n){return(t=ke(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function ke(e){var t=Ae(e,`string`);return Y(t)==`symbol`?t:t+``}function Ae(e,t){if(Y(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(Y(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var je=function(){try{return window.Quill}catch{return null}}(),X={name:`Editor`,extends:Te,inheritAttrs:!1,emits:[`text-change`,`selection-change`,`load`],quill:null,watch:{modelValue:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},d_value:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},readonly:function(){this.handleReadOnlyChange()}},mounted:function(){var e=this,t={modules:De({toolbar:this.$refs.toolbarElement},this.modules),readOnly:this.readonly,theme:`snow`,formats:this.formats,placeholder:this.placeholder};je?(this.quill=new je(this.$refs.editorElement,t),this.initQuill(),this.handleLoad()):j(()=>import(`./quill-BLmY9xB4.js`).then(function(n){n&&D(e.$refs.editorElement)&&(n.default?e.quill=new n.default(e.$refs.editorElement,t):e.quill=new n(e.$refs.editorElement,t),e.initQuill())}),__vite__mapDeps([0,1])).then(function(){e.handleLoad()})},beforeUnmount:function(){this.quill=null},methods:{renderValue:function(e){if(this.quill)if(e){var t=this.quill.clipboard.convert({html:e});this.quill.setContents(t)}else this.quill.setText(``)},initQuill:function(){var e=this;this.renderValue(this.d_value),this.quill.on(`text-change`,function(t,n,r){if(r===`user`){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();i===`<p><br></p>`&&(i=``),e.writeValue(i),e.$emit(`text-change`,{htmlValue:i,textValue:a,delta:t,source:r,instance:e.quill})}}),this.quill.on(`selection-change`,function(t,n,r){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();e.$emit(`selection-change`,{htmlValue:i,textValue:a,range:t,oldRange:n,source:r,instance:e.quill})})},handleLoad:function(){this.quill&&this.quill.getModule(`toolbar`)&&this.$emit(`load`,{instance:this.quill})},handleReadOnlyChange:function(){this.quill&&this.quill.enable(!this.readonly)}}};function Me(e,n,r,i,a,s){return t(),v(`div`,S({class:e.cx(`root`)},e.ptmi(`root`)),[g(`div`,S({ref:`toolbarElement`,class:e.cx(`toolbar`)},e.ptm(`toolbar`)),[c(e.$slots,`toolbar`,{},function(){return[g(`span`,S({class:`ql-formats`},e.ptm(`formats`)),[g(`select`,S({class:`ql-header`,defaultValue:`0`},e.ptm(`header`)),[g(`option`,S({value:`1`},e.ptm(`option`)),`Heading`,16),g(`option`,S({value:`2`},e.ptm(`option`)),`Subheading`,16),g(`option`,S({value:`0`},e.ptm(`option`)),`Normal`,16)],16),g(`select`,S({class:`ql-font`},e.ptm(`font`)),[g(`option`,o(b(e.ptm(`option`))),null,16),g(`option`,S({value:`serif`},e.ptm(`option`)),null,16),g(`option`,S({value:`monospace`},e.ptm(`option`)),null,16)],16)],16),g(`span`,S({class:`ql-formats`},e.ptm(`formats`)),[g(`button`,S({class:`ql-bold`,type:`button`},e.ptm(`bold`)),null,16),g(`button`,S({class:`ql-italic`,type:`button`},e.ptm(`italic`)),null,16),g(`button`,S({class:`ql-underline`,type:`button`},e.ptm(`underline`)),null,16)],16),g(`span`,S({class:`ql-formats`},e.ptm(`formats`)),[g(`select`,S({class:`ql-color`},e.ptm(`color`)),null,16),g(`select`,S({class:`ql-background`},e.ptm(`background`)),null,16)],16),g(`span`,S({class:`ql-formats`},e.ptm(`formats`)),[g(`button`,S({class:`ql-list`,value:`ordered`,type:`button`},e.ptm(`list`)),null,16),g(`button`,S({class:`ql-list`,value:`bullet`,type:`button`},e.ptm(`list`)),null,16),g(`select`,S({class:`ql-align`},e.ptm(`select`)),[g(`option`,S({defaultValue:``},e.ptm(`option`)),null,16),g(`option`,S({value:`center`},e.ptm(`option`)),null,16),g(`option`,S({value:`right`},e.ptm(`option`)),null,16),g(`option`,S({value:`justify`},e.ptm(`option`)),null,16)],16)],16),g(`span`,S({class:`ql-formats`},e.ptm(`formats`)),[g(`button`,S({class:`ql-link`,type:`button`},e.ptm(`link`)),null,16),g(`button`,S({class:`ql-image`,type:`button`},e.ptm(`image`)),null,16),g(`button`,S({class:`ql-code-block`,type:`button`},e.ptm(`codeBlock`)),null,16)],16),g(`span`,S({class:`ql-formats`},e.ptm(`formats`)),[g(`button`,S({class:`ql-clean`,type:`button`},e.ptm(`clean`)),null,16)],16)]})],16),g(`div`,S({ref:`editorElement`,class:e.cx(`content`),style:e.editorStyle},e.ptm(`content`)),null,16)],16)}X.render=Me;var Ne={key:0,class:`text-gray-500 dark:text-gray-400 text-xs`},Pe=[`innerHTML`],Fe={key:2,class:`text-gray-800 dark:text-gray-100`},Ie={class:`flex items-center gap-1`},Z={__name:`DataTableSection`,props:{items:Array,loading:Boolean,total:Number,columns:{type:Array,default:()=>[]},entity:String,orgId:String,onLoad:Function},emits:[`edit`,`delete`],setup(n){let o=n,{t:u}=P(),p=i(1),_=i(15),b=f(()=>(p.value-1)*_.value),S=f(()=>[...o.columns.map(e=>({type:`text`,width:`w-24`,headerWidth:`w-20`})),{type:`icons`,count:2,headerWidth:`w-16`}]);function w(e){p.value=e.page+1,_.value=e.rows,o.onLoad&&o.onLoad(p.value,_.value)}return x(()=>{o.onLoad&&o.onLoad(1,15)}),(i,o)=>{let f=s(`tooltip`);return n.loading?(t(),a(ee,{key:0,columns:S.value,rows:8},null,8,[`columns`])):(t(),a(m(H),{key:1,value:n.items,lazy:``,totalRecords:n.total,first:b.value,rows:_.value,onPage:w,paginator:``,paginatorTemplate:`FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown`,rowsPerPageOptions:[10,15,25,50],size:`small`,class:`!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden`},{empty:C(()=>[c(i.$slots,`empty`)]),default:C(()=>[(t(!0),v(h,null,r(n.columns,e=>(t(),a(m(W),{key:e.field,field:e.field,header:e.header,sortable:``},{body:C(({data:n})=>[e.field.startsWith(`_`)?(t(),v(`span`,Ne,y(n[e.field]||`-`),1)):l(``,!0),e.html?(t(),v(`div`,{key:1,class:`editor-content`,innerHTML:n[e.field]},null,8,Pe)):(t(),v(`span`,Fe,y(n[e.field]||`-`),1))]),_:2},1032,[`field`,`header`]))),128)),d(m(W),{header:m(u)(`common.actions`),style:{width:`90px`},frozen:``,alignFrozen:`right`},{body:C(({data:t})=>[g(`div`,Ie,[e(d(m(E),{icon:`pi pi-pencil`,size:`small`,text:``,severity:`secondary`,onClick:e=>i.$emit(`edit`,t)},null,8,[`onClick`]),[[f,m(u)(`common.edit`),void 0,{left:!0}]]),e(d(m(E),{icon:`pi pi-trash`,size:`small`,text:``,severity:`danger`,onClick:e=>i.$emit(`delete`,t)},null,8,[`onClick`]),[[f,m(u)(`common.delete`),void 0,{left:!0}]])])]),_:1},8,[`header`])]),_:3},8,[`value`,`totalRecords`,`first`,`rows`]))}}},Le={class:`space-y-4`},Re={key:0,class:`bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded p-3 text-xs text-red-700 dark:text-red-300`},Q={__name:`DialogForm`,props:{visible:Boolean,title:String,saving:Boolean,errors:{type:Object,default:()=>({})},width:{type:String,default:`480px`}},emits:[`save`,`cancel`],setup(e){let n=e,{t:i}=P(),o=f(()=>n.width===`maximize`?`90vw`:n.width);return(n,s)=>(t(),a(m(M),{visible:e.visible,"onUpdate:visible":s[2]||=e=>n.$emit(`update:visible`,e),header:e.title,modal:``,style:u({width:o.value}),class:`p-fluid`,closable:!e.saving},{footer:C(()=>[d(m(E),{label:m(i)(`common.cancel`),size:`small`,outlined:``,severity:`secondary`,disabled:e.saving,onClick:s[0]||=e=>n.$emit(`cancel`)},null,8,[`label`,`disabled`]),d(m(E),{label:m(i)(`common.save`),icon:`pi pi-check`,size:`small`,loading:e.saving,onClick:s[1]||=e=>n.$emit(`save`)},null,8,[`label`,`loading`])]),default:C(()=>[g(`div`,Le,[c(n.$slots,`default`),Object.keys(e.errors).length?(t(),v(`div`,Re,[(t(!0),v(h,null,r(e.errors,(e,n)=>(t(),v(`p`,{key:n,class:`mb-1`},[g(`strong`,null,y(n)+`:`,1),p(` `+y(Array.isArray(e)?e.join(`, `):e),1)]))),128))])):l(``,!0)])]),_:3},8,[`visible`,`header`,`style`,`closable`]))}},ze={class:`space-y-4`},Be={class:`flex items-center justify-between`},Ve={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},He={class:`text-sm text-gray-500 dark:text-gray-400`},Ue={class:`flex flex-col items-center justify-center py-10 text-gray-400`},We={class:`text-sm font-medium`},Ge=`/api/v1/tenant/job-management/responsibilities`,Ke={__name:`JobResponsibilitySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:n}){let r=e,a=n,{t:o}=P(),s=A(),c=i([]),u=i(!1),p=i(0),b=i(!1),x=i(!1),S=i(``),T=i(!1),D=i({}),O=i(!1),k=i(!1),j=i(``),M=i(null),N=i({main_task:``,activities:``,outputs:``,success_indicators:``}),I=f(()=>{let e=o(`job_management.responsibilities_title`);return x.value?`${e}`:`${o(`common.create`)} ${e}`}),R=f(()=>[{field:`main_task`,header:o(`job_management.main_task`),html:!0},{field:`activities`,header:o(`job_management.activities`),html:!0},{field:`outputs`,header:o(`job_management.outputs`),html:!0},{field:`success_indicators`,header:o(`job_management.success_indicators`),html:!0}]);async function z(e,t){u.value=!0;try{let n=await w.get(Ge,{params:{page:e,per_page:t,organization_id:r.orgId}}),i=n.data?.data||[];c.value=i.map(e=>({...e,main_task:e.main_task,activities:e.activities,outputs:e.outputs,success_indicators:e.success_indicators})),p.value=n.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{u.value=!1}}function B(){x.value=!1,S.value=``,N.value={main_task:``,activities:``,outputs:``,success_indicators:``},D.value={},b.value=!0}function V(e){x.value=!0,S.value=e.id,N.value={main_task:e.main_task||``,activities:e.activities||``,outputs:e.outputs||``,success_indicators:e.success_indicators||``},D.value={},b.value=!0}async function H(){T.value=!0,D.value={};try{let e={nomenclature:r.orgName||``,full_code:r.orgCode||``,...N.value,organization_id:r.orgId};x.value?await w.put(`${Ge}/${S.value}`,e):await w.post(Ge,e),b.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),z(1,15)}catch(e){let t=F(e);Object.keys(t).length?D.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{T.value=!1}}function U(e){M.value=e,j.value=``,O.value=!0}async function W(){if(M.value){k.value=!0,j.value=``;try{await w.delete(`${Ge}/${M.value.id}`),O.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),z(1,15)}catch(e){j.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{k.value=!1}}}return(n,r)=>(t(),v(`div`,ze,[g(`div`,Be,[g(`div`,null,[g(`h2`,Ve,y(m(o)(`job_management.responsibilities_title`)),1),g(`p`,He,y(m(o)(`job_management.responsibilities_description`)),1)]),d(m(E),{label:m(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:r[0]||=e=>B()},null,8,[`label`])]),d(Z,{items:c.value,loading:u.value,total:p.value,columns:R.value,entity:`responsibilities`,"org-id":e.orgId,"on-load":z,onEdit:V,onDelete:U},{empty:C(()=>[g(`div`,Ue,[r[9]||=g(`i`,{class:`pi pi-list-check text-3xl mb-2 opacity-50`},null,-1),g(`p`,We,y(m(o)(`job_management.empty_responsibilities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),d(Q,{visible:b.value,"onUpdate:visible":r[5]||=e=>b.value=e,title:I.value,saving:T.value,errors:D.value,width:`maximize`,onSave:H,onCancel:r[6]||=e=>b.value=!1},{default:C(()=>[b.value?(t(),v(h,{key:0},[d(L,{label:m(o)(`job_management.main_task`),errors:D.value?.main_task},{default:C(()=>[d(m(X),{modelValue:N.value.main_task,"onUpdate:modelValue":r[1]||=e=>N.value.main_task=e,editorStyle:`height:120px`,class:_({"p-invalid":D.value?.main_task})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.activities`),errors:D.value?.activities},{default:C(()=>[d(m(X),{modelValue:N.value.activities,"onUpdate:modelValue":r[2]||=e=>N.value.activities=e,editorStyle:`height:120px`,class:_({"p-invalid":D.value?.activities})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.outputs`),errors:D.value?.outputs},{default:C(()=>[d(m(X),{modelValue:N.value.outputs,"onUpdate:modelValue":r[3]||=e=>N.value.outputs=e,editorStyle:`height:120px`,class:_({"p-invalid":D.value?.outputs})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.success_indicators`),errors:D.value?.success_indicators},{default:C(()=>[d(m(X),{modelValue:N.value.success_indicators,"onUpdate:modelValue":r[4]||=e=>N.value.success_indicators=e,editorStyle:`height:120px`,class:_({"p-invalid":D.value?.success_indicators})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])],64)):l(``,!0)]),_:1},8,[`visible`,`title`,`saving`,`errors`]),d(K,{visible:O.value,"onUpdate:visible":r[7]||=e=>O.value=e,loading:k.value,"error-msg":j.value,onConfirm:W,onCancel:r[8]||=e=>O.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},qe={class:`space-y-4`},Je={class:`flex items-center justify-between`},Ye={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Xe={class:`text-sm text-gray-500 dark:text-gray-400`},Ze={class:`flex flex-col items-center justify-center py-10 text-gray-400`},Qe={class:`text-sm font-medium`},$e=`/api/v1/tenant/job-management/hr-authorities`,et={__name:`JobHRAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:n}){let r=e,a=n,{t:o}=P(),s=A(),c=i([]),l=i(!1),u=i(0),p=i(!1),h=i(!1),b=i(``),x=i(!1),S=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({description:``}),M=f(()=>{let e=o(`job_management.hr_authorities`);return h.value?`${o(`common.edit`)} ${e}`:`${o(`common.create`)} ${e}`}),N=f(()=>[{field:`description`,header:o(`job_management.description`)}]);async function I(e,t){l.value=!0;try{let n=await w.get($e,{params:{page:e,per_page:t,organization_id:r.orgId}});c.value=n.data?.data||[],u.value=n.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function R(){h.value=!1,b.value=``,j.value={nomenclature:``,full_code:``,description:``},S.value={},p.value=!0}function B(e){h.value=!0,b.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},S.value={},p.value=!0}async function V(){x.value=!0,S.value={};try{let e={...j.value,nomenclature:r.orgName||``,full_code:r.orgCode||``,organization_id:r.orgId};h.value?await w.put(`${$e}/${b.value}`,e):await w.post($e,e),p.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),I(1,15)}catch(e){let t=F(e);Object.keys(t).length?S.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{x.value=!1}}function H(e){k.value=e,O.value=``,T.value=!0}async function U(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${$e}/${k.value.id}`),T.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),I(1,15)}catch(e){O.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{D.value=!1}}}return(n,r)=>(t(),v(`div`,qe,[g(`div`,Je,[g(`div`,null,[g(`h2`,Ye,y(m(o)(`job_management.hr_authorities`)),1),g(`p`,Xe,y(m(o)(`job_management.authority_description`)),1)]),d(m(E),{label:m(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:r[0]||=e=>R()},null,8,[`label`])]),d(Z,{items:c.value,loading:l.value,total:u.value,columns:N.value,entity:`hr-authorities`,"org-id":e.orgId,"on-load":I,onEdit:B,onDelete:H},{empty:C(()=>[g(`div`,Ze,[r[6]||=g(`i`,{class:`pi pi-users text-3xl mb-2 opacity-50`},null,-1),g(`p`,Qe,y(m(o)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),d(Q,{visible:p.value,"onUpdate:visible":r[2]||=e=>p.value=e,title:M.value,saving:x.value,errors:S.value,onSave:V,onCancel:r[3]||=e=>p.value=!1},{default:C(()=>[d(L,{label:m(o)(`job_management.description`),errors:S.value?.description},{default:C(()=>[d(m(z),{modelValue:j.value.description,"onUpdate:modelValue":r[1]||=e=>j.value.description=e,rows:`3`,class:_([`w-full`,{"p-invalid":S.value?.description}])},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),d(K,{visible:T.value,"onUpdate:visible":r[4]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:U,onCancel:r[5]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},tt={class:`space-y-4`},nt={class:`flex items-center justify-between`},rt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},it={class:`text-sm text-gray-500 dark:text-gray-400`},at={class:`flex flex-col items-center justify-center py-10 text-gray-400`},ot={class:`text-sm font-medium`},st=`/api/v1/tenant/job-management/operational-authorities`,ct={__name:`JobOpAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:n}){let r=e,a=n,{t:o}=P(),s=A(),c=i([]),l=i(!1),u=i(0),p=i(!1),h=i(!1),b=i(``),x=i(!1),S=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({description:``}),M=f(()=>{let e=o(`job_management.op_authorities`);return h.value?`${o(`common.edit`)} ${e}`:`${o(`common.create`)} ${e}`}),N=f(()=>[{field:`description`,header:o(`job_management.description`)}]);async function I(e,t){l.value=!0;try{let n=await w.get(st,{params:{page:e,per_page:t,organization_id:r.orgId}});c.value=n.data?.data||[],u.value=n.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function R(){h.value=!1,b.value=``,j.value={nomenclature:``,full_code:``,description:``},S.value={},p.value=!0}function B(e){h.value=!0,b.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},S.value={},p.value=!0}async function V(){x.value=!0,S.value={};try{let e={...j.value,nomenclature:r.orgName||``,full_code:r.orgCode||``,organization_id:r.orgId};h.value?await w.put(`${st}/${b.value}`,e):await w.post(st,e),p.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),I(1,15)}catch(e){let t=F(e);Object.keys(t).length?S.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{x.value=!1}}function H(e){k.value=e,O.value=``,T.value=!0}async function U(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${st}/${k.value.id}`),T.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),I(1,15)}catch(e){O.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{D.value=!1}}}return(n,r)=>(t(),v(`div`,tt,[g(`div`,nt,[g(`div`,null,[g(`h2`,rt,y(m(o)(`job_management.op_authorities`)),1),g(`p`,it,y(m(o)(`job_management.authority_description`)),1)]),d(m(E),{label:m(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:r[0]||=e=>R()},null,8,[`label`])]),d(Z,{items:c.value,loading:l.value,total:u.value,columns:N.value,entity:`operational-authorities`,"org-id":e.orgId,"on-load":I,onEdit:B,onDelete:H},{empty:C(()=>[g(`div`,at,[r[6]||=g(`i`,{class:`pi pi-cog text-3xl mb-2 opacity-50`},null,-1),g(`p`,ot,y(m(o)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),d(Q,{visible:p.value,"onUpdate:visible":r[2]||=e=>p.value=e,title:M.value,saving:x.value,errors:S.value,onSave:V,onCancel:r[3]||=e=>p.value=!1},{default:C(()=>[d(L,{label:m(o)(`job_management.description`),errors:S.value?.description},{default:C(()=>[d(m(z),{modelValue:j.value.description,"onUpdate:modelValue":r[1]||=e=>j.value.description=e,class:_([`w-full`,{"p-invalid":S.value?.description}]),rows:`3`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),d(K,{visible:T.value,"onUpdate:visible":r[4]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:U,onCancel:r[5]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},lt={class:`space-y-4`},ut={class:`flex items-center justify-between`},dt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},ft={class:`text-sm text-gray-500 dark:text-gray-400`},pt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},mt={class:`text-sm font-medium`},ht=`/api/v1/tenant/job-management/working-activities`,gt={__name:`JobActivitySection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:n}){let r=e,a=n,{t:o}=P(),s=A(),c=i([]),l=i(!1),u=i(0),p=i(!1),h=i(!1),b=i(``),x=i(!1),S=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({nomenclature:``,full_code:``,job_management_value_id:``});f(()=>Object.values(r.jobValueMap||{}).flat());let M=f(()=>[{field:`nomenclature`,header:o(`organization.nomenclature`)},{field:`full_code`,header:o(`organization.full_code`)}]);async function N(e,t){l.value=!0;try{let n=await w.get(ht,{params:{page:e,per_page:t,organization_id:r.orgId}});c.value=n.data?.data||[],u.value=n.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function I(){h.value=!1,b.value=``,j.value={nomenclature:``,full_code:``,job_management_value_id:``},S.value={},p.value=!0}function R(e){h.value=!0,b.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_id:e.job_management_value_id||``},S.value={},p.value=!0}async function z(){x.value=!0,S.value={};try{let e={...j.value,organization_id:r.orgId};h.value?await w.put(`${ht}/${b.value}`,e):await w.post(ht,e),p.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),N(1,15)}catch(e){let t=F(e);Object.keys(t).length?S.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{x.value=!1}}function V(e){k.value=e,O.value=``,T.value=!0}async function H(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${ht}/${k.value.id}`),T.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),N(1,15)}catch(e){O.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{D.value=!1}}}return(n,r)=>(t(),v(`div`,lt,[g(`div`,ut,[g(`div`,null,[g(`h2`,dt,y(m(o)(`job_management.activities`)),1),g(`p`,ft,y(m(o)(`job_management.activity_description`)),1)]),d(m(E),{label:m(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:r[0]||=e=>I()},null,8,[`label`])]),d(Z,{items:c.value,loading:l.value,total:u.value,columns:M.value,entity:`working-activities`,"org-id":e.orgId,"on-load":N,onEdit:R,onDelete:V},{empty:C(()=>[g(`div`,pt,[r[8]||=g(`i`,{class:`pi pi-bolt text-3xl mb-2 opacity-50`},null,-1),g(`p`,mt,y(m(o)(`job_management.empty_activities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),d(Q,{visible:p.value,"onUpdate:visible":r[4]||=e=>p.value=e,title:h.value?m(o)(`common.edit`):m(o)(`common.create`),saving:x.value,errors:S.value,onSave:z,onCancel:r[5]||=e=>p.value=!1},{default:C(()=>[d(L,{label:m(o)(`organization.nomenclature`),required:``,errors:S.value?.nomenclature},{default:C(()=>[d(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":r[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:_({"p-invalid":S.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`organization.full_code`),required:``,errors:S.value?.full_code},{default:C(()=>[d(B,{modelValue:j.value.full_code,"onUpdate:modelValue":r[2]||=e=>j.value.full_code=e,maxlength:`20`,class:_({"p-invalid":S.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.activity_type`),errors:S.value?.job_management_value_id},{default:C(()=>[d(G,{modelValue:j.value.job_management_value_id,"onUpdate:modelValue":r[3]||=e=>j.value.job_management_value_id=e,options:n.activityOptions,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),d(K,{visible:T.value,"onUpdate:visible":r[6]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:H,onCancel:r[7]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},_t={class:`space-y-4`},vt={class:`flex items-center justify-between`},yt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},bt={class:`text-sm text-gray-500 dark:text-gray-400`},xt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},St={class:`text-sm font-medium`},$=`/api/v1/tenant/job-management/working-risks`,Ct={__name:`JobRiskSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:n}){let r=e,a=n,{t:o}=P(),s=A(),c=i([]),l=i(!1),u=i(0),p=i(!1),h=i(!1),b=i(``),x=i(!1),S=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({nomenclature:``,full_code:``,job_management_value_environment_id:``,job_management_value_hazard_id:``}),M=f(()=>r.jobValueMap?.environment||[]),N=f(()=>r.jobValueMap?.hazard||[]),I=f(()=>[{field:`nomenclature`,header:o(`organization.nomenclature`)},{field:`full_code`,header:o(`organization.full_code`)}]);async function R(e,t){l.value=!0;try{let n=await w.get($,{params:{page:e,per_page:t,organization_id:r.orgId}});c.value=n.data?.data||[],u.value=n.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function z(){h.value=!1,b.value=``,j.value={nomenclature:``,full_code:``,job_management_value_environment_id:``,job_management_value_hazard_id:``},S.value={},p.value=!0}function V(e){h.value=!0,b.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_environment_id:e.job_management_value_environment_id||``,job_management_value_hazard_id:e.job_management_value_hazard_id||``},S.value={},p.value=!0}async function H(){x.value=!0,S.value={};try{let e={...j.value,organization_id:r.orgId};h.value?await w.put(`${$}/${b.value}`,e):await w.post($,e),p.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),R(1,15)}catch(e){let t=F(e);Object.keys(t).length?S.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{x.value=!1}}function U(e){k.value=e,O.value=``,T.value=!0}async function W(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${$}/${k.value.id}`),T.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),R(1,15)}catch(e){O.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{D.value=!1}}}return(n,r)=>(t(),v(`div`,_t,[g(`div`,vt,[g(`div`,null,[g(`h2`,yt,y(m(o)(`job_management.risks`)),1),g(`p`,bt,y(m(o)(`job_management.risk_description`)),1)]),d(m(E),{label:m(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:r[0]||=e=>z()},null,8,[`label`])]),d(Z,{items:c.value,loading:l.value,total:u.value,columns:I.value,entity:`working-risks`,"org-id":e.orgId,"on-load":R,onEdit:V,onDelete:U},{empty:C(()=>[g(`div`,xt,[r[9]||=g(`i`,{class:`pi pi-exclamation-triangle text-3xl mb-2 opacity-50`},null,-1),g(`p`,St,y(m(o)(`job_management.empty_risks`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),d(Q,{visible:p.value,"onUpdate:visible":r[5]||=e=>p.value=e,title:h.value?m(o)(`common.edit`):m(o)(`common.create`),saving:x.value,errors:S.value,onSave:H,onCancel:r[6]||=e=>p.value=!1},{default:C(()=>[d(L,{label:m(o)(`organization.nomenclature`),required:``,errors:S.value?.nomenclature},{default:C(()=>[d(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":r[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:_({"p-invalid":S.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`organization.full_code`),required:``,errors:S.value?.full_code},{default:C(()=>[d(B,{modelValue:j.value.full_code,"onUpdate:modelValue":r[2]||=e=>j.value.full_code=e,maxlength:`20`,class:_({"p-invalid":S.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.environment_risk`),errors:S.value?.job_management_value_environment_id},{default:C(()=>[d(G,{modelValue:j.value.job_management_value_environment_id,"onUpdate:modelValue":r[3]||=e=>j.value.job_management_value_environment_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.hazard_risk`),errors:S.value?.job_management_value_hazard_id},{default:C(()=>[d(G,{modelValue:j.value.job_management_value_hazard_id,"onUpdate:modelValue":r[4]||=e=>j.value.job_management_value_hazard_id=e,options:N.value,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),d(K,{visible:T.value,"onUpdate:visible":r[7]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:W,onCancel:r[8]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},wt={class:`space-y-4`},Tt={class:`flex items-center justify-between`},Et={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Dt={class:`text-sm text-gray-500 dark:text-gray-400`},Ot={class:`flex flex-col items-center justify-center py-10 text-gray-400`},kt={class:`text-sm font-medium`},At=`/api/v1/tenant/job-management/relationships`,jt={__name:`JobRelationshipSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:n}){let r=e,a=n,{t:o}=P(),s=A(),c=i([]),l=i(!1),u=i(0),p=i(!1),h=i(!1),b=i(``),x=i(!1),S=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({nomenclature:``,full_code:``,job_management_value_relationship_id:``,job_management_value_frequency_id:``}),M=f(()=>r.jobValueMap?.relationship||[]),N=f(()=>r.jobValueMap?.frequency||[]),I=f(()=>[{field:`nomenclature`,header:o(`organization.nomenclature`)},{field:`full_code`,header:o(`organization.full_code`)}]);async function R(e,t){l.value=!0;try{let n=await w.get(At,{params:{page:e,per_page:t,organization_id:r.orgId}});c.value=n.data?.data||[],u.value=n.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function z(){h.value=!1,b.value=``,j.value={nomenclature:``,full_code:``,job_management_value_relationship_id:``,job_management_value_frequency_id:``},S.value={},p.value=!0}function V(e){h.value=!0,b.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_relationship_id:e.job_management_value_relationship_id||``,job_management_value_frequency_id:e.job_management_value_frequency_id||``},S.value={},p.value=!0}async function H(){x.value=!0,S.value={};try{let e={...j.value,organization_id:r.orgId};h.value?await w.put(`${At}/${b.value}`,e):await w.post(At,e),p.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),R(1,15)}catch(e){let t=F(e);Object.keys(t).length?S.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{x.value=!1}}function U(e){k.value=e,O.value=``,T.value=!0}async function W(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${At}/${k.value.id}`),T.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),R(1,15)}catch(e){O.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{D.value=!1}}}return(n,r)=>(t(),v(`div`,wt,[g(`div`,Tt,[g(`div`,null,[g(`h2`,Et,y(m(o)(`job_management.relationships`)),1),g(`p`,Dt,y(m(o)(`job_management.relationship_description`)),1)]),d(m(E),{label:m(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:r[0]||=e=>z()},null,8,[`label`])]),d(Z,{items:c.value,loading:l.value,total:u.value,columns:I.value,entity:`relationships`,"org-id":e.orgId,"on-load":R,onEdit:V,onDelete:U},{empty:C(()=>[g(`div`,Ot,[r[9]||=g(`i`,{class:`pi pi-share-alt text-3xl mb-2 opacity-50`},null,-1),g(`p`,kt,y(m(o)(`job_management.empty_relationships`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),d(Q,{visible:p.value,"onUpdate:visible":r[5]||=e=>p.value=e,title:h.value?m(o)(`common.edit`):m(o)(`common.create`),saving:x.value,errors:S.value,onSave:H,onCancel:r[6]||=e=>p.value=!1},{default:C(()=>[d(L,{label:m(o)(`organization.nomenclature`),required:``,errors:S.value?.nomenclature},{default:C(()=>[d(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":r[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:_({"p-invalid":S.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`organization.full_code`),required:``,errors:S.value?.full_code},{default:C(()=>[d(B,{modelValue:j.value.full_code,"onUpdate:modelValue":r[2]||=e=>j.value.full_code=e,maxlength:`20`,class:_({"p-invalid":S.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.relationship_type`),errors:S.value?.job_management_value_relationship_id},{default:C(()=>[d(G,{modelValue:j.value.job_management_value_relationship_id,"onUpdate:modelValue":r[3]||=e=>j.value.job_management_value_relationship_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.frequency`),errors:S.value?.job_management_value_frequency_id},{default:C(()=>[d(G,{modelValue:j.value.job_management_value_frequency_id,"onUpdate:modelValue":r[4]||=e=>j.value.job_management_value_frequency_id=e,options:N.value,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),d(K,{visible:T.value,"onUpdate:visible":r[7]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:W,onCancel:r[8]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Mt={class:`space-y-4`},Nt={class:`flex items-center justify-between`},Pt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Ft={class:`text-sm text-gray-500 dark:text-gray-400`},It={class:`flex flex-col items-center justify-center py-10 text-gray-400`},Lt={class:`text-sm font-medium`},Rt=`/api/v1/tenant/job-management/subordinate-controls`,zt={__name:`JobSubordinateSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:n}){let r=e,a=n,{t:o}=P(),s=A(),c=i([]),l=i(!1),u=i(0),p=i(!1),h=i(!1),b=i(``),x=i(!1),S=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({nomenclature:``,full_code:``,job_management_value_id:``}),M=f(()=>Object.values(r.jobValueMap||{}).flat()),N=f(()=>[{field:`nomenclature`,header:o(`organization.nomenclature`)},{field:`full_code`,header:o(`organization.full_code`)}]);async function I(e,t){l.value=!0;try{let n=await w.get(Rt,{params:{page:e,per_page:t,organization_id:r.orgId}});c.value=n.data?.data||[],u.value=n.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function R(){h.value=!1,b.value=``,j.value={nomenclature:``,full_code:``,job_management_value_id:``},S.value={},p.value=!0}function z(e){h.value=!0,b.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_id:e.job_management_value_id||``},S.value={},p.value=!0}async function V(){x.value=!0,S.value={};try{let e={...j.value,organization_id:r.orgId};h.value?await w.put(`${Rt}/${b.value}`,e):await w.post(Rt,e),p.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),I(1,15)}catch(e){let t=F(e);Object.keys(t).length?S.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{x.value=!1}}function H(e){k.value=e,O.value=``,T.value=!0}async function U(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${Rt}/${k.value.id}`),T.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),I(1,15)}catch(e){O.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{D.value=!1}}}return(n,r)=>(t(),v(`div`,Mt,[g(`div`,Nt,[g(`div`,null,[g(`h2`,Pt,y(m(o)(`job_management.subordinate_controls`)),1),g(`p`,Ft,y(m(o)(`job_management.subordinate_description`)),1)]),d(m(E),{label:m(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:r[0]||=e=>R()},null,8,[`label`])]),d(Z,{items:c.value,loading:l.value,total:u.value,columns:N.value,entity:`subordinate-controls`,"org-id":e.orgId,"on-load":I,onEdit:z,onDelete:H},{empty:C(()=>[g(`div`,It,[r[8]||=g(`i`,{class:`pi pi-sitemap text-3xl mb-2 opacity-50`},null,-1),g(`p`,Lt,y(m(o)(`job_management.empty_subordinates`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),d(Q,{visible:p.value,"onUpdate:visible":r[4]||=e=>p.value=e,title:h.value?m(o)(`common.edit`):m(o)(`common.create`),saving:x.value,errors:S.value,onSave:V,onCancel:r[5]||=e=>p.value=!1},{default:C(()=>[d(L,{label:m(o)(`organization.nomenclature`),required:``,errors:S.value?.nomenclature},{default:C(()=>[d(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":r[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:_({"p-invalid":S.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`organization.full_code`),required:``,errors:S.value?.full_code},{default:C(()=>[d(B,{modelValue:j.value.full_code,"onUpdate:modelValue":r[2]||=e=>j.value.full_code=e,maxlength:`20`,class:_({"p-invalid":S.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.control_type`),errors:S.value?.job_management_value_id},{default:C(()=>[d(G,{modelValue:j.value.job_management_value_id,"onUpdate:modelValue":r[3]||=e=>j.value.job_management_value_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),d(K,{visible:T.value,"onUpdate:visible":r[6]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:U,onCancel:r[7]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Bt={class:`space-y-4`},Vt={class:`flex items-center justify-between`},Ht={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Ut={class:`text-sm text-gray-500 dark:text-gray-400`},Wt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},Gt={class:`text-sm font-medium`},Kt=`/api/v1/tenant/job-management/assets`,qt={__name:`JobAssetSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:n}){let r=e,a=n,{t:o}=P(),s=A(),c=i([]),l=i(!1),u=i(0),p=i(!1),h=i(!1),b=i(``),x=i(!1),S=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({nomenclature:``,full_code:``,job_management_value_asset_id:``,job_management_value_authority_id:``}),M=f(()=>r.jobValueMap?.asset||[]),N=f(()=>r.jobValueMap?.authority||[]),I=f(()=>[{field:`nomenclature`,header:o(`organization.nomenclature`)},{field:`full_code`,header:o(`organization.full_code`)}]);async function R(e,t){l.value=!0;try{let n=await w.get(Kt,{params:{page:e,per_page:t,organization_id:r.orgId}});c.value=n.data?.data||[],u.value=n.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function z(){h.value=!1,b.value=``,j.value={nomenclature:``,full_code:``,job_management_value_asset_id:``,job_management_value_authority_id:``},S.value={},p.value=!0}function V(e){h.value=!0,b.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,job_management_value_asset_id:e.job_management_value_asset_id||``,job_management_value_authority_id:e.job_management_value_authority_id||``},S.value={},p.value=!0}async function H(){x.value=!0,S.value={};try{let e={...j.value,organization_id:r.orgId};h.value?await w.put(`${Kt}/${b.value}`,e):await w.post(Kt,e),p.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),R(1,15)}catch(e){let t=F(e);Object.keys(t).length?S.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{x.value=!1}}function U(e){k.value=e,O.value=``,T.value=!0}async function W(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${Kt}/${k.value.id}`),T.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),R(1,15)}catch(e){O.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{D.value=!1}}}return(n,r)=>(t(),v(`div`,Bt,[g(`div`,Vt,[g(`div`,null,[g(`h2`,Ht,y(m(o)(`job_management.assets`)),1),g(`p`,Ut,y(m(o)(`job_management.asset_description`)),1)]),d(m(E),{label:m(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:r[0]||=e=>z()},null,8,[`label`])]),d(Z,{items:c.value,loading:l.value,total:u.value,columns:I.value,entity:`assets`,"org-id":e.orgId,"on-load":R,onEdit:V,onDelete:U},{empty:C(()=>[g(`div`,Wt,[r[9]||=g(`i`,{class:`pi pi-box text-3xl mb-2 opacity-50`},null,-1),g(`p`,Gt,y(m(o)(`job_management.empty_assets`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),d(Q,{visible:p.value,"onUpdate:visible":r[5]||=e=>p.value=e,title:h.value?m(o)(`common.edit`):m(o)(`common.create`),saving:x.value,errors:S.value,onSave:H,onCancel:r[6]||=e=>p.value=!1},{default:C(()=>[d(L,{label:m(o)(`organization.nomenclature`),required:``,errors:S.value?.nomenclature},{default:C(()=>[d(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":r[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:_({"p-invalid":S.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`organization.full_code`),required:``,errors:S.value?.full_code},{default:C(()=>[d(B,{modelValue:j.value.full_code,"onUpdate:modelValue":r[2]||=e=>j.value.full_code=e,maxlength:`20`,class:_({"p-invalid":S.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.asset_type`),errors:S.value?.job_management_value_asset_id},{default:C(()=>[d(G,{modelValue:j.value.job_management_value_asset_id,"onUpdate:modelValue":r[3]||=e=>j.value.job_management_value_asset_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.authority_level`),errors:S.value?.job_management_value_authority_id},{default:C(()=>[d(G,{modelValue:j.value.job_management_value_authority_id,"onUpdate:modelValue":r[4]||=e=>j.value.job_management_value_authority_id=e,options:N.value,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),d(K,{visible:T.value,"onUpdate:visible":r[7]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:W,onCancel:r[8]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Jt={class:`space-y-4`},Yt={class:`flex items-center justify-between`},Xt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Zt={class:`text-sm text-gray-500 dark:text-gray-400`},Qt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},$t={class:`text-sm font-medium`},en=`/api/v1/tenant/job-management/financials`,tn={__name:`JobFinancialSection`,props:{orgId:String,jobValueMap:Object},emits:[`saved`],setup(e,{emit:n}){let r=e,a=n,{t:o}=P(),s=A(),c=i([]),l=i(!1),u=i(0),p=i(!1),h=i(!1),b=i(``),x=i(!1),S=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({nomenclature:``,full_code:``,is_authorized:!1,job_management_value_cash_id:``,job_management_value_authority_id:``,job_management_value_impact_id:``}),M=f(()=>r.jobValueMap?.cash||[]),N=f(()=>r.jobValueMap?.authority||[]),I=f(()=>r.jobValueMap?.impact||[]),R=f(()=>[{field:`nomenclature`,header:o(`organization.nomenclature`)},{field:`full_code`,header:o(`organization.full_code`)},{field:`is_authorized`,header:o(`job_management.is_authorized`)}]);async function z(e,t){l.value=!0;try{let n=await w.get(en,{params:{page:e,per_page:t,organization_id:r.orgId}});c.value=n.data?.data||[],u.value=n.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function V(){h.value=!1,b.value=``,j.value={nomenclature:``,full_code:``,is_authorized:!1,job_management_value_cash_id:``,job_management_value_authority_id:``,job_management_value_impact_id:``},S.value={},p.value=!0}function H(e){h.value=!0,b.value=e.id,j.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,is_authorized:!!e.is_authorized,job_management_value_cash_id:e.job_management_value_cash_id||``,job_management_value_authority_id:e.job_management_value_authority_id||``,job_management_value_impact_id:e.job_management_value_impact_id||``},S.value={},p.value=!0}async function U(){x.value=!0,S.value={};try{let e={...j.value,organization_id:r.orgId};h.value?await w.put(`${en}/${b.value}`,e):await w.post(en,e),p.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),z(1,15)}catch(e){let t=F(e);Object.keys(t).length?S.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{x.value=!1}}function W(e){k.value=e,O.value=``,T.value=!0}async function ee(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${en}/${k.value.id}`),T.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),z(1,15)}catch(e){O.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{D.value=!1}}}return(n,r)=>(t(),v(`div`,Jt,[g(`div`,Yt,[g(`div`,null,[g(`h2`,Xt,y(m(o)(`job_management.financials`)),1),g(`p`,Zt,y(m(o)(`job_management.financial_description`)),1)]),d(m(E),{label:m(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:r[0]||=e=>V()},null,8,[`label`])]),d(Z,{items:c.value,loading:l.value,total:u.value,columns:R.value,entity:`financials`,"org-id":e.orgId,"on-load":z,onEdit:H,onDelete:W},{empty:C(()=>[g(`div`,Qt,[r[11]||=g(`i`,{class:`pi pi-money-bill text-3xl mb-2 opacity-50`},null,-1),g(`p`,$t,y(m(o)(`job_management.empty_financials`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),d(Q,{visible:p.value,"onUpdate:visible":r[7]||=e=>p.value=e,title:h.value?m(o)(`common.edit`):m(o)(`common.create`),saving:x.value,errors:S.value,onSave:U,onCancel:r[8]||=e=>p.value=!1},{default:C(()=>[d(L,{label:m(o)(`organization.nomenclature`),required:``,errors:S.value?.nomenclature},{default:C(()=>[d(B,{modelValue:j.value.nomenclature,"onUpdate:modelValue":r[1]||=e=>j.value.nomenclature=e,maxlength:`50`,class:_({"p-invalid":S.value?.nomenclature})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`organization.full_code`),required:``,errors:S.value?.full_code},{default:C(()=>[d(B,{modelValue:j.value.full_code,"onUpdate:modelValue":r[2]||=e=>j.value.full_code=e,maxlength:`20`,class:_({"p-invalid":S.value?.full_code})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.is_authorized`),class:`md:col-span-2`},{default:C(()=>[d(m(te),{modelValue:j.value.is_authorized,"onUpdate:modelValue":r[3]||=e=>j.value.is_authorized=e},null,8,[`modelValue`])]),_:1},8,[`label`]),d(L,{label:m(o)(`job_management.cash_level`),errors:S.value?.job_management_value_cash_id},{default:C(()=>[d(G,{modelValue:j.value.job_management_value_cash_id,"onUpdate:modelValue":r[4]||=e=>j.value.job_management_value_cash_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.authority_level`),errors:S.value?.job_management_value_authority_id},{default:C(()=>[d(G,{modelValue:j.value.job_management_value_authority_id,"onUpdate:modelValue":r[5]||=e=>j.value.job_management_value_authority_id=e,options:N.value,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.impact_level`),errors:S.value?.job_management_value_impact_id},{default:C(()=>[d(G,{modelValue:j.value.job_management_value_impact_id,"onUpdate:modelValue":r[6]||=e=>j.value.job_management_value_impact_id=e,options:I.value,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),d(K,{visible:T.value,"onUpdate:visible":r[9]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:ee,onCancel:r[10]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},nn={class:`space-y-4`},rn={class:`flex items-center justify-between`},an={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},on={class:`text-sm text-gray-500 dark:text-gray-400`},sn={class:`flex flex-col items-center justify-center py-10 text-gray-400`},cn={class:`text-sm font-medium`},ln=`/api/v1/tenant/job-management/potency-competencies`,un={__name:`JobPotencySection`,props:{orgId:String,jobValueMap:Object,competencyOptions:Array},emits:[`saved`],setup(e,{emit:n}){let r=e,a=n,{t:o}=P(),s=A(),c=i([]),l=i(!1),u=i(0),p=i(!1),h=i(!1),b=i(``),x=i(!1),S=i({}),T=i(!1),D=i(!1),O=i(``),k=i(null),j=i({competency_id:``,job_management_value_id:``,weight:null}),M=f(()=>Object.values(r.jobValueMap||{}).flat()),N=f(()=>[{field:`_competency`,header:o(`job_management.competency`)},{field:`weight`,header:o(`job_management.weight`)}]);async function I(e,t){l.value=!0;try{let n=await w.get(ln,{params:{page:e,per_page:t,organization_id:r.orgId}});c.value=(n.data?.data||[]).map(e=>{let t=r.competencyOptions?.find(t=>t.value===e.competency_id);return{...e,_competency:t?.label||e.competency_id}}),u.value=n.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function R(){h.value=!1,b.value=``,j.value={competency_id:``,job_management_value_id:``,weight:null},S.value={},p.value=!0}function z(e){h.value=!0,b.value=e.id,j.value={competency_id:e.competency_id||``,job_management_value_id:e.job_management_value_id||``,weight:e.weight??null},S.value={},p.value=!0}async function B(){x.value=!0,S.value={};try{let e={...j.value,organization_id:r.orgId};h.value?await w.put(`${ln}/${b.value}`,e):await w.post(ln,e),p.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),I(1,15)}catch(e){let t=F(e);Object.keys(t).length?S.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{x.value=!1}}function H(e){k.value=e,O.value=``,T.value=!0}async function U(){if(k.value){D.value=!0,O.value=``;try{await w.delete(`${ln}/${k.value.id}`),T.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),I(1,15)}catch(e){O.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{D.value=!1}}}return(n,r)=>(t(),v(`div`,nn,[g(`div`,rn,[g(`div`,null,[g(`h2`,an,y(m(o)(`job_management.potency_competencies`)),1),g(`p`,on,y(m(o)(`job_management.potency_description`)),1)]),d(m(E),{label:m(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:r[0]||=e=>R()},null,8,[`label`])]),d(Z,{items:c.value,loading:l.value,total:u.value,columns:N.value,entity:`potency-competencies`,"org-id":e.orgId,"on-load":I,onEdit:z,onDelete:H},{empty:C(()=>[g(`div`,sn,[r[8]||=g(`i`,{class:`pi pi-star text-3xl mb-2 opacity-50`},null,-1),g(`p`,cn,y(m(o)(`job_management.empty_potency`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),d(Q,{visible:p.value,"onUpdate:visible":r[4]||=e=>p.value=e,title:h.value?m(o)(`common.edit`):m(o)(`common.create`),saving:x.value,errors:S.value,onSave:B,onCancel:r[5]||=e=>p.value=!1},{default:C(()=>[d(L,{label:m(o)(`job_management.competency`),required:``,errors:S.value?.competency_id},{default:C(()=>[d(G,{modelValue:j.value.competency_id,"onUpdate:modelValue":r[1]||=e=>j.value.competency_id=e,options:e.competencyOptions,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.value_ref`),errors:S.value?.job_management_value_id},{default:C(()=>[d(G,{modelValue:j.value.job_management_value_id,"onUpdate:modelValue":r[2]||=e=>j.value.job_management_value_id=e,options:M.value,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),showClear:``},null,8,[`modelValue`,`options`,`placeholder`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.weight`),errors:S.value?.weight},{default:C(()=>[d(m(V),{modelValue:j.value.weight,"onUpdate:modelValue":r[3]||=e=>j.value.weight=e,min:0,max:100,class:_([{"p-invalid":S.value?.weight},`w-full`]),size:`small`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),d(K,{visible:T.value,"onUpdate:visible":r[6]||=e=>T.value=e,loading:D.value,"error-msg":O.value,onConfirm:U,onCancel:r[7]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},dn={class:`space-y-4`},fn={class:`flex items-center justify-between`},pn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},mn={class:`text-sm text-gray-500 dark:text-gray-400`},hn={class:`flex flex-col items-center justify-center py-10 text-gray-400`},gn={class:`text-sm font-medium`},_n=`/api/v1/tenant/job-management/competency-groups`,vn={__name:`JobCompetencyGroupSection`,props:{orgId:String},emits:[`saved`],setup(e,{emit:n}){let r=e,a=n,{t:o}=P(),s=A(),c=i([]),l=i(!1);i(0);let u=i(!1),p=i(!1),h=i(``),b=i(!1),x=i({}),S=i(!1),T=i(!1),D=i(``),O=i(null),k=i({category:``,weight:null}),j=f(()=>[{label:`${o(`job_management.technical`)} (${o(`job_management.category`)})`,value:`technical`},{label:`${o(`job_management.managerial`)} (${o(`job_management.category`)})`,value:`managerial`}]),M=f(()=>[{field:`category`,header:o(`job_management.category`)},{field:`weight`,header:o(`job_management.weight`)}]);async function N(){l.value=!0;try{let e=await w.get(_n,{params:{organization_id:r.orgId}});c.value=e.data?.data||(Array.isArray(e.data)?e.data:[])}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function I(){p.value=!1,h.value=``,k.value={category:`technical`,weight:null},x.value={},u.value=!0}function R(e){p.value=!0,h.value=e.id,k.value={category:e.category||`technical`,weight:e.weight??null},x.value={},u.value=!0}async function z(){b.value=!0,x.value={};try{let e={...k.value,organization_id:r.orgId};p.value?await w.put(`${_n}/${h.value}`,e):await w.post(_n,e),u.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),N()}catch(e){let t=F(e);Object.keys(t).length?x.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{b.value=!1}}function B(e){O.value=e,D.value=``,S.value=!0}async function H(){if(O.value){T.value=!0,D.value=``;try{await w.delete(`${_n}/${O.value.id}`),S.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),N()}catch(e){D.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{T.value=!1}}}return(n,r)=>(t(),v(`div`,dn,[g(`div`,fn,[g(`div`,null,[g(`h2`,pn,y(m(o)(`job_management.competency_groups`)),1),g(`p`,mn,y(m(o)(`job_management.competency_group_description`)),1)]),d(m(E),{label:m(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:r[0]||=e=>I()},null,8,[`label`])]),d(Z,{items:c.value,loading:l.value,total:c.value.length,columns:M.value,entity:`competency-groups`,"org-id":e.orgId,"on-load":N,onEdit:R,onDelete:B},{empty:C(()=>[g(`div`,hn,[r[7]||=g(`i`,{class:`pi pi-chart-pie text-3xl mb-2 opacity-50`},null,-1),g(`p`,gn,y(m(o)(`job_management.empty_competency_groups`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),d(Q,{visible:u.value,"onUpdate:visible":r[3]||=e=>u.value=e,title:p.value?m(o)(`common.edit`):m(o)(`common.create`),saving:b.value,errors:x.value,onSave:z,onCancel:r[4]||=e=>u.value=!1},{default:C(()=>[d(L,{label:m(o)(`job_management.category`),required:``,errors:x.value?.category},{default:C(()=>[d(G,{modelValue:k.value.category,"onUpdate:modelValue":r[1]||=e=>k.value.category=e,options:j.value,optionLabel:`label`,optionValue:`value`,placeholder:m(o)(`common.select`),class:_({"p-invalid":x.value?.category})},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),d(L,{label:m(o)(`job_management.weight`),required:``,errors:x.value?.weight},{default:C(()=>[d(m(V),{modelValue:k.value.weight,"onUpdate:modelValue":r[2]||=e=>k.value.weight=e,min:0,max:100,suffix:`%`,class:_([{"p-invalid":x.value?.weight},`w-full`]),size:`small`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),d(K,{visible:S.value,"onUpdate:visible":r[5]||=e=>S.value=e,loading:T.value,"error-msg":D.value,onConfirm:H,onCancel:r[6]||=e=>S.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},yn={class:`space-y-6`},bn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},xn={class:`text-sm text-gray-500 dark:text-gray-400`},Sn={key:0,class:`flex items-center justify-center py-12`},Cn={class:`grid grid-cols-1 md:grid-cols-3 gap-4`},wn={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Tn={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},En={class:`text-2xl font-bold text-emerald-600 dark:text-emerald-400`},Dn={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},On={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},kn={class:`text-2xl font-bold text-blue-600 dark:text-blue-400`},An={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},jn={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},Mn={key:0,class:`text-[10px] text-gray-400 mt-2`},Nn={key:0,class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden`},Pn={class:`px-5 py-3 border-b border-gray-200 dark:border-gray-700 font-semibold text-sm text-gray-700 dark:text-gray-300`},Fn={class:`p-5`},In={class:`text-sm text-gray-700 dark:text-gray-300 capitalize`},Ln={class:`text-sm font-semibold text-gray-900 dark:text-gray-100`},Rn={key:2},zn={class:`flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500`},Bn={class:`text-sm font-medium`},Vn={class:`text-xs mt-1`},Hn={class:`flex justify-end gap-3`},Un=`/api/v1/tenant/job-management/scores`,Wn={__name:`JobScoreSection`,props:{orgId:String},emits:[`saved`],setup(e,{emit:n}){let o=e,s=n,{t:c}=P(),u=A(),p=i(!1),_=i(!1),b=i(null),S=f(()=>{if(!b.value?.components)return null;try{return JSON.parse(b.value.components)}catch{return null}});function C(e){return e?.toLocaleString?.(`id-ID`)??`-`}async function T(){if(o.orgId){p.value=!0;try{let e=await w.get(`${Un}/${o.orgId}`);b.value=e.data?.data||null,s(`saved`)}catch{b.value=null}finally{p.value=!1}}}async function D(){if(o.orgId){_.value=!0;try{let e=await w.put(`${Un}/${o.orgId}`,{components:null});b.value=e.data?.data||null,u.add({severity:`success`,summary:c(`message.success`),detail:c(`job_management.score_calculated`),life:2e3})}catch(e){u.add({severity:`error`,summary:c(`message.error`),detail:e.response?.data?.error?.message||c(`message.operation_failed`),life:4e3})}finally{_.value=!1}}}return x(T),(e,n)=>(t(),v(`div`,yn,[g(`div`,null,[g(`h2`,bn,y(m(c)(`job_management.scores`)),1),g(`p`,xn,y(m(c)(`job_management.score_description`)),1)]),p.value?(t(),v(`div`,Sn,[...n[0]||=[g(`i`,{class:`pi pi-spin pi-spinner text-emerald-500 text-2xl`},null,-1)]])):b.value?(t(),v(h,{key:1},[g(`div`,Cn,[g(`div`,wn,[g(`div`,Tn,y(m(c)(`job_management.value_with_financial`)),1),g(`div`,En,y(C(b.value.job_value_with_financial)),1)]),g(`div`,Dn,[g(`div`,On,y(m(c)(`job_management.value_without_financial`)),1),g(`div`,kn,y(C(b.value.job_value_without_financial)),1)]),g(`div`,An,[g(`div`,jn,y(m(c)(`job_management.has_financial_authority`)),1),d(m(I),{value:b.value.has_financial_authority?m(c)(`common.yes`):m(c)(`common.no`),severity:b.value.has_financial_authority?`success`:`danger`,class:`!text-xs`},null,8,[`value`,`severity`]),b.value.calculated_at?(t(),v(`div`,Mn,y(m(c)(`job_management.calculated_at`))+`: `+y(b.value.calculated_at),1)):l(``,!0)])]),S.value?(t(),v(`div`,Nn,[g(`div`,Pn,y(m(c)(`job_management.component_breakdown`)),1),g(`div`,Fn,[(t(!0),v(h,null,r(S.value,(e,n)=>(t(),v(`div`,{key:n,class:`flex items-center justify-between py-1.5 border-b border-gray-100 dark:border-gray-700 last:border-0`},[g(`span`,In,y(n.replace(/_/g,` `)),1),g(`span`,Ln,y(C(e)),1)]))),128))])])):l(``,!0)],64)):(t(),v(`div`,Rn,[g(`div`,zn,[n[1]||=g(`i`,{class:`pi pi-calculator text-4xl mb-3 opacity-50`},null,-1),g(`p`,Bn,y(m(c)(`job_management.no_score`)),1),g(`p`,Vn,y(m(c)(`job_management.score_hint`)),1)])])),g(`div`,Hn,[d(m(E),{label:m(c)(`common.refresh`),icon:`pi pi-refresh`,size:`small`,text:``,onClick:T},null,8,[`label`]),b.value?(t(),a(m(E),{key:0,label:m(c)(`job_management.recalculate`),icon:`pi pi-calculator`,size:`small`,severity:`info`,loading:_.value,onClick:D},null,8,[`label`,`loading`])):l(``,!0)])]))}},Gn={class:`max-w-full mx-auto`},Kn={key:0,class:`flex gap-6`},qn={class:`w-56 space-y-2`},Jn={class:`flex-1 space-y-3`},Yn={key:1,class:`flex gap-6`},Xn={class:`w-56 shrink-0 space-y-1`},Zn=[`onClick`,`onKeydown`],Qn={key:0,class:`pi pi-check text-xs`},$n={class:`flex-1 min-w-0`},er={key:0,class:`pi pi-check-circle text-emerald-400 text-xs shrink-0`},tr={class:`flex-1 min-w-0`},nr={__name:`JobManagementForm`,setup(e){let o=N(),s=k(),{t:c}=P(),u=A(),d=s.query.org_id||``,p=i(0),b=i(!0),S=i(Array(15).fill(!1)),C=i(``),E=i(``),D=i(``),O=i(``),j=i([]),M=i([]),F=i([]),I=i({}),L=i([]),R=[{labelKey:`job_management.identifications`,icon:`pi pi-id-card`,comp:le},{labelKey:`job_management.objectives`,icon:`pi pi-bullseye`,comp:ge},{labelKey:`job_management.responsibilities_title`,icon:`pi pi-list-check`,comp:Ke},{labelKey:`job_management.education_experience`,icon:`pi pi-graduation-cap`,comp:Ce},{labelKey:`job_management.hr_authorities`,icon:`pi pi-users`,comp:et},{labelKey:`job_management.op_authorities`,icon:`pi pi-cog`,comp:ct},{labelKey:`job_management.activities`,icon:`pi pi-bolt`,comp:gt},{labelKey:`job_management.risks`,icon:`pi pi-exclamation-triangle`,comp:Ct},{labelKey:`job_management.relationships`,icon:`pi pi-share-alt`,comp:jt},{labelKey:`job_management.subordinate_controls`,icon:`pi pi-sitemap`,comp:zt},{labelKey:`job_management.assets`,icon:`pi pi-box`,comp:qt},{labelKey:`job_management.financials`,icon:`pi pi-money-bill`,comp:tn},{labelKey:`job_management.potency_competencies`,icon:`pi pi-star`,comp:un},{labelKey:`job_management.competency_groups`,icon:`pi pi-chart-pie`,comp:vn},{labelKey:`job_management.scores`,icon:`pi pi-calculator`,comp:Wn}],z=f(()=>R[p.value]?.comp||null);function B(e){return p.value===e?`bg-emerald-50 dark:bg-emerald-900/20 ring-1 ring-emerald-300 dark:ring-emerald-700`:(S.value[e],`hover:bg-gray-50 dark:hover:bg-gray-800`)}function V(e){return p.value===e?`bg-emerald-600 text-white`:S.value[e]?`bg-emerald-100 dark:bg-emerald-800 text-emerald-600 dark:text-emerald-300`:`bg-gray-200 dark:bg-gray-600 text-gray-600 dark:text-gray-300`}function H(e){return p.value===e?`text-emerald-700 dark:text-emerald-300`:S.value[e]?`text-emerald-600 dark:text-emerald-400`:`text-gray-700 dark:text-gray-300`}function U(e){p.value=e,o.replace({query:{...s.query,section:String(e)}})}function W(e){typeof e==`number`&&(S.value[e]=!0)}async function G(){if(d)try{let e=(await w.get(`/api/v1/tenant/organizations/${d}`)).data?.data;e&&(C.value=e.nomenclature||``,E.value=e.full_code||e.code||``,D.value=e.grading_id||``,O.value=e.job_family_id||``)}catch{}}async function K(){try{let[e,t,n,r]=await Promise.all([w.get(`/api/v1/tenant/settings/gradings?per_page=100`),w.get(`/api/v1/tenant/job-management/values?per_page=200`),w.get(`/api/v1/tenant/competencies?per_page=200`).catch(()=>({data:{data:[]}})),w.get(`/api/v1/tenant/settings/job-families?per_page=100`)]);j.value=(e.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id})),M.value=(r.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id}));let i=t.data?.data||[];F.value=i.map(e=>({label:`${e.type}${e.level?` Lv.`+e.level:``}${e.descriptions?` — `+e.descriptions:``}`,value:e.id,type:e.type,level:e.level,descriptions:e.descriptions}));let a={};i.forEach(e=>{a[e.type]||(a[e.type]=[]),a[e.type].push({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level})}),I.value=a,L.value=(n.data?.data||[]).map(e=>({label:e.name||e.code,value:e.id}))}catch{}}return x(async()=>{try{await Promise.all([G(),K()]);let e=parseInt(s.query.section);!isNaN(e)&&e>=0&&e<R.length&&(p.value=e)}catch(e){u.add({severity:`error`,summary:c(`message.error`),detail:e.response?.data?.error?.message||c(`message.failed_to_load`),life:4e3})}finally{b.value=!1}}),(e,i)=>(t(),v(`div`,Gn,[b.value?(t(),v(`div`,Kn,[g(`div`,qn,[(t(),v(h,null,r(8,e=>g(`div`,{key:e,class:`h-12 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))]),g(`div`,Jn,[(t(),v(h,null,r(6,e=>g(`div`,{key:e,class:`h-9 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))])])):(t(),v(`div`,Yn,[g(`div`,Xn,[(t(),v(h,null,r(R,(e,n)=>g(`div`,{key:n,role:`button`,tabindex:0,class:_([`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150 cursor-pointer select-none`,B(n)]),onClick:e=>U(n),onKeydown:T(e=>U(n),[`enter`])},[g(`div`,{class:_([`w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0 transition-colors duration-150`,V(n)])},[S.value[n]?(t(),v(`i`,Qn)):(t(),v(`i`,{key:1,class:_(e.icon)},null,2))],2),g(`div`,$n,[g(`div`,{class:_([`text-sm font-medium truncate`,H(n)])},y(m(c)(e.labelKey)),3)]),S.value[n]?(t(),v(`i`,er)):l(``,!0)],42,Zn)),64))]),g(`div`,tr,[(t(),a(n(z.value),{key:p.value,"org-id":m(d),"org-name":C.value,"org-code":E.value,"org-grading-id":D.value,"org-job-family-id":O.value,"grading-options":j.value,"job-family-options":M.value,"job-value-options":F.value,"competency-options":L.value,"job-value-map":I.value,onSaved:W},null,40,[`org-id`,`org-name`,`org-code`,`org-grading-id`,`org-job-family-id`,`grading-options`,`job-family-options`,`job-value-options`,`competency-options`,`job-value-map`]))])]))]))}};export{nr as default};