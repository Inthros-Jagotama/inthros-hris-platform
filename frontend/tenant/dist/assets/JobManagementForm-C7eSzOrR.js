const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["assets/quill-BLmY9xB4.js","assets/rolldown-runtime-QTnfLwEv.js"])))=>i.map(i=>d[i]);
import{A as e,C as t,E as n,H as r,K as i,N as a,O as o,P as s,U as c,V as l,X as u,c as d,dt as f,ft as p,h as m,j as h,l as g,m as _,o as v,pt as y,r as b,s as x,u as S,ut as C,y as w}from"./runtime-core.esm-bundler-CFm0BMYx.js";import{a as T,d as E,t as D}from"./button-DVe0lZgG.js";import{A as O,a as k}from"./ripple-FLjJmYYY.js";import{l as A,m as j,s as M,t as N,u as P}from"./index-B7i-gbTX.js";import{t as F}from"./useI18n-CZ6Bzqy-.js";import{n as I}from"./responseHandler-B5MnXl3B.js";import{t as L}from"./tag--2hRZhYy.js";import{t as R}from"./FormRow-D-Jmd1__.js";import{t as z}from"./baseeditableholder-CnB2pcI8.js";import{t as B}from"./textarea-CZKEJ-vM.js";import{t as V}from"./TextInput-DT_yhCva.js";import{a as H,n as U,s as W,t as G}from"./column-CTXbOJTC.js";import{t as K}from"./SelectLabel-7705891Z.js";import{t as q}from"./ConfirmDeleteDialog-DdASfghu.js";import{t as J}from"./SkeletonTable-BS8Q7JD0.js";import{t as ee}from"./toggleswitch-BSOZ6l02.js";import{t as te}from"./multiselect-CyXfkb3k.js";var ne={class:`flex items-end gap-1 h-16 mb-2`},Y={__name:`SkeletonCard`,props:{type:{type:String,default:`kpi`},count:{type:Number,default:null},cols:{type:String,default:null},rows:{type:Number,default:4},padding:{type:String,default:`p-3`},valueWidth:{type:String,default:``},labelWidth:{type:String,default:``}},setup(t){let n=t,r=v(()=>n.count!==null&&n.count!==void 0?n.count:{kpi:6,stat:4,metric:4,alert:3,sparkline:4,detail:4}[n.type]||4),i=v(()=>n.cols?n.cols:{kpi:`grid-cols-2 md:grid-cols-4 lg:grid-cols-6`,stat:`grid-cols-2 md:grid-cols-4`,metric:`grid-cols-1`,alert:`grid-cols-1`,sparkline:`grid-cols-1`,detail:`grid-cols-1`}[n.type]||`grid-cols-1`),a=[45,60,35,70,50,55,40,65];return(n,s)=>(o(),S(`div`,{class:C([`grid gap-3 animate-pulse`,i.value])},[t.type===`kpi`?(o(!0),S(b,{key:0},e(r.value,e=>(o(),S(`div`,{key:e,class:C([t.padding,`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700`])},[...s[0]||=[x(`div`,{class:`w-8 h-8 bg-gray-200 dark:bg-gray-600 rounded-lg mb-2`},null,-1),x(`div`,{class:`h-5 w-3/4 bg-gray-200 dark:bg-gray-600 rounded mb-1`},null,-1),x(`div`,{class:`h-3 w-1/2 bg-gray-200 dark:bg-gray-600 rounded`},null,-1)]],2))),128)):t.type===`stat`?(o(!0),S(b,{key:1},e(r.value,e=>(o(),S(`div`,{key:e,class:C([t.padding,`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700`])},[...s[1]||=[x(`div`,{class:`h-3 w-1/3 bg-gray-200 dark:bg-gray-600 rounded mb-2`},null,-1),x(`div`,{class:`h-6 w-2/3 bg-gray-200 dark:bg-gray-600 rounded`},null,-1)]],2))),128)):t.type===`metric`?(o(!0),S(b,{key:2},e(r.value,e=>(o(),S(`div`,{key:e,class:C([t.padding,`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700`])},[...s[2]||=[x(`div`,{class:`h-8 w-1/3 bg-gray-200 dark:bg-gray-600 rounded mb-1`},null,-1),x(`div`,{class:`h-3 w-1/2 bg-gray-200 dark:bg-gray-600 rounded mb-2`},null,-1),x(`div`,{class:`h-3 w-1/4 bg-gray-200 dark:bg-gray-600 rounded`},null,-1)]],2))),128)):t.type===`alert`?(o(!0),S(b,{key:3},e(r.value,e=>(o(),S(`div`,{key:e,class:C([t.padding,`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 flex items-start gap-2`])},[...s[3]||=[x(`div`,{class:`w-8 h-8 bg-gray-200 dark:bg-gray-600 rounded-full shrink-0`},null,-1),x(`div`,{class:`flex-1 space-y-1.5`},[x(`div`,{class:`h-3 w-3/4 bg-gray-200 dark:bg-gray-600 rounded`}),x(`div`,{class:`h-3 w-1/2 bg-gray-200 dark:bg-gray-600 rounded`})],-1)]],2))),128)):t.type===`sparkline`?(o(!0),S(b,{key:4},e(r.value,n=>(o(),S(`div`,{key:n,class:C([t.padding,`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700`])},[x(`div`,ne,[(o(),S(b,null,e(a,e=>x(`div`,{key:e,class:`flex-1 bg-gray-200 dark:bg-gray-600 rounded-t`,style:p({height:e+`%`})},null,4)),64))]),s[4]||=x(`div`,{class:`h-3 w-1/3 bg-gray-200 dark:bg-gray-600 rounded`},null,-1)],2))),128)):t.type===`detail`?(o(!0),S(b,{key:5},e(r.value,n=>(o(),S(`div`,{key:n,class:C([t.padding,`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700`])},[s[6]||=x(`div`,{class:`h-4 w-1/3 bg-gray-200 dark:bg-gray-600 rounded mb-3`},null,-1),(o(!0),S(b,null,e(t.rows,e=>(o(),S(`div`,{key:e,class:`flex items-center justify-between py-1.5`},[...s[5]||=[x(`div`,{class:`h-3 w-1/4 bg-gray-200 dark:bg-gray-600 rounded`},null,-1),x(`div`,{class:`h-3 w-1/3 bg-gray-200 dark:bg-gray-600 rounded`},null,-1)]]))),128))],2))),128)):g(``,!0)],2))}},re={class:`space-y-4`},X={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},ie={class:`text-sm text-gray-500 dark:text-gray-400`},ae={class:`max-w-2xl`},oe={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},se={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},ce={class:`flex justify-end pt-2`},le=`/api/v1/tenant/job-management/identifications`,ue={__name:`JobIdentificationSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},orgGradingId:{type:String,default:``},orgJobFamilyId:{type:String,default:``},gradingOptions:{type:Array,default:()=>[]},jobFamilyOptions:{type:Array,default:()=>[]}},emits:[`saved`],setup(e,{emit:t}){let a=t,s=e,{t:c}=F(),l=j(),f=i(!1),p=i(!0),h=i(``),_=i({}),b=i(``),C=i({grading_id:``}),w=v(()=>{let e=s.jobFamilyOptions.find(e=>e.value===s.orgJobFamilyId);return e?e.label:s.orgJobFamilyId||`-`});function E(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function O(){if(!s.orgId){p.value=!1;return}try{let e=(await T.get(le,{params:{organization_id:s.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];b.value=t.id,C.value.grading_id=t.grading_id||s.orgGradingId||``}else C.value.grading_id=s.orgGradingId||``}catch{C.value.grading_id=s.orgGradingId||``}finally{p.value=!1}}async function k(){if(h.value=``,_.value={},!C.value.grading_id){h.value=c(`job_management.grading_required`);return}f.value=!0;try{let e={nomenclature:s.orgName||``,full_code:s.orgCode||``,grading_id:C.value.grading_id,organization_id:s.orgId};if(b.value)await T.put(`${le}/${b.value}`,{grading_id:C.value.grading_id});else{let t=await T.post(le,e);b.value=t.data?.data?.id||``}l.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=E(e);Object.keys(t).length>0?(_.value=t,h.value=Object.values(t).join(`, `)):h.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{f.value=!1}}return n(O),(t,n)=>(o(),S(`div`,re,[x(`div`,null,[x(`h2`,X,y(u(c)(`job_management.identifications`)),1),x(`p`,ie,y(u(c)(`job_management.identification_description`)),1)]),x(`div`,ae,[p.value?(o(),d(Y,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(o(),S(`div`,oe,[m(R,{label:u(c)(`organization.nomenclature`)},{default:r(()=>[m(V,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`organization.full_code`)},{default:r(()=>[m(V,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`organization.job_family`)},{default:r(()=>[m(V,{"model-value":w.value,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`organization.grading`)},{default:r(()=>[m(u(W),{modelValue:C.value.grading_id,"onUpdate:modelValue":n[0]||=e=>C.value.grading_id=e,options:e.gradingOptions,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`organization.select_grading`),class:`w-full`,size:`small`,invalid:!!_.value.grading_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`]),h.value?(o(),S(`div`,se,y(h.value),1)):g(``,!0),x(`div`,ce,[m(u(D),{label:u(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:f.value,disabled:!C.value.grading_id,onClick:k},null,8,[`label`,`loading`,`disabled`])])]))])]))}},de={class:`space-y-4`},fe={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},pe={class:`text-sm text-gray-500 dark:text-gray-400`},me={class:`max-w-2xl`},he={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},ge={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},_e={class:`flex justify-end gap-2 pt-2`},ve=`/api/v1/tenant/job-management/objectives`,ye={__name:`JobObjectiveSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let a=t,s=e,{t:c}=F(),l=j(),f=i(!1),p=i(!0),h=i(!1),_=i(``),v=i({}),b=i(``),w=i(!1),E=i(``),O=i({objective:``});function k(e){let t=e?.response?.data?.error?.fields;if(t&&typeof t==`object`){let e={};for(let[n,r]of Object.entries(t))e[n]=Array.isArray(r)?r[0]:r;return e}return{}}async function A(){if(!s.orgId){p.value=!1;return}try{let e=(await T.get(ve,{params:{organization_id:s.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];b.value=t.id,O.value.objective=t.objective||``}}catch{}finally{p.value=!1}}async function M(){_.value=``,v.value={},f.value=!0;try{let e={nomenclature:s.orgName||``,full_code:s.orgCode||``,objective:O.value.objective||``,organization_id:s.orgId};if(b.value)await T.put(`${ve}/${b.value}`,{objective:O.value.objective||``});else{let t=await T.post(ve,e);b.value=t.data?.data?.id||``}l.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=k(e);Object.keys(t).length>0?(v.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{f.value=!1}}async function N(){if(b.value){h.value=!0,E.value=``;try{await T.delete(`${ve}/${b.value}`),w.value=!1,b.value=``,O.value.objective=``,a(`saved`),l.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){E.value=e?.response?.data?.error?.message||c(`message.operation_failed`)}finally{h.value=!1}}}return n(A),(t,n)=>(o(),S(`div`,de,[x(`div`,null,[x(`h2`,fe,y(u(c)(`job_management.objectives`)),1),x(`p`,pe,y(u(c)(`job_management.objective_description`)),1)]),x(`div`,me,[p.value?(o(),d(Y,{key:0,type:`detail`,count:1,rows:3,cols:`grid-cols-1`,padding:`p-5`})):(o(),S(`div`,he,[m(R,{label:u(c)(`organization.nomenclature`)},{default:r(()=>[m(V,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`organization.full_code`)},{default:r(()=>[m(V,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`job_management.objective`)},{default:r(()=>[m(u(B),{modelValue:O.value.objective,"onUpdate:modelValue":n[0]||=e=>O.value.objective=e,rows:`3`,class:C([`w-full`,{"p-invalid":v.value.objective}]),placeholder:u(c)(`job_management.objective`)+`...`},null,8,[`modelValue`,`class`,`placeholder`])]),_:1},8,[`label`]),_.value?(o(),S(`div`,ge,y(_.value),1)):g(``,!0),x(`div`,_e,[b.value?(o(),d(u(D),{key:0,label:u(c)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:n[1]||=e=>w.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:b.value?u(c)(`common.update`):u(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:f.value,disabled:f.value,onClick:M},null,8,[`label`,`loading`,`disabled`])])]))]),m(q,{visible:w.value,"onUpdate:visible":n[2]||=e=>w.value=e,loading:h.value,"error-msg":E.value,onConfirm:N,onCancel:n[3]||=e=>w.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},be={class:`space-y-4`},xe={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Se={class:`text-sm text-gray-500 dark:text-gray-400`},Ce={class:`max-w-2xl`},we={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Te={class:`pt-4 border-t border-gray-200 dark:border-gray-700`},Ee={class:`flex items-center gap-2 mb-3`},De={class:`text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400`},Oe={class:`space-y-4`},ke={class:`pt-4 border-t border-gray-200 dark:border-gray-700`},Ae={class:`flex items-center gap-2 mb-3`},je={class:`text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400`},Me={class:`space-y-4`},Ne={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Pe={class:`flex justify-end gap-2 pt-2`},Fe=`/api/v1/tenant/job-management/education-experiences`,Ie={__name:`JobEduExpSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:t}){let a=t,s=e,{t:c}=F(),l=j(),f=i(!1),p=i(!0),h=i(!1),_=i(``),v=i({}),b=i(``),w=i(!1),E=i(``),O=i({education_id:``,education_major_id:[],job_family_id:[],experience_id:``}),k=i([]),A=i([]),M=i([]),N=i([]);async function P(){try{let[e,t,n,r]=await Promise.all([T.get(`/api/v1/tenant/job-management/values`,{params:{type:`education`,per_page:100}}),T.get(`/api/v1/tenant/job-management/values`,{params:{type:`experience`,per_page:100}}),T.get(`/api/v1/tenant/settings/education-majors?per_page=200`),T.get(`/api/v1/tenant/settings/job-families?per_page=100`)]);A.value=(e.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),k.value=(t.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),M.value=(n.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id})),N.value=(r.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id}))}catch{}}async function L(){if(s.orgId)try{let e=(await T.get(Fe,{params:{organization_id:s.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];b.value=t.id,O.value.education_id=t.education_id||``,O.value.education_major_id=Array.isArray(t.education_major_id)?t.education_major_id:[],O.value.job_family_id=Array.isArray(t.job_family_id)?t.job_family_id:[],O.value.experience_id=t.experience_id||``}}catch{}}async function z(){_.value=``,v.value={},f.value=!0;try{let e={nomenclature:s.orgName||``,full_code:s.orgCode||``,education_id:O.value.education_id||null,education_major_id:O.value.education_major_id||[],job_family_id:O.value.job_family_id||[],experience_id:O.value.experience_id||null,organization_id:s.orgId};if(b.value)await T.put(`${Fe}/${b.value}`,{education_id:O.value.education_id||``,education_major_id:O.value.education_major_id||[],job_family_id:O.value.job_family_id||[],experience_id:O.value.experience_id||``});else{let t=await T.post(Fe,e);b.value=t.data?.data?.id||``}l.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=I(e);Object.keys(t).length>0?(v.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{f.value=!1}}async function B(){if(b.value){h.value=!0,E.value=``;try{await T.delete(`${Fe}/${b.value}`),w.value=!1,b.value=``,O.value.education_id=``,O.value.education_major_id=[],O.value.job_family_id=[],O.value.experience_id=``,a(`saved`),l.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){E.value=e?.response?.data?.error?.message||c(`message.operation_failed`)}finally{h.value=!1}}}return n(async()=>{try{await Promise.all([P(),L()])}finally{p.value=!1}}),(t,n)=>(o(),S(`div`,be,[x(`div`,null,[x(`h2`,xe,y(u(c)(`job_management.education_experience`)),1),x(`p`,Se,y(u(c)(`job_management.education_experience_description`)),1)]),x(`div`,Ce,[p.value?(o(),d(Y,{key:0,type:`detail`,count:1,rows:6,cols:`grid-cols-1`,padding:`p-5`})):(o(),S(`div`,we,[m(R,{label:u(c)(`organization.nomenclature`)},{default:r(()=>[m(V,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`organization.full_code`)},{default:r(()=>[m(V,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),x(`div`,Te,[x(`div`,Ee,[n[7]||=x(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400`},[x(`i`,{class:`pi pi-graduation-cap text-sm`})],-1),x(`h3`,De,y(u(c)(`job_management.group_education`)),1),n[8]||=x(`div`,{class:`flex-1 border-t border-gray-200 dark:border-gray-700`},null,-1)]),x(`div`,Oe,[m(R,{label:u(c)(`job_management.education_level`),errors:v.value?.education_id},{default:r(()=>[m(K,{modelValue:O.value.education_id,"onUpdate:modelValue":n[0]||=e=>O.value.education_id=e,options:A.value,placeholder:u(c)(`job_values.select_education`),class:C({"p-invalid":v.value?.education_id})},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),m(R,{label:u(c)(`job_management.education_major`),errors:v.value?.education_major_id},{default:r(()=>[m(u(te),{modelValue:O.value.education_major_id,"onUpdate:modelValue":n[1]||=e=>O.value.education_major_id=e,options:M.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),class:`w-full`,size:`small`,filter:``,showClear:``,display:`chip`,maxSelectedLabels:2,invalid:!!v.value.education_major_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`])])]),x(`div`,ke,[x(`div`,Ae,[n[9]||=x(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400`},[x(`i`,{class:`pi pi-briefcase text-sm`})],-1),x(`h3`,je,y(u(c)(`job_management.group_experience`)),1),n[10]||=x(`div`,{class:`flex-1 border-t border-gray-200 dark:border-gray-700`},null,-1)]),x(`div`,Me,[m(R,{label:u(c)(`job_management.experience_range`),errors:v.value?.experience_id},{default:r(()=>[m(K,{modelValue:O.value.experience_id,"onUpdate:modelValue":n[2]||=e=>O.value.experience_id=e,options:k.value,placeholder:u(c)(`common.select`),class:C({"p-invalid":v.value?.experience_id})},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),m(R,{label:u(c)(`job_management.job_family`),errors:v.value?.job_family_id},{default:r(()=>[m(u(te),{modelValue:O.value.job_family_id,"onUpdate:modelValue":n[3]||=e=>O.value.job_family_id=e,options:N.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),class:`w-full`,size:`small`,filter:``,showClear:``,display:`chip`,maxSelectedLabels:2,invalid:!!v.value.job_family_id},null,8,[`modelValue`,`options`,`placeholder`,`invalid`])]),_:1},8,[`label`,`errors`])])]),_.value?(o(),S(`div`,Ne,y(_.value),1)):g(``,!0),x(`div`,Pe,[b.value?(o(),d(u(D),{key:0,label:u(c)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:n[4]||=e=>w.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:b.value?u(c)(`common.update`):u(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:f.value,disabled:f.value,onClick:z},null,8,[`label`,`loading`,`disabled`])])]))]),m(q,{visible:w.value,"onUpdate:visible":n[5]||=e=>w.value=e,loading:h.value,"error-msg":E.value,onConfirm:B,onCancel:n[6]||=e=>w.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Le=k.extend({name:`editor`,style:`
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
`,classes:{root:function(e){return[`p-editor`,{"p-invalid":e.instance.$invalid}]},toolbar:`p-editor-toolbar`,content:`p-editor-content`}}),Re={name:`BaseEditor`,extends:z,props:{placeholder:String,readonly:Boolean,formats:Array,editorStyle:null,modules:null},style:Le,provide:function(){return{$pcEditor:this,$parentInstance:this}}};function Z(e){"@babel/helpers - typeof";return Z=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},Z(e)}function ze(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function Be(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?ze(Object(n),!0).forEach(function(t){Ve(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):ze(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function Ve(e,t,n){return(t=He(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function He(e){var t=Ue(e,`string`);return Z(t)==`symbol`?t:t+``}function Ue(e,t){if(Z(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(Z(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var We=function(){try{return window.Quill}catch{return null}}(),Ge={name:`Editor`,extends:Re,inheritAttrs:!1,emits:[`text-change`,`selection-change`,`load`],quill:null,watch:{modelValue:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},d_value:function(e,t){e!==t&&this.quill&&!this.quill.hasFocus()&&this.renderValue(e)},readonly:function(){this.handleReadOnlyChange()}},mounted:function(){var e=this,t={modules:Be({toolbar:this.$refs.toolbarElement},this.modules),readOnly:this.readonly,theme:`snow`,formats:this.formats,placeholder:this.placeholder};We?(this.quill=new We(this.$refs.editorElement,t),this.initQuill(),this.handleLoad()):M(()=>import(`./quill-BLmY9xB4.js`).then(function(n){n&&O(e.$refs.editorElement)&&(n.default?e.quill=new n.default(e.$refs.editorElement,t):e.quill=new n(e.$refs.editorElement,t),e.initQuill())}),__vite__mapDeps([0,1])).then(function(){e.handleLoad()})},beforeUnmount:function(){this.quill=null},methods:{renderValue:function(e){if(this.quill)if(e){var t=this.quill.clipboard.convert({html:e});this.quill.setContents(t)}else this.quill.setText(``)},initQuill:function(){var e=this;this.renderValue(this.d_value),this.quill.on(`text-change`,function(t,n,r){if(r===`user`){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();i===`<p><br></p>`&&(i=``),e.writeValue(i),e.$emit(`text-change`,{htmlValue:i,textValue:a,delta:t,source:r,instance:e.quill})}}),this.quill.on(`selection-change`,function(t,n,r){var i=e.quill.getSemanticHTML(),a=e.quill.getText().trim();e.$emit(`selection-change`,{htmlValue:i,textValue:a,range:t,oldRange:n,source:r,instance:e.quill})})},handleLoad:function(){this.quill&&this.quill.getModule(`toolbar`)&&this.$emit(`load`,{instance:this.quill})},handleReadOnlyChange:function(){this.quill&&this.quill.enable(!this.readonly)}}};function Ke(e,n,r,i,a,s){return o(),S(`div`,t({class:e.cx(`root`)},e.ptmi(`root`)),[x(`div`,t({ref:`toolbarElement`,class:e.cx(`toolbar`)},e.ptm(`toolbar`)),[h(e.$slots,`toolbar`,{},function(){return[x(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[x(`select`,t({class:`ql-header`,defaultValue:`0`},e.ptm(`header`)),[x(`option`,t({value:`1`},e.ptm(`option`)),`Heading`,16),x(`option`,t({value:`2`},e.ptm(`option`)),`Subheading`,16),x(`option`,t({value:`0`},e.ptm(`option`)),`Normal`,16)],16),x(`select`,t({class:`ql-font`},e.ptm(`font`)),[x(`option`,f(w(e.ptm(`option`))),null,16),x(`option`,t({value:`serif`},e.ptm(`option`)),null,16),x(`option`,t({value:`monospace`},e.ptm(`option`)),null,16)],16)],16),x(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[x(`button`,t({class:`ql-bold`,type:`button`},e.ptm(`bold`)),null,16),x(`button`,t({class:`ql-italic`,type:`button`},e.ptm(`italic`)),null,16),x(`button`,t({class:`ql-underline`,type:`button`},e.ptm(`underline`)),null,16)],16),x(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[x(`select`,t({class:`ql-color`},e.ptm(`color`)),null,16),x(`select`,t({class:`ql-background`},e.ptm(`background`)),null,16)],16),x(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[x(`button`,t({class:`ql-list`,value:`ordered`,type:`button`},e.ptm(`list`)),null,16),x(`button`,t({class:`ql-list`,value:`bullet`,type:`button`},e.ptm(`list`)),null,16),x(`select`,t({class:`ql-align`},e.ptm(`select`)),[x(`option`,t({defaultValue:``},e.ptm(`option`)),null,16),x(`option`,t({value:`center`},e.ptm(`option`)),null,16),x(`option`,t({value:`right`},e.ptm(`option`)),null,16),x(`option`,t({value:`justify`},e.ptm(`option`)),null,16)],16)],16),x(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[x(`button`,t({class:`ql-link`,type:`button`},e.ptm(`link`)),null,16),x(`button`,t({class:`ql-image`,type:`button`},e.ptm(`image`)),null,16),x(`button`,t({class:`ql-code-block`,type:`button`},e.ptm(`codeBlock`)),null,16)],16),x(`span`,t({class:`ql-formats`},e.ptm(`formats`)),[x(`button`,t({class:`ql-clean`,type:`button`},e.ptm(`clean`)),null,16)],16)]})],16),x(`div`,t({ref:`editorElement`,class:e.cx(`content`),style:e.editorStyle},e.ptm(`content`)),null,16)],16)}Ge.render=Ke;var qe={key:0,class:`text-gray-500 dark:text-gray-400 text-xs`},Je=[`innerHTML`],Ye={key:2,class:`text-gray-800 dark:text-gray-100`},Xe={class:`flex items-center gap-1`},Ze={__name:`DataTableSection`,props:{items:Array,loading:Boolean,total:Number,columns:{type:Array,default:()=>[]},entity:String,orgId:String,onLoad:Function},emits:[`edit`,`delete`],setup(t){let s=t,{t:l}=F(),f=i(1),p=i(15),_=v(()=>(f.value-1)*p.value),C=v(()=>[...s.columns.map(e=>({type:`text`,width:`w-24`,headerWidth:`w-20`})),{type:`icons`,count:2,headerWidth:`w-16`}]);function w(e){f.value=e.page+1,p.value=e.rows,s.onLoad&&s.onLoad(f.value,p.value)}return n(()=>{s.onLoad&&s.onLoad(1,15)}),(n,i)=>{let s=a(`tooltip`);return t.loading?(o(),d(J,{key:0,columns:C.value,rows:8},null,8,[`columns`])):(o(),d(u(U),{key:1,value:t.items,lazy:``,totalRecords:t.total,first:_.value,rows:p.value,onPage:w,paginator:``,paginatorTemplate:`FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown`,rowsPerPageOptions:[10,15,25,50],size:`small`,class:`!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden`},{empty:r(()=>[h(n.$slots,`empty`)]),default:r(()=>[(o(!0),S(b,null,e(t.columns,e=>(o(),d(u(G),{key:e.field,field:e.field,header:e.header,sortable:``},{body:r(({data:t})=>[e.field.startsWith(`_`)?(o(),S(`span`,qe,y(t[e.field]||`-`),1)):g(``,!0),e.html?(o(),S(`div`,{key:1,class:`editor-content`,innerHTML:t[e.field]},null,8,Je)):(o(),S(`span`,Ye,y(t[e.field]||`-`),1))]),_:2},1032,[`field`,`header`]))),128)),m(u(G),{header:u(l)(`common.actions`),style:{width:`90px`},frozen:``,alignFrozen:`right`},{body:r(({data:e})=>[x(`div`,Xe,[c(m(u(D),{icon:`pi pi-pencil`,size:`small`,text:``,severity:`secondary`,onClick:t=>n.$emit(`edit`,e)},null,8,[`onClick`]),[[s,u(l)(`common.edit`),void 0,{left:!0}]]),c(m(u(D),{icon:`pi pi-trash`,size:`small`,text:``,severity:`danger`,onClick:t=>n.$emit(`delete`,e)},null,8,[`onClick`]),[[s,u(l)(`common.delete`),void 0,{left:!0}]])])]),_:1},8,[`header`])]),_:3},8,[`value`,`totalRecords`,`first`,`rows`]))}}},Qe={class:`space-y-4`},$e={key:0,class:`bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded p-3 text-xs text-red-700 dark:text-red-300`},et={__name:`DialogForm`,props:{visible:Boolean,title:String,saving:Boolean,errors:{type:Object,default:()=>({})},width:{type:String,default:`480px`}},emits:[`save`,`cancel`],setup(t){let n=t,{t:i}=F(),a=v(()=>n.width===`maximize`?`90vw`:n.width);return(n,s)=>(o(),d(u(N),{visible:t.visible,"onUpdate:visible":s[2]||=e=>n.$emit(`update:visible`,e),header:t.title,modal:``,style:p({width:a.value}),class:`p-fluid`,closable:!t.saving},{footer:r(()=>[m(u(D),{label:u(i)(`common.cancel`),size:`small`,outlined:``,severity:`secondary`,disabled:t.saving,onClick:s[0]||=e=>n.$emit(`cancel`)},null,8,[`label`,`disabled`]),m(u(D),{label:u(i)(`common.save`),icon:`pi pi-check`,size:`small`,loading:t.saving,onClick:s[1]||=e=>n.$emit(`save`)},null,8,[`label`,`loading`])]),default:r(()=>[x(`div`,Qe,[h(n.$slots,`default`),Object.keys(t.errors).length?(o(),S(`div`,$e,[(o(!0),S(b,null,e(t.errors,(e,t)=>(o(),S(`p`,{key:t,class:`mb-1`},[x(`strong`,null,y(t)+`:`,1),_(` `+y(Array.isArray(e)?e.join(`, `):e),1)]))),128))])):g(``,!0)])]),_:3},8,[`visible`,`header`,`style`,`closable`]))}},tt={class:`space-y-4`},nt={class:`flex items-center justify-between`},rt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},it={class:`text-sm text-gray-500 dark:text-gray-400`},at={class:`flex flex-col items-center justify-center py-10 text-gray-400`},ot={class:`text-sm font-medium`},st=`/api/v1/tenant/job-management/responsibilities`,ct={__name:`JobResponsibilitySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=F(),c=j(),l=i([]),d=i(!1),f=i(0),p=i(!1),h=i(!1),_=i(``),w=i(!1),E=i({}),O=i(!1),k=i(!1),A=i(``),M=i(null),N=i({main_task:``,activities:``,outputs:``,success_indicators:``}),P=v(()=>{let e=s(`job_management.responsibilities_title`);return h.value?`${e}`:`${s(`common.create`)} ${e}`}),L=v(()=>[{field:`main_task`,header:s(`job_management.main_task`),html:!0},{field:`activities`,header:s(`job_management.activities`),html:!0},{field:`outputs`,header:s(`job_management.outputs`),html:!0},{field:`success_indicators`,header:s(`job_management.success_indicators`),html:!0}]);async function z(e,t){d.value=!0;try{let r=await T.get(st,{params:{page:e,per_page:t,organization_id:n.orgId}}),i=r.data?.data||[];l.value=i.map(e=>({...e,main_task:e.main_task,activities:e.activities,outputs:e.outputs,success_indicators:e.success_indicators})),f.value=r.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function B(){h.value=!1,_.value=``,N.value={main_task:``,activities:``,outputs:``,success_indicators:``},E.value={},p.value=!0}function V(e){h.value=!0,_.value=e.id,N.value={main_task:e.main_task||``,activities:e.activities||``,outputs:e.outputs||``,success_indicators:e.success_indicators||``},E.value={},p.value=!0}async function H(){w.value=!0,E.value={};try{let e={nomenclature:n.orgName||``,full_code:n.orgCode||``,...N.value,organization_id:n.orgId};h.value?await T.put(`${st}/${_.value}`,e):await T.post(st,e),p.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),z(1,15)}catch(e){let t=I(e);Object.keys(t).length?E.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{w.value=!1}}function U(e){M.value=e,A.value=``,O.value=!0}async function W(){if(M.value){k.value=!0,A.value=``;try{await T.delete(`${st}/${M.value.id}`),O.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),z(1,15)}catch(e){A.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{k.value=!1}}}return(t,n)=>(o(),S(`div`,tt,[x(`div`,nt,[x(`div`,null,[x(`h2`,rt,y(u(s)(`job_management.responsibilities_title`)),1),x(`p`,it,y(u(s)(`job_management.responsibilities_description`)),1)]),m(u(D),{label:u(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>B()},null,8,[`label`])]),m(Ze,{items:l.value,loading:d.value,total:f.value,columns:L.value,entity:`responsibilities`,"org-id":e.orgId,"on-load":z,onEdit:V,onDelete:U},{empty:r(()=>[x(`div`,at,[n[9]||=x(`i`,{class:`pi pi-list-check text-3xl mb-2 opacity-50`},null,-1),x(`p`,ot,y(u(s)(`job_management.empty_responsibilities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),m(et,{visible:p.value,"onUpdate:visible":n[5]||=e=>p.value=e,title:P.value,saving:w.value,errors:E.value,width:`maximize`,onSave:H,onCancel:n[6]||=e=>p.value=!1},{default:r(()=>[p.value?(o(),S(b,{key:0},[m(R,{label:u(s)(`job_management.main_task`),errors:E.value?.main_task},{default:r(()=>[m(u(Ge),{modelValue:N.value.main_task,"onUpdate:modelValue":n[1]||=e=>N.value.main_task=e,editorStyle:`height:120px`,class:C({"p-invalid":E.value?.main_task})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),m(R,{label:u(s)(`job_management.activities`),errors:E.value?.activities},{default:r(()=>[m(u(Ge),{modelValue:N.value.activities,"onUpdate:modelValue":n[2]||=e=>N.value.activities=e,editorStyle:`height:120px`,class:C({"p-invalid":E.value?.activities})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),m(R,{label:u(s)(`job_management.outputs`),errors:E.value?.outputs},{default:r(()=>[m(u(Ge),{modelValue:N.value.outputs,"onUpdate:modelValue":n[3]||=e=>N.value.outputs=e,editorStyle:`height:120px`,class:C({"p-invalid":E.value?.outputs})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`]),m(R,{label:u(s)(`job_management.success_indicators`),errors:E.value?.success_indicators},{default:r(()=>[m(u(Ge),{modelValue:N.value.success_indicators,"onUpdate:modelValue":n[4]||=e=>N.value.success_indicators=e,editorStyle:`height:120px`,class:C({"p-invalid":E.value?.success_indicators})},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])],64)):g(``,!0)]),_:1},8,[`visible`,`title`,`saving`,`errors`]),m(q,{visible:O.value,"onUpdate:visible":n[7]||=e=>O.value=e,loading:k.value,"error-msg":A.value,onConfirm:W,onCancel:n[8]||=e=>O.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},lt={class:`space-y-4`},ut={class:`flex items-center justify-between`},dt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},ft={class:`text-sm text-gray-500 dark:text-gray-400`},pt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},mt={class:`text-sm font-medium`},ht=`/api/v1/tenant/job-management/hr-authorities`,gt={__name:`JobHRAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=F(),c=j(),l=i([]),d=i(!1),f=i(0),p=i(!1),h=i(!1),g=i(``),_=i(!1),b=i({}),w=i(!1),E=i(!1),O=i(``),k=i(null),A=i({description:``}),M=v(()=>{let e=s(`job_management.hr_authorities`);return h.value?`${s(`common.edit`)} ${e}`:`${s(`common.create`)} ${e}`}),N=v(()=>[{field:`description`,header:s(`job_management.description`)}]);async function P(e,t){d.value=!0;try{let r=await T.get(ht,{params:{page:e,per_page:t,organization_id:n.orgId}});l.value=r.data?.data||[],f.value=r.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function L(){h.value=!1,g.value=``,A.value={nomenclature:``,full_code:``,description:``},b.value={},p.value=!0}function z(e){h.value=!0,g.value=e.id,A.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},b.value={},p.value=!0}async function V(){_.value=!0,b.value={};try{let e={...A.value,nomenclature:n.orgName||``,full_code:n.orgCode||``,organization_id:n.orgId};h.value?await T.put(`${ht}/${g.value}`,e):await T.post(ht,e),p.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),P(1,15)}catch(e){let t=I(e);Object.keys(t).length?b.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{_.value=!1}}function H(e){k.value=e,O.value=``,w.value=!0}async function U(){if(k.value){E.value=!0,O.value=``;try{await T.delete(`${ht}/${k.value.id}`),w.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),P(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{E.value=!1}}}return(t,n)=>(o(),S(`div`,lt,[x(`div`,ut,[x(`div`,null,[x(`h2`,dt,y(u(s)(`job_management.hr_authorities`)),1),x(`p`,ft,y(u(s)(`job_management.authority_description`)),1)]),m(u(D),{label:u(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>L()},null,8,[`label`])]),m(Ze,{items:l.value,loading:d.value,total:f.value,columns:N.value,entity:`hr-authorities`,"org-id":e.orgId,"on-load":P,onEdit:z,onDelete:H},{empty:r(()=>[x(`div`,pt,[n[6]||=x(`i`,{class:`pi pi-users text-3xl mb-2 opacity-50`},null,-1),x(`p`,mt,y(u(s)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),m(et,{visible:p.value,"onUpdate:visible":n[2]||=e=>p.value=e,title:M.value,saving:_.value,errors:b.value,onSave:V,onCancel:n[3]||=e=>p.value=!1},{default:r(()=>[m(R,{label:u(s)(`job_management.description`),errors:b.value?.description},{default:r(()=>[m(u(B),{modelValue:A.value.description,"onUpdate:modelValue":n[1]||=e=>A.value.description=e,rows:`3`,class:C([`w-full`,{"p-invalid":b.value?.description}])},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),m(q,{visible:w.value,"onUpdate:visible":n[4]||=e=>w.value=e,loading:E.value,"error-msg":O.value,onConfirm:U,onCancel:n[5]||=e=>w.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},_t={class:`space-y-4`},vt={class:`flex items-center justify-between`},yt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},bt={class:`text-sm text-gray-500 dark:text-gray-400`},xt={class:`flex flex-col items-center justify-center py-10 text-gray-400`},St={class:`text-sm font-medium`},Ct=`/api/v1/tenant/job-management/operational-authorities`,wt={__name:`JobOpAuthoritySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``}},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=F(),c=j(),l=i([]),d=i(!1),f=i(0),p=i(!1),h=i(!1),g=i(``),_=i(!1),b=i({}),w=i(!1),E=i(!1),O=i(``),k=i(null),A=i({description:``}),M=v(()=>{let e=s(`job_management.op_authorities`);return h.value?`${s(`common.edit`)} ${e}`:`${s(`common.create`)} ${e}`}),N=v(()=>[{field:`description`,header:s(`job_management.description`)}]);async function P(e,t){d.value=!0;try{let r=await T.get(Ct,{params:{page:e,per_page:t,organization_id:n.orgId}});l.value=r.data?.data||[],f.value=r.data?.total||0}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function L(){h.value=!1,g.value=``,A.value={nomenclature:``,full_code:``,description:``},b.value={},p.value=!0}function z(e){h.value=!0,g.value=e.id,A.value={nomenclature:e.nomenclature||``,full_code:e.full_code||``,description:e.description||``},b.value={},p.value=!0}async function V(){_.value=!0,b.value={};try{let e={...A.value,nomenclature:n.orgName||``,full_code:n.orgCode||``,organization_id:n.orgId};h.value?await T.put(`${Ct}/${g.value}`,e):await T.post(Ct,e),p.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),P(1,15)}catch(e){let t=I(e);Object.keys(t).length?b.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{_.value=!1}}function H(e){k.value=e,O.value=``,w.value=!0}async function U(){if(k.value){E.value=!0,O.value=``;try{await T.delete(`${Ct}/${k.value.id}`),w.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),P(1,15)}catch(e){O.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{E.value=!1}}}return(t,n)=>(o(),S(`div`,_t,[x(`div`,vt,[x(`div`,null,[x(`h2`,yt,y(u(s)(`job_management.op_authorities`)),1),x(`p`,bt,y(u(s)(`job_management.authority_description`)),1)]),m(u(D),{label:u(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>L()},null,8,[`label`])]),m(Ze,{items:l.value,loading:d.value,total:f.value,columns:N.value,entity:`operational-authorities`,"org-id":e.orgId,"on-load":P,onEdit:z,onDelete:H},{empty:r(()=>[x(`div`,xt,[n[6]||=x(`i`,{class:`pi pi-cog text-3xl mb-2 opacity-50`},null,-1),x(`p`,St,y(u(s)(`job_management.empty_authorities`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),m(et,{visible:p.value,"onUpdate:visible":n[2]||=e=>p.value=e,title:M.value,saving:_.value,errors:b.value,onSave:V,onCancel:n[3]||=e=>p.value=!1},{default:r(()=>[m(R,{label:u(s)(`job_management.description`),errors:b.value?.description},{default:r(()=>[m(u(B),{modelValue:A.value.description,"onUpdate:modelValue":n[1]||=e=>A.value.description=e,class:C([`w-full`,{"p-invalid":b.value?.description}]),rows:`3`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),m(q,{visible:w.value,"onUpdate:visible":n[4]||=e=>w.value=e,loading:E.value,"error-msg":O.value,onConfirm:U,onCancel:n[5]||=e=>w.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Tt={class:`space-y-4`},Et={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Dt={class:`text-sm text-gray-500 dark:text-gray-400`},Ot={class:`max-w-2xl`},kt={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},At={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},jt={class:`flex justify-end gap-2 pt-2`},Mt=`/api/v1/tenant/job-management/working-activities`,Nt={__name:`JobActivitySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:t}){let a=t,s=e,{t:c}=F(),l=j(),f=i(!1),p=i(!0),h=i(``),_=i({}),v=i(``),b=i(!1),w=i(!1),E=i(``),O=i({job_management_value_id:``}),k=i([]);async function A(){try{let e=await T.get(`/api/v1/tenant/job-management/values`,{params:{type:`activity`,per_page:100}});k.value=(e.data?.data||[]).map(e=>({label:e.descriptions,value:e.id}))}catch{}}async function M(){if(!s.orgId){p.value=!1;return}try{let e=(await T.get(Mt,{params:{organization_id:s.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];v.value=t.id,O.value.job_management_value_id=t.job_management_value_id||``}}catch{}}async function N(){h.value=``,_.value={},f.value=!0;try{let e={nomenclature:s.orgName||``,full_code:s.orgCode||``,job_management_value_id:O.value.job_management_value_id||null,organization_id:s.orgId};if(v.value)await T.put(`${Mt}/${v.value}`,{job_management_value_id:O.value.job_management_value_id||``});else{let t=await T.post(Mt,e);v.value=t.data?.data?.id||``}l.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=I(e);Object.keys(t).length>0?(_.value=t,h.value=Object.values(t).join(`, `)):h.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{f.value=!1}}async function P(){if(v.value){w.value=!0,E.value=``;try{await T.delete(`${Mt}/${v.value}`),b.value=!1,v.value=``,O.value.job_management_value_id=``,a(`saved`),l.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){E.value=e?.response?.data?.error?.message||c(`message.operation_failed`)}finally{w.value=!1}}}return n(async()=>{try{await Promise.all([A(),M()])}finally{p.value=!1}}),(t,n)=>(o(),S(`div`,Tt,[x(`div`,null,[x(`h2`,Et,y(u(c)(`job_management.activities`)),1),x(`p`,Dt,y(u(c)(`job_management.activity_description`)),1)]),x(`div`,Ot,[p.value?(o(),d(Y,{key:0,type:`detail`,count:1,rows:3,cols:`grid-cols-1`,padding:`p-5`})):(o(),S(`div`,kt,[m(R,{label:u(c)(`organization.nomenclature`)},{default:r(()=>[m(V,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`organization.full_code`)},{default:r(()=>[m(V,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`job_values.types.activity`),errors:_.value?.job_management_value_id},{default:r(()=>[m(K,{modelValue:O.value.job_management_value_id,"onUpdate:modelValue":n[0]||=e=>O.value.job_management_value_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),class:C({"p-invalid":_.value?.job_management_value_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h.value?(o(),S(`div`,At,y(h.value),1)):g(``,!0),x(`div`,jt,[v.value?(o(),d(u(D),{key:0,label:u(c)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:n[1]||=e=>b.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:v.value?u(c)(`common.update`):u(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:f.value,disabled:f.value,onClick:N},null,8,[`label`,`loading`,`disabled`])])]))]),m(q,{visible:b.value,"onUpdate:visible":n[2]||=e=>b.value=e,loading:w.value,"error-msg":E.value,onConfirm:P,onCancel:n[3]||=e=>b.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Pt={class:`space-y-4`},Ft={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},It={class:`text-sm text-gray-500 dark:text-gray-400`},Lt={class:`max-w-2xl`},Rt={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},zt={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Bt={class:`flex justify-end gap-2 pt-2`},Vt=`/api/v1/tenant/job-management/working-risks`,Ht={__name:`JobRiskSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:t}){let a=t,s=e,{t:c}=F(),l=j(),f=i(!1),p=i(!0),h=i(``),_=i({}),v=i(``),b=i(!1),w=i(!1),E=i(``),O=i({job_management_value_environment_id:``,job_management_value_hazard_id:``}),k=i([]),A=i([]);async function M(){try{let[e,t]=await Promise.all([T.get(`/api/v1/tenant/job-management/values`,{params:{type:`environment`,per_page:100}}),T.get(`/api/v1/tenant/job-management/values`,{params:{type:`risk`,per_page:100}})]);k.value=(e.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),A.value=(t.data?.data||[]).map(e=>({label:e.descriptions,value:e.id}))}catch{}}async function N(){if(!s.orgId){p.value=!1;return}try{let e=(await T.get(Vt,{params:{organization_id:s.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];v.value=t.id,O.value.job_management_value_environment_id=t.job_management_value_environment_id||``,O.value.job_management_value_hazard_id=t.job_management_value_hazard_id||``}}catch{}}async function P(){h.value=``,_.value={},f.value=!0;try{let e={nomenclature:s.orgName||``,full_code:s.orgCode||``,job_management_value_environment_id:O.value.job_management_value_environment_id||null,job_management_value_hazard_id:O.value.job_management_value_hazard_id||null,organization_id:s.orgId};if(v.value)await T.put(`${Vt}/${v.value}`,{job_management_value_environment_id:O.value.job_management_value_environment_id||``,job_management_value_hazard_id:O.value.job_management_value_hazard_id||``});else{let t=await T.post(Vt,e);v.value=t.data?.data?.id||``}l.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=I(e);Object.keys(t).length>0?(_.value=t,h.value=Object.values(t).join(`, `)):h.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{f.value=!1}}async function L(){if(v.value){w.value=!0,E.value=``;try{await T.delete(`${Vt}/${v.value}`),b.value=!1,v.value=``,O.value.job_management_value_environment_id=``,O.value.job_management_value_hazard_id=``,a(`saved`),l.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){E.value=e?.response?.data?.error?.message||c(`message.operation_failed`)}finally{w.value=!1}}}return n(async()=>{try{await Promise.all([M(),N()])}finally{p.value=!1}}),(t,n)=>(o(),S(`div`,Pt,[x(`div`,null,[x(`h2`,Ft,y(u(c)(`job_management.risks`)),1),x(`p`,It,y(u(c)(`job_management.risk_description`)),1)]),x(`div`,Lt,[p.value?(o(),d(Y,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(o(),S(`div`,Rt,[m(R,{label:u(c)(`organization.nomenclature`)},{default:r(()=>[m(V,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`organization.full_code`)},{default:r(()=>[m(V,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`job_management.work_environment`),errors:_.value?.job_management_value_environment_id},{default:r(()=>[m(K,{modelValue:O.value.job_management_value_environment_id,"onUpdate:modelValue":n[0]||=e=>O.value.job_management_value_environment_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),class:C({"p-invalid":_.value?.job_management_value_environment_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),m(R,{label:u(c)(`job_management.risk`),errors:_.value?.job_management_value_hazard_id},{default:r(()=>[m(K,{modelValue:O.value.job_management_value_hazard_id,"onUpdate:modelValue":n[1]||=e=>O.value.job_management_value_hazard_id=e,options:A.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),class:C({"p-invalid":_.value?.job_management_value_hazard_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h.value?(o(),S(`div`,zt,y(h.value),1)):g(``,!0),x(`div`,Bt,[v.value?(o(),d(u(D),{key:0,label:u(c)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:n[2]||=e=>b.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:v.value?u(c)(`common.update`):u(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:f.value,disabled:f.value,onClick:P},null,8,[`label`,`loading`,`disabled`])])]))]),m(q,{visible:b.value,"onUpdate:visible":n[3]||=e=>b.value=e,loading:w.value,"error-msg":E.value,onConfirm:L,onCancel:n[4]||=e=>b.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Ut={class:`space-y-4`},Wt={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Gt={class:`text-sm text-gray-500 dark:text-gray-400`},Kt={class:`max-w-2xl`},qt={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Jt={class:`pt-4 border-t border-gray-200 dark:border-gray-700`},Yt={class:`flex items-center gap-2 mb-3`},Xt={class:`text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400`},Zt={class:`space-y-4`},Qt={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},$t={class:`flex justify-end gap-2 pt-2`},en={class:`max-w-3xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5 space-y-4`},tn={class:`flex items-center justify-between gap-2 flex-wrap`},nn={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},rn={class:`text-sm text-gray-500 dark:text-gray-400`},an={key:0,class:`text-sm text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg px-3 py-2`},on={key:1,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},sn={key:2,class:`overflow-x-auto border border-gray-200 dark:border-gray-700 rounded-lg`},cn={class:`w-full text-sm`},ln={class:`bg-gray-50 dark:bg-gray-700/40 text-left`},un={class:`px-3 py-2 font-semibold text-gray-600 dark:text-gray-300`},dn={class:`px-3 py-2 font-semibold text-gray-600 dark:text-gray-300`},fn={class:`px-3 py-2 align-top text-gray-500 dark:text-gray-400`},pn={class:`px-3 py-2`},mn={class:`px-3 py-2`},hn={class:`px-3 py-2 align-top`},gn={key:3,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},_n={key:4,class:`flex justify-end gap-2 pt-2`},Q=`/api/v1/tenant/job-management/relationships`,vn={__name:`JobRelationshipSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},orgSummaryId:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(t,{emit:a}){let s=a,c=t,{t:l}=F(),f=j(),p=i(!1),h=i(!0),_=i(``),v=i({}),w=i(``),E=i(!1),O=i(!1),k=i(``),A=i({job_management_value_relationship_id:``,job_management_value_frequency_id:``}),M=i([]),N=i([]),P=i([]),L=i([]),z=i(!1),B=i(``);async function H(){try{let[e,t,n]=await Promise.all([T.get(`/api/v1/tenant/job-management/values`,{params:{type:`relationship`,per_page:100}}),T.get(`/api/v1/tenant/job-management/values`,{params:{type:`frequency`,per_page:100}}),c.orgSummaryId?T.get(`/api/v1/tenant/organizations`,{params:{summary_id:c.orgSummaryId,per_page:100}}):Promise.resolve({data:{data:[]}})]);M.value=(e.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),N.value=(t.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),P.value=(n.data?.data||[]).filter(e=>e.id!==c.orgId).map(e=>({label:e.full_code?`${e.full_code} - ${e.nomenclature}`:e.nomenclature,value:e.id}))}catch{}}async function U(){if(!c.orgId){h.value=!1;return}try{let e=(await T.get(Q,{params:{organization_id:c.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];w.value=t.id,A.value.job_management_value_relationship_id=t.job_management_value_relationship_id||``,A.value.job_management_value_frequency_id=t.job_management_value_frequency_id||``,await re()}}catch{}}async function W(){_.value=``,v.value={},p.value=!0;try{let e={nomenclature:c.orgName||``,full_code:c.orgCode||``,job_management_value_relationship_id:A.value.job_management_value_relationship_id||null,job_management_value_frequency_id:A.value.job_management_value_frequency_id||null,organization_id:c.orgId};if(w.value)await T.put(`${Q}/${w.value}`,{job_management_value_relationship_id:A.value.job_management_value_relationship_id||``,job_management_value_frequency_id:A.value.job_management_value_frequency_id||``});else{let t=await T.post(Q,e);w.value=t.data?.data?.id||``}f.add({severity:`success`,summary:l(`message.success`),detail:l(`common.saved`),life:2e3}),s(`saved`)}catch(e){let t=I(e);Object.keys(t).length>0?(v.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||l(`message.operation_failed`)}finally{p.value=!1}}async function G(){if(w.value){O.value=!0,k.value=``;try{await T.delete(`${Q}/${w.value}`),E.value=!1,w.value=``,A.value.job_management_value_relationship_id=``,A.value.job_management_value_frequency_id=``,L.value=[],s(`saved`),f.add({severity:`success`,summary:l(`message.success`),detail:l(`message.deleted`),life:2e3})}catch(e){k.value=e?.response?.data?.error?.message||l(`message.operation_failed`)}finally{O.value=!1}}}let J=0;function ee(){w.value&&L.value.push({_key:`new-${++J}`,id:``,organization_id:``,activity:``})}function te(e){let t=L.value[e];t&&(t.id?ne(t.id,e):L.value.splice(e,1))}async function ne(e,t){try{await T.delete(`${Q}/${w.value}/details/${e}`),L.value.splice(t,1),f.add({severity:`success`,summary:l(`message.success`),detail:l(`message.deleted`),life:2e3})}catch(e){f.add({severity:`error`,summary:l(`message.error`),detail:e?.response?.data?.error?.message||l(`message.operation_failed`),life:4e3})}}async function re(){if(w.value)try{let e=await T.get(`${Q}/${w.value}/details`);L.value=(e.data?.data||[]).map(e=>({_key:`db-${++J}`,id:e.id,organization_id:e.organization_id||``,activity:e.activity||``}))}catch{}}async function X(){if(!(!w.value||z.value)){B.value=``,z.value=!0;try{for(let e of L.value){let t={organization_id:e.organization_id||``,activity:e.activity||``};e.id?await T.put(`${Q}/${w.value}/details/${e.id}`,t):e.id=(await T.post(`${Q}/${w.value}/details`,t)).data?.data?.id||``}await re(),f.add({severity:`success`,summary:l(`message.success`),detail:l(`job_management.relationship_details_saved`),life:2e3})}catch(e){let t=I(e);Object.keys(t).length>0?B.value=Object.values(t).join(`, `):B.value=e?.response?.data?.error?.message||e.message||l(`message.operation_failed`)}finally{z.value=!1}}}return n(async()=>{try{await Promise.all([H(),U()])}finally{h.value=!1}}),(n,i)=>(o(),S(`div`,Ut,[x(`div`,null,[x(`h2`,Wt,y(u(l)(`job_management.relationships`)),1),x(`p`,Gt,y(u(l)(`job_management.relationship_description`)),1)]),x(`div`,Kt,[h.value?(o(),d(Y,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(o(),S(`div`,qt,[m(R,{label:u(l)(`organization.nomenclature`)},{default:r(()=>[m(V,{"model-value":t.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(l)(`organization.full_code`)},{default:r(()=>[m(V,{"model-value":t.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),x(`div`,Jt,[x(`div`,Yt,[i[5]||=x(`div`,{class:`w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400`},[x(`i`,{class:`pi pi-compass text-sm`})],-1),x(`h3`,Xt,y(u(l)(`job_management.relationship_group_scope`)),1),i[6]||=x(`div`,{class:`flex-1 border-t border-gray-200 dark:border-gray-700`},null,-1)]),x(`div`,Zt,[m(R,{label:u(l)(`job_management.relationship_type`),errors:v.value?.job_management_value_relationship_id},{default:r(()=>[m(K,{modelValue:A.value.job_management_value_relationship_id,"onUpdate:modelValue":i[0]||=e=>A.value.job_management_value_relationship_id=e,options:M.value,"option-label":`label`,"option-value":`value`,placeholder:u(l)(`common.select`),class:C({"p-invalid":v.value?.job_management_value_relationship_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),m(R,{label:u(l)(`job_management.frequency`),errors:v.value?.job_management_value_frequency_id},{default:r(()=>[m(K,{modelValue:A.value.job_management_value_frequency_id,"onUpdate:modelValue":i[1]||=e=>A.value.job_management_value_frequency_id=e,options:N.value,"option-label":`label`,"option-value":`value`,placeholder:u(l)(`common.select`),class:C({"p-invalid":v.value?.job_management_value_frequency_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`])])]),_.value?(o(),S(`div`,Qt,y(_.value),1)):g(``,!0),x(`div`,$t,[w.value?(o(),d(u(D),{key:0,label:u(l)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:i[2]||=e=>E.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:w.value?u(l)(`common.update`):u(l)(`common.save`),icon:`pi pi-check`,size:`small`,loading:p.value,disabled:p.value,onClick:W},null,8,[`label`,`loading`,`disabled`])])]))]),x(`div`,en,[x(`div`,tn,[x(`div`,null,[x(`h3`,nn,y(u(l)(`job_management.relationship_details`)),1),x(`p`,rn,y(u(l)(`job_management.relationship_details_description`)),1)]),m(u(D),{label:u(l)(`job_management.add_relationship_detail`),icon:`pi pi-plus`,size:`small`,outlined:``,disabled:!w.value||z.value,onClick:ee},null,8,[`label`,`disabled`])]),w.value?L.value.length===0?(o(),S(`div`,on,y(u(l)(`job_management.no_relationship_details`)),1)):g(``,!0):(o(),S(`div`,an,y(u(l)(`job_management.save_relationship_first`)),1)),L.value.length>0?(o(),S(`div`,sn,[x(`table`,cn,[x(`thead`,null,[x(`tr`,ln,[i[7]||=x(`th`,{class:`px-3 py-2 w-10 font-semibold text-gray-600 dark:text-gray-300`},`#`,-1),x(`th`,un,y(u(l)(`job_management.relationship_organization`)),1),x(`th`,dn,y(u(l)(`job_management.relationship_activity`)),1),i[8]||=x(`th`,{class:`px-3 py-2 w-12`},null,-1)])]),x(`tbody`,null,[(o(!0),S(b,null,e(L.value,(e,t)=>(o(),S(`tr`,{key:e._key,class:`border-t border-gray-200 dark:border-gray-700`},[x(`td`,fn,y(t+1),1),x(`td`,pn,[m(K,{modelValue:e.organization_id,"onUpdate:modelValue":t=>e.organization_id=t,options:P.value,"option-label":`label`,"option-value":`value`,placeholder:u(l)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),x(`td`,mn,[m(V,{modelValue:e.activity,"onUpdate:modelValue":t=>e.activity=t,placeholder:u(l)(`job_management.relationship_activity`)},null,8,[`modelValue`,`onUpdate:modelValue`,`placeholder`])]),x(`td`,hn,[m(u(D),{icon:`pi pi-trash`,severity:`danger`,size:`small`,text:``,rounded:``,"aria-label":`Remove`,onClick:e=>te(t)},null,8,[`onClick`])])]))),128))])])])):g(``,!0),B.value?(o(),S(`div`,gn,y(B.value),1)):g(``,!0),L.value.length>0?(o(),S(`div`,_n,[m(u(D),{label:u(l)(`job_management.save_relationship_details`),icon:`pi pi-save`,size:`small`,loading:z.value,disabled:z.value||!w.value,onClick:X},null,8,[`label`,`loading`,`disabled`])])):g(``,!0)]),m(q,{visible:E.value,"onUpdate:visible":i[3]||=e=>E.value=e,loading:O.value,"error-msg":k.value,onConfirm:G,onCancel:i[4]||=e=>E.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},yn={class:`space-y-4`},bn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},xn={class:`text-sm text-gray-500 dark:text-gray-400`},Sn={class:`max-w-2xl`},Cn={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},wn={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Tn={class:`flex justify-end gap-2 pt-2`},En=`/api/v1/tenant/job-management/subordinate-controls`,Dn={__name:`JobSubordinateSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:t}){let a=t,s=e,{t:c}=F(),l=j(),f=i(!1),p=i(!0),h=i(``),_=i({}),v=i(``),b=i(!1),w=i(!1),E=i(``),O=i({job_management_value_id:``}),k=i([]);async function A(){try{let e=await T.get(`/api/v1/tenant/job-management/values`,{params:{type:`subordinate`,per_page:100}});k.value=(e.data?.data||[]).map(e=>({label:e.descriptions,value:e.id}))}catch{}}async function M(){if(!s.orgId){p.value=!1;return}try{let e=(await T.get(En,{params:{organization_id:s.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];v.value=t.id,O.value.job_management_value_id=t.job_management_value_id||``}}catch{}}async function N(){h.value=``,_.value={},f.value=!0;try{let e={nomenclature:s.orgName||``,full_code:s.orgCode||``,job_management_value_id:O.value.job_management_value_id||null,organization_id:s.orgId};if(v.value)await T.put(`${En}/${v.value}`,{job_management_value_id:O.value.job_management_value_id||``});else{let t=await T.post(En,e);v.value=t.data?.data?.id||``}l.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=I(e);Object.keys(t).length>0?(_.value=t,h.value=Object.values(t).join(`, `)):h.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{f.value=!1}}async function P(){if(v.value){w.value=!0,E.value=``;try{await T.delete(`${En}/${v.value}`),b.value=!1,v.value=``,O.value.job_management_value_id=``,a(`saved`),l.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){E.value=e?.response?.data?.error?.message||c(`message.operation_failed`)}finally{w.value=!1}}}return n(async()=>{try{await Promise.all([A(),M()])}finally{p.value=!1}}),(t,n)=>(o(),S(`div`,yn,[x(`div`,null,[x(`h2`,bn,y(u(c)(`job_management.subordinate_controls`)),1),x(`p`,xn,y(u(c)(`job_management.subordinate_description`)),1)]),x(`div`,Sn,[p.value?(o(),d(Y,{key:0,type:`detail`,count:1,rows:3,cols:`grid-cols-1`,padding:`p-5`})):(o(),S(`div`,Cn,[m(R,{label:u(c)(`organization.nomenclature`)},{default:r(()=>[m(V,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`organization.full_code`)},{default:r(()=>[m(V,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`job_management.control_type`),errors:_.value?.job_management_value_id},{default:r(()=>[m(K,{modelValue:O.value.job_management_value_id,"onUpdate:modelValue":n[0]||=e=>O.value.job_management_value_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),class:C({"p-invalid":_.value?.job_management_value_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h.value?(o(),S(`div`,wn,y(h.value),1)):g(``,!0),x(`div`,Tn,[v.value?(o(),d(u(D),{key:0,label:u(c)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:n[1]||=e=>b.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:v.value?u(c)(`common.update`):u(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:f.value,disabled:f.value,onClick:N},null,8,[`label`,`loading`,`disabled`])])]))]),m(q,{visible:b.value,"onUpdate:visible":n[2]||=e=>b.value=e,loading:w.value,"error-msg":E.value,onConfirm:P,onCancel:n[3]||=e=>b.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},On={class:`space-y-4`},kn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},An={class:`text-sm text-gray-500 dark:text-gray-400`},jn={class:`max-w-2xl`},Mn={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Nn={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Pn={class:`flex justify-end gap-2 pt-2`},Fn=`/api/v1/tenant/job-management/assets`,In={__name:`JobAssetSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:t}){let a=t,s=e,{t:c}=F(),l=j(),f=i(!1),p=i(!0),h=i(``),_=i({}),v=i(``),b=i(!1),w=i(!1),E=i(``),O=i({job_management_value_asset_id:``,job_management_value_authority_id:``}),k=i([]),A=i([]);async function M(){try{let[e,t]=await Promise.all([T.get(`/api/v1/tenant/job-management/values`,{params:{type:`asset`,per_page:100}}),T.get(`/api/v1/tenant/job-management/values`,{params:{type:`asset_authority`,per_page:100}})]);k.value=(e.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),A.value=(t.data?.data||[]).map(e=>({label:e.descriptions,value:e.id}))}catch{}}async function N(){if(!s.orgId){p.value=!1;return}try{let e=(await T.get(Fn,{params:{organization_id:s.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];v.value=t.id,O.value.job_management_value_asset_id=t.job_management_value_asset_id||``,O.value.job_management_value_authority_id=t.job_management_value_authority_id||``}}catch{}}async function P(){h.value=``,_.value={},f.value=!0;try{let e={nomenclature:s.orgName||``,full_code:s.orgCode||``,job_management_value_asset_id:O.value.job_management_value_asset_id||null,job_management_value_authority_id:O.value.job_management_value_authority_id||null,organization_id:s.orgId};if(v.value)await T.put(`${Fn}/${v.value}`,{job_management_value_asset_id:O.value.job_management_value_asset_id||``,job_management_value_authority_id:O.value.job_management_value_authority_id||``});else{let t=await T.post(Fn,e);v.value=t.data?.data?.id||``}l.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=I(e);Object.keys(t).length>0?(_.value=t,h.value=Object.values(t).join(`, `)):h.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{f.value=!1}}async function L(){if(v.value){w.value=!0,E.value=``;try{await T.delete(`${Fn}/${v.value}`),b.value=!1,v.value=``,O.value.job_management_value_asset_id=``,O.value.job_management_value_authority_id=``,a(`saved`),l.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){E.value=e?.response?.data?.error?.message||c(`message.operation_failed`)}finally{w.value=!1}}}return n(async()=>{try{await Promise.all([M(),N()])}finally{p.value=!1}}),(t,n)=>(o(),S(`div`,On,[x(`div`,null,[x(`h2`,kn,y(u(c)(`job_management.assets`)),1),x(`p`,An,y(u(c)(`job_management.asset_description`)),1)]),x(`div`,jn,[p.value?(o(),d(Y,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(o(),S(`div`,Mn,[m(R,{label:u(c)(`organization.nomenclature`)},{default:r(()=>[m(V,{"model-value":e.orgName,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`organization.full_code`)},{default:r(()=>[m(V,{"model-value":e.orgCode,disabled:``,class:`!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed`},null,8,[`model-value`])]),_:1},8,[`label`]),m(R,{label:u(c)(`job_management.asset_type`),errors:_.value?.job_management_value_asset_id},{default:r(()=>[m(K,{modelValue:O.value.job_management_value_asset_id,"onUpdate:modelValue":n[0]||=e=>O.value.job_management_value_asset_id=e,options:k.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),class:C({"p-invalid":_.value?.job_management_value_asset_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),m(R,{label:u(c)(`job_management.authority_level`),errors:_.value?.job_management_value_authority_id},{default:r(()=>[m(K,{modelValue:O.value.job_management_value_authority_id,"onUpdate:modelValue":n[1]||=e=>O.value.job_management_value_authority_id=e,options:A.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),class:C({"p-invalid":_.value?.job_management_value_authority_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),h.value?(o(),S(`div`,Nn,y(h.value),1)):g(``,!0),x(`div`,Pn,[v.value?(o(),d(u(D),{key:0,label:u(c)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:n[2]||=e=>b.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:v.value?u(c)(`common.update`):u(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:f.value,disabled:f.value,onClick:P},null,8,[`label`,`loading`,`disabled`])])]))]),m(q,{visible:b.value,"onUpdate:visible":n[3]||=e=>b.value=e,loading:w.value,"error-msg":E.value,onConfirm:L,onCancel:n[4]||=e=>b.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Ln={class:`space-y-4`},Rn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},zn={class:`text-sm text-gray-500 dark:text-gray-400`},Bn={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Vn={class:`flex items-center justify-between gap-4 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3`},Hn={class:`min-w-0`},Un={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},Wn={class:`text-xs text-gray-500 dark:text-gray-400 mt-0.5`},Gn={class:`space-y-4 pt-4 border-t border-gray-200 dark:border-gray-700`},Kn={key:0,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},qn={class:`flex justify-end gap-2 pt-2`},Jn=`/api/v1/tenant/job-management/financials`,Yn={__name:`JobFinancialSection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})}},emits:[`saved`],setup(e,{emit:t}){let a=t,s=e,{t:c}=F(),f=j(),p=i(!1),h=i(!0),_=i(``),b=i({}),w=i(``),E=i(!1),O=i(!1),k=i(``),A=i({is_authorized:!1,job_management_value_cash_id:``,job_management_value_authority_id:``,job_management_value_impact_id:``}),M=i([]),N=i([]),P=i([]),L=i([]),z=i([]),B=v(()=>A.value.is_authorized?N.value:P.value),V=v(()=>A.value.is_authorized?L.value:z.value);async function H(){try{let[e,t,n,r,i]=await Promise.all([T.get(`/api/v1/tenant/job-management/values`,{params:{type:`cash`,per_page:100}}),T.get(`/api/v1/tenant/job-management/values`,{params:{type:`authority`,per_page:100}}),T.get(`/api/v1/tenant/job-management/values`,{params:{type:`authority_unauthorized`,per_page:100}}),T.get(`/api/v1/tenant/job-management/values`,{params:{type:`impact`,per_page:100}}),T.get(`/api/v1/tenant/job-management/values`,{params:{type:`impact_unauthorized`,per_page:100}})]);M.value=(e.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),N.value=(t.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),P.value=(n.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),L.value=(r.data?.data||[]).map(e=>({label:e.descriptions,value:e.id})),z.value=(i.data?.data||[]).map(e=>({label:e.descriptions,value:e.id}))}catch{}}let U=!1;l(()=>A.value.is_authorized,(e,t)=>{U||e===t||(A.value.job_management_value_cash_id=``,A.value.job_management_value_authority_id=``,A.value.job_management_value_impact_id=``)},{flush:`sync`});async function W(){if(!s.orgId){h.value=!1;return}try{let e=(await T.get(Jn,{params:{organization_id:s.orgId,per_page:1}})).data?.data||[];if(e.length>0){let t=e[0];U=!0,w.value=t.id,A.value.is_authorized=!!t.is_authorized,A.value.job_management_value_cash_id=t.job_management_value_cash_id||``,A.value.job_management_value_authority_id=t.job_management_value_authority_id||``,A.value.job_management_value_impact_id=t.job_management_value_impact_id||``,U=!1}}catch{}}async function G(){_.value=``,b.value={},p.value=!0;try{let e=!!A.value.is_authorized,t={nomenclature:s.orgName||``,full_code:s.orgCode||``,is_authorized:e,job_management_value_cash_id:e&&A.value.job_management_value_cash_id||null,job_management_value_authority_id:A.value.job_management_value_authority_id||null,job_management_value_impact_id:A.value.job_management_value_impact_id||null,organization_id:s.orgId};if(w.value)await T.put(`${Jn}/${w.value}`,{is_authorized:e,job_management_value_cash_id:e&&A.value.job_management_value_cash_id||``,job_management_value_authority_id:A.value.job_management_value_authority_id||``,job_management_value_impact_id:A.value.job_management_value_impact_id||``});else{let e=await T.post(Jn,t);w.value=e.data?.data?.id||``}f.add({severity:`success`,summary:c(`message.success`),detail:c(`common.saved`),life:2e3}),a(`saved`)}catch(e){let t=I(e);Object.keys(t).length>0?(b.value=t,_.value=Object.values(t).join(`, `)):_.value=e?.response?.data?.error?.message||e.message||c(`message.operation_failed`)}finally{p.value=!1}}async function J(){if(w.value){O.value=!0,k.value=``;try{await T.delete(`${Jn}/${w.value}`),E.value=!1,w.value=``,A.value.is_authorized=!1,A.value.job_management_value_cash_id=``,A.value.job_management_value_authority_id=``,A.value.job_management_value_impact_id=``,a(`saved`),f.add({severity:`success`,summary:c(`message.success`),detail:c(`message.deleted`),life:2e3})}catch(e){k.value=e?.response?.data?.error?.message||c(`message.operation_failed`)}finally{O.value=!1}}}return n(async()=>{try{await Promise.all([H(),W()])}finally{h.value=!1}}),(e,t)=>(o(),S(`div`,Ln,[x(`div`,null,[x(`h2`,Rn,y(u(c)(`job_management.financials`)),1),x(`p`,zn,y(u(c)(`job_management.financial_description`)),1)]),x(`div`,null,[h.value?(o(),d(Y,{key:0,type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})):(o(),S(`div`,Bn,[x(`div`,Vn,[x(`div`,Hn,[x(`p`,Un,y(u(c)(`job_management.is_authorized`)),1),x(`p`,Wn,y(u(c)(`job_management.is_authorized_description`)),1)]),m(u(ee),{modelValue:A.value.is_authorized,"onUpdate:modelValue":t[0]||=e=>A.value.is_authorized=e},null,8,[`modelValue`])]),x(`div`,Gn,[A.value.is_authorized?(o(),d(R,{key:0,label:u(c)(`job_management.cash_level`),errors:b.value?.job_management_value_cash_id},{default:r(()=>[m(K,{modelValue:A.value.job_management_value_cash_id,"onUpdate:modelValue":t[1]||=e=>A.value.job_management_value_cash_id=e,options:M.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),class:C({"p-invalid":b.value?.job_management_value_cash_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`])):g(``,!0),m(R,{label:u(c)(`job_management.authority_level`),errors:b.value?.job_management_value_authority_id},{default:r(()=>[m(K,{modelValue:A.value.job_management_value_authority_id,"onUpdate:modelValue":t[2]||=e=>A.value.job_management_value_authority_id=e,options:B.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),class:C({"p-invalid":b.value?.job_management_value_authority_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),m(R,{label:u(c)(`job_management.impact_level`),errors:b.value?.job_management_value_impact_id},{default:r(()=>[m(K,{modelValue:A.value.job_management_value_impact_id,"onUpdate:modelValue":t[3]||=e=>A.value.job_management_value_impact_id=e,options:V.value,"option-label":`label`,"option-value":`value`,placeholder:u(c)(`common.select`),class:C({"p-invalid":b.value?.job_management_value_impact_id}),showClear:``},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`])]),_.value?(o(),S(`div`,Kn,y(_.value),1)):g(``,!0),x(`div`,qn,[w.value?(o(),d(u(D),{key:0,label:u(c)(`common.delete`),icon:`pi pi-trash`,severity:`danger`,size:`small`,outlined:``,onClick:t[4]||=e=>E.value=!0},null,8,[`label`])):g(``,!0),m(u(D),{label:w.value?u(c)(`common.update`):u(c)(`common.save`),icon:`pi pi-check`,size:`small`,loading:p.value,disabled:p.value,onClick:G},null,8,[`label`,`loading`,`disabled`])])]))]),m(q,{visible:E.value,"onUpdate:visible":t[5]||=e=>E.value=e,loading:O.value,"error-msg":k.value,onConfirm:J,onCancel:t[6]||=e=>E.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Xn={class:`space-y-4`},Zn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Qn={class:`text-sm text-gray-500 dark:text-gray-400`},$n={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},er={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},tr={class:`text-sm text-gray-500 dark:text-gray-400`},nr={key:0,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},rr={class:`flex-1 min-w-0`},ir={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},ar={key:0,class:`mt-0.5 text-xs text-gray-500 dark:text-gray-400 leading-relaxed`},or={class:`w-full md:w-80 shrink-0`},sr={key:1,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},cr={key:2,class:`flex justify-end gap-2 pt-1`},lr={class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},ur={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},dr={class:`text-sm text-gray-500 dark:text-gray-400`},fr={key:0,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},pr={class:`flex-1 min-w-0`},mr={class:`text-sm font-medium text-gray-800 dark:text-gray-100`},hr={key:0,class:`mt-0.5 text-xs text-gray-500 dark:text-gray-400 leading-relaxed`},gr={class:`w-full md:w-80 shrink-0`},_r={key:1,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},vr={key:2,class:`flex justify-end gap-2 pt-1`},yr={class:`max-w-2xl`},br={key:0,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},xr={key:1,class:`space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Sr={class:`flex items-center justify-between gap-2 flex-wrap`},Cr={class:`text-base font-semibold text-gray-800 dark:text-gray-100`},wr={class:`text-sm text-gray-500 dark:text-gray-400`},Tr={key:0,class:`text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2`},Er={class:`flex items-center justify-between`},Dr={class:`text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400`},Or={class:`grid grid-cols-1 md:grid-cols-2 gap-3`},kr={key:1,class:`text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2`},Ar={key:2,class:`flex justify-end gap-2 pt-2`},$=`/api/v1/tenant/job-management/potency-competencies`,jr={__name:`JobPotencySection`,props:{orgId:String,orgName:{type:String,default:``},orgCode:{type:String,default:``},jobValueMap:{type:Object,default:()=>({})},competencyOptions:{type:Array,default:()=>[]}},emits:[`saved`],setup(t,{emit:a}){let s=a,c=t,{t:l}=F(),f=j(),p=i(!0),h=i(!1),_=i(``),C=i([]),w=[],E=i([]),O=i(!1),k=v(()=>new Set((c.jobValueMap&&c.jobValueMap.kecerdasan||[]).map(e=>e.value))),A={tenacity:`tenacity`,"creativy & innovation":`innovation_creativity`,"creativity & innovation":`innovation_creativity`,"self confidence":`self_confidence`,flexibility:`flexibility`,"continuous learning":`continuous_learning`};function M(e){let t=(e||``).toLowerCase().replace(/\(.*?\)/g,``).replace(/\s+/g,` `).trim();return A[t]||``}function N(){let e=(c.competencyOptions||[]).filter(e=>M(e.label)).map(e=>{let t=M(e.label);return{competency_id:e.value,competency_name:e.label,competency_definition:e.definition||``,type:t,levelOptions:c.jobValueMap&&c.jobValueMap[t]||[],recordId:``,job_management_value_id:``}}),t=c.jobValueMap&&c.jobValueMap.kecerdasan||[];t.length>0&&e.unshift({competency_id:``,competency_name:l(`job_management.potency_kecerdasan`),competency_definition:``,type:`kecerdasan`,levelOptions:t,recordId:``,job_management_value_id:``}),E.value=e}function P(e){return w.find(t=>!t.competency_id&&e.value.has(t.job_management_value_id))||null}function L(){let e={};w.forEach(t=>{t.competency_id&&(e[t.competency_id]=t)}),E.value.forEach(t=>{let n=t.competency_id?e[t.competency_id]||null:P(k);t.recordId=n?n.id:``,t.job_management_value_id=n&&n.job_management_value_id||``})}async function z(e){_.value=``,O.value=!0;try{for(let t of e)if(t.job_management_value_id){let e=t.competency_id?{competency_id:t.competency_id,job_management_value_id:t.job_management_value_id}:{job_management_value_id:t.job_management_value_id};t.recordId?await T.put(`${$}/${t.recordId}`,e):t.recordId=(await T.post($,{organization_id:c.orgId,...e})).data?.data?.id||``}else t.recordId&&=(await T.delete(`${$}/${t.recordId}`),``);await X(),L(),q(),f.add({severity:`success`,summary:l(`message.success`),detail:l(`common.saved`),life:2e3}),s(`saved`)}catch(e){let t=I(e);Object.keys(t).length>0?_.value=Object.values(t).join(`, `):_.value=e?.response?.data?.error?.message||e.message||l(`message.operation_failed`)}finally{O.value=!1}}function B(){z(E.value)}function V(){z(U.value)}let U=i([]),W=v(()=>new Set((c.jobValueMap&&c.jobValueMap.communicating_influencing_skill||[]).map(e=>e.value)));function G(){let e=c.jobValueMap&&c.jobValueMap.communicating_influencing_skill||[];U.value=e.length>0?[{competency_id:``,competency_name:l(`job_management.skill_communicating_influencing`),competency_definition:``,type:`communicating_influencing_skill`,levelOptions:e,recordId:``,job_management_value_id:``}]:[]}function q(){U.value.forEach(e=>{let t=P(W);e.recordId=t?t.id:``,e.job_management_value_id=t&&t.job_management_value_id||``})}let J=v(()=>Object.values(c.jobValueMap||{}).flat()),ee=0;function te(){C.value.push({_key:`new-${++ee}`,id:``,competency_id:``,job_management_value_id:``,weight:null})}function ne(e){let t=C.value[e];t&&(t.id?re(t.id,e):C.value.splice(e,1))}async function re(e,t){try{await T.delete(`${$}/${e}`),C.value.splice(t,1),f.add({severity:`success`,summary:l(`message.success`),detail:l(`message.deleted`),life:2e3})}catch(e){f.add({severity:`error`,summary:l(`message.error`),detail:e?.response?.data?.error?.message||l(`message.operation_failed`),life:4e3})}}async function X(){if(!c.orgId){w=[],C.value=[];return}try{w=(await T.get($,{params:{organization_id:c.orgId,per_page:100}})).data?.data||[];let e=new Set(E.value.filter(e=>e.competency_id).map(e=>e.competency_id));C.value=w.filter(t=>!(t.competency_id&&e.has(t.competency_id))).filter(e=>!(!e.competency_id&&k.value.has(e.job_management_value_id))).filter(e=>!(!e.competency_id&&W.value.has(e.job_management_value_id))).map(e=>({_key:`db-${++ee}`,id:e.id,competency_id:e.competency_id||``,job_management_value_id:e.job_management_value_id||``,weight:e.weight??null}))}catch{w=[],C.value=[]}}async function ie(){_.value=``,h.value=!0;try{for(let e of C.value)e.id?await T.put(`${$}/${e.id}`,{competency_id:e.competency_id||``,job_management_value_id:e.job_management_value_id||``,weight:e.weight??null}):e.id=(await T.post($,{competency_id:e.competency_id||``,job_management_value_id:e.job_management_value_id||``,weight:e.weight??null,organization_id:c.orgId})).data?.data?.id||``;await X(),L(),q(),f.add({severity:`success`,summary:l(`message.success`),detail:l(`common.saved`),life:2e3}),s(`saved`)}catch(e){let t=I(e);Object.keys(t).length>0?_.value=Object.values(t).join(`, `):_.value=e?.response?.data?.error?.message||e.message||l(`message.operation_failed`)}finally{h.value=!1}}return n(async()=>{N(),G();try{await X()}finally{L(),q(),p.value=!1}}),(n,i)=>(o(),S(`div`,Xn,[x(`div`,null,[x(`h2`,Zn,y(u(l)(`job_management.potency_competencies`)),1),x(`p`,Qn,y(u(l)(`job_management.potency_description`)),1)]),x(`div`,$n,[x(`div`,null,[x(`h3`,er,y(u(l)(`job_management.potency_required_title`)),1),x(`p`,tr,y(u(l)(`job_management.potency_required_description`)),1)]),p.value?(o(),d(Y,{key:0,type:`detail`,count:1,rows:5,cols:`grid-cols-1`,padding:`p-5`})):(o(),S(b,{key:1},[E.value.length===0?(o(),S(`div`,nr,y(u(l)(`job_management.potency_required_empty`)),1)):g(``,!0),(o(!0),S(b,null,e(E.value,e=>(o(),S(`div`,{key:e.competency_id||e.type,class:`flex flex-col md:flex-row md:items-center gap-3 md:gap-6 py-3 border-b border-gray-100 dark:border-gray-700 last:border-b-0`},[x(`div`,rr,[x(`div`,ir,y(e.competency_name),1),e.competency_definition?(o(),S(`div`,ar,y(e.competency_definition),1)):g(``,!0)]),x(`div`,or,[m(K,{modelValue:e.job_management_value_id,"onUpdate:modelValue":t=>e.job_management_value_id=t,options:e.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:u(l)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])])]))),128)),_.value?(o(),S(`div`,sr,y(_.value),1)):g(``,!0),E.value.length>0?(o(),S(`div`,cr,[m(u(D),{label:u(l)(`job_management.save_potency_levels`),icon:`pi pi-check`,size:`small`,loading:O.value,disabled:O.value||!t.orgId,onClick:B},null,8,[`label`,`loading`,`disabled`])])):g(``,!0)],64))]),x(`div`,lr,[x(`div`,null,[x(`h3`,ur,y(u(l)(`job_management.skill_communicating_influencing_title`)),1),x(`p`,dr,y(u(l)(`job_management.skill_communicating_influencing_description`)),1)]),p.value?(o(),d(Y,{key:0,type:`detail`,count:1,rows:2,cols:`grid-cols-1`,padding:`p-5`})):(o(),S(b,{key:1},[U.value.length===0?(o(),S(`div`,fr,y(u(l)(`job_management.skill_communicating_influencing_empty`)),1)):g(``,!0),(o(!0),S(b,null,e(U.value,e=>(o(),S(`div`,{key:e.type,class:`flex flex-col md:flex-row md:items-center gap-3 md:gap-6 py-3 border-b border-gray-100 dark:border-gray-700 last:border-b-0`},[x(`div`,pr,[x(`div`,mr,y(e.competency_name),1),e.competency_definition?(o(),S(`div`,hr,y(e.competency_definition),1)):g(``,!0)]),x(`div`,gr,[m(K,{modelValue:e.job_management_value_id,"onUpdate:modelValue":t=>e.job_management_value_id=t,options:e.levelOptions,"option-label":`label`,"option-value":`value`,placeholder:u(l)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])])]))),128)),_.value?(o(),S(`div`,_r,y(_.value),1)):g(``,!0),U.value.length>0?(o(),S(`div`,vr,[m(u(D),{label:u(l)(`job_management.save_skill`),icon:`pi pi-check`,size:`small`,loading:O.value,disabled:O.value||!t.orgId,onClick:V},null,8,[`label`,`loading`,`disabled`])])):g(``,!0)],64))]),x(`div`,yr,[p.value?(o(),S(`div`,br,[m(Y,{type:`detail`,count:1,rows:4,cols:`grid-cols-1`,padding:`p-5`})])):(o(),S(`div`,xr,[x(`div`,Sr,[x(`div`,null,[x(`h3`,Cr,y(u(l)(`job_management.potency_competencies`)),1),x(`p`,wr,y(u(l)(`job_management.potency_description`)),1)]),m(u(D),{label:u(l)(`common.add`),icon:`pi pi-plus`,size:`small`,outlined:``,disabled:h.value,onClick:te},null,8,[`label`,`disabled`])]),C.value.length===0?(o(),S(`div`,Tr,y(u(l)(`job_management.empty_potency`)),1)):g(``,!0),(o(!0),S(b,null,e(C.value,(e,n)=>(o(),S(`div`,{key:e._key,class:`space-y-2`},[x(`div`,Er,[x(`span`,Dr,y(u(l)(`job_management.potency_item`))+` `+y(n+1),1),m(u(D),{icon:`pi pi-trash`,severity:`danger`,size:`small`,text:``,rounded:``,"aria-label":`Remove`,onClick:e=>ne(n)},null,8,[`onClick`])]),x(`div`,Or,[m(R,{label:u(l)(`job_management.competency`)},{default:r(()=>[m(K,{modelValue:e.competency_id,"onUpdate:modelValue":t=>e.competency_id=t,options:t.competencyOptions,"option-label":`label`,"option-value":`value`,placeholder:u(l)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),_:2},1032,[`label`]),m(R,{label:u(l)(`job_management.value_ref`)},{default:r(()=>[m(K,{modelValue:e.job_management_value_id,"onUpdate:modelValue":t=>e.job_management_value_id=t,options:J.value,"option-label":`label`,"option-value":`value`,placeholder:u(l)(`common.select`),showClear:``},null,8,[`modelValue`,`onUpdate:modelValue`,`options`,`placeholder`])]),_:2},1032,[`label`]),m(R,{label:u(l)(`job_management.weight`)},{default:r(()=>[m(u(H),{modelValue:e.weight,"onUpdate:modelValue":t=>e.weight=t,min:0,max:100,size:`small`,class:`w-full`},null,8,[`modelValue`,`onUpdate:modelValue`])]),_:2},1032,[`label`])])]))),128)),_.value?(o(),S(`div`,kr,y(_.value),1)):g(``,!0),C.value.length>0?(o(),S(`div`,Ar,[m(u(D),{label:u(l)(`job_management.save_potency`),icon:`pi pi-check`,size:`small`,loading:h.value,disabled:h.value,onClick:ie},null,8,[`label`,`loading`,`disabled`])])):g(``,!0)]))])]))}},Mr={class:`space-y-4`},Nr={class:`flex items-center justify-between`},Pr={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Fr={class:`text-sm text-gray-500 dark:text-gray-400`},Ir={class:`flex flex-col items-center justify-center py-10 text-gray-400`},Lr={class:`text-sm font-medium`},Rr=`/api/v1/tenant/job-management/competency-groups`,zr={__name:`JobCompetencyGroupSection`,props:{orgId:String},emits:[`saved`],setup(e,{emit:t}){let n=e,a=t,{t:s}=F(),c=j(),l=i([]),d=i(!1);i(0);let f=i(!1),p=i(!1),h=i(``),g=i(!1),_=i({}),b=i(!1),w=i(!1),E=i(``),O=i(null),k=i({category:``,weight:null}),A=v(()=>[{label:`${s(`job_management.technical`)} (${s(`job_management.category`)})`,value:`technical`},{label:`${s(`job_management.managerial`)} (${s(`job_management.category`)})`,value:`managerial`}]),M=v(()=>[{field:`category`,header:s(`job_management.category`)},{field:`weight`,header:s(`job_management.weight`)}]);async function N(){d.value=!0;try{let e=await T.get(Rr,{params:{organization_id:n.orgId}});l.value=e.data?.data||(Array.isArray(e.data)?e.data:[])}catch(e){c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.failed_to_load`),life:4e3})}finally{d.value=!1}}function P(){p.value=!1,h.value=``,k.value={category:`technical`,weight:null},_.value={},f.value=!0}function L(e){p.value=!0,h.value=e.id,k.value={category:e.category||`technical`,weight:e.weight??null},_.value={},f.value=!0}async function z(){g.value=!0,_.value={};try{let e={...k.value,organization_id:n.orgId};p.value?await T.put(`${Rr}/${h.value}`,e):await T.post(Rr,e),f.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.saved`),life:2e3}),N()}catch(e){let t=I(e);Object.keys(t).length?_.value=t:c.add({severity:`error`,detail:e.response?.data?.error?.message||s(`message.operation_failed`),life:4e3})}finally{g.value=!1}}function B(e){O.value=e,E.value=``,b.value=!0}async function V(){if(O.value){w.value=!0,E.value=``;try{await T.delete(`${Rr}/${O.value.id}`),b.value=!1,a(`saved`),c.add({severity:`success`,detail:s(`message.deleted`),life:2e3}),N()}catch(e){E.value=e.response?.data?.error?.message||s(`message.operation_failed`)}finally{w.value=!1}}}return(t,n)=>(o(),S(`div`,Mr,[x(`div`,Nr,[x(`div`,null,[x(`h2`,Pr,y(u(s)(`job_management.competency_groups`)),1),x(`p`,Fr,y(u(s)(`job_management.competency_group_description`)),1)]),m(u(D),{label:u(s)(`common.create`),icon:`pi pi-plus`,size:`small`,onClick:n[0]||=e=>P()},null,8,[`label`])]),m(Ze,{items:l.value,loading:d.value,total:l.value.length,columns:M.value,entity:`competency-groups`,"org-id":e.orgId,"on-load":N,onEdit:L,onDelete:B},{empty:r(()=>[x(`div`,Ir,[n[7]||=x(`i`,{class:`pi pi-chart-pie text-3xl mb-2 opacity-50`},null,-1),x(`p`,Lr,y(u(s)(`job_management.empty_competency_groups`)),1)])]),_:1},8,[`items`,`loading`,`total`,`columns`,`org-id`]),m(et,{visible:f.value,"onUpdate:visible":n[3]||=e=>f.value=e,title:p.value?u(s)(`common.edit`):u(s)(`common.create`),saving:g.value,errors:_.value,onSave:z,onCancel:n[4]||=e=>f.value=!1},{default:r(()=>[m(R,{label:u(s)(`job_management.category`),required:``,errors:_.value?.category},{default:r(()=>[m(K,{modelValue:k.value.category,"onUpdate:modelValue":n[1]||=e=>k.value.category=e,options:A.value,optionLabel:`label`,optionValue:`value`,placeholder:u(s)(`common.select`),class:C({"p-invalid":_.value?.category})},null,8,[`modelValue`,`options`,`placeholder`,`class`])]),_:1},8,[`label`,`errors`]),m(R,{label:u(s)(`job_management.weight`),required:``,errors:_.value?.weight},{default:r(()=>[m(u(H),{modelValue:k.value.weight,"onUpdate:modelValue":n[2]||=e=>k.value.weight=e,min:0,max:100,suffix:`%`,class:C([{"p-invalid":_.value?.weight},`w-full`]),size:`small`},null,8,[`modelValue`,`class`])]),_:1},8,[`label`,`errors`])]),_:1},8,[`visible`,`title`,`saving`,`errors`]),m(q,{visible:b.value,"onUpdate:visible":n[5]||=e=>b.value=e,loading:w.value,"error-msg":E.value,onConfirm:V,onCancel:n[6]||=e=>b.value=!1},null,8,[`visible`,`loading`,`error-msg`])]))}},Br={class:`space-y-6`},Vr={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Hr={class:`text-sm text-gray-500 dark:text-gray-400`},Ur={key:0,class:`flex items-center justify-center py-12`},Wr={class:`grid grid-cols-1 md:grid-cols-3 gap-4`},Gr={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Kr={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},qr={class:`text-2xl font-bold text-emerald-600 dark:text-emerald-400`},Jr={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Yr={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},Xr={class:`text-2xl font-bold text-blue-600 dark:text-blue-400`},Zr={class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5`},Qr={class:`text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1`},$r={key:0,class:`text-[10px] text-gray-400 mt-2`},ei={key:0,class:`bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden`},ti={class:`px-5 py-3 border-b border-gray-200 dark:border-gray-700 font-semibold text-sm text-gray-700 dark:text-gray-300`},ni={class:`p-5`},ri={class:`text-sm text-gray-700 dark:text-gray-300 capitalize`},ii={class:`text-sm font-semibold text-gray-900 dark:text-gray-100`},ai={key:2},oi={class:`flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500`},si={class:`text-sm font-medium`},ci={class:`text-xs mt-1`},li={class:`flex justify-end gap-3`},ui=`/api/v1/tenant/job-management/scores`,di={__name:`JobScoreSection`,props:{orgId:String},emits:[`saved`],setup(t,{emit:r}){let a=t,s=r,{t:c}=F(),l=j(),f=i(!1),p=i(!1),h=i(null),_=v(()=>{if(!h.value?.components)return null;try{return JSON.parse(h.value.components)}catch{return null}});function C(e){return e?.toLocaleString?.(`id-ID`)??`-`}async function w(){if(a.orgId){f.value=!0;try{let e=await T.get(`${ui}/${a.orgId}`);h.value=e.data?.data||null,s(`saved`)}catch{h.value=null}finally{f.value=!1}}}async function E(){if(a.orgId){p.value=!0;try{let e=await T.put(`${ui}/${a.orgId}`,{components:null});h.value=e.data?.data||null,l.add({severity:`success`,summary:c(`message.success`),detail:c(`job_management.score_calculated`),life:2e3})}catch(e){l.add({severity:`error`,summary:c(`message.error`),detail:e.response?.data?.error?.message||c(`message.operation_failed`),life:4e3})}finally{p.value=!1}}}return n(w),(t,n)=>(o(),S(`div`,Br,[x(`div`,null,[x(`h2`,Vr,y(u(c)(`job_management.scores`)),1),x(`p`,Hr,y(u(c)(`job_management.score_description`)),1)]),f.value?(o(),S(`div`,Ur,[...n[0]||=[x(`i`,{class:`pi pi-spin pi-spinner text-emerald-500 text-2xl`},null,-1)]])):h.value?(o(),S(b,{key:1},[x(`div`,Wr,[x(`div`,Gr,[x(`div`,Kr,y(u(c)(`job_management.value_with_financial`)),1),x(`div`,qr,y(C(h.value.job_value_with_financial)),1)]),x(`div`,Jr,[x(`div`,Yr,y(u(c)(`job_management.value_without_financial`)),1),x(`div`,Xr,y(C(h.value.job_value_without_financial)),1)]),x(`div`,Zr,[x(`div`,Qr,y(u(c)(`job_management.has_financial_authority`)),1),m(u(L),{value:h.value.has_financial_authority?u(c)(`common.yes`):u(c)(`common.no`),severity:h.value.has_financial_authority?`success`:`danger`,class:`!text-xs`},null,8,[`value`,`severity`]),h.value.calculated_at?(o(),S(`div`,$r,y(u(c)(`job_management.calculated_at`))+`: `+y(h.value.calculated_at),1)):g(``,!0)])]),_.value?(o(),S(`div`,ei,[x(`div`,ti,y(u(c)(`job_management.component_breakdown`)),1),x(`div`,ni,[(o(!0),S(b,null,e(_.value,(e,t)=>(o(),S(`div`,{key:t,class:`flex items-center justify-between py-1.5 border-b border-gray-100 dark:border-gray-700 last:border-0`},[x(`span`,ri,y(t.replace(/_/g,` `)),1),x(`span`,ii,y(C(e)),1)]))),128))])])):g(``,!0)],64)):(o(),S(`div`,ai,[x(`div`,oi,[n[1]||=x(`i`,{class:`pi pi-calculator text-4xl mb-3 opacity-50`},null,-1),x(`p`,si,y(u(c)(`job_management.no_score`)),1),x(`p`,ci,y(u(c)(`job_management.score_hint`)),1)])])),x(`div`,li,[m(u(D),{label:u(c)(`common.refresh`),icon:`pi pi-refresh`,size:`small`,text:``,onClick:w},null,8,[`label`]),h.value?(o(),d(u(D),{key:0,label:u(c)(`job_management.recalculate`),icon:`pi pi-calculator`,size:`small`,severity:`info`,loading:p.value,onClick:E},null,8,[`label`,`loading`])):g(``,!0)])]))}},fi={class:`max-w-full mx-auto`},pi={key:0,class:`flex gap-6`},mi={class:`w-56 space-y-2`},hi={class:`flex-1 space-y-3`},gi={key:1,class:`flex gap-6`},_i={class:`w-56 shrink-0 space-y-1`},vi=[`onClick`,`onKeydown`],yi={key:0,class:`pi pi-check text-xs`},bi={class:`flex-1 min-w-0`},xi={key:0,class:`pi pi-check-circle text-emerald-400 text-xs shrink-0`},Si={class:`flex-1 min-w-0`},Ci={__name:`JobManagementForm`,setup(t){let r=P(),a=A(),{t:c}=F(),l=j(),f=a.query.org_id||``,p=i(0),m=i(!0),h=i(Array(15).fill(!1)),_=i(``),w=i(``),D=i(``),O=i(``),k=i(``),M=i([]),N=i([]),I=i([]),L=i({}),R=i([]),z=[{labelKey:`job_management.identifications`,icon:`pi pi-id-card`,comp:ue},{labelKey:`job_management.objectives`,icon:`pi pi-bullseye`,comp:ye},{labelKey:`job_management.responsibilities_title`,icon:`pi pi-list-check`,comp:ct},{labelKey:`job_management.education_experience`,icon:`pi pi-graduation-cap`,comp:Ie},{labelKey:`job_management.potency_competencies`,icon:`pi pi-star`,comp:jr},{labelKey:`job_management.financials`,icon:`pi pi-money-bill`,comp:Yn},{labelKey:`job_management.assets`,icon:`pi pi-box`,comp:In},{labelKey:`job_management.subordinate_controls`,icon:`pi pi-sitemap`,comp:Dn},{labelKey:`job_management.relationships`,icon:`pi pi-share-alt`,comp:vn},{labelKey:`job_management.activities`,icon:`pi pi-bolt`,comp:Nt},{labelKey:`job_management.risks`,icon:`pi pi-exclamation-triangle`,comp:Ht},{labelKey:`job_management.hr_authorities`,icon:`pi pi-users`,comp:gt},{labelKey:`job_management.op_authorities`,icon:`pi pi-cog`,comp:wt},{labelKey:`job_management.competency_groups`,icon:`pi pi-chart-pie`,comp:zr},{labelKey:`job_management.scores`,icon:`pi pi-calculator`,comp:di}],B=v(()=>z[p.value]?.comp||null);function V(e){return p.value===e?`bg-emerald-50 dark:bg-emerald-900/20 ring-1 ring-emerald-300 dark:ring-emerald-700`:(h.value[e],`hover:bg-gray-50 dark:hover:bg-gray-800`)}function H(e){return p.value===e?`bg-emerald-600 text-white`:h.value[e]?`bg-emerald-100 dark:bg-emerald-800 text-emerald-600 dark:text-emerald-300`:`bg-gray-200 dark:bg-gray-600 text-gray-600 dark:text-gray-300`}function U(e){return p.value===e?`text-emerald-700 dark:text-emerald-300`:h.value[e]?`text-emerald-600 dark:text-emerald-400`:`text-gray-700 dark:text-gray-300`}function W(e){p.value=e,r.replace({query:{...a.query,section:String(e)}})}function G(e){typeof e==`number`&&(h.value[e]=!0)}async function K(){if(f)try{let e=(await T.get(`/api/v1/tenant/organizations/${f}`)).data?.data;e&&(_.value=e.nomenclature||``,w.value=e.full_code||e.code||``,D.value=e.organization_summary_id||``,O.value=e.grading_id||``,k.value=e.job_family_id||``)}catch{}}async function q(){try{let[e,t,n,r]=await Promise.all([T.get(`/api/v1/tenant/settings/gradings?per_page=100`),T.get(`/api/v1/tenant/job-management/values?per_page=200`),T.get(`/api/v1/tenant/competencies?per_page=200`).catch(()=>({data:{data:[]}})),T.get(`/api/v1/tenant/settings/job-families?per_page=100`)]);M.value=(e.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id})),N.value=(r.data?.data||[]).map(e=>({label:`${e.code} - ${e.name}`,value:e.id}));let i=t.data?.data||[];I.value=i.map(e=>({label:`${e.type}${e.level?` Lv.`+e.level:``}${e.descriptions?` — `+e.descriptions:``}`,value:e.id,type:e.type,level:e.level,descriptions:e.descriptions}));let a={};i.forEach(e=>{a[e.type]||(a[e.type]=[]),a[e.type].push({label:`Lv.${e.level} — ${e.descriptions||``}`,value:e.id,level:e.level})}),L.value=a,R.value=(n.data?.data||[]).map(e=>({label:e.name||e.code,value:e.id,field:e.field||``,definition:e.definition||``}))}catch{}}return n(async()=>{try{await Promise.all([K(),q()]);let e=parseInt(a.query.section);!isNaN(e)&&e>=0&&e<z.length&&(p.value=e)}catch(e){l.add({severity:`error`,summary:c(`message.error`),detail:e.response?.data?.error?.message||c(`message.failed_to_load`),life:4e3})}finally{m.value=!1}}),(t,n)=>(o(),S(`div`,fi,[m.value?(o(),S(`div`,pi,[x(`div`,mi,[(o(),S(b,null,e(8,e=>x(`div`,{key:e,class:`h-12 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))]),x(`div`,hi,[(o(),S(b,null,e(6,e=>x(`div`,{key:e,class:`h-9 bg-gray-200 dark:bg-gray-700 rounded animate-pulse`})),64))])])):(o(),S(`div`,gi,[x(`div`,_i,[(o(),S(b,null,e(z,(e,t)=>x(`div`,{key:t,role:`button`,tabindex:0,class:C([`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150 cursor-pointer select-none`,V(t)]),onClick:e=>W(t),onKeydown:E(e=>W(t),[`enter`])},[x(`div`,{class:C([`w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0 transition-colors duration-150`,H(t)])},[h.value[t]?(o(),S(`i`,yi)):(o(),S(`i`,{key:1,class:C(e.icon)},null,2))],2),x(`div`,bi,[x(`div`,{class:C([`text-sm font-medium truncate`,U(t)])},y(u(c)(e.labelKey)),3)]),h.value[t]?(o(),S(`i`,xi)):g(``,!0)],42,vi)),64))]),x(`div`,Si,[(o(),d(s(B.value),{key:p.value,"org-id":u(f),"org-name":_.value,"org-code":w.value,"org-summary-id":D.value,"org-grading-id":O.value,"org-job-family-id":k.value,"grading-options":M.value,"job-family-options":N.value,"job-value-options":I.value,"competency-options":R.value,"job-value-map":L.value,onSaved:G},null,40,[`org-id`,`org-name`,`org-code`,`org-summary-id`,`org-grading-id`,`org-job-family-id`,`grading-options`,`job-family-options`,`job-value-options`,`competency-options`,`job-value-map`]))])]))]))}};export{Ci as default};