import{n as e}from"./auth-CKDpQU78.js";import{A as t,D as n,E as r,G as i,L as a,O as o,R as s,S as c,V as l,b as u,c as d,ct as f,f as p,k as m,l as h,m as g,o as _,p as v,r as y,s as b,st as x,u as S,ut as C,w}from"./runtime-core.esm-bundler-DVMRdshy.js";import{$ as T,E,F as D,I as O,J as k,N as A,O as j,R as M,V as N,Z as P,_ as ee,a as F,dt as I,et as L,h as R,k as te,lt as z,mt as B,n as V,st as ne,t as H,w as re,x as ie}from"./ripple-iUk5gDdq.js";import{S as ae,_ as oe,a as se,b as ce,c as le,g as ue,h as U,i as W,l as G,m as de,n as fe,o as pe,t as me,w as he}from"./index-3U4bVB7C.js";import{t as ge}from"./useI18n-C-pD3zqu.js";import{n as _e,t as ve}from"./chevronright-BmsVdiZM.js";import{t as ye}from"./_plugin-vue_export-helper-BDNMzG2s.js";import{t as K}from"./inputtext-6fRrjfDU.js";import{n as be,t as xe}from"./tag-DrE29mxQ.js";import{a as Se,c as Ce,i as we,l as Te,n as Ee,o as De,r as Oe,s as ke,t as q,u as Ae}from"./column-Depf2U4z.js";import{t as je}from"./toggleswitch-BGwwhFlI.js";var Me={name:`ChevronLeftIcon`,extends:G};function Ne(e){return Le(e)||Ie(e)||Fe(e)||Pe()}function Pe(){throw TypeError(`Invalid attempt to spread non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}function Fe(e,t){if(e){if(typeof e==`string`)return Re(e,t);var n={}.toString.call(e).slice(8,-1);return n===`Object`&&e.constructor&&(n=e.constructor.name),n===`Map`||n===`Set`?Array.from(e):n===`Arguments`||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?Re(e,t):void 0}}function Ie(e){if(typeof Symbol<`u`&&e[Symbol.iterator]!=null||e[`@@iterator`]!=null)return Array.from(e)}function Le(e){if(Array.isArray(e))return Re(e)}function Re(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=Array(t);n<t;n++)r[n]=e[n];return r}function ze(e,t,n,r,i,a){return w(),S(`svg`,u({width:`14`,height:`14`,viewBox:`0 0 14 14`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},e.pti()),Ne(t[0]||=[b(`path`,{d:`M9.61296 13C9.50997 13.0005 9.40792 12.9804 9.3128 12.9409C9.21767 12.9014 9.13139 12.8433 9.05902 12.7701L3.83313 7.54416C3.68634 7.39718 3.60388 7.19795 3.60388 6.99022C3.60388 6.78249 3.68634 6.58325 3.83313 6.43628L9.05902 1.21039C9.20762 1.07192 9.40416 0.996539 9.60724 1.00012C9.81032 1.00371 10.0041 1.08597 10.1477 1.22959C10.2913 1.37322 10.3736 1.56698 10.3772 1.77005C10.3808 1.97313 10.3054 2.16968 10.1669 2.31827L5.49496 6.99022L10.1669 11.6622C10.3137 11.8091 10.3962 12.0084 10.3962 12.2161C10.3962 12.4238 10.3137 12.6231 10.1669 12.7701C10.0945 12.8433 10.0083 12.9014 9.91313 12.9409C9.81801 12.9804 9.71596 13.0005 9.61296 13Z`,fill:`currentColor`},null,-1)]),16)}Me.render=ze;var Be=F.extend({name:`tabview`,style:`
    .p-tabview-tablist-container {
        position: relative;
    }

    .p-tabview-scrollable > .p-tabview-tablist-container {
        overflow: hidden;
    }

    .p-tabview-tablist-scroll-container {
        overflow-x: auto;
        overflow-y: hidden;
        scroll-behavior: smooth;
        scrollbar-width: none;
        overscroll-behavior: contain auto;
    }

    .p-tabview-tablist-scroll-container::-webkit-scrollbar {
        display: none;
    }

    .p-tabview-tablist {
        display: flex;
        margin: 0;
        padding: 0;
        list-style-type: none;
        flex: 1 1 auto;
        background: dt('tabview.tab.list.background');
        border: 1px solid dt('tabview.tab.list.border.color');
        border-width: 0 0 1px 0;
        position: relative;
    }

    .p-tabview-tab-header {
        cursor: pointer;
        user-select: none;
        display: flex;
        align-items: center;
        text-decoration: none;
        position: relative;
        overflow: hidden;
        border-style: solid;
        border-width: 0 0 1px 0;
        border-color: transparent transparent dt('tabview.tab.border.color') transparent;
        color: dt('tabview.tab.color');
        padding: 1rem 1.125rem;
        font-weight: 600;
        border-top-right-radius: dt('border.radius.md');
        border-top-left-radius: dt('border.radius.md');
        transition:
            color dt('tabview.transition.duration'),
            outline-color dt('tabview.transition.duration');
        margin: 0 0 -1px 0;
        outline-color: transparent;
    }

    .p-tabview-tablist-item:not(.p-disabled) .p-tabview-tab-header:focus-visible {
        outline: dt('focus.ring.width') dt('focus.ring.style') dt('focus.ring.color');
        outline-offset: -1px;
    }

    .p-tabview-tablist-item:not(.p-highlight):not(.p-disabled):hover > .p-tabview-tab-header {
        color: dt('tabview.tab.hover.color');
    }

    .p-tabview-tablist-item.p-highlight > .p-tabview-tab-header {
        color: dt('tabview.tab.active.color');
    }

    .p-tabview-tab-title {
        line-height: 1;
        white-space: nowrap;
    }

    .p-tabview-next-button,
    .p-tabview-prev-button {
        position: absolute;
        top: 0;
        margin: 0;
        padding: 0;
        z-index: 2;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: dt('tabview.nav.button.background');
        color: dt('tabview.nav.button.color');
        width: 2.5rem;
        border-radius: 0;
        outline-color: transparent;
        transition:
            color dt('tabview.transition.duration'),
            outline-color dt('tabview.transition.duration');
        box-shadow: dt('tabview.nav.button.shadow');
        border: none;
        cursor: pointer;
        user-select: none;
    }

    .p-tabview-next-button:focus-visible,
    .p-tabview-prev-button:focus-visible {
        outline: dt('focus.ring.width') dt('focus.ring.style') dt('focus.ring.color');
        outline-offset: dt('focus.ring.offset');
    }

    .p-tabview-next-button:hover,
    .p-tabview-prev-button:hover {
        color: dt('tabview.nav.button.hover.color');
    }

    .p-tabview-prev-button {
        left: 0;
    }

    .p-tabview-next-button {
        right: 0;
    }

    .p-tabview-panels {
        background: dt('tabview.tab.panel.background');
        color: dt('tabview.tab.panel.color');
        padding: 0.875rem 1.125rem 1.125rem 1.125rem;
    }

    .p-tabview-ink-bar {
        z-index: 1;
        display: block;
        position: absolute;
        bottom: -1px;
        height: 1px;
        background: dt('tabview.tab.active.border.color');
        transition: 250ms cubic-bezier(0.35, 0, 0.25, 1);
    }
`,classes:{root:function(e){return[`p-tabview p-component`,{"p-tabview-scrollable":e.props.scrollable}]},navContainer:`p-tabview-tablist-container`,prevButton:`p-tabview-prev-button`,navContent:`p-tabview-tablist-scroll-container`,nav:`p-tabview-tablist`,tab:{header:function(e){var t=e.instance,n=e.tab,r=e.index;return[`p-tabview-tablist-item`,t.getTabProp(n,`headerClass`),{"p-tabview-tablist-item-active":t.d_activeIndex===r,"p-disabled":t.getTabProp(n,`disabled`)}]},headerAction:`p-tabview-tab-header`,headerTitle:`p-tabview-tab-title`,content:function(e){var t=e.instance,n=e.tab;return[`p-tabview-panel`,t.getTabProp(n,`contentClass`)]}},inkbar:`p-tabview-ink-bar`,nextButton:`p-tabview-next-button`,panelContainer:`p-tabview-panels`}}),Ve={name:`TabView`,extends:{name:`BaseTabView`,extends:V,props:{activeIndex:{type:Number,default:0},lazy:{type:Boolean,default:!1},scrollable:{type:Boolean,default:!1},tabindex:{type:Number,default:0},selectOnFocus:{type:Boolean,default:!1},prevButtonProps:{type:null,default:null},nextButtonProps:{type:null,default:null},prevIcon:{type:String,default:void 0},nextIcon:{type:String,default:void 0}},style:Be,provide:function(){return{$pcTabs:void 0,$pcTabView:this,$parentInstance:this}}},inheritAttrs:!1,emits:[`update:activeIndex`,`tab-change`,`tab-click`],data:function(){return{d_activeIndex:this.activeIndex,isPrevButtonDisabled:!0,isNextButtonDisabled:!1}},watch:{activeIndex:function(e){this.d_activeIndex=e,this.scrollInView({index:e})}},mounted:function(){console.warn(`Deprecated since v4. Use Tabs component instead.`),this.updateInkBar(),this.scrollable&&this.updateButtonState()},updated:function(){this.updateInkBar(),this.scrollable&&this.updateButtonState()},methods:{isTabPanel:function(e){return e.type.name===`TabPanel`},isTabActive:function(e){return this.d_activeIndex===e},getTabProp:function(e,t){return e.props?e.props[t]:void 0},getKey:function(e,t){return this.getTabProp(e,`header`)||t},getTabHeaderActionId:function(e){return`${this.$id}_${e}_header_action`},getTabContentId:function(e){return`${this.$id}_${e}_content`},getTabPT:function(e,t,n){var r=this.tabs.length,i={props:e.props,parent:{instance:this,props:this.$props,state:this.$data},context:{index:n,count:r,first:n===0,last:n===r-1,active:this.isTabActive(n)}};return u(this.ptm(`tabpanel.${t}`,{tabpanel:i}),this.ptm(`tabpanel.${t}`,i),this.ptmo(this.getTabProp(e,`pt`),t,i))},onScroll:function(e){this.scrollable&&this.updateButtonState(),e.preventDefault()},onPrevButtonClick:function(){var e=this.$refs.content,t=j(e),n=e.scrollLeft-t;e.scrollLeft=n<=0?0:n},onNextButtonClick:function(){var e=this.$refs.content,t=j(e)-this.getVisibleButtonWidths(),n=e.scrollLeft+t,r=e.scrollWidth-t;e.scrollLeft=n>=r?r:n},onTabClick:function(e,t,n){this.changeActiveIndex(e,t,n),this.$emit(`tab-click`,{originalEvent:e,index:n})},onTabKeyDown:function(e,t,n){switch(e.code){case`ArrowLeft`:this.onTabArrowLeftKey(e);break;case`ArrowRight`:this.onTabArrowRightKey(e);break;case`Home`:this.onTabHomeKey(e);break;case`End`:this.onTabEndKey(e);break;case`PageDown`:this.onPageDownKey(e);break;case`PageUp`:this.onPageUpKey(e);break;case`Enter`:case`NumpadEnter`:case`Space`:this.onTabEnterKey(e,t,n);break}},onTabArrowRightKey:function(e){var t=this.findNextHeaderAction(e.target.parentElement);t?this.changeFocusedTab(e,t):this.onTabHomeKey(e),e.preventDefault()},onTabArrowLeftKey:function(e){var t=this.findPrevHeaderAction(e.target.parentElement);t?this.changeFocusedTab(e,t):this.onTabEndKey(e),e.preventDefault()},onTabHomeKey:function(e){var t=this.findFirstHeaderAction();this.changeFocusedTab(e,t),e.preventDefault()},onTabEndKey:function(e){var t=this.findLastHeaderAction();this.changeFocusedTab(e,t),e.preventDefault()},onPageDownKey:function(e){this.scrollInView({index:this.$refs.nav.children.length-2}),e.preventDefault()},onPageUpKey:function(e){this.scrollInView({index:0}),e.preventDefault()},onTabEnterKey:function(e,t,n){this.changeActiveIndex(e,t,n),e.preventDefault()},findNextHeaderAction:function(e){var t=arguments.length>1&&arguments[1]!==void 0&&arguments[1]?e:e.nextElementSibling;return t?E(t,`data-p-disabled`)||E(t,`data-pc-section`)===`inkbar`?this.findNextHeaderAction(t):T(t,`[data-pc-section="headeraction"]`):null},findPrevHeaderAction:function(e){var t=arguments.length>1&&arguments[1]!==void 0&&arguments[1]?e:e.previousElementSibling;return t?E(t,`data-p-disabled`)||E(t,`data-pc-section`)===`inkbar`?this.findPrevHeaderAction(t):T(t,`[data-pc-section="headeraction"]`):null},findFirstHeaderAction:function(){return this.findNextHeaderAction(this.$refs.nav.firstElementChild,!0)},findLastHeaderAction:function(){return this.findPrevHeaderAction(this.$refs.nav.lastElementChild,!0)},changeActiveIndex:function(e,t,n){!this.getTabProp(t,`disabled`)&&this.d_activeIndex!==n&&(this.d_activeIndex=n,this.$emit(`update:activeIndex`,n),this.$emit(`tab-change`,{originalEvent:e,index:n}),this.scrollInView({index:n}))},changeFocusedTab:function(e,t){if(t&&(N(t),this.scrollInView({element:t}),this.selectOnFocus)){var n=parseInt(t.parentElement.dataset.pcIndex,10),r=this.tabs[n];this.changeActiveIndex(e,r,n)}},scrollInView:function(e){var t=e.element,n=e.index,r=n===void 0?-1:n,i=t||this.$refs.nav.children[r];i&&i.scrollIntoView&&i.scrollIntoView({block:`nearest`})},updateInkBar:function(){var e=this.$refs.nav.children[this.d_activeIndex];this.$refs.inkbar.style.width=j(e)+`px`,this.$refs.inkbar.style.left=ie(e).left-ie(this.$refs.nav).left+`px`},updateButtonState:function(){var e=this.$refs.content,t=e.scrollLeft,n=e.scrollWidth,r=j(e);this.isPrevButtonDisabled=t===0,this.isNextButtonDisabled=parseInt(t)===n-r},getVisibleButtonWidths:function(){var e=this.$refs;return[e.prevBtn,e.nextBtn].reduce(function(e,t){return t?e+j(t):e},0)}},computed:{tabs:function(){var e=this;return this.$slots.default().reduce(function(t,n){return e.isTabPanel(n)?t.push(n):n.children&&n.children instanceof Array&&n.children.forEach(function(n){e.isTabPanel(n)&&t.push(n)}),t},[])},prevButtonAriaLabel:function(){return this.$primevue.config.locale.aria?this.$primevue.config.locale.aria.previous:void 0},nextButtonAriaLabel:function(){return this.$primevue.config.locale.aria?this.$primevue.config.locale.aria.next:void 0}},directives:{ripple:H},components:{ChevronLeftIcon:Me,ChevronRightIcon:ve}};function J(e){"@babel/helpers - typeof";return J=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},J(e)}function He(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function Y(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?He(Object(n),!0).forEach(function(t){Ue(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):He(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function Ue(e,t,n){return(t=We(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function We(e){var t=Ge(e,`string`);return J(t)==`symbol`?t:t+``}function Ge(e,t){if(J(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(J(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var Ke=[`tabindex`,`aria-label`],qe=[`data-p-active`,`data-p-disabled`,`data-pc-index`],Je=[`id`,`tabindex`,`aria-disabled`,`aria-selected`,`aria-controls`,`onClick`,`onKeydown`],Ye=[`tabindex`,`aria-label`],Xe=[`id`,`aria-labelledby`,`data-pc-index`,`data-p-active`];function Ze(e,i,a,o,c,l){var f=m(`ripple`);return w(),S(`div`,u({class:e.cx(`root`),role:`tablist`},e.ptmi(`root`)),[b(`div`,u({class:e.cx(`navContainer`)},e.ptm(`navContainer`)),[e.scrollable&&!c.isPrevButtonDisabled?s((w(),S(`button`,u({key:0,ref:`prevBtn`,type:`button`,class:e.cx(`prevButton`),tabindex:e.tabindex,"aria-label":l.prevButtonAriaLabel,onClick:i[0]||=function(){return l.onPrevButtonClick&&l.onPrevButtonClick.apply(l,arguments)}},Y(Y({},e.prevButtonProps),e.ptm(`prevButton`)),{"data-pc-group-section":`navbutton`}),[n(e.$slots,`previcon`,{},function(){return[(w(),d(t(e.prevIcon?`span`:`ChevronLeftIcon`),u({"aria-hidden":`true`,class:e.prevIcon},e.ptm(`prevIcon`)),null,16,[`class`]))]})],16,Ke)),[[f]]):h(``,!0),b(`div`,u({ref:`content`,class:e.cx(`navContent`),onScroll:i[1]||=function(){return l.onScroll&&l.onScroll.apply(l,arguments)}},e.ptm(`navContent`)),[b(`ul`,u({ref:`nav`,class:e.cx(`nav`)},e.ptm(`nav`)),[(w(!0),S(y,null,r(l.tabs,function(n,r){return w(),S(`li`,u({key:l.getKey(n,r),style:l.getTabProp(n,`headerStyle`),class:e.cx(`tab.header`,{tab:n,index:r}),role:`presentation`},{ref_for:!0},Y(Y(Y({},l.getTabProp(n,`headerProps`)),l.getTabPT(n,`root`,r)),l.getTabPT(n,`header`,r)),{"data-pc-name":`tabpanel`,"data-p-active":c.d_activeIndex===r,"data-p-disabled":l.getTabProp(n,`disabled`),"data-pc-index":r}),[s((w(),S(`a`,u({id:l.getTabHeaderActionId(r),class:e.cx(`tab.headerAction`),tabindex:l.getTabProp(n,`disabled`)||!l.isTabActive(r)?-1:e.tabindex,role:`tab`,"aria-disabled":l.getTabProp(n,`disabled`),"aria-selected":l.isTabActive(r),"aria-controls":l.getTabContentId(r),onClick:function(e){return l.onTabClick(e,n,r)},onKeydown:function(e){return l.onTabKeyDown(e,n,r)}},{ref_for:!0},Y(Y({},l.getTabProp(n,`headerActionProps`)),l.getTabPT(n,`headerAction`,r))),[n.props&&n.props.header?(w(),S(`span`,u({key:0,class:e.cx(`tab.headerTitle`)},{ref_for:!0},l.getTabPT(n,`headerTitle`,r)),C(n.props.header),17)):h(``,!0),n.children&&n.children.header?(w(),d(t(n.children.header),{key:1})):h(``,!0)],16,Je)),[[f]])],16,qe)}),128)),b(`li`,u({ref:`inkbar`,class:e.cx(`inkbar`),role:`presentation`,"aria-hidden":`true`},e.ptm(`inkbar`)),null,16)],16)],16),e.scrollable&&!c.isNextButtonDisabled?s((w(),S(`button`,u({key:1,ref:`nextBtn`,type:`button`,class:e.cx(`nextButton`),tabindex:e.tabindex,"aria-label":l.nextButtonAriaLabel,onClick:i[2]||=function(){return l.onNextButtonClick&&l.onNextButtonClick.apply(l,arguments)}},Y(Y({},e.nextButtonProps),e.ptm(`nextButton`)),{"data-pc-group-section":`navbutton`}),[n(e.$slots,`nexticon`,{},function(){return[(w(),d(t(e.nextIcon?`span`:`ChevronRightIcon`),u({"aria-hidden":`true`,class:e.nextIcon},e.ptm(`nextIcon`)),null,16,[`class`]))]})],16,Ye)),[[f]]):h(``,!0)],16),b(`div`,u({class:e.cx(`panelContainer`)},e.ptm(`panelContainer`)),[(w(!0),S(y,null,r(l.tabs,function(n,r){return w(),S(y,{key:l.getKey(n,r)},[!e.lazy||l.isTabActive(r)?s((w(),S(`div`,u({key:0,id:l.getTabContentId(r),style:l.getTabProp(n,`contentStyle`),class:e.cx(`tab.content`,{tab:n}),role:`tabpanel`,"aria-labelledby":l.getTabHeaderActionId(r)},{ref_for:!0},Y(Y(Y({},l.getTabProp(n,`contentProps`)),l.getTabPT(n,`root`,r)),l.getTabPT(n,`content`,r)),{"data-pc-name":`tabpanel`,"data-pc-index":r,"data-p-active":c.d_activeIndex===r}),[(w(),d(t(n)))],16,Xe)),[[he,e.lazy?!0:l.isTabActive(r)]]):h(``,!0)],64)}),128))],16)],16)}Ve.render=Ze;var Qe=F.extend({name:`tabpanel`,classes:{root:function(e){return[`p-tabpanel`,{"p-tabpanel-active":e.instance.active}]}}}),$e={name:`TabPanel`,extends:{name:`BaseTabPanel`,extends:V,props:{value:{type:[String,Number],default:void 0},as:{type:[String,Object],default:`DIV`},asChild:{type:Boolean,default:!1},header:null,headerStyle:null,headerClass:null,headerProps:null,headerActionProps:null,contentStyle:null,contentClass:null,contentProps:null,disabled:Boolean},style:Qe,provide:function(){return{$pcTabPanel:this,$parentInstance:this}}},inheritAttrs:!1,inject:[`$pcTabs`],computed:{active:function(){return I(this.$pcTabs?.d_value,this.value)},id:function(){return`${this.$pcTabs?.$id}_tabpanel_${this.value}`},ariaLabelledby:function(){return`${this.$pcTabs?.$id}_tab_${this.value}`},attrs:function(){return u(this.a11yAttrs,this.ptmi(`root`,this.ptParams))},a11yAttrs:function(){return{id:this.id,tabindex:this.$pcTabs?.tabindex,role:`tabpanel`,"aria-labelledby":this.ariaLabelledby,"data-pc-name":`tabpanel`,"data-p-active":this.active}},ptParams:function(){return{context:{active:this.active}}}}};function et(e,r,i,o,c,l){var f,p;return l.$pcTabs?(w(),S(y,{key:1},[e.asChild?n(e.$slots,`default`,{key:1,class:x(e.cx(`root`)),active:l.active,a11yAttrs:l.a11yAttrs}):(w(),S(y,{key:0},[!((f=l.$pcTabs)!=null&&f.lazy)||l.active?s((w(),d(t(e.as),u({key:0,class:e.cx(`root`)},l.attrs),{default:a(function(){return[n(e.$slots,`default`)]}),_:3},16,[`class`])),[[he,(p=l.$pcTabs)!=null&&p.lazy?!0:l.active]]):h(``,!0)],64))],64)):n(e.$slots,`default`,{key:0})}$e.render=et;var tt=F.extend({name:`treetable`,style:`
    .p-treetable {
        position: relative;
    }

    .p-treetable-table {
        border-spacing: 0;
        border-collapse: separate;
        width: 100%;
    }

    .p-treetable-scrollable > .p-treetable-table-container {
        position: relative;
    }

    .p-treetable-scrollable-table > .p-treetable-thead {
        inset-block-start: 0;
        z-index: 1;
    }

    .p-treetable-scrollable-table > .p-treetable-frozen-tbody {
        position: sticky;
        z-index: 1;
    }

    .p-treetable-scrollable-table > .p-treetable-tfoot {
        inset-block-end: 0;
        z-index: 1;
    }

    .p-treetable-scrollable .p-treetable-frozen-column {
        position: sticky;
        background: dt('treetable.header.cell.background');
    }

    .p-treetable-scrollable th.p-treetable-frozen-column {
        z-index: 1;
    }

    .p-treetable-scrollable > .p-treetable-table-container > .p-treetable-table > .p-treetable-thead {
        background: dt('treetable.header.cell.background');
    }

    .p-treetable-scrollable > .p-treetable-table-container > .p-treetable-table > .p-treetable-tfoot {
        background: dt('treetable.footer.cell.background');
    }

    .p-treetable-flex-scrollable {
        display: flex;
        flex-direction: column;
        height: 100%;
    }

    .p-treetable-flex-scrollable > .p-treetable-table-container {
        display: flex;
        flex-direction: column;
        flex: 1;
        height: 100%;
    }

    .p-treetable-scrollable-table > .p-treetable-tbody > .p-treetable-row-group-header {
        position: sticky;
        z-index: 1;
    }

    .p-treetable-resizable-table > .p-treetable-thead > tr > th,
    .p-treetable-resizable-table > .p-treetable-tfoot > tr > td,
    .p-treetable-resizable-table > .p-treetable-tbody > tr > td {
        overflow: hidden;
        white-space: nowrap;
    }

    .p-treetable-resizable-table > .p-treetable-thead > tr > th.p-treetable-resizable-column:not(.p-treetable-frozen-column) {
        background-clip: padding-box;
        position: relative;
    }

    .p-treetable-resizable-table-fit > .p-treetable-thead > tr > th.p-treetable-resizable-column:last-child .p-treetable-column-resizer {
        display: none;
    }

    .p-treetable-column-resizer {
        display: block;
        position: absolute;
        inset-block-start: 0;
        inset-inline-end: 0;
        margin: 0;
        width: dt('treetable.column.resizer.width');
        height: 100%;
        padding: 0;
        cursor: col-resize;
        border: 1px solid transparent;
    }

    .p-treetable-column-header-content {
        display: flex;
        align-items: center;
        gap: dt('treetable.header.cell.gap');
    }

    .p-treetable-column-resize-indicator {
        width: dt('treetable.resize.indicator.width');
        position: absolute;
        z-index: 10;
        display: none;
        background: dt('treetable.resize.indicator.color');
    }

    .p-treetable-mask {
        position: absolute;
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 2;
    }

    .p-treetable-paginator-top {
        border-color: dt('treetable.paginator.top.border.color');
        border-style: solid;
        border-width: dt('treetable.paginator.top.border.width');
    }

    .p-treetable-paginator-bottom {
        border-color: dt('treetable.paginator.bottom.border.color');
        border-style: solid;
        border-width: dt('treetable.paginator.bottom.border.width');
    }

    .p-treetable-header {
        background: dt('treetable.header.background');
        color: dt('treetable.header.color');
        border-color: dt('treetable.header.border.color');
        border-style: solid;
        border-width: dt('treetable.header.border.width');
        padding: dt('treetable.header.padding');
    }

    .p-treetable-footer {
        background: dt('treetable.footer.background');
        color: dt('treetable.footer.color');
        border-color: dt('treetable.footer.border.color');
        border-style: solid;
        border-width: dt('treetable.footer.border.width');
        padding: dt('treetable.footer.padding');
    }

    .p-treetable-header-cell {
        padding: dt('treetable.header.cell.padding');
        background: dt('treetable.header.cell.background');
        border-color: dt('treetable.header.cell.border.color');
        border-style: solid;
        border-width: 0 0 1px 0;
        color: dt('treetable.header.cell.color');
        font-weight: normal;
        text-align: start;
        transition:
            background dt('treetable.transition.duration'),
            color dt('treetable.transition.duration'),
            border-color dt('treetable.transition.duration'),
            outline-color dt('treetable.transition.duration'),
            box-shadow dt('treetable.transition.duration');
    }

    .p-treetable-column-title {
        font-weight: dt('treetable.column.title.font.weight');
    }

    .p-treetable-tbody > tr {
        outline-color: transparent;
        background: dt('treetable.row.background');
        color: dt('treetable.row.color');
        transition:
            background dt('treetable.transition.duration'),
            color dt('treetable.transition.duration'),
            border-color dt('treetable.transition.duration'),
            outline-color dt('treetable.transition.duration'),
            box-shadow dt('treetable.transition.duration');
    }

    .p-treetable-tbody > tr > td {
        text-align: start;
        border-color: dt('treetable.body.cell.border.color');
        border-style: solid;
        border-width: 0 0 1px 0;
        padding: dt('treetable.body.cell.padding');
    }

    .p-treetable-hoverable .p-treetable-tbody > tr:not(.p-treetable-row-selected):hover {
        background: dt('treetable.row.hover.background');
        color: dt('treetable.row.hover.color');
    }

    .p-treetable-tbody > tr.p-treetable-row-selected {
        background: dt('treetable.row.selected.background');
        color: dt('treetable.row.selected.color');
    }

    .p-treetable-tbody > tr:has(+ .p-treetable-row-selected) > td {
        border-block-end-color: dt('treetable.body.cell.selected.border.color');
    }

    .p-treetable-tbody > tr.p-treetable-row-selected > td {
        border-block-end-color: dt('treetable.body.cell.selected.border.color');
    }

    .p-treetable-tbody > tr:focus-visible,
    .p-treetable-tbody > tr.p-treetable-contextmenu-row-selected {
        box-shadow: dt('treetable.row.focus.ring.shadow');
        outline: dt('treetable.row.focus.ring.width') dt('treetable.row.focus.ring.style') dt('treetable.row.focus.ring.color');
        outline-offset: dt('treetable.row.focus.ring.offset');
    }

    .p-treetable-tfoot > tr > td {
        text-align: start;
        padding: dt('treetable.footer.cell.padding');
        border-color: dt('treetable.footer.cell.border.color');
        border-style: solid;
        border-width: 0 0 1px 0;
        color: dt('treetable.footer.cell.color');
        background: dt('treetable.footer.cell.background');
    }

    .p-treetable-column-footer {
        font-weight: dt('treetable.column.footer.font.weight');
    }

    .p-treetable-sortable-column {
        cursor: pointer;
        user-select: none;
        outline-color: transparent;
    }

    .p-treetable-column-title,
    .p-treetable-sort-icon,
    .p-treetable-sort-badge {
        vertical-align: middle;
    }

    .p-treetable-sort-icon {
        color: dt('treetable.sort.icon.color');
        font-size: dt('treetable.sort.icon.size');
        width: dt('treetable.sort.icon.size');
        height: dt('treetable.sort.icon.size');
        transition: color dt('treetable.transition.duration');
    }

    .p-treetable-sortable-column:not(.p-treetable-column-sorted):hover {
        background: dt('treetable.header.cell.hover.background');
        color: dt('treetable.header.cell.hover.color');
    }

    .p-treetable-sortable-column:not(.p-treetable-column-sorted):hover .p-treetable-sort-icon {
        color: dt('treetable.sort.icon.hover.color');
    }

    .p-treetable-column-sorted {
        background: dt('treetable.header.cell.selected.background');
        color: dt('treetable.header.cell.selected.color');
    }

    .p-treetable-column-sorted .p-treetable-sort-icon {
        color: dt('treetable.header.cell.selected.color');
    }

    .p-treetable-sortable-column:focus-visible {
        box-shadow: dt('treetable.header.cell.focus.ring.shadow');
        outline: dt('treetable.header.cell.focus.ring.width') dt('treetable.header.cell.focus.ring.style') dt('treetable.header.cell.focus.ring.color');
        outline-offset: dt('treetable.header.cell.focus.ring.offset');
    }

    .p-treetable-hoverable .p-treetable-selectable-row {
        cursor: pointer;
    }

    .p-treetable-loading-icon {
        font-size: dt('treetable.loading.icon.size');
        width: dt('treetable.loading.icon.size');
        height: dt('treetable.loading.icon.size');
    }

    .p-treetable-gridlines .p-treetable-header {
        border-width: 1px 1px 0 1px;
    }

    .p-treetable-gridlines .p-treetable-footer {
        border-width: 0 1px 1px 1px;
    }

    .p-treetable-gridlines .p-treetable-paginator-top {
        border-width: 1px 1px 0 1px;
    }

    .p-treetable-gridlines .p-treetable-paginator-bottom {
        border-width: 0 1px 1px 1px;
    }

    .p-treetable-gridlines .p-treetable-thead > tr > th {
        border-width: 1px 0 1px 1px;
    }

    .p-treetable-gridlines .p-treetable-thead > tr > th:last-child {
        border-width: 1px;
    }

    .p-treetable-gridlines .p-treetable-tbody > tr > td {
        border-width: 1px 0 0 1px;
    }

    .p-treetable-gridlines .p-treetable-tbody > tr > td:last-child {
        border-width: 1px 1px 0 1px;
    }

    .p-treetable-gridlines .p-treetable-tbody > tr:last-child > td {
        border-width: 1px 0 1px 1px;
    }

    .p-treetable-gridlines .p-treetable-tbody > tr:last-child > td:last-child {
        border-width: 1px;
    }

    .p-treetable-gridlines .p-treetable-tfoot > tr > td {
        border-width: 1px 0 1px 1px;
    }

    .p-treetable-gridlines .p-treetable-tfoot > tr > td:last-child {
        border-width: 1px 1px 1px 1px;
    }

    .p-treetable.p-treetable-gridlines .p-treetable-thead + .p-treetable-tfoot > tr > td {
        border-width: 0 0 1px 1px;
    }

    .p-treetable.p-treetable-gridlines .p-treetable-thead + .p-treetable-tfoot > tr > td:last-child {
        border-width: 0 1px 1px 1px;
    }

    .p-treetable.p-treetable-gridlines:has(.p-treetable-thead):has(.p-treetable-tbody) .p-treetable-tbody > tr > td {
        border-width: 0 0 1px 1px;
    }

    .p-treetable.p-treetable-gridlines:has(.p-treetable-thead):has(.p-treetable-tbody) .p-treetable-tbody > tr > td:last-child {
        border-width: 0 1px 1px 1px;
    }

    .p-treetable.p-treetable-gridlines:has(.p-treetable-tbody):has(.p-treetable-tfoot) .p-treetable-tbody > tr:last-child > td {
        border-width: 0 0 0 1px;
    }

    .p-treetable.p-treetable-gridlines:has(.p-treetable-tbody):has(.p-treetable-tfoot) .p-treetable-tbody > tr:last-child > td:last-child {
        border-width: 0 1px 0 1px;
    }

    .p-treetable.p-treetable-sm .p-treetable-header {
        padding: 0.375rem 0.5rem;
    }

    .p-treetable.p-treetable-sm .p-treetable-thead > tr > th {
        padding: 0.375rem 0.5rem;
    }

    .p-treetable.p-treetable-sm .p-treetable-tbody > tr > td {
        padding: 0.375rem 0.5rem;
    }

    .p-treetable.p-treetable-sm .p-treetable-tfoot > tr > td {
        padding: 0.375rem 0.5rem;
    }

    .p-treetable.p-treetable-sm .p-treetable-footer {
        padding: 0.375rem 0.5rem;
    }

    .p-treetable.p-treetable-lg .p-treetable-header {
        padding: 0.9375rem 1.25rem;
    }

    .p-treetable.p-treetable-lg .p-treetable-thead > tr > th {
        padding: 0.9375rem 1.25rem;
    }

    .p-treetable.p-treetable-lg .p-treetable-tbody > tr > td {
        padding: 0.9375rem 1.25rem;
    }

    .p-treetable.p-treetable-lg .p-treetable-tfoot > tr > td {
        padding: 0.9375rem 1.25rem;
    }

    .p-treetable.p-treetable-lg .p-treetable-footer {
        padding: 0.9375rem 1.25rem;
    }

    .p-treetable-body-cell-content {
        display: flex;
        align-items: center;
        gap: dt('treetable.body.cell.gap');
    }

    .p-treetable-tbody > tr.p-treetable-row-selected .p-treetable-node-toggle-button {
        color: inherit;
    }

    .p-treetable-node-toggle-button {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        overflow: hidden;
        position: relative;
        width: dt('treetable.node.toggle.button.size');
        height: dt('treetable.node.toggle.button.size');
        color: dt('treetable.node.toggle.button.color');
        border: 0 none;
        background: transparent;
        cursor: pointer;
        border-radius: dt('treetable.node.toggle.button.border.radius');
        transition:
            background dt('treetable.transition.duration'),
            color dt('treetable.transition.duration'),
            border-color dt('treetable.transition.duration'),
            outline-color dt('treetable.transition.duration'),
            box-shadow dt('treetable.transition.duration');
        outline-color: transparent;
        user-select: none;
    }

    .p-treetable-node-toggle-button:enabled:hover {
        color: dt('treetable.node.toggle.button.hover.color');
        background: dt('treetable.node.toggle.button.hover.background');
    }

    .p-treetable-tbody > tr.p-treetable-row-selected .p-treetable-node-toggle-button:hover {
        background: dt('treetable.node.toggle.button.selected.hover.background');
        color: dt('treetable.node.toggle.button.selected.hover.color');
    }

    .p-treetable-node-toggle-button:focus-visible {
        box-shadow: dt('treetable.node.toggle.button.focus.ring.shadow');
        outline: dt('treetable.node.toggle.button.focus.ring.width') dt('treetable.node.toggle.button.focus.ring.style') dt('treetable.node.toggle.button.focus.ring.color');
        outline-offset: dt('treetable.node.toggle.button.focus.ring.offset');
    }

    .p-treetable-node-toggle-icon:dir(rtl) {
        transform: rotate(180deg);
    }
`,classes:{root:function(e){var t=e.instance,n=e.props;return[`p-treetable p-component`,{"p-treetable-hoverable":n.rowHover||t.rowSelectionMode,"p-treetable-resizable":n.resizableColumns,"p-treetable-resizable-fit":n.resizableColumns&&n.columnResizeMode===`fit`,"p-treetable-scrollable":n.scrollable,"p-treetable-flex-scrollable":n.scrollable&&n.scrollHeight===`flex`,"p-treetable-gridlines":n.showGridlines,"p-treetable-sm":n.size===`small`,"p-treetable-lg":n.size===`large`}]},loading:`p-treetable-loading`,mask:`p-treetable-mask p-overlay-mask`,loadingIcon:`p-treetable-loading-icon`,header:`p-treetable-header`,paginator:function(e){return`p-treetable-paginator-`+e.position},tableContainer:`p-treetable-table-container`,table:function(e){var t=e.props;return[`p-treetable-table`,{"p-treetable-scrollable-table":t.scrollable,"p-treetable-resizable-table":t.resizableColumns,"p-treetable-resizable-table-fit":t.resizableColumns&&t.columnResizeMode===`fit`}]},thead:`p-treetable-thead`,headerCell:function(e){var t=e.instance,n=e.props;return[`p-treetable-header-cell`,{"p-treetable-sortable-column":t.columnProp(`sortable`),"p-treetable-resizable-column":n.resizableColumns,"p-treetable-column-sorted":t.columnProp(`sortable`)?t.isColumnSorted():!1,"p-treetable-frozen-column":t.columnProp(`frozen`)}]},columnResizer:`p-treetable-column-resizer`,columnHeaderContent:`p-treetable-column-header-content`,columnTitle:`p-treetable-column-title`,sortIcon:`p-treetable-sort-icon`,pcSortBadge:`p-treetable-sort-badge`,tbody:`p-treetable-tbody`,row:function(e){var t=e.props,n=e.instance;return[{"p-treetable-selectable-row":n.$parentInstance.rowSelectionMode,"p-treetable-row-selected":n.selected,"p-treetable-contextmenu-row-selected":t.contextMenuSelection&&n.isSelectedWithContextMenu}]},bodyCell:function(e){return[{"p-treetable-frozen-column":e.instance.columnProp(`frozen`)}]},bodyCellContent:function(e){return[`p-treetable-body-cell-content`,{"p-treetable-body-cell-content-expander":e.instance.columnProp(`expander`)}]},nodeToggleButton:`p-treetable-node-toggle-button`,nodeToggleIcon:`p-treetable-node-toggle-icon`,pcNodeCheckbox:`p-treetable-node-checkbox`,emptyMessage:`p-treetable-empty-message`,tfoot:`p-treetable-tfoot`,footerCell:function(e){return[{"p-treetable-frozen-column":e.instance.columnProp(`frozen`)}]},footer:`p-treetable-footer`,columnResizeIndicator:`p-treetable-column-resize-indicator`},inlineStyles:{tableContainer:{overflow:`auto`},thead:{position:`sticky`},tfoot:{position:`sticky`}}}),nt={name:`BaseTreeTable`,extends:V,props:{value:{type:null,default:null},dataKey:{type:[String,Function],default:`key`},expandedKeys:{type:null,default:null},selectionKeys:{type:null,default:null},selectionMode:{type:String,default:null},metaKeySelection:{type:Boolean,default:!1},contextMenu:{type:Boolean,default:!1},contextMenuSelection:{type:Object,default:null},rows:{type:Number,default:0},first:{type:Number,default:0},totalRecords:{type:Number,default:0},paginator:{type:Boolean,default:!1},paginatorPosition:{type:String,default:`bottom`},alwaysShowPaginator:{type:Boolean,default:!0},paginatorTemplate:{type:String,default:`FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown`},pageLinkSize:{type:Number,default:5},rowsPerPageOptions:{type:Array,default:null},currentPageReportTemplate:{type:String,default:`({currentPage} of {totalPages})`},lazy:{type:Boolean,default:!1},loading:{type:Boolean,default:!1},loadingIcon:{type:String,default:void 0},loadingMode:{type:String,default:`mask`},rowHover:{type:Boolean,default:!1},autoLayout:{type:Boolean,default:!1},sortField:{type:[String,Function],default:null},sortOrder:{type:Number,default:null},defaultSortOrder:{type:Number,default:1},multiSortMeta:{type:Array,default:null},sortMode:{type:String,default:`single`},removableSort:{type:Boolean,default:!1},filters:{type:Object,default:null},filterMode:{type:String,default:`lenient`},filterLocale:{type:String,default:void 0},resizableColumns:{type:Boolean,default:!1},columnResizeMode:{type:String,default:`fit`},indentation:{type:Number,default:1},showGridlines:{type:Boolean,default:!1},scrollable:{type:Boolean,default:!1},scrollHeight:{type:String,default:null},size:{type:String,default:null},tableStyle:{type:null,default:null},tableClass:{type:[String,Object],default:null},tableProps:{type:Object,default:null}},style:tt,provide:function(){return{$pcTreeTable:this,$parentInstance:this}}},rt={name:`FooterCell`,hostName:`TreeTable`,extends:V,props:{column:{type:Object,default:null},index:{type:Number,default:null}},data:function(){return{styleObject:{}}},mounted:function(){this.columnProp(`frozen`)&&this.updateStickyPosition()},updated:function(){this.columnProp(`frozen`)&&this.updateStickyPosition()},methods:{columnProp:function(e){return U(this.column,e)},getColumnPT:function(e){var t={props:this.column.props,parent:{instance:this,props:this.$props,state:this.$data},context:{index:this.index,frozen:this.columnProp(`frozen`),size:this.$parentInstance?.size}};return u(this.ptm(`column.${e}`,{column:t}),this.ptm(`column.${e}`,t),this.ptmo(this.getColumnProp(),e,t))},getColumnProp:function(){return this.column.props&&this.column.props.pt?this.column.props.pt:void 0},updateStickyPosition:function(){if(this.columnProp(`frozen`))if(this.columnProp(`alignFrozen`)===`right`){var e=0,t=D(this.$el,`[data-p-frozen-column="true"]`);t&&(e=P(t)+parseFloat(t.style[`inset-inline-end`]||0)),this.styleObject.insetInlineEnd=e+`px`}else{var n=0,r=re(this.$el,`[data-p-frozen-column="true"]`);r&&(n=P(r)+parseFloat(r.style[`inset-inline-start`]||0)),this.styleObject.insetInlineStart=n+`px`}}},computed:{containerClass:function(){return[this.columnProp(`footerClass`),this.columnProp(`class`),this.cx(`footerCell`)]},containerStyle:function(){var e=this.columnProp(`footerStyle`),t=this.columnProp(`style`);return this.columnProp(`frozen`)?[t,e,this.styleObject]:[t,e]}}};function it(e){"@babel/helpers - typeof";return it=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},it(e)}function at(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function ot(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?at(Object(n),!0).forEach(function(t){st(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):at(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function st(e,t,n){return(t=ct(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function ct(e){var t=lt(e,`string`);return it(t)==`symbol`?t:t+``}function lt(e,t){if(it(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(it(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var ut=[`data-p-frozen-column`];function dt(e,n,r,i,a,o){return w(),S(`td`,u({style:o.containerStyle,class:o.containerClass,role:`cell`},ot(ot({},o.getColumnPT(`root`)),o.getColumnPT(`footerCell`)),{"data-p-frozen-column":o.columnProp(`frozen`)}),[r.column.children&&r.column.children.footer?(w(),d(t(r.column.children.footer),{key:0,column:r.column},null,8,[`column`])):h(``,!0),o.columnProp(`footer`)?(w(),S(`span`,u({key:1,class:e.cx(`columnFooter`)},o.getColumnPT(`columnFooter`)),C(o.columnProp(`footer`)),17)):h(``,!0)],16,ut)}rt.render=dt;var ft={name:`HeaderCell`,hostName:`TreeTable`,extends:V,emits:[`column-click`,`column-resizestart`],props:{column:{type:Object,default:null},resizableColumns:{type:Boolean,default:!1},sortField:{type:[String,Function],default:null},sortOrder:{type:Number,default:null},multiSortMeta:{type:Array,default:null},sortMode:{type:String,default:`single`},index:{type:Number,default:null}},data:function(){return{styleObject:{}}},mounted:function(){this.columnProp(`frozen`)&&this.updateStickyPosition()},updated:function(){this.columnProp(`frozen`)&&this.updateStickyPosition()},methods:{columnProp:function(e){return U(this.column,e)},getColumnPT:function(e){var t={props:this.column.props,parent:{instance:this,props:this.$props,state:this.$data},context:{index:this.index,sorted:this.isColumnSorted(),frozen:this.$parentInstance.scrollable&&this.columnProp(`frozen`),resizable:this.resizableColumns,scrollable:this.$parentInstance.scrollable,showGridlines:this.$parentInstance.showGridlines,size:this.$parentInstance?.size}};return u(this.ptm(`column.${e}`,{column:t}),this.ptm(`column.${e}`,t),this.ptmo(this.getColumnProp(),e,t))},getColumnProp:function(){return this.column.props&&this.column.props.pt?this.column.props.pt:void 0},updateStickyPosition:function(){if(this.columnProp(`frozen`)){if(this.columnProp(`alignFrozen`)===`right`){var e=0,t=D(this.$el,`[data-p-frozen-column="true"]`);t&&(e=P(t)+parseFloat(t.style[`inset-inline-end`]||0)),this.styleObject.insetInlineEnd=e+`px`}else{var n=0,r=re(this.$el,`[data-p-frozen-column="true"]`);r&&(n=P(r)+parseFloat(r.style[`inset-inline-start`]||0)),this.styleObject.insetInlineStart=n+`px`}var i=this.$el.parentElement.nextElementSibling;if(i){var a=ee(this.$el);i.children[a].style[`inset-inline-start`]=this.styleObject[`inset-inline-start`],i.children[a].style[`inset-inline-end`]=this.styleObject[`inset-inline-end`]}}},onClick:function(e){this.$emit(`column-click`,{originalEvent:e,column:this.column})},onKeyDown:function(e){(e.code===`Enter`||e.code===`NumpadEnter`||e.code===`Space`)&&e.currentTarget.nodeName===`TH`&&E(e.currentTarget,`data-p-sortable-column`)&&(this.$emit(`column-click`,{originalEvent:e,column:this.column}),e.preventDefault())},onResizeStart:function(e){this.$emit(`column-resizestart`,e)},getMultiSortMetaIndex:function(){for(var e=-1,t=0;t<this.multiSortMeta.length;t++){var n=this.multiSortMeta[t];if(n.field===this.columnProp(`field`)||n.field===this.columnProp(`sortField`)){e=t;break}}return e},isMultiSorted:function(){return this.columnProp(`sortable`)&&this.getMultiSortMetaIndex()>-1},isColumnSorted:function(){return this.sortMode===`single`?this.sortField&&(this.sortField===this.columnProp(`field`)||this.sortField===this.columnProp(`sortField`)):this.isMultiSorted()}},computed:{containerClass:function(){return[this.columnProp(`headerClass`),this.columnProp(`class`),this.cx(`headerCell`)]},containerStyle:function(){var e=this.columnProp(`headerStyle`),t=this.columnProp(`style`);return this.columnProp(`frozen`)?[t,e,this.styleObject]:[t,e]},sortState:function(){var e=!1,t=null;if(this.sortMode===`single`)e=this.sortField&&(this.sortField===this.columnProp(`field`)||this.sortField===this.columnProp(`sortField`)),t=e?this.sortOrder:0;else if(this.sortMode===`multiple`){var n=this.getMultiSortMetaIndex();n>-1&&(e=!0,t=this.multiSortMeta[n].order)}return{sorted:e,sortOrder:t}},sortableColumnIcon:function(){var e=this.sortState,t=e.sorted,n=e.sortOrder;return t?t&&n>0?Se:t&&n<0?De:null:ke},ariaSort:function(){if(this.columnProp(`sortable`)){var e=this.sortState,t=e.sorted,n=e.sortOrder;return t&&n<0?`descending`:t&&n>0?`ascending`:`none`}else return null}},components:{Badge:se,SortAltIcon:ke,SortAmountUpAltIcon:Se,SortAmountDownIcon:De}};function pt(e){"@babel/helpers - typeof";return pt=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},pt(e)}function mt(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function ht(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?mt(Object(n),!0).forEach(function(t){gt(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):mt(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function gt(e,t,n){return(t=_t(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function _t(e){var t=vt(e,`string`);return pt(t)==`symbol`?t:t+``}function vt(e,t){if(pt(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(pt(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var yt=[`tabindex`,`aria-sort`,`data-p-sortable-column`,`data-p-resizable-column`,`data-p-sorted`,`data-p-frozen-column`];function bt(e,n,r,i,a,s){var c=o(`Badge`);return w(),S(`th`,u({class:s.containerClass,style:[s.containerStyle],onClick:n[1]||=function(){return s.onClick&&s.onClick.apply(s,arguments)},onKeydown:n[2]||=function(){return s.onKeyDown&&s.onKeyDown.apply(s,arguments)},tabindex:s.columnProp(`sortable`)?`0`:null,"aria-sort":s.ariaSort,role:`columnheader`},ht(ht({},s.getColumnPT(`root`)),s.getColumnPT(`headerCell`)),{"data-p-sortable-column":s.columnProp(`sortable`),"data-p-resizable-column":r.resizableColumns,"data-p-sorted":s.isColumnSorted(),"data-p-frozen-column":s.columnProp(`frozen`)}),[r.resizableColumns&&!s.columnProp(`frozen`)?(w(),S(`span`,u({key:0,class:e.cx(`columnResizer`),onMousedown:n[0]||=function(){return s.onResizeStart&&s.onResizeStart.apply(s,arguments)}},s.getColumnPT(`columnResizer`)),null,16)):h(``,!0),b(`div`,u({class:e.cx(`columnHeaderContent`)},s.getColumnPT(`columnHeaderContent`)),[r.column.children&&r.column.children.header?(w(),d(t(r.column.children.header),{key:0,column:r.column},null,8,[`column`])):h(``,!0),s.columnProp(`header`)?(w(),S(`span`,u({key:1,class:e.cx(`columnTitle`)},s.getColumnPT(`columnTitle`)),C(s.columnProp(`header`)),17)):h(``,!0),s.columnProp(`sortable`)?(w(),S(`span`,f(u({key:2},s.getColumnPT(`sort`))),[(w(),d(t(r.column.children&&r.column.children.sorticon||s.sortableColumnIcon),u({sorted:s.sortState.sorted,sortOrder:s.sortState.sortOrder,class:e.cx(`sortIcon`)},s.getColumnPT(`sortIcon`)),null,16,[`sorted`,`sortOrder`,`class`]))],16)):h(``,!0),s.isMultiSorted()?(w(),d(c,u({key:3,class:e.cx(`pcSortBadge`)},s.getColumnPT(`pcSortBadge`),{value:s.getMultiSortMetaIndex()+1,size:`small`}),null,16,[`class`,`value`])):h(``,!0)],16)],16,yt)}ft.render=bt;var xt={name:`BodyCell`,hostName:`TreeTable`,extends:V,emits:[`node-toggle`,`checkbox-toggle`],props:{node:{type:Object,default:null},column:{type:Object,default:null},level:{type:Number,default:0},indentation:{type:Number,default:1},leaf:{type:Boolean,default:!1},expanded:{type:Boolean,default:!1},selectionMode:{type:String,default:null},checked:{type:Boolean,default:!1},partialChecked:{type:Boolean,default:!1},templates:{type:Object,default:null},index:{type:Number,default:null},loadingMode:{type:String,default:`mask`}},data:function(){return{styleObject:{}}},mounted:function(){this.columnProp(`frozen`)&&this.updateStickyPosition()},updated:function(){this.columnProp(`frozen`)&&this.updateStickyPosition()},methods:{toggle:function(){this.$emit(`node-toggle`,this.node)},columnProp:function(e){return U(this.column,e)},getColumnPT:function(e){var t={props:this.column.props,parent:{instance:this,props:this.$props,state:this.$data},context:{index:this.index,selectable:this.$parentInstance.rowHover||this.$parentInstance.rowSelectionMode,selected:this.$parent.selected,frozen:this.columnProp(`frozen`),scrollable:this.$parentInstance.scrollable,showGridlines:this.$parentInstance.showGridlines,size:this.$parentInstance?.size,node:this.node}};return u(this.ptm(`column.${e}`,{column:t}),this.ptm(`column.${e}`,t),this.ptmo(this.getColumnProp(),e,t))},getColumnProp:function(){return this.column.props&&this.column.props.pt?this.column.props.pt:void 0},getColumnCheckboxPT:function(e){var t={props:this.column.props,parent:{instance:this,props:this.$props,state:this.$data},context:{checked:this.checked,partialChecked:this.partialChecked,node:this.node}};return u(this.ptm(`column.${e}`,{column:t}),this.ptm(`column.${e}`,t),this.ptmo(this.getColumnProp(),e,t))},updateStickyPosition:function(){if(this.columnProp(`frozen`))if(this.columnProp(`alignFrozen`)===`right`){var e=0,t=D(this.$el,`[data-p-frozen-column="true"]`);t&&(e=P(t)+parseFloat(t.style[`inset-inline-end`]||0)),this.styleObject.insetInlineEnd=e+`px`}else{var n=0,r=re(this.$el,`[data-p-frozen-column="true"]`);r&&(n=P(r)+parseFloat(r.style[`inset-inline-start`]||0)),this.styleObject.insetInlineStart=n+`px`}},resolveFieldData:function(e,t){return B(e,t)},toggleCheckbox:function(){this.$emit(`checkbox-toggle`)}},computed:{containerClass:function(){return[this.columnProp(`bodyClass`),this.columnProp(`class`),this.cx(`bodyCell`)]},containerStyle:function(){var e=this.columnProp(`bodyStyle`),t=this.columnProp(`style`);return this.columnProp(`frozen`)?[t,e,this.styleObject]:[t,e]},togglerStyle:function(){return{marginLeft:this.level*this.indentation+`rem`,visibility:this.leaf?`hidden`:`visible`}},checkboxSelectionMode:function(){return this.selectionMode===`checkbox`}},components:{Checkbox:Oe,ChevronRightIcon:ve,ChevronDownIcon:_e,CheckIcon:le,MinusIcon:we,SpinnerIcon:pe},directives:{ripple:H}};function St(e){"@babel/helpers - typeof";return St=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},St(e)}function Ct(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function wt(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?Ct(Object(n),!0).forEach(function(t){Tt(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):Ct(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function Tt(e,t,n){return(t=Et(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function Et(e){var t=Dt(e,`string`);return St(t)==`symbol`?t:t+``}function Dt(e,t){if(St(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(St(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var Ot=[`data-p-frozen-column`];function kt(e,n,r,i,c,l){var f=o(`SpinnerIcon`),p=o(`Checkbox`),g=m(`ripple`);return w(),S(`td`,u({style:l.containerStyle,class:l.containerClass,role:`cell`},wt(wt({},l.getColumnPT(`root`)),l.getColumnPT(`bodyCell`)),{"data-p-frozen-column":l.columnProp(`frozen`)}),[b(`div`,u({class:e.cx(`bodyCellContent`)},l.getColumnPT(`bodyCellContent`)),[l.columnProp(`expander`)?s((w(),S(`button`,u({key:0,type:`button`,class:e.cx(`nodeToggleButton`),onClick:n[0]||=function(){return l.toggle&&l.toggle.apply(l,arguments)},style:l.togglerStyle,tabindex:`-1`},l.getColumnPT(`nodeToggleButton`),{"data-pc-group-section":`rowactionbutton`}),[r.node.loading&&r.loadingMode===`icon`?(w(),S(y,{key:0},[r.templates.nodetoggleicon?(w(),d(t(r.templates.nodetoggleicon),{key:0})):h(``,!0),r.templates.nodetogglericon?(w(),d(t(r.templates.nodetogglericon),{key:1})):(w(),d(f,u({key:2,spin:``},e.ptm(`nodetoggleicon`)),null,16))],64)):(w(),S(y,{key:1},[r.column.children&&r.column.children.rowtoggleicon?(w(),d(t(r.column.children.rowtoggleicon),{key:0,node:r.node,expanded:r.expanded,class:x(e.cx(`nodeToggleIcon`))},null,8,[`node`,`expanded`,`class`])):r.templates.nodetoggleicon?(w(),d(t(r.templates.nodetoggleicon),{key:1,node:r.node,expanded:r.expanded,class:x(e.cx(`nodeToggleIcon`))},null,8,[`node`,`expanded`,`class`])):r.column.children&&r.column.children.rowtogglericon?(w(),d(t(r.column.children.rowtogglericon),{key:2,node:r.node,expanded:r.expanded,class:x(e.cx(`nodeToggleIcon`))},null,8,[`node`,`expanded`,`class`])):r.expanded?(w(),d(t(r.node.expandedIcon?`span`:`ChevronDownIcon`),u({key:3,class:e.cx(`nodeToggleIcon`)},l.getColumnPT(`nodeToggleIcon`)),null,16,[`class`])):(w(),d(t(r.node.collapsedIcon?`span`:`ChevronRightIcon`),u({key:4,class:e.cx(`nodeToggleIcon`)},l.getColumnPT(`nodeToggleIcon`)),null,16,[`class`]))],64))],16)),[[g]]):h(``,!0),l.checkboxSelectionMode&&l.columnProp(`expander`)?(w(),d(p,{key:1,modelValue:r.checked,binary:!0,class:x(e.cx(`pcNodeCheckbox`)),disabled:r.node.selectable===!1,onChange:l.toggleCheckbox,tabindex:-1,indeterminate:r.partialChecked,unstyled:e.unstyled,pt:l.getColumnCheckboxPT(`pcNodeCheckbox`),"data-p-partialchecked":r.partialChecked},{icon:a(function(e){return[r.templates.checkboxicon?(w(),d(t(r.templates.checkboxicon),{key:0,checked:e.checked,partialChecked:r.partialChecked,class:x(e.class)},null,8,[`checked`,`partialChecked`,`class`])):h(``,!0)]}),_:1},8,[`modelValue`,`class`,`disabled`,`onChange`,`indeterminate`,`unstyled`,`pt`,`data-p-partialchecked`])):h(``,!0),r.column.children&&r.column.children.body?(w(),d(t(r.column.children.body),{key:2,node:r.node,column:r.column},null,8,[`node`,`column`])):(w(),S(y,{key:3},[v(C(l.resolveFieldData(r.node.data,l.columnProp(`field`))),1)],64))],16)],16,Ot)}xt.render=kt;function At(e){"@babel/helpers - typeof";return At=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},At(e)}function jt(e,t){var n=typeof Symbol<`u`&&e[Symbol.iterator]||e[`@@iterator`];if(!n){if(Array.isArray(e)||(n=zt(e))||t){n&&(e=n);var r=0,i=function(){};return{s:i,n:function(){return r>=e.length?{done:!0}:{done:!1,value:e[r++]}},e:function(e){throw e},f:i}}throw TypeError(`Invalid attempt to iterate non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}var a,o=!0,s=!1;return{s:function(){n=n.call(e)},n:function(){var e=n.next();return o=e.done,e},e:function(e){s=!0,a=e},f:function(){try{o||n.return==null||n.return()}finally{if(s)throw a}}}}function Mt(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function Nt(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?Mt(Object(n),!0).forEach(function(t){Pt(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):Mt(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function Pt(e,t,n){return(t=Ft(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function Ft(e){var t=It(e,`string`);return At(t)==`symbol`?t:t+``}function It(e,t){if(At(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(At(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}function Lt(e){return Vt(e)||Bt(e)||zt(e)||Rt()}function Rt(){throw TypeError(`Invalid attempt to spread non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}function zt(e,t){if(e){if(typeof e==`string`)return Ht(e,t);var n={}.toString.call(e).slice(8,-1);return n===`Object`&&e.constructor&&(n=e.constructor.name),n===`Map`||n===`Set`?Array.from(e):n===`Arguments`||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?Ht(e,t):void 0}}function Bt(e){if(typeof Symbol<`u`&&e[Symbol.iterator]!=null||e[`@@iterator`]!=null)return Array.from(e)}function Vt(e){if(Array.isArray(e))return Ht(e)}function Ht(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=Array(t);n<t;n++)r[n]=e[n];return r}var Ut={name:`TreeTableRow`,hostName:`TreeTable`,extends:V,emits:[`node-click`,`node-toggle`,`checkbox-change`,`nodeClick`,`nodeToggle`,`checkboxChange`,`row-rightclick`,`rowRightclick`],props:{node:{type:null,default:null},dataKey:{type:[String,Function],default:`key`},parentNode:{type:null,default:null},columns:{type:null,default:null},expandedKeys:{type:null,default:null},selectionKeys:{type:null,default:null},selectionMode:{type:String,default:null},level:{type:Number,default:0},indentation:{type:Number,default:1},tabindex:{type:Number,default:-1},ariaSetSize:{type:Number,default:null},ariaPosInset:{type:Number,default:null},loadingMode:{type:String,default:`mask`},templates:{type:Object,default:null},contextMenu:{type:Boolean,default:!1},contextMenuSelection:{type:Object,default:null}},nodeTouched:!1,methods:{columnProp:function(e,t){return U(e,t)},toggle:function(){this.$emit(`node-toggle`,this.node)},onClick:function(e){R(e.target)||E(e.target,`data-pc-section`)===`nodetogglebutton`||E(e.target,`data-pc-section`)===`nodetoggleicon`||e.target.tagName===`path`||(this.setTabIndexForSelectionMode(e,this.nodeTouched),this.$emit(`node-click`,{originalEvent:e,nodeTouched:this.nodeTouched,node:this.node}),this.nodeTouched=!1)},onRowRightClick:function(e){this.$emit(`row-rightclick`,{originalEvent:e,node:this.node})},onTouchEnd:function(){this.nodeTouched=!0},nodeKey:function(e){return B(e,this.dataKey)},onKeyDown:function(e,t){switch(e.code){case`ArrowDown`:this.onArrowDownKey(e);break;case`ArrowUp`:this.onArrowUpKey(e);break;case`ArrowLeft`:this.onArrowLeftKey(e);break;case`ArrowRight`:this.onArrowRightKey(e);break;case`Home`:this.onHomeKey(e);break;case`End`:this.onEndKey(e);break;case`Enter`:case`NumpadEnter`:case`Space`:R(e.target)||this.onEnterKey(e,t);break;case`Tab`:this.onTabKey(e);break}},onArrowDownKey:function(e){var t=e.currentTarget.nextElementSibling;t&&this.focusRowChange(e.currentTarget,t),e.preventDefault()},onArrowUpKey:function(e){var t=e.currentTarget.previousElementSibling;t&&this.focusRowChange(e.currentTarget,t),e.preventDefault()},onArrowRightKey:function(e){var t=this,n=T(e.currentTarget,`button`).style.visibility===`hidden`,r=T(this.$refs.node,`[data-pc-section="nodetogglebutton"]`);n||(!this.expanded&&r.click(),this.$nextTick(function(){t.onArrowDownKey(e)}),e.preventDefault())},onArrowLeftKey:function(e){if(!(this.level===0&&!this.expanded)){var t=e.currentTarget,n=T(t,`button`).style.visibility===`hidden`,r=T(t,`[data-pc-section="nodetogglebutton"]`);if(this.expanded&&!n){r.click();return}var i=this.findBeforeClickableNode(t);i&&this.focusRowChange(t,i)}},onHomeKey:function(e){var t=T(e.currentTarget.parentElement,`tr[aria-level="${this.level+1}"]`);t&&N(t),e.preventDefault()},onEndKey:function(e){var t=O(e.currentTarget.parentElement,`tr[aria-level="${this.level+1}"]`),n=t[t.length-1];N(n),e.preventDefault()},onEnterKey:function(e){if(e.preventDefault(),this.setTabIndexForSelectionMode(e,this.nodeTouched),this.selectionMode===`checkbox`){this.toggleCheckbox();return}this.$emit(`node-click`,{originalEvent:e,nodeTouched:this.nodeTouched,node:this.node}),this.nodeTouched=!1},onTabKey:function(){var e=Lt(O(this.$refs.node.parentElement,`tr`)),t=e.some(function(e){return E(e,`data-p-selected`)||e.getAttribute(`aria-checked`)===`true`});if(e.forEach(function(e){e.tabIndex=-1}),t){var n=e.filter(function(e){return E(e,`data-p-selected`)||e.getAttribute(`aria-checked`)===`true`});n[0].tabIndex=0;return}e[0].tabIndex=0},focusRowChange:function(e,t){e.tabIndex=`-1`,t.tabIndex=`0`,N(t)},findBeforeClickableNode:function(e){var t=e.previousElementSibling;if(t){var n=t.querySelector(`button`);return n&&n.style.visibility!==`hidden`?t:this.findBeforeClickableNode(t)}return null},toggleCheckbox:function(){var e=this.selectionKeys?Nt({},this.selectionKeys):{},t=!this.checked;this.propagateDown(this.node,t,e),this.$emit(`checkbox-change`,{node:this.node,check:t,selectionKeys:e})},propagateDown:function(e,t,n){if(t?n[this.nodeKey(e)]={checked:!0,partialChecked:!1}:delete n[this.nodeKey(e)],e.children&&e.children.length){var r=jt(e.children),i;try{for(r.s();!(i=r.n()).done;){var a=i.value;this.propagateDown(a,t,n)}}catch(e){r.e(e)}finally{r.f()}}},propagateUp:function(e){var t=e.check,n=Nt({},e.selectionKeys),r=0,i=!1,a=jt(this.node.children),o;try{for(a.s();!(o=a.n()).done;){var s=o.value;n[this.nodeKey(s)]&&n[this.nodeKey(s)].checked?r++:n[this.nodeKey(s)]&&n[this.nodeKey(s)].partialChecked&&(i=!0)}}catch(e){a.e(e)}finally{a.f()}t&&r===this.node.children.length?n[this.nodeKey(this.node)]={checked:!0,partialChecked:!1}:(t||delete n[this.nodeKey(this.node)],i||r>0&&r!==this.node.children.length?n[this.nodeKey(this.node)]={checked:!1,partialChecked:!0}:n[this.nodeKey(this.node)]={checked:!1,partialChecked:!1}),this.$emit(`checkbox-change`,{node:e.node,check:e.check,selectionKeys:n})},onCheckboxChange:function(e){var t=e.check,n=Nt({},e.selectionKeys),r=0,i=!1,a=jt(this.node.children),o;try{for(a.s();!(o=a.n()).done;){var s=o.value;n[this.nodeKey(s)]&&n[this.nodeKey(s)].checked?r++:n[this.nodeKey(s)]&&n[this.nodeKey(s)].partialChecked&&(i=!0)}}catch(e){a.e(e)}finally{a.f()}t&&r===this.node.children.length?n[this.nodeKey(this.node)]={checked:!0,partialChecked:!1}:(t||delete n[this.nodeKey(this.node)],i||r>0&&r!==this.node.children.length?n[this.nodeKey(this.node)]={checked:!1,partialChecked:!0}:n[this.nodeKey(this.node)]={checked:!1,partialChecked:!1}),this.$emit(`checkbox-change`,{node:e.node,check:e.check,selectionKeys:n})},setTabIndexForSelectionMode:function(e,t){if(this.selectionMode!==null){var n=Lt(O(this.$refs.node.parentElement,`tr`));e.currentTarget.tabIndex=t===!1?-1:0,n.every(function(e){return e.tabIndex===-1})&&(n[0].tabIndex=0)}}},computed:{containerClass:function(){return[this.node.styleClass,this.cx(`row`)]},expanded:function(){return this.expandedKeys&&this.expandedKeys[this.nodeKey(this.node)]===!0},leaf:function(){return this.node.leaf!==!1&&!(this.node.children&&this.node.children.length)},selected:function(){return this.selectionMode&&this.selectionKeys?this.selectionKeys[this.nodeKey(this.node)]===!0:!1},isSelectedWithContextMenu:function(){return this.node&&this.contextMenuSelection?I(this.node,this.contextMenuSelection,this.dataKey):!1},checked:function(){return this.selectionKeys?this.selectionKeys[this.nodeKey(this.node)]&&this.selectionKeys[this.nodeKey(this.node)].checked:!1},partialChecked:function(){return this.selectionKeys?this.selectionKeys[this.nodeKey(this.node)]&&this.selectionKeys[this.nodeKey(this.node)].partialChecked:!1},getAriaSelected:function(){return this.selectionMode===`single`||this.selectionMode===`multiple`?this.selected:null},ptmOptions:function(){return{context:{selectable:this.$parentInstance.rowHover||this.$parentInstance.rowSelectionMode,selected:this.selected,scrollable:this.$parentInstance.scrollable}}}},components:{TTBodyCell:xt}},Wt=[`tabindex`,`aria-expanded`,`aria-level`,`aria-setsize`,`aria-posinset`,`aria-selected`,`aria-checked`,`data-p-selected`,`data-p-selected-contextmenu`];function Gt(e,t,n,i,a,s){var c=o(`TTBodyCell`),l=o(`TreeTableRow`,!0);return w(),S(y,null,[b(`tr`,u({ref:`node`,class:s.containerClass,style:n.node.style,tabindex:n.tabindex,role:`row`,"aria-expanded":n.node.children&&n.node.children.length?s.expanded:void 0,"aria-level":n.level+1,"aria-setsize":n.ariaSetSize,"aria-posinset":n.ariaPosInset,"aria-selected":s.getAriaSelected,"aria-checked":s.checked||void 0,onClick:t[1]||=function(){return s.onClick&&s.onClick.apply(s,arguments)},onKeydown:t[2]||=function(){return s.onKeyDown&&s.onKeyDown.apply(s,arguments)},onTouchend:t[3]||=function(){return s.onTouchEnd&&s.onTouchEnd.apply(s,arguments)},onContextmenu:t[4]||=function(){return s.onRowRightClick&&s.onRowRightClick.apply(s,arguments)}},e.ptm(`row`,s.ptmOptions),{"data-p-selected":s.selected,"data-p-selected-contextmenu":n.contextMenuSelection&&s.isSelectedWithContextMenu}),[(w(!0),S(y,null,r(n.columns,function(r,i){return w(),S(y,{key:s.columnProp(r,`columnKey`)||s.columnProp(r,`field`)||i},[s.columnProp(r,`hidden`)?h(``,!0):(w(),d(c,{key:0,column:r,node:n.node,level:n.level,leaf:s.leaf,indentation:n.indentation,expanded:s.expanded,selectionMode:n.selectionMode,checked:s.checked,partialChecked:s.partialChecked,templates:n.templates,onNodeToggle:t[0]||=function(t){return e.$emit(`node-toggle`,t)},onCheckboxToggle:s.toggleCheckbox,index:i,loadingMode:n.loadingMode,unstyled:e.unstyled,pt:e.pt},null,8,[`column`,`node`,`level`,`leaf`,`indentation`,`expanded`,`selectionMode`,`checked`,`partialChecked`,`templates`,`onCheckboxToggle`,`index`,`loadingMode`,`unstyled`,`pt`]))],64)}),128))],16,Wt),s.expanded&&n.node.children&&n.node.children.length?(w(!0),S(y,{key:0},r(n.node.children,function(r){return w(),d(l,{key:s.nodeKey(r),dataKey:n.dataKey,columns:n.columns,node:r,parentNode:n.node,level:n.level+1,expandedKeys:n.expandedKeys,selectionMode:n.selectionMode,selectionKeys:n.selectionKeys,contextMenu:n.contextMenu,contextMenuSelection:n.contextMenuSelection,indentation:n.indentation,ariaPosInset:n.node.children.indexOf(r)+1,ariaSetSize:n.node.children.length,templates:n.templates,onNodeToggle:t[5]||=function(t){return e.$emit(`node-toggle`,t)},onNodeClick:t[6]||=function(t){return e.$emit(`node-click`,t)},onRowRightclick:t[7]||=function(t){return e.$emit(`row-rightclick`,t)},onCheckboxChange:s.onCheckboxChange,unstyled:e.unstyled,pt:e.pt},null,8,[`dataKey`,`columns`,`node`,`parentNode`,`level`,`expandedKeys`,`selectionMode`,`selectionKeys`,`contextMenu`,`contextMenuSelection`,`indentation`,`ariaPosInset`,`ariaSetSize`,`templates`,`onCheckboxChange`,`unstyled`,`pt`])}),128)):h(``,!0)],64)}Ut.render=Gt;function X(e){"@babel/helpers - typeof";return X=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},X(e)}function Kt(e,t){var n=typeof Symbol<`u`&&e[Symbol.iterator]||e[`@@iterator`];if(!n){if(Array.isArray(e)||(n=$t(e))||t){n&&(e=n);var r=0,i=function(){};return{s:i,n:function(){return r>=e.length?{done:!0}:{done:!1,value:e[r++]}},e:function(e){throw e},f:i}}throw TypeError(`Invalid attempt to iterate non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}var a,o=!0,s=!1;return{s:function(){n=n.call(e)},n:function(){var e=n.next();return o=e.done,e},e:function(e){s=!0,a=e},f:function(){try{o||n.return==null||n.return()}finally{if(s)throw a}}}}function qt(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function Z(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?qt(Object(n),!0).forEach(function(t){Jt(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):qt(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function Jt(e,t,n){return(t=Yt(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function Yt(e){var t=Xt(e,`string`);return X(t)==`symbol`?t:t+``}function Xt(e,t){if(X(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(X(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}function Zt(e){return tn(e)||en(e)||$t(e)||Qt()}function Qt(){throw TypeError(`Invalid attempt to spread non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}function $t(e,t){if(e){if(typeof e==`string`)return nn(e,t);var n={}.toString.call(e).slice(8,-1);return n===`Object`&&e.constructor&&(n=e.constructor.name),n===`Map`||n===`Set`?Array.from(e):n===`Arguments`||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?nn(e,t):void 0}}function en(e){if(typeof Symbol<`u`&&e[Symbol.iterator]!=null||e[`@@iterator`]!=null)return Array.from(e)}function tn(e){if(Array.isArray(e))return nn(e)}function nn(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=Array(t);n<t;n++)r[n]=e[n];return r}var rn={name:`TreeTable`,extends:nt,inheritAttrs:!1,emits:[`node-expand`,`node-collapse`,`update:expandedKeys`,`update:selectionKeys`,`node-select`,`node-unselect`,`update:first`,`update:rows`,`page`,`update:sortField`,`update:sortOrder`,`update:multiSortMeta`,`sort`,`filter`,`column-resize-end`,`update:contextMenuSelection`,`row-contextmenu`],provide:function(){return{$columns:this.d_columns}},data:function(){return{d_expandedKeys:this.expandedKeys||{},d_first:this.first,d_rows:this.rows,d_sortField:this.sortField,d_sortOrder:this.sortOrder,d_multiSortMeta:this.multiSortMeta?Zt(this.multiSortMeta):[],hasASelectedNode:!1,d_columns:new de({type:`Column`})}},documentColumnResizeListener:null,documentColumnResizeEndListener:null,lastResizeHelperX:null,resizeColumnElement:null,watch:{expandedKeys:function(e){this.d_expandedKeys=e},first:function(e){this.d_first=e},rows:function(e){this.d_rows=e},sortField:function(e){this.d_sortField=e},sortOrder:function(e){this.d_sortOrder=e},multiSortMeta:function(e){this.d_multiSortMeta=e}},beforeUnmount:function(){this.destroyStyleElement(),this.d_columns.clear()},methods:{columnProp:function(e,t){return U(e,t)},ptHeaderCellOptions:function(e){return{context:{frozen:this.columnProp(e,`frozen`)}}},onNodeToggle:function(e){var t=this.nodeKey(e);this.d_expandedKeys[t]?(delete this.d_expandedKeys[t],this.$emit(`node-collapse`,e)):(this.d_expandedKeys[t]=!0,this.$emit(`node-expand`,e)),this.d_expandedKeys=Z({},this.d_expandedKeys),this.$emit(`update:expandedKeys`,this.d_expandedKeys)},onNodeClick:function(e){if(this.rowSelectionMode&&e.node.selectable!==!1){var t=!e.nodeTouched&&this.metaKeySelection?this.handleSelectionWithMetaKey(e):this.handleSelectionWithoutMetaKey(e);this.$emit(`update:selectionKeys`,t)}},nodeKey:function(e){return B(e,this.dataKey)},handleSelectionWithMetaKey:function(e){var t=e.originalEvent,n=e.node,r=this.nodeKey(n),i=t.metaKey||t.ctrlKey,a=this.isNodeSelected(n),o;return a&&i?(this.isSingleSelectionMode()?o={}:(o=Z({},this.selectionKeys),delete o[r]),this.$emit(`node-unselect`,n)):(this.isSingleSelectionMode()?o={}:this.isMultipleSelectionMode()&&(o=i&&this.selectionKeys?Z({},this.selectionKeys):{}),o[r]=!0,this.$emit(`node-select`,n)),o},handleSelectionWithoutMetaKey:function(e){var t=e.node,n=this.nodeKey(t),r=this.isNodeSelected(t),i;return this.isSingleSelectionMode()?r?(i={},this.$emit(`node-unselect`,t)):(i={},i[n]=!0,this.$emit(`node-select`,t)):r?(i=Z({},this.selectionKeys),delete i[n],this.$emit(`node-unselect`,t)):(i=this.selectionKeys?Z({},this.selectionKeys):{},i[n]=!0,this.$emit(`node-select`,t)),i},onCheckboxChange:function(e){this.$emit(`update:selectionKeys`,e.selectionKeys),e.check?this.$emit(`node-select`,e.node):this.$emit(`node-unselect`,e.node)},onRowRightClick:function(e){this.contextMenu&&(k(),e.originalEvent.target.focus()),this.$emit(`update:contextMenuSelection`,e.node),this.$emit(`row-contextmenu`,e)},isSingleSelectionMode:function(){return this.selectionMode===`single`},isMultipleSelectionMode:function(){return this.selectionMode===`multiple`},onPage:function(e){this.d_first=e.first,this.d_rows=e.rows;var t=this.createLazyLoadEvent(e);t.pageCount=e.pageCount,t.page=e.page,this.d_expandedKeys={},this.$emit(`update:expandedKeys`,this.d_expandedKeys),this.$emit(`update:first`,this.d_first),this.$emit(`update:rows`,this.d_rows),this.$emit(`page`,t)},resetPage:function(){this.d_first=0,this.$emit(`update:first`,this.d_first)},getFilterColumnHeaderClass:function(e){return[this.cx(`headerCell`,{column:e}),this.columnProp(e,`filterHeaderClass`)]},onColumnHeaderClick:function(e){var t=e.originalEvent,n=e.column;if(this.columnProp(n,`sortable`)){var r=t.target,i=this.columnProp(n,`sortField`)||this.columnProp(n,`field`);(E(r,`data-p-sortable-column`)===!0||E(r,`data-pc-section`)===`columntitle`||E(r,`data-pc-section`)===`columnheadercontent`||E(r,`data-pc-section`)===`sorticon`||E(r.parentElement,`data-pc-section`)===`sorticon`||E(r.parentElement.parentElement,`data-pc-section`)===`sorticon`||r.closest(`[data-p-sortable-column="true"]`))&&(k(),this.sortMode===`single`?(this.d_sortField===i?this.removableSort&&this.d_sortOrder*-1===this.defaultSortOrder?(this.d_sortOrder=null,this.d_sortField=null):this.d_sortOrder*=-1:(this.d_sortOrder=this.defaultSortOrder,this.d_sortField=i),this.$emit(`update:sortField`,this.d_sortField),this.$emit(`update:sortOrder`,this.d_sortOrder),this.resetPage()):this.sortMode===`multiple`&&(t.metaKey||t.ctrlKey||(this.d_multiSortMeta=this.d_multiSortMeta.filter(function(e){return e.field===i})),this.addMultiSortField(i),this.$emit(`update:multiSortMeta`,this.d_multiSortMeta)),this.$emit(`sort`,this.createLazyLoadEvent(t)))}},addMultiSortField:function(e){var t=this.d_multiSortMeta.findIndex(function(t){return t.field===e});t>=0?this.removableSort&&this.d_multiSortMeta[t].order*-1===this.defaultSortOrder?this.d_multiSortMeta.splice(t,1):this.d_multiSortMeta[t]={field:e,order:this.d_multiSortMeta[t].order*-1}:this.d_multiSortMeta.push({field:e,order:this.defaultSortOrder}),this.d_multiSortMeta=Zt(this.d_multiSortMeta)},sortSingle:function(e){return this.sortNodesSingle(e)},sortNodesSingle:function(e){var t=this,n=ne();return Zt(e).sort(function(e,r){return z(B(e.data,t.d_sortField),B(r.data,t.d_sortField),t.d_sortOrder,n)}).map(function(e){return e.children&&e.children.length?Z(Z({},e),{},{children:t.sortNodesSingle(e.children)}):e})},sortMultiple:function(e){return this.sortNodesMultiple(e)},sortNodesMultiple:function(e){var t=this;return Zt(e).sort(function(e,n){return t.multisortField(e,n,0)}).map(function(e){return e.children&&e.children.length?Z(Z({},e),{},{children:t.sortNodesMultiple(e.children)}):e})},multisortField:function(e,t,n){var r=B(e.data,this.d_multiSortMeta[n].field),i=B(t.data,this.d_multiSortMeta[n].field),a=ne();return r===i?this.d_multiSortMeta.length-1>n?this.multisortField(e,t,n+1):0:z(r,i,this.d_multiSortMeta[n].order,a)},filter:function(e){var t=[],n=this.filterMode===`strict`,r=Kt(e),i;try{for(r.s();!(i=r.n()).done;){for(var a=i.value,o=Z({},a),s=!0,c=!1,l=0;l<this.columns.length;l++){var u=this.columns[l],d=this.columnProp(u,`filterField`)||this.columnProp(u,`field`);if(Object.prototype.hasOwnProperty.call(this.filters,d)){var f=this.columnProp(u,`filterMatchMode`)||`startsWith`,p={filterField:d,filterValue:this.filters[d],filterConstraint:ce.filters[f],strict:n};if((n&&!(this.findFilteredNodes(o,p)||this.isFilterMatched(o,p))||!n&&!(this.isFilterMatched(o,p)||this.findFilteredNodes(o,p)))&&(s=!1),!s)break}if(this.hasGlobalFilter()&&!c){var m=Z({},o),h={filterField:d,filterValue:this.filters.global,filterConstraint:ce.filters.contains,strict:n};(n&&(this.findFilteredNodes(m,h)||this.isFilterMatched(m,h))||!n&&(this.isFilterMatched(m,h)||this.findFilteredNodes(m,h)))&&(c=!0,o=m)}}var g=s;this.hasGlobalFilter()&&(g=s&&c),g&&t.push(o)}}catch(e){r.e(e)}finally{r.f()}var _=this.createLazyLoadEvent(event);return _.filteredValue=t,this.$emit(`filter`,_),t},findFilteredNodes:function(e,t){if(e){var n=!1;if(e.children){var r=Zt(e.children);e.children=[];var i=Kt(r),a;try{for(i.s();!(a=i.n()).done;){var o=a.value,s=Z({},o);this.isFilterMatched(s,t)&&(n=!0,e.children.push(s))}}catch(e){i.e(e)}finally{i.f()}}if(n)return!0}},isFilterMatched:function(e,t){var n=t.filterField,r=t.filterValue,i=t.filterConstraint,a=t.strict,o=!1;return i(B(e.data,n),r,this.filterLocale)&&(o=!0),(!o||a&&!this.isNodeLeaf(e))&&(o=this.findFilteredNodes(e,{filterField:n,filterValue:r,filterConstraint:i,strict:a})||o),o},isNodeSelected:function(e){return this.selectionMode&&this.selectionKeys?this.selectionKeys[this.nodeKey(e)]===!0:!1},isNodeLeaf:function(e){return e.leaf!==!1&&!(e.children&&e.children.length)},createLazyLoadEvent:function(e){var t=this,n;return this.hasFilters()&&(n={},this.columns.forEach(function(e){t.columnProp(e,`field`)&&(n[e.props.field]=t.columnProp(e,`filterMatchMode`))})),{originalEvent:e,first:this.d_first,rows:this.d_rows,sortField:this.d_sortField,sortOrder:this.d_sortOrder,multiSortMeta:this.d_multiSortMeta,filters:this.filters,filterMatchModes:n}},onColumnResizeStart:function(e){var t=ie(this.$el).left;this.resizeColumnElement=e.target.parentElement,this.columnResizing=!0,this.lastResizeHelperX=e.pageX-t+this.$el.scrollLeft,this.bindColumnResizeEvents()},onColumnResize:function(e){var t=ie(this.$el).left;this.$el.setAttribute(`data-p-unselectable-text`,`true`),!this.isUnstyled&&te(this.$el,{"user-select":`none`}),this.$refs.resizeHelper.style.height=this.$el.offsetHeight+`px`,this.$refs.resizeHelper.style.top=`0px`,this.$refs.resizeHelper.style.left=e.pageX-t+this.$el.scrollLeft+`px`,this.$refs.resizeHelper.style.display=`block`},onColumnResizeEnd:function(){var e=A(this.$el)?this.lastResizeHelperX-this.$refs.resizeHelper.offsetLeft:this.$refs.resizeHelper.offsetLeft-this.lastResizeHelperX,t=this.resizeColumnElement.offsetWidth,n=t+e,r=this.resizeColumnElement.style.minWidth||15;if(t+e>parseInt(r,10)){if(this.columnResizeMode===`fit`){var i=this.resizeColumnElement.nextElementSibling.offsetWidth-e;n>15&&i>15&&this.resizeTableCells(n,i)}else if(this.columnResizeMode===`expand`){var a=this.$refs.table.offsetWidth+e+`px`;this.resizeTableCells(n),function(e){e&&(e.style.width=e.style.minWidth=a)}(this.$refs.table)}this.$emit(`column-resize-end`,{element:this.resizeColumnElement,delta:e})}this.$refs.resizeHelper.style.display=`none`,this.resizeColumn=null,this.$el.removeAttribute(`data-p-unselectable-text`),!this.isUnstyled&&(this.$el.style[`user-select`]=``),this.unbindColumnResizeEvents()},resizeTableCells:function(e,t){var n=ee(this.resizeColumnElement),r=[];O(this.$refs.table,`thead[data-pc-section="thead"] > tr > th`).forEach(function(e){return r.push(P(e))}),this.destroyStyleElement(),this.createStyleElement();var i=``,a=`[data-pc-name="treetable"][${this.$attrSelector}] > [data-pc-section="tablecontainer"] > table[data-pc-section="table"]`;r.forEach(function(r,o){var s=o===n?e:t&&o===n+1?t:r,c=`width: ${s}px !important; max-width: ${s}px !important`;i+=`
                    ${a} > thead[data-pc-section="thead"] > tr > th:nth-child(${o+1}),
                    ${a} > tbody[data-pc-section="tbody"] > tr > td:nth-child(${o+1}),
                    ${a} > tfoot[data-pc-section="tfoot"] > tr > td:nth-child(${o+1}) {
                        ${c}
                    }
                `}),this.styleElement.innerHTML=i},bindColumnResizeEvents:function(){var e=this;this.documentColumnResizeListener||=document.addEventListener(`mousemove`,function(t){e.columnResizing&&e.onColumnResize(t)}),this.documentColumnResizeEndListener||=document.addEventListener(`mouseup`,function(){e.columnResizing&&(e.columnResizing=!1,e.onColumnResizeEnd())})},unbindColumnResizeEvents:function(){this.documentColumnResizeListener&&=(document.removeEventListener(`document`,this.documentColumnResizeListener),null),this.documentColumnResizeEndListener&&=(document.removeEventListener(`document`,this.documentColumnResizeEndListener),null)},onColumnKeyDown:function(e,t){(e.code===`Enter`||e.code===`NumpadEnter`)&&e.currentTarget.nodeName===`TH`&&E(e.currentTarget,`data-p-sortable-column`)&&this.onColumnHeaderClick(e,t)},hasColumnFilter:function(){if(this.columns){var e=Kt(this.columns),t;try{for(e.s();!(t=e.n()).done;){var n=t.value;if(n.children&&n.children.filter)return!0}}catch(t){e.e(t)}finally{e.f()}}return!1},hasFilters:function(){return this.filters&&Object.keys(this.filters).length>0&&this.filters.constructor===Object},hasGlobalFilter:function(){return this.filters&&Object.prototype.hasOwnProperty.call(this.filters,`global`)},getItemLabel:function(e){return e.data.name},createStyleElement:function(){var e;this.styleElement=document.createElement(`style`),this.styleElement.type=`text/css`,M(this.styleElement,`nonce`,(e=this.$primevue)==null||(e=e.config)==null||(e=e.csp)==null?void 0:e.nonce),document.head.appendChild(this.styleElement)},destroyStyleElement:function(){this.styleElement&&=(document.head.removeChild(this.styleElement),null)},setTabindex:function(e,t){if(this.isNodeSelected(e))return this.hasASelectedNode=!0,0;if(this.selectionMode){if(!this.isNodeSelected(e)&&t===0&&!this.hasASelectedNode)return 0}else if(!this.selectionMode&&t===0)return 0;return-1}},computed:{columns:function(){return this.d_columns.get(this)},processedData:function(){if(this.lazy)return this.value;if(this.value&&this.value.length){var e=this.value;return this.sorted&&(this.sortMode===`single`?e=this.sortSingle(e):this.sortMode===`multiple`&&(e=this.sortMultiple(e))),this.hasFilters()&&(e=this.filter(e)),e}else return null},dataToRender:function(){var e=this.processedData;if(this.paginator){var t=this.lazy?0:this.d_first;return e.slice(t,t+this.d_rows)}else return e},empty:function(){var e=this.processedData;return!e||e.length===0},sorted:function(){return this.d_sortField||this.d_multiSortMeta&&this.d_multiSortMeta.length>0},hasFooter:function(){var e=!1,t=Kt(this.columns),n;try{for(t.s();!(n=t.n()).done;){var r=n.value;if(this.columnProp(r,`footer`)||r.children&&r.children.footer){e=!0;break}}}catch(e){t.e(e)}finally{t.f()}return e},paginatorTop:function(){return this.paginator&&(this.paginatorPosition!==`bottom`||this.paginatorPosition===`both`)},paginatorBottom:function(){return this.paginator&&(this.paginatorPosition!==`top`||this.paginatorPosition===`both`)},singleSelectionMode:function(){return this.selectionMode&&this.selectionMode===`single`},multipleSelectionMode:function(){return this.selectionMode&&this.selectionMode===`multiple`},rowSelectionMode:function(){return this.singleSelectionMode||this.multipleSelectionMode},totalRecordsLength:function(){if(this.lazy)return this.totalRecords;var e=this.processedData;return e?e.length:0},dataP:function(){return L(Jt(Jt(Jt({scrollable:this.scrollable,"flex-scrollable":this.scrollable&&this.scrollHeight===`flex`},this.size,this.size),`loading`,this.loading),`empty`,this.empty))}},components:{TTRow:Ut,TTPaginator:Ce,TTHeaderCell:ft,TTFooterCell:rt,SpinnerIcon:pe}};function an(e){"@babel/helpers - typeof";return an=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},an(e)}function on(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function sn(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?on(Object(n),!0).forEach(function(t){cn(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):on(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function cn(e,t,n){return(t=ln(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function ln(e){var t=un(e,`string`);return an(t)==`symbol`?t:t+``}function un(e,t){if(an(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(an(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var dn=[`data-p`],fn=[`colspan`];function pn(e,i,s,c,l,m){var _=o(`TTPaginator`),v=o(`TTHeaderCell`),C=o(`TTRow`),T=o(`TTFooterCell`);return w(),S(`div`,u({class:e.cx(`root`),"data-scrollselectors":`.p-treetable-scrollable-body`,"data-p":m.dataP},e.ptmi(`root`)),[n(e.$slots,`default`),g(ae,{name:`p-overlay-mask`},{default:a(function(){return[e.loading&&e.loadingMode===`mask`?(w(),S(`div`,u({key:0,class:e.cx(`loading`)},e.ptm(`loading`)),[b(`div`,u({class:e.cx(`mask`)},e.ptm(`mask`)),[n(e.$slots,`loadingicon`,{class:x(e.cx(`loadingIcon`))},function(){return[(w(),d(t(e.loadingIcon?`span`:`SpinnerIcon`),u({spin:``,class:[e.cx(`loadingIcon`),e.loadingIcon]},e.ptm(`loadingIcon`)),null,16,[`class`]))]})],16)],16)):h(``,!0)]}),_:3}),e.$slots.header?(w(),S(`div`,u({key:0,class:e.cx(`header`)},e.ptm(`header`)),[n(e.$slots,`header`)],16)):h(``,!0),m.paginatorTop?(w(),d(_,{key:1,rows:l.d_rows,first:l.d_first,totalRecords:m.totalRecordsLength,pageLinkSize:e.pageLinkSize,template:e.paginatorTemplate,rowsPerPageOptions:e.rowsPerPageOptions,currentPageReportTemplate:e.currentPageReportTemplate,class:x(e.cx(`pcPaginator`,{position:`top`})),onPage:i[0]||=function(e){return m.onPage(e)},alwaysShow:e.alwaysShowPaginator,unstyled:e.unstyled,pt:e.ptm(`pcPaginator`)},p({_:2},[e.$slots.paginatorcontainer?{name:`container`,fn:a(function(t){return[n(e.$slots,`paginatorcontainer`,{first:t.first,last:t.last,rows:t.rows,page:t.page,pageCount:t.pageCount,totalRecords:t.totalRecords,firstPageCallback:t.firstPageCallback,lastPageCallback:t.lastPageCallback,prevPageCallback:t.prevPageCallback,nextPageCallback:t.nextPageCallback,rowChangeCallback:t.rowChangeCallback,pageLinks:t.pageLinks,changePageCallback:t.changePageCallback})]}),key:`0`}:void 0,e.$slots.paginatorstart?{name:`start`,fn:a(function(){return[n(e.$slots,`paginatorstart`)]}),key:`1`}:void 0,e.$slots.paginatorend?{name:`end`,fn:a(function(){return[n(e.$slots,`paginatorend`)]}),key:`2`}:void 0,e.$slots.paginatorfirstpagelinkicon?{name:`firstpagelinkicon`,fn:a(function(t){return[n(e.$slots,`paginatorfirstpagelinkicon`,{class:x(t.class)})]}),key:`3`}:void 0,e.$slots.paginatorprevpagelinkicon?{name:`prevpagelinkicon`,fn:a(function(t){return[n(e.$slots,`paginatorprevpagelinkicon`,{class:x(t.class)})]}),key:`4`}:void 0,e.$slots.paginatornextpagelinkicon?{name:`nextpagelinkicon`,fn:a(function(t){return[n(e.$slots,`paginatornextpagelinkicon`,{class:x(t.class)})]}),key:`5`}:void 0,e.$slots.paginatorlastpagelinkicon?{name:`lastpagelinkicon`,fn:a(function(t){return[n(e.$slots,`paginatorlastpagelinkicon`,{class:x(t.class)})]}),key:`6`}:void 0,e.$slots.paginatorjumptopagedropdownicon?{name:`jumptopagedropdownicon`,fn:a(function(t){return[n(e.$slots,`paginatorjumptopagedropdownicon`,{class:x(t.class)})]}),key:`7`}:void 0,e.$slots.paginatorrowsperpagedropdownicon?{name:`rowsperpagedropdownicon`,fn:a(function(t){return[n(e.$slots,`paginatorrowsperpagedropdownicon`,{class:x(t.class)})]}),key:`8`}:void 0]),1032,[`rows`,`first`,`totalRecords`,`pageLinkSize`,`template`,`rowsPerPageOptions`,`currentPageReportTemplate`,`class`,`alwaysShow`,`unstyled`,`pt`])):h(``,!0),b(`div`,u({class:e.cx(`tableContainer`),style:[e.sx(`tableContainer`),{maxHeight:e.scrollHeight}]},e.ptm(`tableContainer`)),[b(`table`,u({ref:`table`,role:`treegrid`,class:[e.cx(`table`),e.tableClass],style:e.tableStyle},sn(sn({},e.tableProps),e.ptm(`table`))),[b(`thead`,u({class:e.cx(`thead`),style:e.sx(`thead`),role:`rowgroup`},e.ptm(`thead`)),[b(`tr`,u({role:`row`},e.ptm(`headerRow`)),[(w(!0),S(y,null,r(m.columns,function(t,n){return w(),S(y,{key:m.columnProp(t,`columnKey`)||m.columnProp(t,`field`)||n},[m.columnProp(t,`hidden`)?h(``,!0):(w(),d(v,{key:0,column:t,resizableColumns:e.resizableColumns,sortField:l.d_sortField,sortOrder:l.d_sortOrder,multiSortMeta:l.d_multiSortMeta,sortMode:e.sortMode,onColumnClick:i[1]||=function(e){return m.onColumnHeaderClick(e)},onColumnResizestart:i[2]||=function(e){return m.onColumnResizeStart(e)},index:n,unstyled:e.unstyled,pt:e.pt},null,8,[`column`,`resizableColumns`,`sortField`,`sortOrder`,`multiSortMeta`,`sortMode`,`index`,`unstyled`,`pt`]))],64)}),128))],16),m.hasColumnFilter()?(w(),S(`tr`,f(u({key:0},e.ptm(`headerRow`))),[(w(!0),S(y,null,r(m.columns,function(n,r){return w(),S(y,{key:m.columnProp(n,`columnKey`)||m.columnProp(n,`field`)||r},[m.columnProp(n,`hidden`)?h(``,!0):(w(),S(`th`,u({key:0,class:m.getFilterColumnHeaderClass(n),style:[m.columnProp(n,`style`),m.columnProp(n,`filterHeaderStyle`)]},{ref_for:!0},e.ptm(`headerCell`,m.ptHeaderCellOptions(n))),[n.children&&n.children.filter?(w(),d(t(n.children.filter),{key:0,column:n,index:r},null,8,[`column`,`index`])):h(``,!0)],16))],64)}),128))],16)):h(``,!0)],16),b(`tbody`,u({class:e.cx(`tbody`),role:`rowgroup`},e.ptm(`tbody`)),[m.empty?(w(),S(`tr`,u({key:1,class:e.cx(`emptyMessage`)},e.ptm(`emptyMessage`)),[b(`td`,u({colspan:m.columns.length},e.ptm(`emptyMessageCell`)),[n(e.$slots,`empty`)],16,fn)],16)):(w(!0),S(y,{key:0},r(m.dataToRender,function(t,n){return w(),d(C,{key:m.nodeKey(t),dataKey:e.dataKey,columns:m.columns,node:t,level:0,expandedKeys:l.d_expandedKeys,indentation:e.indentation,selectionMode:e.selectionMode,selectionKeys:e.selectionKeys,ariaSetSize:m.dataToRender.length,ariaPosInset:n+1,tabindex:m.setTabindex(t,n),loadingMode:e.loadingMode,contextMenu:e.contextMenu,contextMenuSelection:e.contextMenuSelection,templates:e.$slots,onNodeToggle:m.onNodeToggle,onNodeClick:m.onNodeClick,onCheckboxChange:m.onCheckboxChange,onRowRightclick:i[3]||=function(e){return m.onRowRightClick(e)},unstyled:e.unstyled,pt:e.pt},null,8,[`dataKey`,`columns`,`node`,`expandedKeys`,`indentation`,`selectionMode`,`selectionKeys`,`ariaSetSize`,`ariaPosInset`,`tabindex`,`loadingMode`,`contextMenu`,`contextMenuSelection`,`templates`,`onNodeToggle`,`onNodeClick`,`onCheckboxChange`,`unstyled`,`pt`])}),128))],16),m.hasFooter?(w(),S(`tfoot`,u({key:0,class:e.cx(`tfoot`),style:e.sx(`tfoot`),role:`rowgroup`},e.ptm(`tfoot`)),[b(`tr`,u({role:`row`},e.ptm(`footerRow`)),[(w(!0),S(y,null,r(m.columns,function(t,n){return w(),S(y,{key:m.columnProp(t,`columnKey`)||m.columnProp(t,`field`)||n},[m.columnProp(t,`hidden`)?h(``,!0):(w(),d(T,{key:0,column:t,index:n,unstyled:e.unstyled,pt:e.pt},null,8,[`column`,`index`,`unstyled`,`pt`]))],64)}),128))],16)],16)):h(``,!0)],16)],16),m.paginatorBottom?(w(),d(_,{key:2,rows:l.d_rows,first:l.d_first,totalRecords:m.totalRecordsLength,pageLinkSize:e.pageLinkSize,template:e.paginatorTemplate,rowsPerPageOptions:e.rowsPerPageOptions,currentPageReportTemplate:e.currentPageReportTemplate,class:x(e.cx(`pcPaginator`,{position:`bottom`})),onPage:i[4]||=function(e){return m.onPage(e)},alwaysShow:e.alwaysShowPaginator,unstyled:e.unstyled,pt:e.ptm(`pcPaginator`)},p({_:2},[e.$slots.paginatorcontainer?{name:`container`,fn:a(function(t){return[n(e.$slots,`paginatorcontainer`,{first:t.first,last:t.last,rows:t.rows,page:t.page,pageCount:t.pageCount,pageLinks:t.pageLinks,totalRecords:t.totalRecords,firstPageCallback:t.firstPageCallback,lastPageCallback:t.lastPageCallback,prevPageCallback:t.prevPageCallback,nextPageCallback:t.nextPageCallback,rowChangeCallback:t.rowChangeCallback,changePageCallback:t.changePageCallback})]}),key:`0`}:void 0,e.$slots.paginatorstart?{name:`start`,fn:a(function(){return[n(e.$slots,`paginatorstart`)]}),key:`1`}:void 0,e.$slots.paginatorend?{name:`end`,fn:a(function(){return[n(e.$slots,`paginatorend`)]}),key:`2`}:void 0,e.$slots.paginatorfirstpagelinkicon?{name:`firstpagelinkicon`,fn:a(function(t){return[n(e.$slots,`paginatorfirstpagelinkicon`,{class:x(t.class)})]}),key:`3`}:void 0,e.$slots.paginatorprevpagelinkicon?{name:`prevpagelinkicon`,fn:a(function(t){return[n(e.$slots,`paginatorprevpagelinkicon`,{class:x(t.class)})]}),key:`4`}:void 0,e.$slots.paginatornextpagelinkicon?{name:`nextpagelinkicon`,fn:a(function(t){return[n(e.$slots,`paginatornextpagelinkicon`,{class:x(t.class)})]}),key:`5`}:void 0,e.$slots.paginatorlastpagelinkicon?{name:`lastpagelinkicon`,fn:a(function(t){return[n(e.$slots,`paginatorlastpagelinkicon`,{class:x(t.class)})]}),key:`6`}:void 0,e.$slots.paginatorjumptopagedropdownicon?{name:`jumptopagedropdownicon`,fn:a(function(t){return[n(e.$slots,`paginatorjumptopagedropdownicon`,{class:x(t.class)})]}),key:`7`}:void 0,e.$slots.paginatorrowsperpagedropdownicon?{name:`rowsperpagedropdownicon`,fn:a(function(t){return[n(e.$slots,`paginatorrowsperpagedropdownicon`,{class:x(t.class)})]}),key:`8`}:void 0]),1032,[`rows`,`first`,`totalRecords`,`pageLinkSize`,`template`,`rowsPerPageOptions`,`currentPageReportTemplate`,`class`,`alwaysShow`,`unstyled`,`pt`])):h(``,!0),e.$slots.footer?(w(),S(`div`,u({key:3,class:e.cx(`footer`)},e.ptm(`footer`)),[n(e.$slots,`footer`)],16)):h(``,!0),b(`div`,u({ref:`resizeHelper`,class:e.cx(`columnResizeIndicator`),style:{display:`none`}},e.ptm(`columnResizeIndicator`)),null,16)],16,dn)}rn.render=pn;var mn=F.extend({name:`skeleton`,style:`
    .p-skeleton {
        display: block;
        overflow: hidden;
        background: dt('skeleton.background');
        border-radius: dt('skeleton.border.radius');
    }

    .p-skeleton::after {
        content: '';
        animation: p-skeleton-animation 1.2s infinite;
        height: 100%;
        left: 0;
        position: absolute;
        right: 0;
        top: 0;
        transform: translateX(-100%);
        z-index: 1;
        background: linear-gradient(90deg, rgba(255, 255, 255, 0), dt('skeleton.animation.background'), rgba(255, 255, 255, 0));
    }

    [dir='rtl'] .p-skeleton::after {
        animation-name: p-skeleton-animation-rtl;
    }

    .p-skeleton-circle {
        border-radius: 50%;
    }

    .p-skeleton-animation-none::after {
        animation: none;
    }

    @keyframes p-skeleton-animation {
        from {
            transform: translateX(-100%);
        }
        to {
            transform: translateX(100%);
        }
    }

    @keyframes p-skeleton-animation-rtl {
        from {
            transform: translateX(100%);
        }
        to {
            transform: translateX(-100%);
        }
    }
`,classes:{root:function(e){var t=e.props;return[`p-skeleton p-component`,{"p-skeleton-circle":t.shape===`circle`,"p-skeleton-animation-none":t.animation===`none`}]}},inlineStyles:{root:{position:`relative`}}}),hn={name:`BaseSkeleton`,extends:V,props:{shape:{type:String,default:`rectangle`},size:{type:String,default:null},width:{type:String,default:`100%`},height:{type:String,default:`1rem`},borderRadius:{type:String,default:null},animation:{type:String,default:`wave`}},style:mn,provide:function(){return{$pcSkeleton:this,$parentInstance:this}}};function Q(e){"@babel/helpers - typeof";return Q=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},Q(e)}function gn(e,t,n){return(t=_n(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function _n(e){var t=vn(e,`string`);return Q(t)==`symbol`?t:t+``}function vn(e,t){if(Q(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(Q(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var $={name:`Skeleton`,extends:hn,inheritAttrs:!1,computed:{containerStyle:function(){return this.size?{width:this.size,height:this.size,borderRadius:this.borderRadius}:{width:this.width,height:this.height,borderRadius:this.borderRadius}},dataP:function(){return L(gn({},this.shape,this.shape))}}},yn=[`data-p`];function bn(e,t,n,r,i,a){return w(),S(`div`,u({class:e.cx(`root`),style:[e.sx(`root`),a.containerStyle],"aria-hidden":`true`},e.ptmi(`root`),{"data-p":a.dataP}),null,16,yn)}$.render=bn;var xn={class:`space-y-4`},Sn={class:`flex items-center justify-between`},Cn={class:`text-lg font-semibold text-gray-800 dark:text-gray-100`},wn={class:`text-sm text-gray-500 dark:text-gray-400`},Tn={class:`flex items-center gap-2`},En={class:`pt-2`},Dn={key:0,class:`space-y-2`},On={key:1,class:`flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500`},kn={class:`text-sm font-medium`},An={class:`text-sm mt-1 mb-4`},jn={class:`flex items-center gap-2`},Mn={class:`font-medium text-gray-800 dark:text-gray-100`},Nn={class:`text-gray-500 dark:text-gray-400 text-xs font-mono`},Pn={class:`text-gray-500 dark:text-gray-400`},Fn={class:`text-gray-500 dark:text-gray-400`},In={class:`flex items-center gap-1`},Ln={class:`pt-2`},Rn={key:0,class:`space-y-2`},zn={key:1,class:`flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500`},Bn={class:`text-sm font-medium`},Vn={class:`text-sm mt-1 mb-4`},Hn={class:`text-gray-800 dark:text-gray-100 font-medium`},Un={class:`text-gray-500 dark:text-gray-400`},Wn={class:`text-gray-500 dark:text-gray-400`},Gn={class:`flex items-center gap-1`},Kn={class:`space-y-4`},qn={key:0,class:`bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800 rounded-md px-3 py-2 text-sm text-emerald-700 dark:text-emerald-300`},Jn={class:`block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1`},Yn={key:0,class:`text-red-500 text-xs mt-1 block`},Xn={class:`block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1`},Zn={key:0,class:`text-red-500 text-xs mt-1 block`},Qn={key:1},$n={class:`block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1`},er={class:`block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1`},tr={class:`space-y-4`},nr={class:`block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1`},rr={key:0,class:`text-red-500 text-xs mt-1 block`},ir={class:`block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1`},ar={key:0,class:`text-red-500 text-xs mt-1 block`},or={class:`block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1`},sr={class:`flex items-center justify-between`},cr={class:`block text-sm font-medium text-gray-600 dark:text-gray-300`},lr={class:`block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1`},ur=ye({__name:`Organizations`,setup(t){let{t:n}=ge(),o=oe(),u=ue(),f=l(!1),p=l(!1),T=l([]),E=l(null),D=l(!1),O=l(!1),k=l(null),A=l({}),j=l([]),M=l(0),N=l({code:``,nomenclature:``,parent_id:null,sort_order:0}),P=l([]),ee=l(0),F=l(1),I=l(!1),L=l(!1),R=l(!1),te=l(null),z=l(!1),B=l({}),V=l({code:``,name:``,region:``,is_active:!0,sort_order:0});function ne(e){return e.map(e=>({key:e.id,data:e,children:e.children?ne(e.children):[]}))}async function H(){f.value=!0;try{let t=await e.get(`/api/v1/tenant/organizations?tree=true&per_page=200`),n=t.data?.data||t.data||[];T.value=ne(n);let r=(await e.get(`/api/v1/tenant/organizations?per_page=200`)).data?.data||[];j.value=r}catch(e){o.add({severity:`error`,summary:n(`message.error`),detail:e.response?.data?.error?.message||n(`message.failed_to_load`),life:4e3})}finally{f.value=!1}}let re=_(()=>{let e=[{id:null,label:n(`organization.no_parent`)}];return j.value.forEach(t=>{(!O.value||t.id!==k.value)&&e.push({id:t.id,label:`${t.full_code} — ${t.nomenclature}`})}),e}),ie=_(()=>{if(!N.value.parent_id)return n(`organization.no_parent`);let e=j.value.find(e=>e.id===N.value.parent_id);return e?`${e.full_code} — ${e.nomenclature}`:``});function ae(e){O.value=!1,k.value=null,A.value={},N.value={code:``,nomenclature:``,parent_id:e?.id||null,sort_order:0},D.value=!0}function se(e){O.value=!0,k.value=e.id,A.value={},N.value={code:e.code,nomenclature:e.nomenclature,parent_id:e.parent_id||null,sort_order:e.sort_order||0},D.value=!0}function ce(){N.value={code:``,nomenclature:``,parent_id:null,sort_order:0},A.value={},O.value=!1,k.value=null}async function le(){if(A.value={},!N.value.code?.trim()){A.value.code=[n(`form.required`)];return}if(!N.value.nomenclature?.trim()){A.value.nomenclature=[n(`form.required`)];return}p.value=!0;try{O.value?(await e.put(`/api/v1/tenant/organizations/${k.value}`,{code:N.value.code,nomenclature:N.value.nomenclature,sort_order:N.value.sort_order||0}),o.add({severity:`success`,summary:n(`message.success`),detail:n(`organization.updated`),life:3e3})):(await e.post(`/api/v1/tenant/organizations`,{code:N.value.code,nomenclature:N.value.nomenclature,parent_id:N.value.parent_id,sort_order:N.value.sort_order||0}),o.add({severity:`success`,summary:n(`message.success`),detail:n(`organization.created`),life:3e3})),D.value=!1,await H()}catch(e){let t=be(e);Object.keys(t).length>0?A.value=t:o.add({severity:`error`,summary:n(`message.error`),detail:e.response?.data?.error?.message||n(`message.operation_failed`),life:4e3})}finally{p.value=!1}}function U(t){u.require({header:n(`organization.confirm_delete_title`),message:n(`organization.confirm_delete`,{name:t.nomenclature}),icon:`pi pi-exclamation-triangle`,rejectLabel:n(`common.cancel`),acceptLabel:n(`common.delete`),rejectClass:`p-button-outlined p-button-secondary`,acceptClass:`p-button-danger`,accept:async()=>{try{await e.delete(`/api/v1/tenant/organizations/${t.id}`),o.add({severity:`success`,summary:n(`message.success`),detail:n(`organization.deleted`),life:3e3}),await H()}catch(e){o.add({severity:`error`,summary:n(`message.error`),detail:e.response?.data?.error?.message||n(`message.operation_failed`),life:4e3})}}})}async function G(){I.value=!0;try{let t=await e.get(`/api/v1/tenant/settings/zones?page=${F.value}&per_page=20`);P.value=t.data?.data?.data||t.data?.data||[],ee.value=t.data?.data?.total||t.data?.total||0,F.value=t.data?.data?.page||t.data?.page||1}catch(e){o.add({severity:`error`,summary:n(`message.error`),detail:e.response?.data?.error?.message||n(`message.failed_to_load`),life:4e3})}finally{I.value=!1}}function de(e){R.value=!!e,te.value=e?.id||null,B.value={},V.value={code:e?.code||``,name:e?.name||``,region:e?.region||``,is_active:e?.is_active===void 0||e.is_active,sort_order:e?.sort_order||0},L.value=!0}function pe(){V.value={code:``,name:``,region:``,is_active:!0,sort_order:0},B.value={},R.value=!1,te.value=null}async function he(){if(B.value={},!V.value.code?.trim()){B.value={code:[n(`form.required`)]};return}if(!V.value.name?.trim()){B.value={name:[n(`form.required`)]};return}z.value=!0;try{R.value?(await e.put(`/api/v1/tenant/settings/zones/${te.value}`,{code:V.value.code,name:V.value.name,region:V.value.region||void 0,is_active:V.value.is_active,sort_order:V.value.sort_order||0}),o.add({severity:`success`,summary:n(`message.success`),detail:n(`zones.updated`),life:3e3})):(await e.post(`/api/v1/tenant/settings/zones`,{code:V.value.code,name:V.value.name,region:V.value.region||void 0,is_active:V.value.is_active,sort_order:V.value.sort_order||0}),o.add({severity:`success`,summary:n(`message.success`),detail:n(`zones.created`),life:3e3})),L.value=!1,await G()}catch(e){let t=be(e);Object.keys(t).length>0?B.value=t:o.add({severity:`error`,summary:n(`message.error`),detail:e.response?.data?.error?.message||n(`message.operation_failed`),life:4e3})}finally{z.value=!1}}function _e(t){u.require({header:n(`zones.confirm_delete_title`),message:n(`zones.confirm_delete`,{name:t.name}),icon:`pi pi-exclamation-triangle`,rejectLabel:n(`common.cancel`),acceptLabel:n(`common.delete`),rejectClass:`p-button-outlined p-button-secondary`,acceptClass:`p-button-danger`,accept:async()=>{try{await e.delete(`/api/v1/tenant/settings/zones/${t.id}`),o.add({severity:`success`,summary:n(`message.success`),detail:n(`zones.deleted`),life:3e3}),await G()}catch(e){o.add({severity:`error`,summary:n(`message.error`),detail:e.response?.data?.error?.message||n(`message.operation_failed`),life:4e3})}}})}function ve(e){F.value=e.page+1,G()}return c(()=>{H(),G()}),(e,t)=>{let o=m(`tooltip`);return w(),S(`div`,xn,[b(`div`,Sn,[b(`div`,null,[b(`h1`,Cn,C(M.value===0?i(n)(`organization.title`):i(n)(`zones.title`)),1),b(`p`,wn,C(M.value===0?i(n)(`organization.description`):i(n)(`zones.description`)),1)]),b(`div`,Tn,[M.value===0?(w(),d(i(W),{key:0,label:i(n)(`organization.add_root`),icon:`pi pi-plus`,size:`small`,severity:`secondary`,outlined:``,onClick:t[0]||=e=>ae(null)},null,8,[`label`])):h(``,!0),M.value===1?(w(),d(i(W),{key:1,label:i(n)(`zones.new_zone`),icon:`pi pi-plus`,size:`small`,severity:`secondary`,outlined:``,onClick:t[1]||=e=>de()},null,8,[`label`])):h(``,!0),g(i(W),{label:i(n)(`common.refresh`),icon:`pi pi-refresh`,size:`small`,severity:`secondary`,text:``,onClick:t[2]||=e=>M.value===0?H():G(),loading:f.value},null,8,[`label`,`loading`])])]),g(i(Ve),{activeIndex:M.value,"onUpdate:activeIndex":t[6]||=e=>M.value=e,class:`!text-sm`},{default:a(()=>[g(i($e),{header:i(n)(`organization.tree_view`)},{default:a(()=>[b(`div`,En,[f.value?(w(),S(`div`,Dn,[(w(),S(y,null,r(5,e=>b(`div`,{key:e,class:`flex items-center gap-3 py-1`},[g(i($),{shape:`rectangle`,width:`1.25rem`,height:`1.25rem`,class:`!rounded`}),g(i($),{width:`8rem`,height:`1rem`}),g(i($),{width:`12rem`,height:`1rem`})])),64))])):T.value.length===0?(w(),S(`div`,On,[t[20]||=b(`i`,{class:`pi pi-sitemap text-4xl mb-3 opacity-50`},null,-1),b(`p`,kn,C(i(n)(`organization.empty_title`)),1),b(`p`,An,C(i(n)(`organization.empty_tree`)),1),g(i(W),{label:i(n)(`organization.add_root`),icon:`pi pi-plus`,size:`small`,onClick:t[3]||=e=>ae(null)},null,8,[`label`])])):(w(),d(i(rn),{key:2,value:T.value,class:`!text-sm !border-0`,scrollable:!0,scrollHeight:`flex`,stripedRows:``,selectionMode:`single`,selectionKeys:E.value,"onUpdate:selectionKeys":t[4]||=e=>E.value=e},{default:a(()=>[g(i(q),{field:`nomenclature`,header:i(n)(`organization.nomenclature`),expander:!0},{body:a(({node:e})=>[b(`div`,jn,[t[21]||=b(`i`,{class:`pi pi-folder-open text-amber-500 text-xs`},null,-1),b(`span`,Mn,C(e.data.nomenclature),1)])]),_:1},8,[`header`]),g(i(q),{field:`code`,header:i(n)(`organization.code`),style:{width:`120px`}},{body:a(({node:e})=>[g(i(xe),{value:e.data.code,severity:`info`,class:`!text-xs`},null,8,[`value`])]),_:1},8,[`header`]),g(i(q),{field:`full_code`,header:i(n)(`organization.full_code`),style:{width:`160px`}},{body:a(({node:e})=>[b(`span`,Nn,C(e.data.full_code),1)]),_:1},8,[`header`]),g(i(q),{field:`level`,header:i(n)(`organization.level`),style:{width:`80px`}},{body:a(({node:e})=>[b(`span`,Pn,C(e.data.level),1)]),_:1},8,[`header`]),g(i(q),{field:`sort_order`,header:i(n)(`organization.sort_order`),style:{width:`90px`}},{body:a(({node:e})=>[b(`span`,Fn,C(e.data.sort_order),1)]),_:1},8,[`header`]),g(i(q),{header:i(n)(`common.actions`),style:{width:`140px`},frozen:``,alignFrozen:`right`},{body:a(({node:e})=>[b(`div`,In,[s(g(i(W),{icon:`pi pi-plus`,severity:`secondary`,text:``,size:`small`,class:`!p-1`,onClick:t=>ae(e.data)},null,8,[`onClick`]),[[o,i(n)(`organization.add_child`),void 0,{top:!0}]]),s(g(i(W),{icon:`pi pi-pencil`,severity:`secondary`,text:``,size:`small`,class:`!p-1`,onClick:t=>se(e.data)},null,8,[`onClick`]),[[o,i(n)(`common.edit`),void 0,{top:!0}]]),s(g(i(W),{icon:`pi pi-trash`,severity:`danger`,text:``,size:`small`,class:`!p-1`,onClick:t=>U(e.data)},null,8,[`onClick`]),[[o,i(n)(`common.delete`),void 0,{top:!0}]])])]),_:1},8,[`header`])]),_:1},8,[`value`,`selectionKeys`]))])]),_:1},8,[`header`]),g(i($e),{header:i(n)(`zones.title`)},{default:a(()=>[b(`div`,Ln,[I.value?(w(),S(`div`,Rn,[(w(),S(y,null,r(4,e=>b(`div`,{key:e,class:`flex items-center gap-4 py-2`},[g(i($),{width:`5rem`,height:`1rem`}),g(i($),{width:`10rem`,height:`1rem`}),g(i($),{width:`6rem`,height:`1rem`}),g(i($),{width:`4rem`,height:`1.25rem`})])),64))])):P.value.length===0?(w(),S(`div`,zn,[t[22]||=b(`i`,{class:`pi pi-map-marker text-4xl mb-3 opacity-50`},null,-1),b(`p`,Bn,C(i(n)(`zones.empty_title`)),1),b(`p`,Vn,C(i(n)(`zones.description`)),1),g(i(W),{label:i(n)(`zones.new_zone`),icon:`pi pi-plus`,size:`small`,onClick:t[5]||=e=>de()},null,8,[`label`])])):(w(),d(i(Ee),{key:2,value:P.value,class:`!text-sm`,stripedRows:``,loading:I.value,paginator:``,rows:20,totalRecords:ee.value,lazy:!0,onPage:ve},{default:a(()=>[g(i(q),{field:`code`,header:i(n)(`zones.code`),style:{width:`120px`}},{body:a(({data:e})=>[g(i(xe),{value:e.code,severity:`info`,class:`!text-xs`},null,8,[`value`])]),_:1},8,[`header`]),g(i(q),{field:`name`,header:i(n)(`zones.name`)},{body:a(({data:e})=>[b(`span`,Hn,C(e.name),1)]),_:1},8,[`header`]),g(i(q),{field:`region`,header:i(n)(`zones.region`),style:{width:`150px`}},{body:a(({data:e})=>[b(`span`,Un,C(e.region||`—`),1)]),_:1},8,[`header`]),g(i(q),{field:`is_active`,header:i(n)(`zones.is_active`),style:{width:`100px`}},{body:a(({data:e})=>[g(i(xe),{value:e.is_active?i(n)(`common_status.active`):i(n)(`common_status.inactive`),severity:e.is_active?`success`:`warn`,class:`!text-xs`},null,8,[`value`,`severity`])]),_:1},8,[`header`]),g(i(q),{field:`sort_order`,header:i(n)(`zones.sort_order`),style:{width:`100px`}},{body:a(({data:e})=>[b(`span`,Wn,C(e.sort_order),1)]),_:1},8,[`header`]),g(i(q),{header:i(n)(`common.actions`),style:{width:`100px`},frozen:``,alignFrozen:`right`},{body:a(({data:e})=>[b(`div`,Gn,[s(g(i(W),{icon:`pi pi-pencil`,severity:`secondary`,text:``,size:`small`,class:`!p-1`,onClick:t=>de(e)},null,8,[`onClick`]),[[o,i(n)(`common.edit`),void 0,{top:!0}]]),s(g(i(W),{icon:`pi pi-trash`,severity:`danger`,text:``,size:`small`,class:`!p-1`,onClick:t=>_e(e)},null,8,[`onClick`]),[[o,i(n)(`common.delete`),void 0,{top:!0}]])])]),_:1},8,[`header`])]),_:1},8,[`value`,`loading`,`totalRecords`]))])]),_:1},8,[`header`])]),_:1},8,[`activeIndex`]),g(i(fe),{visible:D.value,"onUpdate:visible":t[12]||=e=>D.value=e,header:O.value?i(n)(`organization.edit`):i(n)(`organization.create`),modal:!0,closable:!0,class:`!w-full !max-w-lg`,onHide:ce},{footer:a(()=>[g(i(W),{label:i(n)(`common.cancel`),severity:`secondary`,outlined:``,size:`small`,onClick:t[11]||=e=>D.value=!1},null,8,[`label`]),g(i(W),{label:O.value?i(n)(`common.update`):i(n)(`common.save`),size:`small`,loading:p.value,disabled:p.value,onClick:le},null,8,[`label`,`loading`,`disabled`])]),default:a(()=>[b(`div`,Kn,[N.value.parent_id?(w(),S(`div`,qn,[t[23]||=b(`i`,{class:`pi pi-arrow-right mr-1`},null,-1),v(` `+C(i(n)(`organization.parent`))+`: `,1),b(`strong`,null,C(ie.value),1)])):h(``,!0),b(`div`,null,[b(`label`,Jn,[v(C(i(n)(`organization.code`))+` `,1),t[24]||=b(`span`,{class:`text-red-500`},`*`,-1)]),g(i(K),{modelValue:N.value.code,"onUpdate:modelValue":t[7]||=e=>N.value.code=e,class:x([`!w-full`,{"p-invalid":A.value?.code}]),maxlength:`10`,placeholder:i(n)(`organization.code`)},null,8,[`modelValue`,`class`,`placeholder`]),A.value?.code?(w(),S(`small`,Yn,C(A.value.code),1)):h(``,!0)]),b(`div`,null,[b(`label`,Xn,[v(C(i(n)(`organization.nomenclature`))+` `,1),t[25]||=b(`span`,{class:`text-red-500`},`*`,-1)]),g(i(K),{modelValue:N.value.nomenclature,"onUpdate:modelValue":t[8]||=e=>N.value.nomenclature=e,class:x([`!w-full`,{"p-invalid":A.value?.nomenclature}]),maxlength:`255`,placeholder:i(n)(`organization.nomenclature`)},null,8,[`modelValue`,`class`,`placeholder`]),A.value?.nomenclature?(w(),S(`small`,Zn,C(A.value.nomenclature),1)):h(``,!0)]),O.value?h(``,!0):(w(),S(`div`,Qn,[b(`label`,$n,C(i(n)(`organization.parent`)),1),g(i(Ae),{modelValue:N.value.parent_id,"onUpdate:modelValue":t[9]||=e=>N.value.parent_id=e,options:re.value,optionValue:`id`,optionLabel:`label`,placeholder:i(n)(`organization.select_parent`),class:`!w-full`,showClear:!0},null,8,[`modelValue`,`options`,`placeholder`])])),b(`div`,null,[b(`label`,er,C(i(n)(`organization.sort_order`)),1),g(i(Te),{modelValue:N.value.sort_order,"onUpdate:modelValue":t[10]||=e=>N.value.sort_order=e,class:`!w-full`,min:0},null,8,[`modelValue`])])])]),_:1},8,[`visible`,`header`]),g(i(fe),{visible:L.value,"onUpdate:visible":t[19]||=e=>L.value=e,header:R.value?i(n)(`zones.edit_zone`):i(n)(`zones.new_zone`),modal:!0,closable:!0,class:`!w-full !max-w-md`,onHide:pe},{footer:a(()=>[g(i(W),{label:i(n)(`common.cancel`),severity:`secondary`,outlined:``,size:`small`,onClick:t[18]||=e=>L.value=!1},null,8,[`label`]),g(i(W),{label:R.value?i(n)(`common.update`):i(n)(`common.save`),size:`small`,loading:z.value,disabled:z.value,onClick:he},null,8,[`label`,`loading`,`disabled`])]),default:a(()=>[b(`div`,tr,[b(`div`,null,[b(`label`,nr,[v(C(i(n)(`zones.code`))+` `,1),t[26]||=b(`span`,{class:`text-red-500`},`*`,-1)]),g(i(K),{modelValue:V.value.code,"onUpdate:modelValue":t[13]||=e=>V.value.code=e,class:x([`!w-full`,{"p-invalid":B.value?.code}]),maxlength:`20`,placeholder:i(n)(`zones.code`)},null,8,[`modelValue`,`class`,`placeholder`]),B.value?.code?(w(),S(`small`,rr,C(B.value.code),1)):h(``,!0)]),b(`div`,null,[b(`label`,ir,[v(C(i(n)(`zones.name`))+` `,1),t[27]||=b(`span`,{class:`text-red-500`},`*`,-1)]),g(i(K),{modelValue:V.value.name,"onUpdate:modelValue":t[14]||=e=>V.value.name=e,class:x([`!w-full`,{"p-invalid":B.value?.name}]),maxlength:`255`,placeholder:i(n)(`zones.name`)},null,8,[`modelValue`,`class`,`placeholder`]),B.value?.name?(w(),S(`small`,ar,C(B.value.name),1)):h(``,!0)]),b(`div`,null,[b(`label`,or,C(i(n)(`zones.region`)),1),g(i(K),{modelValue:V.value.region,"onUpdate:modelValue":t[15]||=e=>V.value.region=e,class:`!w-full`,maxlength:`100`,placeholder:i(n)(`zones.region`)},null,8,[`modelValue`,`placeholder`])]),b(`div`,sr,[b(`label`,cr,C(i(n)(`zones.is_active`)),1),g(i(je),{modelValue:V.value.is_active,"onUpdate:modelValue":t[16]||=e=>V.value.is_active=e},null,8,[`modelValue`])]),b(`div`,null,[b(`label`,lr,C(i(n)(`zones.sort_order`)),1),g(i(Te),{modelValue:V.value.sort_order,"onUpdate:modelValue":t[17]||=e=>V.value.sort_order=e,class:`!w-full`,min:0},null,8,[`modelValue`])])])]),_:1},8,[`visible`,`header`]),g(i(me))])}}},[[`__scopeId`,`data-v-92a66d3e`]]);export{ur as default};