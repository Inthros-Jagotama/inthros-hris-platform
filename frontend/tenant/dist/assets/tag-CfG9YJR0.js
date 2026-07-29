import{O as e,T as t,b as n,c as r,j as i,l as a,lt as o,s,u as c}from"./runtime-core.esm-bundler-CMNNIOjW.js";import{n as l,o as u,tt as d}from"./ripple-wktLe1vL.js";function f(e){if(!e)return{};let t=e.response?.data||e,n=t?.error?.errors||t?.error?.fields||{},r={};for(let[e,t]of Object.entries(n))r[e]=Array.isArray(t)?t.join(`, `):String(t);return r}var p=u.extend({name:`tag`,style:`
    .p-tag {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        background: dt('tag.primary.background');
        color: dt('tag.primary.color');
        font-size: dt('tag.font.size');
        font-weight: dt('tag.font.weight');
        padding: dt('tag.padding');
        border-radius: dt('tag.border.radius');
        gap: dt('tag.gap');
    }

    .p-tag-icon {
        font-size: dt('tag.icon.size');
        width: dt('tag.icon.size');
        height: dt('tag.icon.size');
    }

    .p-tag-rounded {
        border-radius: dt('tag.rounded.border.radius');
    }

    .p-tag-success {
        background: dt('tag.success.background');
        color: dt('tag.success.color');
    }

    .p-tag-info {
        background: dt('tag.info.background');
        color: dt('tag.info.color');
    }

    .p-tag-warn {
        background: dt('tag.warn.background');
        color: dt('tag.warn.color');
    }

    .p-tag-danger {
        background: dt('tag.danger.background');
        color: dt('tag.danger.color');
    }

    .p-tag-secondary {
        background: dt('tag.secondary.background');
        color: dt('tag.secondary.color');
    }

    .p-tag-contrast {
        background: dt('tag.contrast.background');
        color: dt('tag.contrast.color');
    }
`,classes:{root:function(e){var t=e.props;return[`p-tag p-component`,{"p-tag-info":t.severity===`info`,"p-tag-success":t.severity===`success`,"p-tag-warn":t.severity===`warn`,"p-tag-danger":t.severity===`danger`,"p-tag-secondary":t.severity===`secondary`,"p-tag-contrast":t.severity===`contrast`,"p-tag-rounded":t.rounded}]},icon:`p-tag-icon`,label:`p-tag-label`}}),m={name:`BaseTag`,extends:l,props:{value:null,severity:null,rounded:Boolean,icon:String},style:p,provide:function(){return{$pcTag:this,$parentInstance:this}}};function h(e){"@babel/helpers - typeof";return h=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},h(e)}function g(e,t,n){return(t=_(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function _(e){var t=v(e,`string`);return h(t)==`symbol`?t:t+``}function v(e,t){if(h(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(h(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var y={name:`Tag`,extends:m,inheritAttrs:!1,computed:{dataP:function(){return d(g({rounded:this.rounded},this.severity,this.severity))}}},b=[`data-p`];function x(l,u,d,f,p,m){return t(),c(`span`,n({class:l.cx(`root`),"data-p":m.dataP},l.ptmi(`root`)),[l.$slots.icon?(t(),r(i(l.$slots.icon),n({key:0,class:l.cx(`icon`)},l.ptm(`icon`)),null,16,[`class`])):l.icon?(t(),c(`span`,n({key:1,class:[l.cx(`icon`),l.icon]},l.ptm(`icon`)),null,16)):a(``,!0),l.value!=null||l.$slots.default?e(l.$slots,`default`,{key:2},function(){return[s(`span`,n({class:l.cx(`label`)},l.ptm(`label`)),o(l.value),17)]}):a(``,!0)],16,b)}y.render=x;export{f as n,y as t};