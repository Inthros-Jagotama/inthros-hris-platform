import{A as e,D as t,H as n,K as r,O as i,R as a,T as o,b as s,c,ct as l,dt as u,f as d,k as f,l as p,m,o as h,r as g,s as _,u as v,z as ee}from"./runtime-core.esm-bundler-CmryIYX-.js";import{gt as y,o as b,pt as x,t as S,tt as C,vt as w}from"./ripple-DtVIDCWU.js";import{t as T}from"./useI18n-B3Td2MI9.js";import{t as E}from"./baseeditableholder-DkY7pyuN.js";var D=b.extend({name:`togglebutton`,style:`
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
`,classes:{root:function(e){var t=e.instance,n=e.props;return[`p-togglebutton p-component`,{"p-togglebutton-checked":t.active,"p-invalid":t.$invalid,"p-togglebutton-fluid":n.fluid,"p-togglebutton-sm p-inputfield-sm":n.size===`small`,"p-togglebutton-lg p-inputfield-lg":n.size===`large`}]},content:`p-togglebutton-content`,icon:`p-togglebutton-icon`,label:`p-togglebutton-label`}}),O={name:`BaseToggleButton`,extends:E,props:{onIcon:String,offIcon:String,onLabel:{type:String,default:`Yes`},offLabel:{type:String,default:`No`},readonly:{type:Boolean,default:!1},tabindex:{type:Number,default:null},ariaLabelledby:{type:String,default:null},ariaLabel:{type:String,default:null},size:{type:String,default:null},fluid:{type:Boolean,default:null}},style:D,provide:function(){return{$pcToggleButton:this,$parentInstance:this}}};function k(e){"@babel/helpers - typeof";return k=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},k(e)}function A(e,t,n){return(t=j(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function j(e){var t=M(e,`string`);return k(t)==`symbol`?t:t+``}function M(e,t){if(k(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(k(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var N={name:`ToggleButton`,extends:O,inheritAttrs:!1,emits:[`change`],methods:{getPTOptions:function(e){return(e===`root`?this.ptmi:this.ptm)(e,{context:{active:this.active,disabled:this.disabled}})},onChange:function(e){!this.disabled&&!this.readonly&&(this.writeValue(!this.d_value,e),this.$emit(`change`,e))},onBlur:function(e){var t,n;(t=(n=this.formField).onBlur)==null||t.call(n,e)}},computed:{active:function(){return this.d_value===!0},hasLabel:function(){return w(this.onLabel)&&w(this.offLabel)},label:function(){return this.hasLabel?this.d_value?this.onLabel:this.offLabel:`\xA0`},dataP:function(){return C(A({checked:this.active,invalid:this.$invalid},this.size,this.size))}},directives:{ripple:S}},P=[`tabindex`,`disabled`,`aria-pressed`,`aria-label`,`aria-labelledby`,`data-p-checked`,`data-p-disabled`,`data-p`],F=[`data-p`];function I(t,n,r,a,c,d){var f=e(`ripple`);return ee((o(),v(`button`,s({type:`button`,class:t.cx(`root`),tabindex:t.tabindex,disabled:t.disabled,"aria-pressed":t.d_value,onClick:n[0]||=function(){return d.onChange&&d.onChange.apply(d,arguments)},onBlur:n[1]||=function(){return d.onBlur&&d.onBlur.apply(d,arguments)}},d.getPTOptions(`root`),{"aria-label":t.ariaLabel,"aria-labelledby":t.ariaLabelledby,"data-p-checked":d.active,"data-p-disabled":t.disabled,"data-p":d.dataP}),[_(`span`,s({class:t.cx(`content`)},d.getPTOptions(`content`),{"data-p":d.dataP}),[i(t.$slots,`default`,{},function(){return[i(t.$slots,`icon`,{value:t.d_value,class:l(t.cx(`icon`))},function(){return[t.onIcon||t.offIcon?(o(),v(`span`,s({key:0,class:[t.cx(`icon`),t.d_value?t.onIcon:t.offIcon]},d.getPTOptions(`icon`)),null,16)):p(``,!0)]}),_(`span`,s({class:t.cx(`label`)},d.getPTOptions(`label`)),u(d.label),17)]})],16,F)],16,P)),[[f]])}N.render=I;var L=b.extend({name:`selectbutton`,style:`
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
`,classes:{root:function(e){var t=e.props;return[`p-selectbutton p-component`,{"p-invalid":e.instance.$invalid,"p-selectbutton-fluid":t.fluid}]}}}),R={name:`BaseSelectButton`,extends:E,props:{options:Array,optionLabel:null,optionValue:null,optionDisabled:null,multiple:Boolean,allowEmpty:{type:Boolean,default:!0},dataKey:null,ariaLabelledby:{type:String,default:null},size:{type:String,default:null},fluid:{type:Boolean,default:null}},style:L,provide:function(){return{$pcSelectButton:this,$parentInstance:this}}};function z(e,t){var n=typeof Symbol<`u`&&e[Symbol.iterator]||e[`@@iterator`];if(!n){if(Array.isArray(e)||(n=H(e))||t){n&&(e=n);var r=0,i=function(){};return{s:i,n:function(){return r>=e.length?{done:!0}:{done:!1,value:e[r++]}},e:function(e){throw e},f:i}}throw TypeError(`Invalid attempt to iterate non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}var a,o=!0,s=!1;return{s:function(){n=n.call(e)},n:function(){var e=n.next();return o=e.done,e},e:function(e){s=!0,a=e},f:function(){try{o||n.return==null||n.return()}finally{if(s)throw a}}}}function B(e){return W(e)||U(e)||H(e)||V()}function V(){throw TypeError(`Invalid attempt to spread non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}function H(e,t){if(e){if(typeof e==`string`)return G(e,t);var n={}.toString.call(e).slice(8,-1);return n===`Object`&&e.constructor&&(n=e.constructor.name),n===`Map`||n===`Set`?Array.from(e):n===`Arguments`||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?G(e,t):void 0}}function U(e){if(typeof Symbol<`u`&&e[Symbol.iterator]!=null||e[`@@iterator`]!=null)return Array.from(e)}function W(e){if(Array.isArray(e))return G(e)}function G(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=Array(t);n<t;n++)r[n]=e[n];return r}var K={name:`SelectButton`,extends:R,inheritAttrs:!1,emits:[`change`],methods:{getOptionLabel:function(e){return this.optionLabel?y(e,this.optionLabel):e},getOptionValue:function(e){return this.optionValue?y(e,this.optionValue):e},getOptionRenderKey:function(e){return this.dataKey?y(e,this.dataKey):this.getOptionLabel(e)},isOptionDisabled:function(e){return this.optionDisabled?y(e,this.optionDisabled):!1},isOptionReadonly:function(e){if(this.allowEmpty)return!1;var t=this.isSelected(e);return this.multiple?t&&this.d_value.length===1:t},onOptionSelect:function(e,t,n){var r=this;if(!(this.disabled||this.isOptionDisabled(t)||this.isOptionReadonly(t))){var i=this.isSelected(t),a=this.getOptionValue(t),o;if(this.multiple)if(i){if(o=this.d_value.filter(function(e){return!x(e,a,r.equalityKey)}),!this.allowEmpty&&o.length===0)return}else o=this.d_value?[].concat(B(this.d_value),[a]):[a];else{if(i&&!this.allowEmpty)return;o=i?null:a}this.writeValue(o,e),this.$emit(`change`,{originalEvent:e,value:o})}},isSelected:function(e){var t=!1,n=this.getOptionValue(e);if(this.multiple){if(this.d_value){var r=z(this.d_value),i;try{for(r.s();!(i=r.n()).done;){var a=i.value;if(x(a,n,this.equalityKey)){t=!0;break}}}catch(e){r.e(e)}finally{r.f()}}}else t=x(this.d_value,n,this.equalityKey);return t}},computed:{equalityKey:function(){return this.optionValue?null:this.dataKey},dataP:function(){return C({invalid:this.$invalid})}},directives:{ripple:S},components:{ToggleButton:N}},q=[`aria-labelledby`,`data-p`];function J(e,n,r,l,p,m){var h=f(`ToggleButton`);return o(),v(`div`,s({class:e.cx(`root`),role:`group`,"aria-labelledby":e.ariaLabelledby},e.ptmi(`root`),{"data-p":m.dataP}),[(o(!0),v(g,null,t(e.options,function(t,n){return o(),c(h,{key:m.getOptionRenderKey(t),modelValue:m.isSelected(t),onLabel:m.getOptionLabel(t),offLabel:m.getOptionLabel(t),disabled:e.disabled||m.isOptionDisabled(t),unstyled:e.unstyled,size:e.size,readonly:m.isOptionReadonly(t),onChange:function(e){return m.onOptionSelect(e,t,n)},pt:e.ptm(`pcToggleButton`)},d({_:2},[e.$slots.option?{name:`default`,fn:a(function(){return[i(e.$slots,`option`,{option:t,index:n},function(){return[_(`span`,s({ref_for:!0},e.ptm(`pcToggleButton`).label),u(m.getOptionLabel(t)),17)]})]}),key:`0`}:void 0]),1032,[`modelValue`,`onLabel`,`offLabel`,`disabled`,`unstyled`,`size`,`readonly`,`onChange`,`pt`])}),128))],16,q)}K.render=J;var Y={class:`space-y-4`},te={class:`flex items-center justify-between`},X={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},Z={class:`text-sm text-gray-500 dark:text-gray-400 mt-0.5`},Q={class:`flex items-center gap-2`},ne={class:`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3`},re={class:`flex items-center justify-between mb-2`},ie={class:`text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider`},ae={class:`text-xl font-bold text-gray-800 dark:text-gray-100`},oe={class:`flex items-center gap-1 mt-1`},se={class:`text-sm text-gray-400 dark:text-gray-500`},ce={class:`grid grid-cols-1 lg:grid-cols-3 gap-4`},le={class:`lg:col-span-2 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3`},ue={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3`},de={class:`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2`},fe=[`onClick`],pe={class:`text-sm text-gray-600 dark:text-gray-300 text-center leading-tight`},$={class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3`},me={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3`},he={class:`space-y-3`},ge={class:`min-w-0`},_e={class:`text-sm text-gray-700 dark:text-gray-200`},ve={class:`text-[11px] text-gray-400 dark:text-gray-500 mt-0.5`},ye={__name:`Dashboard`,setup(e){let{t:i}=T(),a=n(`this-month`),s=h(()=>[{label:i(`dashboard.this_month`),value:`this-month`},{label:i(`dashboard.this_quarter`),value:`this-quarter`},{label:i(`dashboard.this_year`),value:`this-year`}]),c=h(()=>[{label:i(`dashboard.kpi_total_employees`),value:`1,247`,icon:`pi pi-users`,iconColor:`text-emerald-500`,trend:3.2},{label:i(`dashboard.kpi_active_today`),value:`1,183`,icon:`pi pi-check-circle`,iconColor:`text-blue-500`,trend:1.5},{label:i(`dashboard.kpi_on_leave`),value:`42`,icon:`pi pi-calendar`,iconColor:`text-amber-500`,trend:-2.1},{label:i(`dashboard.kpi_pending_approvals`),value:`28`,icon:`pi pi-clock`,iconColor:`text-rose-500`,trend:12.5}]),d=h(()=>[{name:i(`dashboard.employees`),icon:`pi pi-users`,route:`/employees`,bg:`bg-blue-50`,color:`text-blue-600`},{name:i(`dashboard.attendance`),icon:`pi pi-clock`,route:`/attendance`,bg:`bg-emerald-50`,color:`text-emerald-600`},{name:i(`dashboard.leave`),icon:`pi pi-calendar`,route:`/leave`,bg:`bg-amber-50`,color:`text-amber-600`},{name:i(`dashboard.payroll`),icon:`pi pi-dollar`,route:`/payroll`,bg:`bg-indigo-50`,color:`text-indigo-600`},{name:i(`dashboard.approvals`),icon:`pi pi-check-square`,route:`/approvals`,bg:`bg-violet-50`,color:`text-violet-600`},{name:i(`dashboard.performance`),icon:`pi pi-chart-line`,route:`/performance`,bg:`bg-cyan-50`,color:`text-cyan-600`},{name:i(`dashboard.training`),icon:`pi pi-book`,route:`/training`,bg:`bg-orange-50`,color:`text-orange-600`},{name:i(`dashboard.recruitment`),icon:`pi pi-user-plus`,route:`/recruitment`,bg:`bg-rose-50`,color:`text-rose-600`},{name:i(`dashboard.organization`),icon:`pi pi-sitemap`,route:`/organizations`,bg:`bg-teal-50`,color:`text-teal-600`},{name:i(`dashboard.reimbursement`),icon:`pi pi-credit-card`,route:`/reimbursements`,bg:`bg-sky-50`,color:`text-sky-600`},{name:i(`dashboard.workforce_intel`),icon:`pi pi-chart-bar`,route:`/workforce-intelligence`,bg:`bg-slate-50`,color:`text-slate-600`},{name:i(`dashboard.career_intel`),icon:`pi pi-road`,route:`/career-intelligence`,bg:`bg-pink-50`,color:`text-pink-600`}]),f=[{text:`15 new employees added this week`,time:`2 hours ago`,dotColor:`bg-emerald-400`},{text:`Payroll run for August completed`,time:`5 hours ago`,dotColor:`bg-blue-400`},{text:`3 leave requests pending approval`,time:`1 day ago`,dotColor:`bg-amber-400`},{text:`Performance reviews Q3 initiated`,time:`2 days ago`,dotColor:`bg-violet-400`},{text:`Training session "Leadership 101" scheduled`,time:`3 days ago`,dotColor:`bg-orange-400`}];return(e,n)=>(o(),v(`div`,Y,[_(`div`,te,[_(`div`,null,[_(`h1`,X,u(r(i)(`dashboard.title`)),1),_(`p`,Z,u(r(i)(`dashboard.welcome`)),1)]),_(`div`,Q,[m(r(K),{modelValue:a.value,"onUpdate:modelValue":n[0]||=e=>a.value=e,options:s.value,optionLabel:`label`,optionValue:`value`,size:`small`},null,8,[`modelValue`,`options`])])]),_(`div`,ne,[(o(!0),v(g,null,t(c.value,e=>(o(),v(`div`,{key:e.label,class:`bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 hover:shadow-sm dark:hover:shadow-gray-900/50 transition-shadow`},[_(`div`,re,[_(`span`,ie,u(e.label),1),_(`i`,{class:l([[e.icon,e.iconColor],`text-lg`])},null,2)]),_(`div`,ae,u(e.value),1),_(`div`,oe,[_(`i`,{class:l([e.trend>=0?`pi pi-arrow-up text-emerald-500`:`pi pi-arrow-down text-rose-500`,`text-sm`])},null,2),_(`span`,{class:l([e.trend>=0?`text-emerald-600`:`text-rose-600`,`text-sm font-medium`])},u(Math.abs(e.trend))+`% `,3),_(`span`,se,u(r(i)(`dashboard.vs_last_month`)),1)])]))),128))]),_(`div`,ce,[_(`div`,le,[_(`h2`,ue,u(r(i)(`dashboard.quick_access`)),1),_(`div`,de,[(o(!0),v(g,null,t(d.value,t=>(o(),v(`div`,{key:t.name,class:`flex flex-col items-center gap-1.5 p-2.5 rounded-lg cursor-pointer hover:bg-emerald-50 dark:hover:bg-emerald-900/20 hover:border-emerald-200 dark:hover:border-emerald-700 border border-transparent transition-all`,onClick:n=>e.$router.push(t.route)},[_(`div`,{class:l([t.bg,`w-9 h-9 rounded-lg flex items-center justify-center`])},[_(`i`,{class:l([[t.icon,t.color],`text-sm`])},null,2)],2),_(`span`,pe,u(t.name),1)],8,fe))),128))])]),_(`div`,$,[_(`h2`,me,u(r(i)(`dashboard.recent_activity`)),1),_(`div`,he,[(o(),v(g,null,t(f,(e,t)=>_(`div`,{key:t,class:`flex items-start gap-2.5`},[_(`div`,{class:l([e.dotColor,`w-2 h-2 rounded-full mt-1.5 shrink-0`])},null,2),_(`div`,ge,[_(`p`,_e,u(e.text),1),_(`p`,ve,u(e.time),1)])])),64))])])])]))}};export{ye as default};