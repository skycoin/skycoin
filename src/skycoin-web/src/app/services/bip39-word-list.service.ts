export class Bip39WordListService {

  private wordMap!: Map<string, boolean>;

  constructor() {
    import(`../../assets/bip39-word-list.json`).then ((result: any) => {
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
