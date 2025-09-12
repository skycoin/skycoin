export class SearchPage {
  search(text) {
    return $('.search-bar-container input').setValue(text);
  }
}
