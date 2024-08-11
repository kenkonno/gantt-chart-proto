import {createRouter, createWebHistory, RouteRecordRaw} from 'vue-router'
import {loggedIn} from "@/composable/auth";
import {Api} from "@/api/axios";

const routes: Array<RouteRecordRaw> = [
    {
        path: '/login',
        name: 'login',
        meta: {title: "工程管理ツール | ログイン", requiresAuth: false},
        component: () => import('../views/LoginView.vue'),
    },
    {
        path: '/',
        component: () => import('../views/TopPage.vue'),
        children: [
            {
                path: '/',
                name: 'gantt',
                meta: {title: "工程管理ツール | 案件ビュー", requiresAuth: true},
                component: () => import('../views/GanttFacilityView.vue')
            },
            {
                path: '/all-view',
                name: 'gantt-all-view',
                meta: {title: "工程管理ツール | 全体ビュー", requiresAuth: true},
                component: () => import('../views/GanttAllView.vue')
            },
        ]
    }
]

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes
})
router.beforeEach(async (to, from, next) => {

    // ゲストログインの処理を実施する
    if (to.query.uuid != undefined) {
        console.log("uuid is detected.")
        await Api.postLogin({id: "", password: "", uuid: String(to.query.uuid)});

    }

    if (to.matched.some(record => record.meta.requiresAuth)) {
        const {user} = await loggedIn()
        if (user?.id == undefined) {
            // ログインしていない場合、ログインページへリダイレクトします
            next({
                path: '/login',
                query: {redirect: to.fullPath},
            });
        } else {
            next();
        }
    } else {
        next();
    }
});
router.afterEach((to) => {
    document.title = to.meta.title as string
})
export default router
