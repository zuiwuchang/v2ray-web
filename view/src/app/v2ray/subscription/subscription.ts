export class Element {
    id: number = 0
    name: string = ''
    url: string = ''
    static compare(l: Element, r: Element): number {
        if (l.name == r.name) {
            if (l.id == r.id) {
                return 0
            }
            return l.id > r.id ? 1 : -1
        }
        return l.name > r.name ? 1 : -1
    }
}