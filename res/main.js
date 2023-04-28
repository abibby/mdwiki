//@ts-check

/**
 *
 * @param {ClipboardEvent} event
 */
async function onPaste(event) {
  /** @type {HTMLTextAreaElement} */
  const input = event.target;

  // TODO: look into this: use event.originalEvent.clipboard for newer chrome versions
  if (event.clipboardData === null) {
    return;
  }
  var items = event.clipboardData.items;
  // find pasted image among pasted items
  var blob = null;
  for (const item of items) {
    if (item.type.indexOf("image") === 0) {
      blob = item.getAsFile();
    }
  }
  // load image if there is a pasted image
  if (blob !== null) {
    event.preventDefault();

    const loading = "![](...uploading image)";

    insertText(input, loading);

    // var boundary = Math.random().toString().slice(2);
    const body = new FormData();
    body.set("file", blob);
    const result = await fetch("/upload", {
      method: "POST",
      body: body,
    }).then((r) => r.json());
    replaceText(input, loading, `![alt](${result.file})`, 2, 5);
  }
}

for (const textarea of document.querySelectorAll("textarea")) {
  textarea.addEventListener("paste", onPaste);
}

/**
 *
 * @param {HTMLTextAreaElement} input
 * @param {string} text
 */
function insertText(input, text) {
  // document.execCommand("insertText", false, loading);
  const start = input.selectionStart;
  input.value =
    input.value.slice(0, input.selectionStart) +
    text +
    input.value.slice(input.selectionEnd);

  input.setSelectionRange(start + text.length, null);
}

/**
 *
 * @param {HTMLTextAreaElement} input
 * @param {string} from
 * @param {string} to
 * @param {number} selectStart
 * @param {number} selectEnd
 */
function replaceText(input, from, to, selectStart, selectEnd) {
  // document.execCommand("insertText", false, loading);
  const selectionStart = input.selectionStart;
  const value = input.value;

  const start = value.indexOf(from);

  input.value =
    input.value.slice(0, start) + to + input.value.slice(start + from.length);

  input.setSelectionRange(start + selectStart, start + selectEnd);
}
