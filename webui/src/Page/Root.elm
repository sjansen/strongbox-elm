module Page.Root exposing (Model, Msg, init, subscriptions, toSession, update, view)

import Css exposing (..)
import Html.Styled exposing (Html, br, button, div, img, text)
import Html.Styled.Attributes exposing (css, src)
import Html.Styled.Events exposing (onClick)
import Session exposing (Session)


type alias Model =
    { session : Session
    , unlocked : Bool
    }


type Msg
    = GotSession Session
    | Locked
    | Unlocked


init : Session -> ( Model, Cmd Msg )
init session =
    ( { session = session
      , unlocked = False
      }
    , Cmd.none
    )


subscriptions : Model -> Sub Msg
subscriptions model =
    Session.changes GotSession (Session.navKey model.session)


toSession : Model -> Session
toSession model =
    model.session


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case ( msg, model.unlocked ) of
        ( Locked, True ) ->
            ( { model | unlocked = False }, Cmd.none )

        ( Locked, False ) ->
            ( model, Cmd.none )

        ( Unlocked, True ) ->
            ( model, Cmd.none )

        ( Unlocked, False ) ->
            ( { model | unlocked = True }, Cmd.none )

        ( GotSession session, _ ) ->
            ( { model | session = session }, Cmd.none )


view : Model -> { title : String, content : Html Msg }
view model =
    let
        ( imgSrc, label, msg ) =
            if model.unlocked then
                ( "/unlocked.svg", "Lock", Locked )

            else
                ( "/locked.svg", "Unlock", Unlocked )
    in
    { title = ""
    , content =
        div
            [ css
                [ margin2 (px 50) auto
                , textAlign center
                ]
            ]
            [ img
                [ css
                    [ backgroundColor (rgb 255 255 255)
                    , border3 (px 1) solid (rgb 120 120 120)
                    , borderRadius (px 4)
                    , boxShadow4 (px 0) (px 1) (px 5) (rgba 0 0 0 0.2)
                    , height (px 200)
                    , width (px 200)
                    , padding (px 10)
                    ]
                , src imgSrc
                ]
                []
            , br [] []
            , button
                [ css
                    [ backgroundColor (rgb 64 64 192)
                    , borderRadius (px 2)
                    , border3 (px 1) solid (rgb 120 120 120)
                    , boxSizing borderBox
                    , color (rgb 255 255 255)
                    , cursor pointer
                    , padding (em 1)
                    , margin (em 2)
                    , textTransform uppercase
                    , width (em 8)
                    ]
                , onClick msg
                ]
                [ text label ]
            ]
    }
