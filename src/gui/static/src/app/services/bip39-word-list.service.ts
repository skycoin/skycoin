import { debounceTime, map } from 'rxjs/operators';
import { Subject, Observable } from 'rxjs';


export class Bip39WordListService {

  get searchResults(): Observable<string[]> {
    return this.searchResultsSubject.asObservable().pipe(debounceTime(100), map(searchTerm => {
      if (searchTerm.length > 1) {
        return this.wordList.filter(option => option.startsWith(searchTerm));
      } else {
        return [];
      }
    }));
  }

  private lastSearchTerm = '';
  private searchResultsSubject: Subject<string> = new Subject<string>();
  private wordList: string[] = [];
  private wordMap: Map<string, boolean> = new Map<string, boolean>();

  constructor() {
    System.import(`../../assets/bip39-word-list.json`).then (result => {
      this.wordList = result.list;
      this.wordList.forEach(word => {
        this.wordMap.set(word, true);
      });
      this.searchResultsSubject.next(this.lastSearchTerm);
    });
  }

  setSearchTerm(value: string) {
    this.lastSearchTerm = value;
    this.searchResultsSubject.next(value);
  }

  validateWord(word: string): boolean | null {
    if (this.wordList.length > 0) {
      if (!this.wordMap.has(word)) {
        return false;
      }

      return true;
    } else {
      return null;
    }
  }
}
