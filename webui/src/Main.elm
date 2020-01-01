module Main exposing (Model, Msg(..), init, main, update, view)

import Browser
import Browser.Navigation as Nav
import Css exposing (..)
import Html.Styled exposing (Html, a, div, h1, img, text, toUnstyled)
import Html.Styled.Attributes exposing (css, href, src)
import Url exposing (Url)
import Url.Parser as Url exposing ((</>), Parser)


type Status
    = Locked
    | Unlocked


urlToStatus : Url -> Status
urlToStatus url =
    url
        |> Url.parse urlParser
        |> Maybe.withDefault Locked


urlParser : Parser (Status -> a) a
urlParser =
    Url.oneOf
        [ Url.map Locked Url.top
        , Url.map Unlocked (Url.s "unlocked")
        ]


type alias Model =
    { key : Nav.Key
    , status : Status
    }


init : () -> Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url key =
    ( Model key (urlToStatus url), Cmd.none )


type Msg
    = LinkClicked Browser.UrlRequest
    | UrlChanged Url


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        LinkClicked urlRequest ->
            case urlRequest of
                Browser.Internal url ->
                    ( model, Nav.pushUrl model.key (Url.toString url) )

                Browser.External href ->
                    ( model, Nav.load href )

        UrlChanged url ->
            ( { model | status = urlToStatus url }
            , Cmd.none
            )


view : Model -> Browser.Document Msg
view model =
    let
        body =
            div
                [ css
                    [ backgroundColor (rgb 255 255 255)
                    , border3 (px 1) solid (rgb 120 120 120)
                    , margin2 (px 50) auto
                    , padding (px 10)
                    , width (px 200)
                    , height (px 200)
                    ]
                ]
                (case model.status of
                    Locked ->
                        [ a [ href "/unlocked" ]
                            [ img [ src "/locked.svg" ] []
                            ]
                        ]

                    Unlocked ->
                        [ a [ href "/" ]
                            [ img [ src "/unlocked.svg" ] []
                            ]
                        ]
                )
    in
    { title =
        case model.status of
            Locked ->
                "Strongbox - Locked"

            Unlocked ->
                "Strongbox - Unlocked"
    , body = [ Html.Styled.toUnstyled body ]
    }


main : Program () Model Msg
main =
    Browser.application
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        , onUrlChange = UrlChanged
        , onUrlRequest = LinkClicked
        }
