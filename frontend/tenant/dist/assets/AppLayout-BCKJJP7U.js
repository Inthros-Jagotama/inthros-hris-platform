import{$ as e,A as t,B as n,C as r,E as i,F as a,G as o,H as s,I as c,J as l,K as u,L as d,M as f,N as p,O as m,P as h,Q as g,S as _,T as v,U as y,V as b,W as x,X as S,Y as C,Z as w,_ as T,a as E,b as D,c as O,d as k,f as ee,g as te,h as A,i as j,j as ne,k as M,l as re,m as ie,n as ae,o as oe,p as N,q as P,r as F,s as se,t as I,u as L,v as R,w as z,x as ce,y as B,z as le}from"./index-BNiZ82cV.js";import{t as ue}from"./_plugin-vue_export-helper-BDNMzG2s.js";import{t as de}from"./baseeditableholder-lKcXaRhN.js";var V={name:`ChevronDownIcon`,extends:j};function fe(e){return ge(e)||he(e)||me(e)||pe()}function pe(){throw TypeError(`Invalid attempt to spread non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}function me(e,t){if(e){if(typeof e==`string`)return H(e,t);var n={}.toString.call(e).slice(8,-1);return n===`Object`&&e.constructor&&(n=e.constructor.name),n===`Map`||n===`Set`?Array.from(e):n===`Arguments`||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?H(e,t):void 0}}function he(e){if(typeof Symbol<`u`&&e[Symbol.iterator]!=null||e[`@@iterator`]!=null)return Array.from(e)}function ge(e){if(Array.isArray(e))return H(e)}function H(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=Array(t);n<t;n++)r[n]=e[n];return r}function _e(e,t,n,r,i,a){return s(),d(`svg`,b({width:`14`,height:`14`,viewBox:`0 0 14 14`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},e.pti()),fe(t[0]||=[h(`path`,{d:`M7.01744 10.398C6.91269 10.3985 6.8089 10.378 6.71215 10.3379C6.61541 10.2977 6.52766 10.2386 6.45405 10.1641L1.13907 4.84913C1.03306 4.69404 0.985221 4.5065 1.00399 4.31958C1.02276 4.13266 1.10693 3.95838 1.24166 3.82747C1.37639 3.69655 1.55301 3.61742 1.74039 3.60402C1.92777 3.59062 2.11386 3.64382 2.26584 3.75424L7.01744 8.47394L11.769 3.75424C11.9189 3.65709 12.097 3.61306 12.2748 3.62921C12.4527 3.64535 12.6199 3.72073 12.7498 3.84328C12.8797 3.96582 12.9647 4.12842 12.9912 4.30502C13.0177 4.48162 12.9841 4.662 12.8958 4.81724L7.58083 10.1322C7.50996 10.2125 7.42344 10.2775 7.32656 10.3232C7.22968 10.3689 7.12449 10.3944 7.01744 10.398Z`,fill:`currentColor`},null,-1)]),16)}V.render=_e;var U={name:`ChevronRightIcon`,extends:j};function ve(e){return Se(e)||xe(e)||be(e)||ye()}function ye(){throw TypeError(`Invalid attempt to spread non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}function be(e,t){if(e){if(typeof e==`string`)return W(e,t);var n={}.toString.call(e).slice(8,-1);return n===`Object`&&e.constructor&&(n=e.constructor.name),n===`Map`||n===`Set`?Array.from(e):n===`Arguments`||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?W(e,t):void 0}}function xe(e){if(typeof Symbol<`u`&&e[Symbol.iterator]!=null||e[`@@iterator`]!=null)return Array.from(e)}function Se(e){if(Array.isArray(e))return W(e)}function W(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=Array(t);n<t;n++)r[n]=e[n];return r}function Ce(e,t,n,r,i,a){return s(),d(`svg`,b({width:`14`,height:`14`,viewBox:`0 0 14 14`,fill:`none`,xmlns:`http://www.w3.org/2000/svg`},e.pti()),ve(t[0]||=[h(`path`,{d:`M4.38708 13C4.28408 13.0005 4.18203 12.9804 4.08691 12.9409C3.99178 12.9014 3.9055 12.8433 3.83313 12.7701C3.68634 12.6231 3.60388 12.4238 3.60388 12.2161C3.60388 12.0084 3.68634 11.8091 3.83313 11.6622L8.50507 6.99022L3.83313 2.31827C3.69467 2.16968 3.61928 1.97313 3.62287 1.77005C3.62645 1.56698 3.70872 1.37322 3.85234 1.22959C3.99596 1.08597 4.18972 1.00371 4.3928 1.00012C4.59588 0.996539 4.79242 1.07192 4.94102 1.21039L10.1669 6.43628C10.3137 6.58325 10.3962 6.78249 10.3962 6.99022C10.3962 7.19795 10.3137 7.39718 10.1669 7.54416L4.94102 12.7701C4.86865 12.8433 4.78237 12.9014 4.68724 12.9409C4.59212 12.9804 4.49007 13.0005 4.38708 13Z`,fill:`currentColor`},null,-1)]),16)}U.render=Ce;var we=L.extend({name:`panelmenu`,style:`
    .p-panelmenu {
        display: flex;
        flex-direction: column;
        gap: dt('panelmenu.gap');
    }

    .p-panelmenu-panel {
        background: dt('panelmenu.panel.background');
        border-width: dt('panelmenu.panel.border.width');
        border-style: solid;
        border-color: dt('panelmenu.panel.border.color');
        color: dt('panelmenu.panel.color');
        border-radius: dt('panelmenu.panel.border.radius');
        padding: dt('panelmenu.panel.padding');
    }

    .p-panelmenu-panel:first-child {
        border-width: dt('panelmenu.panel.first.border.width');
        border-start-start-radius: dt('panelmenu.panel.first.top.border.radius');
        border-start-end-radius: dt('panelmenu.panel.first.top.border.radius');
    }

    .p-panelmenu-panel:last-child {
        border-width: dt('panelmenu.panel.last.border.width');
        border-end-start-radius: dt('panelmenu.panel.last.bottom.border.radius');
        border-end-end-radius: dt('panelmenu.panel.last.bottom.border.radius');
    }

    .p-panelmenu-header {
        outline: 0 none;
    }

    .p-panelmenu-header-content {
        border-radius: dt('panelmenu.item.border.radius');
        transition:
            background dt('panelmenu.transition.duration'),
            color dt('panelmenu.transition.duration'),
            outline-color dt('panelmenu.transition.duration'),
            box-shadow dt('panelmenu.transition.duration');
        outline-color: transparent;
        color: dt('panelmenu.item.color');
    }

    .p-panelmenu-header-link {
        display: flex;
        gap: dt('panelmenu.item.gap');
        padding: dt('panelmenu.item.padding');
        align-items: center;
        user-select: none;
        cursor: pointer;
        position: relative;
        text-decoration: none;
        color: inherit;
    }

    .p-panelmenu-header-icon,
    .p-panelmenu-item-icon {
        color: dt('panelmenu.item.icon.color');
    }

    .p-panelmenu-submenu-icon {
        color: dt('panelmenu.submenu.icon.color');
    }

    .p-panelmenu-submenu-icon:dir(rtl) {
        transform: rotate(180deg);
    }

    .p-panelmenu-header:not(.p-disabled):focus-visible .p-panelmenu-header-content {
        background: dt('panelmenu.item.focus.background');
        color: dt('panelmenu.item.focus.color');
    }

    .p-panelmenu-header:not(.p-disabled):focus-visible .p-panelmenu-header-content .p-panelmenu-header-icon {
        color: dt('panelmenu.item.icon.focus.color');
    }

    .p-panelmenu-header:not(.p-disabled):focus-visible .p-panelmenu-header-content .p-panelmenu-submenu-icon {
        color: dt('panelmenu.submenu.icon.focus.color');
    }

    .p-panelmenu-header:not(.p-disabled) .p-panelmenu-header-content:hover {
        background: dt('panelmenu.item.focus.background');
        color: dt('panelmenu.item.focus.color');
    }

    .p-panelmenu-header:not(.p-disabled) .p-panelmenu-header-content:hover .p-panelmenu-header-icon {
        color: dt('panelmenu.item.icon.focus.color');
    }

    .p-panelmenu-header:not(.p-disabled) .p-panelmenu-header-content:hover .p-panelmenu-submenu-icon {
        color: dt('panelmenu.submenu.icon.focus.color');
    }

    .p-panelmenu-submenu {
        margin: 0;
        padding: 0 0 0 dt('panelmenu.submenu.indent');
        outline: 0;
        list-style: none;
    }

    .p-panelmenu-submenu:dir(rtl) {
        padding: 0 dt('panelmenu.submenu.indent') 0 0;
    }

    .p-panelmenu-item-link {
        display: flex;
        gap: dt('panelmenu.item.gap');
        padding: dt('panelmenu.item.padding');
        align-items: center;
        user-select: none;
        cursor: pointer;
        text-decoration: none;
        color: inherit;
        position: relative;
        overflow: hidden;
    }

    .p-panelmenu-item-label {
        line-height: 1;
    }

    .p-panelmenu-item-content {
        border-radius: dt('panelmenu.item.border.radius');
        transition:
            background dt('panelmenu.transition.duration'),
            color dt('panelmenu.transition.duration'),
            outline-color dt('panelmenu.transition.duration'),
            box-shadow dt('panelmenu.transition.duration');
        color: dt('panelmenu.item.color');
        outline-color: transparent;
    }

    .p-panelmenu-item.p-focus > .p-panelmenu-item-content {
        background: dt('panelmenu.item.focus.background');
        color: dt('panelmenu.item.focus.color');
    }

    .p-panelmenu-item.p-focus > .p-panelmenu-item-content .p-panelmenu-item-icon {
        color: dt('panelmenu.item.focus.color');
    }

    .p-panelmenu-item.p-focus > .p-panelmenu-item-content .p-panelmenu-submenu-icon {
        color: dt('panelmenu.submenu.icon.focus.color');
    }

    .p-panelmenu-item:not(.p-disabled) > .p-panelmenu-item-content:hover {
        background: dt('panelmenu.item.focus.background');
        color: dt('panelmenu.item.focus.color');
    }

    .p-panelmenu-item:not(.p-disabled) > .p-panelmenu-item-content:hover .p-panelmenu-item-icon {
        color: dt('panelmenu.item.icon.focus.color');
    }

    .p-panelmenu-item:not(.p-disabled) > .p-panelmenu-item-content:hover .p-panelmenu-submenu-icon {
        color: dt('panelmenu.submenu.icon.focus.color');
    }

    .p-panelmenu-content-container {
        display: grid;
        grid-template-rows: 1fr;
    }

    .p-panelmenu-content-wrapper {
        min-height: 0;
    }
`,classes:{root:`p-panelmenu p-component`,panel:`p-panelmenu-panel`,header:function(e){var t=e.instance,n=e.item;return[`p-panelmenu-header`,{"p-panelmenu-header-active":t.isItemActive(n)&&!!n.items,"p-disabled":t.isItemDisabled(n)}]},headerContent:`p-panelmenu-header-content`,headerLink:`p-panelmenu-header-link`,headerIcon:`p-panelmenu-header-icon`,headerLabel:`p-panelmenu-header-label`,contentContainer:`p-panelmenu-content-container`,contentWrapper:`p-panelmenu-content-wrapper`,content:`p-panelmenu-content`,rootList:`p-panelmenu-root-list`,item:function(e){var t=e.instance,n=e.processedItem;return[`p-panelmenu-item`,{"p-focus":t.isItemFocused(n),"p-disabled":t.isItemDisabled(n)}]},itemContent:`p-panelmenu-item-content`,itemLink:`p-panelmenu-item-link`,itemIcon:`p-panelmenu-item-icon`,itemLabel:`p-panelmenu-item-label`,submenuIcon:`p-panelmenu-submenu-icon`,submenu:`p-panelmenu-submenu`,separator:`p-menuitem-separator`}}),Te={name:`BasePanelMenu`,extends:E,props:{model:{type:Array,default:null},expandedKeys:{type:Object,default:null},multiple:{type:Boolean,default:!1},tabindex:{type:Number,default:0}},style:we,provide:function(){return{$pcPanelMenu:this,$parentInstance:this}}},G={name:`PanelMenuSub`,hostName:`PanelMenu`,extends:E,emits:[`item-toggle`,`item-mousemove`],props:{panelId:{type:String,default:null},focusedItemId:{type:String,default:null},items:{type:Array,default:null},level:{type:Number,default:0},templates:{type:Object,default:null},activeItemPath:{type:Object,default:null},tabindex:{type:Number,default:-1}},methods:{getItemId:function(e){return`${this.panelId}_${e.key}`},getItemKey:function(e){return this.getItemId(e)},getItemProp:function(e,t,n){return e&&e.item?i(e.item[t],n):void 0},getItemLabel:function(e){return this.getItemProp(e,`label`)},getPTOptions:function(e,t,n){return this.ptm(e,{context:{item:t.item,index:n,active:this.isItemActive(t),focused:this.isItemFocused(t),disabled:this.isItemDisabled(t)}})},isItemActive:function(e){return this.activeItemPath.some(function(t){return t.key===e.key})},isItemVisible:function(e){return this.getItemProp(e,`visible`)!==!1},isItemDisabled:function(e){return this.getItemProp(e,`disabled`)},isItemFocused:function(e){return this.focusedItemId===this.getItemId(e)},isItemGroup:function(e){return m(e.items)},onItemClick:function(e,t){this.getItemProp(t,`command`,{originalEvent:e,item:t.item}),this.$emit(`item-toggle`,{processedItem:t,expanded:!this.isItemActive(t)})},onItemToggle:function(e){this.$emit(`item-toggle`,e)},onItemMouseMove:function(e,t){this.$emit(`item-mousemove`,{originalEvent:e,processedItem:t})},getAriaSetSize:function(){var e=this;return this.items.filter(function(t){return e.isItemVisible(t)&&!e.getItemProp(t,`separator`)}).length},getAriaPosInset:function(e){var t=this;return e-this.items.slice(0,e).filter(function(e){return t.isItemVisible(e)&&t.getItemProp(e,`separator`)}).length+1},getMenuItemProps:function(e,t){return{action:b({class:this.cx(`itemLink`),tabindex:-1},this.getPTOptions(`itemLink`,e,t)),icon:b({class:[this.cx(`itemIcon`),this.getItemProp(e,`icon`)]},this.getPTOptions(`itemIcon`,e,t)),label:b({class:this.cx(`itemLabel`)},this.getPTOptions(`itemLabel`,e,t)),submenuicon:b({class:this.cx(`submenuIcon`)},this.getPTOptions(`submenuicon`,e,t))}}},components:{ChevronRightIcon:U,ChevronDownIcon:V},directives:{ripple:F}},Ee=[`tabindex`],De=[`id`,`aria-label`,`aria-expanded`,`aria-level`,`aria-setsize`,`aria-posinset`,`data-p-focused`,`data-p-disabled`],Oe=[`onClick`,`onMousemove`],ke=[`href`,`target`];function Ae(r,i,p,m,_,v){var x=o(`PanelMenuSub`,!0),S=u(`ripple`);return s(),d(`ul`,{class:g(r.cx(`submenu`)),tabindex:p.tabindex},[(s(!0),d(f,null,y(p.items,function(o,u){return s(),d(f,{key:v.getItemKey(o)},[v.isItemVisible(o)&&!v.getItemProp(o,`separator`)?(s(),d(`li`,b({key:0,id:v.getItemId(o),class:[r.cx(`item`,{processedItem:o}),v.getItemProp(o,`class`)],style:v.getItemProp(o,`style`),role:`treeitem`,"aria-label":v.getItemLabel(o),"aria-expanded":v.isItemGroup(o)?v.isItemActive(o):void 0,"aria-level":p.level+1,"aria-setsize":v.getAriaSetSize(),"aria-posinset":v.getAriaPosInset(u)},{ref_for:!0},v.getPTOptions(`item`,o,u),{"data-p-focused":v.isItemFocused(o),"data-p-disabled":v.isItemDisabled(o)}),[h(`div`,b({class:r.cx(`itemContent`),onClick:function(e){return v.onItemClick(e,o)},onMousemove:function(e){return v.onItemMouseMove(e,o)}},{ref_for:!0},v.getPTOptions(`itemContent`,o,u)),[p.templates.item?(s(),a(P(p.templates.item),{key:1,item:o.item,root:!1,active:v.isItemActive(o),hasSubmenu:v.isItemGroup(o),label:v.getItemLabel(o),props:v.getMenuItemProps(o,u)},null,8,[`item`,`active`,`hasSubmenu`,`label`,`props`])):C((s(),d(`a`,b({key:0,href:v.getItemProp(o,`url`),class:r.cx(`itemLink`),target:v.getItemProp(o,`target`),tabindex:`-1`},{ref_for:!0},v.getPTOptions(`itemLink`,o,u)),[v.isItemGroup(o)?(s(),d(f,{key:0},[p.templates.submenuicon?(s(),a(P(p.templates.submenuicon),b({key:0,class:r.cx(`submenuIcon`),active:v.isItemActive(o)},{ref_for:!0},v.getPTOptions(`submenuIcon`,o,u)),null,16,[`class`,`active`])):(s(),a(P(v.isItemActive(o)?`ChevronDownIcon`:`ChevronRightIcon`),b({key:1,class:r.cx(`submenuIcon`)},{ref_for:!0},v.getPTOptions(`submenuIcon`,o,u)),null,16,[`class`]))],64)):c(``,!0),p.templates.itemicon?(s(),a(P(p.templates.itemicon),{key:1,item:o.item,class:g(r.cx(`itemIcon`))},null,8,[`item`,`class`])):v.getItemProp(o,`icon`)?(s(),d(`span`,b({key:2,class:[r.cx(`itemIcon`),v.getItemProp(o,`icon`)]},{ref_for:!0},v.getPTOptions(`itemIcon`,o,u)),null,16)):c(``,!0),h(`span`,b({class:r.cx(`itemLabel`)},{ref_for:!0},v.getPTOptions(`itemLabel`,o,u)),e(v.getItemLabel(o)),17)],16,ke)),[[S]])],16,Oe),n(M,b({name:`p-collapsible`},{ref_for:!0},r.ptm(`transition`)),{default:l(function(){return[C(h(`div`,b({class:r.cx(`contentContainer`)},{ref_for:!0},r.ptm(`contentContainer`)),[h(`div`,b({class:r.cx(`contentWrapper`)},{ref_for:!0},r.ptm(`contentWrapper`)),[v.isItemVisible(o)&&v.isItemGroup(o)?(s(),a(x,b({key:0,id:v.getItemId(o)+`_list`,role:`group`,panelId:p.panelId,focusedItemId:p.focusedItemId,items:o.items,level:p.level+1,templates:p.templates,activeItemPath:p.activeItemPath,onItemToggle:v.onItemToggle,onItemMousemove:i[0]||=function(e){return r.$emit(`item-mousemove`,e)},pt:r.pt,unstyled:r.unstyled},{ref_for:!0},r.ptm(`submenu`)),null,16,[`id`,`panelId`,`focusedItemId`,`items`,`level`,`templates`,`activeItemPath`,`onItemToggle`,`pt`,`unstyled`])):c(``,!0)],16)],16),[[t,v.isItemActive(o)]])]}),_:2},1040)],16,De)):c(``,!0),v.isItemVisible(o)&&v.getItemProp(o,`separator`)?(s(),d(`li`,b({key:1,style:v.getItemProp(o,`style`),class:[r.cx(`separator`),v.getItemProp(o,`class`)],role:`separator`},{ref_for:!0},r.ptm(`separator`)),null,16)):c(``,!0)],64)}),128))],10,Ee)}G.render=Ae;function je(e,t){return Fe(e)||Pe(e,t)||Ne(e,t)||Me()}function Me(){throw TypeError(`Invalid attempt to destructure non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}function Ne(e,t){if(e){if(typeof e==`string`)return K(e,t);var n={}.toString.call(e).slice(8,-1);return n===`Object`&&e.constructor&&(n=e.constructor.name),n===`Map`||n===`Set`?Array.from(e):n===`Arguments`||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?K(e,t):void 0}}function K(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=Array(t);n<t;n++)r[n]=e[n];return r}function Pe(e,t){var n=e==null?null:typeof Symbol<`u`&&e[Symbol.iterator]||e[`@@iterator`];if(n!=null){var r,i,a,o,s=[],c=!0,l=!1;try{if(a=(n=n.call(e)).next,t!==0)for(;!(c=(r=a.call(n)).done)&&(s.push(r.value),s.length!==t);c=!0);}catch(e){l=!0,i=e}finally{try{if(!c&&n.return!=null&&(o=n.return(),Object(o)!==o))return}finally{if(l)throw i}}return s}}function Fe(e){if(Array.isArray(e))return e}var q={name:`PanelMenuList`,hostName:`PanelMenu`,extends:E,emits:[`item-toggle`,`header-focus`],props:{panelId:{type:String,default:null},items:{type:Array,default:null},templates:{type:Object,default:null},expandedKeys:{type:Object,default:null}},searchTimeout:null,searchValue:null,data:function(){return{focused:!1,focusedItem:null,activeItemPath:[]}},watch:{expandedKeys:function(e){this.autoUpdateActiveItemPath(e)}},created:function(){this.autoUpdateActiveItemPath(this.expandedKeys)},methods:{getItemProp:function(e,t){return e&&e.item?i(e.item[t]):void 0},getItemLabel:function(e){return this.getItemProp(e,`label`)},isItemVisible:function(e){return this.getItemProp(e,`visible`)!==!1},isItemDisabled:function(e){return this.getItemProp(e,`disabled`)},isItemActive:function(e){return this.activeItemPath.some(function(t){return t.key===e.parentKey})},isItemGroup:function(e){return m(e.items)},onFocus:function(e){this.focused=!0,this.focusedItem=this.focusedItem||(this.isElementInPanel(e,e.relatedTarget)?this.findFirstItem():this.findLastItem())},onBlur:function(){this.focused=!1,this.focusedItem=null,this.searchValue=``},onKeyDown:function(e){var t=e.metaKey||e.ctrlKey;switch(e.code){case`ArrowDown`:this.onArrowDownKey(e);break;case`ArrowUp`:this.onArrowUpKey(e);break;case`ArrowLeft`:this.onArrowLeftKey(e);break;case`ArrowRight`:this.onArrowRightKey(e);break;case`Home`:this.onHomeKey(e);break;case`End`:this.onEndKey(e);break;case`Space`:this.onSpaceKey(e);break;case`Enter`:case`NumpadEnter`:this.onEnterKey(e);break;case`Escape`:case`Tab`:case`PageDown`:case`PageUp`:case`Backspace`:case`ShiftLeft`:case`ShiftRight`:break;default:!t&&_(e.key)&&this.searchItems(e,e.key);break}},onArrowDownKey:function(e){var t=m(this.focusedItem)?this.findNextItem(this.focusedItem):this.findFirstItem();this.changeFocusedItem({originalEvent:e,processedItem:t,focusOnNext:!0}),e.preventDefault()},onArrowUpKey:function(e){var t=m(this.focusedItem)?this.findPrevItem(this.focusedItem):this.findLastItem();this.changeFocusedItem({originalEvent:e,processedItem:t,selfCheck:!0}),e.preventDefault()},onArrowLeftKey:function(e){var t=this;m(this.focusedItem)&&(this.activeItemPath.some(function(e){return e.key===t.focusedItem.key})?this.activeItemPath=this.activeItemPath.filter(function(e){return e.key!==t.focusedItem.key}):this.focusedItem=m(this.focusedItem.parent)?this.focusedItem.parent:this.focusedItem,e.preventDefault())},onArrowRightKey:function(e){var t=this;m(this.focusedItem)&&(this.isItemGroup(this.focusedItem)&&(this.activeItemPath.some(function(e){return e.key===t.focusedItem.key})?this.onArrowDownKey(e):(this.activeItemPath=this.activeItemPath.filter(function(e){return e.parentKey!==t.focusedItem.parentKey}),this.activeItemPath.push(this.focusedItem))),e.preventDefault())},onHomeKey:function(e){this.changeFocusedItem({originalEvent:e,processedItem:this.findFirstItem(),allowHeaderFocus:!1}),e.preventDefault()},onEndKey:function(e){this.changeFocusedItem({originalEvent:e,processedItem:this.findLastItem(),focusOnNext:!0,allowHeaderFocus:!1}),e.preventDefault()},onEnterKey:function(e){if(m(this.focusedItem)){var t=B(this.$el,`li[id="${`${this.focusedItemId}`}"]`),n=t&&(B(t,`[data-pc-section="itemlink"]`)||B(t,`a,button`));n?n.click():t&&t.click()}e.preventDefault()},onSpaceKey:function(e){this.onEnterKey(e)},onItemToggle:function(e){var t=e.processedItem,n=e.expanded;this.expandedKeys?this.$emit(`item-toggle`,{item:t.item,expanded:n}):(this.activeItemPath=this.activeItemPath.filter(function(e){return e.parentKey!==t.parentKey}),n&&this.activeItemPath.push(t)),this.focusedItem=t,T(this.$el)},onItemMouseMove:function(e){this.focused&&(this.focusedItem=e.processedItem)},isElementInPanel:function(e,t){var n=e.currentTarget.closest(`[data-pc-section="panel"]`);return n&&n.contains(t)},isItemMatched:function(e){return this.isValidItem(e)&&this.getItemLabel(e)?.toLocaleLowerCase(this.searchLocale).startsWith(this.searchValue.toLocaleLowerCase(this.searchLocale))},isVisibleItem:function(e){return!!e&&(e.level===0||this.isItemActive(e))&&this.isItemVisible(e)},isValidItem:function(e){return!!e&&!this.isItemDisabled(e)&&!this.getItemProp(e,`separator`)},findFirstItem:function(){var e=this;return this.visibleItems.find(function(t){return e.isValidItem(t)})},findLastItem:function(){var e=this;return r(this.visibleItems,function(t){return e.isValidItem(t)})},findNextItem:function(e){var t=this,n=this.visibleItems.findIndex(function(t){return t.key===e.key});return(n<this.visibleItems.length-1?this.visibleItems.slice(n+1).find(function(e){return t.isValidItem(e)}):void 0)||e},findPrevItem:function(e){var t=this,n=this.visibleItems.findIndex(function(t){return t.key===e.key});return(n>0?r(this.visibleItems.slice(0,n),function(e){return t.isValidItem(e)}):void 0)||e},searchItems:function(e,t){var n=this;this.searchValue=(this.searchValue||``)+t;var r=null,i=!1;if(m(this.focusedItem)){var a=this.visibleItems.findIndex(function(e){return e.key===n.focusedItem.key});r=this.visibleItems.slice(a).find(function(e){return n.isItemMatched(e)}),r=v(r)?this.visibleItems.slice(0,a).find(function(e){return n.isItemMatched(e)}):r}else r=this.visibleItems.find(function(e){return n.isItemMatched(e)});return m(r)&&(i=!0),v(r)&&v(this.focusedItem)&&(r=this.findFirstItem()),m(r)&&this.changeFocusedItem({originalEvent:e,processedItem:r,allowHeaderFocus:!1}),this.searchTimeout&&clearTimeout(this.searchTimeout),this.searchTimeout=setTimeout(function(){n.searchValue=``,n.searchTimeout=null},500),i},changeFocusedItem:function(e){var t=e.originalEvent,n=e.processedItem,r=e.focusOnNext,i=e.selfCheck,a=e.allowHeaderFocus,o=a===void 0||a;m(this.focusedItem)&&this.focusedItem.key!==n.key?(this.focusedItem=n,this.scrollInView()):o&&this.$emit(`header-focus`,{originalEvent:t,focusOnNext:r,selfCheck:i})},scrollInView:function(){var e=B(this.$el,`li[id="${`${this.focusedItemId}`}"]`);e&&e.scrollIntoView&&e.scrollIntoView({block:`nearest`,inline:`start`})},autoUpdateActiveItemPath:function(e){var t=this;this.activeItemPath=Object.entries(e||{}).reduce(function(e,n){var r=je(n,2),i=r[0];if(r[1]){var a=t.findProcessedItemByItemKey(i);a&&e.push(a)}return e},[])},findProcessedItemByItemKey:function(e,t){var n=arguments.length>2&&arguments[2]!==void 0?arguments[2]:0;if(t||=n===0&&this.processedItems,!t)return null;for(var r=0;r<t.length;r++){var i=t[r];if(this.getItemProp(i,`key`)===e)return i;var a=this.findProcessedItemByItemKey(e,i.items,n+1);if(a)return a}},createProcessedItems:function(e){var t=this,n=arguments.length>1&&arguments[1]!==void 0?arguments[1]:0,r=arguments.length>2&&arguments[2]!==void 0?arguments[2]:{},i=arguments.length>3&&arguments[3]!==void 0?arguments[3]:``,a=[];return e&&e.forEach(function(e,o){var s=(i===``?``:i+`_`)+o,c={item:e,index:o,level:n,key:s,parent:r,parentKey:i};c.items=t.createProcessedItems(e.items,n+1,c,s),a.push(c)}),a},flatItems:function(e){var t=this,n=arguments.length>1&&arguments[1]!==void 0?arguments[1]:[];return e&&e.forEach(function(e){t.isVisibleItem(e)&&(n.push(e),t.flatItems(e.items,n))}),n}},computed:{processedItems:function(){return this.createProcessedItems(this.items||[])},visibleItems:function(){return this.flatItems(this.processedItems)},focusedItemId:function(){return m(this.focusedItem)?`${this.panelId}_${this.focusedItem.key}`:null}},components:{PanelMenuSub:G}};function Ie(e,t,n,r,i,c){var l=o(`PanelMenuSub`);return s(),a(l,b({id:n.panelId+`_list`,class:e.cx(`rootList`),role:`tree`,tabindex:-1,"aria-activedescendant":i.focused?c.focusedItemId:void 0,panelId:n.panelId,focusedItemId:i.focused?c.focusedItemId:void 0,items:c.processedItems,templates:n.templates,activeItemPath:i.activeItemPath,onFocus:c.onFocus,onBlur:c.onBlur,onKeydown:c.onKeyDown,onItemToggle:c.onItemToggle,onItemMousemove:c.onItemMouseMove,pt:e.pt,unstyled:e.unstyled},e.ptm(`rootList`)),null,16,[`id`,`class`,`aria-activedescendant`,`panelId`,`focusedItemId`,`items`,`templates`,`activeItemPath`,`onFocus`,`onBlur`,`onKeydown`,`onItemToggle`,`onItemMousemove`,`pt`,`unstyled`])}q.render=Ie;function J(e){"@babel/helpers - typeof";return J=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},J(e)}function Y(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function Le(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?Y(Object(n),!0).forEach(function(t){Re(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):Y(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function Re(e,t,n){return(t=ze(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function ze(e){var t=Be(e,`string`);return J(t)==`symbol`?t:t+``}function Be(e,t){if(J(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(J(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var X={name:`PanelMenu`,extends:Te,inheritAttrs:!1,emits:[`update:expandedKeys`,`panel-open`,`panel-close`],data:function(){return{activeItem:null,activeItems:[]}},methods:{getItemProp:function(e,t){return e?i(e[t]):void 0},getItemLabel:function(e){return this.getItemProp(e,`label`)},getPTOptions:function(e,t,n){return this.ptm(e,{context:{index:n,active:this.isItemActive(t),focused:this.isItemFocused(t),disabled:this.isItemDisabled(t)}})},isItemActive:function(e){return this.expandedKeys?this.expandedKeys[this.getItemProp(e,`key`)]:this.multiple?this.activeItems.some(function(t){return z(e,t)}):z(e,this.activeItem)},isItemVisible:function(e){return this.getItemProp(e,`visible`)!==!1},isItemDisabled:function(e){return this.getItemProp(e,`disabled`)},isItemFocused:function(e){return z(e,this.activeItem)},isItemGroup:function(e){return m(e.items)},getPanelId:function(e){return`${this.$id}_${e}`},getPanelKey:function(e){return this.getPanelId(e)},getHeaderId:function(e){return`${this.getPanelId(e)}_header`},getContentId:function(e){return`${this.getPanelId(e)}_content`},onHeaderClick:function(e,t){if(this.isItemDisabled(t)){e.preventDefault();return}t.command&&t.command({originalEvent:e,item:t}),this.changeActiveItem(e,t),T(e.currentTarget)},onHeaderKeyDown:function(e,t){switch(e.code){case`ArrowDown`:this.onHeaderArrowDownKey(e);break;case`ArrowUp`:this.onHeaderArrowUpKey(e);break;case`Home`:this.onHeaderHomeKey(e);break;case`End`:this.onHeaderEndKey(e);break;case`Enter`:case`NumpadEnter`:case`Space`:this.onHeaderEnterKey(e,t);break}},onHeaderArrowDownKey:function(e){var t=N(e.currentTarget,`data-p-active`)===!0?B(e.currentTarget.nextElementSibling,`[data-pc-section="rootlist"]`):null;t?T(t):this.updateFocusedHeader({originalEvent:e,focusOnNext:!0}),e.preventDefault()},onHeaderArrowUpKey:function(e){var t=this.findPrevHeader(e.currentTarget.parentElement)||this.findLastHeader(),n=N(t,`data-p-active`)===!0?B(t.nextElementSibling,`[data-pc-section="rootlist"]`):null;n?T(n):this.updateFocusedHeader({originalEvent:e,focusOnNext:!1}),e.preventDefault()},onHeaderHomeKey:function(e){this.changeFocusedHeader(e,this.findFirstHeader()),e.preventDefault()},onHeaderEndKey:function(e){this.changeFocusedHeader(e,this.findLastHeader()),e.preventDefault()},onHeaderEnterKey:function(e,t){var n=B(e.currentTarget,`[data-pc-section="headerlink"]`);n?n.click():this.onHeaderClick(e,t),e.preventDefault()},findNextHeader:function(e){var t=B(arguments.length>1&&arguments[1]!==void 0&&arguments[1]?e:e.nextElementSibling,`[data-pc-section="header"]`);return t?N(t,`data-p-disabled`)?this.findNextHeader(t.parentElement):t:null},findPrevHeader:function(e){var t=B(arguments.length>1&&arguments[1]!==void 0&&arguments[1]?e:e.previousElementSibling,`[data-pc-section="header"]`);return t?N(t,`data-p-disabled`)?this.findPrevHeader(t.parentElement):t:null},findFirstHeader:function(){return this.findNextHeader(this.$el.firstElementChild,!0)},findLastHeader:function(){return this.findPrevHeader(this.$el.lastElementChild,!0)},updateFocusedHeader:function(e){var t=e.originalEvent,n=e.focusOnNext,r=e.selfCheck,i=t.currentTarget.closest(`[data-pc-section="panel"]`),a=r?B(i,`[data-pc-section="header"]`):n?this.findNextHeader(i):this.findPrevHeader(i);a?this.changeFocusedHeader(t,a):n?this.onHeaderHomeKey(t):this.onHeaderEndKey(t)},changeActiveItem:function(e,t){var n=arguments.length>2&&arguments[2]!==void 0&&arguments[2];if(!this.isItemDisabled(t)){var r=this.isItemActive(t),i=r?`panel-close`:`panel-open`;this.activeItem=n?t:this.activeItem&&z(t,this.activeItem)?null:t,this.multiple&&(this.activeItems.some(function(e){return z(t,e)})?this.activeItems=this.activeItems.filter(function(e){return!z(t,e)}):this.activeItems.push(t)),this.changeExpandedKeys({item:t,expanded:!r}),this.$emit(i,{originalEvent:e,item:t})}},changeExpandedKeys:function(e){var t=e.item,n=e.expanded,r=n!==void 0&&n;if(this.expandedKeys){var i=Le({},this.expandedKeys);r?i[t.key]=!0:delete i[t.key],this.$emit(`update:expandedKeys`,i)}},changeFocusedHeader:function(e,t){t&&T(t)},getMenuItemProps:function(e,t){return{icon:b({class:[this.cx(`headerIcon`),this.getItemProp(e,`icon`)]},this.getPTOptions(`headerIcon`,e,t)),label:b({class:this.cx(`headerLabel`)},this.getPTOptions(`headerLabel`,e,t))}}},components:{PanelMenuList:q,ChevronRightIcon:U,ChevronDownIcon:V}},Ve=[`id`],He=[`id`,`tabindex`,`aria-label`,`aria-expanded`,`aria-controls`,`aria-disabled`,`onClick`,`onKeydown`,`data-p-active`,`data-p-disabled`],Ue=[`href`],We=[`id`,`aria-labelledby`];function Ge(r,i,u,p,m,_){var v=o(`PanelMenuList`);return s(),d(`div`,b({id:r.$id,class:r.cx(`root`)},r.ptmi(`root`)),[(s(!0),d(f,null,y(r.model,function(i,o){return s(),d(f,{key:_.getPanelKey(o)},[_.isItemVisible(i)?(s(),d(`div`,b({key:0,style:_.getItemProp(i,`style`),class:[r.cx(`panel`),_.getItemProp(i,`class`)]},{ref_for:!0},r.ptm(`panel`)),[h(`div`,b({id:_.getHeaderId(o),class:[r.cx(`header`,{item:i}),_.getItemProp(i,`headerClass`)],tabindex:_.isItemDisabled(i)?-1:r.tabindex,role:`button`,"aria-label":_.getItemLabel(i),"aria-expanded":_.isItemActive(i),"aria-controls":_.getContentId(o),"aria-disabled":_.isItemDisabled(i),onClick:function(e){return _.onHeaderClick(e,i)},onKeydown:function(e){return _.onHeaderKeyDown(e,i)}},{ref_for:!0},_.getPTOptions(`header`,i,o),{"data-p-active":_.isItemActive(i),"data-p-disabled":_.isItemDisabled(i)}),[h(`div`,b({class:r.cx(`headerContent`)},{ref_for:!0},_.getPTOptions(`headerContent`,i,o)),[r.$slots.item?(s(),a(P(r.$slots.item),{key:1,item:i,root:!0,active:_.isItemActive(i),hasSubmenu:_.isItemGroup(i),label:_.getItemLabel(i),props:_.getMenuItemProps(i,o)},null,8,[`item`,`active`,`hasSubmenu`,`label`,`props`])):(s(),d(`a`,b({key:0,href:_.getItemProp(i,`url`),class:r.cx(`headerLink`),tabindex:-1},{ref_for:!0},_.getPTOptions(`headerLink`,i,o)),[_.getItemProp(i,`items`)?x(r.$slots,`submenuicon`,{key:0,active:_.isItemActive(i)},function(){return[(s(),a(P(_.isItemActive(i)?`ChevronDownIcon`:`ChevronRightIcon`),b({class:r.cx(`submenuIcon`)},{ref_for:!0},_.getPTOptions(`submenuIcon`,i,o)),null,16,[`class`]))]}):c(``,!0),r.$slots.headericon?(s(),a(P(r.$slots.headericon),{key:1,item:i,class:g([r.cx(`headerIcon`),_.getItemProp(i,`icon`)])},null,8,[`item`,`class`])):_.getItemProp(i,`icon`)?(s(),d(`span`,b({key:2,class:[r.cx(`headerIcon`),_.getItemProp(i,`icon`)]},{ref_for:!0},_.getPTOptions(`headerIcon`,i,o)),null,16)):c(``,!0),h(`span`,b({class:r.cx(`headerLabel`)},{ref_for:!0},_.getPTOptions(`headerLabel`,i,o)),e(_.getItemLabel(i)),17)],16,Ue))],16)],16,He),n(M,b({name:`p-collapsible`},{ref_for:!0},r.ptm(`transition`)),{default:l(function(){return[C(h(`div`,b({id:_.getContentId(o),class:r.cx(`contentContainer`),role:`region`,"aria-labelledby":_.getHeaderId(o)},{ref_for:!0},r.ptm(`contentContainer`)),[h(`div`,b({class:r.cx(`contentWrapper`)},{ref_for:!0},r.ptm(`contentWrapper`)),[_.getItemProp(i,`items`)?(s(),d(`div`,b({key:0,class:r.cx(`content`)},{ref_for:!0},r.ptm(`content`)),[n(v,{panelId:_.getPanelId(o),items:_.getItemProp(i,`items`),templates:r.$slots,expandedKeys:r.expandedKeys,onItemToggle:_.changeExpandedKeys,onHeaderFocus:_.updateFocusedHeader,pt:r.pt,unstyled:r.unstyled},null,8,[`panelId`,`items`,`templates`,`expandedKeys`,`onItemToggle`,`onHeaderFocus`,`pt`,`unstyled`])],16)):c(``,!0)],16)],16,We),[[t,_.isItemActive(i)]])]}),_:2},1040)],16)):c(``,!0)],64)}),128))],16,Ve)}X.render=Ge;var Ke={class:`flex items-center h-12 px-4 border-b border-gray-200 shrink-0`},qe={key:0,class:`font-semibold text-sm text-gray-800 truncate`},Je={key:0,class:`flex-1 overflow-y-auto py-2 px-2`},Ye={key:1,class:`flex-1 overflow-y-auto py-3 px-1 flex flex-col items-center gap-1`},Xe=[`onClick`],Ze=ue({__name:`Sidebar`,props:{collapsed:{type:Boolean,default:!1}},emits:[`toggle`],setup(e){let t=O(),r=se(),i=[{label:`Dashboard`,icon:`pi pi-home`,command:()=>t.push(`/dashboard`),class:r.name===`Dashboard`?`bg-emerald-50 text-emerald-700 rounded-md`:``},{label:`Core HR`,icon:`pi pi-building`,items:[{label:`Organization`,icon:`pi pi-sitemap`,command:()=>t.push(`/organizations`)},{label:`Employees`,icon:`pi pi-users`,command:()=>t.push(`/employees`)},{label:`Job Management`,icon:`pi pi-briefcase`,command:()=>t.push(`/job-management`)}]},{label:`Talent`,icon:`pi pi-star`,items:[{label:`Competency`,icon:`pi pi-star`,command:()=>t.push(`/competencies`)},{label:`Performance`,icon:`pi pi-chart-line`,command:()=>t.push(`/performance`)},{label:`Training`,icon:`pi pi-book`,command:()=>t.push(`/training`)},{label:`Recruitment`,icon:`pi pi-user-plus`,command:()=>t.push(`/recruitment`)}]},{label:`Operations`,icon:`pi pi-cog`,items:[{label:`Attendance`,icon:`pi pi-clock`,command:()=>t.push(`/attendance`)},{label:`Leave`,icon:`pi pi-calendar`,command:()=>t.push(`/leave`)},{label:`Movement`,icon:`pi pi-arrows-alt`,command:()=>t.push(`/employee-movements`)},{label:`Approval`,icon:`pi pi-check-square`,command:()=>t.push(`/approvals`)}]},{label:`Finance`,icon:`pi pi-dollar`,items:[{label:`Payroll`,icon:`pi pi-dollar`,command:()=>t.push(`/payroll`)},{label:`Reimbursement`,icon:`pi pi-credit-card`,command:()=>t.push(`/reimbursements`)}]},{label:`Strategic`,icon:`pi pi-chart-bar`,items:[{label:`Workforce Intel`,icon:`pi pi-chart-bar`,command:()=>t.push(`/workforce-intelligence`)},{label:`Career Intel`,icon:`pi pi-road`,command:()=>t.push(`/career-intelligence`)}]}],a=p(()=>[{key:`Dashboard`,label:`Dashboard`,path:`/dashboard`,icon:`pi pi-home`,command:()=>t.push(`/dashboard`)},{key:`CoreHR`,label:`Core HR`,path:`/organizations`,icon:`pi pi-building`,command:()=>t.push(`/organizations`)},{key:`Talent`,label:`Talent`,path:`/competencies`,icon:`pi pi-star`,command:()=>t.push(`/competencies`)},{key:`Operations`,label:`Operations`,path:`/attendance`,icon:`pi pi-cog`,command:()=>t.push(`/attendance`)},{key:`Finance`,label:`Finance`,path:`/payroll`,icon:`pi pi-dollar`,command:()=>t.push(`/payroll`)},{key:`Strategic`,label:`Strategic`,path:`/workforce-intelligence`,icon:`pi pi-chart-bar`,command:()=>t.push(`/workforce-intelligence`)}]);function o(e){return e.path?r.path.startsWith(e.path):!1}return(t,r)=>{let l=u(`tooltip`);return s(),d(`aside`,{class:g([`flex flex-col bg-white border-r border-gray-200 transition-all duration-200 overflow-hidden`,e.collapsed?`w-16`:`w-60`])},[h(`div`,Ke,[r[0]||=h(`i`,{class:`pi pi-building text-emerald-600 text-lg mr-2`},null,-1),e.collapsed?c(``,!0):(s(),d(`span`,qe,`HRIS Platform`))]),e.collapsed?(s(),d(`nav`,Ye,[(s(!0),d(f,null,y(a.value,e=>C((s(),d(`div`,{key:e.key||e.label,class:g([`w-9 h-9 rounded-lg flex items-center justify-center cursor-pointer hover:bg-emerald-100 transition-colors`,{"bg-emerald-100 text-emerald-700":o(e)}]),onClick:t=>e.command?.()},[h(`i`,{class:g([e.icon,`text-sm`])},null,2)],10,Xe)),[[l,e.tooltip||e.label,void 0,{left:!0}]])),128))])):(s(),d(`nav`,Je,[n(w(X),{model:i,class:`border-none !bg-transparent`})])),r[1]||=h(`div`,{class:`border-t border-gray-200 p-3 shrink-0`},[h(`div`,{class:`flex items-center gap-2 text-sm text-gray-500`},[h(`i`,{class:`pi pi-circle-fill text-emerald-400 text-[6px]`}),h(`span`,{class:`truncate`},`Tenant: PT. ABC`)])],-1)],2)}}},[[`__scopeId`,`data-v-c7458717`]]),Qe={name:`BaseInputText`,extends:{name:`BaseInput`,extends:de,props:{size:{type:String,default:null},fluid:{type:Boolean,default:null},variant:{type:String,default:null}},inject:{$parentInstance:{default:void 0},$pcFluid:{default:void 0}},computed:{$variant:function(){return this.variant??(this.$primevue.config.inputStyle||this.$primevue.config.inputVariant)},$fluid:function(){return this.fluid??!!this.$pcFluid},hasFluid:function(){return this.$fluid}}},style:L.extend({name:`inputtext`,style:`
    .p-inputtext {
        font-family: inherit;
        font-feature-settings: inherit;
        font-size: 1rem;
        color: dt('inputtext.color');
        background: dt('inputtext.background');
        padding-block: dt('inputtext.padding.y');
        padding-inline: dt('inputtext.padding.x');
        border: 1px solid dt('inputtext.border.color');
        transition:
            background dt('inputtext.transition.duration'),
            color dt('inputtext.transition.duration'),
            border-color dt('inputtext.transition.duration'),
            outline-color dt('inputtext.transition.duration'),
            box-shadow dt('inputtext.transition.duration');
        appearance: none;
        border-radius: dt('inputtext.border.radius');
        outline-color: transparent;
        box-shadow: dt('inputtext.shadow');
    }

    .p-inputtext:enabled:hover {
        border-color: dt('inputtext.hover.border.color');
    }

    .p-inputtext:enabled:focus {
        border-color: dt('inputtext.focus.border.color');
        box-shadow: dt('inputtext.focus.ring.shadow');
        outline: dt('inputtext.focus.ring.width') dt('inputtext.focus.ring.style') dt('inputtext.focus.ring.color');
        outline-offset: dt('inputtext.focus.ring.offset');
    }

    .p-inputtext.p-invalid {
        border-color: dt('inputtext.invalid.border.color');
    }

    .p-inputtext.p-variant-filled {
        background: dt('inputtext.filled.background');
    }

    .p-inputtext.p-variant-filled:enabled:hover {
        background: dt('inputtext.filled.hover.background');
    }

    .p-inputtext.p-variant-filled:enabled:focus {
        background: dt('inputtext.filled.focus.background');
    }

    .p-inputtext:disabled {
        opacity: 1;
        background: dt('inputtext.disabled.background');
        color: dt('inputtext.disabled.color');
    }

    .p-inputtext::placeholder {
        color: dt('inputtext.placeholder.color');
    }

    .p-inputtext.p-invalid::placeholder {
        color: dt('inputtext.invalid.placeholder.color');
    }

    .p-inputtext-sm {
        font-size: dt('inputtext.sm.font.size');
        padding-block: dt('inputtext.sm.padding.y');
        padding-inline: dt('inputtext.sm.padding.x');
    }

    .p-inputtext-lg {
        font-size: dt('inputtext.lg.font.size');
        padding-block: dt('inputtext.lg.padding.y');
        padding-inline: dt('inputtext.lg.padding.x');
    }

    .p-inputtext-fluid {
        width: 100%;
    }
`,classes:{root:function(e){var t=e.instance,n=e.props;return[`p-inputtext p-component`,{"p-filled":t.$filled,"p-inputtext-sm p-inputfield-sm":n.size===`small`,"p-inputtext-lg p-inputfield-lg":n.size===`large`,"p-invalid":t.$invalid,"p-variant-filled":t.$variant===`filled`,"p-inputtext-fluid":t.$fluid}]}}}),provide:function(){return{$pcInputText:this,$parentInstance:this}}};function Z(e){"@babel/helpers - typeof";return Z=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},Z(e)}function $e(e,t,n){return(t=et(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function et(e){var t=tt(e,`string`);return Z(t)==`symbol`?t:t+``}function tt(e,t){if(Z(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(Z(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var nt={name:`InputText`,extends:Qe,inheritAttrs:!1,methods:{onInput:function(e){this.writeValue(e.target.value,e)}},computed:{attrs:function(){return b(this.ptmi(`root`,{context:{filled:this.$filled,disabled:this.disabled}}),this.formField)},dataP:function(){return D($e({invalid:this.$invalid,fluid:this.$fluid,filled:this.$variant===`filled`},this.size,this.size))}}},rt=[`value`,`name`,`disabled`,`aria-invalid`,`data-p`];function it(e,t,n,r,i,a){return s(),d(`input`,b({type:`text`,class:e.cx(`root`),value:e.d_value,name:e.name,disabled:e.disabled,"aria-invalid":e.$invalid||void 0,"data-p":a.dataP,onInput:t[0]||=function(){return a.onInput&&a.onInput.apply(a,arguments)}},a.attrs),null,16,rt)}nt.render=it;var at={name:`InputIcon`,extends:{name:`BaseInputIcon`,extends:E,style:L.extend({name:`inputicon`,classes:{root:`p-inputicon`}}),props:{class:null},provide:function(){return{$pcInputIcon:this,$parentInstance:this}}},inheritAttrs:!1,computed:{containerClass:function(){return[this.cx(`root`),this.class]}}};function ot(e,t,n,r,i,a){return s(),d(`span`,b({class:a.containerClass},e.ptmi(`root`),{"aria-hidden":`true`}),[x(e.$slots,`default`)],16)}at.render=ot;var st={name:`IconField`,extends:{name:`BaseIconField`,extends:E,style:L.extend({name:`iconfield`,style:`
    .p-iconfield {
        position: relative;
        display: block;
    }

    .p-inputicon {
        position: absolute;
        top: 50%;
        margin-top: calc(-1 * (dt('icon.size') / 2));
        color: dt('iconfield.icon.color');
        line-height: 1;
        z-index: 1;
    }

    .p-iconfield .p-inputicon:first-child {
        inset-inline-start: dt('form.field.padding.x');
    }

    .p-iconfield .p-inputicon:last-child {
        inset-inline-end: dt('form.field.padding.x');
    }

    .p-iconfield .p-inputtext:not(:first-child),
    .p-iconfield .p-inputwrapper:not(:first-child) .p-inputtext {
        padding-inline-start: calc((dt('form.field.padding.x') * 2) + dt('icon.size'));
    }

    .p-iconfield .p-inputtext:not(:last-child) {
        padding-inline-end: calc((dt('form.field.padding.x') * 2) + dt('icon.size'));
    }

    .p-iconfield:has(.p-inputfield-sm) .p-inputicon {
        font-size: dt('form.field.sm.font.size');
        width: dt('form.field.sm.font.size');
        height: dt('form.field.sm.font.size');
        margin-top: calc(-1 * (dt('form.field.sm.font.size') / 2));
    }

    .p-iconfield:has(.p-inputfield-lg) .p-inputicon {
        font-size: dt('form.field.lg.font.size');
        width: dt('form.field.lg.font.size');
        height: dt('form.field.lg.font.size');
        margin-top: calc(-1 * (dt('form.field.lg.font.size') / 2));
    }
`,classes:{root:`p-iconfield`}}),provide:function(){return{$pcIconField:this,$parentInstance:this}}},inheritAttrs:!1};function ct(e,t,n,r,i,a){return s(),d(`div`,b({class:e.cx(`root`)},e.ptmi(`root`)),[x(e.$slots,`default`)],16)}st.render=ct;var lt=L.extend({name:`avatar`,style:`
    .p-avatar {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: dt('avatar.width');
        height: dt('avatar.height');
        font-size: dt('avatar.font.size');
        background: dt('avatar.background');
        color: dt('avatar.color');
        border-radius: dt('avatar.border.radius');
    }

    .p-avatar-image {
        background: transparent;
    }

    .p-avatar-circle {
        border-radius: 50%;
    }

    .p-avatar-circle img {
        border-radius: 50%;
    }

    .p-avatar-icon {
        font-size: dt('avatar.icon.size');
        width: dt('avatar.icon.size');
        height: dt('avatar.icon.size');
    }

    .p-avatar img {
        width: 100%;
        height: 100%;
    }

    .p-avatar-lg {
        width: dt('avatar.lg.width');
        height: dt('avatar.lg.width');
        font-size: dt('avatar.lg.font.size');
    }

    .p-avatar-lg .p-avatar-icon {
        font-size: dt('avatar.lg.icon.size');
        width: dt('avatar.lg.icon.size');
        height: dt('avatar.lg.icon.size');
    }

    .p-avatar-xl {
        width: dt('avatar.xl.width');
        height: dt('avatar.xl.width');
        font-size: dt('avatar.xl.font.size');
    }

    .p-avatar-xl .p-avatar-icon {
        font-size: dt('avatar.xl.icon.size');
        width: dt('avatar.xl.icon.size');
        height: dt('avatar.xl.icon.size');
    }

    .p-avatar-group {
        display: flex;
        align-items: center;
    }

    .p-avatar-group .p-avatar + .p-avatar {
        margin-inline-start: dt('avatar.group.offset');
    }

    .p-avatar-group .p-avatar {
        border: 2px solid dt('avatar.group.border.color');
    }

    .p-avatar-group .p-avatar-lg + .p-avatar-lg {
        margin-inline-start: dt('avatar.lg.group.offset');
    }

    .p-avatar-group .p-avatar-xl + .p-avatar-xl {
        margin-inline-start: dt('avatar.xl.group.offset');
    }
`,classes:{root:function(e){var t=e.props;return[`p-avatar p-component`,{"p-avatar-image":t.image!=null,"p-avatar-circle":t.shape===`circle`,"p-avatar-lg":t.size===`large`,"p-avatar-xl":t.size===`xlarge`}]},label:`p-avatar-label`,icon:`p-avatar-icon`}}),ut={name:`BaseAvatar`,extends:E,props:{label:{type:String,default:null},icon:{type:String,default:null},image:{type:String,default:null},size:{type:String,default:`normal`},shape:{type:String,default:`square`},ariaLabelledby:{type:String,default:null},ariaLabel:{type:String,default:null}},style:lt,provide:function(){return{$pcAvatar:this,$parentInstance:this}}};function Q(e){"@babel/helpers - typeof";return Q=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},Q(e)}function dt(e,t,n){return(t=ft(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function ft(e){var t=pt(e,`string`);return Q(t)==`symbol`?t:t+``}function pt(e,t){if(Q(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(Q(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var mt={name:`Avatar`,extends:ut,inheritAttrs:!1,emits:[`error`],methods:{onError:function(e){this.$emit(`error`,e)}},computed:{dataP:function(){return D(dt(dt({},this.shape,this.shape),this.size,this.size))}}},ht=[`aria-labelledby`,`aria-label`,`data-p`],gt=[`data-p`],_t=[`data-p`],vt=[`src`,`alt`,`data-p`];function yt(t,n,r,i,o,l){return s(),d(`div`,b({class:t.cx(`root`),"aria-labelledby":t.ariaLabelledby,"aria-label":t.ariaLabel},t.ptmi(`root`),{"data-p":l.dataP}),[x(t.$slots,`default`,{},function(){return[t.label?(s(),d(`span`,b({key:0,class:t.cx(`label`)},t.ptm(`label`),{"data-p":l.dataP}),e(t.label),17,gt)):t.$slots.icon?(s(),a(P(t.$slots.icon),{key:1,class:g(t.cx(`icon`))},null,8,[`class`])):t.icon?(s(),d(`span`,b({key:2,class:[t.cx(`icon`),t.icon]},t.ptm(`icon`),{"data-p":l.dataP}),null,16,_t)):t.image?(s(),d(`img`,b({key:3,src:t.image,alt:t.ariaLabel,onError:n[0]||=function(){return l.onError&&l.onError.apply(l,arguments)}},t.ptm(`image`),{"data-p":l.dataP}),null,16,vt)):c(``,!0)]})],16,ht)}mt.render=yt;var bt=ce(),xt=L.extend({name:`menu`,style:`
    .p-menu {
        background: dt('menu.background');
        color: dt('menu.color');
        border: 1px solid dt('menu.border.color');
        border-radius: dt('menu.border.radius');
        min-width: 12.5rem;
    }

    .p-menu-list {
        margin: 0;
        padding: dt('menu.list.padding');
        outline: 0 none;
        list-style: none;
        display: flex;
        flex-direction: column;
        gap: dt('menu.list.gap');
    }

    .p-menu-item-content {
        transition:
            background dt('menu.transition.duration'),
            color dt('menu.transition.duration');
        border-radius: dt('menu.item.border.radius');
        color: dt('menu.item.color');
        overflow: hidden;
    }

    .p-menu-item-link {
        cursor: pointer;
        display: flex;
        align-items: center;
        text-decoration: none;
        overflow: hidden;
        position: relative;
        color: inherit;
        padding: dt('menu.item.padding');
        gap: dt('menu.item.gap');
        user-select: none;
        outline: 0 none;
    }

    .p-menu-item-label {
        line-height: 1;
    }

    .p-menu-item-icon {
        color: dt('menu.item.icon.color');
    }

    .p-menu-item.p-focus .p-menu-item-content {
        color: dt('menu.item.focus.color');
        background: dt('menu.item.focus.background');
    }

    .p-menu-item.p-focus .p-menu-item-icon {
        color: dt('menu.item.icon.focus.color');
    }

    .p-menu-item:not(.p-disabled) .p-menu-item-content:hover {
        color: dt('menu.item.focus.color');
        background: dt('menu.item.focus.background');
    }

    .p-menu-item:not(.p-disabled) .p-menu-item-content:hover .p-menu-item-icon {
        color: dt('menu.item.icon.focus.color');
    }

    .p-menu-overlay {
        box-shadow: dt('menu.shadow');
    }

    .p-menu-submenu-label {
        background: dt('menu.submenu.label.background');
        padding: dt('menu.submenu.label.padding');
        color: dt('menu.submenu.label.color');
        font-weight: dt('menu.submenu.label.font.weight');
    }

    .p-menu-separator {
        border-block-start: 1px solid dt('menu.separator.border.color');
    }
`,classes:{root:function(e){return[`p-menu p-component`,{"p-menu-overlay":e.props.popup}]},start:`p-menu-start`,list:`p-menu-list`,submenuLabel:`p-menu-submenu-label`,separator:`p-menu-separator`,end:`p-menu-end`,item:function(e){var t=e.instance;return[`p-menu-item`,{"p-focus":t.id===t.focusedOptionId,"p-disabled":t.disabled()}]},itemContent:`p-menu-item-content`,itemLink:`p-menu-item-link`,itemIcon:`p-menu-item-icon`,itemLabel:`p-menu-item-label`}}),St={name:`BaseMenu`,extends:E,props:{popup:{type:Boolean,default:!1},model:{type:Array,default:null},appendTo:{type:[String,Object],default:`body`},autoZIndex:{type:Boolean,default:!0},baseZIndex:{type:Number,default:0},tabindex:{type:Number,default:0},ariaLabel:{type:String,default:null},ariaLabelledby:{type:String,default:null}},style:xt,provide:function(){return{$pcMenu:this,$parentInstance:this}}},Ct={name:`Menuitem`,hostName:`Menu`,extends:E,inheritAttrs:!1,emits:[`item-click`,`item-mousemove`],props:{item:null,templates:null,id:null,focusedOptionId:null,index:null},methods:{getItemProp:function(e,t){return e&&e.item?i(e.item[t]):void 0},getPTOptions:function(e){return this.ptm(e,{context:{item:this.item,index:this.index,focused:this.isItemFocused(),disabled:this.disabled()}})},isItemFocused:function(){return this.focusedOptionId===this.id},onItemClick:function(e){var t=this.getItemProp(this.item,`command`);t&&t({originalEvent:e,item:this.item.item}),this.$emit(`item-click`,{originalEvent:e,item:this.item,id:this.id})},onItemMouseMove:function(e){this.$emit(`item-mousemove`,{originalEvent:e,item:this.item,id:this.id})},visible:function(){return typeof this.item.visible==`function`?this.item.visible():this.item.visible!==!1},disabled:function(){return typeof this.item.disabled==`function`?this.item.disabled():this.item.disabled},label:function(){return typeof this.item.label==`function`?this.item.label():this.item.label},getMenuItemProps:function(e){return{action:b({class:this.cx(`itemLink`),tabindex:`-1`},this.getPTOptions(`itemLink`)),icon:b({class:[this.cx(`itemIcon`),e.icon]},this.getPTOptions(`itemIcon`)),label:b({class:this.cx(`itemLabel`)},this.getPTOptions(`itemLabel`))}}},computed:{dataP:function(){return D({focus:this.isItemFocused(),disabled:this.disabled()})}},directives:{ripple:F}},wt=[`id`,`aria-label`,`aria-disabled`,`data-p-focused`,`data-p-disabled`,`data-p`],Tt=[`data-p`],Et=[`href`,`target`],Dt=[`data-p`],Ot=[`data-p`];function kt(t,n,r,i,o,l){var f=u(`ripple`);return l.visible()?(s(),d(`li`,b({key:0,id:r.id,class:[t.cx(`item`),r.item.class],role:`menuitem`,style:r.item.style,"aria-label":l.label(),"aria-disabled":l.disabled(),"data-p-focused":l.isItemFocused(),"data-p-disabled":l.disabled()||!1,"data-p":l.dataP},l.getPTOptions(`item`)),[h(`div`,b({class:t.cx(`itemContent`),onClick:n[0]||=function(e){return l.onItemClick(e)},onMousemove:n[1]||=function(e){return l.onItemMouseMove(e)},"data-p":l.dataP},l.getPTOptions(`itemContent`)),[r.templates.item?r.templates.item?(s(),a(P(r.templates.item),{key:1,item:r.item,label:l.label(),props:l.getMenuItemProps(r.item)},null,8,[`item`,`label`,`props`])):c(``,!0):C((s(),d(`a`,b({key:0,href:r.item.url,class:t.cx(`itemLink`),target:r.item.target,tabindex:`-1`},l.getPTOptions(`itemLink`)),[r.templates.itemicon?(s(),a(P(r.templates.itemicon),{key:0,item:r.item,class:g(t.cx(`itemIcon`))},null,8,[`item`,`class`])):r.item.icon?(s(),d(`span`,b({key:1,class:[t.cx(`itemIcon`),r.item.icon],"data-p":l.dataP},l.getPTOptions(`itemIcon`)),null,16,Dt)):c(``,!0),h(`span`,b({class:t.cx(`itemLabel`),"data-p":l.dataP},l.getPTOptions(`itemLabel`)),e(l.label()),17,Ot)],16,Et)),[[f]])],16,Tt)],16,wt)):c(``,!0)}Ct.render=kt;function At(e){return Pt(e)||Nt(e)||Mt(e)||jt()}function jt(){throw TypeError(`Invalid attempt to spread non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}function Mt(e,t){if(e){if(typeof e==`string`)return $(e,t);var n={}.toString.call(e).slice(8,-1);return n===`Object`&&e.constructor&&(n=e.constructor.name),n===`Map`||n===`Set`?Array.from(e):n===`Arguments`||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?$(e,t):void 0}}function Nt(e){if(typeof Symbol<`u`&&e[Symbol.iterator]!=null||e[`@@iterator`]!=null)return Array.from(e)}function Pt(e){if(Array.isArray(e))return $(e)}function $(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=Array(t);n<t;n++)r[n]=e[n];return r}var Ft={name:`Menu`,extends:St,inheritAttrs:!1,emits:[`show`,`hide`,`focus`,`blur`],data:function(){return{overlayVisible:!1,focused:!1,focusedOptionIndex:-1,selectedOptionIndex:-1}},target:null,outsideClickListener:null,scrollHandler:null,resizeListener:null,container:null,list:null,mounted:function(){this.popup||(this.bindResizeListener(),this.bindOutsideClickListener())},beforeUnmount:function(){this.unbindResizeListener(),this.unbindOutsideClickListener(),this.scrollHandler&&=(this.scrollHandler.destroy(),null),this.target=null,this.container&&this.autoZIndex&&k.clear(this.container),this.container=null},methods:{itemClick:function(e){var t=e.item;this.disabled(t)||(t.command&&t.command(e),this.overlayVisible&&this.hide(),!this.popup&&this.focusedOptionIndex!==e.id&&(this.focusedOptionIndex=e.id))},itemMouseMove:function(e){this.focused&&(this.focusedOptionIndex=e.id)},onListFocus:function(e){this.focused=!0,!this.popup&&this.changeFocusedOptionIndex(0),this.$emit(`focus`,e)},onListBlur:function(e){this.focused=!1,this.focusedOptionIndex=-1,this.$emit(`blur`,e)},onListKeyDown:function(e){switch(e.code){case`ArrowDown`:this.onArrowDownKey(e);break;case`ArrowUp`:this.onArrowUpKey(e);break;case`Home`:this.onHomeKey(e);break;case`End`:this.onEndKey(e);break;case`Enter`:case`NumpadEnter`:this.onEnterKey(e);break;case`Space`:this.onSpaceKey(e);break;case`Escape`:this.popup&&(T(this.target),this.hide());case`Tab`:this.overlayVisible&&this.hide();break}},onArrowDownKey:function(e){var t=this.findNextOptionIndex(this.focusedOptionIndex);this.changeFocusedOptionIndex(t),e.preventDefault()},onArrowUpKey:function(e){if(e.altKey&&this.popup)T(this.target),this.hide(),e.preventDefault();else{var t=this.findPrevOptionIndex(this.focusedOptionIndex);this.changeFocusedOptionIndex(t),e.preventDefault()}},onHomeKey:function(e){this.changeFocusedOptionIndex(0),e.preventDefault()},onEndKey:function(e){this.changeFocusedOptionIndex(A(this.container,`li[data-pc-section="item"][data-p-disabled="false"]`).length-1),e.preventDefault()},onEnterKey:function(e){var t=B(this.list,`li[id="${`${this.focusedOptionIndex}`}"]`),n=t&&B(t,`a[data-pc-section="itemlink"]`);this.popup&&T(this.target),n?n.click():t&&t.click(),e.preventDefault()},onSpaceKey:function(e){this.onEnterKey(e)},findNextOptionIndex:function(e){var t=At(A(this.container,`li[data-pc-section="item"][data-p-disabled="false"]`)).findIndex(function(t){return t.id===e});return t>-1?t+1:0},findPrevOptionIndex:function(e){var t=At(A(this.container,`li[data-pc-section="item"][data-p-disabled="false"]`)).findIndex(function(t){return t.id===e});return t>-1?t-1:0},changeFocusedOptionIndex:function(e){var t=A(this.container,`li[data-pc-section="item"][data-p-disabled="false"]`),n=e>=t.length?t.length-1:e<0?0:e;n>-1&&(this.focusedOptionIndex=t[n].getAttribute(`id`))},toggle:function(e,t){this.overlayVisible?this.hide():this.show(e,t)},show:function(e,t){this.overlayVisible=!0,this.target=t??e.currentTarget},hide:function(){this.overlayVisible=!1,this.target=null},onEnter:function(e){ie(e,{position:`absolute`,top:`0`}),this.alignOverlay(),this.bindOutsideClickListener(),this.bindResizeListener(),this.bindScrollListener(),this.autoZIndex&&k.set(`menu`,e,this.baseZIndex||this.$primevue.config.zIndex.menu),this.popup&&T(this.list),this.$emit(`show`)},onLeave:function(){this.unbindOutsideClickListener(),this.unbindResizeListener(),this.unbindScrollListener(),this.$emit(`hide`)},onAfterLeave:function(e){this.autoZIndex&&k.clear(e)},alignOverlay:function(){ee(this.container,this.target),R(this.target)>R(this.container)&&(this.container.style.minWidth=R(this.target)+`px`)},bindOutsideClickListener:function(){var e=this;this.outsideClickListener||(this.outsideClickListener=function(t){var n=e.container&&!e.container.contains(t.target),r=!(e.target&&(e.target===t.target||e.target.contains(t.target)));e.overlayVisible&&n&&r?e.hide():!e.popup&&n&&r&&(e.focusedOptionIndex=-1)},document.addEventListener(`click`,this.outsideClickListener,!0))},unbindOutsideClickListener:function(){this.outsideClickListener&&=(document.removeEventListener(`click`,this.outsideClickListener,!0),null)},bindScrollListener:function(){var e=this;this.scrollHandler||=new re(this.target,function(){e.overlayVisible&&e.hide()}),this.scrollHandler.bindScrollListener()},unbindScrollListener:function(){this.scrollHandler&&this.scrollHandler.unbindScrollListener()},bindResizeListener:function(){var e=this;this.resizeListener||(this.resizeListener=function(){e.overlayVisible&&!te()&&e.hide()},window.addEventListener(`resize`,this.resizeListener))},unbindResizeListener:function(){this.resizeListener&&=(window.removeEventListener(`resize`,this.resizeListener),null)},visible:function(e){return typeof e.visible==`function`?e.visible():e.visible!==!1},disabled:function(e){return typeof e.disabled==`function`?e.disabled():e.disabled},label:function(e){return typeof e.label==`function`?e.label():e.label},onOverlayClick:function(e){bt.emit(`overlay-click`,{originalEvent:e,target:this.target})},containerRef:function(e){this.container=e},listRef:function(e){this.list=e}},computed:{focusedOptionId:function(){return this.focusedOptionIndex===-1?null:this.focusedOptionIndex},dataP:function(){return D({popup:this.popup})}},components:{PVMenuitem:Ct,Portal:oe}},It=[`id`,`data-p`],Lt=[`id`,`tabindex`,`aria-activedescendant`,`aria-label`,`aria-labelledby`],Rt=[`id`];function zt(t,r,i,u,p,m){var g=o(`PVMenuitem`),_=o(`Portal`);return s(),a(_,{appendTo:t.appendTo,disabled:!t.popup},{default:l(function(){return[n(M,b({name:`p-anchored-overlay`,onEnter:m.onEnter,onLeave:m.onLeave,onAfterLeave:m.onAfterLeave},t.ptm(`transition`)),{default:l(function(){return[!t.popup||p.overlayVisible?(s(),d(`div`,b({key:0,ref:m.containerRef,id:t.$id,class:t.cx(`root`),onClick:r[3]||=function(){return m.onOverlayClick&&m.onOverlayClick.apply(m,arguments)},"data-p":m.dataP},t.ptmi(`root`)),[t.$slots.start?(s(),d(`div`,b({key:0,class:t.cx(`start`)},t.ptm(`start`)),[x(t.$slots,`start`)],16)):c(``,!0),h(`ul`,b({ref:m.listRef,id:t.$id+`_list`,class:t.cx(`list`),role:`menu`,tabindex:t.tabindex,"aria-activedescendant":p.focused?m.focusedOptionId:void 0,"aria-label":t.ariaLabel,"aria-labelledby":t.ariaLabelledby,onFocus:r[0]||=function(){return m.onListFocus&&m.onListFocus.apply(m,arguments)},onBlur:r[1]||=function(){return m.onListBlur&&m.onListBlur.apply(m,arguments)},onKeydown:r[2]||=function(){return m.onListKeyDown&&m.onListKeyDown.apply(m,arguments)}},t.ptm(`list`)),[(s(!0),d(f,null,y(t.model,function(n,r){return s(),d(f,{key:m.label(n)+r.toString()},[n.items&&m.visible(n)&&!n.separator?(s(),d(f,{key:0},[n.items?(s(),d(`li`,b({key:0,id:t.$id+`_`+r,class:[t.cx(`submenuLabel`),n.class],role:`none`},{ref_for:!0},t.ptm(`submenuLabel`)),[x(t.$slots,t.$slots.submenulabel?`submenulabel`:`submenuheader`,{item:n},function(){return[le(e(m.label(n)),1)]})],16,Rt)):c(``,!0),(s(!0),d(f,null,y(n.items,function(e,i){return s(),d(f,{key:e.label+r+`_`+i},[m.visible(e)&&!e.separator?(s(),a(g,{key:0,id:t.$id+`_`+r+`_`+i,item:e,templates:t.$slots,focusedOptionId:m.focusedOptionId,unstyled:t.unstyled,onItemClick:m.itemClick,onItemMousemove:m.itemMouseMove,pt:t.pt},null,8,[`id`,`item`,`templates`,`focusedOptionId`,`unstyled`,`onItemClick`,`onItemMousemove`,`pt`])):m.visible(e)&&e.separator?(s(),d(`li`,b({key:`separator`+r+i,class:[t.cx(`separator`),n.class],style:e.style,role:`separator`},{ref_for:!0},t.ptm(`separator`)),null,16)):c(``,!0)],64)}),128))],64)):m.visible(n)&&n.separator?(s(),d(`li`,b({key:`separator`+r.toString(),class:[t.cx(`separator`),n.class],style:n.style,role:`separator`},{ref_for:!0},t.ptm(`separator`)),null,16)):(s(),a(g,{key:m.label(n)+r.toString(),id:t.$id+`_`+r,item:n,index:r,templates:t.$slots,focusedOptionId:m.focusedOptionId,unstyled:t.unstyled,onItemClick:m.itemClick,onItemMousemove:m.itemMouseMove,pt:t.pt},null,8,[`id`,`item`,`index`,`templates`,`focusedOptionId`,`unstyled`,`onItemClick`,`onItemMousemove`,`pt`]))],64)}),128))],16,Lt),t.$slots.end?(s(),d(`div`,b({key:1,class:t.cx(`end`)},t.ptm(`end`)),[x(t.$slots,`end`)],16)):c(``,!0)],16,It)):c(``,!0)]}),_:3},16,[`onEnter`,`onLeave`,`onAfterLeave`])]}),_:3},8,[`appendTo`,`disabled`])}Ft.render=zt;var Bt={class:`flex items-center h-12 bg-white border-b border-gray-200 px-4 shrink-0`},Vt={class:`flex items-center gap-3 flex-1 min-w-0`},Ht={class:`text-sm text-gray-500 font-medium truncate`},Ut={class:`flex items-center gap-2`},Wt={class:`relative`},Gt={class:`flex items-center gap-2`},Kt={__name:`HeaderBar`,emits:[`toggle-sidebar`],setup(t){let r=se(),i=O(),a=S(``),o=S(null),c=S(!1),u=[{label:`Profile`,icon:`pi pi-user`,command:()=>{}},{label:`Settings`,icon:`pi pi-cog`,command:()=>{}},{separator:!0},{label:`Logout`,icon:`pi pi-sign-out`,command:()=>{}}];function f(){if(!a.value.trim())return;let e=a.value.toLowerCase(),t=i.getRoutes().find(t=>t.meta?.title?.toLowerCase().includes(e)||t.path?.toLowerCase().includes(e));t&&(i.push(t.path),a.value=``)}return(t,i)=>(s(),d(`header`,Bt,[h(`div`,Vt,[n(w(I),{icon:`pi pi-bars`,severity:`secondary`,text:``,size:`small`,onClick:i[0]||=e=>t.$emit(`toggle-sidebar`),class:`!p-1.5`}),h(`span`,Ht,e(w(r).meta?.title||`Dashboard`),1)]),h(`div`,Ut,[n(w(st),null,{default:l(()=>[n(w(at),{class:`pi pi-search`}),n(w(nt),{modelValue:a.value,"onUpdate:modelValue":i[1]||=e=>a.value=e,placeholder:`Search modules...`,class:`!w-48 !h-8 !text-sm`,onKeyup:ne(f,[`enter`])},null,8,[`modelValue`])]),_:1}),h(`div`,Wt,[n(w(I),{icon:`pi pi-bell`,severity:`secondary`,text:``,size:`small`,class:`!p-1.5`}),n(w(ae),{value:`3`,severity:`danger`,class:`!absolute -top-0.5 -right-0.5 !text-xs !min-w-[1.1rem] !h-[1.1rem]`})]),n(w(I),{severity:`secondary`,text:``,size:`small`,class:`!p-1`,onClick:i[2]||=e=>c.value=!c.value},{default:l(()=>[h(`div`,Gt,[n(w(mt),{icon:`pi pi-user`,size:`small`,class:`!w-7 !h-7 !bg-emerald-100 !text-emerald-700`}),i[3]||=h(`span`,{class:`text-sm text-gray-700 hidden sm:inline`},`Admin`,-1),i[4]||=h(`i`,{class:`pi pi-chevron-down text-sm text-gray-400`},null,-1)])]),_:1}),n(w(Ft),{ref_key:`userMenu`,ref:o,model:u,popup:``},null,512)])]))}},qt={class:`flex h-screen overflow-hidden bg-gray-50`},Jt={class:`flex flex-col flex-1 min-w-0`},Yt={class:`flex-1 overflow-auto p-4`},Xt={__name:`AppLayout`,setup(e){let t=S(!1);return(e,r)=>{let i=o(`router-view`);return s(),d(`div`,qt,[n(Ze,{collapsed:t.value,onToggle:r[0]||=e=>t.value=!t.value},null,8,[`collapsed`]),h(`div`,Jt,[n(Kt,{onToggleSidebar:r[1]||=e=>t.value=!t.value}),h(`main`,Yt,[n(i)])])])}}};export{Xt as default};