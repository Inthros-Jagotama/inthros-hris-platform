const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["assets/quill-BLmY9xB4.js","assets/rolldown-runtime-QTnfLwEv.js"])))=>i.map(i=>d[i]);
import{$ as e,A as t,D as n,F as r,G as i,I as a,K as o,M as s,N as c,U as l,W as u,Y as d,b as f,c as p,d as m,g as h,gt as g,h as _,ht as v,i as y,l as b,mt as x,pt as S,s as C,u as w,w as T}from"./language-8RKtU9ID.js";import{D as E,r as D}from"./basecomponent-Dq2wVn8v.js";import{O,_ as k,g as A,p as j,r as M,t as N,u as P,x as F}from"./index-C_5oItlt.js";import{t as I}from"./useI18n-CeYrJQoe.js";import{r as L}from"./responseHandler-BJxA-JZj.js";import{t as R}from"./tag-BToTW5OD.js";import{t as z}from"./FormRow-WmCnbL1x.js";import{t as B}from"./baseeditableholder-DSjqJDX7.js";import{t as V}from"./textarea-mG-wfBcE.js";import{t as H}from"./TextInput-BKI4yAUI.js";import{n as U,t as W}from"./column-LdUBebhW.js";import{t as G}from"./select-BzErv8r3.js";import{t as K}from"./inputnumber-DH7ykmI4.js";import{t as q}from"./multiselect-CQA3ZoTy.js";import{t as ee}from"./toggleswitch-Ds8FSqKu.js";import{t as te}from"./SkeletonTable-CaYzJhGc.js";import{t as J}from"./ConfirmDeleteDialog-qr422Ci3.js";import{t as Y}from"./SkeletonCard-DpG2AJij.js";import{t as X}from"./SelectLabel-DIRF6pqw.js";var ne={class:`space-y-4`},Z={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},re={class:`text-sm text-gray-500 dark:text-gray-400`},ie={class:`max-w-2xl`},ae={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},oe={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},se={class:`flex justify-end pt-2`},Q=`/api/v1/tenant/job-management/identifications`,ce={__name:`JobIdentificationSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},orgGradingId:{type:String,default:``},orgJobFamilyId:{type:String,default:``},gradingOptions:{type:Array,default:()=>[]},jobFamilyOptions:{type:Array,default:()=>[]}},emits:[`saved`],setup(r,{emit:i}){let a=i,o=r,{t:s}=I(),c=F(),l=d(!1),f=d(!0),_=d(``),v=d({}),y=d(``),x=d({grading_id:``}),S=C(()=>{let e=o.jobFamilyOptions.find(e=>e.value===o.orgJobFamilyId);return e?e.label:o.orgJobFamilyId||`-`});function T(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function E(){if(!o.orgId){f.value=!1;return}try{let e=(await j.get(Q,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];y.value=t.id,x.value.grading_id=t.grading_id||o.orgGradingId||``}else x.value.grading_id=o.orgGradingId||``}catch{x.value.grading_id=o.orgGradingId||``}finally{f.value=!1}}async function D(){if(_.value=``,v.value={},!x.value.grading_id){_.value=s(`job_management.grading_required`);return}l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,grading_id:x.value.grading_id,organization_id:o.orgId};if(y.value)await j.put(`${Q}/${y.value}`,{grading_id:x.value.grading_id});else{let t=await j.post(Q,e);y.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=T(e);Object.keys(t).length>0?(v.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}return n(E),(n,i)=>(t(),m(`div`,ne,[p(`div`,null,[p(`h2`,Z,g(e(s)(`job_management.identifications`)),1),p(`p`,re,g(e(s)(`job_management.identification_description`)),1)]),p(`div`,ie,[f.value?(t(),b(Y,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(t(),m(`div`,ae,[h(z,{label:e(s)(`organization.job_family`)},{default:u(()=>[h(H,{"model-value":S.value,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),h(z,{label:e(s)(`organization.grading`)},{default:u(()=>[h(e(G),{modelValue:x.value.grading_id,"onUpdate:modelValue":i[0]||=e=>x.value.grading_id=e,options:r.gradingOptions,"option-label":`label`,"option-value":`value`,placeholder:e(s)(`organization.select_grading`),class:`w-full`,size:`small`,invalid:!!v.value.grading_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`]),_.value?(t(),m(`div`,oe,g(_.value),1)):w(``,!0),p(`div`,se,[h(e(M),{label:e(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:!x.value.grading_id,onClick:D},null,8,[`label`,`loading`,`disabled`])])]))])]))}},le={class:`space-y-4`},ue={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},de={class:`text-sm text-gray-500 dark:text-gray-400`},fe={class:`max-w-2xl`},pe={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},me={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},he={class:`flex justify-end gap-2 pt-2`},ge=`/api/v1/tenant/job-management/objectives`,_e={__name:`JobObjectiveSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(r,{emit:i}){let a=i,o=r,{t:s}=I(),c=F(),l=d(!1),f=d(!0),_=d(!1),v=d(``),y=d({}),x=d(``),C=d(!1),T=d(``),E=d({objective:``});function D(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function O(){if(!o.orgId){f.value=!1;return}try{let e=(await j.get(ge,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];x.value=t.id,E.value.objective=t.objective||``}}catch{}finally{f.value=!1}}async function k(){v.value=``,y.value={},l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,objective:E.value.objective||``,organization_id:o.orgId};if(x.value)await j.put(`${ge}/${x.value}`,{objective:E.value.objective||``});else{let t=await j.post(ge,e);x.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=D(e);Object.keys(t).length>0?(y.value=t,v.value=Object.values(t).join(`, `)):v.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}async function A(){if(x.value){_.value=!0,T.value=``;try{await j.delete(`${ge}/${x.value}`),C.value=!1,x.value=``,E.value.objective=``,a(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{_.value=!1}}}return n(O),(n,r)=>(t(),m(`div`,le,[p(`div`,null,[p(`h2`,ue,g(e(s)(`job_management.objectives`)),1),p(`p`,de,g(e(s)(`job_management.objective_description`)),1)]),p(`div`,fe,[f.value?(t(),b(Y,{key:0,type:`detail`,count:1,rows:3,cols:`grid-cols-1`,padding:`p-5`})):(t(),m(`div`,pe,[h(z,{label:e(s)(`job_management.objective`)},{default:u(()=>[h(e(V),{modelValue:E.value.objective,"onUpdate:modelValue":r[0]||=e=>E.value.objective=e,rows:`3`,class:S([`w-full`,{"p-invalid":y.value.objective}]),placeholder:e(s)(`job_management.objective`)+`...`},null,8,[`modelValue`,`class`,`placeholder`])]),_:1},8,[`label`]),v.value?(t(),m(`div`,me,g(v.value),1)):w(``,!0),p(`div`,he,[x.value?(t(),b(e(M),{key:0,label:e(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:r[1]||=e=>C.value=!0},null,8,[`label`])):w(``,!0),h(e(M),{label:x.value?e(s)(`common.update`):e(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:l.value,onClick:k},null,8,[`label`,`loading`,`disabled`])])]))]),h(J,{visible:C.value,"onUpdate:visible":r[2]||=e=>C.value=e,loading:_.value,"error-msg":T.value,onConfirm:A,onCancel:r[3]||=e=>C.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},ve={class:`space-y-4`},ye={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},be={class:`text-sm text-gray-500 dark:text-gray-400`},xe={class:`max-w-2xl`},Se={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Ce={class:`pt-1`},we={class:`flex items-center gap-2 mb-3`},Te={class:`text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400`},Ee={class:`space-y-4`},De={class:`pt-4 border-t border-gray-200 dark:border-gray-700`},Oe={class:`flex items-center gap-2 mb-3`},ke={class:`text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400`},Ae={class:`space-y-4`},je={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Me={class:`flex justify-end gap-2 pt-2`},Ne=`/api/v1/tenant/job-management/education-experiences`,Pe={__name:`JobEduExpSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(r,{emit:i}){let a=i,o=r,{t:s}=I(),c=F(),l=d(!1),f=d(!0),_=d(!1),v=d(``),y=d({}),x=d(``),C=d(!1),T=d(``),E=d({education_id:``,education_major_id:[],job_family_id:[],experience_id:``}),D=d([]),O=d([]),k=d([]),A=d([]);async function N(){try{let[e,t,n,r]=await Promise.all([j.get(`/api/v1/tenant/job-management/values`,{params:{type:`education`,per_page:100}}),j.get(`/api/v1/tenant/job-management/values`,{params:{type:`experience`,per_page:100}}),j.get(`/api/v1/tenant/settings/education-majors?per_page=200`),j.get(`/api/v1/tenant/settings/job-families?per_page=100`)]);O.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),D.value=(t.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),k.value=(n.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id})),A.value=(r.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id}))}catch{}}async function P(){if(o.orgId)try{let e=(await j.get(Ne,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];x.value=t.id,E.value.education_id=t.education_id||``,E.value.education_major_id=Array.isArray(t.education_major_id)?t.education_major_id:[],E.value.job_family_id=Array.isArray(t.job_family_id)?t.job_family_id:[],E.value.experience_id=t.experience_id||``}}catch{}}async function R(){v.value=``,y.value={},l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,education_id:E.value.education_id||null,education_major_id:E.value.education_major_id||[],job_family_id:E.value.job_family_id||[],experience_id:E.value.experience_id||null,organization_id:o.orgId};if(x.value)await j.put(`${Ne}/${x.value}`,{education_id:E.value.education_id||``,education_major_id:E.value.education_major_id||[],job_family_id:E.value.job_family_id||[],experience_id:E.value.experience_id||``});else{let t=await j.post(Ne,e);x.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(y.value=t,v.value=Object.values(t).join(`, `)):v.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}async function B(){if(x.value){_.value=!0,T.value=``;try{await j.delete(`${Ne}/${x.value}`),C.value=!1,x.value=``,E.value.education_id=``,E.value.education_major_id=[],E.value.job_family_id=[],E.value.experience_id=``,a(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{_.value=!1}}}return n(async()=>{try{await Promise.all([N(),P()])}finally{f.value=!1}}),(n,r)=>(t(),m(`div`,ve,[p(`div`,null,[p(`h2`,ye,g(e(s)(`job_management.education_experience`)),1),p(`p`,be,g(e(s)(`job_management.education_experience_description`)),1)]),p(`div`,xe,[f.value?(t(),b(Y,{key:0,type:`detail`,count:1,rows:6,cols:`grid-cols-1`,padding:`p-5`})):(t(),m(`div`,Se,[p(`div`,Ce,[p(`div`,we,[r[7]||=p(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400`},[p(`i`,{class:`pi pi-graduation-cap text-sm`})],-1),p(`h3`,Te,g(e(s)(`job_management.group_education`)),1),r[8]||=p(`div`,{class:`flex-1 border-t border-gray-200 dark:border-gray-700`},null,-1)]),p(`div`,Ee,[h(z,{label:e(s)(`job_management.education_level`),errors:y.value?.education_id},{default:u(()=>[h(X,{modelValue:E.value.education_id,"onUpdate:modelValue":r[0]||=e=>E.value.education_id=e,options:O.value,placeholder:e(s)(`job_values.select_education`),class:S({"p-invalid":y.value?.education_id})},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h(z,{label:e(s)(`job_management.education_major`),errors:y.value?.education_major_id},{default:u(()=>[h(e(q),{modelValue:E.value.education_major_id,"onUpdate:modelValue":r[1]||=e=>E.value.education_major_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:e(s)(`common.select`),class:`w-full`,size:`small`,filter:``,showClear:``,display:`chip`,maxSelectedLabels:2,invalid:!!y.value.education_major_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`])])]),p(`div`,De,[p(`div`,Oe,[r[9]||=p(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400`},[p(`i`,{class:`pi pi-briefcase text-sm`})],-1),p(`h3`,ke,g(e(s)(`job_management.group_experience`)),1),r[10]||=p(`div`,{class:`flex-1 border-t border-gray-200 dark:border-gray-700`},null,-1)]),p(`div`,Ae,[h(z,{label:e(s)(`job_management.experience_range`),errors:y.value?.experience_id},{default:u(()=>[h(X,{modelValue:E.value.experience_id,"onUpdate:modelValue":r[2]||=e=>E.value.experience_id=e,options:D.value,placeholder:e(s)(`common.select`),class:S({"p-invalid":y.value?.experience_id})},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h(z,{label:e(s)(`job_management.job_family`),errors:y.value?.job_family_id},{default:u(()=>[h(e(q),{modelValue:E.value.job_family_id,"onUpdate:modelValue":r[3]||=e=>E.value.job_family_id=e,options:A.value,"option-label":`label`,"option-value":`value`,placeholder:e(s)(`common.select`),class:`w-full`,size:`small`,filter:``,showClear:``,display:`chip`,maxSelectedLabels:2,invalid:!!y.value.job_family_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`])])]),v.value?(t(),m(`div`,je,g(v.value),1)):w(``,!0),p(`div`,Me,[x.value?(t(),b(e(M),{key:0,label:e(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:r[4]||=e=>C.value=!0},null,8,[`label`])):w(``,!0),h(e(M),{label:x.value?e(s)(`common.update`):e(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:l.value,onClick:R},null,8,[`label`,`loading`,`disabled`])])]))]),h(J,{visible:C.value,"onUpdate:visible":r[5]||=e=>C.value=e,loading:_.value,"error-msg":T.value,onConfirm:B,onCancel:r[6]||=e=>C.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Fe=D.extend({name:`editor`,style:`
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
`,classes:{root:function(e){return[`p-editor`,{"p-invalid":e.instance.$invalid}]},toolbar:`p-editor-toolbar`,content:`p-editor-content`}}),Ie={name:`BaseEditor`,extends:B,props:{placeholder:String,readonly:Boolean,formats:Array,editorStyle:null,modules:null},style:Fe,provide:function(){return{$pcEditor:this,$parentInstance:this}}};function Le(e){"@babel/helpers - typeof";return Le=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},Le(e)}function Re(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function ze(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?Re(Object(n),!0).forEach(function(t){Be(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):Re(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function Be(e,t,n){return(t=Ve(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function Ve(e){var t=He(e,`string`);return Le(t)==`symbol`?t:t+``}function He(e,t){if(Le(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(Le(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var Ue=function(){try{return window.Quill}catch{return null}}(),We={name:`Editor`,extends:Ie,inheritAttrs:!1,emits:[`text-change`,`selection-change`,`load`],quill:null,watch:{modelValue:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},d_value:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},readonly:function(){this.handleReadOnlyChange()}},mounted:function(){var e=this,t={modules:ze({toolbar:this.$refs.toolbarElement},this.modules),readOnly:this.readonly,theme:`snow`,formats:this.formats,placeholder:this.placeholder};Ue?(this.quill=new Ue(this.$refs.editorElement,t),this.initQuill(),this.handleLoad()):P(()=>import(`./quill-BLmY9xB4.js`).then(function(n){n&&E(e.$refs.editorElement)&&(n.default?e.quill=new n.default(e.$refs.editorElement,t):e.quill=new n(e.$refs.editorElement,t),e.initQuill())}),__vite__mapDeps([0,1])).then(function(){e.handleLoad()})},beforeUnmount:function(){this.quill=null},methods:{renderValue:function(e){if(this.quill)if(e){var t=this.quill.clipboard.convert({html:e});this.quill.setContents(t)}else this.quill.setText(``)},initQuill:function(){var e=this;this.renderValue(this.d_value),this.quill.on(`text-change`,function(t,n,r){if(r===`user`){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();i===`<p><br></p>`&&(i=``),e.writeValue(i),e.$emit(`text-change`,{htmlValue:i,textValue:a,delta:t,source:r,instance:e.quill})}}),this.quill.on(`selection-change`,function(t,n,r){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();e.$emit(`selection-change`,{htmlValue:i,textValue:a,range:t,oldRange:n,source:r,instance:e.quill})})},handleLoad:function(){this.quill&&this.quill.getModule(`toolbar`)&&this.$emit(`load`,{instance:this.quill})},handleReadOnlyChange:function(){this.quill&&this.quill.enable(!this.readonly)}}};function Ge(e,n,r,i,a,o){return t(),m(`div`,T({class:e.cx(`root`)},e.ptmi(`root`)),[p(`div`,T({ref:`toolbarElement`,class:e.cx(`toolbar`)},e.ptm(`toolbar`)),[c(e.$slots,`toolbar`,{},function(){return[p(`span`,T({class:`ql-formats`},e.ptm(`formats`)),[p(`select`,T({class:`ql-header`,defaultValue:`0`},e.ptm(`header`)),[p(`option`,T({value:`1`},e.ptm(`option`)),`Heading`,16),p(`option`,T({value:`2`},e.ptm(`option`)),`Subheading`,16),p(`option`,T({value:`0`},e.ptm(`option`)),`Normal`,16)],16),p(`select`,T({class:`ql-font`},e.ptm(`font`)),[p(`option`,x(f(e.ptm(`option`))),null,16),p(`option`,T({value:`serif`},e.ptm(`option`)),null,16),p(`option`,T({value:`monospace`},e.ptm(`option`)),null,16)],16)],16),p(`span`,T({class:`ql-formats`},e.ptm(`formats`)),[p(`button`,T({class:`ql-bold`,type:`button`},e.ptm(`bold`)),null,16),p(`button`,T({class:`ql-italic`,type:`button`},e.ptm(`italic`)),null,16),p(`button`,T({class:`ql-underline`,type:`button`},e.ptm(`underline`)),null,16)],16),p(`span`,T({class:`ql-formats`},e.ptm(`formats`)),[p(`select`,T({class:`ql-color`},e.ptm(`color`)),null,16),p(`select`,T({class:`ql-background`},e.ptm(`background`)),null,16)],16),p(`span`,T({class:`ql-formats`},e.ptm(`formats`)),[p(`button`,T({class:`ql-list`,value:`ordered`,type:`button`},e.ptm(`list`)),null,16),p(`button`,T({class:`ql-list`,value:`bullet`,type:`button`},e.ptm(`list`)),null,16),p(`select`,T({class:`ql-align`},e.ptm(`select`)),[p(`option`,T({defaultValue:``},e.ptm(`option`)),null,16),p(`option`,T({value:`center`},e.ptm(`option`)),null,16),p(`option`,T({value:`right`},e.ptm(`option`)),null,16),p(`option`,T({value:`justify`},e.ptm(`option`)),null,16)],16)],16),p(`span`,T({class:`ql-formats`},e.ptm(`formats`)),[p(`button`,T({class:`ql-link`,type:`button`},e.ptm(`link`)),null,16),p(`button`,T({class:`ql-image`,type:`button`},e.ptm(`image`)),null,16),p(`button`,T({class:`ql-code-block`,type:`button`},e.ptm(`codeBlock`)),null,16)],16),p(`span`,T({class:`ql-formats`},e.ptm(`formats`)),[p(`button`,T({class:`ql-clean`,type:`button`},e.ptm(`clean`)),null,16)],16)]})],16),p(`div`,T({ref:`editorElement`,class:e.cx(`content`),style:e.editorStyle},e.ptm(`content`)),null,16)],16)}We.render=Ge;var Ke={key:0,class:`text-gray-500 dark:text-gray-400 text-xs`},qe=[`innerHTML`],Je={key:2,class:`text-gray-800 dark:text-gray-100`},Ye={class:`flex items-center gap-1`},Xe={__name:`DataTableSection`,props:{items:Array,loading:Boolean,total:Number,columns:{type:Array,default:()=>[]},entity:String,orgId:String,onLoad:Function},emits:[`edit`,`delete`],setup(a){let o=a,{t:l}=I(),f=d(1),_=d(15),v=C(()=>(f.value-1)*_.value),x=C(()=>[...o.columns.map(e=>({type:`text`,width:`w-24`,headerWidth:`w-20`})),{type:`icons`,count:2,headerWidth:`w-16`}]);function S(e){f.value=e.page+1,_.value=e.rows,o.onLoad&&o.onLoad(f.value,_.value)}return n(()=>{o.onLoad&&o.onLoad(1,15)}),(n,o)=>{let d=r(`tooltip`);return a.loading?(t(),b(te,{key:0,columns:x.value,rows:8},null,8,[`columns`])):(t(),b(e(U),{key:1,value:a.items,lazy:``,totalRecords:a.total,first:v.value,rows:_.value,onPage:S,paginator:``,paginatorTemplate:`FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown`,rowsPerPageOptions:[10,15,25,50],size:`small`,class:`!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden`},{empty:u(()=>[c(n.$slots,`empty`)]),default:u(()=>[(t(!0),m(y,null,s(a.columns,n=>(t(),b(e(W),{key:n.field,field:n.field,header:n.header,sortable:``},{body:u(({data:e})=>[n.field.startsWith(`_`)?(t(),m(`span`,Ke,g(e[n.field]||`-`),1)):w(``,!0),n.html?(t(),m(`div`,{key:1,class:`editor-content`,innerHTML:e[n.field]},null,8,qe)):(t(),m(`span`,Je,g(e[n.field]||`-`),1))]),_:2},1032,[`field`,`header`]))),128)),h(e(W),{header:e(l)(`common.actions`),style:{width:`90px`},frozen:``,alignFrozen:`right`},{body:u(({data:t})=>[p(`div`,Ye,[i(h(e(M),{icon:`pi pi-pencil`,size:`small`,text:``,severity:`secondary`,onClick:e=>n.$emit(`edit`,t)},null,8,[`onClick`]),[[d,e(l)(`common.edit`),void 0,{left:!0}]]),i(h(e(M),{icon:`pi pi-trash`,size:`small`,text:``,severity:`danger`,onClick:e=>n.$emit(`delete`,t)},null,8,[`onClick`]),[[d,e(l)(`common.delete`),void 0,{left:!0}]])])]),_:1},8,[`header`])]),_:3},8,[`value`,`totalRecords`,`first`,`rows`]))}}},Ze={class:`space-y-4`},Qe={key:0,class:`bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded p-3 text-xs text-red-700 dark:text-red-300`},$e={__name:`DialogForm`,props:{visible:Boolean,title:String,saving:Boolean,errors:{type:Object,default:()=>({})},width:{type:String,default:`480px`}},emits:[`save`,`cancel`],setup(n){let r=n,{t:i}=I(),a=C(()=>r.width===`maximize`?`90vw`:r.width);return(r,o)=>(t(),b(e(N),{visible:n.visible,"onUpdate:visible":o[2]||=e=>r.$emit(`update:visible`,e),header:n.title,modal:``,style:v({width:a.value}),class:`p-fluid`,closable:!n.saving},{footer:u(()=>[h(e(M),{label:e(i)(`common.cancel`),size:`small`,outlined:``,severity:`secondary`,disabled:n.saving,onClick:o[0]||=e=>r.$emit(`cancel`)},null,8,[`label`,`disabled`]),h(e(M),{label:e(i)(`common.save`),icon:`pi pi-check`,size:`small`,loading:n.saving,onClick:o[1]||=e=>r.$emit(`save`)},null,8,[`label`,`loading`])]),default:u(()=>[p(`div`,Ze,[c(r.$slots,`default`),Object.keys(n.errors).length?(t(),m(`div`,Qe,[(t(!0),m(y,null,s(n.errors,(e,n)=>(t(),m(`p`,{key:n,class:`mb-1`},[p(`strong`,null,g(n)+`:`,1),_(` `+g(Array.isArray(e)?e.join(`, `):e),1)]))),128))])):w(``,!0)])]),_:3},8,[`visible`,`header`,`style`,`closable`]))}},et={class:`space-y-4`},tt={class:`flex items-center justify-between`},nt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},rt={class:`text-sm text-gray-500 dark:text-gray-400`},it={class:`flex flex-col items-center justify-center py-10 text-gray-400`},at={class:`text-sm font-medium`},ot=`/api/v1/tenant/job-management/responsibilities`,st={__name:`JobResponsibilitySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(n,{emit:r}){let i=n,a=r,{t:o}=I(),s=F(),c=d([]),l=d(!1),f=d(0),_=d(!1),v=d(!1),b=d(``),x=d(!1),T=d({}),E=d(!1),D=d(!1),O=d(``),k=d(null),A=d({main_task:``,activities:``,outputs:``,success_indicators:``}),N=C(()=>{let e=o(`job_management.responsibilities_title`);return v.value?`${e}`:`${o(`common.create`)} ${e}`}),P=C(()=>[{field:`main_task`,header:o(`job_management.main_task`),html:!0},{field:`activities`,header:o(`job_management.activities`),html:!0},{field:`outputs`,header:o(`job_management.outputs`),html:!0},{field:`success_indicators`,header:o(`job_management.success_indicators`),html:!0}]);async function R(e,t){l.value=!0;try{let n=await j.get(ot,{params:{page:e,per_page:t,organization_id:i.orgId}}),r=n.data?.data||[];c.value=r.map(e=>({...e,main_task:e.main_task,activities:e.activities,outputs:e.outputs,success_indicators:e.success_indicators})),f.value=n.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function B(){v.value=!1,b.value=``,A.value={main_task:``,activities:``,outputs:``,success_indicators:``},T.value={},_.value=!0}function V(e){v.value=!0,b.value=e.id,A.value={main_task:e.main_task||``,activities:e.activities||``,outputs:e.outputs||``,success_indicators:e.success_indicators||``},T.value={},_.value=!0}async function H(){x.value=!0,T.value={};try{let e={nomenclature:i.orgName||``,full_code:i.orgCode||``,...A.value,organization_id:i.orgId};v.value?await j.put(`${ot}/${b.value}`,e):await j.post(ot,e),_.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),R(1,15)}catch(e){let t=L(e);Object.keys(t).length?T.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{x.value=!1}}function U(e){k.value=e,O.value=``,E.value=!0}async function W(){if(k.value){D.value=!0,O.value=``;try{await j.delete(`${ot}/${k.value.id}`),E.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),R(1,15)}catch(e){O.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{D.value=!1}}}return(r,i)=>(t(),m(`div`,et,[p(`div`,tt,[p(`div`,null,[p(`h2`,nt,g(e(o)(`job_management.responsibilities_title`)),1),p(`p`,rt,g(e(o)(`job_management.responsibilities_description`)),1)]),h(e(M),{label:e(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:i[0]||=e=>B()},null,8,[`label`])]),h(Xe,{items:c.value,loading:l.value,total:f.value,columns:P.value,entity:`responsibilities`,"org-id":n.orgId,"on-load":R,onEdit:V,onDelete:U},{empty:u(()=>[p(`div`,it,[i[9]||=p(`i`,{class:`pi pi-list-check text-3xl mb-2 opacity-50`},null,-1),p(`p`,at,g(e(o)(`job_management.empty_responsibilities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),h($e,{visible:_.value,"onUpdate:visible":i[5]||=e=>_.value=e,title:N.value,saving:x.value,errors:T.value,width:`maximize`,onSave:H,onCancel:i[6]||=e=>_.value=!1},{default:u(()=>[_.value?(t(),m(y,{key:0},[h(z,{label:e(o)(`job_management.main_task`),errors:T.value?.main_task},{default:u(()=>[h(e(We),{modelValue:A.value.main_task,"onUpdate:modelValue":i[1]||=e=>A.value.main_task=e,editorStyle:`height:120px`,class:S({"p-invalid":T.value?.main_task})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),h(z,{label:e(o)(`job_management.activities`),errors:T.value?.activities},{default:u(()=>[h(e(We),{modelValue:A.value.activities,"onUpdate:modelValue":i[2]||=e=>A.value.activities=e,editorStyle:`height:120px`,class:S({"p-invalid":T.value?.activities})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),h(z,{label:e(o)(`job_management.outputs`),errors:T.value?.outputs},{default:u(()=>[h(e(We),{modelValue:A.value.outputs,"onUpdate:modelValue":i[3]||=e=>A.value.outputs=e,editorStyle:`height:120px`,class:S({"p-invalid":T.value?.outputs})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),h(z,{label:e(o)(`job_management.success_indicators`),errors:T.value?.success_indicators},{default:u(()=>[h(e(We),{modelValue:A.value.success_indicators,"onUpdate:modelValue":i[4]||=e=>A.value.success_indicators=e,editorStyle:`height:120px`,class:S({"p-invalid":T.value?.success_indicators})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])],64)):w(``,!0)]),_:1},8,[`visible`,`title`,`saving`,`errors`]),h(J,{visible:E.value,"onUpdate:visible":i[7]||=e=>E.value=e,loading:D.value,"error-msg":O.value,onConfirm:W,onCancel:i[8]||=e=>E.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},ct={class:`space-y-4`},lt={class:`flex items-center justify-between`},ut={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},dt={class:`text-sm text-gray-500 dark:text-gray-400`},ft={class:`flex flex-col items-center justify-center py-10 text-gray-400`},pt={class:`text-sm font-medium`},mt=`/api/v1/tenant/job-management/hr-authorities`,ht={__name:`JobHRAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(n,{emit:r}){let i=n,a=r,{t:o}=I(),s=F(),c=d([]),l=d(!1),f=d(0),_=d(!1),v=d(!1),y=d(``),b=d(!1),x=d({}),w=d(!1),T=d(!1),E=d(``),D=d(null),O=d({description:``}),k=C(()=>{let e=o(`job_management.hr_authorities`);return v.value?`${o(`common.edit`)} ${e}`:`${o(`common.create`)} ${e}`}),A=C(()=>[{field:`description`,header:o(`job_management.description`)}]);async function N(e,t){l.value=!0;try{let n=await j.get(mt,{params:{page:e,per_page:t,organization_id:i.orgId}});c.value=n.data?.data||[],f.value=n.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function P(){v.value=!1,y.value=``,O.value={nomenclature:``,full_code:``,description:``},x.value={},_.value=!0}function R(e){v.value=!0,y.value=e.id,O.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},x.value={},_.value=!0}async function B(){b.value=!0,x.value={};try{let e={...O.value,nomenclature:i.orgName||``,full_code:i.orgCode||``,organization_id:i.orgId};v.value?await j.put(`${mt}/${y.value}`,e):await j.post(mt,e),_.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),N(1,15)}catch(e){let t=L(e);Object.keys(t).length?x.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{b.value=!1}}function H(e){D.value=e,E.value=``,w.value=!0}async function U(){if(D.value){T.value=!0,E.value=``;try{await j.delete(`${mt}/${D.value.id}`),w.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),N(1,15)}catch(e){E.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{T.value=!1}}}return(r,i)=>(t(),m(`div`,ct,[p(`div`,lt,[p(`div`,null,[p(`h2`,ut,g(e(o)(`job_management.hr_authorities`)),1),p(`p`,dt,g(e(o)(`job_management.authority_description`)),1)]),h(e(M),{label:e(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:i[0]||=e=>P()},null,8,[`label`])]),h(Xe,{items:c.value,loading:l.value,total:f.value,columns:A.value,entity:`hr-authorities`,"org-id":n.orgId,"on-load":N,onEdit:R,onDelete:H},{empty:u(()=>[p(`div`,ft,[i[6]||=p(`i`,{class:`pi pi-users text-3xl mb-2 opacity-50`},null,-1),p(`p`,pt,g(e(o)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),h($e,{visible:_.value,"onUpdate:visible":i[2]||=e=>_.value=e,title:k.value,saving:b.value,errors:x.value,onSave:B,onCancel:i[3]||=e=>_.value=!1},{default:u(()=>[h(z,{label:e(o)(`job_management.description`),errors:x.value?.description},{default:u(()=>[h(e(V),{modelValue:O.value.description,"onUpdate:modelValue":i[1]||=e=>O.value.description=e,rows:`3`,class:S([`w-full`,{"p-invalid":x.value?.description}])},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),h(J,{visible:w.value,"onUpdate:visible":i[4]||=e=>w.value=e,loading:T.value,"error-msg":E.value,onConfirm:U,onCancel:i[5]||=e=>w.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},gt={class:`space-y-4`},_t={class:`flex items-center justify-between`},vt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},yt={class:`text-sm text-gray-500 dark:text-gray-400`},bt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},xt={class:`text-sm font-medium`},St=`/api/v1/tenant/job-management/operational-authorities`,Ct={__name:`JobOpAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(n,{emit:r}){let i=n,a=r,{t:o}=I(),s=F(),c=d([]),l=d(!1),f=d(0),_=d(!1),v=d(!1),y=d(``),b=d(!1),x=d({}),w=d(!1),T=d(!1),E=d(``),D=d(null),O=d({description:``}),k=C(()=>{let e=o(`job_management.op_authorities`);return v.value?`${o(`common.edit`)} ${e}`:`${o(`common.create`)} ${e}`}),A=C(()=>[{field:`description`,header:o(`job_management.description`)}]);async function N(e,t){l.value=!0;try{let n=await j.get(St,{params:{page:e,per_page:t,organization_id:i.orgId}});c.value=n.data?.data||[],f.value=n.data?.total||0}catch(e){s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.failed_to_load`),life:4e3})}finally{l.value=!1}}function P(){v.value=!1,y.value=``,O.value={nomenclature:``,full_code:``,description:``},x.value={},_.value=!0}function R(e){v.value=!0,y.value=e.id,O.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},x.value={},_.value=!0}async function B(){b.value=!0,x.value={};try{let e={...O.value,nomenclature:i.orgName||``,full_code:i.orgCode||``,organization_id:i.orgId};v.value?await j.put(`${St}/${y.value}`,e):await j.post(St,e),_.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.saved`),life:2e3}),N(1,15)}catch(e){let t=L(e);Object.keys(t).length?x.value=t:s.add({severity:`error`,detail:e.response?.data?.error?.message||o(`message.operation_failed`),life:4e3})}finally{b.value=!1}}function H(e){D.value=e,E.value=``,w.value=!0}async function U(){if(D.value){T.value=!0,E.value=``;try{await j.delete(`${St}/${D.value.id}`),w.value=!1,a(`saved`),s.add({severity:`success`,detail:o(`message.deleted`),life:2e3}),N(1,15)}catch(e){E.value=e.response?.data?.error?.message||o(`message.operation_failed`)}finally{T.value=!1}}}return(r,i)=>(t(),m(`div`,gt,[p(`div`,_t,[p(`div`,null,[p(`h2`,vt,g(e(o)(`job_management.op_authorities`)),1),p(`p`,yt,g(e(o)(`job_management.authority_description`)),1)]),h(e(M),{label:e(o)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:i[0]||=e=>P()},null,8,[`label`])]),h(Xe,{items:c.value,loading:l.value,total:f.value,columns:A.value,entity:`operational-authorities`,"org-id":n.orgId,"on-load":N,onEdit:R,onDelete:H},{empty:u(()=>[p(`div`,bt,[i[6]||=p(`i`,{class:`pi pi-cog text-3xl mb-2 opacity-50`},null,-1),p(`p`,xt,g(e(o)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),h($e,{visible:_.value,"onUpdate:visible":i[2]||=e=>_.value=e,title:k.value,saving:b.value,errors:x.value,onSave:B,onCancel:i[3]||=e=>_.value=!1},{default:u(()=>[h(z,{label:e(o)(`job_management.description`),errors:x.value?.description},{default:u(()=>[h(e(V),{modelValue:O.value.description,"onUpdate:modelValue":i[1]||=e=>O.value.description=e,class:S([`w-full`,{"p-invalid":x.value?.description}]),rows:`3`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),h(J,{visible:w.value,"onUpdate:visible":i[4]||=e=>w.value=e,loading:T.value,"error-msg":E.value,onConfirm:U,onCancel:i[5]||=e=>w.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},wt={class:`space-y-4`},Tt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Et={class:`text-sm text-gray-500 dark:text-gray-400`},Dt={class:`max-w-2xl`},Ot={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},kt={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},At={class:`flex justify-end gap-2 pt-2`},jt=`/api/v1/tenant/job-management/working-activities`,Mt={__name:`JobActivitySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(r,{emit:i}){let a=i,o=r,{t:s}=I(),c=F(),l=d(!1),f=d(!0),_=d(``),v=d({}),y=d(``),x=d(!1),C=d(!1),T=d(``),E=d({job_management_value_id:``}),D=d([]);async function O(){try{let e=await j.get(`/api/v1/tenant/job-management/values`,{params:{type:`activity`,per_page:100}});D.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id}))}catch{}}async function k(){if(!o.orgId){f.value=!1;return}try{let e=(await j.get(jt,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];y.value=t.id,E.value.job_management_value_id=t.job_management_value_id||``}}catch{}}async function A(){_.value=``,v.value={},l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,job_management_value_id:E.value.job_management_value_id||null,organization_id:o.orgId};if(y.value)await j.put(`${jt}/${y.value}`,{job_management_value_id:E.value.job_management_value_id||``});else{let t=await j.post(jt,e);y.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(v.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}async function N(){if(y.value){C.value=!0,T.value=``;try{await j.delete(`${jt}/${y.value}`),x.value=!1,y.value=``,E.value.job_management_value_id=``,a(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{C.value=!1}}}return n(async()=>{try{await Promise.all([O(),k()])}finally{f.value=!1}}),(n,r)=>(t(),m(`div`,wt,[p(`div`,null,[p(`h2`,Tt,g(e(s)(`job_management.activities`)),1),p(`p`,Et,g(e(s)(`job_management.activity_description`)),1)]),p(`div`,Dt,[f.value?(t(),b(Y,{key:0,type:`detail`,count:1,rows:3,cols:`grid-cols-1`,padding:`p-5`})):(t(),m(`div`,Ot,[h(z,{label:e(s)(`job_values.types.activity`),errors:v.value?.job_management_value_id},{default:u(()=>[h(X,{modelValue:E.value.job_management_value_id,"onUpdate:modelValue":r[0]||=e=>E.value.job_management_value_id=e,options:D.value,"option-label":`label`,"option-value":`value`,placeholder:e(s)(`common.select`),class:S({"p-invalid":v.value?.job_management_value_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),_.value?(t(),m(`div`,kt,g(_.value),1)):w(``,!0),p(`div`,At,[y.value?(t(),b(e(M),{key:0,label:e(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:r[1]||=e=>x.value=!0},null,8,[`label`])):w(``,!0),h(e(M),{label:y.value?e(s)(`common.update`):e(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:l.value,onClick:A},null,8,[`label`,`loading`,`disabled`])])]))]),h(J,{visible:x.value,"onUpdate:visible":r[2]||=e=>x.value=e,loading:C.value,"error-msg":T.value,onConfirm:N,onCancel:r[3]||=e=>x.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Nt={class:`space-y-4`},Pt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Ft={class:`text-sm text-gray-500 dark:text-gray-400`},It={class:`max-w-2xl`},Lt={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Rt={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},zt={class:`flex justify-end gap-2 pt-2`},Bt=`/api/v1/tenant/job-management/working-risks`,Vt={__name:`JobRiskSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(r,{emit:i}){let a=i,o=r,{t:s}=I(),c=F(),l=d(!1),f=d(!0),_=d(``),v=d({}),y=d(``),x=d(!1),C=d(!1),T=d(``),E=d({job_management_value_environment_id:``,job_management_value_hazard_id:``}),D=d([]),O=d([]);async function k(){try{let[e,t]=await Promise.all([j.get(`/api/v1/tenant/job-management/values`,{params:{type:`environment`,per_page:100}}),j.get(`/api/v1/tenant/job-management/values`,{params:{type:`risk`,per_page:100}})]);D.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),O.value=(t.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id}))}catch{}}async function A(){if(!o.orgId){f.value=!1;return}try{let e=(await j.get(Bt,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];y.value=t.id,E.value.job_management_value_environment_id=t.job_management_value_environment_id||``,E.value.job_management_value_hazard_id=t.job_management_value_hazard_id||``}}catch{}}async function N(){_.value=``,v.value={},l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,job_management_value_environment_id:E.value.job_management_value_environment_id||null,job_management_value_hazard_id:E.value.job_management_value_hazard_id||null,organization_id:o.orgId};if(y.value)await j.put(`${Bt}/${y.value}`,{job_management_value_environment_id:E.value.job_management_value_environment_id||``,job_management_value_hazard_id:E.value.job_management_value_hazard_id||``});else{let t=await j.post(Bt,e);y.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(v.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}async function P(){if(y.value){C.value=!0,T.value=``;try{await j.delete(`${Bt}/${y.value}`),x.value=!1,y.value=``,E.value.job_management_value_environment_id=``,E.value.job_management_value_hazard_id=``,a(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{C.value=!1}}}return n(async()=>{try{await Promise.all([k(),A()])}finally{f.value=!1}}),(n,r)=>(t(),m(`div`,Nt,[p(`div`,null,[p(`h2`,Pt,g(e(s)(`job_management.risks`)),1),p(`p`,Ft,g(e(s)(`job_management.risk_description`)),1)]),p(`div`,It,[f.value?(t(),b(Y,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(t(),m(`div`,Lt,[h(z,{label:e(s)(`job_management.work_environment`),errors:v.value?.job_management_value_environment_id},{default:u(()=>[h(X,{modelValue:E.value.job_management_value_environment_id,"onUpdate:modelValue":r[0]||=e=>E.value.job_management_value_environment_id=e,options:D.value,"option-label":`label`,"option-value":`value`,placeholder:e(s)(`common.select`),class:S({"p-invalid":v.value?.job_management_value_environment_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h(z,{label:e(s)(`job_management.risk`),errors:v.value?.job_management_value_hazard_id},{default:u(()=>[h(X,{modelValue:E.value.job_management_value_hazard_id,"onUpdate:modelValue":r[1]||=e=>E.value.job_management_value_hazard_id=e,options:O.value,"option-label":`label`,"option-value":`value`,placeholder:e(s)(`common.select`),class:S({"p-invalid":v.value?.job_management_value_hazard_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),_.value?(t(),m(`div`,Rt,g(_.value),1)):w(``,!0),p(`div`,zt,[y.value?(t(),b(e(M),{key:0,label:e(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:r[2]||=e=>x.value=!0},null,8,[`label`])):w(``,!0),h(e(M),{label:y.value?e(s)(`common.update`):e(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:l.value,onClick:N},null,8,[`label`,`loading`,`disabled`])])]))]),h(J,{visible:x.value,"onUpdate:visible":r[3]||=e=>x.value=e,loading:C.value,"error-msg":T.value,onConfirm:P,onCancel:r[4]||=e=>x.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Ht={class:`space-y-4`},Ut={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Wt={class:`text-sm text-gray-500 dark:text-gray-400`},Gt={class:`max-w-2xl`},Kt={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},qt={class:`pt-1`},Jt={class:`flex items-center gap-2 mb-3`},Yt={class:`text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400`},Xt={class:`space-y-4`},Zt={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Qt={class:`flex justify-end gap-2 pt-2`},$t={class:`max-w-3xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5 space-y-4`},en={class:`flex items-center justify-between gap-2 flex-wrap`},tn={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},nn={class:`text-sm text-gray-500 dark:text-gray-400`},rn={key:0,class:`text-sm text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg px-3 py-2`},an={key:1,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},on={key:2,class:`overflow-x-auto border border-gray-200 dark:border-gray-700 rounded-lg`},sn={class:`w-full text-sm`},cn={class:`bg-gray-50 dark:bg-gray-700/40 text-left`},ln={class:`px-3 py-2 font-semibold text-gray-600 dark:text-gray-300`},un={class:`px-3 py-2 font-semibold text-gray-600 dark:text-gray-300`},dn={class:`px-3 py-2 align-top text-gray-500 dark:text-gray-400`},fn={class:`px-3 py-2`},pn={class:`px-3 py-2`},mn={class:`px-3 py-2 align-top`},hn={key:3,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},gn={key:4,class:`flex justify-end gap-2 pt-2`},$=`/api/v1/tenant/job-management/relationships`,_n={__name:`JobRelationshipSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},orgSummaryId:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(r,{emit:i}){let a=i,o=r,{t:c}=I(),l=F(),f=d(!1),_=d(!0),v=d(``),x=d({}),C=d(``),T=d(!1),E=d(!1),D=d(``),O=d({job_management_value_relationship_id:``,job_management_value_frequency_id:``}),k=d([]),A=d([]),N=d([]),P=d([]),R=d(!1),B=d(``);async function V(){try{let[e,t,n]=await Promise.all([j.get(`/api/v1/tenant/job-management/values`,{params:{type:`relationship`,per_page:100}}),j.get(`/api/v1/tenant/job-management/values`,{params:{type:`frequency`,per_page:100}}),o.orgSummaryId?j.get(`/api/v1/tenant/organizations`,{params:{summary_id:o.orgSummaryId,per_page:100}}):Promise.resolve({data:{data:[]}})]);k.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),A.value=(t.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),N.value=(n.data?.data||[]).filter(e=>e.id!==o.orgId).map(e=>({label:e.full_code?`${e.full_code} - ${e.nomenclature}`:e.nomenclature,value:e.id}))}catch{}}async function U(){if(!o.orgId){_.value=!1;return}try{let e=(await j.get($,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];C.value=t.id,O.value.job_management_value_relationship_id=t.job_management_value_relationship_id||``,O.value.job_management_value_frequency_id=t.job_management_value_frequency_id||``,await ne()}}catch{}}async function W(){v.value=``,x.value={},f.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,job_management_value_relationship_id:O.value.job_management_value_relationship_id||null,job_management_value_frequency_id:O.value.job_management_value_frequency_id||null,organization_id:o.orgId};if(C.value)await j.put(`${$}/${C.value}`,{job_management_value_relationship_id:O.value.job_management_value_relationship_id||``,job_management_value_frequency_id:O.value.job_management_value_frequency_id||``});else{let t=await j.post($,e);C.value=t.data?.data?.id||``}l.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(x.value=t,v.value=Object.values(t).join(`, `)):v.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{f.value=!1}}async function G(){if(C.value){E.value=!0,D.value=``;try{await j.delete(`${$}/${C.value}`),T.value=!1,C.value=``,O.value.job_management_value_relationship_id=``,O.value.job_management_value_frequency_id=``,P.value=[],a(`saved`),l.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){D.value=e?.response?.data?.error?.message||c(`message.operation_failed`)}finally{E.value=!1}}}let K=0;function q(){C.value&&P.value.push({_key:`new-${++K}`,id:``,organization_id:``,activity:``})}function ee(e){let t=P.value[e];t&&(t.id?te(t.id,e):P.value.splice(e,1))}async function te(e,t){try{await j.delete(`${$}/${C.value}/details/${e}`),P.value.splice(t,1),l.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){l.add({severity:`error`,summary:c(`message.error`),detail:e?.response?.data?.error?.message||c(`message.operation_failed`),life:4e3})}}async function ne(){if(C.value)try{let e=await j.get(`${$}/${C.value}/details`);P.value=(e.data?.data||[]).map(e=>({_key:`db-${++K}`,id:e.id,organization_id:e.organization_id||``,activity:e.activity||``}))}catch{}}async function Z(){if(!(!C.value||R.value)){B.value=``,R.value=!0;try{for(let e of P.value){let t={organization_id:e.organization_id||``,activity:e.activity||``};e.id?await j.put(`${$}/${C.value}/details/${e.id}`,t):e.id=(await j.post(`${$}/${C.value}/details`,t)).data?.data?.id||``}await ne(),l.add({severity:`success`,summary:c(`message.success`),detail:c(`job_management.relationship_details_saved`),life:2e3})}catch(e){let t=L(e);Object.keys(t).length>0?B.value=Object.values(t).join(`, `):B.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{R.value=!1}}}return n(async()=>{try{await Promise.all([V(),U()])}finally{_.value=!1}}),(n,r)=>(t(),m(`div`,Ht,[p(`div`,null,[p(`h2`,Ut,g(e(c)(`job_management.relationships`)),1),p(`p`,Wt,g(e(c)(`job_management.relationship_description`)),1)]),p(`div`,Gt,[_.value?(t(),b(Y,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(t(),m(`div`,Kt,[p(`div`,qt,[p(`div`,Jt,[r[5]||=p(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400`},[p(`i`,{class:`pi pi-compass text-sm`})],-1),p(`h3`,Yt,g(e(c)(`job_management.relationship_group_scope`)),1),r[6]||=p(`div`,{class:`flex-1 border-t border-gray-200 dark:border-gray-700`},null,-1)]),p(`div`,Xt,[h(z,{label:e(c)(`job_management.relationship_type`),errors:x.value?.job_management_value_relationship_id},{default:u(()=>[h(X,{modelValue:O.value.job_management_value_relationship_id,"onUpdate:modelValue":r[0]||=e=>O.value.job_management_value_relationship_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:e(c)(`common.select`),class:S({"p-invalid":x.value?.job_management_value_relationship_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h(z,{label:e(c)(`job_management.frequency`),errors:x.value?.job_management_value_frequency_id},{default:u(()=>[h(X,{modelValue:O.value.job_management_value_frequency_id,"onUpdate:modelValue":r[1]||=e=>O.value.job_management_value_frequency_id=e,options:A.value,"option-label":`label`,"option-value":`value`,placeholder:e(c)(`common.select`),class:S({"p-invalid":x.value?.job_management_value_frequency_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`])])]),v.value?(t(),m(`div`,Zt,g(v.value),1)):w(``,!0),p(`div`,Qt,[C.value?(t(),b(e(M),{key:0,label:e(c)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:r[2]||=e=>T.value=!0},null,8,[`label`])):w(``,!0),h(e(M),{label:C.value?e(c)(`common.update`):e(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:f.value,disabled:f.value,onClick:W},null,8,[`label`,`loading`,`disabled`])])]))]),p(`div`,$t,[p(`div`,en,[p(`div`,null,[p(`h3`,tn,g(e(c)(`job_management.relationship_details`)),1),p(`p`,nn,g(e(c)(`job_management.relationship_details_description`)),1)]),h(e(M),{label:e(c)(`job_management.add_relationship_detail`),icon:`pi pi-plus`,size:`small`,outlined:``,disabled:!C.value||R.value,onClick:q},null,8,[`label`,`disabled`])]),C.value?P.value.length===0?(t(),m(`div`,an,g(e(c)(`job_management.no_relationship_details`)),1)):w(``,!0):(t(),m(`div`,rn,g(e(c)(`job_management.save_relationship_first`)),1)),P.value.length>0?(t(),m(`div`,on,[p(`table`,sn,[p(`thead`,null,[p(`tr`,cn,[r[7]||=p(`th`,{class:`px-3 py-2 w-10 font-semibold text-gray-600 dark:text-gray-300`},`#`,-1),p(`th`,ln,g(e(c)(`job_management.relationship_organization`)),1),p(`th`,un,g(e(c)(`job_management.relationship_activity`)),1),r[8]||=p(`th`,{class:`px-3 py-2 w-12`},null,-1)])]),p(`tbody`,null,[(t(!0),m(y,null,s(P.value,(n,r)=>(t(),m(`tr`,{key:n._key,class:`border-t border-gray-200 dark:border-gray-700`},[p(`td`,dn,g(r+1),1),p(`td`,fn,[h(X,{modelValue:n.organization_id,"onUpdate:modelValue":e=>n.organization_id=e,options:N.value,"option-label":`label`,"option-value":`value`,placeholder:e(c)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),p(`td`,pn,[h(H,{modelValue:n.activity,"onUpdate:modelValue":e=>n.activity=e,placeholder:e(c)(`job_management.relationship_activity`)},null,8,[`modelValue`,`onUpdate:modelValue`,`placeholder`])]),p(`td`,mn,[h(e(M),{icon:`pi pi-trash`,severity:`danger`,size:`small`,text:``,rounded:``,"aria-label":`Remove`,onClick:e=>ee(r)},null,8,[`onClick`])])]))),128))])])])):w(``,!0),B.value?(t(),m(`div`,hn,g(B.value),1)):w(``,!0),P.value.length>0?(t(),m(`div`,gn,[h(e(M),{label:e(c)(`job_management.save_relationship_details`),icon:`pi pi-save`,size:`small`,loading:R.value,disabled:R.value||!C.value,onClick:Z},null,8,[`label`,`loading`,`disabled`])])):w(``,!0)]),h(J,{visible:T.value,"onUpdate:visible":r[3]||=e=>T.value=e,loading:E.value,"error-msg":D.value,onConfirm:G,onCancel:r[4]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},vn={class:`space-y-4`},yn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},bn={class:`text-sm text-gray-500 dark:text-gray-400`},xn={class:`max-w-2xl`},Sn={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Cn={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},wn={class:`flex justify-end gap-2 pt-2`},Tn=`/api/v1/tenant/job-management/subordinate-controls`,En={__name:`JobSubordinateSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(r,{emit:i}){let a=i,o=r,{t:s}=I(),c=F(),l=d(!1),f=d(!0),_=d(``),v=d({}),y=d(``),x=d(!1),C=d(!1),T=d(``),E=d({job_management_value_id:``}),D=d([]);async function O(){try{let e=await j.get(`/api/v1/tenant/job-management/values`,{params:{type:`subordinate`,per_page:100}});D.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id}))}catch{}}async function k(){if(!o.orgId){f.value=!1;return}try{let e=(await j.get(Tn,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];y.value=t.id,E.value.job_management_value_id=t.job_management_value_id||``}}catch{}}async function A(){_.value=``,v.value={},l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,job_management_value_id:E.value.job_management_value_id||null,organization_id:o.orgId};if(y.value)await j.put(`${Tn}/${y.value}`,{job_management_value_id:E.value.job_management_value_id||``});else{let t=await j.post(Tn,e);y.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(v.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}async function N(){if(y.value){C.value=!0,T.value=``;try{await j.delete(`${Tn}/${y.value}`),x.value=!1,y.value=``,E.value.job_management_value_id=``,a(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{C.value=!1}}}return n(async()=>{try{await Promise.all([O(),k()])}finally{f.value=!1}}),(n,r)=>(t(),m(`div`,vn,[p(`div`,null,[p(`h2`,yn,g(e(s)(`job_management.subordinate_controls`)),1),p(`p`,bn,g(e(s)(`job_management.subordinate_description`)),1)]),p(`div`,xn,[f.value?(t(),b(Y,{key:0,type:`detail`,count:1,rows:3,cols:`grid-cols-1`,padding:`p-5`})):(t(),m(`div`,Sn,[h(z,{label:e(s)(`job_management.control_type`),errors:v.value?.job_management_value_id},{default:u(()=>[h(X,{modelValue:E.value.job_management_value_id,"onUpdate:modelValue":r[0]||=e=>E.value.job_management_value_id=e,options:D.value,"option-label":`label`,"option-value":`value`,placeholder:e(s)(`common.select`),class:S({"p-invalid":v.value?.job_management_value_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),_.value?(t(),m(`div`,Cn,g(_.value),1)):w(``,!0),p(`div`,wn,[y.value?(t(),b(e(M),{key:0,label:e(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:r[1]||=e=>x.value=!0},null,8,[`label`])):w(``,!0),h(e(M),{label:y.value?e(s)(`common.update`):e(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:l.value,onClick:A},null,8,[`label`,`loading`,`disabled`])])]))]),h(J,{visible:x.value,"onUpdate:visible":r[2]||=e=>x.value=e,loading:C.value,"error-msg":T.value,onConfirm:N,onCancel:r[3]||=e=>x.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Dn={class:`space-y-4`},On={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},kn={class:`text-sm text-gray-500 dark:text-gray-400`},An={class:`max-w-2xl`},jn={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Mn={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Nn={class:`flex justify-end gap-2 pt-2`},Pn=`/api/v1/tenant/job-management/assets`,Fn={__name:`JobAssetSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(r,{emit:i}){let a=i,o=r,{t:s}=I(),c=F(),l=d(!1),f=d(!0),_=d(``),v=d({}),y=d(``),x=d(!1),C=d(!1),T=d(``),E=d({job_management_value_asset_id:``,job_management_value_authority_id:``}),D=d([]),O=d([]);async function k(){try{let[e,t]=await Promise.all([j.get(`/api/v1/tenant/job-management/values`,{params:{type:`asset`,per_page:100}}),j.get(`/api/v1/tenant/job-management/values`,{params:{type:`asset_authority`,per_page:100}})]);D.value=(e.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id})),O.value=(t.data?.data||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions}`,value:e.id}))}catch{}}async function A(){if(!o.orgId){f.value=!1;return}try{let e=(await j.get(Pn,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];y.value=t.id,E.value.job_management_value_asset_id=t.job_management_value_asset_id||``,E.value.job_management_value_authority_id=t.job_management_value_authority_id||``}}catch{}}async function N(){_.value=``,v.value={},l.value=!0;try{let e={nomenclature:o.orgName||``,full_code:o.orgCode||``,job_management_value_asset_id:E.value.job_management_value_asset_id||null,job_management_value_authority_id:E.value.job_management_value_authority_id||null,organization_id:o.orgId};if(y.value)await j.put(`${Pn}/${y.value}`,{job_management_value_asset_id:E.value.job_management_value_asset_id||``,job_management_value_authority_id:E.value.job_management_value_authority_id||``});else{let t=await j.post(Pn,e);y.value=t.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(v.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{l.value=!1}}async function P(){if(y.value){C.value=!0,T.value=``;try{await j.delete(`${Pn}/${y.value}`),x.value=!1,y.value=``,E.value.job_management_value_asset_id=``,E.value.job_management_value_authority_id=``,a(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){T.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{C.value=!1}}}return n(async()=>{try{await Promise.all([k(),A()])}finally{f.value=!1}}),(n,r)=>(t(),m(`div`,Dn,[p(`div`,null,[p(`h2`,On,g(e(s)(`job_management.assets`)),1),p(`p`,kn,g(e(s)(`job_management.asset_description`)),1)]),p(`div`,An,[f.value?(t(),b(Y,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(t(),m(`div`,jn,[h(z,{label:e(s)(`job_management.asset_type`),errors:v.value?.job_management_value_asset_id},{default:u(()=>[h(X,{modelValue:E.value.job_management_value_asset_id,"onUpdate:modelValue":r[0]||=e=>E.value.job_management_value_asset_id=e,options:D.value,"option-label":`label`,"option-value":`value`,placeholder:e(s)(`common.select`),class:S({"p-invalid":v.value?.job_management_value_asset_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h(z,{label:e(s)(`job_management.authority_level`),errors:v.value?.job_management_value_authority_id},{default:u(()=>[h(X,{modelValue:E.value.job_management_value_authority_id,"onUpdate:modelValue":r[1]||=e=>E.value.job_management_value_authority_id=e,options:O.value,"option-label":`label`,"option-value":`value`,placeholder:e(s)(`common.select`),class:S({"p-invalid":v.value?.job_management_value_authority_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),_.value?(t(),m(`div`,Mn,g(_.value),1)):w(``,!0),p(`div`,Nn,[y.value?(t(),b(e(M),{key:0,label:e(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:r[2]||=e=>x.value=!0},null,8,[`label`])):w(``,!0),h(e(M),{label:y.value?e(s)(`common.update`):e(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:l.value,disabled:l.value,onClick:N},null,8,[`label`,`loading`,`disabled`])])]))]),h(J,{visible:x.value,"onUpdate:visible":r[3]||=e=>x.value=e,loading:C.value,"error-msg":T.value,onConfirm:P,onCancel:r[4]||=e=>x.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},In={class:`space-y-4`},Ln={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Rn={class:`text-sm text-gray-500 dark:text-gray-400`},zn={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Bn={class:`flex items-center justify-between gap-4 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3`},Vn={class:`min-w-0`},Hn={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},Un={class:`text-xs text-gray-500 dark:text-gray-400 mt-0.5`},Wn={class:`space-y-4 pt-4 border-t border-gray-200 dark:border-gray-700`},Gn={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Kn={class:`flex justify-end gap-2 pt-2`},qn=`/api/v1/tenant/job-management/financials`,Jn={__name:`JobFinancialSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(r,{emit:i}){let a=i,o=r,{t:s}=I(),c=F(),f=d(!1),_=d(!0),v=d(``),y=d({}),x=d(``),T=d(!1),E=d(!1),D=d(``),O=d({is_authorized:!1,job_management_value_cash_id:``,job_management_value_authority_id:``,job_management_value_impact_id:``}),k=d([]),A=d([]),N=d([]),P=d([]),R=d([]),B=C(()=>O.value.is_authorized?A.value:N.value),V=C(()=>O.value.is_authorized?P.value:R.value);async function H(){try{let[e,t,n,r,i]=await Promise.all([j.get(`/api/v1/tenant/job-management/values`,{params:{type:`cash`,per_page:100}}),j.get(`/api/v1/tenant/job-management/values`,{params:{type:`authority`,per_page:100}}),j.get(`/api/v1/tenant/job-management/values`,{params:{type:`authority_unauthorized`,per_page:100}}),j.get(`/api/v1/tenant/job-management/values`,{params:{type:`impact`,per_page:100}}),j.get(`/api/v1/tenant/job-management/values`,{params:{type:`impact_unauthorized`,per_page:100}})]);k.value=(e.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),A.value=(t.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),N.value=(n.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),P.value=(r.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),R.value=(i.data?.data||[]).map(e=>({label:e.descriptions,value:e.id}))}catch{}}let U=!1;l(()=>O.value.is_authorized,(e,t)=>{U||e===t||(O.value.job_management_value_cash_id=``,O.value.job_management_value_authority_id=``,O.value.job_management_value_impact_id=``)},{flush:`sync`});async function W(){if(!o.orgId){_.value=!1;return}try{let e=(await j.get(qn,{params:{organization_id:o.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];U=!0,x.value=t.id,O.value.is_authorized=!!t.is_authorized,O.value.job_management_value_cash_id=t.job_management_value_cash_id||``,O.value.job_management_value_authority_id=t.job_management_value_authority_id||``,O.value.job_management_value_impact_id=t.job_management_value_impact_id||``,U=!1}}catch{}}async function G(){v.value=``,y.value={},f.value=!0;try{let e=!!O.value.is_authorized,t={nomenclature:o.orgName||``,full_code:o.orgCode||``,is_authorized:e,job_management_value_cash_id:e&&O.value.job_management_value_cash_id||null,job_management_value_authority_id:O.value.job_management_value_authority_id||null,job_management_value_impact_id:O.value.job_management_value_impact_id||null,organization_id:o.orgId};if(x.value)await j.put(`${qn}/${x.value}`,{is_authorized:e,job_management_value_cash_id:e&&O.value.job_management_value_cash_id||``,job_management_value_authority_id:O.value.job_management_value_authority_id||``,job_management_value_impact_id:O.value.job_management_value_impact_id||``});else{let e=await j.post(qn,t);x.value=e.data?.data?.id||``}c.add({severity:`success`,summary:s(`message.success`),detail:s(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=L(e);Object.keys(t).length>0?(y.value=t,v.value=Object.values(t).join(`, `)):v.value=e?.response?.data?.error?.message||e.message||s(`message.operation_failed`)}finally{f.value=!1}}async function K(){if(x.value){E.value=!0,D.value=``;try{await j.delete(`${qn}/${x.value}`),T.value=!1,x.value=``,O.value.is_authorized=!1,O.value.job_management_value_cash_id=``,O.value.job_management_value_authority_id=``,O.value.job_management_value_impact_id=``,a(`saved`),c.add({severity:`success`,summary:s(`message.success`),detail:s(`message.deleted`),life:2e3})}catch(e){D.value=e?.response?.data?.error?.message||s(`message.operation_failed`)}finally{E.value=!1}}}return n(async()=>{try{await Promise.all([H(),W()])}finally{_.value=!1}}),(n,r)=>(t(),m(`div`,In,[p(`div`,null,[p(`h2`,Ln,g(e(s)(`job_management.financials`)),1),p(`p`,Rn,g(e(s)(`job_management.financial_description`)),1)]),p(`div`,null,[_.value?(t(),b(Y,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(t(),m(`div`,zn,[p(`div`,Bn,[p(`div`,Vn,[p(`p`,Hn,g(e(s)(`job_management.is_authorized`)),1),p(`p`,Un,g(e(s)(`job_management.is_authorized_description`)),1)]),h(e(ee),{modelValue:O.value.is_authorized,"onUpdate:modelValue":r[0]||=e=>O.value.is_authorized=e},null,8,[`modelValue`])]),p(`div`,Wn,[O.value.is_authorized?(t(),b(z,{key:0,label:e(s)(`job_management.cash_level`),errors:y.value?.job_management_value_cash_id},{default:u(()=>[h(X,{modelValue:O.value.job_management_value_cash_id,"onUpdate:modelValue":r[1]||=e=>O.value.job_management_value_cash_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:e(s)(`common.select`),class:S({"p-invalid":y.value?.job_management_value_cash_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`])):w(``,!0),h(z,{label:e(s)(`job_management.authority_level`),errors:y.value?.job_management_value_authority_id},{default:u(()=>[h(X,{modelValue:O.value.job_management_value_authority_id,"onUpdate:modelValue":r[2]||=e=>O.value.job_management_value_authority_id=e,options:B.value,"option-label":`label`,"option-value":`value`,placeholder:e(s)(`common.select`),class:S({"p-invalid":y.value?.job_management_value_authority_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h(z,{label:e(s)(`job_management.impact_level`),errors:y.value?.job_management_value_impact_id},{default:u(()=>[h(X,{modelValue:O.value.job_management_value_impact_id,"onUpdate:modelValue":r[3]||=e=>O.value.job_management_value_impact_id=e,options:V.value,"option-label":`label`,"option-value":`value`,placeholder:e(s)(`common.select`),class:S({"p-invalid":y.value?.job_management_value_impact_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`])]),v.value?(t(),m(`div`,Gn,g(v.value),1)):w(``,!0),p(`div`,Kn,[x.value?(t(),b(e(M),{key:0,label:e(s)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:r[4]||=e=>T.value=!0},null,8,[`label`])):w(``,!0),h(e(M),{label:x.value?e(s)(`common.update`):e(s)(`common.save`),icon:`pi pi-check`,size:`small`,loading:f.value,disabled:f.value,onClick:G},null,8,[`label`,`loading`,`disabled`])])]))]),h(J,{visible:T.value,"onUpdate:visible":r[5]||=e=>T.value=e,loading:E.value,"error-msg":D.value,onConfirm:K,onCancel:r[6]||=e=>T.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Yn=`/api/v1/tenant/job-management/potency-competencies`;function Xn({orgId:e,rows:t,afterDelete:n,onSaved:r,matchBy:i=`value`,descriptionField:a=`descriptions`}){let{t:o}=I(),s=F(),c=d(!1),l=d(``),u=d(!1),f=d(!1),p=d(``),m=d(null),h=d([]);function g(e){let t=(e.levelOptions||[]).find(t=>t.value===e.job_management_value_id);return t&&t[a]||``}function _(e){if(i===`competency`)return e.competency_id&&h.value.find(t=>t.competency_id&&t.competency_id===e.competency_id)||null;let t=new Set((e.levelOptions||[]).map(e=>e.value));return h.value.find(e=>e.job_management_value_id&&t.has(e.job_management_value_id))||null}function v(){t.value.forEach(e=>{let t=_(e);e.recordId=t?t.id:``,e.job_management_value_id=t&&t.job_management_value_id||``,e.weight!==void 0&&(e.weight=t?t.weight??e.weight:e.weight)})}async function y(){if(!e.value){h.value=[];return}try{let t=await j.get(Yn,{params:{organization_id:e.value,per_page:100}});h.value=t.data?.data||[]}catch{h.value=[]}}function b(e){m.value=e,p.value=``,u.value=!0}async function x(){let e=m.value;if(e){f.value=!0,p.value=``;try{e.recordId&&await j.delete(`${Yn}/${e.recordId}`),n&&n(e),u.value=!1,await y(),v(),s.add({severity:`success`,summary:o(`message.success`),detail:o(`message.deleted`),life:2e3}),r&&r()}catch(e){p.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{f.value=!1,m.value=null}}}async function S(){l.value=``,c.value=!0;try{for(let n of t.value)if(n.job_management_value_id){let t=n.competency_id?{competency_id:n.competency_id,job_management_value_id:n.job_management_value_id}:{job_management_value_id:n.job_management_value_id};n.weight!==void 0&&n.weight!==null&&n.weight!==``&&(t.weight=n.weight),n.recordId?await j.put(`${Yn}/${n.recordId}`,t):n.recordId=(await j.post(Yn,{organization_id:e.value,...t})).data?.data?.id||``}else n.recordId&&=(await j.delete(`${Yn}/${n.recordId}`),``);await y(),v(),s.add({severity:`success`,summary:o(`message.success`),detail:o(`common.saved`),life:2e3}),r&&r()}catch(e){let t=L(e);Object.keys(t).length>0?l.value=Object.values(t).join(`, `):l.value=e?.response?.data?.error?.message||e.message||o(`message.operation_failed`)}finally{c.value=!1}}return{savingCard:c,errorMsg:l,deleteVisible:u,deleting:f,deleteError:p,deleteTarget:m,records:h,levelDescription:g,hydrateRows:v,loadData:y,askDeleteRow:b,handleDelete:x,handleSave:S}}var Zn={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Qn={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},$n={class:`text-sm text-gray-500 dark:text-gray-400`},er={key:0,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},tr={key:1,class:`overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700`},nr={class:`w-full text-sm`},rr={class:`bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400`},ir={class:`px-4 py-3 font-semibold min-w-[220px]`},ar={class:`px-4 py-3 font-semibold min-w-[260px]`},or={class:`px-4 py-3 font-semibold min-w-[260px]`},sr={class:`px-4 py-3 font-semibold w-16 text-right`},cr={class:`px-4 py-3`},lr={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},ur={class:`px-4 py-3`},dr={class:`px-4 py-3 text-sm text-gray-600 dark:text-gray-300`},fr={class:`px-4 py-3 text-right`},pr={key:2,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},mr={key:3,class:`flex justify-end gap-2 pt-1`},hr={__name:`SelectablePotencyCard`,props:{orgId:String,typeGroup:{type:String,required:!0},skeletonRows:{type:Number,default:5},titleKey:{type:String,required:!0},descriptionKey:{type:String,required:!0},emptyKey:{type:String,required:!0},saveLabelKey:{type:String,required:!0},deleteTitleKey:{type:String,required:!0},deleteMessageKey:{type:String,required:!0}},emits:[`saved`],setup(r,{emit:i}){let a=i,c=r,{t:u}=I(),f=d(!0),_=d([]),v=d([]),x=d([]),{savingCard:S,errorMsg:T,deleteVisible:E,deleting:D,deleteError:O,deleteTarget:k,records:A,levelDescription:N,hydrateRows:P,loadData:F,askDeleteRow:L,handleDelete:R,handleSave:z}=Xn({orgId:C(()=>c.orgId),rows:v,afterDelete:e=>{let t=Array.isArray(x.value)?x.value:[];x.value=t.filter(t=>t!==e.type)},onSaved:()=>a(`saved`)}),B=C(()=>(_.value||[]).find(e=>e.type_group===c.typeGroup));function V(e){let t=`job_values.types.${e.type}`,n=u(t);return n===t?e.description_group||e.type:n}let H=C(()=>(B.value?.types||[]).map(e=>({label:V(e),value:e.type})));function U(){let e={};(B.value?.types||[]).forEach(t=>{e[t.type]=t});let t=Array.isArray(x.value)?x.value:x.value?[x.value]:[];v.value=t.filter(t=>e[t]).map(t=>{let n=e[t];return{competency_id:``,competency_name:V(n),competency_definition:``,type:n.type,levelOptions:(n.options||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``})),recordId:``,job_management_value_id:``}})}async function W(){try{let e=await j.get(`/api/v1/tenant/job-management/values/tree`);_.value=e.data?.data||[],U()}catch{_.value=[],v.value=[]}}function G(){let e={};(B.value?.types||[]).forEach(t=>{(t.options||[]).forEach(n=>{e[n.id]=t.type})});let t=[];A.value.forEach(n=>{let r=n.job_management_value_id&&e[n.job_management_value_id];r&&!t.includes(r)&&t.push(r)}),x.value=t,U(),P()}return l(x,()=>{U(),P()}),n(async()=>{try{await Promise.all([W(),F()])}finally{G(),f.value=!1}}),(n,i)=>(t(),m(`div`,Zn,[p(`div`,null,[p(`h3`,Qn,g(e(u)(r.titleKey)),1),p(`p`,$n,g(e(u)(r.descriptionKey)),1)]),f.value?(t(),b(Y,{key:0,type:`detail`,count:1,rows:r.skeletonRows,cols:`grid-cols-1`,padding:`p-5`},null,8,[`rows`])):(t(),m(y,{key:1},[h(X,{modelValue:x.value,"onUpdate:modelValue":i[0]||=e=>x.value=e,options:H.value,"option-label":`label`,"option-value":`value`,placeholder:e(u)(`common.select`),showClear:``,multiple:``},null,8,[`modelValue`,`options`,`placeholder`]),v.value.length===0?(t(),m(`div`,er,g(e(u)(r.emptyKey)),1)):(t(),m(`div`,tr,[p(`table`,nr,[p(`thead`,null,[p(`tr`,rr,[p(`th`,ir,g(e(u)(`job_management.potency_table_name`)),1),p(`th`,ar,g(e(u)(`job_management.potency_table_level`)),1),p(`th`,or,g(e(u)(`job_management.potency_table_description`)),1),p(`th`,sr,g(e(u)(`common.actions`)),1)])]),p(`tbody`,null,[(t(!0),m(y,null,s(v.value,n=>(t(),m(`tr`,{key:n.type,class:`border-t border-gray-100 dark:border-gray-700 align-top`},[p(`td`,cr,[p(`div`,lr,g(n.competency_name),1)]),p(`td`,ur,[h(X,{modelValue:n.job_management_value_id,"onUpdate:modelValue":e=>n.job_management_value_id=e,options:n.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:e(u)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),p(`td`,dr,g(e(N)(n)),1),p(`td`,fr,[h(e(M),{icon:`pi pi-trash`,severity:`danger`,text:``,rounded:``,size:`small`,disabled:e(S),"aria-label":e(u)(`common.delete`),onClick:t=>e(L)(n)},null,8,[`disabled`,`aria-label`,`onClick`])])]))),128))])])])),e(T)?(t(),m(`div`,pr,g(e(T)),1)):w(``,!0),v.value.length>0?(t(),m(`div`,mr,[h(e(M),{label:e(u)(r.saveLabelKey),icon:`pi pi-check`,size:`small`,loading:e(S),disabled:e(S)||!r.orgId,onClick:e(z)},null,8,[`label`,`loading`,`disabled`,`onClick`])])):w(``,!0)],64)),h(J,{visible:e(E),"onUpdate:visible":i[1]||=e=>o(E)?E.value=e:null,title:e(u)(r.deleteTitleKey),message:e(u)(r.deleteMessageKey,{name:e(k)?.competency_name||``}),loading:e(D),"error-msg":e(O),onConfirm:e(R),onCancel:i[2]||=e=>E.value=!1},null,8,[`visible`,`title`,`message`,`loading`,`error-msg`,`onConfirm`])]))}},gr={__name:`PsychologicalPotencyCard`,props:{orgId:String},emits:[`saved`],setup(e,{emit:n}){let r=n;return(n,i)=>(t(),b(hr,{"org-id":e.orgId,"type-group":`psychological`,"skeleton-rows":5,"title-key":`job_management.potency_required_title`,"description-key":`job_management.potency_required_description`,"empty-key":`job_management.potency_required_empty`,"save-label-key":`job_management.save_potency_levels`,"delete-title-key":`job_management.potency_confirm_delete_title`,"delete-message-key":`job_management.potency_confirm_delete`,onSaved:i[0]||=e=>r(`saved`)},null,8,[`org-id`]))}},_r={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},vr={class:`flex items-start justify-between gap-4`},yr={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},br={class:`text-sm text-gray-500 dark:text-gray-400`},xr={class:`flex flex-col items-end gap-1 shrink-0`},Sr={class:`flex items-center gap-2`},Cr={class:`text-xs font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap`},wr={class:`w-24 shrink-0`},Tr={key:0,class:`text-xs text-red-500 dark:text-red-400`},Er={key:0,class:`text-sm text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg px-3 py-2`},Dr={key:1,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},Or={key:2,class:`overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700`},kr={class:`w-full text-sm`},Ar={class:`bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400`},jr={class:`px-4 py-3 font-semibold min-w-[220px]`},Mr={class:`px-4 py-3 font-semibold min-w-[260px]`},Nr={class:`px-4 py-3 font-semibold min-w-[130px]`},Pr={class:`px-4 py-3 font-semibold min-w-[260px]`},Fr={class:`px-4 py-3 font-semibold w-16 text-right`},Ir={class:`px-4 py-3`},Lr={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},Rr={key:0,class:`text-xs text-gray-400 dark:text-gray-500`},zr={class:`px-4 py-3`},Br={class:`px-4 py-3`},Vr={class:`px-4 py-3 text-sm text-gray-600 dark:text-gray-300`},Hr={class:`px-4 py-3 text-right`},Ur={key:3,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Wr={key:4,class:`flex justify-end gap-2 pt-1`},Gr={__name:`TechnicalPotencyCard`,props:{orgId:String},emits:[`saved`,`weight-saved`],setup(r,{emit:i}){let a=i,c=r,{t:u}=I(),f=F(),_=d(!0),v=d([]),x=d([]),S=d([]),T=d([]),E=d([]),D=C(()=>E.value.length>0),O=d(``),k=d(``),A=d(``),N=d(!1),P=d(``),{savingCard:L,errorMsg:R,deleteVisible:z,deleting:B,deleteError:V,deleteTarget:H,records:U,levelDescription:W,hydrateRows:G,loadData:q,askDeleteRow:ee,handleDelete:te,handleSave:ne}=Xn({orgId:C(()=>c.orgId),rows:S,matchBy:`competency`,descriptionField:`note`,afterDelete:e=>{let t=Array.isArray(T.value)?T.value:[];T.value=t.filter(t=>t!==e.competency_id)},onSaved:()=>a(`saved`)}),Z=C(()=>(v.value||[]).map(e=>({id:e.id,name:e.name,cluster:e.cluster||``}))),re=C(()=>{let e={};return(Z.value||[]).forEach(t=>{(e[t.cluster]=e[t.cluster]||[]).push(t)}),Object.keys(e).sort().map(t=>({label:t,items:e[t].sort((e,t)=>e.name.localeCompare(t.name))}))}),ie=C(()=>(x.value||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``,note:e.note||``})));function ae(){let e={};(Z.value||[]).forEach(t=>{e[t.id]=t});let t=(Array.isArray(T.value)?T.value:T.value?[T.value]:[]).filter(t=>e[t]),n=t.length>0?Math.round(100/t.length*100)/100:0;S.value=t.map(t=>{let r=e[t];return{competency_id:t,competency_name:r.name,cluster:r.cluster,levelOptions:ie.value,recordId:``,job_management_value_id:``,weight:n}})}async function oe(){try{let[e,t]=await Promise.all([j.get(`/api/v1/tenant/settings/competencies`,{params:{per_page:500}}),j.get(`/api/v1/tenant/job-management/values/clusters/technical`)]);E.value=t.data?.data?.clusters||[];let n=new Set(E.value),r=e.data?.data||[];v.value=r.filter(e=>e.cluster&&n.has(e.cluster))}catch{v.value=[]}}async function se(){try{let e=await j.get(`/api/v1/tenant/job-management/values`,{params:{type:`technical`,per_page:100}});x.value=e.data?.data||[]}catch{x.value=[]}}async function Q(){if(!c.orgId){O.value=``,k.value=``;return}try{let e=((await j.get(`/api/v1/tenant/job-management/competency-groups`,{params:{organization_id:c.orgId}})).data?.data||[]).find(e=>e.category===`technical`);k.value=e?e.id:``,O.value=e?e.weight:``,A.value=O.value}catch{O.value=``,k.value=``,A.value=``}}async function ce(){if(O.value===``||O.value===null||O.value===void 0){P.value=u(`job_management.potency_technical_weight_required`);return}if(!(A.value!==``&&O.value===A.value)){N.value=!0,P.value=``;try{let e=k.value;if(!e)try{let t=((await j.get(`/api/v1/tenant/job-management/competency-groups`,{params:{organization_id:c.orgId}})).data?.data||[]).find(e=>e.category===`technical`);t&&(e=t.id)}catch{}let t={weight:O.value};e?await j.put(`/api/v1/tenant/job-management/competency-groups/${e}`,t):await j.post(`/api/v1/tenant/job-management/competency-groups`,{organization_id:c.orgId,category:`technical`,weight:O.value}),A.value=O.value,f.add({severity:`success`,summary:u(`message.success`),detail:u(`job_management.potency_technical_weight_saved`),life:2e3}),a(`saved`),a(`weight-saved`,O.value),await Q()}catch(e){P.value=e?.response?.data?.error?.message||e.message||u(`message.operation_failed`)}finally{N.value=!1}}}function le(){let e={};(Z.value||[]).forEach(t=>{e[t.id]=t});let t=[];U.value.forEach(n=>{n.competency_id&&e[n.competency_id]&&!t.includes(n.competency_id)&&t.push(n.competency_id)}),T.value=t,ae(),G()}return l(T,()=>{ae(),G()}),n(async()=>{try{await Promise.all([oe(),se(),q(),Q()])}finally{le(),_.value=!1}}),(n,i)=>(t(),m(`div`,_r,[p(`div`,vr,[p(`div`,null,[p(`h3`,yr,g(e(u)(`job_management.potency_technical_title`)),1),p(`p`,br,g(e(u)(`job_management.potency_technical_description`)),1)]),p(`div`,xr,[p(`div`,Sr,[p(`label`,Cr,g(e(u)(`job_management.potency_technical_weight_label`)),1),p(`div`,wr,[h(e(K),{modelValue:O.value,"onUpdate:modelValue":i[0]||=e=>O.value=e,fluid:``,min:0,max:100,suffix:`%`,size:`small`,disabled:N.value||!r.orgId,onBlur:ce},null,8,[`modelValue`,`disabled`])])]),P.value?(t(),m(`div`,Tr,g(P.value),1)):w(``,!0)])]),_.value?(t(),b(Y,{key:0,type:`detail`,count:1,rows:8,cols:`grid-cols-1`,padding:`p-5`})):(t(),m(y,{key:1},[h(X,{modelValue:T.value,"onUpdate:modelValue":i[1]||=e=>T.value=e,options:re.value,"option-label":`name`,"option-value":`id`,"option-group-label":`label`,"option-group-children":`items`,placeholder:e(u)(`job_management.potency_technical_placeholder`),showClear:``,multiple:``},null,8,[`modelValue`,`options`,`placeholder`]),D.value?S.value.length===0?(t(),m(`div`,Dr,g(e(u)(`job_management.potency_technical_empty`)),1)):(t(),m(`div`,Or,[p(`table`,kr,[p(`thead`,null,[p(`tr`,Ar,[p(`th`,jr,g(e(u)(`job_management.potency_table_name`)),1),p(`th`,Mr,g(e(u)(`job_management.potency_table_level`)),1),p(`th`,Nr,g(e(u)(`job_management.potency_table_weight`)),1),p(`th`,Pr,g(e(u)(`job_management.potency_table_description`)),1),p(`th`,Fr,g(e(u)(`common.actions`)),1)])]),p(`tbody`,null,[(t(!0),m(y,null,s(S.value,n=>(t(),m(`tr`,{key:n.competency_id,class:`border-t border-gray-100 dark:border-gray-700 align-top`},[p(`td`,Ir,[p(`div`,Lr,g(n.competency_name),1),n.cluster?(t(),m(`div`,Rr,g(n.cluster),1)):w(``,!0)]),p(`td`,zr,[h(X,{modelValue:n.job_management_value_id,"onUpdate:modelValue":e=>n.job_management_value_id=e,options:n.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:e(u)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),p(`td`,Br,[h(e(K),{modelValue:n.weight,"onUpdate:modelValue":e=>n.weight=e,class:`!w-full`,min:0,max:100,suffix:`%`,size:`small`},null,8,[`modelValue`,`onUpdate:modelValue`])]),p(`td`,Vr,g(e(W)(n)),1),p(`td`,Hr,[h(e(M),{icon:`pi pi-trash`,severity:`danger`,text:``,rounded:``,size:`small`,disabled:e(L),"aria-label":e(u)(`common.delete`),onClick:t=>e(ee)(n)},null,8,[`disabled`,`aria-label`,`onClick`])])]))),128))])])])):(t(),m(`div`,Er,g(e(u)(`job_management.potency_technical_no_mapping`)),1)),e(R)?(t(),m(`div`,Ur,g(e(R)),1)):w(``,!0),S.value.length>0?(t(),m(`div`,Wr,[h(e(M),{label:e(u)(`job_management.save_technical`),icon:`pi pi-check`,size:`small`,loading:e(L),disabled:e(L)||!r.orgId,onClick:e(ne)},null,8,[`label`,`loading`,`disabled`,`onClick`])])):w(``,!0)],64)),h(J,{visible:e(z),"onUpdate:visible":i[2]||=e=>o(z)?z.value=e:null,title:e(u)(`job_management.potency_confirm_delete_title`),message:e(u)(`job_management.potency_confirm_delete`,{name:e(H)?.competency_name||``}),loading:e(B),"error-msg":e(V),onConfirm:e(te),onCancel:i[3]||=e=>z.value=!1},null,8,[`visible`,`title`,`message`,`loading`,`error-msg`,`onConfirm`])]))}},Kr={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},qr={class:`flex items-start justify-between gap-4`},Jr={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},Yr={class:`text-sm text-gray-500 dark:text-gray-400`},Xr={class:`flex flex-col items-end gap-1 shrink-0`},Zr={class:`flex items-center gap-2`},Qr={class:`text-xs font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap`},$r={class:`w-24 shrink-0 text-right`},ei={key:0,class:`text-sm font-semibold text-gray-700 dark:text-gray-200`},ti={key:1,class:`text-sm text-gray-400 dark:text-gray-500`},ni={key:0,class:`pi pi-spin pi-spinner text-sm text-gray-400`},ri={key:0,class:`text-xs text-red-500 dark:text-red-400`},ii={key:0,class:`text-sm text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg px-3 py-2`},ai={key:1,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},oi={key:2,class:`overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700`},si={class:`w-full text-sm`},ci={class:`bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400`},li={class:`px-4 py-3 font-semibold min-w-[220px]`},ui={class:`px-4 py-3 font-semibold min-w-[260px]`},di={class:`px-4 py-3 font-semibold min-w-[130px]`},fi={class:`px-4 py-3 font-semibold min-w-[260px]`},pi={class:`px-4 py-3 font-semibold w-16 text-right`},mi={class:`px-4 py-3`},hi={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},gi={key:0,class:`text-xs text-gray-400 dark:text-gray-500`},_i={class:`px-4 py-3`},vi={class:`px-4 py-3`},yi={class:`px-4 py-3 text-sm text-gray-600 dark:text-gray-300`},bi={class:`px-4 py-3 text-right`},xi={key:3,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Si={key:4,class:`flex justify-end gap-2 pt-1`},Ci={__name:`ManagerialPotencyCard`,props:{orgId:String,technicalWeight:{type:Number,default:null}},emits:[`saved`],setup(r,{emit:i}){let a=i,c=r,{t:u}=I(),f=F(),_=d(!0),v=d([]),x=d([]),S=d([]),T=d([]),E=d([]),D=C(()=>E.value.length>0),O=d(``),k=d(``),A=d(``),N=d(!1),P=d(``),L=C(()=>{let e=O.value;return e===``||e==null?null:Math.round((100-e)*100)/100}),{savingCard:R,errorMsg:z,deleteVisible:B,deleting:V,deleteError:H,deleteTarget:U,records:W,levelDescription:G,hydrateRows:q,loadData:ee,askDeleteRow:te,handleDelete:ne,handleSave:Z}=Xn({orgId:C(()=>c.orgId),rows:S,matchBy:`competency`,descriptionField:`note`,afterDelete:e=>{let t=Array.isArray(T.value)?T.value:[];T.value=t.filter(t=>t!==e.competency_id)},onSaved:()=>a(`saved`)}),re=C(()=>(v.value||[]).map(e=>({id:e.id,name:e.name,cluster:e.cluster||``}))),ie=C(()=>{let e={};return(re.value||[]).forEach(t=>{(e[t.cluster]=e[t.cluster]||[]).push(t)}),Object.keys(e).sort().map(t=>({label:t,items:e[t].sort((e,t)=>e.name.localeCompare(t.name))}))}),ae=C(()=>(x.value||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``,note:e.note||``})));function oe(){let e={};(re.value||[]).forEach(t=>{e[t.id]=t});let t=(Array.isArray(T.value)?T.value:T.value?[T.value]:[]).filter(t=>e[t]),n=t.length>0?Math.round(100/t.length*100)/100:0;S.value=t.map(t=>{let r=e[t];return{competency_id:t,competency_name:r.name,cluster:r.cluster,levelOptions:ae.value,recordId:``,job_management_value_id:``,weight:n}})}async function se(){try{let[e,t]=await Promise.all([j.get(`/api/v1/tenant/settings/competencies`,{params:{per_page:500}}),j.get(`/api/v1/tenant/job-management/values/clusters/managerial`)]);E.value=t.data?.data?.clusters||[];let n=new Set(E.value),r=e.data?.data||[];v.value=r.filter(e=>e.cluster&&n.has(e.cluster))}catch{v.value=[]}}async function Q(){try{let e=await j.get(`/api/v1/tenant/job-management/values`,{params:{type:`managerial`,per_page:100}});x.value=e.data?.data||[]}catch{x.value=[]}}async function ce(){if(!c.orgId){O.value=``,k.value=``;return}try{let e=(await j.get(`/api/v1/tenant/job-management/competency-groups`,{params:{organization_id:c.orgId}})).data?.data||[],t=e.find(e=>e.category===`technical`),n=e.find(e=>e.category===`managerial`);O.value=t?t.weight:``,k.value=n?n.id:``,A.value=n?n.weight:``}catch{O.value=``,k.value=``,A.value=``}}async function le({silent:e=!1}={}){let t=L.value;if(!(t===null||!c.orgId)){N.value=!0,P.value=``;try{let n=k.value;if(!n)try{let e=((await j.get(`/api/v1/tenant/job-management/competency-groups`,{params:{organization_id:c.orgId}})).data?.data||[]).find(e=>e.category===`managerial`);e&&(n=e.id)}catch{}let r={weight:t};n?await j.put(`/api/v1/tenant/job-management/competency-groups/${n}`,r):await j.post(`/api/v1/tenant/job-management/competency-groups`,{organization_id:c.orgId,category:`managerial`,weight:t}),k.value=n||``,A.value=t,e||f.add({severity:`success`,summary:u(`message.success`),detail:u(`job_management.potency_managerial_weight_saved`),life:2e3}),a(`saved`)}catch(e){P.value=e?.response?.data?.error?.message||e.message||u(`message.operation_failed`)}finally{N.value=!1}}}function ue(){let e={};(re.value||[]).forEach(t=>{e[t.id]=t});let t=[];W.value.forEach(n=>{n.competency_id&&e[n.competency_id]&&!t.includes(n.competency_id)&&t.push(n.competency_id)}),T.value=t,oe(),q()}return l(T,()=>{oe(),q()}),l(()=>c.technicalWeight,e=>{e!=null&&e!==``&&(O.value=e,le())}),n(async()=>{try{await Promise.all([se(),Q(),ee(),ce()])}finally{ue(),_.value=!1;let e=L.value;if(e!==null&&c.orgId){let t=A.value;(t===``||Math.abs(t-e)>.005)&&le({silent:!0})}}}),(n,i)=>(t(),m(`div`,Kr,[p(`div`,qr,[p(`div`,null,[p(`h3`,Jr,g(e(u)(`job_management.potency_managerial_title`)),1),p(`p`,Yr,g(e(u)(`job_management.potency_managerial_description`)),1)]),p(`div`,Xr,[p(`div`,Zr,[p(`label`,Qr,g(e(u)(`job_management.potency_managerial_weight_label`)),1),p(`div`,$r,[L.value===null?(t(),m(`span`,ti,`—`)):(t(),m(`span`,ei,g(L.value)+`%`,1))]),N.value?(t(),m(`i`,ni)):w(``,!0)]),P.value?(t(),m(`div`,ri,g(P.value),1)):w(``,!0)])]),_.value?(t(),b(Y,{key:0,type:`detail`,count:1,rows:5,cols:`grid-cols-1`,padding:`p-5`})):(t(),m(y,{key:1},[h(X,{modelValue:T.value,"onUpdate:modelValue":i[0]||=e=>T.value=e,options:ie.value,"option-label":`name`,"option-value":`id`,"option-group-label":`label`,"option-group-children":`items`,placeholder:e(u)(`job_management.potency_managerial_placeholder`),showClear:``,multiple:``},null,8,[`modelValue`,`options`,`placeholder`]),D.value?S.value.length===0?(t(),m(`div`,ai,g(e(u)(`job_management.potency_managerial_empty`)),1)):(t(),m(`div`,oi,[p(`table`,si,[p(`thead`,null,[p(`tr`,ci,[p(`th`,li,g(e(u)(`job_management.potency_table_name`)),1),p(`th`,ui,g(e(u)(`job_management.potency_table_level`)),1),p(`th`,di,g(e(u)(`job_management.potency_table_weight`)),1),p(`th`,fi,g(e(u)(`job_management.potency_table_description`)),1),p(`th`,pi,g(e(u)(`common.actions`)),1)])]),p(`tbody`,null,[(t(!0),m(y,null,s(S.value,n=>(t(),m(`tr`,{key:n.competency_id,class:`border-t border-gray-100 dark:border-gray-700 align-top`},[p(`td`,mi,[p(`div`,hi,g(n.competency_name),1),n.cluster?(t(),m(`div`,gi,g(n.cluster),1)):w(``,!0)]),p(`td`,_i,[h(X,{modelValue:n.job_management_value_id,"onUpdate:modelValue":e=>n.job_management_value_id=e,options:n.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:e(u)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),p(`td`,vi,[h(e(K),{modelValue:n.weight,"onUpdate:modelValue":e=>n.weight=e,class:`!w-full`,min:0,max:100,suffix:`%`,size:`small`},null,8,[`modelValue`,`onUpdate:modelValue`])]),p(`td`,yi,g(e(G)(n)),1),p(`td`,bi,[h(e(M),{icon:`pi pi-trash`,severity:`danger`,text:``,rounded:``,size:`small`,disabled:e(R),"aria-label":e(u)(`common.delete`),onClick:t=>e(te)(n)},null,8,[`disabled`,`aria-label`,`onClick`])])]))),128))])])])):(t(),m(`div`,ii,g(e(u)(`job_management.potency_managerial_no_mapping`)),1)),e(z)?(t(),m(`div`,xi,g(e(z)),1)):w(``,!0),S.value.length>0?(t(),m(`div`,Si,[h(e(M),{label:e(u)(`job_management.save_managerial`),icon:`pi pi-check`,size:`small`,loading:e(R),disabled:e(R)||!r.orgId,onClick:e(Z)},null,8,[`label`,`loading`,`disabled`,`onClick`])])):w(``,!0)],64)),h(J,{visible:e(B),"onUpdate:visible":i[1]||=e=>o(B)?B.value=e:null,title:e(u)(`job_management.potency_confirm_delete_title`),message:e(u)(`job_management.potency_confirm_delete`,{name:e(U)?.competency_name||``}),loading:e(V),"error-msg":e(H),onConfirm:e(ne),onCancel:i[2]||=e=>B.value=!1},null,8,[`visible`,`title`,`message`,`loading`,`error-msg`,`onConfirm`])]))}},wi={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Ti={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},Ei={class:`text-sm text-gray-500 dark:text-gray-400`},Di={key:0,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},Oi={key:1,class:`overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700`},ki={class:`w-full text-sm`},Ai={class:`bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400`},ji={class:`px-4 py-3 font-semibold min-w-[220px]`},Mi={class:`px-4 py-3 font-semibold min-w-[260px]`},Ni={class:`px-4 py-3 font-semibold min-w-[260px]`},Pi={class:`px-4 py-3 font-semibold w-16 text-right`},Fi={class:`px-4 py-3`},Ii={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},Li={key:0,class:`mt-0.5 text-xs text-gray-500 dark:text-gray-400 leading-relaxed`},Ri={class:`px-4 py-3`},zi={class:`px-4 py-3 text-sm text-gray-600 dark:text-gray-300`},Bi={class:`px-4 py-3 text-right`},Vi={key:2,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Hi={key:3,class:`flex justify-end gap-2 pt-1`},Ui={__name:`PotencyLevelsCard`,props:{orgId:String,rows:{type:Array,default:()=>[]},optionsReady:{type:Boolean,default:!1},skeletonRows:{type:Number,default:2},titleKey:{type:String,required:!0},descriptionKey:{type:String,required:!0},emptyKey:{type:String,required:!0},saveLabelKey:{type:String,required:!0},deleteTitleKey:{type:String,required:!0},deleteMessageKey:{type:String,required:!0}},emits:[`saved`],setup(n,{emit:r}){let i=r,a=n,{t:c}=I(),u=d(!0),f=C(()=>a.rows),{savingCard:_,errorMsg:v,deleteVisible:x,deleting:S,deleteError:T,deleteTarget:E,levelDescription:D,hydrateRows:O,loadData:k,askDeleteRow:A,handleDelete:j,handleSave:N}=Xn({orgId:C(()=>a.orgId),rows:f,onSaved:()=>i(`saved`)}),P=!1;return l(()=>a.optionsReady,async e=>{if(!(!e||P)){P=!0;try{await k()}finally{O(),u.value=!1}}},{immediate:!0}),(r,i)=>(t(),m(`div`,wi,[p(`div`,null,[p(`h3`,Ti,g(e(c)(n.titleKey)),1),p(`p`,Ei,g(e(c)(n.descriptionKey)),1)]),u.value?(t(),b(Y,{key:0,type:`detail`,count:1,rows:n.skeletonRows,cols:`grid-cols-1`,padding:`p-5`},null,8,[`rows`])):(t(),m(y,{key:1},[n.rows.length===0?(t(),m(`div`,Di,g(e(c)(n.emptyKey)),1)):(t(),m(`div`,Oi,[p(`table`,ki,[p(`thead`,null,[p(`tr`,Ai,[p(`th`,ji,g(e(c)(`job_management.potency_table_name`)),1),p(`th`,Mi,g(e(c)(`job_management.potency_table_level`)),1),p(`th`,Ni,g(e(c)(`job_management.potency_table_description`)),1),p(`th`,Pi,g(e(c)(`common.actions`)),1)])]),p(`tbody`,null,[(t(!0),m(y,null,s(n.rows,n=>(t(),m(`tr`,{key:n.type,class:`border-t border-gray-100 dark:border-gray-700 align-top`},[p(`td`,Fi,[p(`div`,Ii,g(n.competency_name),1),n.competency_definition?(t(),m(`div`,Li,g(n.competency_definition),1)):w(``,!0)]),p(`td`,Ri,[h(X,{modelValue:n.job_management_value_id,"onUpdate:modelValue":e=>n.job_management_value_id=e,options:n.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:e(c)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),p(`td`,zi,g(e(D)(n)),1),p(`td`,Bi,[h(e(M),{icon:`pi pi-trash`,severity:`danger`,text:``,rounded:``,size:`small`,disabled:e(_),"aria-label":e(c)(`common.delete`),onClick:t=>e(A)(n)},null,8,[`disabled`,`aria-label`,`onClick`])])]))),128))])])])),e(v)?(t(),m(`div`,Vi,g(e(v)),1)):w(``,!0),n.rows.length>0?(t(),m(`div`,Hi,[h(e(M),{label:e(c)(n.saveLabelKey),icon:`pi pi-check`,size:`small`,loading:e(_),disabled:e(_)||!n.orgId,onClick:e(N)},null,8,[`label`,`loading`,`disabled`,`onClick`])])):w(``,!0)],64)),h(J,{visible:e(x),"onUpdate:visible":i[0]||=e=>o(x)?x.value=e:null,title:e(c)(n.deleteTitleKey),message:e(c)(n.deleteMessageKey,{name:e(E)?.competency_name||``}),loading:e(S),"error-msg":e(T),onConfirm:e(j),onCancel:i[1]||=e=>x.value=!1},null,8,[`visible`,`title`,`message`,`loading`,`error-msg`,`onConfirm`])]))}},Wi={__name:`TypesPotencyCard`,props:{orgId:String,types:{type:Array,required:!0},skeletonRows:{type:Number,default:2},titleKey:{type:String,required:!0},descriptionKey:{type:String,required:!0},emptyKey:{type:String,required:!0},saveLabelKey:{type:String,required:!0},deleteTitleKey:{type:String,required:!0},deleteMessageKey:{type:String,required:!0}},emits:[`saved`],setup(e,{emit:r}){let i=r,a=e,{t:o}=I(),s=d([]),c=d(!1);function l(e){s.value=a.types.filter(t=>(e[t.type]||[]).length>0).map(t=>({competency_id:``,competency_name:o(t.nameKey),competency_definition:``,type:t.type,levelOptions:e[t.type]||[],recordId:``,job_management_value_id:``}))}async function u(){try{let e=await j.get(`/api/v1/tenant/job-management/values/tree`),t={};(e.data?.data||[]).forEach(e=>{(e.types||[]).forEach(e=>{a.types.some(t=>t.type===e.type)&&(t[e.type]=(e.options||[]).map(e=>({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``})))})}),l(t)}catch{s.value=[]}}return n(async()=>{await u(),c.value=!0}),(n,r)=>(t(),b(Ui,{"org-id":e.orgId,rows:s.value,"options-ready":c.value,"skeleton-rows":e.skeletonRows,"title-key":e.titleKey,"description-key":e.descriptionKey,"empty-key":e.emptyKey,"save-label-key":e.saveLabelKey,"delete-title-key":e.deleteTitleKey,"delete-message-key":e.deleteMessageKey,onSaved:r[0]||=e=>i(`saved`)},null,8,[`org-id`,`rows`,`options-ready`,`skeleton-rows`,`title-key`,`description-key`,`empty-key`,`save-label-key`,`delete-title-key`,`delete-message-key`]))}},Gi={__name:`ProblemSolvingPotencyCard`,props:{orgId:String},emits:[`saved`],setup(e,{emit:n}){let r=n,i=[{type:`thinking_environment`,nameKey:`job_management.problem_solving_environment`},{type:`thinking_chalenge`,nameKey:`job_management.problem_solving_challenge`}];return(n,a)=>(t(),b(Wi,{"org-id":e.orgId,types:i,"skeleton-rows":2,"title-key":`job_management.problem_solving_title`,"description-key":`job_management.problem_solving_description`,"empty-key":`job_management.problem_solving_empty`,"save-label-key":`job_management.save_problem_solving`,"delete-title-key":`job_management.problem_solving_confirm_delete_title`,"delete-message-key":`job_management.problem_solving_confirm_delete`,onSaved:a[0]||=e=>r(`saved`)},null,8,[`org-id`]))}},Ki={__name:`SkillPotencyCard`,props:{orgId:String},emits:[`saved`],setup(e,{emit:n}){let r=n,i=[{type:`communicating_influencing_skill`,nameKey:`job_management.skill_communicating_influencing`}];return(n,a)=>(t(),b(Wi,{"org-id":e.orgId,types:i,"skeleton-rows":2,"title-key":`job_management.skill_communicating_influencing_title`,"description-key":`job_management.skill_communicating_influencing_description`,"empty-key":`job_management.skill_communicating_influencing_empty`,"save-label-key":`job_management.save_skill`,"delete-title-key":`job_management.skill_confirm_delete_title`,"delete-message-key":`job_management.skill_confirm_delete`,onSaved:a[0]||=e=>r(`saved`)},null,8,[`org-id`]))}},qi={class:`space-y-4`},Ji={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Yi={class:`text-sm text-gray-500 dark:text-gray-400`},Xi={__name:`JobPotencySection`,props:{orgId:String,jobValueMap:{type:Object,default:()=>({})},competencyOptions:{type:Array,default:()=>[]}},emits:[`saved`],setup(n,{emit:r}){let i=r,{t:a}=I(),o=d(null);return(r,s)=>(t(),m(`div`,qi,[p(`div`,null,[p(`h2`,Ji,g(e(a)(`job_management.potency_competencies`)),1),p(`p`,Yi,g(e(a)(`job_management.potency_description`)),1)]),h(gr,{"org-id":n.orgId,onSaved:s[0]||=e=>i(`saved`)},null,8,[`org-id`]),h(Gr,{"org-id":n.orgId,onSaved:s[1]||=e=>i(`saved`),onWeightSaved:s[2]||=e=>o.value=e},null,8,[`org-id`]),h(Ci,{"org-id":n.orgId,"technical-weight":o.value,onSaved:s[3]||=e=>i(`saved`)},null,8,[`org-id`,`technical-weight`]),h(Gi,{"org-id":n.orgId,onSaved:s[4]||=e=>i(`saved`)},null,8,[`org-id`]),h(Ki,{"org-id":n.orgId,onSaved:s[5]||=e=>i(`saved`)},null,8,[`org-id`])]))}},Zi={class:`space-y-6`},Qi={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},$i={class:`text-sm text-gray-500 dark:text-gray-400`},ea={key:0,class:`flex items-center justify-center py-12`},ta={key:0,class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden`},na={class:`px-5 py-3 border-b border-gray-200 dark:border-gray-700 font-semibold text-sm text-gray-700 dark:text-gray-300`},ra={class:`divide-y divide-gray-100 dark:divide-gray-700`},ia={class:`hidden md:grid grid-cols-[minmax(0,2fr)_minmax(0,3fr)_auto] gap-4 px-5 py-2.5 bg-gray-50 dark:bg-gray-900/40 text-[11px] uppercase tracking-wider text-gray-400 dark:text-gray-500 font-medium`},aa={class:`text-right`},oa={class:`grid grid-cols-1 md:grid-cols-[minmax(0,2fr)_minmax(0,3fr)_auto] md:items-center gap-2`},sa={class:`text-sm font-semibold text-gray-700 dark:text-gray-200`},ca={class:`flex flex-wrap gap-1.5`},la={class:`font-medium`},ua={key:0,class:`font-mono`},da={class:`font-bold text-emerald-600 dark:text-emerald-400`},fa={class:`text-right`},pa={class:`text-sm font-bold text-gray-900 dark:text-gray-100`},ma={class:`px-5 py-4 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/40`},ha={class:`grid grid-cols-1 sm:grid-cols-2 gap-4`},ga={class:`flex items-center justify-between`},_a={class:`text-xs text-gray-500 dark:text-gray-400`},va={class:`text-sm font-bold text-emerald-600 dark:text-emerald-400`},ya={class:`flex items-center justify-between`},ba={class:`text-xs text-gray-500 dark:text-gray-400`},xa={class:`text-sm font-bold text-blue-600 dark:text-blue-400`},Sa={key:2},Ca={class:`flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500`},wa={class:`text-sm font-medium`},Ta={class:`text-xs mt-1`},Ea={class:`flex justify-end gap-3`},Da=`/api/v1/tenant/job-management/scores/org`,Oa={__name:`JobScoreSection`,props:{orgId:String},emits:[`saved`],setup(r,{emit:i}){let a=r,o=i,{t:c}=I(),l=F(),u=d(!1),f=d(!1),v=d(null),x=[{key:`education_experience`,labelKey:`job_management.education_experience`,points:[{labelKey:`job_management.group_education`,level:`education_level`,pts:`education_points`},{labelKey:`job_management.group_experience`,level:`experience_level`,pts:`experience_points`}]},{key:`potentials`,labelKey:`job_management.score_potentials`,points:[{labelKey:`job_management.average_level`,level:`average_level`,pts:null}]},{key:`competencies`,labelKey:`job_management.potency_competencies`,score:`base_score`,points:[{labelKey:`job_management.potency_technical_title`,level:`technical_average_level`,pts:`technical_points`},{labelKey:`job_management.potency_managerial_title`,level:`managerial_average_level`,pts:`managerial_points`},{labelKey:`job_management.skill_communicating_influencing`,level:`communication_level`,pts:`communication_points`}]},{key:`problem_solving`,labelKey:`job_management.problem_solving_title`,points:[{labelKey:`job_management.problem_solving_environment`,level:`environment_level`,pts:`environment_points`},{labelKey:`job_management.problem_solving_challenge`,level:`challenge_level`,pts:`challenge_points`}]},{key:`financial_authority`,labelKey:`job_management.financials`,points:[{labelKey:`job_management.cash_level`,level:`money_level`,pts:`money_points`},{labelKey:`job_management.authority_level`,level:`authority_level`,pts:`authority_points`},{labelKey:`job_management.impact_level`,level:`impact_level`,pts:`impact_points`}]},{key:`asset_authority`,labelKey:`job_management.assets`,points:[{labelKey:`job_management.asset_type`,level:`asset_value_level`,pts:`asset_value_points`},{labelKey:`job_management.authority_level`,level:`asset_authority_level`,pts:`asset_authority_points`}]},{key:`subordinate_control`,labelKey:`job_management.subordinate_controls`,points:[{labelKey:`job_management.score_level`,level:`level`,pts:`points`}]},{key:`work_scope`,labelKey:`job_management.relationships`,points:[{labelKey:`job_management.relationship_group_scope`,level:`scope_level`,pts:`scope_points`},{labelKey:`job_management.frequency`,level:`frequency_level`,pts:`frequency_points`}]},{key:`work_activity`,labelKey:`job_management.activities`,points:[{labelKey:`job_management.score_level`,level:`level`,pts:`points`}]},{key:`work_risk`,labelKey:`job_management.risks`,points:[{labelKey:`job_management.environment_risk`,level:`environment_level`,pts:`environment_points`},{labelKey:`job_management.hazard`,level:`hazard_level`,pts:`hazard_points`}]}],T=C(()=>{if(!v.value?.components)return null;try{return JSON.parse(v.value.components)}catch{return null}}),E=C(()=>T.value?x.map(e=>{let t=T.value[e.key]||{};return{key:e.key,labelKey:e.labelKey,score:t[e.score||`score`]??0,points:e.points.map(e=>({labelKey:e.labelKey,level:t[e.level]??null,points:e.pts==null?null:t[e.pts]??0}))}}):[]);function D(e){return e?.toLocaleString?.(`id-ID`)??`-`}function O(e){return e==null?`-`:String(e)}async function k(){if(a.orgId){u.value=!0;try{let e=await j.get(`${Da}/${a.orgId}`);v.value=e.data?.data||null,o(`saved`)}catch{v.value=null}finally{u.value=!1}}}async function A(){if(a.orgId){f.value=!0;try{let e=await j.put(`${Da}/${a.orgId}`,{components:null});v.value=e.data?.data||null,l.add({severity:`success`,summary:c(`message.success`),detail:c(`job_management.score_calculated`),life:2e3})}catch(e){l.add({severity:`error`,summary:c(`message.error`),detail:e.response?.data?.error?.message||c(`message.operation_failed`),life:4e3})}finally{f.value=!1}}}return n(k),(n,r)=>(t(),m(`div`,Zi,[p(`div`,null,[p(`h2`,Qi,g(e(c)(`job_management.scores`)),1),p(`p`,$i,g(e(c)(`job_management.score_description`)),1)]),u.value?(t(),m(`div`,ea,[...r[0]||=[p(`i`,{class:`pi pi-spin pi-spinner text-emerald-500 text-2xl`},null,-1)]])):v.value?(t(),m(y,{key:1},[E.value.length?(t(),m(`div`,ta,[p(`div`,na,g(e(c)(`job_management.component_breakdown`)),1),p(`div`,ra,[p(`div`,ia,[p(`span`,null,g(e(c)(`job_management.score_component`)),1),p(`span`,null,g(e(c)(`job_management.score_points`)),1),p(`span`,aa,g(e(c)(`job_management.score_score`)),1)]),(t(!0),m(y,null,s(E.value,n=>(t(),m(`div`,{key:n.key,class:`px-5 py-1`},[p(`div`,oa,[p(`div`,sa,g(e(c)(n.labelKey)),1),p(`div`,ca,[(t(!0),m(y,null,s(n.points,n=>(t(),m(`span`,{key:n.labelKey,class:S([`inline-flex items-center gap-1 text-[11px] px-2 py-0.5 rounded-md border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-900/40 text-gray-600 dark:text-gray-300`,{"opacity-50":n.level==null}])},[p(`span`,la,g(e(c)(n.labelKey)),1),n.level==null?w(``,!0):(t(),m(`span`,ua,`Lv.`+g(O(n.level)),1)),n.points==null?n.level==null?(t(),m(y,{key:2},[_(`—`)],64)):w(``,!0):(t(),m(y,{key:1},[r[1]||=p(`i`,{class:`pi pi-arrow-right text-[8px] opacity-60`},null,-1),p(`span`,da,g(D(n.points)),1)],64))],2))),128))]),p(`div`,fa,[p(`span`,pa,g(D(n.score)),1)])])]))),128))]),p(`div`,ma,[p(`div`,ha,[p(`div`,ga,[p(`span`,_a,g(e(c)(`job_management.value_with_financial`)),1),p(`span`,va,g(D(v.value.job_value_with_financial)),1)]),p(`div`,ya,[p(`span`,ba,g(e(c)(`job_management.value_without_financial`)),1),p(`span`,xa,g(D(v.value.job_value_without_financial)),1)])])])])):w(``,!0)],64)):(t(),m(`div`,Sa,[p(`div`,Ca,[r[2]||=p(`i`,{class:`pi pi-calculator text-4xl mb-3 opacity-50`},null,-1),p(`p`,wa,g(e(c)(`job_management.no_score`)),1),p(`p`,Ta,g(e(c)(`job_management.score_hint`)),1)])])),p(`div`,Ea,[h(e(M),{label:e(c)(`common.refresh`),icon:`pi pi-refresh`,size:`small`,text:``,onClick:k},null,8,[`label`]),v.value?(t(),b(e(M),{key:0,label:e(c)(`job_management.recalculate`),icon:`pi pi-calculator`,size:`small`,severity:`info`,loading:f.value,onClick:A},null,8,[`label`,`loading`])):w(``,!0)])]))}},ka={class:`grid grid-cols-1 md:grid-cols-3 gap-4`},Aa={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4`},ja={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},Ma={class:`text-2xl font-bold text-emerald-600 dark:text-emerald-400`},Na={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4`},Pa={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},Fa={class:`text-2xl font-bold text-blue-600 dark:text-blue-400`},Ia={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4`},La={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},Ra={class:`mt-3 flex flex-wrap items-center gap-2`},za={key:1,class:`text-[10px] text-gray-400`},Ba={key:0,class:`text-[10px] text-gray-400 mt-2`},Va=`/api/v1/tenant/job-management/scores/org`,Ha={__name:`JobScoreSummary`,props:{orgId:String},setup(r,{expose:i}){let a=r,{t:o}=I(),c=d(!0),u=d(null);function f(e){return e?.toLocaleString?.(`id-ID`)??`-`}async function _(){if(a.orgId){c.value=!0;try{let e=await j.get(`${Va}/${a.orgId}`);u.value=e.data?.data||null}catch{u.value=null}finally{c.value=!1}}}return i({refresh:_}),l(()=>a.orgId,_),n(_),(n,r)=>(t(),m(`div`,ka,[c.value&&!u.value?(t(),m(y,{key:0},s(3,e=>p(`div`,{key:e,class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4 animate-pulse`},[...r[0]||=[p(`div`,{class:`h-3 w-24 bg-gray-200 dark:bg-gray-700 rounded mb-2`},null,-1),p(`div`,{class:`h-7 w-16 bg-gray-200 dark:bg-gray-700 rounded`},null,-1)]])),64)):(t(),m(y,{key:1},[p(`div`,Aa,[p(`div`,ja,g(e(o)(`job_management.value_with_financial`)),1),p(`div`,Ma,g(f(u.value?.job_value_with_financial)),1)]),p(`div`,Na,[p(`div`,Pa,g(e(o)(`job_management.value_without_financial`)),1),p(`div`,Fa,g(f(u.value?.job_value_without_financial)),1)]),p(`div`,Ia,[p(`div`,La,g(e(o)(`job_management.has_financial_authority`)),1),h(e(R),{value:u.value?u.value.has_financial_authority?e(o)(`common.yes`):e(o)(`common.no`):`-`,severity:u.value?.has_financial_authority?`success`:`danger`,class:`!text-xs`},null,8,[`value`,`severity`]),p(`div`,Ra,[u.value?(t(),b(e(R),{key:0,value:u.value.is_complete?e(o)(`job_management.score_complete`):e(o)(`job_management.score_incomplete`),severity:u.value.is_complete?`success`:`warning`,icon:u.value.is_complete?`pi pi-check-circle`:`pi pi-exclamation-triangle`,class:`!text-xs`},null,8,[`value`,`severity`,`icon`])):w(``,!0),u.value?.is_complete&&u.value.completed_at?(t(),m(`span`,za,g(e(o)(`job_management.completed_at`))+`: `+g(u.value.completed_at),1)):w(``,!0)]),u.value?.calculated_at?(t(),m(`div`,Ba,g(e(o)(`job_management.calculated_at`))+`: `+g(u.value.calculated_at),1)):w(``,!0)])],64))]))}},Ua={class:`max-w-full mx-auto`},Wa={key:0,class:`flex gap-6`},Ga={class:`w-56 space-y-2`},Ka={class:`flex-1 space-y-3`},qa={key:1,class:`flex gap-6`},Ja={class:`w-56 shrink-0 space-y-1`},Ya=[`onClick`,`onKeydown`],Xa={key:0,class:`pi pi-check text-xs`},Za={class:`flex-1 min-w-0`},Qa={key:0,class:`pi pi-check-circle text-emerald-400 text-xs shrink-0`},$a={class:`flex-1 min-w-0 space-y-4`},eo={class:`flex flex-col md:flex-row gap-4`},to={class:`md:w-72 shrink-0 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4`},no={class:`flex items-center gap-2 mb-3`},ro={class:`text-sm font-semibold text-gray-800 dark:text-gray-100 truncate`},io={class:`flex items-center justify-between gap-2`},ao={class:`text-[10px] uppercase tracking-wider text-gray-400 dark:text-gray-500`},oo={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 font-mono truncate`},so={class:`flex-1 min-w-0`},co={__name:`JobManagementForm`,setup(r){let i=k(),o=A(),{t:c}=I(),u=F(),f=C(()=>o.query.org_id||``),_=d(0),v=d(!0),x=d(Array(15).fill(!1)),T=d(``),E=d(``),D=d(``),M=d(``),N=d(``),P=d([]),L=d([]),R=d([]),z=d({}),B=d([]),V=d(null),H=[{labelKey:`job_management.identifications`,icon:`pi pi-id-card`,comp:ce},{labelKey:`job_management.objectives`,icon:`pi pi-bullseye`,comp:_e},{labelKey:`job_management.responsibilities_title`,icon:`pi pi-list-check`,comp:st},{labelKey:`job_management.education_experience`,icon:`pi pi-graduation-cap`,comp:Pe},{labelKey:`job_management.potency_competencies`,icon:`pi pi-star`,comp:Xi},{labelKey:`job_management.financials`,icon:`pi pi-money-bill`,comp:Jn},{labelKey:`job_management.assets`,icon:`pi pi-box`,comp:Fn},{labelKey:`job_management.subordinate_controls`,icon:`pi pi-sitemap`,comp:En},{labelKey:`job_management.relationships`,icon:`pi pi-share-alt`,comp:_n},{labelKey:`job_management.activities`,icon:`pi pi-bolt`,comp:Mt},{labelKey:`job_management.risks`,icon:`pi pi-exclamation-triangle`,comp:Vt},{labelKey:`job_management.hr_authorities`,icon:`pi pi-users`,comp:ht},{labelKey:`job_management.op_authorities`,icon:`pi pi-cog`,comp:Ct},{labelKey:`job_management.scores`,icon:`pi pi-calculator`,comp:Oa}],U=C(()=>H[_.value]?.comp||null);function W(e){return _.value===e?`bg-emerald-50 dark:bg-emerald-900/20 ring-1 ring-emerald-300 dark:ring-emerald-700`:(x.value[e],`hover:bg-gray-50 dark:hover:bg-gray-800`)}function G(e){return _.value===e?`bg-emerald-600 text-white`:x.value[e]?`bg-emerald-100 dark:bg-emerald-800 text-emerald-600 dark:text-emerald-300`:`bg-gray-200 dark:bg-gray-600 text-gray-600 dark:text-gray-300`}function K(e){return _.value===e?`text-emerald-700 dark:text-emerald-300`:x.value[e]?`text-emerald-600 dark:text-emerald-400`:`text-gray-700 dark:text-gray-300`}function q(e){_.value=e,i.replace({query:{...o.query,section:String(e)}})}function ee(e){typeof e==`number`&&(x.value[e]=!0),V.value?.refresh()}async function te(){if(f.value)try{let e=(await j.get(`/api/v1/tenant/organizations/${f.value}`)).data?.data;e&&(T.value=e.nomenclature||``,E.value=e.full_code||e.code||``,D.value=e.organization_summary_id||``,M.value=e.grading_id||``,N.value=e.job_family_id||``)}catch{}}async function J(){try{let[e,t,n,r]=await Promise.all([j.get(`/api/v1/tenant/settings/gradings?per_page=100`),j.get(`/api/v1/tenant/job-management/values?per_page=500`),j.get(`/api/v1/tenant/competency/competencies?per_page=200`).catch(()=>({data:{data:[]}})),j.get(`/api/v1/tenant/settings/job-families?per_page=100`)]);P.value=(e.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id})),L.value=(r.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id}));let i=t.data?.data||[];R.value=i.map(e=>({label:`${e.type}${e.level?` Lv.`+e.level:``}${e.descriptions?` — `+e.descriptions:``}`,value:e.id,type:e.type,level:e.level,descriptions:e.descriptions}));let a={};i.forEach(e=>{a[e.type]||(a[e.type]=[]),a[e.type].push({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level,descriptions:e.descriptions||``,type_group:e.type_group||``,description_group:e.description_group||``})}),z.value=a,B.value=(n.data?.data||[]).map(e=>({label:e.name||e.code,value:e.id,field:e.field||``,definition:e.definition||``}))}catch{}}return l(f,(e,t)=>{e!==t&&(x.value=Array(H.length).fill(!1),T.value=``,E.value=``,D.value=``,M.value=``,N.value=``,te())}),n(async()=>{try{await Promise.all([te(),J()]);let e=parseInt(o.query.section);!isNaN(e)&&e>=0&&e<H.length&&(_.value=e)}catch(e){u.add({severity:`error`,summary:c(`message.error`),detail:e.response?.data?.error?.message||c(`message.failed_to_load`),life:4e3})}finally{v.value=!1}}),(n,r)=>(t(),m(`div`,Ua,[v.value?(t(),m(`div`,Wa,[p(`div`,Ga,[(t(),m(y,null,s(8,e=>p(`div`,{key:e,class:`h-12 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))]),p(`div`,Ka,[(t(),m(y,null,s(6,e=>p(`div`,{key:e,class:`h-9 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))])])):(t(),m(`div`,qa,[p(`div`,Ja,[(t(),m(y,null,s(H,(n,r)=>p(`div`,{key:r,role:`button`,tabindex:0,class:S([`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150 cursor-pointer select-none`,W(r)]),onClick:e=>q(r),onKeydown:O(e=>q(r),[`enter`])},[p(`div`,{class:S([`w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0 transition-colors duration-150`,G(r)])},[x.value[r]?(t(),m(`i`,Xa)):(t(),m(`i`,{key:1,class:S(n.icon)},null,2))],2),p(`div`,Za,[p(`div`,{class:S([`text-sm font-medium truncate`,K(r)])},g(e(c)(n.labelKey)),3)]),x.value[r]?(t(),m(`i`,Qa)):w(``,!0)],42,Ya)),64))]),p(`div`,$a,[(t(),m(`div`,{key:`summary-${f.value}`,class:`sticky top-0 z-10 bg-white dark:bg-gray-900 pt-1 pb-3`},[p(`div`,eo,[p(`div`,to,[p(`div`,no,[r[0]||=p(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400`},[p(`i`,{class:`pi pi-briefcase text-sm`})],-1),p(`h3`,ro,g(T.value||e(c)(`job_management.job_info_untitled`)),1)]),p(`div`,io,[p(`span`,ao,g(e(c)(`organization.full_code`)),1),p(`span`,oo,g(E.value||`-`),1)])]),p(`div`,so,[h(Ha,{ref_key:`scoreSummaryRef`,ref:V,"org-id":f.value},null,8,[`org-id`])])])])),(t(),b(a(U.value),{key:`${_.value}-${f.value}`,"org-id":f.value,"org-name":T.value,"org-code":E.value,"org-summary-id":D.value,"org-grading-id":M.value,"org-job-family-id":N.value,"grading-options":P.value,"job-family-options":L.value,"job-value-options":R.value,"competency-options":B.value,"job-value-map":z.value,onSaved:ee},null,40,[`org-id`,`org-name`,`org-code`,`org-summary-id`,`org-grading-id`,`org-job-family-id`,`grading-options`,`job-family-options`,`job-value-options`,`competency-options`,`job-value-map`]))])]))]))}};export{co as default};