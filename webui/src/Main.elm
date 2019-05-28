module Main exposing (Model, Msg(..), init, main, update, view)

import Browser
import Css exposing (..)
import Html.Styled exposing (Html, div, h1, img, text, toUnstyled)
import Html.Styled.Attributes exposing (css, src)


type alias Model =
    {}


init : ( Model, Cmd Msg )
init =
    ( {}, Cmd.none )


type Msg
    = NoOp


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    ( model, Cmd.none )


view : Model -> Html Msg
view model =
    div
        [ css
            [ backgroundColor (rgb 255 255 255)
            , border3 (px 1) solid (rgb 120 120 120)
            , margin2 (px 50) auto
            , padding (px 10)
            , width (px 200)
            ]
        ]
        [ img
            [ src "/locked.svg" ]
            []
        ]


main : Program () Model Msg
main =
    Browser.element
        { view = view >> toUnstyled
        , init = \_ -> init
        , update = update
        , subscriptions = always Sub.none
        }
