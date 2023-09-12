/**
 * Submit the form containing the element
 * @param {HTMLElement} el
 */
window.cvSubmitNearestForm = function (el) {
  var form = null;
  var parent = el;
  for (var i = 0; i < 10; i += 1) {
    parent = parent.parentNode;
    if (!parent) {
      break;
    }
    if (parent.tagName === 'FORM') {
      form = parent;
      break;
    }
  }
  if (!form) {
    console.error('Element: ', el);
    throw new Error('No form found around element');
  }

  form.submit();
};
