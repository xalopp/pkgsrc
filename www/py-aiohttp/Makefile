# $NetBSD: Makefile,v 1.11 2018/01/18 10:29:54 adam Exp $

DISTNAME=	aiohttp-2.3.9
PKGNAME=	${PYPKGPREFIX}-${DISTNAME}
CATEGORIES=	www python
MASTER_SITES=	${MASTER_SITE_PYPI:=a/aiohttp/}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	https://github.com/aio-libs/aiohttp/
COMMENT=	Async http client/server framework
LICENSE=	apache-2.0

DEPENDS+=	${PYPKGPREFIX}-chardet-[0-9]*:../../converters/py-chardet
DEPENDS+=	${PYPKGPREFIX}-multidict>=3.0.0:../../databases/py-multidict
DEPENDS+=	${PYPKGPREFIX}-async-timeout>=1.2.0:../../devel/py-async-timeout
DEPENDS+=	${PYPKGPREFIX}-yarl>=0.11:../../www/py-yarl
BUILD_DEPENDS+=	${PYPKGPREFIX}-gunicorn-[0-9]*:../../www/py-gunicorn
BUILD_DEPENDS+=	${PYPKGPREFIX}-test-[0-9]*:../../devel/py-test
BUILD_DEPENDS+=	${PYPKGPREFIX}-test-timeout-[0-9]*:../../devel/py-test-timeout

PYTHON_VERSIONS_INCOMPATIBLE=	27

.include "../../lang/python/egg.mk"
.include "../../mk/bsd.pkg.mk"
