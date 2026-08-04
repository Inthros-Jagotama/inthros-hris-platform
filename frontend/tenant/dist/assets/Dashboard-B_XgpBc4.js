import{A as e,C as t,H as n,K as r,M as i,N as a,O as o,U as s,X as c,c as l,f as u,h as d,j as f,l as p,o as m,pt as h,r as g,s as _,u as v,ut as y}from"./runtime-core.esm-bundler-CFm0BMYx.js";import{a as b,et as x,gt as S,pt as C,t as w,vt as T}from"./ripple-FLjJmYYY.js";import{t as E}from"./useI18n-kftu6zy2.js";import{t as D}from"./baseeditableholder-CnB2pcI8.js";var O=b.extend({name:`togglebutton`,style:`
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
`,classes:{root:function(e){var t=e.instance,n=e.props;return[`p-togglebutton p-component`,{"p-togglebutton-checked":t.active,"p-invalid":t.$invalid,"p-togglebutton-fluid":n.fluid,"p-togglebutton-sm p-inputfield-sm":n.size===`small`,"p-togglebutton-lg p-inputfield-lg":n.size===`large`}]},content:`p-togglebutton-content`,icon:`p-togglebutton-icon`,label:`p-togglebutton-label`}}),k={name:`BaseToggleButton`,extends:D,props:{onIcon:String,offIcon:String,onLabel:{type:String,default:`Yes`},offLabel:{type:String,default:`No`},readonly:{type:Boolean,default:!1},tabindex:{type:Number,default:null},ariaLabelledby:{type:String,default:null},ariaLabel:{type:String,default:null},size:{type:String,default:null},fluid:{type:Boolean,default:null}},style:O,provide:function(){return{$pcToggleButton:this,$parentInstance:this}}};function A(e){"@babel/helpers - typeof";return A=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},A(e)}function ee(e,t,n){return(t=j(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function j(e){var t=M(e,`string`);return A(t)==`symbol`?t:t+``}function M(e,t){if(A(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(A(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var N={name:`ToggleButton`,extends:k,inheritAttrs:!1,emits:[`change`],methods:{getPTOptions:function(e){return(e===`root`?this.ptmi:this.ptm)(e,{context:{active:this.active,disabled:this.disabled}})},onChange:function(e){!this.disabled&&!this.readonly&&(this.writeValue(!this.d_value,e),this.$emit(`change`,e))},onBlur:function(e){var t,n;(t=(n=this.formField).onBlur)==null||t.call(n,e)}},computed:{active:function(){return this.d_value===!0},hasLabel:function(){return T(this.onLabel)&&T(this.offLabel)},label:function(){return this.hasLabel?this.d_value?this.onLabel:this.offLabel:`\xA0`},dataP:function(){return x(ee({checked:this.active,invalid:this.$invalid},this.size,this.size))}},directives:{ripple:w}},P=[`tabindex`,`disabled`,`aria-pressed`,`aria-label`,`aria-labelledby`,`data-p-checked`,`data-p-disabled`,`data-p`],F=[`data-p`];function I(e,n,r,i,c,l){var u=a(`ripple`);return s((o(),v(`button`,t({type:`button`,class:e.cx(`root`),tabindex:e.tabindex,disabled:e.disabled,"aria-pressed":e.d_value,onClick:n[0]||=function(){return l.onChange&&l.onChange.apply(l,arguments)},onBlur:n[1]||=function(){return l.onBlur&&l.onBlur.apply(l,arguments)}},l.getPTOptions(`root`),{"aria-label":e.ariaLabel,"aria-labelledby":e.ariaLabelledby,"data-p-checked":l.active,"data-p-disabled":e.disabled,"data-p":l.dataP}),[_(`span`,t({class:e.cx(`content`)},l.getPTOptions(`content`),{"data-p":l.dataP}),[f(e.$slots,`default`,{},function(){return[f(e.$slots,`icon`,{value:e.d_value,class:y(e.cx(`icon`))},function(){return[e.onIcon||e.offIcon?(o(),v(`span`,t({key:0,class:[e.cx(`icon`),e.d_value?e.onIcon:e.offIcon]},l.getPTOptions(`icon`)),null,16)):p(``,!0)]}),_(`span`,t({class:e.cx(`label`)},l.getPTOptions(`label`)),h(l.label),17)]})],16,F)],16,P)),[[u]])}N.render=I;var L=b.extend({name:`selectbutton`,style:`
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
`,classes:{root:function(e){var t=e.props;return[`p-selectbutton p-component`,{"p-invalid":e.instance.$invalid,"p-selectbutton-fluid":t.fluid}]}}}),R={name:`BaseSelectButton`,extends:D,props:{options:Array,optionLabel:null,optionValue:null,optionDisabled:null,multiple:Boolean,allowEmpty:{type:Boolean,default:!0},dataKey:null,ariaLabelledby:{type:String,default:null},size:{type:String,default:null},fluid:{type:Boolean,default:null}},style:L,provide:function(){return{$pcSelectButton:this,$parentInstance:this}}};function z(e,t){var n=typeof Symbol<`u`&&e[Symbol.iterator]||e[`@@iterator`];if(!n){if(Array.isArray(e)||(n=H(e))||t){n&&(e=n);var r=0,i=function(){};return{s:i,n:function(){return r>=e.length?{done:!0}:{done:!1,value:e[r++]}},e:function(e){throw e},f:i}}throw TypeError(`Invalid attempt to iterate non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}var a,o=!0,s=!1;return{s:function(){n=n.call(e)},n:function(){var e=n.next();return o=e.done,e},e:function(e){s=!0,a=e},f:function(){try{o||n.return==null||n.return()}finally{if(s)throw a}}}}function B(e){return W(e)||U(e)||H(e)||V()}function V(){throw TypeError(`Invalid attempt to spread non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}function H(e,t){if(e){if(typeof e==`string`)return G(e,t);var n={}.toString.call(e).slice(8,-1);return n===`Object`&&e.constructor&&(n=e.constructor.name),n===`Map`||n===`Set`?Array.from(e):n===`Arguments`||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?G(e,t):void 0}}function U(e){if(typeof Symbol<`u`&&e[Symbol.iterator]!=null||e[`@@iterator`]!=null)return Array.from(e)}function W(e){if(Array.isArray(e))return G(e)}function G(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=Array(t);n<t;n++)r[n]=e[n];return r}var K={name:`SelectButton`,extends:R,inheritAttrs:!1,emits:[`change`],methods:{getOptionLabel:function(e){return this.optionLabel?S(e,this.optionLabel):e},getOptionValue:function(e){return this.optionValue?S(e,this.optionValue):e},getOptionRenderKey:function(e){return this.dataKey?S(e,this.dataKey):this.getOptionLabel(e)},isOptionDisabled:function(e){return this.optionDisabled?S(e,this.optionDisabled):!1},isOptionReadonly:function(e){if(this.allowEmpty)return!1;var t=this.isSelected(e);return this.multiple?t&&this.d_value.length===1:t},onOptionSelect:function(e,t,n){var r=this;if(!(this.disabled||this.isOptionDisabled(t)||this.isOptionReadonly(t))){var i=this.isSelected(t),a=this.getOptionValue(t),o;if(this.multiple)if(i){if(o=this.d_value.filter(function(e){return!C(e,a,r.equalityKey)}),!this.allowEmpty&&o.length===0)return}else o=this.d_value?[].concat(B(this.d_value),[a]):[a];else{if(i&&!this.allowEmpty)return;o=i?null:a}this.writeValue(o,e),this.$emit(`change`,{originalEvent:e,value:o})}},isSelected:function(e){var t=!1,n=this.getOptionValue(e);if(this.multiple){if(this.d_value){var r=z(this.d_value),i;try{for(r.s();!(i=r.n()).done;){var a=i.value;if(C(a,n,this.equalityKey)){t=!0;break}}}catch(e){r.e(e)}finally{r.f()}}}else t=C(this.d_value,n,this.equalityKey);return t}},computed:{equalityKey:function(){return this.optionValue?null:this.dataKey},dataP:function(){return x({invalid:this.$invalid})}},directives:{ripple:w},components:{ToggleButton:N}},q=[`aria-labelledby`,`data-p`];function J(r,a,s,c,d,p){var m=i(`ToggleButton`);return o(),v(`div`,t({class:r.cx(`root`),role:`group`,"aria-labelledby":r.ariaLabelledby},r.ptmi(`root`),{"data-p":p.dataP}),[(o(!0),v(g,null,e(r.options,function(e,i){return o(),l(m,{key:p.getOptionRenderKey(e),modelValue:p.isSelected(e),onLabel:p.getOptionLabel(e),offLabel:p.getOptionLabel(e),disabled:r.disabled||p.isOptionDisabled(e),unstyled:r.unstyled,size:r.size,readonly:p.isOptionReadonly(e),onChange:function(t){return p.onOptionSelect(t,e,i)},pt:r.ptm(`pcToggleButton`)},u({_:2},[r.$slots.option?{name:`default`,fn:n(function(){return[f(r.$slots,`option`,{option:e,index:i},function(){return[_(`span`,t({ref_for:!0},r.ptm(`pcToggleButton`).label),h(p.getOptionLabel(e)),17)]})]}),key:`0`}:void 0]),1032,[`modelValue`,`onLabel`,`offLabel`,`disabled`,`unstyled`,`size`,`readonly`,`onChange`,`pt`])}),128))],16,q)}K.render=J;var Y={class:`space-y-4`},X={class:`flex items-center justify-end`},te={class:`flex items-center gap-2`},Z={class:`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3`},Q={class:`flex items-center justify-between mb-2`},ne={class:`text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider`},re={class:`text-xl font-bold text-gray-800 dark:text-gray-100`},ie={class:`flex items-center gap-1 mt-1`},ae={class:`text-sm text-gray-400 dark:text-gray-500`},oe={class:`grid grid-cols-1 lg:grid-cols-3 gap-4`},se={class:`lg:col-span-2 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3`},ce={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3`},le={class:`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2`},ue=[`onClick`],de={class:`text-sm text-gray-600 dark:text-gray-300 text-center leading-tight`},fe={class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3`},pe={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3`},$={class:`space-y-3`},me={class:`min-w-0`},he={class:`text-sm text-gray-700 dark:text-gray-200`},ge={class:`text-[11px] text-gray-400 dark:text-gray-500 mt-0.5`},_e={__name:`Dashboard`,setup(t){let{t:n}=E(),i=r(`this-month`),a=m(()=>[{label:n(`dashboard.this_month`),value:`this-month`},{label:n(`dashboard.this_quarter`),value:`this-quarter`},{label:n(`dashboard.this_year`),value:`this-year`}]),s=m(()=>[{label:n(`dashboard.kpi_total_employees`),value:`1,247`,icon:`pi pi-users`,iconColor:`text-emerald-500`,trend:3.2},{label:n(`dashboard.kpi_active_today`),value:`1,183`,icon:`pi pi-check-circle`,iconColor:`text-blue-500`,trend:1.5},{label:n(`dashboard.kpi_on_leave`),value:`42`,icon:`pi pi-calendar`,iconColor:`text-amber-500`,trend:-2.1},{label:n(`dashboard.kpi_pending_approvals`),value:`28`,icon:`pi pi-clock`,iconColor:`text-rose-500`,trend:12.5}]),l=m(()=>[{name:n(`dashboard.employees`),icon:`pi pi-users`,route:`/employees`,bg:`bg-blue-50`,color:`text-blue-600`},{name:n(`dashboard.attendance`),icon:`pi pi-clock`,route:`/attendance`,bg:`bg-emerald-50`,color:`text-emerald-600`},{name:n(`dashboard.leave`),icon:`pi pi-calendar`,route:`/leave`,bg:`bg-amber-50`,color:`text-amber-600`},{name:n(`dashboard.payroll`),icon:`pi pi-dollar`,route:`/payroll`,bg:`bg-indigo-50`,color:`text-indigo-600`},{name:n(`dashboard.approvals`),icon:`pi pi-check-square`,route:`/approvals`,bg:`bg-violet-50`,color:`text-violet-600`},{name:n(`dashboard.performance`),icon:`pi pi-chart-line`,route:`/performance`,bg:`bg-cyan-50`,color:`text-cyan-600`},{name:n(`dashboard.training`),icon:`pi pi-book`,route:`/training`,bg:`bg-orange-50`,color:`text-orange-600`},{name:n(`dashboard.recruitment`),icon:`pi pi-user-plus`,route:`/recruitment`,bg:`bg-rose-50`,color:`text-rose-600`},{name:n(`dashboard.organization`),icon:`pi pi-sitemap`,route:`/organizations`,bg:`bg-teal-50`,color:`text-teal-600`},{name:n(`dashboard.reimbursement`),icon:`pi pi-credit-card`,route:`/reimbursements`,bg:`bg-sky-50`,color:`text-sky-600`},{name:n(`dashboard.workforce_intel`),icon:`pi pi-chart-bar`,route:`/workforce-intelligence`,bg:`bg-slate-50`,color:`text-slate-600`},{name:n(`dashboard.career_intel`),icon:`pi pi-road`,route:`/career-intelligence`,bg:`bg-pink-50`,color:`text-pink-600`}]),u=[{text:`15 new employees added this week`,time:`2 hours ago`,dotColor:`bg-emerald-400`},{text:`Payroll run for August completed`,time:`5 hours ago`,dotColor:`bg-blue-400`},{text:`3 leave requests pending approval`,time:`1 day ago`,dotColor:`bg-amber-400`},{text:`Performance reviews Q3 initiated`,time:`2 days ago`,dotColor:`bg-violet-400`},{text:`Training session "Leadership 101" scheduled`,time:`3 days ago`,dotColor:`bg-orange-400`}];return(t,r)=>(o(),v(`div`,Y,[_(`div`,X,[_(`div`,te,[d(c(K),{modelValue:i.value,"onUpdate:modelValue":r[0]||=e=>i.value=e,options:a.value,optionLabel:`label`,optionValue:`value`,size:`small`},null,8,[`modelValue`,`options`])])]),_(`div`,Z,[(o(!0),v(g,null,e(s.value,e=>(o(),v(`div`,{key:e.label,class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 hover:shadow-sm dark:hover:shadow-gray-900/50 transition-shadow`},[_(`div`,Q,[_(`span`,ne,h(e.label),1),_(`i`,{class:y([[e.icon,e.iconColor],`text-lg`])},null,2)]),_(`div`,re,h(e.value),1),_(`div`,ie,[_(`i`,{class:y([e.trend>=0?`pi pi-arrow-up text-emerald-500`:`pi pi-arrow-down text-rose-500`,`text-sm`])},null,2),_(`span`,{class:y([e.trend>=0?`text-emerald-600`:`text-rose-600`,`text-sm font-medium`])},h(Math.abs(e.trend))+`% `,3),_(`span`,ae,h(c(n)(`dashboard.vs_last_month`)),1)])]))),128))]),_(`div`,oe,[_(`div`,se,[_(`h2`,ce,h(c(n)(`dashboard.quick_access`)),1),_(`div`,le,[(o(!0),v(g,null,e(l.value,e=>(o(),v(`div`,{key:e.name,class:`flex flex-col items-center gap-1.5 p-2.5 rounded-lg cursor-pointer hover:bg-emerald-50 dark:hover:bg-emerald-900/20 hover:border-emerald-200 dark:hover:border-emerald-700 border border-transparent transition-all`,onClick:n=>t.$router.push(e.route)},[_(`div`,{class:y([e.bg,`w-9 h-9 rounded-lg flex items-center justify-center`])},[_(`i`,{class:y([[e.icon,e.color],`text-sm`])},null,2)],2),_(`span`,de,h(e.name),1)],8,ue))),128))])]),_(`div`,fe,[_(`h2`,pe,h(c(n)(`dashboard.recent_activity`)),1),_(`div`,$,[(o(),v(g,null,e(u,(e,t)=>_(`div`,{key:t,class:`flex items-start gap-2.5`},[_(`div`,{class:y([e.dotColor,`w-2 h-2 rounded-full mt-1.5 shrink-0`])},null,2),_(`div`,me,[_(`p`,he,h(e.text),1),_(`p`,ge,h(e.time),1)])])),64))])])])]))}};export{_e as default};