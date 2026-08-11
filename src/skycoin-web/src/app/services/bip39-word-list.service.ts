export class Bip39WordListService {

  private wordMap!: Map<string, boolean>;

  constructor() {
    // Unwrap the module's default export: esbuild turns a JSON import into a
    // module with one named export per top-level key plus `default`, and keys
    // that are not valid JS identifiers cannot become named exports. This one
    // reads `list`, which happens to be a valid identifier, so it works by luck
    // rather than by construction.
    import(`../../assets/bip39-word-list.json`).then((m: any) => m.default ?? m).then((result: any) => {
      this.wordMap = new Map<string, boolean>();
      result.list.forEach((word: string) => {
        this.wordMap.set(word, true);
      });
    });
  }

  validateWord(word: string): boolean | null {
    if (this.wordMap) {
      if (!this.wordMap.has(word)) {
        return false;
      }

      return true;
    } else {
      return null;
    }
  }
}
