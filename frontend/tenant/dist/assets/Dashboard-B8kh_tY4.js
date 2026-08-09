import{C as e,J as t,M as n,N as r,P as i,Q as a,U as o,W as s,c,f as l,ft as u,h as d,ht as f,j as p,k as m,l as h,o as g,r as _,s as v,u as y}from"./runtime-core.esm-bundler-huI9Rd5Y.js";import{a as b,et as x,gt as S,pt as C,t as w,vt as T}from"./ripple-Cuh3jQqx.js";import{t as E}from"./useI18n-m28xTjhW.js";import{t as D}from"./baseeditableholder-6RjpXrFr.js";var O=b.extend({name:`togglebutton`,style:`
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
`,classes:{root:function(e){var t=e.instance,n=e.props;return[`p-togglebutton p-component`,{"p-togglebutton-checked":t.active,"p-invalid":t.$invalid,"p-togglebutton-fluid":n.fluid,"p-togglebutton-sm p-inputfield-sm":n.size===`small`,"p-togglebutton-lg p-inputfield-lg":n.size===`large`}]},content:`p-togglebutton-content`,icon:`p-togglebutton-icon`,label:`p-togglebutton-label`}}),k={name:`BaseToggleButton`,extends:D,props:{onIcon:String,offIcon:String,onLabel:{type:String,default:`Yes`},offLabel:{type:String,default:`No`},readonly:{type:Boolean,default:!1},tabindex:{type:Number,default:null},ariaLabelledby:{type:String,default:null},ariaLabel:{type:String,default:null},size:{type:String,default:null},fluid:{type:Boolean,default:null}},style:O,provide:function(){return{$pcToggleButton:this,$parentInstance:this}}};function A(e){"@babel/helpers - typeof";return A=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},A(e)}function j(e,t,n){return(t=ee(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function ee(e){var t=M(e,`string`);return A(t)==`symbol`?t:t+``}function M(e,t){if(A(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(A(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var N={name:`ToggleButton`,extends:k,inheritAttrs:!1,emits:[`change`],methods:{getPTOptions:function(e){return(e===`root`?this.ptmi:this.ptm)(e,{context:{active:this.active,disabled:this.disabled}})},onChange:function(e){!this.disabled&&!this.readonly&&(this.writeValue(!this.d_value,e),this.$emit(`change`,e))},onBlur:function(e){var t,n;(t=(n=this.formField).onBlur)==null||t.call(n,e)}},computed:{active:function(){return this.d_value===!0},hasLabel:function(){return T(this.onLabel)&&T(this.offLabel)},label:function(){return this.hasLabel?this.d_value?this.onLabel:this.offLabel:`\xA0`},dataP:function(){return x(j({checked:this.active,invalid:this.$invalid},this.size,this.size))}},directives:{ripple:w}},P=[`tabindex`,`disabled`,`aria-pressed`,`aria-label`,`aria-labelledby`,`data-p-checked`,`data-p-disabled`,`data-p`],F=[`data-p`];function I(t,r,a,o,c,l){var d=i(`ripple`);return s((m(),y(`button`,e({type:`button`,class:t.cx(`root`),tabindex:t.tabindex,disabled:t.disabled,"aria-pressed":t.d_value,onClick:r[0]||=function(){return l.onChange&&l.onChange.apply(l,arguments)},onBlur:r[1]||=function(){return l.onBlur&&l.onBlur.apply(l,arguments)}},l.getPTOptions(`root`),{"aria-label":t.ariaLabel,"aria-labelledby":t.ariaLabelledby,"data-p-checked":l.active,"data-p-disabled":t.disabled,"data-p":l.dataP}),[v(`span`,e({class:t.cx(`content`)},l.getPTOptions(`content`),{"data-p":l.dataP}),[n(t.$slots,`default`,{},function(){return[n(t.$slots,`icon`,{value:t.d_value,class:u(t.cx(`icon`))},function(){return[t.onIcon||t.offIcon?(m(),y(`span`,e({key:0,class:[t.cx(`icon`),t.d_value?t.onIcon:t.offIcon]},l.getPTOptions(`icon`)),null,16)):h(``,!0)]}),v(`span`,e({class:t.cx(`label`)},l.getPTOptions(`label`)),f(l.label),17)]})],16,F)],16,P)),[[d]])}N.render=I;var L=b.extend({name:`selectbutton`,style:`
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
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}function H(e,t){if(e){if(typeof e==`string`)return G(e,t);var n={}.toString.call(e).slice(8,-1);return n===`Object`&&e.constructor&&(n=e.constructor.name),n===`Map`||n===`Set`?Array.from(e):n===`Arguments`||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?G(e,t):void 0}}function U(e){if(typeof Symbol<`u`&&e[Symbol.iterator]!=null||e[`@@iterator`]!=null)return Array.from(e)}function W(e){if(Array.isArray(e))return G(e)}function G(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=Array(t);n<t;n++)r[n]=e[n];return r}var K={name:`SelectButton`,extends:R,inheritAttrs:!1,emits:[`change`],methods:{getOptionLabel:function(e){return this.optionLabel?S(e,this.optionLabel):e},getOptionValue:function(e){return this.optionValue?S(e,this.optionValue):e},getOptionRenderKey:function(e){return this.dataKey?S(e,this.dataKey):this.getOptionLabel(e)},isOptionDisabled:function(e){return this.optionDisabled?S(e,this.optionDisabled):!1},isOptionReadonly:function(e){if(this.allowEmpty)return!1;var t=this.isSelected(e);return this.multiple?t&&this.d_value.length===1:t},onOptionSelect:function(e,t,n){var r=this;if(!(this.disabled||this.isOptionDisabled(t)||this.isOptionReadonly(t))){var i=this.isSelected(t),a=this.getOptionValue(t),o;if(this.multiple)if(i){if(o=this.d_value.filter(function(e){return!C(e,a,r.equalityKey)}),!this.allowEmpty&&o.length===0)return}else o=this.d_value?[].concat(B(this.d_value),[a]):[a];else{if(i&&!this.allowEmpty)return;o=i?null:a}this.writeValue(o,e),this.$emit(`change`,{originalEvent:e,value:o})}},isSelected:function(e){var t=!1,n=this.getOptionValue(e);if(this.multiple){if(this.d_value){var r=z(this.d_value),i;try{for(r.s();!(i=r.n()).done;){var a=i.value;if(C(a,n,this.equalityKey)){t=!0;break}}}catch(e){r.e(e)}finally{r.f()}}}else t=C(this.d_value,n,this.equalityKey);return t}},computed:{equalityKey:function(){return this.optionValue?null:this.dataKey},dataP:function(){return x({invalid:this.$invalid})}},directives:{ripple:w},components:{ToggleButton:N}},q=[`aria-labelledby`,`data-p`];function te(t,i,a,s,u,d){var h=r(`ToggleButton`);return m(),y(`div`,e({class:t.cx(`root`),role:`group`,"aria-labelledby":t.ariaLabelledby},t.ptmi(`root`),{"data-p":d.dataP}),[(m(!0),y(_,null,p(t.options,function(r,i){return m(),c(h,{key:d.getOptionRenderKey(r),modelValue:d.isSelected(r),onLabel:d.getOptionLabel(r),offLabel:d.getOptionLabel(r),disabled:t.disabled||d.isOptionDisabled(r),unstyled:t.unstyled,size:t.size,readonly:d.isOptionReadonly(r),onChange:function(e){return d.onOptionSelect(e,r,i)},pt:t.ptm(`pcToggleButton`)},l({_:2},[t.$slots.option?{name:`default`,fn:o(function(){return[n(t.$slots,`option`,{option:r,index:i},function(){return[v(`span`,e({ref_for:!0},t.ptm(`pcToggleButton`).label),f(d.getOptionLabel(r)),17)]})]}),key:`0`}:void 0]),1032,[`modelValue`,`onLabel`,`offLabel`,`disabled`,`unstyled`,`size`,`readonly`,`onChange`,`pt`])}),128))],16,q)}K.render=te;var J={class:`space-y-4`},Y={class:`flex items-center justify-end`},X={class:`flex items-center gap-2`},Z={class:`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3`},Q={class:`flex items-center justify-between mb-2`},ne={class:`text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider`},re={class:`text-xl font-bold text-gray-800 dark:text-gray-100`},ie={class:`flex items-center gap-1 mt-1`},ae={class:`text-sm text-gray-400 dark:text-gray-500`},oe={class:`grid grid-cols-1 lg:grid-cols-3 gap-4`},se={class:`lg:col-span-2 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3`},ce={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3`},le={class:`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2`},ue=[`onClick`],de={class:`text-sm text-gray-600 dark:text-gray-300 text-center leading-tight`},fe={class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3`},pe={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3`},$={class:`space-y-3`},me={class:`min-w-0`},he={class:`text-sm text-gray-700 dark:text-gray-200`},ge={class:`text-[11px] text-gray-400 dark:text-gray-500 mt-0.5`},_e={__name:`Dashboard`,setup(e){let{t:n}=E(),r=t(`this-month`),i=g(()=>[{label:n(`dashboard.this_month`),value:`this-month`},{label:n(`dashboard.this_quarter`),value:`this-quarter`},{label:n(`dashboard.this_year`),value:`this-year`}]),o=g(()=>[{label:n(`dashboard.kpi_total_employees`),value:`1,247`,icon:`pi pi-users`,iconColor:`text-emerald-500`,trend:3.2},{label:n(`dashboard.kpi_active_today`),value:`1,183`,icon:`pi pi-check-circle`,iconColor:`text-blue-500`,trend:1.5},{label:n(`dashboard.kpi_on_leave`),value:`42`,icon:`pi pi-calendar`,iconColor:`text-amber-500`,trend:-2.1},{label:n(`dashboard.kpi_pending_approvals`),value:`28`,icon:`pi pi-clock`,iconColor:`text-rose-500`,trend:12.5}]),s=g(()=>[{name:n(`dashboard.employees`),icon:`pi pi-users`,route:`/employees`,bg:`bg-blue-50`,color:`text-blue-600`},{name:n(`dashboard.attendance`),icon:`pi pi-clock`,route:`/attendance`,bg:`bg-emerald-50`,color:`text-emerald-600`},{name:n(`dashboard.leave`),icon:`pi pi-calendar`,route:`/leave`,bg:`bg-amber-50`,color:`text-amber-600`},{name:n(`dashboard.payroll`),icon:`pi pi-dollar`,route:`/payroll`,bg:`bg-indigo-50`,color:`text-indigo-600`},{name:n(`dashboard.approvals`),icon:`pi pi-check-square`,route:`/approvals`,bg:`bg-violet-50`,color:`text-violet-600`},{name:n(`dashboard.performance`),icon:`pi pi-chart-line`,route:`/performance`,bg:`bg-cyan-50`,color:`text-cyan-600`},{name:n(`dashboard.training`),icon:`pi pi-book`,route:`/training`,bg:`bg-orange-50`,color:`text-orange-600`},{name:n(`dashboard.recruitment`),icon:`pi pi-user-plus`,route:`/recruitment`,bg:`bg-rose-50`,color:`text-rose-600`},{name:n(`dashboard.organization`),icon:`pi pi-sitemap`,route:`/organizations`,bg:`bg-teal-50`,color:`text-teal-600`},{name:n(`dashboard.reimbursement`),icon:`pi pi-credit-card`,route:`/reimbursements`,bg:`bg-sky-50`,color:`text-sky-600`},{name:n(`dashboard.workforce_intel`),icon:`pi pi-chart-bar`,route:`/workforce-intelligence`,bg:`bg-slate-50`,color:`text-slate-600`},{name:n(`dashboard.career_intel`),icon:`pi pi-road`,route:`/career-intelligence`,bg:`bg-pink-50`,color:`text-pink-600`}]),c=[{text:`15 new employees added this week`,time:`2 hours ago`,dotColor:`bg-emerald-400`},{text:`Payroll run for August completed`,time:`5 hours ago`,dotColor:`bg-blue-400`},{text:`3 leave requests pending approval`,time:`1 day ago`,dotColor:`bg-amber-400`},{text:`Performance reviews Q3 initiated`,time:`2 days ago`,dotColor:`bg-violet-400`},{text:`Training session "Leadership 101" scheduled`,time:`3 days ago`,dotColor:`bg-orange-400`}];return(e,t)=>(m(),y(`div`,J,[v(`div`,Y,[v(`div`,X,[d(a(K),{modelValue:r.value,"onUpdate:modelValue":t[0]||=e=>r.value=e,options:i.value,optionLabel:`label`,optionValue:`value`,size:`small`},null,8,[`modelValue`,`options`])])]),v(`div`,Z,[(m(!0),y(_,null,p(o.value,e=>(m(),y(`div`,{key:e.label,class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 hover:shadow-sm dark:hover:shadow-gray-900/50 transition-shadow`},[v(`div`,Q,[v(`span`,ne,f(e.label),1),v(`i`,{class:u([[e.icon,e.iconColor],`text-lg`])},null,2)]),v(`div`,re,f(e.value),1),v(`div`,ie,[v(`i`,{class:u([e.trend>=0?`pi pi-arrow-up text-emerald-500`:`pi pi-arrow-down text-rose-500`,`text-sm`])},null,2),v(`span`,{class:u([e.trend>=0?`text-emerald-600`:`text-rose-600`,`text-sm font-medium`])},f(Math.abs(e.trend))+`% `,3),v(`span`,ae,f(a(n)(`dashboard.vs_last_month`)),1)])]))),128))]),v(`div`,oe,[v(`div`,se,[v(`h2`,ce,f(a(n)(`dashboard.quick_access`)),1),v(`div`,le,[(m(!0),y(_,null,p(s.value,t=>(m(),y(`div`,{key:t.name,class:`flex flex-col items-center gap-1.5 p-2.5 rounded-lg cursor-pointer hover:bg-emerald-50 dark:hover:bg-emerald-900/20 hover:border-emerald-200 dark:hover:border-emerald-700 border border-transparent transition-all`,onClick:n=>e.$router.push(t.route)},[v(`div`,{class:u([t.bg,`w-9 h-9 rounded-lg flex items-center justify-center`])},[v(`i`,{class:u([[t.icon,t.color],`text-sm`])},null,2)],2),v(`span`,de,f(t.name),1)],8,ue))),128))])]),v(`div`,fe,[v(`h2`,pe,f(a(n)(`dashboard.recent_activity`)),1),v(`div`,$,[(m(),y(_,null,p(c,(e,t)=>v(`div`,{key:t,class:`flex items-start gap-2.5`},[v(`div`,{class:u([e.dotColor,`w-2 h-2 rounded-full mt-1.5 shrink-0`])},null,2),v(`div`,me,[v(`p`,he,f(e.text),1),v(`p`,ge,f(e.time),1)])])),64))])])])]))}};export{_e as default};