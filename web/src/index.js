import { Application } from 'stimulus'
import { definitionsFromContext } from 'stimulus/webpack-helpers'
import './css/style.scss'

import jquery from 'jquery'
import 'bootstrap'
window.jQuery = jquery
window.$ = jquery

const application = Application.start()
const context = require.context('./controllers', true, /\.js$/)
application.load(definitionsFromContext(context))
