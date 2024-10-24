package chain

func contains[T comparable](list []T, elem T) bool {
    for _, item := range list {
        if item == elem {
            return true
        }
    }
    return false
}