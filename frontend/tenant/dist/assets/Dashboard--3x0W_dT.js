import{$ as e,B as t,D as n,F as r,G as i,H as a,I as o,J as s,K as c,L as l,M as u,O as d,P as f,Q as p,R as m,U as h,V as g,W as _,X as v,Y as y,Z as ee,b,r as x,u as S,w as C}from"./index-BNiZ82cV.js";import{t as w}from"./baseeditableholder-lKcXaRhN.js";var T=S.extend({name:`togglebutton`,style:`
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
`,classes:{root:function(e){var t=e.instance,n=e.props;return[`p-togglebutton p-component`,{"p-togglebutton-checked":t.active,"p-invalid":t.$invalid,"p-togglebutton-fluid":n.fluid,"p-togglebutton-sm p-inputfield-sm":n.size===`small`,"p-togglebutton-lg p-inputfield-lg":n.size===`large`}]},content:`p-togglebutton-content`,icon:`p-togglebutton-icon`,label:`p-togglebutton-label`}}),te={name:`BaseToggleButton`,extends:w,props:{onIcon:String,offIcon:String,onLabel:{type:String,default:`Yes`},offLabel:{type:String,default:`No`},readonly:{type:Boolean,default:!1},tabindex:{type:Number,default:null},ariaLabelledby:{type:String,default:null},ariaLabel:{type:String,default:null},size:{type:String,default:null},fluid:{type:Boolean,default:null}},style:T,provide:function(){return{$pcToggleButton:this,$parentInstance:this}}};function E(e){"@babel/helpers - typeof";return E=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},E(e)}function D(e,t,n){return(t=O(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function O(e){var t=k(e,`string`);return E(t)==`symbol`?t:t+``}function k(e,t){if(E(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(E(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var A={name:`ToggleButton`,extends:te,inheritAttrs:!1,emits:[`change`],methods:{getPTOptions:function(e){return(e===`root`?this.ptmi:this.ptm)(e,{context:{active:this.active,disabled:this.disabled}})},onChange:function(e){!this.disabled&&!this.readonly&&(this.writeValue(!this.d_value,e),this.$emit(`change`,e))},onBlur:function(e){var t,n;(t=(n=this.formField).onBlur)==null||t.call(n,e)}},computed:{active:function(){return this.d_value===!0},hasLabel:function(){return d(this.onLabel)&&d(this.offLabel)},label:function(){return this.hasLabel?this.d_value?this.onLabel:this.offLabel:`\xA0`},dataP:function(){return b(D({checked:this.active,invalid:this.$invalid},this.size,this.size))}},directives:{ripple:x}},j=[`tabindex`,`disabled`,`aria-pressed`,`aria-label`,`aria-labelledby`,`data-p-checked`,`data-p-disabled`,`data-p`],M=[`data-p`];function N(t,n,r,i,s,u){var d=c(`ripple`);return y((a(),l(`button`,g({type:`button`,class:t.cx(`root`),tabindex:t.tabindex,disabled:t.disabled,"aria-pressed":t.d_value,onClick:n[0]||=function(){return u.onChange&&u.onChange.apply(u,arguments)},onBlur:n[1]||=function(){return u.onBlur&&u.onBlur.apply(u,arguments)}},u.getPTOptions(`root`),{"aria-label":t.ariaLabel,"aria-labelledby":t.ariaLabelledby,"data-p-checked":u.active,"data-p-disabled":t.disabled,"data-p":u.dataP}),[f(`span`,g({class:t.cx(`content`)},u.getPTOptions(`content`),{"data-p":u.dataP}),[_(t.$slots,`default`,{},function(){return[_(t.$slots,`icon`,{value:t.d_value,class:p(t.cx(`icon`))},function(){return[t.onIcon||t.offIcon?(a(),l(`span`,g({key:0,class:[t.cx(`icon`),t.d_value?t.onIcon:t.offIcon]},u.getPTOptions(`icon`)),null,16)):o(``,!0)]}),f(`span`,g({class:t.cx(`label`)},u.getPTOptions(`label`)),e(u.label),17)]})],16,M)],16,j)),[[d]])}A.render=N;var P=S.extend({name:`selectbutton`,style:`
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
`,classes:{root:function(e){var t=e.props;return[`p-selectbutton p-component`,{"p-invalid":e.instance.$invalid,"p-selectbutton-fluid":t.fluid}]}}}),F={name:`BaseSelectButton`,extends:w,props:{options:Array,optionLabel:null,optionValue:null,optionDisabled:null,multiple:Boolean,allowEmpty:{type:Boolean,default:!0},dataKey:null,ariaLabelledby:{type:String,default:null},size:{type:String,default:null},fluid:{type:Boolean,default:null}},style:P,provide:function(){return{$pcSelectButton:this,$parentInstance:this}}};function I(e,t){var n=typeof Symbol<`u`&&e[Symbol.iterator]||e[`@@iterator`];if(!n){if(Array.isArray(e)||(n=z(e))||t){n&&(e=n);var r=0,i=function(){};return{s:i,n:function(){return r>=e.length?{done:!0}:{done:!1,value:e[r++]}},e:function(e){throw e},f:i}}throw TypeError(`Invalid attempt to iterate non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}var a,o=!0,s=!1;return{s:function(){n=n.call(e)},n:function(){var e=n.next();return o=e.done,e},e:function(e){s=!0,a=e},f:function(){try{o||n.return==null||n.return()}finally{if(s)throw a}}}}function L(e){return V(e)||B(e)||z(e)||R()}function R(){throw TypeError(`Invalid attempt to spread non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}function z(e,t){if(e){if(typeof e==`string`)return H(e,t);var n={}.toString.call(e).slice(8,-1);return n===`Object`&&e.constructor&&(n=e.constructor.name),n===`Map`||n===`Set`?Array.from(e):n===`Arguments`||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?H(e,t):void 0}}function B(e){if(typeof Symbol<`u`&&e[Symbol.iterator]!=null||e[`@@iterator`]!=null)return Array.from(e)}function V(e){if(Array.isArray(e))return H(e)}function H(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=Array(t);n<t;n++)r[n]=e[n];return r}var U={name:`SelectButton`,extends:F,inheritAttrs:!1,emits:[`change`],methods:{getOptionLabel:function(e){return this.optionLabel?n(e,this.optionLabel):e},getOptionValue:function(e){return this.optionValue?n(e,this.optionValue):e},getOptionRenderKey:function(e){return this.dataKey?n(e,this.dataKey):this.getOptionLabel(e)},isOptionDisabled:function(e){return this.optionDisabled?n(e,this.optionDisabled):!1},isOptionReadonly:function(e){if(this.allowEmpty)return!1;var t=this.isSelected(e);return this.multiple?t&&this.d_value.length===1:t},onOptionSelect:function(e,t,n){var r=this;if(!(this.disabled||this.isOptionDisabled(t)||this.isOptionReadonly(t))){var i=this.isSelected(t),a=this.getOptionValue(t),o;if(this.multiple)if(i){if(o=this.d_value.filter(function(e){return!C(e,a,r.equalityKey)}),!this.allowEmpty&&o.length===0)return}else o=this.d_value?[].concat(L(this.d_value),[a]):[a];else{if(i&&!this.allowEmpty)return;o=i?null:a}this.writeValue(o,e),this.$emit(`change`,{originalEvent:e,value:o})}},isSelected:function(e){var t=!1,n=this.getOptionValue(e);if(this.multiple){if(this.d_value){var r=I(this.d_value),i;try{for(r.s();!(i=r.n()).done;){var a=i.value;if(C(a,n,this.equalityKey)){t=!0;break}}}catch(e){r.e(e)}finally{r.f()}}}else t=C(this.d_value,n,this.equalityKey);return t}},computed:{equalityKey:function(){return this.optionValue?null:this.dataKey},dataP:function(){return b({invalid:this.$invalid})}},directives:{ripple:x},components:{ToggleButton:A}},W=[`aria-labelledby`,`data-p`];function G(t,n,o,c,d,p){var v=i(`ToggleButton`);return a(),l(`div`,g({class:t.cx(`root`),role:`group`,"aria-labelledby":t.ariaLabelledby},t.ptmi(`root`),{"data-p":p.dataP}),[(a(!0),l(u,null,h(t.options,function(n,i){return a(),r(v,{key:p.getOptionRenderKey(n),modelValue:p.isSelected(n),onLabel:p.getOptionLabel(n),offLabel:p.getOptionLabel(n),disabled:t.disabled||p.isOptionDisabled(n),unstyled:t.unstyled,size:t.size,readonly:p.isOptionReadonly(n),onChange:function(e){return p.onOptionSelect(e,n,i)},pt:t.ptm(`pcToggleButton`)},m({_:2},[t.$slots.option?{name:`default`,fn:s(function(){return[_(t.$slots,`option`,{option:n,index:i},function(){return[f(`span`,g({ref_for:!0},t.ptm(`pcToggleButton`).label),e(p.getOptionLabel(n)),17)]})]}),key:`0`}:void 0]),1032,[`modelValue`,`onLabel`,`offLabel`,`disabled`,`unstyled`,`size`,`readonly`,`onChange`,`pt`])}),128))],16,W)}U.render=G;var K={class:`space-y-4`},q={class:`flex items-center justify-between`},J={class:`flex items-center gap-2`},Y={class:`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3`},X={class:`flex items-center justify-between mb-2`},Z={class:`text-sm font-medium text-gray-500 uppercase tracking-wider`},Q={class:`text-xl font-bold text-gray-800`},ne={class:`flex items-center gap-1 mt-1`},re={class:`grid grid-cols-1 lg:grid-cols-3 gap-4`},ie={class:`lg:col-span-2 bg-white rounded-lg border border-gray-200 p-3`},ae={class:`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2`},oe=[`onClick`],$={class:`text-sm text-gray-600 text-center leading-tight`},se={class:`bg-white rounded-lg border border-gray-200 p-3`},ce={class:`space-y-3`},le={class:`min-w-0`},ue={class:`text-sm text-gray-700`},de={class:`text-[11px] text-gray-400 mt-0.5`},fe={__name:`Dashboard`,setup(n){let r=v(`this-month`),i=[{label:`This Month`,value:`this-month`},{label:`This Quarter`,value:`this-quarter`},{label:`This Year`,value:`this-year`}],o=[{label:`Total Employees`,value:`1,247`,icon:`pi pi-users`,iconColor:`text-emerald-500`,trend:3.2},{label:`Active Today`,value:`1,183`,icon:`pi pi-check-circle`,iconColor:`text-blue-500`,trend:1.5},{label:`On Leave`,value:`42`,icon:`pi pi-calendar`,iconColor:`text-amber-500`,trend:-2.1},{label:`Pending Approvals`,value:`28`,icon:`pi pi-clock`,iconColor:`text-rose-500`,trend:12.5}],s=[{name:`Employees`,icon:`pi pi-users`,route:`/employees`,bg:`bg-blue-50`,color:`text-blue-600`},{name:`Attendance`,icon:`pi pi-clock`,route:`/attendance`,bg:`bg-emerald-50`,color:`text-emerald-600`},{name:`Leave`,icon:`pi pi-calendar`,route:`/leave`,bg:`bg-amber-50`,color:`text-amber-600`},{name:`Payroll`,icon:`pi pi-dollar`,route:`/payroll`,bg:`bg-indigo-50`,color:`text-indigo-600`},{name:`Approvals`,icon:`pi pi-check-square`,route:`/approvals`,bg:`bg-violet-50`,color:`text-violet-600`},{name:`Performance`,icon:`pi pi-chart-line`,route:`/performance`,bg:`bg-cyan-50`,color:`text-cyan-600`},{name:`Training`,icon:`pi pi-book`,route:`/training`,bg:`bg-orange-50`,color:`text-orange-600`},{name:`Recruitment`,icon:`pi pi-user-plus`,route:`/recruitment`,bg:`bg-rose-50`,color:`text-rose-600`},{name:`Organization`,icon:`pi pi-sitemap`,route:`/organizations`,bg:`bg-teal-50`,color:`text-teal-600`},{name:`Reimbursement`,icon:`pi pi-credit-card`,route:`/reimbursements`,bg:`bg-sky-50`,color:`text-sky-600`},{name:`Workforce Intel`,icon:`pi pi-chart-bar`,route:`/workforce-intelligence`,bg:`bg-slate-50`,color:`text-slate-600`},{name:`Career Intel`,icon:`pi pi-road`,route:`/career-intelligence`,bg:`bg-pink-50`,color:`text-pink-600`}],c=[{text:`15 new employees added this week`,time:`2 hours ago`,dotColor:`bg-emerald-400`},{text:`Payroll run for August completed`,time:`5 hours ago`,dotColor:`bg-blue-400`},{text:`3 leave requests pending approval`,time:`1 day ago`,dotColor:`bg-amber-400`},{text:`Performance reviews Q3 initiated`,time:`2 days ago`,dotColor:`bg-violet-400`},{text:`Training session "Leadership 101" scheduled`,time:`3 days ago`,dotColor:`bg-orange-400`}];return(n,d)=>(a(),l(`div`,K,[f(`div`,q,[d[1]||=f(`div`,null,[f(`h1`,{class:`text-lg font-semibold text-gray-800`},`Dashboard`),f(`p`,{class:`text-sm text-gray-500 mt-0.5`},`Selamat datang di HRIS Platform`)],-1),f(`div`,J,[t(ee(U),{modelValue:r.value,"onUpdate:modelValue":d[0]||=e=>r.value=e,options:i,optionLabel:`label`,optionValue:`value`,size:`small`},null,8,[`modelValue`])])]),f(`div`,Y,[(a(),l(u,null,h(o,t=>f(`div`,{key:t.label,class:`bg-white rounded-lg border border-gray-200 p-3 hover:shadow-sm transition-shadow`},[f(`div`,X,[f(`span`,Z,e(t.label),1),f(`i`,{class:p([[t.icon,t.iconColor],`text-lg`])},null,2)]),f(`div`,Q,e(t.value),1),f(`div`,ne,[f(`i`,{class:p([t.trend>=0?`pi pi-arrow-up text-emerald-500`:`pi pi-arrow-down text-rose-500`,`text-sm`])},null,2),f(`span`,{class:p([t.trend>=0?`text-emerald-600`:`text-rose-600`,`text-sm font-medium`])},e(Math.abs(t.trend))+`% `,3),d[2]||=f(`span`,{class:`text-sm text-gray-400`},`vs last month`,-1)])])),64))]),f(`div`,re,[f(`div`,ie,[d[3]||=f(`h2`,{class:`text-sm font-semibold text-gray-700 mb-3`},`Module Quick Access`,-1),f(`div`,ae,[(a(),l(u,null,h(s,t=>f(`div`,{key:t.name,class:`flex flex-col items-center gap-1.5 p-2.5 rounded-lg cursor-pointer hover:bg-emerald-50 hover:border-emerald-200 border border-transparent transition-all`,onClick:e=>n.$router.push(t.route)},[f(`div`,{class:p([t.bg,`w-9 h-9 rounded-lg flex items-center justify-center`])},[f(`i`,{class:p([[t.icon,t.color],`text-sm`])},null,2)],2),f(`span`,$,e(t.name),1)],8,oe)),64))])]),f(`div`,se,[d[4]||=f(`h2`,{class:`text-sm font-semibold text-gray-700 mb-3`},`Recent Activity`,-1),f(`div`,ce,[(a(),l(u,null,h(c,(t,n)=>f(`div`,{key:n,class:`flex items-start gap-2.5`},[f(`div`,{class:p([t.dotColor,`w-2 h-2 rounded-full mt-1.5 shrink-0`])},null,2),f(`div`,le,[f(`p`,ue,e(t.text),1),f(`p`,de,e(t.time),1)])])),64))])])])]))}};export{fe as default};