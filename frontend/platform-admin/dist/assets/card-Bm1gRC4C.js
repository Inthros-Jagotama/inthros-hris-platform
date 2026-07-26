import{At as e,Jt as t,Mt as n,Nt as r,Ut as i,Zt as a,f as o,o as s}from"./button-CU9bjc3B.js";var c={name:`Card`,extends:{name:`BaseCard`,extends:s,style:o.extend({name:`card`,style:`
    .p-card {
        background: dt('card.background');
        color: dt('card.color');
        box-shadow: dt('card.shadow');
        border-radius: dt('card.border.radius');
        display: flex;
        flex-direction: column;
    }

    .p-card-caption {
        display: flex;
        flex-direction: column;
        gap: dt('card.caption.gap');
    }

    .p-card-body {
        padding: dt('card.body.padding');
        display: flex;
        flex-direction: column;
        gap: dt('card.body.gap');
    }

    .p-card-title {
        font-size: dt('card.title.font.size');
        font-weight: dt('card.title.font.weight');
    }

    .p-card-subtitle {
        color: dt('card.subtitle.color');
    }
`,classes:{root:`p-card p-component`,header:`p-card-header`,body:`p-card-body`,caption:`p-card-caption`,title:`p-card-title`,subtitle:`p-card-subtitle`,content:`p-card-content`,footer:`p-card-footer`}}),provide:function(){return{$pcCard:this,$parentInstance:this}}},inheritAttrs:!1};function l(o,s,c,l,u,d){return t(),r(`div`,i({class:o.cx(`root`)},o.ptmi(`root`)),[o.$slots.header?(t(),r(`div`,i({key:0,class:o.cx(`header`)},o.ptm(`header`)),[a(o.$slots,`header`)],16)):n(``,!0),e(`div`,i({class:o.cx(`body`)},o.ptm(`body`)),[o.$slots.title||o.$slots.subtitle?(t(),r(`div`,i({key:0,class:o.cx(`caption`)},o.ptm(`caption`)),[o.$slots.title?(t(),r(`div`,i({key:0,class:o.cx(`title`)},o.ptm(`title`)),[a(o.$slots,`title`)],16)):n(``,!0),o.$slots.subtitle?(t(),r(`div`,i({key:1,class:o.cx(`subtitle`)},o.ptm(`subtitle`)),[a(o.$slots,`subtitle`)],16)):n(``,!0)],16)):n(``,!0),e(`div`,i({class:o.cx(`content`)},o.ptm(`content`)),[a(o.$slots,`content`)],16),o.$slots.footer?(t(),r(`div`,i({key:1,class:o.cx(`footer`)},o.ptm(`footer`)),[a(o.$slots,`footer`)],16)):n(``,!0)],16)],16)}c.render=l;export{c as t};