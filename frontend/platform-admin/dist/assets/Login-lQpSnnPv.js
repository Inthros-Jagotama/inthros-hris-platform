import{At as e,Ft as t,Ht as n,It as r,Kt as i,On as a,Ot as o,Pt as s,Yt as c,d as l,dn as u,in as d,jt as f,o as p,s as m,sn as h,t as g}from"./button-DwvXYdIe.js";import{s as _,v}from"./index-nWODPmgs.js";import{i as y,r as b,t as x}from"./inputtext-2C0VN0p8.js";var S={name:`Card`,extends:{name:`BaseCard`,extends:p,style:l.extend({name:`card`,style:`
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
`,classes:{root:`p-card p-component`,header:`p-card-header`,body:`p-card-body`,caption:`p-card-caption`,title:`p-card-title`,subtitle:`p-card-subtitle`,content:`p-card-content`,footer:`p-card-footer`}}),provide:function(){return{$pcCard:this,$parentInstance:this}}},inheritAttrs:!1};function C(t,r,a,s,l,u){return i(),f(`div`,n({class:t.cx(`root`)},t.ptmi(`root`)),[t.$slots.header?(i(),f(`div`,n({key:0,class:t.cx(`header`)},t.ptm(`header`)),[c(t.$slots,`header`)],16)):e(``,!0),o(`div`,n({class:t.cx(`body`)},t.ptm(`body`)),[t.$slots.title||t.$slots.subtitle?(i(),f(`div`,n({key:0,class:t.cx(`caption`)},t.ptm(`caption`)),[t.$slots.title?(i(),f(`div`,n({key:0,class:t.cx(`title`)},t.ptm(`title`)),[c(t.$slots,`title`)],16)):e(``,!0),t.$slots.subtitle?(i(),f(`div`,n({key:1,class:t.cx(`subtitle`)},t.ptm(`subtitle`)),[c(t.$slots,`subtitle`)],16)):e(``,!0)],16)):e(``,!0),o(`div`,n({class:t.cx(`content`)},t.ptm(`content`)),[c(t.$slots,`content`)],16),t.$slots.footer?(i(),f(`div`,n({key:1,class:t.cx(`footer`)},t.ptm(`footer`)),[c(t.$slots,`footer`)],16)):e(``,!0)],16)],16)}S.render=C;var w={class:`min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50 via-white to-indigo-100 p-4`},T={class:`w-full max-w-sm`},E={key:0,class:`text-sm text-rose-600 bg-rose-50 border border-rose-200 rounded-md px-3 py-2`},D={__name:`Login`,setup(n){let c=_(),{login:l}=m(),p=h(``),C=h(``),D=h(!1),O=h(``);async function k(){if(!p.value||!C.value){O.value=`Please enter email and password`;return}D.value=!0,O.value=``;try{await l(p.value,C.value),c.push(`/dashboard`)}catch(e){O.value=e.response?.data?.message||`Invalid credentials. Please try again.`}finally{D.value=!1}}return(n,c)=>(i(),f(`div`,w,[o(`div`,T,[c[5]||=s(`<div class="text-center mb-8"><div class="inline-flex items-center justify-center w-14 h-14 rounded-xl bg-indigo-600 text-white mb-4"><i class="pi pi-shield text-2xl"></i></div><h1 class="text-xl font-bold text-gray-900">HRIS Platform Admin</h1><p class="text-sm text-gray-500 mt-1">Sign in to manage your platform</p></div>`,1),r(u(S),{class:`!shadow-lg !rounded-xl`},{content:d(()=>[o(`form`,{onSubmit:v(k,[`prevent`]),class:`space-y-4`},[o(`div`,null,[c[2]||=o(`label`,{class:`block text-sm font-medium text-gray-600 mb-1`},`Email`,-1),r(u(y),null,{default:d(()=>[r(u(b),{class:`pi pi-envelope`}),r(u(x),{modelValue:p.value,"onUpdate:modelValue":c[0]||=e=>p.value=e,type:`email`,placeholder:`admin@company.com`,class:`!w-full`,disabled:D.value,autocomplete:`email`},null,8,[`modelValue`,`disabled`])]),_:1})]),o(`div`,null,[c[3]||=o(`label`,{class:`block text-sm font-medium text-gray-600 mb-1`},`Password`,-1),r(u(y),null,{default:d(()=>[r(u(b),{class:`pi pi-lock`}),r(u(x),{modelValue:C.value,"onUpdate:modelValue":c[1]||=e=>C.value=e,type:`password`,placeholder:`••••••••`,class:`!w-full`,disabled:D.value,autocomplete:`current-password`},null,8,[`modelValue`,`disabled`])]),_:1})]),O.value?(i(),f(`div`,E,[c[4]||=o(`i`,{class:`pi pi-exclamation-circle mr-1`},null,-1),t(` `+a(O.value),1)])):e(``,!0),r(u(g),{type:`submit`,label:`Sign In`,icon:`pi pi-sign-in`,class:`!w-full`,loading:D.value},null,8,[`loading`])],32)]),_:1}),c[6]||=o(`p`,{class:`text-center text-sm text-gray-400 mt-6`},` HRIS Platform v1.6.3 — Enterprise Edition `,-1)])]))}};export{D as default};