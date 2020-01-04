module Main exposing (Model, Msg(..), init, main, update, view)

import Browser exposing (Document)
import Browser.Navigation as Nav
import Html
import Json.Decode exposing (Value)
import LocalStorage exposing (decodeStorage, login, logout)
import Page
import Page.Blank as Blank
import Page.NotFound as NotFound
import Page.Root as Root
import Route exposing (Route)
import Session exposing (Session)
import Url exposing (Url)
import Viewer exposing (Viewer)


type Model
    = Redirect Session
    | NotFound Session
    | Main Root.Model


type Msg
    = LinkClicked Browser.UrlRequest
    | UrlChanged Url
    | GotRootMsg Root.Msg
    | GotSession Session


init : Value -> Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url navKey =
    let
        maybeViewer =
            decodeStorage flags
    in
    changeRouteTo (Route.fromUrl url)
        (Redirect (Session.fromViewer navKey maybeViewer))


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case ( msg, model ) of
        ( LinkClicked urlRequest, _ ) ->
            case urlRequest of
                Browser.Internal url ->
                    ( model
                    , Nav.pushUrl
                        (Session.navKey (toSession model))
                        (Url.toString url)
                    )

                Browser.External href ->
                    ( model, Nav.load href )

        ( UrlChanged url, _ ) ->
            changeRouteTo (Route.fromUrl url) model

        ( GotRootMsg subMsg, Main subModel ) ->
            Root.update subMsg subModel
                |> updateWith Main GotRootMsg model

        ( GotSession session, _ ) ->
            ( model
            , Route.replaceUrl (Session.navKey session) Route.Root
            )

        ( _, _ ) ->
            ( model, Cmd.none )


view : Model -> Document Msg
view model =
    let
        viewer =
            Session.viewer (toSession model)

        viewPage page toMsg config =
            let
                { title, body } =
                    Page.view viewer page config
            in
            { title = title
            , body = List.map (Html.map toMsg) body
            }
    in
    case model of
        NotFound _ ->
            Page.view viewer Page.Other NotFound.view

        Redirect _ ->
            Page.view viewer Page.Other Blank.view

        Main model_ ->
            viewPage Page.Main GotRootMsg (Root.view model_)


subscriptions : Model -> Sub Msg
subscriptions model =
    case model of
        NotFound _ ->
            Sub.none

        Redirect _ ->
            Session.changes GotSession (Session.navKey (toSession model))

        Main m ->
            Sub.map GotRootMsg (Root.subscriptions m)


main : Program Value Model Msg
main =
    Browser.application
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        , onUrlChange = UrlChanged
        , onUrlRequest = LinkClicked
        }



-- INTERNAL


changeRouteTo : Maybe Route -> Model -> ( Model, Cmd Msg )
changeRouteTo maybeRoute model =
    let
        session =
            toSession model
    in
    case maybeRoute of
        Just Route.Root ->
            Root.init session
                |> updateWith Main GotRootMsg model

        Just Route.Login ->
            ( model
            , Cmd.batch
                [ login (Viewer "jdoe@example.com" "Doe" "John")
                , Nav.replaceUrl (Session.navKey session) "/"
                ]
            )

        Just Route.Logout ->
            ( model
            , Cmd.batch
                [ logout
                , Nav.replaceUrl (Session.navKey session) "/"
                ]
            )

        Nothing ->
            ( NotFound session, Cmd.none )


toSession : Model -> Session
toSession page =
    case page of
        NotFound session ->
            session

        Redirect session ->
            session

        Main model ->
            Root.toSession model


updateWith : (subModel -> Model) -> (subMsg -> Msg) -> Model -> ( subModel, Cmd subMsg ) -> ( Model, Cmd Msg )
updateWith toModel toMsg _ ( subModel, subCmd ) =
    ( toModel subModel
    , Cmd.map toMsg subCmd
    )
