import{$ as e,A as t,D as n,F as r,G as i,M as a,N as o,P as s,W as c,Y as l,c as u,d,g as f,gt as p,h as m,i as h,l as g,p as _,pt as v,s as y,u as b,w as x}from"./language-8RKtU9ID.js";import{Z as S,bt as C,mt as w,r as T,vt as E}from"./basecomponent-Dq2wVn8v.js";import{t as D}from"./ripple-DWGM_HLq.js";import{f as O,p as k,r as A,x as j}from"./index-C_5oItlt.js";import{t as M}from"./useI18n-CeYrJQoe.js";import{n as N}from"./responseHandler-BJxA-JZj.js";import{t as P}from"./baseeditableholder-DSjqJDX7.js";var F=T.extend({name:`togglebutton`,style:`
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
`,classes:{root:function(e){var t=e.instance,n=e.props;return[`p-togglebutton p-component`,{"p-togglebutton-checked":t.active,"p-invalid":t.$invalid,"p-togglebutton-fluid":n.fluid,"p-togglebutton-sm p-inputfield-sm":n.size===`small`,"p-togglebutton-lg p-inputfield-lg":n.size===`large`}]},content:`p-togglebutton-content`,icon:`p-togglebutton-icon`,label:`p-togglebutton-label`}}),I={name:`BaseToggleButton`,extends:P,props:{onIcon:String,offIcon:String,onLabel:{type:String,default:`Yes`},offLabel:{type:String,default:`No`},readonly:{type:Boolean,default:!1},tabindex:{type:Number,default:null},ariaLabelledby:{type:String,default:null},ariaLabel:{type:String,default:null},size:{type:String,default:null},fluid:{type:Boolean,default:null}},style:F,provide:function(){return{$pcToggleButton:this,$parentInstance:this}}};function L(e){"@babel/helpers - typeof";return L=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},L(e)}function R(e,t,n){return(t=z(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function z(e){var t=B(e,`string`);return L(t)==`symbol`?t:t+``}function B(e,t){if(L(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(L(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var V={name:`ToggleButton`,extends:I,inheritAttrs:!1,emits:[`change`],methods:{getPTOptions:function(e){return(e===`root`?this.ptmi:this.ptm)(e,{context:{active:this.active,disabled:this.disabled}})},onChange:function(e){!this.disabled&&!this.readonly&&(this.writeValue(!this.d_value,e),this.$emit(`change`,e))},onBlur:function(e){var t,n;(t=(n=this.formField).onBlur)==null||t.call(n,e)}},computed:{active:function(){return this.d_value===!0},hasLabel:function(){return C(this.onLabel)&&C(this.offLabel)},label:function(){return this.hasLabel?this.d_value?this.onLabel:this.offLabel:`\xA0`},dataP:function(){return S(R({checked:this.active,invalid:this.$invalid},this.size,this.size))}},directives:{ripple:D}},H=[`tabindex`,`disabled`,`aria-pressed`,`aria-label`,`aria-labelledby`,`data-p-checked`,`data-p-disabled`,`data-p`],ee=[`data-p`];function te(e,n,a,s,c,l){var f=r(`ripple`);return i((t(),d(`button`,x({type:`button`,class:e.cx(`root`),tabindex:e.tabindex,disabled:e.disabled,"aria-pressed":e.d_value,onClick:n[0]||=function(){return l.onChange&&l.onChange.apply(l,arguments)},onBlur:n[1]||=function(){return l.onBlur&&l.onBlur.apply(l,arguments)}},l.getPTOptions(`root`),{"aria-label":e.ariaLabel,"aria-labelledby":e.ariaLabelledby,"data-p-checked":l.active,"data-p-disabled":e.disabled,"data-p":l.dataP}),[u(`span`,x({class:e.cx(`content`)},l.getPTOptions(`content`),{"data-p":l.dataP}),[o(e.$slots,`default`,{},function(){return[o(e.$slots,`icon`,{value:e.d_value,class:v(e.cx(`icon`))},function(){return[e.onIcon||e.offIcon?(t(),d(`span`,x({key:0,class:[e.cx(`icon`),e.d_value?e.onIcon:e.offIcon]},l.getPTOptions(`icon`)),null,16)):b(``,!0)]}),u(`span`,x({class:e.cx(`label`)},l.getPTOptions(`label`)),p(l.label),17)]})],16,ee)],16,H)),[[f]])}V.render=te;var U=T.extend({name:`selectbutton`,style:`
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
`,classes:{root:function(e){var t=e.props;return[`p-selectbutton p-component`,{"p-invalid":e.instance.$invalid,"p-selectbutton-fluid":t.fluid}]}}}),W={name:`BaseSelectButton`,extends:P,props:{options:Array,optionLabel:null,optionValue:null,optionDisabled:null,multiple:Boolean,allowEmpty:{type:Boolean,default:!0},dataKey:null,ariaLabelledby:{type:String,default:null},size:{type:String,default:null},fluid:{type:Boolean,default:null}},style:U,provide:function(){return{$pcSelectButton:this,$parentInstance:this}}};function G(e,t){var n=typeof Symbol<`u`&&e[Symbol.iterator]||e[`@@iterator`];if(!n){if(Array.isArray(e)||(n=J(e))||t){n&&(e=n);var r=0,i=function(){};return{s:i,n:function(){return r>=e.length?{done:!0}:{done:!1,value:e[r++]}},e:function(e){throw e},f:i}}throw TypeError(`Invalid attempt to iterate non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}var a,o=!0,s=!1;return{s:function(){n=n.call(e)},n:function(){var e=n.next();return o=e.done,e},e:function(e){s=!0,a=e},f:function(){try{o||n.return==null||n.return()}finally{if(s)throw a}}}}function K(e){return X(e)||Y(e)||J(e)||q()}function q(){throw TypeError(`Invalid attempt to spread non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}function J(e,t){if(e){if(typeof e==`string`)return Z(e,t);var n={}.toString.call(e).slice(8,-1);return n===`Object`&&e.constructor&&(n=e.constructor.name),n===`Map`||n===`Set`?Array.from(e):n===`Arguments`||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?Z(e,t):void 0}}function Y(e){if(typeof Symbol<`u`&&e[Symbol.iterator]!=null||e[`@@iterator`]!=null)return Array.from(e)}function X(e){if(Array.isArray(e))return Z(e)}function Z(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=Array(t);n<t;n++)r[n]=e[n];return r}var Q={name:`SelectButton`,extends:W,inheritAttrs:!1,emits:[`change`],methods:{getOptionLabel:function(e){return this.optionLabel?E(e,this.optionLabel):e},getOptionValue:function(e){return this.optionValue?E(e,this.optionValue):e},getOptionRenderKey:function(e){return this.dataKey?E(e,this.dataKey):this.getOptionLabel(e)},isOptionDisabled:function(e){return this.optionDisabled?E(e,this.optionDisabled):!1},isOptionReadonly:function(e){if(this.allowEmpty)return!1;var t=this.isSelected(e);return this.multiple?t&&this.d_value.length===1:t},onOptionSelect:function(e,t,n){var r=this;if(!(this.disabled||this.isOptionDisabled(t)||this.isOptionReadonly(t))){var i=this.isSelected(t),a=this.getOptionValue(t),o;if(this.multiple)if(i){if(o=this.d_value.filter(function(e){return!w(e,a,r.equalityKey)}),!this.allowEmpty&&o.length===0)return}else o=this.d_value?[].concat(K(this.d_value),[a]):[a];else{if(i&&!this.allowEmpty)return;o=i?null:a}this.writeValue(o,e),this.$emit(`change`,{originalEvent:e,value:o})}},isSelected:function(e){var t=!1,n=this.getOptionValue(e);if(this.multiple){if(this.d_value){var r=G(this.d_value),i;try{for(r.s();!(i=r.n()).done;){var a=i.value;if(w(a,n,this.equalityKey)){t=!0;break}}}catch(e){r.e(e)}finally{r.f()}}}else t=w(this.d_value,n,this.equalityKey);return t}},computed:{equalityKey:function(){return this.optionValue?null:this.dataKey},dataP:function(){return S({invalid:this.$invalid})}},directives:{ripple:D},components:{ToggleButton:V}},ne=[`aria-labelledby`,`data-p`];function re(e,n,r,i,l,f){var m=s(`ToggleButton`);return t(),d(`div`,x({class:e.cx(`root`),role:`group`,"aria-labelledby":e.ariaLabelledby},e.ptmi(`root`),{"data-p":f.dataP}),[(t(!0),d(h,null,a(e.options,function(n,r){return t(),g(m,{key:f.getOptionRenderKey(n),modelValue:f.isSelected(n),onLabel:f.getOptionLabel(n),offLabel:f.getOptionLabel(n),disabled:e.disabled||f.isOptionDisabled(n),unstyled:e.unstyled,size:e.size,readonly:f.isOptionReadonly(n),onChange:function(e){return f.onOptionSelect(e,n,r)},pt:e.ptm(`pcToggleButton`)},_({_:2},[e.$slots.option?{name:`default`,fn:c(function(){return[o(e.$slots,`option`,{option:n,index:r},function(){return[u(`span`,x({ref_for:!0},e.ptm(`pcToggleButton`).label),p(f.getOptionLabel(n)),17)]})]}),key:`0`}:void 0]),1032,[`modelValue`,`onLabel`,`offLabel`,`disabled`,`unstyled`,`size`,`readonly`,`onChange`,`pt`])}),128))],16,ne)}Q.render=re;var ie={class:`space-y-4`},ae={class:`flex items-center justify-end`},oe={class:`flex items-center gap-2`},se={class:`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3`},ce={class:`flex items-center justify-between mb-2`},le={class:`text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider`},ue={class:`text-xl font-bold text-gray-800 dark:text-gray-100`},de={class:`flex items-center gap-1 mt-1`},fe={class:`text-sm text-gray-400 dark:text-gray-500`},pe={key:0,class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 animate-pulse`},me={class:`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2`},he={key:1,class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3`},ge={class:`flex items-center justify-between gap-2 flex-wrap mb-3`},_e={class:`flex items-center gap-2`},ve={class:`text-sm font-semibold text-gray-700 dark:text-gray-200`},ye={class:`grid grid-cols-1 lg:grid-cols-3 gap-4`},be={class:`lg:col-span-2`},xe={class:`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2`},Se={class:`min-w-0`},Ce={class:`text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate`},we={class:`text-lg font-bold text-gray-800 dark:text-gray-100`},Te={class:`space-y-3`},Ee={class:`rounded-lg border border-amber-300 dark:border-amber-700/50 bg-amber-50/50 dark:bg-amber-900/10 p-3`},De={class:`text-xs font-medium text-amber-600 dark:text-amber-400 uppercase tracking-wider flex items-center gap-1.5`},Oe={class:`text-2xl font-bold text-amber-700 dark:text-amber-300`},ke={class:`rounded-lg border border-sky-300 dark:border-sky-700/50 bg-sky-50/50 dark:bg-sky-900/10 p-3`},Ae={class:`text-xs font-medium text-sky-600 dark:text-sky-400 uppercase tracking-wider flex items-center gap-1.5`},je={class:`text-2xl font-bold text-sky-700 dark:text-sky-300`},Me={class:`mt-3 pt-3 border-t border-gray-100 dark:border-gray-800`},Ne={class:`flex items-center gap-2 mb-2`},Pe={class:`text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider`},Fe={class:`grid grid-cols-3 gap-2`},Ie={class:`rounded-lg border border-gray-200 dark:border-gray-700 p-2.5`},Le={class:`text-[11px] text-gray-500 dark:text-gray-400`},Re={class:`text-lg font-bold text-gray-800 dark:text-gray-100`},ze={class:`rounded-lg border border-amber-300 dark:border-amber-700/50 bg-amber-50/50 dark:bg-amber-900/10 p-2.5`},Be={class:`text-[11px] text-amber-600 dark:text-amber-400`},Ve={class:`text-lg font-bold text-amber-700 dark:text-amber-300`},He={class:`rounded-lg border border-gray-200 dark:border-gray-700 p-2.5`},Ue={class:`text-[11px] text-gray-500 dark:text-gray-400`},We={class:`text-lg font-bold text-gray-800 dark:text-gray-100`},Ge={key:2,class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 animate-pulse`},Ke={class:`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-2`},qe={key:3,class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3`},Je={class:`flex items-center justify-between gap-2 flex-wrap mb-3`},Ye={class:`flex items-center gap-2`},Xe={class:`text-sm font-semibold text-gray-700 dark:text-gray-200`},Ze={class:`text-[11px] text-gray-400 dark:text-gray-500 hidden sm:inline`},Qe={class:`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-2`},$e={class:`rounded-lg border border-emerald-200 dark:border-emerald-700/50 bg-emerald-50/50 dark:bg-emerald-900/10 p-2.5`},et={class:`text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider`},tt={class:`text-lg font-bold text-emerald-600 dark:text-emerald-400`},nt={class:`rounded-lg border border-gray-200 dark:border-gray-700 p-2.5`},rt={class:`text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider`},it={class:`text-lg font-bold text-gray-800 dark:text-gray-100`},at={class:`rounded-lg border border-gray-200 dark:border-gray-700 p-2.5`},ot={class:`text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider`},st={class:`text-lg font-bold text-gray-800 dark:text-gray-100`},ct={class:`rounded-lg border border-gray-200 dark:border-gray-700 p-2.5`},lt={class:`text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider`},ut={class:`text-lg font-bold text-gray-800 dark:text-gray-100`},dt={class:`rounded-lg border border-gray-200 dark:border-gray-700 p-2.5`},ft={class:`text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider`},$={class:`text-lg font-bold text-gray-800 dark:text-gray-100`},pt={class:`grid grid-cols-1 lg:grid-cols-3 gap-4`},mt={class:`lg:col-span-2 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3`},ht={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3`},gt={class:`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2`},_t=[`onClick`],vt={class:`text-sm text-gray-600 dark:text-gray-300 text-center leading-tight`},yt={class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3`},bt={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3`},xt={class:`space-y-3`},St={class:`min-w-0`},Ct={class:`text-sm text-gray-700 dark:text-gray-200`},wt={class:`text-[11px] text-gray-400 dark:text-gray-500 mt-0.5`},Tt={__name:`Dashboard`,setup(r){let{t:i}=M(),o=j(),s=O(),c=l(`this-month`),g=y(()=>[{label:i(`dashboard.this_month`),value:`this-month`},{label:i(`dashboard.this_quarter`),value:`this-quarter`},{label:i(`dashboard.this_year`),value:`this-year`}]),_=y(()=>[{label:i(`dashboard.kpi_total_employees`),value:`1,247`,icon:`pi pi-users`,iconColor:`text-emerald-500`,trend:3.2},{label:i(`dashboard.kpi_active_today`),value:`1,183`,icon:`pi pi-check-circle`,iconColor:`text-blue-500`,trend:1.5},{label:i(`dashboard.kpi_on_leave`),value:`42`,icon:`pi pi-calendar`,iconColor:`text-amber-500`,trend:-2.1},{label:i(`dashboard.kpi_pending_approvals`),value:`28`,icon:`pi pi-clock`,iconColor:`text-rose-500`,trend:12.5}]),x=y(()=>[{name:i(`dashboard.employees`),icon:`pi pi-users`,route:`/employees`,bg:`bg-blue-50`,color:`text-blue-600`},{name:i(`dashboard.attendance`),icon:`pi pi-clock`,route:`/attendance`,bg:`bg-emerald-50`,color:`text-emerald-600`},{name:i(`dashboard.leave`),icon:`pi pi-calendar`,route:`/leave`,bg:`bg-amber-50`,color:`text-amber-600`},{name:i(`dashboard.payroll`),icon:`pi pi-dollar`,route:`/payroll`,bg:`bg-indigo-50`,color:`text-indigo-600`},{name:i(`dashboard.approvals`),icon:`pi pi-check-square`,route:`/approvals`,bg:`bg-violet-50`,color:`text-violet-600`},{name:i(`dashboard.performance`),icon:`pi pi-chart-line`,route:`/performance`,bg:`bg-cyan-50`,color:`text-cyan-600`},{name:i(`dashboard.training`),icon:`pi pi-book`,route:`/training`,bg:`bg-orange-50`,color:`text-orange-600`},{name:i(`dashboard.recruitment`),icon:`pi pi-user-plus`,route:`/recruitment`,bg:`bg-rose-50`,color:`text-rose-600`},{name:i(`dashboard.organization`),icon:`pi pi-sitemap`,route:`/organizations`,bg:`bg-teal-50`,color:`text-teal-600`},{name:i(`dashboard.reimbursement`),icon:`pi pi-credit-card`,route:`/reimbursements`,bg:`bg-sky-50`,color:`text-sky-600`},{name:i(`dashboard.workforce_intel`),icon:`pi pi-chart-bar`,route:`/workforce-intelligence`,bg:`bg-slate-50`,color:`text-slate-600`},{name:i(`dashboard.career_intel`),icon:`pi pi-chart-line`,route:`/career-intelligence`,bg:`bg-pink-50`,color:`text-pink-600`}]),S=[{text:`15 new employees added this week`,time:`2 hours ago`,dotColor:`bg-emerald-400`},{text:`Payroll run for August completed`,time:`5 hours ago`,dotColor:`bg-blue-400`},{text:`3 leave requests pending approval`,time:`1 day ago`,dotColor:`bg-amber-400`},{text:`Performance reviews Q3 initiated`,time:`2 days ago`,dotColor:`bg-violet-400`},{text:`Training session "Leadership 101" scheduled`,time:`3 days ago`,dotColor:`bg-orange-400`}],C=l(!1),w=l(!1),T=l({overall_score:0,hires_analyzed:0,interview_score:0,onboarding_completion_rate:0,retention_rate:0});function E(e){return(Number(e)||0).toFixed(1)}function D(e){return`${(Number(e)||0).toFixed(1)}%`}async function P(){w.value=!0;try{let e=await k.get(`/api/v1/tenant/workforce-intelligence/analytics/quality-of-hire`);T.value=e.data?.data||T.value}catch(e){o.add({severity:`error`,summary:i(`message.error`),detail:N(e,i(`message.failed_to_load`)),life:4e3})}finally{w.value=!1}}let F=l(!1),I=l(!1),L=l({movement_by_type:{},pending_approval:0,effective_this_month:0,contracts:{}}),R=y(()=>[`promotion`,`demotion`,`mutation`,`contract_extension`,`status_change`,`retirement`,`offboarding`,`other`].map(e=>({label:z(e),value:e})));function z(e){let t=`employee_movement.type_${e}`;return i(t)===t?e:i(t)}function B(e){switch(e){case`promotion`:return`pi pi-arrow-up`;case`demotion`:return`pi pi-arrow-down`;case`mutation`:return`pi pi-shuffle`;case`contract_extension`:return`pi pi-file-edit`;case`status_change`:return`pi pi-id-card`;case`retirement`:return`pi pi-sun`;case`offboarding`:return`pi pi-sign-out`;default:return`pi pi-circle`}}function V(e){switch(e){case`promotion`:return`text-emerald-500`;case`demotion`:return`text-red-500`;case`mutation`:return`text-sky-500`;case`contract_extension`:return`text-amber-500`;case`status_change`:return`text-indigo-500`;case`retirement`:return`text-gray-400`;case`offboarding`:return`text-red-400`;default:return`text-gray-400`}}async function H(){I.value=!0;try{let e=await k.get(`/api/v1/tenant/employee-movements/dashboard`);L.value=e.data?.data||L.value}catch(e){o.add({severity:`error`,summary:i(`message.error`),detail:N(e,i(`message.failed_to_load`)),life:4e3})}finally{I.value=!1}}return n(async()=>{await s.fetchActiveModules(),F.value=s.hasModule(`employeemovement`),F.value&&H(),C.value=s.hasModule(`workforce-intelligence`),C.value&&P()}),(n,r)=>(t(),d(`div`,ie,[u(`div`,ae,[u(`div`,oe,[f(e(Q),{modelValue:c.value,"onUpdate:modelValue":r[0]||=e=>c.value=e,options:g.value,optionLabel:`label`,optionValue:`value`,size:`small`},null,8,[`modelValue`,`options`])])]),u(`div`,se,[(t(!0),d(h,null,a(_.value,n=>(t(),d(`div`,{key:n.label,class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 hover:shadow-sm dark:hover:shadow-gray-900/50 transition-shadow`},[u(`div`,ce,[u(`span`,le,p(n.label),1),u(`i`,{class:v([[n.icon,n.iconColor],`text-lg`])},null,2)]),u(`div`,ue,p(n.value),1),u(`div`,de,[u(`i`,{class:v([n.trend>=0?`pi pi-arrow-up text-emerald-500`:`pi pi-arrow-down text-rose-500`,`text-sm`])},null,2),u(`span`,{class:v([n.trend>=0?`text-emerald-600`:`text-rose-600`,`text-sm font-medium`])},p(Math.abs(n.trend))+`% `,3),u(`span`,fe,p(e(i)(`dashboard.vs_last_month`)),1)])]))),128))]),F.value&&I.value?(t(),d(`div`,pe,[r[3]||=u(`div`,{class:`flex items-center gap-2 mb-3`},[u(`div`,{class:`w-4 h-4 rounded bg-gray-200 dark:bg-gray-700`}),u(`div`,{class:`h-4 w-40 rounded bg-gray-200 dark:bg-gray-700`})],-1),u(`div`,me,[(t(),d(h,null,a(8,e=>u(`div`,{key:e,class:`h-16 rounded-lg bg-gray-100 dark:bg-gray-700/50`})),64))])])):b(``,!0),F.value&&!I.value?(t(),d(`div`,he,[u(`div`,ge,[u(`div`,_e,[r[4]||=u(`i`,{class:`pi pi-arrows-alt text-sm text-emerald-500`},null,-1),u(`h2`,ve,p(e(i)(`dashboard.movement_title`)),1)]),f(e(A),{label:e(i)(`dashboard.view_reports`),icon:`pi pi-chart-bar`,size:`small`,text:``,class:`!text-xs`,onClick:r[1]||=e=>n.$router.push(`/admin/career/reports`)},null,8,[`label`])]),u(`div`,ye,[u(`div`,be,[u(`div`,xe,[(t(!0),d(h,null,a(R.value,e=>(t(),d(`div`,{key:e.value,class:`rounded-lg border border-gray-200 dark:border-gray-700 p-2.5 flex items-center justify-between gap-2 hover:shadow-sm dark:hover:shadow-gray-900/50 transition-shadow`},[u(`div`,Se,[u(`p`,Ce,p(e.label),1),u(`p`,we,p(L.value.movement_by_type?.[e.value]||0),1)]),u(`i`,{class:v([[B(e.value),V(e.value)],`text-base shrink-0`])},null,2)]))),128))])]),u(`div`,Te,[u(`div`,Ee,[u(`p`,De,[r[5]||=u(`i`,{class:`pi pi-clock text-xs`},null,-1),m(p(e(i)(`dashboard.movement_pending`)),1)]),u(`p`,Oe,p(L.value.pending_approval||0),1)]),u(`div`,ke,[u(`p`,Ae,[r[6]||=u(`i`,{class:`pi pi-calendar text-xs`},null,-1),m(p(e(i)(`dashboard.movement_effective_month`)),1)]),u(`p`,je,p(L.value.effective_this_month||0),1)])])]),u(`div`,Me,[u(`div`,Ne,[r[7]||=u(`i`,{class:`pi pi-file-edit text-xs text-gray-400`},null,-1),u(`span`,Pe,p(e(i)(`dashboard.contract_title`)),1)]),u(`div`,Fe,[u(`div`,Ie,[u(`p`,Le,p(e(i)(`dashboard.contract_active`)),1),u(`p`,Re,p(L.value.contracts?.active||0),1)]),u(`div`,ze,[u(`p`,Be,p(e(i)(`dashboard.contract_expiring`)),1),u(`p`,Ve,p(L.value.contracts?.expiring||0),1)]),u(`div`,He,[u(`p`,Ue,p(e(i)(`dashboard.contract_expired`)),1),u(`p`,We,p(L.value.contracts?.expired||0),1)])])])])):b(``,!0),C.value&&w.value?(t(),d(`div`,Ge,[r[8]||=u(`div`,{class:`flex items-center gap-2 mb-3`},[u(`div`,{class:`w-4 h-4 rounded bg-gray-200 dark:bg-gray-700`}),u(`div`,{class:`h-4 w-44 rounded bg-gray-200 dark:bg-gray-700`})],-1),u(`div`,Ke,[(t(),d(h,null,a(5,e=>u(`div`,{key:e,class:`h-16 rounded-lg bg-gray-100 dark:bg-gray-700/50`})),64))])])):b(``,!0),C.value&&!w.value&&T.value.hires_analyzed>0?(t(),d(`div`,qe,[u(`div`,Je,[u(`div`,Ye,[r[9]||=u(`i`,{class:`pi pi-bullseye text-sm text-emerald-500`},null,-1),u(`h2`,Xe,p(e(i)(`dashboard.quality_of_hire_title`)),1),u(`span`,Ze,p(e(i)(`dashboard.quality_of_hire_desc`)),1)]),f(e(A),{label:e(i)(`dashboard.view_analytics`),icon:`pi pi-chart-bar`,size:`small`,text:``,class:`!text-xs`,onClick:r[2]||=e=>n.$router.push(`/workforce-intelligence/quality-of-hire`)},null,8,[`label`])]),u(`div`,Qe,[u(`div`,$e,[u(`p`,et,p(e(i)(`quality_of_hire.overall_score`)),1),u(`p`,tt,p(E(T.value.overall_score)),1)]),u(`div`,nt,[u(`p`,rt,p(e(i)(`quality_of_hire.hires_analyzed`)),1),u(`p`,it,p(T.value.hires_analyzed),1)]),u(`div`,at,[u(`p`,ot,p(e(i)(`quality_of_hire.interview_score`)),1),u(`p`,st,p(E(T.value.interview_score)),1)]),u(`div`,ct,[u(`p`,lt,p(e(i)(`quality_of_hire.onboarding_completion`)),1),u(`p`,ut,p(D(T.value.onboarding_completion_rate)),1)]),u(`div`,dt,[u(`p`,ft,p(e(i)(`quality_of_hire.retention_rate`)),1),u(`p`,$,p(D(T.value.retention_rate)),1)])])])):b(``,!0),u(`div`,pt,[u(`div`,mt,[u(`h2`,ht,p(e(i)(`dashboard.quick_access`)),1),u(`div`,gt,[(t(!0),d(h,null,a(x.value,e=>(t(),d(`div`,{key:e.name,class:`flex flex-col items-center gap-1.5 p-2.5 rounded-lg cursor-pointer hover:bg-emerald-50 dark:hover:bg-emerald-900/20 hover:border-emerald-200 dark:hover:border-emerald-700 border border-transparent transition-all`,onClick:t=>n.$router.push(e.route)},[u(`div`,{class:v([e.bg,`w-9 h-9 rounded-lg flex items-center justify-center`])},[u(`i`,{class:v([[e.icon,e.color],`text-sm`])},null,2)],2),u(`span`,vt,p(e.name),1)],8,_t))),128))])]),u(`div`,yt,[u(`h2`,bt,p(e(i)(`dashboard.recent_activity`)),1),u(`div`,xt,[(t(),d(h,null,a(S,(e,t)=>u(`div`,{key:t,class:`flex items-start gap-2.5`},[u(`div`,{class:v([e.dotColor,`w-2 h-2 rounded-full mt-1.5 shrink-0`])},null,2),u(`div`,St,[u(`p`,Ct,p(e.text),1),u(`p`,wt,p(e.time),1)])])),64))])])])]))}};export{Tt as default};