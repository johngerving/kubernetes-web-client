type User = {
    id : number,
    email : string,
}

type Workspace = {
    id : number,
    name : string,
    owner : number,
}

type PostWorkspaceFormErrors = {
    name?: string,
}
interface PostWorkspaceFormData extends FormData {
    name: string,
    errors: PostWorkspaceFormErrors
}