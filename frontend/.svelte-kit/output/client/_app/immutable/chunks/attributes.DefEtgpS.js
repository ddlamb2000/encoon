import{a8 as e,a9 as i,A as t,aa as v,ab as u}from"./runtime.BU7-Umqy.js";import{a as h}from"./render.D-xr2h7m.js";function A(r){if(t){var s=!1,a=()=>{if(!s){if(s=!0,r.hasAttribute("value")){var o=r.value;d(r,"value",null),r.value=o}if(r.hasAttribute("checked")){var _=r.checked;d(r,"checked",null),r.checked=_}}};r.__on_r=a,u(a),h()}}function y(r,s){var a=r.__attributes??(r.__attributes={});a.value===(a.value=s??void 0)||r.value===s&&(s!==0||r.nodeName!=="PROGRESS")||(r.value=s)}function d(r,s,a,o){var _=r.__attributes??(r.__attributes={});t&&(_[s]=r.getAttribute(s),s==="src"||s==="srcset"||s==="href"&&r.nodeName==="LINK")||_[s]!==(_[s]=a)&&(s==="style"&&"__styles"in r&&(r.__styles={}),s==="loading"&&(r[e]=a),a==null?r.removeAttribute(s):typeof a!="string"&&g(r).includes(s)?r[s]=a:r.setAttribute(s,a))}var c=new Map;function g(r){var s=c.get(r.nodeName);if(s)return s;c.set(r.nodeName,s=[]);for(var a,o=i(r),_=Element.prototype;_!==o;){a=v(o);for(var f in a)a[f].set&&s.push(f);o=i(o)}return s}export{d as a,A as r,y as s};
