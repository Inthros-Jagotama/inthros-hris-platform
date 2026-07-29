import{A as e,C as t,D as n,H as r,K as i,O as a,R as o,T as s,b as c,c as l,ct as u,dt as d,l as f,m as p,o as m,p as h,r as ee,s as g,u as _,z as v}from"./runtime-core.esm-bundler-CmryIYX-.js";import{o as y,tt as b}from"./ripple-DtVIDCWU.js";import{b as x,i as S,n as te,p as C,t as w,x as ne}from"./index-DVa_2WOz.js";import{t as re}from"./useI18n-B3Td2MI9.js";import{t as T}from"./baseeditableholder-DkY7pyuN.js";import{t as E}from"./inputtext-DK6WE_d5.js";import{n as ie,t as D}from"./tag-DavKWYzI.js";import{n as ae,r as oe,t as O}from"./column-c1YkNctB.js";import{t as se}from"./SkeletonTable-DGNqAjaG.js";var k=y.extend({name:`toggleswitch`,style:`
    .p-toggleswitch {
        display: inline-block;
        width: dt('toggleswitch.width');
        height: dt('toggleswitch.height');
    }

    .p-toggleswitch-input {
        cursor: pointer;
        appearance: none;
        position: absolute;
        top: 0;
        inset-inline-start: 0;
        width: 100%;
        height: 100%;
        padding: 0;
        margin: 0;
        opacity: 0;
        z-index: 1;
        outline: 0 none;
        border-radius: dt('toggleswitch.border.radius');
    }

    .p-toggleswitch-slider {
        cursor: pointer;
        width: 100%;
        height: 100%;
        border-width: dt('toggleswitch.border.width');
        border-style: solid;
        border-color: dt('toggleswitch.border.color');
        background: dt('toggleswitch.background');
        transition:
            background dt('toggleswitch.transition.duration'),
            color dt('toggleswitch.transition.duration'),
            border-color dt('toggleswitch.transition.duration'),
            outline-color dt('toggleswitch.transition.duration'),
            box-shadow dt('toggleswitch.transition.duration');
        border-radius: dt('toggleswitch.border.radius');
        outline-color: transparent;
        box-shadow: dt('toggleswitch.shadow');
    }

    .p-toggleswitch-handle {
        position: absolute;
        top: 50%;
        display: flex;
        justify-content: center;
        align-items: center;
        background: dt('toggleswitch.handle.background');
        color: dt('toggleswitch.handle.color');
        width: dt('toggleswitch.handle.size');
        height: dt('toggleswitch.handle.size');
        inset-inline-start: dt('toggleswitch.gap');
        margin-block-start: calc(-1 * calc(dt('toggleswitch.handle.size') / 2));
        border-radius: dt('toggleswitch.handle.border.radius');
        transition:
            background dt('toggleswitch.transition.duration'),
            color dt('toggleswitch.transition.duration'),
            inset-inline-start dt('toggleswitch.slide.duration'),
            box-shadow dt('toggleswitch.slide.duration');
    }

    .p-toggleswitch.p-toggleswitch-checked .p-toggleswitch-slider {
        background: dt('toggleswitch.checked.background');
        border-color: dt('toggleswitch.checked.border.color');
    }

    .p-toggleswitch.p-toggleswitch-checked .p-toggleswitch-handle {
        background: dt('toggleswitch.handle.checked.background');
        color: dt('toggleswitch.handle.checked.color');
        inset-inline-start: calc(dt('toggleswitch.width') - calc(dt('toggleswitch.handle.size') + dt('toggleswitch.gap')));
    }

    .p-toggleswitch:not(.p-disabled):has(.p-toggleswitch-input:hover) .p-toggleswitch-slider {
        background: dt('toggleswitch.hover.background');
        border-color: dt('toggleswitch.hover.border.color');
    }

    .p-toggleswitch:not(.p-disabled):has(.p-toggleswitch-input:hover) .p-toggleswitch-handle {
        background: dt('toggleswitch.handle.hover.background');
        color: dt('toggleswitch.handle.hover.color');
    }

    .p-toggleswitch:not(.p-disabled):has(.p-toggleswitch-input:hover).p-toggleswitch-checked .p-toggleswitch-slider {
        background: dt('toggleswitch.checked.hover.background');
        border-color: dt('toggleswitch.checked.hover.border.color');
    }

    .p-toggleswitch:not(.p-disabled):has(.p-toggleswitch-input:hover).p-toggleswitch-checked .p-toggleswitch-handle {
        background: dt('toggleswitch.handle.checked.hover.background');
        color: dt('toggleswitch.handle.checked.hover.color');
    }

    .p-toggleswitch:not(.p-disabled):has(.p-toggleswitch-input:focus-visible) .p-toggleswitch-slider {
        box-shadow: dt('toggleswitch.focus.ring.shadow');
        outline: dt('toggleswitch.focus.ring.width') dt('toggleswitch.focus.ring.style') dt('toggleswitch.focus.ring.color');
        outline-offset: dt('toggleswitch.focus.ring.offset');
    }

    .p-toggleswitch.p-invalid > .p-toggleswitch-slider {
        border-color: dt('toggleswitch.invalid.border.color');
    }

    .p-toggleswitch.p-disabled {
        opacity: 1;
    }

    .p-toggleswitch.p-disabled .p-toggleswitch-slider {
        background: dt('toggleswitch.disabled.background');
    }

    .p-toggleswitch.p-disabled .p-toggleswitch-handle {
        background: dt('toggleswitch.handle.disabled.background');
    }
`,classes:{root:function(e){var t=e.instance,n=e.props;return[`p-toggleswitch p-component`,{"p-toggleswitch-checked":t.checked,"p-disabled":n.disabled,"p-invalid":t.$invalid}]},input:`p-toggleswitch-input`,slider:`p-toggleswitch-slider`,handle:`p-toggleswitch-handle`},inlineStyles:{root:{position:`relative`}}}),A={name:`ToggleSwitch`,extends:{name:`BaseToggleSwitch`,extends:T,props:{trueValue:{type:null,default:!0},falseValue:{type:null,default:!1},readonly:{type:Boolean,default:!1},tabindex:{type:Number,default:null},inputId:{type:String,default:null},inputClass:{type:[String,Object],default:null},inputStyle:{type:Object,default:null},ariaLabelledby:{type:String,default:null},ariaLabel:{type:String,default:null}},style:k,provide:function(){return{$pcToggleSwitch:this,$parentInstance:this}}},inheritAttrs:!1,emits:[`change`,`focus`,`blur`],methods:{getPTOptions:function(e){return(e===`root`?this.ptmi:this.ptm)(e,{context:{checked:this.checked,disabled:this.disabled}})},onChange:function(e){if(!this.disabled&&!this.readonly){var t=this.checked?this.falseValue:this.trueValue;this.writeValue(t,e),this.$emit(`change`,e)}},onFocus:function(e){this.$emit(`focus`,e)},onBlur:function(e){var t,n;this.$emit(`blur`,e),(t=(n=this.formField).onBlur)==null||t.call(n,e)}},computed:{checked:function(){return this.d_value===this.trueValue},dataP:function(){return b({checked:this.checked,disabled:this.disabled,invalid:this.$invalid})}}},j=[`data-p-checked`,`data-p-disabled`,`data-p`],M=[`id`,`checked`,`tabindex`,`disabled`,`readonly`,`aria-checked`,`aria-labelledby`,`aria-label`,`aria-invalid`],N=[`data-p`],P=[`data-p`];function F(e,t,n,r,i,o){return s(),_(`div`,c({class:e.cx(`root`),style:e.sx(`root`)},o.getPTOptions(`root`),{"data-p-checked":o.checked,"data-p-disabled":e.disabled,"data-p":o.dataP}),[g(`input`,c({id:e.inputId,type:`checkbox`,role:`switch`,class:[e.cx(`input`),e.inputClass],style:e.inputStyle,checked:o.checked,tabindex:e.tabindex,disabled:e.disabled,readonly:e.readonly,"aria-checked":o.checked,"aria-labelledby":e.ariaLabelledby,"aria-label":e.ariaLabel,"aria-invalid":e.invalid||void 0,onFocus:t[0]||=function(){return o.onFocus&&o.onFocus.apply(o,arguments)},onBlur:t[1]||=function(){return o.onBlur&&o.onBlur.apply(o,arguments)},onChange:t[2]||=function(){return o.onChange&&o.onChange.apply(o,arguments)}},o.getPTOptions(`input`)),null,16,M),g(`div`,c({class:e.cx(`slider`)},o.getPTOptions(`slider`),{"data-p":o.dataP}),[g(`div`,c({class:e.cx(`handle`)},o.getPTOptions(`handle`),{"data-p":o.dataP}),[a(e.$slots,`handle`,{checked:o.checked})],16,P)],16,N)],16,j)}A.render=F;var I={class:`space-y-1`},L={class:`flex items-center justify-between gap-2 flex-wrap`},R={class:`flex items-center gap-1.5`},z={class:`flex items-center gap-2`},B={key:0,class:`text-xs text-gray-400 dark:text-gray-500`},V={class:`flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500`},H={class:`text-sm font-medium`},U={class:`text-gray-800 dark:text-gray-100 font-medium`},W={class:`text-gray-500 dark:text-gray-400`},G={class:`text-gray-500 dark:text-gray-400`},K={class:`flex items-center gap-1`},ce={class:`space-y-4`},le={class:`text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2 flex items-center gap-1.5`},ue={class:`space-y-2`},de={class:`block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1`},fe={key:0,class:`text-red-500 text-xs mt-1 block`},pe={class:`block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1`},me={key:0,class:`text-red-500 text-xs mt-1 block`},he={class:`block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1`},ge={class:`flex items-center justify-between`},_e={class:`block text-sm font-medium text-gray-600 dark:text-gray-300`},ve={class:`block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1`},ye={class:`flex items-center justify-between`},be={class:`flex items-center gap-2 ml-auto`},q={__name:`ZonesView`,setup(a){let{t:c}=re(),y=ne(),b=x(),T=r([]),k=r(!1),j=r(0),M=r(1),N=r(15),P=r(null),F=r(!1),q=r(!1),J=r(null),Y=r(!1),X=r({}),Z=r({code:``,name:``,region:``,is_active:!0,sort_order:0}),xe=[{type:`tag`,width:`w-20`,headerWidth:`w-16`},{type:`text`,width:`w-32`,headerWidth:`w-16`},{type:`text`,width:`w-24`,headerWidth:`w-16`},{type:`tag`,width:`w-16`,headerWidth:`w-12`},{type:`text`,width:`w-12`,headerWidth:`w-16`},{type:`icons`,count:2,headerWidth:`w-16`}],Se=m(()=>[{label:c(`common.all`),value:null,severity:`info`},{label:c(`common_status.active`),value:`active`,severity:`success`},{label:c(`common_status.inactive`),value:`inactive`,severity:`warn`}]),Ce=m(()=>{let e=T.value;return P.value===`active`?e=e.filter(e=>e.is_active===!0):P.value===`inactive`&&(e=e.filter(e=>e.is_active===!1)),e}),we=m(()=>(M.value-1)*N.value);async function Q(){k.value=!0;try{let e=(await C.get(`/api/v1/tenant/settings/zones`,{params:{page:M.value,per_page:N.value}})).data;T.value=e?.data||[],j.value=e?.total||0,e?.page&&(M.value=e.page)}catch(e){y.add({severity:`error`,summary:c(`message.error`),detail:e.response?.data?.error?.message||c(`message.failed_to_load`),life:4e3})}finally{k.value=!1}}function Te(e){M.value=e.page+1,N.value=e.rows,Q()}function $(e){q.value=!!e,J.value=e?.id||null,X.value={},Z.value={code:e?.code||``,name:e?.name||``,region:e?.region||``,is_active:e?.is_active===void 0||e.is_active,sort_order:e?.sort_order||0},F.value=!0}function Ee(){Z.value={code:``,name:``,region:``,is_active:!0,sort_order:0},X.value={},q.value=!1,J.value=null}async function De(){if(X.value={},!Z.value.code?.trim()){X.value={code:[c(`form.required`)]};return}if(!Z.value.name?.trim()){X.value={name:[c(`form.required`)]};return}Y.value=!0;try{let e={code:Z.value.code,name:Z.value.name,region:Z.value.region||void 0,is_active:Z.value.is_active,sort_order:Z.value.sort_order||0};q.value?(await C.put(`/api/v1/tenant/settings/zones/${J.value}`,e),y.add({severity:`success`,summary:c(`message.success`),detail:c(`zones.updated`),life:3e3})):(await C.post(`/api/v1/tenant/settings/zones`,e),y.add({severity:`success`,summary:c(`message.success`),detail:c(`zones.created`),life:3e3})),F.value=!1,await Q()}catch(e){let t=ie(e);Object.keys(t).length>0?X.value=t:y.add({severity:`error`,summary:c(`message.error`),detail:e.response?.data?.error?.message||c(`message.operation_failed`),life:4e3})}finally{Y.value=!1}}function Oe(e){b.require({header:c(`zones.confirm_delete_title`),message:c(`zones.confirm_delete`,{name:e.name}),icon:`pi pi-exclamation-triangle`,rejectLabel:c(`common.cancel`),acceptLabel:c(`common.delete`),rejectClass:`p-button-outlined p-button-secondary`,acceptClass:`p-button-danger`,accept:async()=>{try{await C.delete(`/api/v1/tenant/settings/zones/${e.id}`),y.add({severity:`success`,summary:c(`message.success`),detail:c(`zones.deleted`),life:3e3}),await Q()}catch(e){y.add({severity:`error`,summary:c(`message.error`),detail:e.response?.data?.error?.message||c(`message.operation_failed`),life:4e3})}}})}return t(Q),(t,r)=>{let a=e(`tooltip`);return s(),_(`div`,I,[g(`div`,L,[g(`div`,R,[(s(!0),_(ee,null,n(Se.value,e=>(s(),l(i(S),{key:e.value,label:e.label,severity:P.value===e.value&&e.severity||`secondary`,outlined:P.value!==e.value,size:`small`,class:`!text-xs !px-2 !py-1`,onClick:t=>P.value=e.value},null,8,[`label`,`severity`,`outlined`,`onClick`]))),128))]),g(`div`,z,[j.value>0?(s(),_(`span`,B,d(j.value)+` `+d(i(c)(`common.items`)),1)):f(``,!0),p(i(S),{label:i(c)(`zones.new_zone`),icon:`pi pi-plus`,size:`small`,onClick:r[0]||=e=>$()},null,8,[`label`])])]),k.value?(s(),l(se,{key:0,columns:xe,rows:8})):(s(),l(i(ae),{key:1,value:Ce.value,lazy:``,totalRecords:j.value,first:we.value,rows:N.value,onPage:r[1]||=e=>Te(e),paginator:``,paginatorTemplate:`FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown`,rowsPerPageOptions:[10,15,25,50],size:`small`,class:`!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden`,sortField:`sort_order`,sortOrder:1},{empty:o(()=>[g(`div`,V,[r[9]||=g(`i`,{class:`pi pi-map-marker text-3xl mb-2 opacity-50`},null,-1),g(`p`,H,d(i(c)(`zones.empty_title`)),1)])]),default:o(()=>[p(i(O),{field:`code`,header:i(c)(`zones.code`),sortable:``,style:{width:`120px`}},{body:o(({data:e})=>[p(i(D),{value:e.code,severity:`info`,class:`!text-xs !px-1.5 !py-0.5`},null,8,[`value`])]),_:1},8,[`header`]),p(i(O),{field:`name`,header:i(c)(`zones.name`),sortable:``},{body:o(({data:e})=>[g(`span`,U,d(e.name),1)]),_:1},8,[`header`]),p(i(O),{field:`region`,header:i(c)(`zones.region`),sortable:``,style:{width:`150px`}},{body:o(({data:e})=>[g(`span`,W,d(e.region||`—`),1)]),_:1},8,[`header`]),p(i(O),{field:`is_active`,header:i(c)(`common.status`),sortable:``,style:{width:`100px`}},{body:o(({data:e})=>[p(i(D),{value:e.is_active?i(c)(`common_status.active`):i(c)(`common_status.inactive`),severity:e.is_active?`success`:`warn`,class:`!text-xs !px-1.5 !py-0.5`},null,8,[`value`,`severity`])]),_:1},8,[`header`]),p(i(O),{field:`sort_order`,header:i(c)(`zones.sort_order`),sortable:``,style:{width:`100px`}},{body:o(({data:e})=>[g(`span`,G,d(e.sort_order),1)]),_:1},8,[`header`]),p(i(O),{header:i(c)(`common.actions`),style:{width:`100px`},frozen:``,alignFrozen:`right`},{body:o(({data:e})=>[g(`div`,K,[v(p(i(S),{icon:`pi pi-pencil`,size:`small`,text:``,severity:`secondary`,onClick:t=>$(e)},null,8,[`onClick`]),[[a,i(c)(`common.edit`),void 0,{left:!0}]]),v(p(i(S),{icon:`pi pi-trash`,size:`small`,text:``,severity:`danger`,onClick:t=>Oe(e)},null,8,[`onClick`]),[[a,i(c)(`common.delete`),void 0,{left:!0}]])])]),_:1},8,[`header`])]),_:1},8,[`value`,`totalRecords`,`first`,`rows`])),p(i(te),{visible:F.value,"onUpdate:visible":r[8]||=e=>F.value=e,header:q.value?i(c)(`zones.edit_zone`):i(c)(`zones.new_zone`),modal:``,style:{width:`520px`},closable:!0,onHide:Ee},{footer:o(()=>[g(`div`,ye,[g(`div`,be,[p(i(S),{label:i(c)(`common.cancel`),severity:`secondary`,outlined:``,size:`small`,onClick:r[7]||=e=>F.value=!1},null,8,[`label`]),p(i(S),{label:q.value?i(c)(`common.update`):i(c)(`common.save`),size:`small`,loading:Y.value,disabled:Y.value,onClick:De},null,8,[`label`,`loading`,`disabled`])])])]),default:o(()=>[g(`div`,ce,[g(`div`,null,[g(`h3`,le,[r[10]||=g(`i`,{class:`pi pi-map-marker text-indigo-400 text-sm`},null,-1),h(` `+d(q.value?i(c)(`zones.edit_zone`):i(c)(`zones.new_zone`)),1)]),g(`div`,ue,[g(`div`,null,[g(`label`,de,[h(d(i(c)(`zones.code`))+` `,1),r[11]||=g(`span`,{class:`text-red-500`},`*`,-1)]),p(i(E),{modelValue:Z.value.code,"onUpdate:modelValue":r[2]||=e=>Z.value.code=e,class:u([`!w-full`,{"p-invalid":X.value?.code}]),maxlength:`20`,autofocus:``,placeholder:i(c)(`zones.code`)},null,8,[`modelValue`,`class`,`placeholder`]),X.value?.code?(s(),_(`small`,fe,d(X.value.code),1)):f(``,!0)]),g(`div`,null,[g(`label`,pe,[h(d(i(c)(`zones.name`))+` `,1),r[12]||=g(`span`,{class:`text-red-500`},`*`,-1)]),p(i(E),{modelValue:Z.value.name,"onUpdate:modelValue":r[3]||=e=>Z.value.name=e,class:u([`!w-full`,{"p-invalid":X.value?.name}]),maxlength:`255`,placeholder:i(c)(`zones.name`)},null,8,[`modelValue`,`class`,`placeholder`]),X.value?.name?(s(),_(`small`,me,d(X.value.name),1)):f(``,!0)]),g(`div`,null,[g(`label`,he,d(i(c)(`zones.region`)),1),p(i(E),{modelValue:Z.value.region,"onUpdate:modelValue":r[4]||=e=>Z.value.region=e,class:`!w-full`,maxlength:`100`,placeholder:i(c)(`zones.region`)},null,8,[`modelValue`,`placeholder`])]),g(`div`,ge,[g(`label`,_e,d(i(c)(`zones.is_active`)),1),p(i(A),{modelValue:Z.value.is_active,"onUpdate:modelValue":r[5]||=e=>Z.value.is_active=e},null,8,[`modelValue`])]),g(`div`,null,[g(`label`,ve,d(i(c)(`zones.sort_order`)),1),p(i(oe),{modelValue:Z.value.sort_order,"onUpdate:modelValue":r[6]||=e=>Z.value.sort_order=e,class:`!w-full`,min:0},null,8,[`modelValue`])])])])])]),_:1},8,[`visible`,`header`]),p(i(w))])}}};export{q as default};