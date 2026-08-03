import{At as e,Ft as t,It as n,Mt as r,Nt as i,Rt as a,Ut as o,Xt as s,Yt as c,dn as l,f as u,fn as d,ln as f,o as p,qt as m,st as h}from"./button-BrzTdEG-.js";import{o as g}from"./index-Dz8sog5S.js";var _=u.extend({name:`chart`,classes:{root:`p-chart`},inlineStyles:{root:{position:`relative`}}}),v={name:`Chart`,extends:{name:`BaseChart`,extends:p,props:{type:String,data:null,options:null,plugins:null,width:{type:Number,default:300},height:{type:Number,default:150},canvasProps:{type:null,default:null}},style:_,provide:function(){return{$pcChart:this,$parentInstance:this}}},inheritAttrs:!1,emits:[`select`,`loaded`],chart:null,watch:{data:{handler:function(){this.reinit()},deep:!0},type:function(){this.reinit()},options:function(){this.reinit()}},mounted:function(){this.initChart()},beforeUnmount:function(){this.chart&&=(this.chart.destroy(),null)},methods:{initChart:function(){var e=this;g(()=>import(`./auto-B5fFbl-6.js`).then(function(t){e.chart&&=(e.chart.destroy(),null),t&&t.default&&(e.chart=new t.default(e.$refs.canvas,{type:e.type,data:e.data,options:e.options,plugins:e.plugins})),e.$emit(`loaded`,e.chart)}),[])},getCanvas:function(){return this.$canvas},getChart:function(){return this.chart},getBase64Image:function(){return this.chart.toBase64Image()},refresh:function(){this.chart&&this.chart.update()},reinit:function(){this.initChart()},onCanvasClick:function(e){if(this.chart){var t=this.chart.getElementsAtEventForMode(e,`nearest`,{intersect:!0},!1),n=this.chart.getElementsAtEventForMode(e,`dataset`,{intersect:!0},!1);t&&t[0]&&n&&this.$emit(`select`,{originalEvent:e,element:t[0],dataset:n})}},generateLegend:function(){if(this.chart)return this.chart.generateLegend()}}};function y(e){"@babel/helpers - typeof";return y=typeof Symbol==`function`&&typeof Symbol.iterator==`symbol`?function(e){return typeof e}:function(e){return e&&typeof Symbol==`function`&&e.constructor===Symbol&&e!==Symbol.prototype?`symbol`:typeof e},y(e)}function b(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable})),n.push.apply(n,r)}return n}function x(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]==null?{}:arguments[t];t%2?b(Object(n),!0).forEach(function(t){S(e,t,n[t])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):b(Object(n)).forEach(function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))})}return e}function S(e,t,n){return(t=C(t))in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function C(e){var t=w(e,`string`);return y(t)==`symbol`?t:t+``}function w(e,t){if(y(e)!=`object`||!e)return e;var n=e[Symbol.toPrimitive];if(n!==void 0){var r=n.call(e,t);if(y(r)!=`object`)return r;throw TypeError(`@@toPrimitive must return a primitive value.`)}return(t===`string`?String:Number)(e)}var T=[`width`,`height`];function E(e,t,r,a,s,c){return m(),n(`div`,o({class:e.cx(`root`),style:e.sx(`root`)},e.ptmi(`root`)),[i(`canvas`,o({ref:`canvas`,width:e.width,height:e.height,onClick:t[0]||=function(e){return c.onCanvasClick(e)}},x(x({},e.canvasProps),e.ptm(`canvas`))),null,16,T)],16)}v.render=E;var D=u.extend({name:`progressbar`,style:`
    .p-progressbar {
        display: block;
        position: relative;
        overflow: hidden;
        height: dt('progressbar.height');
        background: dt('progressbar.background');
        border-radius: dt('progressbar.border.radius');
    }

    .p-progressbar-value {
        margin: 0;
        background: dt('progressbar.value.background');
    }

    .p-progressbar-label {
        color: dt('progressbar.label.color');
        font-size: dt('progressbar.label.font.size');
        font-weight: dt('progressbar.label.font.weight');
    }

    .p-progressbar-determinate .p-progressbar-value {
        height: 100%;
        width: 0%;
        position: absolute;
        display: none;
        display: flex;
        align-items: center;
        justify-content: center;
        overflow: hidden;
        transition: width 1s ease-in-out;
    }

    .p-progressbar-determinate .p-progressbar-label {
        display: inline-flex;
    }

    .p-progressbar-indeterminate .p-progressbar-value::before {
        content: '';
        position: absolute;
        background: inherit;
        inset-block-start: 0;
        inset-inline-start: 0;
        inset-block-end: 0;
        will-change: inset-inline-start, inset-inline-end;
        animation: p-progressbar-indeterminate-anim 2.1s cubic-bezier(0.65, 0.815, 0.735, 0.395) infinite;
    }

    .p-progressbar-indeterminate .p-progressbar-value::after {
        content: '';
        position: absolute;
        background: inherit;
        inset-block-start: 0;
        inset-inline-start: 0;
        inset-block-end: 0;
        will-change: inset-inline-start, inset-inline-end;
        animation: p-progressbar-indeterminate-anim-short 2.1s cubic-bezier(0.165, 0.84, 0.44, 1) infinite;
        animation-delay: 1.15s;
    }

    @keyframes p-progressbar-indeterminate-anim {
        0% {
            inset-inline-start: -35%;
            inset-inline-end: 100%;
        }
        60% {
            inset-inline-start: 100%;
            inset-inline-end: -90%;
        }
        100% {
            inset-inline-start: 100%;
            inset-inline-end: -90%;
        }
    }
    @-webkit-keyframes p-progressbar-indeterminate-anim {
        0% {
            inset-inline-start: -35%;
            inset-inline-end: 100%;
        }
        60% {
            inset-inline-start: 100%;
            inset-inline-end: -90%;
        }
        100% {
            inset-inline-start: 100%;
            inset-inline-end: -90%;
        }
    }

    @keyframes p-progressbar-indeterminate-anim-short {
        0% {
            inset-inline-start: -200%;
            inset-inline-end: 100%;
        }
        60% {
            inset-inline-start: 107%;
            inset-inline-end: -8%;
        }
        100% {
            inset-inline-start: 107%;
            inset-inline-end: -8%;
        }
    }
    @-webkit-keyframes p-progressbar-indeterminate-anim-short {
        0% {
            inset-inline-start: -200%;
            inset-inline-end: 100%;
        }
        60% {
            inset-inline-start: 107%;
            inset-inline-end: -8%;
        }
        100% {
            inset-inline-start: 107%;
            inset-inline-end: -8%;
        }
    }
`,classes:{root:function(e){var t=e.instance;return[`p-progressbar p-component`,{"p-progressbar-determinate":t.determinate,"p-progressbar-indeterminate":t.indeterminate}]},value:`p-progressbar-value`,label:`p-progressbar-label`}}),O={name:`ProgressBar`,extends:{name:`BaseProgressBar`,extends:p,props:{value:{type:Number,default:null},mode:{type:String,default:`determinate`},showValue:{type:Boolean,default:!0}},style:D,provide:function(){return{$pcProgressBar:this,$parentInstance:this}}},inheritAttrs:!1,computed:{progressStyle:function(){return{width:this.value+`%`,display:`flex`}},indeterminate:function(){return this.mode===`indeterminate`},determinate:function(){return this.mode===`determinate`},dataP:function(){return h({determinate:this.determinate,indeterminate:this.indeterminate})}}},k=[`aria-valuenow`,`data-p`],A=[`data-p`],j=[`data-p`],M=[`data-p`];function N(e,r,i,c,l,u){return m(),n(`div`,o({role:`progressbar`,class:e.cx(`root`),"aria-valuemin":`0`,"aria-valuenow":e.value,"aria-valuemax":`100`,"data-p":u.dataP},e.ptmi(`root`)),[u.determinate?(m(),n(`div`,o({key:0,class:e.cx(`value`),style:u.progressStyle,"data-p":u.dataP},e.ptm(`value`)),[e.value!=null&&e.value!==0&&e.showValue?(m(),n(`div`,o({key:0,class:e.cx(`label`),"data-p":u.dataP},e.ptm(`label`)),[s(e.$slots,`default`,{},function(){return[a(d(e.value+`%`),1)]})],16,j)):t(``,!0)],16,A)):u.indeterminate?(m(),n(`div`,o({key:1,class:e.cx(`value`),"data-p":u.dataP},e.ptm(`value`)),null,16,M)):t(``,!0)],16,k)}O.render=N;var P={key:3,class:`flex items-center gap-2`},F={class:`flex-1 space-y-1.5`},I={class:`flex items-end gap-1 h-20 mb-2`},L={class:`space-y-2.5`},R={__name:`SkeletonCard`,props:{type:{type:String,default:`kpi`,validator:e=>[`kpi`,`stat`,`metric`,`alert`,`sparkline`,`detail`].includes(e)},count:{type:Number,default:null},cols:{type:String,default:null},rows:{type:Number,default:4},padding:{type:String,default:`p-3`},valueWidth:{type:String,default:null},labelWidth:{type:String,default:null},sparkHeights:{type:Array,default:()=>[40,60,30,70,45,55,35,65]}},setup(a){let o=a,s=r(()=>o.count===null?o.type===`kpi`?6:o.type===`stat`?4:o.type===`alert`?3:4:o.count),u=r(()=>o.cols===null?o.type===`kpi`?`grid-cols-2 md:grid-cols-4 lg:grid-cols-6`:o.type===`stat`?`grid-cols-2 md:grid-cols-4`:`grid-cols-1`:o.cols);return(r,o)=>(m(),n(`div`,{class:f([`grid gap-3 animate-pulse`,u.value])},[(m(!0),n(e,null,c(s.value,r=>(m(),n(`div`,{key:r,class:f([`bg-white rounded-lg border border-gray-200`,a.padding])},[a.type===`kpi`?(m(),n(e,{key:0},[o[0]||=i(`div`,{class:`w-8 h-8 rounded-lg bg-gray-200 mb-2`},null,-1),i(`div`,{class:f([`h-5 rounded bg-gray-200 mb-1`,a.valueWidth])},null,2),i(`div`,{class:f([`h-3 rounded bg-gray-200`,a.labelWidth])},null,2)],64)):a.type===`stat`?(m(),n(e,{key:1},[i(`div`,{class:f([`h-3 rounded bg-gray-200 mb-2`,a.labelWidth])},null,2),i(`div`,{class:f([`h-5 rounded bg-gray-200`,a.valueWidth])},null,2)],64)):a.type===`metric`?(m(),n(e,{key:2},[i(`div`,{class:f([`h-4 rounded bg-gray-200 mb-3`,a.valueWidth||`w-24`])},null,2),i(`div`,{class:f([`h-6 rounded bg-gray-200 mb-1`,a.labelWidth||`w-16`])},null,2),o[1]||=i(`div`,{class:`h-2 w-12 rounded bg-gray-100`},null,-1)],64)):a.type===`alert`?(m(),n(`div`,P,[o[3]||=i(`div`,{class:`w-5 h-5 rounded-full bg-gray-200 shrink-0`},null,-1),i(`div`,F,[i(`div`,{class:f([`h-3 rounded bg-gray-200`,a.labelWidth||`w-full`])},null,2),o[2]||=i(`div`,{class:`h-3 rounded bg-gray-100 w-3/4`},null,-1)])])):a.type===`sparkline`?(m(),n(e,{key:4},[i(`div`,I,[(m(!0),n(e,null,c(a.sparkHeights,e=>(m(),n(`div`,{key:e,class:`flex-1 bg-gray-100 rounded-t`,style:l({height:e+`%`})},null,4))),128))]),i(`div`,{class:f([`h-3 rounded bg-gray-200`,a.labelWidth||`w-16`])},null,2)],64)):a.type===`detail`?(m(),n(e,{key:5},[i(`div`,{class:f([`h-4 rounded bg-gray-200 mb-3`,a.valueWidth||`w-32`])},null,2),i(`div`,L,[(m(!0),n(e,null,c(a.rows||4,e=>(m(),n(`div`,{key:e,class:`flex items-center justify-between`},[...o[4]||=[i(`div`,{class:`h-3 rounded bg-gray-200 w-20`},null,-1),i(`div`,{class:`h-3 rounded bg-gray-200 w-12`},null,-1)]]))),128))])],64)):t(``,!0)],2))),128))],2))}},z=(e,t)=>{let n=e.__vccOpts||e;for(let[e,r]of t)n[e]=r;return n};export{v as i,R as n,O as r,z as t};