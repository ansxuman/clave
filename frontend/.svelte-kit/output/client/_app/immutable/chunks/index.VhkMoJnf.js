var V=Object.defineProperty;var W=(t,e,n)=>e in t?V(t,e,{enumerable:!0,configurable:!0,writable:!0,value:n}):t[e]=n;var C=(t,e,n)=>W(t,typeof e!="symbol"?e+"":e,n);import{n as w,r as E,i as P,f as j,h as G,j as I,k as J,l as K,m as Q,p as X,q as O,v as Y,w as Z,x as tt}from"./scheduler.CgIr7hM3.js";const z=typeof window<"u";let et=z?()=>window.performance.now():()=>Date.now(),R=z?t=>requestAnimationFrame(t):w;const x=new Set;function L(t){x.forEach(e=>{e.c(t)||(x.delete(e),e.f())}),x.size!==0&&R(L)}function nt(t){let e;return x.size===0&&R(L),{promise:new Promise(n=>{x.add(e={c:t,f:n})}),abort(){x.delete(e)}}}let A=!1;function it(){A=!0}function rt(){A=!1}function st(t,e,n,i){for(;t<e;){const r=t+(e-t>>1);n(r)<=i?t=r+1:e=r}return t}function at(t){if(t.hydrate_init)return;t.hydrate_init=!0;let e=t.childNodes;if(t.nodeName==="HEAD"){const s=[];for(let a=0;a<e.length;a++){const u=e[a];u.claim_order!==void 0&&s.push(u)}e=s}const n=new Int32Array(e.length+1),i=new Int32Array(e.length);n[0]=-1;let r=0;for(let s=0;s<e.length;s++){const a=e[s].claim_order,u=(r>0&&e[n[r]].claim_order<=a?r+1:st(1,r,_=>e[n[_]].claim_order,a))-1;i[s]=n[u]+1;const f=u+1;n[f]=s,r=Math.max(f,r)}const o=[],l=[];let c=e.length-1;for(let s=n[r]+1;s!=0;s=i[s-1]){for(o.push(e[s-1]);c>=s;c--)l.push(e[c]);c--}for(;c>=0;c--)l.push(e[c]);o.reverse(),l.sort((s,a)=>s.claim_order-a.claim_order);for(let s=0,a=0;s<l.length;s++){for(;a<o.length&&l[s].claim_order>=o[a].claim_order;)a++;const u=a<o.length?o[a]:null;t.insertBefore(l[s],u)}}function lt(t,e){t.appendChild(e)}function M(t){if(!t)return document;const e=t.getRootNode?t.getRootNode():t.ownerDocument;return e&&e.host?e:t.ownerDocument}function ot(t){const e=H("style");return e.textContent="/* empty */",ct(M(t),e),e.sheet}function ct(t,e){return lt(t.head||t,e),e.sheet}function ft(t,e){if(A){for(at(t),(t.actual_end_child===void 0||t.actual_end_child!==null&&t.actual_end_child.parentNode!==t)&&(t.actual_end_child=t.firstChild);t.actual_end_child!==null&&t.actual_end_child.claim_order===void 0;)t.actual_end_child=t.actual_end_child.nextSibling;e!==t.actual_end_child?(e.claim_order!==void 0||e.parentNode!==t)&&t.insertBefore(e,t.actual_end_child):t.actual_end_child=e.nextSibling}else(e.parentNode!==t||e.nextSibling!==null)&&t.appendChild(e)}function jt(t,e,n){A&&!n?ft(t,e):(e.parentNode!==t||e.nextSibling!=n)&&t.insertBefore(e,n||null)}function T(t){t.parentNode&&t.parentNode.removeChild(t)}function H(t){return document.createElement(t)}function B(t){return document.createTextNode(t)}function Pt(){return B(" ")}function Rt(){return B("")}function Bt(t,e,n,i){return t.addEventListener(e,n,i),()=>t.removeEventListener(e,n,i)}function Dt(t,e,n){n==null?t.removeAttribute(e):t.getAttribute(e)!==n&&t.setAttribute(e,n)}function It(t){return t.dataset.svelteH}function ut(t){return Array.from(t.childNodes)}function _t(t){t.claim_info===void 0&&(t.claim_info={last_index:0,total_claimed:0})}function F(t,e,n,i,r=!1){_t(t);const o=(()=>{for(let l=t.claim_info.last_index;l<t.length;l++){const c=t[l];if(e(c)){const s=n(c);return s===void 0?t.splice(l,1):t[l]=s,r||(t.claim_info.last_index=l),c}}for(let l=t.claim_info.last_index-1;l>=0;l--){const c=t[l];if(e(c)){const s=n(c);return s===void 0?t.splice(l,1):t[l]=s,r?s===void 0&&t.claim_info.last_index--:t.claim_info.last_index=l,c}}return i()})();return o.claim_order=t.claim_info.total_claimed,t.claim_info.total_claimed+=1,o}function dt(t,e,n,i){return F(t,r=>r.nodeName===e,r=>{const o=[];for(let l=0;l<r.attributes.length;l++){const c=r.attributes[l];n[c.name]||o.push(c.name)}o.forEach(l=>r.removeAttribute(l))},()=>i(e))}function Ot(t,e,n){return dt(t,e,n,H)}function mt(t,e){return F(t,n=>n.nodeType===3,n=>{const i=""+e;if(n.data.startsWith(i)){if(n.data.length!==i.length)return n.splitText(i.length)}else n.data=i},()=>B(e),!0)}function kt(t){return mt(t," ")}function qt(t,e){e=""+e,t.data!==e&&(t.data=e)}function zt(t,e,n,i){n==null?t.style.removeProperty(e):t.style.setProperty(e,n,"")}function ht(t,e,{bubbles:n=!1,cancelable:i=!1}={}){return new CustomEvent(t,{detail:e,bubbles:n,cancelable:i})}function Lt(t,e){return new t(e)}const N=new Map;let b=0;function pt(t){let e=5381,n=t.length;for(;n--;)e=(e<<5)-e^t.charCodeAt(n);return e>>>0}function $t(t,e){const n={stylesheet:ot(e),rules:{}};return N.set(t,n),n}function yt(t,e,n,i,r,o,l,c=0){const s=16.666/i;let a=`{
`;for(let h=0;h<=1;h+=s){const y=e+(n-e)*o(h);a+=h*100+`%{${l(y,1-y)}}
`}const u=a+`100% {${l(n,1-n)}}
}`,f=`__svelte_${pt(u)}_${c}`,_=M(t),{stylesheet:d,rules:m}=N.get(_)||$t(_,t);m[f]||(m[f]=!0,d.insertRule(`@keyframes ${f} ${u}`,d.cssRules.length));const $=t.style.animation||"";return t.style.animation=`${$?`${$}, `:""}${f} ${i}ms linear ${r}ms 1 both`,b+=1,f}function k(t,e){const n=(t.style.animation||"").split(", "),i=n.filter(e?o=>o.indexOf(e)<0:o=>o.indexOf("__svelte")===-1),r=n.length-i.length;r&&(t.style.animation=i.join(", "),b-=r,b||xt())}function xt(){R(()=>{b||(N.forEach(t=>{const{ownerNode:e}=t.stylesheet;e&&T(e)}),N.clear())})}let g;function gt(){return g||(g=Promise.resolve(),g.then(()=>{g=null})),g}function q(t,e,n){t.dispatchEvent(ht(`intro${n}`))}const v=new Set;let p;function Mt(){p={r:0,c:[],p}}function Tt(){p.r||E(p.c),p=p.p}function wt(t,e){t&&t.i&&(v.delete(t),t.i(e))}function Ht(t,e,n,i){if(t&&t.o){if(v.has(t))return;v.add(t),p.c.push(()=>{v.delete(t),i&&(n&&t.d(1),i())}),t.o(e)}else i&&i()}const vt={duration:0};function Ft(t,e,n){const i={direction:"in"};let r=e(t,n,i),o=!1,l,c,s=0;function a(){l&&k(t,l)}function u(){const{delay:_=0,duration:d=300,easing:m=G,tick:$=w,css:h}=r||vt;h&&(l=yt(t,0,1,d,_,m,h,s++)),$(0,1);const y=et()+_,U=y+d;c&&c.abort(),o=!0,j(()=>q(t,!0,"start")),c=nt(S=>{if(o){if(S>=U)return $(1,0),q(t,!0,"end"),a(),o=!1;if(S>=y){const D=m((S-y)/d);$(D,1-D)}}return o})}let f=!1;return{start(){f||(f=!0,k(t),P(r)?(r=r(i),gt().then(u)):u())},invalidate(){f=!1},end(){o&&(a(),o=!1)}}}function Ut(t){t&&t.c()}function Vt(t,e){t&&t.l(e)}function Nt(t,e,n){const{fragment:i,after_update:r}=t.$$;i&&i.m(e,n),j(()=>{const o=t.$$.on_mount.map(Y).filter(P);t.$$.on_destroy?t.$$.on_destroy.push(...o):E(o),t.$$.on_mount=[]}),r.forEach(j)}function bt(t,e){const n=t.$$;n.fragment!==null&&(Q(n.after_update),E(n.on_destroy),n.fragment&&n.fragment.d(e),n.on_destroy=n.fragment=null,n.ctx=[])}function Et(t,e){t.$$.dirty[0]===-1&&(Z.push(t),tt(),t.$$.dirty.fill(0)),t.$$.dirty[e/31|0]|=1<<e%31}function Wt(t,e,n,i,r,o,l=null,c=[-1]){const s=X;O(t);const a=t.$$={fragment:null,ctx:[],props:o,update:w,not_equal:r,bound:I(),on_mount:[],on_destroy:[],on_disconnect:[],before_update:[],after_update:[],context:new Map(e.context||(s?s.$$.context:[])),callbacks:I(),dirty:c,skip_bound:!1,root:e.target||s.$$.root};l&&l(a.root);let u=!1;if(a.ctx=n?n(t,e.props||{},(f,_,...d)=>{const m=d.length?d[0]:_;return a.ctx&&r(a.ctx[f],a.ctx[f]=m)&&(!a.skip_bound&&a.bound[f]&&a.bound[f](m),u&&Et(t,f)),_}):[],a.update(),u=!0,E(a.before_update),a.fragment=i?i(a.ctx):!1,e.target){if(e.hydrate){it();const f=ut(e.target);a.fragment&&a.fragment.l(f),f.forEach(T)}else a.fragment&&a.fragment.c();e.intro&&wt(t.$$.fragment),Nt(t,e.target,e.anchor),rt(),J()}O(s)}class Gt{constructor(){C(this,"$$");C(this,"$$set")}$destroy(){bt(this,1),this.$destroy=w}$on(e,n){if(!P(n))return w;const i=this.$$.callbacks[e]||(this.$$.callbacks[e]=[]);return i.push(n),()=>{const r=i.indexOf(n);r!==-1&&i.splice(r,1)}}$set(e){this.$$set&&!K(e)&&(this.$$.skip_bound=!0,this.$$set(e),this.$$.skip_bound=!1)}}const At="4";typeof window<"u"&&(window.__svelte||(window.__svelte={v:new Set})).v.add(At);export{Ft as A,Gt as S,ut as a,mt as b,Ot as c,T as d,H as e,kt as f,jt as g,ft as h,Wt as i,qt as j,Dt as k,wt as l,Ht as m,Rt as n,Tt as o,zt as p,Mt as q,Lt as r,Pt as s,B as t,Ut as u,Vt as v,Nt as w,bt as x,It as y,Bt as z};
