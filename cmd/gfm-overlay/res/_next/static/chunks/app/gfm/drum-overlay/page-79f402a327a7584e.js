(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[301],{9065:function(e,s,t){Promise.resolve().then(t.bind(t,1947))},1947:function(e,s,t){"use strict";t.r(s),t.d(s,{default:function(){return y}});var a=t(9268);t(1642);var r=t(6006),u=t(6008),l=t(4791),n=t(1202);let i={UP:"⭡","UP,LEFT":"⭦","LEFT,UP":"⭦",LEFT:"⭠","LEFT,DOWN":"⭩","DOWN,LEFT":"⭩",DOWN:"⭣","DOWN,RIGHT":"⭨","RIGHT,DOWN":"⭨",RIGHT:"⭢","RIGHT,UP":"⭧","UP,RIGHT":"⭧","NEXT-ACTION":"⟳","PREV-ACTION":"⟲","RESET-ACTION":"⮌","CHANGE-MODE":"\uD83D\uDDD8",ATTACK:"⚔️",JUMP:"\uD83E\uDDB6",SKILL:"\uD83E\uDE84",SPELL:"\uD83D\uDCA3",PAUSE:"=",SWAP:"✨"};function o(){var e,s,t,a,l,n;let i=(0,u.useSearchParams)(),o=null!==(e=i.get("action_key"))&&void 0!==e?e:"Z",c=null!==(s=i.get("jump_key"))&&void 0!==s?s:"SPACE",d=null!==(t=i.get("skill_key"))&&void 0!==t?t:"X",y=null!==(a=i.get("spell_key"))&&void 0!==a?a:"C",m=null!==(l=i.get("pause_key"))&&void 0!==l?l:"V",x=null!==(n=i.get("swap_key"))&&void 0!==n?n:"S",k=(0,r.useMemo)(()=>({[o]:"ATTACK",[c]:"JUMP",[d]:"SKILL",[y]:"SPELL",[m]:"PAUSE",[x]:"SWAP"}),[o,c,d,y,m,x]);return{key2action:k}}function c(e){var s;let{x:t,y:r,keys:u}=e,{key2action:l}=o(),n=null!==(s=l[u])&&void 0!==s?s:u,c=i[n];if("SWAP"==n){let e=i["NEXT-ACTION"];return(0,a.jsxs)(a.Fragment,{children:[(0,a.jsx)("p",{className:"action ".concat(n),style:{left:t,top:r},children:c}),(0,a.jsx)("p",{className:"action ".concat(n),style:{left:t,top:r},children:e})]})}return(0,a.jsx)("p",{className:"action ".concat(n),style:{left:t,top:r},children:c})}function d(e){let[s,t]=(0,r.useState)(!1),{className:u,lastPress:l,duration:n}=e,i=null!=n?n:200;return((0,r.useEffect)(()=>{if(!l)return;let e=Date.parse(l);if(!e){console.error("failed to parse the press date: ".concat(l));return}let s=i+e-Date.now();if(s<=0){console.warn("timeout already expired");return}let a=setTimeout(()=>t(!1),s);return t(!0),()=>{clearTimeout(a)}},[t,l,i]),s)?(0,a.jsx)("div",{className:u}):null}function y(){var e;let s=(0,u.useSearchParams)(),t=(0,l.Z)(s.get("refresh_rate_ms"),25),i=null!==(e=s.get("drum_url"))&&void 0!==e?e:"/ram_store/drums",{dateSetter:y,bassDate:m,crash1Date:x,crash2Date:k,hihatPedalDate:f,hihatDate:v,rideDate:S,snareDate:j,tom1Date:N,tom2Date:P,tom3Date:_,tom4Date:p,keysSetter:h,bassKeys:E,crash1Keys:T,crash2Keys:b,hihatPedalKeys:g,hihatKeys:L,rideKeys:O,snareKeys:A,tom1Keys:D,tom2Keys:I,tom3Keys:C,tom4Keys:U}=function(){var e,s,t,a,l,n,i,o,c,d,y;let[m,x]=(0,r.useState)(""),[k,f]=(0,r.useState)(""),[v,S]=(0,r.useState)(""),[j,N]=(0,r.useState)(""),[P,_]=(0,r.useState)(""),[p,h]=(0,r.useState)(""),[E,T]=(0,r.useState)(""),[b,g]=(0,r.useState)(""),[L,O]=(0,r.useState)(""),[A,D]=(0,r.useState)(""),[I,C]=(0,r.useState)(""),[U,R]=(0,r.useState)(""),[W,w]=(0,r.useState)(""),[F,G]=(0,r.useState)(""),[H,K]=(0,r.useState)(""),[M,Z]=(0,r.useState)(""),[J,X]=(0,r.useState)(""),[B,V]=(0,r.useState)(""),[$,Y]=(0,r.useState)(""),[q,z]=(0,r.useState)(""),[Q,ee]=(0,r.useState)(""),[es,et]=(0,r.useState)(""),ea=(0,u.useSearchParams)(),er=null!==(e=ea.get("bass_event"))&&void 0!==e?e:"0x0924",eu=null!==(s=ea.get("crash1_event"))&&void 0!==s?s:"0x0931",el=null!==(t=ea.get("crash2_event"))&&void 0!==t?t:"0x0939",en=null!==(a=ea.get("hihat_pedal_event"))&&void 0!==a?a:"0x092c",ei=null!==(l=ea.get("hihat_event"))&&void 0!==l?l:"0x092a",eo=null!==(n=ea.get("ride_event"))&&void 0!==n?n:"0x0933",ec=null!==(i=ea.get("snare_event"))&&void 0!==i?i:"0x0926",ed=null!==(o=ea.get("tom1_event"))&&void 0!==o?o:"0x0930",ey=null!==(c=ea.get("tom2_event"))&&void 0!==c?c:"0x092d",em=null!==(d=ea.get("tom3_event"))&&void 0!==d?d:"0x092b",ex=null!==(y=ea.get("tom4_event"))&&void 0!==y?y:"0x0929",[ek,ef]=(0,r.useMemo)(()=>[{[er]:R,[eu]:w,[el]:G,[en]:K,[ei]:Z,[eo]:X,[ec]:V,[ed]:Y,[ey]:z,[em]:ee,[ex]:et},{[er]:x,[eu]:f,[el]:S,[en]:N,[ei]:_,[eo]:h,[ec]:T,[ed]:g,[ey]:O,[em]:D,[ex]:C}],[x,f,S,N,_,h,T,g,O,D,C,R,w,G,K,Z,X,V,Y,z,ee,et,er,eu,el,en,ei,eo,ec,ed,ey,em,ex]);return{bassDate:m,crash1Date:k,crash2Date:v,hihatPedalDate:j,hihatDate:P,rideDate:p,snareDate:E,tom1Date:b,tom2Date:L,tom3Date:A,tom4Date:I,bassKeys:U,crash1Keys:W,crash2Keys:F,hihatPedalKeys:H,hihatKeys:M,rideKeys:J,snareKeys:B,tom1Keys:$,tom2Keys:q,tom3Keys:Q,tom4Keys:es,keysSetter:ek,dateSetter:ef}}(),{keysSetter:R,attack:W,skill:w,spell:F,jump:G,swap:H,pause:K,up:M,down:Z,left:J,right:X}=function(){let[e,s]=(0,r.useState)(!1),[t,a]=(0,r.useState)(!1),[u,l]=(0,r.useState)(!1),[n,i]=(0,r.useState)(!1),[c,d]=(0,r.useState)(!1),[y,m]=(0,r.useState)(!1),[x,k]=(0,r.useState)(!1),[f,v]=(0,r.useState)(!1),[S,j]=(0,r.useState)(!1),[N,P]=(0,r.useState)(!1),{key2action:_}=o(),p=(0,r.useMemo)(()=>{let e={UP:k,DOWN:v,LEFT:j,RIGHT:P};return Object.entries(_).forEach(t=>{let[r,u]=t;switch(u){case"ATTACK":e[r]=s;break;case"SKILL":e[r]=a;break;case"SPELL":e[r]=l;break;case"JUMP":e[r]=i;break;case"SWAP":e[r]=d;break;case"PAUSE":e[r]=m}}),e},[_,s,a,l,i,d,m,k,v,j,P]);return{keysSetter:p,attack:e,skill:t,spell:u,jump:n,swap:c,pause:y,up:x,down:f,left:S,right:N}}(),B=(0,r.useCallback)(()=>{(0,n.Z)(i,void 0,e=>{e.midi&&Object.entries(e.midi).forEach(e=>{let[s,t]=e,a=y[s];if(!a){console.error("missing date setter for evevent ".concat(s));return}a(t)}),e.map&&Object.entries(e.map).forEach(e=>{let[s,t]=e,a=h[s];if(!a){console.error("missing date setter for evevent ".concat(s));return}a(t)}),e.keys&&Object.entries(e.keys).forEach(e=>{let[s,t]=e,a=R[s];if(!a){console.error("missing setter for key ".concat(s));return}a(t)})})},[i,y,h,R]);return(0,r.useEffect)(()=>{let e=setInterval(B,t);return()=>{clearInterval(e)}},[B,t]),(0,a.jsxs)(a.Fragment,{children:[(0,a.jsx)("div",{className:"base-drum drum-kit"}),(0,a.jsx)(d,{className:"bass-drum drum-kit",lastPress:m}),(0,a.jsx)(c,{x:168,y:208,keys:E}),(0,a.jsx)(d,{className:"crash1-drum drum-kit",lastPress:x}),(0,a.jsx)(c,{x:112,y:16,keys:T}),(0,a.jsx)(d,{className:"crash2-drum drum-kit",lastPress:k}),(0,a.jsx)(c,{x:358,y:184,keys:b}),(0,a.jsx)(d,{className:"hihat-pedal-drum drum-kit",lastPress:f}),(0,a.jsx)(c,{x:8,y:240,keys:g}),(0,a.jsx)(d,{className:"hihat-drum drum-kit",lastPress:v}),(0,a.jsx)(c,{x:40,y:109,keys:L}),(0,a.jsx)(d,{className:"ride-drum drum-kit",lastPress:S}),(0,a.jsx)(c,{x:304,y:56,keys:O}),(0,a.jsx)(d,{className:"snare-drum drum-kit",lastPress:j}),(0,a.jsx)(c,{x:56,y:192,keys:A}),(0,a.jsx)(d,{className:"tom1-drum drum-kit",lastPress:N}),(0,a.jsx)(c,{x:128,y:104,keys:D}),(0,a.jsx)(d,{className:"tom2-drum drum-kit",lastPress:P}),(0,a.jsx)(c,{x:208,y:104,keys:I}),(0,a.jsx)(d,{className:"tom3-drum drum-kit",lastPress:_}),(0,a.jsx)(c,{x:272,y:152,keys:C}),(0,a.jsx)(d,{className:"tom4-drum drum-kit",lastPress:p}),(0,a.jsx)(c,{x:280,y:232,keys:U}),(0,a.jsxs)("div",{className:"keys base-key",children:[W?(0,a.jsx)("div",{className:"attack-key base-key"}):null,(0,a.jsx)(c,{x:0,y:0,keys:"ATTACK"}),w?(0,a.jsx)("div",{className:"skill-key base-key"}):null,(0,a.jsx)(c,{x:66,y:0,keys:"SKILL"}),F?(0,a.jsx)("div",{className:"spell-key base-key"}):null,(0,a.jsx)(c,{x:132,y:0,keys:"SPELL"}),G?(0,a.jsx)("div",{className:"jump-key base-key"}):null,(0,a.jsx)(c,{x:0,y:66,keys:"JUMP"}),H?(0,a.jsx)("div",{className:"swap-key base-key"}):null,(0,a.jsx)(c,{x:66,y:66,keys:"SWAP"}),K?(0,a.jsx)("div",{className:"pause-key base-key"}):null,(0,a.jsx)(c,{x:132,y:66,keys:"PAUSE"}),M?(0,a.jsx)("div",{className:"up-key base-key"}):null,(0,a.jsx)(c,{x:66,y:146,keys:"UP"}),Z?(0,a.jsx)("div",{className:"down-key base-key"}):null,(0,a.jsx)(c,{x:66,y:212,keys:"DOWN"}),J?(0,a.jsx)("div",{className:"left-key base-key"}):null,(0,a.jsx)(c,{x:0,y:212,keys:"LEFT"}),X?(0,a.jsx)("div",{className:"right-key base-key"}):null,(0,a.jsx)(c,{x:132,y:212,keys:"RIGHT"})]})]})}},4791:function(e,s,t){"use strict";function a(e){let s=arguments.length>1&&void 0!==arguments[1]?arguments[1]:.25,t=parseFloat(e+"");return t||s}t.d(s,{Z:function(){return a}})},1202:function(e,s,t){"use strict";function a(e,s,t,a){let r={...s,cache:"no-store",headers:{"content-type":"application/json"}};fetch(e,r).then(e=>{if(!e.ok){console.error(e.statusText),null==a||a(e.status);return}200==e.status&&t?e.json().then(e=>{t(e)}):204==e.status&&t&&t()})}t.d(s,{Z:function(){return a}})},1642:function(){},3177:function(e,s,t){"use strict";/**
 * @license React
 * react-jsx-runtime.production.min.js
 *
 * Copyright (c) Meta Platforms, Inc. and affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */var a=t(6006),r=Symbol.for("react.element"),u=Symbol.for("react.fragment"),l=Object.prototype.hasOwnProperty,n=a.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED.ReactCurrentOwner,i={key:!0,ref:!0,__self:!0,__source:!0};function o(e,s,t){var a,u={},o=null,c=null;for(a in void 0!==t&&(o=""+t),void 0!==s.key&&(o=""+s.key),void 0!==s.ref&&(c=s.ref),s)l.call(s,a)&&!i.hasOwnProperty(a)&&(u[a]=s[a]);if(e&&e.defaultProps)for(a in s=e.defaultProps)void 0===u[a]&&(u[a]=s[a]);return{$$typeof:r,type:e,key:o,ref:c,props:u,_owner:n.current}}s.Fragment=u,s.jsx=o,s.jsxs=o},9268:function(e,s,t){"use strict";e.exports=t(3177)},6008:function(e,s,t){e.exports=t(3027)}},function(e){e.O(0,[667,488,744],function(){return e(e.s=9065)}),_N_E=e.O()}]);