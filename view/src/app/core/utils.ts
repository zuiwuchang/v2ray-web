export function sortString(l: string, r: string): number {
    if (l == r) {
        return 0
    }
    return l < r ? -1 : 1
}
export function sortNameValue(l: { name: string, value: number }, r: { name: string, value: number }): number {
    if (l.value == r.value) {
        return sortString(l.name, r.name)
    }
    return l.value - r.value
}