export class User {
    name: string = ''
    password: string = ''
    static compare(l: User, r: User): number {
        if (l.name == r.name) {
            return 0
        }
        return l.name > r.name ? 1 : -1
    }
}