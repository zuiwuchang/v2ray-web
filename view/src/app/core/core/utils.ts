import { BehaviorSubject } from 'rxjs';
import { filter } from 'rxjs/operators';
export class Closed {
    private closed_ = new BehaviorSubject<boolean>(false)
    public readonly observable = this.closed_.pipe(
        filter((ok) => ok),
    )
    close(): boolean {
        if (this.closed_.value) {
            return false
        }
        this.closed_.next(true)
        return true
    }
    get isClosed(): boolean {
        return this.closed_.value
    }
    get isNotClosed(): boolean {
        return !this.closed_.value
    }
}
