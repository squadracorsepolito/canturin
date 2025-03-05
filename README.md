# Canturin

## Status

-   Bus (Done)
-   Node (Done)
-   Message (WIP)
-   Signal Types (Done)
-   Signal Unit (Done)
-   Signal Enum (Done)

## Prerequirements

-   [go](https://go.dev/) at least v1.22.4
-   [node js](https://nodejs.org/en) at least v20.16.0
-   [pnpm](https://pnpm.io/installation) at least v9.x

## Run the Project

1. Install wails3:

    ```
    go install -v github.com/wailsapp/wails/v3/cmd/wails3@latest
    ```

> [!IMPORTANT]
> Check wails3 is at version `v3.0.0-alpha.9` with the command `wails3 version`.

2. Clone this repo.

3. Navigate to the root folder.

4. Install go dependencies:

    ```
    go mod tidy
    ```

    Also, it may be useful to install in advance the frontend dependencies:

    ```
    cd frontend
    pnpm install
    ```

5. Run the application in development mode (in the root folder):

    ```
    wails3 dev
    ```

    This will start the application and enable hot-reloading for both frontend and backend changes.

    In MacOS/Linux it is better to run the application with sudo:

    ```
    sudo wails3 dev
    ```
