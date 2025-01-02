import { Z as current_component, R as ensure_array_like, Q as pop, O as push, Y as head } from "../../../../chunks/index3.js";
import { e as escape_html } from "../../../../chunks/escaping.js";
import { a as attr } from "../../../../chunks/attributes.js";
function onDestroy(fn) {
  var context = (
    /** @type {Component} */
    current_component
  );
  (context.d ??= []).push(fn);
}
const seedData = [
  {
    uuid: "colors",
    title: "Colors",
    cols: [
      { uuid: "colors-col-1", title: "Color", type: "coltypes-row-3" },
      { uuid: "colors-col-2", title: "Hex", type: "coltypes-row-3" },
      { uuid: "colors-col-3", title: "Red", type: "coltypes-row-4" },
      { uuid: "colors-col-4", title: "Green", type: "coltypes-row-4" },
      { uuid: "colors-col-5", title: "Blue", type: "coltypes-row-4" }
    ],
    rows: [
      { uuid: "colors-row-1", data: ["IndianRed", "#CD5C5C", "205", "92", "92"] },
      { uuid: "colors-row-2", data: ["LightCoral", "#F08080", "240", "128", "128"] },
      { uuid: "colors-row-2", data: ["Salmon", "#FA8072", "250", "128", "114"] }
    ]
  },
  {
    uuid: "coltypes",
    title: "Column types",
    cols: [
      { uuid: "coltypes-col-1", title: "Type", type: "coltypes-row-3" }
    ],
    rows: [
      { uuid: "coltypes-row-1", data: ["Any"] },
      { uuid: "coltypes-row-2", data: ["Title"] },
      { uuid: "coltypes-row-3", data: ["String"] },
      { uuid: "coltypes-row-4", data: ["Integer"] },
      { uuid: "coltypes-row-5", data: ["Decimal"] },
      { uuid: "coltypes-row-6", data: ["Date"] },
      { uuid: "coltypes-row-7", data: ["Boolean"] },
      { uuid: "coltypes-row-8", data: ["Text"] },
      { uuid: "coltypes-row-9", data: ["Grid"] },
      { uuid: "coltypes-row-10", data: ["View"] },
      { uuid: "coltypes-row-11", data: ["Image"] },
      { uuid: "coltypes-row-12", data: ["Video"] },
      { uuid: "coltypes-row-13", data: ["Sound"] }
    ]
  }
];
function Info($$payload, $$props) {
  push();
  let {
    focus,
    messageStack,
    isSending,
    messageStatus,
    isStreaming
  } = $$props;
  const each_array_1 = ensure_array_like(messageStack);
  $$payload.out += `<aside><div><p>`;
  if (isStreaming) {
    $$payload.out += "<!--[-->";
    $$payload.out += `Streaming messages`;
  } else {
    $$payload.out += "<!--[!-->";
  }
  $$payload.out += `<!--]--></p> <p>`;
  if (isSending) {
    $$payload.out += "<!--[-->";
    $$payload.out += `Sending message`;
  } else {
    $$payload.out += "<!--[!-->";
  }
  $$payload.out += `<!--]--> `;
  if (messageStatus) {
    $$payload.out += "<!--[-->";
    $$payload.out += `${escape_html(messageStatus)}`;
  } else {
    $$payload.out += "<!--[!-->";
  }
  $$payload.out += `<!--]--></p></div> `;
  if (focus.grid !== null) {
    $$payload.out += "<!--[-->";
    const each_array = ensure_array_like(focus.grid.cols);
    $$payload.out += `<ul><li class="svelte-1u94c12">Grid: ${escape_html(focus.grid.title)}</li> <li class="svelte-1u94c12">i: ${escape_html(focus.i)}</li> <li class="svelte-1u94c12">j: ${escape_html(focus.j)}</li> <li class="svelte-1u94c12">Columns <ul><!--[-->`;
    for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
      let col = each_array[$$index];
      $$payload.out += `<li class="svelte-1u94c12">${escape_html(col.title)}</li>`;
    }
    $$payload.out += `<!--]--></ul></li> <li class="svelte-1u94c12">Content: ${escape_html(focus.grid.rows[focus.i].data[focus.j])}</li></ul>`;
  } else {
    $$payload.out += "<!--[!-->";
  }
  $$payload.out += `<!--]--> <ul><!--[-->`;
  for (let $$index_1 = 0, $$length = each_array_1.length; $$index_1 < $$length; $$index_1++) {
    let message = each_array_1[$$index_1];
    if (message.request) {
      $$payload.out += "<!--[-->";
      $$payload.out += `<li class="request svelte-1u94c12">→ ${escape_html(message.request.messageKey)} ${escape_html(message.request.message.substring(0, 200))}</li>`;
    } else {
      $$payload.out += "<!--[!-->";
    }
    $$payload.out += `<!--]--> `;
    if (message.response) {
      $$payload.out += "<!--[-->";
      $$payload.out += `<li class="svelte-1u94c12">← ${escape_html(message.response.messageKey)} ${escape_html(message.response.message.substring(0, 200))}</li>`;
    } else {
      $$payload.out += "<!--[!-->";
    }
    $$payload.out += `<!--]-->`;
  }
  $$payload.out += `<!--]--></ul></aside>`;
  pop();
}
function _page($$payload, $$props) {
  push();
  let { data } = $$props;
  const dbName = data.dbName;
  data.gridUuid;
  data.url;
  const grids = seedData;
  let focus = { grid: null, i: -1, j: -1 };
  let isSending = false;
  let messageStatus = "";
  let isStreaming = false;
  const messageStack = [{}];
  let loginId = "";
  let loginPassword = "";
  onDestroy(() => {
  });
  function initGrid(grid) {
    grid.search = "";
    grid.columnSeq = grid.cols.length;
    applyFilters(grid);
  }
  function applyFilters(grid) {
    if (grid.search === "") grid.rows.forEach((row) => row.filtered = true);
    else {
      const regex = new RegExp(grid.search, "i");
      grid.rows.forEach((row) => row.filtered = regex.test(row.data[0]));
    }
  }
  grids.forEach((grid) => initGrid(grid));
  function findGrid(uuid) {
    return grids.find((grid) => grid.uuid === uuid);
  }
  findGrid("coltypes");
  head($$payload, ($$payload2) => {
    $$payload2.title = `<title>εncooη - ${escape_html(data.dbName)}</title>`;
  });
  $$payload.out += `<div class="layout svelte-h8te8x"><main><h1>${escape_html(dbName)}</h1> `;
  {
    $$payload.out += "<!--[!-->";
    $$payload.out += `<form><label>Username<input${attr("value", loginId)}></label> <label>Passphrase<input${attr("value", loginPassword)} type="password"></label> <button type="submit">Log in</button></form>`;
  }
  $$payload.out += `<!--]--></main> `;
  Info($$payload, {
    focus,
    messageStack,
    isSending,
    messageStatus,
    isStreaming
  });
  $$payload.out += `<!----></div>`;
  pop();
}
export {
  _page as default
};
