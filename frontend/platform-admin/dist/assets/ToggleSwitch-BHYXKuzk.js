import{An as e,At as t,Dn as n,Gt as r,Jt as i,Lt as a,Mt as o,Nt as s,Ut as c,Zt as l,f as u,ln as d,pn as f,st as p}from"./button-CU9bjc3B.js";import{r as m}from"./inputtext-CadJn40T.js";var h=u.extend({name:`toggleswitch`,style:`
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
`,classes:{root:function(e){var t=e.instance,n=e.props;return[`p-toggleswitch p-component`,{"p-toggleswitch-checked":t.checked,"p-disabled":n.disabled,"p-invalid":t.$invalid}]},input:`p-toggleswitch-input`,slider:`p-toggleswitch-slider`,handle:`p-toggleswitch-handle`},inlineStyles:{root:{position:`relative`}}}),g={name:`ToggleSwitch`,extends:{name:`BaseToggleSwitch`,extends:m,props:{trueValue:{type:null,default:!0},falseValue:{type:null,default:!1},readonly:{type:Boolean,default:!1},tabindex:{type:Number,default:null},inputId:{type:String,default:null},inputClass:{type:[String,Object],default:null},inputStyle:{type:Object,default:null},ariaLabelledby:{type:String,default:null},ariaLabel:{type:String,default:null}},style:h,provide:function(){return{$pcToggleSwitch:this,$parentInstance:this}}},inheritAttrs:!1,emits:[`change`,`focus`,`blur`],methods:{getPTOptions:function(e){return(e===`root`?this.ptmi:this.ptm)(e,{context:{checked:this.checked,disabled:this.disabled}})},onChange:function(e){if(!this.disabled&&!this.readonly){var t=this.checked?this.falseValue:this.trueValue;this.writeValue(t,e),this.$emit(`change`,e)}},onFocus:function(e){this.$emit(`focus`,e)},onBlur:function(e){var t,n;this.$emit(`blur`,e),(t=(n=this.formField).onBlur)==null||t.call(n,e)}},computed:{checked:function(){return this.d_value===this.trueValue},dataP:function(){return p({checked:this.checked,disabled:this.disabled,invalid:this.$invalid})}}},_=[`data-p-checked`,`data-p-disabled`,`data-p`],v=[`id`,`checked`,`tabindex`,`disabled`,`readonly`,`aria-checked`,`aria-labelledby`,`aria-label`,`aria-invalid`],y=[`data-p`],b=[`data-p`];function x(e,n,r,a,o,u){return i(),s(`div`,c({class:e.cx(`root`),style:e.sx(`root`)},u.getPTOptions(`root`),{"data-p-checked":u.checked,"data-p-disabled":e.disabled,"data-p":u.dataP}),[t(`input`,c({id:e.inputId,type:`checkbox`,role:`switch`,class:[e.cx(`input`),e.inputClass],style:e.inputStyle,checked:u.checked,tabindex:e.tabindex,disabled:e.disabled,readonly:e.readonly,"aria-checked":u.checked,"aria-labelledby":e.ariaLabelledby,"aria-label":e.ariaLabel,"aria-invalid":e.invalid||void 0,onFocus:n[0]||=function(){return u.onFocus&&u.onFocus.apply(u,arguments)},onBlur:n[1]||=function(){return u.onBlur&&u.onBlur.apply(u,arguments)},onChange:n[2]||=function(){return u.onChange&&u.onChange.apply(u,arguments)}},u.getPTOptions(`input`)),null,16,v),t(`div`,c({class:e.cx(`slider`)},u.getPTOptions(`slider`),{"data-p":u.dataP}),[t(`div`,c({class:e.cx(`handle`)},u.getPTOptions(`handle`),{"data-p":u.dataP}),[l(e.$slots,`handle`,{checked:u.checked})],16,b)],16,y)],16,_)}g.render=x;var S={class:`inline-flex items-center gap-3 select-none`},C={__name:`ToggleSwitch`,props:{modelValue:{type:[Boolean,Number,String],default:!1},disabled:{type:Boolean,default:!1},label:{type:String,default:``},trueValue:{type:[Boolean,Number,String],default:!0},falseValue:{type:[Boolean,Number,String],default:!1}},emits:[`update:modelValue`],setup(t,{expose:c,emit:l}){let u=t,p=l,m=d(null),h=()=>{let e=u.modelValue==u.trueValue?u.falseValue:u.trueValue;p(`update:modelValue`,e)};return r(()=>{m.value?.$el?.hasAttribute(`autofocus`)&&m.value.$el.focus()}),c({focus:()=>m.value?.$el?.focus()}),(r,c)=>(i(),s(`div`,S,[a(f(g),{ref_key:`switchRef`,ref:m,"model-value":t.modelValue,disabled:t.disabled,"true-value":t.trueValue,"false-value":t.falseValue,"onUpdate:modelValue":c[0]||=e=>r.$emit(`update:modelValue`,e)},null,8,[`model-value`,`disabled`,`true-value`,`false-value`]),t.label?(i(),s(`label`,{key:0,class:n([`text-sm font-medium text-gray-700`,t.disabled?`opacity-50 cursor-not-allowed`:`cursor-pointer`]),onClick:c[1]||=e=>!t.disabled&&h()},e(t.label),3)):o(``,!0)]))}};export{C as t};