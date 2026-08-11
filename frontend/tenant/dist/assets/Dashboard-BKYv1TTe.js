import{C as e,E as t,J as n,M as r,N as i,P as a,Q as o,U as s,W as c,c as l,f as u,ft as d,h as f,ht as p,j as m,k as h,l as g,m as _,o as v,r as y,s as b,u as x}from"./runtime-core.esm-bundler-huI9Rd5Y.js";import{Z as S,bt as C,mt as w,r as T,vt as E}from"./basecomponent-DqsHbXkj.js";import{i as D,t as O}from"./button-D_kLtROP.js";import{l as k,u as ee,y as te}from"./index-DIUIXv-J.js";import{t as A}from"./useI18n-BykaV4WF.js";import{n as j}from"./responseHandler-BJxA-JZj.js";import{t as M}from"./baseeditableholder-9FLZGK7-.js";var N=T.extend({name:`togglebutton`,style:`
    .p-togglebutton {
        display: inline-flex;
        cursor: pointer;
        user-select: none;
        overflow: hidden;
        position: relative;
        color: dt('togglebutton.color');
        background: dt('togglebutton.background');
        border: 1px solid dt('togglebutton.border.color');
        padding: dt('togglebutton.padding');
        font-size: 1rem;
        font-family: inherit;
        font-feature-settings: inherit;
        transition:
            background dt('togglebutton.transition.duration'),
            color dt('togglebutton.transition.duration'),
            border-color dt('togglebutton.transition.duration'),
            outline-color dt('togglebutton.transition.duration'),
            box-shadow dt('togglebutton.transition.duration');
        border-radius: dt('togglebutton.border.radius');
        outline-color: transparent;
        font-weight: dt('togglebutton.font.weight');
    }

    .p-togglebutton-content {
        display: inline-flex;
        flex: 1 1 auto;
        align-items: center;
        justify-content: center;
        gap: dt('togglebutton.gap');
        padding: dt('togglebutton.content.padding');
        background: transparent;
        border-radius: dt('togglebutton.content.border.radius');
        transition:
            background dt('togglebutton.transition.duration'),
            color dt('togglebutton.transition.duration'),
            border-color dt('togglebutton.transition.duration'),
            outline-color dt('togglebutton.transition.duration'),
            box-shadow dt('togglebutton.transition.duration');
    }

    .p-togglebutton:not(:disabled):not(.p-togglebutton-checked):hover {
        background: dt('togglebutton.hover.background');
        color: dt('togglebutton.hover.color');
    }

    .p-togglebutton.p-togglebutton-checked {
        background: dt('togglebutton.checked.background');
        border-color: dt('togglebutton.checked.border.color');
        color: dt('togglebutton.checked.color');
    }

    .p-togglebutton-checked .p-togglebutton-content {
        background: dt('togglebutton.content.checked.background');
        box-shadow: dt('togglebutton.content.checked.shadow');
    }

    .p-togglebutton:focus-visible {
        box-shadow: dt('togglebutton.focus.ring.shadow');
        outline: dt('togglebutton.focus.ring.width') dt('togglebutton.focus.ring.style') dt('togglebutton.focus.ring.color');
        outline-offset: dt('togglebutton.focus.ring.offset');
    }

    .p-togglebutton.p-invalid {
        border-color: dt('togglebutton.invalid.border.color');
    }

    .p-togglebutton:disabled {
        opacity: 1;
        cursor: default;
        background: dt('togglebutton.disabled.background');
        border-color: dt('togglebutton.disabled.border.color');
        color: dt('togglebutton.disabled.color');
    }

    .p-togglebutton-label,
    .p-togglebutton-icon {
        position: relative;
        transition: none;
    }

    .p-togglebutton-icon {
        color: dt('togglebutton.icon.color');
    }

    .p-togglebutton:not(:disabled):not(.p-togglebutton-checked):hover .p-togglebutton-icon {
        color: dt('togglebutton.icon.hover.color');
    }

    .p-togglebutton.p-togglebutton-checked .p-togglebutton-icon {
        color: dt('togglebutton.icon.checked.color');
    }

    .p-togglebutton:disabled .p-togglebutton-icon {
        color: dt('togglebutton.icon.disabled.color');
    }

    .p-togglebutton-sm {
        padding: dt('togglebutton.sm.padding');
        font-size: dt('togglebutton.sm.font.size');
    }

    .p-togglebutton-sm .p-togglebutton-content {
        padding: dt('togglebutton.content.sm.padding');
    }

    .p-togglebutton-lg {
        padding: dt('togglebutton.lg.padding');
        font-size: dt('togglebutton.lg.font.size');
    }

    .p-togglebutton-lg .p-togglebutton-content {
        padding: dt('togglebutton.content.lg.padding');
    }

    .p-togglebutton-fluid {
        width: 100%;
    }
`,classes:{root:function(e){var t=e.instance,n=e.props;return[`p-togglebutton p-component`,{"p-togglebutton-checked":t.active,"p-invalid":t.$invalid,"p-togglebutton-fluid":n.fluid,"p-togglebutton-sm p-inputfield-sm":n.size===`small`,"p-togglebutton-lg p-inputfield-lg":n.size===`large`}]},content:`p-togglebutton-content`,icon:`p-togglebutton-icon`,label:`p-togglebutton-label`}}),P={name:`BaseToggleButton`,extends:M,props:{onIcon:String,offIcon:String,onLabel:{type:String,default:`Yes`},offLabel:{type:String,default:`No`},readonly:{type:Boolean,default:!1},tabindex:{type:Number,default:null},ariaLabelledby:{type:String,default:null},ariaLabel:{type:String,default:null},size:{type:String,default:null},fluid:{type:Boolean,default:null}},style:N,provide:function(){return{$pcToggleButton:this,$parentInstance:this}}};function F(e){"@babel/helpers - typeof";return F=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},F(e)}function I(e,t,n){return(t=L(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function L(e){var t=R(e,`string`);return F(t)==`symbol`?t:t+``}function R(e,t){if(F(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(F(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var z={name:`ToggleButton`,extends:P,inheritAttrs:!1,emits:[`change`],methods:{getPTOptions:function(e){return(e===`root`?this.ptmi:this.ptm)(e,{context:{active:this.active,disabled:this.disabled}})},onChange:function(e){!this.disabled&&!this.readonly&&(this.writeValue(!this.d_value,e),this.$emit(`change`,e))},onBlur:function(e){var t,n;(t=(n=this.formField).onBlur)==null||t.call(n,e)}},computed:{active:function(){return this.d_value===!0},hasLabel:function(){return C(this.onLabel)&&C(this.offLabel)},label:function(){return this.hasLabel?this.d_value?this.onLabel:this.offLabel:`\xA0`},dataP:function(){return S(I({checked:this.active,invalid:this.$invalid},this.size,this.size))}},directives:{ripple:D}},B=[`tabindex`,`disabled`,`aria-pressed`,`aria-label`,`aria-labelledby`,`data-p-checked`,`data-p-disabled`,`data-p`],V=[`data-p`];function H(t,n,i,o,s,l){var u=a(`ripple`);return c((h(),x(`button`,e({type:`button`,class:t.cx(`root`),tabindex:t.tabindex,disabled:t.disabled,"aria-pressed":t.d_value,onClick:n[0]||=function(){return l.onChange&&l.onChange.apply(l,arguments)},onBlur:n[1]||=function(){return l.onBlur&&l.onBlur.apply(l,arguments)}},l.getPTOptions(`root`),{"aria-label":t.ariaLabel,"aria-labelledby":t.ariaLabelledby,"data-p-checked":l.active,"data-p-disabled":t.disabled,"data-p":l.dataP}),[b(`span`,e({class:t.cx(`content`)},l.getPTOptions(`content`),{"data-p":l.dataP}),[r(t.$slots,`default`,{},function(){return[r(t.$slots,`icon`,{value:t.d_value,class:d(t.cx(`icon`))},function(){return[t.onIcon||t.offIcon?(h(),x(`span`,e({key:0,class:[t.cx(`icon`),t.d_value?t.onIcon:t.offIcon]},l.getPTOptions(`icon`)),null,16)):g(``,!0)]}),b(`span`,e({class:t.cx(`label`)},l.getPTOptions(`label`)),p(l.label),17)]})],16,V)],16,B)),[[u]])}z.render=H;var U=T.extend({name:`selectbutton`,style:`
    .p-selectbutton {
        display: inline-flex;
        user-select: none;
        vertical-align: bottom;
        outline-color: transparent;
        border-radius: dt('selectbutton.border.radius');
    }

    .p-selectbutton .p-togglebutton {
        border-radius: 0;
        border-width: 1px 1px 1px 0;
    }

    .p-selectbutton .p-togglebutton:focus-visible {
        position: relative;
        z-index: 1;
    }

    .p-selectbutton .p-togglebutton:first-child {
        border-inline-start-width: 1px;
        border-start-start-radius: dt('selectbutton.border.radius');
        border-end-start-radius: dt('selectbutton.border.radius');
    }

    .p-selectbutton .p-togglebutton:last-child {
        border-start-end-radius: dt('selectbutton.border.radius');
        border-end-end-radius: dt('selectbutton.border.radius');
    }

    .p-selectbutton.p-invalid {
        outline: 1px solid dt('selectbutton.invalid.border.color');
        outline-offset: 0;
    }

    .p-selectbutton-fluid {
        width: 100%;
    }
    
    .p-selectbutton-fluid .p-togglebutton {
        flex: 1 1 0;
    }
`,classes:{root:function(e){var t=e.props;return[`p-selectbutton p-component`,{"p-invalid":e.instance.$invalid,"p-selectbutton-fluid":t.fluid}]}}}),W={name:`BaseSelectButton`,extends:M,props:{options:Array,optionLabel:null,optionValue:null,optionDisabled:null,multiple:Boolean,allowEmpty:{type:Boolean,default:!0},dataKey:null,ariaLabelledby:{type:String,default:null},size:{type:String,default:null},fluid:{type:Boolean,default:null}},style:U,provide:function(){return{$pcSelectButton:this,$parentInstance:this}}};function G(e,t){var n=typeof Symbol<`u`&&e[Symbol.iterator]||e[`@@iterator`];if(!n){if(Array.isArray(e)||(n=J(e))||t){n&&(e=n);var r=0,i=function(){};return{s:i,n:function(){return r>=e.length?{done:!0}:{done:!1,value:e[r++]}},e:function(e){throw e},f:i}}throw TypeError(`Invalid attempt to iterate non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}var a,o=!0,s=!1;return{s:function(){n=n.call(e)},n:function(){var e=n.next();return o=e.done,e},e:function(e){s=!0,a=e},f:function(){try{o||n.return==null||n.return()}finally{if(s)throw a}}}}function K(e){return X(e)||Y(e)||J(e)||q()}function q(){throw TypeError(`Invalid attempt to spread non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}function J(e,t){if(e){if(typeof e==`string`)return Z(e,t);var n={}.toString.call(e).slice(8,-1);return n===`Object`&&e.constructor&&(n=e.constructor.name),n===`Map`||n===`Set`?Array.from(e):n===`Arguments`||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?Z(e,t):void 0}}function Y(e){if(typeof Symbol<`u`&&e[Symbol.iterator]!=null||e[`@@iterator`]!=null)return Array.from(e)}function X(e){if(Array.isArray(e))return Z(e)}function Z(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=Array(t);n<t;n++)r[n]=e[n];return r}var Q={name:`SelectButton`,extends:W,inheritAttrs:!1,emits:[`change`],methods:{getOptionLabel:function(e){return this.optionLabel?E(e,this.optionLabel):e},getOptionValue:function(e){return this.optionValue?E(e,this.optionValue):e},getOptionRenderKey:function(e){return this.dataKey?E(e,this.dataKey):this.getOptionLabel(e)},isOptionDisabled:function(e){return this.optionDisabled?E(e,this.optionDisabled):!1},isOptionReadonly:function(e){if(this.allowEmpty)return!1;var t=this.isSelected(e);return this.multiple?t&&this.d_value.length===1:t},onOptionSelect:function(e,t,n){var r=this;if(!(this.disabled||this.isOptionDisabled(t)||this.isOptionReadonly(t))){var i=this.isSelected(t),a=this.getOptionValue(t),o;if(this.multiple)if(i){if(o=this.d_value.filter(function(e){return!w(e,a,r.equalityKey)}),!this.allowEmpty&&o.length===0)return}else o=this.d_value?[].concat(K(this.d_value),[a]):[a];else{if(i&&!this.allowEmpty)return;o=i?null:a}this.writeValue(o,e),this.$emit(`change`,{originalEvent:e,value:o})}},isSelected:function(e){var t=!1,n=this.getOptionValue(e);if(this.multiple){if(this.d_value){var r=G(this.d_value),i;try{for(r.s();!(i=r.n()).done;){var a=i.value;if(w(a,n,this.equalityKey)){t=!0;break}}}catch(e){r.e(e)}finally{r.f()}}}else t=w(this.d_value,n,this.equalityKey);return t}},computed:{equalityKey:function(){return this.optionValue?null:this.dataKey},dataP:function(){return S({invalid:this.$invalid})}},directives:{ripple:D},components:{ToggleButton:z}},ne=[`aria-labelledby`,`data-p`];function re(t,n,a,o,c,d){var f=i(`ToggleButton`);return h(),x(`div`,e({class:t.cx(`root`),role:`group`,"aria-labelledby":t.ariaLabelledby},t.ptmi(`root`),{"data-p":d.dataP}),[(h(!0),x(y,null,m(t.options,function(n,i){return h(),l(f,{key:d.getOptionRenderKey(n),modelValue:d.isSelected(n),onLabel:d.getOptionLabel(n),offLabel:d.getOptionLabel(n),disabled:t.disabled||d.isOptionDisabled(n),unstyled:t.unstyled,size:t.size,readonly:d.isOptionReadonly(n),onChange:function(e){return d.onOptionSelect(e,n,i)},pt:t.ptm(`pcToggleButton`)},u({_:2},[t.$slots.option?{name:`default`,fn:s(function(){return[r(t.$slots,`option`,{option:n,index:i},function(){return[b(`span`,e({ref_for:!0},t.ptm(`pcToggleButton`).label),p(d.getOptionLabel(n)),17)]})]}),key:`0`}:void 0]),1032,[`modelValue`,`onLabel`,`offLabel`,`disabled`,`unstyled`,`size`,`readonly`,`onChange`,`pt`])}),128))],16,ne)}Q.render=re;var ie={class:`space-y-4`},ae={class:`flex items-center justify-end`},oe={class:`flex items-center gap-2`},se={class:`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3`},ce={class:`flex items-center justify-between mb-2`},le={class:`text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider`},ue={class:`text-xl font-bold text-gray-800 dark:text-gray-100`},de={class:`flex items-center gap-1 mt-1`},fe={class:`text-sm text-gray-400 dark:text-gray-500`},pe={key:0,class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 animate-pulse`},me={class:`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2`},he={key:1,class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3`},ge={class:`flex items-center justify-between gap-2 flex-wrap mb-3`},_e={class:`flex items-center gap-2`},ve={class:`text-sm font-semibold text-gray-700 dark:text-gray-200`},ye={class:`grid grid-cols-1 lg:grid-cols-3 gap-4`},be={class:`lg:col-span-2`},xe={class:`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2`},Se={class:`min-w-0`},Ce={class:`text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate`},we={class:`text-lg font-bold text-gray-800 dark:text-gray-100`},Te={class:`space-y-3`},Ee={class:`rounded-lg border border-amber-300 dark:border-amber-700/50 bg-amber-50/50 dark:bg-amber-900/10 p-3`},De={class:`text-xs font-medium text-amber-600 dark:text-amber-400 uppercase tracking-wider flex items-center gap-1.5`},Oe={class:`text-2xl font-bold text-amber-700 dark:text-amber-300`},ke={class:`rounded-lg border border-sky-300 dark:border-sky-700/50 bg-sky-50/50 dark:bg-sky-900/10 p-3`},Ae={class:`text-xs font-medium text-sky-600 dark:text-sky-400 uppercase tracking-wider flex items-center gap-1.5`},je={class:`text-2xl font-bold text-sky-700 dark:text-sky-300`},Me={class:`mt-3 pt-3 border-t border-gray-100 dark:border-gray-800`},Ne={class:`flex items-center gap-2 mb-2`},Pe={class:`text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider`},Fe={class:`grid grid-cols-3 gap-2`},Ie={class:`rounded-lg border border-gray-200 dark:border-gray-700 p-2.5`},Le={class:`text-[11px] text-gray-500 dark:text-gray-400`},Re={class:`text-lg font-bold text-gray-800 dark:text-gray-100`},ze={class:`rounded-lg border border-amber-300 dark:border-amber-700/50 bg-amber-50/50 dark:bg-amber-900/10 p-2.5`},Be={class:`text-[11px] text-amber-600 dark:text-amber-400`},Ve={class:`text-lg font-bold text-amber-700 dark:text-amber-300`},$={class:`rounded-lg border border-gray-200 dark:border-gray-700 p-2.5`},He={class:`text-[11px] text-gray-500 dark:text-gray-400`},Ue={class:`text-lg font-bold text-gray-800 dark:text-gray-100`},We={class:`grid grid-cols-1 lg:grid-cols-3 gap-4`},Ge={class:`lg:col-span-2 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3`},Ke={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3`},qe={class:`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2`},Je=[`onClick`],Ye={class:`text-sm text-gray-600 dark:text-gray-300 text-center leading-tight`},Xe={class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3`},Ze={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3`},Qe={class:`space-y-3`},$e={class:`min-w-0`},et={class:`text-sm text-gray-700 dark:text-gray-200`},tt={class:`text-[11px] text-gray-400 dark:text-gray-500 mt-0.5`},nt={__name:`Dashboard`,setup(e){let{t:r}=A(),i=te(),a=k(),s=n(`this-month`),c=v(()=>[{label:r(`dashboard.this_month`),value:`this-month`},{label:r(`dashboard.this_quarter`),value:`this-quarter`},{label:r(`dashboard.this_year`),value:`this-year`}]),l=v(()=>[{label:r(`dashboard.kpi_total_employees`),value:`1,247`,icon:`pi pi-users`,iconColor:`text-emerald-500`,trend:3.2},{label:r(`dashboard.kpi_active_today`),value:`1,183`,icon:`pi pi-check-circle`,iconColor:`text-blue-500`,trend:1.5},{label:r(`dashboard.kpi_on_leave`),value:`42`,icon:`pi pi-calendar`,iconColor:`text-amber-500`,trend:-2.1},{label:r(`dashboard.kpi_pending_approvals`),value:`28`,icon:`pi pi-clock`,iconColor:`text-rose-500`,trend:12.5}]),u=v(()=>[{name:r(`dashboard.employees`),icon:`pi pi-users`,route:`/employees`,bg:`bg-blue-50`,color:`text-blue-600`},{name:r(`dashboard.attendance`),icon:`pi pi-clock`,route:`/attendance`,bg:`bg-emerald-50`,color:`text-emerald-600`},{name:r(`dashboard.leave`),icon:`pi pi-calendar`,route:`/leave`,bg:`bg-amber-50`,color:`text-amber-600`},{name:r(`dashboard.payroll`),icon:`pi pi-dollar`,route:`/payroll`,bg:`bg-indigo-50`,color:`text-indigo-600`},{name:r(`dashboard.approvals`),icon:`pi pi-check-square`,route:`/approvals`,bg:`bg-violet-50`,color:`text-violet-600`},{name:r(`dashboard.performance`),icon:`pi pi-chart-line`,route:`/performance`,bg:`bg-cyan-50`,color:`text-cyan-600`},{name:r(`dashboard.training`),icon:`pi pi-book`,route:`/training`,bg:`bg-orange-50`,color:`text-orange-600`},{name:r(`dashboard.recruitment`),icon:`pi pi-user-plus`,route:`/recruitment`,bg:`bg-rose-50`,color:`text-rose-600`},{name:r(`dashboard.organization`),icon:`pi pi-sitemap`,route:`/organizations`,bg:`bg-teal-50`,color:`text-teal-600`},{name:r(`dashboard.reimbursement`),icon:`pi pi-credit-card`,route:`/reimbursements`,bg:`bg-sky-50`,color:`text-sky-600`},{name:r(`dashboard.workforce_intel`),icon:`pi pi-chart-bar`,route:`/workforce-intelligence`,bg:`bg-slate-50`,color:`text-slate-600`},{name:r(`dashboard.career_intel`),icon:`pi pi-chart-line`,route:`/career-intelligence`,bg:`bg-pink-50`,color:`text-pink-600`}]),S=[{text:`15 new employees added this week`,time:`2 hours ago`,dotColor:`bg-emerald-400`},{text:`Payroll run for August completed`,time:`5 hours ago`,dotColor:`bg-blue-400`},{text:`3 leave requests pending approval`,time:`1 day ago`,dotColor:`bg-amber-400`},{text:`Performance reviews Q3 initiated`,time:`2 days ago`,dotColor:`bg-violet-400`},{text:`Training session "Leadership 101" scheduled`,time:`3 days ago`,dotColor:`bg-orange-400`}],C=n(!1),w=n(!1),T=n({movement_by_type:{},pending_approval:0,effective_this_month:0,contracts:{}}),E=v(()=>[`promotion`,`demotion`,`mutation`,`contract_extension`,`status_change`,`retirement`,`offboarding`,`other`].map(e=>({label:D(e),value:e})));function D(e){let t=`employee_movement.type_${e}`;return r(t)===t?e:r(t)}function M(e){switch(e){case`promotion`:return`pi pi-arrow-up`;case`demotion`:return`pi pi-arrow-down`;case`mutation`:return`pi pi-shuffle`;case`contract_extension`:return`pi pi-file-edit`;case`status_change`:return`pi pi-id-card`;case`retirement`:return`pi pi-sun`;case`offboarding`:return`pi pi-sign-out`;default:return`pi pi-circle`}}function N(e){switch(e){case`promotion`:return`text-emerald-500`;case`demotion`:return`text-red-500`;case`mutation`:return`text-sky-500`;case`contract_extension`:return`text-amber-500`;case`status_change`:return`text-indigo-500`;case`retirement`:return`text-gray-400`;case`offboarding`:return`text-red-400`;default:return`text-gray-400`}}async function P(){w.value=!0;try{let e=await ee.get(`/api/v1/tenant/employee-movements/dashboard`);T.value=e.data?.data||T.value}catch(e){i.add({severity:`error`,summary:r(`message.error`),detail:j(e,r(`message.failed_to_load`)),life:4e3})}finally{w.value=!1}}return t(async()=>{await a.fetchActiveModules(),C.value=a.hasModule(`employeemovement`),C.value&&P()}),(e,t)=>(h(),x(`div`,ie,[b(`div`,ae,[b(`div`,oe,[f(o(Q),{modelValue:s.value,"onUpdate:modelValue":t[0]||=e=>s.value=e,options:c.value,optionLabel:`label`,optionValue:`value`,size:`small`},null,8,[`modelValue`,`options`])])]),b(`div`,se,[(h(!0),x(y,null,m(l.value,e=>(h(),x(`div`,{key:e.label,class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 hover:shadow-sm dark:hover:shadow-gray-900/50 transition-shadow`},[b(`div`,ce,[b(`span`,le,p(e.label),1),b(`i`,{class:d([[e.icon,e.iconColor],`text-lg`])},null,2)]),b(`div`,ue,p(e.value),1),b(`div`,de,[b(`i`,{class:d([e.trend>=0?`pi pi-arrow-up text-emerald-500`:`pi pi-arrow-down text-rose-500`,`text-sm`])},null,2),b(`span`,{class:d([e.trend>=0?`text-emerald-600`:`text-rose-600`,`text-sm font-medium`])},p(Math.abs(e.trend))+`% `,3),b(`span`,fe,p(o(r)(`dashboard.vs_last_month`)),1)])]))),128))]),C.value&&w.value?(h(),x(`div`,pe,[t[2]||=b(`div`,{class:`flex items-center gap-2 mb-3`},[b(`div`,{class:`w-4 h-4 rounded bg-gray-200 dark:bg-gray-700`}),b(`div`,{class:`h-4 w-40 rounded bg-gray-200 dark:bg-gray-700`})],-1),b(`div`,me,[(h(),x(y,null,m(8,e=>b(`div`,{key:e,class:`h-16 rounded-lg bg-gray-100 dark:bg-gray-700/50`})),64))])])):g(``,!0),C.value&&!w.value?(h(),x(`div`,he,[b(`div`,ge,[b(`div`,_e,[t[3]||=b(`i`,{class:`pi pi-arrows-alt text-sm text-emerald-500`},null,-1),b(`h2`,ve,p(o(r)(`dashboard.movement_title`)),1)]),f(o(O),{label:o(r)(`dashboard.view_reports`),icon:`pi pi-chart-bar`,size:`small`,text:``,class:`!text-xs`,onClick:t[1]||=t=>e.$router.push(`/admin/career/reports`)},null,8,[`label`])]),b(`div`,ye,[b(`div`,be,[b(`div`,xe,[(h(!0),x(y,null,m(E.value,e=>(h(),x(`div`,{key:e.value,class:`rounded-lg border border-gray-200 dark:border-gray-700 p-2.5 flex items-center justify-between gap-2 hover:shadow-sm dark:hover:shadow-gray-900/50 transition-shadow`},[b(`div`,Se,[b(`p`,Ce,p(e.label),1),b(`p`,we,p(T.value.movement_by_type?.[e.value]||0),1)]),b(`i`,{class:d([[M(e.value),N(e.value)],`text-base shrink-0`])},null,2)]))),128))])]),b(`div`,Te,[b(`div`,Ee,[b(`p`,De,[t[4]||=b(`i`,{class:`pi pi-clock text-xs`},null,-1),_(p(o(r)(`dashboard.movement_pending`)),1)]),b(`p`,Oe,p(T.value.pending_approval||0),1)]),b(`div`,ke,[b(`p`,Ae,[t[5]||=b(`i`,{class:`pi pi-calendar text-xs`},null,-1),_(p(o(r)(`dashboard.movement_effective_month`)),1)]),b(`p`,je,p(T.value.effective_this_month||0),1)])])]),b(`div`,Me,[b(`div`,Ne,[t[6]||=b(`i`,{class:`pi pi-file-edit text-xs text-gray-400`},null,-1),b(`span`,Pe,p(o(r)(`dashboard.contract_title`)),1)]),b(`div`,Fe,[b(`div`,Ie,[b(`p`,Le,p(o(r)(`dashboard.contract_active`)),1),b(`p`,Re,p(T.value.contracts?.active||0),1)]),b(`div`,ze,[b(`p`,Be,p(o(r)(`dashboard.contract_expiring`)),1),b(`p`,Ve,p(T.value.contracts?.expiring||0),1)]),b(`div`,$,[b(`p`,He,p(o(r)(`dashboard.contract_expired`)),1),b(`p`,Ue,p(T.value.contracts?.expired||0),1)])])])])):g(``,!0),b(`div`,We,[b(`div`,Ge,[b(`h2`,Ke,p(o(r)(`dashboard.quick_access`)),1),b(`div`,qe,[(h(!0),x(y,null,m(u.value,t=>(h(),x(`div`,{key:t.name,class:`flex flex-col items-center gap-1.5 p-2.5 rounded-lg cursor-pointer hover:bg-emerald-50 dark:hover:bg-emerald-900/20 hover:border-emerald-200 dark:hover:border-emerald-700 border border-transparent transition-all`,onClick:n=>e.$router.push(t.route)},[b(`div`,{class:d([t.bg,`w-9 h-9 rounded-lg flex items-center justify-center`])},[b(`i`,{class:d([[t.icon,t.color],`text-sm`])},null,2)],2),b(`span`,Ye,p(t.name),1)],8,Je))),128))])]),b(`div`,Xe,[b(`h2`,Ze,p(o(r)(`dashboard.recent_activity`)),1),b(`div`,Qe,[(h(),x(y,null,m(S,(e,t)=>b(`div`,{key:t,class:`flex items-start gap-2.5`},[b(`div`,{class:d([e.dotColor,`w-2 h-2 rounded-full mt-1.5 shrink-0`])},null,2),b(`div`,$e,[b(`p`,et,p(e.text),1),b(`p`,tt,p(e.time),1)])])),64))])])])]))}};export{nt as default};